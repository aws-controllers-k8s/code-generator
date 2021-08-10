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
	"errors"
	"fmt"

	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
)

// FindIdentifierInShape returns the identifier field of a given shape which
// can be singular or plural. Errors iff multiple identifier fields detected
// in the shape.
func FindIdentifierInShape(
	r *model.CRD,
	shape *awssdkmodel.Shape) (string, error) {
	identifierLookup := []string{
		"Name",
		"Names",
		r.Names.Original + "Name",
		r.Names.Original + "Names",
		r.Names.Original + "Id",
		r.Names.Original + "Ids",
	}
	identifier := ""

	for _, memberName := range shape.MemberNames() {
		if util.InStrings(memberName, identifierLookup) {
			if identifier == "" {
				identifier = memberName
			} else {
				return "", errors.New(fmt.Sprintf(
					"Found multiple possible identifiers for %s: %s, %s ",
					r.Names.Original, identifier, memberName))
			}
		}
	}

	return identifier, nil
}

// FindIdentifierInCRD returns the identifier field of a given CRD which
// can be singular or plural. Errors iff multiple identifier fields detected
// in the CRD.
func FindIdentifierInCRD(
	r *model.CRD) (string, error) {
	identifierLookup := []string{
		"Name",
		"Names",
		r.Names.Original + "Name",
		r.Names.Original + "Names",
		r.Names.Original + "Id",
		r.Names.Original + "Ids",
	}
	identifier := ""

	for _, id := range identifierLookup {
		_, found := r.SpecFields[id]
		if !found {
			_, found = r.StatusFields[id]
		}
		if found {
			if identifier == "" {
				identifier = id
			} else {
				return "", errors.New(fmt.Sprintf(
					"Found multiple possible identifiers for %s: %s, %s ",
					r.Names.Original, identifier, id))
			}
		}
	}

	return identifier, nil
}
