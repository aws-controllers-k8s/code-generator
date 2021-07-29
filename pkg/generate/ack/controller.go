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
	"path/filepath"
	"strings"
	ttpl "text/template"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
)

var (
	controllerConfigTemplatePaths = []string{
		"config/controller/deployment.yaml.tpl",
		"config/controller/service.yaml.tpl",
		"config/controller/kustomization.yaml.tpl",
		"config/default/kustomization.yaml.tpl",
		"config/rbac/cluster-role-binding.yaml.tpl",
		"config/rbac/role-reader.yaml.tpl",
		"config/rbac/role-writer.yaml.tpl",
		"config/rbac/kustomization.yaml.tpl",
		"config/crd/kustomization.yaml.tpl",
	}
	controllerIncludePaths = []string{
		"config/controller/kustomization_def.yaml.tpl",
		"boilerplate.go.tpl",
		"pkg/resource/sdk_find_read_one.go.tpl",
		"pkg/resource/sdk_find_get_attributes.go.tpl",
		"pkg/resource/sdk_find_read_many.go.tpl",
		"pkg/resource/sdk_find_not_implemented.go.tpl",
		"pkg/resource/sdk_update.go.tpl",
		"pkg/resource/sdk_update_custom.go.tpl",
		"pkg/resource/sdk_update_set_attributes.go.tpl",
		"pkg/resource/sdk_update_not_implemented.go.tpl",
	}
	controllerCopyPaths = []string{}
	controllerFuncMap   = ttpl.FuncMap{
		"ToLower": strings.ToLower,
		"ResourceExceptionCode": func(r *ackmodel.CRD, httpStatusCode int) string {
			return r.ExceptionCode(httpStatusCode)
		},
		"GoCodeSetExceptionMessageCheck": func(r *ackmodel.CRD, httpStatusCode int) string {
			return code.CheckExceptionMessage(r.Config(), r, httpStatusCode)
		},
		"GoCodeSetReadOneOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetResource(r.Config(), r, ackmodel.OpTypeGet, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetReadOneInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeGet, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetReadManyOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetResource(r.Config(), r, ackmodel.OpTypeList, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetReadManyInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeList, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeGetAttributesSetInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetSDKGetAttributes(r.Config(), r, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetAttributesSetInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetSDKSetAttributes(r.Config(), r, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeGetAttributesSetOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetResourceGetAttributes(r.Config(), r, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetCreateOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetResource(r.Config(), r, ackmodel.OpTypeCreate, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetCreateInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeCreate, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetUpdateOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetResource(r.Config(), r, ackmodel.OpTypeUpdate, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetUpdateInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeUpdate, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetDeleteInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeDelete, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetOperationStruct": func(r *ackmodel.CRD, targetFieldName string, targetVarName string, targetShapeRef *awssdkmodel.ShapeRef, sourceFieldPath string, sourceVarName string, indentLevel int) string {
			return code.SetSDKForStruct(r.Config(), r, targetFieldName, targetVarName, targetShapeRef, sourceFieldPath, sourceVarName, indentLevel)
		},
		"GoCodeSetResourceStruct": func(r *ackmodel.CRD, targetFieldName string, targetVarName string, targetShapeRef *awssdkmodel.ShapeRef, sourceVarName string, sourceShapeRef *awssdkmodel.ShapeRef, indentLevel int) string {
			return code.SetResourceForStruct(r.Config(), r, targetFieldName, targetVarName, targetShapeRef, sourceVarName, sourceShapeRef, indentLevel)
		},
		"GoCodeCompare": func(r *ackmodel.CRD, deltaVarName string, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.CompareResource(r.Config(), r, deltaVarName, sourceVarName, targetVarName, indentLevel)
		},
		"Empty": func(subject string) bool {
			return strings.TrimSpace(subject) == ""
		},
		"GoCodeRequiredFieldsMissingFromReadOneInput": func(r *ackmodel.CRD, koVarName string, indentLevel int) string {
			return code.CheckRequiredFieldsMissingFromShape(r, ackmodel.OpTypeGet, koVarName, indentLevel)
		},
		"GoCodeRequiredFieldsMissingFromGetAttributesInput": func(r *ackmodel.CRD, koVarName string, indentLevel int) string {
			return code.CheckRequiredFieldsMissingFromShape(r, ackmodel.OpTypeGetAttributes, koVarName, indentLevel)
		},
		"GoCodeRequiredFieldsMissingFromSetAttributesInput": func(r *ackmodel.CRD, koVarName string, indentLevel int) string {
			return code.CheckRequiredFieldsMissingFromShape(r, ackmodel.OpTypeSetAttributes, koVarName, indentLevel)
		},
		"GoCodeSetResourceIdentifiers": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.SetResourceIdentifiers(r.Config(), r, sourceVarName, targetVarName, indentLevel)
		},
	}
)

