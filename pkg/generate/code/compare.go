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

// CompareResource returns the Go code that traverses a set of two Resources,
// adding differences between the two Resources to an `ackcompare.Delta`
//
// By default, we produce Go code that only looks at the fields in a resource's
// Spec, since those are the fields that represent the desired state of a
// resource. When we make a ReadOne/ReadMany/GetAttributes call to a backend
// AWS API, we construct a Resource and set the Spec fields to values contained
// in the ReadOne/ReadMany/GetAttributes Output shape. This Resource,
// constructed from the Read operation, is compared to the Resource we got from
// the Kubernetes API server's event bus. The code that is returned from this
// function is the code that compares those two Resources.
func CompareResource(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like "a.ko". See
	// `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like "b.ko". See
	// `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := "\n"

	fieldConfigs := cfg.ResourceFields(r.Names.Original)

	// We need a deterministic order to traverse our top-level fields...
	specFieldNames := []string{}
	for fieldName := range r.SpecFields {
		specFieldNames = append(specFieldNames, fieldName)
	}
	sort.Strings(specFieldNames)

	for _, fieldName := range specFieldNames {
		specField := r.SpecFields[fieldName]
		indent := strings.Repeat("\t", indentLevel)
		firstResAdaptedVarName := firstResVarName + cfg.PrefixConfig.SpecField
		firstResAdaptedVarName += "." + specField.Names.Camel
		secondResAdaptedVarName := secondResVarName + cfg.PrefixConfig.SpecField
		secondResAdaptedVarName += "." + specField.Names.Camel

		var compareConfig *ackgenconfig.CompareFieldConfig
		fieldConfig := fieldConfigs[fieldName]
		if fieldConfig != nil {
			compareConfig = fieldConfig.Compare
		}

		if compareConfig != nil && compareConfig.IsIgnored {
			continue
		}

		// this is the "path" to the field within the structs being compared.
		// This is passed down into the compareXXX functions recursively and
		// appended to with each level of nested structs we recurse into.
		fieldPath := strings.TrimPrefix(
			cfg.PrefixConfig.SpecField+"."+specField.Names.Camel, ".",
		)

		memberShapeRef := specField.ShapeRef
		memberShape := memberShapeRef.Shape

		// if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name == nil) {
		//   delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
		// }
		nilCode := compareNil(
			compareConfig,
			memberShape,
			deltaVarName,
			firstResAdaptedVarName,
			secondResAdaptedVarName,
			fieldPath,
			indentLevel,
		)

		if nilCode != "" {
			// else {
			out += nilCode + " else {\n"
			indentLevel++
		} else {
			out += "\n"
		}

		switch memberShape.Type {
		case "structure":
			// Recurse through all the struct's fields and subfields, building
			// nested conditionals and calls to `delta.Add()`...
			out += compareStruct(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				fieldPath,
				indentLevel,
			)
		case "list":
			// Returns Go code that compares all the elements of the slice fields...
			out += compareSlice(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				fieldPath,
				indentLevel,
			)
		case "map":
			// Returns Go code that compares all the elements of the map fields...
			out += compareMap(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				fieldPath,
				indentLevel,
			)
		default:
			//   if *a.ko.Spec.Name != *b.ko.Spec.Name) {
			//     delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
			//   }
			out += compareScalar(
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				fieldPath,
				indentLevel,
			)
		}
		// }
		out += fmt.Sprintf(
			"%s}\n", indent,
		)
		indentLevel--
	}
	return out
}

// compareNil outputs Go code that compares two field values for nullability
// and, if there is a nil difference, adds the difference to a variable
// represeting the `ackcompare.Delta`
//
// Output code will look something like this:
//
// if ackcompare.HasNilDifferenceStringP(a.ko.Spec.Name, b.ko.Spec.Name == nil) {
//   delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
// }
func compareNil(
	// struct informing code generator how to compare the field values
	compareConfig *ackgenconfig.CompareFieldConfig,
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
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
			indent, firstResVarName, secondResVarName,
		)
	default:
		panic("Unsupported shape type in generate.code.compareNil: " + shape.Type)
	}
	//   delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	out += fmt.Sprintf(
		"%s\t%s.Add(\"%s\", %s, %s)\n",
		indent, deltaVarName, fieldPath, firstResVarName, secondResVarName,
	)
	// }
	out += fmt.Sprintf(
		"%s}", indent,
	)

	return out
}

