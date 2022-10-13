# Copyright 2022 Mia-Platform

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#    http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# General variables

# Set Output Directory Path
PROJECT_DIR := $(shell pwd -P)
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(PROJECT_DIR)/bin
endif
CMDNAME := vab
TOOLS_DIR := $(PROJECT_DIR)/tools
TOOLS_BIN := $(TOOLS_DIR)/bin

KIND_CLUSTER_1_NAME := vab-cluster-1
KIND_CLUSTER_2_NAME := vab-cluster-2

KIND_IMAGE_VERSION ?= kindest/node:v1.22.15@sha256:bfd5eaae36849bfb3c1e3b9442f3da17d730718248939d9d547e86bbac5da586
GOLANCI_LINT_VERSION := 1.49.0
DEEPCOPY_GEN_VERSION := 0.25.1

SUPPORTED_PLATFORMS := linux/386,linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6
$(eval BUILD_PLATFORMS = $(shell echo "$(SUPPORTED_PLATFORMS)" | sed "s#,# #g;s#/#.#g"))
MULTIARCH_BUILD := $(addprefix build.,$(BUILD_PLATFORMS))

# Golang variables
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

##@ Build

# Set the version number.
VERSION := $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
DATE_FMT := +%Y-%m-%d
COMPLETE_FMT := +%Y-%m-%dT%H:%M:%SZ
BUILD_DATE := $(shell date -u "$(DATE_FMT)")
COMPLETE_BUILD_DATE := $(shell date -u "$(COMPLETE_FMT)")

GO_LDFLAGS := -X github.com/mia-platform/vab/internal/cmd.Version=$(VERSION) $(GO_LDFLAGS)
GO_LDFLAGS := -X github.com/mia-platform/vab/internal/cmd.BuildDate=$(BUILD_DATE) $(GO_LDFLAGS)

.PHONY: build build-all
build: build.$(GOOS).$(GOARCH)
build-all: $(MULTIARCH_BUILD)

build.%:
	$(eval OS := $(word 1,$(subst ., ,$*)))
	$(eval ARCH := $(word 2,$(subst ., ,$*)))
	$(eval VARM := $(word 3,$(subst ., ,$*)))
	@echo "Building cli for ${OS} ${ARCH}${ARM}..."

	$(eval OUTDIR := $(shell if [ "${VARM}" ]; then echo "${OUTPUT_DIR}/${OS}/${ARCH}/${VARM}/${CMDNAME}"; else echo "${OUTPUT_DIR}/${OS}/${ARCH}/${CMDNAME}"; fi))
	$(eval ARM = $(subst v,,$(VARM)))
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) GOARM=$(ARM) go build -ldflags "$(GO_LDFLAGS)" \
		-o $(OUTDIR) $(PROJECT_DIR)/cmd/$(CMDNAME)

##@ Test

TEST_VERBOSE ?= "false"
.PHONY: test
test:
ifneq ($(TEST_VERBOSE), "false")
	go test -test.v ./...
else
	go test ./...
endif

.PHONY: test-e2e
test-e2e: start-kind e2e stop-kind

e2e:
ifneq ($(TEST_VERBOSE), "false")
	go test -test.v -tags=e2e ./... || true
else
	go test ./... -tags=e2e || true
endif

start-kind:
	@./tools/start-kind.sh ${KIND_IMAGE_VERSION} ${KIND_CLUSTER_1_NAME} ${KIND_CLUSTER_2_NAME}

stop-kind:
	@kind delete cluster --name ${KIND_CLUSTER_1_NAME}
	@kind delete cluster --name ${KIND_CLUSTER_2_NAME}

.PHONY: test-coverage
test-coverage:
	go test ./... -race -coverprofile=coverage.xml -covermode=atomic

.PHONY: show-coverage
show-coverage: test-coverage
	go tool cover -func=coverage.xml

##@ Clean project

.PHONY: clean
clean:
	@echo "Clean all artifact files..."
	@rm -fr $(OUTPUT_DIR)
	@rm -fr coverage.xml

.PHONY: clean-go
clean-go:
	@echo "Clean golang cache..."
	@go clean -cache

.PHONY: clean-tools
clean-tools:
	@echo "Clean tools folder..."
	@rm -fr .tools/bin

.PHONY: clean-all
clean-all: clean clean-go clean-tools

##@ Code generation

