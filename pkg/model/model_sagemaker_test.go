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

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestSageMaker_ARN_Field_Override(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sagemaker")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("DataQualityJobDefinition", crds)
	require.NotNil(crd)

	// The CreateDataQualityJobDefinition has the following definition:
	//
	//   "CreateDataQualityJobDefinition":{
	//   "name":"CreateDataQualityJobDefinition",
	//   "http":{
	//     "method":"POST",
	//     "requestUri":"/"
	//   },
	//   "input":{"shape":"CreateDataQualityJobDefinitionRequest"},
	//   "output":{"shape":"CreateDataQualityJobDefinitionResponse"},
	//   "errors":[
	//     {"shape":"ResourceLimitExceeded"},
	//     {"shape":"ResourceInUse"}
	//   ]
	// }
	//
	// Where the CreateDataQualityJobDefinitionResponse shape looks like this:
	//
	// "CreateDataQualityJobDefinitionResponse":{
	// 	"type":"structure",
	// 	"required":["JobDefinitionArn"],
	// 	"members":{
	// 	  "JobDefinitionArn":{"shape":"MonitoringJobDefinitionArn"}
	// 	}
	// }
	//
	// So, we expect that the logic in crd.IsPrimaryARNField() parses through
	// field config and identifies the JobDefinitionArn as the primaryARNField
	// for the resource
	assert.Equal(true, crd.IsPrimaryARNField("JobDefinitionArn"))

}

func TestSageMaker_Error_Prefix_Message(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sagemaker")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("TrainingJob", crds)
	require.NotNil(crd)

	require.NotNil(crd.Ops)

	assert.NotNil(crd.Ops.ReadOne)

	// "DescribeTrainingJob":{
	// 	"name":"DescribeTrainingJob",
	// 	"http":{
	// 	  "method":"POST",
	// 	  "requestUri":"/"
	// 	},
	// 	"input":{"shape":"DescribeTrainingJobRequest"},
	// 	"output":{"shape":"DescribeTrainingJobResponse"},
	// 	"errors":[
	// 	  {"shape":"ResourceNotFound"}
	// 	]
	//   },

	// Which does not indicate that the error is a 404 :( So, the logic in the
	// CRD.ExceptionCode(404) method needs to get its override from the
	// generate.yaml configuration file.
	assert.Equal("ValidationException", crd.ExceptionCode(404))

	// Validation Exception has prefix Requested resource not found.
	assert.Equal("&& strings.HasPrefix(awsErr.ErrorMessage(), \"Requested resource not found\") ", code.CheckExceptionMessage(crd.Config(), crd, 404))
}

func TestSageMaker_Error_Suffix_Message(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sagemaker")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("ModelPackageGroup", crds)
	require.NotNil(crd)

	require.NotNil(crd.Ops)
	assert.NotNil(crd.Ops.ReadOne)

	// 	"DescribeModelPackageGroup":{
	// 	"name":"DescribeModelPackageGroup",
	// 	"http":{
	// 	  "method":"POST",
	// 	  "requestUri":"/"
	// 	},
	// 	"input":{"shape":"DescribeModelPackageGroupInput"},
	// 	"output":{"shape":"DescribeModelPackageGroupOutput"}
	//   }

	// Does not list an error however a ValidationException can occur
	// Which does not indicate that the error is a 404 :( So, the logic in the
	// CRD.ExceptionCode(404) method needs to get its override from the
	// generate.yaml configuration file.
	assert.Equal("ValidationException", crd.ExceptionCode(404))

	// Validation Exception has suffix ModelPackageGroup arn:aws:sagemaker:/ does not exist
	assert.Equal("&& strings.HasSuffix(awsErr.ErrorMessage(), \"does not exist.\") ", code.CheckExceptionMessage(crd.Config(), crd, 404))
}

