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

package command

import (
	"fmt"
	"sort"
	"strings"

	k8sversion "k8s.io/apimachinery/pkg/version"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	ackmetadata "github.com/aws-controllers-k8s/code-generator/pkg/metadata"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	acksdk "github.com/aws-controllers-k8s/code-generator/pkg/sdk"
)

// loadModelWithLatestAPIVersion finds the AWS SDK for a given service alias and
// creates a new model with the latest API version.
func loadModelWithLatestAPIVersion(svcAlias string, metadata *ackmetadata.ServiceMetadata) (*ackmodel.Model, error) {
	latestAPIVersion, err := getLatestAPIVersion(metadata.APIVersions)
	if err != nil {
		return nil, err
	}
	return loadModel(svcAlias, latestAPIVersion, "", ackgenerate.DefaultConfig)
}

// loadModel finds the AWS SDK for a given service alias and creates a new model
// with the given API version.
func loadModel(svcAlias string, apiVersion string, apiGroup string, defaultCfg ackgenconfig.Config) (*ackmodel.Model, error) {
	cfg, err := ackgenconfig.New(optGeneratorConfigPath, defaultCfg)
	if err != nil {
		return nil, err
	}

	modelName := strings.ToLower(cfg.ModelName)
	if modelName == "" {
		modelName = svcAlias
	}

	sdkHelper := acksdk.NewHelper(sdkDir, cfg)
	sdkAPI, err := sdkHelper.API(modelName)
	if err != nil {
		retryModelName, err := FallBackFindServiceID(sdkDir, svcAlias)
		if err != nil {
			return nil, err
		}
		// Retry using path found by querying service ID
		sdkAPI, err = sdkHelper.API(retryModelName)
		if err != nil {
			return nil, fmt.Errorf("service %s not found", svcAlias)
		}
	}

	if apiGroup != "" {
		sdkAPI.APIGroupSuffix = apiGroup
	}

	m, err := ackmodel.New(
		sdkAPI, svcAlias, apiVersion, cfg,
	)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// getLatestAPIVersion looks in the controller metadata file to determine what
// the latest Kubernetes API version for CRDs exposed by the generated service
// controller.
func getLatestAPIVersion(apiVersions []ackmetadata.ServiceVersion) (string, error) {
	versions := []string{}

	for _, version := range apiVersions {
		versions = append(versions, version.APIVersion)
	}
	sort.Slice(versions, func(i, j int) bool {
		return k8sversion.CompareKubeAwareVersionStrings(versions[i], versions[j]) < 0
	})
	return versions[len(versions)-1], nil
}

// getServiceAccountName gets the service account name from the optional flag passed into ack-generate
func getServiceAccountName() (string, error) {

	if optServiceAccountName != "" {
		return optServiceAccountName, nil
	}

	return "", fmt.Errorf("service account name not set")
}
