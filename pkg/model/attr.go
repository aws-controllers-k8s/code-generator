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
