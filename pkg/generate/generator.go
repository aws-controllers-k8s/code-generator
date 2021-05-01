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

package generate

import (
	"fmt"
	"sort"
	"strings"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/names"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

// Generator contains the ACK model for the generator to process and apply
// templates against
type Generator struct {
	SDKAPI       *ackmodel.SDKAPI
	serviceAlias string
	apiVersion   string
	crds         []*ackmodel.CRD
	typeDefs     []*ackmodel.TypeDef
	typeImports  map[string]string
	typeRenames  map[string]string
	// Instructions to the code generator how to handle the API and its
	// resources
	cfg *ackgenconfig.Config
}

// MetaVars returns a MetaVars struct populated with metadata about the AWS
// service API
func (g *Generator) MetaVars() templateset.MetaVars {
	return templateset.MetaVars{
		ServiceAlias:            g.serviceAlias,
		ServiceID:               g.SDKAPI.ServiceID(),
		ServiceIDClean:          g.SDKAPI.ServiceIDClean(),
		APIGroup:                g.SDKAPI.APIGroup(),
		APIVersion:              g.apiVersion,
		SDKAPIInterfaceTypeName: g.SDKAPI.SDKAPIInterfaceTypeName(),
		CRDNames:                g.crdNames(),
	}
}

// crdNames returns all crd names lowercased and in plural
func (g *Generator) crdNames() []string {
	var crdConfigs []string

	crds, _ := g.GetCRDs()
	for _, crd := range crds {
		crdConfigs = append(crdConfigs, strings.ToLower(crd.Plural))
	}

	return crdConfigs
}

// GetCRDs returns a slice of `ackmodel.CRD` structs that describe the
// top-level resources discovered by the code generator for an AWS service API
func (g *Generator) GetCRDs() ([]*ackmodel.CRD, error) {
	if g.crds != nil {
		return g.crds, nil
	}
	crds := []*ackmodel.CRD{}

	opMap := g.SDKAPI.GetOperationMap(g.cfg)

	createOps := (*opMap)[ackmodel.OpTypeCreate]
	readOneOps := (*opMap)[ackmodel.OpTypeGet]
	readManyOps := (*opMap)[ackmodel.OpTypeList]
	updateOps := (*opMap)[ackmodel.OpTypeUpdate]
	deleteOps := (*opMap)[ackmodel.OpTypeDelete]
	getAttributesOps := (*opMap)[ackmodel.OpTypeGetAttributes]
	setAttributesOps := (*opMap)[ackmodel.OpTypeSetAttributes]

	for crdName, createOp := range createOps {
		if g.cfg.IsIgnoredResource(crdName) {
			continue
		}
		crdNames := names.New(crdName)
		ops := ackmodel.Ops{
			Create:        createOps[crdName],
			ReadOne:       readOneOps[crdName],
			ReadMany:      readManyOps[crdName],
			Update:        updateOps[crdName],
			Delete:        deleteOps[crdName],
			GetAttributes: getAttributesOps[crdName],
			SetAttributes: setAttributesOps[crdName],
		}
		g.RemoveIgnoredOperations(&ops)
		crd := ackmodel.NewCRD(g.SDKAPI, g.cfg, crdNames, ops)

		// OK, begin to gather the CRDFields that will go into the Spec struct.
		// These fields are those members of the Create operation's Input
		// Shape.
		inputShape := createOp.InputRef.Shape
		if inputShape == nil {
			return nil, ErrNilShapePointer
		}
		for memberName, memberShapeRef := range inputShape.MemberRefs {
			if memberShapeRef.Shape == nil {
				return nil, ErrNilShapePointer
			}
			renamedName, _ := crd.InputFieldRename(
				createOp.Name, memberName,
			)
			memberNames := names.New(renamedName)
			memberNames.ModelOriginal = memberName
			if memberName == "Attributes" && g.cfg.UnpacksAttributesMap(crdName) {
				crd.UnpackAttributes()
				continue
			}
			crd.AddSpecField(memberNames, memberShapeRef)
		}

		// Now any additional Spec fields that are required from other API
		// operations.
		for targetFieldName, fieldConfig := range g.cfg.ResourceFields(crdName) {
			if fieldConfig.IsReadOnly {
				// It's a Status field...
				continue
			}
			if fieldConfig.From == nil {
				// Isn't an additional Spec field...
				continue
			}
			from := fieldConfig.From
			memberShapeRef, found := g.SDKAPI.GetInputShapeRef(
				from.Operation, from.Path,
			)
			if found {
				memberNames := names.New(targetFieldName)
				crd.AddSpecField(memberNames, memberShapeRef)
			} else {
				// This is a compile-time failure, just bomb out...
				msg := fmt.Sprintf(
					"unknown additional Spec field with Op: %s and Path: %s",
					from.Operation, from.Path,
				)
				panic(msg)
			}
		}

		// Now process the fields that will go into the Status struct. We want
		// fields that are in the Create operation's Output Shape but that are
		// not in the Input Shape.
		outputShape := createOp.OutputRef.Shape
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
			// fields in the input shape (where the input shape has potentially
			// been renamed)
			renamedName, _ := crd.InputFieldRename(
				createOp.Name, memberName,
			)
			memberNames := names.New(renamedName)
			if _, found := crd.SpecFields[renamedName]; found {
				// We don't put fields that are already in the Spec struct into
				// the Status struct
				continue
			}
			if memberName == "Attributes" && g.cfg.UnpacksAttributesMap(crdName) {
				continue
			}
			if crd.IsPrimaryARNField(memberName) {
				// We automatically place the primary resource ARN value into
				// the Status.ACKResourceMetadata.ARN field
				continue
			}
			crd.AddStatusField(memberNames, memberShapeRef)
		}

		// Now add the additional Status fields that are required from other
		// API operations.
		for targetFieldName, fieldConfig := range g.cfg.ResourceFields(crdName) {
			if !fieldConfig.IsReadOnly {
				// It's a Spec field...
				continue
			}
			if fieldConfig.From == nil {
				// Isn't an additional Status field...
				continue
			}
			from := fieldConfig.From
			memberShapeRef, found := g.SDKAPI.GetOutputShapeRef(
				from.Operation, from.Path,
			)
			if found {
				memberNames := names.New(targetFieldName)
				crd.AddStatusField(memberNames, memberShapeRef)
			} else {
				// This is a compile-time failure, just bomb out...
				msg := fmt.Sprintf(
					"unknown additional Status field with Op: %s and Path: %s",
					from.Operation, from.Path,
				)
				panic(msg)
			}
		}

		crds = append(crds, crd)
	}
	sort.Slice(crds, func(i, j int) bool {
		return crds[i].Names.Camel < crds[j].Names.Camel
	})
	// This is the place that we build out the CRD.Fields map with
	// `pkg/model.Field` objects that represent the non-top-level Spec and
	// Status fields.
	g.processNestedFields(crds)
	g.crds = crds
	return crds, nil
}

