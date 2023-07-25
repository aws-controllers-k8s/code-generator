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

"""Integration tests for the CloudWatch API MetricAlarm resource
"""

import time

import pytest

from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_cloudwatch_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e import condition
from e2e import metric_alarm

RESOURCE_PLURAL = 'metricalarms'

CHECK_STATUS_WAIT_SECONDS = 10
MODIFY_WAIT_AFTER_SECONDS = 10
DELETE_WAIT_AFTER_SECONDS = 5

@pytest.fixture
def _metric_alarm():
    metric_alarm_name = random_suffix_name("ack-test-metric-alarm", 24)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["METRIC_ALARM_NAME"] = metric_alarm_name
    resource_data = load_cloudwatch_resource(
        "metric_alarm",
        additional_replacements=replacements,
    )

    # Create the k8s resource
    ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
        metric_alarm_name, namespace="default",
    )
    k8s.create_custom_resource(ref, resource_data)
    cr = k8s.wait_resource_consumed_by_controller(ref)

    assert cr is not None
    assert k8s.get_resource_exists(ref)

    yield (ref, cr)

    # Try to delete, if doesn't already exist
    _, deleted = k8s.delete_custom_resource(
        ref,
        period_length=DELETE_WAIT_AFTER_SECONDS,
    )
    assert deleted

    metric_alarm.wait_until_deleted(metric_alarm_name)


@service_marker
@pytest.mark.canary
class TestMetricAlarm:
    def test_crud(self, _metric_alarm):
        (ref, cr) = _metric_alarm
        metric_alarm_name = ref.name

        time.sleep(CHECK_STATUS_WAIT_SECONDS)

        condition.assert_synced(ref)

        assert metric_alarm.exists(metric_alarm_name)
        assert k8s.get_resource_exists(ref)
