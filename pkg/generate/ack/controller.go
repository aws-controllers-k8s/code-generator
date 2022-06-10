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
	"sort"
	"strings"
	ttpl "text/template"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/fieldpath"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
)

var (
	controllerConfigTemplatePaths = []string{
		"config/controller/deployment.yaml.tpl",
		"config/controller/service.yaml.tpl",
		"config/default/kustomization.yaml.tpl",
		"config/rbac/cluster-role-binding.yaml.tpl",
		"config/rbac/role-reader.yaml.tpl",
		"config/rbac/role-writer.yaml.tpl",
		"config/rbac/service-account.yaml.tpl",
		"config/rbac/kustomization.yaml.tpl",
		"config/crd/kustomization.yaml.tpl",
		"config/overlays/namespaced/kustomization.yaml.tpl",
	}
	controllerIncludePaths = []string{
		"boilerplate.go.tpl",
		"pkg/resource/references_read_referenced_resource.go.tpl",
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
		"Dereference": func(s *string) string {
			return *s
		},
		"ResourceExceptionCode": func(r *ackmodel.CRD, httpStatusCode int) string {
			return r.ExceptionCode(httpStatusCode)
		},
		"ConstructFieldPath": func(targetFieldPath string) *fieldpath.Path {
			return fieldpath.FromString(targetFieldPath)
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
		"GoCodeSetSDKForStruct": func(r *ackmodel.CRD, targetFieldName string, targetVarName string, targetShapeRef *awssdkmodel.ShapeRef, sourceFieldPath string, sourceVarName string, indentLevel int) string {
			return code.SetSDKForStruct(r.Config(), r, targetFieldName, targetVarName, targetShapeRef, sourceFieldPath, sourceVarName, indentLevel)
		},
		"GoCodeSetResourceForStruct": func(r *ackmodel.CRD, targetFieldName string, targetVarName string, targetShapeRef *awssdkmodel.ShapeRef, sourceVarName string, sourceShapeRef *awssdkmodel.ShapeRef, indentLevel int) string {
			f, ok := r.Fields[targetFieldName]
			if !ok {
				return ""
			}
			// We may have some special instructions for how to handle setting the
			// field value...
			setCfg := f.GetSetterConfig(model.OpTypeList)

			if setCfg != nil && setCfg.Ignore {
				return ""
			}
			return code.SetResourceForStruct(r.Config(), r, targetFieldName, targetVarName, targetShapeRef, setCfg, sourceVarName, sourceShapeRef, "", model.OpTypeList, indentLevel)
		},
		"GoCodeCompare": func(r *ackmodel.CRD, deltaVarName string, sourceVarName string, targetVarName string, indentLevel int) string {
			return code.CompareResource(r.Config(), r, deltaVarName, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeIsSynced": func(r *ackmodel.CRD, resVarName string, indentLevel int) string {
			return code.ResourceIsSynced(r.Config(), r, resVarName, indentLevel)
		},
		"GoCodeCompareStruct": func(r *ackmodel.CRD, shape *awssdkmodel.Shape, deltaVarName string, sourceVarName string, targetVarName string, fieldPath string, indentLevel int) string {
			return code.CompareStruct(r.Config(), r, nil, shape, deltaVarName, sourceVarName, targetVarName, fieldPath, indentLevel)
		},
		"Empty": func(subject string) bool {
			return strings.TrimSpace(subject) == ""
		},
		"GoCodeRequiredFieldsMissingFromReadOneInput": func(r *ackmodel.CRD, koVarName string, indentLevel int) string {
			return code.CheckRequiredFieldsMissingFromShape(r, ackmodel.OpTypeGet, koVarName, indentLevel)
		},
		"GoCodeRequiredFieldsMissingFromReadManyInput": func(r *ackmodel.CRD, koVarName string, indentLevel int) string {
			return code.CheckRequiredFieldsMissingFromShape(r, ackmodel.OpTypeList, koVarName, indentLevel)
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
		"GoCodeFindLateInitializedFieldNames": func(r *ackmodel.CRD, resVarName string, indentLevel int) string {
			return code.FindLateInitializedFieldNames(r.Config(), r, resVarName, indentLevel)
		},
		"GoCodeLateInitializeFromReadOne": func(r *ackmodel.CRD, sourceResVarName string, targetResVarName string, indentLevel int) string {
			return code.LateInitializeFromReadOne(r.Config(), r, sourceResVarName, targetResVarName, indentLevel)
		},
		"GoCodeIncompleteLateInitialization": func(r *ackmodel.CRD, resVarName string, indentLevel int) string {
			return code.IncompleteLateInitialization(r.Config(), r, resVarName, indentLevel)
		},
		"GoCodeReferencesValidation": func(r *ackmodel.CRD, sourceVarName string, indentLevel int) string {
			return code.ReferenceFieldsValidation(r, sourceVarName, indentLevel)
		},
		"GoCodeContainsReferences": func(r *ackmodel.CRD, sourceVarName string) string {
			return code.ReferenceFieldsPresent(r, sourceVarName)
		},
		"CheckNilFieldPath": func(f *ackmodel.Field, sourceVarName string) string {
			return code.CheckNilFieldPath(f, sourceVarName)
		},
		"CheckNilReferencesPath": func(f *ackmodel.Field, sourceVarName string) string {
			return code.CheckNilReferencesPath(f, sourceVarName)
		},
                "Each": func (args ...interface{}) []interface{} {
                        return args
                },
	}
)

// Controller returns a pointer to a TemplateSet containing all the templates
// for generating ACK service controller implementations
func Controller(
	m *ackmodel.Model,
	templateBasePaths []string,
	// serviceAccountName is the name of the ServiceAccount used in the Helm chart
	serviceAccountName string,
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
		"references.go.tpl",
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
		serviceAccountName,
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
	// using Map to implement the Set
	referencedServiceNamesMap := make(map[string]struct{})
	for _, crd := range crds {
		snakeCasedCRDNames = append(snakeCasedCRDNames, crd.Names.Snake)
		for _, serviceName := range crd.ReferencedServiceNames() {
			referencedServiceNamesMap[serviceName] = struct{}{}
		}
	}
	referencedServiceNames := make([]string, 0)
	for serviceName := range referencedServiceNamesMap {
		referencedServiceNames = append(referencedServiceNames, serviceName)
	}
	sort.Strings(referencedServiceNames)
	cmdVars := &templateCmdVars{
		metaVars,
		snakeCasedCRDNames,
		referencedServiceNames,
	}
	if err = ts.Add("cmd/controller/main.go", "cmd/controller/main.go.tpl", cmdVars); err != nil {
		return nil, err
	}

	// Finally, add the configuration YAML file templates
	for _, path := range controllerConfigTemplatePaths {
		outPath := strings.TrimSuffix(path, ".tpl")
		if err = ts.Add(outPath, path, configVars); err != nil {
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
	// ReferencedServiceNames contains the service name for ACK controllers whose
	// resources are referenced inside the CRDs.
	// Service name is go package name of AWS service in aws-sdk-go.
	ReferencedServiceNames []string
}

// templateConfigVars contains template variables for the templates that require
// access to the generator configuration definition
type templateConfigVars struct {
	templateset.MetaVars
	GeneratorConfig    *ackgenconfig.Config
	ServiceAccountName string
}
