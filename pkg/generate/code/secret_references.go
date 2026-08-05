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

	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// SecretReferences returns Go code that collects all SecretKeyReferences from
// a resource's Spec and returns them as a slice. This is used to implement the
// SecretReferenceResolver interface.
//
// Sample output:
//
//	var refs []*ackv1alpha1.SecretKeyReference
//	if r.ko.Spec.MasterUserPassword != nil {
//	    refs = append(refs, r.ko.Spec.MasterUserPassword)
//	}
//	if r.ko.Spec.Users != nil {
//	    for _, elem := range r.ko.Spec.Users {
//	        if elem != nil && elem.Password != nil {
//	            refs = append(refs, elem.Password)
//	        }
//	    }
//	}
//	return refs
func SecretReferences(
	r *model.CRD,
	// koVarName is the variable name for the resource's ko object
	koVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	out += fmt.Sprintf("%svar refs []*ackv1alpha1.SecretKeyReference\n", indent)

	for _, field := range r.SpecFields {
		out += secretReferencesForField(field, koVarName+".Spec", indentLevel)
	}

	out += fmt.Sprintf("%sreturn refs\n", indent)
	return out
}

func secretReferencesForField(
	field *model.Field,
	parentAccessor string,
	indentLevel int,
) string {
	if field.FieldConfig != nil && field.FieldConfig.IsSecret {
		return secretReferenceDirectField(field, parentAccessor, indentLevel)
	}

	if field.MemberFields != nil {
		return secretReferencesForStruct(field, parentAccessor, indentLevel)
	}

	return ""
}

// secretReferenceDirectField handles fields that are directly a
// SecretKeyReference or a slice of SecretKeyReferences.
func secretReferenceDirectField(
	field *model.Field,
	parentAccessor string,
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	accessor := fmt.Sprintf("%s.%s", parentAccessor, field.Names.Camel)

	switch {
	case field.GoType == "*ackv1alpha1.SecretKeyReference":
		out += fmt.Sprintf("%sif %s != nil {\n", indent, accessor)
		out += fmt.Sprintf("%s\trefs = append(refs, %s)\n", indent, accessor)
		out += fmt.Sprintf("%s}\n", indent)

	case field.GoType == "[]*ackv1alpha1.SecretKeyReference":
		out += fmt.Sprintf("%sif %s != nil {\n", indent, accessor)
		out += fmt.Sprintf("%s\trefs = append(refs, %s...)\n", indent, accessor)
		out += fmt.Sprintf("%s}\n", indent)

	case field.GoType == "map[string]*ackv1alpha1.SecretKeyReference":
		out += fmt.Sprintf("%sfor _, v := range %s {\n", indent, accessor)
		out += fmt.Sprintf("%s\tif v != nil {\n", indent)
		out += fmt.Sprintf("%s\t\trefs = append(refs, v)\n", indent)
		out += fmt.Sprintf("%s\t}\n", indent)
		out += fmt.Sprintf("%s}\n", indent)
	}

	return out
}

// secretReferencesForStruct handles struct fields that may contain secret
// fields as members — either the struct itself or when wrapped in a slice.
func secretReferencesForStruct(
	field *model.Field,
	parentAccessor string,
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	accessor := fmt.Sprintf("%s.%s", parentAccessor, field.Names.Camel)

	hasSecretMember := false
	for _, mf := range field.MemberFields {
		if mf.FieldConfig != nil && mf.FieldConfig.IsSecret {
			hasSecretMember = true
			break
		}
	}
	if !hasSecretMember {
		return ""
	}

	isSlice := strings.HasPrefix(field.GoType, "[]")

	if isSlice {
		out += fmt.Sprintf("%sif %s != nil {\n", indent, accessor)
		out += fmt.Sprintf("%s\tfor _, elem := range %s {\n", indent, accessor)
		out += fmt.Sprintf("%s\t\tif elem == nil {\n", indent)
		out += fmt.Sprintf("%s\t\t\tcontinue\n", indent)
		out += fmt.Sprintf("%s\t\t}\n", indent)
		for _, mf := range field.MemberFields {
			if mf.FieldConfig != nil && mf.FieldConfig.IsSecret {
				memberAccessor := fmt.Sprintf("elem.%s", mf.Names.Camel)
				switch {
				case mf.GoType == "*ackv1alpha1.SecretKeyReference":
					out += fmt.Sprintf("%s\t\tif %s != nil {\n", indent, memberAccessor)
					out += fmt.Sprintf("%s\t\t\trefs = append(refs, %s)\n", indent, memberAccessor)
					out += fmt.Sprintf("%s\t\t}\n", indent)
				case mf.GoType == "[]*ackv1alpha1.SecretKeyReference":
					out += fmt.Sprintf("%s\t\tif %s != nil {\n", indent, memberAccessor)
					out += fmt.Sprintf("%s\t\t\trefs = append(refs, %s...)\n", indent, memberAccessor)
					out += fmt.Sprintf("%s\t\t}\n", indent)
				case mf.GoType == "map[string]*ackv1alpha1.SecretKeyReference":
					out += fmt.Sprintf("%s\t\tfor _, v := range %s {\n", indent, memberAccessor)
					out += fmt.Sprintf("%s\t\t\tif v != nil {\n", indent)
					out += fmt.Sprintf("%s\t\t\t\trefs = append(refs, v)\n", indent)
					out += fmt.Sprintf("%s\t\t\t}\n", indent)
					out += fmt.Sprintf("%s\t\t}\n", indent)
				}
			}
		}
		out += fmt.Sprintf("%s\t}\n", indent)
		out += fmt.Sprintf("%s}\n", indent)
	} else {
		out += fmt.Sprintf("%sif %s != nil {\n", indent, accessor)
		for _, mf := range field.MemberFields {
			if mf.FieldConfig != nil && mf.FieldConfig.IsSecret {
				memberAccessor := fmt.Sprintf("%s.%s", accessor, mf.Names.Camel)
				switch {
				case mf.GoType == "*ackv1alpha1.SecretKeyReference":
					out += fmt.Sprintf("%s\tif %s != nil {\n", indent, memberAccessor)
					out += fmt.Sprintf("%s\t\trefs = append(refs, %s)\n", indent, memberAccessor)
					out += fmt.Sprintf("%s\t}\n", indent)
				case mf.GoType == "[]*ackv1alpha1.SecretKeyReference":
					out += fmt.Sprintf("%s\tif %s != nil {\n", indent, memberAccessor)
					out += fmt.Sprintf("%s\t\trefs = append(refs, %s...)\n", indent, memberAccessor)
					out += fmt.Sprintf("%s\t}\n", indent)
				case mf.GoType == "map[string]*ackv1alpha1.SecretKeyReference":
					out += fmt.Sprintf("%s\tfor _, v := range %s {\n", indent, memberAccessor)
					out += fmt.Sprintf("%s\t\tif v != nil {\n", indent)
					out += fmt.Sprintf("%s\t\t\trefs = append(refs, v)\n", indent)
					out += fmt.Sprintf("%s\t\t}\n", indent)
					out += fmt.Sprintf("%s\t}\n", indent)
				}
			}
		}
		out += fmt.Sprintf("%s}\n", indent)
	}

	return out
}
