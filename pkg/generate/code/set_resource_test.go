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
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestSetResource_APIGWv2_Route_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "Route")
	require.NotNil(crd)

	expected := `
	if resp.ApiGatewayManaged != nil {
		ko.Status.APIGatewayManaged = resp.ApiGatewayManaged
	} else {
		ko.Status.APIGatewayManaged = nil
	}
	if resp.ApiKeyRequired != nil {
		ko.Spec.APIKeyRequired = resp.ApiKeyRequired
	} else {
		ko.Spec.APIKeyRequired = nil
	}
	if resp.AuthorizationScopes != nil {
		f2 := []*string{}
		for _, f2iter := range resp.AuthorizationScopes {
			var f2elem string
			f2elem = *f2iter
			f2 = append(f2, &f2elem)
		}
		ko.Spec.AuthorizationScopes = f2
	} else {
		ko.Spec.AuthorizationScopes = nil
	}
	if resp.AuthorizationType != nil {
		ko.Spec.AuthorizationType = resp.AuthorizationType
	} else {
		ko.Spec.AuthorizationType = nil
	}
	if resp.AuthorizerId != nil {
		ko.Spec.AuthorizerID = resp.AuthorizerId
	} else {
		ko.Spec.AuthorizerID = nil
	}
	if resp.ModelSelectionExpression != nil {
		ko.Spec.ModelSelectionExpression = resp.ModelSelectionExpression
	} else {
		ko.Spec.ModelSelectionExpression = nil
	}
	if resp.OperationName != nil {
		ko.Spec.OperationName = resp.OperationName
	} else {
		ko.Spec.OperationName = nil
	}
	if resp.RequestModels != nil {
		f7 := map[string]*string{}
		for f7key, f7valiter := range resp.RequestModels {
			var f7val string
			f7val = *f7valiter
			f7[f7key] = &f7val
		}
		ko.Spec.RequestModels = f7
	} else {
		ko.Spec.RequestModels = nil
	}
	if resp.RequestParameters != nil {
		f8 := map[string]*svcapitypes.ParameterConstraints{}
		for f8key, f8valiter := range resp.RequestParameters {
			f8val := &svcapitypes.ParameterConstraints{}
			if f8valiter.Required != nil {
				f8val.Required = f8valiter.Required
			}
			f8[f8key] = f8val
		}
		ko.Spec.RequestParameters = f8
	} else {
		ko.Spec.RequestParameters = nil
	}
	if resp.RouteId != nil {
		ko.Status.RouteID = resp.RouteId
	} else {
		ko.Status.RouteID = nil
	}
	if resp.RouteKey != nil {
		ko.Spec.RouteKey = resp.RouteKey
	} else {
		ko.Spec.RouteKey = nil
	}
	if resp.RouteResponseSelectionExpression != nil {
		ko.Spec.RouteResponseSelectionExpression = resp.RouteResponseSelectionExpression
	} else {
		ko.Spec.RouteResponseSelectionExpression = nil
	}
	if resp.Target != nil {
		ko.Spec.Target = resp.Target
	} else {
		ko.Spec.Target = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}

func TestSetResource_APIGWv2_Route_ReadOne(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "Route")
	require.NotNil(crd)

	expected := `
	if resp.ApiGatewayManaged != nil {
		ko.Status.APIGatewayManaged = resp.ApiGatewayManaged
	} else {
		ko.Status.APIGatewayManaged = nil
	}
	if resp.ApiKeyRequired != nil {
		ko.Spec.APIKeyRequired = resp.ApiKeyRequired
	} else {
		ko.Spec.APIKeyRequired = nil
	}
	if resp.AuthorizationScopes != nil {
		f2 := []*string{}
		for _, f2iter := range resp.AuthorizationScopes {
			var f2elem string
			f2elem = *f2iter
			f2 = append(f2, &f2elem)
		}
		ko.Spec.AuthorizationScopes = f2
	} else {
		ko.Spec.AuthorizationScopes = nil
	}
	if resp.AuthorizationType != nil {
		ko.Spec.AuthorizationType = resp.AuthorizationType
	} else {
		ko.Spec.AuthorizationType = nil
	}
	if resp.AuthorizerId != nil {
		ko.Spec.AuthorizerID = resp.AuthorizerId
	} else {
		ko.Spec.AuthorizerID = nil
	}
	if resp.ModelSelectionExpression != nil {
		ko.Spec.ModelSelectionExpression = resp.ModelSelectionExpression
	} else {
		ko.Spec.ModelSelectionExpression = nil
	}
	if resp.OperationName != nil {
		ko.Spec.OperationName = resp.OperationName
	} else {
		ko.Spec.OperationName = nil
	}
	if resp.RequestModels != nil {
		f7 := map[string]*string{}
		for f7key, f7valiter := range resp.RequestModels {
			var f7val string
			f7val = *f7valiter
			f7[f7key] = &f7val
		}
		ko.Spec.RequestModels = f7
	} else {
		ko.Spec.RequestModels = nil
	}
	if resp.RequestParameters != nil {
		f8 := map[string]*svcapitypes.ParameterConstraints{}
		for f8key, f8valiter := range resp.RequestParameters {
			f8val := &svcapitypes.ParameterConstraints{}
			if f8valiter.Required != nil {
				f8val.Required = f8valiter.Required
			}
			f8[f8key] = f8val
		}
		ko.Spec.RequestParameters = f8
	} else {
		ko.Spec.RequestParameters = nil
	}
	if resp.RouteId != nil {
		ko.Status.RouteID = resp.RouteId
	} else {
		ko.Status.RouteID = nil
	}
	if resp.RouteKey != nil {
		ko.Spec.RouteKey = resp.RouteKey
	} else {
		ko.Spec.RouteKey = nil
	}
	if resp.RouteResponseSelectionExpression != nil {
		ko.Spec.RouteResponseSelectionExpression = resp.RouteResponseSelectionExpression
	} else {
		ko.Spec.RouteResponseSelectionExpression = nil
	}
	if resp.Target != nil {
		ko.Spec.Target = resp.Target
	} else {
		ko.Spec.Target = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeGet, "resp", "ko", 1),
	)
}

func TestSetResource_DynamoDB_Backup_ReadOne(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "dynamodb")

	crd := testutil.GetCRDByName(t, g, "Backup")
	require.NotNil(crd)

	expected := `
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.BackupDescription.BackupDetails.BackupArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.BackupDescription.BackupDetails.BackupArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.BackupDescription.BackupDetails.BackupCreationDateTime != nil {
		ko.Status.BackupCreationDateTime = &metav1.Time{*resp.BackupDescription.BackupDetails.BackupCreationDateTime}
	} else {
		ko.Status.BackupCreationDateTime = nil
	}
	if resp.BackupDescription.BackupDetails.BackupExpiryDateTime != nil {
		ko.Status.BackupExpiryDateTime = &metav1.Time{*resp.BackupDescription.BackupDetails.BackupExpiryDateTime}
	} else {
		ko.Status.BackupExpiryDateTime = nil
	}
	if resp.BackupDescription.BackupDetails.BackupName != nil {
		ko.Spec.BackupName = resp.BackupDescription.BackupDetails.BackupName
	} else {
		ko.Spec.BackupName = nil
	}
	if resp.BackupDescription.BackupDetails.BackupSizeBytes != nil {
		ko.Status.BackupSizeBytes = resp.BackupDescription.BackupDetails.BackupSizeBytes
	} else {
		ko.Status.BackupSizeBytes = nil
	}
	if resp.BackupDescription.BackupDetails.BackupStatus != nil {
		ko.Status.BackupStatus = resp.BackupDescription.BackupDetails.BackupStatus
	} else {
		ko.Status.BackupStatus = nil
	}
	if resp.BackupDescription.BackupDetails.BackupType != nil {
		ko.Status.BackupType = resp.BackupDescription.BackupDetails.BackupType
	} else {
		ko.Status.BackupType = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeGet, "resp", "ko", 1),
	)
}

func TestSetResource_CodeDeploy_Deployment_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "codedeploy")

	crd := testutil.GetCRDByName(t, g, "Deployment")
	require.NotNil(crd)

	// However, all of the fields in the Deployment resource's
	// CreateDeploymentInput shape are returned in the GetDeploymentOutput
	// shape, and there is a DeploymentInfo wrapper struct for the output
	// shape, so the readOne accessor contains the wrapper struct's name.
	expected := `
	if resp.DeploymentId != nil {
		ko.Status.DeploymentID = resp.DeploymentId
	} else {
		ko.Status.DeploymentID = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}

