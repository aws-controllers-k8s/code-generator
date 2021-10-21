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

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/names"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

const (
	// ConflictingNameSuffix is appended to type names when they overlap with
	// well-known common struct names for things like a CRD itself, or its
	// Spec/Status subfield struct type name.
	ConflictingNameSuffix = "_SDK"
)

// SDKAPI contains an API model for a single AWS service API
type SDKAPI struct {
	API            *awssdkmodel.API
	APIGroupSuffix string
	// A map of operation type and resource name to
	// aws-sdk-go/private/model/api.Operation structs
	opMap *OperationMap
	// Map, keyed by original Shape GoTypeElem(), with the values being a
	// renamed type name (due to conflicting names)
	typeRenames map[string]string
	// Default is "services.k8s.aws"
}

// GetPayloads returns a slice of strings of Shape names representing input and
// output request/response payloads
func (a *SDKAPI) GetPayloads() []string {
	res := []string{}
	for _, op := range a.API.Operations {
		res = append(res, op.InputRef.ShapeName)
		res = append(res, op.OutputRef.ShapeName)
	}
	return res
}

// GetOperationMap returns a map, keyed by the operation type and operation
// ID/name, of aws-sdk-go private/model/api.Operation struct pointers
func (a *SDKAPI) GetOperationMap(cfg *ackgenconfig.Config) *OperationMap {
	if a.opMap != nil {
		return a.opMap
	}
	// create an index of Operations by operation types and resource name
	opMap := OperationMap{}
	for opID, op := range a.API.Operations {
		opTypeArray, resName := getOpTypeAndResourceName(opID, cfg)
		for _, opType := range opTypeArray {
			if _, found := opMap[opType]; !found {
				opMap[opType] = map[string]*awssdkmodel.Operation{}
			}
			opMap[opType][resName] = op
		}
	}
	a.opMap = &opMap
	return &opMap
}

// GetInputShapeRef finds a ShapeRef for a supplied member path (dot-notation)
// for given API operation
func (a *SDKAPI) GetInputShapeRef(
	opID string,
	path string,
) (*awssdkmodel.ShapeRef, bool) {
	op, ok := a.API.Operations[opID]
	if !ok {
		return nil, false
	}
	return getMemberByPath(op.InputRef.Shape, path)
}

// GetOutputShapeRef finds a ShapeRef for a supplied member path (dot-notation)
// for given API operation
func (a *SDKAPI) GetOutputShapeRef(
	opID string,
	path string,
) (*awssdkmodel.ShapeRef, bool) {
	op, ok := a.API.Operations[opID]
	if !ok {
		return nil, false
	}
	return getMemberByPath(op.OutputRef.Shape, path)
}

// CRDNames returns a slice of names structs for all top-level resources in the
// API
func (a *SDKAPI) CRDNames(cfg *ackgenconfig.Config) []names.Names {
	opMap := a.GetOperationMap(cfg)
	createOps := (*opMap)[OpTypeCreate]
	crdNames := []names.Names{}
	for crdName := range createOps {
		if cfg.IsIgnoredResource(crdName) {
			continue
		}
		crdNames = append(crdNames, names.New(crdName))
	}
	return crdNames
}

// GetTypeRenames returns a map of original type name to renamed name (some
// type definition names conflict with generated names)
func (a *SDKAPI) GetTypeRenames(cfg *ackgenconfig.Config) map[string]string {
	if a.typeRenames != nil {
		return a.typeRenames
	}

	trenames := map[string]string{}

	payloads := a.GetPayloads()

	for shapeName, shape := range a.API.Shapes {
		if util.InStrings(shapeName, payloads) {
			// Payloads are not type defs
			continue
		}
		if shape.Type != "structure" {
			continue
		}
		if shape.Exception {
			// Neither are exceptions
			continue
		}
		tdefNames := names.New(shapeName)
		if a.HasConflictingTypeName(shapeName, cfg) {
			tdefNames.Camel += ConflictingNameSuffix
			trenames[shapeName] = tdefNames.Camel
		}
	}
	a.typeRenames = trenames
	return trenames
}

