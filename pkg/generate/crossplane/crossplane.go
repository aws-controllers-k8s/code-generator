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

package crossplane

import (
	"path/filepath"
	"strings"
	ttpl "text/template"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/iancoleman/strcase"
)

var (
	apisGenericTemplatesPaths = []string{
		"crossplane/apis/doc.go.tpl",
		"crossplane/apis/enums.go.tpl",
		"crossplane/apis/groupversion_info.go.tpl",
		"crossplane/apis/types.go.tpl",
	}
	crdTemplatePath     = "crossplane/apis/crd.go.tpl"
	controllerTmplPath  = "crossplane/pkg/controller.go.tpl"
	conversionsTmplPath = "crossplane/pkg/conversions.go.tpl"
	includePaths        = []string{
		"crossplane/boilerplate.go.tpl",
		"crossplane/apis/enum_def.go.tpl",
		"crossplane/apis/type_def.go.tpl",
		"crossplane/pkg/sdk_find_read_one.go.tpl",
		"crossplane/pkg/sdk_find_read_many.go.tpl",
		"crossplane/pkg/sdk_find_get_attributes.go.tpl",
	}
	copyPaths = []string{}
	funcMap   = ttpl.FuncMap{
		"ToLower": strings.ToLower,
		"ResourceExceptionCode": func(r *ackmodel.CRD, httpStatusCode int) string {
			return r.ExceptionCode(httpStatusCode)
		},
		"GoCodeSetExceptionMessageCheck": func(r *ackmodel.CRD, httpStatusCode int) string {
			return code.CheckExceptionMessage(r.Config(), r, httpStatusCode)
		},
		"GoCodeSetReadOneOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetResource(r.Config(), r, ackmodel.OpTypeGet, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetReadOneInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeGet, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetReadManyOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetResource(r.Config(), r, ackmodel.OpTypeList, sourceVarName, targetVarName, indentLevel)
		},
		"ListMemberNameInReadManyOutput": func(r *ackmodel.CRD) (string, error) {
			return code.ListMemberNameInReadManyOutput(r)
		},
		"GoCodeSetReadManyInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeList, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeGetAttributesSetInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetSDKGetAttributes(r.Config(), r, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetAttributesSetInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetSDKSetAttributes(r.Config(), r, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeGetAttributesSetOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetResourceGetAttributes(r.Config(), r, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetCreateOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetResource(r.Config(), r, ackmodel.OpTypeCreate, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetCreateInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeCreate, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetUpdateInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeUpdate, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetDeleteInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeDelete, sourceVarName, targetVarName, indentLevel)
		},
		"Empty": func(subject string) bool {
			return strings.TrimSpace(subject) == ""
		},
	}
)

// templateAPIVars contains template variables for templates that output Go
// code in the /services/$SERVICE/apis/$API_VERSION directory
type templateAPIVars struct {
	templateset.MetaVars
	EnumDefs []*ackmodel.EnumDef
	TypeDefs []*ackmodel.TypeDef
}

// templateCRDVars contains template variables for the template that outputs Go
// code for a single top-level resource's API definition
type templateCRDVars struct {
	templateset.MetaVars
	CRD *ackmodel.CRD
}

// Crossplane returns a pointer to a TemplateSet containing all the templates for
// generating Crossplane API types and controller code for an AWS service API
func Crossplane(
	m *ackmodel.Model,
	templateBasePaths []string,
) (*templateset.TemplateSet, error) {
	enumDefs, err := m.GetEnumDefs()
	if err != nil {
		return nil, err
	}
	typeDefs, err := m.GetTypeDefs()
	if err != nil {
		return nil, err
	}
	crds, err := m.GetCRDs()
	if err != nil {
		return nil, err
	}

	ts := templateset.New(
		templateBasePaths,
		includePaths,
		copyPaths,
		funcMap,
	)

	metaVars := m.MetaVars()
	detectedCRDVersions := make(map[string]bool)
	for _, crd := range crds {
		v, err := crd.GetStorageVersion(metaVars.APIVersion)
		if err != nil {
			return nil, err
		}
		detectedCRDVersions[v] = true
	}

	// First add all the CRDs and API types
	for apiVersion := range detectedCRDVersions {
		for _, path := range apisGenericTemplatesPaths {
			apiVars := &templateAPIVars{
				metaVars,
				enumDefs,
				typeDefs,
			}
			apiVars.APIVersion = apiVersion
			outPath := filepath.Join(
				"apis",
				metaVars.ServicePackageName,
				apiVersion,
				"zz_"+strings.TrimSuffix(filepath.Base(path), ".tpl"),
			)
			if err = ts.Add(outPath, path, apiVars); err != nil {
				return nil, err
			}
		}
	}
	for _, crd := range crds {
		v, err := crd.GetStorageVersion(metaVars.APIVersion)
		if err != nil {
			return nil, err
		}
		crdFileName := filepath.Join(
			"apis", metaVars.ServicePackageName, v,
			"zz_"+strcase.ToSnake(crd.Kind)+".go",
		)
		crdVars := &templateCRDVars{
			metaVars,
			crd,
		}
		crdVars.APIVersion = v
		if err = ts.Add(crdFileName, crdTemplatePath, crdVars); err != nil {
			return nil, err
		}
	}

	// Next add the controller package for each CRD
	for _, crd := range crds {
		outPath := filepath.Join(
			"pkg", "controller", metaVars.ServicePackageName, crd.Names.Lower,
			"zz_controller.go",
		)
		crdVars := &templateCRDVars{
			metaVars,
			crd,
		}
		if crdVars.APIVersion, err = crd.GetStorageVersion(metaVars.APIVersion); err != nil {
			return nil, err
		}
		if err = ts.Add(outPath, controllerTmplPath, crdVars); err != nil {
			return nil, err
		}
		outPath = filepath.Join(
			"pkg", "controller", metaVars.ServicePackageName, crd.Names.Lower,
			"zz_conversions.go",
		)
		crdVars = &templateCRDVars{
			metaVars,
			crd,
		}
		if crdVars.APIVersion, err = crd.GetStorageVersion(metaVars.APIVersion); err != nil {
			return nil, err
		}
		if err = ts.Add(outPath, conversionsTmplPath, crdVars); err != nil {
			return nil, err
		}
	}

	return ts, nil
}
