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

	"github.com/aws-controllers-k8s/code-generator/pkg/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestCheckRequiredFields_Attributes_ARNField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sns")

	crd := testutil.GetCRDByName(t, g, "Topic")
	require.NotNil(crd)

	// The Go code for checking the GetTopicAttributes Input shape's required
	// fields needs to return false when any required field is missing in the
	// corresponding Spec or Status. The GetTopicAttributesInput shape has a
	// required TopicArn field which corresponds to the resource's ARN which is
	// stored in ACKMetadata.ARN, so the primary resource ARN field if
	// condition is a bit special.
	expReqFieldsInShape := `
	return (ko.Status.ACKResourceMetadata == nil || ko.Status.ACKResourceMetadata.ARN == nil)
`
	gotCode, err := code.CheckRequiredFieldsMissingFromShape(
		crd, model.OpTypeGetAttributes, "ko", 1,
	)
	require.NoError(err)
	assert.Equal(
		strings.TrimSpace(expReqFieldsInShape),
		strings.TrimSpace(gotCode),
	)
}

func TestCheckRequiredFields_Attributes_StatusField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sqs")

	crd := testutil.GetCRDByName(t, g, "Queue")
	require.NotNil(crd)

	expRequiredFieldsCode := `
	return r.ko.Status.QueueURL == nil
`
	gotCode, err := code.CheckRequiredFieldsMissingFromShape(
		crd, model.OpTypeGetAttributes, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal(
		strings.TrimSpace(expRequiredFieldsCode),
		strings.TrimSpace(gotCode),
	)
}

func TestCheckRequiredFields_Attributes_StatusAndSpecField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "Route")
	require.NotNil(crd)

	expRequiredFieldsCode := `
	return r.ko.Spec.APIID == nil || r.ko.Status.RouteID == nil
`
	gotCode, err := code.CheckRequiredFieldsMissingFromShape(
		crd, model.OpTypeGet, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal(
		strings.TrimSpace(expRequiredFieldsCode),
		strings.TrimSpace(gotCode),
	)
}

func TestCheckRequiredFields_RenamedSpecField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "eks")

	crd := testutil.GetCRDByName(t, g, "FargateProfile")
	require.NotNil(crd)

	expRequiredFieldsCode := `
	return r.ko.Spec.ClusterName == nil || r.ko.Spec.Name == nil
`
	gotCode, err := code.CheckRequiredFieldsMissingFromShape(
		crd, model.OpTypeGet, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal(
		strings.TrimSpace(expRequiredFieldsCode),
		strings.TrimSpace(gotCode),
	)
}

func TestCheckRequiredFields_StatusField_ReadMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "Vpc")
	require.NotNil(crd)

	expRequiredFieldsCode := `
	return r.ko.Status.VPCID == nil
`
	gotCode, err := code.CheckRequiredFieldsMissingFromShape(
		crd, model.OpTypeList, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal(
		strings.TrimSpace(expRequiredFieldsCode),
		strings.TrimSpace(gotCode),
	)
}

func TestCheckRequiredFields_StatusField_ReadMany_EgressOnlyIGW(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ec2", &testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-egress.yaml"})

	crd := testutil.GetCRDByName(t, g, "EgressOnlyInternetGateway")
	require.NotNil(crd)

	expRequiredFieldsCode := `
	return r.ko.Status.ID == nil
`
	gotCode, err := code.CheckRequiredFieldsMissingFromShape(
		crd, model.OpTypeList, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal(
		strings.TrimSpace(expRequiredFieldsCode),
		strings.TrimSpace(gotCode),
	)
}

func TestCheckRequiredFields_MutuallyExclusiveIdentifiers_ReadOne(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "opensearchserverless", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-mutually-exclusive-identifiers.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "SecurityPolicy")
	require.NotNil(crd)
	require.True(crd.HasMutuallyExclusiveIdentifiers())

	// GetSecurityPolicy (ReadOne) marks both `name` and `type` required.
	// Declaring them mutually exclusive drops each from the per-field `||` list
	// and collapses them into a single grouped condition that is true only when
	// neither identifier is set.
	expRequiredFieldsCode := `
	return (r.ko.Spec.Name == nil && r.ko.Spec.Type == nil)
`
	gotCode, err := code.CheckRequiredFieldsMissingFromShape(
		crd, model.OpTypeGet, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal(
		strings.TrimSpace(expRequiredFieldsCode),
		strings.TrimSpace(gotCode),
	)
}

