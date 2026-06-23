#!/usr/bin/env bash

# A script that generates Go e2e test scaffolds for an ACK service controller

set -eo pipefail

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ROOT_DIR="$SCRIPTS_DIR/.."

source "$SCRIPTS_DIR/lib/common.sh"

ACK_GENERATE_CACHE_DIR=${ACK_GENERATE_CACHE_DIR:-"$HOME/.cache/aws-controllers-k8s"}
DEFAULT_ACK_GENERATE_BIN_PATH="$ROOT_DIR/bin/ack-generate"
ACK_GENERATE_BIN_PATH=${ACK_GENERATE_BIN_PATH:-$DEFAULT_ACK_GENERATE_BIN_PATH}
ACK_GENERATE_CONFIG_PATH=${ACK_GENERATE_CONFIG_PATH:-""}
ACK_METADATA_CONFIG_PATH=${ACK_METADATA_CONFIG_PATH:-""}
AWS_SDK_GO_VERSION=${AWS_SDK_GO_VERSION:-""}

USAGE="
Usage:
  $(basename "$0") <service>

<service> should be an AWS service API alias that you wish to generate e2e tests
for -- e.g. 's3' 'sns' or 'sqs'

Environment variables:
  ACK_GENERATE_CACHE_DIR:               Overrides the directory used for caching AWS API
                                        models used by the ack-generate tool.
                                        Default: $ACK_GENERATE_CACHE_DIR
  ACK_GENERATE_BIN_PATH:                Overrides the path to the ack-generate binary.
                                        Default: $ACK_GENERATE_BIN_PATH
  ACK_GENERATE_CONFIG_PATH:             Specify a path to the generator config YAML file.
                                        Default: generator.yaml in controller source path
  ACK_METADATA_CONFIG_PATH:             Specify a path to the metadata config YAML file.
                                        Default: metadata.yaml in controller source path
  AWS_SDK_GO_VERSION:                   Overrides the version of github.com/aws/aws-sdk-go used
                                        by 'ack-generate' to fetch the service API Specifications.
                                        Default: Version of aws/aws-sdk-go in service go.mod
  SERVICE_CONTROLLER_SOURCE_PATH:       Path to the service controller source directory.
                                        Default: ../<service>-controller
  TEMPLATE_DIRS:                        Overrides the list of directories containing ack-generate
                                        templates.
  TEST_CONFIG_PATH:                     Path to testconfig.yaml.
                                        Default: testconfig.yaml in controller source path
"

if [ $# -ne 1 ]; then
    echo "ERROR: $(basename "$0") only accepts a single parameter" 1>&2
    echo "$USAGE"
    exit 1
fi

if [ ! -f "$ACK_GENERATE_BIN_PATH" ]; then
    if is_installed "ack-generate"; then
        ACK_GENERATE_BIN_PATH=$(which "ack-generate")
    else
        echo "ERROR: Unable to find an ack-generate binary.
Either set the ACK_GENERATE_BIN_PATH to a valid location or
run:

   make build-ack-generate

from the root directory or install ack-generate using:

   go get -u -tags codegen github.com/aws-controllers-k8s/code-generator/cmd/ack-generate" 1>&2
        exit 1;
    fi
fi
SERVICE=$(echo "$1" | tr '[:upper:]' '[:lower:]')

DEFAULT_SERVICE_CONTROLLER_SOURCE_PATH="$ROOT_DIR/../$SERVICE-controller"
SERVICE_CONTROLLER_SOURCE_PATH=${SERVICE_CONTROLLER_SOURCE_PATH:-$DEFAULT_SERVICE_CONTROLLER_SOURCE_PATH}

if [[ ! -d $SERVICE_CONTROLLER_SOURCE_PATH ]]; then
    echo "Error evaluating SERVICE_CONTROLLER_SOURCE_PATH environment variable:" 1>&2
    echo "$SERVICE_CONTROLLER_SOURCE_PATH is not a directory." 1>&2
    echo "${USAGE}"
    exit 1
fi

DEFAULT_TEMPLATE_DIRS="$ROOT_DIR/templates"
if [[ -d "$SERVICE_CONTROLLER_SOURCE_PATH/templates" ]]; then
    DEFAULT_TEMPLATE_DIRS="$SERVICE_CONTROLLER_SOURCE_PATH/templates,$DEFAULT_TEMPLATE_DIRS"
fi
TEMPLATE_DIRS=${TEMPLATE_DIRS:-$DEFAULT_TEMPLATE_DIRS}

# Determine the API version from the controller source
ACK_GENERATE_API_VERSION=${ACK_GENERATE_API_VERSION:-"v1alpha1"}
if [[ -d "$SERVICE_CONTROLLER_SOURCE_PATH/apis" ]]; then
    LATEST_API_VERSION=$(ls "$SERVICE_CONTROLLER_SOURCE_PATH/apis/" | sort -V | tail -1)
    if [[ -n "$LATEST_API_VERSION" ]]; then
        ACK_GENERATE_API_VERSION="$LATEST_API_VERSION"
    fi
fi

# Build args for ack-generate e2e-tests
ag_args=("e2e-tests" "$SERVICE" -o "$SERVICE_CONTROLLER_SOURCE_PATH" --template-dirs "$TEMPLATE_DIRS")

if [ -n "$ACK_GENERATE_CACHE_DIR" ]; then
    ag_args=("${ag_args[@]}" --cache-dir "$ACK_GENERATE_CACHE_DIR")
fi

if [ -n "$ACK_GENERATE_CONFIG_PATH" ]; then
    ag_args=("${ag_args[@]}" --generator-config-path "$ACK_GENERATE_CONFIG_PATH")
elif [ -f "$SERVICE_CONTROLLER_SOURCE_PATH/generator.yaml" ]; then
    ag_args=("${ag_args[@]}" --generator-config-path "$SERVICE_CONTROLLER_SOURCE_PATH/generator.yaml")
fi

if [ -n "$ACK_METADATA_CONFIG_PATH" ]; then
    ag_args=("${ag_args[@]}" --metadata-config-path "$ACK_METADATA_CONFIG_PATH")
elif [ -f "$SERVICE_CONTROLLER_SOURCE_PATH/metadata.yaml" ]; then
    ag_args=("${ag_args[@]}" --metadata-config-path "$SERVICE_CONTROLLER_SOURCE_PATH/metadata.yaml")
fi

if [ -n "$AWS_SDK_GO_VERSION" ]; then
    ag_args=("${ag_args[@]}" --aws-sdk-go-version "$AWS_SDK_GO_VERSION")
fi

TEST_CONFIG_PATH=${TEST_CONFIG_PATH:-"$SERVICE_CONTROLLER_SOURCE_PATH/testconfig.yaml"}
if [ -f "$TEST_CONFIG_PATH" ]; then
    ag_args=("${ag_args[@]}" --test-config "$TEST_CONFIG_PATH")
else
    echo "ERROR: testconfig.yaml not found at $TEST_CONFIG_PATH" 1>&2
    echo "Create a testconfig.yaml to define test values for resources." 1>&2
    exit 1
fi

echo "Generating e2e tests for $SERVICE"
$ACK_GENERATE_BIN_PATH "${ag_args[@]}"

pushd "$SERVICE_CONTROLLER_SOURCE_PATH" 1>/dev/null
gofmt -w .
popd 1>/dev/null

echo "e2e test generation complete."
