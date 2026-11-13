// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package model_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func fieldNames(fields []*ackmodel.Field) []string {
	names := make([]string, len(fields))
	for i, f := range fields {
		names[i] = f.Names.Camel
	}
	sort.Strings(names)
	return names
}

func TestFieldGroupOperations_UpdateAutoDetect(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-field-groups.yaml",
	})

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Repository", crds)
	require.NotNil(crd)

	// Should have 2 update field groups
	assert.True(crd.HasFieldGroupUpdates())
	require.Len(crd.UpdateFieldGroups, 2)

	// PutImageScanningConfiguration
	fg0 := crd.UpdateFieldGroups[0]
	assert.Equal("PutImageScanningConfiguration", fg0.OperationID)
	assert.Equal(ackmodel.FieldGroupOpTypeUpdate, fg0.OpType)
	assert.NotNil(fg0.Operation)
	assert.True(fg0.Config.RequeueOnSuccess)

	// Identifier fields: RegistryId and RepositoryName (shared with Delete)
	idNames0 := fieldNames(fg0.IdentifierFields)
	assert.Equal([]string{"RegistryID", "RepositoryName"}, idNames0)

	// Payload field: ImageScanningConfiguration
	payloadNames0 := fieldNames(fg0.PayloadFields)
	assert.Equal([]string{"ImageScanningConfiguration"}, payloadNames0)

	// PutImageTagMutability
	fg1 := crd.UpdateFieldGroups[1]
	assert.Equal("PutImageTagMutability", fg1.OperationID)
	assert.False(fg1.Config.RequeueOnSuccess)

	idNames1 := fieldNames(fg1.IdentifierFields)
	assert.Equal([]string{"RegistryID", "RepositoryName"}, idNames1)

	payloadNames1 := fieldNames(fg1.PayloadFields)
	assert.Equal([]string{"ImageTagMutability"}, payloadNames1)
}

func TestFieldGroupOperations_ReadAutoDetect(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-field-groups.yaml",
	})

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Repository", crds)
	require.NotNil(crd)

	// Should have 1 read field group
	assert.True(crd.HasFieldGroupReads())
	require.Len(crd.ReadFieldGroups, 1)

	fg := crd.ReadFieldGroups[0]
	assert.Equal("GetRepositoryPolicy", fg.OperationID)
	assert.Equal(ackmodel.FieldGroupOpTypeRead, fg.OpType)
	assert.NotNil(fg.Operation)

	// Identifier fields from Input: RegistryId and RepositoryName
	idNames := fieldNames(fg.IdentifierFields)
	assert.Equal([]string{"RegistryID", "RepositoryName"}, idNames)

	// Payload fields from Output: PolicyText is not a CRD field, so
	// auto-detection yields no payload. In practice, a read operation
	// like this would use explicit `fields` config or have matching CRD fields.
	assert.Len(fg.PayloadFields, 0)
}

func TestFieldGroupOperations_NoConfig(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	// Default ECR generator.yaml has no field group operations
	g := testutil.NewModelForService(t, "ecr")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Repository", crds)
	require.NotNil(crd)

	assert.False(crd.HasFieldGroupUpdates())
	assert.False(crd.HasFieldGroupReads())
	assert.Len(crd.UpdateFieldGroups, 0)
	assert.Len(crd.ReadFieldGroups, 0)
}

func TestFieldGroupOperations_PayloadFieldNames(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-field-groups.yaml",
	})

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Repository", crds)
	require.NotNil(crd)

	// All update payload field names across both groups
	updatePayloadNames := crd.FieldGroupPayloadFieldNames(ackmodel.FieldGroupOpTypeUpdate)
	assert.Equal([]string{"ImageScanningConfiguration", "ImageTagMutability"}, updatePayloadNames)

	// Read payload has no CRD-matching fields in this test
	readPayloadNames := crd.FieldGroupPayloadFieldNames(ackmodel.FieldGroupOpTypeRead)
	assert.Len(readPayloadNames, 0)
}
