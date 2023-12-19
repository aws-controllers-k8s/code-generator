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

	"github.com/aws-controllers-k8s/pkg/names"
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/gertd/go-pluralize"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/fieldpath"
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
	docCfg *ackgenconfig.DocumentationConfig
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

// GetStorageVersion returns the configured storage API version for the CRD, or
// the specified default version.
func (r *CRD) GetStorageVersion(defaultVersion string) (string, error) {
	apiVersions := r.cfg.GetAPIVersions(r.Names.Original)
	// if not configured
	if len(apiVersions) == 0 {
		return defaultVersion, nil
	}

	for _, v := range apiVersions {
		if v.Storage != nil && *v.Storage {
			return v.Name, nil
		}
	}
	return "", fmt.Errorf("exactly one configured version must be marked as the storage version for the %q CRD",
		r.Names.Original)
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
	for _, field := range r.SpecFields {
		if field.ShapeRef != nil && shapeHasMember(field.ShapeRef.Shape, toFind) {
			return true
		}
	}
	for _, field := range r.StatusFields {
		if field.ShapeRef != nil && shapeHasMember(field.ShapeRef.Shape, toFind) {
			return true
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

// AddSpecField adds a new Field of a given name and shape into the Spec
// field of a CRD
func (r *CRD) AddSpecField(
	memberNames names.Names,
	shapeRef *awssdkmodel.ShapeRef,
) {
	fPath := memberNames.Camel
	fConfig := r.cfg.GetFieldConfigByPath(r.Names.Original, fPath)
	f := NewField(r, fPath, memberNames, shapeRef, fConfig)
	if fConfig != nil && fConfig.Print != nil {
		r.addSpecPrintableColumn(f)
	}
	r.SpecFields[memberNames.Original] = f
	r.Fields[fPath] = f

	// If this field has a ReferencesConfig, Add the new
	// Reference field inside Spec as well
	if fConfig != nil && fConfig.References != nil {
		referenceFieldNames := f.GetReferenceFieldName()
		rf := NewReferenceField(r, referenceFieldNames, shapeRef)
		r.SpecFields[referenceFieldNames.Original] = rf
		r.Fields[referenceFieldNames.Camel] = rf
	}
}

// AddStatusField adds a new Field of a given name and shape into the Status
// field of a CRD
func (r *CRD) AddStatusField(
	memberNames names.Names,
	shapeRef *awssdkmodel.ShapeRef,
) {
	fPath := memberNames.Camel
	fConfig := r.cfg.GetFieldConfigByPath(r.Names.Original, fPath)
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
	return r.cfg.ResourceContainsAttributesMap(r.Names.Original)
}

// CompareIgnoredFields returns the list of fields compare logic should ignore
func (r *CRD) CompareIgnoredFields() []string {
	return r.cfg.GetCompareIgnoredFieldPaths(r.Names.Original)
}

// SetAttributesSingleAttribute returns true if the supplied resource name has
// a SetAttributes operation that only actually changes a single attribute at a
// time. See: SNS SetTopicAttributes API call, which is entirely different from
// the SNS SetPlatformApplicationAttributes API call, which sets multiple
// attributes at once. :shrug:
func (r *CRD) SetAttributesSingleAttribute() bool {
	return r.cfg.ResourceSetsSingleAttribute(r.Names.Original)
}

// UnpackAttributes grabs instructions about fields that are represented in the
// AWS API as a `map[string]*string` but are actually real, schema'd fields and
// adds Field definitions for those fields.
func (r *CRD) UnpackAttributes() {
	if !r.cfg.ResourceContainsAttributesMap(r.Names.Original) {
		return
	}
	fieldConfigs := r.cfg.GetFieldConfigs(r.Names.Original)
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
	if !r.cfg.IncludeACKMetadata {
		return false
	}
	fieldConfigs := r.cfg.GetFieldConfigs(r.Names.Original)
	for fName, fConfig := range fieldConfigs {
		if fConfig.IsARN {
			return strings.EqualFold(fieldName, fName)
		}
	}
	return strings.EqualFold(fieldName, "arn") ||
		strings.EqualFold(fieldName, r.Names.Original+"arn")
}

// IsSecretField returns true if the supplied field *path* refers to a Field
// that is a SecretKeyReference
func (r *CRD) IsSecretField(path string) bool {
	fConfigs := r.cfg.GetFieldConfigs(r.Names.Original)
	fConfig, found := fConfigs[path]
	if found {
		return fConfig.IsSecret
	}
	return false
}

// GetImmutableFieldPaths returns list of immutable field paths present in CRD
func (r *CRD) GetImmutableFieldPaths() []string {
	fConfigs := r.cfg.GetFieldConfigs(r.Names.Original)
	var immutableFields []string

	for field, fieldConfig := range fConfigs {
		if fieldConfig.IsImmutable {
			immutableFields = append(immutableFields, field)
		}
	}

	// We need a deterministic order to traverse the immutable fields
	sort.Strings(immutableFields)
	return immutableFields
}

// HasImmutableFieldChanges helper function that return true if there are any immutable field changes
func (r *CRD) HasImmutableFieldChanges() bool {
	fConfigs := r.cfg.GetFieldConfigs(r.Names.Original)
	for _, fieldConfig := range fConfigs {
		if fieldConfig.IsImmutable {
			return true
		}
	}
	return false
}

// OmitUnchangedFieldsOnUpdate returns whether the controller needs to omit
// unchanged fields from an update request or not.
func (r *CRD) OmitUnchangedFieldsOnUpdate() bool {
	if r.Config() == nil {
		return false
	}
	rConfig, found := r.Config().Resources[r.Names.Original]
	if found {
		if rConfig.UpdateOperation != nil {
			return rConfig.UpdateOperation.OmitUnchangedFields
		}
	}
	return false
}

// IsARNPrimaryKey returns true if the CRD uses its ARN as its primary key in
// ReadOne calls.
func (r *CRD) IsARNPrimaryKey() bool {
	resGenConfig := r.cfg.GetResourceConfig(r.Names.Original)
	if resGenConfig == nil {
		return false
	}
	return resGenConfig.IsARNPrimaryKey
}

// GetPrimaryKeyField returns the field designated as the primary key, nil if
// none are specified or an error if multiple are designated.
func (r *CRD) GetPrimaryKeyField() (*Field, error) {
	fConfigs := r.cfg.GetFieldConfigs(r.Names.Original)

	var primaryField *Field
	for fieldName, fieldConfig := range fConfigs {
		if !fieldConfig.IsPrimaryKey {
			continue
		}

		// Multiple primary fields
		if primaryField != nil {
			return nil, fmt.Errorf("multiple fields are marked with is_primary_key")
		}

		fieldNames := names.New(fieldName)
		fPath := fieldNames.Camel
		var found bool
		primaryField, found = r.Fields[fPath]
		if !found {
			return nil, fmt.Errorf("could not find field with path " + fPath +
				" for primary key " + fieldName)
		}
	}
	return primaryField, nil
}

// GetMatchingInputShapeFieldName returns the name of the field in the Input shape.
// For simplicity, we assume that there will be only one setConfig for the
// any unique sdkField, per operation. Which means that we will never set
// two different sdk field from the same
func (r *CRD) GetMatchingInputShapeFieldName(opType OpType, sdkField string) string {
	// At this stage nil-checks for r.cfg is not necessary
	for _, f := range r.Fields {
		if f.FieldConfig == nil {
			continue
		}
		rmMethod := ResourceManagerMethodFromOpType(opType)
		for _, setCfg := range f.FieldConfig.Set {
			if setCfg == nil {
				continue
			}
			if setCfg.Ignore == true || setCfg.To == nil {
				continue
			}
			// If the Method attribute is nil, that means the setter config applies to
			// all resource manager methods for this field.
			if setCfg.Method == nil || strings.EqualFold(rmMethod, *setCfg.Method) {
				if setCfg.To != nil && *setCfg.To == sdkField {
					return f.Names.Camel
				}
			}
		}
	}
	return ""
}

// SetOutputCustomMethodName returns custom set output operation as *string for
// given operation on custom resource
func (r *CRD) SetOutputCustomMethodName(
	// The operation to look for the Output shape
	op *awssdkmodel.Operation,
) *string {
	return r.cfg.GetSetOutputCustomMethodName(op)
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
// as *string for a given operation.
func (r *CRD) GetOutputWrapperFieldPath(
	op *awssdkmodel.Operation,
) *string {
	return r.cfg.GetOutputWrapperFieldPath(op)
}

// GetOutputShape returns the Output shape for a given operation and applies
// wrapper field path overrides.
func (r *CRD) GetOutputShape(
	op *awssdkmodel.Operation,
) (*awssdkmodel.Shape, error) {
	if op == nil {
		return nil, errors.New("no output shape for nil operation")
	}

	outputShape := op.OutputRef.Shape
	if outputShape == nil {
		return nil, errors.New("output shape not found")
	}

	// Check for wrapper field path overrides
	wrapperFieldPath := r.GetOutputWrapperFieldPath(op)
	if wrapperFieldPath != nil {
		wrapperOutputShape, err := r.getWrapperOutputShape(outputShape,
			*wrapperFieldPath)
		if err != nil {
			msg := fmt.Sprintf("Unable to unwrap the output shape: %s "+
				"with field path override: %s. error: %v",
				outputShape.OrigShapeName, *wrapperFieldPath, err)
			panic(msg)
		}
		outputShape = wrapperOutputShape
	}
	return outputShape, nil
}

// getWrapperOutputShape returns the shape of the last element of a given field
// Path. It unwraps the output shape and verifies that every element of the
// field path exists in their corresponding parent shape and that they are
// structures.
func (r *CRD) getWrapperOutputShape(
	shape *awssdkmodel.Shape,
	fieldPath string,
) (*awssdkmodel.Shape, error) {
	if fieldPath == "" {
		return shape, nil
	}
	fp := fieldpath.FromString(fieldPath)
	wrapperField := fp.PopFront()

	memberRef, ok := shape.MemberRefs[wrapperField]
	if !ok {
		return nil, fmt.Errorf(
			"could not find wrapper override field %s in Shape %s",
			wrapperField, shape.ShapeName)
	}

	// wrapper field must be list or structure; otherwise cannot unpack
	if memberRef.Shape.Type == "list" {
		memberRef = &memberRef.Shape.MemberRef
	}
	if memberRef.Shape.Type != "structure" {
		return nil, fmt.Errorf(
			"output wrapper overrides can only contain fields of type"+
				" 'structure'. Found wrapper override field %s of type '%s'",
			wrapperField, memberRef.Shape.Type)
	}
	return r.getWrapperOutputShape(memberRef.Shape, fp.String())
}

// GetCustomImplementation returns custom implementation method name for the
// supplied operation as specified in generator config
func (r *CRD) GetCustomImplementation(
	// The type of operation
	op *awssdkmodel.Operation,
) string {
	return r.cfg.GetCustomImplementation(op)
}

// UpdateConditionsCustomMethodName returns custom update conditions operation
// as *string for custom resource
func (r *CRD) UpdateConditionsCustomMethodName() string {
	return r.cfg.GetUpdateConditionsCustomMethodName(r.Names.Original)
}

// GetCustomCheckRequiredFieldsMissingMethod returns custom check required fields missing method
// as string for custom resource
func (r *CRD) GetCustomCheckRequiredFieldsMissingMethod(
	// The type of operation
	op *awssdkmodel.Operation,
) string {
	return r.cfg.GetCustomCheckRequiredFieldsMissingMethod(op)
}

// SpecIdentifierField returns the name of the "Name" or string identifier field
// in the Spec.
func (r *CRD) SpecIdentifierField() *string {
	rConfig := r.cfg.GetResourceConfig(r.Names.Original)
	if rConfig != nil {
		for fName, fConfig := range rConfig.Fields {
			if fConfig.IsPrimaryKey {
				return &fName
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
	return r.cfg.ResourceDisplaysAgeColumn(r.Names.Camel)
}

// PrintSyncedColumn returns whether the code generator should append 'Sync'
// kubebuilder:printcolumn comment marker
func (r *CRD) PrintSyncedColumn() bool {
	return r.cfg.ResourceDisplaysSyncedColumn(r.Names.Camel)
}

func (r *CRD) addAdditionalPrinterColumns(additionalColumns []*ackgenconfig.AdditionalColumnConfig) {
	for _, additionalColumn := range additionalColumns {
		printerColumn := &PrinterColumn{}
		printerColumn.Name = additionalColumn.Name
		printerColumn.JSONPath = additionalColumn.JSONPath
		printerColumn.Type = additionalColumn.Type
		printerColumn.Priority = additionalColumn.Priority
		printerColumn.Index = additionalColumn.Index
		r.additionalPrinterColumns = append(r.additionalPrinterColumns, printerColumn)
	}
}

// checkSpecOrStatus checks whether the new nested field is for Spec or Status struct
// and returns the top level field accordingly
func (crd *CRD) checkSpecOrStatus(
	field string,
) (*Field, bool) {
	fSpec, okSpec := crd.SpecFields[field]
	if okSpec {
		return fSpec, okSpec
	}
	fStatus, okStatus := crd.StatusFields[field]
	if okStatus {
		return fStatus, okStatus
	}
	return nil, false
}

// addCustomNestedFields injects the provided customNestedFields into the
// Spec or Status struct as a nested field. The customNestedFields are
// identified by the field path. The field path is a dot separated string
// that represents the path to the nested field. For example, if we want to
// inject a field called "Password" into the "User" struct, the field path
// would be "User.Password". The field path can be as deep as needed.
func (crd *CRD) addCustomNestedFields(customNestedFields map[string]*ackgenconfig.FieldConfig) {
	//  We have collected all the nested fields in `customNestedFields` map and now we can process
	//  and validate that they indeed are inject-able (i.e. the parent field is of type struct)
	// and inject them into the Spec or Status struct.
	for customNestedField, customNestedFieldConfig := range customNestedFields {
		fieldParts := strings.Split(customNestedField, ".")
		// we know that the length of fieldParts is at least 2
		// it is safe to access the first element.
		topLevelField := fieldParts[0]

		f, ok := crd.checkSpecOrStatus(topLevelField)

		if ok && f.ShapeRef.Shape.Type != "structure" {
			// We need to panic here because the user is providing wrong configuration.
			msg := fmt.Sprintf("Expected parent field to be of type structure, but found %s", f.ShapeRef.Shape.Type)
			panic(msg)
		}

		// If the provided top level field is not in the crd.SpecFields or crd.StatusFields...
		if !ok {
			// We need to panic here because the user is providing wrong configuration.
			msg := fmt.Sprintf("Expected top level field %s to be present in Spec or Status", topLevelField)
			panic(msg)
		}

		// We will have to keep track of the previous field in the path
		// to check it's member fields.
		parentField := f

		// loop over the all left fieldParts except the last one
		for _, currentFieldName := range fieldParts[1 : len(fieldParts)-1] {
			// Check if parentField contains current field
			currentField, ok := parentField.MemberFields[currentFieldName]
			if !ok || currentField.ShapeRef.Shape.Type != "structure" {
				// Check if the field exists AND is of type structure
				msg := fmt.Sprintf("Cannot inject field, %s member doesn't exist or isn't a structure", currentFieldName)
				panic(msg)
			}
			parentField = currentField
		}

		// arriving here means that successfully walked the path and
		// parentField is the parent of the new field.

		// the last part is the field name
		fieldName := fieldParts[len(fieldParts)-1]
		typeOverride := customNestedFieldConfig.Type
		shapeRef := crd.sdkAPI.GetShapeRefFromType(*typeOverride)

		// Create a new field with the provided field name and shapeRef
		newCustomNestedField := NewField(crd, fieldName, names.New(fieldName), shapeRef, customNestedFieldConfig)

		// Add the new field to the parentField
		parentField.MemberFields[fieldName] = newCustomNestedField
		// Add the new field to the parentField's shapeRef
		parentField.ShapeRef.Shape.MemberRefs[fieldName] = crd.sdkAPI.GetShapeRefFromType(*customNestedFieldConfig.Type)
	}
}

// ReconcileRequeuOnSuccessSeconds returns the duration after which to requeue
// the custom resource as int
func (r *CRD) ReconcileRequeuOnSuccessSeconds() int {
	return r.cfg.GetReconcileRequeueOnSuccessSeconds(r.Names.Original)
}

// CustomUpdateMethodName returns the name of the custom resourceManager method
// for updating the resource state, if any has been specified in the generator
// config
func (r *CRD) CustomUpdateMethodName() string {
	return r.cfg.GetCustomUpdateMethodName(r.Names.Original)
}

func (r *CRD) CustomFindMethodName() string {
	return r.cfg.GetCustomFindMethodName(r.Names.Original)
}

// ListOpMatchFieldNames returns a slice of strings representing the field
// names in the List operation's Output shape's element Shape that we should
// check a corresponding value in the target Spec exists.
func (r *CRD) ListOpMatchFieldNames() []string {
	return r.cfg.GetListOpMatchFieldNames(r.Names.Original)
}

// GetAllRenames returns all the field renames observed in the generator config
// for a given OpType.
func (r *CRD) GetAllRenames(op OpType) map[string]string {
	opMap := r.sdkAPI.GetOperationMap(r.cfg)
	operations := (*opMap)[op]
	return r.cfg.GetAllRenames(r.Names.Original, operations)
}

// GetIdentifiers returns the identifier fields of a given CRD which
// can be singular or plural. Note, these fields will be the *original* field
// names from the API model shape, not renamed field names.
func (r *CRD) GetIdentifiers() []string {
	var identifiers []string
	if r == nil {
		return identifiers
	}
	identifierLookup := []string{
		"Id",
		"Ids",
		r.Names.Original + "Id",
		r.Names.Original + "Ids",
		"Name",
		"Names",
		r.Names.Original + "Name",
		r.Names.Original + "Names",
	}

	for _, id := range identifierLookup {
		_, found := r.SpecFields[id]
		if !found {
			_, found = r.StatusFields[id]
		}
		if found {
			identifiers = append(identifiers, id)
		}
	}

	return identifiers
}

// GetSanitizedMemberPath takes a shape member field, checks for renames, checks
// for existence in Spec and Status, then constructs and returns the var path.
// Returns error if memberName is not present in either Spec or Status.
func (r *CRD) GetSanitizedMemberPath(
	memberName string,
	op *awssdkmodel.Operation,
	koVarName string) (string, error) {
	resVarPath := koVarName
	cfg := r.Config()

	// Handles field renames, if applicable
	fieldName := cfg.GetResourceFieldName(
		r.Names.Original,
		op.ExportedName,
		memberName,
	)
	cleanFieldNames := names.New(fieldName)
	pathFieldName := cleanFieldNames.Camel

	inSpec, inStatus := r.HasMember(fieldName, op.ExportedName)
	if inSpec {
		resVarPath = resVarPath + cfg.PrefixConfig.SpecField + "." + pathFieldName
	} else if inStatus {
		resVarPath = resVarPath + cfg.PrefixConfig.StatusField + "." + pathFieldName
	} else {
		return "", fmt.Errorf(
			"the required field %s is NOT present in CR's Spec or Status", memberName)
	}
	return resVarPath, nil
}

// HasMember returns true in the respective field if Spec xor Status field
// contains memberName or rename
func (r *CRD) HasMember(
	memberName string,
	operationName string,
) (inSpec bool, inStatus bool) {
	fieldName := r.Config().GetResourceFieldName(
		r.Names.Original,
		operationName,
		memberName,
	)
	if _, found := r.SpecFields[fieldName]; found {
		inSpec = true
	} else if _, found := r.StatusFields[fieldName]; found {
		inStatus = true
	}
	return inSpec, inStatus
}

// HasReferenceFields returns true if any of the fields in CRD is a reference
// field. Otherwise returns false
func (r *CRD) HasReferenceFields() bool {
	for _, field := range r.Fields {
		if field.HasReference() {
			return true
		}
	}
	return false
}

// ReferencedServiceNames returns the set of service names for ACK controllers
// whose resources are referenced inside the CRD. The service name is
// the go package name for the AWS service inside aws-sdk-go.
//
// If a CRD has no reference fields, nil is returned(zero vale of slice)
func (r *CRD) ReferencedServiceNames() (serviceNames []string) {
	// We are using Map to implement a Set of service names
	serviceNamesMap := make(map[string]struct{})
	existsValue := struct{}{}

	for _, field := range r.Fields {
		if serviceName := field.ReferencedServiceName(); serviceName != "" {
			serviceNamesMap[serviceName] = existsValue
		}
	}

	for serviceName, _ := range serviceNamesMap {
		serviceNames = append(serviceNames, serviceName)
	}
	sort.Strings(serviceNames)
	return serviceNames
}

// SortedFieldNames returns the fieldNames of the CRD in a sorted
// order.
func (r *CRD) SortedFieldNames() []string {
	fieldNames := make([]string, 0, len(r.Fields))
	for fieldName := range r.Fields {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Strings(fieldNames)
	return fieldNames
}

// NewCRD returns a pointer to a new `ackmodel.CRD` struct that describes a
// single top-level resource in an AWS service API
func NewCRD(
	sdkAPI *SDKAPI,
	cfg *ackgenconfig.Config,
	docCfg *ackgenconfig.DocumentationConfig,
	crdNames names.Names,
	ops Ops,
) *CRD {
	pluralize := pluralize.NewClient()
	kind := crdNames.Camel
	plural := pluralize.Plural(kind)
	return &CRD{
		sdkAPI:                   sdkAPI,
		cfg:                      cfg,
		docCfg:                   docCfg,
		Names:                    crdNames,
		Kind:                     kind,
		Plural:                   plural,
		Ops:                      ops,
		additionalPrinterColumns: make([]*PrinterColumn, 0),
		SpecFields:               map[string]*Field{},
		StatusFields:             map[string]*Field{},
		Fields:                   map[string]*Field{},
		ShortNames:               cfg.GetResourceShortNames(kind),
	}
}
