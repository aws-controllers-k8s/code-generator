#!/usr/bin/env bash

# ./scripts/install-helm.sh
#
# Installs the latest version helm if not installed.
#
# NOTE: helm will be installed to /usr/local/bin/helm

set -eo pipefail

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ROOT_DIR="$SCRIPTS_DIR/.."

source "$SCRIPTS_DIR/lib/common.sh"

if ! is_installed helm ; then
    __helm_url="https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3"
    echo -n "installing helm from $__helm_url ... "
    curl --silent "$__helm_url" | bash 1>/dev/null
    echo "ok."
fi
