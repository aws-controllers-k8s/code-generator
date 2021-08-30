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
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

func TestFindIdentifiersInShape_EC2_VPC_ReadMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "Vpc")
	require.NotNil(crd)

	expIdentifier := "VpcIds"
	actualIdentifiers := code.FindIdentifiersInShape(crd,
		crd.Ops.ReadMany.InputRef.Shape)
	assert.Len(actualIdentifiers, 1)
	assert.Equal(
		strings.TrimSpace(expIdentifier),
		strings.TrimSpace(actualIdentifiers[0]),
	)
}

func TestFindIdentifiersInCRD_S3_Bucket_ReadMany_Empty(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "s3")

	crd := testutil.GetCRDByName(t, g, "Bucket")
	require.NotNil(crd)

	actualIdentifiers := code.FindIdentifiersInShape(crd,
		crd.Ops.ReadMany.InputRef.Shape)
	assert.Len(actualIdentifiers, 0)
}

func TestGetIdentifiers_EC2_VPC_StatusField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "Vpc")
	require.NotNil(crd)

	expIdentifier := "VpcId"
	actualIdentifiers := crd.GetIdentifiers()
	assert.Len(actualIdentifiers, 1)
	assert.Equal(
		strings.TrimSpace(expIdentifier),
		strings.TrimSpace(actualIdentifiers[0]),
	)
}

func TestGetIdentifiers_S3_Bucket_SpecField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "s3")

	crd := testutil.GetCRDByName(t, g, "Bucket")
	require.NotNil(crd)

	expIdentifier := "Name"
	actualIdentifiers := crd.GetIdentifiers()
	assert.Len(actualIdentifiers, 1)
	assert.Equal(
		strings.TrimSpace(expIdentifier),
		strings.TrimSpace(actualIdentifiers[0]),
	)
}

func TestGetIdentifiers_APIGatewayV2_API_Multiple(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")
	crd := testutil.GetCRDByName(t, g, "Api")
	require.NotNil(crd)

	expIdentifiers := []string{"ApiId", "Name"}
	actualIdentifiers := crd.GetIdentifiers()
	assert.Len(actualIdentifiers, 2)
	assert.True(ackcompare.SliceStringEqual(expIdentifiers, actualIdentifiers))
}