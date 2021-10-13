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
	"github.com/aws-controllers-k8s/code-generator/pkg/operations"
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
	assert.Equal(
		strings.TrimSpace(expReqFieldsInShape),
		strings.TrimSpace(
			code.CheckRequiredFieldsMissingFromShape(
				crd, operations.OpTypeGetAttributes, "ko", 1,
			),
		),
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
	gotCode := code.CheckRequiredFieldsMissingFromShape(
		crd, operations.OpTypeGetAttributes, "r.ko", 1,
	)
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
	gotCode := code.CheckRequiredFieldsMissingFromShape(
		crd, operations.OpTypeGet, "r.ko", 1,
	)
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
	gotCode := code.CheckRequiredFieldsMissingFromShape(
		crd, operations.OpTypeGet, "r.ko", 1,
	)
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
	gotCode := code.CheckRequiredFieldsMissingFromShape(
		crd, operations.OpTypeList, "r.ko", 1,
	)
	assert.Equal(
		strings.TrimSpace(expRequiredFieldsCode),
		strings.TrimSpace(gotCode),
	)
}
