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

func ResolveReferencesForField(r *model.CRD, field *model.Field, indentLevel int) string {
	out := ""

	fp := fieldpath.FromString(field.Path)
	fp.Pop()

	isNested := fp.Size() > 0
	isList := field.ShapeRef.Shape.Type == "list"

	if !isList && !isNested {
		return resolveSingleReference(r, field, indentLevel)
	} else if !isNested {
		return resolveSliceOfReferences(r, field, indentLevel)
	} else {
		parentField, ok := r.Fields[fp.String()]
		if !ok {
			panic(fmt.Sprintf("unable to find parent field with path %s", fp.String()))
		}

		if parentField.ShapeRef.Shape.Type == "list" {
			return resolveNestedSliceOfReferences(r, field, indentLevel)
		} else {
			return resolveNestedSingleReference(r, field, indentLevel)
		}
	}

	return out
}

func resolveSingleReference(r *model.CRD, field *model.Field, indentLevel int) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	rfp := field.ReferenceFieldPath()

	out += fmt.Sprintf("%sif ko.Spec.%s != nil && ko.Spec.%s.From != nil {\n", indent, rfp, rfp)
	out += fmt.Sprintf("%s\tarr := ko.Spec.%s.From\n", indent, rfp)
	out += readReferenceAndValidate(field, indentLevel+1)
	out += fmt.Sprintf("%s\treferencedValue := string(*obj.%s)\n", indent, field.FieldConfig.References.Path)
	out += fmt.Sprintf("%s\tko.Spec.%s = &referencedValue\n", indent, field.Path)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%sreturn nil", indent)

	return out
}

func resolveSliceOfReferences(r *model.CRD, field *model.Field, indentLevel int) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	rfp := field.ReferenceFieldPath()

	out += fmt.Sprintf("%sif ko.Spec.%s != nil &&\n", indent, rfp)
	out += fmt.Sprintf("%s\tlen(ko.Spec.%s) > 0 {\n", indent, rfp)
	out += fmt.Sprintf("%s\tresolvedReferences := []*string{}\n", indent)
	out += fmt.Sprintf("%s\tfor _, arrw := range ko.Spec.%s {\n", indent, rfp)
	out += fmt.Sprintf("%s\t\tarr := arrw.From\n", indent)
	out += readReferenceAndValidate(field, indentLevel+2)
	out += fmt.Sprintf("%s\t\treferencedValue := string(*obj.%s)\n", indent, field.FieldConfig.References.Path)
	out += fmt.Sprintf("%s\t\tresolvedReferences = append(resolvedReferences, &referencedValue)\n", indent)
	out += fmt.Sprintf("%s\t}\n", indent)
	out += fmt.Sprintf("%s\tko.Spec.%s = resolvedReferences\n", indent, field.Path)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%sreturn nil", indent)

	return out
}

func resolveNestedSingleReference(r *model.CRD, field *model.Field, indentLevel int) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	rfp := field.ReferenceFieldPath()

	out += fmt.Sprintf("%sif ko.Spec.%s != nil &&\n", indent, rfp)
	out += fmt.Sprintf("%s\tlen(ko.Spec.%s) > 0 {\n", indent, rfp)
	out += fmt.Sprintf("%s\tresolvedReferences := []*string{}\n", indent)
	out += fmt.Sprintf("%s\tfor _, arrw := range ko.Spec.%s {\n", indent, rfp)
	out += fmt.Sprintf("%s\t\tarr := arrw.From\n", indent)
	out += readReferenceAndValidate(field, indentLevel+2)
	out += fmt.Sprintf("%s\t\treferencedValue := string(*obj.%s)\n", indent, field.FieldConfig.References.Path)
	out += fmt.Sprintf("%s\t\tresolvedReferences = append(resolvedReferences, &referencedValue)\n", indent)
	out += fmt.Sprintf("%s\t}\n", indent)
	out += fmt.Sprintf("%s\tko.Spec.%s = resolvedReferences\n", indent, field.Path)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%sreturn nil", indent)

	return out
}

