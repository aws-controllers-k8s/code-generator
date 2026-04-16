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

package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestBedrockAgentCoreControl_GatewayTarget_CustomNestedField_ListOfStruct(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "bedrock-agentcore-control", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-custom-nested-fields.yaml",
	})

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("GatewayTarget", crds)
	require.NotNil(crd)

	// Walk the nested path to the InlinePayload field:
	// TargetConfiguration -> Mcp -> Lambda -> ToolSchema -> InlinePayload
	// InlinePayload is a list of ToolDefinition structs.
	targetConfigField := crd.Fields["TargetConfiguration"]
	require.NotNil(targetConfigField, "TargetConfiguration field should exist")

	mcpField := targetConfigField.MemberFields["Mcp"]
	require.NotNil(mcpField, "Mcp member field should exist")

	lambdaField := mcpField.MemberFields["Lambda"]
	require.NotNil(lambdaField, "Lambda member field should exist")

	toolSchemaField := lambdaField.MemberFields["ToolSchema"]
	require.NotNil(toolSchemaField, "ToolSchema member field should exist")

	inlinePayloadField := toolSchemaField.MemberFields["InlinePayload"]
	require.NotNil(inlinePayloadField, "InlinePayload member field should exist")
	require.NotNil(inlinePayloadField.ShapeRef)
	require.NotNil(inlinePayloadField.ShapeRef.Shape)
	assert.Equal("list", inlinePayloadField.ShapeRef.Shape.Type)

	// The list element shape (ToolDefinition) should be a structure
	toolDefShape := inlinePayloadField.ShapeRef.Shape.MemberRef.Shape
	require.NotNil(toolDefShape)
	assert.Equal("structure", toolDefShape.Type)

	// Verify the custom nested fields "InputSchema" and "OutputSchema" were
	// added to the ToolDefinition struct shape as string types
	require.Contains(toolDefShape.MemberRefs, "InputSchema")
	inputSchemaRef := toolDefShape.MemberRefs["InputSchema"]
	require.NotNil(inputSchemaRef.Shape)
	assert.Equal("string", inputSchemaRef.Shape.Type)

	require.Contains(toolDefShape.MemberRefs, "OutputSchema")
	outputSchemaRef := toolDefShape.MemberRefs["OutputSchema"]
	require.NotNil(outputSchemaRef.Shape)
	assert.Equal("string", outputSchemaRef.Shape.Type)

	// Verify the original ToolDefinition members are still present
	assert.Contains(toolDefShape.MemberRefs, "Name")
	assert.Contains(toolDefShape.MemberRefs, "Description")

	// Verify the custom nested fields appear as MemberFields on the
	// InlinePayload field
	require.NotNil(inlinePayloadField.MemberFields)
	inputSchemaMember := inlinePayloadField.MemberFields["InputSchema"]
	require.NotNil(inputSchemaMember)

	outputSchemaMember := inlinePayloadField.MemberFields["OutputSchema"]
	require.NotNil(outputSchemaMember)
}
