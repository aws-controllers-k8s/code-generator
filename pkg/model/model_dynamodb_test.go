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

func TestDynamoDB_Table(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "dynamodb")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Table", crds)
	require.NotNil(crd)

	// The DynamoDB Table API has these operations:
	//
	// * CreateTable
	// * DeleteTable
	// * DescribeTable
	// * ListTables
	// * UpdateTable
	require.NotNil(crd.Ops)

	assert.NotNil(crd.Ops.Create)
	assert.NotNil(crd.Ops.Delete)
	assert.NotNil(crd.Ops.ReadOne)
	assert.NotNil(crd.Ops.ReadMany)
	assert.NotNil(crd.Ops.Update)

	assert.Nil(crd.Ops.GetAttributes)
	assert.Nil(crd.Ops.SetAttributes)

	// The DescribeTable operation has the following definition:
	//
	//    "DescribeTable":{
	//      "name":"DescribeTable",
	//      "http":{
	//        "method":"POST",
	//        "requestUri":"/"
	//      },
	//      "input":{"shape":"DescribeTableInput"},
	//      "output":{"shape":"DescribeTableOutput"},
	//      "errors":[
	//        {"shape":"ResourceNotFoundException"},
	//        {"shape":"InternalServerError"}
	//      ],
	//      "endpointdiscovery":{
	//      }
	//    },
	//
	// Where the ResourceNotFoundException shape looks like this:
	//
	//    "ResourceNotFoundException":{
	//      "type":"structure",
	//      "members":{
	//        "message":{"shape":"ErrorMessage"}
	//      },
	//      "exception":true
	//    },
	//
	//
	// Which does not indicate that the error is a 404 :( So, the logic in the
	// CRD.ExceptionCode(404) method needs to get its override from the
	// generate.yaml configuration file.
	assert.Equal("ResourceNotFoundException", crd.ExceptionCode(404))

	specFields := crd.SpecFields

	expSpecFieldCamel := []string{
		"AttributeDefinitions",
		"BillingMode",
		"DeletionProtectionEnabled",
		"GlobalSecondaryIndexes",
		"KeySchema",
		"LocalSecondaryIndexes",
		"OnDemandThroughput",
		"ProvisionedThroughput",
		"ResourcePolicy",
		"SSESpecification",
		"StreamSpecification",
		"TableClass",
		"TableName",
		"Tags",
		"WarmThroughput",
	}
	assert.Equal(expSpecFieldCamel, attrCamelNames(specFields))
}

func TestDynamoDB_CustomShape_ReplicasState(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "dynamodb", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-custom-shapes.yaml",
	})

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Table", crds)
	require.NotNil(crd)

	// Verify the ReplicaStates field exists 
	assert.Contains(crd.StatusFields, "ReplicaStates")
	replicasDescField := crd.StatusFields["ReplicaStates"]
	require.NotNil(replicasDescField)


	replicasStateShape := replicasDescField.ShapeRef.Shape.MemberRef.Shape
	require.NotNil(replicasStateShape)

	// Verify all the expected fields exist in the RepicasState shape
	expectedFields := []string{
		"RegionName",
		"RegionStatus",
		"RegionStatusDescription",
		"RegionStatusPercentProgress",
		"RegionInaccessibleDateTime",
	}

	for _, fieldName := range expectedFields {
		assert.Contains(replicasStateShape.MemberRefs, fieldName, "RepicasState shape is missing field: "+fieldName)
		field := replicasStateShape.MemberRefs[fieldName]
		assert.Equal("string", field.Shape.Type, "Field "+fieldName+" should be of type string")
	}
}
