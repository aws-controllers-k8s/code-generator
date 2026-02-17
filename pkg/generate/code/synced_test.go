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

package code_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestSyncedLambdaFunction(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "lambda")

	crd := testutil.GetCRDByName(t, g, "Function")
	require.NotNil(crd)

	expectedSyncedConditions := `
	if r.ko.Status.State == nil {
		return false, nil
	}
	stateCandidates := []string{"AVAILABLE", "ACTIVE"}
	if !ackutil.InStrings(*r.ko.Status.State, stateCandidates) {
		return false, nil
	}
	if r.ko.Status.LastUpdateStatus == nil {
		return false, nil
	}
	lastUpdateStatusCandidates := []string{"AVAILABLE", "ACTIVE"}
	if !ackutil.InStrings(*r.ko.Status.LastUpdateStatus, lastUpdateStatusCandidates) {
		return false, nil
	}
	if r.ko.Status.CodeSize == nil {
		return false, nil
	}
	codeSizeCandidates := []int{1, 2}
	if !ackutil.InStrings(*r.ko.Status.CodeSize, codeSizeCandidates) {
		return false, nil
	}
`
	got, err := code.ResourceIsSynced(
		crd.Config(), crd, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal(expectedSyncedConditions, got)
}

func TestSyncedDynamodbTable(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "dynamodb")

	crd := testutil.GetCRDByName(t, g, "Table")
	require.NotNil(crd)

	expectedSyncedConditions := `
	if r.ko.Status.TableStatus == nil {
		return false, nil
	}
	tableStatusCandidates := []string{"AVAILABLE", "ACTIVE"}
	if !ackutil.InStrings(*r.ko.Status.TableStatus, tableStatusCandidates) {
		return false, nil
	}
	if r.ko.Spec.ProvisionedThroughput == nil {
		return false, nil
	}
	if r.ko.Spec.ProvisionedThroughput.ReadCapacityUnits == nil {
		return false, nil
	}
	provisionedThroughputCandidates := []int{0, 10}
	if !ackutil.InStrings(*r.ko.Spec.ProvisionedThroughput.ReadCapacityUnits, provisionedThroughputCandidates) {
		return false, nil
	}
	if r.ko.Status.ItemCount == nil {
		return false, nil
	}
	itemCountCandidates := []int{0}
	if !ackutil.InStrings(*r.ko.Status.ItemCount, itemCountCandidates) {
		return false, nil
	}
`
	got, err := code.ResourceIsSynced(
		crd.Config(), crd, "r.ko", 1,
	)
	require.NoError(err)
	assert.Equal(expectedSyncedConditions, got)
}
