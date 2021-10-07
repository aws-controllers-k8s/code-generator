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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/src-d/go-git.v4"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/names"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
)

var (
	ErrInvalidVersionDirectory = errors.New(
		"expected to find only directories in api model directory but found non-directory",
	)
	ErrNoValidVersionDirectory = errors.New(
		"no valid version directories found",
	)
	ErrServiceNotFound = errors.New(
		"no such service",
	)
	ErrAPIVersionNotFound = errors.New(
		"no such api version",
	)
)

// SDKHelper is a helper struct that helps work with the aws-sdk-go models and
// API model loader
type SDKHelper struct {
	gitRepository *git.Repository
	basePath      string
	loader        *awssdkmodel.Loader
	// Default is set by `FirstAPIVersion`
	apiVersion string
	// Default is "services.k8s.aws"
	APIGroupSuffix string
}

// NewSDKHelper returns a new SDKHelper object
func NewSDKHelper(basePath string) *SDKHelper {
	return &SDKHelper{
		basePath: basePath,
		loader: &awssdkmodel.Loader{
			BaseImport:            basePath,
			IgnoreUnsupportedAPIs: true,
		},
	}
}

// WithSDKVersion checks out the sdk git repository to the provided version. To use
// this function h.basePath should point to a git repository.
func (h *SDKHelper) WithSDKVersion(version string) error {
	if h.gitRepository == nil {
		gitRepository, err := util.LoadRepository(h.basePath)
		if err != nil {
			return fmt.Errorf("error loading repository from %s: %v", h.basePath, err)
		}
		h.gitRepository = gitRepository
	}

	err := util.CheckoutRepositoryTag(h.gitRepository, version)
	if err != nil {
		return fmt.Errorf("cannot checkout tag %s: %v", version, err)
	}
	return nil
}

// WithAPIVersion sets the `apiVersion` field.
func (h *SDKHelper) WithAPIVersion(apiVersion string) {
	h.apiVersion = apiVersion
}

// API returns the aws-sdk-go API model for a supplied service alias
func (h *SDKHelper) API(serviceAlias string) (*SDKAPI, error) {
	modelPath, _, err := h.ModelAndDocsPath(serviceAlias)
	if err != nil {
		return nil, err
	}
	apis, err := h.loader.Load([]string{modelPath})
	if err != nil {
		return nil, err
	}
	// apis is a map, keyed by the service alias, of pointers to aws-sdk-go
	// model API objects
	for _, api := range apis {
		// If we don't do this, we can end up with panic()'s like this:
		// panic: assignment to entry in nil map
		// when trying to execute Shape.GoType().
		//
		// Calling API.ServicePackageDoc() ends up resetting the API.imports
		// unexported map variable...
		_ = api.ServicePackageDoc()
		return &SDKAPI{api, nil, nil, h.APIGroupSuffix}, nil
	}
	return nil, ErrServiceNotFound
}

// ModelAndDocsPath returns two string paths to the supplied service alias'
// model and doc JSON files
func (h *SDKHelper) ModelAndDocsPath(
	serviceAlias string,
) (string, string, error) {
	if h.apiVersion == "" {
		apiVersion, err := h.FirstAPIVersion(serviceAlias)
		if err != nil {
			return "", "", err
		}
		h.apiVersion = apiVersion
	}
	versionPath := filepath.Join(
		h.basePath, "models", "apis", serviceAlias, h.apiVersion,
	)
	modelPath := filepath.Join(versionPath, "api-2.json")
	docsPath := filepath.Join(versionPath, "docs-2.json")
	return modelPath, docsPath, nil
}

// FirstAPIVersion returns the first found API version for a service API.
// (e.h. "2012-10-03")
func (h *SDKHelper) FirstAPIVersion(serviceAlias string) (string, error) {
	versions, err := h.GetAPIVersions(serviceAlias)
	if err != nil {
		return "", err
	}
	sort.Strings(versions)
	return versions[0], nil
}

// GetAPIVersions returns the list of API Versions found in a service directory.
func (h *SDKHelper) GetAPIVersions(serviceAlias string) ([]string, error) {
	apiPath := filepath.Join(h.basePath, "models", "apis", serviceAlias)
	versionDirs, err := ioutil.ReadDir(apiPath)
	if err != nil {
		return nil, err
	}
	versions := []string{}
	for _, f := range versionDirs {
		version := f.Name()
		fp := filepath.Join(apiPath, version)
		fi, err := os.Lstat(fp)
		if err != nil {
			return nil, err
		}
		if !fi.IsDir() {
			return nil, fmt.Errorf("found %s: %v", version, ErrInvalidVersionDirectory)
		}
		versions = append(versions, version)
	}
	if len(versions) == 0 {
		return nil, ErrNoValidVersionDirectory
	}
	return versions, nil
}

// SDKAPI contains an API model for a single AWS service API
type SDKAPI struct {
	API *awssdkmodel.API
	// A map of operation type and resource name to
	// aws-sdk-go/private/model/api.Operation structs
	opMap *OperationMap
	// Map, keyed by original Shape GoTypeElem(), with the values being a
	// renamed type name (due to conflicting names)
	typeRenames map[string]string
	// Default is "services.k8s.aws"
	apiGroupSuffix string
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
// service APi's api-2.json file
func (a *SDKAPI) ServiceID() string {
	if a == nil || a.API == nil {
		return ""
	}
	return awssdkmodel.ServiceID(a.API)
}

// ServiceIDClean returns a lowercased, whitespace-stripped ServiceID
func (a *SDKAPI) ServiceIDClean() string {
	serviceID := strings.ToLower(a.ServiceID())
	return strings.Replace(serviceID, " ", "", -1)
}

func (a *SDKAPI) GetServiceFullName() string {
	if a == nil || a.API == nil {
		return ""
	}
	return a.API.Metadata.ServiceFullName
}

// APIGroup returns the normalized Kubernetes APIGroup for the AWS service API,
// e.g. "sns.services.k8s.aws"
func (a *SDKAPI) APIGroup() string {
	serviceID := a.ServiceIDClean()
	suffix := "services.k8s.aws"
	if a.apiGroupSuffix != "" {
		suffix = a.apiGroupSuffix
	}
	return fmt.Sprintf("%s.%s", serviceID, suffix)
}

// SDKAPIInterfaceTypeName returns the name of the aws-sdk-go primary API
// interface type name.
func (a *SDKAPI) SDKAPIInterfaceTypeName() string {
	if a == nil || a.API == nil {
		return ""
	}
	return a.API.StructName()
}

// Override the operation type and/or resource name if specified in config
func getOpTypeAndResourceName(opID string, cfg *ackgenconfig.Config) ([]OpType, string) {
	opType, resName := GetOpTypeAndResourceNameFromOpID(opID)
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
