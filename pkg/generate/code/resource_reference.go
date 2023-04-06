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
	"github.com/samber/lo"
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
	field *model.Field,
	sourceVarName string,
	indentLevel int,
) (out string) {
	out = IterReferenceValues(field, indentLevel, sourceVarName, false, false, true,
		func(fieldAccessPrefix string, indexVarFmt string, innerIndentLevel int) (innerOut string) {
			innerIndent := strings.Repeat("\t", innerIndentLevel)

			// Get parent field path, in order to get to both the refs and concretes
			parentFP := fieldpath.FromString(fieldAccessPrefix)
			parentFP.Pop()

			// Validation to make sure both target field and reference are
			// not present at the same time in desired resource
			innerOut += fmt.Sprintf("%sif %s.%s != nil"+
				" && %s.%s != nil {\n", innerIndent, parentFP.String(), field.GetReferenceFieldName().Camel, parentFP.String(), field.Names.Camel)
			innerOut += fmt.Sprintf("%s\treturn "+
				"ackerr.ResourceReferenceAndIDNotSupportedFor(%q, %q)\n",
				innerIndent, field.Path, field.ReferenceFieldPath())
			innerOut += fmt.Sprintf("%s}\n", innerIndent)

			// If the field is required, make sure either Ref or original
			// field is present in the resource
			if field.IsRequired() {
				innerOut += fmt.Sprintf("%sif %s.%s == nil &&"+
					" %s.%s == nil {\n", innerIndent, parentFP.String(),
					field.GetReferenceFieldName().Camel, parentFP.String(), field.Names.Camel)
				innerOut += fmt.Sprintf("%s\treturn "+
					"ackerr.ResourceReferenceOrIDRequiredFor(%q, %q)\n",
					innerIndent, field.Names.Camel, field.GetReferenceFieldName().Camel)
				innerOut += fmt.Sprintf("%s}\n", innerIndent)
			}

			return innerOut
		})

	return out
}

// ResolveReferencesForField produces Go code for accessing all references that
// are related to the given concrete field, determining whether its in a valid
// condition and updating the concrete field with the referenced value.
// Sample code (resolving a nested singular reference):
//
// ```
// var resolved *string
//
//	if ko.Spec.KMSKeyRef != nil && ko.Spec.KMSKeyRef.From != nil {
//		arr := ko.Spec.KMSKeyRef.From
//		if arr.Name == nil || *arr.Name == "" {
//			return fmt.Errorf("provided resource reference is nil or empty: KMSKeyRef")
//		}
//		obj := &kmsapitypes.Key{}
//		if err := getReferencedResourceState_Key(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
//			return err
//		}
//		resolved = (*string)(obj.Status.ACKResourceMetadata.ARN)
//	}
//
// rm.setReferencedValue(ctx, "KMSKeyARN", resolved)
// ```
//
// Sample code (resolving a list of references):
// ```
// var resolved []*string
//
//	if ko.Spec.VPCConfig != nil {
//		for f0idx, f0iter := range ko.Spec.VPCConfig.SecurityGroupRefs {
//			if f0iter != nil && f0iter.From != nil {
//				arr := f0iter.From
//				if arr.Name == nil || *arr.Name == "" {
//					return fmt.Errorf("provided resource reference is nil or empty: VPCConfig.SecurityGroupRefs")
//				}
//				obj := &ec2apitypes.SecurityGroup{}
//				if err := getReferencedResourceState_SecurityGroup(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
//					return err
//				}
//				resolved[f0idx] = (*string)(obj.Status.ID)
//			}
//		}
//	}
//
// rm.setReferencedValue(ctx, "VPCConfig.SecurityGroupIDs", resolved)
// ```
//
// Sample code (resolving a list of structs containing references):
// ```
// var resolved []*string
// resolved = make([]*string, len(ko.Spec.Routes))
//
//	for f0idx, f0iter := range ko.Spec.Routes {
//		if f0iter.VPCEndpointRef != nil && f0iter.VPCEndpointRef.From != nil {
//			arr := f0iter.VPCEndpointRef.From
//			if arr.Name == nil || *arr.Name == "" {
//				return fmt.Errorf("provided resource reference is nil or empty: Routes.VPCEndpointRef")
//			}
//			obj := &svcapitypes.VPCEndpoint{}
//			if err := getReferencedResourceState_VPCEndpoint(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
//				return err
//			}
//			resolved[f0idx] = (*string)(obj.Status.VPCEndpointID)
//		}
//	}
//
// rm.setReferencedValue(ctx, "Routes.VPCEndpointID", resolved)
// ```
func ResolveReferencesForField(field *model.Field, sourceVarName string, indentLevel int) string {
	isListOfRefs := field.ShapeRef.Shape.Type == "list"

	resRefElemType := field.ShapeRef.GoType()
	if isListOfRefs {
		resRefElemType = field.ShapeRef.Shape.MemberRef.GoType()
	}

	iterOut := IterReferenceValues(field, indentLevel, sourceVarName, true, true, true,
		func(fieldAccessPrefix string, indexVarFmt string, innerIndentLevel int) (innerOut string) {
			innerIndent := strings.Repeat("\t", innerIndentLevel)

			innerOut += fmt.Sprintf("%shasReferences = true\n", innerIndent)

			innerOut += fmt.Sprintf("%sarr := %s.From\n", innerIndent, fieldAccessPrefix)
			innerOut += fmt.Sprintf("%sif arr.Name == nil || *arr.Name == \"\" {\n", innerIndent)
			innerOut += fmt.Sprintf("%s\treturn hasReferences, fmt.Errorf(\"provided resource reference is nil or empty: %s\")\n", innerIndent, field.ReferenceFieldPath())
			innerOut += fmt.Sprintf("%s}\n", innerIndent)

			innerOut += getReferencedStateForField(field, innerIndentLevel)

			concreteValueAccessor := buildIndexBasedFieldAccessor(field, sourceVarName, indexVarFmt)
			if isListOfRefs {
				innerOut += fmt.Sprintf("%sif %s == nil {\n", innerIndent, concreteValueAccessor)
				innerOut += fmt.Sprintf("%s\t%s = make([]%s, 0, 1)\n", innerIndent, concreteValueAccessor, resRefElemType)
				innerOut += fmt.Sprintf("%s}\n", innerIndent)
				innerOut += fmt.Sprintf("%s%s = append(%s, (%s)(obj.%s))\n", innerIndent, concreteValueAccessor, concreteValueAccessor, resRefElemType, field.FieldConfig.References.Path)
			} else {
				innerOut += fmt.Sprintf("%s%s = (%s)(obj.%s)\n", innerIndent, concreteValueAccessor, resRefElemType, field.FieldConfig.References.Path)
			}

			return innerOut
		})

	return iterOut
}

