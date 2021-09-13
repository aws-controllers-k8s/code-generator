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
	shape *awssdkmodel.Shape) []string {
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

	for _, memberName := range shape.MemberNames() {
		if util.InStrings(memberName, identifierLookup) {
			identifiers = append(identifiers, memberName)
		}
	}

	return identifiers
}

// FindIdentifiersInShape returns the identifier fields of a given shape which
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
// the corresponding shape field name. This method will returns the original
// field name - renames will need to be handled separately.
// For example, DescribeVpcsInput has a `VpcIds` field which would be matched
// to the `Status.VPCID` CRD field - the return value would be
// "VPCID", "VpcIds".
func FindPluralizedIdentifiersInShape(
	r *model.CRD,
	shape *awssdkmodel.Shape,
) (crField string, shapeField string) {
	shapeIdentifiers := FindIdentifiersInShape(r, shape)
	crIdentifiers := r.GetIdentifiers()
	if len(shapeIdentifiers) == 0 || len(crIdentifiers) == 0 {
		return "", ""
	}

	pluralize := pluralize.NewClient()
	for _, si := range shapeIdentifiers {
		for _, ci := range crIdentifiers {
			if strings.EqualFold(pluralize.Singular(si),
				pluralize.Singular(ci)) {
				// The CRD identifiers being used for comparison reflect the
				// *original* field names in the API model shape.
				if crField == "" {
					crField = ci
					shapeField = si
				} else {
					// If there are multiple identifiers, then prioritize the
					// 'Id' identifier. Checking 'Id' to determine resource
					// creation should be safe as the field is usually
					// present in CR.Status.
					if !strings.HasSuffix(crField, "Id") ||
						!strings.HasSuffix(crField, "Ids") {
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
			identifiers := FindIdentifiersInShape(r, shape)
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
			crField, shapeField = FindPluralizedIdentifiersInShape(r, shape)
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
		renamedName, _ := r.InputFieldRename(
			op.Name, shapeField,
		)

		_, inSpec := r.SpecFields[renamedName]
		_, inStatus := r.StatusFields[renamedName]
		if inSpec || inStatus {
			crField = renamedName
		} else {
			panic("Could not find corresponding spec or status field for primary identifier " + shapeField)
		}
	}

	return crField, shapeField
}