func TestSetResource_DynamoDB_Table_ReadOne(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "dynamodb")

	crd := testutil.GetCRDByName(t, g, "Table")
	require.NotNil(crd)

	// The DynamoDB API uses an API that uses "wrapper" single-member objects
	// in the JSON response for the create/describe calls. In other words, the
	// returned result from the DescribeTable API looks like this:
	//
	// {
	//   "table": {
	//	 .. bunch of fields for the table ..
	//   }
	// }
	//
	// However, the *ShapeName* of the "table" field is actually
	// TableDescription. This tests that we're properly outputting the
	// memberName (which is "Table" and not "TableDescription") when we build
	// the Table CRD's Status field from the DescribeTableOutput shape.
	expected := `
	if resp.Table.ArchivalSummary != nil {
		f0 := &svcapitypes.ArchivalSummary{}
		if resp.Table.ArchivalSummary.ArchivalBackupArn != nil {
			f0.ArchivalBackupARN = resp.Table.ArchivalSummary.ArchivalBackupArn
		}
		if resp.Table.ArchivalSummary.ArchivalDateTime != nil {
			f0.ArchivalDateTime = &metav1.Time{*resp.Table.ArchivalSummary.ArchivalDateTime}
		}
		if resp.Table.ArchivalSummary.ArchivalReason != nil {
			f0.ArchivalReason = resp.Table.ArchivalSummary.ArchivalReason
		}
		ko.Status.ArchivalSummary = f0
	} else {
		ko.Status.ArchivalSummary = nil
	}
	if resp.Table.AttributeDefinitions != nil {
		f1 := []*svcapitypes.AttributeDefinition{}
		for _, f1iter := range resp.Table.AttributeDefinitions {
			f1elem := &svcapitypes.AttributeDefinition{}
			if f1iter.AttributeName != nil {
				f1elem.AttributeName = f1iter.AttributeName
			}
			if f1iter.AttributeType != nil {
				f1elem.AttributeType = f1iter.AttributeType
			}
			f1 = append(f1, f1elem)
		}
		ko.Spec.AttributeDefinitions = f1
	} else {
		ko.Spec.AttributeDefinitions = nil
	}
	if resp.Table.BillingModeSummary != nil {
		f2 := &svcapitypes.BillingModeSummary{}
		if resp.Table.BillingModeSummary.BillingMode != nil {
			f2.BillingMode = resp.Table.BillingModeSummary.BillingMode
		}
		if resp.Table.BillingModeSummary.LastUpdateToPayPerRequestDateTime != nil {
			f2.LastUpdateToPayPerRequestDateTime = &metav1.Time{*resp.Table.BillingModeSummary.LastUpdateToPayPerRequestDateTime}
		}
		ko.Status.BillingModeSummary = f2
	} else {
		ko.Status.BillingModeSummary = nil
	}
	if resp.Table.CreationDateTime != nil {
		ko.Status.CreationDateTime = &metav1.Time{*resp.Table.CreationDateTime}
	} else {
		ko.Status.CreationDateTime = nil
	}
	if resp.Table.GlobalSecondaryIndexes != nil {
		f4 := []*svcapitypes.GlobalSecondaryIndex{}
		for _, f4iter := range resp.Table.GlobalSecondaryIndexes {
			f4elem := &svcapitypes.GlobalSecondaryIndex{}
			if f4iter.IndexName != nil {
				f4elem.IndexName = f4iter.IndexName
			}
			if f4iter.KeySchema != nil {
				f4elemf6 := []*svcapitypes.KeySchemaElement{}
				for _, f4elemf6iter := range f4iter.KeySchema {
					f4elemf6elem := &svcapitypes.KeySchemaElement{}
					if f4elemf6iter.AttributeName != nil {
						f4elemf6elem.AttributeName = f4elemf6iter.AttributeName
					}
					if f4elemf6iter.KeyType != nil {
						f4elemf6elem.KeyType = f4elemf6iter.KeyType
					}
					f4elemf6 = append(f4elemf6, f4elemf6elem)
				}
				f4elem.KeySchema = f4elemf6
			}
			if f4iter.Projection != nil {
				f4elemf7 := &svcapitypes.Projection{}
				if f4iter.Projection.NonKeyAttributes != nil {
					f4elemf7f0 := []*string{}
					for _, f4elemf7f0iter := range f4iter.Projection.NonKeyAttributes {
						var f4elemf7f0elem string
						f4elemf7f0elem = *f4elemf7f0iter
						f4elemf7f0 = append(f4elemf7f0, &f4elemf7f0elem)
					}
					f4elemf7.NonKeyAttributes = f4elemf7f0
				}
				if f4iter.Projection.ProjectionType != nil {
					f4elemf7.ProjectionType = f4iter.Projection.ProjectionType
				}
				f4elem.Projection = f4elemf7
			}
			if f4iter.ProvisionedThroughput != nil {
				f4elemf8 := &svcapitypes.ProvisionedThroughput{}
				if f4iter.ProvisionedThroughput.ReadCapacityUnits != nil {
					f4elemf8.ReadCapacityUnits = f4iter.ProvisionedThroughput.ReadCapacityUnits
				}
				if f4iter.ProvisionedThroughput.WriteCapacityUnits != nil {
					f4elemf8.WriteCapacityUnits = f4iter.ProvisionedThroughput.WriteCapacityUnits
				}
				f4elem.ProvisionedThroughput = f4elemf8
			}
			f4 = append(f4, f4elem)
		}
		ko.Spec.GlobalSecondaryIndexes = f4
	} else {
		ko.Spec.GlobalSecondaryIndexes = nil
	}
	if resp.Table.GlobalTableVersion != nil {
		ko.Status.GlobalTableVersion = resp.Table.GlobalTableVersion
	} else {
		ko.Status.GlobalTableVersion = nil
	}
	if resp.Table.ItemCount != nil {
		ko.Status.ItemCount = resp.Table.ItemCount
	} else {
		ko.Status.ItemCount = nil
	}
	if resp.Table.KeySchema != nil {
		f7 := []*svcapitypes.KeySchemaElement{}
		for _, f7iter := range resp.Table.KeySchema {
			f7elem := &svcapitypes.KeySchemaElement{}
			if f7iter.AttributeName != nil {
				f7elem.AttributeName = f7iter.AttributeName
			}
			if f7iter.KeyType != nil {
				f7elem.KeyType = f7iter.KeyType
			}
			f7 = append(f7, f7elem)
		}
		ko.Spec.KeySchema = f7
	} else {
		ko.Spec.KeySchema = nil
	}
	if resp.Table.LatestStreamArn != nil {
		ko.Status.LatestStreamARN = resp.Table.LatestStreamArn
	} else {
		ko.Status.LatestStreamARN = nil
	}
	if resp.Table.LatestStreamLabel != nil {
		ko.Status.LatestStreamLabel = resp.Table.LatestStreamLabel
	} else {
		ko.Status.LatestStreamLabel = nil
	}
	if resp.Table.LocalSecondaryIndexes != nil {
		f10 := []*svcapitypes.LocalSecondaryIndex{}
		for _, f10iter := range resp.Table.LocalSecondaryIndexes {
			f10elem := &svcapitypes.LocalSecondaryIndex{}
			if f10iter.IndexName != nil {
				f10elem.IndexName = f10iter.IndexName
			}
			if f10iter.KeySchema != nil {
				f10elemf4 := []*svcapitypes.KeySchemaElement{}
				for _, f10elemf4iter := range f10iter.KeySchema {
					f10elemf4elem := &svcapitypes.KeySchemaElement{}
					if f10elemf4iter.AttributeName != nil {
						f10elemf4elem.AttributeName = f10elemf4iter.AttributeName
					}
					if f10elemf4iter.KeyType != nil {
						f10elemf4elem.KeyType = f10elemf4iter.KeyType
					}
					f10elemf4 = append(f10elemf4, f10elemf4elem)
				}
				f10elem.KeySchema = f10elemf4
			}
			if f10iter.Projection != nil {
				f10elemf5 := &svcapitypes.Projection{}
				if f10iter.Projection.NonKeyAttributes != nil {
					f10elemf5f0 := []*string{}
					for _, f10elemf5f0iter := range f10iter.Projection.NonKeyAttributes {
						var f10elemf5f0elem string
						f10elemf5f0elem = *f10elemf5f0iter
						f10elemf5f0 = append(f10elemf5f0, &f10elemf5f0elem)
					}
					f10elemf5.NonKeyAttributes = f10elemf5f0
				}
				if f10iter.Projection.ProjectionType != nil {
					f10elemf5.ProjectionType = f10iter.Projection.ProjectionType
				}
				f10elem.Projection = f10elemf5
			}
			f10 = append(f10, f10elem)
		}
		ko.Spec.LocalSecondaryIndexes = f10
	} else {
		ko.Spec.LocalSecondaryIndexes = nil
	}
	if resp.Table.ProvisionedThroughput != nil {
		f11 := &svcapitypes.ProvisionedThroughput{}
		if resp.Table.ProvisionedThroughput.ReadCapacityUnits != nil {
			f11.ReadCapacityUnits = resp.Table.ProvisionedThroughput.ReadCapacityUnits
		}
		if resp.Table.ProvisionedThroughput.WriteCapacityUnits != nil {
			f11.WriteCapacityUnits = resp.Table.ProvisionedThroughput.WriteCapacityUnits
		}
		ko.Spec.ProvisionedThroughput = f11
	} else {
		ko.Spec.ProvisionedThroughput = nil
	}
	if resp.Table.Replicas != nil {
		f12 := []*svcapitypes.ReplicaDescription{}
		for _, f12iter := range resp.Table.Replicas {
			f12elem := &svcapitypes.ReplicaDescription{}
			if f12iter.GlobalSecondaryIndexes != nil {
				f12elemf0 := []*svcapitypes.ReplicaGlobalSecondaryIndexDescription{}
				for _, f12elemf0iter := range f12iter.GlobalSecondaryIndexes {
					f12elemf0elem := &svcapitypes.ReplicaGlobalSecondaryIndexDescription{}
					if f12elemf0iter.IndexName != nil {
						f12elemf0elem.IndexName = f12elemf0iter.IndexName
					}
					if f12elemf0iter.ProvisionedThroughputOverride != nil {
						f12elemf0elemf1 := &svcapitypes.ProvisionedThroughputOverride{}
						if f12elemf0iter.ProvisionedThroughputOverride.ReadCapacityUnits != nil {
							f12elemf0elemf1.ReadCapacityUnits = f12elemf0iter.ProvisionedThroughputOverride.ReadCapacityUnits
						}
						f12elemf0elem.ProvisionedThroughputOverride = f12elemf0elemf1
					}
					f12elemf0 = append(f12elemf0, f12elemf0elem)
				}
				f12elem.GlobalSecondaryIndexes = f12elemf0
			}
			if f12iter.KMSMasterKeyId != nil {
				f12elem.KMSMasterKeyID = f12iter.KMSMasterKeyId
			}
			if f12iter.ProvisionedThroughputOverride != nil {
				f12elemf2 := &svcapitypes.ProvisionedThroughputOverride{}
				if f12iter.ProvisionedThroughputOverride.ReadCapacityUnits != nil {
					f12elemf2.ReadCapacityUnits = f12iter.ProvisionedThroughputOverride.ReadCapacityUnits
				}
				f12elem.ProvisionedThroughputOverride = f12elemf2
			}
			if f12iter.RegionName != nil {
				f12elem.RegionName = f12iter.RegionName
			}
			if f12iter.ReplicaStatus != nil {
				f12elem.ReplicaStatus = f12iter.ReplicaStatus
			}
			if f12iter.ReplicaStatusDescription != nil {
				f12elem.ReplicaStatusDescription = f12iter.ReplicaStatusDescription
			}
			if f12iter.ReplicaStatusPercentProgress != nil {
				f12elem.ReplicaStatusPercentProgress = f12iter.ReplicaStatusPercentProgress
			}
			f12 = append(f12, f12elem)
		}
		ko.Status.Replicas = f12
	} else {
		ko.Status.Replicas = nil
	}
	if resp.Table.RestoreSummary != nil {
		f13 := &svcapitypes.RestoreSummary{}
		if resp.Table.RestoreSummary.RestoreDateTime != nil {
			f13.RestoreDateTime = &metav1.Time{*resp.Table.RestoreSummary.RestoreDateTime}
		}
		if resp.Table.RestoreSummary.RestoreInProgress != nil {
			f13.RestoreInProgress = resp.Table.RestoreSummary.RestoreInProgress
		}
		if resp.Table.RestoreSummary.SourceBackupArn != nil {
			f13.SourceBackupARN = resp.Table.RestoreSummary.SourceBackupArn
		}
		if resp.Table.RestoreSummary.SourceTableArn != nil {
			f13.SourceTableARN = resp.Table.RestoreSummary.SourceTableArn
		}
		ko.Status.RestoreSummary = f13
	} else {
		ko.Status.RestoreSummary = nil
	}
	if resp.Table.SSEDescription != nil {
		f14 := &svcapitypes.SSEDescription{}
		if resp.Table.SSEDescription.InaccessibleEncryptionDateTime != nil {
			f14.InaccessibleEncryptionDateTime = &metav1.Time{*resp.Table.SSEDescription.InaccessibleEncryptionDateTime}
		}
		if resp.Table.SSEDescription.KMSMasterKeyArn != nil {
			f14.KMSMasterKeyARN = resp.Table.SSEDescription.KMSMasterKeyArn
		}
		if resp.Table.SSEDescription.SSEType != nil {
			f14.SSEType = resp.Table.SSEDescription.SSEType
		}
		if resp.Table.SSEDescription.Status != nil {
			f14.Status = resp.Table.SSEDescription.Status
		}
		ko.Status.SSEDescription = f14
	} else {
		ko.Status.SSEDescription = nil
	}
	if resp.Table.StreamSpecification != nil {
		f15 := &svcapitypes.StreamSpecification{}
		if resp.Table.StreamSpecification.StreamEnabled != nil {
			f15.StreamEnabled = resp.Table.StreamSpecification.StreamEnabled
		}
		if resp.Table.StreamSpecification.StreamViewType != nil {
			f15.StreamViewType = resp.Table.StreamSpecification.StreamViewType
		}
		ko.Spec.StreamSpecification = f15
	} else {
		ko.Spec.StreamSpecification = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Table.TableArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Table.TableArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.Table.TableId != nil {
		ko.Status.TableID = resp.Table.TableId
	} else {
		ko.Status.TableID = nil
	}
	if resp.Table.TableName != nil {
		ko.Spec.TableName = resp.Table.TableName
	} else {
		ko.Spec.TableName = nil
	}
	if resp.Table.TableSizeBytes != nil {
		ko.Status.TableSizeBytes = resp.Table.TableSizeBytes
	} else {
		ko.Status.TableSizeBytes = nil
	}
	if resp.Table.TableStatus != nil {
		ko.Status.TableStatus = resp.Table.TableStatus
	} else {
		ko.Status.TableStatus = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeGet, "resp", "ko", 1),
	)
}

func TestSetResource_EC2_LaunchTemplate_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "LaunchTemplate")
	require.NotNil(crd)

	// Check that we properly determined how to find the CreatedBy attribute
	// within the CreateLaunchTemplateResult shape, which has a single field called
	// "LaunchTemplate" that contains the CreatedBy field.
	expected := `
	if resp.LaunchTemplate.CreateTime != nil {
		ko.Status.CreateTime = &metav1.Time{*resp.LaunchTemplate.CreateTime}
	} else {
		ko.Status.CreateTime = nil
	}
	if resp.LaunchTemplate.CreatedBy != nil {
		ko.Status.CreatedBy = resp.LaunchTemplate.CreatedBy
	} else {
		ko.Status.CreatedBy = nil
	}
	if resp.LaunchTemplate.DefaultVersionNumber != nil {
		ko.Status.DefaultVersionNumber = resp.LaunchTemplate.DefaultVersionNumber
	} else {
		ko.Status.DefaultVersionNumber = nil
	}
	if resp.LaunchTemplate.LatestVersionNumber != nil {
		ko.Status.LatestVersionNumber = resp.LaunchTemplate.LatestVersionNumber
	} else {
		ko.Status.LatestVersionNumber = nil
	}
	if resp.LaunchTemplate.LaunchTemplateId != nil {
		ko.Status.LaunchTemplateID = resp.LaunchTemplate.LaunchTemplateId
	} else {
		ko.Status.LaunchTemplateID = nil
	}
	if resp.LaunchTemplate.LaunchTemplateName != nil {
		ko.Spec.LaunchTemplateName = resp.LaunchTemplate.LaunchTemplateName
	} else {
		ko.Spec.LaunchTemplateName = nil
	}
	if resp.LaunchTemplate.Tags != nil {
		f6 := []*svcapitypes.Tag{}
		for _, f6iter := range resp.LaunchTemplate.Tags {
			f6elem := &svcapitypes.Tag{}
			if f6iter.Key != nil {
				f6elem.Key = f6iter.Key
			}
			if f6iter.Value != nil {
				f6elem.Value = f6iter.Value
			}
			f6 = append(f6, f6elem)
		}
		ko.Status.Tags = f6
	} else {
		ko.Status.Tags = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}

func TestSetResource_ECR_Repository_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ecr")

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)

	// Check that we properly determined how to find the RegistryID attribute
	// within the CreateRepositoryOutput shape, which has a single field called
	// "Repository" that contains the RegistryId field.
	expected := `
	if resp.Repository.CreatedAt != nil {
		ko.Status.CreatedAt = &metav1.Time{*resp.Repository.CreatedAt}
	} else {
		ko.Status.CreatedAt = nil
	}
	if resp.Repository.ImageScanningConfiguration != nil {
		f1 := &svcapitypes.ImageScanningConfiguration{}
		if resp.Repository.ImageScanningConfiguration.ScanOnPush != nil {
			f1.ScanOnPush = resp.Repository.ImageScanningConfiguration.ScanOnPush
		}
		ko.Spec.ImageScanningConfiguration = f1
	} else {
		ko.Spec.ImageScanningConfiguration = nil
	}
	if resp.Repository.ImageTagMutability != nil {
		ko.Spec.ImageTagMutability = resp.Repository.ImageTagMutability
	} else {
		ko.Spec.ImageTagMutability = nil
	}
	if resp.Repository.RegistryId != nil {
		ko.Status.RegistryID = resp.Repository.RegistryId
	} else {
		ko.Status.RegistryID = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Repository.RepositoryArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Repository.RepositoryArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.Repository.RepositoryName != nil {
		ko.Spec.RepositoryName = resp.Repository.RepositoryName
	} else {
		ko.Spec.RepositoryName = nil
	}
	if resp.Repository.RepositoryUri != nil {
		ko.Status.RepositoryURI = resp.Repository.RepositoryUri
	} else {
		ko.Status.RepositoryURI = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}

func TestSetResource_ECR_Repository_ReadMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ecr")

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)

	// Check that the DescribeRepositories output is filtered by the
	// RepositoryName field of the CR's Spec, since there is no ReadOne
	// operation for ECR and we have no yet implemented a heuristic that allows
	// the DescribeRepositoriesInput.RepositoryNames filter from the CR's
	// Spec.RepositoryName field.
	expected := `
	found := false
	for _, elem := range resp.Repositories {
		if elem.CreatedAt != nil {
			ko.Status.CreatedAt = &metav1.Time{*elem.CreatedAt}
		} else {
			ko.Status.CreatedAt = nil
		}
		if elem.ImageScanningConfiguration != nil {
			f1 := &svcapitypes.ImageScanningConfiguration{}
			if elem.ImageScanningConfiguration.ScanOnPush != nil {
				f1.ScanOnPush = elem.ImageScanningConfiguration.ScanOnPush
			}
			ko.Spec.ImageScanningConfiguration = f1
		} else {
			ko.Spec.ImageScanningConfiguration = nil
		}
		if elem.ImageTagMutability != nil {
			ko.Spec.ImageTagMutability = elem.ImageTagMutability
		} else {
			ko.Spec.ImageTagMutability = nil
		}
		if elem.RegistryId != nil {
			ko.Status.RegistryID = elem.RegistryId
		} else {
			ko.Status.RegistryID = nil
		}
		if elem.RepositoryArn != nil {
			if ko.Status.ACKResourceMetadata == nil {
				ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			tmpARN := ackv1alpha1.AWSResourceName(*elem.RepositoryArn)
			ko.Status.ACKResourceMetadata.ARN = &tmpARN
		}
		if elem.RepositoryName != nil {
			if ko.Spec.RepositoryName != nil {
				if *elem.RepositoryName != *ko.Spec.RepositoryName {
					continue
				}
			}
			ko.Spec.RepositoryName = elem.RepositoryName
		} else {
			ko.Spec.RepositoryName = nil
		}
		if elem.RepositoryUri != nil {
			ko.Status.RepositoryURI = elem.RepositoryUri
		} else {
			ko.Status.RepositoryURI = nil
		}
		found = true
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeList, "resp", "ko", 1),
	)
}

