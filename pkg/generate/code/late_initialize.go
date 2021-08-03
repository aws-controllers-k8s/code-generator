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

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
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
	sortedFieldNames, _ := getSortedLateInitFieldsAndConfig(cfg, r)
	if len(sortedFieldNames) > 0 {
		out += fmt.Sprintf("%svar %s = []string{", indent, resVarName)
		for _, fName := range sortedFieldNames {
			out += fmt.Sprintf("%q,", fName)
		}
		out += "}\n"
	}
	return out
}

// getSortedLateInitFieldsAndConfig returns the field names in alphabetically sorted order which have LateInitialization
// configuration inside generator config and also a map from fieldName to LateInitializationConfig.
func getSortedLateInitFieldsAndConfig(
	cfg *ackgenconfig.Config,
	r *model.CRD,
) ([]string, map[string]*ackgenconfig.LateInitializeConfig) {
	fieldNameToConfig := cfg.ResourceFields(r.Names.Original)
	fieldNameToLateInitConfig := make(map[string]*ackgenconfig.LateInitializeConfig)
	sortedLateInitFieldNames := make([]string, 0)
	if len(fieldNameToConfig) > 0 {
		for fName, fConfig := range fieldNameToConfig {
			if fConfig != nil && fConfig.LateInitialize != nil {
				fieldNameToLateInitConfig[fName] = fConfig.LateInitialize
				sortedLateInitFieldNames = append(sortedLateInitFieldNames, fName)
			}
		}
		sort.Strings(sortedLateInitFieldNames)
	}
	return sortedLateInitFieldNames, fieldNameToLateInitConfig
}

// LateInitializeFromReadOne returns the gocode to set LateInitialization fields from the ReadOne output
// Field path separated by '.' indicates members in a struct
// Field path separated by '..' indicates member/key in a map
// Note: Unlike Map, updating individual element of a list is not supported. LateInitializing complete list is supported.
//
// Sample generator config:
// fields:
//      Name:
//        late_initialize: {}
//      ImageScanningConfiguration.ScanOnPush:
//        late_initialize:
//          min_backoff_seconds: 5
//          max_backoff_seconds: 15
//      map..subfield.x:
//        late_initialize:
//          min_backoff_seconds: 5
//      another.map..lastfield:
//        late_initialize:
//          min_backoff_seconds: 5
//      some.list:
//        late_initialize:
//          min_backoff_seconds: 10
//      structA.mapB..structC.valueD:
//        late_initialize:
//          min_backoff_seconds: 20
//
// Sample output:
//if observed.Spec.ImageScanningConfiguration != nil && koWithDefaults.Spec.ImageScanningConfiguration != nil {
//	if observed.Spec.ImageScanningConfiguration.ScanOnPush != nil && koWithDefaults.Spec.ImageScanningConfiguration.ScanOnPush == nil {
//		koWithDefaults.Spec.ImageScanningConfiguration.ScanOnPush = observed.Spec.ImageScanningConfiguration.ScanOnPush
//	}
//}
//if observed.Spec.Name != nil && koWithDefaults.Spec.Name == nil {
//	koWithDefaults.Spec.Name = observed.Spec.Name
//}
//if observed.Spec.another != nil && koWithDefaults.Spec.another != nil {
//	if observed.Spec.another.map != nil && koWithDefaults.Spec.another.map != nil {
//		if observed.Spec.another.map[lastfield] != nil && koWithDefaults.Spec.another.map[lastfield] == nil {
//		koWithDefaults.Spec.another.map[lastfield] = observed.Spec.another.map[lastfield]
//	}
//	}
//}
//if observed.Spec.map != nil && koWithDefaults.Spec.map != nil {
//	if observed.Spec.map[subfield] != nil && koWithDefaults.Spec.map[subfield] != nil {
//	if observed.Spec.map[subfield].x != nil && koWithDefaults.Spec.map[subfield].x == nil {
//	koWithDefaults.Spec.map[subfield].x = observed.Spec.map[subfield].x
//}
//}
//}
//if observed.Spec.some != nil && koWithDefaults.Spec.some != nil {
//	if observed.Spec.some.list != nil && koWithDefaults.Spec.some.list == nil {
//		koWithDefaults.Spec.some.list = observed.Spec.some.list
//	}
//}
func LateInitializeFromReadOne(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	sourceKoVarName string,
	targetKoVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	lateInitializedFieldNames, _ := getSortedLateInitFieldsAndConfig(cfg, r)
	// TODO(vijat@): Add validation for correct field path in lateInitializedFieldNames
	for _, fName := range lateInitializedFieldNames {
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
				out += fmt.Sprintf("%sif %s.%s != nil && %s.%s != nil {\n", indent, sourceKoVarName, fNamePartAccesor, targetKoVarName, fNamePartAccesor)
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
				out += fmt.Sprintf("%sif %s.%s != nil && %s.%s == nil {\n", indent, sourceKoVarName, fNamePartAccesor, targetKoVarName, fNamePartAccesor)
				fNameIndentLevel = fNameIndentLevel + 1
				indent = strings.Repeat("\t", fNameIndentLevel)
				out += fmt.Sprintf("%s%s.%s = %s.%s\n", indent, targetKoVarName, fNamePartAccesor, sourceKoVarName, fNamePartAccesor)
			}
		}
		// Close all if blocks with proper indentation
		fNameIndentLevel = fNameIndentLevel - 1
		for fNameIndentLevel >= indentLevel {
			out += fmt.Sprintf("%s}\n", strings.Repeat("\t", fNameIndentLevel))
			fNameIndentLevel = fNameIndentLevel - 1
		}
	}
	return out
}

// IncompleteLateInitialization returns the go code which checks whether all the fields are late initialized.
// If all the fields are not late initialized, this method also returns the requeue delay needed to attempt
// late initialization again.
//
// Sample GeneratorConfig:
// fields:
//      Name:
//        late_initialize: {}
//      ImageScanningConfiguration.ScanOnPush:
//        late_initialize:
//          min_backoff_seconds: 5
//          max_backoff_seconds: 15
//      map..subfield.x:
//        late_initialize:
//          min_backoff_seconds: 5
//      another.map..lastfield:
//        late_initialize:
//          min_backoff_seconds: 5
//      some.list:
//        late_initialize:
//          min_backoff_seconds: 10
//      structA.mapB..structC.valueD:
//        late_initialize:
//          min_backoff_seconds: 20
//
//
// Sample Output:
//	ko := latestWithDefaults.ko
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
//
func IncompleteLateInitialization(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	resVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	sortedLateInitFieldNames, _ := getSortedLateInitFieldsAndConfig(cfg, r)
	if len(sortedLateInitFieldNames) == 0 {
		out += fmt.Sprintf("%sreturn false\n", indent)
		return out
	}
	out += fmt.Sprintf("%sko := %s.ko\n", indent, resVarName)
	for _, fName := range sortedLateInitFieldNames {
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
	out += fmt.Sprintf("%sreturn false\n", indent)
	return out
}
