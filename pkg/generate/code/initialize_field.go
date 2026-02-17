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

// InitializeNestedStructField returns the go code for initializing a nested
// struct field. Currently this method only supports the struct shape for
// nested elements.
//
// TODO(vijtrip2): Refactor the code out of set_resource.go for generating
// constructors and reuse here. This method is currently being used for handling
// nested Tagging fields.
//
// Example: generated code for "Logging.LoggingEnabled.TargetBucket" field
// inside "s3" "bucket" crd looks like:
//
//	 ```
//		r.ko.Spec.Logging = &svcapitypes.BucketLoggingStatus{}
//		r.ko.Spec.Logging.LoggingEnabled = &svcapitypes.LoggingEnabled{}
//	 ```
func InitializeNestedStructField(
	r *model.CRD,
	sourceVarName string,
	field *model.Field,
	// apiPkgAlias contains the imported package alias where the type definition
	// for nested structs is present.
	// ex: svcapitypes "github.com/aws-controllers-k8s/s3-controller/apis/v1alpha1"
	apiPkgAlias string,
	// Number of levels of indentation to use
	indentLevel int,
) (string, error) {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	fieldPath := field.Path
	if fieldPath != "" {
		fp := fieldpath.FromString(fieldPath)
		if fp.Size() > 1 {
			// replace the front field name with front field shape name inside
			// the field path to construct the fieldShapePath
			front := fp.Front()
			frontField := r.Fields[front]
			if frontField == nil {
				return "", fmt.Errorf(
					"resource %q, field %q: unable to find field %q in fieldpath",
					r.Names.Original, fieldPath, front,
				)
			}
			if frontField.ShapeRef == nil {
				return "", fmt.Errorf(
					"resource %q, field %q: nil ShapeRef for field %q",
					r.Names.Original, fieldPath, front,
				)
			}
			fieldShapePath := strings.Replace(fieldPath, front,
				frontField.ShapeRef.ShapeName, 1)
			fsp := fieldpath.FromString(fieldShapePath)
			var index int
			// Build the prefix to access elements in field path.
			// Use the front of fieldpath to determine whether the field is
			// a spec field or status field.
			elemAccessPrefix := sourceVarName
			if _, found := r.SpecFields[front]; found {
				elemAccessPrefix = fmt.Sprintf("%s%s", elemAccessPrefix,
					r.Config().PrefixConfig.SpecField)
			} else {
				elemAccessPrefix = fmt.Sprintf("%s%s", elemAccessPrefix,
					r.Config().PrefixConfig.StatusField)
			}
			var importPath string
			if apiPkgAlias != "" {
				importPath = fmt.Sprintf("%s.", apiPkgAlias)
			}
			// traverse over the fieldShapePath and initialize every element
			// except the last.
			for index < fsp.Size()-1 {
				elemName := fp.At(index)
				elemShapeRef := fsp.ShapeRefAt(frontField.ShapeRef, index)
				if elemShapeRef.Shape.Type != "structure" {
					return "", fmt.Errorf(
						"resource %q, field %q: only nested structures are supported, but %q has shape type %q",
						r.Names.Original, fieldPath, elemName, elemShapeRef.Shape.Type,
					)
				}
				out += fmt.Sprintf("%s%s.%s = &%s%s{}\n",
					indent, elemAccessPrefix, elemName, importPath,
					elemShapeRef.GoTypeElem())
				elemAccessPrefix = fmt.Sprintf("%s.%s", elemAccessPrefix,
					elemName)
				index++
			}
		}
	}
	return out, nil
}
