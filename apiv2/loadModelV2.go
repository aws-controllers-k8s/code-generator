package apiv2

import (
	"strings"

	sdkmodel "github.com/aws/aws-sdk-go/private/model/api"
)

type ModelV2Shape struct {
	Name         string
	Raw_shape    map[string]interface{}
	ServiceAlias string
	ServiceObjV2
}

func NewModelV2Shape(raw_shape interface{}, key string, serviceAlias string) *ModelV2Shape {

	return &ModelV2Shape{
		Name:         key,
		Raw_shape:    raw_shape.(map[string]interface{}),
		ServiceAlias: serviceAlias,
	}
}

func (m *ModelV2Shape) GetShapeType() string {

	for key, value := range m.Raw_shape {
		if key == "type" {
			switch value.(string) {
			case "structure":
				return "structure"
			case "map":
				return "map"
			case "list":
				return "list"
			case "string":
				return "string"
			case "integer":
				return "integer"
			case "long":
				return "long"
			case "double":
				return "double"
			case "boolean":
				return "boolean"
			case "timestamp":
				return "timestamp"
			case "enum":
				return "enum"
			case "operation":
				return "operation"
			case "blob":
				return "blob"
			default:
				return "unknown"
			}
		}
	}
	return "unknown"
}

func (m *ModelV2Shape) CreateShape(name string) *sdkmodel.Shape {

	api := sdkmodel.API{}

	api.AddImport("time")

	new_shape := &sdkmodel.Shape{
		ShapeName:  *m.RemovePrefix(name, m.ServiceAlias),
		Exception:  false,
		MemberRefs: make(map[string]*sdkmodel.ShapeRef),
		MemberRef:  sdkmodel.ShapeRef{},
		API:        &api,
	}

	if Type, ok := m.Raw_shape["type"]; ok {
		if Type.(string) == "enum" {
			new_shape.Type = "string"
		} else {
			new_shape.Type = Type.(string)
		}
	}

	return new_shape

}

func (m *ModelV2Shape) AddDefaultMembers(new_shape *sdkmodel.Shape) *sdkmodel.Shape {

	if m.GetShapeType() != "structure" || m.GetShapeType() != "list" || m.GetShapeType() != "enum" {
		for key, _ := range m.Raw_shape {
			new_shape.MemberRefs[key] = &sdkmodel.ShapeRef{
				Shape:     m.CreateShape(*m.RemovePrefix(key, m.ServiceAlias)),
				ShapeName: *m.RemovePrefix(key, m.ServiceAlias),
			}
		}
	}

	return new_shape
}

func (m *ModelV2Shape) AddStructMembers(new_shape *sdkmodel.Shape) *sdkmodel.Shape {

	if m.GetShapeType() == "structure" {

		for key, members := range m.Raw_shape["members"].(map[string]interface{}) {

			if value, ok := members.(map[string]interface{})["target"]; ok {

				valueRawShape := GetRawShape(value.(string))
				memberModel := NewModelV2Shape(valueRawShape, value.(string), m.ServiceAlias)
				key = memberModel.CapitaliseName(&key)

				new_shape.MemberRefs[key] = &sdkmodel.ShapeRef{
					Shape:     memberModel.CreateShape(*memberModel.RemovePrefix(value.(string), memberModel.ServiceAlias)),
					ShapeName: *memberModel.RemovePrefix(value.(string), memberModel.ServiceAlias),
				}

				switch new_shape.MemberRefs[key].Shape.Type {
				case "list":
					new_shape.MemberRefs[key].Shape = memberModel.AddListRef(new_shape.MemberRefs[key].Shape)
				case "structure":
					new_shape.MemberRefs[key].Shape = memberModel.AddStructMembers(new_shape.MemberRefs[key].Shape)
				}
				if memberModel.GetShapeType() == "enum" {
					new_shape.MemberRefs[key].Shape = memberModel.AddEnumValues(new_shape.MemberRefs[key].Shape)
				}

			}

		}
	}
	return new_shape
}

