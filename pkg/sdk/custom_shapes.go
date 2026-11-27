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

package sdk

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	simpleschema "github.com/kubernetes-sigs/kro/pkg/simpleschema"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
)

var (
	ErrMemberShapeNotFound = errors.New("base shape not found")
)

const (
	ShapeNameTemplateList = "%sList"
	ShapeNameTemplateMap  = "%sMap"
	ShapeNameTemplateKey  = "%sKey"
)

type customShapeInjector struct {
	sdkAPI *ackmodel.SDKAPI
}

// InjectCustomShapes will create custom shapes for each of the spec and status
// fields that contain CustomFieldConfig values. It will append these values
// into the list of shapes in the API and update the list of custom shapes in
// the SDKAPI object.
func (h *Helper) InjectCustomShapes(sdkapi *ackmodel.SDKAPI) error {
	injector := customShapeInjector{sdkapi}

	if err := injector.injectSimpleSchemaShapes(h.cfg.CustomShapes); err != nil {
		return err
	}

	for _, memberShape := range h.cfg.GetCustomMapFieldMembers() {
		customShape, err := injector.newMap(memberShape)
		if err != nil {
			return err
		}

		sdkapi.API.Shapes[customShape.Shape.ShapeName] = customShape.Shape
		sdkapi.CustomShapes = append(sdkapi.CustomShapes, customShape)
	}
	for _, memberShape := range h.cfg.GetCustomListFieldMembers() {
		customShape, err := injector.newList(memberShape)
		if err != nil {
			return err
		}

		sdkapi.API.Shapes[customShape.Shape.ShapeName] = customShape.Shape
		sdkapi.CustomShapes = append(sdkapi.CustomShapes, customShape)
	}
	return nil
}

// injectSimpleSchemaShapes processes custom shapes from the top-level custom_shapes
// section and injects them into the SDK API model.
//
// A member's type is a scalar (string, boolean, integer, float), the name of
// another custom shape, or a list or map of either ("[]Name",
// "map[string]Name"). Shapes are created before their members are resolved,
// so custom shapes can refer to each other in any order.
func (i *customShapeInjector) injectSimpleSchemaShapes(customShapes map[string]map[string]interface{}) error {
	if len(customShapes) == 0 {
		return nil
	}

	apiShapeNames := i.sdkAPI.API.ShapeNames()
	shapeNames := make([]string, 0, len(customShapes))
	for shapeName := range customShapes {
		for _, as := range apiShapeNames {
			if as == shapeName {
				return fmt.Errorf("CustomType name %s already exists in the API", shapeName)
			}
		}
		shapeNames = append(shapeNames, shapeName)
	}
	sort.Strings(shapeNames)

	// Validate the simpleschema syntax, with the other custom shapes as
	// known types.
	customTypes := make(map[string]interface{}, len(customShapes))
	for name, fields := range customShapes {
		customTypes[name] = fields
	}
	for _, shapeName := range shapeNames {
		if _, err := simpleschema.ToOpenAPISpec(customShapes[shapeName], customTypes); err != nil {
			return fmt.Errorf("custom shape %s: %w", shapeName, err)
		}
	}

	shapes := make(map[string]*awssdkmodel.Shape, len(shapeNames))
	for _, shapeName := range shapeNames {
		shape := &awssdkmodel.Shape{
			API:           i.sdkAPI.API,
			ShapeName:     shapeName,
			Type:          "structure",
			Documentation: "// Custom ACK type for " + shapeName,
			MemberRefs:    make(map[string]*awssdkmodel.ShapeRef),
		}
		shapes[shapeName] = shape
		i.sdkAPI.API.Shapes[shapeName] = shape
		i.sdkAPI.CustomShapes = append(i.sdkAPI.CustomShapes, &ackmodel.CustomShape{
			Shape:           shape,
			ShapeRef:        i.createShapeRefForMember(shape),
			MemberShapeName: nil,
			ValueShapeName:  nil,
		})
	}

	for _, shapeName := range shapeNames {
		shape := shapes[shapeName]
		fieldNames := make([]string, 0, len(customShapes[shapeName]))
		for fieldName := range customShapes[shapeName] {
			fieldNames = append(fieldNames, fieldName)
		}
		sort.Strings(fieldNames)
		for _, fieldName := range fieldNames {
			spec, ok := customShapes[shapeName][fieldName].(string)
			if !ok {
				return fmt.Errorf("custom shape %s, field %s: nested objects are not supported, declare a named custom shape instead", shapeName, fieldName)
			}
			typ := strings.TrimSpace(strings.SplitN(spec, "|", 2)[0])
			memberShape, err := i.memberShape(shapes, shapeName+fieldName, typ)
			if err != nil {
				return fmt.Errorf("custom shape %s, field %s: %w", shapeName, fieldName, err)
			}
			shape.MemberRefs[fieldName] = i.createShapeRefForMember(memberShape)
		}
	}
	return nil
}

