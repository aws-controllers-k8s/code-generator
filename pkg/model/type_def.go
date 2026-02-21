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
	"sort"
	"strings"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/pkg/names"
)

// TypeDef is a Go type definition for structs that are member fields of the
// Spec or Status structs in Custom Resource Definitions (CRDs).
type TypeDef struct {
	Names names.Names
	Attrs map[string]*Attr
	Shape *awssdkmodel.Shape
}

// SortedAttrNames returns the attribute names in sorted order.
func (td *TypeDef) SortedAttrNames() []string {
	names := make([]string, 0, len(td.Attrs))
	for name := range td.Attrs {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetAttribute returns the Attribute with name "attrName".
// This method performs case-insensitive matching to find the Attribute.
func (td *TypeDef) GetAttribute(attrName string) *Attr {
	if td == nil {
		return nil
	}
	for aName, attr := range td.Attrs {
		if strings.EqualFold(aName, attrName) {
			return attr
		}
	}
	return nil
}
