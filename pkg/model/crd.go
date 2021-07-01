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

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/gertd/go-pluralize"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/names"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

// Ops are the CRUD operations controlling a particular resource
type Ops struct {
	Create        *awssdkmodel.Operation
	ReadOne       *awssdkmodel.Operation
	ReadMany      *awssdkmodel.Operation
	Update        *awssdkmodel.Operation
	Delete        *awssdkmodel.Operation
	GetAttributes *awssdkmodel.Operation
	SetAttributes *awssdkmodel.Operation
}

// IterOps returns a slice of Operations for a resource
func (ops Ops) IterOps() []*awssdkmodel.Operation {
	res := []*awssdkmodel.Operation{}
	if ops.Create != nil {
		res = append(res, ops.Create)
	}
	if ops.ReadOne != nil {
		res = append(res, ops.ReadOne)
	}
	if ops.ReadMany != nil {
		res = append(res, ops.ReadMany)
	}
	if ops.Update != nil {
		res = append(res, ops.Update)
	}
	if ops.Delete != nil {
		res = append(res, ops.Delete)
	}
	return res
}

// CRD describes a single top-level resource in an AWS service API
type CRD struct {
	sdkAPI *SDKAPI
	cfg    *ackgenconfig.Config
	Names  names.Names
	Kind   string
	Plural string
	// Ops are the CRUD operations controlling this resource
	Ops Ops
	// additionalPrinterColumns is an array of PrinterColumn objects
	// representing the printer column settings for the CRD
	additionalPrinterColumns []*PrinterColumn
	// SpecFields is a map, keyed by the **original SDK member name** of
	// Field objects representing those fields in the CRD's Spec struct
	// field.
	SpecFields map[string]*Field
	// StatusFields is a map, keyed by the **original SDK member name** of
	// Field objects representing those fields in the CRD's Status struct
	// field. Note that there are no fields in StatusFields that are also in
	// SpecFields.
	StatusFields map[string]*Field
	// Fields is a map, keyed by the **renamed/normalized field path**, of
	// Field objects representing a field in the CRD's Spec or Status objects.
	Fields map[string]*Field
	// TypeImports is a map, keyed by an import string, with the map value
	// being the import alias
	TypeImports map[string]string
	// ShortNames represent the CRD list of aliases. Short names allow shorter
	// strings to match a CR on the CLI.
	ShortNames []string
}

// Config returns a pointer to the generator config
func (r *CRD) Config() *ackgenconfig.Config {
	return r.cfg
}

// SDKAPIPackageName returns the aws-sdk-go package name used for this
// resource's API
func (r *CRD) SDKAPIPackageName() string {
	return r.sdkAPI.API.PackageName()
}

// TypeRenames returns a map of original type name to renamed name (some
// type definition names conflict with generated names)
func (r *CRD) TypeRenames() map[string]string {
	return r.sdkAPI.GetTypeRenames(r.cfg)
}

// Documentation returns the base documentation string for the API formatted as
// a Go code comment block
func (r *CRD) Documentation() string {
	docString := fmt.Sprintf("// %sSpec defines the desired state of %s.", r.Names.Original, r.Names.Original)
	shape, ok := r.sdkAPI.API.Shapes[r.Names.Original]
	if ok {
		// Separate with a double newline to force a newline in the CRD base
		docString += "\n//\n" + shape.Documentation
	}
	return docString
}

// HasShapeAsMember returns true if the supplied Shape name appears in *any*
// payload shape of *any* Operation for the resource. It recurses down through
// the resource's Operation Input and Output shapes and their member shapes
// looking for a shape with the supplied name
func (r *CRD) HasShapeAsMember(toFind string) bool {
	for _, op := range r.Ops.IterOps() {
		if op.InputRef.Shape != nil {
			inShape := op.InputRef.Shape
			for _, memberShapeRef := range inShape.MemberRefs {
				if shapeHasMember(memberShapeRef.Shape, toFind) {
					return true
				}
			}
		}
		if op.OutputRef.Shape != nil {
			outShape := op.OutputRef.Shape
			for _, memberShapeRef := range outShape.MemberRefs {
				if shapeHasMember(memberShapeRef.Shape, toFind) {
					return true
				}
			}
		}
	}
	return false
}

