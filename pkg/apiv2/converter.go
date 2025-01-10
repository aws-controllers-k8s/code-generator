package apiv2

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/aws-controllers-k8s/code-generator/pkg/api"
)

// API holds all the shapes defined in the <service>.json
// api model file provided by aws-sdk-go-v2
type API struct {
	Shapes map[string]Shape `json:"shapes"`
}

// Shape contains the definition of a resource, field,
// operation, etc.
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

// ShapeRef defines the usage of a shape within the API
type ShapeRef struct {
	API       *API   `json:"-"`
	Shape     *Shape `json:"-"`
	ShapeName string `json:"target"`
	Traits    map[string]interface{}
}

// ConvertAPIV2Shapes loads the V2 api model file, and later translates that
// structure into the v1 API structure.
func ConvertApiV2Shapes(serviceAlias, modelPath string) (map[string]*api.API, error) {
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

	// build V1 API structure
	newApi, err := buildAPI(customAPI.Shapes, serviceAlias)
	if err != nil {
		return nil, fmt.Errorf("error building api: %v", err)
	}

	// Setup the API (make sure the shapes of refs are in the right place)
	// (documentation is clear among the shapes)
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
func buildAPI(shapes map[string]Shape, serviceAlias string) (*api.API, error) {

	newApi := api.API{
		Metadata:        api.Metadata{},
		Operations:      map[string]*api.Operation{},
		Shapes:          map[string]*api.Shape{},
		StrictServiceId: true,
	}

	for shapeName, shape := range shapes {
		name := removeShapeNamePrefix(shapeName, serviceAlias)
		// handling the creation of service and operation types
		// differently
		if shape.Type != "service" && shape.Type != "operation" {
			newShape, err := createApiShape(shape)
			if err != nil {
				return nil, err
			}
			newApi.Shapes[name] = newShape
		}

		// Ignoring types String, Integer, and boolean since they will be added as MemberRefs
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
			addMemberRefs(newApi.Shapes[name], shape, serviceAlias)
		case "list":
			addMemberRef(newApi.Shapes[name], shape, serviceAlias)
		case "map":
			addKeyAndValueRef(newApi.Shapes[name], shape, serviceAlias)
		case "enum":
			newApi.Shapes[name].Type = "string"
			addEnumRef(newApi.Shapes[name], shape)
		// Union, introduced in S3 model file, is a structure
		case "union":
			newApi.Shapes[name].Type = "structure"
			addMemberRefs(newApi.Shapes[name], shape, serviceAlias)
		}

	}

	return &newApi, nil
}

func createApiOperation(shape Shape, name, serviceAlias string) *api.Operation {
	// Some operations may not have documentation
	doc, _ := shape.Traits["smithy.api#documentation"].(string)

	newOperation := &api.Operation{
		Name:          name,
		Documentation: api.AppendDocstring("", doc),
	}

	if hasPrefix(shape.InputRef.ShapeName, serviceAlias) {
		inputName := removeShapeNamePrefix(shape.InputRef.ShapeName, serviceAlias)
		newOperation.InputRef = api.ShapeRef{
			ShapeName: inputName,
		}
	}
	if hasPrefix(shape.OutputRef.ShapeName, serviceAlias) {
		outputName := removeShapeNamePrefix(shape.OutputRef.ShapeName, serviceAlias)
		newOperation.OutputRef = api.ShapeRef{
			ShapeName: outputName,
		}
	}

	for _, err := range shape.ErrorRefs {
		newOperation.ErrorRefs = append(newOperation.ErrorRefs, api.ShapeRef{
			ShapeName: removeShapeNamePrefix(err.ShapeName, serviceAlias),
		})
	}

	return newOperation
}

func createApiShape(shape Shape) (*api.Shape, error) {
	isException := shape.IsException()
	apiShape := &api.Shape{
		Type:       shape.Type,
		Exception:  isException,
		MemberRefs: make(map[string]*api.ShapeRef),
		MemberRef:  api.ShapeRef{},
		KeyRef:     api.ShapeRef{},
		ValueRef:   api.ShapeRef{},
		Required:   []string{},
	}
	// Shapes that are default so far are booleans and
	// integers that are non pointers in the service API
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
	documentation, _ := shape.Traits["smithy.api#documentation"].(string)
	apiShape.Documentation = api.AppendDocstring("", documentation)

	return apiShape, nil
}

func addMemberRefs(apiShape *api.Shape, shape Shape, serviceAlias string) {
	for memberName, member := range shape.MemberRefs {
		if !hasPrefix(member.ShapeName, serviceAlias) {
			continue
		}
		shapeNameClean := removeShapeNamePrefix(member.ShapeName, serviceAlias)
		documentation, _ := member.Traits["smithy.api#documentation"].(string)
		if member.isRequired() {
			apiShape.Required = append(apiShape.Required, memberName)
		}
		apiShape.MemberRefs[memberName] = &api.ShapeRef{
			ShapeName:     shapeNameClean,
			Documentation: api.AppendDocstring("", documentation),
		}
	}
	slices.Sort(apiShape.Required)
}

func addMemberRef(apiShape *api.Shape, shape Shape, serviceAlias string) {

	apiShape.MemberRef = api.ShapeRef{
		ShapeName: removeShapeNamePrefix(shape.MemberRef.ShapeName, serviceAlias),
	}
}

func addKeyAndValueRef(apiShape *api.Shape, shape Shape, serviceAlias string) {

	apiShape.KeyRef = api.ShapeRef{
		ShapeName: removeShapeNamePrefix(shape.KeyRef.ShapeName, serviceAlias),
	}
	apiShape.ValueRef = api.ShapeRef{
		ShapeName: removeShapeNamePrefix(shape.ValueRef.ShapeName, serviceAlias),
	}
}

func addEnumRef(apiShape *api.Shape, shape Shape) {
	for memberName := range shape.MemberRefs {
		apiShape.Enum = append(apiShape.Enum, memberName)
	}
	slices.Sort(apiShape.Enum)
}

func (s ShapeRef) isRequired() bool {
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

func removeShapeNamePrefix(name, alias string) string {
	toTrim := fmt.Sprintf("com.amazonaws.%s#", alias)

	newName := strings.TrimPrefix(name, toTrim)

	return newName
}
