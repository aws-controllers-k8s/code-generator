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

package sdk_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	config "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/sdk"
)

func sdkAPIWithCustomShapes(t *testing.T, shapes map[string]map[string]interface{}) (*model.SDKAPI, error) {
	t.Helper()
	helper := sdk.NewHelper(filepath.Clean("../testdata"), config.Config{CustomShapes: shapes})
	return helper.API("s3")
}

func TestCustomShapesTypedMembers(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	api, err := sdkAPIWithCustomShapes(t, map[string]map[string]interface{}{
		// Declared before the shapes it refers to.
		"Outer": {
			"Inner":    "InnerShape",
			"Enabled":  "boolean",
			"Items":    "[]InnerShape",
			"ByName":   "map[string]InnerShape",
			"Tags":     "[]string",
			"Required": "string | required=true",
		},
		"InnerShape": {
			"Name":  "string",
			"Delay": "integer",
			"Ratio": "float",
		},
	})
	require.NoError(err)

	outer := api.API.Shapes["Outer"]
	inner := api.API.Shapes["InnerShape"]
	require.NotNil(outer)
	require.NotNil(inner)
	assert.Equal("structure", outer.Type)

	assert.Same(inner, outer.MemberRefs["Inner"].Shape)
	assert.Equal("boolean", outer.MemberRefs["Enabled"].Shape.Type)
	assert.Equal("list", outer.MemberRefs["Items"].Shape.Type)
	assert.Same(inner, outer.MemberRefs["Items"].Shape.MemberRef.Shape)
	assert.Equal("map", outer.MemberRefs["ByName"].Shape.Type)
	assert.Same(inner, outer.MemberRefs["ByName"].Shape.ValueRef.Shape)
	assert.Equal("string", outer.MemberRefs["Tags"].Shape.MemberRef.Shape.Type)
	assert.Equal("string", outer.MemberRefs["Required"].Shape.Type)

	assert.Equal("string", inner.MemberRefs["Name"].Shape.Type)
	assert.Equal("long", inner.MemberRefs["Delay"].Shape.Type)
	assert.Equal("double", inner.MemberRefs["Ratio"].Shape.Type)

	// A field's type override can name a custom shape.
	ref, err := api.GetShapeRefFromType("Outer")
	require.NoError(err)
	assert.Same(outer, ref.Shape)
	ref, err = api.GetShapeRefFromType("[]InnerShape")
	require.NoError(err)
	assert.Equal("list", ref.Shape.Type)
	assert.Same(inner, ref.Shape.MemberRef.Shape)
}

func TestCustomShapesRejectsUnknownType(t *testing.T) {
	_, err := sdkAPIWithCustomShapes(t, map[string]map[string]interface{}{
		"Outer": {"Inner": "NoSuchShape"},
	})
	assert.Error(t, err)
}

func TestCustomShapesRejectsAnonymousObject(t *testing.T) {
	_, err := sdkAPIWithCustomShapes(t, map[string]map[string]interface{}{
		"Outer": {"Inner": map[string]interface{}{"Name": "string"}},
	})
	assert.Error(t, err)
}
