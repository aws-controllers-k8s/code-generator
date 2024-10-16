package apiv2

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	sdkmodel "github.com/aws/aws-sdk-go/private/model/api"
)

var (
	serviceAlias string
	shapes       = make(map[string]*sdkmodel.Shape)
	operations   = make(map[string]*sdkmodel.Operation)
	raw_shapes   map[string]interface{}
)

type ServiceObjV2 map[string]interface{}
type APIs map[string]*sdkmodel.API

func (Obj *ServiceObjV2) RemovePrefix(key string, serviceAlias string) *string {

	to_Trim := "com.amazonaws." + serviceAlias + "#"
	new_key := strings.TrimPrefix(key, to_Trim)

	return &new_key
}

func (Obj *ServiceObjV2) ReplaceShapeSuffixRequest(key string, serviceAlias string) *string {
	if strings.HasSuffix(key, "Request") {
		new_key := strings.Replace(key, "Request", "Input", 1)
		return &new_key
	} else if strings.HasSuffix(key, "Response") {
		new_key := strings.Replace(key, "Response", "Output", 1)
		return &new_key
	}

	return &key
}

func (Obj *ServiceObjV2) CapitaliseName(name *string) string {

	runes := []rune(*name)
	if unicode.IsUpper(runes[0]) {
		return *name
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)

}

func ReadV2file(serviceAlias string) (*ServiceObjV2, error) {

	var Obj ServiceObjV2

	dir, _ := os.Getwd()
	filepath := filepath.Join(dir + "/apiv2/" + serviceAlias + ".json")

	file, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &Obj)
	if err != nil {
		return nil, err
	}
	return &Obj, nil
}

func CheckShapes(Obj *ServiceObjV2) (interface{}, error) {
	if value, ok := (*Obj)["shapes"]; ok {
		raw_shapes = value.(map[string]interface{})
		return raw_shapes, nil
	}
	return nil, errors.New("shapes not found")
}

func GetRawShape(key string) map[string]interface{} {

	if raw_shape, ok := raw_shapes[key]; ok {

		return raw_shape.(map[string]interface{})
	}

	return nil
}

func BuildModelV2(m *ModelV2Shape) *sdkmodel.Shape {

	if m.GetShapeType() == "structure" {
		newShape := m.CreateShape(m.Name)
		newShape = m.AddStructMembers(newShape)

		return newShape

	} else if m.GetShapeType() == "list" {
		newShape := m.CreateShape(m.Name)

		newShape = m.AddListRef(newShape)

		return newShape

	} else if m.GetShapeType() == "enum" {
		newShape := m.CreateShape(m.Name)
		newShape = m.AddEnumValues(newShape)

		return newShape

	} else if m.GetShapeType() == "integer" || m.GetShapeType() == "long" || m.GetShapeType() == "string" || m.GetShapeType() == "boolean" {
		newShape := m.CreateShape(m.Name)

		return newShape

	}
	return nil
}

func LoadModelV2(shapes map[string]*sdkmodel.Shape, operations map[string]*sdkmodel.Operation, serviceAlias string) map[string]*sdkmodel.API {

	api := sdkmodel.API{
		Shapes:     shapes,
		Operations: operations,
	}

	api.StrictServiceId = true
	api.Metadata.ServiceID = serviceAlias

	apis := APIs{}
	/// Need to add pull path here as key in APIsmap
	apis[serviceAlias] = &api

	return apis

}

// Write a main function to call ReadV2file and PrintV2Obj
func CollectApis(modelPath string) map[string]*sdkmodel.API {

	Obj, err := ReadV2file(modelPath)

	if err != nil {
		panic(err.Error())

	}

	if raw_shapes, err := CheckShapes(Obj); err != nil {

		panic(err.Error())

	} else {

		for key, raw_shape := range raw_shapes.(map[string]interface{}) {

			name := Obj.RemovePrefix(key, serviceAlias)
			name = Obj.ReplaceShapeSuffixRequest(*name, serviceAlias)

			m := NewModelV2Shape(raw_shape, *name, serviceAlias)
			newShape := BuildModelV2(m)

			if newShape != nil {

				shapes[newShape.ShapeName] = newShape

			}

		}

		for key, raw_shape := range raw_shapes.(map[string]interface{}) {

			name := Obj.RemovePrefix(key, serviceAlias)
			name = Obj.ReplaceShapeSuffixRequest(*name, serviceAlias)

			m := NewModelV2Shape(raw_shape, *name, serviceAlias)

			if m.GetShapeType() == "operation" {

				op := m.CreateOperation(m.Name)

				if input, ok := m.Raw_shape["input"]; ok {
					op = m.AddInputShapeRef(op, input)

				}
				if output, ok := m.Raw_shape["output"]; ok {
					op = m.AddOutputShapeRef(op, output)

				}

				operations[op.Name] = op

			}

		}

	}
	ApisV2 := LoadModelV2(shapes, operations, serviceAlias)
	return ApisV2
}
