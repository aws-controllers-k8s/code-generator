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
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// ServiceMetadata consists of information about the service and relative links as well
// as a list of supported/deprecated versions
type ServiceMetadata struct {
	Service ServiceDetails `json:"service"`
	// CRDs to ignore. ACK generator would skip these resources.
	Versions []ServiceVersion `json:"versions"`
}

// ServiceDetails contains string identifiers and relevant links for the
// service
type ServiceDetails struct {
	// The full display name for the service. eg. Amazon Elastic Kubernetes
	// Service
	FullName string `json:"full_name"`
	// The short name (abbreviation) for the service. eg. S3
	ShortName string `json:"short_name"`
	// The URL of the service's homepage
	Link string `json:"link"`
	// The URL of the service's main documentation/user guide
	Documentation string `json:"documentation"`
}

// ServiceVersion describes the status of all existing version of the controller
type ServiceVersion struct {
	APIVersion string `json:"api_version"`
	//TODO(RedbackThomson): Update with APIStatus after merging #121
	Status string `json:"status"`
}

// New returns a new Metadata object given a supplied
// path to a metadata file
func New(
	metadataPath string,
) (ServiceMetadata, error) {
	if metadataPath == "" {
		return ServiceMetadata{}, fmt.Errorf("expected metadata file path, none provided")
	}
	content, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return ServiceMetadata{}, err
	}
	gc := ServiceMetadata{}
	if err = yaml.Unmarshal(content, &gc); err != nil {
		return ServiceMetadata{}, err
	}
	return gc, nil
}
