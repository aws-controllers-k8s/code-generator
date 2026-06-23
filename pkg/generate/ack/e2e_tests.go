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

	"github.com/aws-controllers-k8s/pkg/names"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/templateset"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
)

var (
	e2eIncludePaths = []string{}
	e2eCopyPaths    = []string{}
	e2eFuncMap      = ttpl.FuncMap{
		"ToLower": strings.ToLower,
		"TruncateTo": func(maxLen int, s string) string {
			if len(s) > maxLen {
				return s[:maxLen]
			}
			return s
		},
	}
)

// VerifyCall represents a single AWS SDK verification call to make after
// create or update.
type VerifyCall struct {
	// OperationName is the SDK operation (e.g., "GetBucketVersioning")
	OperationName string
	// InputTypeName is the SDK input struct (e.g., "GetBucketVersioningInput")
	InputTypeName string
	// IdentifierAssignments maps SDK input struct field names to Go expressions
	// (e.g., "Bucket": "aws.String(name)")
	IdentifierAssignments map[string]string
	// Assertions are the field-level assertions on the output
	Assertions []VerifyAssertion
}

// VerifyAssertion represents a single assertion on a verification response.
type VerifyAssertion struct {
	// ResponsePath is the Go access path on the output (e.g., "Status")
	ResponsePath string
	// Expected is the Go expression for the expected value
	Expected string
	// Message is the failure message
	Message string
	// DerefFunc is an optional dereference function to wrap the response access
	// (e.g., "aws.ToBool" for *bool fields)
	DerefFunc string
}

// templateE2EMainVars holds template variables for the main_test.go template.
type templateE2EMainVars struct {
	templateset.MetaVars
	TestConfig *ackgenconfig.TestConfig
}

// templateE2EFileVars holds template variables for the per-resource test file.
type templateE2EFileVars struct {
	templateset.MetaVars
	CRD   *ackmodel.CRD
	Tests []templateE2ETestVars
}

// templateE2ETestVars holds template variables for a single named test within
// a resource test file.
type templateE2ETestVars struct {
	TestName           string // snake_case test name from testconfig key
	TestNameCamel      string // CamelCase for function name
	FieldAssignments   map[string]string
	UpdateAssignments  map[string]string
	HasUpdate          bool
	CreateWaitDuration string
	HasVerification    bool
	CreateVerifyCalls  []VerifyCall
	UpdateVerifyCalls  []VerifyCall
}

// E2ETests returns a pointer to a TemplateSet containing all templates for
// generating Go e2e test files for an ACK service controller.
func E2ETests(
	m *ackmodel.Model,
	templateBasePaths []string,
	testCfg *ackgenconfig.TestConfig,
) (*templateset.TemplateSet, error) {
	crds, err := m.GetCRDs()
	if err != nil {
		return nil, err
	}

	metaVars := m.MetaVars()

	ts := templateset.New(
		templateBasePaths,
		e2eIncludePaths,
		e2eCopyPaths,
		e2eFuncMap,
	)

	// Add service-level templates
	mainVars := &templateE2EMainVars{
		MetaVars:   metaVars,
		TestConfig: testCfg,
	}

	serviceTemplates := map[string]string{
		"test/e2e-go/main_test.go": "test/e2e/main_test.go.tpl",
		"Makefile":                 "test/e2e/Makefile.tpl",
	}
	for outPath, tplPath := range serviceTemplates {
		if err := ts.Add(outPath, tplPath, mainVars); err != nil {
			return nil, err
		}
	}

	// Add per-CRD test templates
	for _, crd := range crds {
		testsMap, ok := testCfg.Resources[crd.Names.Original]
		if !ok {
			continue
		}

		// Sort test names for deterministic output
		testNames := make([]string, 0, len(testsMap))
		for name := range testsMap {
			testNames = append(testNames, name)
		}
		sort.Strings(testNames)

		var testEntries []templateE2ETestVars
		for _, testName := range testNames {
			resourceCfg := testsMap[testName]
			if resourceCfg.Skip {
				continue
			}

			fieldAssignments := resolveFieldAssignments(crd, resourceCfg, testCfg)
			updateAssignments := resolveUpdateAssignments(crd, resourceCfg)
			hasUpdate := len(updateAssignments) > 0

			createWait := resourceCfg.GetCreateWait()
			waitDuration := fmt.Sprintf("%d*time.Second", createWait)

			// Resolve verification calls
			createVerifyCalls := resolveVerifyCalls(crd, resourceCfg.CreateValues)
			updateVerifyCalls := resolveVerifyCalls(crd, resourceCfg.UpdateValues)
			hasVerification := len(createVerifyCalls) > 0

			testEntries = append(testEntries, templateE2ETestVars{
				TestName:           testName,
				TestNameCamel:      snakeToCamel(testName),
				FieldAssignments:   fieldAssignments,
				UpdateAssignments:  updateAssignments,
				HasUpdate:          hasUpdate,
				CreateWaitDuration: waitDuration,
				HasVerification:    hasVerification,
				CreateVerifyCalls:  createVerifyCalls,
				UpdateVerifyCalls:  updateVerifyCalls,
			})
		}

		if len(testEntries) == 0 {
			continue
		}

		fileVars := &templateE2EFileVars{
			MetaVars: metaVars,
			CRD:      crd,
			Tests:    testEntries,
		}

		outPath := filepath.Join("test/e2e-go", fmt.Sprintf("resource_%s_test.go", crd.Names.Snake))
		tplPath := "test/e2e/resource_test.go.tpl"
		if err := ts.Add(outPath, tplPath, fileVars); err != nil {
			return nil, err
		}
	}

	return ts, nil
}

