package apiv2

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/pkg/names"
)

// API holds all the shapes defined in the <service>.json
// api model file provided by aws-sdk-go-v2
type API struct {
	Shapes map[string]Shape `json:"shapes"`
}

// Shape contains the definition of a resource, field,
// operation, service, etc.
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
	ShapeName string `json:"target"`
	Traits    map[string]interface{}
}

// ConvertAPIV2Shapes loads the V2 api model file, and later translates that
// structure into the v1 API structure.
func ConvertApiV2Shapes(modelPath string) (map[string]*awssdkmodel.API, error) {
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
	newApi, serviceAlias, err := buildAPI(customAPI.Shapes)
	if err != nil {
		return nil, fmt.Errorf("error building api: %v", err)
	}

	// Setup the API: ensures the shapes of refs (memberRef, keyRef) are in the right place
	// and documentation is clear among the shapes
	err = newApi.Setup()
	if err != nil {
		return nil, fmt.Errorf("error setting up api: %v", err)
	}

	// fmt.Println(*newApi.Shapes["AIMLOptionsInput_"])

	return map[string]*awssdkmodel.API{
		serviceAlias: newApi,
	}, nil
}

// This function tries to translate the API from sdk go v2
// into the API struct from sdk go v1 (now moved to code-gen pkg)
func buildAPI(shapes map[string]Shape) (*awssdkmodel.API, string, error) {

	newApi := awssdkmodel.API{
		Metadata:        awssdkmodel.Metadata{},
		Operations:      map[string]*awssdkmodel.Operation{},
		Shapes:          injectCustomPrimitiveShapes(),
		StrictServiceId: true,
	}

	var serviceAlias string

	for shapeName, shape := range shapes {
		name, err := removeShapeNamePrefix(shapeName)
		if err != nil {
			return nil, "", err
		}
		if serviceAlias == "" {
			serviceAlias = extractServiceAlias(shapeName)
		}
		
		// Since in V2, operation and service are still considered Shape,
		// handling them separately to create 
		if shape.Type != "service" && shape.Type != "operation" {
			newShape, err := createApiShape(shape)
			if err != nil {
				return nil, "", err
			}
			newApi.Shapes[name] = newShape
		}

		// Ignoring types String, Integer, and boolean since they will be added as MemberRefs
		switch shape.Type {
		case "service":
			serviceId, ok := shape.Traits["aws.api#service"].(map[string]interface{})["sdkId"]
			if !ok {
				return nil, "", errors.New("service id not found")
			}
			newApi.Metadata.ServiceID = serviceId.(string)
			doc, ok := shape.Traits["smithy.api#documentation"]
			if !ok {
				return nil, "", errors.New("service documentation not found")
			}
			newApi.Documentation = awssdkmodel.AppendDocstring("", doc.(string))
		case "operation":
			newApi.Operations[name] = createApiOperation(shape, name, serviceAlias)
		case "structure":
			addMemberRefs(newApi.Shapes[name], shape, serviceAlias)
		case "list":
			addMemberRef(newApi.Shapes[name], shape)
		case "map":
			addKeyAndValueRef(newApi.Shapes[name], shape)
		case "enum":
			newApi.Shapes[name].Type = "string"
			addEnumRef(newApi.Shapes[name], shape)
		// Union, introduced in S3 model file, is a structure
		case "union":
			newApi.Shapes[name].Type = "structure"
			addMemberRefs(newApi.Shapes[name], shape, serviceAlias)
		case "string":
			val, ok := shape.Traits["smithy.api#enum"]
			if ok {
				addEnumValues(newApi.Shapes[name], val)
			}
		}

	}

	return &newApi, serviceAlias, nil
}

// createApiOperation creates an awssdkmodel type Operation.
func createApiOperation(shape Shape, name, serviceAlias string) *awssdkmodel.Operation {
	// Some operations may not have documentation
	doc, _ := shape.Traits["smithy.api#documentation"].(string)

	newOperation := &awssdkmodel.Operation{
		Name:          name,
		Documentation: awssdkmodel.AppendDocstring("", doc),
	}

	if hasPrefix(shape.InputRef.ShapeName, serviceAlias) {
		inputName, _:= removeShapeNamePrefix(shape.InputRef.ShapeName)
		newOperation.InputRef = awssdkmodel.ShapeRef{
			ShapeName: inputName,
		}
	}
	if hasPrefix(shape.OutputRef.ShapeName, serviceAlias) {
		outputName, _:= removeShapeNamePrefix(shape.OutputRef.ShapeName)
		newOperation.OutputRef = awssdkmodel.ShapeRef{
			ShapeName: outputName,
		}
	}

	for _, err := range shape.ErrorRefs {
		sn, _:= removeShapeNamePrefix(err.ShapeName)
		newOperation.ErrorRefs = append(newOperation.ErrorRefs, awssdkmodel.ShapeRef{
			ShapeName: sn,
		})
	}

	return newOperation
}

