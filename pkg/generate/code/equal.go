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

	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/names"
)

// IsEqualTypeDef returns Go code that checks the equality of two struct types.
// The generated code returns false at the first spotted difference.
//
// Output code will look something like this:
//
//  if ackcompare.HasNilDifference(a.CreatedAt, b.CreatedAt) {
//  	return false
//  } else if a.CreatedAt != nil && b.CreatedAt != nil {
//  	if *a.CreatedAt != *b.CreatedAt {
//  		return false
//  	}
//  }
//  if ackcompare.HasNilDifference(a.EncryptionConfiguration, b.EncryptionConfiguration) {
//  	return false
//  } else if a.EncryptionConfiguration != nil && b.EncryptionConfiguration != nil {
//  	if !IsEqualEncryptionConfiguration(a.EncryptionConfiguration, b.EncryptionConfiguration) {
//  		return false
//  	}
//  }
func IsEqualTypeDef(
	typeDef *ackmodel.TypeDef,
	// String representing the name of the variable that represents the first
	// object under comparison. This will typically be something like "a"
	firstVarName string,
	// String representing the name of the variable that represents the second
	// object under comparison. This will typically be something like "b".
	secondVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""

	// We need a deterministic order to loop over attributes
	attrsNames := []string{}
	for attrName := range typeDef.Attrs {
		attrsNames = append(attrsNames, attrName)
	}
	sort.Strings(attrsNames)

	// For each attribute generate the equality check logic
	for _, attrName := range attrsNames {
		attr := typeDef.Attrs[attrName]
		out += isEqualField(attrName, attr.Shape, firstVarName, secondVarName, indentLevel)
	}
	return out
}

// isEqualField outputs Go code that compares two similar fields from different objects.
// The generated code first checks the nility of the fields and then procede on comparing
// their value.
//
// Output code will look something like this:
//
//	if ackcompare.HasNilDifference(a.RepositoryName, b.RepositoryName) {
//		return false
//	} else if a.RepositoryName != nil && b.RepositoryName != nil {
//		if *a.RepositoryName != *b.RepositoryName {
//			return false
//		}
//	}
func isEqualField(
	// fieldName is the field to generate the comparison logic for
	fieldName string,
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

	memberNames := names.New(fieldName)
	memberNameClean := memberNames.Camel
	firstAdaptedVarName := firstVarName + "." + memberNameClean
	secondAdaptedVarName := secondVarName + "." + memberNameClean

	nilCode := isEqualNil(
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

	indent := strings.Repeat("\t", indentLevel)

	switch shape.Type {
	case "structure":
		// We just re-use the generated functions to compare field of `struct` type
		out += fmt.Sprintf(
			"%sif !IsEqual%s(%s, %s) {\n",
			indent,
			shape.ShapeName,
			firstAdaptedVarName,
			secondAdaptedVarName,
		)
		out += fmt.Sprintf("%s\treturn false\n", indent)
		out += fmt.Sprintf("%s}\n", indent)

	case "list":
		out += isEqualSlice(
			shape,
			firstAdaptedVarName,
			secondAdaptedVarName,
			indentLevel,
		)
	case "map":
		out += isEqualMap(
			shape,
			firstAdaptedVarName,
			secondAdaptedVarName,
			indentLevel,
		)
	default:
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
	return out
}

// isEqualNil outputs Go code that compares pointer to struct types for nullability,
// if there is a nil difference the code return false.
//
// Output code will look something like this:
//
// if ackcompare.HasNilDifference(a.DataTraceEnabled, b.DataTraceEnabled) {
//     return false
// }
func isEqualNil(
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

// isEqualSlice outputs Go code that compares two Go slices of the the same value type.
// at the first spotted difference the code return false, return true if the element
// are equal.
// TODO(hilalymh): Modify this function to be configurable: Ordered/Non-Ordered/ExactCountRepeatedElements
func isEqualSlice(
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
	indent := strings.Repeat("\t", indentLevel)
	out := fmt.Sprintf(
		"%s//TODO(a-hilaly): equality check for slices\n", indent,
	)
	return out
}

// isEqualMap outputs Go code that compares two Go maps of the the same value type.
// at the first spotted difference the code return false, return true if the maps
// key/values are equal.
func isEqualMap(
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
	indent := strings.Repeat("\t", indentLevel)
	out := fmt.Sprintf(
		"%s//TODO(a-hilaly): equality check for maps\n", indent,
	)
	return out
}
