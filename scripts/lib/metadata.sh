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

# write_new_metadata will write all fields of a new metadata to file 
# write_new_metadata accepts the following parameters:
#    $1: template path
#    $2: output path
#    $3: full_name
#    $4: short_name
#    $5: homepage link
#    $6: documentation link
write_new_metadata() {
  if [[ "$#" -ne 6 ]]; then
    echo "[FAIL] Usage: write_new_metadata template_path output_path full_name short_name homepage_link documentation_link"
    exit 1
  fi

  __template_path="$1"
  __output_path="$2"
  export __full_name="$3"
  export __short_name="$4"
  export __link="$5"
  export __documentation="$6"

  cp "$__template_path" "$__output_path"

  yq e '.service.full_name=env(__full_name)' -i $__output_path
  yq e '.service.short_name=env(__short_name)' -i $__output_path
  yq e '.service.link=env(__link)' -i $__output_path
  yq e '.service.documentation=env(__documentation)' -i $__output_path
}