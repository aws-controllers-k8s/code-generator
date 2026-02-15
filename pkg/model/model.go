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

package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/pkg/names"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	ackfp "github.com/aws-controllers-k8s/code-generator/pkg/fieldpath"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

var (
	// ErrNilShapePointer indicates an unexpected nil Shape pointer
	ErrNilShapePointer = errors.New("found nil Shape pointer")
)

// Model contains the ACK model for the generator to process and apply
// templates against.
type Model struct {
	SDKAPI             *SDKAPI
	servicePackageName string
	apiVersion         string
	crds               []*CRD
	typeDefs           []*TypeDef
	typeImports        map[string]string
	typeRenames        map[string]string
	// Instructions to the code generator how to handle the API and its
	// resources
	cfg    *ackgenconfig.Config
	docCfg *ackgenconfig.DocumentationConfig
}

// MetaVars returns a MetaVars struct populated with metadata about the AWS
// service API
func (m *Model) MetaVars() templateset.MetaVars {
	controllerName := m.cfg.ControllerName
	if controllerName == "" {
		controllerName = m.servicePackageName
	}
	// NOTE(a-hilaly): I know this is a bit of a hack and it's confusing, but
	// long time ago, we assumed that model_name is always equal to the service
	// name. This is not the case anymore, prometheusservice and documentdb
	// are examples of services that have different model names.
	//
	// TODO(a-hilaly): We should probably rework all this naming stuff to be
	// more consistent. To whoever is reading this, I'm sorry.
	servicePackageName := m.servicePackageName
	if m.cfg.SDKNames.Package != "" {
		servicePackageName = m.cfg.SDKNames.Package
	}
	return templateset.MetaVars{
		ControllerName:          controllerName,
		ServicePackageName:      servicePackageName,
		ServiceID:               m.SDKAPI.ServiceID(),
		ServiceModelName:        m.cfg.SDKNames.Model,
		APIGroup:                m.APIGroup(),
		APIVersion:              m.apiVersion,
		ClientInterfaceTypeName: m.ClientInterfaceTypeName(),
		ClientStructTypeName:    m.ClientStructTypeName(),
		CRDNames:                m.crdNames(),
	}
}

// crdNames returns all crd names lowercased and in plural
func (m *Model) crdNames() []string {
	var crdConfigs []string

	crds, _ := m.GetCRDs()
	for _, crd := range crds {
		crdConfigs = append(crdConfigs, strings.ToLower(crd.Plural))
	}

	return crdConfigs
}