// HasConflictingTypeName returns true if the supplied type name will conflict
// with any generated type in the service's API package
func (a *SDKAPI) HasConflictingTypeName(typeName string, cfg *ackgenconfig.Config) bool {
	// First grab the set of CRD struct names and the names of their Spec and
	// Status structs
	cleanTypeName := names.New(typeName).Camel
	crdNames := a.CRDNames(cfg)
	crdResourceNames := []string{}
	crdListResourceNames := []string{}
	crdSpecNames := []string{}
	crdStatusNames := []string{}

	for _, crdName := range crdNames {
		cleanResourceName := crdName.Camel
		crdResourceNames = append(crdResourceNames, cleanResourceName)
		crdSpecNames = append(crdSpecNames, cleanResourceName+"Spec")
		crdStatusNames = append(crdStatusNames, cleanResourceName+"Status")
		crdListResourceNames = append(crdListResourceNames, cleanResourceName+"List")
	}
	return util.InStrings(cleanTypeName, crdResourceNames) ||
		util.InStrings(cleanTypeName, crdSpecNames) ||
		util.InStrings(cleanTypeName, crdStatusNames) ||
		util.InStrings(cleanTypeName, crdListResourceNames)
}

// ServiceID returns the exact `metadata.serviceId` attribute for the AWS
// service APi's api-2.json file.
// This MAY NOT MATCH the AWS SDK Go package used by the service. For example:
// AWS SDK Go uses `opensearchservice` whereas the service ID is `opensearch`
func (a *SDKAPI) ServiceID() string {
	if a == nil || a.API == nil {
		return ""
	}
	return awssdkmodel.ServiceID(a.API)
}

func (a *SDKAPI) GetServiceFullName() string {
	if a == nil || a.API == nil {
		return ""
	}
	return a.API.Metadata.ServiceFullName
}

// APIInterfaceTypeName returns the name of the aws-sdk-go primary API
// interface type name.
func (a *SDKAPI) APIInterfaceTypeName() string {
	if a == nil || a.API == nil {
		return ""
	}
	return a.API.StructName()
}

// NewSDKAPI returns a pointer to a new `ackmodel.SDKAPI` struct that describes
// the AWS SDK API and its respective groupings, mappings and renamings.
func NewSDKAPI(api *awssdkmodel.API, apiGroupSuffix string) *SDKAPI {
	return &SDKAPI{
		API:            api,
		APIGroupSuffix: apiGroupSuffix,
		opMap:          nil,
		typeRenames:    nil,
	}
}

// Override the operation type and/or resource name if specified in config
func getOpTypeAndResourceName(opID string, cfg *ackgenconfig.Config) ([]OpType, string) {
	opType, resName := GetOpTypeAndResourceNameFromOpID(opID, cfg)
	opTypes := []OpType{opType}

	if cfg == nil {
		return opTypes, resName
	}
	if operationConfig, exists := cfg.Operations[opID]; exists {
		if operationConfig.ResourceName != "" {
			resName = operationConfig.ResourceName
		}

		for _, operationType := range operationConfig.OperationType {
			opType = OpTypeFromString(operationType)
			opTypes = append(opTypes, opType)
		}
	}
	return opTypes, resName
}

// getMemberByPath returns a ShapeRef given a root Shape and a dot-notation
// object search path. Given the explicit type check for list type members
// both ".." and "." notations work currently.
// TODO: Add support for other types such as map.
func getMemberByPath(
	shape *awssdkmodel.Shape,
	path string,
) (*awssdkmodel.ShapeRef, bool) {
	elements := strings.Split(path, ".")
	last := len(elements) - 1
	for x, elem := range elements {
		if elem == "" {
			continue
		}
		if shape == nil {
			return nil, false
		}
		shapeRef, ok := shape.MemberRefs[elem]
		if !ok {
			return nil, false
		}
		if x == last {
			return shapeRef, true
		}
		elemType := shapeRef.Shape.Type
		switch elemType {
		case "list":
			shape = shapeRef.Shape.MemberRef.Shape
		default:
			shape = shapeRef.Shape
		}
	}
	return nil, false
}
