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

func TestElasticache_ReplicationGroup(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("ReplicationGroup", crds)
	require.NotNil(crd)

	assert.Equal("ReplicationGroup", crd.Names.Camel)
	assert.Equal("replicationGroup", crd.Names.CamelLower)
	assert.Equal("replication_group", crd.Names.Snake)

	// The DescribeCacheClusters operation has the following definition:
	//
	//	"DescribeReplicationGroups":{
	//		"name":"DescribeReplicationGroups",
	//			"http":{
	//			"method":"POST",
	//				"requestUri":"/"
	//		},
	//		"input":{"shape":"DescribeReplicationGroupsMessage"},
	//		"output":{
	//			"shape":"ReplicationGroupMessage",
	//				"resultWrapper":"DescribeReplicationGroupsResult"
	//		},
	//		"errors":[
	//			{"shape":"ReplicationGroupNotFoundFault"},
	//			{"shape":"InvalidParameterValueException"},
	//			{"shape":"InvalidParameterCombinationException"}
	//		]
	//	}
	//
	// Where the ReplicationGroupNotFoundFault shape looks like this:
	//
	//    "ReplicationGroupNotFoundFault":{
	//      "type":"structure",
	//      "members":{
	//      },
	//      "error":{
	//        "code":"ReplicationGroupNotFoundFault",
	//        "httpStatusCode":404,
	//        "senderFault":true
	//      },
	//      "exception":true
	//    },
	//
	// Which indicates that the error is a 404 and is our NotFoundException
	// error with a "code" value of ReplicationGroupNotFoundFault
	assert.Equal("ReplicationGroupNotFoundFault", crd.ExceptionCode(404))

	// The Elasticache ReplicationGroup API has CUD+L operations:
	//
	// * CreateReplicationGroup
	// * DeleteReplicationGroup
	// * ModifyReplicationGroup
	// * DescribeReplicationGroup
	require.NotNil(crd.Ops)

	assert.NotNil(crd.Ops.Create)
	assert.NotNil(crd.Ops.Delete)
	assert.NotNil(crd.Ops.Update)
	assert.NotNil(crd.Ops.ReadMany)

	// But no ReadOne operation...
	assert.Nil(crd.Ops.ReadOne)

	// And no separate get/set attributes calls.
	assert.Nil(crd.Ops.GetAttributes)
	assert.Nil(crd.Ops.SetAttributes)

	specFields := crd.SpecFields
	statusFields := crd.StatusFields

	expSpecFieldCamel := []string{
		"AtRestEncryptionEnabled",
		"AuthToken",
		"AutoMinorVersionUpgrade",
		"AutomaticFailoverEnabled",
		"CacheNodeType",
		"CacheParameterGroupName",
		"CacheSecurityGroupNames",
		"CacheSubnetGroupName",
		"ClusterMode",
		"DataTieringEnabled",
		"Engine",
		"EngineVersion",
		"IPDiscovery",
		"KMSKeyID",
		"LogDeliveryConfigurations",
		"MultiAZEnabled",
		"NetworkType",
		"NodeGroupConfiguration",
		"NotificationTopicARN",
		"NumCacheClusters",
		"NumNodeGroups",
		"Port",
		"PreferredCacheClusterAZs",
		"PreferredMaintenanceWindow",
		"PrimaryClusterID",
		"ReplicasPerNodeGroup",
		"ReplicationGroupDescription",
		"ReplicationGroupID",
		"SecurityGroupIDs",
		"ServerlessCacheSnapshotName",
		"SnapshotARNs",
		"SnapshotName",
		"SnapshotRetentionLimit",
		"SnapshotWindow",
		"Tags",
		"TransitEncryptionEnabled",
		"TransitEncryptionMode",
		"UserGroupIDs",
	}
	assert.Equal(expSpecFieldCamel, attrCamelNames(specFields))

	expStatusFieldCamel := []string{
		"AllowedScaleDownModifications",
		"AllowedScaleUpModifications",
		"AuthTokenEnabled",
		"AuthTokenLastModifiedDate",
		"AutomaticFailover",
		"ClusterEnabled",
		"ConfigurationEndpoint",
		"DataTiering",
		"Description",
		"Events",
		"GlobalReplicationGroupInfo",
		"MemberClusters",
		"MemberClustersOutpostARNs",
		"MultiAZ",
		"NodeGroups",
		"PendingModifiedValues",
		"ReplicationGroupCreateTime",
		"SnapshottingClusterID",
		"Status",
	}
	assert.Equal(expStatusFieldCamel, attrCamelNames(statusFields))
}