func shapeHasMember(shape *awssdkmodel.Shape, toFind string) bool {
	if shape.ShapeName == toFind {
		return true
	}
	switch shape.Type {
	case "structure":
		for _, memberShapeRef := range shape.MemberRefs {
			if shapeHasMember(memberShapeRef.Shape, toFind) {
				return true
			}
		}
	case "list":
		return shapeHasMember(shape.MemberRef.Shape, toFind)
	case "map":
		return shapeHasMember(shape.ValueRef.Shape, toFind)
	}
	return false
}

// InputFieldRename returns the renamed field for a supplied Operation ID and
// original field name and whether or not a renamed override field name was
// found
func (r *CRD) InputFieldRename(
	opID string,
	origFieldName string,
) (string, bool) {
	if r.cfg == nil {
		return origFieldName, false
	}
	return r.cfg.ResourceInputFieldRename(
		r.Names.Original, opID, origFieldName,
	)
}

// AddSpecField adds a new Field of a given name and shape into the Spec
// field of a CRD
func (r *CRD) AddSpecField(
	memberNames names.Names,
	shapeRef *awssdkmodel.ShapeRef,
) {
	fPath := memberNames.Camel
	fConfigs := r.cfg.ResourceFields(r.Names.Original)
	fConfig := fConfigs[memberNames.Original]
	f := NewField(r, fPath, memberNames, shapeRef, fConfig)
	if fConfig != nil && fConfig.Print != nil {
		r.addSpecPrintableColumn(f)
	}
	r.SpecFields[memberNames.Original] = f
	r.Fields[fPath] = f
}

// AddStatusField adds a new Field of a given name and shape into the Status
// field of a CRD
func (r *CRD) AddStatusField(
	memberNames names.Names,
	shapeRef *awssdkmodel.ShapeRef,
) {
	fPath := memberNames.Camel
	fConfigs := r.cfg.ResourceFields(r.Names.Original)
	fConfig := fConfigs[memberNames.Original]
	f := NewField(r, fPath, memberNames, shapeRef, fConfig)
	if fConfig != nil && fConfig.Print != nil {
		r.addStatusPrintableColumn(f)
	}
	r.StatusFields[memberNames.Original] = f
	r.Fields[fPath] = f
}

// AddTypeImport adds an entry in the CRD's TypeImports map for an import line
// and optional alias
func (r *CRD) AddTypeImport(
	packagePath string,
	alias string,
) {
	if r.TypeImports == nil {
		r.TypeImports = map[string]string{}
	}
	r.TypeImports[packagePath] = alias
}

// SpecFieldNames returns a sorted slice of field names for the Spec fields
func (r *CRD) SpecFieldNames() []string {
	res := make([]string, 0, len(r.SpecFields))
	for fieldName := range r.SpecFields {
		res = append(res, fieldName)
	}
	sort.Strings(res)
	return res
}

// UnpacksAttributesMap returns true if the underlying API has
// Get{Resource}Attributes/Set{Resource}Attributes API calls that map real,
// schema'd fields to a raw `map[string]*string` for this resource (see SNS and
// SQS APIs)
func (r *CRD) UnpacksAttributesMap() bool {
	return r.cfg.UnpacksAttributesMap(r.Names.Original)
}

// CompareIgnoredFields returns the list of fields compare logic should ignore
func (r *CRD) CompareIgnoredFields() []string {
	return r.cfg.GetCompareIgnoredFields(r.Names.Original)
}

// SetAttributesSingleAttribute returns true if the supplied resource name has
// a SetAttributes operation that only actually changes a single attribute at a
// time. See: SNS SetTopicAttributes API call, which is entirely different from
// the SNS SetPlatformApplicationAttributes API call, which sets multiple
// attributes at once. :shrug:
func (r *CRD) SetAttributesSingleAttribute() bool {
	return r.cfg.SetAttributesSingleAttribute(r.Names.Original)
}

