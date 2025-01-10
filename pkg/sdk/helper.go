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
	"path/filepath"

	"github.com/go-git/go-git/v5"

	"github.com/aws-controllers-k8s/code-generator/pkg/apiv2"
	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
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

// Helper is a helper struct that helps work with the aws-sdk-go-v2 models and
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

// APIV2 returns the aws-sdk-go-v2 API model for a supplied service model name.
func (h *Helper) APIV2(serviceModelName string) (*model.SDKAPI, error) {
	modelPath := h.ModelAndDocsPath(serviceModelName)
	apis, err := apiv2.ConvertApiV2Shapes(serviceModelName, modelPath)
	if err != nil {
		return nil, err
	}
	// apis is a map, keyed by the service alias, of pointers to aws-sdk-go-v2
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
func (h *Helper) ModelAndDocsPath(serviceModelName string) string {
	modelPath := filepath.Join(
		h.basePath,
		"codegen",
		"sdk-codegen",
		"aws-models",
		fmt.Sprintf("%s.json", serviceModelName),
	)
	return modelPath
}