// createApiShape creates a shape of awssdkmodel.Shape type
// from the apiv2 Shape.
func createApiShape(shape Shape) (*awssdkmodel.Shape, error) {
	isException := shape.isException()
	apiShape := &awssdkmodel.Shape{
		Type:       shape.Type,
		Exception:  isException,
		MemberRefs: make(map[string]*awssdkmodel.ShapeRef),
		MemberRef:  awssdkmodel.ShapeRef{},
		KeyRef:     awssdkmodel.ShapeRef{},
		ValueRef:   awssdkmodel.ShapeRef{},
		Required:   []string{},
	}

	if isException {
		code, ok := shape.Traits["smithy.api#httpError"]
		if ok {
			switch code := code.(type) {
			case float64:
				apiShape.ErrorInfo = awssdkmodel.ErrorInfo{
					HTTPStatusCode: int(code),
				}
			case int:
				apiShape.ErrorInfo = awssdkmodel.ErrorInfo{
					HTTPStatusCode: code,
				}
			case int64:
				apiShape.ErrorInfo = awssdkmodel.ErrorInfo{
					HTTPStatusCode: int(code),
				}
			default:
				return nil, fmt.Errorf("status code type not found for exception")
			}
		}
	}

	// Shapes that are default so far are booleans and
	// integers that are non pointers in the service API
	val, ok := shape.Traits["smithy.api#default"]
	if ok {
		apiShape.DefaultValue = fmt.Sprintf("%v", val)
	}
	documentation, _ := shape.Traits["smithy.api#documentation"].(string)
	apiShape.Documentation = awssdkmodel.AppendDocstring("", documentation)

	return apiShape, nil
}

// addMemberRefs adds member fields of structures to the structure shape
func addMemberRefs(apiShape *awssdkmodel.Shape, shape Shape, serviceAlias string) {
	for memberName, member := range shape.MemberRefs {
		// Here we make an assumption that we will not get an error
		shapeNameClean, _ := removeShapeNamePrefix(member.ShapeName)
		documentation, _ := member.Traits["smithy.api#documentation"].(string)
		if member.isRequired() {
			apiShape.Required = append(apiShape.Required, memberName)
		}
		apiShape.MemberRefs[memberName] = &awssdkmodel.ShapeRef{
			ShapeName:     shapeNameClean,
			Documentation: awssdkmodel.AppendDocstring("", documentation),
		}
		val, ok := member.Traits["smithy.api#default"]
		if ok {
			apiShape.MemberRefs[memberName].DefaultValue = fmt.Sprintf("%v", val)
		}
	}
	slices.Sort(apiShape.Required)
}

func addMemberRef(apiShape *awssdkmodel.Shape, shape Shape) {

	sn, _:= removeShapeNamePrefix(shape.MemberRef.ShapeName)
	apiShape.MemberRef = awssdkmodel.ShapeRef{
		ShapeName: sn,
	}
}

func addKeyAndValueRef(apiShape *awssdkmodel.Shape, shape Shape) {

	sn, _:= removeShapeNamePrefix(shape.KeyRef.ShapeName)
	apiShape.KeyRef = awssdkmodel.ShapeRef{
		ShapeName: sn,
	}
	sn, _= removeShapeNamePrefix(shape.ValueRef.ShapeName)
	apiShape.ValueRef = awssdkmodel.ShapeRef{
		ShapeName: sn,
	}
}

func addEnumRef(apiShape *awssdkmodel.Shape, shape Shape) {
	for memberName := range shape.MemberRefs {
		apiShape.Enum = append(apiShape.Enum, memberName)
	}
	slices.Sort(apiShape.Enum)
}

func (s ShapeRef) isRequired() bool {
	_, ok := s.Traits["smithy.api#required"]
	return ok
}

func (s Shape) isException() bool {
	_, ok := s.Traits["smithy.api#error"]
	return ok
}

// hasPrefix ensures a shapeRef has a corresponding shape
// (since all shapes start with the prefix from belos)
func hasPrefix(name, alias string) bool {
	prefix := fmt.Sprintf("com.amazonaws.%s#", alias)
	return strings.HasPrefix(name, prefix)
}

// removeShapeNamePrefix removes the prefix from the shapeName.
// The prefix format of a shape in v2 is com.amazonaws.<serviceAlias>#shapeName
func removeShapeNamePrefix(name string) (string, error) {
	temp := strings.Split(name, "#")
	if len(temp) != 2 {
		return "", fmt.Errorf("%s shape name is not formatted correctly, expected format: <url>:<shapeName>", name)
	}
	newName := temp[1]

	return newName, nil
}

// extractServiceAlias extracts the service alias from a shapeName
// (see removeShapeNamePrefix)
func extractServiceAlias(name string) (string) {
	temp := strings.Split(name, ".")
	anotherTemp := strings.Split(temp[len(temp)-1], "#")
	if len(anotherTemp) != 2 {
		return ""
	}
	alias := anotherTemp[0]
	return alias
}

// injectCustomPrimitiveShapes injects custom shapes of primitive
// types in the api. Almost all the sdk models ensure that we have
// shapes for these primitive types. If that's the case, these shapes
// will be replaced with the ones in the model.
//
// For example:
// amp(prometheus service) 
// has target: "smithy.api#String", but no corresponding shape
// named String.
func injectCustomPrimitiveShapes() map[string]*awssdkmodel.Shape {
	shapes := map[string]*awssdkmodel.Shape{}

	types := []string{
		"Blob",
		"Boolean",
		"Document",
		"Double",
		"Float",
		"Integer",
		"Long",
		"PrimitiveBoolean",
		"PrimitiveLong",
		"String",
		"Timestamp",
	}

	for _, t := range types {
		name := names.New(t)
		shapes[t] = &awssdkmodel.Shape{
			Type: name.CamelLower,
		}
		switch t {
		case "PrimitiveBoolean":
			shapes[t].DefaultValue = "false"
		case "PrimitiveLong":
			shapes[t].DefaultValue = "0"
		}
	}
	return shapes
}

// addEnumValues adds the enum values to the shape.
// Some shapes in the model are not of type enum, and instead are strings
func addEnumValues(shape *awssdkmodel.Shape, val interface{}) {

	enumArr := val.([]interface{})
	for _, enumHolder := range enumArr {
		enum, ok := enumHolder.(map[string]interface{})["value"]
		if ok {
			shape.Enum = append(shape.Enum, enum.(string))
		}
	}
}