// simpleSchemaScalarTypes maps simpleschema scalar types to SDK shape types.
var simpleSchemaScalarTypes = map[string]string{
	"string":  "string",
	"boolean": "boolean",
	"integer": "long",
	"float":   "double",
}

// memberShape returns the SDK shape for a custom shape member's type. Shapes
// built here are named name, which must be unique: the field model detects cycles
// by shape name.
func (i *customShapeInjector) memberShape(
	shapes map[string]*awssdkmodel.Shape,
	name string,
	typ string,
) (*awssdkmodel.Shape, error) {
	switch {
	case strings.HasPrefix(typ, "[]"):
		elem, err := i.memberShape(shapes, name+"Member", typ[2:])
		if err != nil {
			return nil, err
		}
		return &awssdkmodel.Shape{
			API:       i.sdkAPI.API,
			ShapeName: name,
			Type:      "list",
			MemberRef: *i.createShapeRefForMember(elem),
		}, nil
	case strings.HasPrefix(typ, "map[string]"):
		elem, err := i.memberShape(shapes, name+"Value", typ[len("map[string]"):])
		if err != nil {
			return nil, err
		}
		return &awssdkmodel.Shape{
			API:       i.sdkAPI.API,
			ShapeName: name,
			Type:      "map",
			KeyRef:    *i.createShapeRefForMember(&awssdkmodel.Shape{API: i.sdkAPI.API, ShapeName: name + "Key", Type: "string"}),
			ValueRef:  *i.createShapeRefForMember(elem),
		}, nil
	}
	if sdkType, ok := simpleSchemaScalarTypes[typ]; ok {
		return &awssdkmodel.Shape{API: i.sdkAPI.API, ShapeName: name, Type: sdkType}, nil
	}
	if shape, ok := shapes[typ]; ok {
		return shape, nil
	}
	return nil, fmt.Errorf("unsupported type %q", typ)
}

// createShapeRefForMember creates a minimal ShapeRef type to encapsulate a
// shape.
func (i *customShapeInjector) createShapeRefForMember(shape *awssdkmodel.Shape) *awssdkmodel.ShapeRef {
	return &awssdkmodel.ShapeRef{
		API:           i.sdkAPI.API,
		Shape:         shape,
		Documentation: shape.Documentation,
		ShapeName:     shape.ShapeName,
	}
}

// createKeyShape creates a Shape that acts as the string key shape for a
// custom map.
func (i *customShapeInjector) createKeyShape(shapeName string) *awssdkmodel.Shape {
	return &awssdkmodel.Shape{
		API:       i.sdkAPI.API,
		ShapeName: fmt.Sprintf(ShapeNameTemplateKey, shapeName),
		Type:      "string",
	}
}

// newMap loads a shape given its name and creates a custom shape that is a
// map with strings as keys and that shape as the value.
func (i *customShapeInjector) newMap(valueShapeName string) (*ackmodel.CustomShape, error) {
	valueShape, exists := i.sdkAPI.API.Shapes[valueShapeName]
	if !exists {
		return nil, ErrMemberShapeNotFound
	}
	valueShapeRef := i.createShapeRefForMember(valueShape)

	shapeName := fmt.Sprintf(ShapeNameTemplateMap, valueShape.ShapeName)
	documentation := ""

	keyShape := i.createKeyShape(shapeName)
	keyShapeRef := i.createShapeRefForMember(keyShape)

	shape := &awssdkmodel.Shape{
		API:       i.sdkAPI.API,
		ShapeName: shapeName,
		// TODO (RedbackThomson): Support documentation for custom shapes
		Documentation: documentation,
		KeyRef:        *keyShapeRef,
		ValueRef:      *valueShapeRef,
		Type:          "map",
	}

	shapeRef := &awssdkmodel.ShapeRef{
		API:           i.sdkAPI.API,
		Shape:         shape,
		Documentation: documentation,
		ShapeName:     shapeName,
	}

	return ackmodel.NewCustomMapShape(shape, shapeRef, valueShapeName), nil
}

// newList loads a shape given its name and creates a custom shape that is a
// list of that shape.
func (i *customShapeInjector) newList(memberShapeName string) (*ackmodel.CustomShape, error) {
	memberShape, exists := i.sdkAPI.API.Shapes[memberShapeName]
	if !exists {
		return nil, ErrMemberShapeNotFound
	}
	memberShapeRef := i.createShapeRefForMember(memberShape)

	shapeName := fmt.Sprintf(ShapeNameTemplateList, memberShape.ShapeName)
	documentation := ""

	shape := &awssdkmodel.Shape{
		API:       i.sdkAPI.API,
		ShapeName: shapeName,
		// TODO (RedbackThomson): Support documentation for custom shapes
		Documentation: documentation,
		MemberRef:     *memberShapeRef,
		Type:          "list",
	}

	shapeRef := &awssdkmodel.ShapeRef{
		API:           i.sdkAPI.API,
		Shape:         shape,
		Documentation: documentation,
		ShapeName:     shapeName,
	}

	return ackmodel.NewCustomListShape(shape, shapeRef, memberShapeName), nil
}
