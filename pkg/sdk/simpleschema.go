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
	"fmt"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	simpleschema "github.com/kro-run/kro/pkg/simpleschema"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// InjectSimpleSchemaShapes processes custom shapes from the top-level custom_shapes
// section and injects them into the SDK API model.
// Only string types are supported - all other types will cause a panic.
func (h *Helper) InjectSimpleSchemaShapes(sdkapi *ackmodel.SDKAPI) error {
	if len(h.cfg.CustomShapes) == 0 {
		return nil
	}

	apiShapeNames := sdkapi.API.ShapeNames()
	for shapeName, fieldsMap := range h.cfg.CustomShapes {
		// check dublicates
		for _, as := range apiShapeNames {
			if as == shapeName {
				return fmt.Errorf("shapeName %s already exists in the API", shapeName)
			}
		}
		schemaObj := convertMapValues(fieldsMap)
		openAPISchema, err := simpleschema.ToOpenAPISpec(schemaObj)
		if err != nil {
			return err
		}

		// Create and register the base structure shape
		shape, shapeRef := h.createBaseShape(sdkapi, shapeName, openAPISchema)
		sdkapi.API.Shapes[shape.ShapeName] = shape
		sdkapi.CustomShapes = append(sdkapi.CustomShapes, &ackmodel.CustomShape{
			Shape:           shape,
			ShapeRef:        shapeRef,
			MemberShapeName: nil,
			ValueShapeName:  nil,
		})
	}

	return nil
}

// convertFieldsMapToSchemaObj converts a fields map to a schema object for OpenAPI spec generation
func convertMapValues(fieldsMap map[string]string) map[string]interface{} {
	schemaObj := make(map[string]interface{})
	for fieldName, fieldType := range fieldsMap {
		schemaObj[fieldName] = fieldType
	}
	return schemaObj
}

// createBaseShape creates a base structure shape with its member fields
func (h *Helper) createBaseShape(
	sdkapi *ackmodel.SDKAPI,
	shapeName string,
	openAPISchema *apiextv1.JSONSchemaProps,
) (*awssdkmodel.Shape, *awssdkmodel.ShapeRef) {
	injector := customShapeInjector{sdkapi}
	shape := &awssdkmodel.Shape{
		API:           sdkapi.API,
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
		addStringFieldToShape(sdkapi, shape, fieldName, shapeName)
	}

	shapeRef := injector.createShapeRefForMember(shape)
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
