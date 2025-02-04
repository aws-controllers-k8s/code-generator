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
	"fmt"
	"strings"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/pkg/names"
	"github.com/gertd/go-pluralize"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

// simpleStringShapeRef is used for attribute fields and fields where there was
// no found ShapeRef
var simpleStringShapeRef *awssdkmodel.ShapeRef = &awssdkmodel.ShapeRef{
	Shape: &awssdkmodel.Shape{
		Type: "string",
	},
}

// Field represents a single field in the CRD's Spec or Status objects. The
// field may be a direct field of the Spec or Status object or may be a field
// of a list or struct-type field of the Spec or Status object. We call these
// latter fields "nested fields" and they are identified by the Field.Path
// attribute.
type Field struct {
	// CRD is the a pointer to the top-level custom resource definition
	// descriptor for the field or field's parent (if a nested field)
	CRD *CRD
	// Names is a set of normalized names for the field
	Names names.Names
	// Path is a "field path" that indicates where the field is within the CRD.
	// For example "Spec.Name" or "Status.BrokerInstances..Endpoint". Note for
	// the latter example, the field path indicates that the field `Endpoint`
	// is an attribute of the `Status.BrokerInstances` top-level field and the
	// double dot (`..` indicates that BrokerInstances is a list type).
	Path string
	// GoType is a string containing the Go data type for the field
	GoType string
	// GoTypeElem indicates the Go data type for the type of list element if
	// the field is a list type
	GoTypeElem        string
	GoTypeWithPkgName string
	ShapeRef          *awssdkmodel.ShapeRef
	FieldConfig       *ackgenconfig.FieldConfig

	// MemberFields is a map, keyed by the *renamed, cleaned, camel-cased name* of
	// member fields when this Field is a struct type.
	MemberFields map[string]*Field
}

// GetDocumentation returns a string containing the field's
// description/docstring. The ShapeRef from the AWS SDK is used as the default
// value for the documentation of the field. If a documentation override has
// been provided (in the `documentation.yaml`) file, then it will either be
// prependend or appended, or entirely overridden.
//
// For example, if there is a field with a ShapeRef that has a Documentation
// string containing:
//
// "// This field contains the identifier for the cluster
//
//	running the cache services"
//
// and the field has a documentation config specifies the following should be
// appended:
//
// "please note that this field is updated on the service
//
//	side"
//
// then the string returned from this method will be:
//
// "// This field contains the identifier for the cluster
//
//	// running the cache services
//	// please note that this field is updated on the service
//	// side"
func (f *Field) GetDocumentation() string {
	cfg := f.GetFieldDocsConfig()

	hasShapeDoc := false
	var prepend strings.Builder
	var out strings.Builder

	if f.ShapeRef != nil {
		if f.ShapeRef.Documentation != "" {
			hasShapeDoc = true
			out.WriteString(f.ShapeRef.Documentation)
		}
	}
	if cfg == nil {
		return out.String()
	}

	// Ensure configuration has exclusive options
	if cfg.Override != nil && (cfg.Append != nil || cfg.Prepend != nil) {
		panic("Documentation cannot contain override and prepend/append. Field: " + f.CRD.Kind)
	}

	if cfg.Override != nil {
		return f.formatUserProvidedDocstring(*cfg.Override)
	}

	if cfg.Append != nil {
		if hasShapeDoc {
			out.WriteString("\n//\n")
		}
		out.WriteString(f.formatUserProvidedDocstring(*cfg.Append))
	}

	if cfg.Prepend != nil {
		prepend.WriteString(f.formatUserProvidedDocstring(*cfg.Prepend))
		if hasShapeDoc {
			prepend.WriteString("\n//\n")
		}
	}

	return prepend.String() + out.String()
}

