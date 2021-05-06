package code

import (
	"fmt"
	"sort"
	"strings"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/names"
)

// CompareResourceHelpers returns Go code that define the helper functions used by
// newResourceDelta function. The helper functions return true if two given parameters
// are equal,false otherwise.
func CompareResourceHelpers(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := "\n"

	specFieldNames := []string{}
	for fieldName := range r.SpecFields {
		specFieldNames = append(specFieldNames, fieldName)
	}
	sort.Strings(specFieldNames)

	for _, fieldName := range specFieldNames {
		specField := r.SpecFields[fieldName]
		memberShapeRef := specField.ShapeRef
		memberShape := memberShapeRef.Shape

		// A map storing visited shapes to avoid collecting duplicate shapes
		visitedShapes := make(map[string]*struct{})
		shapes := collectSubShapes(memberShape, visitedShapes)

		for _, shape := range shapes {
			shapeGoType := shape.GoTypeWithPkgName()
			goType := model.ReplacePkgName(shapeGoType, r.SDKAPIPackageName(), "svcapitypes", true)
			firstFuncParam := "a"
			secondFuncParam := "b"

			// setting helper function name and signature
			// func equalRouttingSettingsMap(a, b map[string]*svcapitypes.RouteSettings) bool {
			out += fmt.Sprintf(
				"func equal%s(%s, %s %s) bool {\n",
				shape.ShapeName,
				firstFuncParam,
				secondFuncParam,
				goType,
			)
			switch shape.Type {
			case "map":
				// Returns Go code that compares two maps
				out += equalMap(
					shape,
					firstFuncParam,
					secondFuncParam,
					1,
				)
			case "list":
				// Returns Go code that compares two slices
				out += equalSlice(
					shape,
					firstFuncParam,
					secondFuncParam,
					1,
				)
			default:
				msg := "Should not generate a helper function for any type different" +
					"than string to struct/list/map map and struct/list/map slices"
				panic(msg)
			}

			// return true
			out += fmt.Sprintf("\treturn true\n")
			// }
			out += "}\n"
		}
	}
	return out
}

// expectHelperFunction returns true if a given shape needs a helper function, false
// otherwise.
func expectHelperFunction(
	shape *awssdkmodel.Shape,
) bool {
	if shape.Type == "map" {
		valueType := shape.ValueRef.Shape.Type
		return valueType == "structure" || valueType == "list" || valueType == "map"
	}

	if shape.Type == "list" {
		memberType := shape.MemberRef.Shape.Type
		return memberType == "structure" || memberType == "list" || memberType == "map"
	}
	return false
}

// collectSubShapes takes a shape and returns all the nested member and value
// shapes that needs a helper function.
func collectSubShapes(shape *awssdkmodel.Shape, visitedShapes map[string]*struct{}) []*awssdkmodel.Shape {
	// exit if shape was already visited
	if _, ok := visitedShapes[shape.ShapeName]; ok {
		return nil
	}
	// mark shape as visited
	visitedShapes[shape.ShapeName] = nil

	shapes := make([]*awssdkmodel.Shape, 0, 0)
	if expectHelperFunction(shape) {
		shapes = append(shapes, shape)
	}

	// set target shape
	var targetShape *awssdkmodel.Shape
	if shape.Type == "map" {
		targetShape = shape.ValueRef.Shape
	} else if shape.Type == "list" {
		targetShape = shape.MemberRef.Shape
	}

	switch shape.Type {
	case "map", "list":
		switch targetShape.Type {
		case "map":
			valueShape := targetShape.ValueRef.Shape
			shapes = append(shapes, collectSubShapes(valueShape, visitedShapes)...)
		case "list":
			memberShape := targetShape.MemberRef.Shape
			shapes = append(shapes, collectSubShapes(memberShape, visitedShapes)...)
		case "structure":
			for _, memberName := range targetShape.MemberNames() {
				memberShape := targetShape.MemberRefs[memberName].Shape
				shapes = append(shapes, collectSubShapes(memberShape, visitedShapes)...)
			}
		}
	case "structure":
		for _, memberName := range shape.MemberNames() {
			memberShapeRef := shape.MemberRefs[memberName]
			memberShape := memberShapeRef.Shape
			shapes = append(shapes, collectSubShapes(memberShape, visitedShapes)...)
		}
	}

	return shapes
}