func TestSetResource_Elasticache_ReplicationGroup_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")

	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	expected := `
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.ReplicationGroup.ARN != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.ReplicationGroup.ARN)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.ReplicationGroup.AtRestEncryptionEnabled != nil {
		ko.Spec.AtRestEncryptionEnabled = resp.ReplicationGroup.AtRestEncryptionEnabled
	} else {
		ko.Spec.AtRestEncryptionEnabled = nil
	}
	if resp.ReplicationGroup.AuthTokenEnabled != nil {
		ko.Status.AuthTokenEnabled = resp.ReplicationGroup.AuthTokenEnabled
	} else {
		ko.Status.AuthTokenEnabled = nil
	}
	if resp.ReplicationGroup.AuthTokenLastModifiedDate != nil {
		ko.Status.AuthTokenLastModifiedDate = &metav1.Time{*resp.ReplicationGroup.AuthTokenLastModifiedDate}
	} else {
		ko.Status.AuthTokenLastModifiedDate = nil
	}
	if resp.ReplicationGroup.AutomaticFailover != nil {
		ko.Status.AutomaticFailover = resp.ReplicationGroup.AutomaticFailover
	} else {
		ko.Status.AutomaticFailover = nil
	}
	if resp.ReplicationGroup.CacheNodeType != nil {
		ko.Spec.CacheNodeType = resp.ReplicationGroup.CacheNodeType
	} else {
		ko.Spec.CacheNodeType = nil
	}
	if resp.ReplicationGroup.ClusterEnabled != nil {
		ko.Status.ClusterEnabled = resp.ReplicationGroup.ClusterEnabled
	} else {
		ko.Status.ClusterEnabled = nil
	}
	if resp.ReplicationGroup.ConfigurationEndpoint != nil {
		f7 := &svcapitypes.Endpoint{}
		if resp.ReplicationGroup.ConfigurationEndpoint.Address != nil {
			f7.Address = resp.ReplicationGroup.ConfigurationEndpoint.Address
		}
		if resp.ReplicationGroup.ConfigurationEndpoint.Port != nil {
			f7.Port = resp.ReplicationGroup.ConfigurationEndpoint.Port
		}
		ko.Status.ConfigurationEndpoint = f7
	} else {
		ko.Status.ConfigurationEndpoint = nil
	}
	if resp.ReplicationGroup.Description != nil {
		ko.Status.Description = resp.ReplicationGroup.Description
	} else {
		ko.Status.Description = nil
	}
	if resp.ReplicationGroup.GlobalReplicationGroupInfo != nil {
		f9 := &svcapitypes.GlobalReplicationGroupInfo{}
		if resp.ReplicationGroup.GlobalReplicationGroupInfo.GlobalReplicationGroupId != nil {
			f9.GlobalReplicationGroupID = resp.ReplicationGroup.GlobalReplicationGroupInfo.GlobalReplicationGroupId
		}
		if resp.ReplicationGroup.GlobalReplicationGroupInfo.GlobalReplicationGroupMemberRole != nil {
			f9.GlobalReplicationGroupMemberRole = resp.ReplicationGroup.GlobalReplicationGroupInfo.GlobalReplicationGroupMemberRole
		}
		ko.Status.GlobalReplicationGroupInfo = f9
	} else {
		ko.Status.GlobalReplicationGroupInfo = nil
	}
	if resp.ReplicationGroup.KmsKeyId != nil {
		ko.Spec.KMSKeyID = resp.ReplicationGroup.KmsKeyId
	} else {
		ko.Spec.KMSKeyID = nil
	}
	if resp.ReplicationGroup.MemberClusters != nil {
		f12 := []*string{}
		for _, f12iter := range resp.ReplicationGroup.MemberClusters {
			var f12elem string
			f12elem = *f12iter
			f12 = append(f12, &f12elem)
		}
		ko.Status.MemberClusters = f12
	} else {
		ko.Status.MemberClusters = nil
	}
	if resp.ReplicationGroup.MemberClustersOutpostArns != nil {
		f13 := []*string{}
		for _, f13iter := range resp.ReplicationGroup.MemberClustersOutpostArns {
			var f13elem string
			f13elem = *f13iter
			f13 = append(f13, &f13elem)
		}
		ko.Status.MemberClustersOutpostARNs = f13
	} else {
		ko.Status.MemberClustersOutpostARNs = nil
	}
	if resp.ReplicationGroup.MultiAZ != nil {
		ko.Status.MultiAZ = resp.ReplicationGroup.MultiAZ
	} else {
		ko.Status.MultiAZ = nil
	}
	if resp.ReplicationGroup.NodeGroups != nil {
		f15 := []*svcapitypes.NodeGroup{}
		for _, f15iter := range resp.ReplicationGroup.NodeGroups {
			f15elem := &svcapitypes.NodeGroup{}
			if f15iter.NodeGroupId != nil {
				f15elem.NodeGroupID = f15iter.NodeGroupId
			}
			if f15iter.NodeGroupMembers != nil {
				f15elemf1 := []*svcapitypes.NodeGroupMember{}
				for _, f15elemf1iter := range f15iter.NodeGroupMembers {
					f15elemf1elem := &svcapitypes.NodeGroupMember{}
					if f15elemf1iter.CacheClusterId != nil {
						f15elemf1elem.CacheClusterID = f15elemf1iter.CacheClusterId
					}
					if f15elemf1iter.CacheNodeId != nil {
						f15elemf1elem.CacheNodeID = f15elemf1iter.CacheNodeId
					}
					if f15elemf1iter.CurrentRole != nil {
						f15elemf1elem.CurrentRole = f15elemf1iter.CurrentRole
					}
					if f15elemf1iter.PreferredAvailabilityZone != nil {
						f15elemf1elem.PreferredAvailabilityZone = f15elemf1iter.PreferredAvailabilityZone
					}
					if f15elemf1iter.PreferredOutpostArn != nil {
						f15elemf1elem.PreferredOutpostARN = f15elemf1iter.PreferredOutpostArn
					}
					if f15elemf1iter.ReadEndpoint != nil {
						f15elemf1elemf5 := &svcapitypes.Endpoint{}
						if f15elemf1iter.ReadEndpoint.Address != nil {
							f15elemf1elemf5.Address = f15elemf1iter.ReadEndpoint.Address
						}
						if f15elemf1iter.ReadEndpoint.Port != nil {
							f15elemf1elemf5.Port = f15elemf1iter.ReadEndpoint.Port
						}
						f15elemf1elem.ReadEndpoint = f15elemf1elemf5
					}
					f15elemf1 = append(f15elemf1, f15elemf1elem)
				}
				f15elem.NodeGroupMembers = f15elemf1
			}
			if f15iter.PrimaryEndpoint != nil {
				f15elemf2 := &svcapitypes.Endpoint{}
				if f15iter.PrimaryEndpoint.Address != nil {
					f15elemf2.Address = f15iter.PrimaryEndpoint.Address
				}
				if f15iter.PrimaryEndpoint.Port != nil {
					f15elemf2.Port = f15iter.PrimaryEndpoint.Port
				}
				f15elem.PrimaryEndpoint = f15elemf2
			}
			if f15iter.ReaderEndpoint != nil {
				f15elemf3 := &svcapitypes.Endpoint{}
				if f15iter.ReaderEndpoint.Address != nil {
					f15elemf3.Address = f15iter.ReaderEndpoint.Address
				}
				if f15iter.ReaderEndpoint.Port != nil {
					f15elemf3.Port = f15iter.ReaderEndpoint.Port
				}
				f15elem.ReaderEndpoint = f15elemf3
			}
			if f15iter.Slots != nil {
				f15elem.Slots = f15iter.Slots
			}
			if f15iter.Status != nil {
				f15elem.Status = f15iter.Status
			}
			f15 = append(f15, f15elem)
		}
		ko.Status.NodeGroups = f15
	} else {
		ko.Status.NodeGroups = nil
	}
	if resp.ReplicationGroup.PendingModifiedValues != nil {
		f16 := &svcapitypes.ReplicationGroupPendingModifiedValues{}
		if resp.ReplicationGroup.PendingModifiedValues.AuthTokenStatus != nil {
			f16.AuthTokenStatus = resp.ReplicationGroup.PendingModifiedValues.AuthTokenStatus
		}
		if resp.ReplicationGroup.PendingModifiedValues.AutomaticFailoverStatus != nil {
			f16.AutomaticFailoverStatus = resp.ReplicationGroup.PendingModifiedValues.AutomaticFailoverStatus
		}
		if resp.ReplicationGroup.PendingModifiedValues.LogDeliveryConfigurations != nil {
			f16f2 := []*svcapitypes.PendingLogDeliveryConfiguration{}
			for _, f16f2iter := range resp.ReplicationGroup.PendingModifiedValues.LogDeliveryConfigurations {
				f16f2elem := &svcapitypes.PendingLogDeliveryConfiguration{}
				if f16f2iter.DestinationDetails != nil {
					f16f2elemf0 := &svcapitypes.DestinationDetails{}
					if f16f2iter.DestinationDetails.CloudWatchLogsDetails != nil {
						f16f2elemf0f0 := &svcapitypes.CloudWatchLogsDestinationDetails{}
						if f16f2iter.DestinationDetails.CloudWatchLogsDetails.LogGroup != nil {
							f16f2elemf0f0.LogGroup = f16f2iter.DestinationDetails.CloudWatchLogsDetails.LogGroup
						}
						f16f2elemf0.CloudWatchLogsDetails = f16f2elemf0f0
					}
					if f16f2iter.DestinationDetails.KinesisFirehoseDetails != nil {
						f16f2elemf0f1 := &svcapitypes.KinesisFirehoseDestinationDetails{}
						if f16f2iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream != nil {
							f16f2elemf0f1.DeliveryStream = f16f2iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream
						}
						f16f2elemf0.KinesisFirehoseDetails = f16f2elemf0f1
					}
					f16f2elem.DestinationDetails = f16f2elemf0
				}
				if f16f2iter.DestinationType != nil {
					f16f2elem.DestinationType = f16f2iter.DestinationType
				}
				if f16f2iter.LogFormat != nil {
					f16f2elem.LogFormat = f16f2iter.LogFormat
				}
				if f16f2iter.LogType != nil {
					f16f2elem.LogType = f16f2iter.LogType
				}
				f16f2 = append(f16f2, f16f2elem)
			}
			f16.LogDeliveryConfigurations = f16f2
		}
		if resp.ReplicationGroup.PendingModifiedValues.PrimaryClusterId != nil {
			f16.PrimaryClusterID = resp.ReplicationGroup.PendingModifiedValues.PrimaryClusterId
		}
		if resp.ReplicationGroup.PendingModifiedValues.Resharding != nil {
			f16f4 := &svcapitypes.ReshardingStatus{}
			if resp.ReplicationGroup.PendingModifiedValues.Resharding.SlotMigration != nil {
				f16f4f0 := &svcapitypes.SlotMigration{}
				if resp.ReplicationGroup.PendingModifiedValues.Resharding.SlotMigration.ProgressPercentage != nil {
					f16f4f0.ProgressPercentage = resp.ReplicationGroup.PendingModifiedValues.Resharding.SlotMigration.ProgressPercentage
				}
				f16f4.SlotMigration = f16f4f0
			}
			f16.Resharding = f16f4
		}
		if resp.ReplicationGroup.PendingModifiedValues.UserGroups != nil {
			f16f5 := &svcapitypes.UserGroupsUpdateStatus{}
			if resp.ReplicationGroup.PendingModifiedValues.UserGroups.UserGroupIdsToAdd != nil {
				f16f5f0 := []*string{}
				for _, f16f5f0iter := range resp.ReplicationGroup.PendingModifiedValues.UserGroups.UserGroupIdsToAdd {
					var f16f5f0elem string
					f16f5f0elem = *f16f5f0iter
					f16f5f0 = append(f16f5f0, &f16f5f0elem)
				}
				f16f5.UserGroupIDsToAdd = f16f5f0
			}
			if resp.ReplicationGroup.PendingModifiedValues.UserGroups.UserGroupIdsToRemove != nil {
				f16f5f1 := []*string{}
				for _, f16f5f1iter := range resp.ReplicationGroup.PendingModifiedValues.UserGroups.UserGroupIdsToRemove {
					var f16f5f1elem string
					f16f5f1elem = *f16f5f1iter
					f16f5f1 = append(f16f5f1, &f16f5f1elem)
				}
				f16f5.UserGroupIDsToRemove = f16f5f1
			}
			f16.UserGroups = f16f5
		}
		ko.Status.PendingModifiedValues = f16
	} else {
		ko.Status.PendingModifiedValues = nil
	}
	if resp.ReplicationGroup.ReplicationGroupId != nil {
		ko.Spec.ReplicationGroupID = resp.ReplicationGroup.ReplicationGroupId
	} else {
		ko.Spec.ReplicationGroupID = nil
	}
	if resp.ReplicationGroup.SnapshotRetentionLimit != nil {
		ko.Spec.SnapshotRetentionLimit = resp.ReplicationGroup.SnapshotRetentionLimit
	} else {
		ko.Spec.SnapshotRetentionLimit = nil
	}
	if resp.ReplicationGroup.SnapshotWindow != nil {
		ko.Spec.SnapshotWindow = resp.ReplicationGroup.SnapshotWindow
	} else {
		ko.Spec.SnapshotWindow = nil
	}
	if resp.ReplicationGroup.SnapshottingClusterId != nil {
		ko.Status.SnapshottingClusterID = resp.ReplicationGroup.SnapshottingClusterId
	} else {
		ko.Status.SnapshottingClusterID = nil
	}
	if resp.ReplicationGroup.Status != nil {
		ko.Status.Status = resp.ReplicationGroup.Status
	} else {
		ko.Status.Status = nil
	}
	if resp.ReplicationGroup.TransitEncryptionEnabled != nil {
		ko.Spec.TransitEncryptionEnabled = resp.ReplicationGroup.TransitEncryptionEnabled
	} else {
		ko.Spec.TransitEncryptionEnabled = nil
	}
	if resp.ReplicationGroup.UserGroupIds != nil {
		f23 := []*string{}
		for _, f23iter := range resp.ReplicationGroup.UserGroupIds {
			var f23elem string
			f23elem = *f23iter
			f23 = append(f23, &f23elem)
		}
		ko.Spec.UserGroupIDs = f23
	} else {
		ko.Spec.UserGroupIDs = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}

func TestSetResource_Elasticache_ReplicationGroup_ReadMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")

	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	expected := `
	found := false
	for _, elem := range resp.ReplicationGroups {
		if elem.ARN != nil {
			if ko.Status.ACKResourceMetadata == nil {
				ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			tmpARN := ackv1alpha1.AWSResourceName(*elem.ARN)
			ko.Status.ACKResourceMetadata.ARN = &tmpARN
		}
		if elem.AtRestEncryptionEnabled != nil {
			ko.Spec.AtRestEncryptionEnabled = elem.AtRestEncryptionEnabled
		} else {
			ko.Spec.AtRestEncryptionEnabled = nil
		}
		if elem.AuthTokenEnabled != nil {
			ko.Status.AuthTokenEnabled = elem.AuthTokenEnabled
		} else {
			ko.Status.AuthTokenEnabled = nil
		}
		if elem.AuthTokenLastModifiedDate != nil {
			ko.Status.AuthTokenLastModifiedDate = &metav1.Time{*elem.AuthTokenLastModifiedDate}
		} else {
			ko.Status.AuthTokenLastModifiedDate = nil
		}
		if elem.AutomaticFailover != nil {
			ko.Status.AutomaticFailover = elem.AutomaticFailover
		} else {
			ko.Status.AutomaticFailover = nil
		}
		if elem.CacheNodeType != nil {
			ko.Spec.CacheNodeType = elem.CacheNodeType
		} else {
			ko.Spec.CacheNodeType = nil
		}
		if elem.ClusterEnabled != nil {
			ko.Status.ClusterEnabled = elem.ClusterEnabled
		} else {
			ko.Status.ClusterEnabled = nil
		}
		if elem.ConfigurationEndpoint != nil {
			f7 := &svcapitypes.Endpoint{}
			if elem.ConfigurationEndpoint.Address != nil {
				f7.Address = elem.ConfigurationEndpoint.Address
			}
			if elem.ConfigurationEndpoint.Port != nil {
				f7.Port = elem.ConfigurationEndpoint.Port
			}
			ko.Status.ConfigurationEndpoint = f7
		} else {
			ko.Status.ConfigurationEndpoint = nil
		}
		if elem.Description != nil {
			ko.Status.Description = elem.Description
		} else {
			ko.Status.Description = nil
		}
		if elem.GlobalReplicationGroupInfo != nil {
			f9 := &svcapitypes.GlobalReplicationGroupInfo{}
			if elem.GlobalReplicationGroupInfo.GlobalReplicationGroupId != nil {
				f9.GlobalReplicationGroupID = elem.GlobalReplicationGroupInfo.GlobalReplicationGroupId
			}
			if elem.GlobalReplicationGroupInfo.GlobalReplicationGroupMemberRole != nil {
				f9.GlobalReplicationGroupMemberRole = elem.GlobalReplicationGroupInfo.GlobalReplicationGroupMemberRole
			}
			ko.Status.GlobalReplicationGroupInfo = f9
		} else {
			ko.Status.GlobalReplicationGroupInfo = nil
		}
		if elem.KmsKeyId != nil {
			ko.Spec.KMSKeyID = elem.KmsKeyId
		} else {
			ko.Spec.KMSKeyID = nil
		}
		if elem.LogDeliveryConfigurations != nil {
			f11 := []*svcapitypes.LogDeliveryConfigurationRequest{}
			for _, f11iter := range elem.LogDeliveryConfigurations {
				f11elem := &svcapitypes.LogDeliveryConfigurationRequest{}
				if f11iter.DestinationDetails != nil {
					f11elemf0 := &svcapitypes.DestinationDetails{}
					if f11iter.DestinationDetails.CloudWatchLogsDetails != nil {
						f11elemf0f0 := &svcapitypes.CloudWatchLogsDestinationDetails{}
						if f11iter.DestinationDetails.CloudWatchLogsDetails.LogGroup != nil {
							f11elemf0f0.LogGroup = f11iter.DestinationDetails.CloudWatchLogsDetails.LogGroup
						}
						f11elemf0.CloudWatchLogsDetails = f11elemf0f0
					}
					if f11iter.DestinationDetails.KinesisFirehoseDetails != nil {
						f11elemf0f1 := &svcapitypes.KinesisFirehoseDestinationDetails{}
						if f11iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream != nil {
							f11elemf0f1.DeliveryStream = f11iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream
						}
						f11elemf0.KinesisFirehoseDetails = f11elemf0f1
					}
					f11elem.DestinationDetails = f11elemf0
				}
				if f11iter.DestinationType != nil {
					f11elem.DestinationType = f11iter.DestinationType
				}
				if f11iter.LogFormat != nil {
					f11elem.LogFormat = f11iter.LogFormat
				}
				if f11iter.LogType != nil {
					f11elem.LogType = f11iter.LogType
				}
				f11 = append(f11, f11elem)
			}
			ko.Spec.LogDeliveryConfigurations = f11
		} else {
			ko.Spec.LogDeliveryConfigurations = nil
		}
		if elem.MemberClusters != nil {
			f12 := []*string{}
			for _, f12iter := range elem.MemberClusters {
				var f12elem string
				f12elem = *f12iter
				f12 = append(f12, &f12elem)
			}
			ko.Status.MemberClusters = f12
		} else {
			ko.Status.MemberClusters = nil
		}
		if elem.MemberClustersOutpostArns != nil {
			f13 := []*string{}
			for _, f13iter := range elem.MemberClustersOutpostArns {
				var f13elem string
				f13elem = *f13iter
				f13 = append(f13, &f13elem)
			}
			ko.Status.MemberClustersOutpostARNs = f13
		} else {
			ko.Status.MemberClustersOutpostARNs = nil
		}
		if elem.MultiAZ != nil {
			ko.Status.MultiAZ = elem.MultiAZ
		} else {
			ko.Status.MultiAZ = nil
		}
		if elem.NodeGroups != nil {
			f15 := []*svcapitypes.NodeGroup{}
			for _, f15iter := range elem.NodeGroups {
				f15elem := &svcapitypes.NodeGroup{}
				if f15iter.NodeGroupId != nil {
					f15elem.NodeGroupID = f15iter.NodeGroupId
				}
				if f15iter.NodeGroupMembers != nil {
					f15elemf1 := []*svcapitypes.NodeGroupMember{}
					for _, f15elemf1iter := range f15iter.NodeGroupMembers {
						f15elemf1elem := &svcapitypes.NodeGroupMember{}
						if f15elemf1iter.CacheClusterId != nil {
							f15elemf1elem.CacheClusterID = f15elemf1iter.CacheClusterId
						}
						if f15elemf1iter.CacheNodeId != nil {
							f15elemf1elem.CacheNodeID = f15elemf1iter.CacheNodeId
						}
						if f15elemf1iter.CurrentRole != nil {
							f15elemf1elem.CurrentRole = f15elemf1iter.CurrentRole
						}
						if f15elemf1iter.PreferredAvailabilityZone != nil {
							f15elemf1elem.PreferredAvailabilityZone = f15elemf1iter.PreferredAvailabilityZone
						}
						if f15elemf1iter.PreferredOutpostArn != nil {
							f15elemf1elem.PreferredOutpostARN = f15elemf1iter.PreferredOutpostArn
						}
						if f15elemf1iter.ReadEndpoint != nil {
							f15elemf1elemf5 := &svcapitypes.Endpoint{}
							if f15elemf1iter.ReadEndpoint.Address != nil {
								f15elemf1elemf5.Address = f15elemf1iter.ReadEndpoint.Address
							}
							if f15elemf1iter.ReadEndpoint.Port != nil {
								f15elemf1elemf5.Port = f15elemf1iter.ReadEndpoint.Port
							}
							f15elemf1elem.ReadEndpoint = f15elemf1elemf5
						}
						f15elemf1 = append(f15elemf1, f15elemf1elem)
					}
					f15elem.NodeGroupMembers = f15elemf1
				}
				if f15iter.PrimaryEndpoint != nil {
					f15elemf2 := &svcapitypes.Endpoint{}
					if f15iter.PrimaryEndpoint.Address != nil {
						f15elemf2.Address = f15iter.PrimaryEndpoint.Address
					}
					if f15iter.PrimaryEndpoint.Port != nil {
						f15elemf2.Port = f15iter.PrimaryEndpoint.Port
					}
					f15elem.PrimaryEndpoint = f15elemf2
				}
				if f15iter.ReaderEndpoint != nil {
					f15elemf3 := &svcapitypes.Endpoint{}
					if f15iter.ReaderEndpoint.Address != nil {
						f15elemf3.Address = f15iter.ReaderEndpoint.Address
					}
					if f15iter.ReaderEndpoint.Port != nil {
						f15elemf3.Port = f15iter.ReaderEndpoint.Port
					}
					f15elem.ReaderEndpoint = f15elemf3
				}
				if f15iter.Slots != nil {
					f15elem.Slots = f15iter.Slots
				}
				if f15iter.Status != nil {
					f15elem.Status = f15iter.Status
				}
				f15 = append(f15, f15elem)
			}
			ko.Status.NodeGroups = f15
		} else {
			ko.Status.NodeGroups = nil
		}
		if elem.PendingModifiedValues != nil {
			f16 := &svcapitypes.ReplicationGroupPendingModifiedValues{}
			if elem.PendingModifiedValues.AuthTokenStatus != nil {
				f16.AuthTokenStatus = elem.PendingModifiedValues.AuthTokenStatus
			}
			if elem.PendingModifiedValues.AutomaticFailoverStatus != nil {
				f16.AutomaticFailoverStatus = elem.PendingModifiedValues.AutomaticFailoverStatus
			}
			if elem.PendingModifiedValues.LogDeliveryConfigurations != nil {
				f16f2 := []*svcapitypes.PendingLogDeliveryConfiguration{}
				for _, f16f2iter := range elem.PendingModifiedValues.LogDeliveryConfigurations {
					f16f2elem := &svcapitypes.PendingLogDeliveryConfiguration{}
					if f16f2iter.DestinationDetails != nil {
						f16f2elemf0 := &svcapitypes.DestinationDetails{}
						if f16f2iter.DestinationDetails.CloudWatchLogsDetails != nil {
							f16f2elemf0f0 := &svcapitypes.CloudWatchLogsDestinationDetails{}
							if f16f2iter.DestinationDetails.CloudWatchLogsDetails.LogGroup != nil {
								f16f2elemf0f0.LogGroup = f16f2iter.DestinationDetails.CloudWatchLogsDetails.LogGroup
							}
							f16f2elemf0.CloudWatchLogsDetails = f16f2elemf0f0
						}
						if f16f2iter.DestinationDetails.KinesisFirehoseDetails != nil {
							f16f2elemf0f1 := &svcapitypes.KinesisFirehoseDestinationDetails{}
							if f16f2iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream != nil {
								f16f2elemf0f1.DeliveryStream = f16f2iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream
							}
							f16f2elemf0.KinesisFirehoseDetails = f16f2elemf0f1
						}
						f16f2elem.DestinationDetails = f16f2elemf0
					}
					if f16f2iter.DestinationType != nil {
						f16f2elem.DestinationType = f16f2iter.DestinationType
					}
					if f16f2iter.LogFormat != nil {
						f16f2elem.LogFormat = f16f2iter.LogFormat
					}
					if f16f2iter.LogType != nil {
						f16f2elem.LogType = f16f2iter.LogType
					}
					f16f2 = append(f16f2, f16f2elem)
				}
				f16.LogDeliveryConfigurations = f16f2
			}
			if elem.PendingModifiedValues.PrimaryClusterId != nil {
				f16.PrimaryClusterID = elem.PendingModifiedValues.PrimaryClusterId
			}
			if elem.PendingModifiedValues.Resharding != nil {
				f16f4 := &svcapitypes.ReshardingStatus{}
				if elem.PendingModifiedValues.Resharding.SlotMigration != nil {
					f16f4f0 := &svcapitypes.SlotMigration{}
					if elem.PendingModifiedValues.Resharding.SlotMigration.ProgressPercentage != nil {
						f16f4f0.ProgressPercentage = elem.PendingModifiedValues.Resharding.SlotMigration.ProgressPercentage
					}
					f16f4.SlotMigration = f16f4f0
				}
				f16.Resharding = f16f4
			}
			if elem.PendingModifiedValues.UserGroups != nil {
				f16f5 := &svcapitypes.UserGroupsUpdateStatus{}
				if elem.PendingModifiedValues.UserGroups.UserGroupIdsToAdd != nil {
					f16f5f0 := []*string{}
					for _, f16f5f0iter := range elem.PendingModifiedValues.UserGroups.UserGroupIdsToAdd {
						var f16f5f0elem string
						f16f5f0elem = *f16f5f0iter
						f16f5f0 = append(f16f5f0, &f16f5f0elem)
					}
					f16f5.UserGroupIDsToAdd = f16f5f0
				}
				if elem.PendingModifiedValues.UserGroups.UserGroupIdsToRemove != nil {
					f16f5f1 := []*string{}
					for _, f16f5f1iter := range elem.PendingModifiedValues.UserGroups.UserGroupIdsToRemove {
						var f16f5f1elem string
						f16f5f1elem = *f16f5f1iter
						f16f5f1 = append(f16f5f1, &f16f5f1elem)
					}
					f16f5.UserGroupIDsToRemove = f16f5f1
				}
				f16.UserGroups = f16f5
			}
			ko.Status.PendingModifiedValues = f16
		} else {
			ko.Status.PendingModifiedValues = nil
		}
		if elem.ReplicationGroupId != nil {
			ko.Spec.ReplicationGroupID = elem.ReplicationGroupId
		} else {
			ko.Spec.ReplicationGroupID = nil
		}
		if elem.SnapshotRetentionLimit != nil {
			ko.Spec.SnapshotRetentionLimit = elem.SnapshotRetentionLimit
		} else {
			ko.Spec.SnapshotRetentionLimit = nil
		}
		if elem.SnapshotWindow != nil {
			ko.Spec.SnapshotWindow = elem.SnapshotWindow
		} else {
			ko.Spec.SnapshotWindow = nil
		}
		if elem.SnapshottingClusterId != nil {
			ko.Status.SnapshottingClusterID = elem.SnapshottingClusterId
		} else {
			ko.Status.SnapshottingClusterID = nil
		}
		if elem.Status != nil {
			ko.Status.Status = elem.Status
		} else {
			ko.Status.Status = nil
		}
		if elem.TransitEncryptionEnabled != nil {
			ko.Spec.TransitEncryptionEnabled = elem.TransitEncryptionEnabled
		} else {
			ko.Spec.TransitEncryptionEnabled = nil
		}
		if elem.UserGroupIds != nil {
			f23 := []*string{}
			for _, f23iter := range elem.UserGroupIds {
				var f23elem string
				f23elem = *f23iter
				f23 = append(f23, &f23elem)
			}
			ko.Spec.UserGroupIDs = f23
		} else {
			ko.Spec.UserGroupIDs = nil
		}
		found = true
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeList, "resp", "ko", 1),
	)
}