func TestSageMaker_RequeueOnSuccessSeconds(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sagemaker")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Endpoint", crds)
	require.NotNil(crd)

	// The CreateEndpoint has the following definition:
	//
	// "CreateEndpoint":{
	// 	"name":"CreateEndpoint",
	// 	"http":{
	// 	  "method":"POST",
	// 	  "requestUri":"/"
	// 	},
	// 	"input":{"shape":"CreateEndpointInput"},
	// 	"output":{"shape":"CreateEndpointOutput"},
	// 	"errors":[
	// 	  {"shape":"ResourceLimitExceeded"}
	// 	]
	//   }
	//
	// Where the CreateEndpointOutput shape looks like this:
	//
	// "CreateEndpointOutput":{
	// 	"type":"structure",
	// 	"required":["EndpointArn"],
	// 	"members":{
	// 	  "EndpointArn":{"shape":"EndpointArn"}
	// 	}
	//   }
	//
	// So, we expect that crd.ReconcileRequeuOnSuccessSeconds() returns the requeue
	// duration specified in the config file
	assert.Equal(10, crd.ReconcileRequeuOnSuccessSeconds())
}

func TestSageMaker_RequeueOnSuccessSeconds_Default(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sagemaker")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("DataQualityJobDefinition", crds)
	require.NotNil(crd)

	// The CreateDataQualityJobDefinition has the following definition:
	//
	//   "CreateDataQualityJobDefinition":{
	//   "name":"CreateDataQualityJobDefinition",
	//   "http":{
	//     "method":"POST",
	//     "requestUri":"/"
	//   },
	//   "input":{"shape":"CreateDataQualityJobDefinitionRequest"},
	//   "output":{"shape":"CreateDataQualityJobDefinitionResponse"},
	//   "errors":[
	//     {"shape":"ResourceLimitExceeded"},
	//     {"shape":"ResourceInUse"}
	//   ]
	// }
	//
	// Where the CreateDataQualityJobDefinitionResponse shape looks like this:
	//
	// "CreateDataQualityJobDefinitionResponse":{
	// 	"type":"structure",
	// 	"required":["JobDefinitionArn"],
	// 	"members":{
	// 	  "JobDefinitionArn":{"shape":"MonitoringJobDefinitionArn"}
	// 	}
	// }
	//
	// So, we expect that crd.ReconcileRequeuOnSuccessSeconds() returns the default
	// requeue duration of 0 because it is not specified in the config file
	assert.Equal(0, crd.ReconcileRequeuOnSuccessSeconds())

}

func TestSageMaker_TrialComponent_CustomNestedField_MapOfStruct(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "sagemaker", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-custom-nested-fields.yaml",
	})

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("TrialComponent", crds)
	require.NotNil(crd)

	// InputArtifacts is a map-of-struct field (map[string]*TrialComponentArtifact).
	// The generator config adds a custom nested field "InputArtifacts.CustomTag"
	// of type *string. Verify it was injected into the map value's struct shape.
	inputArtifactsField := crd.Fields["InputArtifacts"]
	require.NotNil(inputArtifactsField)
	require.NotNil(inputArtifactsField.ShapeRef)
	require.NotNil(inputArtifactsField.ShapeRef.Shape)
	assert.Equal("map", inputArtifactsField.ShapeRef.Shape.Type)

	// The map value shape should be a structure
	valueShape := inputArtifactsField.ShapeRef.Shape.ValueRef.Shape
	require.NotNil(valueShape)
	assert.Equal("structure", valueShape.Type)

	// Verify the custom nested field "CustomTag" was added to the map value's struct
	require.Contains(valueShape.MemberRefs, "CustomTag")
	customTagRef := valueShape.MemberRefs["CustomTag"]
	require.NotNil(customTagRef)
	require.NotNil(customTagRef.Shape)
	assert.Equal("string", customTagRef.Shape.Type)

	// Verify the original members are still present
	assert.Contains(valueShape.MemberRefs, "Value")
	assert.Contains(valueShape.MemberRefs, "MediaType")

	// Verify the custom nested field also appears as a MemberField on the CRD field
	require.NotNil(inputArtifactsField.MemberFields)
	customTagMemberField := inputArtifactsField.MemberFields["CustomTag"]
	require.NotNil(customTagMemberField)
}
