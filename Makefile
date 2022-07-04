# Golang variables
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

##@ Build

.DEFAULT_GOAL = build

.PHONY: build build-linux-amd64 build-linux-arm64 build-all
build: build.${GOOS}.${GOARCH}
build-linux-amd64: build.linux.amd64
build-linux-arm64: build.linux.arm64
build-all: build-linux-amd64 build-linux-arm64

build.%:
	$(eval OS := $(word 1,$(subst ., ,$*)))
	$(eval ARCH := $(word 2,$(subst ., ,$*)))
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) go build -o bin/$(OS)/$(ARCH)/vab ./cmd/vab

.PHONY: test
test:
	go test ./...
