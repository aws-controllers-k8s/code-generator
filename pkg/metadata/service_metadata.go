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
	"errors"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

var (
	ErrNoServiceMetadataFile      = errors.New("expected metadata file path, none provided")
	ErrNoAvailableVersions        = errors.New("service metadata contains no available versions")
	ErrNoDevelopmentVersion       = errors.New("service metadata contains no development versions")
	ErrMultipleDevelopmentVersion = errors.New("service metadata contains multiple development versions")
)

// ServiceMetadata consists of information about the service and relative links as well
// as a list of supported/deprecated versions
type ServiceMetadata struct {
	Service ServiceDetails `json:"service"`
	// A list of all generated API versions of the service
	APIVersions []ServiceVersion `json:"api_versions"`
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
	APIVersion string    `json:"api_version"`
	Status     APIStatus `json:"status"`
}

func (m *ServiceMetadata) GetLatestAPIVersion() (string, error) {
	availableVersions := m.GetAvailableAPIVersions()

	if len(availableVersions) == 0 {
		return "", ErrNoAvailableVersions
	}

	return availableVersions[len(availableVersions)-1], nil
}

func (m *ServiceMetadata) GetDeprecatedAPIVersions() []string {
	return m.getVersionsByStatus(APIStatusDeprecated)
}

func (m *ServiceMetadata) GetRemovedAPIVersions() []string {
	return m.getVersionsByStatus(APIStatusRemoved)
}

func (m *ServiceMetadata) GetAvailableAPIVersions() []string {
	return m.getVersionsByStatus(APIStatusAvailable)
}

func (m *ServiceMetadata) GetDevelopmentAPIVersion() (string, error) {
	devVersions := m.getVersionsByStatus(APIStatusInDevelopment)

	if len(devVersions) == 0 {
		return "", ErrNoDevelopmentVersion
	}

	if len(devVersions) > 1 {
		return "", ErrMultipleDevelopmentVersion
	}

	return devVersions[0], nil
}

func (m *ServiceMetadata) getVersionsByStatus(status APIStatus) []string {
	if len(m.APIVersions) == 0 {
		return []string{}
	}

	versions := []string{}
	for _, v := range m.APIVersions {
		if v.Status == status {
			versions = append(versions, v.APIVersion)
		}
	}
	return versions
}

// NewServiceMetadata returns a new Metadata object given a supplied
// path to a metadata file
func NewServiceMetadata(
	metadataPath string,
) (ServiceMetadata, error) {
	if metadataPath == "" {
		return ServiceMetadata{}, ErrNoServiceMetadataFile
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
