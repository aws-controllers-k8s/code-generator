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
	"sort"
	"strings"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/names"
)

// Late-initialization code uses returned result from AWS and sets the values of
// each field in custom resource, if the corresponding field is empty. This is
// especially useful in cases where user doesn't give a value for that field and
// AWS assigns a default one.
// See the following for more details about late-initialization:
// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#late-initialization

// LateInitializeReadOne generates the code that will late initialize the custom
// resource using information fetched from AWS, using ReadOne type output.
func LateInitializeReadOne(
	cfg *ackgenconfig.Config,
	r *model.CRD,
// String representing the name of the variable that we will grab the
// Output shape from. This will likely be "resp" since in the templates
// that call this method, the "source variable" is the response struct
// returned by the aws-sdk-go's SDK API call corresponding to the Operation
	sourceVarName string,
// String representing the name of the variable that we will be **setting**
// with values we get from the Output shape. This will likely be
// "ko.Status" since that is the name of the "target variable" that the
// templates that call this method use.
	targetPath string,
) string {
	out := ""
	path := sourceVarName
	shape := r.Ops.ReadOne.OutputRef.Shape
	if len(shape.MemberRefs) == 1 {
		for key, val := range shape.MemberRefs {
			if val.Shape.Type == "structure" {
				path = fmt.Sprintf("%s.%s", path, key)
				shape = val.Shape
				break
			}
		}
	}
	for _, memberName := range shape.MemberNames() {
		f, found := r.SpecFields[memberName]
		if !found {
			continue
		}
		responsePath := fmt.Sprintf("%s.%s", path, f.Names.Original)
		crPath := fmt.Sprintf("%s%s.%s", targetPath, cfg.PrefixConfig.SpecField, f.Names.Camel)
		out += lateInit(responsePath, crPath, r, f.ShapeRef, 0)
	}
	return strings.TrimSuffix(out, "\n")
}

// LateInitializeReadMany generates the code that will late initialize the custom
// resource using information fetched from AWS, using ReadMany type output.
func LateInitializeReadMany(
	cfg *ackgenconfig.Config,
	r *model.CRD,
// String representing the name of the variable that we will grab the
// Output shape from. This will likely be "resp" since in the templates
// that call this method, the "source variable" is the response struct
// returned by the aws-sdk-go's SDK API call corresponding to the Operation
	sourceVarName string,
// String representing the name of the variable that we will be **setting**
// with values we get from the Output shape. This will likely be
// "ko.Status" since that is the name of the "target variable" that the
// templates that call this method use.
	targetPath string,
) string {
	listFieldName := ""
	var respShapeRef *awssdkmodel.ShapeRef
	for name, ref := range r.Ops.ReadMany.OutputRef.Shape.MemberRefs {
		if ref.Shape.Type == "list"{
			listFieldName = name
			respShapeRef = &ref.Shape.MemberRef
			break
		}
	}
	if respShapeRef == nil {
		panic("could not find a list shaped member in readmany output shape")
	}
	out := fmt.Sprintf("for _, resource := range %s.%s {\n", sourceVarName, listFieldName)
	path := "resource"
	shape := respShapeRef.Shape
	for _, memberName := range shape.MemberNames() {
		f, found := r.SpecFields[memberName]
		if !found {
			continue
		}
		responsePath := fmt.Sprintf("%s.%s", path, f.Names.Original)
		crPath := fmt.Sprintf("%s%s.%s", targetPath, cfg.PrefixConfig.SpecField, f.Names.Camel)
		out += lateInit(responsePath, crPath, r, f.ShapeRef, 0)
	}
	out += fmt.Sprintf("}")
	return out
}

// LateInitializeGetAttributes generates the code that will late initialize the custom
// resource using information fetched from AWS, using GetAttributes type output.
func LateInitializeGetAttributes(
	cfg *ackgenconfig.Config,
	r *model.CRD,
// String representing the name of the variable that we will grab the
// Output shape from. This will likely be "resp" since in the templates
// that call this method, the "source variable" is the response struct
// returned by the aws-sdk-go's SDK API call corresponding to the Operation
	path string,
// String representing the name of the variable that we will be **setting**
// with values we get from the Output shape. This will likely be
// "ko.Status" since that is the name of the "target variable" that the
// templates that call this method use.
	targetPath string,
) string {
	if !r.UnpacksAttributesMap() {
		// This is a bug in the code generation if this occurs...
		msg := fmt.Sprintf(
			"called SetResourceGetAttributes for a resource '%s' that doesn't unpack attributes map",
			r.Ops.GetAttributes.Name,
		)
		panic(msg)
	}
	out := ""
	var fieldList []string
	for name, config := range cfg.ResourceFields(r.Names.Original) {
		if config.IsAttribute && !config.IsReadOnly {
			fieldList = append(fieldList, name)
		}
	}
	sort.Strings(fieldList)
	for _, name := range fieldList {
		n := names.New(name)
		respFieldPath := fmt.Sprintf("%s.Attributes[\"%s\"]", path, n.Original)
		crFieldPath := fmt.Sprintf("%s%s.%s", targetPath, cfg.PrefixConfig.SpecField, n.Camel)
		out += fmt.Sprintf("%s = li.LateInitializeStringPtr(%s, %s)\n", crFieldPath, crFieldPath, respFieldPath)
	}
	return strings.TrimSuffix(out, "\n")
}