// GetCRDs returns a slice of `CRD` structs that describe the
// top-level resources discovered by the code generator for an AWS service API
func (m *Model) GetCRDs() ([]*CRD, error) {
	if m.crds != nil {
		return m.crds, nil
	}
	crds := []*CRD{}

	opMap, err := m.SDKAPI.GetOperationMap(m.cfg)
	if err != nil {
		return nil, err
	}

	createOps := (*opMap)[OpTypeCreate]
	readOneOps := (*opMap)[OpTypeGet]
	readManyOps := (*opMap)[OpTypeList]
	updateOps := (*opMap)[OpTypeUpdate]
	deleteOps := (*opMap)[OpTypeDelete]
	getAttributesOps := (*opMap)[OpTypeGetAttributes]
	setAttributesOps := (*opMap)[OpTypeSetAttributes]

	for crdName, createOp := range createOps {
		if m.cfg.ResourceIsIgnored(crdName) {
			continue
		}
		crdNames := names.New(crdName)
		ops := Ops{
			Create:        createOps[crdName],
			ReadOne:       readOneOps[crdName],
			ReadMany:      readManyOps[crdName],
			Update:        updateOps[crdName],
			Delete:        deleteOps[crdName],
			GetAttributes: getAttributesOps[crdName],
			SetAttributes: setAttributesOps[crdName],
		}
		m.RemoveIgnoredOperations(&ops)
		crd := NewCRD(m.SDKAPI, m.cfg, m.docCfg, crdNames, ops)

		// OK, begin to gather the CRDFields that will go into the Spec struct.
		// These fields are those members of the Create operation's Input
		// Shape.
		inputShape := createOp.InputRef.Shape
		if inputShape == nil {
			return nil, ErrNilShapePointer
		}

		// Check if there's an input wrapper field path configured. If so, we
		// flatten the wrapper's fields into the Spec instead of creating a
		// nested structure.
		//
		// NOTE(input_wrapper_field_path): We intentionally don't reuse
		// CRD.GetInputShape() here because this code runs during CRD construction
		// before the CRD is fully initialized. GetInputShape is designed for use
		// after CRD construction (e.g., in code generation). The logic is similar
		// but operates at different lifecycle stages.
		inputWrapperFieldPath := m.cfg.GetInputWrapperFieldPath(createOp)

		for memberName, memberShapeRef := range inputShape.MemberRefs {
			if memberShapeRef.Shape == nil {
				return nil, ErrNilShapePointer
			}

			// If this is the wrapper field and we have input_wrapper_field_path
			// configured, add the wrapper's member fields instead of the wrapper
			if inputWrapperFieldPath != nil && memberName == *inputWrapperFieldPath {
				wrapperShape := memberShapeRef.Shape
				// NOTE(input_wrapper_field_path): Currently only structure wrappers
				// are supported. If needed, support for nested paths (e.g., a.b.c
				// where b is a list) could be added in a future PR by extending
				// this logic to handle list/map types similar to getWrapperShape
				// in crd.go.
				if wrapperShape.Type == "structure" {
					for wrapperMemberName, wrapperMemberShapeRef := range wrapperShape.MemberRefs {
						if wrapperMemberShapeRef.Shape == nil {
							return nil, ErrNilShapePointer
						}
						// Handles field renames, if applicable
						fieldName := m.cfg.GetResourceFieldName(
							crd.Names.Original,
							createOp.Name,
							wrapperMemberName,
						)
						wrapperMemberNames := names.New(fieldName)
						if err := crd.AddSpecField(wrapperMemberNames, wrapperMemberShapeRef); err != nil {
							return nil, err
						}
					}
				}
				continue
			}

			// When input_wrapper_field_path is configured, skip fields that are
			// not part of the wrapper. This is consistent with output_wrapper_field_path
			// behavior - only the wrapper's fields are flattened into the CRD.
			if inputWrapperFieldPath != nil {
				continue
			}

			// Handles field renames, if applicable
			fieldName := m.cfg.GetResourceFieldName(
				crd.Names.Original,
				createOp.Name,
				memberName,
			)
			memberNames := names.New(fieldName)
			if memberName == "Attributes" && m.cfg.ResourceContainsAttributesMap(crdName) {
				if err := crd.UnpackAttributes(); err != nil {
					return nil, err
				}
				continue
			}
			if err := crd.AddSpecField(memberNames, memberShapeRef); err != nil {
				return nil, err
			}
		}

		// A list of fields that should be processed after gathering
		// the Spec and Status top level fields. The customNestedFields will be
		// injected into the Spec or Status struct as a nested field.
		//
		// Note that we could reuse the Field struct here, but we don't because
		// we don't need all the fields that the Field struct provides. We only
		// need the field path and the FieldConfig. Using Field could lead to
		// confusion.
		customNestedFields := make(map[string]*ackgenconfig.FieldConfig)

		for targetFieldName, fieldConfig := range m.cfg.GetFieldConfigs(crdName) {
			if fieldConfig.IsReadOnly {
				// It's a Status field...
				continue
			}

			var found bool
			var memberShapeRef *awssdkmodel.ShapeRef

			if fieldConfig.From != nil {
				from := fieldConfig.From
				memberShapeRef, found = m.SDKAPI.GetInputShapeRef(
					from.Operation, from.Path,
				)
				// allowing getting spec fields from output shape
				if !found {
					memberShapeRef, found = m.SDKAPI.GetOutputShapeRef(
						from.Operation, from.Path,
					)
				}
				if !found {
					return nil, fmt.Errorf(
						"resource %q: unknown Spec field source — operation %q path %q not found in input or output shapes",
						crdName, from.Operation, from.Path,
					)
				}
			} else if fieldConfig.CustomField != nil {
				customField := fieldConfig.CustomField
				if customField.ListOf != "" {
					memberShapeRef = m.SDKAPI.GetCustomShapeRef(customField.ListOf)
				} else {
					memberShapeRef = m.SDKAPI.GetCustomShapeRef(customField.MapOf)
				}
				if memberShapeRef == nil {
					return nil, fmt.Errorf(
						"resource %q: unknown custom Spec field — custom field %+v has no matching shape",
						crdName, customField,
					)
				}
			} else if fieldConfig.Type != nil {
				// A nested field will always have a "." in the field path.
				// Let's collect those fields and process them after we've
				// gathered all the top level fields.
				if strings.Contains(targetFieldName, ".") {
					// This is a nested field
					customNestedFields[targetFieldName] = fieldConfig
					continue
				}
				// If we're here, we have a custom top level field (non-nested).

				// We have a custom field that has a type override and has not
				// been inferred via the normal Create Input shape or via the
				// SourceFieldConfig. Manually construct the field and its
				// shape reference here.
				typeOverride := *fieldConfig.Type
				var err error
				memberShapeRef, err = m.SDKAPI.GetShapeRefFromType(typeOverride)
				if err != nil {
					return nil, fmt.Errorf("resource %q, field %q: %w", crdName, targetFieldName, err)
				}
			} else {
				// Spec field is not well defined
				continue
			}

			memberNames := names.New(targetFieldName)
			if err := crd.AddSpecField(memberNames, memberShapeRef); err != nil {
				return nil, err
			}

		}

		// Now process the fields that will go into the Status struct. We want
		// fields that are in the Create operation's Output Shape but that are
		// not in the Input Shape.
		outputShape, err := crd.GetOutputShape(createOp)
		if err != nil {
			return nil, err
		}
		if outputShape.UsedAsOutput && len(outputShape.MemberRefs) == 1 {
			// We might be in a "wrapper" shape. Unwrap it to find the real object
			// representation for the CRD's createOp. If there is a single member
			// shape and that member shape is a structure, unwrap it.
			for _, memberRef := range outputShape.MemberRefs {
				if memberRef.Shape.Type == "structure" {
					outputShape = memberRef.Shape
				}
			}
		}
		for memberName, memberShapeRef := range outputShape.MemberRefs {
			if memberShapeRef.Shape == nil {
				return nil, ErrNilShapePointer
			}
			// Check that the field in the output shape isn't the same as
			// fields in the input shape (handles field renames, if applicable)
			fieldName := m.cfg.GetResourceFieldName(
				crd.Names.Original,
				createOp.Name,
				memberName,
			)
			if inSpec, _ := crd.HasMember(fieldName, createOp.Name); inSpec {
				// We don't put fields that are already in the Spec struct into
				// the Status struct
				continue
			}
			memberNames := names.New(fieldName)

			//TODO:(brycahta) should we support overriding these fields?
			if memberName == "Attributes" && m.cfg.ResourceContainsAttributesMap(crdName) {
				continue
			}
			if crd.IsPrimaryARNField(memberName) {
				// We automatically place the primary resource ARN value into
				// the Status.ACKResourceMetadata.ARN field
				continue
			}
			if err := crd.AddStatusField(memberNames, memberShapeRef); err != nil {
				return nil, err
			}
		}

		// Now add the additional Status fields that are required from other
		// API operations.
		for targetFieldName, fieldConfig := range m.cfg.GetFieldConfigs(crdName) {
			if !fieldConfig.IsReadOnly {
				// It's a Spec field...
				continue
			}

			var found bool
			var memberShapeRef *awssdkmodel.ShapeRef

			if fieldConfig.From != nil {
				from := fieldConfig.From
				memberShapeRef, found = m.SDKAPI.GetOutputShapeRef(
					from.Operation, from.Path,
				)
				// allowing to get status fields from output shapes
				if !found {
					memberShapeRef, found = m.SDKAPI.GetInputShapeRef(
						from.Operation, from.Path,
					)
				}
				if !found {
					return nil, fmt.Errorf(
						"resource %q: unknown Status field source — operation %q path %q not found in input or output shapes",
						crdName, from.Operation, from.Path,
					)
				}
			} else if fieldConfig.CustomField != nil {
				customField := fieldConfig.CustomField
				if customField.ListOf != "" {
					memberShapeRef = m.SDKAPI.GetCustomShapeRef(customField.ListOf)
				} else {
					memberShapeRef = m.SDKAPI.GetCustomShapeRef(customField.MapOf)
				}
				if memberShapeRef == nil {
					return nil, fmt.Errorf(
						"resource %q: unknown custom Status field — custom field %+v has no matching shape",
						crdName, customField,
					)
				}
			} else if fieldConfig.Type != nil {
				// A nested field will always have a "." in the field path.
				// Let's collect those fields and process them after we've
				// gathered all the top level fields.
				if strings.Contains(targetFieldName, ".") {
					// This is a nested field
					customNestedFields[targetFieldName] = fieldConfig
					continue
				}
				// If we're here, we have a custom top level field (non-nested).

				// We have a custom field that has a type override and has not
				// been inferred via the normal Create Input shape or via the
				// SourceFieldConfig. Manually construct the field and its
				// shape reference here.
				typeOverride := *fieldConfig.Type
				var err error
				memberShapeRef, err = m.SDKAPI.GetShapeRefFromType(typeOverride)
				if err != nil {
					return nil, fmt.Errorf("resource %q, field %q: %w", crdName, targetFieldName, err)
				}
			} else {
				// Status field is not well defined
				continue
			}

			memberNames := names.New(targetFieldName)
			if err := crd.AddStatusField(memberNames, memberShapeRef); err != nil {
				return nil, err
			}
		}

		// Now add the additional printer columns that have been defined explicitly
		// in additional_columns
		crd.addAdditionalPrinterColumns(m.cfg.GetAdditionalColumns(crdName))
		// Process the custom nested fields
		if err := crd.addCustomNestedFields(customNestedFields); err != nil {
			return nil, err
		}
		crds = append(crds, crd)
	}
	sort.Slice(crds, func(i, j int) bool {
		return crds[i].Names.Camel < crds[j].Names.Camel
	})
	// This is the place that we build out the CRD.Fields map with
	// `pkg/model.Field` objects that represent the non-top-level Spec and
	// Status fields.
	if err := m.processFields(crds); err != nil {
		return nil, err
	}
	m.crds = crds
	return crds, nil
}

