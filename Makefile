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
build: bin/${GOOS}/${GOARCH}/vab
build-linux-amd64: bin/linux/amd64/vab
build-linux-arm64: bin/linux/arm64/vab
build-all: build-linux-amd64 build-linux-arm64

bin/%/vab:
	CGO_ENABLED=0 GOOS=$(word 1,$(subst /, ,$*)) GOARCH=$(word 2,$(subst /, ,$*)) go build -o $@ ./cmd/vab

.PHONY: test
test:
	go test ./...
