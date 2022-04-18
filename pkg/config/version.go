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

// APIVersion represents an API version of the generated CRD.
type APIVersion struct {
	// Name of the API version, e.g. v1beta1
	Name string `json:"name"`
	// Served whether this version is enabled or not
	Served *bool `json:"served,omitempty"`
	// Storage whether this version is the storage version.
	// One and only one version can be set as the storage version.
	Storage *bool `json:"storage,omitempty"`
}

// GetAPIVersions returns the API version(s) for a CRD
func (c *Config) GetAPIVersions(crdName string) []APIVersion {
	res := []APIVersion{}
	if c == nil {
		return res
	}
	resConfig, found := c.Resources[crdName]
	if !found {
		return res
	}
	return resConfig.APIVersions
}