// RemoveIgnoredOperations updates Ops argument by setting those
// operations to nil that are configured to be ignored in generator config for
// the AWS service
func (m *Model) RemoveIgnoredOperations(ops *Ops) {
	if m.cfg.OperationIsIgnored(ops.Create) {
		ops.Create = nil
	}
	if m.cfg.OperationIsIgnored(ops.ReadOne) {
		ops.ReadOne = nil
	}
	if m.cfg.OperationIsIgnored(ops.ReadMany) {
		ops.ReadMany = nil
	}
	if m.cfg.OperationIsIgnored(ops.Update) {
		ops.Update = nil
	}
	if m.cfg.OperationIsIgnored(ops.Delete) {
		ops.Delete = nil
	}
	if m.cfg.OperationIsIgnored(ops.GetAttributes) {
		ops.GetAttributes = nil
	}
	if m.cfg.OperationIsIgnored(ops.SetAttributes) {
		ops.SetAttributes = nil
	}
}

// IsShapeUsedInCRDs returns true if the supplied shape name is a member of amy
// CRD's payloads or those payloads sub-member shapes
func (m *Model) IsShapeUsedInCRDs(shapeName string) bool {
	crds, _ := m.GetCRDs()
	for _, crd := range crds {
		if crd.HasShapeAsMember(shapeName) {
			return true
		}
	}
	return false
}