// Controller returns a pointer to a TemplateSet containing all the templates
// for generating ACK service controller implementations
func Controller(
	m *ackmodel.Model,
	templateBasePaths []string,
) (*templateset.TemplateSet, error) {
	crds, err := m.GetCRDs()
	if err != nil {
		return nil, err
	}

	metaVars := m.MetaVars()

	// Hook code can reference a template path, and we can look up the template
	// in any of our base paths...
	controllerFuncMap["Hook"] = func(r *ackmodel.CRD, hookID string) string {
		crdVars := &templateCRDVars{
			metaVars,
			m.SDKAPI,
			r,
		}
		code, err := ResourceHookCode(templateBasePaths, r, hookID, crdVars, controllerFuncMap)
		if err != nil {
			// It's a compile-time error, so just panic...
			panic(err)
		}
		return code
	}

	ts := templateset.New(
		templateBasePaths,
		controllerIncludePaths,
		controllerCopyPaths,
		controllerFuncMap,
	)

	// First add all the CRD pkg/resource templates
	targets := []string{
		"delta.go.tpl",
		"descriptor.go.tpl",
		"identifiers.go.tpl",
		"manager.go.tpl",
		"manager_factory.go.tpl",
		"resource.go.tpl",
		"sdk.go.tpl",
	}
	for _, crd := range crds {
		for _, target := range targets {
			outPath := filepath.Join("pkg/resource", crd.Names.Snake, strings.TrimSuffix(target, ".tpl"))
			tplPath := filepath.Join("pkg/resource", target)
			crdVars := &templateCRDVars{
				metaVars,
				m.SDKAPI,
				crd,
			}
			if err = ts.Add(outPath, tplPath, crdVars); err != nil {
				return nil, err
			}
		}
	}

	configVars := &templateConfigVars{
		metaVars,
		m.GetConfig(),
	}
	if err = ts.Add("pkg/resource/registry.go", "pkg/resource/registry.go.tpl", configVars); err != nil {
		return nil, err
	}

	// Next add the template for pkg/version/version.go file
	if err = ts.Add("pkg/version/version.go", "pkg/version/version.go.tpl", nil); err != nil {
		return nil, err
	}

	// Next add the template for the main.go file
	snakeCasedCRDNames := make([]string, 0)
	for _, crd := range crds {
		snakeCasedCRDNames = append(snakeCasedCRDNames, crd.Names.Snake)
	}
	cmdVars := &templateCmdVars{
		metaVars,
		snakeCasedCRDNames,
	}
	if err = ts.Add("cmd/controller/main.go", "cmd/controller/main.go.tpl", cmdVars); err != nil {
		return nil, err
	}

	// Finally, add the configuration YAML file templates
	for _, path := range controllerConfigTemplatePaths {
		outPath := strings.TrimSuffix(path, ".tpl")
		if err = ts.Add(outPath, path, metaVars); err != nil {
			return nil, err
		}
	}
	return ts, nil
}

// templateCmdVars contains template variables for the template that outputs Go
// code for a single top-level resource's API definition
type templateCmdVars struct {
	templateset.MetaVars
	SnakeCasedCRDNames []string
}

// templateConfigVars contains template variables for the templates that require
// access to the generator configuration definition
type templateConfigVars struct {
	templateset.MetaVars
	GeneratorConfig *ackgenconfig.Config
}
