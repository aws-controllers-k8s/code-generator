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

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// FindLateInitializedFieldNames outputs the code to create a sorted slice of fieldNames to
// late initialize. This slice helps with short circuiting the AWSResourceManager.LateInitialize()
// method if there are no fields to late initialize.
//
// Sample Output:
// var lateInitializeFieldNames = []string{"Name"}
func FindLateInitializedFieldNames(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	resVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	var lateInitFieldNames []string
	lateInitConfigs := cfg.GetLateInitConfigs(r.Names.Original)
	for fieldName := range lateInitConfigs {
		lateInitFieldNames = append(lateInitFieldNames, fieldName)
	}
	if len(lateInitFieldNames) == 0 {
		return fmt.Sprintf("%svar %s = []string{}\n", indent, resVarName)
	}
	// sort the slice to help with short circuiting AWSResourceManager.LateInitialize()
	sort.Strings(lateInitFieldNames)
	out += fmt.Sprintf("%svar %s = []string{", indent, resVarName)
	for _, fName := range lateInitFieldNames {
		out += fmt.Sprintf("%q,", fName)
	}
	out += "}\n"
	return out
}

// LateInitializeFromReadOne returns the gocode to set LateInitialization fields from the ReadOne output
// Field path separated by '.' indicates members in a struct
// Field path separated by '..' indicates member/key in a map
// Note: Unlike Map, updating individual element of a list is not supported. LateInitializing complete list is supported.
//
// Sample generator config:
// fields:
//
//	Name:
//	  late_initialize: {}
//	ImageScanningConfiguration.ScanOnPush:
//	  late_initialize:
//	    min_backoff_seconds: 5
//	    max_backoff_seconds: 15
//	map..subfield.x:
//	  late_initialize:
//	    min_backoff_seconds: 5
//	another.map..lastfield:
//	  late_initialize:
//	    min_backoff_seconds: 5
//	some.list:
//	  late_initialize:
//	    min_backoff_seconds: 10
//	structA.mapB..structC.valueD:
//	  late_initialize:
//	    min_backoff_seconds: 20
//
// Sample output:
//
//	observedKo := rm.concreteResource(observed).ko
//	latestKo := rm.concreteResource(latest).ko
//	if observedKo.Spec.ImageScanningConfiguration != nil && latestKo.Spec.ImageScanningConfiguration != nil {
//		if observedKo.Spec.ImageScanningConfiguration.ScanOnPush != nil && latestKo.Spec.ImageScanningConfiguration.ScanOnPush == nil {
//			latestKo.Spec.ImageScanningConfiguration.ScanOnPush = observedKo.Spec.ImageScanningConfiguration.ScanOnPush
//		}
//	}
//	if observedKo.Spec.Name != nil && latestKo.Spec.Name == nil {
//		latestKo.Spec.Name = observedKo.Spec.Name
//	}
//	if observedKo.Spec.another != nil && latestKo.Spec.another != nil {
//		if observedKo.Spec.another.map != nil && latestKo.Spec.another.map != nil {
//			if observedKo.Spec.another.map["lastfield"] != nil && latestKo.Spec.another.map["lastfield"] == nil {
//				latestKo.Spec.another.map["lastfield"] = observedKo.Spec.another.map["lastfield"]
//			}
//		}
//	}
//	if observedKo.Spec.map != nil && latestKo.Spec.map != nil {
//		if observedKo.Spec.map["subfield"] != nil && latestKo.Spec.map["subfield"] != nil {
//			if observedKo.Spec.map["subfield"].x != nil && latestKo.Spec.map["subfield"].x == nil {
//				latestKo.Spec.map["subfield"].x = observedKo.Spec.map["subfield"].x
//			}
//		}
//	}
//	if observedKo.Spec.some != nil && latestKo.Spec.some != nil {
//		if observedKo.Spec.some.list != nil && latestKo.Spec.some.list == nil {
//			latestKo.Spec.some.list = observedKo.Spec.some.list
//		}
//	}
//	if observedKo.Spec.structA != nil && latestKo.Spec.structA != nil {
//		if observedKo.Spec.structA.mapB != nil && latestKo.Spec.structA.mapB != nil {
//			if observedKo.Spec.structA.mapB["structC"] != nil && latestKo.Spec.structA.mapB["structC"] != nil {
//				if observedKo.Spec.structA.mapB["structC"].valueD != nil && latestKo.Spec.structA.mapB["structC"].valueD == nil {
//					latestKo.Spec.structA.mapB["structC"].valueD = observedKo.Spec.structA.mapB["structC"].valueD
//				}
//			}
//		}
//	}
//	return &resource{latestKo}
func LateInitializeFromReadOne(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	sourceResVarName string,
	targetResVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	var lateInitFieldNames []string
	lateInitConfigs := cfg.GetLateInitConfigs(r.Names.Original)
	for fieldName := range lateInitConfigs {
		lateInitFieldNames = append(lateInitFieldNames, fieldName)
	}
	if len(lateInitFieldNames) == 0 {
		return fmt.Sprintf("%sreturn %s", indent, targetResVarName)
	}
	sort.Strings(lateInitFieldNames)
	out += fmt.Sprintf("%sobservedKo := rm.concreteResource(%s).ko.DeepCopy()\n", indent, sourceResVarName)
	out += fmt.Sprintf("%slatestKo := rm.concreteResource(%s).ko.DeepCopy()\n", indent, targetResVarName)
	// TODO(vijat@): Add validation for correct field path in lateInitializedFieldNames
	for _, fName := range lateInitFieldNames {
		// split the field name by period
		// each substring represents a field.
		fNameParts := strings.Split(fName, ".")
		// fNameIndentLevel tracks the indentation level for every new line added
		// This variable is incremented when building nested if blocks and decremented when closing those if blocks.
		fNameIndentLevel := indentLevel
		// fParentPath keeps track of parent path for any fNamePart
		fParentPath := ""
		mapShapedParent := false
		// for every part except last, perform the nil check
		// entries in both source and target koVarName should not be nil
		for i, fNamePart := range fNameParts {
			if fNamePart == "" {
				mapShapedParent = true
				continue
			}
			indent := strings.Repeat("\t", fNameIndentLevel)
			fNamePartAccesor := fmt.Sprintf("Spec%s.%s", fParentPath, fNamePart)
			if mapShapedParent {
				fNamePartAccesor = fmt.Sprintf("Spec%s[%q]", fParentPath, fNamePart)
			}
			// Handling for all parts except last one
			if i != len(fNameParts)-1 {
				out += fmt.Sprintf("%sif observedKo.%s != nil && latestKo.%s != nil {\n", indent, fNamePartAccesor, fNamePartAccesor)
				// update fParentPath and fNameIndentLevel for next iteration
				if mapShapedParent {
					fParentPath = fmt.Sprintf("%s[%q]", fParentPath, fNamePart)
					mapShapedParent = false
				} else {
					fParentPath = fmt.Sprintf("%s.%s", fParentPath, fNamePart)
				}
				fNameIndentLevel = fNameIndentLevel + 1
			} else {
				// handle last part here
				// for last part, set the lateInitialized field if user did not specify field value and readOne has server side defaulted value.
				// i.e. field is not nil in sourceKoVarName but is nil in targetkoVarName
				out += fmt.Sprintf("%sif observedKo.%s != nil && latestKo.%s == nil {\n", indent, fNamePartAccesor, fNamePartAccesor)
				fNameIndentLevel = fNameIndentLevel + 1
				indent = strings.Repeat("\t", fNameIndentLevel)
				out += fmt.Sprintf("%slatestKo.%s = observedKo.%s\n", indent, fNamePartAccesor, fNamePartAccesor)
			}
		}
		// Close all if blocks with proper indentation
		fNameIndentLevel = fNameIndentLevel - 1
		for fNameIndentLevel >= indentLevel {
			out += fmt.Sprintf("%s}\n", strings.Repeat("\t", fNameIndentLevel))
			fNameIndentLevel = fNameIndentLevel - 1
		}
	}
	out += fmt.Sprintf("%sreturn &resource{latestKo}", indent)
	return out
}

