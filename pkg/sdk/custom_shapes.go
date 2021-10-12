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

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

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

	for _, memberShape := range h.getCustomMapFieldMembers() {
		customShape, err := injector.newMap(memberShape)
		if err != nil {
			return err
		}

		sdkapi.API.Shapes[customShape.Shape.ShapeName] = customShape.Shape
		sdkapi.CustomShapes = append(sdkapi.CustomShapes, customShape)
	}

	for _, memberShape := range h.getCustomListFieldMembers() {
		customShape, err := injector.newList(memberShape)
		if err != nil {
			return err
		}

		sdkapi.API.Shapes[customShape.Shape.ShapeName] = customShape.Shape
		sdkapi.CustomShapes = append(sdkapi.CustomShapes, customShape)
	}

	return nil
}

// getCustomListFieldMembers finds all of the custom list fields that need to
// be generated as defined in the generator config.
func (h *Helper) getCustomListFieldMembers() []string {
	members := []string{}

	for _, resource := range h.cfg.Resources {
		for _, field := range resource.Fields {
			if field.CustomField != nil && field.CustomField.ListOf != "" {
				members = append(members, field.CustomField.ListOf)
			}
		}
	}

	return members
}

// getCustomMapFieldMembers finds all of the custom map fields that need to be
// generated as defined in the generator config.
func (h *Helper) getCustomMapFieldMembers() []string {
	members := []string{}

	for _, resource := range h.cfg.Resources {
		for _, field := range resource.Fields {
			if field.CustomField != nil && field.CustomField.MapOf != "" {
				members = append(members, field.CustomField.MapOf)
			}
		}
	}

	return members
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

// newList loads a shape given its name and creates a custom shape that is a
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
