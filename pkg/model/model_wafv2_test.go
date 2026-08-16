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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
