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

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/gertd/go-pluralize"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

var (
	PrimaryIdentifierARNOverride = "ARN"
)

// FindIdentifiersInShape returns the identifier fields of a given shape which
// can be singular or plural.
func FindIdentifiersInShape(
	r *model.CRD,
	shape *awssdkmodel.Shape,
	op *awssdkmodel.Operation,
) []string {
	var identifiers []string
	if r == nil || shape == nil {
		return identifiers
	}

	identifierLookup := []string{
		"Id",
		"Ids",
		r.Names.Original + "Id",
		r.Names.Original + "Ids",
		"Name",
		"Names",
		r.Names.Original + "Name",
		r.Names.Original + "Names",
	}

	// Handles field renames
	opType, _ := model.GetOpTypeAndResourceNameFromOpID(op.Name, r.Config())
	renames, _ := r.GetAllRenames(opType)
	for _, memberName := range shape.MemberNames() {
		lookupName := memberName
		if renamedName, found := renames[memberName]; found {
			lookupName = renamedName
		}
		if util.InStrings(lookupName, identifierLookup) {
			identifiers = append(identifiers, lookupName)
		}
	}

	return identifiers
}

// FindARNIdentifiersInShape returns the identifier fields of a given shape which
// fit expect an ARN.
func FindARNIdentifiersInShape(
	r *model.CRD,
	shape *awssdkmodel.Shape,
) []string {
	var identifiers []string
	if r == nil || shape == nil {
		return identifiers
	}

	for _, memberName := range shape.MemberNames() {
		if r.IsPrimaryARNField(memberName) {
			identifiers = append(identifiers, memberName)
		}
	}

	return identifiers
}

// FindPluralizedIdentifiersInShape returns the name of a Spec OR Status field
// that has a matching pluralized field in the given shape and the name of
// the corresponding shape field name. This method handles identifier field
// renames and will return the same, when applicable.
// For example, DescribeVpcsInput has a `VpcIds` field which would be matched
// to the `Status.VPCID` CRD field - the return value would be
// "VPCID", "VpcIds".
func FindPluralizedIdentifiersInShape(
	r *model.CRD,
	shape *awssdkmodel.Shape,
	op *awssdkmodel.Operation,
) (crField string, shapeField string) {
	shapeIdentifiers := FindIdentifiersInShape(r, shape, op)
	crIdentifiers := r.GetIdentifiers()
	if len(shapeIdentifiers) == 0 || len(crIdentifiers) == 0 {
		return "", ""
	}

	pluralize := pluralize.NewClient()
	for _, si := range shapeIdentifiers {
		for _, ci := range crIdentifiers {
			// If the identifier field is renamed, we must take that into
			// consideration in order to find the corresponding matching
			// shapeIdentifier.
			siRenamed, _ := r.Config().ResourceFieldRename(
				r.Names.Original,
				op.Name,
				pluralize.Singular(si),
			)
			if strings.EqualFold(
				siRenamed,
				pluralize.Singular(ci),
			) {
				// The CRD identifiers being used for comparison reflect any
				// renamed field names in the API model shape.
				if crField == "" {
					crField = ci
					shapeField = si
				} else {
					// If there are multiple identifiers, then prioritize the
					// 'Id' identifier. Checking 'Id' to determine resource
					// creation should be safe as the field is usually
					// present in CR.Status.
					if !strings.Contains(crField, "Id") {
						crField = ci
						shapeField = si
					}
				}
			}
		}
	}
	return crField, shapeField
}

// FindPrimaryIdentifierFieldNames returns the resource identifier field name
// for the primary identifier used in a given operation and its corresponding
// shape field name.
func FindPrimaryIdentifierFieldNames(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	op *awssdkmodel.Operation,
) (crField string, shapeField string) {
	shape := op.InputRef.Shape

	if shapeField == "" {
		// For ReadOne, search for a direct identifier
		if op == r.Ops.ReadOne {
			identifiers := FindIdentifiersInShape(r, shape, op)
			identifiers = append(identifiers, FindARNIdentifiersInShape(r, shape)...)

			switch len(identifiers) {
			case 0:
				break
			case 1:
				shapeField = identifiers[0]
			default:
				panic("Found multiple possible primary identifiers for " +
					r.Names.Original + ". Set " +
					"`is_primary_key` for the primary field in the " +
					r.Names.Camel + " resource.")
			}
		} else {
			// For ReadMany, search for pluralized identifiers
			crField, shapeField = FindPluralizedIdentifiersInShape(r, shape, op)
		}

		// Require override if still can't find any identifiers
		if shapeField == "" {
			panic("Could not find primary identifier for " + r.Names.Original +
				". Set `is_primary_key` for the primary field in the " +
				r.Names.Camel + " resource.")
		}
	}

	if r.IsPrimaryARNField(shapeField) {
		return "", PrimaryIdentifierARNOverride
	}

	if crField == "" {
		if inSpec, inStat := r.HasMember(shapeField, op.Name); !inSpec && !inStat {
			panic("Could not find corresponding spec or status field " +
				"for primary identifier " + shapeField)
		}
		crField, _ = cfg.ResourceFieldRename(
			r.Names.Original,
			op.Name,
			shapeField,
		)
	}
	return crField, shapeField
}

// GetMemberIndex returns the index of memberName by iterating over
// shape's slice of member names for deterministic ordering
func GetMemberIndex(shape *awssdkmodel.Shape, memberName string) (int, error) {
	for index, shapeMemberName := range shape.MemberNames() {
		if strings.EqualFold(shapeMemberName, memberName) {
			return index, nil
		}
	}
	return -1, fmt.Errorf("Could not find %s in shape %s", memberName, shape.ShapeName)
}
