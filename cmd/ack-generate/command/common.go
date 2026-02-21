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
	"time"

	k8sversion "k8s.io/apimachinery/pkg/version"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackmetadata "github.com/aws-controllers-k8s/code-generator/pkg/metadata"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	acksdk "github.com/aws-controllers-k8s/code-generator/pkg/sdk"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

// resolveModelName returns the SDK model name for a service, checking the
// generator config for an override.
func resolveModelName(svcAlias string, cfg ackgenconfig.Config) string {
	modelName := strings.ToLower(cfg.SDKNames.Model)
	if modelName == "" {
		modelName = svcAlias
	}
	return modelName
}

// loadModelWithLatestAPIVersion finds the AWS SDK for a given service alias and
// creates a new model with the latest API version.
func loadModelWithLatestAPIVersion(svcAlias string, metadata *ackmetadata.ServiceMetadata, cfg ackgenconfig.Config) (*ackmodel.Model, error) {
	latestAPIVersion, err := getLatestAPIVersion(metadata.APIVersions)
	if err != nil {
		return nil, err
	}
	return loadModel(svcAlias, latestAPIVersion, "", cfg)
}

// loadModel finds the AWS SDK for a given service alias and creates a new model
// with the given API version. The cfg parameter should be pre-loaded by the
// caller so that the model name can be resolved before fetching.
func loadModel(svcAlias string, apiVersion string, apiGroup string, cfg ackgenconfig.Config) (*ackmodel.Model, error) {
	totalStart := time.Now()

	modelName := resolveModelName(svcAlias, cfg)

	apiStart := time.Now()
	sdkHelper := acksdk.NewHelper(sdkDir, cfg)
	sdkAPI, err := sdkHelper.API(modelName)
	if err != nil {
		return nil, err
	}
	util.Tracef("SDK API loading (model=%s): %s\n", modelName, time.Since(apiStart))

	if apiGroup != "" {
		sdkAPI.APIGroupSuffix = apiGroup
	}

	docCfg, err := ackgenconfig.NewDocumentationConfig(optDocumentationConfigPath)
	if err != nil {
		return nil, err
	}

	modelStart := time.Now()
	m, err := ackmodel.New(
		sdkAPI, svcAlias, apiVersion, cfg, docCfg,
	)
	if err != nil {
		return nil, err
	}
	util.Tracef("model.New(): %s\n", time.Since(modelStart))
	util.Tracef("loadModel total: %s\n", time.Since(totalStart))
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
