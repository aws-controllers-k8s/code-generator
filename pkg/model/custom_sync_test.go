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

package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

// TestCustomSyncFields verifies that the model resolves custom_sync fields in a
// deterministic order and derives the sync method name correctly, both from the
// field name and from an explicit `method` override.
func TestCustomSyncFields(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "elasticache",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-custom-sync.yaml",
		})

	crds, err := g.GetCRDs()
	require.NoError(err)

	crd := getCRDByName("ReplicationGroup", crds)
	require.NotNil(crd)

	assert.True(crd.HasCustomSyncFields())

	fields := crd.CustomSyncFields()
	require.Len(fields, 2)

	// Sorted by field name so generated output is stable across runs.
	assert.Equal("LogDeliveryConfigurations", fields[0].Names.Camel)
	assert.Equal("syncLogDelivery", fields[0].CustomSyncMethodName())

	assert.Equal("Tags", fields[1].Names.Camel)
	// No `method` override, so the name is derived from the field name.
	assert.Equal("syncTags", fields[1].CustomSyncMethodName())
}

// TestCustomSyncFields_NotConfigured confirms the accessors are inert for a
// resource with no custom_sync fields, which is the case for nearly every
// resource.
func TestCustomSyncFields_NotConfigured(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")

	crds, err := g.GetCRDs()
	require.NoError(err)

	crd := getCRDByName("ReplicationGroup", crds)
	require.NotNil(crd)

	assert.False(crd.HasCustomSyncFields())
	assert.Empty(crd.CustomSyncFields())

	// A field without the config reports no sync method.
	tags := crd.SpecFields["Tags"]
	require.NotNil(tags)
	assert.False(tags.HasCustomSync())
	assert.Equal("", tags.CustomSyncMethodName())
}

// TestCustomSyncInvalid_CompareIgnored rejects custom_sync on an ignored field.
// Such a field never enters the delta, so the sync would never run and the
// generated DifferentExcept short-circuit would fire on every reconcile,
// stopping the resource from updating at all.
func TestCustomSyncInvalid_CompareIgnored(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "elasticache",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-custom-sync-compare-ignored.yaml",
		})

	_, err := g.GetCRDs()
	require.Error(err)
	assert.Contains(err.Error(), "resources.ReplicationGroup.fields.Tags.custom_sync")
	assert.Contains(err.Error(), "compare.is_ignored")
}

// TestCustomSyncInvalid_NestedField rejects custom_sync on a nested field. The
// emitted code builds a "Spec.<Field>" delta path and a nil check directly off
// ko.Spec, neither of which is correct below the top level.
func TestCustomSyncInvalid_NestedField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "elasticache",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-custom-sync-nested.yaml",
		})

	_, err := g.GetCRDs()
	require.Error(err)
	assert.Contains(err.Error(), "only supported on top-level Spec fields")
	assert.Contains(err.Error(), "LogDeliveryConfigurations.DestinationType")
}

// TestCustomSyncInvalid_ReadOnly rejects custom_sync on a Status field, where
// there is no user-supplied desired value to sync toward.
func TestCustomSyncInvalid_ReadOnly(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "elasticache",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-custom-sync-read-only.yaml",
		})

	_, err := g.GetCRDs()
	require.Error(err)
	assert.Contains(err.Error(), "is_read_only")
}
