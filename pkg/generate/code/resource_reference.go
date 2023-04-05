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

	lo "github.com/samber/lo"

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
//		if ko.Spec.APIRef != nil && ko.Spec.APIRef.From != nil {
//		arr := ko.Spec.APIRef.From
//		if arr == nil || arr.Name == nil || *arr.Name == "" {
//			return fmt.Errorf("provided resource reference is nil or empty: APIRef")
//		}
//		obj := &svcapitypes.API{}
//		if err := getReferencedResourceState_API(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
//			return err
//		}
//		ko.Spec.APIID = (*string)(obj.Status.APIID)
//	}
//
// ```
func ResolveReferencesForField(field *model.Field, sourceVarName string, indentLevel int) string {
	r := field.CRD
	fp := fieldpath.FromString(field.Path)

	outPrefix := ""
	outSuffix := ""

	indent := strings.Repeat("\t", indentLevel)

	numLists := field.GetNumberParentLists()
	currentListDepth := 0 // Stores how many slices we are iterating within

	resRefVar := "resolved"
	nestedIndexVarPrefixFmt := "f%didx"

	resRefElemType := field.ShapeRef.GoType()

	// If the field is within a list, or nested lists, then use a
	// multi-dimensional array to hold the relative indexes of the reference
	// within those parent lists
	resRefType := fmt.Sprintf("%s%s", strings.Repeat("[]", numLists), resRefElemType)

	outPrefix += fmt.Sprintf("%svar %s %s\n", indent, resRefVar, resRefType)
	outSuffix = fmt.Sprintf("%srm.setReferencedValue(ctx, %q, %s)\n%s", indent, field.Path, resRefVar, outSuffix)

	fieldAccessPrefix := fmt.Sprintf("%s%s", sourceVarName, r.Config().PrefixConfig.SpecField)

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
			iterVarName := fmt.Sprintf("f%diter", idx)
			idxVarName := fmt.Sprintf(nestedIndexVarPrefixFmt, idx)

			fieldAccessPrefix = fmt.Sprintf("%s.%s", fieldAccessPrefix, fp.At(idx))

			// Initialise the nested slice within resolved
			outPrefix += fmt.Sprintf("%s%s = make(%s%s, len(%s))\n", indent, resRefVar, strings.Repeat("[]", numLists-currentListDepth), resRefElemType, fieldAccessPrefix)
			outPrefix += fmt.Sprintf("%sfor %s, %s := range %s {\n", indent, idxVarName, iterVarName, fieldAccessPrefix)

			outSuffix = fmt.Sprintf("%s}\n%s", indent, outSuffix)

			currentListDepth++

			fieldAccessPrefix = iterVarName
		case ("map"):
			panic("references cannot be within a map")
		default:
			// base case for single references
			fieldAccessPrefix = fmt.Sprintf("%s.%s", fieldAccessPrefix, cur.GetReferenceFieldName().Camel)

			outPrefix += fmt.Sprintf("%sif %s != nil && %s.From != nil {\n", indent, fieldAccessPrefix, fieldAccessPrefix)
			outPrefix += fmt.Sprintf("%s\tarr := %s.From\n", indent, fieldAccessPrefix)
			outPrefix += fmt.Sprintf("%s\tif arr == nil || arr.Name == nil || *arr.Name == \"\" {\n", indent)
			outPrefix += fmt.Sprintf("%s\t\treturn fmt.Errorf(\"provided resource reference is nil or empty: %s\")\n", indent, field.ReferenceFieldPath())
			outPrefix += fmt.Sprintf("%s\t}\n", indent)

			outPrefix += getReferencedStateForField(field, indentLevel+idx)

			// if we are inside lists, set it to the appropriate indexes
			indexList := strings.Join(lo.Times(numLists, func(indx int) string {
				return fmt.Sprintf("[%s]", fmt.Sprintf(nestedIndexVarPrefixFmt, indx))
			}), "")
			outPrefix += fmt.Sprintf("%s\t\t%s%s = obj.%s\n", indent, resRefVar, indexList, field.FieldConfig.References.Path)
			outPrefix += fmt.Sprintf("%s}\n", indent)
		}
	}

	return outPrefix + outSuffix
}

func CopyWithResolvedReferences(field *model.Field, targetVarName string, indentLevel int) string {
	r := field.CRD

	numLists := field.GetNumberParentLists()
	nestedIndexVarPrefixFmt := "f%didx"

	outPrefix, outSuffix := IterResolvedReferenceValues(field, indentLevel, true, nestedIndexVarPrefixFmt, func(nestedIndent string, fieldNamePrefix string, castResolvedVar string) (outPrefix string, outSuffix string) {
		indexList := strings.Join(lo.Times(numLists, func(indx int) string {
			return fmt.Sprintf("[%s]", fmt.Sprintf(nestedIndexVarPrefixFmt, indx))
		}), "")
		outPrefix += fmt.Sprintf("%s%s%s%s = %s%s\n", nestedIndent, targetVarName, r.Config().PrefixConfig.SpecField, fieldNamePrefix, castResolvedVar, indexList)
		return outPrefix, outSuffix
	})

	return outPrefix + outSuffix
}