func TestElasticache_Ignored_Resources(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("GlobalReplicationGroup", crds)
	require.Nil(crd)
}

func TestElasticache_Additional_Snapshot_Spec(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")
	crds, err := g.GetCRDs()

	require.Nil(err)

	crd := getCRDByName("Snapshot", crds)
	require.NotNil(crd)

	assert := assert.New(t)
	assert.Contains(crd.SpecFields, "SourceSnapshotName")
}

func TestElasticache_Additional_CacheParameterGroup_Spec(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")
	crds, err := g.GetCRDs()

	require.Nil(err)

	crd := getCRDByName("CacheParameterGroup", crds)
	require.NotNil(crd)

	assert := assert.New(t)
	assert.Contains(crd.SpecFields, "ParameterNameValues")
}

func TestElasticache_Additional_CacheParameterGroup_Status(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")
	crds, err := g.GetCRDs()

	require.Nil(err)

	crd := getCRDByName("CacheParameterGroup", crds)
	require.NotNil(crd)

	assert := assert.New(t)
	assert.Contains(crd.StatusFields, "Parameters")
	assert.Contains(crd.StatusFields, "Events")
}

func TestElasticache_Additional_ReplicationGroup_Status(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")
	crds, err := g.GetCRDs()

	require.Nil(err)

	crd := getCRDByName("ReplicationGroup", crds)
	require.NotNil(crd)

	assert := assert.New(t)
	assert.Contains(crd.StatusFields, "Events")
}

func TestElasticache_Additional_CacheSubnetGroup_Status(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")
	crds, err := g.GetCRDs()

	require.Nil(err)

	crd := getCRDByName("CacheSubnetGroup", crds)
	require.NotNil(crd)

	assert := assert.New(t)
	assert.Contains(crd.StatusFields, "Events")
}

func TestElasticache_Additional_ReplicationGroup_Status_RenameField(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")
	crds, err := g.GetCRDs()

	require.Nil(err)

	crd := getCRDByName("ReplicationGroup", crds)
	require.NotNil(crd)

	assert := assert.New(t)
	assert.Contains(crd.StatusFields, "AllowedScaleUpModifications")
	assert.Contains(crd.StatusFields, "AllowedScaleDownModifications")
}

func TestElasticache_ValidateAuthTokenIsSecret(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")
	crds, err := g.GetCRDs()

	require.Nil(err)

	crd := getCRDByName("ReplicationGroup", crds)
	require.NotNil(crd)

	assert := assert.New(t)
	assert.Equal("*ackv1alpha1.SecretKeyReference", crd.SpecFields["AuthToken"].GoType)
	assert.Equal("SecretKeyReference", crd.SpecFields["AuthToken"].GoTypeElem)
	assert.Equal("*ackv1alpha1.SecretKeyReference", crd.SpecFields["AuthToken"].GoTypeWithPkgName)

	crd = getCRDByName("User", crds)
	require.NotNil(crd)

	assert.Equal("[]*ackv1alpha1.SecretKeyReference", crd.SpecFields["Passwords"].GoType)
	assert.Equal("SecretKeyReference", crd.SpecFields["Passwords"].GoTypeElem)
	assert.Equal("[]*ackv1alpha1.SecretKeyReference", crd.SpecFields["Passwords"].GoTypeWithPkgName)
}
