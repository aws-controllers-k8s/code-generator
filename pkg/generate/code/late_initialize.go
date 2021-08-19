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
// lateInitializeFieldNames = []string{"Name"}
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
		out += fmt.Sprintf("%s%s = []string{", indent, resVarName)
		for _, fName := range sortedFieldNames {
			out += fmt.Sprintf("%q,", fName)
		}
		out += "}\n"
	} else {
		out += fmt.Sprintf("%s%s = []string{}\n", indent, resVarName)
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
//	return latest
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
	lateInitializedFieldNames, _ := getSortedLateInitFieldsAndConfig(cfg, r)
	if len(lateInitializedFieldNames) == 0 {
		return fmt.Sprintf("%sreturn %s", indent, targetResVarName)
	}
	out += fmt.Sprintf("%sobservedKo := rm.concreteResource(%s).ko\n", indent, sourceResVarName)
	out += fmt.Sprintf("%slatestKo := rm.concreteResource(%s).ko\n", indent, targetResVarName)
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
	out += fmt.Sprintf("%sreturn %s", indent, targetResVarName)
	return out
}

// CalculateRequeueDelay returns the go code which
// a) checks whether all the fields are late initialized and
// b) if any fields are not initialized, updates the 'delayVarNameInt' and 'incompleteInitializationVarNameBool', which
//    are used to requeue the requests based on the delay configured in LateInitializationConfig.
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
//ko := rm.concreteResource(latest).ko
//numLateInitializationAttempt := ackannotation.GetNumLateInitializationAttempt(latest.MetaObject())
//requeueDelay := time.Duration(0)*time.Second
//incompleteInitialization := false
//if ko.Spec.ImageScanningConfiguration != nil {
//	if ko.Spec.ImageScanningConfiguration.ScanOnPush == nil {
//		fDelay := (&acktypes.Exponential{Initial:time.Duration(5)*time.Second, Factor: 2, MaxDelay: time.Duration(15)*time.Second,}).GetBackoff(numLateInitializationAttempt)
//		requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
//		incompleteInitialization= true
//	}
//}
//if ko.Spec.Name == nil {
//	fDelay := (&acktypes.Exponential{Initial:time.Duration(0)*time.Second, Factor: 2, MaxDelay: time.Duration(0)*time.Second,}).GetBackoff(numLateInitializationAttempt)
//	requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
//	incompleteInitialization= true
//}
//if ko.Spec.another != nil {
//	if ko.Spec.another.map != nil {
//		if ko.Spec.another.map["lastfield"] == nil {
//			fDelay := (&acktypes.Exponential{Initial:time.Duration(5)*time.Second, Factor: 2, MaxDelay: time.Duration(0)*time.Second,}).GetBackoff(numLateInitializationAttempt)
//			requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
//			incompleteInitialization= true
//		}
//	}
//}
//if ko.Spec.map != nil {
//	if ko.Spec.map["subfield"] != nil {
//		if ko.Spec.map["subfield"].x == nil {
//			fDelay := (&acktypes.Exponential{Initial:time.Duration(5)*time.Second, Factor: 2, MaxDelay: time.Duration(0)*time.Second,}).GetBackoff(numLateInitializationAttempt)
//			requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
//			incompleteInitialization= true
//		}
//	}
//}
//if ko.Spec.some != nil {
//	if ko.Spec.some.list == nil {
//		fDelay := (&acktypes.Exponential{Initial:time.Duration(10)*time.Second, Factor: 2, MaxDelay: time.Duration(0)*time.Second,}).GetBackoff(numLateInitializationAttempt)
//		requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
//		incompleteInitialization= true
//	}
//}
//if ko.Spec.structA != nil {
//	if ko.Spec.structA.mapB != nil {
//		if ko.Spec.structA.mapB["structC"] != nil {
//			if ko.Spec.structA.mapB["structC"].valueD == nil {
//				fDelay := (&acktypes.Exponential{Initial:time.Duration(20)*time.Second, Factor: 2, MaxDelay: time.Duration(0)*time.Second,}).GetBackoff(numLateInitializationAttempt)
//				requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
//				incompleteInitialization= true
//			}
//		}
//	}
//}
//return requeueDelay, incompleteInitialization
func CalculateRequeueDelay(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	resVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	sortedLateInitFieldNames, fieldNameToLateInitConfig := getSortedLateInitFieldsAndConfig(cfg, r)
	if len(sortedLateInitFieldNames) == 0 {
		out += fmt.Sprintf("%sreturn time.Duration(0), false", indent)
		return out
	}
	out += fmt.Sprintf("%sko := rm.concreteResource(%s).ko\n", indent, resVarName)
	out += fmt.Sprintf("%snumLateInitializationAttempt := ackannotation.GetNumLateInitializationAttempt(latest.MetaObject())\n", indent)
	out += fmt.Sprintf("%srequeueDelay := time.Duration(0)*time.Second\n", indent)
	out += fmt.Sprintf("%sincompleteInitialization := false\n", indent)
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
				minBackoffSeconds := fieldNameToLateInitConfig[fName].MinBackoffSeconds
				maxBackoffSeconds := fieldNameToLateInitConfig[fName].MaxBackoffSeconds
				out += fmt.Sprintf("%sfDelay := (&acktypes.Exponential{Initial:time.Duration(%d)*time.Second, Factor: 2, MaxDelay: time.Duration(%d)*time.Second,}).GetBackoff(numLateInitializationAttempt)\n", indent, minBackoffSeconds, maxBackoffSeconds)
				out += fmt.Sprintf("%srequeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second\n", indent)
				out += fmt.Sprintf("%sincompleteInitialization= true\n", indent)
			}
		}
		// Close all if blocks with proper indentation
		fNameIndentLevel = fNameIndentLevel - 1
		for fNameIndentLevel >= indentLevel {
			out += fmt.Sprintf("%s}\n", strings.Repeat("\t", fNameIndentLevel))
			fNameIndentLevel = fNameIndentLevel - 1
		}
	}
	out += fmt.Sprintf("%sreturn requeueDelay, incompleteInitialization", indent)
	return out
}