func (m *ModelV2Shape) AddListRef(new_shape *sdkmodel.Shape) *sdkmodel.Shape {

	if m.GetShapeType() == "list" {
		for _, members := range m.Raw_shape["member"].(map[string]interface{}) {

			if value, ok := members.(string); ok {

				valueRawShape := GetRawShape(value)
				memberModel := NewModelV2Shape(valueRawShape, new_shape.ShapeName, m.ServiceAlias)

				new_shape.MemberRef = sdkmodel.ShapeRef{
					Shape:     memberModel.CreateShape(*memberModel.RemovePrefix(value, memberModel.ServiceAlias)),
					ShapeName: *memberModel.RemovePrefix(value, memberModel.ServiceAlias),
				}

				switch new_shape.MemberRef.Shape.Type {
				case "structure":
					new_shape.MemberRef.Shape = memberModel.AddStructMembers(new_shape.MemberRef.Shape)
				default:
					new_shape.MemberRef.Shape = memberModel.AddDefaultMembers(new_shape.MemberRef.Shape)

				}

			}

		}
	}
	return new_shape
}

func (m *ModelV2Shape) AddEnumValues(new_shape *sdkmodel.Shape) *sdkmodel.Shape {

	if m.GetShapeType() == "enum" {
		for key := range m.Raw_shape["members"].(map[string]interface{}) {
			new_shape.Enum = append(new_shape.Enum, key)
		}
	}

	return new_shape
}

func (m *ModelV2Shape) CreateOperation(name string) *sdkmodel.Operation {

	if m.GetShapeType() == "operation" {
		new_operation := &sdkmodel.Operation{
			Name:         name,
			ExportedName: *m.RemovePrefix(name, serviceAlias),
		}
		return new_operation

	}
	return nil

}

func (m *ModelV2Shape) AddInputShapeRef(newOperation *sdkmodel.Operation, input interface{}) *sdkmodel.Operation {

	inputShapeName := input.(map[string]interface{})["target"]
	inputRawShape := GetRawShape(inputShapeName.(string))

	name := m.RemovePrefix(inputShapeName.(string), m.ServiceAlias)
	name = m.ReplaceShapeSuffixRequest(*name, m.ServiceAlias)

	inputShapeModel := NewModelV2Shape(inputRawShape, *name, m.ServiceAlias)

	if _, ok := shapes[inputShapeModel.Name]; ok {
		newOperation.InputRef = sdkmodel.ShapeRef{
			Shape:     shapes[inputShapeModel.Name],
			ShapeName: inputShapeModel.Name,
		}
	} else {
		newOperation.InputRef = sdkmodel.ShapeRef{
			Shape:     inputShapeModel.CreateShape(inputShapeModel.Name),
			ShapeName: inputShapeModel.Name,
		}
		shapes[newOperation.InputRef.ShapeName] = newOperation.InputRef.Shape
	}

	return newOperation

}

func (m *ModelV2Shape) AddOutputShapeRef(newOperation *sdkmodel.Operation, output interface{}) *sdkmodel.Operation {

	outputShapeName := output.(map[string]interface{})["target"]
	outputRawShape := GetRawShape(outputShapeName.(string))

	name := m.RemovePrefix(outputShapeName.(string), m.ServiceAlias)
	name = m.ReplaceShapeSuffixRequest(*name, m.ServiceAlias)
	outputShapeModel := NewModelV2Shape(outputRawShape, *name, m.ServiceAlias)

	if outputShapeModel.Name == "smithy.api#Unit" {

		newOperation.OutputRef = sdkmodel.ShapeRef{
			Shape:     outputShapeModel.CreateShape(newOperation.Name + "Output"),
			ShapeName: newOperation.Name + "Output",
		}
		newOperation.OutputRef.Shape.Type = "structure"
		shapes[newOperation.OutputRef.ShapeName] = newOperation.OutputRef.Shape

		return newOperation

	}

	if OpOutputRef, ok := shapes[outputShapeModel.Name]; ok {

		var OutputShapeName string
		if strings.HasSuffix(OpOutputRef.ShapeName, "Description") {

			OutputShapeName = newOperation.Name + "Output"
			outputShapeModel.Name = OutputShapeName
			outPutShape := BuildModelV2(outputShapeModel)
			newOperation.OutputRef.Shape = outPutShape
			return newOperation

		}
		newOperation.OutputRef.Shape = OpOutputRef

	}

	return newOperation

}
