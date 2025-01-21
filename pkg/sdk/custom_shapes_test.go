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

var (
	s3 *model.SDKAPI
)

func customListConfig(fieldName string, shapeName string) config.Config {
	return config.Config{
		Resources: map[string]config.ResourceConfig{
			"Bucket": {
				Fields: map[string]*config.FieldConfig{
					fieldName: {
						CustomField: &config.CustomFieldConfig{
							ListOf: shapeName,
						},
					},
				},
			},
		},
	}
}

func customMapConfig(fieldName string, shapeName string) config.Config {
	return config.Config{
		Resources: map[string]config.ResourceConfig{
			"Bucket": {
				Fields: map[string]*config.FieldConfig{
					fieldName: {
						CustomField: &config.CustomFieldConfig{
							MapOf: shapeName,
						},
					},
				},
			},
		},
	}
}

func s3SDKAPI(t *testing.T, cfg config.Config) *model.SDKAPI {
	if s3 != nil {
		return s3
	}
	path := filepath.Clean("../testdata")
	sdkHelper := sdk.NewHelper(path, cfg)
	s3, err := sdkHelper.API("s3")
	if err != nil {
		t.Fatal(err)
	}
	return s3
}

func TestCustomListField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fieldName := "MyCustomListField"
	shapeName := "Tag"

	api := s3SDKAPI(t, customListConfig(fieldName, shapeName))

	// Assert custom shape was registered with SDKAPI
	shapeRef := api.GetCustomShapeRef(shapeName)
	assert.NotNil(shapeRef)

	memberShape, exists := api.API.Shapes[shapeName]
	require.True(exists)

	// Assert custom shape was well formed
	assert.Equal(shapeRef.Shape.MemberRef.Shape, memberShape)
	assert.Nil(shapeRef.Shape.KeyRef.Shape)
	assert.Nil(shapeRef.Shape.ValueRef.Shape)
	assert.Empty(shapeRef.Shape.MemberRefs)

	// Assert custom shape was registered into API shapes
	_, exists = api.API.Shapes[shapeRef.ShapeName]
	assert.True(exists)
}

func TestCustomMapField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	fieldName := "MyCustomMapField"
	shapeName := "Tag"

	api := s3SDKAPI(t, customMapConfig(fieldName, shapeName))

	// Assert custom shape was registered with SDKAPI
	shapeRef := api.GetCustomShapeRef(shapeName)
	assert.NotNil(shapeRef)

	memberShape, exists := api.API.Shapes[shapeName]
	require.True(exists)

	// Assert custom shape was well formed
	assert.Equal(shapeRef.Shape.ValueRef.Shape, memberShape)
	assert.Nil(shapeRef.Shape.MemberRef.Shape)
	assert.Empty(shapeRef.Shape.MemberRefs)

	// Assert custom key shape was created
	keyRef := shapeRef.Shape.KeyRef
	assert.NotNil(keyRef.Shape)
	assert.Equal(keyRef.Shape.Type, "string")

	// Assert custom shape was registered into API shapes
	_, exists = api.API.Shapes[shapeRef.ShapeName]
	assert.True(exists)
}
