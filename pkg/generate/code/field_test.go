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
	"github.com/aws-controllers-k8s/code-generator/pkg/model"

	"github.com/stretchr/testify/assert"
)

func Test_NilFieldPathCheck(t *testing.T) {
	// Empty FieldPath
	field := model.Field{Path: ""}
	assert.Equal(t, "", code.NilFieldPathCheck(&field, "ko.Spec"))
	// FieldPath only has fieldName
	field.Path = "JWTConfiguration"
	assert.Equal(t, "", code.NilFieldPathCheck(&field, "ko.Spec"))
	// Nested FieldPath
	field.Path = "JWTConfiguration.Issuer"
	assert.Equal(t,
		"ko.Spec.JWTConfiguration == nil",
		code.NilFieldPathCheck(&field, "ko.Spec"))
	// Multi Level Nested FieldPath
	field.Path = "JWTConfiguration.Issuer.FieldName"
	assert.Equal(t,
		"ko.Spec.JWTConfiguration == nil || ko.Spec.JWTConfiguration.Issuer == nil",
		code.NilFieldPathCheck(&field, "ko.Spec"))
}
