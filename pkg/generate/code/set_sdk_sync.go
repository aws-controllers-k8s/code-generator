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

// SetSyncInput generates Go code that sets the SDK Input struct fields for a
// sync-list operation's add or remove call. It sets identifier fields from the
// resource and the per-item scalar field from the loop variable.
//
// Generated code looks like:
//
//	input.RoleName = desired.ko.Spec.Name
//	input.PolicyArn = item
func SetSyncInput(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	fg *model.FieldGroupOperation,
	// Name of the resource variable (e.g., "desired")
	resourceVarName string,
	// Name of the SDK input variable (e.g., "input")
	targetVarName string,
	// Name of the loop item variable (e.g., "item")
	itemVarName string,
	indentLevel int,
) string {
	op := fg.Operation
	if op == nil || op.InputRef.Shape == nil {
		return ""
	}
	inputShape := op.InputRef.Shape
	identifierMemberNames := buildIdentifierMemberNamesForOp(cfg, r, fg, op.ExportedName)

	out := ""
	indent := strings.Repeat("\t", indentLevel)

	for _, memberName := range inputShape.MemberNames() {
		if identifierMemberNames[memberName] {
			// Identifier field — set from the resource
			fieldName := cfg.GetResourceFieldName(
				r.Names.Original, op.ExportedName, memberName,
			)
			f, ok := r.SpecFields[fieldName]
			if !ok {
				continue
			}
			out += fmt.Sprintf(
				"%s%s.%s = %s.ko.Spec.%s\n",
				indent, targetVarName, memberName,
				resourceVarName, f.Names.Camel,
			)
		} else {
			// Per-item field — set from the loop variable
			out += fmt.Sprintf(
				"%s%s.%s = %s\n",
				indent, targetVarName, memberName, itemVarName,
			)
		}
	}
	return out
}

// SetSyncMapAddInput generates Go code that sets the SDK Input struct fields
// for a sync-map operation's add/update call. It sets identifier fields from
// the resource, the map key field, and the map value field from loop variables.
//
// Generated code looks like:
//
//	input.RoleName = desired.ko.Spec.Name
//	input.PolicyName = &key
//	input.PolicyDocument = val
func SetSyncMapAddInput(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	fg *model.FieldGroupOperation,
	resourceVarName string,
	targetVarName string,
	keyVarName string,
	valVarName string,
	indentLevel int,
) string {
	op := fg.Operation
	if op == nil || op.InputRef.Shape == nil {
		return ""
	}
	inputShape := op.InputRef.Shape
	identifierMemberNames := buildIdentifierMemberNamesForOp(cfg, r, fg, op.ExportedName)

	out := ""
	indent := strings.Repeat("\t", indentLevel)

	// Collect non-identifier members — these are the key and value fields.
	// By convention, the first non-identifier string field is the key and
	// the second is the value.
	var nonIDMembers []string
	for _, memberName := range inputShape.MemberNames() {
		if !identifierMemberNames[memberName] {
			nonIDMembers = append(nonIDMembers, memberName)
		}
	}

	for _, memberName := range inputShape.MemberNames() {
		if identifierMemberNames[memberName] {
			fieldName := cfg.GetResourceFieldName(
				r.Names.Original, op.ExportedName, memberName,
			)
			f, ok := r.SpecFields[fieldName]
			if !ok {
				continue
			}
			out += fmt.Sprintf(
				"%s%s.%s = %s.ko.Spec.%s\n",
				indent, targetVarName, memberName,
				resourceVarName, f.Names.Camel,
			)
		} else {
			// Determine if this is the key or value field
			if len(nonIDMembers) > 0 && memberName == nonIDMembers[0] {
				out += fmt.Sprintf(
					"%s%s.%s = &%s\n",
					indent, targetVarName, memberName, keyVarName,
				)
			} else {
				out += fmt.Sprintf(
					"%s%s.%s = %s\n",
					indent, targetVarName, memberName, valVarName,
				)
			}
		}
	}
	return out
}

// SetSyncMapRemoveInput generates Go code that sets the SDK Input struct fields
// for a sync-map operation's remove call. It sets identifier fields from the
// resource and the map key field from the loop variable.
//
// Generated code looks like:
//
//	input.RoleName = desired.ko.Spec.Name
//	input.PolicyName = &key
func SetSyncMapRemoveInput(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	fg *model.FieldGroupOperation,
	resourceVarName string,
	targetVarName string,
	keyVarName string,
	indentLevel int,
) string {
	removeOp := fg.RemoveOperation
	if removeOp == nil || removeOp.InputRef.Shape == nil {
		return ""
	}
	inputShape := removeOp.InputRef.Shape
	identifierMemberNames := buildIdentifierMemberNamesForOp(cfg, r, fg, removeOp.ExportedName)

	out := ""
	indent := strings.Repeat("\t", indentLevel)

	for _, memberName := range inputShape.MemberNames() {
		if identifierMemberNames[memberName] {
			fieldName := cfg.GetResourceFieldName(
				r.Names.Original, removeOp.ExportedName, memberName,
			)
			f, ok := r.SpecFields[fieldName]
			if !ok {
				continue
			}
			out += fmt.Sprintf(
				"%s%s.%s = %s.ko.Spec.%s\n",
				indent, targetVarName, memberName,
				resourceVarName, f.Names.Camel,
			)
		} else {
			out += fmt.Sprintf(
				"%s%s.%s = &%s\n",
				indent, targetVarName, memberName, keyVarName,
			)
		}
	}
	return out
}

// buildIdentifierMemberNamesForOp returns a set of SDK Input shape member
// names that correspond to identifier fields for the given operation. It
// reverses the rename mapping to get the original SDK member name.
func buildIdentifierMemberNamesForOp(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	fg *model.FieldGroupOperation,
	opID string,
) map[string]bool {
	s := make(map[string]bool, len(fg.IdentifierFields))
	for _, f := range fg.IdentifierFields {
		memberName := cfg.GetOriginalMemberName(
			r.Names.Original, opID, f.Names.Camel,
		)
		s[memberName] = true
	}
	return s
}
