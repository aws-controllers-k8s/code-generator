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
				"ackerr.ResourceReferenceAndIDNotSupportedFor(%q, %q)\n",
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
					"ackerr.ResourceReferenceOrIDRequiredFor(%q, %q)\n",
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

// ResolveReferencesForField produces Go code for accessing all references that
// are related to the given concrete field, determining whether its in a valid
// condition and updating the concrete field with the referenced value.
// Sample code (resolving a nested singular reference):
//
// ```
//
//	if ko.Spec.JWTConfiguration != nil {
//		if ko.Spec.JWTConfiguration.IssuerRef != nil && ko.Spec.JWTConfiguration.IssuerRef.From != nil {
//			arr := ko.Spec.JWTConfiguration.IssuerRef.From
//			if arr == nil || arr.Name == nil || *arr.Name == "" {
//				return fmt.Errorf("provided resource reference is nil or empty: \"JWTConfiguration.IssuerRef"\")
//			}
//			obj := &svcapitypes.API{}
//			if err := getReferencedResourceState_API(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
//				return err
//			}
//			ko.Spec.JWTConfiguration.Issuer = obj.Status.APIID
//		}
//	}
//
// ```
func ResolveReferencesForField(field *model.Field, sourceVarName string, indentLevel int) string {
	r := field.CRD
	fp := fieldpath.FromString(field.Path)

	outPrefix := ""
	outSuffix := ""

	fieldAccessPrefix := fmt.Sprintf("%s%s", sourceVarName, r.Config().PrefixConfig.SpecField)
	targetVarName := fmt.Sprintf("%s.%s", fieldAccessPrefix, field.Path)

	for idx := 0; idx < fp.Size(); idx++ {
		curFP := fp.CopyAt(idx).String()
		cur, ok := r.Fields[curFP]
		if !ok {
			panic(fmt.Sprintf("unable to find field with path %q. crd: %q", curFP, r.Kind))
		}

		ref := cur.ShapeRef

		indent := strings.Repeat("\t", indentLevel+idx)

		switch ref.Shape.Type {
		case ("structure"):
			fieldAccessPrefix = fmt.Sprintf("%s.%s", fieldAccessPrefix, fp.At(idx))

			outPrefix += fmt.Sprintf("%sif %s != nil {\n", indent, fieldAccessPrefix)
			outSuffix = fmt.Sprintf("%s}\n%s", indent, outSuffix)
		case ("list"):
			if (fp.Size() - idx) > 1 {
				// TODO(nithomso): add support for structs nested within lists
				// The logic for structs nested within lists needs to not only
				// be added here, but also in a custom patching solution since
				// it isn't supported by `StrategicMergePatch`
				// see https://github.com/aws-controllers-k8s/community/issues/1291
				panic(fmt.Errorf("references within lists inside lists aren't supported. crd: %q, path: %q", r.Kind, field.Path))
			}
			fieldAccessPrefix = fmt.Sprintf("%s.%s", fieldAccessPrefix, fp.At(idx))

			iterVarName := fmt.Sprintf("iter%d", idx)

			// base case for references in a list
			outPrefix += fmt.Sprintf("%s%s = %s{}\n", indent, targetVarName, field.GoType)
			outPrefix += fmt.Sprintf("%sfor _, %s := range %s {\n", indent, iterVarName, fieldAccessPrefix)

			fieldAccessPrefix = iterVarName
			outPrefix += fmt.Sprintf("%s\tarr := %s.From\n", indent, fieldAccessPrefix)
			outPrefix += fmt.Sprintf("%s\tif arr == nil || arr.Name == nil || *arr.Name == \"\" {\n", indent)
			outPrefix += fmt.Sprintf("%s\t\treturn fmt.Errorf(\"provided resource reference is nil or empty: \\%q\\\")\n", indent, field.ReferenceFieldPath())
			outPrefix += fmt.Sprintf("%s\t}\n", indent)

			outPrefix += fmt.Sprintf("%s\tif err := getReferencedResourceState_%s(ctx, apiReader, obj, *arr.Name, namespace); err != nil {\n", indent, field.FieldConfig.References.Resource)
			outPrefix += fmt.Sprintf("%s\t\treturn err\n", indent)
			outPrefix += fmt.Sprintf("%s\t}\n", indent)
			outPrefix += fmt.Sprintf("%s\t%s = append(%s, obj.%s)\n", indent, targetVarName, targetVarName, field.FieldConfig.References.Path)
			outPrefix += fmt.Sprintf("%s}\n", indent)
		case ("map"):
			panic("references cannot be within a map")
		default:
			// base case for single references
			fieldAccessPrefix = fmt.Sprintf("%s.%s", fieldAccessPrefix, cur.GetReferenceFieldName().Camel)

			outPrefix += fmt.Sprintf("%sif %s != nil && %s.From != nil {\n", indent, fieldAccessPrefix, fieldAccessPrefix)
			outPrefix += fmt.Sprintf("%s\tarr := %s.From\n", indent, fieldAccessPrefix)
			outPrefix += fmt.Sprintf("%s\tif arr == nil || arr.Name == nil || *arr.Name == \"\" {\n", indent)
			outPrefix += fmt.Sprintf("%s\t\treturn fmt.Errorf(\"provided resource reference is nil or empty: \\%q\\\")\n", indent, field.ReferenceFieldPath())
			outPrefix += fmt.Sprintf("%s\t}\n", indent)

			if field.FieldConfig.References.ServiceName == "" {
				outPrefix += fmt.Sprintf("%s\tobj := &svcapitypes.%s{}\n", indent, field.FieldConfig.References.Resource)
			} else {
				outPrefix += fmt.Sprintf("%s\tobj := &%sapitypes.%s{}\n", indent, field.ReferencedServiceName(), field.FieldConfig.References.Resource)
			}
			outPrefix += fmt.Sprintf("%s\tif err := getReferencedResourceState_%s(ctx, apiReader, obj, *arr.Name, namespace); err != nil {\n", indent, field.FieldConfig.References.Resource)
			outPrefix += fmt.Sprintf("%s\t\treturn err\n", indent)
			outPrefix += fmt.Sprintf("%s\t}\n", indent)
			outPrefix += fmt.Sprintf("%s\t%s = obj.%s\n", indent, targetVarName, field.FieldConfig.References.Path)
			outPrefix += fmt.Sprintf("%s}\n", indent)
		}
	}

	return outPrefix + outSuffix
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
