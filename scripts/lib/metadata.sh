#!/usr/bin/env bash

THIS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ROOT_DIR="$THIS_DIR/../.."
SCRIPTS_DIR="$ROOT_DIR/scripts"

# get_service_metadata returns the path of the metadata file for a given service
# get_service_metadata accepts the following parameters:
#   $1: service_name
get_service_metadata() {
  if [[ "$#" -ne 1 ]]; then
    echo "[FAIL] Usage: get_service_metadata service_name"
    exit 1
  fi

  local __service_name="$1"

  echo "$ROOT_DIR/../$__service_name-controller/metadata.yaml"
}

# get_service_short_name returns the short name for a given service
# get_service_short_name accepts the following parameters:
#   $1: service_name
get_service_short_name() {
  if [[ "$#" -ne 1 ]]; then
    echo "[FAIL] Usage: get_service_short_name service_name"
    exit 1
  fi

  local __service_name="$1"

  metadata="$(get_service_metadata $__service_name)"
  yq eval '.service.short_name' $metadata
}

# get_service_full_name returns the full name for a given service
# get_service_full_name accepts the following parameters:
#   $1: service_name
get_service_full_name() {
  if [[ "$#" -ne 1 ]]; then
    echo "[FAIL] Usage: get_service_full_name service_name"
    exit 1
  fi

  local __service_name="$1"

  metadata="$(get_service_metadata $__service_name)"
  yq eval '.service.full_name' $metadata
}

# get_service_versions returns the list of api versions supported by the given
# controller
# get_service_versions accepts the following parameters:
#   $1: service_name
get_service_versions() {
  if [[ "$#" -ne 1 ]]; then
    echo "[FAIL] Usage: get_service_versions service_name"
    exit 1
  fi

  local __service_name="$1"

  metadata="$(get_service_metadata $__service_name)"
  yq eval '.versions[].api_version' $metadata
}

# get_latest_version returns the latest api version supported by the given
# controller
# get_latest_version accepts the following parameters:
#   $1: service_name
get_latest_version() {
  if [[ "$#" -ne 1 ]]; then
    echo "[FAIL] Usage: get_latest_version service_name"
    exit 1
  fi

  local __service_name="$1"

  metadata="$(get_service_metadata $__service_name)"
  yq eval '[.versions[].api_version] | .[-1]' $metadata
}