// compareScalar outputs Go code that compares two scalar values from two
// resource fields and, if there is a difference, adds the difference to a
// variable representing an `ackcompare.Delta`.
//
// Output code will look something like this:
//
//   if *a.ko.Spec.Name != *b.ko.Spec.Name) {
//     delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
//   }
func compareScalar(
	// struct informing code generator how to compare the field values
	compareConfig *ackgenconfig.CompareFieldConfig,
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	switch shape.Type {
	case "boolean", "string", "character", "byte", "short", "integer", "long", "float", "double":
		// if *a.ko.Spec.Name != *b.ko.Spec.Name {
		out += fmt.Sprintf(
			"%sif *%s != *%s {\n",
			indent, firstResVarName, secondResVarName,
		)
	case "timestamp":
		// if !a.ko.Spec.CreatedAt.Equal(b.ko.Spec.CreatedAt) {
		out += fmt.Sprintf(
			"%sif !%s.Equal(%s) {\n",
			indent, firstResVarName, secondResVarName,
		)
	default:
		panic("Unsupported shape type in generate.code.compareScalar: " + shape.Type)
	}
	//   delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	out += fmt.Sprintf(
		"%s\t%s.Add(\"%s\", %s, %s)\n",
		indent, deltaVarName, fieldPath, firstResVarName, secondResVarName,
	)
	// }
	out += fmt.Sprintf(
		"%s}\n", indent,
	)

	return out
}

// compareMap outputs Go code that compares two map values from two resource
// fields and, if there is a difference, adds the difference to a variable
// representing an `ackcompare.Delta`.
//
// Output code will look something like this:
//
//   if !ackcompare.MapStringStringEqual(a.ko.Spec.Tags, b.ko.Spec.Tags) {
//     delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
//   }
func compareMap(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// struct informing code generator how to compare the field values
	compareConfig *ackgenconfig.CompareFieldConfig,
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	keyType := shape.KeyRef.Shape.Type

	if keyType != "string" {
		panic("generate.code.compareMap cannot deal with non-string key types: " + keyType)
	}

	valType := shape.ValueRef.Shape.Type

	switch valType {
	case "string":
		// if !ackcompare.MapStringStringEqual(a.ko.Spec.Tags, b.ko.Spec.Tags) {
		out += fmt.Sprintf(
			"%sif !ackcompare.MapStringStringEqual(%s, %s) {\n",
			indent, firstResVarName, secondResVarName,
		)
	case "structure":
		// TODO(jaypipes): Implement this by walking the keys and struct values
		// and comparing each struct individually, building up the fieldPath
		// appropriately...
	default:
		panic("Unsupported shape type in generate.code.compareMap: " + shape.Type)
	}
	//   delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	out += fmt.Sprintf(
		"%s\t%s.Add(\"%s\", %s, %s)\n",
		indent, deltaVarName, fieldPath, firstResVarName, secondResVarName,
	)
	// }
	out += fmt.Sprintf(
		"%s}\n", indent,
	)

	return out
}

