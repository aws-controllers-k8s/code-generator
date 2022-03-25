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

	"github.com/aws-controllers-k8s/code-generator/pkg/fieldpath"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// NilFieldPathCheck returns the condition statement for Nil check
// on a field path. This nil check on field path is useful to avoid
// nil pointer panics when accessing a field value.
//
// This function only outputs the logical condition and not the "if" block
// so that the output can be reused in many templates, where
// logic inside "if" block can be different.
//
// Example Output for fieldpath "JWTConfiguration.Issuer.SomeField" is
// "ko.Spec.JWTConfiguration == nil || ko.Spec.JWTConfiguration.Issuer == nil"
func NilFieldPathCheck(field *model.Field, sourceVarName string) string {
	out := ""
	fp := fieldpath.FromString(field.Path)
	// remove fieldName from fieldPath before adding nil checks
	fp.Pop()
	fieldNamePrefix := ""
	for fp.Size() > 0 {
		fieldNamePrefix = fmt.Sprintf("%s.%s", fieldNamePrefix, fp.PopFront())
		out += fmt.Sprintf(" || %s%s == nil", sourceVarName, fieldNamePrefix)
	}
	return strings.TrimPrefix(out, " || ")
}
