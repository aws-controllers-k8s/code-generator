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

package sdk_test

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
		expOpType  sdk.OpType
		expResName string
	}{
		{
			"CreateTopic",
			sdk.OpTypeCreate,
			"Topic",
		},
		{
			"CreateOrUpdateTopic",
			sdk.OpTypeReplace,
			"Topic",
		},
		{
			"CreateBatchTopics",
			sdk.OpTypeCreateBatch,
			"Topic",
		},
		{
			"CreateBatchTopic",
			sdk.OpTypeCreateBatch,
			"Topic",
		},
		{
			"BatchCreateTopics",
			sdk.OpTypeCreateBatch,
			"Topic",
		},
		{
			"BatchCreateTopic",
			sdk.OpTypeCreateBatch,
			"Topic",
		},
		{
			"CreateTopics",
			sdk.OpTypeCreateBatch,
			"Topic",
		},
		{
			"DescribeEC2Instances",
			sdk.OpTypeList,
			"EC2Instance",
		},
		{
			"DescribeEC2Instance",
			sdk.OpTypeGet,
			"EC2Instance",
		},
		{
			"UpdateTopic",
			sdk.OpTypeUpdate,
			"Topic",
		},
		{
			"DeleteTopic",
			sdk.OpTypeDelete,
			"Topic",
		},
		{
			"DescribeInstances",
			sdk.OpTypeList,
			"Instance",
		},
		{
			"ListDeploymentGroups",
			sdk.OpTypeList,
			"DeploymentGroup",
		},
		{
			"GetDeployment",
			sdk.OpTypeGet,
			"Deployment",
		},
		{
			"PauseEC2Instance",
			sdk.OpTypeUnknown,
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
