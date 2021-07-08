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

func TestAPIGatewayV2_GetTypeDefs(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	// There is an "Api" Shape that is a struct that is an element of the
	// GetApis Operation. Its name conflicts with the CRD called API and thus
	// we need to check that its cleaned name is set to API_SDK (the _SDK
	// suffix is appended to the type name to avoid the conflict with
	// CRD-specific structs.
	tdef := testutil.GetTypeDefByName(t, g, "Api")
	require.NotNil(tdef)

	assert.Equal("API_SDK", tdef.Names.Camel)
}

func TestAPIGatewayV2_Api(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Api", crds)
	require.NotNil(crd)

	assert.Equal("API", crd.Names.Camel)
	assert.Equal("api", crd.Names.CamelLower)
	assert.Equal("api", crd.Names.Snake)

	assert.NotNil(crd.SpecFields["Name"])
	assert.NotNil(crd.SpecFields["ProtocolType"])
	// Body, Basepath and FailOnWarnings fields from ImportApi operation should get added to APISpec
	assert.NotNil(crd.SpecFields["Body"])
	assert.NotNil(crd.SpecFields["Basepath"])
	assert.NotNil(crd.SpecFields["FailOnWarnings"])

	// The required property should get overriden for Name and ProtocolType fields.
	assert.False(crd.SpecFields["Name"].IsRequired())
	assert.False(crd.SpecFields["ProtocolType"].IsRequired())
}

func TestAPIGatewayV2_Route(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Route", crds)
	require.NotNil(crd)

	assert.Equal("Route", crd.Names.Camel)
	assert.Equal("route", crd.Names.CamelLower)
	assert.Equal("route", crd.Names.Snake)

	// The GetRoute operation has the following definition:
	//
	//    "GetRoute" : {
	//      "name" : "GetRoute",
	//      "http" : {
	//        "method" : "GET",
	//        "requestUri" : "/v2/apis/{apiId}/routes/{routeId}",
	//        "responseCode" : 200
	//      },
	//      "input" : {
	//        "shape" : "GetRouteRequest"
	//      },
	//      "output" : {
	//        "shape" : "GetRouteResult"
	//      },
	//      "errors" : [ {
	//        "shape" : "NotFoundException"
	//      }, {
	//        "shape" : "TooManyRequestsException"
	//      } ]
	//    },
	//
	// Where the NotFoundException shape looks like this:
	//
	//    "NotFoundException" : {
	//      "type" : "structure",
	//      "members" : {
	//        "Message" : {
	//          "shape" : "__string",
	//          "locationName" : "message"
	//        },
	//        "ResourceType" : {
	//          "shape" : "__string",
	//          "locationName" : "resourceType"
	//        }
	//      },
	//      "exception" : true,
	//      "error" : {
	//        "httpStatusCode" : 404
	//      }
	//    },
	//
	// Which indicates that the error is a 404 and is our NotFoundException
	// code but the "code" attribute of the ErrorInfo struct is empty so
	// instead of returning a blank string, we need to use the name of the
	// shape itself...
	assert.Equal("NotFoundException", crd.ExceptionCode(404))

	// The APIGatewayV2 Route API has CRUD+L operations:
	//
	// * CreateRoute
	// * DeleteRoute
	// * UpdateRoute
	// * GetRoute
	// * GetRoutes
	require.NotNil(crd.Ops)

	assert.NotNil(crd.Ops.Create)
	assert.NotNil(crd.Ops.Delete)
	assert.NotNil(crd.Ops.Update)
	assert.NotNil(crd.Ops.ReadOne)
	assert.NotNil(crd.Ops.ReadMany)

	// And no separate get/set attributes calls.
	assert.Nil(crd.Ops.GetAttributes)
	assert.Nil(crd.Ops.SetAttributes)

	specFields := crd.SpecFields
	statusFields := crd.StatusFields

	expSpecFieldCamel := []string{
		"APIID",
		"APIKeyRequired",
		"AuthorizationScopes",
		"AuthorizationType",
		"AuthorizerID",
		"ModelSelectionExpression",
		"OperationName",
		"RequestModels",
		"RequestParameters",
		"RouteKey",
		"RouteResponseSelectionExpression",
		"Target",
	}
	assert.Equal(expSpecFieldCamel, attrCamelNames(specFields))

	expStatusFieldCamel := []string{
		"APIGatewayManaged",
		"RouteID",
	}
	assert.Equal(expStatusFieldCamel, attrCamelNames(statusFields))
}
