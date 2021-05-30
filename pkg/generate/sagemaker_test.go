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

package generate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestSageMaker_ARN_Field_Override(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "sagemaker", ackgenconfig.DefaultConfig)

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
