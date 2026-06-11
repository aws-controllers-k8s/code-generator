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

package metadata

import (
	"strings"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerationMetadata_RoundTrip validates that serializing GenerationMetadata
// to YAML and deserializing back produces identical field values, including the
// aws_service_sdk_version field when present and its absence when empty.
func TestGenerationMetadata_RoundTrip(t *testing.T) {
	tests := []struct {
		name                 string
		awsServiceSDKVersion string
		expectFieldInYAML    bool
	}{
		{
			name:                 "with per-service SDK version",
			awsServiceSDKVersion: "v1.0.0",
			expectFieldInYAML:    true,
		},
		{
			name:                 "without per-service SDK version",
			awsServiceSDKVersion: "",
			expectFieldInYAML:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			original := GenerationMetadata{
				APIVersion:           "v1alpha1",
				APIDirectoryChecksum: "abc123",
				LastModification: lastModificationInfo{
					Reason: UpdateReasonAPIGeneration,
				},
				AWSSDKGoVersion:      "v1.41.5",
				AWSServiceSDKVersion: tc.awsServiceSDKVersion,
				ACKGenerateInfo: ackGenerateInfo{
					Version:   "v0.58.0",
					GoVersion: "go1.25.3",
					BuildDate: "2026-04-14T18:02:39Z",
					BuildHash: "a9e2cea",
				},
				GeneratorConfigInfo: generatorConfigInfo{
					OriginalFileName: "generator.yaml",
					FileChecksum:     "5b522171",
				},
			}

			data, err := yaml.Marshal(original)
			require.NoError(t, err)

			// Verify the aws_service_sdk_version field presence/absence in raw YAML
			yamlStr := string(data)
			if tc.expectFieldInYAML {
				assert.True(t, strings.Contains(yamlStr, "aws_service_sdk_version"),
					"expected aws_service_sdk_version in YAML output")
			} else {
				assert.False(t, strings.Contains(yamlStr, "aws_service_sdk_version"),
					"expected aws_service_sdk_version to be omitted from YAML output")
			}

			// Unmarshal back and verify round-trip equality
			var restored GenerationMetadata
			err = yaml.Unmarshal(data, &restored)
			require.NoError(t, err)

			assert.Equal(t, original.APIVersion, restored.APIVersion)
			assert.Equal(t, original.APIDirectoryChecksum, restored.APIDirectoryChecksum)
			assert.Equal(t, original.LastModification, restored.LastModification)
			assert.Equal(t, original.AWSSDKGoVersion, restored.AWSSDKGoVersion)
			assert.Equal(t, original.AWSServiceSDKVersion, restored.AWSServiceSDKVersion)
			assert.Equal(t, original.ACKGenerateInfo, restored.ACKGenerateInfo)
			assert.Equal(t, original.GeneratorConfigInfo, restored.GeneratorConfigInfo)
		})
	}
}

// TestGenerationMetadata_BackwardCompatibility verifies that YAML without the
// aws_service_sdk_version field deserializes without error and the field
// defaults to empty string (Requirement 4.3).
func TestGenerationMetadata_BackwardCompatibility(t *testing.T) {
	// YAML that predates the aws_service_sdk_version field
	oldYAML := `
api_version: v1alpha1
api_directory_checksum: abc123
aws_sdk_go_version: v1.41.5
last_modification:
  reason: API generation
ack_generate_info:
  version: v0.58.0
  go_version: go1.25.3
  build_date: "2026-04-14T18:02:39Z"
  build_hash: a9e2cea
generator_config_info:
  original_file_name: generator.yaml
  file_checksum: "5b522171"
`
	var gm GenerationMetadata
	err := yaml.Unmarshal([]byte(oldYAML), &gm)
	require.NoError(t, err)

	assert.Equal(t, "v1alpha1", gm.APIVersion)
	assert.Equal(t, "v1.41.5", gm.AWSSDKGoVersion)
	assert.Equal(t, "", gm.AWSServiceSDKVersion, "missing field should default to empty string")
}
