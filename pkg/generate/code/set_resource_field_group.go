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

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/fieldpath"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// SetResourceFieldGroup returns the Go code that sets a CRD's field values
// from the Output shape of a field-group operation. Only members that
// correspond to the field group's payload fields are set.
//
// This is the field-group analog of SetResource. It generates code like:
//
//	if resp.ImageScanningConfiguration != nil {
//	    f0 := &svcapitypes.ImageScanningConfiguration{}
//	    ...
//	    ko.Spec.ImageScanningConfiguration = f0
//	}
func SetResourceFieldGroup(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	fg *model.FieldGroupOperation,
	// String representing the name of the variable that we will grab the
	// Output values from (typically "resp")
	sourceVarName string,
	// String representing the name of the variable that we will be setting
	// (typically "ko")
	targetVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) (string, error) {
	op := fg.Operation
	if op == nil {
		return "", nil
	}
	outputShape := op.OutputRef.Shape
	if outputShape == nil {
		return "", nil
	}

	// For field-group operations, the payload fields are what we read back.
	// Build a set of their CRD field names for filtering.
	payloadFieldNames := make(map[string]bool, len(fg.PayloadFields))
	for _, f := range fg.PayloadFields {
		payloadFieldNames[f.Names.Camel] = true
	}

	out := "\n"
	indent := strings.Repeat("\t", indentLevel)

	// Determine the opType for setter config lookup
	opType := model.OpTypeUpdate
	if fg.OpType == model.FieldGroupOpTypeRead {
		opType = model.OpTypeGet
	}

	for memberIndex, memberName := range outputShape.MemberNames() {
		sourceAdaptedVarName := sourceVarName + "." + memberName

		// Resolve the CRD field name
		fieldName := cfg.GetResourceFieldName(
			r.Names.Original, op.ExportedName, memberName,
		)

		// Look up the field in Spec or Status
		var f *model.Field
		targetAdaptedVarName := targetVarName
		inSpec, inStatus := r.HasMember(fieldName, op.ExportedName)
		if inSpec {
			f = r.SpecFields[fieldName]
			targetAdaptedVarName += cfg.PrefixConfig.SpecField
		} else if inStatus {
			f = r.StatusFields[fieldName]
			targetAdaptedVarName += cfg.PrefixConfig.StatusField
		} else {
			continue
		}

		// Only process fields that are in this field group's payload
		if !payloadFieldNames[f.Names.Camel] {
			continue
		}

		targetMemberShapeRef := f.ShapeRef
		setCfg := f.GetSetterConfig(opType)
		if setCfg != nil && setCfg.IgnoreResourceSetter() {
			continue
		}

		sourceMemberShapeRef := outputShape.MemberRefs[memberName]
		if sourceMemberShapeRef.Shape == nil {
			if setCfg != nil && setCfg.From != nil {
				fp := fieldpath.FromString(*setCfg.From)
				sourceMemberShapeRef = fp.ShapeRef(sourceMemberShapeRef)
			}
			if sourceMemberShapeRef == nil || sourceMemberShapeRef.Shape == nil {
				return "", fmt.Errorf(
					"resource %q, field %q: expected .Shape to not be nil for ShapeRef",
					r.Names.Original, memberName,
				)
			}
		}

		if sourceMemberShapeRef.Shape.RealType == "union" {
			sourceMemberShapeRef.Shape.Type = "union"
		}

		targetMemberShape := targetMemberShapeRef.Shape

		if sourceMemberShapeRef.Shape.IsEnum() {
			out += fmt.Sprintf(
				"%sif %s != \"\" {\n", indent, sourceAdaptedVarName,
			)
		} else if !sourceMemberShapeRef.HasDefaultValue() {
			out += fmt.Sprintf(
				"%sif %s != nil {\n", indent, sourceAdaptedVarName,
			)
		} else {
			indentLevel -= 1
		}

		qualifiedTargetVar := fmt.Sprintf(
			"%s.%s", targetAdaptedVarName, f.Names.Camel,
		)

		switch targetMemberShape.Type {
		case "list", "map", "structure", "union":
			adaption := setResourceAdaptPrimitiveCollection(
				sourceMemberShapeRef.Shape, qualifiedTargetVar,
				sourceAdaptedVarName, indent, r.IsSecretField(memberName),
			)
			out += adaption
			if adaption != "" {
				break
			}
			{
				memberVarName := fmt.Sprintf("f%d", memberIndex)
				out += varEmptyConstructorK8sType(
					cfg, r, memberVarName,
					targetMemberShapeRef.Shape, indentLevel+1,
				)
				containerOut, err := setResourceForContainer(
					cfg, r, f.Names.Camel, memberVarName,
					targetMemberShapeRef, setCfg,
					sourceAdaptedVarName, sourceMemberShapeRef,
					f.Names.Camel, false, opType, indentLevel+1,
				)
				if err != nil {
					return "", err
				}
				out += containerOut
				out += setResourceForScalar(
					qualifiedTargetVar, memberVarName,
					sourceMemberShapeRef, indentLevel+1, false, false,
				)
			}
		default:
			if setCfg != nil && setCfg.From != nil {
				sourceAdaptedVarName = sourceVarName + "." + *setCfg.From
			}
			out += setResourceForScalar(
				qualifiedTargetVar, sourceAdaptedVarName,
				sourceMemberShapeRef, indentLevel+1, false, false,
			)
		}
		if sourceMemberShapeRef.Shape.RealType == "union" {
			sourceMemberShapeRef.Shape.Type = "structure"
		}
		if sourceMemberShapeRef.Shape.IsEnum() || !sourceMemberShapeRef.HasDefaultValue() {
			out += fmt.Sprintf(
				"%s} else {\n", indent,
			)
			out += fmt.Sprintf(
				"%s%s%s.%s = nil\n", indent, indent,
				targetAdaptedVarName, f.Names.Camel,
			)
			out += fmt.Sprintf(
				"%s}\n", indent,
			)
		} else {
			indentLevel += 1
		}
	}
	return out, nil
}