// UnpackAttributes grabs instructions about fields that are represented in the
// AWS API as a `map[string]*string` but are actually real, schema'd fields and
// adds Field definitions for those fields.
func (r *CRD) UnpackAttributes() {
	if !r.cfg.UnpacksAttributesMap(r.Names.Original) {
		return
	}
	fieldConfigs := r.cfg.ResourceFields(r.Names.Original)
	for fieldName, fieldConfig := range fieldConfigs {
		if !fieldConfig.IsAttribute {
			continue
		}
		if r.IsPrimaryARNField(fieldName) {
			// ignore since this is handled by Status.ACKResourceMetadata.ARN
			continue
		}
		fieldNames := names.New(fieldName)
		fPath := fieldNames.Camel

		f := NewField(r, fPath, fieldNames, nil, fieldConfig)
		if !fieldConfig.IsReadOnly {
			r.SpecFields[fieldName] = f
		} else {
			r.StatusFields[fieldName] = f
		}
		r.Fields[fPath] = f
	}
}

// IsPrimaryARNField returns true if the supplied field name is likely the resource's
// ARN identifier field.
func (r *CRD) IsPrimaryARNField(fieldName string) bool {
	if r.cfg != nil && !r.cfg.IncludeACKMetadata {
		return false
	}
	rConfig, found := r.cfg.Resources[r.Names.Original]
	if found {
		for fName, fConfig := range rConfig.Fields {
			if fConfig.IsARN {
				return strings.EqualFold(fieldName, fName)
			}
		}
	}

	return strings.EqualFold(fieldName, "arn") ||
		strings.EqualFold(fieldName, r.Names.Original+"arn")
}

// IsSecretField returns true if the supplied field *path* refers to a Field
// that is a SecretKeyReference
func (r *CRD) IsSecretField(path string) bool {
	fConfigs := r.cfg.ResourceFields(r.Names.Original)
	fConfig, found := fConfigs[path]
	if found {
		return fConfig.IsSecret
	}
	return false
}

// GetImmutableFieldPaths returns list of immutable field paths present in CRD
func (r *CRD) GetImmutableFieldPaths() []string {
	fConfigs := r.cfg.ResourceFields(r.Names.Original)
	var immutableFields []string

	for field, fieldConfig := range fConfigs {
		if fieldConfig.IsImmutable {
			immutableFields = append(immutableFields, field)
		}
	}

	return immutableFields
}

// HasImmutableFieldChanges helper function that return true if there are any immutable field changes
func (r *CRD) HasImmutableFieldChanges() bool {
	fConfigs := r.cfg.ResourceFields(r.Names.Original)

	for _, fieldConfig := range fConfigs {
		if fieldConfig.IsImmutable {
			return true
		}
	}
	return false
}

// SetOutputCustomMethodName returns custom set output operation as *string for
// given operation on custom resource, if specified in generator config
func (r *CRD) SetOutputCustomMethodName(
	// The operation to look for the Output shape
	op *awssdkmodel.Operation,
) *string {
	if op == nil {
		return nil
	}
	if r.cfg == nil {
		return nil
	}
	resGenConfig, found := r.cfg.Operations[op.Name]
	if !found {
		return nil
	}

	if resGenConfig.SetOutputCustomMethodName == "" {
		return nil
	}
	return &resGenConfig.SetOutputCustomMethodName
}

// GetOutputShapeGoType returns the Go type of the supplied operation's Output
// shape, renamed to use the standardized svcsdk alias.
func (r *CRD) GetOutputShapeGoType(
	op *awssdkmodel.Operation,
) string {
	if op == nil {
		panic("called GetOutputShapeGoType on nil operation.")
	}
	orig := op.OutputRef.GoType()
	// orig will contain "*<OutputShape>" with no package specifier
	return "*svcsdk." + orig[1:]
}

