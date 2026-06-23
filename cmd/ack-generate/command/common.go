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
	"os"
	"path/filepath"
	"sort"
	"strings"

	k8sversion "k8s.io/apimachinery/pkg/version"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	ackmetadata "github.com/aws-controllers-k8s/code-generator/pkg/metadata"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	acksdk "github.com/aws-controllers-k8s/code-generator/pkg/sdk"
)

// controllerContext holds the resolved state shared between the controller
// and release commands.
type controllerContext struct {
	SvcAlias       string
	ControllerPath string
	APIVersion     string
	RuntimeVersion string
	RBACRoleName   string
	Model          *ackmodel.Model
	Metadata       *ackmetadata.ServiceMetadata
}

// resolveControllerContext performs the common setup for commands that run
// from a controller repo: resolves config paths, loads the model, and
// determines the API version and runtime version.
func resolveControllerContext(svcAlias string) (*controllerContext, error) {
	controllerPath := optOutputPath
	if controllerPath == "" {
		var err error
		controllerPath, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("determining current directory: %w", err)
		}
	}
	controllerPath, _ = filepath.Abs(controllerPath)
	optOutputPath = controllerPath

	if _, err := os.Stat(filepath.Join(controllerPath, "go.mod")); err != nil {
		return nil, fmt.Errorf("path does not appear to be a controller repo (no go.mod found): %s", controllerPath)
	}

	resolved := ackgenerate.ResolveConfigPaths(controllerPath)
	if optGeneratorConfigPath == "" {
		optGeneratorConfigPath = resolved.GeneratorConfigPath
	}
	if optMetadataConfigPath == "" {
		optMetadataConfigPath = resolved.MetadataConfigPath
	}
	if optDocumentationConfigPath == "" {
		optDocumentationConfigPath = resolved.DocumentationConfigPath
	}
	if optServiceAccountName == "" {
		optServiceAccountName = fmt.Sprintf("ack-%s-controller", svcAlias)
	}

	// Detect template overrides in the controller repo
	svcTemplatesDir := filepath.Join(controllerPath, "templates")
	if fi, err := os.Stat(svcTemplatesDir); err == nil && fi.IsDir() {
		optTemplateDirs = append([]string{svcTemplatesDir}, optTemplateDirs...)
	}

	// Load service metadata early so we can determine the API version
	// and read the last-generation SDK version from ack-generate-metadata.yaml.
	metadata, err := ackmetadata.NewServiceMetadata(optMetadataConfigPath)
	if err != nil {
		return nil, err
	}

	apiVersion := "v1alpha1"
	if len(metadata.APIVersions) > 0 {
		av, err := getLatestAPIVersion(metadata.APIVersions)
		if err == nil {
			apiVersion = av
		}
	}

	cfg, err := setupGenerator(svcAlias)
	if err != nil {
		return nil, err
	}

	m, err := loadModelWithLatestAPIVersion(svcAlias, metadata, cfg)
	if err != nil {
		return nil, err
	}

	runtimeVersion, err := acksdk.GetRuntimeVersion(controllerPath)
	if err != nil {
		return nil, fmt.Errorf("resolving runtime version: %w", err)
	}

	return &controllerContext{
		SvcAlias:       svcAlias,
		ControllerPath: controllerPath,
		APIVersion:     apiVersion,
		RuntimeVersion: runtimeVersion,
		RBACRoleName:   fmt.Sprintf("ack-%s-controller", svcAlias),
		Model:          m,
		Metadata:       metadata,
	}, nil
}

// writeTemplateSet executes a template set and writes the output files to
// the given base directory. If optDryRun is set, prints to stdout instead.
func writeTemplateSet(ts *templateset.TemplateSet, baseDir string) error {
	if err := ts.Execute(); err != nil {
		return err
	}
	for path, contents := range ts.Executed() {
		if optDryRun {
			fmt.Printf("============================= %s ======================================\n", path)
			fmt.Println(strings.TrimSpace(contents.String()))
			continue
		}
		outPath := filepath.Join(baseDir, path)
		if _, err := acksdk.EnsureDir(filepath.Dir(outPath)); err != nil {
			return err
		}
		if err := os.WriteFile(outPath, contents.Bytes(), 0666); err != nil {
			return err
		}
	}
	return nil
}

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
// with the given API version.
func loadModel(svcAlias string, apiVersion string, apiGroup string, cfg ackgenconfig.Config) (*ackmodel.Model, error) {
	modelName := resolveModelName(svcAlias, cfg)

	sdkHelper := acksdk.NewHelper(sdkDir, cfg)
	sdkAPI, err := sdkHelper.API(modelName)
	if err != nil {
		return nil, err
	}

	if apiGroup != "" {
		sdkAPI.APIGroupSuffix = apiGroup
	}

	docCfg, err := ackgenconfig.NewDocumentationConfig(optDocumentationConfigPath)
	if err != nil {
		return nil, err
	}

	return ackmodel.New(sdkAPI, svcAlias, apiVersion, cfg, docCfg)
}

// getLatestAPIVersion looks in the controller metadata file to determine the
// latest Kubernetes API version for CRDs exposed by the generated service
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

// setupGenerator loads the generator configuration, resolves the SDK version
// and fetches the model file.
func setupGenerator(svcAlias string) (ackgenconfig.Config, error) {
	cfg, err := ackgenconfig.New(optGeneratorConfigPath, ackgenerate.DefaultConfig)
	if err != nil {
		return cfg, err
	}

	// Read the last-used SDK version from ack-generate-metadata.yaml
	lastGenVersion := findLastSDKVersion(optOutputPath)

	resolvedVersion, err := acksdk.GetSDKVersion(optAWSSDKGoVersion, lastGenVersion, optOutputPath)
	if err != nil {
		return cfg, err
	}
	resolvedVersion = acksdk.EnsureSemverPrefix(resolvedVersion)
	fmt.Printf("Using AWS SDK version %s\n", resolvedVersion)

	modelName := resolveModelName(svcAlias, cfg)
	ctx, cancel := acksdk.ContextWithSigterm(context.Background())
	defer cancel()
	basePath, err := acksdk.EnsureModel(ctx, optCacheDir, resolvedVersion, modelName)
	if err != nil {
		return cfg, err
	}
	sdkDir = basePath
	sdkVersion = resolvedVersion

	return cfg, nil
}

// findLastSDKVersion scans the apis/ directory for an ack-generate-metadata.yaml
// file and returns the AWS SDK version recorded in it. Returns "" if not found.
func findLastSDKVersion(outputPath string) string {
	if outputPath == "" {
		return ""
	}
	apisPath := filepath.Join(outputPath, "apis")
	entries, err := os.ReadDir(apisPath)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		meta, err := ackmetadata.ReadGenerationMetadata(apisPath, entry.Name())
		if err == nil && meta.AWSSDKGoVersion != "" {
			return meta.AWSSDKGoVersion
		}
	}
	return ""
}
