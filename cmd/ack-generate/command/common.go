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
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	k8sversion "k8s.io/apimachinery/pkg/version"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	ackmetadata "github.com/aws-controllers-k8s/code-generator/pkg/metadata"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/sdk"
	acksdk "github.com/aws-controllers-k8s/code-generator/pkg/sdk"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

// svcSDKVersion holds the resolved per-service SDK version for use by
// saveGeneratedMetadata.
var svcSDKVersion string

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

// setupGenerator loads the generator configuration, resolves the SDK version and fetches the
// model file
func setupGenerator(svcAlias string) (ackgenconfig.Config, error) {
	var cfg ackgenconfig.Config

	// Mutual exclusivity: both explicit CLI flags is an error
	if optAWSSDKGoVersion != "" && optAWSServiceSDKVersion != "" {
		return cfg, fmt.Errorf(
			"--aws-sdk-go-version and --aws-service-sdk-version are mutually exclusive; provide only one",
		)
	}

	// Load generator config to resolve model name before fetching
	cfg, err := ackgenconfig.New(optGeneratorConfigPath, ackgenerate.DefaultConfig)
	if err != nil {
		return cfg, err
	}

	// Load existing generation metadata (used for both per-service and core
	// version resolution fallbacks).
	var metadataSvcSDKVersion string
	var metadataCoreSDKVersion string
	existingMetadata, err := ackmetadata.LoadGenerationMetadata(
		filepath.Join(optOutputPath, "apis"), optGenVersion,
	)
	if err != nil {
		return cfg, fmt.Errorf("cannot load existing generation metadata: %v", err)
	}
	if existingMetadata != nil {
		metadataSvcSDKVersion = existingMetadata.AWSServiceSDKVersion
		metadataCoreSDKVersion = existingMetadata.AWSSDKGoVersion
	}

	// Resolve per-service SDK version from priority chain:
	// CLI flag → metadata YAML → empty
	svcSDKVersion = sdk.GetServiceSDKVersion(optAWSServiceSDKVersion, metadataSvcSDKVersion)

	// Resolve SDK version and fetch the model file
	fetchStart := time.Now()

	// When a per-service SDK version is set, the core version is still needed
	// for the sdkVersion variable (used by metadata saving and other callers),
	// but it is resolved from metadata/go.mod as a fallback — not as the
	// primary fetch source. A resolution failure is non-fatal when the
	// per-service version drives the EnsureModel fetch strategy.
	resolvedVersion, err := sdk.GetSDKVersion(optAWSSDKGoVersion, metadataCoreSDKVersion, optOutputPath)
	if err != nil {
		if svcSDKVersion == "" {
			return cfg, err
		}
		// Per-service version is set; core version is best-effort.
		resolvedVersion = ""
	}
	if resolvedVersion != "" {
		resolvedVersion = sdk.EnsureSemverPrefix(resolvedVersion)
	}

	modelName := resolveModelName(svcAlias, cfg)

	ctx, cancel := sdk.ContextWithSigterm(context.Background())
	defer cancel()
	basePath, err := sdk.EnsureModel(ctx, optCacheDir, resolvedVersion, modelName, svcSDKVersion)
	if err != nil {
		return cfg, err
	}
	sdkDir = basePath
	sdkVersion = resolvedVersion
	util.Tracef("EnsureModel: %s\n", time.Since(fetchStart))

	return cfg, nil
}
