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

.PHONY: all build-ack-generate build-controller test \
	build-controller-image \
	local-build-controller-image lint-shell \
	check-crd-compatibility

all: test

build-ack-generate:	## Build ack-generate binary
	@echo -n "building ack-generate ... "
	@go build ${GO_CMD_FLAGS} ${GO_LDFLAGS} -o bin/ack-generate cmd/ack-generate/main.go
	@echo "ok."

CONTROLLER_PATH ?= ../$(AWS_SERVICE)-controller
RELEASE_VERSION ?= $(shell cd $(CONTROLLER_PATH) && git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0-non-release-version")

ACK_GENERATE_OLM ?= false
OLM_CONFIG_PATH ?= $(CONTROLLER_PATH)/olm/olmconfig.yaml

build-controller: build-ack-generate	## Generate a service controller (SERVICE=s3)
	@$(CURDIR)/bin/ack-generate controller $(AWS_SERVICE) -o $(CONTROLLER_PATH)
	@$(CURDIR)/bin/ack-generate release $(AWS_SERVICE) $(RELEASE_VERSION) -o $(CONTROLLER_PATH)
ifeq ($(ACK_GENERATE_OLM),true)
	@$(CURDIR)/bin/ack-generate olm $(AWS_SERVICE) $(RELEASE_VERSION) -o $(CONTROLLER_PATH) --olm-config $(OLM_CONFIG_PATH)
endif

build-controller-image: export LOCAL_MODULES = false
build-controller-image:	## Build container image for SERVICE
	@./scripts/build-controller-image.sh $(AWS_SERVICE)

local-build-controller-image: export LOCAL_MODULES = true
local-build-controller-image:	## Build container image for SERVICE allowing local modules
	@./scripts/build-controller-image.sh $(AWS_SERVICE)

BASE_REF ?= main
CRD_PATHS ?= config/crd/bases,helm/crds

check-crd-compatibility: build-ack-generate	## Check CRDs for breaking changes against BASE_REF
	@bin/ack-generate crd-compat-check --base-ref=$(BASE_REF) --crd-paths=$(CRD_PATHS)

test: 				## Run code tests
	go test ${GO_CMD_FLAGS} ./...

lint-shell:	## Run linters against all of the bash scripts
	@find . -type f -name "*.sh" | xargs shellcheck -e SC1091

help:           	## Show this help.
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep | sed -e 's/\\$$//' \
		| awk -F'[:#]' '{print $$1 = sprintf("%-30s", $$1), $$4}'
