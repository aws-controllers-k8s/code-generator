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

package model

import (
	"fmt"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/pkg/names"
)

type Attr struct {
	Names       names.Names
	GoType      string
	Shape       *awssdkmodel.Shape
	GoTag       string
	IsImmutable bool
	// ImmutabilityAllowsFirstWrite selects the weaker once-set immutability
	// rule for this member, because the controller itself may legitimately
	// perform the member's first write. Only meaningful alongside IsImmutable.
	//
	// Two configurations cause it, both set in setTypeDefAttributeImmutable:
	// late_initialize, and a nested reference whose resolved value the
	// controller can materialize into the stored spec. See immutabilityCELRule.
	ImmutabilityAllowsFirstWrite bool
}

func NewAttr(
	names names.Names,
	goType string,
	shape *awssdkmodel.Shape,
) *Attr {
	return &Attr{
		Names:  names,
		GoType: goType,
		Shape:  shape,
	}
}

// GetGoTag returns the Go Tag to inject for this attribute. If the GoTag
// field is not empty, it will be used. Otherwise, one will be generated
// from the attribute's name.
func (a *Attr) GetGoTag() string {
	if a.GoTag != "" {
		return a.GoTag
	}
	return fmt.Sprintf("`json:\"%s,omitempty\"`", a.Names.CamelLower)
}

// GetCRDJSONFieldName returns the JSON name this attribute is serialized as in
// the generated CRD type, i.e. the name a Kubernetes user sees in the
// resource's YAML. It derives the name by parsing the attribute's own struct
// tag as produced by GetGoTag, so it always agrees with what the generated Go
// struct actually serializes -- including any `go_tag` override from
// generator.yaml.
//
// This is the Attr analogue of Field.GetCRDJSONFieldName.
func (a *Attr) GetCRDJSONFieldName() string {
	return jsonNameFromGoTag(a.GetGoTag(), a.Names.CamelLower)
}

// ImmutabilityCELRule returns the CEL expression that a containing struct must
// carry in order to freeze this attribute. It is only meaningful when
// IsImmutable is true, and is consumed by the apis/type_def.go.tpl template,
// which renders the rule as an XValidation marker on the struct that owns this
// attribute rather than on the attribute itself.
//
// See immutabilityCELRule for why the rule lives on the parent, and why a
// late-initialized attribute gets the weaker once-set form.
func (a *Attr) ImmutabilityCELRule() string {
	return immutabilityCELRule(a.GetCRDJSONFieldName(), a.ImmutabilityAllowsFirstWrite)
}

// ImmutabilityCELFieldPath returns the XValidation fieldPath that attributes an
// immutability rejection to this attribute rather than to its containing struct.
func (a *Attr) ImmutabilityCELFieldPath() string {
	return immutabilityCELFieldPath(a.GetCRDJSONFieldName())
}
