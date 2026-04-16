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

package sdk_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws-controllers-k8s/code-generator/pkg/sdk"
)

// TestGetServiceSDKVersion validates the version resolution priority chain:
// CLI flag if non-empty, else metadata YAML value if non-empty, else empty
// string.
func TestGetServiceSDKVersion(t *testing.T) {
	tests := []struct {
		name                  string
		awsServiceSDKVersion  string
		lastGenerationVersion string
		expected              string
	}{
		{
			name:                  "flag only",
			awsServiceSDKVersion:  "v1.2.0",
			lastGenerationVersion: "",
			expected:              "v1.2.0",
		},
		{
			name:                  "metadata only",
			awsServiceSDKVersion:  "",
			lastGenerationVersion: "v1.0.0",
			expected:              "v1.0.0",
		},
		{
			name:                  "both provided, flag wins",
			awsServiceSDKVersion:  "v2.0.0",
			lastGenerationVersion: "v1.0.0",
			expected:              "v2.0.0",
		},
		{
			name:                  "neither provided, returns empty",
			awsServiceSDKVersion:  "",
			lastGenerationVersion: "",
			expected:              "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sdk.GetServiceSDKVersion(tc.awsServiceSDKVersion, tc.lastGenerationVersion)
			assert.Equal(t, tc.expected, got)
		})
	}
}