// GetTypeDefs returns a slice of `TypeDef` pointers
func (m *Model) GetTypeDefs() ([]*TypeDef, error) {
	if m.typeDefs != nil {
		return m.typeDefs, nil
	}

	tdefs := []*TypeDef{}
	// Map, keyed by original Shape GoTypeElem(), with the values being a
	// renamed type name (due to conflicting names)
	trenames := map[string]string{}

	payloads := m.SDKAPI.GetPayloads()

	for shapeName, shape := range m.SDKAPI.API.Shapes {
		if util.InStrings(shapeName, payloads) && !m.IsShapeUsedInCRDs(shapeName) {
			// Payloads are not type defs, unless explicitly used
			continue
		}
		if shape.Type != "structure" {
			continue
		}
		if shape.Exception {
			// Neither are exceptions
			continue
		}
		tdefNames := names.New(shapeName)
		if m.SDKAPI.HasConflictingTypeName(shapeName, m.cfg) {
			tdefNames.Camel += ConflictingNameSuffix
			trenames[shapeName] = tdefNames.Camel
		}

		attrs := map[string]*Attr{}
		for memberName, memberRef := range shape.MemberRefs {
			memberNames := names.New(memberName)
			memberShape := memberRef.Shape
			if !m.IsShapeUsedInCRDs(memberShape.ShapeName) {
				continue
			}
			gt, err := m.getShapeCleanGoType(memberShape)
			if err != nil {
				return nil, err
			}
			attrs[memberName] = NewAttr(memberNames, gt, memberShape)
		}
		if len(attrs) == 0 {
			// Just ignore these...
			continue
		}
		tdefs = append(tdefs, &TypeDef{
			Shape: shape,
			Names: tdefNames,
			Attrs: attrs,
		})
	}
	sort.Slice(tdefs, func(i, j int) bool {
		return tdefs[i].Names.Camel < tdefs[j].Names.Camel
	})
	if err := m.processNestedFieldTypeDefs(tdefs); err != nil {
		return nil, err
	}
	m.typeDefs = tdefs
	m.typeRenames = trenames
	return tdefs, nil
}

// getShapeCleanGoType returns a cleaned-up and Camel-cased GoType name for a given shape.
func (m *Model) getShapeCleanGoType(shape *awssdkmodel.Shape) (string, error) {
	switch shape.Type {
	case "map":
		// If it's a map type we need to set the GoType to the cleaned-up
		// Camel-cased name
		gt, err := m.getShapeCleanGoType(shape.ValueRef.Shape)
		if err != nil {
			return "", err
		}
		return "map[string]" + gt, nil
	case "list", "array":
		// If it's a list type, we need to set the GoType to the cleaned-up
		// Camel-cased name
		gt, err := m.getShapeCleanGoType(shape.MemberRef.Shape)
		if err != nil {
			return "", err
		}
		return "[]" + gt, nil
	case "timestamp":
		// time.Time needs to be converted to apimachinery/metav1.Time
		// otherwise there is no DeepCopy support
		return "*metav1.Time", nil
	case "structure":
		if len(shape.MemberRefs) == 0 {
			if m.cfg.HasEmptyShape(shape.ShapeName) {
				return "map[string]*string", nil
			}
			return "", fmt.Errorf(
				"structure %q has no fields — configure it as an empty_shape or manually set the field type in generator.yaml",
				shape.ShapeName,
			)
		}
		// There are shapes that are called things like DBProxyStatus that are
		// fields in a DBProxy CRD... we need to ensure the type names don't
		// conflict. Also, the name of the Go type in the generated code is
		// Camel-cased and normalized, so we use that as the Go type
		goType := shape.GoType()
		typeNames := names.New(goType)
		if m.SDKAPI.HasConflictingTypeName(goType, m.cfg) {
			typeNames.Camel += ConflictingNameSuffix
		}
		return "*" + typeNames.Camel, nil
	default:
		return shape.GoType(), nil
	}
}

