// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	 http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package code_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

// TestCustomSyncUpdate verifies the sdkUpdate boilerplate for a resource with
// two custom_sync fields. Both are collected into a single DifferentExcept
// call, which is the property that makes generating this worthwhile: a
// hand-written hook has to be widened by hand every time a field is added, and
// forgetting to do so silently short-circuits legitimate updates.
//
// Each method name is derived from its field name, so the calls are rm.syncTags
// and rm.syncLogDeliveryConfigurations.
func TestCustomSyncUpdate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "elasticache",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-custom-sync.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	// Fields are emitted in sorted order, so LogDeliveryConfigurations precedes
	// Tags regardless of the order they appear in generator.yaml.
	expected := `
	updatedDesired := desired.DeepCopy()
	updatedDesired.SetStatus(latest)
	if delta.DifferentAt("Spec.LogDeliveryConfigurations") {
		err = rm.syncLogDeliveryConfigurations(ctx, desired, latest)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.Tags") {
		err = rm.syncTags(ctx, desired, latest)
		if err != nil {
			return nil, err
		}
	}
	if !delta.DifferentExcept("Spec.LogDeliveryConfigurations", "Spec.Tags") {
		return rm.concreteResource(updatedDesired), nil
	}
`
	assert.Equal(
		strings.TrimSpace(expected),
		strings.TrimSpace(code.CustomSyncUpdate(crd, "desired", "latest", "delta", 1)),
	)
}

// TestCustomSyncCreate verifies the post-create marker. The resource is marked
// unsynced only when at least one custom_sync field is actually set, since a
// resource created without any of them has nothing for the follow-up reconcile
// to do. The condition message is generic, so it stays accurate no matter which
// subset of the fields the user populated.
func TestCustomSyncCreate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "elasticache",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-custom-sync.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	expected := `
	if ko.Spec.LogDeliveryConfigurations != nil || ko.Spec.Tags != nil {
		msg := "Secondary sync required; resource will be requeued"
		ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, &msg, nil)
	}
`
	assert.Equal(
		strings.TrimSpace(expected),
		strings.TrimSpace(code.CustomSyncCreate(crd, "ko", 1)),
	)
}

// TestCustomSyncCreate_AppliedOnCreate verifies that a field marked
// `applied_on_create` is left out of the post-create marker. Create already
// carried its value, so there is nothing for the follow-up reconcile to apply,
// and including it would cost a needless requeue while reporting a synced
// resource as unsynced.
//
// Only LogDeliveryConfigurations remains in the guard, so the marker fires on
// exactly the field that is still pending.
func TestCustomSyncCreate_AppliedOnCreate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "elasticache",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-custom-sync-applied-on-create.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	expected := `
	if ko.Spec.LogDeliveryConfigurations != nil {
		msg := "Secondary sync required; resource will be requeued"
		ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, &msg, nil)
	}
`
	assert.Equal(
		strings.TrimSpace(expected),
		strings.TrimSpace(code.CustomSyncCreate(crd, "ko", 1)),
	)
}

// TestCustomSyncUpdate_AppliedOnCreate verifies that `applied_on_create` does
// NOT change the update path. Tags are applied by create but still updated
// out-of-band, so sdkUpdate must keep calling rm.syncTags and must keep counting
// Spec.Tags in the DifferentExcept short-circuit - otherwise a tags-only change
// would fall through to the Update operation, which cannot apply it.
func TestCustomSyncUpdate_AppliedOnCreate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "elasticache",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-custom-sync-applied-on-create.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	expected := `
	updatedDesired := desired.DeepCopy()
	updatedDesired.SetStatus(latest)
	if delta.DifferentAt("Spec.LogDeliveryConfigurations") {
		err = rm.syncLogDeliveryConfigurations(ctx, desired, latest)
		if err != nil {
			return nil, err
		}
	}
	if delta.DifferentAt("Spec.Tags") {
		err = rm.syncTags(ctx, desired, latest)
		if err != nil {
			return nil, err
		}
	}
	if !delta.DifferentExcept("Spec.LogDeliveryConfigurations", "Spec.Tags") {
		return rm.concreteResource(updatedDesired), nil
	}
`
	assert.Equal(
		strings.TrimSpace(expected),
		strings.TrimSpace(code.CustomSyncUpdate(crd, "desired", "latest", "delta", 1)),
	)
}

// TestCustomSyncNoFields confirms both emitters are inert for the overwhelming
// majority of resources, which have no custom_sync fields at all.
func TestCustomSyncNoFields(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")

	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	assert.Equal("", code.CustomSyncUpdate(crd, "desired", "latest", "delta", 1))
	assert.Equal("", code.CustomSyncCreate(crd, "ko", 1))
}
