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

package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestTypedCustomShapes(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "dynamodb", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-typed-custom-shapes.yaml",
	})

	crds, err := g.GetCRDs()
	require.NoError(err)
	crd := getCRDByName("Table", crds)
	require.NotNil(crd)

	// Fields typed as a custom shape, and as a map of one.
	cfg := crd.SpecFields["TypedConfig"]
	require.NotNil(cfg)
	assert.Equal("*TypedConfig", cfg.GoType)
	byName := crd.SpecFields["TypedConfigByName"]
	require.NotNil(byName)
	assert.Equal("map[string]*TypedInner", byName.GoType)

	// Member fields reach into nested custom shapes.
	require.Contains(cfg.MemberFields, "Items")
	items := cfg.MemberFields["Items"]
	assert.Equal("[]*TypedInner", items.GoType)
	require.Contains(items.MemberFields, "TableName")
	assert.Equal("*string", items.MemberFields["TableName"].GoType)

	// A reference on a nested custom shape path.
	tableName := crd.Fields["TypedConfig.Items.TableName"]
	require.NotNil(tableName)
	assert.True(tableName.HasReference())
	assert.True(crd.HasReferenceFields())

	// Both custom shapes get type definitions with typed members.
	outer := testutil.GetTypeDefByName(t, g, "TypedConfig")
	require.NotNil(outer)
	assert.Equal("*bool", outer.GetAttribute("Enabled").GoType)
	assert.Equal("*int64", outer.GetAttribute("Limit").GoType)
	assert.Equal("*float64", outer.GetAttribute("Ratio").GoType)
	assert.Equal("*TypedInner", outer.GetAttribute("Inner").GoType)
	assert.Equal("[]*TypedInner", outer.GetAttribute("Items").GoType)

	inner := testutil.GetTypeDefByName(t, g, "TypedInner")
	require.NotNil(inner)
	assert.Equal("*string", inner.GetAttribute("TableName").GoType)
	assert.Equal("map[string]*string", inner.GetAttribute("Labels").GoType)
	assert.NotNil(inner.GetAttribute("TableRef"))
}