// processNestedFieldTypeDefs updates the supplied TypeDef structs' if a nested
// field has been configured with a type overriding FieldConfig -- such as
// FieldConfig.IsSecret.
func (m *Model) processNestedFieldTypeDefs(
	tdefs []*TypeDef,
) error {
	crds, err := m.GetCRDs()
	if err != nil {
		return err
	}
	for _, crd := range crds {
		for fieldPath, field := range crd.Fields {
			if !strings.Contains(fieldPath, ".") {
				// top-level fields have already had their structure
				// transformed during the CRD.AddSpecField and
				// CRD.AddStatusField methods. All we need to do here is look
				// at nested fields, which are identifiable as fields with
				// field paths contains a dot (".")
				continue
			}
			if field.FieldConfig == nil {
				// Likewise, we don't need to transform any TypeDef if the
				// nested field doesn't have a FieldConfig instructing us to
				// treat this field differently.
				continue
			}
			if field.FieldConfig.IsSecret {
				// Find the TypeDef that was created for the *containing*
				// secret field struct. For example, assume the nested field
				// path `Users..Password`, we'd want to find the TypeDef that
				// was created for the `Users` field's element type (which is a
				// struct)
				if err := replaceSecretAttrGoType(crd, field, tdefs); err != nil {
					return fmt.Errorf("resource %q, field %q: %w", crd.Names.Original, fieldPath, err)
				}
			}
			if field.FieldConfig.References != nil {
				if err := updateTypeDefAttributeWithReference(crd, fieldPath, tdefs); err != nil {
					return fmt.Errorf("resource %q, field %q: %w", crd.Names.Original, fieldPath, err)
				}
			}
			if field.FieldConfig.GoTag != nil {
				if err := setTypeDefAttributeGoTag(crd, fieldPath, field, tdefs); err != nil {
					return fmt.Errorf("resource %q, field %q: %w", crd.Names.Original, fieldPath, err)
				}
			}
			if field.IsImmutable() {
				if err := setTypeDefAttributeImmutable(crd, fieldPath, tdefs); err != nil {
					return fmt.Errorf("resource %q, field %q: %w", crd.Names.Original, fieldPath, err)
				}
			}
		}
	}
	return nil
}

// getAttributeFromPath extracts the parent TypeDef and the target attribute for
// the corresponding fieldPath of nested field. This function should only be
// called for nested fieldPath. Non-nested fieldPath should be handled by higher
// level functions.
func getAttributeFromPath(crd *CRD, fieldPath string, tdefs []*TypeDef) (parentTypeDef *TypeDef, target *Attr, err error) {
	fp := ackfp.FromString(fieldPath)
	if fp.Size() < 2 {
		// This function should only be called for nested fieldPath. Non-nested
		// fieldPath should be handled by higher level functions.
		return nil, nil, nil
	}
	// First part of nested reference fieldPath is the name of top level Spec
	// field. Ex: For 'ResourcesVpcConfig.SecurityGroupIds' fieldpath the
	// topLevelFieldName is 'ResourcesVpcConfig'
	topLevelFieldName := fp.Front()
	var topLevelField *Field
	foundInSpec := false
	foundInStatus := false
	for fName, field := range crd.SpecFields {
		if strings.EqualFold(fName, topLevelFieldName) {
			topLevelField = field
			foundInSpec = true
			break
		}
	}
	if !foundInSpec {
		for fName, field := range crd.StatusFields {
			if strings.EqualFold(fName, topLevelFieldName) {
				topLevelField = field
				foundInStatus = true
			}
		}
	}
	if !foundInSpec && !foundInStatus {
		return nil, nil, fmt.Errorf(
			"unable to find a spec or status field with name %s for path %s",
			topLevelFieldName, fieldPath,
		)
	}
	if foundInSpec && foundInStatus {
		return nil, nil, fmt.Errorf(
			"field %s in path %s exists in both spec and status",
			topLevelFieldName, fieldPath,
		)
	}

	// Create a new fieldPath starting with ShapeName of Spec Field
	// to determine the shape of typedef which will contain the reference
	// attribute. We replace the spec-field Name with spec-field ShapeName in
	// the beginning of field path and leave rest of nested member names as is.
	// Ex: ResourcesVpcConfig.SecurityGroupIDs will become VPCConfigRequest.SecurityGroupIDs
	// for Cluster resource in eks-controller.
	specFieldShapeRef := topLevelField.ShapeRef
	specFieldShapeName := specFieldShapeRef.ShapeName
	switch shapeType := specFieldShapeRef.Shape.Type; shapeType {
	case "list":
		specFieldShapeName = topLevelField.ShapeRef.Shape.MemberRef.ShapeName
		specFieldShapeRef = &topLevelField.ShapeRef.Shape.MemberRef
	case "map":
		specFieldShapeName = topLevelField.ShapeRef.Shape.ValueRef.ShapeName
		specFieldShapeRef = &topLevelField.ShapeRef.Shape.ValueRef
	}
	fieldShapePath := strings.Replace(fieldPath, topLevelFieldName, specFieldShapeName, 1)
	fsp := ackfp.FromString(fieldShapePath)

	// "fieldName" is the member name for which reference field will be created.
	// Ex: SecurityGroupIDs in ResourcesVpcConfig.SecurityGroupIDs
	fieldName := fsp.Pop()
	// "parentFieldName" is the Shape/Member name whose "TypeDef" contains the
	// "fieldName" as attribute. To add a corresponding reference for "fieldName"
	// , we will add new attribute in TypeDef for "parentFieldName".
	parentFieldName := fsp.Back()
	parentFieldShapeRef := fsp.ShapeRef(specFieldShapeRef)
	if parentFieldShapeRef == nil {
		return nil, nil, fmt.Errorf(
			"unable to find shape member %s for path %s",
			parentFieldName, fieldPath,
		)
	}
	parentFieldTypeDefName := parentFieldShapeRef.ShapeName

	var parentFieldTypeDef *TypeDef
	for _, td := range tdefs {
		fallbackName := ""
		switch parentFieldShapeRef.Shape.Type {
		case "list":
			// e.g FunctionAssociationsList in CloudFront DistributionConfig.DefaultCacheBehavior.FunctionAssociations
			fallbackName = parentFieldShapeRef.Shape.MemberRef.ShapeName
			fallbackName = strings.TrimSuffix(fallbackName, "List")
		default:
			// NOTE(a-hilaly): Very likely that we will need to add more cases here
			// as we encounter more special APIs in the future.
		}

		if strings.EqualFold(td.Names.Original, parentFieldTypeDefName) ||
			(fallbackName != "" && strings.EqualFold(td.Names.Original, fallbackName)) {
			parentFieldTypeDef = td
			break
		}
	}
	if parentFieldTypeDef == nil {
		return nil, nil, fmt.Errorf(
			"unable to find TypeDef %s in service model for path %s",
			parentFieldTypeDefName, fieldPath,
		)
	}

	fieldAttr := parentFieldTypeDef.GetAttribute(fieldName)
	if fieldAttr == nil {
		return nil, nil, fmt.Errorf(
			"unable to find member %s in TypeDef %s for path %s",
			fieldName, parentFieldTypeDefName, fieldPath,
		)
	}
	return parentFieldTypeDef, fieldAttr, nil
}