func lateInitStruct(responsePath, crPath string, r *model.CRD, str *awssdkmodel.ShapeRef, level int) string {
	out := fmt.Sprintf("if %s != nil {\n", responsePath)
	out += fmt.Sprintf("if %s == nil {\n", crPath)
	out += fmt.Sprintf("%s = &%s{}\n", crPath, GetCRDStructType(str.Shape, r, false))
	out += fmt.Sprintf("}\n")
	// Keys need to be sorted so that we get deterministic output.
	var keys []string
	for key := range str.Shape.MemberRefs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, name := range keys {
		n := names.New(name)
		respFieldPath := fmt.Sprintf("%s.%s", responsePath, n.Original)
		crFieldPath := fmt.Sprintf("%s.%s", crPath, n.Camel)
		out += lateInit(respFieldPath, crFieldPath, r, str.Shape.MemberRefs[name], level+1)
	}
	out += fmt.Sprintf("}\n")
	return out
}

func lateInitMap(responsePath, crPath string, r *model.CRD, str *awssdkmodel.ShapeRef, level int) string {
	out := fmt.Sprintf("if %s != nil {\n", responsePath)
	out += fmt.Sprintf("if %s == nil {\n", crPath)
	out += fmt.Sprintf("%s = %s{}\n", crPath, GetCRDStructType(str.Shape, r, true))
	out += fmt.Sprintf("}\n")

	out += fmt.Sprintf("for key%d := range %s {\n", level, responsePath)
	respFieldPath := fmt.Sprintf("%s[key%d]", responsePath, level)
	crFieldPath := fmt.Sprintf("%s[key%d]", crPath, level)
	out += lateInit(respFieldPath, crFieldPath, r, &str.Shape.ValueRef, level+1)
	out += fmt.Sprintf("}\n")
	out += fmt.Sprintf("}\n")
	return out
}

func lateInitSlice(responsePath, crPath string, r *model.CRD, respShapeRef *awssdkmodel.ShapeRef, level int) string {
	// NOTE(muvaf): Late initialization for slices is not a full-featured late-init.
	// It works only if the slice in CRD is empty. If there is even one element,
	// we don't late-init since slices do not have unique identifier like map/struct.
	// We're playing conservative and not rely on index so that we don't risk
	// overriding desired state with what's observed.
	// TODO(muvaf): Kubernetes Apply uses comment markers on fields to specify an
	// identifier (`name` most of the time). If we can do such thing, we can implement
	// full late-init for slice elements as well, just like maps.
	out := fmt.Sprintf(
		"if len(%s) != 0 && len(%s) == 0 {\n", responsePath, crPath,
	)
	out += fmt.Sprintf("%s = make([]%s, len(%s))\n", crPath, GetCRDStructType(respShapeRef.Shape.MemberRef.Shape, r, true), responsePath)
	out += fmt.Sprintf("for i%d := range %s {\n", level, responsePath)
	respFieldPath := fmt.Sprintf("%s[i%d]", responsePath, level)
	crFieldPath := fmt.Sprintf("%s[i%d]", crPath, level)
	out += lateInit(respFieldPath, crFieldPath, r, &respShapeRef.Shape.MemberRef, level+1)
	out += fmt.Sprintf("}\n")
	out += fmt.Sprintf("}\n")
	return out
}

func lateInit(responsePath, crPath string, r *model.CRD, str *awssdkmodel.ShapeRef, level int) string {
	switch str.Shape.Type {
	case "string":
		return fmt.Sprintf("%s = li.LateInitializeStringPtr(%s, %s)\n", crPath, crPath, responsePath)
	case "long", "integer":
		return fmt.Sprintf("%s = li.LateInitializeInt64Ptr(%s, %s)\n", crPath, crPath, responsePath)
	case "boolean":
		return fmt.Sprintf("%s = li.LateInitializeBoolPtr(%s, %s)\n", crPath, crPath, responsePath)
	case "timestamp":
		return fmt.Sprintf("%s = li.LateInitializeTimePtr(%s, %s)\n", crPath, crPath, responsePath)
	case "list":
		return lateInitSlice(responsePath, crPath, r, str, level)
	case "structure":
		return lateInitStruct(responsePath, crPath, r, str, level)
	case "map":
		return lateInitMap(responsePath, crPath, r, str, level)
	case "double":
		return fmt.Sprintf("// Please handle %s manually. The check for double type is not implemented yet.\n", crPath)
	case "blob":
		return fmt.Sprintf("// Please handle %s manually. The check for blob type is not implemented yet.\n", crPath)
	default:
		panic(fmt.Sprintf("unknown shape type %s", str.Shape.Type))
	}
}

func GetCRDStructType(s *awssdkmodel.Shape, r *model.CRD, keepPointer bool) string {
	goType := model.ReplacePkgName(s.GoTypeWithPkgName(), r.SDKAPIPackageName(), "svcapitypes", keepPointer)
	if !strings.Contains(goType, ".") {
		return goType
	}
	goTypeNoPkg := strings.Split(goType, ".")[1]
	goPkg := strings.Split(goType, ".")[0]
	if r.TypeRenames()[goTypeNoPkg] != "" {
		goTypeNoPkg = r.TypeRenames()[goTypeNoPkg]
	} else {
		goTypeNoPkg = names.New(goTypeNoPkg).Camel
	}
	return fmt.Sprintf("%s.%s", goPkg, goTypeNoPkg)
}