func TestSetResource_RDS_DBInstance_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "rds")

	crd := testutil.GetCRDByName(t, g, "DBInstance")
	require.NotNil(crd)

	expected := `
	if resp.DBInstance.AllocatedStorage != nil {
		ko.Spec.AllocatedStorage = resp.DBInstance.AllocatedStorage
	} else {
		ko.Spec.AllocatedStorage = nil
	}
	if resp.DBInstance.AssociatedRoles != nil {
		f1 := []*svcapitypes.DBInstanceRole{}
		for _, f1iter := range resp.DBInstance.AssociatedRoles {
			f1elem := &svcapitypes.DBInstanceRole{}
			if f1iter.FeatureName != nil {
				f1elem.FeatureName = f1iter.FeatureName
			}
			if f1iter.RoleArn != nil {
				f1elem.RoleARN = f1iter.RoleArn
			}
			if f1iter.Status != nil {
				f1elem.Status = f1iter.Status
			}
			f1 = append(f1, f1elem)
		}
		ko.Status.AssociatedRoles = f1
	} else {
		ko.Status.AssociatedRoles = nil
	}
	if resp.DBInstance.AutoMinorVersionUpgrade != nil {
		ko.Spec.AutoMinorVersionUpgrade = resp.DBInstance.AutoMinorVersionUpgrade
	} else {
		ko.Spec.AutoMinorVersionUpgrade = nil
	}
	if resp.DBInstance.AvailabilityZone != nil {
		ko.Spec.AvailabilityZone = resp.DBInstance.AvailabilityZone
	} else {
		ko.Spec.AvailabilityZone = nil
	}
	if resp.DBInstance.BackupRetentionPeriod != nil {
		ko.Spec.BackupRetentionPeriod = resp.DBInstance.BackupRetentionPeriod
	} else {
		ko.Spec.BackupRetentionPeriod = nil
	}
	if resp.DBInstance.CACertificateIdentifier != nil {
		ko.Status.CACertificateIdentifier = resp.DBInstance.CACertificateIdentifier
	} else {
		ko.Status.CACertificateIdentifier = nil
	}
	if resp.DBInstance.CharacterSetName != nil {
		ko.Spec.CharacterSetName = resp.DBInstance.CharacterSetName
	} else {
		ko.Spec.CharacterSetName = nil
	}
	if resp.DBInstance.CopyTagsToSnapshot != nil {
		ko.Spec.CopyTagsToSnapshot = resp.DBInstance.CopyTagsToSnapshot
	} else {
		ko.Spec.CopyTagsToSnapshot = nil
	}
	if resp.DBInstance.DBClusterIdentifier != nil {
		ko.Spec.DBClusterIdentifier = resp.DBInstance.DBClusterIdentifier
	} else {
		ko.Spec.DBClusterIdentifier = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.DBInstance.DBInstanceArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.DBInstance.DBInstanceArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.DBInstance.DBInstanceClass != nil {
		ko.Spec.DBInstanceClass = resp.DBInstance.DBInstanceClass
	} else {
		ko.Spec.DBInstanceClass = nil
	}
	if resp.DBInstance.DBInstanceIdentifier != nil {
		ko.Spec.DBInstanceIdentifier = resp.DBInstance.DBInstanceIdentifier
	} else {
		ko.Spec.DBInstanceIdentifier = nil
	}
	if resp.DBInstance.DBInstanceStatus != nil {
		ko.Status.DBInstanceStatus = resp.DBInstance.DBInstanceStatus
	} else {
		ko.Status.DBInstanceStatus = nil
	}
	if resp.DBInstance.DBName != nil {
		ko.Spec.DBName = resp.DBInstance.DBName
	} else {
		ko.Spec.DBName = nil
	}
	if resp.DBInstance.DBParameterGroups != nil {
		f14 := []*svcapitypes.DBParameterGroupStatus_SDK{}
		for _, f14iter := range resp.DBInstance.DBParameterGroups {
			f14elem := &svcapitypes.DBParameterGroupStatus_SDK{}
			if f14iter.DBParameterGroupName != nil {
				f14elem.DBParameterGroupName = f14iter.DBParameterGroupName
			}
			if f14iter.ParameterApplyStatus != nil {
				f14elem.ParameterApplyStatus = f14iter.ParameterApplyStatus
			}
			f14 = append(f14, f14elem)
		}
		ko.Status.DBParameterGroups = f14
	} else {
		ko.Status.DBParameterGroups = nil
	}
	if resp.DBInstance.DBSecurityGroups != nil {
		f15 := []*string{}
		for _, f15iter := range resp.DBInstance.DBSecurityGroups {
			var f15elem string
			f15elem = *f15iter.DBSecurityGroupName
			f15 = append(f15, &f15elem)
		}
		ko.Spec.DBSecurityGroups = f15
	} else {
		ko.Spec.DBSecurityGroups = nil
	}
	if resp.DBInstance.DBSubnetGroup != nil {
		f16 := &svcapitypes.DBSubnetGroup_SDK{}
		if resp.DBInstance.DBSubnetGroup.DBSubnetGroupArn != nil {
			f16.DBSubnetGroupARN = resp.DBInstance.DBSubnetGroup.DBSubnetGroupArn
		}
		if resp.DBInstance.DBSubnetGroup.DBSubnetGroupDescription != nil {
			f16.DBSubnetGroupDescription = resp.DBInstance.DBSubnetGroup.DBSubnetGroupDescription
		}
		if resp.DBInstance.DBSubnetGroup.DBSubnetGroupName != nil {
			f16.DBSubnetGroupName = resp.DBInstance.DBSubnetGroup.DBSubnetGroupName
		}
		if resp.DBInstance.DBSubnetGroup.SubnetGroupStatus != nil {
			f16.SubnetGroupStatus = resp.DBInstance.DBSubnetGroup.SubnetGroupStatus
		}
		if resp.DBInstance.DBSubnetGroup.Subnets != nil {
			f16f4 := []*svcapitypes.Subnet{}
			for _, f16f4iter := range resp.DBInstance.DBSubnetGroup.Subnets {
				f16f4elem := &svcapitypes.Subnet{}
				if f16f4iter.SubnetAvailabilityZone != nil {
					f16f4elemf0 := &svcapitypes.AvailabilityZone{}
					if f16f4iter.SubnetAvailabilityZone.Name != nil {
						f16f4elemf0.Name = f16f4iter.SubnetAvailabilityZone.Name
					}
					f16f4elem.SubnetAvailabilityZone = f16f4elemf0
				}
				if f16f4iter.SubnetIdentifier != nil {
					f16f4elem.SubnetIdentifier = f16f4iter.SubnetIdentifier
				}
				if f16f4iter.SubnetOutpost != nil {
					f16f4elemf2 := &svcapitypes.Outpost{}
					if f16f4iter.SubnetOutpost.Arn != nil {
						f16f4elemf2.ARN = f16f4iter.SubnetOutpost.Arn
					}
					f16f4elem.SubnetOutpost = f16f4elemf2
				}
				if f16f4iter.SubnetStatus != nil {
					f16f4elem.SubnetStatus = f16f4iter.SubnetStatus
				}
				f16f4 = append(f16f4, f16f4elem)
			}
			f16.Subnets = f16f4
		}
		if resp.DBInstance.DBSubnetGroup.VpcId != nil {
			f16.VPCID = resp.DBInstance.DBSubnetGroup.VpcId
		}
		ko.Status.DBSubnetGroup = f16
	} else {
		ko.Status.DBSubnetGroup = nil
	}
	if resp.DBInstance.DbInstancePort != nil {
		ko.Status.DBInstancePort = resp.DBInstance.DbInstancePort
	} else {
		ko.Status.DBInstancePort = nil
	}
	if resp.DBInstance.DbiResourceId != nil {
		ko.Status.DBIResourceID = resp.DBInstance.DbiResourceId
	} else {
		ko.Status.DBIResourceID = nil
	}
	if resp.DBInstance.DeletionProtection != nil {
		ko.Spec.DeletionProtection = resp.DBInstance.DeletionProtection
	} else {
		ko.Spec.DeletionProtection = nil
	}
	if resp.DBInstance.DomainMemberships != nil {
		f20 := []*svcapitypes.DomainMembership{}
		for _, f20iter := range resp.DBInstance.DomainMemberships {
			f20elem := &svcapitypes.DomainMembership{}
			if f20iter.Domain != nil {
				f20elem.Domain = f20iter.Domain
			}
			if f20iter.FQDN != nil {
				f20elem.FQDN = f20iter.FQDN
			}
			if f20iter.IAMRoleName != nil {
				f20elem.IAMRoleName = f20iter.IAMRoleName
			}
			if f20iter.Status != nil {
				f20elem.Status = f20iter.Status
			}
			f20 = append(f20, f20elem)
		}
		ko.Status.DomainMemberships = f20
	} else {
		ko.Status.DomainMemberships = nil
	}
	if resp.DBInstance.EnabledCloudwatchLogsExports != nil {
		f21 := []*string{}
		for _, f21iter := range resp.DBInstance.EnabledCloudwatchLogsExports {
			var f21elem string
			f21elem = *f21iter
			f21 = append(f21, &f21elem)
		}
		ko.Status.EnabledCloudwatchLogsExports = f21
	} else {
		ko.Status.EnabledCloudwatchLogsExports = nil
	}
	if resp.DBInstance.Endpoint != nil {
		f22 := &svcapitypes.Endpoint{}
		if resp.DBInstance.Endpoint.Address != nil {
			f22.Address = resp.DBInstance.Endpoint.Address
		}
		if resp.DBInstance.Endpoint.HostedZoneId != nil {
			f22.HostedZoneID = resp.DBInstance.Endpoint.HostedZoneId
		}
		if resp.DBInstance.Endpoint.Port != nil {
			f22.Port = resp.DBInstance.Endpoint.Port
		}
		ko.Status.Endpoint = f22
	} else {
		ko.Status.Endpoint = nil
	}
	if resp.DBInstance.Engine != nil {
		ko.Spec.Engine = resp.DBInstance.Engine
	} else {
		ko.Spec.Engine = nil
	}
	if resp.DBInstance.EngineVersion != nil {
		ko.Spec.EngineVersion = resp.DBInstance.EngineVersion
	} else {
		ko.Spec.EngineVersion = nil
	}
	if resp.DBInstance.EnhancedMonitoringResourceArn != nil {
		ko.Status.EnhancedMonitoringResourceARN = resp.DBInstance.EnhancedMonitoringResourceArn
	} else {
		ko.Status.EnhancedMonitoringResourceARN = nil
	}
	if resp.DBInstance.IAMDatabaseAuthenticationEnabled != nil {
		ko.Status.IAMDatabaseAuthenticationEnabled = resp.DBInstance.IAMDatabaseAuthenticationEnabled
	} else {
		ko.Status.IAMDatabaseAuthenticationEnabled = nil
	}
	if resp.DBInstance.InstanceCreateTime != nil {
		ko.Status.InstanceCreateTime = &metav1.Time{*resp.DBInstance.InstanceCreateTime}
	} else {
		ko.Status.InstanceCreateTime = nil
	}
	if resp.DBInstance.Iops != nil {
		ko.Spec.IOPS = resp.DBInstance.Iops
	} else {
		ko.Spec.IOPS = nil
	}
	if resp.DBInstance.KmsKeyId != nil {
		ko.Spec.KMSKeyID = resp.DBInstance.KmsKeyId
	} else {
		ko.Spec.KMSKeyID = nil
	}
	if resp.DBInstance.LatestRestorableTime != nil {
		ko.Status.LatestRestorableTime = &metav1.Time{*resp.DBInstance.LatestRestorableTime}
	} else {
		ko.Status.LatestRestorableTime = nil
	}
	if resp.DBInstance.LicenseModel != nil {
		ko.Spec.LicenseModel = resp.DBInstance.LicenseModel
	} else {
		ko.Spec.LicenseModel = nil
	}
	if resp.DBInstance.ListenerEndpoint != nil {
		f32 := &svcapitypes.Endpoint{}
		if resp.DBInstance.ListenerEndpoint.Address != nil {
			f32.Address = resp.DBInstance.ListenerEndpoint.Address
		}
		if resp.DBInstance.ListenerEndpoint.HostedZoneId != nil {
			f32.HostedZoneID = resp.DBInstance.ListenerEndpoint.HostedZoneId
		}
		if resp.DBInstance.ListenerEndpoint.Port != nil {
			f32.Port = resp.DBInstance.ListenerEndpoint.Port
		}
		ko.Status.ListenerEndpoint = f32
	} else {
		ko.Status.ListenerEndpoint = nil
	}
	if resp.DBInstance.MasterUsername != nil {
		ko.Spec.MasterUsername = resp.DBInstance.MasterUsername
	} else {
		ko.Spec.MasterUsername = nil
	}
	if resp.DBInstance.MaxAllocatedStorage != nil {
		ko.Spec.MaxAllocatedStorage = resp.DBInstance.MaxAllocatedStorage
	} else {
		ko.Spec.MaxAllocatedStorage = nil
	}
	if resp.DBInstance.MonitoringInterval != nil {
		ko.Spec.MonitoringInterval = resp.DBInstance.MonitoringInterval
	} else {
		ko.Spec.MonitoringInterval = nil
	}
	if resp.DBInstance.MonitoringRoleArn != nil {
		ko.Spec.MonitoringRoleARN = resp.DBInstance.MonitoringRoleArn
	} else {
		ko.Spec.MonitoringRoleARN = nil
	}
	if resp.DBInstance.MultiAZ != nil {
		ko.Spec.MultiAZ = resp.DBInstance.MultiAZ
	} else {
		ko.Spec.MultiAZ = nil
	}
	if resp.DBInstance.OptionGroupMemberships != nil {
		f38 := []*svcapitypes.OptionGroupMembership{}
		for _, f38iter := range resp.DBInstance.OptionGroupMemberships {
			f38elem := &svcapitypes.OptionGroupMembership{}
			if f38iter.OptionGroupName != nil {
				f38elem.OptionGroupName = f38iter.OptionGroupName
			}
			if f38iter.Status != nil {
				f38elem.Status = f38iter.Status
			}
			f38 = append(f38, f38elem)
		}
		ko.Status.OptionGroupMemberships = f38
	} else {
		ko.Status.OptionGroupMemberships = nil
	}
	if resp.DBInstance.PendingModifiedValues != nil {
		f39 := &svcapitypes.PendingModifiedValues{}
		if resp.DBInstance.PendingModifiedValues.AllocatedStorage != nil {
			f39.AllocatedStorage = resp.DBInstance.PendingModifiedValues.AllocatedStorage
		}
		if resp.DBInstance.PendingModifiedValues.BackupRetentionPeriod != nil {
			f39.BackupRetentionPeriod = resp.DBInstance.PendingModifiedValues.BackupRetentionPeriod
		}
		if resp.DBInstance.PendingModifiedValues.CACertificateIdentifier != nil {
			f39.CACertificateIdentifier = resp.DBInstance.PendingModifiedValues.CACertificateIdentifier
		}
		if resp.DBInstance.PendingModifiedValues.DBInstanceClass != nil {
			f39.DBInstanceClass = resp.DBInstance.PendingModifiedValues.DBInstanceClass
		}
		if resp.DBInstance.PendingModifiedValues.DBInstanceIdentifier != nil {
			f39.DBInstanceIdentifier = resp.DBInstance.PendingModifiedValues.DBInstanceIdentifier
		}
		if resp.DBInstance.PendingModifiedValues.DBSubnetGroupName != nil {
			f39.DBSubnetGroupName = resp.DBInstance.PendingModifiedValues.DBSubnetGroupName
		}
		if resp.DBInstance.PendingModifiedValues.EngineVersion != nil {
			f39.EngineVersion = resp.DBInstance.PendingModifiedValues.EngineVersion
		}
		if resp.DBInstance.PendingModifiedValues.Iops != nil {
			f39.IOPS = resp.DBInstance.PendingModifiedValues.Iops
		}
		if resp.DBInstance.PendingModifiedValues.LicenseModel != nil {
			f39.LicenseModel = resp.DBInstance.PendingModifiedValues.LicenseModel
		}
		if resp.DBInstance.PendingModifiedValues.MasterUserPassword != nil {
			f39.MasterUserPassword = resp.DBInstance.PendingModifiedValues.MasterUserPassword
		}
		if resp.DBInstance.PendingModifiedValues.MultiAZ != nil {
			f39.MultiAZ = resp.DBInstance.PendingModifiedValues.MultiAZ
		}
		if resp.DBInstance.PendingModifiedValues.PendingCloudwatchLogsExports != nil {
			f39f11 := &svcapitypes.PendingCloudwatchLogsExports{}
			if resp.DBInstance.PendingModifiedValues.PendingCloudwatchLogsExports.LogTypesToDisable != nil {
				f39f11f0 := []*string{}
				for _, f39f11f0iter := range resp.DBInstance.PendingModifiedValues.PendingCloudwatchLogsExports.LogTypesToDisable {
					var f39f11f0elem string
					f39f11f0elem = *f39f11f0iter
					f39f11f0 = append(f39f11f0, &f39f11f0elem)
				}
				f39f11.LogTypesToDisable = f39f11f0
			}
			if resp.DBInstance.PendingModifiedValues.PendingCloudwatchLogsExports.LogTypesToEnable != nil {
				f39f11f1 := []*string{}
				for _, f39f11f1iter := range resp.DBInstance.PendingModifiedValues.PendingCloudwatchLogsExports.LogTypesToEnable {
					var f39f11f1elem string
					f39f11f1elem = *f39f11f1iter
					f39f11f1 = append(f39f11f1, &f39f11f1elem)
				}
				f39f11.LogTypesToEnable = f39f11f1
			}
			f39.PendingCloudwatchLogsExports = f39f11
		}
		if resp.DBInstance.PendingModifiedValues.Port != nil {
			f39.Port = resp.DBInstance.PendingModifiedValues.Port
		}
		if resp.DBInstance.PendingModifiedValues.ProcessorFeatures != nil {
			f39f13 := []*svcapitypes.ProcessorFeature{}
			for _, f39f13iter := range resp.DBInstance.PendingModifiedValues.ProcessorFeatures {
				f39f13elem := &svcapitypes.ProcessorFeature{}
				if f39f13iter.Name != nil {
					f39f13elem.Name = f39f13iter.Name
				}
				if f39f13iter.Value != nil {
					f39f13elem.Value = f39f13iter.Value
				}
				f39f13 = append(f39f13, f39f13elem)
			}
			f39.ProcessorFeatures = f39f13
		}
		if resp.DBInstance.PendingModifiedValues.StorageType != nil {
			f39.StorageType = resp.DBInstance.PendingModifiedValues.StorageType
		}
		ko.Status.PendingModifiedValues = f39
	} else {
		ko.Status.PendingModifiedValues = nil
	}
	if resp.DBInstance.PerformanceInsightsEnabled != nil {
		ko.Status.PerformanceInsightsEnabled = resp.DBInstance.PerformanceInsightsEnabled
	} else {
		ko.Status.PerformanceInsightsEnabled = nil
	}
	if resp.DBInstance.PerformanceInsightsKMSKeyId != nil {
		ko.Spec.PerformanceInsightsKMSKeyID = resp.DBInstance.PerformanceInsightsKMSKeyId
	} else {
		ko.Spec.PerformanceInsightsKMSKeyID = nil
	}
	if resp.DBInstance.PerformanceInsightsRetentionPeriod != nil {
		ko.Spec.PerformanceInsightsRetentionPeriod = resp.DBInstance.PerformanceInsightsRetentionPeriod
	} else {
		ko.Spec.PerformanceInsightsRetentionPeriod = nil
	}
	if resp.DBInstance.PreferredBackupWindow != nil {
		ko.Spec.PreferredBackupWindow = resp.DBInstance.PreferredBackupWindow
	} else {
		ko.Spec.PreferredBackupWindow = nil
	}
	if resp.DBInstance.PreferredMaintenanceWindow != nil {
		ko.Spec.PreferredMaintenanceWindow = resp.DBInstance.PreferredMaintenanceWindow
	} else {
		ko.Spec.PreferredMaintenanceWindow = nil
	}
	if resp.DBInstance.ProcessorFeatures != nil {
		f45 := []*svcapitypes.ProcessorFeature{}
		for _, f45iter := range resp.DBInstance.ProcessorFeatures {
			f45elem := &svcapitypes.ProcessorFeature{}
			if f45iter.Name != nil {
				f45elem.Name = f45iter.Name
			}
			if f45iter.Value != nil {
				f45elem.Value = f45iter.Value
			}
			f45 = append(f45, f45elem)
		}
		ko.Spec.ProcessorFeatures = f45
	} else {
		ko.Spec.ProcessorFeatures = nil
	}
	if resp.DBInstance.PromotionTier != nil {
		ko.Spec.PromotionTier = resp.DBInstance.PromotionTier
	} else {
		ko.Spec.PromotionTier = nil
	}
	if resp.DBInstance.PubliclyAccessible != nil {
		ko.Spec.PubliclyAccessible = resp.DBInstance.PubliclyAccessible
	} else {
		ko.Spec.PubliclyAccessible = nil
	}
	if resp.DBInstance.ReadReplicaDBClusterIdentifiers != nil {
		f48 := []*string{}
		for _, f48iter := range resp.DBInstance.ReadReplicaDBClusterIdentifiers {
			var f48elem string
			f48elem = *f48iter
			f48 = append(f48, &f48elem)
		}
		ko.Status.ReadReplicaDBClusterIdentifiers = f48
	} else {
		ko.Status.ReadReplicaDBClusterIdentifiers = nil
	}
	if resp.DBInstance.ReadReplicaDBInstanceIdentifiers != nil {
		f49 := []*string{}
		for _, f49iter := range resp.DBInstance.ReadReplicaDBInstanceIdentifiers {
			var f49elem string
			f49elem = *f49iter
			f49 = append(f49, &f49elem)
		}
		ko.Status.ReadReplicaDBInstanceIdentifiers = f49
	} else {
		ko.Status.ReadReplicaDBInstanceIdentifiers = nil
	}
	if resp.DBInstance.ReadReplicaSourceDBInstanceIdentifier != nil {
		ko.Status.ReadReplicaSourceDBInstanceIdentifier = resp.DBInstance.ReadReplicaSourceDBInstanceIdentifier
	} else {
		ko.Status.ReadReplicaSourceDBInstanceIdentifier = nil
	}
	if resp.DBInstance.SecondaryAvailabilityZone != nil {
		ko.Status.SecondaryAvailabilityZone = resp.DBInstance.SecondaryAvailabilityZone
	} else {
		ko.Status.SecondaryAvailabilityZone = nil
	}
	if resp.DBInstance.StatusInfos != nil {
		f52 := []*svcapitypes.DBInstanceStatusInfo{}
		for _, f52iter := range resp.DBInstance.StatusInfos {
			f52elem := &svcapitypes.DBInstanceStatusInfo{}
			if f52iter.Message != nil {
				f52elem.Message = f52iter.Message
			}
			if f52iter.Normal != nil {
				f52elem.Normal = f52iter.Normal
			}
			if f52iter.Status != nil {
				f52elem.Status = f52iter.Status
			}
			if f52iter.StatusType != nil {
				f52elem.StatusType = f52iter.StatusType
			}
			f52 = append(f52, f52elem)
		}
		ko.Status.StatusInfos = f52
	} else {
		ko.Status.StatusInfos = nil
	}
	if resp.DBInstance.StorageEncrypted != nil {
		ko.Spec.StorageEncrypted = resp.DBInstance.StorageEncrypted
	} else {
		ko.Spec.StorageEncrypted = nil
	}
	if resp.DBInstance.StorageType != nil {
		ko.Spec.StorageType = resp.DBInstance.StorageType
	} else {
		ko.Spec.StorageType = nil
	}
	if resp.DBInstance.TdeCredentialArn != nil {
		ko.Spec.TDECredentialARN = resp.DBInstance.TdeCredentialArn
	} else {
		ko.Spec.TDECredentialARN = nil
	}
	if resp.DBInstance.Timezone != nil {
		ko.Spec.Timezone = resp.DBInstance.Timezone
	} else {
		ko.Spec.Timezone = nil
	}
	if resp.DBInstance.VpcSecurityGroups != nil {
		f57 := []*svcapitypes.VPCSecurityGroupMembership{}
		for _, f57iter := range resp.DBInstance.VpcSecurityGroups {
			f57elem := &svcapitypes.VPCSecurityGroupMembership{}
			if f57iter.Status != nil {
				f57elem.Status = f57iter.Status
			}
			if f57iter.VpcSecurityGroupId != nil {
				f57elem.VPCSecurityGroupID = f57iter.VpcSecurityGroupId
			}
			f57 = append(f57, f57elem)
		}
		ko.Status.VPCSecurityGroups = f57
	} else {
		ko.Status.VPCSecurityGroups = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}

func TestSetResource_RDS_DBInstance_ReadMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "rds")

	crd := testutil.GetCRDByName(t, g, "DBInstance")
	require.NotNil(crd)

	// This asserts that the fields of the Spec and Status structs of the
	// target variable are constructed with cleaned, renamed-friendly names
	// referring to the generated Kubernetes API type definitions
	expected := `
	found := false
	for _, elem := range resp.DBInstances {
		if elem.AllocatedStorage != nil {
			ko.Spec.AllocatedStorage = elem.AllocatedStorage
		} else {
			ko.Spec.AllocatedStorage = nil
		}
		if elem.AssociatedRoles != nil {
			f1 := []*svcapitypes.DBInstanceRole{}
			for _, f1iter := range elem.AssociatedRoles {
				f1elem := &svcapitypes.DBInstanceRole{}
				if f1iter.FeatureName != nil {
					f1elem.FeatureName = f1iter.FeatureName
				}
				if f1iter.RoleArn != nil {
					f1elem.RoleARN = f1iter.RoleArn
				}
				if f1iter.Status != nil {
					f1elem.Status = f1iter.Status
				}
				f1 = append(f1, f1elem)
			}
			ko.Status.AssociatedRoles = f1
		} else {
			ko.Status.AssociatedRoles = nil
		}
		if elem.AutoMinorVersionUpgrade != nil {
			ko.Spec.AutoMinorVersionUpgrade = elem.AutoMinorVersionUpgrade
		} else {
			ko.Spec.AutoMinorVersionUpgrade = nil
		}
		if elem.AvailabilityZone != nil {
			ko.Spec.AvailabilityZone = elem.AvailabilityZone
		} else {
			ko.Spec.AvailabilityZone = nil
		}
		if elem.BackupRetentionPeriod != nil {
			ko.Spec.BackupRetentionPeriod = elem.BackupRetentionPeriod
		} else {
			ko.Spec.BackupRetentionPeriod = nil
		}
		if elem.CACertificateIdentifier != nil {
			ko.Status.CACertificateIdentifier = elem.CACertificateIdentifier
		} else {
			ko.Status.CACertificateIdentifier = nil
		}
		if elem.CharacterSetName != nil {
			ko.Spec.CharacterSetName = elem.CharacterSetName
		} else {
			ko.Spec.CharacterSetName = nil
		}
		if elem.CopyTagsToSnapshot != nil {
			ko.Spec.CopyTagsToSnapshot = elem.CopyTagsToSnapshot
		} else {
			ko.Spec.CopyTagsToSnapshot = nil
		}
		if elem.DBClusterIdentifier != nil {
			ko.Spec.DBClusterIdentifier = elem.DBClusterIdentifier
		} else {
			ko.Spec.DBClusterIdentifier = nil
		}
		if elem.DBInstanceArn != nil {
			if ko.Status.ACKResourceMetadata == nil {
				ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			tmpARN := ackv1alpha1.AWSResourceName(*elem.DBInstanceArn)
			ko.Status.ACKResourceMetadata.ARN = &tmpARN
		}
		if elem.DBInstanceClass != nil {
			ko.Spec.DBInstanceClass = elem.DBInstanceClass
		} else {
			ko.Spec.DBInstanceClass = nil
		}
		if elem.DBInstanceIdentifier != nil {
			ko.Spec.DBInstanceIdentifier = elem.DBInstanceIdentifier
		} else {
			ko.Spec.DBInstanceIdentifier = nil
		}
		if elem.DBInstanceStatus != nil {
			ko.Status.DBInstanceStatus = elem.DBInstanceStatus
		} else {
			ko.Status.DBInstanceStatus = nil
		}
		if elem.DBName != nil {
			ko.Spec.DBName = elem.DBName
		} else {
			ko.Spec.DBName = nil
		}
		if elem.DBParameterGroups != nil {
			f14 := []*svcapitypes.DBParameterGroupStatus_SDK{}
			for _, f14iter := range elem.DBParameterGroups {
				f14elem := &svcapitypes.DBParameterGroupStatus_SDK{}
				if f14iter.DBParameterGroupName != nil {
					f14elem.DBParameterGroupName = f14iter.DBParameterGroupName
				}
				if f14iter.ParameterApplyStatus != nil {
					f14elem.ParameterApplyStatus = f14iter.ParameterApplyStatus
				}
				f14 = append(f14, f14elem)
			}
			ko.Status.DBParameterGroups = f14
		} else {
			ko.Status.DBParameterGroups = nil
		}
		if elem.DBSecurityGroups != nil {
			f15 := []*string{}
			for _, f15iter := range elem.DBSecurityGroups {
				var f15elem string
				f15elem = *f15iter.DBSecurityGroupName
				f15 = append(f15, &f15elem)
			}
			ko.Spec.DBSecurityGroups = f15
		} else {
			ko.Spec.DBSecurityGroups = nil
		}
		if elem.DBSubnetGroup != nil {
			f16 := &svcapitypes.DBSubnetGroup_SDK{}
			if elem.DBSubnetGroup.DBSubnetGroupArn != nil {
				f16.DBSubnetGroupARN = elem.DBSubnetGroup.DBSubnetGroupArn
			}
			if elem.DBSubnetGroup.DBSubnetGroupDescription != nil {
				f16.DBSubnetGroupDescription = elem.DBSubnetGroup.DBSubnetGroupDescription
			}
			if elem.DBSubnetGroup.DBSubnetGroupName != nil {
				f16.DBSubnetGroupName = elem.DBSubnetGroup.DBSubnetGroupName
			}
			if elem.DBSubnetGroup.SubnetGroupStatus != nil {
				f16.SubnetGroupStatus = elem.DBSubnetGroup.SubnetGroupStatus
			}
			if elem.DBSubnetGroup.Subnets != nil {
				f16f4 := []*svcapitypes.Subnet{}
				for _, f16f4iter := range elem.DBSubnetGroup.Subnets {
					f16f4elem := &svcapitypes.Subnet{}
					if f16f4iter.SubnetAvailabilityZone != nil {
						f16f4elemf0 := &svcapitypes.AvailabilityZone{}
						if f16f4iter.SubnetAvailabilityZone.Name != nil {
							f16f4elemf0.Name = f16f4iter.SubnetAvailabilityZone.Name
						}
						f16f4elem.SubnetAvailabilityZone = f16f4elemf0
					}
					if f16f4iter.SubnetIdentifier != nil {
						f16f4elem.SubnetIdentifier = f16f4iter.SubnetIdentifier
					}
					if f16f4iter.SubnetOutpost != nil {
						f16f4elemf2 := &svcapitypes.Outpost{}
						if f16f4iter.SubnetOutpost.Arn != nil {
							f16f4elemf2.ARN = f16f4iter.SubnetOutpost.Arn
						}
						f16f4elem.SubnetOutpost = f16f4elemf2
					}
					if f16f4iter.SubnetStatus != nil {
						f16f4elem.SubnetStatus = f16f4iter.SubnetStatus
					}
					f16f4 = append(f16f4, f16f4elem)
				}
				f16.Subnets = f16f4
			}
			if elem.DBSubnetGroup.VpcId != nil {
				f16.VPCID = elem.DBSubnetGroup.VpcId
			}
			ko.Status.DBSubnetGroup = f16
		} else {
			ko.Status.DBSubnetGroup = nil
		}
		if elem.DbInstancePort != nil {
			ko.Status.DBInstancePort = elem.DbInstancePort
		} else {
			ko.Status.DBInstancePort = nil
		}
		if elem.DbiResourceId != nil {
			ko.Status.DBIResourceID = elem.DbiResourceId
		} else {
			ko.Status.DBIResourceID = nil
		}
		if elem.DeletionProtection != nil {
			ko.Spec.DeletionProtection = elem.DeletionProtection
		} else {
			ko.Spec.DeletionProtection = nil
		}
		if elem.DomainMemberships != nil {
			f20 := []*svcapitypes.DomainMembership{}
			for _, f20iter := range elem.DomainMemberships {
				f20elem := &svcapitypes.DomainMembership{}
				if f20iter.Domain != nil {
					f20elem.Domain = f20iter.Domain
				}
				if f20iter.FQDN != nil {
					f20elem.FQDN = f20iter.FQDN
				}
				if f20iter.IAMRoleName != nil {
					f20elem.IAMRoleName = f20iter.IAMRoleName
				}
				if f20iter.Status != nil {
					f20elem.Status = f20iter.Status
				}
				f20 = append(f20, f20elem)
			}
			ko.Status.DomainMemberships = f20
		} else {
			ko.Status.DomainMemberships = nil
		}
		if elem.EnabledCloudwatchLogsExports != nil {
			f21 := []*string{}
			for _, f21iter := range elem.EnabledCloudwatchLogsExports {
				var f21elem string
				f21elem = *f21iter
				f21 = append(f21, &f21elem)
			}
			ko.Status.EnabledCloudwatchLogsExports = f21
		} else {
			ko.Status.EnabledCloudwatchLogsExports = nil
		}
		if elem.Endpoint != nil {
			f22 := &svcapitypes.Endpoint{}
			if elem.Endpoint.Address != nil {
				f22.Address = elem.Endpoint.Address
			}
			if elem.Endpoint.HostedZoneId != nil {
				f22.HostedZoneID = elem.Endpoint.HostedZoneId
			}
			if elem.Endpoint.Port != nil {
				f22.Port = elem.Endpoint.Port
			}
			ko.Status.Endpoint = f22
		} else {
			ko.Status.Endpoint = nil
		}
		if elem.Engine != nil {
			ko.Spec.Engine = elem.Engine
		} else {
			ko.Spec.Engine = nil
		}
		if elem.EngineVersion != nil {
			ko.Spec.EngineVersion = elem.EngineVersion
		} else {
			ko.Spec.EngineVersion = nil
		}
		if elem.EnhancedMonitoringResourceArn != nil {
			ko.Status.EnhancedMonitoringResourceARN = elem.EnhancedMonitoringResourceArn
		} else {
			ko.Status.EnhancedMonitoringResourceARN = nil
		}
		if elem.IAMDatabaseAuthenticationEnabled != nil {
			ko.Status.IAMDatabaseAuthenticationEnabled = elem.IAMDatabaseAuthenticationEnabled
		} else {
			ko.Status.IAMDatabaseAuthenticationEnabled = nil
		}
		if elem.InstanceCreateTime != nil {
			ko.Status.InstanceCreateTime = &metav1.Time{*elem.InstanceCreateTime}
		} else {
			ko.Status.InstanceCreateTime = nil
		}
		if elem.Iops != nil {
			ko.Spec.IOPS = elem.Iops
		} else {
			ko.Spec.IOPS = nil
		}
		if elem.KmsKeyId != nil {
			ko.Spec.KMSKeyID = elem.KmsKeyId
		} else {
			ko.Spec.KMSKeyID = nil
		}
		if elem.LatestRestorableTime != nil {
			ko.Status.LatestRestorableTime = &metav1.Time{*elem.LatestRestorableTime}
		} else {
			ko.Status.LatestRestorableTime = nil
		}
		if elem.LicenseModel != nil {
			ko.Spec.LicenseModel = elem.LicenseModel
		} else {
			ko.Spec.LicenseModel = nil
		}
		if elem.ListenerEndpoint != nil {
			f32 := &svcapitypes.Endpoint{}
			if elem.ListenerEndpoint.Address != nil {
				f32.Address = elem.ListenerEndpoint.Address
			}
			if elem.ListenerEndpoint.HostedZoneId != nil {
				f32.HostedZoneID = elem.ListenerEndpoint.HostedZoneId
			}
			if elem.ListenerEndpoint.Port != nil {
				f32.Port = elem.ListenerEndpoint.Port
			}
			ko.Status.ListenerEndpoint = f32
		} else {
			ko.Status.ListenerEndpoint = nil
		}
		if elem.MasterUsername != nil {
			ko.Spec.MasterUsername = elem.MasterUsername
		} else {
			ko.Spec.MasterUsername = nil
		}
		if elem.MaxAllocatedStorage != nil {
			ko.Spec.MaxAllocatedStorage = elem.MaxAllocatedStorage
		} else {
			ko.Spec.MaxAllocatedStorage = nil
		}
		if elem.MonitoringInterval != nil {
			ko.Spec.MonitoringInterval = elem.MonitoringInterval
		} else {
			ko.Spec.MonitoringInterval = nil
		}
		if elem.MonitoringRoleArn != nil {
			ko.Spec.MonitoringRoleARN = elem.MonitoringRoleArn
		} else {
			ko.Spec.MonitoringRoleARN = nil
		}
		if elem.MultiAZ != nil {
			ko.Spec.MultiAZ = elem.MultiAZ
		} else {
			ko.Spec.MultiAZ = nil
		}
		if elem.OptionGroupMemberships != nil {
			f38 := []*svcapitypes.OptionGroupMembership{}
			for _, f38iter := range elem.OptionGroupMemberships {
				f38elem := &svcapitypes.OptionGroupMembership{}
				if f38iter.OptionGroupName != nil {
					f38elem.OptionGroupName = f38iter.OptionGroupName
				}
				if f38iter.Status != nil {
					f38elem.Status = f38iter.Status
				}
				f38 = append(f38, f38elem)
			}
			ko.Status.OptionGroupMemberships = f38
		} else {
			ko.Status.OptionGroupMemberships = nil
		}
		if elem.PendingModifiedValues != nil {
			f39 := &svcapitypes.PendingModifiedValues{}
			if elem.PendingModifiedValues.AllocatedStorage != nil {
				f39.AllocatedStorage = elem.PendingModifiedValues.AllocatedStorage
			}
			if elem.PendingModifiedValues.BackupRetentionPeriod != nil {
				f39.BackupRetentionPeriod = elem.PendingModifiedValues.BackupRetentionPeriod
			}
			if elem.PendingModifiedValues.CACertificateIdentifier != nil {
				f39.CACertificateIdentifier = elem.PendingModifiedValues.CACertificateIdentifier
			}
			if elem.PendingModifiedValues.DBInstanceClass != nil {
				f39.DBInstanceClass = elem.PendingModifiedValues.DBInstanceClass
			}
			if elem.PendingModifiedValues.DBInstanceIdentifier != nil {
				f39.DBInstanceIdentifier = elem.PendingModifiedValues.DBInstanceIdentifier
			}
			if elem.PendingModifiedValues.DBSubnetGroupName != nil {
				f39.DBSubnetGroupName = elem.PendingModifiedValues.DBSubnetGroupName
			}
			if elem.PendingModifiedValues.EngineVersion != nil {
				f39.EngineVersion = elem.PendingModifiedValues.EngineVersion
			}
			if elem.PendingModifiedValues.Iops != nil {
				f39.IOPS = elem.PendingModifiedValues.Iops
			}
			if elem.PendingModifiedValues.LicenseModel != nil {
				f39.LicenseModel = elem.PendingModifiedValues.LicenseModel
			}
			if elem.PendingModifiedValues.MasterUserPassword != nil {
				f39.MasterUserPassword = elem.PendingModifiedValues.MasterUserPassword
			}
			if elem.PendingModifiedValues.MultiAZ != nil {
				f39.MultiAZ = elem.PendingModifiedValues.MultiAZ
			}
			if elem.PendingModifiedValues.PendingCloudwatchLogsExports != nil {
				f39f11 := &svcapitypes.PendingCloudwatchLogsExports{}
				if elem.PendingModifiedValues.PendingCloudwatchLogsExports.LogTypesToDisable != nil {
					f39f11f0 := []*string{}
					for _, f39f11f0iter := range elem.PendingModifiedValues.PendingCloudwatchLogsExports.LogTypesToDisable {
						var f39f11f0elem string
						f39f11f0elem = *f39f11f0iter
						f39f11f0 = append(f39f11f0, &f39f11f0elem)
					}
					f39f11.LogTypesToDisable = f39f11f0
				}
				if elem.PendingModifiedValues.PendingCloudwatchLogsExports.LogTypesToEnable != nil {
					f39f11f1 := []*string{}
					for _, f39f11f1iter := range elem.PendingModifiedValues.PendingCloudwatchLogsExports.LogTypesToEnable {
						var f39f11f1elem string
						f39f11f1elem = *f39f11f1iter
						f39f11f1 = append(f39f11f1, &f39f11f1elem)
					}
					f39f11.LogTypesToEnable = f39f11f1
				}
				f39.PendingCloudwatchLogsExports = f39f11
			}
			if elem.PendingModifiedValues.Port != nil {
				f39.Port = elem.PendingModifiedValues.Port
			}
			if elem.PendingModifiedValues.ProcessorFeatures != nil {
				f39f13 := []*svcapitypes.ProcessorFeature{}
				for _, f39f13iter := range elem.PendingModifiedValues.ProcessorFeatures {
					f39f13elem := &svcapitypes.ProcessorFeature{}
					if f39f13iter.Name != nil {
						f39f13elem.Name = f39f13iter.Name
					}
					if f39f13iter.Value != nil {
						f39f13elem.Value = f39f13iter.Value
					}
					f39f13 = append(f39f13, f39f13elem)
				}
				f39.ProcessorFeatures = f39f13
			}
			if elem.PendingModifiedValues.StorageType != nil {
				f39.StorageType = elem.PendingModifiedValues.StorageType
			}
			ko.Status.PendingModifiedValues = f39
		} else {
			ko.Status.PendingModifiedValues = nil
		}
		if elem.PerformanceInsightsEnabled != nil {
			ko.Status.PerformanceInsightsEnabled = elem.PerformanceInsightsEnabled
		} else {
			ko.Status.PerformanceInsightsEnabled = nil
		}
		if elem.PerformanceInsightsKMSKeyId != nil {
			ko.Spec.PerformanceInsightsKMSKeyID = elem.PerformanceInsightsKMSKeyId
		} else {
			ko.Spec.PerformanceInsightsKMSKeyID = nil
		}
		if elem.PerformanceInsightsRetentionPeriod != nil {
			ko.Spec.PerformanceInsightsRetentionPeriod = elem.PerformanceInsightsRetentionPeriod
		} else {
			ko.Spec.PerformanceInsightsRetentionPeriod = nil
		}
		if elem.PreferredBackupWindow != nil {
			ko.Spec.PreferredBackupWindow = elem.PreferredBackupWindow
		} else {
			ko.Spec.PreferredBackupWindow = nil
		}
		if elem.PreferredMaintenanceWindow != nil {
			ko.Spec.PreferredMaintenanceWindow = elem.PreferredMaintenanceWindow
		} else {
			ko.Spec.PreferredMaintenanceWindow = nil
		}
		if elem.ProcessorFeatures != nil {
			f45 := []*svcapitypes.ProcessorFeature{}
			for _, f45iter := range elem.ProcessorFeatures {
				f45elem := &svcapitypes.ProcessorFeature{}
				if f45iter.Name != nil {
					f45elem.Name = f45iter.Name
				}
				if f45iter.Value != nil {
					f45elem.Value = f45iter.Value
				}
				f45 = append(f45, f45elem)
			}
			ko.Spec.ProcessorFeatures = f45
		} else {
			ko.Spec.ProcessorFeatures = nil
		}
		if elem.PromotionTier != nil {
			ko.Spec.PromotionTier = elem.PromotionTier
		} else {
			ko.Spec.PromotionTier = nil
		}
		if elem.PubliclyAccessible != nil {
			ko.Spec.PubliclyAccessible = elem.PubliclyAccessible
		} else {
			ko.Spec.PubliclyAccessible = nil
		}
		if elem.ReadReplicaDBClusterIdentifiers != nil {
			f48 := []*string{}
			for _, f48iter := range elem.ReadReplicaDBClusterIdentifiers {
				var f48elem string
				f48elem = *f48iter
				f48 = append(f48, &f48elem)
			}
			ko.Status.ReadReplicaDBClusterIdentifiers = f48
		} else {
			ko.Status.ReadReplicaDBClusterIdentifiers = nil
		}
		if elem.ReadReplicaDBInstanceIdentifiers != nil {
			f49 := []*string{}
			for _, f49iter := range elem.ReadReplicaDBInstanceIdentifiers {
				var f49elem string
				f49elem = *f49iter
				f49 = append(f49, &f49elem)
			}
			ko.Status.ReadReplicaDBInstanceIdentifiers = f49
		} else {
			ko.Status.ReadReplicaDBInstanceIdentifiers = nil
		}
		if elem.ReadReplicaSourceDBInstanceIdentifier != nil {
			ko.Status.ReadReplicaSourceDBInstanceIdentifier = elem.ReadReplicaSourceDBInstanceIdentifier
		} else {
			ko.Status.ReadReplicaSourceDBInstanceIdentifier = nil
		}
		if elem.SecondaryAvailabilityZone != nil {
			ko.Status.SecondaryAvailabilityZone = elem.SecondaryAvailabilityZone
		} else {
			ko.Status.SecondaryAvailabilityZone = nil
		}
		if elem.StatusInfos != nil {
			f52 := []*svcapitypes.DBInstanceStatusInfo{}
			for _, f52iter := range elem.StatusInfos {
				f52elem := &svcapitypes.DBInstanceStatusInfo{}
				if f52iter.Message != nil {
					f52elem.Message = f52iter.Message
				}
				if f52iter.Normal != nil {
					f52elem.Normal = f52iter.Normal
				}
				if f52iter.Status != nil {
					f52elem.Status = f52iter.Status
				}
				if f52iter.StatusType != nil {
					f52elem.StatusType = f52iter.StatusType
				}
				f52 = append(f52, f52elem)
			}
			ko.Status.StatusInfos = f52
		} else {
			ko.Status.StatusInfos = nil
		}
		if elem.StorageEncrypted != nil {
			ko.Spec.StorageEncrypted = elem.StorageEncrypted
		} else {
			ko.Spec.StorageEncrypted = nil
		}
		if elem.StorageType != nil {
			ko.Spec.StorageType = elem.StorageType
		} else {
			ko.Spec.StorageType = nil
		}
		if elem.TdeCredentialArn != nil {
			ko.Spec.TDECredentialARN = elem.TdeCredentialArn
		} else {
			ko.Spec.TDECredentialARN = nil
		}
		if elem.Timezone != nil {
			ko.Spec.Timezone = elem.Timezone
		} else {
			ko.Spec.Timezone = nil
		}
		if elem.VpcSecurityGroups != nil {
			f57 := []*svcapitypes.VPCSecurityGroupMembership{}
			for _, f57iter := range elem.VpcSecurityGroups {
				f57elem := &svcapitypes.VPCSecurityGroupMembership{}
				if f57iter.Status != nil {
					f57elem.Status = f57iter.Status
				}
				if f57iter.VpcSecurityGroupId != nil {
					f57elem.VPCSecurityGroupID = f57iter.VpcSecurityGroupId
				}
				f57 = append(f57, f57elem)
			}
			ko.Status.VPCSecurityGroups = f57
		} else {
			ko.Status.VPCSecurityGroups = nil
		}
		found = true
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeList, "resp", "ko", 1),
	)
}

