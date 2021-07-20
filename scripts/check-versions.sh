#!/usr/bin/env bash

# A script that ensures the code-generator binary version matches its associated
# runtime version

set -eo pipefail

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ROOT_DIR="$SCRIPTS_DIR/.."

DEFAULT_GO_MOD_PATH="$ROOT_DIR/go.mod"
GO_MOD_PATH=${GO_MOD_PATH:-$DEFAULT_GO_MOD_PATH}

USAGE="
Usage:
  $(basename "$0") <Makefile version>

Checks that the code-generator version located in the Makefile matches the
runtime version pinned in go.mod

Environment variables:
  GO_MOD_PATH:              Overrides the path to the go.mod file containing the
                            pinned runtime version.
                            Default: $DEFAULT_GO_MOD_PATH
"

if [ $# -ne 1 ]; then
    echo "ERROR: $(basename "$0") only accepts a single parameter" 1>&2
    echo "$USAGE"
    exit 1
fi

RUNTIME_DEPENDENCY_VERSION="aws-controllers-k8s/runtime"
GO_MOD_VERSION=$(go list -m -f '{{ .Version }}' github.com/aws-controllers-k8s/runtime)

if [[ "$1" != "$GO_MOD_VERSION" ]]; then
    echo "Code-generator version in Makefile $1 does not match $RUNTIME_DEPENDENCY_VERSION version $GO_MOD_VERSION"
    exit 1
fi

exit 0