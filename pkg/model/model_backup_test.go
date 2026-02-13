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

package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestBackup_BackupPlan_InputWrapperFieldPath(t *testing.T) {
	// CreateBackupPlanInput has a BackupPlan struct that wraps
	// the actual plan fields (BackupPlanName, Rules, etc.).
	// The input_wrapper_field_path config flattens these fields
	// into the CRD Spec so users don't need to nest them.
	//
	// Config used:
	//   operations:
	//     CreateBackupPlan:
	//       input_wrapper_field_path: BackupPlan
	//     UpdateBackupPlan:
	//       input_wrapper_field_path: BackupPlan
	//     GetBackupPlan:
	//       output_wrapper_field_path: BackupPlan
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "backup")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("BackupPlan", crds)
	require.NotNil(crd)

	assert.Equal("BackupPlan", crd.Names.Camel)
	assert.Equal("backupPlan", crd.Names.CamelLower)
	assert.Equal("backup_plan", crd.Names.Snake)

	specFields := crd.SpecFields
	statusFields := crd.StatusFields

	// Verify that fields from the BackupPlan wrapper are flattened into Spec.
	// The wrapper contains: BackupPlanName (renamed to Name), Rules, AdvancedBackupSettings, ScanSettings
	// Fields outside the wrapper (BackupPlanTags, CreatorRequestId) are NOT included -
	// this is consistent with output_wrapper_field_path behavior.
	expSpecFieldCamel := []string{
		"AdvancedBackupSettings",
		"Name",
		"Rules",
		"ScanSettings",
	}
	assert.Equal(expSpecFieldCamel, attrCamelNames(specFields))

	// Verify status fields come from the output (GetBackupPlan response)
	// which also uses output_wrapper_field_path: BackupPlan
	expStatusFieldCamel := []string{
		"BackupPlanID",
		"CreationDate",
		"VersionID",
	}
	assert.Equal(expStatusFieldCamel, attrCamelNames(statusFields))

	// Verify the CRD has the expected operations configured
	require.NotNil(crd.Ops)
	assert.NotNil(crd.Ops.Create)
	assert.NotNil(crd.Ops.Update)
	assert.NotNil(crd.Ops.Delete)
	assert.NotNil(crd.Ops.ReadOne)

	// Verify input_wrapper_field_path is set for Create operation
	createInputWrapper := crd.GetInputWrapperFieldPath(crd.Ops.Create)
	require.NotNil(createInputWrapper)
	assert.Equal("BackupPlan", *createInputWrapper)

	// Verify input_wrapper_field_path is set for Update operation
	updateInputWrapper := crd.GetInputWrapperFieldPath(crd.Ops.Update)
	require.NotNil(updateInputWrapper)
	assert.Equal("BackupPlan", *updateInputWrapper)
}
