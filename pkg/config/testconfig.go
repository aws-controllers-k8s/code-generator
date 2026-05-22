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

package config

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

// TestConfig represents the testconfig.yaml file that provides domain-specific
// test values for e2e test generation.
type TestConfig struct {
	// Service is the ACK service alias (e.g., "rds", "s3").
	Service string `json:"service"`
	// Resources maps CRD Kind names to a map of named test configurations.
	// The outer key is the CRD Kind (e.g., "Bucket"), the inner key is the
	// test name (e.g., "versioning_lifecycle") which becomes the test function
	// suffix.
	Resources map[string]map[string]TestResourceConfig `json:"resources"`
	// Bootstrap lists shared AWS resources to create before tests and destroy after.
	Bootstrap []BootstrapResourceConfig `json:"bootstrap,omitempty"`
}

// TestResourceConfig provides test values and configuration for a single CRD resource.
type TestResourceConfig struct {
	// CreateValues maps Go struct field names (PascalCase) to values used when
	// creating the resource in tests.
	CreateValues map[string]interface{} `json:"create_values,omitempty"`
	// UpdateValues maps Go struct field names to new values used in update tests.
	UpdateValues map[string]interface{} `json:"update_values,omitempty"`
	// CreateWait is the timeout in seconds to wait for the resource to become synced
	// after creation. Defaults to 60.
	CreateWait int `json:"create_wait,omitempty"`
	// DeleteWait is the timeout in seconds to wait for deletion. Defaults to 60.
	DeleteWait int `json:"delete_wait,omitempty"`
	// Skip when true, no test is generated for this resource.
	Skip bool `json:"skip,omitempty"`
	// BootstrapFields maps Go struct field names to bootstrap resource field paths
	// (e.g., "SharedVPC.SecurityGroupID"). At runtime, the value is retrieved from
	// the named bootstrap resource.
	BootstrapFields map[string]string `json:"bootstrap_fields,omitempty"`
}

// BootstrapResourceConfig describes a shared AWS resource to be bootstrapped.
type BootstrapResourceConfig struct {
	// Type is the bootstrap resource type (e.g., "VPC", "IAMRole", "S3Bucket").
	Type string `json:"type"`
	// Name is the logical name used to reference this resource in bootstrap_fields.
	Name string `json:"name"`
	// DependsOn lists logical names of other bootstrap resources that must be
	// created before this one.
	DependsOn []string `json:"depends_on,omitempty"`
}

// GetCreateWait returns the create wait timeout in seconds, defaulting to 60.
func (r TestResourceConfig) GetCreateWait() int {
	if r.CreateWait > 0 {
		return r.CreateWait
	}
	return 60
}

// GetDeleteWait returns the delete wait timeout in seconds, defaulting to 60.
func (r TestResourceConfig) GetDeleteWait() int {
	if r.DeleteWait > 0 {
		return r.DeleteWait
	}
	return 60
}

// NewTestConfig reads and parses a testconfig.yaml file at the given path.
// Returns nil (not an error) if the file does not exist.
func NewTestConfig(path string) (*TestConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading test config %s: %w", path, err)
	}
	tc := &TestConfig{}
	if err := yaml.Unmarshal(data, tc); err != nil {
		return nil, fmt.Errorf("parsing test config %s: %w", path, err)
	}
	return tc, nil
}
