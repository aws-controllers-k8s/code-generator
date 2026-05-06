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
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// FieldGroupDeltaCheck returns a Go boolean expression string that evaluates
// to true if any of the field group's payload fields have changed according
// to the delta. Used by update orchestrator templates to decide whether to
// call a specific field-group update operation.
//
// The generated expression looks like:
//
//	delta.DifferentAt("Spec.ImageScanningConfiguration")
//
// For multiple payload fields:
//
//	delta.DifferentAt("Spec.ImageScanningConfiguration") || delta.DifferentAt("Spec.ImageTagMutability")
func FieldGroupDeltaCheck(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	fg *model.FieldGroupOperation,
	deltaVarName string,
) string {
	if len(fg.PayloadFields) == 0 {
		return "false"
	}

	parts := make([]string, 0, len(fg.PayloadFields))
	for _, f := range fg.PayloadFields {
		fieldPath := fmt.Sprintf("%s.%s",
			strings.TrimPrefix(cfg.PrefixConfig.SpecField, "."),
			f.Names.Camel,
		)
		parts = append(parts, fmt.Sprintf(
			"%s.DifferentAt(%q)", deltaVarName, fieldPath,
		))
	}
	return strings.Join(parts, " || ")
}
