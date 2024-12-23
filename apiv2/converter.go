package apiv2

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aws-controllers-k8s/code-generator/pkg/api"
)

type API struct {
	Shapes map[string]Shape
}

type Shape struct {
	Type       string
	Traits     map[string]interface{}
	MemberRefs map[string]*ShapeRef `json:"members"`
	MemberRef  *ShapeRef            `json:"member"`
	KeyRef     ShapeRef             `json:"key"`
	ValueRef   ShapeRef             `json:"value"`
	InputRef   ShapeRef             `json:"input"`
	OutputRef  ShapeRef             `json:"output"`
	ErrorRefs  []ShapeRef           `json:"errors"`
}

type ShapeRef struct {
	API       *API   `json:"-"`
	Shape     *Shape `json:"-"`
	ShapeName string `json:"target"`
	Traits    map[string]interface{}
}

func ConvertApiV2Shapes(modelPath string) (map[string]*api.API, error) {

	// Read the json file
	file, err := os.ReadFile(modelPath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// unmarshal the file
	var customAPI API
	err = json.Unmarshal(file, &customAPI)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling file: %v", err)
	}

	serviceAlias := extractServiceAlias(modelPath)

	newApi, err := BuildAPI(customAPI.Shapes, serviceAlias)
	if err != nil {
		return nil, fmt.Errorf("error building api: %v", err)
	}

	newApi.StrictServiceId = true

	err = newApi.Setup()
	if err != nil {
		return nil, fmt.Errorf("error setting up api: %v", err)
	}

	return map[string]*api.API{
		serviceAlias: newApi,
	}, nil
}

// This function tries to translate the API from sdk go v2
// into the struct for sdk go v1
func BuildAPI(shapes map[string]Shape, serviceAlias string) (*api.API, error) {

	newApi := api.API{
		Metadata:   api.Metadata{},
		Operations: map[string]*api.Operation{},
		Shapes:     map[string]*api.Shape{},
	}

	for shapeName, shape := range shapes {

		name := removeNamePrefix(shapeName, serviceAlias)
		if shape.Type != "service" && shape.Type != "operation" {
			newShape, err := createApiShape(shape)
			if err != nil {
				return nil, err
			}
			newApi.Shapes[name] = newShape
		}

		switch shape.Type {
		case "service":
			serviceId, ok := shape.Traits["aws.api#service"].(map[string]interface{})["sdkId"]
			if !ok {
				return nil, errors.New("service id not found")
			}
			newApi.Metadata.ServiceID = serviceId.(string)
			doc, ok := shape.Traits["smithy.api#documentation"]
			if !ok {
				return nil, errors.New("service documentation not found")
			}
			newApi.Documentation = api.AppendDocstring("", doc.(string))
		case "operation":
			newApi.Operations[name] = createApiOperation(shape, name, serviceAlias)
		case "structure":
			AddMemberRefs(newApi.Shapes[name], shape, serviceAlias)
		case "list":
			AddMemberRef(newApi.Shapes[name], name, shape, serviceAlias)
		case "map":
			AddKeyAndValueRef(newApi.Shapes[name], name, shape, serviceAlias)
		case "enum":
			AddEnumRef(newApi.Shapes[name], shape)
		}

	}

	return &newApi, nil
}

func createApiOperation(shape Shape, name, serviceAlias string) *api.Operation {

	newOperation := &api.Operation{
		Name:          name,
		Documentation: api.AppendDocstring("", shape.Traits["smithy.api#documentation"].(string)),
	}

	if hasPrefix(shape.InputRef.ShapeName, serviceAlias) {
		inputName := removeNamePrefix(shape.InputRef.ShapeName, serviceAlias)
		newOperation.InputRef = api.ShapeRef{
			ShapeName: inputName,
		}
	}
	if hasPrefix(shape.OutputRef.ShapeName, serviceAlias) {
		outputName := removeNamePrefix(shape.OutputRef.ShapeName, serviceAlias)
		newOperation.OutputRef = api.ShapeRef{
			ShapeName: outputName,
		}
	}

	for _, err := range shape.ErrorRefs {
		newOperation.ErrorRefs = append(newOperation.ErrorRefs, api.ShapeRef{
			ShapeName: removeNamePrefix(err.ShapeName, serviceAlias),
		})
	}

	return newOperation
}

