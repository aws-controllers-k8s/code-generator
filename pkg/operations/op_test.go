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

package operations_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestGetOpTypeAndResourceNameFromOpID(t *testing.T) {
	assert := assert.New(t)

	g := testutil.NewModelForService(t, "s3")

	tests := []struct {
		opID       string
		expOpType  operations.OpType
		expResName string
	}{
		{
			"CreateTopic",
			operations.OpTypeCreate,
			"Topic",
		},
		{
			"CreateOrUpdateTopic",
			operations.OpTypeReplace,
			"Topic",
		},
		{
			"CreateBatchTopics",
			operations.OpTypeCreateBatch,
			"Topic",
		},
		{
			"CreateBatchTopic",
			operations.OpTypeCreateBatch,
			"Topic",
		},
		{
			"BatchCreateTopics",
			operations.OpTypeCreateBatch,
			"Topic",
		},
		{
			"BatchCreateTopic",
			operations.OpTypeCreateBatch,
			"Topic",
		},
		{
			"CreateTopics",
			operations.OpTypeCreateBatch,
			"Topic",
		},
		{
			"DescribeEC2Instances",
			operations.OpTypeList,
			"EC2Instance",
		},
		{
			"DescribeEC2Instance",
			operations.OpTypeGet,
			"EC2Instance",
		},
		{
			"UpdateTopic",
			operations.OpTypeUpdate,
			"Topic",
		},
		{
			"DeleteTopic",
			operations.OpTypeDelete,
			"Topic",
		},
		{
			"DescribeInstances",
			operations.OpTypeList,
			"Instance",
		},
		{
			"ListDeploymentGroups",
			operations.OpTypeList,
			"DeploymentGroup",
		},
		{
			"GetDeployment",
			operations.OpTypeGet,
			"Deployment",
		},
		{
			"PauseEC2Instance",
			operations.OpTypeUnknown,
			"PauseEC2Instance",
		},
		// Heuristic should incorrectly parse DhcpOptions ops
		// due to resource not being in s3's generator config
		{
			"CreateDhcpOptions",
			model.OpTypeCreateBatch,
			"DhcpOption",
		},
		{
			"DescribeDhcpOptions",
			model.OpTypeList,
			"DhcpOption",
		},
	}
	for _, test := range tests {
		ot, resName := model.GetOpTypeAndResourceNameFromOpID(test.opID, g.GetConfig())
		assert.Equal(test.expOpType, ot, test.opID)
		assert.Equal(test.expResName, resName, test.opID)
	}
}

func TestGetOpTypeAndResourceNameFromOpID_PluralSingular(t *testing.T) {
	assert := assert.New(t)

	g := testutil.NewModelForService(t, "ec2")

	tests := []struct {
		opID       string
		expOpType  model.OpType
		expResName string
	}{
		{
			"CreateDhcpOptions",
			model.OpTypeCreate,
			"DhcpOptions",
		},
		{
			"DescribeDhcpOptions",
			model.OpTypeList,
			"DhcpOptions",
		},
		{
			"DeleteDhcpOptions",
			model.OpTypeDelete,
			"DhcpOptions",
		},
	}
	for _, test := range tests {
		ot, resName := model.GetOpTypeAndResourceNameFromOpID(test.opID, g.GetConfig())
		assert.Equal(test.expOpType, ot, test.opID)
		assert.Equal(test.expResName, resName, test.opID)
	}
}