func ClearResolvedReferencesForField(field *model.Field, targetVarName string, indentLevel int) string {
	isListOfRefs := field.ShapeRef.Shape.Type == "list"

	iterOut := IterReferenceValues(field, indentLevel, targetVarName, false, true, false,
		func(fieldAccessPrefix string, indexVarFmt string, innerIndentLevel int) (innerOut string) {
			innerIndent := strings.Repeat("\t", innerIndentLevel)

			// If we are dealing with a list of references, then we don't need to
			// iterate over all of the references individually. We know that if the list
			// has >0 elements, then the entire concrete value list should be made nil.
			// To deal with this, we should iterate only to the parent of the list and
			// then check these conditions.
			if isListOfRefs {
				innerOut += fmt.Sprintf("%sif len(%s) > 0 {\n", innerIndent, fieldAccessPrefix)
			} else {
				innerOut += fmt.Sprintf("%sif %s != nil {\n", innerIndent, fieldAccessPrefix)
			}
			innerOut += fmt.Sprintf("%s\t%s = nil\n", innerIndent, buildIndexBasedFieldAccessor(field, targetVarName, indexVarFmt))
			innerOut += fmt.Sprintf("%s}\n", innerIndent)

			return innerOut
		})

	return iterOut
}

func IterReferenceValues(field *model.Field, indentLevel int, sourceVarName string, shouldValidateRef bool, shouldRenderIndexes, shouldRenderIterators bool, innerRender func(fieldAccessPrefix string, indexVarFmt string, indentLevel int) (innerOut string)) (out string) {
	r := field.CRD
	fp := fieldpath.FromString(field.ReferenceFieldPath())

	outPrefix := ""
	outSuffix := ""

	currentListDepth := 0 // Stores how many slices we are iterating within

	nestedIterVarPrefixFmt := "f%diter"
	nestedIndexVarPrefixFmt := "f%didx"

	fieldAccessPrefix := fmt.Sprintf("%s%s", sourceVarName, r.Config().PrefixConfig.SpecField)

	fpDepth := 0
	for fpDepth = 0; fpDepth < fp.Size()-1; fpDepth++ {
		curFP := fp.CopyAt(fpDepth).String()
		cur, ok := r.Fields[curFP]
		if !ok {
			panic(fmt.Sprintf("unable to find field with path %q. crd: %q", curFP, r.Kind))
		}

		ref := cur.ShapeRef

		indent := strings.Repeat("\t", indentLevel+fpDepth)

		switch ref.Shape.Type {
		case ("map"):
			panic("references cannot be within a map")
		case ("structure"):
			fieldAccessPrefix = fmt.Sprintf("%s.%s", fieldAccessPrefix, fp.At(fpDepth))

			outPrefix += fmt.Sprintf("%sif %s != nil {\n", indent, fieldAccessPrefix)
			outSuffix = fmt.Sprintf("%s}\n%s", indent, outSuffix)
		case ("list"):
			iterVarName := fmt.Sprintf(nestedIterVarPrefixFmt, currentListDepth)
			idxVarName := fmt.Sprintf(nestedIndexVarPrefixFmt, currentListDepth)

			fieldAccessPrefix = fmt.Sprintf("%s.%s", fieldAccessPrefix, fp.At(fpDepth))

			outPrefix += fmt.Sprintf("%sfor %s, %s := range %s {\n", indent,
				lo.Ternary(shouldRenderIndexes, idxVarName, "_"),
				iterVarName,
				fieldAccessPrefix,
			)

			outSuffix = fmt.Sprintf("%s}\n%s", indent, outSuffix)

			fieldAccessPrefix = iterVarName

			currentListDepth++
		}
	}

	fieldAccessPrefix = fmt.Sprintf("%s.%s", fieldAccessPrefix, field.GetReferenceFieldName().Camel)
	innerIndentLevel := indentLevel + fpDepth

	if shouldValidateRef {
		// If the reference field is a list of primitives, iterate through each.
		// We need to duplicate some code from above, here, because `*Ref` fields
		// aren't registered as fields, so we can't use the same common logic
		if field.ShapeRef.Shape.Type == "list" {
			iterVarName := fmt.Sprintf(nestedIterVarPrefixFmt, currentListDepth)

			outPrefix += fmt.Sprintf("%sfor _, %s := range %s {\n", strings.Repeat("\t", innerIndentLevel),
				lo.Ternary(shouldRenderIterators, iterVarName, "_"),
				fieldAccessPrefix,
			)
			outSuffix = fmt.Sprintf("%s}\n%s", strings.Repeat("\t", innerIndentLevel), outSuffix)
			fieldAccessPrefix = iterVarName

			innerIndentLevel++
		}

		outPrefix += fmt.Sprintf("%sif %s != nil && %s.From != nil {\n", strings.Repeat("\t", innerIndentLevel), fieldAccessPrefix, fieldAccessPrefix)
		outSuffix = fmt.Sprintf("%s}\n%s", strings.Repeat("\t", innerIndentLevel), outSuffix)
		innerIndentLevel++
	}

	innerPrefix := innerRender(fieldAccessPrefix, nestedIndexVarPrefixFmt, innerIndentLevel)

	return outPrefix + innerPrefix + outSuffix
}

