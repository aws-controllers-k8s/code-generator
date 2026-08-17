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
// deterministic order and derives each sync method name from its field name.
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
	assert.Equal("syncLogDeliveryConfigurations", fields[0].CustomSyncMethodName())

	assert.Equal("Tags", fields[1].Names.Camel)
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

// TestCustomSync_CompareIgnored accepts custom_sync on a field that is also
// compare.is_ignored.
//
// `compare.is_ignored` suppresses only the generated comparison. A resource with
// a `delta_pre_compare` hook adds the same path by hand, which is how the
// out-of-band tag pattern this feature generalizes is written in the controllers
// that have it today. Whether the path reaches the delta is therefore a property
// of hand-written code, and the generator does not attempt to judge it.
func TestCustomSync_CompareIgnored(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "elasticache",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-custom-sync-compare-ignored.yaml",
		})

	crds, err := g.GetCRDs()
	require.NoError(err)

	crd := getCRDByName("ReplicationGroup", crds)
	require.NotNil(crd)

	// The config is honored, not silently dropped: the field is still collected
	// as a custom_sync field and still yields a sync method name.
	fields := crd.CustomSyncFields()
	require.Len(fields, 1)
	assert.Equal("Tags", fields[0].Names.Camel)
	assert.Equal("syncTags", fields[0].CustomSyncMethodName())
}

// TestCustomSyncAppliedOnCreate covers the model accessors for
// `applied_on_create`: it partitions the custom_sync fields into those a
// successful create leaves pending and those it does not, while leaving the full
// set - which the update path uses - untouched.
func TestCustomSyncAppliedOnCreate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "elasticache",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-custom-sync-applied-on-create.yaml",
		})

	crds, err := g.GetCRDs()
	require.NoError(err)

	crd := getCRDByName("ReplicationGroup", crds)
	require.NotNil(crd)

	// Both fields are still custom_sync fields, so both are synced on update.
	all := crd.CustomSyncFields()
	require.Len(all, 2)
	assert.Equal("LogDeliveryConfigurations", all[0].Names.Camel)
	assert.False(all[0].CustomSyncAppliedOnCreate())
	assert.Equal("Tags", all[1].Names.Camel)
	assert.True(all[1].CustomSyncAppliedOnCreate())

	// Only the field create does not apply is pending afterwards.
	pending := crd.CustomSyncFieldsPendingAfterCreate()
	require.Len(pending, 1)
	assert.Equal("LogDeliveryConfigurations", pending[0].Names.Camel)
}

// TestCustomSyncAppliedOnCreate_Default confirms the option defaults to false,
// so a field configured with a bare `custom_sync: {}` is still treated as
// pending after create.
func TestCustomSyncAppliedOnCreate_Default(t *testing.T) {
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

	for _, f := range crd.CustomSyncFields() {
		assert.False(f.CustomSyncAppliedOnCreate())
	}
	assert.Len(crd.CustomSyncFieldsPendingAfterCreate(), 2)
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