func createApiShape(shape Shape) (*api.Shape, error) {

	isException := shape.IsException()

	shapeType := shape.Type
	if shapeType == "enum" {
		shapeType = "string"
	}

	apiShape := &api.Shape{
		Type:       shapeType,
		Exception:  isException,
		MemberRefs: make(map[string]*api.ShapeRef),
		MemberRef:  api.ShapeRef{},
		KeyRef:     api.ShapeRef{},
		ValueRef:   api.ShapeRef{},
		Required:   []string{},
	}
	val, ok := shape.Traits["smithy.api#default"]
	if ok {
		apiShape.DefaultValue = &val
	}

	if isException {
		code, ok := shape.Traits["smithy.api#httpError"]
		if ok {
			switch code := code.(type) {
			case float64:
				apiShape.ErrorInfo = api.ErrorInfo{
					HTTPStatusCode: int(code),
				}
			case int:
				apiShape.ErrorInfo = api.ErrorInfo{
					HTTPStatusCode: code,
				}
			case int64:
				apiShape.ErrorInfo = api.ErrorInfo{
					HTTPStatusCode: int(code),
				}
			default:
				return nil, fmt.Errorf("status code type not found for exception")
			}
		}
	}

	return apiShape, nil
}

func AddMemberRefs(apiShape *api.Shape, shape Shape, serviceAlias string) {

	var documentation string
	for memberName, member := range shape.MemberRefs {
		if !hasPrefix(member.ShapeName, serviceAlias) {
			continue
		}
		shapeNameClean := removeNamePrefix(member.ShapeName, serviceAlias)
		if member.Traits["smithy.api#documentation"] != nil {
			documentation = api.AppendDocstring("", member.Traits["smithy.api#documentation"].(string))
		}
		if member.IsRequired() {
			apiShape.Required = append(apiShape.Required, memberName)
		}
		apiShape.MemberRefs[memberName] = &api.ShapeRef{
			ShapeName:     shapeNameClean,
			Documentation: documentation,
		}
	}

	if shape.Traits["smithy.api#documentation"] != nil {
		documentation = api.AppendDocstring("", shape.Traits["smithy.api#documentation"].(string))
	}
	// Add the documentation to the shape
	apiShape.Documentation = documentation
}

func AddMemberRef(apiShape *api.Shape, shapeName string, shape Shape, serviceAlias string) {

	apiShape.MemberRef = api.ShapeRef{
		ShapeName: removeNamePrefix(shape.MemberRef.ShapeName, serviceAlias),
	}
}

func AddKeyAndValueRef(apiShape *api.Shape, shapeName string, shape Shape, serviceAlias string) {

	apiShape.KeyRef = api.ShapeRef{
		ShapeName: removeNamePrefix(shape.KeyRef.ShapeName, serviceAlias),
	}
	apiShape.ValueRef = api.ShapeRef{
		ShapeName: removeNamePrefix(shape.ValueRef.ShapeName, serviceAlias),
	}
}

func AddEnumRef(apiShape *api.Shape, shape Shape) {
	for memberName := range shape.MemberRefs {
		apiShape.Enum = append(apiShape.Enum, memberName)
	}
}

func (s ShapeRef) IsRequired() bool {
	_, ok := s.Traits["smithy.api#required"]
	return ok
}

func (s Shape) IsException() bool {
	_, ok := s.Traits["smithy.api#error"]
	return ok
}

func hasPrefix(name, alias string) bool {

	prefix := fmt.Sprintf("com.amazonaws.%s#", alias)

	return strings.HasPrefix(name, prefix)
}

func removeNamePrefix(name, alias string) string {

	toTrim := fmt.Sprintf("com.amazonaws.%s#", alias)

	newName := strings.TrimPrefix(name, toTrim)

	return newName
}

func extractServiceAlias(modelPath string) string {
	// Split the path into parts
	parts := strings.Split(modelPath, "/")

	// Get the last part
	lastPart := parts[len(parts)-1]

	// Split the last part by "." to get the service alias
	serviceAlias := strings.Split(lastPart, ".")[0]

	return serviceAlias
}