// buildNestedFieldAccessor generates Go code that accesses an inner struct,
// using slice indexes where necessary.
//
// `indexVarFmt` should be a format string that takes a single integer and
// returns the name of a variable which holds the index for the n-th parent
// slice. For example, f%didx will be used to create f0idx, f1idx, etc. for the
// parent slices in the accessors.
//
// By default, this method will iterate through every field in the field path.
// Supplying a `parentOffset` will only iterate through the first `fp.Size() -
// parentOffset` number of paths.
func buildIndexBasedFieldAccessorWithOffset(field *model.Field, sourceVarName, indexVarFmt string, parentOffset int) string {
	r := field.CRD
	fp := fieldpath.FromString(field.Path)

	isList := field.ShapeRef.Shape.Type == "list"

	fieldNamePrefix := ""
	nestedFieldDepth := 0
	for idx := 0; idx < fp.Size()-parentOffset; idx++ {
		curFP := fp.CopyAt(idx)

		cur, ok := r.Fields[curFP.String()]
		if !ok {
			panic(fmt.Sprintf("unable to find field with path %q. crd: %q", curFP.String(), r.Kind))
		}

		fieldName := curFP.Pop()
		indexList := ""

		if cur.ShapeRef.Shape.Type == "list" {

			// We want to access indexes when iterating through lists of
			// structs. If we find a list at the end of the field path, then we
			// know the initial field must be a list. We want to pass that back
			// as a full array, rather than accessing the individual values.
			// This only applies for when there is no offset, since any offset >
			// 0 will cut off the initial field from the path
			if idx != (fp.Size()-1) && !isList {
				indexList = fmt.Sprintf("[%s]", fmt.Sprintf(indexVarFmt, nestedFieldDepth))
				nestedFieldDepth++
			}
		}

		fieldNamePrefix = fmt.Sprintf("%s.%s%s", fieldNamePrefix, fieldName, indexList)
	}

	return fmt.Sprintf("%s%s%s", sourceVarName, r.Config().PrefixConfig.SpecField, fieldNamePrefix)
}

// buildIndexBasedFieldAccessor calls buildNestedFieldAccessorWithOffset with an
// offset of 0.
func buildIndexBasedFieldAccessor(field *model.Field, sourceVarName, indexVarFmt string) string {
	return buildIndexBasedFieldAccessorWithOffset(field, sourceVarName, indexVarFmt, 0)
}

func getReferencedStateForField(field *model.Field, indentLevel int) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	if field.FieldConfig.References.ServiceName == "" {
		out += fmt.Sprintf("%sobj := &svcapitypes.%s{}\n", indent, field.FieldConfig.References.Resource)
	} else {
		out += fmt.Sprintf("%sobj := &%sapitypes.%s{}\n", indent, field.ReferencedServiceName(), field.FieldConfig.References.Resource)
	}
	out += fmt.Sprintf("%sif err := getReferencedResourceState_%s(ctx, apiReader, obj, *arr.Name, namespace); err != nil {\n", indent, field.FieldConfig.References.Resource)
	out += fmt.Sprintf("%s\treturn hasReferences, err\n", indent)
	out += fmt.Sprintf("%s}\n", indent)

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
