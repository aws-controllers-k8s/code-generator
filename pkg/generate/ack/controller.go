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
	"path/filepath"
	"sort"
	"strings"
	ttpl "text/template"
	"time"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/fieldpath"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
	"github.com/aws-controllers-k8s/pkg/names"
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
		"config/rbac/leader-election-role-binding.yaml.tpl",
		"config/rbac/leader-election-role.yaml.tpl",
		"config/rbac/kustomization.yaml.tpl",
		"config/crd/kustomization.yaml.tpl",
		"config/overlays/namespaced/kustomization.yaml.tpl",
	}
	controllerIncludePaths = []string{
		"boilerplate.go.tpl",
		"pkg/resource/references_read_referenced_resource.go.tpl",
		"pkg/resource/sdk_delete_custom.go.tpl",
		"pkg/resource/sdk_find_custom.go.tpl",
		"pkg/resource/sdk_find_read_one.go.tpl",
		"pkg/resource/sdk_find_get_attributes.go.tpl",
		"pkg/resource/sdk_find_read_many.go.tpl",
		"pkg/resource/sdk_find_not_implemented.go.tpl",
		"pkg/resource/sdk_update.go.tpl",
		"pkg/resource/sdk_update_custom.go.tpl",
		"pkg/resource/sdk_update_set_attributes.go.tpl",
		"pkg/resource/sdk_update_not_implemented.go.tpl",
		"pkg/resource/sdk_update_sub_resource_sync.go.tpl",
		"pkg/resource/sdk_delete_sub_resource_sync.go.tpl",
		"pkg/resource/sdk_find_sub_resource_get.go.tpl",
		"pkg/resource/sub_resource_manager_scalar.go.tpl",
		"pkg/resource/sub_resource_manager_struct.go.tpl",
		"pkg/resource/sub_resource_manager_list_scalar.go.tpl",
		"pkg/resource/sub_resource_manager_list_struct.go.tpl",
		"pkg/resource/sub_resource_manager_map.go.tpl",
		"pkg/resource/sub_resource_manager_map_scalar.go.tpl",
		"pkg/resource/sub_resource_manager_map_struct.go.tpl",
	}
	controllerCopyPaths = []string{}
	controllerFuncMap   = ttpl.FuncMap{
		"ToLower": strings.ToLower,
		"TrimPrefix": func(s string, prefix string) string {
			return strings.TrimPrefix(s, prefix)
		},
		"HasPrefix": func(s string, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"ManagerMapper": func(r *ackmodel.CRD) []*ackgenconfig.MapperConfig {
			return r.Config().GetManagerMapper(r.Names.Original)
		},
		"ManagerReadFieldPath": func(r *ackmodel.CRD) string {
			return r.Config().GetManagerReadFieldPath(r.Names.Original)
		},
		// ManagerPrimaryKey derives the primary key for a sub-resource by
		// finding the mapper entry whose From matches the parent field path
		// or a special token ("$item.*" for struct lists, "$item" for scalar
		// lists, "$key" for maps), and returning its To. Struct field
		// patterns ($item.Field) take priority over bare $item.
		"ManagerPrimaryKey": func(r *ackmodel.CRD) string {
			fieldPath := r.Config().GetManagerParentFieldPath(r.Names.Original)
			mapper := r.Config().GetManagerMapper(r.Names.Original)
			// First pass: prefer $item.Field (struct field access)
			for _, m := range mapper {
				if strings.HasPrefix(m.From, "$item.") {
					return m.To
				}
			}
			// Second pass: bare $item, $key, or exact field path match
			for _, m := range mapper {
				if m.From == fieldPath || m.From == "$item" || m.From == "$key" {
					return m.To
				}
			}
			return ""
		},
		"IsSubResource": func(r *ackmodel.CRD) bool {
			return r.Config().IsSubResource(r.Names.Original)
		},
		// SubResourceManagerInfos returns a slice of SubResourceManagerInfo
		// for each sub-resource with a manager defined under the given CRD.
		// Returns nil if the CRD has no sub-resources with managers.
		"SubResourceManagerInfos": func(r *ackmodel.CRD) []SubResourceManagerInfo {
			subResources := r.Config().GetSubResources(r.Names.Original)
			if len(subResources) == 0 {
				return nil
			}
			var infos []SubResourceManagerInfo
			for internalName, subResCfg := range subResources {
				if subResCfg.Manager == nil {
					continue
				}
				fieldPath := r.Config().GetManagerParentFieldPath(internalName)
				snakeName := names.New(internalName).Snake
				infos = append(infos, SubResourceManagerInfo{
					FieldPath:   fieldPath,
					PackageName: snakeName,
				})
			}
			return infos
		},
		"Dereference": func(s *string) string {
			return *s
		},
		"AddToMap": func(m map[string]interface{}, k string, v interface{}) map[string]interface{} {
			if len(m) == 0 {
				m = make(map[string]interface{})
			}
			m[k] = v
			return m
		},
		"Nil": func() interface{} {
			return nil
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
		"GoCodeSetReadOneOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetResource(r.Config(), r, ackmodel.OpTypeGet, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetReadOneInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeGet, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetReadManyOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetResource(r.Config(), r, ackmodel.OpTypeList, sourceVarName, targetVarName, indentLevel)
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
		"GoCodeSetUpdateOutput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetResource(r.Config(), r, ackmodel.OpTypeUpdate, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetUpdateInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeUpdate, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetDeleteInput": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetSDK(r.Config(), r, ackmodel.OpTypeDelete, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeSetSDKForStruct": func(r *ackmodel.CRD, targetFieldName string, targetVarName string, targetShapeRef *awssdkmodel.ShapeRef, sourceFieldPath string, sourceVarName string, indentLevel int) (string, error) {
			return code.SetSDKForStruct(r.Config(), r, targetFieldName, targetVarName, targetShapeRef, sourceFieldPath, sourceVarName, model.OpTypeList, indentLevel)
		},
		"GoCodeSetResourceForStruct": func(r *ackmodel.CRD, targetFieldName string, targetVarName string, targetShapeRef *awssdkmodel.ShapeRef, sourceVarName string, sourceShapeRef *awssdkmodel.ShapeRef, indentLevel int) (string, error) {
			var setCfg *ackgenconfig.SetFieldConfig = nil
			// We may have some special instructions for how to handle setting the
			// field value...
			//
			// We do not want to return an empty string if the setConfig is not provided,
			// so that we can allow non top-level fields to be used with this function.
			f, ok := r.Fields[targetFieldName]
			if ok {
				setCfg = f.GetSetterConfig(ackmodel.OpTypeList)
			}
			if setCfg != nil && setCfg.IgnoreResourceSetter() {
				return "", nil
			}
			return code.SetResourceForStruct(r.Config(), r, targetVarName, targetShapeRef, setCfg, sourceVarName, sourceShapeRef, "", model.OpTypeList, indentLevel)
		},
		"GoCodeCompare": func(r *ackmodel.CRD, deltaVarName string, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.CompareResource(r.Config(), r, deltaVarName, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodeIsSynced": func(r *ackmodel.CRD, resVarName string, indentLevel int) (string, error) {
			return code.ResourceIsSynced(r.Config(), r, resVarName, indentLevel)
		},
		"GoCodeCompareStruct": func(r *ackmodel.CRD, shape *awssdkmodel.Shape, deltaVarName string, sourceVarName string, targetVarName string, fieldPath string, indentLevel int) (string, error) {
			return code.CompareStruct(r.Config(), r, nil, shape, deltaVarName, sourceVarName, targetVarName, fieldPath, indentLevel)
		},
		"Empty": func(subject string) bool {
			return strings.TrimSpace(subject) == ""
		},
		"GoCodeRequiredFieldsMissingFromReadOneInput": func(r *ackmodel.CRD, koVarName string, indentLevel int) (string, error) {
			return code.CheckRequiredFieldsMissingFromShape(r, ackmodel.OpTypeGet, koVarName, indentLevel)
		},
		"GoCodeRequiredFieldsMissingFromReadManyInput": func(r *ackmodel.CRD, koVarName string, indentLevel int) (string, error) {
			return code.CheckRequiredFieldsMissingFromShape(r, ackmodel.OpTypeList, koVarName, indentLevel)
		},
		"GoCodeRequiredFieldsMissingFromGetAttributesInput": func(r *ackmodel.CRD, koVarName string, indentLevel int) (string, error) {
			return code.CheckRequiredFieldsMissingFromShape(r, ackmodel.OpTypeGetAttributes, koVarName, indentLevel)
		},
		"GoCodeRequiredFieldsMissingFromSetAttributesInput": func(r *ackmodel.CRD, koVarName string, indentLevel int) (string, error) {
			return code.CheckRequiredFieldsMissingFromShape(r, ackmodel.OpTypeSetAttributes, koVarName, indentLevel)
		},
		"GoCodeSetResourceIdentifiers": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.SetResourceIdentifiers(r.Config(), r, sourceVarName, targetVarName, indentLevel)
		},
		"GoCodePopulateResourceFromAnnotation": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.PopulateResourceFromAnnotation(r.Config(), r, sourceVarName, targetVarName, indentLevel)
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
		"GoCodeReferencesValidation": func(f *ackmodel.Field, sourceVarName string, indentLevel int) (string, error) {
			return code.ReferenceFieldsValidation(f, sourceVarName, indentLevel)
		},
		"CheckNilFieldPath": func(f *ackmodel.Field, sourceVarName string) string {
			return code.CheckNilFieldPath(f, sourceVarName)
		},
		"CheckNilReferencesPath": func(f *ackmodel.Field, sourceVarName string) string {
			return code.CheckNilReferencesPath(f, sourceVarName)
		},
		"Each": func(args ...interface{}) []interface{} {
			return args
		},
		"GoCodeInitializeNestedStructField": func(r *ackmodel.CRD,
			sourceVarName string, f *ackmodel.Field, apiPkgImportName string,
			indentLevel int) (string, error) {
			return code.InitializeNestedStructField(r, sourceVarName, f,
				apiPkgImportName, indentLevel)
		},
		"GoCodeResolveReference": func(f *ackmodel.Field, sourceVarName string, indentLevel int) (string, error) {
			return code.ResolveReferencesForField(f, sourceVarName, indentLevel)
		},
		"GoCodeClearResolvedReferences": func(f *ackmodel.Field, targetVarName string, indentLevel int) (string, error) {
			return code.ClearResolvedReferencesForField(f, targetVarName, indentLevel)
		},
		"GoCodeConvertToACKTags": func(r *ackmodel.CRD, sourceVarName string, targetVarName string, keyOrderVarName string, indentLevel int) (string, error) {
			return code.GoCodeConvertToACKTags(r, sourceVarName, targetVarName, keyOrderVarName, indentLevel)
		},
		"GoCodeFromACKTags": func(r *ackmodel.CRD, sourceVarName string, orderVarName string, targetVarName string, indentLevel int) (string, error) {
			return code.GoCodeFromACKTags(r, sourceVarName, orderVarName, targetVarName, indentLevel)
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
	totalStart := time.Now()

	crdStart := time.Now()
	crds, err := m.GetCRDs()
	if err != nil {
		return nil, err
	}
	util.Tracef("GetCRDs (%d CRDs): %s\n", len(crds), time.Since(crdStart))

	tplStart := time.Now()
	metaVars := m.MetaVars()

	// Build a lookup of CRDs by name so ManagerSourceType can resolve
	// the parent CRD's field shape type.
	crdsByName := make(map[string]*ackmodel.CRD, len(crds))
	for _, crd := range crds {
		crdsByName[crd.Names.Original] = crd
	}

	// ManagerSourceType returns a SourceTypeInfo for the parent field that
	// feeds the sub-resource conversion. It derives the field path from the
	// sub-resource key name and inspects the parent CRD's spec field shape to
	// determine whether the source is a list, map, or string.
	controllerFuncMap["ManagerSourceType"] = func(r *ackmodel.CRD) (*ackgenconfig.SourceTypeInfo, error) {
		fieldPath := r.Config().GetManagerParentFieldPath(r.Names.Original)
		parentName := r.Config().GetParentResourceName(r.Names.Original)
		if parentName == "" {
			return nil, fmt.Errorf("manager source type not implemented: no parent resource for %s", r.Names.Original)
		}
		parentCRD, ok := crdsByName[parentName]
		if !ok {
			return nil, fmt.Errorf("manager source type not implemented: parent CRD %s not found for %s", parentName, r.Names.Original)
		}
		// fieldPath is "Spec.FieldName" — extract the field name
		parts := strings.SplitN(fieldPath, ".", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("manager source type not implemented: invalid field path %q for %s", fieldPath, r.Names.Original)
		}
		f, ok := parentCRD.SpecFields[parts[1]]
		if !ok || f.ShapeRef == nil || f.ShapeRef.Shape == nil {
			return nil, fmt.Errorf("manager source type not implemented: field %s not found in parent %s for %s", parts[1], parentName, r.Names.Original)
		}
		info := &ackgenconfig.SourceTypeInfo{FieldPath: fieldPath, ParentKind: parentCRD.Kind}
		switch f.ShapeRef.Shape.Type {
		case "map":
			if f.ShapeRef.Shape.ValueRef.Shape != nil && f.ShapeRef.Shape.ValueRef.Shape.Type == "structure" {
				info.Type = ackgenconfig.SourceTypeMapStruct
			} else {
				info.Type = ackgenconfig.SourceTypeMapScalar
			}
		case "list":
			if f.ShapeRef.Shape.MemberRef.Shape != nil && f.ShapeRef.Shape.MemberRef.Shape.Type == "structure" {
				info.Type = ackgenconfig.SourceTypeListStruct
			} else {
				info.Type = ackgenconfig.SourceTypeListScalar
			}
		case "structure":
			info.Type = ackgenconfig.SourceTypeStruct
		case "string", "boolean", "integer", "long", "float", "double", "timestamp", "blob":
			info.Type = ackgenconfig.SourceTypeScalar
		default:
			return nil, fmt.Errorf("manager source type not implemented for shape type %q at %s in %s", f.ShapeRef.Shape.Type, fieldPath, r.Names.Original)
		}

		// Derive BatchFieldPath: if the mapper has a "$item" entry whose
		// target field is a list on the sub-resource CRD, the SDK operation
		// accepts multiple items per call and we can batch.
		mapper := r.Config().GetManagerMapper(r.Names.Original)
		for _, m := range mapper {
			if m.From == "$item" {
				// Check if the target field is a list on the sub-resource CRD.
				// The To path is "Spec.FieldName" — extract the field name.
				toParts := strings.SplitN(m.To, ".", 2)
				if len(toParts) == 2 {
					if sf, ok := r.SpecFields[toParts[1]]; ok {
						if sf.ShapeRef != nil && sf.ShapeRef.Shape != nil && sf.ShapeRef.Shape.Type == "list" {
							info.BatchFieldPath = m.To
						}
					}
				}
				break
			}
		}

		return info, nil
	}

	// Hook code can reference a template path, and we can look up the template
	// in any of our base paths...
	controllerFuncMap["Hook"] = func(r *ackmodel.CRD, hookID string) (string, error) {
		crdVars := &templateCRDVars{
			metaVars,
			m.SDKAPI,
			r,
		}
		code, err := ResourceHookCode(templateBasePaths, r, hookID, crdVars, controllerFuncMap)
		if err != nil {
			return "", err
		}
		return code, nil
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
		"tags.go.tpl",
	}
	for _, crd := range crds {
		isSubRes := crd.Config().IsSubResource(crd.Names.Original)

		// Determine output base path: sub-resources nest under parent
		var outBase string
		if isSubRes {
			parentName := crd.Config().GetParentResourceName(crd.Names.Original)
			parentSnake := names.New(parentName).Snake
			outBase = filepath.Join("pkg/resource", parentSnake, crd.Names.Snake)
		} else {
			outBase = filepath.Join("pkg/resource", crd.Names.Snake)
		}

		for _, target := range targets {
			// skip adding "tags.go.tpl" file if tagging is ignored for a crd
			if target == "tags.go.tpl" && crd.Config().TagsAreIgnored(crd.Names.Original) {
				continue
			}
			// Sub-resources only need delta.go and sdk.go — the resource
			// struct, resourceManager, and all other scaffolding are
			// generated directly in the sub-resource manager file.
			if isSubRes && target != "delta.go.tpl" && target != "sdk.go.tpl" {
				continue
			}
			outPath := filepath.Join(outBase, strings.TrimSuffix(target, ".tpl"))
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
		// Generate manager.go for sub-resources when the manager field is specified
		if crd.Config().GetManagerName(crd.Names.Original) != "" {
			outPath := filepath.Join(outBase, "manager.go")
			tplPath := filepath.Join("pkg/resource", "sub_resource_manager.go.tpl")
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
		// Exclude sub-resource packages from main.go imports — they are
		// imported by the parent resource's hooks, not by the main controller.
		if crd.Config().IsSubResource(crd.Names.Original) {
			continue
		}
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
	util.Tracef("template setup (%d templates): %s\n", len(ts.Executed())+len(controllerConfigTemplatePaths), time.Since(tplStart))
	util.Tracef("Controller() total: %s\n", time.Since(totalStart))
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

// SubResourceManagerInfo holds the information needed by templates to generate
// sub-resource manager sync code in the parent resource's sdkUpdate function.
type SubResourceManagerInfo struct {
	// FieldPath is the parent spec field path (e.g. "Spec.Policies")
	FieldPath string
	// PackageName is the snake_case package name for the sub-resource manager
	// (e.g. "role_policies")
	PackageName string
}