// equalMap outputs Go code that compares two Go maps of the the same value type.
// at the first spotted difference the code return false, return true if the maps
// key/values are equal.
func equalMap(
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that represents the object
	// under comparison.
	firstVarName string,
	// String representing the name of the variable that represents the second
	// object under comparison.
	secondVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	memberValueShape := shape.ValueRef.Shape

	// 	if len(a) != len(b) {
	out += fmt.Sprintf(
		"%sif len(%s) != len(%s) {\n",
		indent,
		firstVarName,
		secondVarName,
	)
	// return false
	out += fmt.Sprintf(
		"%s\treturn false\n", indent,
	)
	// }
	out += fmt.Sprintf("%s}\n", indent)

	// for ka := range a {
	out += fmt.Sprintf(
		"%sfor ka := range %s {\n",
		indent,
		firstVarName,
	)
	// _, ok := b[ka]
	out += fmt.Sprintf(
		"%s\t_, ok := %s[ka]\n",
		indent,
		secondVarName,
	)
	// if !ok {
	out += fmt.Sprintf("%s\tif !ok {\n", indent)
	// return false
	out += fmt.Sprintf(
		"%s\t\treturn false\n", indent,
	)
	// }
	out += fmt.Sprintf("%s\t}\n", indent)
	// closing the for loop
	out += fmt.Sprintf("%s}\n", indent)

	firstVarAdaptedName := firstVarName + "X"
	secondVarAdaptedName := secondVarName + "Y"
	// for ka, kv := range b {
	out += fmt.Sprintf("%sfor ka, %s := range a {\n", indent, firstVarAdaptedName)
	// vb := b[ka]
	out += fmt.Sprintf("%s\t%s := b[ka]\n", indent, secondVarAdaptedName)

	switch memberValueShape.Type {
	case "string":
		// if !ackcompare.MapStringStringPEqual(a.ko.Spec.Tags, b.ko.Spec.Tags) {
		out += fmt.Sprintf(
			"%sif !ackcompare.MapStringStringPEqual(%s, %s) {\n",
			indent, firstVarAdaptedName, secondVarAdaptedName,
		)
		out += fmt.Sprintf(
			"%s\n", indent,
		)
	case "structure":
		out += equalStruct(
			memberValueShape,
			firstVarAdaptedName,
			secondVarAdaptedName,
			indentLevel+1,
		)
	case "list", "blob":
		out += equalSlice(
			memberValueShape,
			firstVarAdaptedName,
			secondVarAdaptedName,
			indentLevel+1,
		)
	case "map":
		out += equalMap(
			memberValueShape,
			firstVarAdaptedName,
			secondVarAdaptedName,
			indentLevel+1,
		)
	default:
		panic("Unsupported shape type in generate.code.equalMap: " + shape.Type)
	}

	out += fmt.Sprintf("%s}\n", indent)
	return out
}

// equalNil outputs Go code that compares pointer to struct types for nullability,
// if there is a nil difference the code return false.
//
// Output code will look something like this:
//
// if ackcompare.HasNilDifference(va.DataTraceEnabled, vb.DataTraceEnabled) {
//     return false
// }
func equalNil(
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that represents the object
	// under comparison.
	firstVarName string,
	// String representing the name of the variable that represents the second
	// object under comparison.
	secondVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	switch shape.Type {
	case "list", "blob":
		// for slice types, there is no nilability test. Instead, the normal
		// value test checks length of slices.
		return ""
	case "boolean", "string", "character", "byte", "short", "integer", "long",
		"float", "double", "timestamp", "structure", "map", "jsonvalue":
		// if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name) {
		out += fmt.Sprintf(
			"%sif ackcompare.HasNilDifference(%s, %s) {\n",
			indent, firstVarName, secondVarName,
		)
	default:
		panic("Unsupported shape type in generate.code.compareNil: " + shape.Type)
	}
	// return false
	out += fmt.Sprintf(
		"%s\treturn false\n",
		indent,
	)
	// }
	out += fmt.Sprintf(
		"%s}", indent,
	)
	return out
}

// equalStruct outputs Go code that compares two struct values from two
// resource fields and, if there is a difference code will return false.
func equalStruct(
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that represents the object
	// under comparison.
	firstVarName string,
	// String representing the name of the variable that represents the second
	// object under comparison.
	secondVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	//indent := strings.Repeat("\t", indentLevel)

	for _, memberName := range shape.MemberNames() {
		memberShapeRef := shape.MemberRefs[memberName]
		memberNames := names.New(memberName)
		memberNameClean := memberNames.Camel
		firstAdaptedVarName := firstVarName + "." + memberNameClean
		secondAdaptedVarName := secondVarName + "." + memberNameClean

		nilCode := equalNil(
			shape,
			firstAdaptedVarName,
			secondAdaptedVarName,
			indentLevel,
		)
		if nilCode != "" {
			out += fmt.Sprintf(
				"%s else if %s != nil && %s != nil {\n",
				nilCode, firstAdaptedVarName, secondAdaptedVarName,
			)
			indentLevel++
		} else {
			out += "\n"
		}

		memberShape := memberShapeRef.Shape
		switch memberShape.Type {
		case "structure":
			out += equalStruct(
				memberShape,
				firstAdaptedVarName,
				secondAdaptedVarName,
				indentLevel,
			)
		case "list":
			out += equalSlice(
				memberShape,
				firstAdaptedVarName,
				secondAdaptedVarName,
				indentLevel,
			)
		case "map":
			out += equalMap(
				memberShape,
				firstAdaptedVarName,
				secondAdaptedVarName,
				indentLevel,
			)
		default:
			indent := strings.Repeat("\t", indentLevel)
			out += fmt.Sprintf(
				"%sif *%s != *%s {\n",
				indent,
				firstAdaptedVarName,
				secondAdaptedVarName,
			)
			out += fmt.Sprintf("%s\treturn false\n", indent)
			out += fmt.Sprintf("%s}\n", indent)
		}

		if nilCode != "" {
			indentLevel--
			indent := strings.Repeat("\t", indentLevel)
			out += fmt.Sprintf("%s}\n", indent)
		}

	}
	return out
}

// equalSlice outputs Go code that compares two Go slices of the the same value type.
// at the first spotted difference the code return false, return true if the element
// are equal.
// TODO(hilalymh): Modify this function to be configurable: Ordered/Non-Ordered/ExactCountRepeatedElements
func equalSlice(
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that represents the object
	// under comparison.
	firstVarName string,
	// String representing the name of the variable that represents the second
	// object under comparison.
	secondVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	//TODO(a-hilaly)
	return "\t//TODO(a-hilaly) implement this function\n"
}
