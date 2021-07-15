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
	ttpl "text/template"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/model/multiversion"
)

var (
	webhooksIncludePaths = []string{
		"boilerplate.go.tpl",
		"apis/webhooks/conversion.go.tpl",
	}
	webhookCopyPaths = []string{}
	webhooksFuncMap  = ttpl.FuncMap{
		"GoCodeConvert": func(
			src *ackmodel.CRD,
			dst *ackmodel.CRD,
			convertingToHub bool,
			hubImportPath string,
			sourceVarName string,
			targetVarName string,
			indentLevel int,
		) string {
			return code.Convert(src, dst, convertingToHub, hubImportPath, sourceVarName, targetVarName, indentLevel)
		},
	}
)

// ConversionWebhooks returns a pointer to a TemplateSet containing all the templates
// for generating conversion webhooks.
func ConversionWebhooks(
	mgr *multiversion.APIVersionManager,
	templateBasePaths []string,
) (*templateset.TemplateSet, error) {
	ts := templateset.New(
		templateBasePaths,
		webhooksIncludePaths,
		webhookCopyPaths,
		webhooksFuncMap,
	)
	hubVersion := mgr.GetHubVersion()
	hubModel, err := mgr.GetModel(hubVersion)
	if err != nil {
		return nil, err
	}

	hubMetaVars := hubModel.MetaVars()
	hubCRDs, err := hubModel.GetCRDs()
	if err != nil {
		return nil, err
	}

	for _, crd := range hubCRDs {
		convertVars := conversionVars{
			MetaVars:  hubMetaVars,
			SourceCRD: crd,
			IsHub:     true,
		}
		// Add the hub version template
		target := fmt.Sprintf("apis/%s/%s_conversion.go", hubVersion, crd.Names.Snake)
		if err = ts.Add(target, "apis/webhooks/conversion.go.tpl", convertVars); err != nil {
			return nil, err
		}
	}

	// Add spoke version templates
	for _, spokeVersion := range mgr.GetSpokeVersions() {
		model, err := mgr.GetModel(spokeVersion)
		if err != nil {
			return nil, err
		}

		metaVars := model.MetaVars()
		crds, err := model.GetCRDs()
		if err != nil {
			return nil, err
		}

		for i, crd := range crds {
			convertVars := conversionVars{
				MetaVars:   metaVars,
				SourceCRD:  crd,
				DestCRD:    hubCRDs[i],
				IsHub:      false,
				HubVersion: hubVersion,
			}

			target := fmt.Sprintf("apis/%s/%s_conversion.go", spokeVersion, crd.Names.Snake)
			if err = ts.Add(target, "apis/webhooks/conversion.go.tpl", convertVars); err != nil {
				return nil, err
			}
		}
	}

	return ts, nil
}

// conversionVars contains template variables for templates that output
// Go conversion functions.
type conversionVars struct {
	templateset.MetaVars
	SourceCRD  *ackmodel.CRD
	DestCRD    *ackmodel.CRD
	HubVersion string
	IsHub      bool
}