// resolveVerifyCalls builds verification calls for the given field values by
// finding the appropriate read operation for each field.
func resolveVerifyCalls(crd *ackmodel.CRD, fieldValues map[string]interface{}) []VerifyCall {
	if len(fieldValues) == 0 {
		return nil
	}

	// Group fields by their read operation
	type opGroup struct {
		fg     *ackmodel.FieldGroupOperation
		fields map[string]interface{}
	}
	groups := make(map[string]*opGroup)

	for fieldName, val := range fieldValues {
		fg := findReadOpForField(crd, fieldName)
		if fg == nil {
			continue
		}
		key := fg.OperationID
		if _, ok := groups[key]; !ok {
			groups[key] = &opGroup{fg: fg, fields: make(map[string]interface{})}
		}
		groups[key].fields[fieldName] = val
	}

	// Sort by operation name for deterministic output
	opNames := make([]string, 0, len(groups))
	for name := range groups {
		opNames = append(opNames, name)
	}
	sort.Strings(opNames)

	var calls []VerifyCall
	for _, opName := range opNames {
		group := groups[opName]
		fg := group.fg

		identAssignments := buildIdentifierAssignments(fg)
		assertions := buildAssertionsFromShape(fg, group.fields)

		if len(assertions) == 0 {
			continue
		}

		calls = append(calls, VerifyCall{
			OperationName:         fg.OperationID,
			InputTypeName:         fg.OperationID + "Input",
			IdentifierAssignments: identAssignments,
			Assertions:            assertions,
		})
	}

	return calls
}

