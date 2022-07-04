# General variables

# Set Output Directory Path
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(shell pwd -P)/bin
endif
CMDNAME := vab

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
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o $(OUTPUT_DIR)/$(OS)/$(ARCH)/$(CMDNAME) ./cmd/$(CMDNAME)

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

.PHONY: clean-all
clean-all: clean clean-go
