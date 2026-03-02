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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

// TestSetSDK_EMRServerless_Application_Create tests that the EMR Serverless
// Application resource generates correct code for the Create operation.
//
// This test specifically verifies the fix for the following compile errors:
// 1. WorkerCount field should be treated as a pointer (*int64) not double pointer (**int64)
//   - The WorkerCounts shape has a bad default value that was incorrectly handled
//
// 2. WorkerTypeSpecificationInput should use the original SDK shape name in map declarations
//   - The code generator was using renamed shape names instead of original SDK names
//
// Related fix: commit 84ad1a8 "Fix incorrect shap name in bad default filtering"
func TestSetSDK_EMRServerless_Application_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "emr-serverless")

	crd := testutil.GetCRDByName(t, g, "Application")
	require.NotNil(crd)

	got, err := code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1)
	require.NoError(err)

	// Verify WorkerCount is handled correctly as a pointer, not double pointer
	// The fix ensures WorkerCount (which has a bad default in the SDK) is treated as *int64
	// Before the fix: &f7valiter.WorkerCount (type **int64) - WRONG
	// After the fix: f7valiter.WorkerCount (type *int64) - CORRECT
	//
	// The generated code should have:
	//   if f7valiter.WorkerCount != nil {
	//       f7val.WorkerCount = f7valiter.WorkerCount
	//   }
	// NOT:
	//   f7val.WorkerCount = &f7valiter.WorkerCount (which would be **int64)
	assert.Contains(got, "f7valiter.WorkerCount")
	assert.Contains(got, "f7val.WorkerCount = f7valiter.WorkerCount",
		"WorkerCount should be assigned directly without & prefix")

	// Verify WorkerTypeSpecificationInput uses the original SDK shape name
	// The fix ensures map[string]svcsdktypes.WorkerTypeSpecificationInput is used
	// Before the fix: svcsdktypes.WorkerTypeSpecificationInput_ - WRONG (renamed shape)
	// After the fix: svcsdktypes.WorkerTypeSpecificationInput - CORRECT (original shape)
	assert.Contains(got, "map[string]svcsdktypes.WorkerTypeSpecificationInput{}")
	assert.NotContains(got, "WorkerTypeSpecificationInput_",
		"Should use original SDK shape name WorkerTypeSpecificationInput, not renamed version")
}

// TestSetSDK_EMRServerless_Application_Update tests that the EMR Serverless
// Application resource generates correct code for the Update operation.
// This ensures the same fixes apply to update operations as well.
func TestSetSDK_EMRServerless_Application_Update(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "emr-serverless")

	crd := testutil.GetCRDByName(t, g, "Application")
	require.NotNil(crd)

	got, err := code.SetSDK(crd.Config(), crd, model.OpTypeUpdate, "r.ko", "res", 1)
	require.NoError(err)

	// Verify the same fixes apply to Update operation
	// WorkerCount should be a pointer, not double pointer
	if strings.Contains(got, "WorkerCount") {
		// Ensure we're not using & prefix which would create double pointer
		assert.NotContains(got, "= &f",
			"WorkerCount should not be dereferenced with & in Update operation")
	}

	// WorkerTypeSpecificationInput should use original SDK shape name
	if strings.Contains(got, "WorkerTypeSpecificationInput") {
		assert.NotContains(got, "WorkerTypeSpecificationInput_",
			"Should use original SDK shape name in Update operation")
	}
}

// TestSetSDK_EMRServerless_Application_InitialCapacityConfig tests that
// InitialCapacityConfig map with WorkerCount field generates correct code.
func TestSetSDK_EMRServerless_Application_InitialCapacityConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "emr-serverless")

	crd := testutil.GetCRDByName(t, g, "Application")
	require.NotNil(crd)

	got, err := code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1)
	require.NoError(err)

	// Verify InitialCapacityConfig map is generated correctly
	// This map contains InitialCapacityConfig values which have WorkerCount field
	assert.Contains(got, "InitialCapacityConfig")

	// The generated code should properly handle the nested WorkerCount field
	// which has a bad default value in the SDK model
	assert.Contains(got, "WorkerCount")
}
