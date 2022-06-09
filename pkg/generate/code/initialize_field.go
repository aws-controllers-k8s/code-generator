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
//  ```
// 	r.ko.BucketLoggingStatus = &svcapitypes.BucketLoggingStatus{}
//	r.ko.BucketLoggingStatus.LoggingEnabled = &svcapitypes.LoggingEnabled{}
//  ```

func InitializeNestedStructField(
	r *model.CRD,
	sourceVarName string,
	field *model.Field,
	// apiPkgImportName contains the imported package name where the type definition
	// for nested structs is present.
	// ex: svcapitypes "github.com/aws-controllers-k8s/s3-controller/apis/v1alpha1"
	apiPkgImportName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
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
				panic(fmt.Sprintf("unable to find the field with name %s"+
					" for fieldpath %s", front, fieldPath))
			}
			if frontField.ShapeRef == nil {
				panic(fmt.Sprintf("nil ShapeRef for field %s", front))
			}
			fieldShapePath := strings.Replace(fieldPath, front,
				frontField.ShapeRef.ShapeName, 1)
			fsp := fieldpath.FromString(fieldShapePath)
			var index int
			elemAccessPrefix := sourceVarName
			var importPath string
			if apiPkgImportName != "" {
				importPath = fmt.Sprintf("%s.", apiPkgImportName)
			}
			// traverse over the fieldShapePath and initialize every element
			// except the last.
			for index < fsp.Size()-1 {
				elemName := fsp.At(index)
				elemShapeRef := fsp.ShapeRefAt(frontField.ShapeRef, index)
				if elemShapeRef.Shape.Type != "structure" {
					panic(fmt.Sprintf("only nested structures are supported."+
						" Shape type for %s is %s inside fieldpath %s", elemName,
						elemShapeRef.Shape.Type, fieldPath))
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
	return out
}
