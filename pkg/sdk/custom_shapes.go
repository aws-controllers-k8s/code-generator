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

	simpleschema "github.com/kubernetes-sigs/kro/pkg/simpleschema"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

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
// Only string types are supported - all other types will cause a panic.
func (i *customShapeInjector) injectSimpleSchemaShapes(customShapes map[string]map[string]interface{}) error {
	if len(customShapes) == 0 {
		return nil
	}

	apiShapeNames := i.sdkAPI.API.ShapeNames()
	for shapeName, fieldsMap := range customShapes {
		// check for duplicates
		for _, as := range apiShapeNames {
			if as == shapeName {
				return fmt.Errorf("CustomType name %s already exists in the API", shapeName)
			}
		}
		openAPISchema, err := simpleschema.ToOpenAPISpec(fieldsMap, nil)
		if err != nil {
			return err
		}

		// Create and register the base structure shape
		shape, shapeRef := i.newStructureShape(shapeName, openAPISchema)
		i.sdkAPI.API.Shapes[shape.ShapeName] = shape
		i.sdkAPI.CustomShapes = append(i.sdkAPI.CustomShapes, &ackmodel.CustomShape{
			Shape:           shape,
			ShapeRef:        shapeRef,
			MemberShapeName: nil,
			ValueShapeName:  nil,
		})
	}

	return nil
}

// newStructureShape creates a base shape with its member fields
func (i *customShapeInjector) newStructureShape(
	shapeName string,
	openAPISchema *apiextv1.JSONSchemaProps,
) (*awssdkmodel.Shape, *awssdkmodel.ShapeRef) {
	shape := &awssdkmodel.Shape{
		API:           i.sdkAPI.API,
		ShapeName:     shapeName,
		Type:          "structure",
		Documentation: "// Custom ACK type for " + shapeName,
		MemberRefs:    make(map[string]*awssdkmodel.ShapeRef),
	}

	properties := openAPISchema.Properties
	for fieldName, propObj := range properties {
		propType := propObj.Type
		if propType != "string" {
			panic(fmt.Sprintf("Field %s in shape %s has non-string type '%s'",
				fieldName, shapeName, propType))
		}
		addStringFieldToShape(i.sdkAPI, shape, fieldName, shapeName)
	}

	shapeRef := i.createShapeRefForMember(shape)
	return shape, shapeRef
}

// addStringFieldToShape adds a string field to the parent shape
func addStringFieldToShape(
	sdkapi *ackmodel.SDKAPI,
	parentShape *awssdkmodel.Shape,
	fieldName string,
	shapeName string,
) {
	injector := customShapeInjector{sdkapi}
	fieldShape := &awssdkmodel.Shape{
		API:       sdkapi.API,
		ShapeName: fieldName,
		Type:      "string",
	}

	sdkapi.API.Shapes[fieldShape.ShapeName] = fieldShape
	parentShape.MemberRefs[fieldName] = injector.createShapeRefForMember(fieldShape)
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
