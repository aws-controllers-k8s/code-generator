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
	"strings"

	"github.com/aws-controllers-k8s/code-generator/pkg/fieldpath"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// ReferenceFieldsValidation produces the go code to validate reference field and
// corresponding identifier field. Iterates through all references within
// slices, if necessary.
//
//	for _, iter0 := range ko.Spec.Routes {
//	  if iter0.GatewayRef != nil && iter0.GatewayID != nil {
//	    return ackerr.ResourceReferenceAndIDNotSupportedFor("Routes.GatewayID", "Routes.GatewayRef")
//	  }
//	}
//
// Sample code:
//
//	if ko.Spec.APIRef != nil && ko.Spec.APIID != nil {
//	  return ackerr.ResourceReferenceAndIDNotSupportedFor("APIID", "APIRef")
//	}
//
//	if ko.Spec.APIRef == nil && ko.Spec.APIID == nil {
//	  return ackerr.ResourceReferenceOrIDRequiredFor("APIID", "APIRef")
//	}
func ReferenceFieldsValidation(
	crd *model.CRD,
	sourceVarName string,
	indentLevel int,
) string {
	out := ""
	// Sorted fieldnames are used for consistent code-generation
	for _, fieldName := range crd.SortedFieldNames() {
		field := crd.Fields[fieldName]
		var fIndent string
		if field.HasReference() {
			fIndentLevel := indentLevel
			fp := fieldpath.FromString(field.Path)
			// remove fieldName from fieldPath before adding nil checks
			fp.Pop()

			// prefix of the field path for referencing in the model
			fieldNamePrefix := ""
			// prefix of the field path for the generated code
			pathVarPrefix := fmt.Sprintf("%s%s", sourceVarName, crd.Config().PrefixConfig.SpecField)

			// this loop outputs a nil-guard for each level of nested field path
			// or an iterator for any level that is a slice
			fieldDepth := 0
			for fp.Size() > 0 {
				fIndent = strings.Repeat("\t", fIndentLevel)
				currentField := fp.PopFront()

				if fieldNamePrefix == "" {
					fieldNamePrefix = currentField
				} else {
					fieldNamePrefix = fmt.Sprintf("%s.%s", fieldNamePrefix, currentField)
				}
				pathVarPrefix = fmt.Sprintf("%s.%s", pathVarPrefix, currentField)

				fieldConfig, ok := crd.Fields[fieldNamePrefix]
				if !ok {
					panic(fmt.Sprintf("CRD %s has no Field with path %s", crd.Kind, fieldNamePrefix))
				}

				if fieldConfig.ShapeRef.Shape.Type == "list" {
					out += fmt.Sprintf("%sfor _, iter%d := range %s {\n", fIndent, fieldDepth, pathVarPrefix)
					// reset the path variable name
					pathVarPrefix = fmt.Sprintf("iter%d", fieldDepth)
				} else {
					out += fmt.Sprintf("%sif %s != nil {\n", fIndent, pathVarPrefix)
				}

				fIndentLevel++
				fieldDepth++
			}

			fIndent = strings.Repeat("\t", fIndentLevel)
			// Validation to make sure both target field and reference are
			// not present at the same time in desired resource
			out += fmt.Sprintf("%sif %s.%s != nil"+
				" && %s.%s != nil {\n", fIndent, pathVarPrefix, field.GetReferenceFieldName().Camel, pathVarPrefix, field.Names.Camel)
			out += fmt.Sprintf("%s\treturn "+
				"ackerr.ResourceReferenceAndIDNotSupportedFor(\"%s\", \"%s\")\n",
				fIndent, field.Path, field.ReferenceFieldPath())

			// Close out all the curly braces with proper indentation
			for fIndentLevel >= indentLevel {
				fIndent = strings.Repeat("\t", fIndentLevel)
				out += fmt.Sprintf("%s}\n", fIndent)
				fIndentLevel--
			}

			fIndent = strings.Repeat("\t", indentLevel)

			// If the field is required, make sure either Ref or original
			// field is present in the resource
			if field.IsRequired() {
				out += fmt.Sprintf("%sif %s.%s == nil &&"+
					" %s.%s == nil {\n", fIndent, pathVarPrefix,
					field.ReferenceFieldPath(), pathVarPrefix, field.Path)
				out += fmt.Sprintf("%s\treturn "+
					"ackerr.ResourceReferenceOrIDRequiredFor(\"%s\", \"%s\")\n",
					fIndent, field.Path, field.ReferenceFieldPath())
				out += fmt.Sprintf("%s}\n", fIndent)
			}
		}
	}
	return out
}

// ReferenceFieldsPresent produces go code(logical condition) for finding whether
// a non-nil reference field is present in a resource. This checks helps in deciding
// whether ACK.ReferencesResolved condition should be added to resource status
// Sample Code:
//
//	if ko.Spec.Routes != nil {
//	  for _, iter35 := range ko.Spec.Routes {
//	    if iter35.GatewayRef != nil {
//	      return true
//	    }
//	  }
//	}
//
// return false || (ko.Spec.APIRef != nil)
func ReferenceFieldsPresent(
	crd *model.CRD,
	sourceVarName string,
) string {
	iteratorsOut := ""
	returnOut := "return false"
	fieldAccessPrefix := fmt.Sprintf("%s%s", sourceVarName,
		crd.Config().PrefixConfig.SpecField)
	// Sorted fieldnames are used for consistent code-generation
	for fieldIndex, fieldName := range crd.SortedFieldNames() {
		field := crd.Fields[fieldName]
		if field.HasReference() {
			fp := fieldpath.FromString(field.Path)
			// remove fieldName from fieldPath before adding nil checks
			// for nested fieldPath
			fp.Pop()

			// Determine whether the field is nested
			if fp.Size() > 0 {
				// Determine whether the field is inside a slice
				parentField, ok := crd.Fields[fp.String()]
				if !ok {
					panic(fmt.Sprintf("CRD %s has no Field with path %s", crd.Kind, fp.String()))
				}

				if parentField.ShapeRef.Shape.Type == "list" {
					iteratorsOut += fmt.Sprintf("if %s {\n", nestedStructNilCheck(*fp.Copy(), fieldAccessPrefix))
					iteratorsOut += fmt.Sprintf("\tfor _, iter%d := range %s.%s {\n", fieldIndex, fieldAccessPrefix, parentField.Path)
					iteratorsOut += fmt.Sprintf("\t\tif iter%d.%s != nil {\n", fieldIndex, field.GetReferenceFieldName().Camel)
					iteratorsOut += fmt.Sprintf("\t\t\treturn true\n")
					iteratorsOut += fmt.Sprintf("\t\t}\n")
					iteratorsOut += fmt.Sprintf("\t}\n")
					iteratorsOut += fmt.Sprintf("}\n")
					continue
				}
			}

			nilCheck := nestedStructNilCheck(*fp.Copy(), fieldAccessPrefix) + " && " + fmt.Sprintf("%s.%s != nil", fieldAccessPrefix,
				field.ReferenceFieldPath())
			returnOut += " || (" + strings.TrimPrefix(nilCheck, " && ") + ")"
		}
	}
	return iteratorsOut + returnOut
}

func nestedStructNilCheck(path fieldpath.Path, fieldAccessPrefix string) string {
	out := ""
	fieldNamePrefix := ""
	for path.Size() > 0 {
		fieldNamePrefix = fmt.Sprintf("%s.%s", fieldNamePrefix, path.PopFront())
		out += fmt.Sprintf("%s%s != nil && ", fieldAccessPrefix, fieldNamePrefix)
	}
	return strings.TrimSuffix(out, " && ")
}