// findReadOpForField finds the FieldGroupOperation that reads the given field.
// Search order:
// 1. ReadFieldGroups (explicitly configured read operations)
// 2. Infer Get* from UpdateFieldGroups Put* operations
// 3. Infer Get* from field's From.Operation (Put*/Update* → Get*/Describe*)
func findReadOpForField(crd *ackmodel.CRD, fieldName string) *ackmodel.FieldGroupOperation {
	// 1. Check ReadFieldGroups
	for _, fg := range crd.ReadFieldGroups {
		for _, f := range fg.PayloadFields {
			if f.Names.Camel == fieldName {
				return fg
			}
		}
	}

	// 2. Infer from UpdateFieldGroups: Put* → Get*
	for _, fg := range crd.UpdateFieldGroups {
		for _, f := range fg.PayloadFields {
			if f.Names.Camel != fieldName {
				continue
			}
			getOpID := inferGetFromPut(fg.OperationID)
			if getOpID == "" {
				continue
			}
			op, found := crd.GetSDKAPI().API.Operations[getOpID]
			if !found {
				continue
			}
			return &ackmodel.FieldGroupOperation{
				OpType:           ackmodel.FieldGroupOpTypeRead,
				Kind:             ackmodel.FieldGroupKindDirect,
				OperationID:      getOpID,
				Names:            names.New(getOpID),
				Operation:        op,
				IdentifierFields: fg.IdentifierFields,
				PayloadFields:    fg.PayloadFields,
				Config:           fg.Config,
			}
		}
	}

	// 3. Fallback: use the field's From.Operation config to infer the Get operation.
	// This handles controllers that use custom update methods but still have
	// fields configured with from.operation (e.g., S3's committed generator.yaml).
	field, found := crd.Fields[fieldName]
	if !found {
		return nil
	}
	if field.FieldConfig == nil || field.FieldConfig.From == nil {
		return nil
	}
	putOpID := field.FieldConfig.From.Operation
	getOpID := inferGetFromPut(putOpID)
	if getOpID == "" {
		return nil
	}
	op, found := crd.GetSDKAPI().API.Operations[getOpID]
	if !found {
		return nil
	}
	// Build identifier fields from the primary key
	var identifierFields []*ackmodel.Field
	pkField, err := crd.GetPrimaryKeyField()
	if err == nil && pkField != nil {
		identifierFields = []*ackmodel.Field{pkField}
	}
	return &ackmodel.FieldGroupOperation{
		OpType:           ackmodel.FieldGroupOpTypeRead,
		Kind:             ackmodel.FieldGroupKindDirect,
		OperationID:      getOpID,
		Names:            names.New(getOpID),
		Operation:        op,
		IdentifierFields: identifierFields,
		PayloadFields:    []*ackmodel.Field{field},
	}
}

// inferGetFromPut derives the Get operation name from a Put operation name.
func inferGetFromPut(putOpID string) string {
	if strings.HasPrefix(putOpID, "Put") {
		return "Get" + putOpID[3:]
	}
	if strings.HasPrefix(putOpID, "Update") {
		return "Describe" + putOpID[6:]
	}
	return ""
}

// buildIdentifierAssignments generates Go expressions for each identifier field
// in the read operation's input.
func buildIdentifierAssignments(fg *ackmodel.FieldGroupOperation) map[string]string {
	assignments := make(map[string]string)
	if fg.Operation == nil || fg.Operation.InputRef.Shape == nil {
		return assignments
	}

	inputShape := fg.Operation.InputRef.Shape
	for _, memberName := range inputShape.MemberNames() {
		// Only include identifier fields (skip optional fields like ExpectedBucketOwner)
		isIdentifier := false
		for _, idf := range fg.IdentifierFields {
			if idf.Names.Camel == memberName || idf.Names.Original == memberName {
				isIdentifier = true
				break
			}
		}
		if !isIdentifier {
			// Check if this member is required (it might be the bucket name with a different name)
			memberRef := inputShape.MemberRefs[memberName]
			if memberRef == nil || !inputShape.IsRequired(memberName) {
				continue
			}
		}

		memberRef := inputShape.MemberRefs[memberName]
		if memberRef == nil || memberRef.Shape == nil {
			continue
		}
		switch memberRef.Shape.Type {
		case "string":
			assignments[memberName] = "aws.String(name)"
		default:
			assignments[memberName] = "aws.String(name)"
		}
	}

	return assignments
}

// buildAssertionsFromShape generates assertions by matching field values against
// the SDK output shape members.
func buildAssertionsFromShape(fg *ackmodel.FieldGroupOperation, fieldValues map[string]interface{}) []VerifyAssertion {
	if fg.Operation == nil || fg.Operation.OutputRef.Shape == nil {
		return nil
	}

	outputShape := fg.Operation.OutputRef.Shape
	var assertions []VerifyAssertion

	for fieldName, val := range fieldValues {
		// The field value from testconfig may be a nested struct (e.g., {Status: "Enabled"})
		// or a scalar. We need to match it against the output shape.
		mapVal, isMap := normalizeToMap(val)
		if isMap {
			// Nested struct: each key in the map corresponds to an output shape member
			subAssertions := buildNestedAssertions(outputShape, "", mapVal, fieldName)
			assertions = append(assertions, subAssertions...)
		} else {
			// Scalar: find the output member that corresponds to this field
			member := findOutputMemberForField(outputShape, fieldName)
			if member != "" {
				memberRef := outputShape.MemberRefs[member]
				expr := goAssertionExpr(memberRef, val)
				assertions = append(assertions, VerifyAssertion{
					ResponsePath: member,
					Expected:     expr,
					DerefFunc:    derefFuncForShape(memberRef),
					Message:      fmt.Sprintf("%s mismatch", fieldName),
				})
			}
		}
	}

	sort.Slice(assertions, func(i, j int) bool {
		return assertions[i].ResponsePath < assertions[j].ResponsePath
	})
	return assertions
}

