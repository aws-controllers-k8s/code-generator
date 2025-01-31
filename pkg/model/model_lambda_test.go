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

func TestLambda_Function(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "lambda")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Function", crds)
	require.NotNil(crd)

	assert.Equal("Function", crd.Names.Camel)
	assert.Equal("function", crd.Names.CamelLower)
	assert.Equal("function", crd.Names.Snake)

	// The Lambda Function API has Create, Delete, ReadOne and ReadMany
	// operations, however has no single Update operation. Instead, there are
	// multiple Update operations, depending on the attributes of the function
	// being changed...
	require.NotNil(crd.Ops)

	assert.NotNil(crd.Ops.Create)
	assert.NotNil(crd.Ops.Delete)
	assert.NotNil(crd.Ops.ReadOne)
	assert.NotNil(crd.Ops.ReadMany)

	assert.Nil(crd.Ops.GetAttributes)
	assert.Nil(crd.Ops.SetAttributes)
	assert.Nil(crd.Ops.Update)

	specFields := crd.SpecFields
	statusFields := crd.StatusFields

	expSpecFieldCamel := []string{
		"Code",
		"CodeSigningConfigARN",
		"DeadLetterConfig",
		"Description",
		"Environment",
		"FileSystemConfigs",
		"FunctionName",
		"Handler",
		"ImageConfig",
		"KMSKeyARN",
		"Layers",
		"MemorySize",
		"PackageType",
		"Publish",
		"Role",
		"Runtime",
		"Tags",
		"Timeout",
		"TracingConfig",
		"VPCConfig",
	}
	assert.Equal(expSpecFieldCamel, attrCamelNames(specFields))

	expStatusFieldCamel := []string{
		// Added from generator.yaml
		"Architectures",
		"CodeLocation",
		"CodeRepositoryType",
		"CodeSHA256",
		"CodeSize",
		"EphemeralStorage",
		// "FunctionArn", <-- ACKMetadata.ARN
		"ImageConfigResponse",
		"LastModified",
		"LastUpdateStatus",
		"LastUpdateStatusReason",
		"LastUpdateStatusReasonCode",
		"LoggingConfig",
		"MasterARN",
		"RevisionID",
		"RuntimeVersionConfig",
		"SigningJobARN",
		"SigningProfileVersionARN",
		"SnapStart",
		"State",
		"StateReason",
		"StateReasonCode",
		"Version",
	}
	assert.Equal(expStatusFieldCamel, attrCamelNames(statusFields))
}

func TestLambda_customNestedFields_Spec_Depth2(t *testing.T) {
	// This test is to check if a custom field
	// defined as a nestedField using `type:`,
	// is nested properly inside its parentField

	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "lambda", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-custom-nested-types.yaml",
	})

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Function", crds)
	require.NotNil(crd)

	assert.Contains(crd.SpecFields, "Code")
	codeField := crd.SpecFields["Code"]

	// Check if Nested Field is inside its Parent Field
	assert.Contains(codeField.MemberFields, "S3SHA256")
	assert.Contains(codeField.ShapeRef.Shape.MemberRefs, "S3SHA256")
}
func TestLambda_customNestedFields_Spec_Depth3(t *testing.T) {
	// This test is to check if a custom field
	// defined as a nestedField using `type:`,
	// is nested properly inside its parentField

	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "lambda", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-custom-nested-types.yaml",
	})

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("EventSourceMapping", crds)
	require.NotNil(crd)

	assert.Contains(crd.SpecFields, "DestinationConfig")
	OnSuccessField := crd.SpecFields["DestinationConfig"].MemberFields["OnSuccess"]

	assert.Contains(OnSuccessField.MemberFields, "New")
	assert.Contains(OnSuccessField.ShapeRef.Shape.MemberRefs, "New")
}

func TestLambda_customNestedFields_Status_Depth3(t *testing.T) {
	// This test is to check if a custom field
	// defined as a nestedField using `type:`,
	// is nested properly inside its parentField

	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "lambda", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-custom-nested-types.yaml",
	})

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Function", crds)
	require.NotNil(crd)

	assert.Contains(crd.StatusFields, "ImageConfigResponse")
	ErrorField := crd.StatusFields["ImageConfigResponse"].MemberFields["Error"]

	assert.Contains(ErrorField.MemberFields, "New")
	assert.Contains(ErrorField.ShapeRef.Shape.MemberRefs, "New")
}
