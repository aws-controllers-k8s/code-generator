// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
// http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package code_test

import (
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/crossplane"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func Test_LateInitializeReadOne(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "apigatewayv2", crossplane.DefaultConfig)

	crd := testutil.GetCRDByName(t, g, "Route")
	require.NotNil(crd)

	expected := `cr.Spec.ForProvider.APIKeyRequired = li.LateInitializeBoolPtr(cr.Spec.ForProvider.APIKeyRequired, resp.ApiKeyRequired)
if len(resp.AuthorizationScopes) != 0 && len(cr.Spec.ForProvider.AuthorizationScopes) == 0 {
cr.Spec.ForProvider.AuthorizationScopes = make([]*string, len(resp.AuthorizationScopes))
for i0 := range resp.AuthorizationScopes {
cr.Spec.ForProvider.AuthorizationScopes[i0] = li.LateInitializeStringPtr(cr.Spec.ForProvider.AuthorizationScopes[i0], resp.AuthorizationScopes[i0])
}
}
cr.Spec.ForProvider.AuthorizationType = li.LateInitializeStringPtr(cr.Spec.ForProvider.AuthorizationType, resp.AuthorizationType)
cr.Spec.ForProvider.AuthorizerID = li.LateInitializeStringPtr(cr.Spec.ForProvider.AuthorizerID, resp.AuthorizerId)
cr.Spec.ForProvider.ModelSelectionExpression = li.LateInitializeStringPtr(cr.Spec.ForProvider.ModelSelectionExpression, resp.ModelSelectionExpression)
cr.Spec.ForProvider.OperationName = li.LateInitializeStringPtr(cr.Spec.ForProvider.OperationName, resp.OperationName)
if resp.RequestModels != nil {
if cr.Spec.ForProvider.RequestModels == nil {
cr.Spec.ForProvider.RequestModels = map[string]*string{}
}
for key0 := range resp.RequestModels {
cr.Spec.ForProvider.RequestModels[key0] = li.LateInitializeStringPtr(cr.Spec.ForProvider.RequestModels[key0], resp.RequestModels[key0])
}
}
if resp.RequestParameters != nil {
if cr.Spec.ForProvider.RequestParameters == nil {
cr.Spec.ForProvider.RequestParameters = map[string]*svcapitypes.ParameterConstraints{}
}
for key0 := range resp.RequestParameters {
if resp.RequestParameters[key0] != nil {
if cr.Spec.ForProvider.RequestParameters[key0] == nil {
cr.Spec.ForProvider.RequestParameters[key0] = &svcapitypes.ParameterConstraints{}
}
cr.Spec.ForProvider.RequestParameters[key0].Required = li.LateInitializeBoolPtr(cr.Spec.ForProvider.RequestParameters[key0].Required, resp.RequestParameters[key0].Required)
}
}
}
cr.Spec.ForProvider.RouteKey = li.LateInitializeStringPtr(cr.Spec.ForProvider.RouteKey, resp.RouteKey)
cr.Spec.ForProvider.RouteResponseSelectionExpression = li.LateInitializeStringPtr(cr.Spec.ForProvider.RouteResponseSelectionExpression, resp.RouteResponseSelectionExpression)
cr.Spec.ForProvider.Target = li.LateInitializeStringPtr(cr.Spec.ForProvider.Target, resp.Target)`
	assert.Equal(
		expected,
		code.LateInitializeReadOne(crd.Config(), crd, "resp", "cr"),
	)
}

