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

package model

import (
	"testing"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddMemberShapRef_Structure(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	structShape := &awssdkmodel.Shape{
		Type:       "structure",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{},
	}
	shapeRef := &awssdkmodel.ShapeRef{
		Shape: structShape,
	}

	memberShapeRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "string"},
	}

	err := addMemberShapRef(shapeRef, memberShapeRef, "NewField")

	require.NoError(err)
	require.Contains(structShape.MemberRefs, "NewField")
	assert.Equal(memberShapeRef, structShape.MemberRefs["NewField"])
}

func TestAddMemberShapRef_List(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	innerStructShape := &awssdkmodel.Shape{
		Type:       "structure",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{},
	}
	listShape := &awssdkmodel.Shape{
		Type: "list",
		MemberRef: awssdkmodel.ShapeRef{
			Shape: innerStructShape,
		},
	}
	shapeRef := &awssdkmodel.ShapeRef{
		Shape: listShape,
	}

	memberShapeRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "integer"},
	}

	err := addMemberShapRef(shapeRef, memberShapeRef, "Count")

	require.NoError(err)
	require.Contains(innerStructShape.MemberRefs, "Count")
	assert.Equal(memberShapeRef, innerStructShape.MemberRefs["Count"])
}

func TestAddMemberShapRef_Map(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	valueStructShape := &awssdkmodel.Shape{
		Type:       "structure",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{},
	}
	mapShape := &awssdkmodel.Shape{
		Type: "map",
		ValueRef: awssdkmodel.ShapeRef{
			Shape: valueStructShape,
		},
	}
	shapeRef := &awssdkmodel.ShapeRef{
		Shape: mapShape,
	}

	memberShapeRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "boolean"},
	}

	err := addMemberShapRef(shapeRef, memberShapeRef, "Enabled")

	require.NoError(err)
	require.Contains(valueStructShape.MemberRefs, "Enabled")
	assert.Equal(memberShapeRef, valueStructShape.MemberRefs["Enabled"])
}

func TestAddMemberShapRef_UnsupportedType(t *testing.T) {
	assert := assert.New(t)

	scalarShape := &awssdkmodel.Shape{
		Type:       "string",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{},
	}
	shapeRef := &awssdkmodel.ShapeRef{
		Shape: scalarShape,
	}

	memberShapeRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "string"},
	}

	err := addMemberShapRef(shapeRef, memberShapeRef, "ShouldNotExist")

	assert.Error(err)
	assert.Contains(err.Error(), "unsupported shape type")
	assert.Contains(err.Error(), "string")
	assert.Contains(err.Error(), "ShouldNotExist")
	assert.NotContains(scalarShape.MemberRefs, "ShouldNotExist")
}

func TestAddMemberShapRef_Structure_PreservesExistingMembers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	existingMemberRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "string"},
	}
	structShape := &awssdkmodel.Shape{
		Type: "structure",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{
			"ExistingField": existingMemberRef,
		},
	}
	shapeRef := &awssdkmodel.ShapeRef{
		Shape: structShape,
	}

	newMemberShapeRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "integer"},
	}

	err := addMemberShapRef(shapeRef, newMemberShapeRef, "NewField")

	require.NoError(err)
	require.Len(structShape.MemberRefs, 2)
	assert.Equal(existingMemberRef, structShape.MemberRefs["ExistingField"])
	assert.Equal(newMemberShapeRef, structShape.MemberRefs["NewField"])
}

func TestAddMemberShapRef_List_PreservesExistingMembers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	existingMemberRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "string"},
	}
	innerStructShape := &awssdkmodel.Shape{
		Type: "structure",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{
			"ExistingField": existingMemberRef,
		},
	}
	listShape := &awssdkmodel.Shape{
		Type: "list",
		MemberRef: awssdkmodel.ShapeRef{
			Shape: innerStructShape,
		},
	}
	shapeRef := &awssdkmodel.ShapeRef{
		Shape: listShape,
	}

	newMemberShapeRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "boolean"},
	}

	err := addMemberShapRef(shapeRef, newMemberShapeRef, "Active")

	require.NoError(err)
	require.Len(innerStructShape.MemberRefs, 2)
	assert.Equal(existingMemberRef, innerStructShape.MemberRefs["ExistingField"])
	assert.Equal(newMemberShapeRef, innerStructShape.MemberRefs["Active"])
}

