#!/usr/bin/env bash

DEFAULT_AWS_CLI_VERSION="2.0.52"

# daws() executes the AWS Python CLI tool from a Docker container.
#
# Instead of relying on developers having a particular version of the AWS
# Python CLI tool, this method allows a specific version of the CLI tool to be
# executed within a Docker container.
#
# You call the daws function just like you were calling the `aws` CLI tool.
#
# Usage:
#
#   daws SERVICE COMMAND [OPTIONS]
#
# Example:
#
#   daws ecr describe-repositories --repository-name my-repo
#
# To use a specific version of the AWS CLI, set the ACK_AWS_CLI_IMAGE_VERSION
# environment variable, otherwise the value of DEFAULT_AWS_CLI_VERSION is used.
daws() {
    aws_cli_profile_env=()
    if [ -n "$AWS_PROFILE" ]; then
        aws_cli_profile_env=("--env AWS_PROFILE=$AWS_PROFILE")
    fi
    aws_cli_img_version=${ACK_AWS_CLI_IMAGE_VERSION:-$DEFAULT_AWS_CLI_VERSION}
    aws_cli_img="amazon/aws-cli:$aws_cli_img_version"
    docker run --rm -v ~/.aws:/root/.aws:z "${aws_cli_profile_env[@]}" -v "$(pwd)":/aws "$aws_cli_img" "$@"
}

# aws_check_credentials() calls the STS::GetCallerIdentity API call and
# verifies that there is a local identity for running AWS commands
aws_check_credentials() {
    echo -n "checking AWS credentials ... "
    daws sts get-caller-identity --query "Account" >/dev/null ||
        ( printf "\nFATAL: No AWS credentials found. Please run \`aws configure\` to set up the CLI for your credentials." && exit 1)
    echo "ok."
}

aws_account_id() {
    JSON=$(daws sts get-caller-identity --output json || exit 1)
    echo "${JSON}" | jq --raw-output ".Account"
}