// buildNestedAssertions handles map values from testconfig (e.g., {Status: "Enabled"})
// and generates assertions for each leaf value.
func buildNestedAssertions(shape *awssdkmodel.Shape, pathPrefix string, vals map[string]interface{}, contextFieldName string) []VerifyAssertion {
	var assertions []VerifyAssertion

	// If none of the vals keys match the shape's direct members, check if the
	// shape has a single struct wrapper member we should descend into.
	// E.g., GetPublicAccessBlockOutput has only "PublicAccessBlockConfiguration"
	// but the test values have "BlockPublicACLs" etc.
	if !anyMemberMatches(shape, vals) {
		if wrapperRef, wrapperName := findSingleStructMember(shape); wrapperRef != nil {
			wrapperPath := wrapperName
			if pathPrefix != "" {
				wrapperPath = pathPrefix + "." + wrapperName
			}
			return buildNestedAssertions(wrapperRef.Shape, wrapperPath, vals, contextFieldName)
		}
	}

	for key, val := range vals {
		memberRef, memberName := findMemberByName(shape, key)
		if memberRef == nil {
			continue
		}

		fullPath := memberName
		if pathPrefix != "" {
			fullPath = pathPrefix + "." + memberName
		}

		subMap, isSubMap := normalizeToMap(val)
		if isSubMap && memberRef.Shape != nil && memberRef.Shape.Type == "structure" {
			nested := buildNestedAssertions(memberRef.Shape, fullPath, subMap, contextFieldName+"."+key)
			assertions = append(assertions, nested...)
		} else {
			expr := goAssertionExpr(memberRef, val)
			assertions = append(assertions, VerifyAssertion{
				ResponsePath: fullPath,
				Expected:     expr,
				DerefFunc:    derefFuncForShape(memberRef),
				Message:      fmt.Sprintf("%s mismatch", contextFieldName+"."+key),
			})
		}
	}

	return assertions
}

// anyMemberMatches returns true if any key in vals matches a member of the shape.
func anyMemberMatches(shape *awssdkmodel.Shape, vals map[string]interface{}) bool {
	for key := range vals {
		if ref, _ := findMemberByName(shape, key); ref != nil {
			return true
		}
	}
	return false
}

// findSingleStructMember returns the sole struct-type member of a shape, if
// there is exactly one (ignoring metadata/result fields). Returns nil otherwise.
func findSingleStructMember(shape *awssdkmodel.Shape) (*awssdkmodel.ShapeRef, string) {
	if shape == nil {
		return nil, ""
	}
	var found *awssdkmodel.ShapeRef
	var foundName string
	for name, ref := range shape.MemberRefs {
		if name == "ResultMetadata" {
			continue
		}
		if ref.Shape != nil && ref.Shape.Type == "structure" {
			if found != nil {
				return nil, ""
			}
			found = ref
			foundName = name
		}
	}
	return found, foundName
}

// findMemberByName looks up a shape member by camelCase or original name.
func findMemberByName(shape *awssdkmodel.Shape, name string) (*awssdkmodel.ShapeRef, string) {
	if shape == nil {
		return nil, ""
	}
	// Direct match
	if ref, ok := shape.MemberRefs[name]; ok {
		return ref, name
	}
	// Case-insensitive match
	lower := strings.ToLower(name)
	for memberName, ref := range shape.MemberRefs {
		if strings.ToLower(memberName) == lower {
			return ref, memberName
		}
	}
	return nil, ""
}

