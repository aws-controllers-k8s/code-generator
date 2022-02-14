SHELL := /bin/bash # Use bash syntax

# Set up variables
GO111MODULE=on

AWS_SERVICE=$(shell echo $(SERVICE) | tr '[:upper:]' '[:lower:]')

# Build ldflags
VERSION ?= "v0.17.0"
GITCOMMIT=$(shell git rev-parse HEAD)
BUILDDATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
IMPORT_PATH=github.com/aws-controllers-k8s/code-generator
GO_LDFLAGS=-ldflags "-X $(IMPORT_PATH)/pkg/version.Version=$(VERSION) \
			-X $(IMPORT_PATH)/pkg/version.BuildHash=$(GITCOMMIT) \
			-X $(IMPORT_PATH)/pkg/version.BuildDate=$(BUILDDATE)"

# We need to use the codegen tag when building and testing because the
# aws-sdk-go/private/model/api package is gated behind a build tag "codegen"...
GO_TAGS=-tags codegen

.PHONY: all build-ack-generate build-controller test \
	build-controller-image local-build-controller-image

all: test

build-ack-generate:	## Build ack-generate binary
	@echo -n "building ack-generate ... "
	@go build ${GO_TAGS} ${GO_LDFLAGS} -o bin/ack-generate cmd/ack-generate/main.go
	@echo "ok."

build-controller: build-ack-generate ## Generate controller code for SERVICE
	@./scripts/install-controller-gen.sh
	@echo "==== building $(AWS_SERVICE)-controller ===="
	@./scripts/build-controller.sh $(AWS_SERVICE)
	@echo "==== building $(AWS_SERVICE)-controller release artifacts ===="
	@./scripts/build-controller-release.sh $(AWS_SERVICE)

build-controller-image: export LOCAL_MODULES = false
build-controller-image:	## Build container image for SERVICE
	@./scripts/build-controller-image.sh $(AWS_SERVICE)

local-build-controller-image: export LOCAL_MODULES = true
local-build-controller-image:	## Build container image for SERVICE allowing local modules
	@./scripts/build-controller-image.sh $(AWS_SERVICE)

check-versions: ## Checks the code-generator version matches the runtime dependency version
	@./scripts/check-versions.sh $(VERSION)

test: 				## Run code tests
	go test ${GO_TAGS} ./...

help:           	## Show this help.
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep | sed -e 's/\\$$//' \
		| awk -F'[:#]' '{print $$1 = sprintf("%-30s", $$1), $$4}'