func TestSetResource_SNS_Topic_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sns")

	crd := testutil.GetCRDByName(t, g, "Topic")
	require.NotNil(crd)

	// None of the fields in the Topic resource's CreateTopicInput shape are
	// returned in the CreateTopicOutput shape, so none of them return any Go
	// code for setting a Status struct field to a corresponding Create Output
	// Shape member. However, the returned output shape DOES include the
	// Topic's ARN field (TopicArn), which we should be storing in the
	// ACKResourceMetadata.ARN standardized field
	expected := `
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.TopicArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.TopicArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}

func TestSetResource_SNS_Topic_GetAttributes(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sns")

	crd := testutil.GetCRDByName(t, g, "Topic")
	require.NotNil(crd)

	// The output shape for the GetAttributes operation contains a single field
	// "Attributes" that must be unpacked into the Topic CRD's Status fields.
	// There are only three attribute keys that are *not* in the Input shape
	// (and thus in the Spec fields). Two of them are the tesource's ARN and
	// AWS Owner account ID, both of which are handled specially.
	expected := `
	ko.Status.EffectiveDeliveryPolicy = resp.Attributes["EffectiveDeliveryPolicy"]
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	tmpOwnerID := ackv1alpha1.AWSAccountID(*resp.Attributes["Owner"])
	ko.Status.ACKResourceMetadata.OwnerAccountID = &tmpOwnerID
	tmpARN := ackv1alpha1.AWSResourceName(*resp.Attributes["TopicArn"])
	ko.Status.ACKResourceMetadata.ARN = &tmpARN