func Test_LateInitializeReadMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "rds", crossplane.DefaultConfig)

	crd := testutil.GetCRDByName(t, g, "DBCluster")
	require.NotNil(crd)

	expected := `for _, resource := range resp.DBClusters {
if len(resource.AvailabilityZones) != 0 && len(cr.Spec.ForProvider.AvailabilityZones) == 0 {
cr.Spec.ForProvider.AvailabilityZones = make([]*string, len(resource.AvailabilityZones))
for i0 := range resource.AvailabilityZones {
cr.Spec.ForProvider.AvailabilityZones[i0] = li.LateInitializeStringPtr(cr.Spec.ForProvider.AvailabilityZones[i0], resource.AvailabilityZones[i0])
}
}
cr.Spec.ForProvider.BacktrackWindow = li.LateInitializeInt64Ptr(cr.Spec.ForProvider.BacktrackWindow, resource.BacktrackWindow)
cr.Spec.ForProvider.BackupRetentionPeriod = li.LateInitializeInt64Ptr(cr.Spec.ForProvider.BackupRetentionPeriod, resource.BackupRetentionPeriod)
cr.Spec.ForProvider.CharacterSetName = li.LateInitializeStringPtr(cr.Spec.ForProvider.CharacterSetName, resource.CharacterSetName)
cr.Spec.ForProvider.CopyTagsToSnapshot = li.LateInitializeBoolPtr(cr.Spec.ForProvider.CopyTagsToSnapshot, resource.CopyTagsToSnapshot)
cr.Spec.ForProvider.DBClusterIdentifier = li.LateInitializeStringPtr(cr.Spec.ForProvider.DBClusterIdentifier, resource.DBClusterIdentifier)
cr.Spec.ForProvider.DatabaseName = li.LateInitializeStringPtr(cr.Spec.ForProvider.DatabaseName, resource.DatabaseName)
cr.Spec.ForProvider.DeletionProtection = li.LateInitializeBoolPtr(cr.Spec.ForProvider.DeletionProtection, resource.DeletionProtection)
cr.Spec.ForProvider.Engine = li.LateInitializeStringPtr(cr.Spec.ForProvider.Engine, resource.Engine)
cr.Spec.ForProvider.EngineMode = li.LateInitializeStringPtr(cr.Spec.ForProvider.EngineMode, resource.EngineMode)
cr.Spec.ForProvider.EngineVersion = li.LateInitializeStringPtr(cr.Spec.ForProvider.EngineVersion, resource.EngineVersion)
cr.Spec.ForProvider.KMSKeyID = li.LateInitializeStringPtr(cr.Spec.ForProvider.KMSKeyID, resource.KmsKeyId)
cr.Spec.ForProvider.MasterUsername = li.LateInitializeStringPtr(cr.Spec.ForProvider.MasterUsername, resource.MasterUsername)
cr.Spec.ForProvider.Port = li.LateInitializeInt64Ptr(cr.Spec.ForProvider.Port, resource.Port)
cr.Spec.ForProvider.PreferredBackupWindow = li.LateInitializeStringPtr(cr.Spec.ForProvider.PreferredBackupWindow, resource.PreferredBackupWindow)
cr.Spec.ForProvider.PreferredMaintenanceWindow = li.LateInitializeStringPtr(cr.Spec.ForProvider.PreferredMaintenanceWindow, resource.PreferredMaintenanceWindow)
cr.Spec.ForProvider.ReplicationSourceIdentifier = li.LateInitializeStringPtr(cr.Spec.ForProvider.ReplicationSourceIdentifier, resource.ReplicationSourceIdentifier)
cr.Spec.ForProvider.StorageEncrypted = li.LateInitializeBoolPtr(cr.Spec.ForProvider.StorageEncrypted, resource.StorageEncrypted)
}`
	assert.Equal(
		expected,
		code.LateInitializeReadMany(crd.Config(), crd, "resp", "cr"),
	)
}

func Test_LateInitializeGetAttributes(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "sqs", crossplane.DefaultConfig)

	crd := testutil.GetCRDByName(t, g, "Queue")
	require.NotNil(crd)

	expected := `cr.Spec.ForProvider.ContentBasedDeduplication = li.LateInitializeStringPtr(cr.Spec.ForProvider.ContentBasedDeduplication, resp.Attributes["ContentBasedDeduplication"])
cr.Spec.ForProvider.DelaySeconds = li.LateInitializeStringPtr(cr.Spec.ForProvider.DelaySeconds, resp.Attributes["DelaySeconds"])
cr.Spec.ForProvider.FifoQueue = li.LateInitializeStringPtr(cr.Spec.ForProvider.FifoQueue, resp.Attributes["FifoQueue"])
cr.Spec.ForProvider.KMSDataKeyReusePeriodSeconds = li.LateInitializeStringPtr(cr.Spec.ForProvider.KMSDataKeyReusePeriodSeconds, resp.Attributes["KmsDataKeyReusePeriodSeconds"])
cr.Spec.ForProvider.KMSMasterKeyID = li.LateInitializeStringPtr(cr.Spec.ForProvider.KMSMasterKeyID, resp.Attributes["KmsMasterKeyId"])
cr.Spec.ForProvider.MaximumMessageSize = li.LateInitializeStringPtr(cr.Spec.ForProvider.MaximumMessageSize, resp.Attributes["MaximumMessageSize"])
cr.Spec.ForProvider.MessageRetentionPeriod = li.LateInitializeStringPtr(cr.Spec.ForProvider.MessageRetentionPeriod, resp.Attributes["MessageRetentionPeriod"])
cr.Spec.ForProvider.Policy = li.LateInitializeStringPtr(cr.Spec.ForProvider.Policy, resp.Attributes["Policy"])
cr.Spec.ForProvider.ReceiveMessageWaitTimeSeconds = li.LateInitializeStringPtr(cr.Spec.ForProvider.ReceiveMessageWaitTimeSeconds, resp.Attributes["ReceiveMessageWaitTimeSeconds"])
cr.Spec.ForProvider.RedrivePolicy = li.LateInitializeStringPtr(cr.Spec.ForProvider.RedrivePolicy, resp.Attributes["RedrivePolicy"])
cr.Spec.ForProvider.VisibilityTimeout = li.LateInitializeStringPtr(cr.Spec.ForProvider.VisibilityTimeout, resp.Attributes["VisibilityTimeout"])`
	assert.Equal(
		expected,
		code.LateInitializeGetAttributes(crd.Config(), crd, "resp", "cr"),
	)
}
