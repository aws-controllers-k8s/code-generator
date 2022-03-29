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
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestEKS_WithNestedReference(t *testing.T) {
	_ = assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "eks", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-nested-reference.yaml",
	})

	tds, err := g.GetTypeDefs()
	require.Nil(err)
	require.NotNil(tds)

	// ResourcesVpcConfig field's type is VpcConfigRequest
	var vpcConfigRequestTD *model.TypeDef

	for _, td := range tds {
		if td != nil && strings.EqualFold(td.Names.Original, "vpcConfigRequest") {
			vpcConfigRequestTD = td
			break
		}
	}
	assert.NotNil(t, vpcConfigRequestTD)
	securityGroupIdsAttr := vpcConfigRequestTD.GetAttributeIgnoreCase("SecurityGroupIds")
	securityGroupRefsAttr := vpcConfigRequestTD.GetAttributeIgnoreCase("SecurityGroupRefs")

	assert.Equal(t, "SecurityGroupIDs", securityGroupIdsAttr.Names.Camel)
	assert.Equal(t, "SecurityGroupRefs", securityGroupRefsAttr.Names.Camel)
	assert.Equal(t, "[]*ackv1alpha1.AWSResourceReferenceWrapper", securityGroupRefsAttr.GoType)
}