// setTypeDefAttributeGoTag sets the GoTag for the corresponding attribute
// represented by fieldPath of nested field.
func setTypeDefAttributeGoTag(crd *CRD, fieldPath string, f *Field, tdefs []*TypeDef) error {
	_, fieldAttr, err := getAttributeFromPath(crd, fieldPath, tdefs)
	if err != nil {
		return err
	}
	if fieldAttr != nil {
		fieldAttr.GoTag = f.GetGoTag()
	}
	return nil
}

// setTypeDefAttributeImmutable sets the IsImmutable flag for the corresponding
// attribute represented by fieldPath of nested field.
func setTypeDefAttributeImmutable(crd *CRD, fieldPath string, tdefs []*TypeDef) error {
	_, fieldAttr, err := getAttributeFromPath(crd, fieldPath, tdefs)
	if err != nil {
		return err
	}
	if fieldAttr != nil {
		fieldAttr.IsImmutable = true
	}
	return nil
}

// updateTypeDefAttributeWithReference adds a new AWSResourceReference attribute
// for the corresponding attribute represented by fieldPath of nested field
func updateTypeDefAttributeWithReference(crd *CRD, fieldPath string, tdefs []*TypeDef) error {
	parentFieldTypeDef, fieldAttr, err := getAttributeFromPath(crd, fieldPath, tdefs)
	if err != nil {
		return err
	}
	if fieldAttr != nil && parentFieldTypeDef != nil {
		if err := addReferenceAttribute(parentFieldTypeDef, fieldAttr); err != nil {
			return err
		}
	}
	return nil
}

// addReferenceAttribute creates a corresponding reference attribute for
// "attr" attribute and adds it to "td" TypeDef
func addReferenceAttribute(td *TypeDef, attr *Attr) error {
	// Create a custom "model.Field" to generate ReferenceFieldName and reuse
	// the existing method for generating top-level reference fields
	fieldShapeRef := awssdkmodel.ShapeRef{Shape: attr.Shape}
	field := &Field{
		Names:    attr.Names,
		ShapeRef: &fieldShapeRef,
	}
	refAttrName, err := field.GetReferenceFieldName()
	if err != nil {
		return err
	}
	refAttrShape := &awssdkmodel.Shape{
		Documentation: "// Reference field for " + attr.Names.Camel,
	}
	refAttrGoType := "*ackv1alpha1.AWSResourceReferenceWrapper"
	if attr.Shape.Type == "list" {
		refAttrGoType = fmt.Sprintf("[]%s", refAttrGoType)
	}
	refAttr := NewAttr(refAttrName, refAttrGoType, refAttrShape)
	// Add reference attribute to the parent field typedef
	td.Attrs[refAttrName.Original] = refAttr
	return nil
}

// replaceSecretAttrGoType replaces a nested field Attr's GoType with
// `*ackv1alpha1.SecretKeyReference`.
func replaceSecretAttrGoType(
	crd *CRD,
	field *Field,
	tdefs []*TypeDef,
) error {
	fieldPath := ackfp.FromString(field.Path)
	parentFieldPath := fieldPath.Copy()
	parentFieldPath.Pop()
	parentField, ok := crd.Fields[parentFieldPath.String()]
	if !ok {
		return fmt.Errorf(
			"cannot find parent field at path %s for %s",
			parentFieldPath, fieldPath,
		)
	}
	if parentField.ShapeRef == nil {
		return fmt.Errorf(
			"parent field at path %s has a nil ShapeRef",
			parentFieldPath,
		)
	}
	parentFieldShape := parentField.ShapeRef.Shape
	parentFieldShapeName := parentField.ShapeRef.ShapeName
	parentFieldShapeType := parentFieldShape.Type
	// For list and map types, we need to grab the element/value
	// type, since that's the type def we need to modify.
	if parentFieldShapeType == "list" {
		if parentFieldShape.MemberRef.Shape.Type != "structure" {
			return fmt.Errorf(
				"parent field at path %s is a list with non-structure element type %s",
				parentFieldPath, parentFieldShape.MemberRef.Shape.Type,
			)
		}
		parentFieldShapeName = parentField.ShapeRef.Shape.MemberRef.ShapeName
	} else if parentFieldShapeType == "map" {
		if parentFieldShape.ValueRef.Shape.Type != "structure" {
			return fmt.Errorf(
				"parent field at path %s is a map with non-structure value type %s",
				parentFieldPath, parentFieldShape.ValueRef.Shape.Type,
			)
		}
		parentFieldShapeName = parentField.ShapeRef.Shape.ValueRef.ShapeName
	}
	var parentTypeDef *TypeDef
	for _, tdef := range tdefs {
		if tdef.Names.Original == parentFieldShapeName {
			parentTypeDef = tdef
		}
	}
	if parentTypeDef == nil {
		return fmt.Errorf(
			"unable to find TypeDef for parent field at path %s (shape %s)",
			parentFieldPath, parentFieldShapeName,
		)
	}
	// Now we modify the parent type def's Attr that corresponds to
	// the secret field...
	attr, found := parentTypeDef.Attrs[field.Names.Camel]
	if !found {
		return fmt.Errorf(
			"unable to find attr %s in TypeDef %s at path %s",
			field.Names.Camel, parentTypeDef.Names.Original, parentFieldPath,
		)
	}
	attr.GoType = field.GoType
	return nil
}

