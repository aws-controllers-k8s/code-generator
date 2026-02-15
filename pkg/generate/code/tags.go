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

// GoCodeToACKTags returns Go code that converts Resource Tags
// to ACK Tags. If Resource Tags field is of type list, we
// also maintain and return the order of the list as a []string
// 
//
//
// Sample output:
//
//	 for _, k := range keyOrder {
//	 	v, ok := tags[k]
//	 	if ok {
//	 		tag := svcapitypes.Tag{Key: &k, Value: &v}
//	 		result = append(result, &tag)
//	 		delete(tags, k)
//	 	}
//	 }
//	 for k, v := range tags {
//	 	tag := svcapitypes.Tag{Key: &k, Value: &v}
//	 	result = append(result, &tag)
//	 }
func GoCodeConvertToACKTags(r *model.CRD, sourceVarName string, targetVarName string, keyOrderVarName string, indentLevel int) (string, error) {

	out := "\n"
	indent := strings.Repeat("\t", indentLevel)
	tagField, err := r.GetTagField()
	if err != nil {
		return "", fmt.Errorf("resource %q: does not have tags — ignore in generator.yaml: %w", r.Names.Original, err)
	}

	if tagField == nil {
		return "", nil
	}

	tagFieldShapeType := tagField.ShapeRef.Shape.Type
	keyMemberName := r.GetTagKeyMemberName()
	valueMemberName := r.GetTagValueMemberName()

	out += fmt.Sprintf("%sif len(%s) == 0 {\n", indent, sourceVarName)
	out += fmt.Sprintf("%s\treturn %s, %s\n", indent, targetVarName, keyOrderVarName)
	out += fmt.Sprintf("%s}\n", indent)

	switch tagFieldShapeType {
	case "list":
		out += fmt.Sprintf("%sfor _, t := range %s {\n", indent, sourceVarName)
		out += fmt.Sprintf("%s\tif t.%s != nil {\n", indent, keyMemberName)
		out += fmt.Sprintf("%s\t\t%s = append(%s, *t.%s)\n", indent, keyOrderVarName, keyOrderVarName, keyMemberName)
		out += fmt.Sprintf("%s\t\tif t.%s != nil {\n", indent, valueMemberName)
		out += fmt.Sprintf("%s\t\t\t%s[*t.%s] = *t.%s\n", indent, targetVarName, keyMemberName, valueMemberName)
		out += fmt.Sprintf("%s\t\t} else {\n", indent)
		out += fmt.Sprintf("%s\t\t\t%s[*t.%s] = \"\"\n", indent, targetVarName, keyMemberName)
		out += fmt.Sprintf("%s\t\t}\n", indent)
		out += fmt.Sprintf("%s\t}\n", indent)
		out += fmt.Sprintf("%s}\n", indent)

	case "map":
		out += fmt.Sprintf("%sfor k, v := range %s {\n", indent, sourceVarName)
		out += fmt.Sprintf("%s\tif v == nil {\n", indent)
		out += fmt.Sprintf("%s\t\t%s[k] = \"\"\n", indent, targetVarName)
		out += fmt.Sprintf("%s\t} else {\n", indent)
		out += fmt.Sprintf("%s\t\t%s[k] = *v\n", indent, targetVarName)
		out += fmt.Sprintf("%s\t}\n", indent)
		out += fmt.Sprintf("%s}\n", indent)
	default:
		return "", fmt.Errorf("resource %q: tag type can only be a list or a map, got %q", r.Names.Original, tagFieldShapeType)
	}

	return out, nil
}

// GoCodeFromACKTags returns Go code that converts ACKTags
// to the Resource Tag shape type. Tag fields can only be
// maps or lists of Tag Go type. If Tag field is a list,
// when converting from ACK Tags, we try to preserve the 
// original order
// 
//
//
// Sample output:
//
//	 for _, k := range keyOrder {
//	 	v, ok := tags[k]
//	 	if ok {
//	 		tag := svcapitypes.Tag{Key: &k, Value: &v}
//	 		result = append(result, &tag)
//	 		delete(tags, k)
//	 	}
//	 }
//	 for k, v := range tags {
//	 	tag := svcapitypes.Tag{Key: &k, Value: &v}
//	 	result = append(result, &tag)
//	 }
func GoCodeFromACKTags(r *model.CRD, tagsSourceVarName string, orderVarName string, targetVarName string, indentLevel int) (string, error) {
	out := "\n"
	indent := strings.Repeat("\t", indentLevel)
	tagField, _ := r.GetTagField()

	if tagField == nil {
		return "", nil
	}

	tagFieldShapeType := tagField.ShapeRef.Shape.Type
	tagFieldGoType := tagField.GoTypeElem
	keyMemberName := r.GetTagKeyMemberName()
	valueMemberName := r.GetTagValueMemberName()

	switch tagFieldShapeType {
	case "list":
		out += fmt.Sprintf("%sfor _, k := range %s {\n", indent, orderVarName)
		out += fmt.Sprintf("%s\tv, ok := %s[k]\n", indent, tagsSourceVarName)
		out += fmt.Sprintf("%s\tif ok {\n", indent)
		out += fmt.Sprintf("%s\t\ttag := svcapitypes.%s{%s: &k, %s: &v}\n", indent, tagFieldGoType, keyMemberName, valueMemberName)
		out += fmt.Sprintf("%s\t\t%s = append(%s, &tag)\n", indent, targetVarName, targetVarName)
		out += fmt.Sprintf("%s\t\tdelete(%s, k)\n", indent, tagsSourceVarName)
		out += fmt.Sprintf("%s\t}\n", indent)
		out += fmt.Sprintf("%s}\n", indent)
	case "map":
		out += fmt.Sprintf("%s_ = %s\n", indent, orderVarName)
	default:
		return "", fmt.Errorf("resource %q: tag type can only be a list or a map, got %q", r.Names.Original, tagFieldShapeType)
	}

	out += fmt.Sprintf("%sfor k, v := range %s {\n", indent, tagsSourceVarName)
	switch tagFieldShapeType {
	case "list":
		out += fmt.Sprintf("%s\ttag := svcapitypes.%s{%s: &k, %s: &v}\n", indent, tagFieldGoType, keyMemberName, valueMemberName)
		out += fmt.Sprintf("%s\t%s = append(%s, &tag)\n", indent, targetVarName, targetVarName)
	case "map":
		out += fmt.Sprintf("%s\t%s[k] = &v\n", indent, targetVarName)
	}
	out += fmt.Sprintf("%s}\n", indent)

	return out, nil
}
