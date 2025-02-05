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

package multiversion_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/model/multiversion"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestAreEqualShapes_APIGatewayV2_DomainName(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	model := testutil.NewModelForServiceWithOptions(t, "apigatewayv2", &testutil.TestingModelOptions{})
	crds, err := model.GetCRDs()
	require.Nil(err)
	require.Len(crds, 12)
	domainNameCRD := crds[4]

	for _, fieldNameX := range domainNameCRD.SpecFieldNames() {
		for _, fieldNameY := range domainNameCRD.SpecFieldNames() {
			testName := fmt.Sprintf("comparing %s with %s", fieldNameX, fieldNameY)
			t.Run(testName, func(t *testing.T) {

				fieldX := domainNameCRD.SpecFields[fieldNameX]
				fieldY := domainNameCRD.SpecFields[fieldNameY]
				equal, _ := multiversion.AreEqualShapes(fieldX.ShapeRef.Shape, fieldY.ShapeRef.Shape, false)
				if fieldNameY == fieldNameX {
					assert.True(equal)
				} else {
					assert.False(equal)
				}
			})
		}
	}
}

func TestComputeCRDDeltas_APIGatewayV2_DomainName(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	model := testutil.NewModelForServiceWithOptions(t, "apigatewayv2", &testutil.TestingModelOptions{})
	crds, err := model.GetCRDs()
	require.Nil(err)
	require.Len(crds, 12)
	domainNameCRD := crds[4]

	deltas, err := multiversion.ComputeCRDFieldDeltas(domainNameCRD, domainNameCRD)
	require.Nil(err)
	assert.Len(deltas.SpecDeltas, len(domainNameCRD.SpecFields))
	assert.Len(deltas.StatusDeltas, len(domainNameCRD.StatusFields))

	for _, delta := range deltas.SpecDeltas {
		assert.Equal(delta.ChangeType, multiversion.FieldChangeTypeNone)
	}
	for _, delta := range deltas.StatusDeltas {
		assert.Equal(delta.ChangeType, multiversion.FieldChangeTypeNone)
	}
}
