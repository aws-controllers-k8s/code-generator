#!/usr/bin/env bash

set -eo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
SCRIPTS_DIR=$DIR
ROOT_DIR=$DIR/..

source "$SCRIPTS_DIR/lib/metadata.sh"

USAGE="
Usage:
  $(basename "$0") <aws_service>

Constructs a new metadata.yaml file for a fresh service. 

Example: $(basename "$0") ecr

<aws_service> should be an AWS Service name (ecr, sns, sqs, petstore, bookstore)

Environment variables:
  SERVICE_CONTROLLER_SOURCE_PATH:   Directory to find the service controller to build an image for.
                                    Default: ../\$AWS_SERVICE-controller
"

if [ $# -ne 1 ]; then
    echo "AWS_SERVICE is not defined. Script accepts one parameter, the <AWS_SERVICE> to build a container image of that service" 1>&2
    echo "${USAGE}"
    exit 1
fi

AWS_SERVICE=$(echo "$1" | tr '[:upper:]' '[:lower:]')

# Source code for the controller will be in a separate repo, typically in
# $GOPATH/src/github.com/aws-controllers-k8s/$AWS_SERVICE-controller/
DEFAULT_SERVICE_CONTROLLER_SOURCE_PATH="$ROOT_DIR/../$AWS_SERVICE-controller"
SERVICE_CONTROLLER_SOURCE_PATH=${SERVICE_CONTROLLER_SOURCE_PATH:-$DEFAULT_SERVICE_CONTROLLER_SOURCE_PATH}

if [[ ! -d $SERVICE_CONTROLLER_SOURCE_PATH ]]; then
    echo "Error evaluating SERVICE_CONTROLLER_SOURCE_PATH environment variable:" 1>&2
    echo "$SERVICE_CONTROLLER_SOURCE_PATH is not a directory." 1>&2
    echo "${USAGE}"
    exit 1
fi

METADATA_TEMPLATE_PATH="$ROOT_DIR/templates/metadata.yaml"

DEFAULT_METADATA_OUTPUT_PATH="$SERVICE_CONTROLLER_SOURCE_PATH/metadata.yaml"
METADATA_OUTPUT_PATH="${METADATA_OUTPUT_PATH:-$DEFAULT_METADATA_OUTPUT_PATH}"

echo "🧙 Welcome to the service metadata setup wizard"
echo "⚠️ WARNING: This script will overwrite the metadata.yaml file in your service controller path"

echo -n "Enter the service's full name [e.g. Amazon Elastic Kubernetes Service]: "
read -r full_name

echo -n "Enter the service's acronym [e.g. EKS, S3, EC2]: "
read -r short_name

echo -n "Enter the URL for the service homepage [e.g. https://aws.amazon.com/eks/]: "
read -r link

echo -n "Enter the URL for the service's documentation [e.g. https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html]: "
read -r documentation

echo -n "Generating metadata ... "
write_new_metadata "$METADATA_TEMPLATE_PATH" "$METADATA_OUTPUT_PATH" "$full_name" "$short_name" "$link" "$documentation"
echo "Success!"