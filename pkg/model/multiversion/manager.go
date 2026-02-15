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
	"errors"
	"fmt"
	"sort"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmetadata "github.com/aws-controllers-k8s/code-generator/pkg/metadata"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	acksdk "github.com/aws-controllers-k8s/code-generator/pkg/sdk"
)

var (
	ErrAPIVersionNotFound   = errors.New("api version not found")
	ErrAPIVersionRemoved    = errors.New("api version removed")
	ErrAPIVersionDeprecated = errors.New("api version deprecated")
)

// APIVersionManager is a API versions manager. It contains the mapping
// of each non-deprecated version with their corresponding ackmodel.Model
// and APIInfos.
type APIVersionManager struct {
	metadata *ackmetadata.ServiceMetadata

	hubVersion    string
	spokeVersions []string

	apiInfos map[string]ackmetadata.APIInfo
	models   map[string]*ackmodel.Model
}

// NewAPIVersionManager initialises and returns a new APIVersionManager.
func NewAPIVersionManager(
	sdkCacheDir string,
	metadataPath string,
	servicePackageName string,
	hubVersion string,
	apisInfo map[string]ackmetadata.APIInfo,
	defaultConfig ackgenconfig.Config,
) (*APIVersionManager, error) {
	if len(apisInfo) == 0 {
		return nil, fmt.Errorf("empty apisInfo")
	}

	metadata, err := ackmetadata.NewServiceMetadata(metadataPath)
	if err != nil {
		return nil, err
	}

	spokeVersions := []string{}

	// create model for each non-deprecated api version
	models := map[string]*ackmodel.Model{}
	for _, version := range metadata.APIVersions {
		if version.Status == ackmetadata.APIStatusDeprecated || version.Status == ackmetadata.APIStatusRemoved {
			continue
		}

		if version.APIVersion != hubVersion {
			spokeVersions = append(spokeVersions, version.APIVersion)
		}

		apiInfo, ok := apisInfo[version.APIVersion]
		if !ok {
			return nil, fmt.Errorf("could not find API info for API version %s", version.APIVersion)
		}

		cfg, err := ackgenconfig.New(apiInfo.GeneratorConfigPath, defaultConfig)
		if err != nil {
			return nil, err
		}

		sdkAPIHelper := acksdk.NewHelper(sdkCacheDir, cfg)
		err = sdkAPIHelper.WithSDKVersion(apiInfo.AWSSDKVersion)
		if err != nil {
			return nil, err
		}

		sdkAPI, err := sdkAPIHelper.API(servicePackageName)
		if err != nil {
			return nil, err
		}

		docCfg, err := ackgenconfig.NewDocumentationConfig(apiInfo.DocumentationConfigPath)
		if err != nil {
			return nil, err
		}

		i, err := ackmodel.New(
			sdkAPI,
			servicePackageName,
			version.APIVersion,
			cfg,
			docCfg,
		)
		if err != nil {
			return nil, fmt.Errorf("cannot create model for API version %s: %v", version.APIVersion, err)
		}
		models[version.APIVersion] = i
	}

	sort.Strings(spokeVersions)
	model := &APIVersionManager{
		metadata:      metadata,
		hubVersion:    hubVersion,
		spokeVersions: spokeVersions,
		apiInfos:      apisInfo,
		models:        models,
	}

	return model, nil
}

// GetModel returns the model of a given api version.
func (m *APIVersionManager) GetModel(apiVersion string) (*ackmodel.Model, error) {
	if err := m.VerifyAPIVersions(apiVersion); err != nil {
		return nil, fmt.Errorf("cannot verify API version %s: %v", apiVersion, err)
	}
	return m.models[apiVersion], nil
}

// GetSpokeVersions returns the spokes versions list.
func (m *APIVersionManager) GetSpokeVersions() []string {
	return m.spokeVersions
}

// GetHubVersion returns the hub version.
func (m *APIVersionManager) GetHubVersion() string {
	return m.hubVersion
}

// CompareHubWith compares a given api version with the hub version and returns
// a string to *CRDDelta map.
func (m *APIVersionManager) CompareHubWith(apiVersion string) (map[string]*CRDDelta, error) {
	return m.CompareAPIVersions(apiVersion, m.hubVersion)
}

// VerifyAPIVersions verifies that an API version exists and is not deprecated.
func (m *APIVersionManager) VerifyAPIVersions(apiVersions ...string) error {
	for _, apiVersion := range apiVersions {
		apiInfo, ok := m.apiInfos[apiVersion]
		if !ok {
			return fmt.Errorf("%v: %s", ErrAPIVersionNotFound, apiVersion)
		}
		if apiInfo.Status == ackmetadata.APIStatusDeprecated {
			return fmt.Errorf("%v: %s", ErrAPIVersionDeprecated, apiVersion)
		}
		if apiInfo.Status == ackmetadata.APIStatusRemoved {
			return fmt.Errorf("%v: %s", ErrAPIVersionRemoved, apiVersion)
		}
	}
	return nil
}

// CompareAPIVersions compares two api versions and returns a slice of FieldDeltas
// representing the diff between CRDs status and spec fields.
func (m *APIVersionManager) CompareAPIVersions(srcAPIVersion, dstAPIVersion string) (
	map[string]*CRDDelta,
	error,
) {
	if srcAPIVersion == dstAPIVersion {
		return nil, fmt.Errorf("cannot compare an API version with it self")
	}

	// get source CRDs
	srcModel, err := m.GetModel(srcAPIVersion)
	if err != nil {
		return nil, err
	}
	srcCRDs, err := srcModel.GetCRDs()
	if err != nil {
		return nil, fmt.Errorf("error getting crds for %s: %v", srcAPIVersion, err)
	}

	// get destination crds
	dstModel, err := m.GetModel(dstAPIVersion)
	if err != nil {
		return nil, err
	}
	dstCRDs, err := dstModel.GetCRDs()
	if err != nil {
		return nil, fmt.Errorf("error getting crds for %s: %v", dstAPIVersion, err)
	}

	// compute FieldDeltas for each CRD
	apiDeltas := make(map[string]*CRDDelta)
	if len(srcCRDs) != len(dstCRDs) {
		// TODO(a-hilaly) handle added/removed CRDs
		return nil, fmt.Errorf("source and destination API versions don't have the same number of CRDs")
	}
	for i, crd := range dstCRDs {
		crdDelta, err := ComputeCRDFieldDeltas(srcCRDs[i], dstCRDs[i])
		if err != nil {
			return nil, fmt.Errorf("cannot compute crd field deltas: %v", err)
		}
		apiDeltas[crd.Names.Camel] = crdDelta
	}
	return apiDeltas, nil
}
