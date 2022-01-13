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

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/gertd/go-pluralize"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
)

type OpType int

const (
	OpTypeUnknown OpType = iota
	OpTypeCreate
	OpTypeCreateBatch
	OpTypeDelete
	OpTypeReplace
	OpTypeUpdate
	OpTypeAddChild
	OpTypeAddChildren
	OpTypeRemoveChild
	OpTypeRemoveChildren
	OpTypeGet
	OpTypeList
	OpTypeGetAttributes
	OpTypeSetAttributes
)

type OperationMap map[OpType]map[string]*awssdkmodel.Operation

// resourceExistsInConfig returns true if the supplied resource name exists in
// the supplied Config object's Resources collection, case-insensitive
// matching.
func resourceExistsInConfig(
	subject string,
	cfg *ackgenconfig.Config,
) bool {
	for resName := range cfg.Resources {
		if strings.EqualFold(resName, subject) {
			return true
		}
	}
	return false
}

// GetOpTypeAndResourceNameFromOpID guesses the resource name and type of
// operation from the OperationID
func GetOpTypeAndResourceNameFromOpID(
	opID string,
	cfg *ackgenconfig.Config,
) (OpType, string) {
	pluralize := pluralize.NewClient()
	if strings.HasPrefix(opID, "CreateOrUpdate") {
		return OpTypeReplace, strings.TrimPrefix(opID, "CreateOrUpdate")
	} else if strings.HasPrefix(opID, "BatchCreate") {
		resName := strings.TrimPrefix(opID, "BatchCreate")
		if pluralize.IsPlural(resName) {
			// Do not singularize "pluralized singular" resources
			// like EC2's DhcpOptions, if defined in generator config's list of
			// resources.
			if resourceExistsInConfig(resName, cfg) {
				return OpTypeCreateBatch, resName
			}
			return OpTypeCreateBatch, pluralize.Singular(resName)
		}
		return OpTypeCreateBatch, resName
	} else if strings.HasPrefix(opID, "CreateBatch") {
		resName := strings.TrimPrefix(opID, "CreateBatch")
		if pluralize.IsPlural(resName) {
			if resourceExistsInConfig(resName, cfg) {
				return OpTypeCreateBatch, resName
			}
			return OpTypeCreateBatch, pluralize.Singular(resName)
		}
		return OpTypeCreateBatch, resName
	} else if strings.HasPrefix(opID, "Create") {
		resName := strings.TrimPrefix(opID, "Create")
		if pluralize.IsPlural(resName) {
			// If resName exists in the generator configuration's list of
			// resources, then just return OpTypeCreate and the resource name.
			// This handles "pluralized singular" resource names like EC2's
			// DhcpOptions.
			if resourceExistsInConfig(resName, cfg) {
				return OpTypeCreate, resName
			}
			return OpTypeCreateBatch, pluralize.Singular(resName)
		}
		return OpTypeCreate, resName
	} else if strings.HasPrefix(opID, "Modify") {
		return OpTypeUpdate, strings.TrimPrefix(opID, "Modify")
	} else if strings.HasPrefix(opID, "Update") {
		return OpTypeUpdate, strings.TrimPrefix(opID, "Update")
	} else if strings.HasPrefix(opID, "Delete") {
		return OpTypeDelete, strings.TrimPrefix(opID, "Delete")
	} else if strings.HasPrefix(opID, "Describe") {
		resName := strings.TrimPrefix(opID, "Describe")
		if pluralize.IsPlural(resName) {
			if resourceExistsInConfig(resName, cfg) {
				return OpTypeGet, resName
			}
			return OpTypeList, pluralize.Singular(resName)
		}
		return OpTypeGet, resName
	} else if strings.HasPrefix(opID, "Get") {
		if strings.HasSuffix(opID, "Attributes") {
			resName := strings.TrimPrefix(opID, "Get")
			resName = strings.TrimSuffix(resName, "Attributes")
			return OpTypeGetAttributes, resName
		}
		resName := strings.TrimPrefix(opID, "Get")
		if pluralize.IsPlural(resName) {
			if resourceExistsInConfig(resName, cfg) {
				return OpTypeGet, resName
			}
			return OpTypeList, pluralize.Singular(resName)
		}
		return OpTypeGet, resName
	} else if strings.HasPrefix(opID, "List") {
		resName := strings.TrimPrefix(opID, "List")
		if pluralize.IsPlural(resName) {
			if resourceExistsInConfig(resName, cfg) {
				return OpTypeList, resName
			}
			return OpTypeList, pluralize.Singular(resName)
		}
		return OpTypeList, resName
	} else if strings.HasPrefix(opID, "Set") {
		if strings.HasSuffix(opID, "Attributes") {
			resName := strings.TrimPrefix(opID, "Set")
			resName = strings.TrimSuffix(resName, "Attributes")
			return OpTypeSetAttributes, resName
		}
	}
	return OpTypeUnknown, opID
}

func OpTypeFromString(s string) OpType {
	switch s {
	case "Create":
		return OpTypeCreate
	case "CreateBatch":
		return OpTypeCreateBatch
	case "Delete":
		return OpTypeDelete
	case "Replace":
		return OpTypeReplace
	case "Update":
		return OpTypeUpdate
	case "AddChild":
		return OpTypeAddChild
	case "AddChildren":
		return OpTypeAddChildren
	case "RemoveChild":
		return OpTypeRemoveChild
	case "RemoveChildren":
		return OpTypeRemoveChildren
	case "Get":
		return OpTypeGet
	case "List":
		return OpTypeList
	case "GetAttributes":
		return OpTypeGetAttributes
	case "SetAttributes":
		return OpTypeSetAttributes
	}

	return OpTypeUnknown
}

// ResourceManagerMethodFromOpType returns the string representing the
// AWSResourceManager method ("Create", "Update", "Delete" or "ReadOne")
// corresponding to the supplied OpType.
func ResourceManagerMethodFromOpType(opType OpType) string {
	switch opType {
	case OpTypeCreate:
		return "Create"
	case OpTypeUpdate:
		return "Update"
	case OpTypeDelete:
		return "Delete"
	case OpTypeGet, OpTypeList:
		return "ReadOne"
	default:
		return ""
	}
}
