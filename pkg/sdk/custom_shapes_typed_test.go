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
			"Inner":   "InnerShape",
			"Enabled": "boolean",
			"Items":   "[]InnerShape",
			"ByName":  "map[string]InnerShape",
			"Tags":    "[]string",
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
	// Member refs don't repeat the custom shape's own comment.
	assert.Empty(outer.MemberRefs["Inner"].Documentation)
	assert.Empty(outer.MemberRefs["Items"].Shape.MemberRef.Documentation)

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

// String members must not register global shapes named after the member,
// which would replace a custom or SDK shape of the same name.
func TestCustomShapesStringMemberNameDoesNotReplaceShapes(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	api, err := sdkAPIWithCustomShapes(t, map[string]map[string]interface{}{
		"Outer": {
			"InnerShape": "string",
			"Bucket":     "string",
			"Inner":      "InnerShape",
			"Items":      "[]InnerShape",
		},
		"InnerShape": {"Name": "string"},
	})
	require.NoError(err)

	inner := api.API.Shapes["InnerShape"]
	require.NotNil(inner)
	assert.Equal("structure", inner.Type)
	assert.Equal("structure", api.API.Shapes["Bucket"].Type)
	assert.Equal("string", api.API.Shapes["Outer"].MemberRefs["InnerShape"].Shape.Type)

	ref, err := api.GetShapeRefFromType("[]InnerShape")
	require.NoError(err)
	require.NotNil(ref)
	assert.Same(inner, ref.Shape.MemberRef.Shape)
}

func TestCustomShapesRejectsUnsupportedTypes(t *testing.T) {
	for _, typ := range []string{"NoSuchShape", "long", "int64", "bool", "timestamp", "[]long", "map[string]NoSuchShape"} {
		t.Run(typ, func(t *testing.T) {
			_, err := sdkAPIWithCustomShapes(t, map[string]map[string]interface{}{
				"Outer": {"Inner": typ},
			})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "custom shape Outer, field Inner: unsupported type")
			assert.Contains(t, err.Error(), "want string, boolean, integer, float")
		})
	}
}

func TestGetShapeRefFromTypeStructures(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	api, err := sdkAPIWithCustomShapes(t, nil)
	require.NoError(err)

	// SDK structure shapes resolve.
	ref, err := api.GetShapeRefFromType("Bucket")
	require.NoError(err)
	assert.Same(api.API.Shapes["Bucket"], ref.Shape)
	ref, err = api.GetShapeRefFromType("map[string]Bucket")
	require.NoError(err)
	assert.Equal("map", ref.Shape.Type)
	assert.Same(api.API.Shapes["Bucket"], ref.Shape.ValueRef.Shape)

	// Unknown names and non-structure shapes are errors, not *string.
	require.NotNil(api.API.Shapes["BucketName"])
	require.NotEqual("structure", api.API.Shapes["BucketName"].Type)
	for _, typ := range []string{"NoSuchShape", "[]NoSuchShape", "BucketName"} {
		_, err := api.GetShapeRefFromType(typ)
		assert.ErrorContains(err, "unsupported type override", typ)
	}
}

func TestCustomShapesRejectsCycles(t *testing.T) {
	cases := map[string]map[string]map[string]interface{}{
		"self":       {"A": {"Next": "A"}},
		"self list":  {"A": {"Children": "[]A"}},
		"self map":   {"A": {"Children": "map[string]A"}},
		"two shapes": {"A": {"B": "B"}, "B": {"A": "A"}},
	}
	for name, shapes := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := sdkAPIWithCustomShapes(t, shapes)
			require.Error(t, err)
			assert.Regexp(t, "self references are not allowed|contains a cycle", err.Error())
		})
	}
}

func TestCustomShapesRejectsMarkers(t *testing.T) {
	_, err := sdkAPIWithCustomShapes(t, map[string]map[string]interface{}{
		"Outer": {"Name": "string | required=true"},
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "custom shape Outer, field Name")
	assert.Contains(t, err.Error(), "markers")
}

func TestCustomShapesRejectsAnonymousObject(t *testing.T) {
	_, err := sdkAPIWithCustomShapes(t, map[string]map[string]interface{}{
		"Outer": {"Inner": map[string]interface{}{"Name": "string"}},
	})
	assert.Error(t, err)
}
