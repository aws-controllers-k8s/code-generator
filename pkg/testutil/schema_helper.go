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

	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
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

// NewModelForService returns a new *ackmodel.Model used for testing purposes.
func NewModelForService(t *testing.T, servicePackageName string) *ackmodel.Model {
	return NewModelForServiceWithOptions(t, servicePackageName, &TestingModelOptions{})
}

// NewModelForServiceWithOptions returns a new *ackmodel.Model used for testing purposes.
func NewModelForServiceWithOptions(t *testing.T, servicePackageName string, options *TestingModelOptions) *ackmodel.Model {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	// We have subdirectories in pkg/generate and pkg/model that rely on the testdata
	// in pkg/generate. This code simply detects if we're running from one of
	// those subdirectories and if so, rebuilds the path to the API model files
	// in pkg/generate/testdata
	pathParts := strings.Split(path, "/")
	for x, pathPart := range pathParts {
		if pathPart == "generate" || pathPart == "model" {
			path = filepath.Join(pathParts[0:x]...)
			path = filepath.Join("/", path, "testdata")
			break
		}
	}
	options.SetDefaults()

	generatorConfigPath := filepath.Join(path, "models", "apis", servicePackageName, options.ServiceAPIVersion, options.GeneratorConfigFile)
	if _, err := os.Stat(generatorConfigPath); os.IsNotExist(err) {
		generatorConfigPath = ""
	}
	cfg, err := ackgenconfig.New(generatorConfigPath, ackgenerate.DefaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	sdkHelper := acksdk.NewHelper(path, cfg)
	sdkHelper.WithAPIVersion(options.ServiceAPIVersion)
	sdkAPI, err := sdkHelper.API(servicePackageName)
	if err != nil {
		t.Fatal(err)
	}
	m, err := ackmodel.New(sdkAPI, servicePackageName, options.APIVersion, cfg)
	if err != nil {
		t.Fatal(err)
	}
	return m
}
