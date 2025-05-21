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

package code

import (
	"testing"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestSetSDKForScalar(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name            string
		targetFieldName string
		targetVarName   string
		targetVarType   string
		sourceFieldPath string
		sourceVarName   string
		isListMember    bool
		shapeRef        *awssdkmodel.ShapeRef
		indentLevel     int
		expected        string
	}{
		{
			name:            "string scalar",
			targetFieldName: "BucketName",
			targetVarName:   "res",
			targetVarType:   "structure",
			sourceFieldPath: "Name",
			sourceVarName:   "ko.Spec.Name",
			isListMember:    false,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type: "string",
				},
				OriginalMemberName: "BucketName",
			},
			indentLevel: 1,
			expected:    "\tres.BucketName = ko.Spec.Name\n",
		},
		{
			name:            "boolean scalar",
			targetFieldName: "Enabled",
			targetVarName:   "res",
			targetVarType:   "structure",
			sourceFieldPath: "Enabled",
			sourceVarName:   "ko.Spec.Enabled",
			isListMember:    false,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type: "boolean",
				},
				OriginalMemberName: "Enabled",
			},
			indentLevel: 1,
			expected:    "\tres.Enabled = ko.Spec.Enabled\n",
		},
		{
			name:            "integer scalar",
			targetFieldName: "MaxKeys",
			targetVarName:   "res",
			targetVarType:   "structure",
			sourceFieldPath: "MaxKeys",
			sourceVarName:   "ko.Spec.MaxKeys",
			isListMember:    false,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type: "integer",
				},
				OriginalMemberName: "MaxKeys",
			},
			indentLevel: 1,
			expected: `	maxKeysCopy0 := *ko.Spec.MaxKeys
	if maxKeysCopy0 > math.MaxInt32 || maxKeysCopy0 < math.MinInt32 {
		return nil, fmt.Errorf("error: field MaxKeys is of type int32")
	}
	maxKeysCopy := int32(maxKeysCopy0)
	res.MaxKeys = &maxKeysCopy
`,
		},
		{
			name:            "float scalar",
			targetFieldName: "Temperature",
			targetVarName:   "res",
			targetVarType:   "structure",
			sourceFieldPath: "Temperature",
			sourceVarName:   "ko.Spec.Temperature",
			isListMember:    false,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type: "float",
				},
				OriginalMemberName: "Temperature",
			},
			indentLevel: 1,
			expected: `	temperatureCopy0 := *ko.Spec.Temperature
	if temperatureCopy0 > math.MaxFloat32 || temperatureCopy0 < -math.MaxFloat32 || (temperatureCopy0 < math.SmallestNonzeroFloat32 && !(temperatureCopy0 <= 0)) || (temperatureCopy0 > -math.SmallestNonzeroFloat32 && !(temperatureCopy0 >= 0)) {
		return nil, fmt.Errorf("error: field Temperature is of type float32")
	}
	temperatureCopy := float32(temperatureCopy0)
	res.Temperature = &temperatureCopy
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := setSDKForScalar(
				tc.targetFieldName,
				tc.targetVarName,
				tc.targetVarType,
				tc.sourceFieldPath,
				tc.sourceVarName,
				tc.isListMember,
				tc.shapeRef,
				tc.indentLevel,
			)

			assert.Equal(tc.expected, result, "setSDKForScalar() did not return expected result for %s", tc.name)
		})
	}
}