`
	assert.Equal(
		expected,
		code.SetResourceGetAttributes(crd.Config(), crd, "resp", "ko", 1),
	)
}

func TestSetResource_SQS_Queue_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sqs")

	crd := testutil.GetCRDByName(t, g, "Queue")
	require.NotNil(crd)

	// There are no fields other than QueueID in the returned CreateQueueResult
	// shape
	expected := `
	if resp.QueueUrl != nil {
		ko.Status.QueueURL = resp.QueueUrl
	} else {
		ko.Status.QueueURL = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}

func TestSetResource_SQS_Queue_GetAttributes(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sqs")

	crd := testutil.GetCRDByName(t, g, "Queue")
	require.NotNil(crd)

	// The output shape for the GetAttributes operation contains a single field
	// "Attributes" that must be unpacked into the Queue CRD's Status fields.
	// There are only three attribute keys that are *not* in the Input shape
	// (and thus in the Spec fields). One of them is the resource's ARN which
	// is handled specially.
	expected := `
	ko.Status.CreatedTimestamp = resp.Attributes["CreatedTimestamp"]
	ko.Status.LastModifiedTimestamp = resp.Attributes["LastModifiedTimestamp"]
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	tmpARN := ackv1alpha1.AWSResourceName(*resp.Attributes["QueueArn"])
	ko.Status.ACKResourceMetadata.ARN = &tmpARN
`
	assert.Equal(
		expected,
		code.SetResourceGetAttributes(crd.Config(), crd, "resp", "ko", 1),
	)
}