func TestAddMemberShapRef_Map_PreservesExistingMembers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	existingMemberRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "string"},
	}
	valueStructShape := &awssdkmodel.Shape{
		Type: "structure",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{
			"ExistingField": existingMemberRef,
		},
	}
	mapShape := &awssdkmodel.Shape{
		Type: "map",
		ValueRef: awssdkmodel.ShapeRef{
			Shape: valueStructShape,
		},
	}
	shapeRef := &awssdkmodel.ShapeRef{
		Shape: mapShape,
	}

	newMemberShapeRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "integer"},
	}

	err := addMemberShapRef(shapeRef, newMemberShapeRef, "Priority")

	require.NoError(err)
	require.Len(valueStructShape.MemberRefs, 2)
	assert.Equal(existingMemberRef, valueStructShape.MemberRefs["ExistingField"])
	assert.Equal(newMemberShapeRef, valueStructShape.MemberRefs["Priority"])
}

func TestAddMemberShapRef_DuplicateField_Structure(t *testing.T) {
	assert := assert.New(t)

	oldMemberRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "string"},
	}
	structShape := &awssdkmodel.Shape{
		Type: "structure",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{
			"Field": oldMemberRef,
		},
	}
	shapeRef := &awssdkmodel.ShapeRef{
		Shape: structShape,
	}

	newMemberRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "integer"},
	}

	err := addMemberShapRef(shapeRef, newMemberRef, "Field")

	assert.Error(err)
	assert.Contains(err.Error(), "Field")
	assert.Contains(err.Error(), "already exists")
	// Original member should be unchanged
	assert.Equal(oldMemberRef, structShape.MemberRefs["Field"])
}

func TestAddMemberShapRef_DuplicateField_List(t *testing.T) {
	assert := assert.New(t)

	oldMemberRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "string"},
	}
	innerStructShape := &awssdkmodel.Shape{
		Type: "structure",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{
			"Field": oldMemberRef,
		},
	}
	listShape := &awssdkmodel.Shape{
		Type: "list",
		MemberRef: awssdkmodel.ShapeRef{
			Shape: innerStructShape,
		},
	}
	shapeRef := &awssdkmodel.ShapeRef{
		Shape: listShape,
	}

	newMemberRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "integer"},
	}

	err := addMemberShapRef(shapeRef, newMemberRef, "Field")

	assert.Error(err)
	assert.Contains(err.Error(), "Field")
	assert.Contains(err.Error(), "already exists")
	assert.Equal(oldMemberRef, innerStructShape.MemberRefs["Field"])
}

func TestAddMemberShapRef_DuplicateField_Map(t *testing.T) {
	assert := assert.New(t)

	oldMemberRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "string"},
	}
	valueStructShape := &awssdkmodel.Shape{
		Type: "structure",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{
			"Field": oldMemberRef,
		},
	}
	mapShape := &awssdkmodel.Shape{
		Type: "map",
		ValueRef: awssdkmodel.ShapeRef{
			Shape: valueStructShape,
		},
	}
	shapeRef := &awssdkmodel.ShapeRef{
		Shape: mapShape,
	}

	newMemberRef := &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "integer"},
	}

	err := addMemberShapRef(shapeRef, newMemberRef, "Field")

	assert.Error(err)
	assert.Contains(err.Error(), "Field")
	assert.Contains(err.Error(), "already exists")
	assert.Equal(oldMemberRef, valueStructShape.MemberRefs["Field"])
}
