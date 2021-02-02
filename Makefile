SHELL := /bin/bash # Use bash syntax

# Set up variables
GO111MODULE=on

AWS_SERVICE=$(shell echo $(SERVICE) | tr '[:upper:]' '[:lower:]')

# Build ldflags
VERSION ?= "v0.0.0"
GITCOMMIT=$(shell git rev-parse HEAD)
BUILDDATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GO_LDFLAGS=-ldflags "-X main.version=$(VERSION) \
			-X main.buildHash=$(GITCOMMIT) \
			-X main.buildDate=$(BUILDDATE)"

# We need to use the codegen tag when building and testing because the
# aws-sdk-go/private/model/api package is gated behind a build tag "codegen"...
GO_TAGS=-tags codegen

.PHONY: all build-ack-generate build-controller test

all: test

build-ack-generate:	## Build ack-generate binary
	@echo -n "building ack-generate ... "
	@go build ${GO_TAGS} ${GO_LDFLAGS} -o bin/ack-generate cmd/ack-generate/main.go
	@echo "ok."

build-controller:   ## Generate controller code for SERVICE
	@./scripts/install-controller-gen.sh 
	@./scripts/build-controller.sh $(AWS_SERVICE)

test: 				## Run code tests
	go test ${GO_TAGS} ./...

help:           	## Show this help.
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep | sed -e 's/\\$$//' \
		| awk -F'[:#]' '{print $$1 = sprintf("%-30s", $$1), $$4}'
