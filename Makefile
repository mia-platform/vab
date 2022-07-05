# Copyright 2022 Mia-Platform

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

# 	http://www.apache.org/licenses/LICENSE-2.0

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

.PHONY: build build-linux-amd64 build-linux-arm64 build-all
build: build.${GOOS}.${GOARCH}
build-linux-amd64: build.linux.amd64
build-linux-arm64: build.linux.arm64
build-all: build-linux-amd64 build-linux-arm64

build.%:
	$(eval OS := $(word 1,$(subst ., ,$*)))
	$(eval ARCH := $(word 2,$(subst ., ,$*)))
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build \
		-o $(OUTPUT_DIR)/$(OS)/$(ARCH)/$(CMDNAME) $(PROJECT_DIR)/cmd/$(CMDNAME)

##@ Test

.PHONY: test
test:
	go test ./...

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

generate-dep:
	@GOBIN=$(TOOLS_BIN) go install k8s.io/code-generator/cmd/deepcopy-gen@v0.24.2