// processFields is responsible for walking all of the CRDs' Spec and
// Status fields' Shape objects and adding `pkg/model.Field` objects for all
// nested fields along with that `Field`'s `Config` object that allows us to
// determine if the TypeDef associated with that nested field should have its
// data type overridden (e.g. for SecretKeyReferences)
func (m *Model) processFields(crds []*CRD) error {
	for _, crd := range crds {
		for _, field := range crd.SpecFields {
			if err := m.processTopLevelField(crd, field); err != nil {
				return err
			}
		}
		for _, field := range crd.StatusFields {
			if err := m.processTopLevelField(crd, field); err != nil {
				return err
			}
		}
	}
	return nil
}

// processTopLevelField processes any nested fields (non-scalar fields associated
// with the Spec and Status objects)
func (m *Model) processTopLevelField(
	crd *CRD,
	field *Field,
) error {
	if field.ShapeRef == nil && !field.IsReference() && (field.FieldConfig == nil || !field.FieldConfig.IsAttribute) {
		fmt.Printf(
			"WARNING: Field %s:%s has nil ShapeRef and is not defined as an Attribute-based Field!\n",
			crd.Names.Original,
			field.Names.Original,
		)
		return nil
	}
	if field.ShapeRef != nil {
		fieldShape := field.ShapeRef.Shape
		fieldType := fieldShape.Type
		switch fieldType {
		case "structure":
			if err := m.processStructField(crd, field.Path+".", field); err != nil {
				return err
			}
		case "list":
			if err := m.processListField(crd, field.Path+".", field); err != nil {
				return err
			}
		case "map":
			if err := m.processMapField(crd, field.Path+".", field); err != nil {
				return err
			}
		}
	}
	return nil
}

// processField adds a new Field definition for a field within the CR
func (m *Model) processField(
	crd *CRD,
	parentFieldPath string,
	parentField *Field,
	fieldName string,
	fieldShapeRef *awssdkmodel.ShapeRef,
) error {
	fieldNames := names.New(fieldName)
	fieldShape := fieldShapeRef.Shape
	fieldShapeType := fieldShape.Type
	fieldPath := parentFieldPath + fieldNames.Camel
	fieldConfig := crd.Config().GetFieldConfigByPath(crd.Names.Original, fieldPath)
	field, err := NewField(crd, fieldPath, fieldNames, fieldShapeRef, fieldConfig)
	if err != nil {
		return fmt.Errorf("resource %q, field %q: %w", crd.Names.Original, fieldPath, err)
	}
	switch fieldShapeType {
	case "structure":
		if err := m.processStructField(crd, fieldPath+".", field); err != nil {
			return err
		}
	case "list":
		if err := m.processListField(crd, fieldPath+".", field); err != nil {
			return err
		}
	case "map":
		if err := m.processMapField(crd, fieldPath+".", field); err != nil {
			return err
		}
	}
	crd.Fields[fieldPath] = field
	return nil
}

// processStructField recurses through the members of a nested field that
// is a struct type and adds any Field objects to the supplied CRD.
func (m *Model) processStructField(
	crd *CRD,
	fieldPath string,
	field *Field,
) error {
	fieldShape := field.ShapeRef.Shape
	for memberName, memberRef := range fieldShape.MemberRefs {
		if err := m.processField(crd, fieldPath, field, memberName, memberRef); err != nil {
			return err
		}
	}
	return nil
}

// processListField recurses through the members of a nested field that
// is a list type that has a struct element type and adds any Field objects to
// the supplied CRD.
func (m *Model) processListField(
	crd *CRD,
	fieldPath string,
	field *Field,
) error {
	fieldShape := field.ShapeRef.Shape
	elementFieldShape := fieldShape.MemberRef.Shape
	if elementFieldShape.Type != "structure" {
		return nil
	}
	for memberName, memberRef := range elementFieldShape.MemberRefs {
		if err := m.processField(crd, fieldPath, field, memberName, memberRef); err != nil {
			return err
		}
	}
	return nil
}