func TestSetResource_RDS_DBSubnetGroup_ReadMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "rds")

	crd := testutil.GetCRDByName(t, g, "DBSubnetGroup")
	require.NotNil(crd)

	// There are renamed fields for Name and Description in order to
	// "de-stutter" the field names. We want to verify that the SetResource for
	// the DescribeDBSubnetGroups API operation sets these fields in the Spec
	// properly
	expected := `
	found := false
	for _, elem := range resp.DBSubnetGroups {
		if elem.DBSubnetGroupArn != nil {
			if ko.Status.ACKResourceMetadata == nil {
				ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			tmpARN := ackv1alpha1.AWSResourceName(*elem.DBSubnetGroupArn)
			ko.Status.ACKResourceMetadata.ARN = &tmpARN
		}
		if elem.DBSubnetGroupDescription != nil {
			ko.Spec.Description = elem.DBSubnetGroupDescription
		} else {
			ko.Spec.Description = nil
		}
		if elem.DBSubnetGroupName != nil {
			ko.Spec.Name = elem.DBSubnetGroupName
		} else {
			ko.Spec.Name = nil
		}
		if elem.SubnetGroupStatus != nil {
			ko.Status.SubnetGroupStatus = elem.SubnetGroupStatus
		} else {
			ko.Status.SubnetGroupStatus = nil
		}
		if elem.Subnets != nil {
			f4 := []*svcapitypes.Subnet{}
			for _, f4iter := range elem.Subnets {
				f4elem := &svcapitypes.Subnet{}
				if f4iter.SubnetAvailabilityZone != nil {
					f4elemf0 := &svcapitypes.AvailabilityZone{}
					if f4iter.SubnetAvailabilityZone.Name != nil {
						f4elemf0.Name = f4iter.SubnetAvailabilityZone.Name
					}
					f4elem.SubnetAvailabilityZone = f4elemf0
				}
				if f4iter.SubnetIdentifier != nil {
					f4elem.SubnetIdentifier = f4iter.SubnetIdentifier
				}
				if f4iter.SubnetOutpost != nil {
					f4elemf2 := &svcapitypes.Outpost{}
					if f4iter.SubnetOutpost.Arn != nil {
						f4elemf2.ARN = f4iter.SubnetOutpost.Arn
					}
					f4elem.SubnetOutpost = f4elemf2
				}
				if f4iter.SubnetStatus != nil {
					f4elem.SubnetStatus = f4iter.SubnetStatus
				}
				f4 = append(f4, f4elem)
			}
			ko.Status.Subnets = f4
		} else {
			ko.Status.Subnets = nil
		}
		if elem.VpcId != nil {
			ko.Status.VPCID = elem.VpcId
		} else {
			ko.Status.VPCID = nil
		}
		found = true
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeList, "resp", "ko", 1),
	)
}

func TestGetOutputShape_VPC_No_Override(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "Vpc")
	require.NotNil(crd)

	expectedShape := crd.Ops.ReadMany.OutputRef.Shape
	outputShape, _ := crd.GetOutputShape(crd.Ops.ReadMany)
	assert.Equal(
		expectedShape,
		outputShape)
}

func TestGetOutputShape_DynamoDB_Override(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "dynamodb")

	crd := testutil.GetCRDByName(t, g, "Backup")
	require.NotNil(crd)

	outputShape, _ := crd.GetOutputShape(crd.Ops.ReadOne)
	assert.Equal(
		"BackupDetails",
		outputShape.ShapeName)
}

func TestGetOutputShape_VPCEndpoint_Override(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "VpcEndpoint")
	require.NotNil(crd)

	outputShape, _ := crd.GetOutputShape(crd.Ops.Create)
	assert.Equal(
		"VpcEndpoint",
		outputShape.ShapeName)
}

func TestSetResource_MQ_Broker_SetResourceIdentifiers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "mq")

	crd := testutil.GetCRDByName(t, g, "Broker")
	require.NotNil(crd)

	expected := `
	if identifier.NameOrID == "" {
		return ackerrors.MissingNameIdentifier
	}
	r.ko.Status.BrokerID = &identifier.NameOrID