// findOutputMemberForField finds the output shape member that corresponds to
// a CRD field name.
func findOutputMemberForField(outputShape *awssdkmodel.Shape, fieldName string) string {
	if outputShape == nil {
		return ""
	}
	// Direct match
	if _, ok := outputShape.MemberRefs[fieldName]; ok {
		return fieldName
	}
	// Case-insensitive match
	lower := strings.ToLower(fieldName)
	for memberName := range outputShape.MemberRefs {
		if strings.ToLower(memberName) == lower {
			return memberName
		}
	}
	return ""
}

// goAssertionExpr generates a Go expression for an expected value based on the
// SDK shape reference type.
func goAssertionExpr(ref *awssdkmodel.ShapeRef, val interface{}) string {
	if ref == nil || ref.Shape == nil {
		return fmt.Sprintf("%q", fmt.Sprintf("%v", val))
	}
	shape := ref.Shape

	switch shape.Type {
	case "string", "character":
		s := fmt.Sprintf("%v", val)
		// If the shape has an enum, use the typed enum
		if len(shape.Enum) > 0 {
			typeName := shape.ShapeName
			return fmt.Sprintf("svcsdktypes.%s(%q)", typeName, s)
		}
		return fmt.Sprintf("%q", s)
	case "boolean", "primitiveBoolean":
		return fmt.Sprintf("%t", val)
	case "byte", "short", "integer", "primitiveInteger":
		switch v := val.(type) {
		case float64:
			return fmt.Sprintf("int32(%d)", int32(v))
		default:
			return fmt.Sprintf("int32(%v)", v)
		}
	case "long":
		switch v := val.(type) {
		case float64:
			return fmt.Sprintf("int64(%d)", int64(v))
		default:
			return fmt.Sprintf("int64(%v)", v)
		}
	case "float", "double":
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%q", fmt.Sprintf("%v", val))
	}
}

// derefFuncForShape returns the aws.ToX function name needed to dereference
// pointer-typed SDK struct members, or "" if no dereference is needed.
func derefFuncForShape(ref *awssdkmodel.ShapeRef) string {
	if ref == nil || ref.Shape == nil {
		return ""
	}
	switch ref.Shape.Type {
	case "boolean", "primitiveBoolean":
		return "aws.ToBool"
	case "string", "character":
		if len(ref.Shape.Enum) > 0 {
			return ""
		}
		return "aws.ToString"
	case "byte", "short", "integer", "primitiveInteger":
		return "aws.ToInt32"
	case "long":
		return "aws.ToInt64"
	case "float":
		return "aws.ToFloat32"
	case "double":
		return "aws.ToFloat64"
	default:
		return ""
	}
}

// normalizeToMap converts a value to map[string]interface{} if it's a map type.
func normalizeToMap(val interface{}) (map[string]interface{}, bool) {
	switch v := val.(type) {
	case map[string]interface{}:
		return v, true
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, vv := range v {
			m[fmt.Sprintf("%v", k)] = vv
		}
		return m, true
	}
	return nil, false
}

// snakeToCamel converts a snake_case string to CamelCase.
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

// resolveFieldAssignments builds the map of Go code expressions for each field
// in the CRD Spec, following the value resolution order:
// 1. bootstrap_fields → bootstrapResources.X.Y
// 2. create_values → literal value
// 3. auto-derived identifier → name variable reference
// 4. auto-derived shape defaults (only for required fields)
// 5. omit optional fields with no explicit value
func resolveFieldAssignments(
	crd *ackmodel.CRD,
	resourceCfg ackgenconfig.TestResourceConfig,
	testCfg *ackgenconfig.TestConfig,
) map[string]string {
	assignments := make(map[string]string)

	for fieldName, field := range crd.SpecFields {
		goFieldName := field.Names.Camel

		// 1. Bootstrap fields (always included if configured)
		if path, ok := resourceCfg.BootstrapFields[goFieldName]; ok {
			assignments[goFieldName] = bootstrapFieldExpr(field, path)
			continue
		}

		// 2. Explicit create_values (always included if configured)
		if val, ok := resourceCfg.CreateValues[goFieldName]; ok {
			assignments[goFieldName] = goLiteralForValue(field, val)
			continue
		}

		// 3. Auto-derived: primary identifier fields get the random name
		if isPrimaryIdentifierField(crd, fieldName) {
			assignments[goFieldName] = identifierExpr(field)
			continue
		}

		// For optional fields with no explicit config, skip them
		if !field.IsRequired() {
			continue
		}

		// 4. Auto-derived from shape (required fields only)
		if expr := shapeDefaultExpr(field); expr != "" {
			assignments[goFieldName] = expr
			continue
		}

		// 5. Required field with no derivable value — placeholder
		assignments[goFieldName] = fmt.Sprintf(`new("TODO_%s")`, goFieldName)
	}

	return assignments
}

