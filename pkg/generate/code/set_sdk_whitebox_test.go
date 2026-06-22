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
			// An intEnum is surfaced on the ACK side as a named string enum,
			// while the SDK field is an int32 alias. The write path emits a
			// switch mapping each human-friendly name to its integer value.
			name:            "intEnum scalar (name -> int value)",
			targetFieldName: "EngineVersion",
			targetVarName:   "res",
			targetVarType:   "structure",
			sourceFieldPath: "EngineVersion",
			sourceVarName:   "ko.Spec.EngineVersion",
			isListMember:    false,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type:          "string",
					ShapeName:     "EngineVersion",
					Enum:          []string{"ONE", "TWO"},
					IntEnumValues: map[string]int64{"ONE": 1, "TWO": 2},
				},
				OriginalMemberName: "EngineVersion",
			},
			indentLevel: 1,
			expected: `	switch *ko.Spec.EngineVersion {
	case "ONE":
		res.EngineVersion = svcsdktypes.EngineVersion(1)
	case "TWO":
		res.EngineVersion = svcsdktypes.EngineVersion(2)
	}
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
		{
			// "double" maps to float64 natively — no range check or cast
			// needed. In non-list context, the CRD source is *float64 and
			// the SDK target is *float64, so we strip the dereference.
			name:            "double scalar non-list",
			targetFieldName: "Value",
			targetVarName:   "res",
			targetVarType:   "structure",
			sourceFieldPath: "Value",
			sourceVarName:   "ko.Spec.Value",
			isListMember:    false,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type: "double",
				},
				OriginalMemberName: "Value",
			},
			indentLevel: 1,
			expected:    "\tres.Value = ko.Spec.Value\n",
		},
		{
			// In a list context, the CRD source element is *float64 but
			// the SDK target element is float64 (value type). The generated
			// code must dereference with *.
			name:            "double scalar list member",
			targetFieldName: "",
			targetVarName:   "f0elem",
			targetVarType:   "",
			sourceFieldPath: "Values",
			sourceVarName:   "f0iter",
			isListMember:    true,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type:      "double",
					ShapeName: "Double",
				},
				OriginalMemberName: "Values",
				OrigShapeName:      "Double",
			},
			indentLevel: 1,
			expected:    "\tf0elem = *f0iter\n",
		},
		{
			// "long" maps to int64 natively — no range check or cast needed.
			// In a list context, the CRD source element is *int64 but the
			// SDK target element is int64 (value type). Must dereference.
			name:            "long scalar list member",
			targetFieldName: "",
			targetVarName:   "f0elem",
			targetVarType:   "",
			sourceFieldPath: "Values",
			sourceVarName:   "f0iter",
			isListMember:    true,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type:      "long",
					ShapeName: "Long",
				},
				OriginalMemberName: "Values",
				OrigShapeName:      "Long",
			},
			indentLevel: 1,
			expected:    "\tf0elem = *f0iter\n",
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

func TestSetResourceForScalar(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		name        string
		targetVar   string
		sourceVar   string
		shapeRef    *awssdkmodel.ShapeRef
		indentLevel int
		isList      bool
		isUnion     bool
		expected    string
	}{
		{
			// Plain integer: the SDK field is a *pointer*
			// (HasDefaultValue() == false), so the read path dereferences the
			// source before the int64 cast: `int64(*resp.MaxKeys)`.
			name:        "integer scalar (pointer SDK field)",
			targetVar:   "ko.Spec.MaxKeys",
			sourceVar:   "resp.MaxKeys",
			indentLevel: 1,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type: "integer",
				},
				OriginalMemberName: "MaxKeys",
			},
			expected: "\tmaxKeysCopy := int64(*resp.MaxKeys)\n\tko.Spec.MaxKeys = &maxKeysCopy\n",
		},
		{
			// An intEnum's SDK field is a non-pointer int32 alias while the ACK
			// field is a human-friendly string name. The read path emits a
			// self-contained switch mapping each integer value back to its name,
			// with a default that clears the field to nil. There is no outer
			// "!= 0" guard, since zero can be a legitimate intEnum member.
			name:        "intEnum scalar (int value -> name)",
			targetVar:   "ko.Spec.EngineVersion",
			sourceVar:   "resp.EngineVersion",
			indentLevel: 1,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type:          "string",
					ShapeName:     "EngineVersion",
					Enum:          []string{"ONE", "TWO"},
					IntEnumValues: map[string]int64{"ONE": 1, "TWO": 2},
				},
				OriginalMemberName: "EngineVersion",
			},
			expected: `	switch resp.EngineVersion {
	case 1:
		engineVersionName := "ONE"
		ko.Spec.EngineVersion = &engineVersionName
	case 2:
		engineVersionName := "TWO"
		ko.Spec.EngineVersion = &engineVersionName
	default:
		ko.Spec.EngineVersion = nil
	}
`,
		},
		{
			// Regression: a zero-valued intEnum member must produce a real
			// `case 0:` arm and must NOT be swallowed by a `!= 0` guard. Since
			// the SDK field is a non-pointer int32, zero is indistinguishable
			// from "unset", so the only safe behavior is to map it to its name
			// like any other member.
			name:        "intEnum scalar with zero-valued member",
			targetVar:   "ko.Spec.EngineVersion",
			sourceVar:   "resp.EngineVersion",
			indentLevel: 1,
			shapeRef: &awssdkmodel.ShapeRef{
				Shape: &awssdkmodel.Shape{
					Type:          "string",
					ShapeName:     "EngineVersion",
					Enum:          []string{"ZERO", "ONE"},
					IntEnumValues: map[string]int64{"ZERO": 0, "ONE": 1},
				},
				OriginalMemberName: "EngineVersion",
			},
			expected: `	switch resp.EngineVersion {
	case 0:
		engineVersionName := "ZERO"
		ko.Spec.EngineVersion = &engineVersionName
	case 1:
		engineVersionName := "ONE"
		ko.Spec.EngineVersion = &engineVersionName
	default:
		ko.Spec.EngineVersion = nil
	}
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := setResourceForScalar(
				tc.targetVar,
				tc.sourceVar,
				tc.shapeRef,
				tc.indentLevel,
				tc.isList,
				tc.isUnion,
			)

			assert.Equal(tc.expected, result, "setResourceForScalar() did not return expected result for %s", tc.name)
		})
	}
}