`
	assert.Equal(
		expected,
		code.SetResourceIdentifiers(crd.Config(), crd, "identifier", "r.ko", 1),
	)
}

func TestSetResource_RDS_DBInstances_SetResourceIdentifiers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "rds")

	crd := testutil.GetCRDByName(t, g, "DBInstance")
	require.NotNil(crd)

	expected := `
	if identifier.NameOrID == "" {
		return ackerrors.MissingNameIdentifier
	}
	r.ko.Spec.DBInstanceIdentifier = &identifier.NameOrID

`
	assert.Equal(
		expected,
		code.SetResourceIdentifiers(crd.Config(), crd, "identifier", "r.ko", 1),
	)
}

func TestSetResource_RDS_DBSubnetGroup_SetResourceIdentifiers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "rds")

	crd := testutil.GetCRDByName(t, g, "DBSubnetGroup")
	require.NotNil(crd)

	// In our testdata generator.yaml file, we've renamed the original
	// `DBSubnetGroupName` to just `Name`
	expected := `
	if identifier.NameOrID == "" {
		return ackerrors.MissingNameIdentifier
	}
	r.ko.Spec.Name = &identifier.NameOrID

`
	assert.Equal(
		expected,
		code.SetResourceIdentifiers(crd.Config(), crd, "identifier", "r.ko", 1),
	)
}

func TestSetResource_APIGWV2_ApiMapping_SetResourceIdentifiers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "ApiMapping")
	require.NotNil(crd)

	expected := `
	if identifier.NameOrID == "" {
		return ackerrors.MissingNameIdentifier
	}
	r.ko.Status.APIMappingID = &identifier.NameOrID

	f1, f1ok := identifier.AdditionalKeys["domainName"]
	if f1ok {
		r.ko.Spec.DomainName = &f1
	}
`
	assert.Equal(
		expected,
		code.SetResourceIdentifiers(crd.Config(), crd, "identifier", "r.ko", 1),
	)
}

func TestSetResource_SageMaker_ModelPackage_SetResourceIdentifiers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sagemaker")

	crd := testutil.GetCRDByName(t, g, "ModelPackage")
	require.NotNil(crd)

	expected := `
	if r.ko.Status.ACKResourceMetadata == nil {
		r.ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	r.ko.Status.ACKResourceMetadata.ARN = identifier.ARN
`
	assert.Equal(
		expected,
		code.SetResourceIdentifiers(crd.Config(), crd, "identifier", "r.ko", 1),
	)
}

func TestSetResource_EC2_VPC_SetResourceIdentifiers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "Vpc")
	require.NotNil(crd)

	expected := `
	if identifier.NameOrID == "" {
		return ackerrors.MissingNameIdentifier
	}
	r.ko.Status.VPCID = &identifier.NameOrID

`
	assert.Equal(
		expected,
		code.SetResourceIdentifiers(crd.Config(), crd, "identifier", "r.ko", 1),
	)
}

func TestSetResource_EC2_SecurityGroups_SetResourceIdentifiers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "SecurityGroup")
	require.NotNil(crd)

	// CreateSecurityGroup Output returns a GroupId,
	// which has been renamed to ID in the CRD
	expected := `
	if resp.GroupId != nil {
		ko.Status.ID = resp.GroupId
	} else {
		ko.Status.ID = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}

func TestSetResource_IAM_Role_NestedSetConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "iam")

	crd := testutil.GetCRDByName(t, g, "Role")
	require.NotNil(crd)

	// The input and output shapes for the PermissionsBoundary are different
	// and we have a custom SetConfig for this field in our generator.yaml file
	// to configure the SetResource function to set the value of the resource's
	// Spec.PermissionsBoundary to the value of the (nested)
	// GetRoleResponse.Role.PermissionsBoundary.PermissionsBoundaryArn field
	expected := `
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Role.Arn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Role.Arn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.Role.AssumeRolePolicyDocument != nil {
		ko.Spec.AssumeRolePolicyDocument = resp.Role.AssumeRolePolicyDocument
	} else {
		ko.Spec.AssumeRolePolicyDocument = nil
	}
	if resp.Role.CreateDate != nil {
		ko.Status.CreateDate = &metav1.Time{*resp.Role.CreateDate}
	} else {
		ko.Status.CreateDate = nil
	}
	if resp.Role.Description != nil {
		ko.Spec.Description = resp.Role.Description
	} else {
		ko.Spec.Description = nil
	}
	if resp.Role.MaxSessionDuration != nil {
		ko.Spec.MaxSessionDuration = resp.Role.MaxSessionDuration
	} else {
		ko.Spec.MaxSessionDuration = nil
	}
	if resp.Role.Path != nil {
		ko.Spec.Path = resp.Role.Path
	} else {
		ko.Spec.Path = nil
	}
	if resp.Role.PermissionsBoundary != nil {
		ko.Spec.PermissionsBoundary = resp.Role.PermissionsBoundary.PermissionsBoundaryArn
	} else {
		ko.Spec.PermissionsBoundary = nil
	}
	if resp.Role.RoleId != nil {
		ko.Status.RoleID = resp.Role.RoleId
	} else {
		ko.Status.RoleID = nil
	}
	if resp.Role.RoleLastUsed != nil {
		f8 := &svcapitypes.RoleLastUsed{}
		if resp.Role.RoleLastUsed.LastUsedDate != nil {
			f8.LastUsedDate = &metav1.Time{*resp.Role.RoleLastUsed.LastUsedDate}
		}
		if resp.Role.RoleLastUsed.Region != nil {
			f8.Region = resp.Role.RoleLastUsed.Region
		}
		ko.Status.RoleLastUsed = f8
	} else {
		ko.Status.RoleLastUsed = nil
	}
	if resp.Role.RoleName != nil {
		ko.Spec.Name = resp.Role.RoleName
	} else {
		ko.Spec.Name = nil
	}
	if resp.Role.Tags != nil {
		f10 := []*svcapitypes.Tag{}
		for _, f10iter := range resp.Role.Tags {
			f10elem := &svcapitypes.Tag{}
			if f10iter.Key != nil {
				f10elem.Key = f10iter.Key
			}
			if f10iter.Value != nil {
				f10elem.Value = f10iter.Value
			}
			f10 = append(f10, f10elem)
		}
		ko.Spec.Tags = f10
	} else {
		ko.Spec.Tags = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeGet, "resp", "ko", 1),
	)
}

func TestSetResource_EC2_DHCPOptions_NestedSetConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "DhcpOptions")
	require.NotNil(crd)

	// The input and output shapes for the DhcpConfigurations.Values are different
	// and we have a custom SetConfig for this field in our generator.yaml file
	// to configure the SetResource function to set the value of the resource's
	// (nested) Spec.DHCPConfigurations.Values to the value of the (nested)
	// GetDhcpOptionsResponse.DhcpOptions.DhcpConfigurations.Values.Value field
	expected := `
	if resp.DhcpOptions.DhcpConfigurations != nil {
		f0 := []*svcapitypes.NewDHCPConfiguration{}
		for _, f0iter := range resp.DhcpOptions.DhcpConfigurations {
			f0elem := &svcapitypes.NewDHCPConfiguration{}
			if f0iter.Key != nil {
				f0elem.Key = f0iter.Key
			}
			if f0iter.Values != nil {
				f0elemf1 := []*string{}
				for _, f0elemf1iter := range f0iter.Values {
					var f0elemf1elem string
					if f0elemf1iter.Value != nil {
						f0elemf1elem = *f0elemf1iter.Value
					}
					f0elemf1 = append(f0elemf1, &f0elemf1elem)
				}
				f0elem.Values = f0elemf1
			}
			f0 = append(f0, f0elem)
		}
		ko.Spec.DHCPConfigurations = f0
	} else {
		ko.Spec.DHCPConfigurations = nil
	}
	if resp.DhcpOptions.DhcpOptionsId != nil {
		ko.Status.DHCPOptionsID = resp.DhcpOptions.DhcpOptionsId
	} else {
		ko.Status.DHCPOptionsID = nil
	}
	if resp.DhcpOptions.OwnerId != nil {
		ko.Status.OwnerID = resp.DhcpOptions.OwnerId
	} else {
		ko.Status.OwnerID = nil
	}
	if resp.DhcpOptions.Tags != nil {
		f3 := []*svcapitypes.Tag{}
		for _, f3iter := range resp.DhcpOptions.Tags {
			f3elem := &svcapitypes.Tag{}
			if f3iter.Key != nil {
				f3elem.Key = f3iter.Key
			}
			if f3iter.Value != nil {
				f3elem.Value = f3iter.Value
			}
			f3 = append(f3, f3elem)
		}
		ko.Status.Tags = f3
	} else {
		ko.Status.Tags = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}
