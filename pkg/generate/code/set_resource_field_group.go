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

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
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

	// Track which payload fields were successfully matched from the output
	// shape members. Any unmatched fields with a from.path config will be
	// handled by the unwrapped-output fallback below.
	matchedFields := make(map[string]bool)

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
		matchedFields[f.Names.Camel] = true
	}

	// Fallback: for payload fields that weren't matched by any output member,
	// check if the field has a from.path config. If so, the output response
	// members are the *unwrapped* content of that shape, and we generate code
	// that constructs the wrapper struct from the flat output members.
	for _, f := range fg.PayloadFields {
		if matchedFields[f.Names.Camel] {
			continue
		}
		if f.FieldConfig == nil || f.FieldConfig.From == nil {
			continue
		}
		wrapperOut, err := setResourceFieldGroupWrapped(
			cfg, r, f, fg, outputShape, sourceVarName, targetVarName,
			indentLevel,
		)
		if err != nil {
			return "", err
		}
		out += wrapperOut
	}

	return out, nil
}

// setResourceFieldGroupWrapped generates code for a payload field whose Get
// output is the unwrapped content of the field's from.path shape. It
// constructs the wrapper struct and maps output members into it.
//
// For example, if the CRD field is Accelerate (*AccelerateConfiguration) with
// from.path = "AccelerateConfiguration", and the Get output has member
// "Status" (which is a member of AccelerateConfiguration), we generate:
//
//	ko.Spec.Accelerate = &svcapitypes.AccelerateConfiguration{}
//	ko.Spec.Accelerate.Status = aws.String(string(resp.Status))
func setResourceFieldGroupWrapped(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	f *model.Field,
	fg *model.FieldGroupOperation,
	outputShape *awssdkmodel.Shape,
	sourceVarName string,
	targetVarName string,
	indentLevel int,
) (string, error) {
	fromCfg := f.FieldConfig.From
	// Look up the wrapper shape from the from.operation's Input shape
	fromOpID := fromCfg.Operation
	fromOp, ok := r.GetSDKAPI().API.Operations[fromOpID]
	if !ok {
		return "", fmt.Errorf(
			"field %q: from.operation %q not found in SDK",
			f.Names.Camel, fromOpID,
		)
	}
	fromInputShape := fromOp.InputRef.Shape
	if fromInputShape == nil {
		return "", fmt.Errorf(
			"field %q: from.operation %q has nil Input shape",
			f.Names.Camel, fromOpID,
		)
	}

	// Navigate from.path to find the wrapper shape
	wrapperShapeRef, ok := fromInputShape.MemberRefs[fromCfg.Path]
	if !ok {
		return "", fmt.Errorf(
			"field %q: from.path %q not found in %s Input shape",
			f.Names.Camel, fromCfg.Path, fromOpID,
		)
	}
	wrapperShape := wrapperShapeRef.Shape
	if wrapperShape == nil {
		return "", fmt.Errorf(
			"field %q: from.path %q has nil shape",
			f.Names.Camel, fromCfg.Path,
		)
	}

	indent := strings.Repeat("\t", indentLevel)
	targetPrefix := targetVarName + cfg.PrefixConfig.SpecField + "." + f.Names.Camel

	out := ""
	// Generate the wrapper struct construction using the existing
	// setResourceForContainer machinery. The output response serves as the
	// source, and the wrapper shape defines both source and target member
	// structure.
	out += fmt.Sprintf(
		"%s%s = &svcapitypes.%s{}\n",
		indent, targetPrefix, wrapperShape.ShapeName,
	)

	// Map each wrapper shape member from the output response
	for _, wrapperMemberName := range wrapperShape.MemberNames() {
		// Check if the output shape has this member
		outputMemberRef, outputHasMember := outputShape.MemberRefs[wrapperMemberName]
		if !outputHasMember || outputMemberRef.Shape == nil {
			continue
		}
		sourceMemberVar := sourceVarName + "." + wrapperMemberName
		targetMemberVar := targetPrefix + "." + wrapperMemberName

		outputMemberShape := outputMemberRef.Shape
		switch outputMemberShape.Type {
		case "string":
			if outputMemberShape.IsEnum() {
				out += fmt.Sprintf(
					"%sif %s != \"\" {\n", indent, sourceMemberVar,
				)
				out += fmt.Sprintf(
					"%s\t%s = aws.String(string(%s))\n",
					indent, targetMemberVar, sourceMemberVar,
				)
				out += fmt.Sprintf("%s}\n", indent)
			} else {
				out += fmt.Sprintf(
					"%sif %s != nil {\n", indent, sourceMemberVar,
				)
				out += fmt.Sprintf(
					"%s\t%s = %s\n",
					indent, targetMemberVar, sourceMemberVar,
				)
				out += fmt.Sprintf("%s}\n", indent)
			}
		case "integer", "long":
			out += fmt.Sprintf(
				"%sif %s != nil {\n", indent, sourceMemberVar,
			)
			out += fmt.Sprintf(
				"%s\t%s = %s\n",
				indent, targetMemberVar, sourceMemberVar,
			)
			out += fmt.Sprintf("%s}\n", indent)
		case "boolean":
			out += fmt.Sprintf(
				"%sif %s != nil {\n", indent, sourceMemberVar,
			)
			out += fmt.Sprintf(
				"%s\t%s = %s\n",
				indent, targetMemberVar, sourceMemberVar,
			)
			out += fmt.Sprintf("%s}\n", indent)
		case "list", "map", "structure":
			out += fmt.Sprintf(
				"%sif %s != nil {\n", indent, sourceMemberVar,
			)
			out += fmt.Sprintf(
				"%s\t%s = %s\n",
				indent, targetMemberVar, sourceMemberVar,
			)
			out += fmt.Sprintf("%s}\n", indent)
		}
	}

	return out, nil
}
