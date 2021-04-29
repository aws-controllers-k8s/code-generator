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
	"strings"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/names"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
)

// Field represents a single field in the CRD's Spec or Status objects. The
// field may be a direct field of the Spec or Status object or may be a field
// of a list or struct-type field of the Spec or Status object. We call these
// latter fields "nested fields" and they are identified by the Field.Path
// attribute.
type Field struct {
	// CRD is the a pointer to the top-level custom resource definition
	// descriptor for the field or field's parent (if a nested field)
	CRD *CRD
	// Names is a set of normalized names for the field
	Names names.Names
	// Path is a "field path" that indicates where the field is within the CRD.
	// For example "Spec.Name" or "Status.BrokerInstances..Endpoint". Note for
	// the latter example, the field path indicates that the field `Endpoint`
	// is an attribute of the `Status.BrokerInstances` top-level field and the
	// double dot (`..` indicates that BrokerInstances is a list type).
	Path string
	// GoType is a string containing the Go data type for the field
	GoType string
	// GoTypeElem indicates the Go data type for the type of list element if
	// the field is a list type
	GoTypeElem        string
	GoTypeWithPkgName string
	ShapeRef          *awssdkmodel.ShapeRef
	FieldConfig       *ackgenconfig.FieldConfig
}

// IsRequired checks the FieldConfig for Field and returns if the field is
// marked as required or not.A
//
// If there is no required override present for this field in FieldConfig,
// IsRequired will return if the shape is marked as required in AWS SDK Private
// model We use this to append kubebuilder:validation:Required markers to
// validate using the CRD validation schema
func (f *Field) IsRequired() bool {
	if f.FieldConfig != nil && f.FieldConfig.IsRequired != nil {
		return *f.FieldConfig.IsRequired
	}
	return util.InStrings(f.Names.ModelOriginal, f.CRD.Ops.Create.InputRef.Shape.Required)
}

// ParentFieldPath takes a field path and returns the field path of the
// containing "parent" field. For example, if the field path
// `Users..Credentials.Login` is passed in, this function returns
// `Users..Credentials`. If `Users..Password` is supplied, this function
// returns `Users`, etc.
func ParentFieldPath(path string) string {
	parts := strings.Split(path, ".")
	// Pop the last element of the supplied field path
	parts = parts[0 : len(parts)-1]
	// If the parent field's type is a list or map, there will be two dots ".."
	// in the supplied field path. We don't want the returned field path to end
	// in a dot, since that would be invalid, so we trim it off here
	if parts[len(parts)-1] == "" {
		parts = parts[0 : len(parts)-1]
	}
	return strings.Join(parts, ".")
}

// NewField returns a pointer to a new Field object
func NewField(
	crd *CRD,
	path string,
	fieldNames names.Names,
	shapeRef *awssdkmodel.ShapeRef,
	cfg *ackgenconfig.FieldConfig,
) *Field {
	var gte, gt, gtwp string
	var shape *awssdkmodel.Shape
	if shapeRef != nil {
		shape = shapeRef.Shape
	}

	if cfg != nil && cfg.IsSecret {
		gt = "*ackv1alpha1.SecretKeyReference"
		gte = "ackv1alpha1.SecretKeyReference"
		gtwp = "*ackv1alpha1.SecretKeyReference"
	} else if shape != nil {
		gte, gt, gtwp = cleanGoType(crd.sdkAPI, crd.cfg, shape)
	} else {
		gte = "string"
		gt = "*string"
		gtwp = "*string"
	}
	return &Field{
		CRD:               crd,
		Names:             fieldNames,
		Path:              path,
		ShapeRef:          shapeRef,
		GoType:            gt,
		GoTypeElem:        gte,
		GoTypeWithPkgName: gtwp,
		FieldConfig:       cfg,
	}
}