// resolveUpdateAssignments builds Go code expressions for update test fields.
func resolveUpdateAssignments(
	crd *ackmodel.CRD,
	resourceCfg ackgenconfig.TestResourceConfig,
) map[string]string {
	assignments := make(map[string]string)
	for goFieldName, val := range resourceCfg.UpdateValues {
		field := findFieldByCamelName(crd, goFieldName)
		if field == nil {
			continue
		}
		assignments[goFieldName] = goLiteralForValue(field, val)
	}
	return assignments
}

// findFieldByCamelName looks up a spec field by its Go struct name (Names.Camel).
func findFieldByCamelName(crd *ackmodel.CRD, camelName string) *ackmodel.Field {
	for _, field := range crd.SpecFields {
		if field.Names.Camel == camelName {
			return field
		}
	}
	return nil
}

// bootstrapFieldExpr generates a Go expression that reads from the bootstrapResources
// global variable. path is like "SharedVPC.SecurityGroupID".
func bootstrapFieldExpr(field *ackmodel.Field, path string) string {
	expr := "bootstrapResources." + path
	goType := field.GoType
	if strings.HasPrefix(goType, "*") {
		return fmt.Sprintf("new(%s)", expr)
	}
	if strings.HasPrefix(goType, "[]*string") {
		return fmt.Sprintf("[]*string{new(%s)}", expr)
	}
	return expr
}

// identifierExpr generates a Go expression for an identifier field that uses
// the test's random name variable.
func identifierExpr(field *ackmodel.Field) string {
	goType := field.GoType
	if strings.HasPrefix(goType, "*") {
		return "&name"
	}
	return "name"
}

// isPrimaryIdentifierField returns true if the field is the resource's primary
// identifier — either marked as is_primary_key in config, or matching strict
// naming patterns for common identifier fields.
func isPrimaryIdentifierField(crd *ackmodel.CRD, fieldName string) bool {
	// Check if field is explicitly marked as primary key
	field, ok := crd.SpecFields[fieldName]
	if ok && field.FieldConfig != nil && field.FieldConfig.IsPrimaryKey {
		return true
	}

	lower := strings.ToLower(fieldName)
	resourceLower := strings.ToLower(crd.Names.Original)

	// Strict matches: the field name IS the resource name + identifier suffix
	identifierPatterns := []string{
		resourceLower + "name",
		resourceLower + "id",
		resourceLower + "identifier",
	}
	for _, p := range identifierPatterns {
		if lower == p {
			return true
		}
	}

	// Exact match for the generic "Name" field (very common pattern)
	if lower == "name" {
		return true
	}

	return false
}

// shapeDefaultExpr generates a Go expression using shape metadata (enums, etc.)
func shapeDefaultExpr(field *ackmodel.Field) string {
	if field.ShapeRef == nil || field.ShapeRef.Shape == nil {
		return ""
	}
	shape := field.ShapeRef.Shape

	// Enum: use first value
	if len(shape.Enum) > 0 {
		goType := field.GoType
		if strings.HasPrefix(goType, "*") {
			return fmt.Sprintf(`new("%s")`, shape.Enum[0])
		}
		return fmt.Sprintf(`"%s"`, shape.Enum[0])
	}

	return ""
}

