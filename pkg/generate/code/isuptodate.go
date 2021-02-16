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
//

package code

import (
	"fmt"
	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/names"
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"sort"
	"strings"
)

func IsUpToDateReadOne(
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
	if r.Ops.Update == nil {
		return ""
	}
	out := ""
	path := sourceVarName
	shape := r.Ops.Update.InputRef.Shape
	pathShape := r.Ops.ReadOne.OutputRef.Shape
	if len(pathShape.MemberRefs) == 1 {
		for key, val := range pathShape.MemberRefs {
			if val.Shape.Type == "structure" {
				path = fmt.Sprintf("%s.%s", path, key)
				pathShape = val.Shape
				break
			}
			panic("there has to be a structure field in unwrapped readone output shape")
		}
	}
	for _, memberName := range shape.MemberNames() {
		f, found := r.SpecFields[memberName]
		if !found {
			continue
		}
		responsePath := fmt.Sprintf("%s.%s", path, f.Names.Original)
		crPath := fmt.Sprintf("%s%s.%s", targetPath, cfg.PrefixConfig.SpecField, f.Names.Camel)
		if !isInDescribe(memberName, pathShape) {
			out += fmt.Sprintf("// Please handle %s manually.\n", crPath)
			continue
		}
		out += upToDate(responsePath, crPath, r, f.ShapeRef, 0)
	}
	return strings.TrimSuffix(out, "\n")
}

func IsUpToDateReadMany(
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
	if r.Ops.Update == nil {
		return ""
	}
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
	path := "resource"
	out := fmt.Sprintf("for _, %s := range %s.%s {\n", path, sourceVarName, listFieldName)
	shape := r.Ops.Update.InputRef.Shape
	pathShape := respShapeRef.Shape
	for _, memberName := range shape.MemberNames() {
		f, found := r.SpecFields[memberName]
		if !found {
			continue
		}
		responsePath := fmt.Sprintf("%s.%s", path, f.Names.Original)
		crPath := fmt.Sprintf("%s%s.%s", targetPath, cfg.PrefixConfig.SpecField, f.Names.Camel)
		if !isInDescribe(memberName, pathShape) {
			out += fmt.Sprintf("// Please handle %s manually.\n", crPath)
			continue
		}
		out += upToDate(responsePath, crPath, r, f.ShapeRef, 0)
	}
	out += fmt.Sprintf("}")
	return out
}

func isInDescribe(shapeName string, describeShape *awssdkmodel.Shape) bool {
	for _, memberName := range describeShape.MemberNames() {
		if shapeName == memberName {
			return true
		}
	}
	return false
}

func upToDateStruct(responsePath, crPath string, r *model.CRD, str *awssdkmodel.ShapeRef, level int) string {
	out := fmt.Sprintf("if (%s != nil && %s == nil) || (%s == nil && %s != nil) {\n return false \n}\n", responsePath, crPath, responsePath, crPath)
	out += fmt.Sprintf("if %s != nil && %s != nil {\n", responsePath, crPath)
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
		out += upToDate(respFieldPath, crFieldPath, r, str.Shape.MemberRefs[name], level+1)
	}
	out += fmt.Sprintf("}\n")
	return out
}

func upToDateMap(responsePath, crPath string, r *model.CRD, str *awssdkmodel.ShapeRef, level int) string {
	out := fmt.Sprintf(
		"if len(%s) != len(%s) {\n return false\n}\n", responsePath, crPath,
	)

	out += fmt.Sprintf("for key%d := range %s {\n", level, responsePath)
	respFieldPath := fmt.Sprintf("%s[key%d]", responsePath, level)
	crFieldPath := fmt.Sprintf("%s[key%d]", crPath, level)
	out += upToDate(respFieldPath, crFieldPath, r, &str.Shape.ValueRef, level+1)
	out += fmt.Sprintf("}\n")
	return out
}

func upToDateSlice(responsePath, crPath string, r *model.CRD, respShapeRef *awssdkmodel.ShapeRef, level int) string {
	out := fmt.Sprintf(
		"if len(%s) != len(%s) {\n return false\n}\n", responsePath, crPath,
	)
	//out += fmt.Sprintf("%s = make([]%s, len(%s))\n", crPath, GetCorrespondingCRDType(respShapeRef.Shape.MemberRef.Shape, r, true), responsePath)
	out += fmt.Sprintf("for i%d := range %s {\n", level, responsePath)
	respFieldPath := fmt.Sprintf("%s[i%d]", responsePath, level)
	crFieldPath := fmt.Sprintf("%s[i%d]", crPath, level)
	out += upToDate(respFieldPath, crFieldPath, r, &respShapeRef.Shape.MemberRef, level+1)
	out += fmt.Sprintf("}\n")
	return out
}

func upToDate(responsePath, crPath string, r *model.CRD, str *awssdkmodel.ShapeRef, level int) string {
	switch str.Shape.Type {
	case "string":
		return fmt.Sprintf("if awsclients.StringValue(%s) != awsclients.StringValue(%s) {\n return false\n}\n", crPath, responsePath)
	case "long", "integer":
		return fmt.Sprintf("if awsclients.Int64Value(%s) != awsclients.Int64Value(%s) {\n return false\n}\n", crPath, responsePath)
	case "boolean":
		return fmt.Sprintf("if awsclients.BoolValue(%s) != awsclients.BoolValue(%s) {\n return false\n}\n", crPath, responsePath)
	case "list":
		return upToDateSlice(responsePath, crPath, r, str, level)
	case "structure":
		return upToDateStruct(responsePath, crPath, r, str, level)
	case "map":
		return upToDateMap(responsePath, crPath, r, str, level)
	case "timestamp":
		return fmt.Sprintf("// Please handle %s manually. The check for timestamp type is not implemented yet.\n", crPath)
	case "double":
		return fmt.Sprintf("// Please handle %s manually. The check for double type is not implemented yet.\n", crPath)
	case "blob":
		return fmt.Sprintf("// Please handle %s manually. The check for blob type is not implemented yet.\n", crPath)
	default:
		panic(fmt.Sprintf("unknown shape type %s", str.Shape.Type))
	}
}
