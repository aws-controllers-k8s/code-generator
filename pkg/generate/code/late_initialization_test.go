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

	g := testutil.NewGeneratorForService(t, "dynamodb", crossplane.DefaultConfig)

	crd := testutil.GetCRDByName(t, g, "Table")
	require.NotNil(crd)

	expected := `if len(resp.Table.AttributeDefinitions) != 0 && len(cr.Spec.ForProvider.AttributeDefinitions) == 0 {
cr.Spec.ForProvider.AttributeDefinitions = make([]*svcapitypes.AttributeDefinition, len(resp.Table.AttributeDefinitions))
for i0 := range resp.Table.AttributeDefinitions {
if resp.Table.AttributeDefinitions[i0] != nil {
if cr.Spec.ForProvider.AttributeDefinitions[i0] == nil {
cr.Spec.ForProvider.AttributeDefinitions[i0] = &svcapitypes.AttributeDefinition{}
}
cr.Spec.ForProvider.AttributeDefinitions[i0].AttributeName = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.AttributeDefinitions[i0].AttributeName, resp.Table.AttributeDefinitions[i0].AttributeName)
cr.Spec.ForProvider.AttributeDefinitions[i0].AttributeType = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.AttributeDefinitions[i0].AttributeType, resp.Table.AttributeDefinitions[i0].AttributeType)
}
}
}
if len(resp.Table.GlobalSecondaryIndexes) != 0 && len(cr.Spec.ForProvider.GlobalSecondaryIndexes) == 0 {
cr.Spec.ForProvider.GlobalSecondaryIndexes = make([]*svcapitypes.GlobalSecondaryIndex, len(resp.Table.GlobalSecondaryIndexes))
for i0 := range resp.Table.GlobalSecondaryIndexes {
if resp.Table.GlobalSecondaryIndexes[i0] != nil {
if cr.Spec.ForProvider.GlobalSecondaryIndexes[i0] == nil {
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0] = &svcapitypes.GlobalSecondaryIndex{}
}
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].IndexName = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].IndexName, resp.Table.GlobalSecondaryIndexes[i0].IndexName)
if len(resp.Table.GlobalSecondaryIndexes[i0].KeySchema) != 0 && len(cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].KeySchema) == 0 {
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].KeySchema = make([]*svcapitypes.KeySchemaElement, len(resp.Table.GlobalSecondaryIndexes[i0].KeySchema))
for i2 := range resp.Table.GlobalSecondaryIndexes[i0].KeySchema {
if resp.Table.GlobalSecondaryIndexes[i0].KeySchema[i2] != nil {
if cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].KeySchema[i2] == nil {
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].KeySchema[i2] = &svcapitypes.KeySchemaElement{}
}
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].KeySchema[i2].AttributeName = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].KeySchema[i2].AttributeName, resp.Table.GlobalSecondaryIndexes[i0].KeySchema[i2].AttributeName)
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].KeySchema[i2].KeyType = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].KeySchema[i2].KeyType, resp.Table.GlobalSecondaryIndexes[i0].KeySchema[i2].KeyType)
}
}
}
if resp.Table.GlobalSecondaryIndexes[i0].Projection != nil {
if cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].Projection == nil {
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].Projection = &svcapitypes.Projection{}
}
if len(resp.Table.GlobalSecondaryIndexes[i0].Projection.NonKeyAttributes) != 0 && len(cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].Projection.NonKeyAttributes) == 0 {
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].Projection.NonKeyAttributes = make([]*string, len(resp.Table.GlobalSecondaryIndexes[i0].Projection.NonKeyAttributes))
for i3 := range resp.Table.GlobalSecondaryIndexes[i0].Projection.NonKeyAttributes {
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].Projection.NonKeyAttributes[i3] = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].Projection.NonKeyAttributes[i3], resp.Table.GlobalSecondaryIndexes[i0].Projection.NonKeyAttributes[i3])
}
}
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].Projection.ProjectionType = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].Projection.ProjectionType, resp.Table.GlobalSecondaryIndexes[i0].Projection.ProjectionType)
}
if resp.Table.GlobalSecondaryIndexes[i0].ProvisionedThroughput != nil {
if cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].ProvisionedThroughput == nil {
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].ProvisionedThroughput = &svcapitypes.ProvisionedThroughput{}
}
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].ProvisionedThroughput.ReadCapacityUnits = awsclients.LateInitializeInt64Ptr(cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].ProvisionedThroughput.ReadCapacityUnits, resp.Table.GlobalSecondaryIndexes[i0].ProvisionedThroughput.ReadCapacityUnits)
cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].ProvisionedThroughput.WriteCapacityUnits = awsclients.LateInitializeInt64Ptr(cr.Spec.ForProvider.GlobalSecondaryIndexes[i0].ProvisionedThroughput.WriteCapacityUnits, resp.Table.GlobalSecondaryIndexes[i0].ProvisionedThroughput.WriteCapacityUnits)
}
}
}
}
if len(resp.Table.KeySchema) != 0 && len(cr.Spec.ForProvider.KeySchema) == 0 {
cr.Spec.ForProvider.KeySchema = make([]*svcapitypes.KeySchemaElement, len(resp.Table.KeySchema))
for i0 := range resp.Table.KeySchema {
if resp.Table.KeySchema[i0] != nil {
if cr.Spec.ForProvider.KeySchema[i0] == nil {
cr.Spec.ForProvider.KeySchema[i0] = &svcapitypes.KeySchemaElement{}
}
cr.Spec.ForProvider.KeySchema[i0].AttributeName = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.KeySchema[i0].AttributeName, resp.Table.KeySchema[i0].AttributeName)
cr.Spec.ForProvider.KeySchema[i0].KeyType = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.KeySchema[i0].KeyType, resp.Table.KeySchema[i0].KeyType)
}
}
}
if len(resp.Table.LocalSecondaryIndexes) != 0 && len(cr.Spec.ForProvider.LocalSecondaryIndexes) == 0 {
cr.Spec.ForProvider.LocalSecondaryIndexes = make([]*svcapitypes.LocalSecondaryIndex, len(resp.Table.LocalSecondaryIndexes))
for i0 := range resp.Table.LocalSecondaryIndexes {
if resp.Table.LocalSecondaryIndexes[i0] != nil {
if cr.Spec.ForProvider.LocalSecondaryIndexes[i0] == nil {
cr.Spec.ForProvider.LocalSecondaryIndexes[i0] = &svcapitypes.LocalSecondaryIndex{}
}
cr.Spec.ForProvider.LocalSecondaryIndexes[i0].IndexName = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.LocalSecondaryIndexes[i0].IndexName, resp.Table.LocalSecondaryIndexes[i0].IndexName)
if len(resp.Table.LocalSecondaryIndexes[i0].KeySchema) != 0 && len(cr.Spec.ForProvider.LocalSecondaryIndexes[i0].KeySchema) == 0 {
cr.Spec.ForProvider.LocalSecondaryIndexes[i0].KeySchema = make([]*svcapitypes.KeySchemaElement, len(resp.Table.LocalSecondaryIndexes[i0].KeySchema))
for i2 := range resp.Table.LocalSecondaryIndexes[i0].KeySchema {
if resp.Table.LocalSecondaryIndexes[i0].KeySchema[i2] != nil {
if cr.Spec.ForProvider.LocalSecondaryIndexes[i0].KeySchema[i2] == nil {
cr.Spec.ForProvider.LocalSecondaryIndexes[i0].KeySchema[i2] = &svcapitypes.KeySchemaElement{}
}
cr.Spec.ForProvider.LocalSecondaryIndexes[i0].KeySchema[i2].AttributeName = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.LocalSecondaryIndexes[i0].KeySchema[i2].AttributeName, resp.Table.LocalSecondaryIndexes[i0].KeySchema[i2].AttributeName)
cr.Spec.ForProvider.LocalSecondaryIndexes[i0].KeySchema[i2].KeyType = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.LocalSecondaryIndexes[i0].KeySchema[i2].KeyType, resp.Table.LocalSecondaryIndexes[i0].KeySchema[i2].KeyType)
}
}
}
if resp.Table.LocalSecondaryIndexes[i0].Projection != nil {
if cr.Spec.ForProvider.LocalSecondaryIndexes[i0].Projection == nil {
cr.Spec.ForProvider.LocalSecondaryIndexes[i0].Projection = &svcapitypes.Projection{}
}
if len(resp.Table.LocalSecondaryIndexes[i0].Projection.NonKeyAttributes) != 0 && len(cr.Spec.ForProvider.LocalSecondaryIndexes[i0].Projection.NonKeyAttributes) == 0 {
cr.Spec.ForProvider.LocalSecondaryIndexes[i0].Projection.NonKeyAttributes = make([]*string, len(resp.Table.LocalSecondaryIndexes[i0].Projection.NonKeyAttributes))
for i3 := range resp.Table.LocalSecondaryIndexes[i0].Projection.NonKeyAttributes {
cr.Spec.ForProvider.LocalSecondaryIndexes[i0].Projection.NonKeyAttributes[i3] = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.LocalSecondaryIndexes[i0].Projection.NonKeyAttributes[i3], resp.Table.LocalSecondaryIndexes[i0].Projection.NonKeyAttributes[i3])
}
}
cr.Spec.ForProvider.LocalSecondaryIndexes[i0].Projection.ProjectionType = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.LocalSecondaryIndexes[i0].Projection.ProjectionType, resp.Table.LocalSecondaryIndexes[i0].Projection.ProjectionType)
}
}
}
}
if resp.Table.ProvisionedThroughput != nil {
if cr.Spec.ForProvider.ProvisionedThroughput == nil {
cr.Spec.ForProvider.ProvisionedThroughput = &svcapitypes.ProvisionedThroughput{}
}
cr.Spec.ForProvider.ProvisionedThroughput.ReadCapacityUnits = awsclients.LateInitializeInt64Ptr(cr.Spec.ForProvider.ProvisionedThroughput.ReadCapacityUnits, resp.Table.ProvisionedThroughput.ReadCapacityUnits)
cr.Spec.ForProvider.ProvisionedThroughput.WriteCapacityUnits = awsclients.LateInitializeInt64Ptr(cr.Spec.ForProvider.ProvisionedThroughput.WriteCapacityUnits, resp.Table.ProvisionedThroughput.WriteCapacityUnits)
}
if resp.Table.StreamSpecification != nil {
if cr.Spec.ForProvider.StreamSpecification == nil {
cr.Spec.ForProvider.StreamSpecification = &svcapitypes.StreamSpecification{}
}
cr.Spec.ForProvider.StreamSpecification.StreamEnabled = awsclients.LateInitializeBoolPtr(cr.Spec.ForProvider.StreamSpecification.StreamEnabled, resp.Table.StreamSpecification.StreamEnabled)
cr.Spec.ForProvider.StreamSpecification.StreamViewType = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.StreamSpecification.StreamViewType, resp.Table.StreamSpecification.StreamViewType)
}
cr.Spec.ForProvider.TableName = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.TableName, resp.Table.TableName)`
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
cr.Spec.ForProvider.AvailabilityZones[i0] = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.AvailabilityZones[i0], resource.AvailabilityZones[i0])
}
}
cr.Spec.ForProvider.BacktrackWindow = awsclients.LateInitializeInt64Ptr(cr.Spec.ForProvider.BacktrackWindow, resource.BacktrackWindow)
cr.Spec.ForProvider.BackupRetentionPeriod = awsclients.LateInitializeInt64Ptr(cr.Spec.ForProvider.BackupRetentionPeriod, resource.BackupRetentionPeriod)
cr.Spec.ForProvider.CharacterSetName = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.CharacterSetName, resource.CharacterSetName)
cr.Spec.ForProvider.CopyTagsToSnapshot = awsclients.LateInitializeBoolPtr(cr.Spec.ForProvider.CopyTagsToSnapshot, resource.CopyTagsToSnapshot)
cr.Spec.ForProvider.DBClusterIdentifier = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.DBClusterIdentifier, resource.DBClusterIdentifier)
cr.Spec.ForProvider.DatabaseName = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.DatabaseName, resource.DatabaseName)
cr.Spec.ForProvider.DeletionProtection = awsclients.LateInitializeBoolPtr(cr.Spec.ForProvider.DeletionProtection, resource.DeletionProtection)
cr.Spec.ForProvider.Engine = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.Engine, resource.Engine)
cr.Spec.ForProvider.EngineMode = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.EngineMode, resource.EngineMode)
cr.Spec.ForProvider.EngineVersion = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.EngineVersion, resource.EngineVersion)
cr.Spec.ForProvider.KMSKeyID = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.KMSKeyID, resource.KmsKeyId)
cr.Spec.ForProvider.MasterUsername = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.MasterUsername, resource.MasterUsername)
cr.Spec.ForProvider.Port = awsclients.LateInitializeInt64Ptr(cr.Spec.ForProvider.Port, resource.Port)
cr.Spec.ForProvider.PreferredBackupWindow = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.PreferredBackupWindow, resource.PreferredBackupWindow)
cr.Spec.ForProvider.PreferredMaintenanceWindow = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.PreferredMaintenanceWindow, resource.PreferredMaintenanceWindow)
cr.Spec.ForProvider.ReplicationSourceIdentifier = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.ReplicationSourceIdentifier, resource.ReplicationSourceIdentifier)
cr.Spec.ForProvider.StorageEncrypted = awsclients.LateInitializeBoolPtr(cr.Spec.ForProvider.StorageEncrypted, resource.StorageEncrypted)
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

	expected := `cr.Spec.ForProvider.ContentBasedDeduplication = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.ContentBasedDeduplication, resp.Attributes["ContentBasedDeduplication"])
cr.Spec.ForProvider.DelaySeconds = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.DelaySeconds, resp.Attributes["DelaySeconds"])
cr.Spec.ForProvider.FifoQueue = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.FifoQueue, resp.Attributes["FifoQueue"])
cr.Spec.ForProvider.KMSDataKeyReusePeriodSeconds = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.KMSDataKeyReusePeriodSeconds, resp.Attributes["KmsDataKeyReusePeriodSeconds"])
cr.Spec.ForProvider.KMSMasterKeyID = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.KMSMasterKeyID, resp.Attributes["KmsMasterKeyId"])
cr.Spec.ForProvider.MaximumMessageSize = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.MaximumMessageSize, resp.Attributes["MaximumMessageSize"])
cr.Spec.ForProvider.MessageRetentionPeriod = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.MessageRetentionPeriod, resp.Attributes["MessageRetentionPeriod"])
cr.Spec.ForProvider.Policy = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.Policy, resp.Attributes["Policy"])
cr.Spec.ForProvider.ReceiveMessageWaitTimeSeconds = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.ReceiveMessageWaitTimeSeconds, resp.Attributes["ReceiveMessageWaitTimeSeconds"])
cr.Spec.ForProvider.RedrivePolicy = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.RedrivePolicy, resp.Attributes["RedrivePolicy"])
cr.Spec.ForProvider.VisibilityTimeout = awsclients.LateInitializeStringPtr(cr.Spec.ForProvider.VisibilityTimeout, resp.Attributes["VisibilityTimeout"])`
	assert.Equal(
		expected,
		code.LateInitializeGetAttributes(crd.Config(), crd, "resp", "cr"),
	)
}