// RemoveIgnoredOperations updates Ops argument by setting those
// operations to nil that are configured to be ignored in generator config for
// the AWS service
func (g *Generator) RemoveIgnoredOperations(ops *ackmodel.Ops) {
	if g.cfg.IsIgnoredOperation(ops.Create) {
		ops.Create = nil
	}
	if g.cfg.IsIgnoredOperation(ops.ReadOne) {
		ops.ReadOne = nil
	}
	if g.cfg.IsIgnoredOperation(ops.ReadMany) {
		ops.ReadMany = nil
	}
	if g.cfg.IsIgnoredOperation(ops.Update) {
		ops.Update = nil
	}
	if g.cfg.IsIgnoredOperation(ops.Delete) {
		ops.Delete = nil
	}
	if g.cfg.IsIgnoredOperation(ops.GetAttributes) {
		ops.GetAttributes = nil
	}
	if g.cfg.IsIgnoredOperation(ops.SetAttributes) {
		ops.SetAttributes = nil
	}
}

// IsShapeUsedInCRDs returns true if the supplied shape name is a member of amy
// CRD's payloads or those payloads sub-member shapes
func (g *Generator) IsShapeUsedInCRDs(shapeName string) bool {
	crds, _ := g.GetCRDs()
	for _, crd := range crds {
		if crd.HasShapeAsMember(shapeName) {
			return true
		}
	}
	return false
}

