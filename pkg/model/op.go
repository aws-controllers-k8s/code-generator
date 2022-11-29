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

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
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
				return OpTypeList, resName
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
	switch strings.ToLower(s) {
	case "create":
		return OpTypeCreate
	case "createbatch":
		return OpTypeCreateBatch
	case "delete":
		return OpTypeDelete
	case "replace":
		return OpTypeReplace
	case "update":
		return OpTypeUpdate
	case "addchild":
		return OpTypeAddChild
	case "addchildren":
		return OpTypeAddChildren
	case "removechild":
		return OpTypeRemoveChild
	case "removechildren":
		return OpTypeRemoveChildren
	case "get", "readone", "read_one":
		return OpTypeGet
	case "list", "readmany", "read_many":
		return OpTypeList
	case "getattributes", "get_attributes":
		return OpTypeGetAttributes
	case "setattributes", "set_attributes":
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
