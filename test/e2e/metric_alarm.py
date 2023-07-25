# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Utilities for working with Metric Alarm resources"""

import datetime
import time

import boto3
import pytest

DEFAULT_WAIT_UNTIL_DELETED_TIMEOUT_SECONDS = 60*20
DEFAULT_WAIT_UNTIL_DELETED_INTERVAL_SECONDS = 15


def wait_until_deleted(
        metric_alarm_name: str,
        timeout_seconds: int = DEFAULT_WAIT_UNTIL_DELETED_TIMEOUT_SECONDS,
        interval_seconds: int = DEFAULT_WAIT_UNTIL_DELETED_INTERVAL_SECONDS,
    ) -> None:
    """Waits until a Metric Alarm with a supplied name is no longer returned from
    the CloudWatch API.

    Usage:
        from e2e.metric_alarm import wait_until_deleted

        wait_until_deleted(alarm_name)

    Raises:
        pytest.fail upon timeout or if the Metric Alarm goes to any other status
        other than 'deleting'
    """
    now = datetime.datetime.now()
    timeout = now + datetime.timedelta(seconds=timeout_seconds)

    while True:
        if datetime.datetime.now() >= timeout:
            pytest.fail(
                "Timed out waiting for Metric Alarm to be "
                "deleted in CloudWatch API"
            )
        time.sleep(interval_seconds)

        latest = get(metric_alarm_name)
        if latest is None:
            break


def exists(metric_alarm_name):
    """Returns True if the supplied Metric Alarm exists, False otherwise.
    """
    return get(metric_alarm_name) is not None


def get(metric_alarm_name):
    """Returns a dict containing the Metric Alarm record from the CloudWatch API.

    If no such Metric Alarm exists, returns None.
    """
    c = boto3.client('cloudwatch')
    resp = c.describe_alarms(AlarmNames=[metric_alarm_name])
    if len(resp['MetricAlarms']) == 1:
        return resp['MetricAlarms'][0]
    return None


def get_tags(metric_alarm_arn):
    """Returns a dict containing the Metric Alarm's tag records from the
    CloudWatch API.

    If no such Metric Alarm exists, returns None.
    """
    c = boto3.client('cloudwatch')
    try:
        resp = c.list_tags_for_resource(
            ResourceName=metric_alarm_arn,
        )
        return resp['Tags']
    except c.exceptions.ResourceNotFoundException:
        return None
