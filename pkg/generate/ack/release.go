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
	"fmt"
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
		"helm/templates/cluster-role-controller.yaml.tpl",
		"helm/templates/_helpers.tpl.tpl",
		"helm/Chart.yaml.tpl",
		"helm/values.yaml.tpl",
		"helm/values.schema.json",
		"helm/templates/NOTES.txt.tpl",
		"helm/templates/role-reader.yaml.tpl",
		"helm/templates/role-writer.yaml.tpl",
		"helm/templates/caches-role.yaml.tpl",
		"helm/templates/caches-role-binding.yaml.tpl",
		"helm/templates/leader-election-role.yaml.tpl",
		"helm/templates/leader-election-role-binding.yaml.tpl",
		"helm/templates/deployment.yaml.tpl",
		"helm/templates/metrics-service.yaml.tpl",
		"helm/templates/service-account.yaml.tpl",
	}
	releaseIncludePaths = []string{
		"config/controller/kustomization_def.yaml.tpl",
	}
	releaseCopyPaths = []string{}
	releaseFuncMap   = func(serviceName string) ttpl.FuncMap {
		return ttpl.FuncMap{
			"ToLower": strings.ToLower,
			"Empty": func(subject string) bool {
				return strings.TrimSpace(subject) == ""
			},
			"DefineTemplate": func(templateName string) string {
				// Returnes a statement that defines a new template name with unique
				// prefix for the ACK controller.
				// For example, if serviceName is "s3" and templateName is "app.name"
				// it will return {{- define "ack-s3-controller.app.name" -}}
				return fmt.Sprintf("{{- define \"%s\" -}}", prefixServiceTemplateName(serviceName, templateName))
			},
			"IncludeTemplate": func(templateName string) string {
				// Returns a statement that includes a template defined with DefineTemplate.
				// For example, if serviceName is "s3" and templateName is "app.name"
				// it will return {{- include "ack-s3-controller.app.name" . -}}
				return fmt.Sprintf("{{ include \"%s\" . }}", prefixServiceTemplateName(serviceName, templateName))
			},
			"VarIncludeTemplate": func(variableName, templateName string) string {
				// Returns a statement that declares a variable and includes a template defined with
				// DefineTemplate.
				// For example, if variableName is appName, serviceName is "s3", and templateName is "app.name"
				// it will return {{- $variable := include "ack-s3-controller.app.name" .app.name -}}
				return fmt.Sprintf("{{ $%s := include \"%s\" . }}", variableName, prefixServiceTemplateName(serviceName, templateName))
			},
		}
	}
)

func prefixServiceTemplateName(serviceName, templateName string) string {
	return fmt.Sprintf("ack-%s-controller.%s", serviceName, templateName)
}

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
		releaseFuncMap(m.MetaVars().ServicePackageName),
	)
	metaVars := m.MetaVars()

	releaseVars := &templateReleaseVars{
		metaVars,
		ImageReleaseVars{
			ReleaseVersion:  strings.TrimPrefix(releaseVersion, "v"),
			ImageRepository: imageRepository,
		},
		metadata,
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

type ImageReleaseVars struct {
	// ReleaseVersion is the semver release tag (or Git SHA1 commit) that is
	// used for the binary image artifacts and Helm release version
	ReleaseVersion string
	// ImageRepository is the Docker image repository to inject into the Helm
	// values template
	ImageRepository string
}

// templateReleaseVars contains template variables for the template that
// outputs Go code for a release artifact
type templateReleaseVars struct {
	templateset.MetaVars
	ImageReleaseVars
	Metadata *ackmetadata.ServiceMetadata
	// ServiceAccountName is the name of the ServiceAccount used in the Helm chart
	ServiceAccountName string
}
