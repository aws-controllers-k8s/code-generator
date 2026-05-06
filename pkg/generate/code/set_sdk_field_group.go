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
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// SetSDKFieldGroup returns the Go code that sets an SDK Input shape's member
// fields from a CRD's fields for a specific field-group operation. Only
// members that are identifier or payload fields of the field group are set.
//
// This is the field-group analog of SetSDK. It generates code like:
//
//	res.RegistryId = r.ko.Spec.RegistryID
//	res.RepositoryName = r.ko.Spec.RepositoryName
//	if r.ko.Spec.ImageScanningConfiguration != nil {
//	    f0 := &svcsdk.ImageScanningConfiguration{}
//	    ...
//	    res.ImageScanningConfiguration = f0
//	}
func SetSDKFieldGroup(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	fg *model.FieldGroupOperation,
	// String representing the name of the variable that we will grab the
	// Input values from (typically "r.ko")
	sourceVarName string,
	// String representing the name of the variable that we will be setting
	// (typically "res")
	targetVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) (string, error) {
	op := fg.Operation
	if op == nil || op.InputRef.Shape == nil {
		return "", nil
	}

	inputShape := op.InputRef.Shape

	// Build a set of field names that belong to this field group (both
	// identifiers and payload) for fast membership checks.
	fgFieldNames := fieldGroupFieldNameSet(fg)

	out := "\n"
	indent := strings.Repeat("\t", indentLevel)

	for memberIndex, memberName := range inputShape.MemberNames() {
		// Resolve the CRD field name for this Input member
		fieldName := cfg.GetResourceFieldName(
			r.Names.Original, op.ExportedName, memberName,
		)

		// Check if this member's field is part of the field group
		f, inSpec := r.SpecFields[fieldName]
		if !inSpec || !fgFieldNames[f.Names.Camel] {
			continue
		}

		sourceAdaptedVarName := sourceVarName + cfg.PrefixConfig.SpecField + "." + f.Names.Camel
		sourceFieldPath := f.Names.Camel

		setCfg := f.GetSetterConfig(model.OpTypeUpdate)
		if setCfg != nil && setCfg.IgnoreSDKSetter() {
			continue
		}

		memberShapeRef := inputShape.MemberRefs[memberName]
		memberShape := memberShapeRef.Shape
		if memberShape.RealType == "union" {
			memberShape.Type = "union"
		}

		out += fmt.Sprintf(
			"%sif %s != nil {\n", indent, sourceAdaptedVarName,
		)

		switch memberShape.Type {
		case "list", "structure", "map", "union":
			adaptiveCollection := setSDKAdaptiveResourceCollection(
				memberShape, targetVarName, memberName,
				sourceAdaptedVarName, indent, r.IsSecretField(memberName),
			)
			out += adaptiveCollection
			if adaptiveCollection != "" {
				break
			}
			{
				memberVarName := fmt.Sprintf("f%d", memberIndex)
				out += varEmptyConstructorSDKType(
					cfg, r, memberVarName, memberShape, indentLevel+1,
				)
				containerOut, err := setSDKForContainer(
					cfg, r, memberName, memberVarName,
					sourceFieldPath, sourceAdaptedVarName,
					memberShapeRef, false, model.OpTypeUpdate, indentLevel+1,
				)
				if err != nil {
					return "", err
				}
				out += containerOut
				out += setSDKForScalar(
					memberName, targetVarName, inputShape.Type,
					sourceFieldPath, memberVarName, false,
					memberShapeRef, indentLevel+1,
				)
			}
		default:
			if r.IsSecretField(memberName) {
				out += setSDKForSecret(
					cfg, r, memberName, targetVarName,
					sourceAdaptedVarName, indentLevel,
				)
			} else {
				out += setSDKForScalar(
					memberName, targetVarName, inputShape.Type,
					sourceFieldPath, sourceAdaptedVarName, false,
					memberShapeRef, indentLevel+1,
				)
			}
		}
		if memberShape.RealType == "union" {
			memberShape.Type = "structure"
		}
		out += fmt.Sprintf("%s}\n", indent)
	}

	return out, nil
}

// fieldGroupFieldNameSet returns a set of CRD field names (Camel) that belong
// to the given field group operation (identifiers + payload).
func fieldGroupFieldNameSet(fg *model.FieldGroupOperation) map[string]bool {
	s := make(map[string]bool, len(fg.IdentifierFields)+len(fg.PayloadFields))
	for _, f := range fg.IdentifierFields {
		s[f.Names.Camel] = true
	}
	for _, f := range fg.PayloadFields {
		s[f.Names.Camel] = true
	}
	return s
}
