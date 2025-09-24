#!/usr/bin/env bash

export CONTROLLER_TOOLS_VERSION="v0.19.0"
export HELM_VERSION="v3.7"

# setting the -x option if debugging is true
if [[ "${DEBUG:-"false"}" = "true" ]]; then
    set -x
fi

# check_is_installed checks to see if the supplied executable is installed and
# exits if not. An optional second argument is an extra message to display when
# the supplied executable is not installed.
#
# Usage:
#
#   check_is_installed PROGRAM [ MSG ]
#
# Example:
#
#   check_is_installed kind "You can install kind with the helper scripts/install-kind.sh"
check_is_installed() {
    local __name="$1"
    local __extra_msg="$2"
    if ! is_installed "$__name"; then
        echo "FATAL: Missing requirement '$__name'"
        echo "Please install $__name before running this script."
        if [[ -n $__extra_msg ]]; then
            echo ""
            echo "$__extra_msg"
            echo ""
        fi
        exit 1
    fi
}

is_installed() {
    local __name="$1"
    if which "$__name" >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

display_timelines() {
    echo ""
    echo "Displaying all step durations."
    echo "TIMELINE: Docker build took $DOCKER_BUILD_DURATION seconds."
    echo "TIMELINE: Upping test cluster took $UP_CLUSTER_DURATION seconds."
    echo "TIMELINE: Base image integration tests took $BASE_INTEGRATION_DURATION seconds."
    echo "TIMELINE: Current image integration tests took $LATEST_INTEGRATION_DURATION seconds."
    echo "TIMELINE: Down processes took $DOWN_DURATION seconds."
}

should_execute() {
  if [[ "$TEST_PASS" -ne 0 ]]; then
    echo "NOTE: Skipping operation '$1'. Test is already marked as failed."
    return 1
  else
    return 0
  fi
}

# filenoext returns just the name of the supplied filename without the
# extension
filenoext() {
    local __name="$1"
    local __filename
    __filename="$( basename "$__name" )"
    # How much do I despise Bash?!
    echo "${__filename%.*}"
}

DEFAULT_DEBUG_PREFIX="DEBUG: "

# debug_msg prints out a supplied message if the DEBUG environs variable is
# set. An optional second argument indicates the "indentation level" for the
# message. If the indentation level argument is missing, we look for the
# existence of an environs variable called "indent_level" and use that.
debug_msg() {
    local __msg=${1:-}
    local __indent_level=${2:-}
    local __debug="${DEBUG:-""}"
    local __debug_prefix="${DEBUG_PREFIX:-$DEFAULT_DEBUG_PREFIX}"
    if [ -z "$__debug" ]; then
        return 0
    fi
    __indent=""
    if [ -n "$__indent_level" ]; then
        __indent="$( for _ in $( seq 0 "$__indent_level" ); do printf " "; done )"
    fi
    echo "$__debug_prefix$__indent$__msg"
}

# k8s_controller_gen_version_equals accepts a string version and returns 0 if the
# installed version of controller-gen matches the supplied version, otherwise
# returns 1
#
# Usage:
#
#   if k8s_controller_gen_version_equals "v0.4.0"; then
#       echo "controller-gen is at version 0.4.0"
#   fi
k8s_controller_gen_version_equals() {
    local currentver
    currentver="$(controller-gen --version | cut -d' ' -f2 | tr -d '\n')";
    local requiredver="$1";
    if [ "$currentver" = "$requiredver" ]; then
        return 0
    else
        return 1
    fi;
}

# helm_version_equals_or_greater accepts a string version and returns 0 if the
# installed version of helm matches or greater than the supplied version, 
# otherwise returns 1
#
# Usage:
#
#   if helm_version_equals_or_greater "v3.9.0"; then
#       echo "Installed helm version is greater than or equal to version v3.9.0"
#   fi
helm_version_equals_or_greater() {
    local currentver
    currentver="$(helm version --template='Version: {{.Version}}'| cut -d' ' -f2 | tr -d '\n')"
    local requiredver="$1"
    printf '%s\n%s\n' "$requiredver" "$currentver" | sort -C -V
    return $?
}

# is_public_ecr_logged_in returns 0 if the Docker client is authenticated
# with ECR public and therefore can pull and push to ECR public, otherwise
# returns 1
#
# Usage:
#
# if ! is_public_ecr_logged_in; then
#   aws ecr-public get-login-password --region us-east-1 \
#   | docker login --username AWS --password-stdin public.ecr.aws
# fi
is_public_ecr_logged_in() {
    local public_ecr_url="public.ecr.aws"

    # Load the auth string
    # Base64 decode it
    # Parse it as <Username>:<B64 Payload>, and take only the payload
    # Base64 decode it
    # Read the "expiration" value
    local expiration_time
    auth_string=$(jq -r --arg url $public_ecr_url '.auths[$url].auth' ~/.docker/config.json)

    [ -z "$auth_string" ] && return 1
    [ "$auth_string" = "null" ] && return 1

    expiration_time=$(echo $auth_string | base64 -d | cut -d":" -f2 | base64 -d | jq -r ".expiration")

    # If any part of this doesn't exist, the user isn't logged in
    [ -z "$expiration_time" ] && return 1

    local current_time
    current_time=$(date +%s)

    # If the credentials have expired, the user isn't logged in
    [ "$expiration_time" -lt "$current_time" ] && return 1

    return 0
}