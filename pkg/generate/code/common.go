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
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
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
		"Name",
		"Names",
		r.Names.Original + "Name",
		r.Names.Original + "Names",
		r.Names.Original + "Id",
		r.Names.Original + "Ids",
	}

	for _, memberName := range shape.MemberNames() {
		if util.InStrings(memberName, identifierLookup) {
			identifiers = append(identifiers, memberName)
		}
	}

	return identifiers
}

// FindIdentifiersInCRD returns the identifier fields of a given CRD which
// can be singular or plural.
func FindIdentifiersInCRD(
	r *model.CRD) []string {
	var identifiers []string
	if r == nil {
		return identifiers
	}
	identifierLookup := []string{
		"Id",
		"Ids",
		"Name",
		"Names",
		r.Names.Original + "Name",
		r.Names.Original + "Names",
		r.Names.Original + "Id",
		r.Names.Original + "Ids",
	}

	for _, id := range identifierLookup {
		_, found := r.SpecFields[id]
		if !found {
			_, found = r.StatusFields[id]
		}
		if found {
			identifiers = append(identifiers, id)
		}
	}

	return identifiers
}
