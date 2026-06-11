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

package testutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	acksdk "github.com/aws-controllers-k8s/code-generator/pkg/sdk"
)

// TestingModelOptions contains optional variables that are passed to
// `NewModelForServiceWithOptions`.
type TestingModelOptions struct {
	// The CR API Version. Defaults to v1alpha1
	APIVersion string
	// The generator config file. Defaults to generator.yaml
	GeneratorConfigFile string
	// The documentation config file. No default value
	DocumentationConfigFile string
	// The AWS Service's API version. Defaults to 00-00-0000
	ServiceAPIVersion string
}

// SetDefaults sets the empty fields to a default value.
func (o *TestingModelOptions) SetDefaults() {
	if o.APIVersion == "" {
		o.APIVersion = "v1alpha1"
	}
	if o.GeneratorConfigFile == "" {
		o.GeneratorConfigFile = "generator.yaml"
	}
	if o.ServiceAPIVersion == "" {
		o.ServiceAPIVersion = "0000-00-00"
	}
}

// repoRootFromWD returns the repository root directory by looking for "generate"
// or "model" in the working directory path (both live directly under pkg/).
// Returns ("", false) if the root cannot be determined.
func repoRootFromWD(wd string) (string, bool) {
	pathParts := strings.Split(wd, "/")
	for x, pathPart := range pathParts {
		if pathPart == "generate" || pathPart == "model" {
			// x-1 points to "pkg"; everything before that is the repo root
			return filepath.Join("/", filepath.Join(pathParts[0:x-1]...)), true
		}
	}
	return "", false
}

// TemplatesBasePath returns the absolute path to the templates/ directory at
// the root of the repository.
func TemplatesBasePath(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	root, ok := repoRootFromWD(wd)
	if !ok {
		t.Fatalf("could not determine templates/ path from working directory %q", wd)
	}
	return filepath.Join(root, "templates")
}

// NewModelForService returns a new *ackmodel.Model used for testing purposes.
func NewModelForService(t *testing.T, servicePackageName string) *ackmodel.Model {
	return NewModelForServiceWithOptions(t, servicePackageName, &TestingModelOptions{})
}

// NewModelForServiceWithOptions returns a new *ackmodel.Model used for testing purposes.
func NewModelForServiceWithOptions(t *testing.T, servicePackageName string, options *TestingModelOptions) *ackmodel.Model {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// We have subdirectories in pkg/generate and pkg/model that rely on the testdata
	// in pkg/generate. This code simply detects if we're running from one of
	// those subdirectories and if so, rebuilds the path to the API model files
	// in pkg/testdata
	testdataPath := wd
	if root, ok := repoRootFromWD(wd); ok {
		testdataPath = filepath.Join(root, "pkg", "testdata")
	}
	options.SetDefaults()

	generatorConfigPath := filepath.Join(testdataPath, "models", "apis", servicePackageName, options.ServiceAPIVersion, options.GeneratorConfigFile)
	if _, err := os.Stat(generatorConfigPath); os.IsNotExist(err) {
		t.Fatalf("Could not find generator file %q", generatorConfigPath)
	}
	cfg, err := ackgenconfig.New(generatorConfigPath, ackgenerate.DefaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	sdkHelper := acksdk.NewHelper(testdataPath, cfg)
	sdkHelper.WithAPIVersion(options.ServiceAPIVersion)
	sdkAPI, err := sdkHelper.API(servicePackageName)
	if err != nil {
		t.Fatal(err)
	}

	docConfigPath := ""
	if options.DocumentationConfigFile != "" {
		docConfigPath = filepath.Join(testdataPath, "models", "apis", servicePackageName, options.ServiceAPIVersion, options.DocumentationConfigFile)
	}
	docCfg, err := ackgenconfig.NewDocumentationConfig(docConfigPath)
	if err != nil {
		t.Fatal(err)
	}

	m, err := ackmodel.New(sdkAPI, servicePackageName, options.APIVersion, cfg, docCfg)
	if err != nil {
		t.Fatal(err)
	}
	return m
}