func resolveNestedSliceOfReferences(r *model.CRD, field *model.Field, indentLevel int) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	fp := fieldpath.FromString(field.Path)
	fp.Pop()

	parent, ok := r.Fields[fp.String()]
	if !ok {
		panic("")
	}

	out += fmt.Sprintf("%sif len(ko.Spec.%s) > 0 {\n", indent, parent.Path)
	out += fmt.Sprintf("%s\tfor _, elem := range ko.Spec.%s {\n", indent, parent.Path)
	out += fmt.Sprintf("%s\t\tarrw := elem.%s\n", indent, field.GetReferenceFieldName().Camel)
	out += fmt.Sprintf("%s\n", indent)
	out += fmt.Sprintf("%s\t\tif arrw == nil || arrw.From == nil {\n", indent)
	out += fmt.Sprintf("%s\t\t\tcontinue\n", indent)
	out += fmt.Sprintf("%s\t\t}\n", indent)
	out += fmt.Sprintf("%s\n", indent)
	out += fmt.Sprintf("%s\t\tarr := arrw.From\n", indent)
	out += fmt.Sprintf("%s\t\tif arr.Name == nil || *arr.Name == \"\" {\n", indent)
	out += fmt.Sprintf("%s\t\t\treturn fmt.Errorf(\"provided resource reference is nil or empty\")\n", indent)
	out += fmt.Sprintf("%s\t\t}\n", indent)
	out += fmt.Sprintf("%s\n", indent)
	out += readReferenceAndValidate(field, indentLevel+2)
	out += fmt.Sprintf("%s\t\treferencedValue := string(*obj.%s)\n", indent, field.FieldConfig.References.Path)
	out += fmt.Sprintf("%s\t\telem.%s = &referencedValue\n", indent, field.Names.Camel)
	out += fmt.Sprintf("%s\t}\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%sreturn nil", indent)

	return out
}

// readReferenceAndValidate produces Go code that attempts to fetch a referenced
// object from the K8s API server and validates whether it is in synced or
// terminal conditions, returning the appropriate errors.
func readReferenceAndValidate(field *model.Field, indentLevel int) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	out += fmt.Sprintf("%sif arr == nil || arr.Name == nil || *arr.Name == \"\" {\n", indent)
	out += fmt.Sprintf("%s\treturn fmt.Errorf(\"provided resource reference is nil or empty\")\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%snamespacedName := types.NamespacedName{\n", indent)
	out += fmt.Sprintf("%s\tNamespace: namespace,\n", indent)
	out += fmt.Sprintf("%s\tName: *arr.Name,\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	if field.FieldConfig.References.ServiceName == "" {
		out += fmt.Sprintf("%sobj := svcapitypes.%s{}\n", indent, field.FieldConfig.References.Resource)
	} else {
		out += fmt.Sprintf("%sobj := %sapitypes.%s{}\n", indent, field.ReferencedServiceName(), field.FieldConfig.References.Resource)
	}
	out += fmt.Sprintf("%serr := apiReader.Get(ctx, namespacedName, &obj)\n", indent)
	out += fmt.Sprintf("%sif err != nil {\n", indent)
	out += fmt.Sprintf("%s\treturn err\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%svar refResourceSynced, refResourceTerminal bool\n", indent)
	out += fmt.Sprintf("%sfor _, cond := range obj.Status.Conditions {\n", indent)
	out += fmt.Sprintf("%s\tif cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&\n", indent)
	out += fmt.Sprintf("%s\t\tcond.Status == corev1.ConditionTrue {\n", indent)
	out += fmt.Sprintf("%s\t\trefResourceSynced = true\n", indent)
	out += fmt.Sprintf("%s\t}\n", indent)
	out += fmt.Sprintf("%s\tif cond.Type == ackv1alpha1.ConditionTypeTerminal &&\n", indent)
	out += fmt.Sprintf("%s\t\tcond.Status == corev1.ConditionTrue {\n", indent)
	out += fmt.Sprintf("%s\t\trefResourceTerminal = true\n", indent)
	out += fmt.Sprintf("%s\t}\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%sif refResourceTerminal {\n", indent)
	out += fmt.Sprintf("%s\treturn ackerr.ResourceReferenceTerminalFor(\n", indent)
	out += fmt.Sprintf("%s\t\t\"%s\",\n", indent, field.FieldConfig.References.Resource)
	out += fmt.Sprintf("%s\t\tnamespace, *arr.Name)\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%sif !refResourceSynced {\n", indent)
	out += fmt.Sprintf("%s\treturn ackerr.ResourceReferenceNotSyncedFor(\n", indent)
	out += fmt.Sprintf("%s\t\t\"%s\",\n", indent, field.FieldConfig.References.Resource)
	out += fmt.Sprintf("%s\t\tnamespace, *arr.Name)\n", indent)
	out += fmt.Sprintf("%s}\n", indent)

	nilCheck := CheckNilReferencesPath(field, "obj")
	if nilCheck != "" {
		out += fmt.Sprintf("%sif %s {\n", indent, nilCheck)
		out += fmt.Sprintf("%s\treturn ackerr.ResourceReferenceMissingTargetFieldFor(\n", indent)
		out += fmt.Sprintf("%s\t\t\"%s\",\n", indent, field.FieldConfig.References.Resource)
		out += fmt.Sprintf("%s\t\tnamespace, *arr.Name,\n", indent)
		out += fmt.Sprintf("%s\t\t\"%s\")\n", indent, field.FieldConfig.References.Path)
		out += fmt.Sprintf("%s}\n", indent)
	}

	return out
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
