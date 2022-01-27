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

package ack

import (
	"strings"
	ttpl "text/template"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	ackmetadata "github.com/aws-controllers-k8s/code-generator/pkg/metadata"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
)

var (
	releaseTemplatePaths = []string{
		"config/controller/kustomization.yaml.tpl",
		"helm/templates/cluster-role-binding.yaml.tpl",
		"helm/Chart.yaml.tpl",
		"helm/values.yaml.tpl",
		"helm/values.schema.json",
		"helm/templates/NOTES.txt.tpl",
		"helm/templates/role-reader.yaml.tpl",
		"helm/templates/role-writer.yaml.tpl",
		"helm/templates/_controller-role-kind-patch.yaml.tpl",
	}
	releaseIncludePaths = []string{
		"config/controller/kustomization_def.yaml.tpl",
	}
	releaseCopyPaths = []string{
		"helm/templates/_helpers.tpl",
		"helm/templates/deployment.yaml",
		"helm/templates/metrics-service.yaml",
		"helm/templates/service-account.yaml",
	}
	releaseFuncMap = ttpl.FuncMap{
		"ToLower": strings.ToLower,
		"Empty": func(subject string) bool {
			return strings.TrimSpace(subject) == ""
		},
	}
)

// Release returns a pointer to a TemplateSet containing all the templates for
// generating an ACK service controller release (Helm artifacts, etc)
func Release(
	m *ackmodel.Model,
	metadata *ackmetadata.ServiceMetadata,
	templateBasePaths []string,
	// releaseVersion is the SemVer string describing the release that the Helm
	// chart will install
	releaseVersion string,
	// imageRepository is the Docker image repository to use when generating
	// release files
	imageRepository string,
	// serviceAccountName is the name of the ServiceAccount used in the Helm chart
	serviceAccountName string,
) (*templateset.TemplateSet, error) {
	ts := templateset.New(
		templateBasePaths,
		releaseIncludePaths,
		releaseCopyPaths,
		releaseFuncMap,
	)

	metaVars := m.MetaVars()
	releaseVars := &templateReleaseVars{
		metaVars,
		metadata,
		releaseVersion,
		imageRepository,
		serviceAccountName,
	}
	for _, path := range releaseTemplatePaths {
		outPath := strings.TrimSuffix(path, ".tpl")
		if err := ts.Add(outPath, path, releaseVars); err != nil {
			return nil, err
		}
	}

	return ts, nil
}

// templateReleaseVars contains template variables for the template that
// outputs Go code for a release artifact
type templateReleaseVars struct {
	templateset.MetaVars
	Metadata *ackmetadata.ServiceMetadata
	// ReleaseVersion is the semver release tag (or Git SHA1 commit) that is
	// used for the binary image artifacts and Helm release version
	ReleaseVersion string
	// ImageRepository is the Docker image repository to inject into the Helm
	// values template
	ImageRepository string
	// ServiceAccountName is the name of the ServiceAccount used in the Helm chart
	ServiceAccountName string
}