// IncompleteLateInitialization returns the go code which checks whether all the fields are late initialized.
// If all the fields are not late initialized, this method also returns the requeue delay needed to attempt
// late initialization again.
//
// Sample GeneratorConfig:
// fields:
//
//	Name:
//	  late_initialize: {}
//	ImageScanningConfiguration.ScanOnPush:
//	  late_initialize:
//	    min_backoff_seconds: 5
//	    max_backoff_seconds: 15
//	map..subfield.x:
//	  late_initialize:
//	    min_backoff_seconds: 5
//	another.map..lastfield:
//	  late_initialize:
//	    min_backoff_seconds: 5
//	some.list:
//	  late_initialize:
//	    min_backoff_seconds: 10
//	structA.mapB..structC.valueD:
//	  late_initialize:
//	    min_backoff_seconds: 20
//
// Sample Output:
//
//	ko := rm.concreteResource(latest).ko.DeepCopy()
//	if ko.Spec.ImageScanningConfiguration != nil {
//		if ko.Spec.ImageScanningConfiguration.ScanOnPush == nil {
//			return true
//		}
//	}
//	if ko.Spec.Name == nil {
//		return true
//	}
//	if ko.Spec.another != nil {
//		if ko.Spec.another.map != nil {
//			if ko.Spec.another.map["lastfield"] == nil {
//				return true
//			}
//		}
//	}
//	if ko.Spec.map != nil {
//		if ko.Spec.map["subfield"] != nil {
//			if ko.Spec.map["subfield"].x == nil {
//				return true
//			}
//		}
//	}
//	if ko.Spec.some != nil {
//		if ko.Spec.some.list == nil {
//			return true
//		}
//	}
//	if ko.Spec.structA != nil {
//		if ko.Spec.structA.mapB != nil {
//			if ko.Spec.structA.mapB["structC"] != nil {
//				if ko.Spec.structA.mapB["structC"].valueD == nil {
//					return true
//				}
//			}
//		}
//	}
//	return false
func IncompleteLateInitialization(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	resVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	var lateInitFieldNames []string
	lateInitConfigs := cfg.GetLateInitConfigs(r.Names.Original)
	for fieldName, lateInitConfig := range lateInitConfigs {
		if lateInitConfig.SkipIncompleteCheck != nil {
			continue
		}
		lateInitFieldNames = append(lateInitFieldNames, fieldName)
	}
	if len(lateInitFieldNames) == 0 {
		return fmt.Sprintf("%sreturn false", indent)
	}
	sort.Strings(lateInitFieldNames)
	out += fmt.Sprintf("%sko := rm.concreteResource(%s).ko.DeepCopy()\n", indent, resVarName)
	for _, fName := range lateInitFieldNames {
		// split the field name by period
		// each substring represents a field.
		fNameParts := strings.Split(fName, ".")
		// fNameIndentLevel tracks the indentation level for every new line added
		// This variable is incremented when building nested if blocks and decremented when closing those if blocks.
		fNameIndentLevel := indentLevel
		// fParentPath keeps track of parent path for any fNamePart
		fParentPath := ""
		mapShapedParent := false
		for i, fNamePart := range fNameParts {
			if fNamePart == "" {
				mapShapedParent = true
				continue
			}
			indent := strings.Repeat("\t", fNameIndentLevel)
			fNamePartAccesor := fmt.Sprintf("Spec%s.%s", fParentPath, fNamePart)
			if mapShapedParent {
				fNamePartAccesor = fmt.Sprintf("Spec%s[%q]", fParentPath, fNamePart)
			}
			// Handling for all parts except last one
			if i != len(fNameParts)-1 {
				out += fmt.Sprintf("%sif ko.%s != nil {\n", indent, fNamePartAccesor)
				// update fParentPath and fNameIndentLevel for next iteration
				if mapShapedParent {
					fParentPath = fmt.Sprintf("%s[%q]", fParentPath, fNamePart)
					mapShapedParent = false
				} else {
					fParentPath = fmt.Sprintf("%s.%s", fParentPath, fNamePart)
				}
				fNameIndentLevel = fNameIndentLevel + 1
			} else {
				// handle last part here
				// for last part, if the late initialized field is still nil, calculate the retry backoff using
				// acktypes.LateInitializationRetryConfig abstraction and set the incompleteInitialization flag to true
				out += fmt.Sprintf("%sif ko.%s == nil {\n", indent, fNamePartAccesor)
				fNameIndentLevel = fNameIndentLevel + 1
				indent = strings.Repeat("\t", fNameIndentLevel)
				out += fmt.Sprintf("%sreturn true\n", indent)
			}
		}
		// Close all if blocks with proper indentation
		fNameIndentLevel = fNameIndentLevel - 1
		for fNameIndentLevel >= indentLevel {
			out += fmt.Sprintf("%s}\n", strings.Repeat("\t", fNameIndentLevel))
			fNameIndentLevel = fNameIndentLevel - 1
		}
	}
	out += fmt.Sprintf("%sreturn false", indent)
	return out
}
