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

package multiversion

import (
	"fmt"
	"sort"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// FieldChangeType represents the type of field modification.
//
//   - FieldChangeTypeUnknown is used when ChangeType cannot be computed.
//   - FieldChangeTypeNone is used when a field name and structure didn't change.
//   - FieldChangeTypeAdded is used when a new field is introduced in a CRD.
//   - FieldChangeTypeRemoved is used a when a field is removed from a CRD.
//   - FieldChangeTypeRenamed is used when a field is renamed.
//   - FieldChangeTypeShapeChanged is used when a field shape has changed.
//   - FieldChangeTypeShapeChangedFromStringToSecret is used when a field change to
//     a k8s secret type.
//   - FieldChangeTypeShapeChangedFromSecretToString is used when a field changed from
//     a k8s secret to a Go string.
type FieldChangeType string

const (
	FieldChangeTypeUnknown                        FieldChangeType = "unknown"
	FieldChangeTypeNone                           FieldChangeType = "none"
	FieldChangeTypeAdded                          FieldChangeType = "added"
	FieldChangeTypeRemoved                        FieldChangeType = "removed"
	FieldChangeTypeRenamed                        FieldChangeType = "renamed"
	FieldChangeTypeShapeChanged                   FieldChangeType = "shape-changed"
	FieldChangeTypeShapeChangedFromStringToSecret FieldChangeType = "shape-changed-from-string-to-secret"
	FieldChangeTypeShapeChangedFromSecretToString FieldChangeType = "shape-changed-from-secret-to-string"
)

// FieldDelta represents the delta between the same original field in two
// different CRD versions. If a field is removed in the Destination version
// the Destination value will be nil. If a field is new in the Destination
// version, the Source value will be nil.
type FieldDelta struct {
	ChangeType FieldChangeType
	// Field from the destination CRD
	Destination *ackmodel.Field
	// Field from the source CRD
	Source *ackmodel.Field
}

// CRDDelta stores the spec and status deltas for a custom resource.
type CRDDelta struct {
	SpecDeltas   []FieldDelta
	StatusDeltas []FieldDelta
}

// ComputeCRDFieldDeltas compares two ackmodel.CRD instances and returns the
// spec and status fields deltas. src is the CRD of the spoke (source) version
// and dst is the CRD of the hub (destination) version.
func ComputeCRDFieldDeltas(src, dst *ackmodel.CRD) (*CRDDelta, error) {
	dstRenames := dst.GetAllRenames(ackmodel.OpTypeCreate)
	srcRenames := src.GetAllRenames(ackmodel.OpTypeCreate)

	renames, err := ComputeRenamesDelta(srcRenames, dstRenames)
	if err != nil {
		return nil, fmt.Errorf("cannot compute the field renames delta: %v", err)
	}

	specDeltas, err := ComputeFieldDeltas(src.SpecFields, dst.SpecFields, renames)
	if err != nil {
		return nil, fmt.Errorf("cannot compute spec fields deltas: %s", err)
	}

	statusDeltas, err := ComputeFieldDeltas(src.StatusFields, dst.StatusFields, renames)
	if err != nil {
		return nil, fmt.Errorf("cannot compute status fields deltas: %s", err)
	}

	return &CRDDelta{
		specDeltas,
		statusDeltas,
	}, nil
}

// fieldChangedToSecret returns true if field changed from string to secret.
func fieldChangedToSecret(src, dst *ackmodel.Field) bool {
	return (src.FieldConfig == nil ||
		(src.FieldConfig != nil && !src.FieldConfig.IsSecret)) &&
		(dst.FieldConfig != nil && dst.FieldConfig.IsSecret)
}

// ComputeFieldDeltas computes the difference between two maps of fields. It returns a list
// of FieldDelta's that contains the ChangeType and at least one field reference.
func ComputeFieldDeltas(
	srcFields map[string]*ackmodel.Field,
	dstFields map[string]*ackmodel.Field,
	// the renames delta renames
	renames map[string]string,
) ([]FieldDelta, error) {
	deltas := []FieldDelta{}

	// collect field names and sort them to ensure a determenistic output order.
	srcNames := []string{}
	for name := range srcFields {
		srcNames = append(srcNames, name)
	}
	sort.Strings(srcNames)

	dstNames := []string{}
	for name := range dstFields {
		dstNames = append(dstNames, name)
	}
	sort.Strings(dstNames)

	// let's make sure we don't visit fields more than once - especially
	// when fields are renamed.
	visitedFields := map[string]bool{}

	// first let's loop over the srcNames array and see if we can find
	// the same field name in dstNames.
	for _, srcName := range srcNames {
		srcField, _ := srcFields[srcName]
		dstField, ok := dstFields[srcName]
		// If a field is found in both arrays only three changes are possible:
		// None, TypeChange and ChangeTypeShapeChangedToSecret.
		// NOTE(a-hilaly): carefull about X -> Y then Z -> X renames. It should
		// not be allowed.
		if ok {
			// mark field as visited.
			visitedFields[srcName] = true
			// check if field change from string to secret
			if fieldChangedToSecret(srcField, dstField) {
				deltas = append(deltas, FieldDelta{
					Source:      srcField,
					Destination: dstField,
					ChangeType:  FieldChangeTypeShapeChangedFromStringToSecret,
				})
				continue
			}
			// check if field changed from secret to string
			if fieldChangedToSecret(dstField, srcField) {
				deltas = append(deltas, FieldDelta{
					Source:      srcField,
					Destination: dstField,
					ChangeType:  FieldChangeTypeShapeChangedFromSecretToString,
				})
				continue
			}

			equalShapes, _ := AreEqualShapes(srcField.ShapeRef.Shape, dstField.ShapeRef.Shape, true)
			if equalShapes {
				// if the fields have equal names and types the change is intact
				deltas = append(deltas, FieldDelta{
					Source:      srcField,
					Destination: dstField,
					ChangeType:  FieldChangeTypeNone,
				})
				continue
			}

			// at this point we know that the fields kept the same name but have different
			// shapes
			deltas = append(deltas, FieldDelta{
				Source:      srcField,
				Destination: dstField,
				ChangeType:  FieldChangeTypeShapeChanged,
			})
			continue
		}

		// if a field is not found in the dstNames, there are three
		// possible changes: Removed, Added or Renamed.

		// First let's check if field was renamed
		newName, ok := renames[srcName]
		if ok {
			dstField, ok2 := dstFields[newName]
			if !ok2 {
				// if a field was renamed and we can't find it in dstNames, something
				// very wrong happend during CRD loading.
				return nil, fmt.Errorf("cannot find renamed field %s " + newName)
			}

			// mark field as visited, both old and new names.
			visitedFields[newName] = true
			visitedFields[srcName] = true

			// this will mostlikely be always true, but let's double check.
			if newName == dstField.Names.Camel {
				// field was renamed
				deltas = append(deltas, FieldDelta{
					Source:      srcField,
					Destination: dstField,
					ChangeType:  FieldChangeTypeRenamed,
				})
				continue
			}
			return nil, fmt.Errorf("renamed field unmatching: %v != %v", newName, dstField.Names.Camel)
		}

		// If the field was not renamed nor left intact nor it shape changed, it's
		// a removed field.
		deltas = append(deltas, FieldDelta{
			Source:      srcField,
			Destination: nil,
			ChangeType:  FieldChangeTypeRemoved,
		})
	}

	// At this point we collected every type of changes except added fields.
	// To find added fields we just look for fields that are in dstNames and
	// were not visited before (are not in srcNames).
	for _, dstName := range dstNames {
		if _, visited := visitedFields[dstName]; visited {
			continue
		}
		dstField, _ := dstFields[dstName]
		deltas = append(deltas, FieldDelta{
			Source:      nil,
			Destination: dstField,
			ChangeType:  FieldChangeTypeAdded,
		})
	}

	sort.Slice(deltas, func(i, j int) bool {
		return getFieldNameFromDelta(deltas[i]) < getFieldNameFromDelta(deltas[j])
	})
	return deltas, nil
}

// getFieldNameFromDelta retrieves the field name from a FieldDelta. If a field is
// renamed it will return the new name.
func getFieldNameFromDelta(delta FieldDelta) string {
	if delta.ChangeType == FieldChangeTypeRemoved {
		return delta.Source.Names.Camel
	}
	return delta.Destination.Names.Camel
}

// AreEqualShapes returns whether two awssdkmodel.ShapeRef are equal or not. When the two
// given shapes are not equal, it will return an error representing the first type mismatch
// detected.
func AreEqualShapes(a, b *awssdkmodel.Shape, allowMemberNamesInequality bool) (bool, error) {
	if a.Type != b.Type {
		return false, fmt.Errorf("found different shape types (%s and %s)", a.ShapeName, a.ShapeName)
	}
	if !allowMemberNamesInequality && a.ShapeName != b.ShapeName {
		return false, fmt.Errorf("found different shape names (%s and %s)", a.ShapeName, a.ShapeName)
	}

	switch a.Type {
	case "structure":
		// verify that both structs have the same member names
		if len(a.MemberNames()) != len(a.MemberNames()) {
			return false, fmt.Errorf("found different MemberNames size in %s", a.ShapeName)
		}

		// loop over the struct members and verify that if they have the same
		// name they should have the same shape.
		for _, memberName := range a.MemberNames() {
			memberRefA := a.MemberRefs[memberName]
			memberRefB, ok := b.MemberRefs[memberName]
			if !ok {
				return false, fmt.Errorf("missing member %s in %s", memberName, memberRefA.ShapeName)
			}
			// if two members with the same name doesn't have the same shape
			// return false.
			if equal, err := AreEqualShapes(memberRefA.Shape, memberRefB.Shape, false); !equal {
				return false, fmt.Errorf("member %s have two different shapes in %s: %v", memberName, memberRefA.ShapeName, err)
			}
		}
	case "map":
		// for maps we check that the keys and values have the same types
		if equal, err := AreEqualShapes(a.KeyRef.Shape, b.KeyRef.Shape, false); !equal {
			return false, fmt.Errorf("map key shape mismatch in %s: %v", a.ShapeName, err)
		}
		if equal, err := AreEqualShapes(a.ValueRef.Shape, b.ValueRef.Shape, false); !equal {
			return false, fmt.Errorf("map value shape mismatch in %s: %v", a.ShapeName, err)
		}
	case "list":
		// for lists we check that the members have the same types
		if equal, err := AreEqualShapes(a.MemberRef.Shape, b.MemberRef.Shape, false); !equal {
			return false, fmt.Errorf("member shape mismatch in %s: %v", a.ShapeName, err)
		}
	}

	return true, nil
}
