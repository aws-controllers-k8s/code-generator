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
	"strings"
)

// GetTagFieldName returns the name of field containing AWS tags. The default
// name is "Tags". If no tag field is found inside the CRD, an error is returned
func (r *CRD) GetTagFieldName() (string, error) {
	tagFieldPath := "Tags" // default tag field path

	// If there is an explicit tag field path mentioned inside generator config,
	// retrieve it.
	if resConfig := r.cfg.GetResourceConfig(r.Names.Original); resConfig != nil {
		if tagConfig := resConfig.TagConfig; tagConfig != nil {
			if tagConfig.Path != nil {
				tagFieldPath = *tagConfig.Path
			}
		}
	}
	// verify that the tagFieldPath exists inside CRD
	for fName, field := range r.Fields {
		if strings.EqualFold(field.Path, tagFieldPath) {
			return fName, nil
		}
	}
	// If the tagFieldPath did not exist in CRD, return an error
	return "", fmt.Errorf("tag field path %s does not exist inside %s"+
		" crd", tagFieldPath, r.Names.Original)
}

// GetTagField return the model.Field representing the Tag field for CRD. If no
// such field is found an error is returned.
func (r *CRD) GetTagField() (*Field, error) {
	if r.cfg.TagsAreIgnored(r.Names.Original) {
		return nil, nil
	}
	tagFieldName, err := r.GetTagFieldName()
	if err != nil {
		return nil, err
	}
	tagField := r.Fields[tagFieldName]
	return tagField, nil
}

// GetTagKeyMemberName returns the member name which represents tag key.
// The default value is "Key". This is only applicable for tag fields with
// shape list of struct.If the tag field shape is not list of struct, an empty
// string is returned.
func (r *CRD) GetTagKeyMemberName() (keyMemberName string) {
	tagField, _ := r.GetTagField()
	if tagField == nil {
		return keyMemberName
	}
	// TagKey member field will only be present when the tag field shape is
	// list of struct
	if isListOfStruct(tagField) {
		keyMemberName = "Key"
		if resConfig := r.Config().GetResourceConfig(r.Names.Original); resConfig != nil {
			if tagsConfig := resConfig.TagConfig; tagsConfig != nil {
				if tagsConfig.KeyMemberName != nil && *tagsConfig.KeyMemberName != "" {
					keyMemberName = *tagsConfig.KeyMemberName
				}
			}
		}
	}
	return keyMemberName
}

// GetTagValueMemberName returns the member name which represents tag value.
// The default value is "Value". This is only applicable for tag fields with
// shape list of struct.If the tag field shape is not list of struct, an empty
// string is returned.
func (r *CRD) GetTagValueMemberName() (valueMemberName string) {
	tagField, _ := r.GetTagField()
	if tagField == nil {
		return valueMemberName
	}
	// TagValue member field will only be present when the tag field shape is
	// list of struct.
	if isListOfStruct(tagField) {
		valueMemberName = "Value"
		if resConfig := r.Config().GetResourceConfig(r.Names.Original); resConfig != nil {
			if tagsConfig := resConfig.TagConfig; tagsConfig != nil {
				if tagsConfig.ValueMemberName != nil && *tagsConfig.ValueMemberName != "" {
					valueMemberName = *tagsConfig.ValueMemberName
				}
			}
		}
	}
	return valueMemberName
}

// isListOfStruct method returns true is the shape of field is list of struct,
// false otherwise
func isListOfStruct(f *Field) bool {
	if f.ShapeRef != nil && f.ShapeRef.Shape != nil {
		if f.ShapeRef.Shape.Type == "list" {
			if f.ShapeRef.Shape.MemberRef.Shape != nil {
				if f.ShapeRef.Shape.MemberRef.Shape.Type == "structure" {
					return true
				}
			}
		}
	}
	return false
}