// GetTypeDefs returns a slice of `ackmodel.TypeDef` pointers
func (g *Generator) GetTypeDefs() ([]*ackmodel.TypeDef, error) {
	if g.typeDefs != nil {
		return g.typeDefs, nil
	}

	tdefs := []*ackmodel.TypeDef{}
	// Map, keyed by original Shape GoTypeElem(), with the values being a
	// renamed type name (due to conflicting names)
	trenames := map[string]string{}

	payloads := g.SDKAPI.GetPayloads()

	for shapeName, shape := range g.SDKAPI.API.Shapes {
		if util.InStrings(shapeName, payloads) {
			// Payloads are not type defs
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
		if g.SDKAPI.HasConflictingTypeName(shapeName, g.cfg) {
			tdefNames.Camel += ackmodel.ConflictingNameSuffix
			trenames[shapeName] = tdefNames.Camel
		}

		attrs := map[string]*ackmodel.Attr{}
		for memberName, memberRef := range shape.MemberRefs {
			memberNames := names.New(memberName)
			memberShape := memberRef.Shape
			if !g.IsShapeUsedInCRDs(memberShape.ShapeName) {
				continue
			}
			// There are shapes that are called things like DBProxyStatus that are
			// fields in a DBProxy CRD... we need to ensure the type names don't
			// conflict. Also, the name of the Go type in the generated code is
			// Camel-cased and normalized, so we use that as the Go type
			gt := memberShape.GoType()
			if memberShape.Type == "structure" {
				typeNames := names.New(memberShape.ShapeName)
				if g.SDKAPI.HasConflictingTypeName(memberShape.ShapeName, g.cfg) {
					typeNames.Camel += ackmodel.ConflictingNameSuffix
				}
				gt = "*" + typeNames.Camel
			} else if memberShape.Type == "list" {
				// If it's a list type, where the element is a structure, we need to
				// set the GoType to the cleaned-up Camel-cased name
				if memberShape.MemberRef.Shape.Type == "structure" {
					elemType := memberShape.MemberRef.Shape.GoTypeElem()
					typeNames := names.New(elemType)
					if g.SDKAPI.HasConflictingTypeName(elemType, g.cfg) {
						typeNames.Camel += ackmodel.ConflictingNameSuffix
					}
					gt = "[]*" + typeNames.Camel
				}
			} else if memberShape.Type == "map" {
				// If it's a map type, where the value element is a structure,
				// we need to set the GoType to the cleaned-up Camel-cased name
				if memberShape.ValueRef.Shape.Type == "structure" {
					valType := memberShape.ValueRef.Shape.GoTypeElem()
					typeNames := names.New(valType)
					if g.SDKAPI.HasConflictingTypeName(valType, g.cfg) {
						typeNames.Camel += ackmodel.ConflictingNameSuffix
					}
					gt = "[]map[string]*" + typeNames.Camel
				}
			} else if memberShape.Type == "timestamp" {
				// time.Time needs to be converted to apimachinery/metav1.Time
				// otherwise there is no DeepCopy support
				gt = "*metav1.Time"
			}
			attrs[memberName] = ackmodel.NewAttr(memberNames, gt, memberShape)
		}
		if len(attrs) == 0 {
			// Just ignore these...
			continue
		}
		tdefs = append(tdefs, &ackmodel.TypeDef{
			Names: tdefNames,
			Attrs: attrs,
		})
	}
	sort.Slice(tdefs, func(i, j int) bool {
		return tdefs[i].Names.Camel < tdefs[j].Names.Camel
	})
	g.processNestedFieldTypeDefs(tdefs)
	g.typeDefs = tdefs
	g.typeRenames = trenames
	return tdefs, nil
}

// processNestedFieldTypeDefs updates the supplied TypeDef structs' if a nested
// field has been configured with a type overriding FieldConfig -- such as
// FieldConfig.IsSecret.
func (g *Generator) processNestedFieldTypeDefs(
	tdefs []*ackmodel.TypeDef,
) {
	crds, _ := g.GetCRDs()
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
				replaceSecretAttrGoType(crd, field, tdefs)
			}
		}
	}
}