func ClearResolvedReferences(field *model.Field, targetVarName string, indentLevel int) string {
	r := field.CRD
	numLists := field.GetNumberParentLists()

	nestedIndexVarPrefixFmt := "f%didx"
	shouldCast := numLists > 0

	outPrefix, outSuffix := IterResolvedReferenceValues(field, indentLevel, shouldCast, nestedIndexVarPrefixFmt, func(nestedIndent string, fieldNamePrefix string, castResolvedVar string) (outPrefix string, outSuffix string) {
		outPrefix += fmt.Sprintf("%s%s%s%s = nil\n", nestedIndent, targetVarName, r.Config().PrefixConfig.SpecField, fieldNamePrefix)
		return outPrefix, outSuffix
	})

	return outPrefix + outSuffix
}

func IterResolvedReferenceValues(field *model.Field, indentLevel int, shouldCast bool, nestedIndexVarPrefixFmt string, innerRender func(nestedIndent string, fieldNamePrefix string, castResolvedVar string) (outPrefix string, outSuffix string)) (outPrefix string, outSuffix string) {
	r := field.CRD
	fp := fieldpath.FromString(field.Path)

	indent := strings.Repeat("\t", indentLevel)

	numLists := field.GetNumberParentLists()

	castResolvedVar := "castResRef"

	if shouldCast {
		outPrefix += fmt.Sprintf("%sif val, ok := rm.getReferencedValue(%q); ok {\n", indent, field.Path)
	} else {
		// Since we aren't using the value in the cast, we need to ignore it
		// in the response value
		outPrefix += fmt.Sprintf("%sif _, ok := rm.getReferencedValue(%q); ok {\n", indent, field.Path)
	}
	outSuffix = fmt.Sprintf("%s}\n%s", indent, outSuffix)

	if shouldCast {
		resRefElemType := field.ShapeRef.GoType()
		resRefType := fmt.Sprintf("%s%s", strings.Repeat("[]", numLists), resRefElemType)
		outPrefix += fmt.Sprintf("%s%s, ok := (val).(%s)\n", indent, castResolvedVar, resRefType)
		outPrefix += fmt.Sprintf("%sif !ok {\n", indent)
		outPrefix += fmt.Sprintf("%s\treturn nil, ackerr.ResourceReferenceValueCastFailedFor(%q)\n", indent, field.Path)
		outPrefix += fmt.Sprintf("%s}\n", indent)
	}

	// Iterate through nested lists if necessary
	for idx := 0; idx < numLists; idx++ {
		innerIndent := fmt.Sprintf("%s%s", indent, strings.Repeat("\t", idx))
		idxVarName := fmt.Sprintf(nestedIndexVarPrefixFmt, idx)

		parentIndexList := strings.Join(lo.Times(idx, func(indx int) string {
			return fmt.Sprintf("[%s]", fmt.Sprintf(nestedIndexVarPrefixFmt, indx))
		}), "")

		outPrefix += fmt.Sprintf("%sfor %s, _ := range %s%s {\n", innerIndent, idxVarName, castResolvedVar, parentIndexList)
		outSuffix = fmt.Sprintf("%s}\n%s", indent, outSuffix)
	}

	nestedIndent := strings.Repeat("\t", indentLevel+numLists)

	// Build the target field accessor
	fieldNamePrefix := ""
	nestedFieldDepth := 0
	for idx := 0; idx < fp.Size(); idx++ {
		curFP := fp.CopyAt(idx)

		cur, ok := r.Fields[curFP.String()]
		if !ok {
			panic(fmt.Sprintf("unable to find field with path %q. crd: %q", curFP.String(), r.Kind))
		}

		fieldName := curFP.Pop()
		indexList := ""

		if cur.ShapeRef.Shape.Type == "list" {
			indexList = fmt.Sprintf("[%s]", fmt.Sprintf(nestedIndexVarPrefixFmt, nestedFieldDepth))
			nestedFieldDepth++
		}

		fieldNamePrefix = fmt.Sprintf("%s.%s%s", fieldNamePrefix, fieldName, indexList)
	}

	innerPrefix, innerSuffix := innerRender(nestedIndent, fieldNamePrefix, castResolvedVar)

	return (outPrefix + innerPrefix), (innerSuffix + outSuffix)
}

func getReferencedStateForField(field *model.Field, indentLevel int) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	if field.FieldConfig.References.ServiceName == "" {
		out += fmt.Sprintf("%s\tobj := &svcapitypes.%s{}\n", indent, field.FieldConfig.References.Resource)
	} else {
		out += fmt.Sprintf("%s\tobj := &%sapitypes.%s{}\n", indent, field.ReferencedServiceName(), field.FieldConfig.References.Resource)
	}
	out += fmt.Sprintf("%s\tif err := getReferencedResourceState_%s(ctx, apiReader, obj, *arr.Name, namespace); err != nil {\n", indent, field.FieldConfig.References.Resource)
	out += fmt.Sprintf("%s\t\treturn err\n", indent)
	out += fmt.Sprintf("%s\t}\n", indent)

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