// goLiteralForValue converts a testconfig.yaml value to a Go code expression.
func goLiteralForValue(field *ackmodel.Field, val interface{}) string {
	goType := field.GoType

	switch v := val.(type) {
	case string:
		if strings.HasPrefix(goType, "*string") {
			return fmt.Sprintf("new(%q)", v)
		}
		return fmt.Sprintf("%q", v)
	case int:
		if strings.Contains(goType, "int64") {
			if strings.HasPrefix(goType, "*") {
				return fmt.Sprintf("new(%d)", v)
			}
			return fmt.Sprintf("int64(%d)", v)
		}
		if strings.HasPrefix(goType, "*") {
			return fmt.Sprintf("new(%d)", v)
		}
		return fmt.Sprintf("%d", v)
	case float64:
		// YAML numbers often decode as float64
		intVal := int64(v)
		if v == float64(intVal) {
			if strings.HasPrefix(goType, "*") {
				return fmt.Sprintf("new(%d)", intVal)
			}
			return fmt.Sprintf("int64(%d)", intVal)
		}
		if strings.HasPrefix(goType, "*") {
			return fmt.Sprintf("new(%g)", v)
		}
		return fmt.Sprintf("%g", v)
	case bool:
		if strings.HasPrefix(goType, "*") {
			return fmt.Sprintf("new(%t)", v)
		}
		return fmt.Sprintf("%t", v)
	case map[string]interface{}:
		return goLiteralForStruct(field, v)
	case map[interface{}]interface{}:
		normalized := make(map[string]interface{})
		for k, val := range v {
			normalized[fmt.Sprintf("%v", k)] = val
		}
		return goLiteralForStruct(field, normalized)
	default:
		return fmt.Sprintf("%#v", v)
	}
}

// goLiteralForStruct generates a Go struct literal from a map of field values.
func goLiteralForStruct(field *ackmodel.Field, vals map[string]interface{}) string {
	if field.ShapeRef == nil || field.ShapeRef.Shape == nil {
		return fmt.Sprintf("%#v", vals)
	}
	return goLiteralForShapeStruct(field.ShapeRef, vals)
}

// goLiteralForShapeStruct generates a struct literal for a given shape reference.
func goLiteralForShapeStruct(ref *awssdkmodel.ShapeRef, vals map[string]interface{}) string {
	shape := ref.Shape
	typeName := shape.ShapeName

	members := make([]string, 0)
	for memberName, memberRef := range shape.MemberRefs {
		goMemberName := names.New(memberName).Camel
		val, ok := vals[goMemberName]
		if !ok {
			val, ok = vals[memberName]
		}
		if !ok {
			continue
		}
		memberExpr := goLiteralForShapeRef(memberRef, val)
		members = append(members, fmt.Sprintf("%s: %s", goMemberName, memberExpr))
	}

	if len(members) == 1 {
		return fmt.Sprintf("&svcapitypes.%s{%s}", typeName, members[0])
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("&svcapitypes.%s{", typeName))
	for _, m := range members {
		sb.WriteString("\n\t" + m + ",")
	}
	sb.WriteString("\n}")
	return sb.String()
}

// goLiteralForShapeRef converts a value to a Go expression based on the AWS SDK shape.
func goLiteralForShapeRef(ref *awssdkmodel.ShapeRef, val interface{}) string {
	if ref == nil || ref.Shape == nil {
		return fmt.Sprintf("%#v", val)
	}
	shape := ref.Shape

	switch shape.Type {
	case "string", "character":
		s := fmt.Sprintf("%v", val)
		return fmt.Sprintf("new(%q)", s)
	case "boolean", "primitiveBoolean":
		return fmt.Sprintf("new(bool(%t))", val)
	case "byte", "short", "integer", "long", "primitiveInteger":
		switch v := val.(type) {
		case float64:
			return fmt.Sprintf("new(%d)", int64(v))
		default:
			return fmt.Sprintf("new(%v)", v)
		}
	case "float", "double":
		return fmt.Sprintf("new(%v)", val)
	case "structure", "union":
		mapVal, ok := val.(map[string]interface{})
		if !ok {
			if m, ok2 := val.(map[interface{}]interface{}); ok2 {
				mapVal = make(map[string]interface{})
				for k, v := range m {
					mapVal[fmt.Sprintf("%v", k)] = v
				}
			}
		}
		if mapVal == nil {
			return fmt.Sprintf("%#v", val)
		}
		return goLiteralForShapeStruct(ref, mapVal)
	default:
		return fmt.Sprintf("%#v", val)
	}
}