// processMapField recurses through the members of a nested field that
// is a map type that has a struct value type and adds any Field objects to
// the supplied CRD.
func (m *Model) processMapField(
	crd *CRD,
	fieldPath string,
	field *Field,
) error {
	fieldShape := field.ShapeRef.Shape
	valueFieldShape := fieldShape.ValueRef.Shape
	if valueFieldShape.Type != "structure" {
		return nil
	}
	for memberName, memberRef := range valueFieldShape.MemberRefs {
		if err := m.processField(crd, fieldPath, field, memberName, memberRef); err != nil {
			return err
		}
	}
	return nil
}

// GetEnumDefs returns a slice of pointers to `EnumDef` structs which
// represent string fields whose value is constrained to one or more specific
// string values.
func (m *Model) GetEnumDefs() ([]*EnumDef, error) {
	edefs := []*EnumDef{}

	for shapeName, shape := range m.SDKAPI.API.Shapes {
		if !shape.IsEnum() {
			continue
		}
		enumNames := names.New(shapeName)
		// Handle name conflicts with top-level CRD.Spec or CRD.Status
		// types
		if m.SDKAPI.HasConflictingTypeName(shapeName, m.cfg) {
			enumNames.Camel += ConflictingNameSuffix
		}
		edef, err := NewEnumDef(enumNames, shape.Enum)
		if err != nil {
			return nil, err
		}
		edefs = append(edefs, edef)
	}
	sort.Slice(edefs, func(i, j int) bool {
		return edefs[i].Names.Camel < edefs[j].Names.Camel
	})
	return edefs, nil
}

// ApplyShapeIgnoreRules removes the ignored shapes and fields from the API object
// so that they are not considered in any of the calculations of code generator.
func (m *Model) ApplyShapeIgnoreRules() error {
	if m.cfg == nil || m.SDKAPI == nil {
		return nil
	}
	for sdkShapeID, shape := range m.SDKAPI.API.Shapes {
		for _, sn := range m.cfg.Ignore.ShapeNames {
			if shape.ShapeName == sn {
				delete(m.SDKAPI.API.Shapes, sdkShapeID)
				continue
			}
			// NOTE(muvaf): We need to remove the usage of the shape as well.
			for sdkMemberID, memberRef := range shape.MemberRefs {
				if memberRef.ShapeName == sn {
					delete(shape.MemberRefs, sdkMemberID)
				}
			}
		}
	}
	for _, fieldpath := range m.cfg.Ignore.FieldPaths {
		fp := ackfp.FromString(fieldpath)
		sn := fp.At(0)
		if shape, found := m.SDKAPI.API.Shapes[sn]; !found {
			return fmt.Errorf(
				"ignore.field_paths refers to unknown shape %q — check generator.yaml", sn,
			)
		} else {
			// This is just some tomfoolery to make the Input and Output shapes
			// into ShapeRefs because the SDKAPI.Shapes is a map of Shape
			// pointers not a map of ShapeRefs...
			wrapper := &awssdkmodel.ShapeRef{
				ShapeName: sn,
				Shape:     shape,
			}
			// The last element of the fieldpath is the field/shape we want to
			// ignore...
			ignoreShape := fp.Pop()
			parentShapeRef := fp.ShapeRef(wrapper)
			// OK, now we delete the ignored shape by removing the shape from
			// the parent's member references...
			delete(parentShapeRef.Shape.MemberRefs, ignoreShape)
		}
	}
	return nil
}

// GetConfig returns the configuration option used to define the current
// generator.
func (m *Model) GetConfig() *ackgenconfig.Config {
	return m.cfg
}

// APIGroup returns the normalized Kubernetes APIGroup for the AWS service API,
// e.g. "sns.services.k8s.aws"
func (m *Model) APIGroup() string {
	suffix := "services.k8s.aws"
	if m.SDKAPI.APIGroupSuffix != "" {
		suffix = m.SDKAPI.APIGroupSuffix
	}
	name := m.GetConfig().ControllerName
	if name == "" {
		name = m.servicePackageName
	}
	return fmt.Sprintf("%s.%s", name, suffix)
}

// ClientInterfaceTypeName returns the name of the aws-sdk-go primary API
// interface type name.
func (m *Model) ClientInterfaceTypeName() string {
	if m.cfg.SDKNames.ClientInterface != "" {
		return m.cfg.SDKNames.ClientInterface
	}
	return m.SDKAPI.ClientInterfaceTypeName()
}

// ClientStructTypeName returns the name of the aws-sdk-go primary client
// struct type name.
func (m *Model) ClientStructTypeName() string {
	if m.cfg.SDKNames.ClientStruct != "" {
		return m.cfg.SDKNames.ClientStruct
	}
	return m.SDKAPI.ClientStructTypeName()
}

// New returns a new Model struct for a supplied API model.
// Optionally, pass a file path to a generator config file that can be used to
// instruct the code generator how to handle the API properly
func New(
	SDKAPI *SDKAPI,
	servicePackageName string,
	apiVersion string,
	cfg ackgenconfig.Config,
	docCfg ackgenconfig.DocumentationConfig,
) (*Model, error) {
	m := &Model{
		SDKAPI:             SDKAPI,
		servicePackageName: servicePackageName,
		apiVersion:         apiVersion,
		cfg:                &cfg,
		docCfg:             &docCfg,
	}
	if err := m.ApplyShapeIgnoreRules(); err != nil {
		return nil, err
	}
	return m, nil
}
