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

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
)

// ResourceIsSynced returns the Go code that verifies whether a resource is synced or
// not. This code is generated using ack-generate configuration.
// See ack-generate/pkg/config.SyncedConfiguration.
//
//  Sample output:
//
//  	candidates0 := []string{"AVAILABLE", "ACTIVE"}
//  	if !ackutil.InStrings(*r.ko.Status.TableStatus, candidates0) {
//  		return false, nil
//  	}
//  	if r.ko.Spec.ProvisionedThroughput == nil {
//  		return false, nil
//  	}
//  	if r.ko.Spec.ProvisionedThroughput.ReadCapacityUnits == nil {
//  		return false, nil
//  	}
//  	candidates1 := []int{0, 10}
//  	if !ackutil.InStrings(*r.ko.Spec.ProvisionedThroughput.ReadCapacityUnits, candidates1) {
//  		return false, nil
//  	}
//  	candidates2 := []int{0}
//  	if !ackutil.InStrings(*r.ko.Status.ItemCount, candidates2) {
//  		return false, nil
//  	}
func ResourceIsSynced(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// resource variable name
	resVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := "\n"
	resConfig := cfg.GetResourceConfig(r.Names.Original)
	if resConfig == nil || resConfig.Synced == nil || len(resConfig.Synced.When) == 0 {
		return out
	}

	for _, condCfg := range resConfig.Synced.When {
		if condCfg.Path == nil || *condCfg.Path == "" {
			panic("Received an empty sync condition path. 'SyncCondition.Path' must be provided.")
		}
		if len(condCfg.In) == 0 {
			panic("'SyncCondition.In' must be provided.")
		}
		fp := fieldpath.FromString(*condCfg.Path)
		field, err := getTopLevelField(r, *condCfg.Path)
		if err != nil {
			msg := fmt.Sprintf("cannot find top level field of path '%s': %v", *condCfg.Path, err)
			panic(msg)
		}
		candidatesVarName := fmt.Sprintf("%sCandidates", field.Names.CamelLower)
		if fp.Size() == 2 {
			out += scalarFieldEqual(resVarName, candidatesVarName, field.ShapeRef.GoTypeElem(), condCfg)
		} else {
			out += fieldPathSafeEqual(resVarName, candidatesVarName, field, condCfg)
		}
	}

	return out
}

func getTopLevelField(r *model.CRD, fieldPath string) (*model.Field, error) {
	fp := fieldpath.FromString(fieldPath)
	if fp.Size() < 2 {
		return nil, fmt.Errorf("fieldPath must contain at least two elements, received: %s", fieldPath)
	}

	head := fp.PopFront()
	fieldName := fp.PopFront()
	switch head {
	case "Spec":
		field, ok := r.Fields[fieldName]
		if !ok {
			return nil, fmt.Errorf("field not found in Spec: %v", fieldName)
		}
		return field, nil
	case "Status":
		field, ok := r.Fields[fieldName]
		if !ok {
			return nil, fmt.Errorf("field not found in Status: %v", fieldName)
		}
		return field, nil
	default:
		return nil, fmt.Errorf("fieldPath must start with 'Spec' or 'Status', received: %v", head)
	}
}

// scalarFieldEqual returns Go code that compares a scalar field to a given set of values.
func scalarFieldEqual(
	resVarName string,
	candidatesVarName string,
	goType string,
	condCfg ackgenconfig.SyncedCondition,
) string {
	out := ""
	fieldPath := fmt.Sprintf("%s.%s", resVarName, *condCfg.Path)
	// if r.ko.Status.Status == nil
	out += fmt.Sprintf("\tif %s == nil {\n", fieldPath)
	// return false, nil
	out += "\t\treturn false, nil\n"
	// }
	out += "\t}\n"
	valuesSlice := ""
	switch goType {
	case "string":
		// []string{"AVAILABLE", "ACTIVE"}
		valuesSlice = fmt.Sprintf("[]string{\"%s\"}", strings.Join(condCfg.In, "\", \""))
	case "int64", "PositiveLongObject", "Long":
		// []int64{1, 2}
		valuesSlice = fmt.Sprintf("[]int{%s}", strings.Join(condCfg.In, ", "))
	case "bool":
		// []bool{false}
		valuesSlice = fmt.Sprintf("[]bool{%s}", condCfg.In)
	default:
		panic("not supported type " + goType)
	}

	// candidates1 := []string{"AVAILABLE", "ACTIVE"}
	out += fmt.Sprintf(
		"\t%s := %v\n",
		candidatesVarName,
		valuesSlice,
	)
	// 	if !ackutil.InStrings(*r.ko.Status.State, candidates1) {
	out += fmt.Sprintf(
		"\tif !ackutil.InStrings(*%s, %s) {\n",
		fieldPath,
		candidatesVarName,
	)

	// return false, nil
	out += "\t\treturn false, nil\n"
	// }
	out += "\t}\n"
	return out
}

// fieldPathSafeEqual returns go code that safely compares a resource field to value
func fieldPathSafeEqual(
	resVarName string,
	candidatesVarName string,
	field *model.Field,
	condCfg ackgenconfig.SyncedCondition,
) string {
	out := ""
	rootPath := fmt.Sprintf("%s.%s", resVarName, strings.Split(*condCfg.Path, ".")[0])
	knownShapesPath := strings.Join(strings.Split(*condCfg.Path, ".")[1:], ".")

	fp := fieldpath.FromString(knownShapesPath)
	shapes := fp.IterShapeRefs(field.ShapeRef)

	subFieldPath := rootPath
	for index, shape := range shapes {
		if index == len(shapes)-1 {
			// We would check for nil in scalarFieldEqual method so no need to loop anymore
			break
		} else {
			subFieldPath += "." + shape.Shape.ShapeName
		}
		// if r.ko.Spec.ProvisionedThroughput == nil
		out += fmt.Sprintf("\tif %s == nil {\n", subFieldPath)
		// return false, nil
		out += "\t\treturn false, nil\n"
		// }
		out += "\t}\n"
	}
	out += scalarFieldEqual(resVarName, candidatesVarName, shapes[len(shapes)-1].GoTypeElem(), condCfg)
	return out
}

func fieldPathContainsMapOrArray(fieldPath string, shapeRef *awssdkmodel.ShapeRef) bool {
	fp := fieldpath.FromString(fieldPath)
	sr := fp.ShapeRef(shapeRef)

	if sr == nil {
		return false
	}
	if sr.ShapeName == "map" || sr.ShapeName == "list" {
		return true
	}
	if sr.ShapeName == "structure" {
		fieldName := fp.PopFront()
		return fieldPathContainsMapOrArray(fp.Copy().At(1), sr.Shape.MemberRefs[fieldName])
	}
	return false
}