// GetOutputWrapperFieldPath returns the JSON-Path of the output wrapper field
// as *string for a given operation, if specified in generator config.
func (r *CRD) GetOutputWrapperFieldPath(
	op *awssdkmodel.Operation,
) *string {
	if op == nil {
		return nil
	}
	if r.cfg == nil {
		return nil
	}
	opConfig, found := r.cfg.Operations[op.Name]
	if !found {
		return nil
	}

	if opConfig.OutputWrapperFieldPath == "" {
		return nil
	}
	return &opConfig.OutputWrapperFieldPath
}

// GetOutputShape returns the Output shape for given operation.
func (r *CRD) GetOutputShape(
	// The operation to look for the Output shape
	op *awssdkmodel.Operation,
) (*awssdkmodel.Shape, error) {
	if op == nil {
		return nil, errors.New("no output shape for nil operation")
	}

	outputShape := op.OutputRef.Shape
	if outputShape == nil {
		return nil, errors.New("output shape not found")
	}

	// We might be in a "wrapper" shape. Unwrap it to find the real object
	// representation for the CRD's createOp/DescribeOP.

	// Use the wrapper field path if it's given in the ack-generate config file.
	wrapperFieldPath := r.GetOutputWrapperFieldPath(op)
	if wrapperFieldPath != nil {
		wrapperOutputShape, err := r.GetWrapperOutputShape(outputShape, *wrapperFieldPath)
		if err != nil {
			return nil, fmt.Errorf("unable to unwrap the output shape: %v", err)
		}
		outputShape = wrapperOutputShape
	} else {
		// If the wrapper field path is not specified in the config file and if
		// there is a single member shape and that member shape is a structure,
		// unwrap it.
		if outputShape.UsedAsOutput && len(outputShape.MemberRefs) == 1 {
			for _, memberRef := range outputShape.MemberRefs {
				if memberRef.Shape.Type == "structure" {
					outputShape = memberRef.Shape
				}
			}
		}
	}
	return outputShape, nil
}

// GetWrapperOutputShape returns the shape of the last element of a given field
// Path. It carefully unwraps the output shape and verifies that every element
// of the field path exists in their correspanding parent shape and that they are
// structures.
func (r *CRD) GetWrapperOutputShape(
	shape *awssdkmodel.Shape,
	fieldPath string,
) (*awssdkmodel.Shape, error) {
	if fieldPath == "" {
		return shape, nil
	}
	fieldPathParts := strings.Split(fieldPath, ".")
	for x, wrapperField := range fieldPathParts {
		for memberName, memberRef := range shape.MemberRefs {
			if memberName == wrapperField {
				if memberRef.Shape.Type != "structure" {
					// All the mentionned shapes must be structure
					return nil, fmt.Errorf(
						"Expected SetOutput.WrapperFieldPath to only contain fields of type 'structure'."+
							" Found %s of type '%s'",
						memberName, memberRef.Shape.Type,
					)
				}
				remainPath := strings.Join(fieldPathParts[x+1:], ".")
				return r.GetWrapperOutputShape(memberRef.Shape, remainPath)
			}
		}
		return nil, fmt.Errorf(
			"Incorrect SetOutput.WrapperFieldPath. Could not find %s in Shape %s",
			wrapperField, shape.ShapeName,
		)
	}
	return shape, nil
}

// GetCustomImplementation returns custom implementation method name for the
// supplied operation as specified in generator config
func (r *CRD) GetCustomImplementation(
	// The type of operation
	op *awssdkmodel.Operation,
) string {
	if op == nil || r.cfg == nil {
		return ""
	}

	operationConfig, found := r.cfg.Operations[op.Name]
	if !found {
		return ""
	}

	return operationConfig.CustomImplementation
}

// UpdateConditionsCustomMethodName returns custom update conditions operation
// as *string for custom resource, if specified in generator config
func (r *CRD) UpdateConditionsCustomMethodName() string {
	if r.cfg == nil {
		return ""
	}
	resGenConfig, found := r.cfg.Resources[r.Names.Original]
	if !found {
		return ""
	}
	return resGenConfig.UpdateConditionsCustomMethodName
}

