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

package code_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestSetSDKFieldGroup_ECR_PutImageScanningConfiguration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-field-groups.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	require.True(crd.HasFieldGroupUpdates())
	require.Len(crd.UpdateFieldGroups, 2)

	fg := crd.UpdateFieldGroups[0]
	require.Equal("PutImageScanningConfiguration", fg.OperationID)

	got, err := code.SetSDKFieldGroup(crd.Config(), crd, fg, "r.ko", "res", 1)
	require.NoError(err)

	// Should set identifier fields (RegistryId, RepositoryName) and
	// payload field (ImageScanningConfiguration)
	assert.Contains(got, "r.ko.Spec.ImageScanningConfiguration")
	assert.Contains(got, "r.ko.Spec.RegistryID")
	assert.Contains(got, "r.ko.Spec.RepositoryName")
	assert.Contains(got, "res.")

	// Should NOT contain Tags (not part of this field group)
	assert.NotContains(got, "Tags")
	// Should NOT contain ImageTagMutability (belongs to different field group)
	assert.NotContains(got, "ImageTagMutability")
}

func TestSetSDKFieldGroup_ECR_PutImageTagMutability(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-field-groups.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)

	fg := crd.UpdateFieldGroups[1]
	require.Equal("PutImageTagMutability", fg.OperationID)

	got, err := code.SetSDKFieldGroup(crd.Config(), crd, fg, "r.ko", "res", 1)
	require.NoError(err)

	assert.Contains(got, "r.ko.Spec.ImageTagMutability")
	assert.Contains(got, "r.ko.Spec.RegistryID")
	assert.Contains(got, "r.ko.Spec.RepositoryName")
	// Should NOT contain ImageScanningConfiguration
	assert.NotContains(got, "ImageScanningConfiguration")
}

func TestSetResourceFieldGroup_ECR_PutImageScanningConfiguration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-field-groups.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)

	fg := crd.UpdateFieldGroups[0]
	require.Equal("PutImageScanningConfiguration", fg.OperationID)

	got, err := code.SetResourceFieldGroup(crd.Config(), crd, fg, "resp", "ko", 1)
	require.NoError(err)

	// Should set the payload field from output
	assert.Contains(got, "resp.ImageScanningConfiguration")
	assert.Contains(got, "ko.Spec.ImageScanningConfiguration")

	// Should NOT contain identifier fields in output (they're not payload)
	// RegistryId and RepositoryName are identifiers, not payload
	assert.NotContains(got, "ko.Spec.RegistryID")
	assert.NotContains(got, "ko.Spec.RepositoryName")
}

func TestFieldGroupDeltaCheck_SingleField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-field-groups.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)

	// PutImageScanningConfiguration has one payload field
	fg := crd.UpdateFieldGroups[0]
	got := code.FieldGroupDeltaCheck(crd.Config(), crd, fg, "delta")

	expected := `delta.DifferentAt("Spec.ImageScanningConfiguration")`
	assert.Equal(
		strings.TrimSpace(expected),
		strings.TrimSpace(got),
	)
}

func TestFieldGroupDeltaCheck_EmptyPayload(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-field-groups.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)

	// Read operation has no payload fields (PolicyText isn't a CRD field)
	fg := crd.ReadFieldGroups[0]
	got := code.FieldGroupDeltaCheck(crd.Config(), crd, fg, "delta")

	assert.Equal("false", got)
}