func TestCheckRequiredFields_MutuallyExclusiveIdentifiers_ReadMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "cloudwatch-logs", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "ResourcePolicy")
	require.NotNil(crd)
	require.True(crd.HasMutuallyExclusiveIdentifiers())

	// DescribeResourcePolicies has no required input members, so without
	// mutually_exclusive_identifiers this would generate `return false` and let
	// sdkFind match an arbitrary policy. Instead, the read input is treated as
	// incomplete unless at least one of the declared identifiers is set, which
	// is what the controller previously had to supply via a custom method.
	expRequiredFieldsCode := `
	return (r.ko.Spec.PolicyName == nil && r.ko.Spec.ResourceARN == nil)
`
	gotCode, err := code.CheckRequiredFieldsMissingFromShape(
		crd, model.OpTypeList, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal(
		strings.TrimSpace(expRequiredFieldsCode),
		strings.TrimSpace(gotCode),
	)
}

func TestCheckNilFieldPath(t *testing.T) {
	// Empty FieldPath
	field := model.Field{Path: ""}
	assert.Equal(t, "", code.CheckNilFieldPath(&field, "ko.Spec"))
	// FieldPath only has fieldName
	field.Path = "JWTConfiguration"
	assert.Equal(t, "", code.CheckNilFieldPath(&field, "ko.Spec"))
	// Nested FieldPath
	field.Path = "JWTConfiguration.Issuer"
	assert.Equal(t,
		"ko.Spec.JWTConfiguration == nil",
		code.CheckNilFieldPath(&field, "ko.Spec"))
	// Multi Level Nested FieldPath
	field.Path = "JWTConfiguration.Issuer.FieldName"
	assert.Equal(t,
		"ko.Spec.JWTConfiguration == nil || ko.Spec.JWTConfiguration.Issuer == nil",
		code.CheckNilFieldPath(&field, "ko.Spec"))
}

func TestCheckNilReferencesPath(t *testing.T) {
	field := model.Field{}
	// Empty ReferencesPath
	referenceFieldConfig := config.ReferencesConfig{Path: ""}
	fieldConfig := config.FieldConfig{References: &referenceFieldConfig}
	field.FieldConfig = &fieldConfig
	assert.Equal(t, "", code.CheckNilReferencesPath(&field, "obj"))
	// Non nested ReferencesPath
	referenceFieldConfig.Path = "Status"
	assert.Equal(t, "", code.CheckNilReferencesPath(&field, "obj"))
	// Nested ReferencesPath
	referenceFieldConfig.Path = "Status.ACKResourceMetadata"
	assert.Equal(t,
		"obj.Status.ACKResourceMetadata == nil",
		code.CheckNilReferencesPath(&field, "obj"))
	// Multi Level Nested ReferencesPath
	referenceFieldConfig.Path = "Status.ACKResourceMetadata.ARN"
	assert.Equal(t,
		"obj.Status.ACKResourceMetadata == nil || obj.Status.ACKResourceMetadata.ARN == nil",
		code.CheckNilReferencesPath(&field, "obj"))
}

// TestCheckRequiredFields_OpenSearchServerless_LifecyclePolicy_ReadMany_InputWrapper
// exercises input_wrapper_field_path on a ReadMany operation. Promoting
// BatchGetLifecyclePolicy to ReadMany makes its input the source for the
// required-fields check. That input nests the identifier fields inside the
// "identifiers" list-of-structure wrapper; with input_wrapper_field_path
// configured, the check sees the unwrapped identifier members and emits a real
// nil-check on EVERY required identifier component of the composite key (Name
// AND Type), rather than the permissive "return false" fallback it would
// otherwise produce or a partial check that guards only the primary key.
func TestCheckRequiredFields_OpenSearchServerless_LifecyclePolicy_ReadMany_InputWrapper(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "opensearchserverless",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-input-wrapper-field-path.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "LifecyclePolicy")
	require.NotNil(crd)

		expReqFieldsInShape := `
	return r.ko.Spec.Name == nil || r.ko.Spec.Type == nil
`
	gotCode, err := code.CheckRequiredFieldsMissingFromShape(
		crd, model.OpTypeList, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal(
		strings.TrimSpace(expReqFieldsInShape),
		strings.TrimSpace(gotCode),
	)
}

// TestCheckRequiredFields_OpenSearchServerless_LifecyclePolicy_ReadMany_NoInputWrapper
// is the regression baseline: without input_wrapper_field_path the nested
// identifiers member matches no scalar identifier, so the required-fields check
// falls back to the permissive "return false".
func TestCheckRequiredFields_OpenSearchServerless_LifecyclePolicy_ReadMany_NoInputWrapper(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "opensearchserverless",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-batchget-readmany-no-wrapper.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "LifecyclePolicy")
	require.NotNil(crd)

	gotCode, err := code.CheckRequiredFieldsMissingFromShape(
		crd, model.OpTypeList, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal("return false", strings.TrimSpace(gotCode))
}