.PHONY: generate
generate: generate-dep
	@echo "Generating deepcopy code..."
	@$(TOOLS_BIN)/deepcopy-gen -i ./pkg/apis/vab.mia-platform.eu/v1alpha1 \
		-o "$(PROJECT_DIR)" -O zz_generated.deepcopy --go-header-file $(TOOLS_DIR)/boilerplate.go.txt

##@ Lint

MODE ?= "colored-line-number"

.PHONY: lint
lint: lintgo-dep lint-mod lint-ci

lint-ci:
	@echo "Linting go files..."
	$(TOOLS_BIN)/golangci-lint run --out-format=$(MODE) --config=$(TOOLS_DIR)/.golangci.yml

lint-mod:
	@echo "Run go mod tidy"
	@go mod tidy -compat=1.18
## ensure all changes have been committed
	git diff --exit-code -- go.mod
	git diff --exit-code -- go.sum

##@ Dependencies

.PHONY: install-dep
install-dep: generate-dep lintgo-dep

lintgo-dep:
	@GOBIN=$(TOOLS_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANCI_LINT_VERSION)

generate-dep:
	@GOBIN=$(TOOLS_BIN) go install k8s.io/code-generator/cmd/deepcopy-gen@v$(DEEPCOPY_GEN_VERSION)

##@ Images

# REGISTRY is the image registry to use for build and push image targets, default to docker hub
REGISTRY ?= docker.io/miaplatform

# TAG is the tag to use for build and push image targets, use git tag or latest
TAG ?= $(shell git describe --tags 2>/dev/null || echo latest)

# Force to use buildkit as engine
DOCKER := DOCKER_BUILDKIT=1 docker
DESCRIPTION := cli for handling installation and upgrade of Mia-Platform unofficial distribution
DOCKER_LABELS := --label "org.opencontainers.image.title=$(CMDNAME)"
DOCKER_LABELS := $(DOCKER_LABELS) --label "org.opencontainers.image.description=$(DESCRIPTION)"
DOCKER_LABELS := $(DOCKER_LABELS) --label "org.opencontainers.image.url=https://github.com/mia-platform/vab"
DOCKER_LABELS := $(DOCKER_LABELS) --label "org.opencontainers.image.source=https://github.com/mia-platform/vab"
DOCKER_LABELS := $(DOCKER_LABELS) --label "org.opencontainers.image.version=$(VERSION)"
DOCKER_LABELS := $(DOCKER_LABELS) --label "org.opencontainers.image.created=$(COMPLETE_BUILD_DATE)"
DOCKER_LABELS := $(DOCKER_LABELS) --label "org.opencontainers.image.revision=$(shell git rev-parse HEAD)"
DOCKER_LABELS := $(DOCKER_LABELS) --label "org.opencontainers.image.licenses=Apache-2.0"
DOCKER_LABELS := $(DOCKER_LABELS) --label "org.opencontainers.image.documentation=https://docs.mia-platform.eu"
DOCKER_LABELS := $(DOCKER_LABELS) --label "org.opencontainers.image.vendor=Mia s.r.l."

.PHONY: build-image
# Force linux OS because alpine don't have darwin specific slices
build-image: build-image.linux.$(GOARCH)
build-image.%: build.%
	$(eval OS := $(word 1,$(subst ., ,$*)))
	$(eval ARCH := $(word 2,$(subst ., ,$*)))
	@echo "Building image for ${OS} ${ARCH}..."
	$(DOCKER) build --pull --platform $(OS)/$(ARCH) \
		--tag $(REGISTRY)/vab:$(TAG) \
		--file ./Dockerfile bin

verify-buildx:
# For building multiarch images for now we need to use buildx so we stop if we don't find it
	@$(DOCKER) buildx &>/dev/null || { echo "Docker buildx command appear to not be installed."; exit 1; }

CONTEXT_NAME := vab-build-context

.PHONY: build-image-multiarch
build-image-multiarch: verify-buildx build-all
# WARNING FOR NOW IS ONLY WORNING ON DARWIN WITH DOCKER DESKTOP INSTALLED!
	@echo "Building image for following architectures: $(SUPPORTED_PLATFORMS)"
	@$(DOCKER) buildx rm $(CONTEXT_NAME) &>/dev/null || true
	$(DOCKER) buildx create --use --name $(CONTEXT_NAME) --platform "$(SUPPORTED_PLATFORMS)"
	$(DOCKER) buildx build --pull --platform "$(SUPPORTED_PLATFORMS)" \
		$(DOCKER_LABELS) \
		--tag docker.io/miaplatform/vab:$(TAG) \
		--tag ghcr.io/mia-platform/vab:$(TAG) \
		--file ./Dockerfile bin
