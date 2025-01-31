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

func TestRDS_DBInstance(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "rds")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("DBInstance", crds)
	require.NotNil(crd)

	assert.Equal("DBInstance", crd.Names.Camel)
	assert.Equal("dbInstance", crd.Names.CamelLower)
	assert.Equal("db_instance", crd.Names.Snake)

	// The RDS DBInstance API has the following operations:
	// - CreateDBInstance
	// - DescribeDBInstances
	// - ModifyDBInstance
	// - DeleteDBInstance
	require.NotNil(crd.Ops)

	assert.NotNil(crd.Ops.Create)
	assert.NotNil(crd.Ops.Delete)
	assert.NotNil(crd.Ops.ReadMany)
	assert.NotNil(crd.Ops.Update)

	assert.Nil(crd.Ops.ReadOne)
	assert.Nil(crd.Ops.GetAttributes)
	assert.Nil(crd.Ops.SetAttributes)

	// The DescribeDBInstances operation has the following definition:
	//
	//    "DescribeDBInstances":{
	//      "name":"DescribeDBInstances",
	//      "http":{
	//        "method":"POST",
	//        "requestUri":"/"
	//      },
	//      "input":{"shape":"DescribeDBInstancesMessage"},
	//      "output":{
	//        "shape":"DBInstanceMessage",
	//        "resultWrapper":"DescribeDBInstancesResult"
	//      },
	//      "errors":[
	//        {"shape":"DBInstanceNotFoundFault"}
	//      ]
	//    },
	//
	// NOTE: This is UNUSUAL for List operation to return a 404 Not Found
	// instead of a 200 OK with an empty array of results.
	//
	// Where the DBInstanceNotFoundFault shape looks like this:
	//
	//    "DBInstanceNotFoundFault":{
	//      "type":"structure",
	//      "members":{
	//      },
	//      "error":{
	//        "code":"DBInstanceNotFound",
	//        "httpStatusCode":404,
	//        "senderFault":true
	//      },
	//      "exception":true
	//    },
	//
	// Which clearly indicates it is the 404 HTTP fault for this resource even
	// though the "code" is "DBInstanceNotFound"
	assert.Equal("DBInstanceNotFound", crd.ExceptionCode(404))

	specFields := crd.SpecFields
	statusFields := crd.StatusFields

	expSpecFieldCamel := []string{
		"AllocatedStorage",
		"AutoMinorVersionUpgrade",
		"AvailabilityZone",
		"BackupRetentionPeriod",
		"BackupTarget",
		"CACertificateIdentifier",
		"CharacterSetName",
		"CopyTagsToSnapshot",
		"CustomIAMInstanceProfile",
		"DBClusterIdentifier",
		"DBClusterSnapshotIdentifier",
		"DBInstanceClass",
		"DBInstanceIdentifier",
		"DBName",
		"DBParameterGroupName",
		"DBParameterGroupRef",
		"DBSnapshotIdentifier",
		"DBSubnetGroupName",
		"DBSubnetGroupRef",
		"DBSystemID",
		"DeletionProtection",
		"DestinationRegion",
		"Domain",
		"DomainFqdn",
		"DomainIAMRoleName",
		"DomainOu",
		"EnableCloudwatchLogsExports",
		"EnableCustomerOwnedIP",
		"EnableIAMDatabaseAuthentication",
		"Engine",
		"EngineLifecycleSupport",
		"EngineVersion",
		"IOPS",
		"KMSKeyID",
		"KMSKeyRef",
		"LicenseModel",
		"ManageMasterUserPassword",
		"MasterUserPassword",
		"MasterUserSecretKMSKeyID",
		"MasterUserSecretKMSKeyRef",
		"MasterUsername",
		"MaxAllocatedStorage",
		"MonitoringInterval",
		"MonitoringRoleARN",
		"MultiAZ",
		"NcharCharacterSetName",
		"NetworkType",
		"OptionGroupName",
		"PerformanceInsightsEnabled",
		"PerformanceInsightsKMSKeyID",
		"PerformanceInsightsRetentionPeriod",
		"Port",
		"PreSignedURL",
		"PreferredBackupWindow",
		"PreferredMaintenanceWindow",
		"ProcessorFeatures",
		"PromotionTier",
		"PubliclyAccessible",
		"ReplicaMode",
		"SourceDBInstanceIdentifier",
		"SourceRegion",
		"StorageEncrypted",
		"StorageThroughput",
		"StorageType",
		"TDECredentialARN",
		"TDECredentialPassword",
		"Tags",
		"Timezone",
		"UseDefaultProcessorFeatures",
		"VPCSecurityGroupIDs",
		"VPCSecurityGroupRefs",
	}
	assert.Equal(expSpecFieldCamel, attrCamelNames(specFields))

	expStatusFieldCamel := []string{
		"AWSBackupRecoveryPointARN",
		"ActivityStreamEngineNativeAuditFieldsIncluded",
		"ActivityStreamKMSKeyID",
		"ActivityStreamKinesisStreamName",
		"ActivityStreamMode",
		"ActivityStreamPolicyStatus",
		"ActivityStreamStatus",
		"AssociatedRoles",
		"AutomaticRestartTime",
		"AutomationMode",
		"CertificateDetails",
		"CustomerOwnedIPEnabled",
		"DBIResourceID",
		"DBInstanceAutomatedBackupsReplications",
		"DBInstancePort",
		"DBInstanceStatus",
		"DBParameterGroups",
		"DBSubnetGroup",
		"DatabaseInsightsMode",
		"DedicatedLogVolume",
		"DomainMemberships",
		"EnabledCloudwatchLogsExports",
		"Endpoint",
		"EnhancedMonitoringResourceARN",
		"IAMDatabaseAuthenticationEnabled",
		"InstanceCreateTime",
		"IsStorageConfigUpgradeAvailable",
		"LatestRestorableTime",
		"ListenerEndpoint",
		"MasterUserSecret",
		"MultiTenant",
		"OptionGroupMemberships",
		"PendingModifiedValues",
		"ReadReplicaDBClusterIdentifiers",
		"ReadReplicaDBInstanceIdentifiers",
		"ReadReplicaSourceDBClusterIdentifier",
		"ReadReplicaSourceDBInstanceIdentifier",
		"ResumeFullAutomationModeTime",
		"SecondaryAvailabilityZone",
		"StatusInfos",
		"VPCSecurityGroups",
	}
	assert.Equal(expStatusFieldCamel, attrCamelNames(statusFields))
}
