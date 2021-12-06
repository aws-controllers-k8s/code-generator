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
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReferenceFieldsValidation_NoReferenceConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-reference.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "Api")
	require.NotNil(crd)
	expected := ""
	assert.Equal(expected, code.ReferenceFieldsValidation(crd, "ko", "Ref", 1))
}

func Test_ReferenceFieldsValidation_SingleReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Integration")
	require.NotNil(crd)
	expected :=
		`	if ko.Spec.APIIDRef != nil && ko.Spec.APIID != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("APIID", "APIIDRef")
	}
	if ko.Spec.APIIDRef == nil && ko.Spec.APIID == nil {
		return ackerr.ResourceReferenceOrIDRequiredFor("APIID", "APIIDRef")
	}
`
	assert.Equal(expected, code.ReferenceFieldsValidation(crd, "ko", "Ref", 1))
}

func Test_ReferenceFieldsValidation_SliceOfReferences(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	//NOTE: For the moment, we are substituting SecurityGroupId with ApiId
	// just to test code generation for slices of reference
	crd := testutil.GetCRDByName(t, g, "VpcLink")
	require.NotNil(crd)
	expected :=
		`	if ko.Spec.SecurityGroupIDsRef != nil && ko.Spec.SecurityGroupIDs != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("SecurityGroupIDs", "SecurityGroupIDsRef")
	}
`
	assert.Equal(expected, code.ReferenceFieldsValidation(crd, "ko", "Ref", 1))
}

func Test_ReferenceFieldsPresent_NoReferenceConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Api")
	require.NotNil(crd)
	expected := "false"
	assert.Equal(expected, code.ReferenceFieldsPresent(crd, "ko"))
}

func Test_ReferenceFieldsPresent_SingleReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Integration")
	require.NotNil(crd)
	expected := "false || ko.Spec.APIIDRef != nil"
	assert.Equal(expected, code.ReferenceFieldsPresent(crd, "ko"))
}

func Test_ReferenceFieldsPresent_SliceOfReferences(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	//NOTE: For the moment, we are substituting SecurityGroupId with ApiId
	// just to test code generation for slices of reference
	crd := testutil.GetCRDByName(t, g, "VpcLink")
	require.NotNil(crd)
	expected := "false || ko.Spec.SecurityGroupIDsRef != nil"
	assert.Equal(expected, code.ReferenceFieldsPresent(crd, "ko"))
}