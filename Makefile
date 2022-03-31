SHELL := /bin/bash # Use bash syntax

# Set up variables
GO111MODULE=on

AWS_SERVICE=$(shell echo $(SERVICE) | tr '[:upper:]' '[:lower:]')

# Build ldflags
VERSION ?= $(shell git describe --tags --always --dirty)
GITCOMMIT=$(shell git rev-parse HEAD)
BUILDDATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
IMPORT_PATH=github.com/aws-controllers-k8s/code-generator
GO_LDFLAGS=-ldflags "-X $(IMPORT_PATH)/pkg/version.Version=$(VERSION) \
			-X $(IMPORT_PATH)/pkg/version.BuildHash=$(GITCOMMIT) \
			-X $(IMPORT_PATH)/pkg/version.BuildDate=$(BUILDDATE)"

# We need to use the codegen tag when building and testing because the
# aws-sdk-go/private/model/api package is gated behind a build tag "codegen"...
GO_CMD_FLAGS=-tags codegen
GO_CMD_LOCAL_FLAGS=-modfile=go.local.mod $(GO_CMD_FLAGS)

.PHONY: all local-build-ack-generate build-ack-generate local-build-controller \
	build-controller test local-test build-controller-image \
	local-build-controller-image

all: test

build-ack-generate:	## Build ack-generate binary
	@echo -n "building ack-generate ... "
	@go build ${GO_CMD_FLAGS} ${GO_LDFLAGS} -o bin/ack-generate cmd/ack-generate/main.go
	@echo "ok."

local-build-ack-generate:	## Build ack-generate binary using the local go.mod
	@echo -n "building ack-generate ... "
	@go build ${GO_CMD_LOCAL_FLAGS} ${GO_LDFLAGS} -o bin/ack-generate cmd/ack-generate/main.go
	@echo "ok."

build-controller: build-ack-generate ## Generate controller code for SERVICE
	@./scripts/install-controller-gen.sh
	@echo "==== building $(AWS_SERVICE)-controller ===="
	@./scripts/build-controller.sh $(AWS_SERVICE)
	@echo "==== building $(AWS_SERVICE)-controller release artifacts ===="
	@./scripts/build-controller-release.sh $(AWS_SERVICE)

local-build-controller: local-build-ack-generate build-controller ## Generate controller code for SERVICE using the local go.mod

build-controller-image: export LOCAL_MODULES = false
build-controller-image:	## Build container image for SERVICE
	@./scripts/build-controller-image.sh $(AWS_SERVICE)

local-build-controller-image: export LOCAL_MODULES = true
local-build-controller-image:	## Build container image for SERVICE allowing local modules
	@./scripts/build-controller-image.sh $(AWS_SERVICE)

test: 				## Run code tests
	go test ${GO_CMD_FLAGS} ./...

local-test:			## Run code tests using the local go.mod
	go test ${GO_CMD_LOCAL_FLAGS} ./...

help:           	## Show this help.
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep | sed -e 's/\\$$//' \
		| awk -F'[:#]' '{print $$1 = sprintf("%-30s", $$1), $$4}'
