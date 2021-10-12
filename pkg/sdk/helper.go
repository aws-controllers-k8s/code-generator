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

package sdk

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/src-d/go-git.v4"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
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

// Helper is a helper struct that helps work with the aws-sdk-go models and
// API model loader
type Helper struct {
	// Default is "services.k8s.aws"
	APIGroupSuffix string
	cfg            ackgenconfig.Config
	gitRepository  *git.Repository
	basePath       string
	loader         *awssdkmodel.Loader
	// Default is set by `FirstAPIVersion`
	apiVersion string
}

// NewHelper returns a new SDKHelper object
func NewHelper(basePath string, cfg ackgenconfig.Config) *Helper {
	return &Helper{
		cfg:      cfg,
		basePath: basePath,
		loader: &awssdkmodel.Loader{
			BaseImport:            basePath,
			IgnoreUnsupportedAPIs: true,
		},
	}
}

// WithSDKVersion checks out the sdk git repository to the provided version. To use
// this function h.basePath should point to a git repository.
func (h *Helper) WithSDKVersion(version string) error {
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
func (h *Helper) WithAPIVersion(apiVersion string) {
	h.apiVersion = apiVersion
}

// API returns the aws-sdk-go API model for a supplied service model name.
func (h *Helper) API(serviceModelName string) (*model.SDKAPI, error) {
	modelPath, _, err := h.ModelAndDocsPath(serviceModelName)
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
		sdkapi := model.NewSDKAPI(api, h.APIGroupSuffix)

		h.InjectCustomShapes(sdkapi)

		return sdkapi, nil
	}
	return nil, ErrServiceNotFound
}

// ModelAndDocsPath returns two string paths to the supplied service's API and
// doc JSON files
func (h *Helper) ModelAndDocsPath(
	serviceModelName string,
) (string, string, error) {
	if h.apiVersion == "" {
		apiVersion, err := h.FirstAPIVersion(serviceModelName)
		if err != nil {
			return "", "", err
		}
		h.apiVersion = apiVersion
	}
	versionPath := filepath.Join(
		h.basePath, "models", "apis", serviceModelName, h.apiVersion,
	)
	modelPath := filepath.Join(versionPath, "api-2.json")
	docsPath := filepath.Join(versionPath, "docs-2.json")
	return modelPath, docsPath, nil
}

// FirstAPIVersion returns the first found API version for a service API.
// (e.h. "2012-10-03")
func (h *Helper) FirstAPIVersion(serviceModelName string) (string, error) {
	versions, err := h.GetAPIVersions(serviceModelName)
	if err != nil {
		return "", err
	}
	sort.Strings(versions)
	return versions[0], nil
}

// GetAPIVersions returns the list of API Versions found in a service directory.
func (h *Helper) GetAPIVersions(serviceModelName string) ([]string, error) {
	apiPath := filepath.Join(h.basePath, "models", "apis", serviceModelName)
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
