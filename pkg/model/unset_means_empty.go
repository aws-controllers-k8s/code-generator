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
	"sort"
	"strings"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/pkg/names"
)

// UnsetMeansEmptyPaths returns the runtime paths of the Spec fields for which
// an unset desired value still means "must be empty" when the runtime's
// UnsetFieldsUnmanaged feature gate is enabled: the fields configured with
// `compare.unset_means_empty`, plus the resource's tag field. A runtime path
// uses Go field names and starts with "Spec", e.g. "Spec.Policy". It is
// emitted into the generated resource descriptor and read by the ACK runtime.
func (r *CRD) UnsetMeansEmptyPaths() ([]string, error) {
	var paths []string
	for fieldPath, fc := range r.cfg.GetFieldConfigs(r.Names.Original) {
		if fc == nil || fc.Compare == nil || !fc.Compare.UnsetMeansEmpty {
			continue
		}
		p, err := r.specRuntimePath(fieldPath)
		if err != nil {
			return nil, fmt.Errorf("resource %q: compare.unset_means_empty: %w", r.Names.Original, err)
		}
		paths = appendUniquePath(paths, p)
	}
	if tagField, err := r.GetTagField(); err == nil && tagField != nil {
		if p, err := r.specRuntimePath(tagField.Path); err == nil {
			paths = appendUniquePath(paths, p)
		}
	}
	sort.Strings(paths)
	return paths, nil
}

// specRuntimePath turns a generator.yaml field path ("A.B") into a runtime
// path ("Spec.A.B") by walking the API shapes through structs. A path that
// goes through a list is rejected: the runtime compares lists as a whole, so
// a field inside list elements cannot be matched.
func (r *CRD) specRuntimePath(fieldPath string) (string, error) {
	parts := strings.Split(fieldPath, ".")
	top, ok := r.Fields[parts[0]]
	if !ok || top.ShapeRef == nil || top.ShapeRef.Shape == nil {
		return "", fmt.Errorf("field %q: no such Spec field", fieldPath)
	}
	if _, isSpec := r.SpecFields[top.Names.Original]; !isSpec {
		return "", fmt.Errorf("field %q: not a Spec field", fieldPath)
	}
	out := "Spec." + top.Names.Camel
	shape := top.ShapeRef.Shape
	for _, part := range parts[1:] {
		if shape.Type == "list" {
			return "", fmt.Errorf("field %q: goes through the list %q", fieldPath, out)
		}
		if shape.Type != "structure" {
			return "", fmt.Errorf("field %q: %q is not inside a struct", fieldPath, part)
		}
		var next *awssdkmodel.Shape
		for memberName, memberRef := range shape.MemberRefs {
			if n := names.New(memberName); n.Camel == part || strings.EqualFold(memberName, part) {
				next = memberRef.Shape
				out += "." + n.Camel
				break
			}
		}
		if next == nil {
			return "", fmt.Errorf("field %q: no member %q", fieldPath, part)
		}
		shape = next
	}
	return out, nil
}

func appendUniquePath(paths []string, p string) []string {
	for _, existing := range paths {
		if existing == p {
			return paths
		}
	}
	return append(paths, p)
}
