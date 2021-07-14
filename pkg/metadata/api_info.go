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

type APIStatus string

const (
	APIStatusAvailable  APIStatus = "available"
	APIStatusRemoved              = "removed"
	APIStatusDeprecated           = "deprecated"
)

// APIInfo contains information related a specific apiVersion.
type APIInfo struct {
	// The API status. Can be one of Available, Removed and Deprecated.
	Status APIStatus
	// the aws-sdk-go version used to generated the apiVersion.
	AWSSDKVersion string
	// Full path of the generator config file.
	GeneratorConfigPath string
	// The API version.
	APIVersion string
}
