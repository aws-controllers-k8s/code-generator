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

// FindLateInitializedFieldsWithDelay outputs the code to create a map of fieldName to
// late intialization delay in seconds.
func FindLateInitializedFieldsWithDelay(
	cfg *ackgenconfig.Config,
	r *model.CRD,
// Number of levels of indentation to use
	indentLevel int,
) string {
	//Sample output
	//var lateInitializeFieldToDelaySeconds = map[string]int{"Name": 0}
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	fieldNameToConfig := cfg.ResourceFields(r.Names.Original)
	if len(fieldNameToConfig) > 0 {
		fieldNameToDelaySeconds := make(map[string]int)
		sortedFieldNames := make([]string, 0)
		for fName, fConfig := range fieldNameToConfig {
			if fConfig != nil && fConfig.LateInitialize != nil {
				fieldNameToDelaySeconds[fName] = fConfig.LateInitialize.DelaySeconds
				sortedFieldNames = append(sortedFieldNames, fName)
			}
		}
		sort.Strings(sortedFieldNames)
		lateInitFieldToDelayValues := ""
		if len(sortedFieldNames) > 0 {
			for _, fName := range sortedFieldNames {
				lateInitFieldToDelayValues += fmt.Sprintf("\"%s\":%d,",fName, fieldNameToDelaySeconds[fName])
			}
			out += fmt.Sprintf("%svar lateInitializeFieldToDelaySeconds = map[string]int{%s}\n", indent, lateInitFieldToDelayValues)
		}
	}
	return out
}

// lateInitializedFieldNames returns the field names which have LateInitialization configuration inside generator config
func lateInitializedFieldNames(
	cfg *ackgenconfig.Config,
	r *model.CRD,
) []string {
	fieldNames := make([]string, 0)
	fieldNameToConfig := cfg.ResourceFields(r.Names.Original)
	if len(fieldNameToConfig) > 0 {
		for fName, fConfig := range fieldNameToConfig {
			if fConfig != nil && fConfig.LateInitialize != nil {
				fieldNames = append(fieldNames, fName)
			}
		}
	}
	return fieldNames
}

// LateInitializeFromReadOne returns the gocode to set LateInitialization fields from the ReadOne output
// TODO(vijat@): add support for Map and list. Currently only structs are supported.
func LateInitializeFromReadOne (
	cfg *ackgenconfig.Config,
	r *model.CRD,
	sourceKoVarName string,
	targetKoVarName string,
// Number of levels of indentation to use
	indentLevel int,
) string {
	//Sample output
	//if observed.Spec.ImageScanningConfiguration != nil && koWithDefaults.Spec.ImageScanningConfiguration != nil {
	//	if observed.Spec.ImageScanningConfiguration.ScanOnPush != nil && koWithDefaults.Spec.ImageScanningConfiguration.ScanOnPush == nil {
	//		koWithDefaults.Spec.ImageScanningConfiguration.ScanOnPush = observed.Spec.ImageScanningConfiguration.ScanOnPush
	//	}
	//}
	//if observed.Spec.Name != nil && koWithDefaults.Spec.Name == nil {
	//	koWithDefaults.Spec.Name = observed.Spec.Name
	//}
	out := ""
	lateInitializedFieldNames := lateInitializedFieldNames(cfg, r)
	// sorting helps produce consistent output for unit test reliability
	sort.Strings(lateInitializedFieldNames)
	for _, fName := range lateInitializedFieldNames {
		// split the field name by period
		// each substring represents a field. No support for '..' currently
		fNameParts := strings.Split(fName, ".")
		// fNameIndentLevel tracks the indentation level for every new line added
		// This variable is incremented when building nested if blocks and decremented when closing those if blocks.
		fNameIndentLevel := indentLevel
		// fParentPath keeps track of parent path for any fNamePart
		fParentPath := ""
		// for every part except last, perform the nil check
		// entries in both source and target koVarName should not be nil
		for i,fNamePart := range fNameParts {
			indent := strings.Repeat("\t", fNameIndentLevel)
			// ignore last part
			if i != len(fNameParts)-1 {
				out += fmt.Sprintf("%sif %s.Spec%s.%s != nil && %s.Spec%s.%s != nil {\n", indent, sourceKoVarName, fParentPath, fNamePart, targetKoVarName, fParentPath, fNamePart)
				// update fParentPath and fNameIndentLevel for next iteration
				fParentPath = fmt.Sprintf("%s.%s", fParentPath, fNamePart)
				fNameIndentLevel = fNameIndentLevel + 1
			}
		}
		// for last part, set the lateInitialized field if user did not specify field value and readOne has server side defaulted value.
		// i.e. field is not nil in sourceKoVarName but is nil in targetkoVarName
		indent := strings.Repeat("\t", fNameIndentLevel)
		lastfNamePart := fNameParts[len(fNameParts)-1]
		out += fmt.Sprintf("%sif %s.Spec%s.%s != nil && %s.Spec%s.%s == nil {\n", indent, sourceKoVarName, fParentPath, lastfNamePart,targetKoVarName, fParentPath, lastfNamePart)
		fNameIndentLevel = fNameIndentLevel + 1
		indent = strings.Repeat("\t", fNameIndentLevel)
		out += fmt.Sprintf("%s%s.Spec%s.%s = %s.Spec%s.%s\n", indent, targetKoVarName, fParentPath,lastfNamePart, sourceKoVarName, fParentPath, lastfNamePart)
		// Close all if blocks with proper indentation
		fNameIndentLevel = fNameIndentLevel - 1
		for fNameIndentLevel >= indentLevel {
			out += fmt.Sprintf("%s}\n", strings.Repeat("\t", fNameIndentLevel))
			fNameIndentLevel = fNameIndentLevel - 1
		}
	}
	return out
}