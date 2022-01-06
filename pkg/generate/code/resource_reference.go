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

// ReferenceFieldsValidation produces the go code to validate reference field and
// corresponding identifier field.
// Sample code:
// if ko.Spec.APIIDRef != nil && ko.Spec.APIID != nil {
//		return ackerr.ResourceReferenceAndIDNotSupportedFor("APIID", "APIIDRef")
//	}
//	if ko.Spec.APIIDRef == nil && ko.Spec.APIID == nil {
//		return ackerr.ResourceReferenceOrIDRequiredFor("APIID", "APIIDRef")
//	}
func ReferenceFieldsValidation(
	crd *model.CRD,
	sourceVarName string,
	referenceFieldSuffix string,
	indentLevel int,
) string {
	out := ""
	for _, field := range crd.Fields {
		if field.HasReference() {
			indent := strings.Repeat("\t", indentLevel)
			// Validation to make sure both target field and reference are
			// not present at the same time in desired resource
			out += fmt.Sprintf("%sif %s.Spec.%s%s != nil"+
				" && %s.Spec.%s != nil {\n", indent, sourceVarName, field.Names.Camel,
				referenceFieldSuffix, sourceVarName, field.Names.Camel)
			out += fmt.Sprintf("%s\treturn "+
				"ackerr.ResourceReferenceAndIDNotSupportedFor(\"%s\", \"%s%s\")\n",
				indent,field.Names.Camel, field.Names.Camel, referenceFieldSuffix)
			out += fmt.Sprintf("%s}\n", indent)

			// If the field is required, make sure either Ref or original
			// field is present in the resource
			if field.IsRequired() {
				out += fmt.Sprintf("%sif %s.Spec.%s%s == nil &&"+
					" %s.Spec.%s == nil {\n", indent, sourceVarName,
					field.Names.Camel, referenceFieldSuffix, sourceVarName,
					field.Names.Camel)
				out += fmt.Sprintf("%s\treturn "+
					"ackerr.ResourceReferenceOrIDRequiredFor(\"%s\", \"%s%s\")\n",
					indent, field.Names.Camel, field.Names.Camel,
					referenceFieldSuffix)
				out += fmt.Sprintf("%s}\n", indent)
			}
		}
	}
	return out
}

// ReferenceFieldsPresent produces go code(logical condition) for finding whether
// a non-nil reference field is present in a resource. This checks helps in deciding
// whether ACK.ReferencesResolved condition should be added to resource status
// Sample Code:
// return false || ko.Spec.APIIDRef != nil
func ReferenceFieldsPresent(
	crd *model.CRD,
	sourceVarName string,
) string {
	out := "false"
	for _, field := range crd.Fields {
		if field.IsReference() {
			out += fmt.Sprintf(" || %s.Spec.%s != nil", sourceVarName,
				field.Names.Camel)
		}
	}
	return out
}
