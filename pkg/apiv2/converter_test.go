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

package apiv2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
)

// TestInjectCustomPrimitiveShapes_Unit verifies that the Smithy prelude Unit
// type is injected as an empty structure shape so that union/structure members
// targeting smithy.api#Unit resolve instead of panicking.
func TestInjectCustomPrimitiveShapes_Unit(t *testing.T) {
	shapes := injectCustomPrimitiveShapes()

	unit, ok := shapes["Unit"]
	require.True(t, ok, "expected an injected \"Unit\" shape")
	assert.Equal(t, "structure", unit.Type)
	assert.Empty(t, unit.MemberRefs, "Unit must be an empty structure")
	assert.NotNil(t, unit.MemberRefs, "Unit.MemberRefs must be initialized (non-nil)")
}

// TestBuildAPI_UnionMemberTargetingUnit reproduces the panic seen when a model
// (e.g. lambda-microvms' PortSpecification.allPorts) has a union member that
// targets smithy.api#Unit. Before the Unit shape was injected, reference
// resolution during Setup() panicked with "unable resolve reference, Unit".
func TestBuildAPI_UnionMemberTargetingUnit(t *testing.T) {
	const alias = "examplesvc"
	prefix := func(name string) string {
		return "com.amazonaws." + alias + "#" + name
	}

	shapes := map[string]Shape{
		prefix("ExampleService"): {
			Type: "service",
			Traits: map[string]interface{}{
				"aws.api#service": map[string]interface{}{
					"sdkId": "Example",
				},
				"smithy.api#documentation": "Example service.",
			},
		},
		prefix("DoThing"): {
			Type:      "operation",
			InputRef:  ShapeRef{ShapeName: prefix("DoThingInput")},
			OutputRef: ShapeRef{ShapeName: prefix("DoThingOutput")},
		},
		prefix("DoThingInput"): {
			Type: "structure",
			MemberRefs: map[string]*ShapeRef{
				"PortSpec": {ShapeName: prefix("PortSpecification")},
			},
		},
		prefix("DoThingOutput"): {
			Type:       "structure",
			MemberRefs: map[string]*ShapeRef{},
		},
		// The union under test: one variant carries data, the other is a
		// valueless variant targeting the Smithy prelude Unit type.
		prefix("PortSpecification"): {
			Type: "union",
			MemberRefs: map[string]*ShapeRef{
				"allPorts": {ShapeName: "smithy.api#Unit"},
				"port":     {ShapeName: prefix("Port")},
			},
		},
		prefix("Port"): {Type: "integer"},
	}

	var (
		api *awssdkmodel.API
		err error
	)
	require.NotPanics(t, func() {
		api, _, err = buildAPI(shapes)
	}, "buildAPI must not panic on a union member targeting smithy.api#Unit")
	require.NoError(t, err)

	// Setup() runs reference resolution, which is where the original panic
	// occurred. It must complete cleanly now that Unit is a known shape.
	require.NotPanics(t, func() {
		err = api.Setup()
	}, "API.Setup() must not panic resolving the Unit reference")
	require.NoError(t, err)

	// The allPorts variant must resolve to the injected empty Unit structure.
	// Setup() exports (capitalizes) member names, so "allPorts" becomes
	// "AllPorts".
	portSpec, ok := api.Shapes["PortSpecification"]
	require.True(t, ok, "PortSpecification shape should exist")
	allPorts, ok := portSpec.MemberRefs["AllPorts"]
	require.True(t, ok, "AllPorts member should be retained")
	require.NotNil(t, allPorts.Shape, "AllPorts member ref should be resolved")
	assert.Equal(t, "structure", allPorts.Shape.Type)
	assert.Empty(t, allPorts.Shape.MemberRefs, "resolved Unit shape should be an empty structure")
}

// TestCreateApiShape_IntEnum verifies a Smithy `intEnum` shape is preserved
// as-is (Type stays "intEnum", no DefaultValue injected). Downstream code
// (IsNonPointerInSDK, setSDKForScalar, etc.) recognizes "intEnum" directly.
func TestCreateApiShape_IntEnum(t *testing.T) {
	assert := assert.New(t)

	intEnumShape := Shape{
		Type: "intEnum",
		MemberRefs: map[string]*ShapeRef{
			"ONE": {
				ShapeName: "smithy.api#Unit",
				Traits: map[string]interface{}{
					"smithy.api#enumValue": float64(1),
				},
			},
		},
	}

	apiShape, err := createApiShape(intEnumShape)
	assert.NoError(err)
	// intEnum type is preserved — no mutation to "integer".
	assert.Equal("intEnum", apiShape.Type)
	// No synthetic DefaultValue is injected; IsNonPointerInSDK recognizes
	// intEnum directly.
	assert.Equal("", apiShape.DefaultValue)
	assert.False(apiShape.HasDefaultValue())
}

// TestCreateApiShape_Integer verifies a plain integer shape is unaffected by
// the intEnum normalization (no spurious DefaultValue is injected, so it stays
// a pointer SDK field).
func TestCreateApiShape_Integer(t *testing.T) {
	assert := assert.New(t)

	apiShape, err := createApiShape(Shape{Type: "integer"})
	assert.NoError(err)
	assert.Equal("integer", apiShape.Type)
	assert.Equal("", apiShape.DefaultValue)
	assert.False(apiShape.HasDefaultValue())
}