// replaceSecretAttrGoType replaces a nested field ackmodel.Attr's GoType with
// `*ackv1alpha1.SecretKeyReference`.
func replaceSecretAttrGoType(
	crd *ackmodel.CRD,
	field *ackmodel.Field,
	tdefs []*ackmodel.TypeDef,
) {
	fieldPath := field.Path
	parentFieldPath := ackmodel.ParentFieldPath(field.Path)
	parentField, ok := crd.Fields[parentFieldPath]
	if !ok {
		msg := fmt.Sprintf(
			"Cannot find parent field at parent path %s for %s",
			parentFieldPath,
			fieldPath,
		)
		panic(msg)
	}
	if parentField.ShapeRef == nil {
		msg := fmt.Sprintf(
			"parent field at parent path %s has a nil ShapeRef!",
			parentFieldPath,
		)
		panic(msg)
	}
	parentFieldShape := parentField.ShapeRef.Shape
	parentFieldShapeName := parentField.ShapeRef.ShapeName
	parentFieldShapeType := parentFieldShape.Type
	// For list and map types, we need to grab the element/value
	// type, since that's the type def we need to modify.
	if parentFieldShapeType == "list" {
		if parentFieldShape.MemberRef.Shape.Type != "structure" {
			msg := fmt.Sprintf(
				"parent field at parent path %s is a list type with a non-structure element member shape %s!",
				parentFieldPath,
				parentFieldShape.MemberRef.Shape.Type,
			)
			panic(msg)
		}
		parentFieldShapeName = parentField.ShapeRef.Shape.MemberRef.ShapeName
	} else if parentFieldShapeType == "map" {
		if parentFieldShape.ValueRef.Shape.Type != "structure" {
			msg := fmt.Sprintf(
				"parent field at parent path %s is a map type with a non-structure value member shape %s!",
				parentFieldPath,
				parentFieldShape.ValueRef.Shape.Type,
			)
			panic(msg)
		}
		parentFieldShapeName = parentField.ShapeRef.Shape.ValueRef.ShapeName
	}
	var parentTypeDef *ackmodel.TypeDef
	for _, tdef := range tdefs {
		if tdef.Names.Original == parentFieldShapeName {
			parentTypeDef = tdef
		}
	}
	if parentTypeDef == nil {
		msg := fmt.Sprintf(
			"unable to find associated TypeDef for parent field "+
				"at parent path %s!",
			parentFieldPath,
		)
		panic(msg)
	}
	// Now we modify the parent type def's Attr that corresponds to
	// the secret field...
	attr, found := parentTypeDef.Attrs[field.Names.Camel]
	if !found {
		msg := fmt.Sprintf(
			"unable to find attr %s in parent TypeDef %s "+
				"at parent path %s!",
			field.Names.Camel,
			parentTypeDef.Names.Original,
			parentFieldPath,
		)
		panic(msg)
	}
	attr.GoType = "*ackv1alpha1.SecretKeyReference"
}

// processNestedFields is responsible for walking all of the CRDs' Spec and
// Status fields' Shape objects and adding `pkg/model.Field` objects for all
// nested fields along with that `Field`'s `Config` object that allows us to
// determine if the TypeDef associated with that nested field should have its
// data type overridden (e.g. for SecretKeyReferences)
func (g *Generator) processNestedFields(crds []*ackmodel.CRD) {
	for _, crd := range crds {
		for _, field := range crd.SpecFields {
			g.processNestedField(crd, field)
		}
		for _, field := range crd.StatusFields {
			g.processNestedField(crd, field)
		}
	}
}

// processNestedField processes any nested fields (non-scalar fields associated
// with the Spec and Status objects)
func (g *Generator) processNestedField(
	crd *ackmodel.CRD,
	field *ackmodel.Field,
) {
	if field.ShapeRef == nil && (field.FieldConfig == nil || !field.FieldConfig.IsAttribute) {
		fmt.Printf(
			"WARNING: Field %s:%s has nil ShapeRef and is not defined as an Attribute-based Field!\n",
			crd.Names.Original,
			field.Names.Original,
		)
		return
	}
	if field.ShapeRef != nil {
		fieldShape := field.ShapeRef.Shape
		fieldType := fieldShape.Type
		switch fieldType {
		case "structure":
			g.processNestedStructField(crd, field.Path+".", field)
		case "list":
			g.processNestedListField(crd, field.Path+"..", field)
		case "map":
			g.processNestedMapField(crd, field.Path+"..", field)
		}
	}
}

// processNestedStructField recurses through the members of a nested field that
// is a struct type and adds any Field objects to the supplied CRD.
func (g *Generator) processNestedStructField(
	crd *ackmodel.CRD,
	baseFieldPath string,
	baseField *ackmodel.Field,
) {
	fieldConfigs := crd.Config().ResourceFields(crd.Names.Original)
	baseFieldShape := baseField.ShapeRef.Shape
	for memberName, memberRef := range baseFieldShape.MemberRefs {
		memberNames := names.New(memberName)
		memberShape := memberRef.Shape
		memberShapeType := memberShape.Type
		fieldPath := baseFieldPath + memberNames.Camel
		fieldConfig := fieldConfigs[fieldPath]
		field := ackmodel.NewField(crd, fieldPath, memberNames, memberRef, fieldConfig)
		switch memberShapeType {
		case "structure":
			g.processNestedStructField(crd, fieldPath+".", field)
		case "list":
			g.processNestedListField(crd, fieldPath+"..", field)
		case "map":
			g.processNestedMapField(crd, fieldPath+"..", field)
		}
		crd.Fields[fieldPath] = field
	}
}

