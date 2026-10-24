// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	 http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package model_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

// TestWAFv2_SharedShape_CustomNestedFields verifies that custom nested fields
// can be injected into a shape (Statement) that is shared across multiple
// resources (RuleGroup and WebACL). Both resources reference the same
// underlying SDK Rules/Statement shape, so the second resource's injection
// must not fail with "member already exists".
func TestWAFv2_SharedShape_CustomNestedFields(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "wafv2")

	crds, err := g.GetCRDs()
	require.Nil(err)

	ruleGroupCRD := getCRDByName("RuleGroup", crds)
	require.NotNil(ruleGroupCRD)

	webACLCRD := getCRDByName("WebACL", crds)
	require.NotNil(webACLCRD)

	// Both resources inject custom nested fields (AndStatement, OrStatement,
	// NotStatement) into the shared Statement shape. Verify they were injected
	// successfully on both CRDs.
	for _, crd := range []struct {
		name string
		fields map[string]*struct{}
	}{
		{name: "RuleGroup"},
		{name: "WebACL"},
	} {
		c := getCRDByName(crd.name, crds)
		require.NotNil(c, "CRD %s not found", crd.name)

		rulesField := c.Fields["Rules"]
		require.NotNil(rulesField, "%s.Rules field not found", crd.name)

		// Rules is a list of Rule structs, each containing a Statement
		stmtField := rulesField.MemberFields["Statement"]
		require.NotNil(stmtField, "%s.Rules.Statement field not found", crd.name)

		// Verify the custom nested fields were injected
		assert.Contains(stmtField.MemberFields, "AndStatement",
			"%s.Rules.Statement should have AndStatement", crd.name)
		assert.Contains(stmtField.MemberFields, "OrStatement",
			"%s.Rules.Statement should have OrStatement", crd.name)
		assert.Contains(stmtField.MemberFields, "NotStatement",
			"%s.Rules.Statement should have NotStatement", crd.name)

		// Verify the injected fields have string type (as configured)
		andStmt := stmtField.MemberFields["AndStatement"]
		require.NotNil(andStmt)
		require.NotNil(andStmt.ShapeRef)
		require.NotNil(andStmt.ShapeRef.Shape)
		assert.Equal("string", andStmt.ShapeRef.Shape.Type,
			"%s.Rules.Statement.AndStatement should be string type", crd.name)
	}
}

// TestWAFv2_NestedReferences verifies that references configured on nested
// ARN fields generate the expected reference field names for RuleGroup and
// WebACL.
func TestWAFv2_NestedReferences(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "wafv2")

	crds, err := g.GetCRDs()
	require.NoError(err)

	testCases := []struct {
		resource  string
		fieldPath string
		refName   string
	}{
		{
			resource:  "RuleGroup",
			fieldPath: "Rules.Statement.IPSetReferenceStatement.ARN",
			refName:   "IPSetRef",
		},
		{
			resource:  "WebACL",
			fieldPath: "Rules.Statement.IPSetReferenceStatement.ARN",
			refName:   "IPSetRef",
		},
		{
			resource:  "WebACL",
			fieldPath: "Rules.Statement.RuleGroupReferenceStatement.ARN",
			refName:   "RuleGroupRef",
		},
	}

	for _, tc := range testCases {
		crd := getCRDByName(tc.resource, crds)
		require.NotNil(crd, "CRD %s not found", tc.resource)

		field := crd.Fields[tc.fieldPath]
		require.NotNil(field, "field %s.%s not found", tc.resource, tc.fieldPath)
		assert.True(field.HasReference(), "field %s.%s should have a reference", tc.resource, tc.fieldPath)

		refName, err := field.GetReferenceFieldName()
		require.NoError(err)
		assert.Equal(tc.refName, refName.Camel)
	}
}

func TestWAFv2_NestedFieldDocumentationOverride(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "wafv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile:     "generator.yaml",
			DocumentationConfigFile: "documentation.yaml",
		},
	)

	typeDefs, err := g.GetTypeDefs()
	require.NoError(err)

	var resourceARN *model.Attr
	for _, typeDef := range typeDefs {
		for _, attr := range typeDef.Attrs {
			if attr.Names.Original == "ResourceArn" && attr.Documentation != "" {
				resourceARN = attr
				break
			}
		}
	}
	require.NotNil(resourceARN)
	require.Equal(
		"\t// The Amazon Resource Name (ARN) of the resource to associate with the web ACL.",
		resourceARN.Documentation,
	)

	templatePath, err := filepath.Abs("../../templates")
	require.NoError(err)
	apiTemplates, err := ackgenerate.APIs(g, []string{templatePath})
	require.NoError(err)
	require.NoError(apiTemplates.Execute())

	typesFile, ok := apiTemplates.Executed()["types.go"]
	require.True(ok)
	require.Contains(typesFile.String(), resourceARN.Documentation)
}