// formatUserProvidedDocstring sanitises a doc string provided by the user so
// that it fits the format of the existing AWS SDK GoDocs. This method will
// split it by lines, strip it of leading whitespace and slashes and prepend
// `//` to each line.
func (f *Field) formatUserProvidedDocstring(in string) string {
	var sb strings.Builder
	// Strip any leading comment slashes from the config option
	// docstring since we'll be automatically adding the Go comment
	// slashes to the beginning of each new line
	lines := strings.Split(in, "\n")

	// Traverse in reverse order until, deleting empty lines
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) != "" {
			break
		}
		lines = lines[:i]
	}

	numLines := len(lines)
	for idx, line := range lines {
		sb.WriteString("// ")
		sb.WriteString(strings.TrimLeft(line, " \\\t/"))
		if idx < (numLines - 1) {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

// IsRequired checks the FieldConfig for Field and returns if the field is
// marked as required or not.A
//
// If there is no required override present for this field in FieldConfig,
// IsRequired will return if the shape is marked as required in AWS SDK Private
// model We use this to append kubebuilder:validation:Required markers to
// validate using the CRD validation schema
func (f *Field) IsRequired() bool {
	if f.FieldConfig != nil && f.FieldConfig.IsRequired != nil {
		return *f.FieldConfig.IsRequired
	}
	// We need to look up the original member name in the input struct
	// otherwise renamed fields will not be discovered as required.
	originalMember := f.CRD.Config().GetOriginalMemberName(
		f.CRD.Names.Original,
		f.CRD.Ops.Create.Name,
		f.Names.Original,
	)
	return util.InStrings(
		originalMember,
		f.CRD.Ops.Create.InputRef.Shape.Required,
	)
}

// GetSetterConfig returns the SetFieldConfig object associated with this field
// and a supplied operation type, or nil if none exists.
func (f *Field) GetSetterConfig(opType OpType) *ackgenconfig.SetFieldConfig {
	if f.FieldConfig == nil {
		return nil
	}
	rmMethod := ResourceManagerMethodFromOpType(opType)
	for _, setCfg := range f.FieldConfig.Set {
		if setCfg == nil {
			continue
		}
		// If the Method attribute is nil, that means the setter config applies to
		// all resource manager methods for this field.
		if setCfg.Method == nil || strings.EqualFold(rmMethod, *setCfg.Method) {
			return setCfg
		}
	}
	return nil
}

// GetFieldDocsConfig returns the field documentation configuration for the
// current field if it exists, otherwise it returns nil.
func (f *Field) GetFieldDocsConfig() *ackgenconfig.FieldDocsConfig {
	resourceConfig, exists := f.CRD.docCfg.Resources[f.CRD.Names.Camel]
	if !exists {
		return nil
	}

	return resourceConfig.Fields[f.Path]
}

// GetGoTag returns the json tag for the field. If the field has a
// FieldConfig with a GoTag attribute, the value of GoTag is returned.
// Otherwise, we evaluate the field type and return the appropriate
// json tag.
func (f *Field) GetGoTag() string {
	// First check if the field has a GoTag attribute in the FieldConfig
	// a.k.a generator.yaml
	if f.FieldConfig != nil && f.FieldConfig.GoTag != nil {
		return fmt.Sprintf("`%s`", *f.FieldConfig.GoTag)
	}

	// If the field is not required, a reference field or part of the status
	// object, we need to inject the `omitempty`` directive into the json tag.
	if !f.IsRequired() || f.HasReference() || f.CRD.StatusFields[f.Names.Camel] != nil {
		return fmt.Sprintf("`json:\"%s,omitempty\"`", f.Names.CamelLower)

	}

	return fmt.Sprintf("`json:\"%s\"`", f.Names.CamelLower)
}

// HasReference returns true if the supplied field *path* refers to a Field
// that contains 'ReferencesConfig' i.e. has a corresponding reference field.
// Ex:
// ```
// Integration:
//
//	fields:
//	  ApiId:
//	    references:
//	      resource: API
//	      path: Status.APIID
//
// ```
// For the above configuration, 'HasReference' for 'ApiId'(Original name) field
// will return true because a corresponding 'APIRef' field will be generated
// by ACK code-generator
func (f *Field) HasReference() bool {
	return f.FieldConfig != nil && f.FieldConfig.References != nil
}

// IsReference returns true if the Field has type '*ackv1alpha1.AWSResourceReferenceWrapper'
// or '[]*ackv1alpha1.AWSResourceReferenceWrapper'.
// These fields are not part of aws-sdk-go model and they are generated by
// ACK code-generator to accept references of other resource(s).
func (f *Field) IsReference() bool {
	trimmedGoType := strings.TrimPrefix(f.GoType, "[]")
	return trimmedGoType == "*ackv1alpha1.AWSResourceReferenceWrapper"
}

// GetReferenceFieldName returns the corresponding Reference field name
// Reference field name is generated by removing the identifier suffix like
// 'Id', 'Arn', 'Name' etc and adding 'Ref(s)' as suffix.
func (f *Field) GetReferenceFieldName() names.Names {
	if f.Names.Original == "" {
		panic(fmt.Sprintf("field with empty name inside crd %s", f.CRD.Names.Original))
	}
	refNamePrefix := f.Names.Original
	identifierSuffixes := []string{
		"id", "ids", "Id", "Ids", "ID", "IDs", "IDS",
		"Name", "Names", "NAME", "NAMEs", "NAMES",
		"Arn", "Arns", "ARN", "ARNs", "ARNS",
	}
	for _, suffix := range identifierSuffixes {
		if strings.HasSuffix(f.Names.Original, suffix) {
			refNamePrefix = strings.TrimSuffix(refNamePrefix, suffix)
			break
		}
	}
	if refNamePrefix == "" {
		panic("The corresponding field name for a reference field cannot" +
			" be just 'id(s)' or 'arn(s)' or 'name(s)'. Current value: " +
			f.Names.Original)
	}
	refName := refNamePrefix
	// If the shape of corresponding field is a list, singularize the refNamePrefix
	// and add Refs at the end
	if f.ShapeRef != nil && f.ShapeRef.Shape != nil && f.ShapeRef.Shape.Type == "list" {
		refName = fmt.Sprintf("%sRefs", pluralize.NewClient().Singular(refNamePrefix))
	} else {
		refName = fmt.Sprintf("%sRef", refNamePrefix)
	}
	return names.New(refName)
}

// ReferencedServiceName returns the serviceName for the referenced resource
// when the field has 'ReferencesConfig'
// If the field does not have 'ReferencesConfig', empty string is returned
func (f *Field) ReferencedServiceName() (referencedServiceName string) {
	if f.FieldConfig != nil && f.FieldConfig.References != nil {
		if f.FieldConfig.References.ServiceName != "" {
			return f.FieldConfig.References.ServiceName
		} else {
			return f.CRD.sdkAPI.API.PackageName()
		}
	}
	return referencedServiceName
}

// ReferencedResourceNamePlural returns the plural of referenced resource
// when the field has a 'ReferencesConfig'
// If the field does not have 'ReferencesConfig', empty string is returned
func (f *Field) ReferencedResourceNamePlural() string {
	var referencedResourceName string
	pluralize := pluralize.NewClient()
	if f.FieldConfig != nil && f.FieldConfig.References != nil {
		referencedResourceName = f.FieldConfig.References.Resource
	}
	if referencedResourceName != "" {
		return pluralize.Plural(referencedResourceName)
	}
	return referencedResourceName
}

// ReferenceFieldPath returns the fieldPath for the corresponding
// Reference field. It replaces the fieldName with ReferenceFieldName
// at the end of fieldPath
func (f *Field) ReferenceFieldPath() string {
	fieldPathPrefix := strings.TrimSuffix(f.Path, f.Names.Camel)
	return fmt.Sprintf("%s%s", fieldPathPrefix, f.GetReferenceFieldName().Camel)
}

// FieldPathWithUnderscore replaces the period in fieldPath with
// underscore. This method is useful for generating go method
// name from the fieldPath.
func (f *Field) FieldPathWithUnderscore() string {
	return strings.ReplaceAll(f.Path, ".", "_")
}

// NewReferenceField returns a pointer to a new Field object.
// The go-type of field is either slice of '*AWSResourceReferenceWrapper' or
// '*AWSResourceReferenceWrapper' depending on whether 'shapeRef' parameter
// has 'list' type or not, respectively
func NewReferenceField(
	crd *CRD,
	fieldNames names.Names,
	shapeRef *awssdkmodel.ShapeRef,
) *Field {
	gt := "*ackv1alpha1.AWSResourceReferenceWrapper"
	gtp := "*github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1.AWSResourceReferenceWrapper"
	gte := ""
	if shapeRef.Shape.Type == "list" {
		gt = "[]" + gt
		gtp = "[]" + gtp
		gte = "*ackv1alpha1.AWSResourceReferenceWrapper"
	}
	return &Field{
		CRD:               crd,
		Names:             fieldNames,
		Path:              fieldNames.Original,
		ShapeRef:          nil,
		GoType:            gt,
		GoTypeElem:        gte,
		GoTypeWithPkgName: gtp,
		FieldConfig:       nil,
	}
}

// NewField returns a pointer to a new Field object
func NewField(
	crd *CRD,
	path string,
	fieldNames names.Names,
	shapeRef *awssdkmodel.ShapeRef,
	cfg *ackgenconfig.FieldConfig,
) *Field {
	return newFieldRecurse(crd, path, make(map[string]struct{}, 0), fieldNames, shapeRef, cfg)
}

// newFieldRecurse recursively calls itself with protection against infinite
// cycles
func newFieldRecurse(
	crd *CRD,
	path string,
	// A map of shape name to Fields. Keeps track of every shape that has been
	// recursively discovered so that we can detect cycles. For example, the
	// emrcontainers `Configuration` object contains a property of type
	// `Configuration`.
	parentFields map[string]struct{},
	fieldNames names.Names,
	shapeRef *awssdkmodel.ShapeRef,
	cfg *ackgenconfig.FieldConfig,
) *Field {
	memberFields := map[string]*Field{}
	var gte, gt, gtwp string
	var shape *awssdkmodel.Shape
	if shapeRef != nil {
		shape = shapeRef.Shape
	}
	// this is a pointer to the "parent" containing Shape when the field being
	// processed here is a structure or a list/map of structures.
	var containerShape *awssdkmodel.Shape = shape

	if shape != nil {
		gte, gt, gtwp = CleanGoType(crd.sdkAPI, crd.cfg, shape, cfg)
		for {
			// If the field is a slice or map of structs, we want to add
			// MemberFields that describe the list or value struct elements so
			// that a field path can be used to "find" nested struct member
			// fields.
			//
			// For example, the EC2 resource DHCPOptions has a Field called
			// DHCPConfigurations which is of type []*NewDHCPConfiguration
			// where the NewDHCPConfiguration struct contains two fields, Key
			// and Values. If we want to be able to refer to the
			// DHCPOptions.DHCPConfigurations.Values field by field path, we
			// need a Field.MemberField that describes the
			// NewDHCPConfiguration.Values field.
			//
			// Here, we essentially dive down into list or map fields,
			// searching for whether the list or map fields have structure list
			// element or value element types and then rely on the code below
			// to "unpack" those struct member fields.
			if containerShape.Type == "list" {
				containerShape = containerShape.MemberRef.Shape
				continue
			} else if containerShape.Type == "map" {
				containerShape = containerShape.ValueRef.Shape
				continue
			}
			break
		}

		// Copy the parent fields map so that we can use it recursively without
		// modifying it for other call stacks
		nestedParentFields := make(map[string]struct{}, len(parentFields))
		for k, v := range parentFields {
			nestedParentFields[k] = v
		}
		nestedParentFields[shapeRef.ShapeName] = struct{}{}

		if containerShape.Type == "structure" {
			// "unpack" the member fields composing this struct field...
			for _, memberName := range containerShape.MemberNames() {
				cleanMemberNames := names.New(memberName)
				memberPath := path + "." + cleanMemberNames.Camel
				memberShape := containerShape.MemberRefs[memberName]

				// Check to see if we have seen this shape before in the stack
				// and panic. Cyclic references are not supported
				if _, ok := nestedParentFields[memberShape.ShapeName]; ok {
					panic(fmt.Sprintf("Detected a cyclic type reference in %s"+
						". Add the corresponding path to `ignore.field_paths`"+
						" in the generator config to continue.",
						containerShape.ShapeName))
				}

				fConfigs := crd.cfg.GetFieldConfigs(crd.Names.Original)
				memberField := newFieldRecurse(
					crd, memberPath, nestedParentFields, cleanMemberNames, memberShape, fConfigs[memberPath],
				)

				memberFields[cleanMemberNames.Camel] = memberField
			}
		}
	} else {
		gte = "string"
		gt = "*string"
		gtwp = "*string"
		shapeRef = simpleStringShapeRef
	}

	return &Field{
		CRD:               crd,
		Names:             fieldNames,
		Path:              path,
		ShapeRef:          shapeRef,
		GoType:            gt,
		GoTypeElem:        gte,
		GoTypeWithPkgName: gtwp,
		FieldConfig:       cfg,
		MemberFields:      memberFields,
	}
}