// compareSlice outputs Go code that compares two slice values from two
// resource fields and, if there is a difference, adds the difference to a
// variable representing an `ackcompare.Delta`.
//
// Output code will look something like this:
//
//   if !ackcompare.SliceStringPEqual(a.ko.Spec.SecurityGroupIDs, b.ko.Spec.SecurityGroupIDs) {
//     delta.Add("Spec.SecurityGroupIDs", a.ko.Spec.SecurityGroupIDs, b.ko.Spec.SecurityGroupIDs)
//   }
func compareSlice(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// struct informing code generator how to compare the field values
	compareConfig *ackgenconfig.CompareFieldConfig,
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	elemType := shape.MemberRef.Shape.Type

	switch elemType {
	case "string":
		// if !ackcompare.SliceStringPEqual(a.ko.Spec.SecurityGroupIDs, b.ko.Spec.SecurityGroupIDs) {
		out += fmt.Sprintf(
			"%sif !ackcompare.SliceStringPEqual(%s, %s) {\n",
			indent, firstResVarName, secondResVarName,
		)
	case "structure":
		// TODO(jaypipes): Implement this by walking the slice of struct values
		// and comparing each struct individually, building up the fieldPath
		// appropriately...
	default:
		panic("Unsupported shape type in generate.code.compareSlice: " + shape.Type)
	}
	//   delta.Add("Spec.SecurityGroupIDs", a.ko.Spec.SecurityGroupIDs, b.ko.Spec.SecurityGroupIDs)
	out += fmt.Sprintf(
		"%s\t%s.Add(\"%s\", %s, %s)\n",
		indent, deltaVarName, fieldPath, firstResVarName, secondResVarName,
	)
	// }
	out += fmt.Sprintf(
		"%s}\n", indent,
	)

	return out
}

// compareStruct outputs Go code that compares two struct values from two
// resource fields and, if there is a difference, adds the difference to a
// variable representing an `ackcompare.Delta`.
func compareStruct(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// struct informing code generator how to compare the field values
	compareConfig *ackgenconfig.CompareFieldConfig,
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""

	fieldConfigs := cfg.ResourceFields(r.Names.Original)

	for _, memberName := range shape.MemberNames() {
		memberShapeRef := shape.MemberRefs[memberName]
		// TODO(jaypipes): This is fragile. We actually need to have a way of
		// normalizing names in a lossless fashion...
		memberNames := names.New(memberName)
		memberNameClean := memberNames.Camel

		// this is the "path" to the field within the structs being compared.
		// This is passed down into the compareXXX functions recursively and
		// appended to with each level of nested structs we recurse into.
		memberFieldPath := fieldPath + "." + memberNameClean
		indent := strings.Repeat("\t", indentLevel)
		firstResAdaptedVarName := firstResVarName + "." + memberNameClean
		secondResAdaptedVarName := secondResVarName + "." + memberNameClean

		var compareConfig *ackgenconfig.CompareFieldConfig
		fieldConfig := fieldConfigs[memberFieldPath]
		if fieldConfig != nil {
			compareConfig = fieldConfig.Compare
		}

		if compareConfig != nil && compareConfig.IsIgnored {
			continue
		}

		memberShape := memberShapeRef.Shape

		// if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name == nil) {
		//   delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
		// }
		nilCode := compareNil(
			compareConfig,
			memberShape,
			deltaVarName,
			firstResAdaptedVarName,
			secondResAdaptedVarName,
			memberFieldPath,
			indentLevel,
		)

		if nilCode != "" {
			// else {
			out += nilCode + " else {\n"
			indentLevel++
		} else {
			out += "\n"
		}
		switch memberShape.Type {
		case "structure":
			// Recurse through all the struct's fields and subfields, building
			// nested conditionals and calls to `delta.Add()`...
			out += compareStruct(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				memberFieldPath,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				indentLevel,
			)
		case "list":
			// Returns Go code that compares all the elements of the slice fields...
			out += compareSlice(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				memberFieldPath,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				indentLevel,
			)
		case "map":
			// Returns Go code that compares all the elements of the map fields...
			out += compareMap(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				memberFieldPath,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				indentLevel,
			)
		default:
			//   if *a.ko.Spec.Name != *b.ko.Spec.Name {
			//     delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
			//   }
			out += compareScalar(
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				memberFieldPath,
				indentLevel,
			)
		}
		// }
		out += fmt.Sprintf(
			"%s}\n", indent,
		)
		indentLevel--
	}
	return out
}
