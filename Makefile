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
VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
DATE_FMT = +%Y-%m-%d
ifdef SOURCE_DATE_EPOCH
    BUILD_DATE ?= $(shell date -u -d "@$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u -r "$(SOURCE_DATE_EPOCH)" "$(DATE_FMT)" 2>/dev/null || date -u "$(DATE_FMT)")
else
    BUILD_DATE ?= $(shell date "$(DATE_FMT)")
endif

GO_LDFLAGS := -X github.com/mia-platform/vab/internal/cmd.Version=$(VERSION) $(GO_LDFLAGS)
GO_LDFLAGS := -X github.com/mia-platform/vab/internal/cmd.BuildDate=$(BUILD_DATE) $(GO_LDFLAGS)

.PHONY: build build-all
build: build.${GOOS}.${GOARCH}
build-linux-amd64: build.linux.amd64
build-linux-arm64: build.linux.arm64
build-darwin-amd64: build.darwin.amd64
build-darwin-arm64: build.darwin.arm64
build-all: build-linux-amd64 build-linux-arm64 build.darwin.amd64 build.darwin.arm64

build.%:
	$(eval OS := $(word 1,$(subst ., ,$*)))
	$(eval ARCH := $(word 2,$(subst ., ,$*)))
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -ldflags "$(GO_LDFLAGS)" \
		-o $(OUTPUT_DIR)/$(OS)/$(ARCH)/$(CMDNAME) $(PROJECT_DIR)/cmd/$(CMDNAME)

##@ Test

TEST_VERBOSE ?= "false"
.PHONY: test
test:
ifneq ($(TEST_VERBOSE), "false")
	go test -test.v ./...
else
	go test ./...
endif

.PHONY: test-coverage
test-coverage:
	go test ./... -race -coverprofile=coverage.xml -covermode=atomic

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
	rm -fr .tools/bin

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
	@GOBIN=$(TOOLS_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2

generate-dep:
	@GOBIN=$(TOOLS_BIN) go install k8s.io/code-generator/cmd/deepcopy-gen@v0.24.2

##@ Images

# REGISTRY is the image registry to use for build and push image targets, default to docker hub
REGISTRY ?= docker.io/miaplatform
IMAGE = ${REGISTRY}/vab

# TAG is the tag to use for build and push image targets, use git tag or latest
TAG ?= $(shell git describe --tags 2>/dev/null || echo latest)

# Force to use buildkit as engine
DOCKER := DOCKER_BUILDKIT=1 docker

build-image: build.linux.${GOARCH}
	@echo "Building image for ${GOARCH}..."
# Force linux OS because alpine don't have darwin specific slices
	DOCKER build --platform linux/${GOARCH} --pull -t $(IMAGE):$(TAG) -f . bin