// CustomCheckRequiredFieldsMissingMethod returns custom check required fields missing method
// as *string for custom resource, if specified in generator config
func (r *CRD) CustomCheckRequiredFieldsMissingMethod(op *awssdkmodel.Operation) string {
	if op == nil || r.cfg == nil {
		return ""
	}

	operationConfig, found := r.cfg.Operations[op.Name]
	if !found {
		return ""
	}

	return operationConfig.CustomCheckRequiredFieldsMissingMethod
}

// SpecIdentifierField returns the name of the "Name" or string identifier field in the Spec
func (r *CRD) SpecIdentifierField() *string {
	if r.cfg != nil {
		rConfig, found := r.cfg.Resources[r.Names.Original]
		if found {
			for fName, fConfig := range rConfig.Fields {
				if fConfig.IsName {
					return &fName
				}
			}
		}
	}
	lookup := []string{
		"Name",
		r.Names.Original + "Name",
		r.Names.Original + "Id",
	}
	for _, memberName := range r.SpecFieldNames() {
		if util.InStrings(memberName, lookup) {
			return &r.SpecFields[memberName].Names.Camel
		}
	}
	return nil
}

// IsAdoptable returns true if the resource can be adopted
func (r *CRD) IsAdoptable() bool {
	if r.cfg == nil {
		// Should never reach this condition
		return false
	}
	return r.cfg.ResourceIsAdoptable(r.Names.Original)
}

// GetResourcePrintOrderByName returns the Printer Column order-by field name
func (r *CRD) GetResourcePrintOrderByName() string {
	orderBy := r.cfg.GetResourcePrintOrderByName(r.Names.Camel)
	if orderBy == "" {
		return "name"
	}
	return orderBy
}

// PrintAgeColumn returns whether the code generator should append 'Age'
// kubebuilder:printcolumn comment marker
func (r *CRD) PrintAgeColumn() bool {
	return r.cfg.GetResourcePrintAddAgeColumn(r.Names.Camel)
}

// ReconcileRequeuOnSuccessSeconds returns the duration after which to requeue
// the custom resource as int, if specified in generator config.
func (r *CRD) ReconcileRequeuOnSuccessSeconds() int {
	if r.cfg == nil {
		return 0
	}
	resGenConfig, found := r.cfg.Resources[r.Names.Original]
	if !found {
		return 0
	}
	reconcile := resGenConfig.Reconcile
	if reconcile != nil {
		return reconcile.RequeueOnSuccessSeconds
	}
	// handles the default case
	return 0
}

// CustomUpdateMethodName returns the name of the custom resourceManager method
// for updating the resource state, if any has been specified in the generator
// config
func (r *CRD) CustomUpdateMethodName() string {
	if r.cfg == nil {
		return ""
	}
	rConfig, found := r.cfg.Resources[r.Names.Original]
	if found {
		if rConfig.UpdateOperation != nil {
			return rConfig.UpdateOperation.CustomMethodName
		}
	}
	return ""
}

// ListOpMatchFieldNames returns a slice of strings representing the field
// names in the List operation's Output shape's element Shape that we should
// check a corresponding value in the target Spec exists.
func (r *CRD) ListOpMatchFieldNames() []string {
	return r.cfg.ListOpMatchFieldNames(r.Names.Original)
}

// NewCRD returns a pointer to a new `ackmodel.CRD` struct that describes a
// single top-level resource in an AWS service API
func NewCRD(
	sdkAPI *SDKAPI,
	cfg *ackgenconfig.Config,
	crdNames names.Names,
	ops Ops,
) *CRD {
	pluralize := pluralize.NewClient()
	kind := crdNames.Camel
	plural := pluralize.Plural(kind)
	return &CRD{
		sdkAPI:                   sdkAPI,
		cfg:                      cfg,
		Names:                    crdNames,
		Kind:                     kind,
		Plural:                   plural,
		Ops:                      ops,
		additionalPrinterColumns: make([]*PrinterColumn, 0),
		SpecFields:               map[string]*Field{},
		StatusFields:             map[string]*Field{},
		Fields:                   map[string]*Field{},
		ShortNames:               cfg.ResourceShortNames(kind),
	}
}