// processNestedListField recurses through the members of a nested field that
// is a list type that has a struct element type and adds any Field objects to
// the supplied CRD.
func (g *Generator) processNestedListField(
	crd *ackmodel.CRD,
	baseFieldPath string,
	baseField *ackmodel.Field,
) {
	baseFieldShape := baseField.ShapeRef.Shape
	elementFieldShape := baseFieldShape.MemberRef.Shape
	if elementFieldShape.Type != "structure" {
		return
	}
	fieldConfigs := crd.Config().ResourceFields(crd.Names.Original)
	for memberName, memberRef := range elementFieldShape.MemberRefs {
		memberNames := names.New(memberName)
		memberShape := memberRef.Shape
		memberShapeType := memberShape.Type
		fieldPath := baseFieldPath + memberNames.Camel
		fieldConfig := fieldConfigs[fieldPath]
		field := ackmodel.NewField(crd, fieldPath, memberNames, memberRef, fieldConfig)
		switch memberShapeType {
		case "structure":
			g.processNestedStructField(crd, fieldPath+".", field)
		case "list":
			g.processNestedListField(crd, fieldPath+"..", field)
		case "map":
			g.processNestedMapField(crd, fieldPath+"..", field)
		}
		crd.Fields[fieldPath] = field
	}
}

// processNestedMapField recurses through the members of a nested field that
// is a map type that has a struct value type and adds any Field objects to
// the supplied CRD.
func (g *Generator) processNestedMapField(
	crd *ackmodel.CRD,
	baseFieldPath string,
	baseField *ackmodel.Field,
) {
	baseFieldShape := baseField.ShapeRef.Shape
	valueFieldShape := baseFieldShape.ValueRef.Shape
	if valueFieldShape.Type != "structure" {
		return
	}
	fieldConfigs := crd.Config().ResourceFields(crd.Names.Original)
	for memberName, memberRef := range valueFieldShape.MemberRefs {
		memberNames := names.New(memberName)
		memberShape := memberRef.Shape
		memberShapeType := memberShape.Type
		fieldPath := baseFieldPath + memberNames.Camel
		fieldConfig := fieldConfigs[fieldPath]
		field := ackmodel.NewField(crd, fieldPath, memberNames, memberRef, fieldConfig)
		switch memberShapeType {
		case "structure":
			g.processNestedStructField(crd, fieldPath+".", field)
		case "list":
			g.processNestedListField(crd, fieldPath+"..", field)
		case "map":
			g.processNestedMapField(crd, fieldPath+"..", field)
		}
		crd.Fields[fieldPath] = field
	}
}

// GetEnumDefs returns a slice of pointers to `ackmodel.EnumDef` structs which
// represent string fields whose value is constrained to one or more specific
// string values.
func (g *Generator) GetEnumDefs() ([]*ackmodel.EnumDef, error) {
	edefs := []*ackmodel.EnumDef{}

	for shapeName, shape := range g.SDKAPI.API.Shapes {
		if !shape.IsEnum() {
			continue
		}
		enumNames := names.New(shapeName)
		// Handle name conflicts with top-level CRD.Spec or CRD.Status
		// types
		if g.SDKAPI.HasConflictingTypeName(shapeName, g.cfg) {
			enumNames.Camel += ackmodel.ConflictingNameSuffix
		}
		edef, err := ackmodel.NewEnumDef(enumNames, shape.Enum)
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
func (g *Generator) ApplyShapeIgnoreRules() {
	if g.cfg == nil || g.SDKAPI == nil {
		return
	}
	for sdkShapeID, shape := range g.SDKAPI.API.Shapes {
		for _, fieldpath := range g.cfg.Ignore.FieldPaths {
			sn := strings.Split(fieldpath, ".")[0]
			fn := strings.Split(fieldpath, ".")[1]
			if shape.ShapeName != sn {
				continue
			}
			delete(shape.MemberRefs, fn)
		}
		for _, sn := range g.cfg.Ignore.ShapeNames {
			if shape.ShapeName == sn {
				delete(g.SDKAPI.API.Shapes, sdkShapeID)
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
}

// New returns a new Generator struct for a supplied API model.
// Optionally, pass a file path to a generator config file that can be used to
// instruct the code generator how to handle the API properly
func New(
	SDKAPI *ackmodel.SDKAPI,
	apiVersion string,
	configPath string,
	defaultConfig ackgenconfig.Config,
) (*Generator, error) {
	cfg, err := ackgenconfig.New(configPath, defaultConfig)
	if err != nil {
		return nil, err
	}
	g := &Generator{
		SDKAPI: SDKAPI,
		// TODO(jaypipes): Handle cases where service alias and service ID
		// don't match (Step Functions)
		serviceAlias: SDKAPI.ServiceID(),
		apiVersion:   apiVersion,
		cfg:          &cfg,
	}
	g.ApplyShapeIgnoreRules()
	return g, nil
}
