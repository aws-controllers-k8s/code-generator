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

func TestSetSDK_APIGWv2_Route_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "Route")
	require.NotNil(crd)

	expected := `
	if r.ko.Spec.APIID != nil {
		res.SetApiId(*r.ko.Spec.APIID)
	}
	if r.ko.Spec.APIKeyRequired != nil {
		res.SetApiKeyRequired(*r.ko.Spec.APIKeyRequired)
	}
	if r.ko.Spec.AuthorizationScopes != nil {
		f2 := []*string{}
		for _, f2iter := range r.ko.Spec.AuthorizationScopes {
			var f2elem string
			f2elem = *f2iter
			f2 = append(f2, &f2elem)
		}
		res.SetAuthorizationScopes(f2)
	}
	if r.ko.Spec.AuthorizationType != nil {
		res.SetAuthorizationType(*r.ko.Spec.AuthorizationType)
	}
	if r.ko.Spec.AuthorizerID != nil {
		res.SetAuthorizerId(*r.ko.Spec.AuthorizerID)
	}
	if r.ko.Spec.ModelSelectionExpression != nil {
		res.SetModelSelectionExpression(*r.ko.Spec.ModelSelectionExpression)
	}
	if r.ko.Spec.OperationName != nil {
		res.SetOperationName(*r.ko.Spec.OperationName)
	}
	if r.ko.Spec.RequestModels != nil {
		f7 := map[string]*string{}
		for f7key, f7valiter := range r.ko.Spec.RequestModels {
			var f7val string
			f7val = *f7valiter
			f7[f7key] = &f7val
		}
		res.SetRequestModels(f7)
	}
	if r.ko.Spec.RequestParameters != nil {
		f8 := map[string]*svcsdk.ParameterConstraints{}
		for f8key, f8valiter := range r.ko.Spec.RequestParameters {
			f8val := &svcsdk.ParameterConstraints{}
			if f8valiter.Required != nil {
				f8val.SetRequired(*f8valiter.Required)
			}
			f8[f8key] = f8val
		}
		res.SetRequestParameters(f8)
	}
	if r.ko.Spec.RouteKey != nil {
		res.SetRouteKey(*r.ko.Spec.RouteKey)
	}
	if r.ko.Spec.RouteResponseSelectionExpression != nil {
		res.SetRouteResponseSelectionExpression(*r.ko.Spec.RouteResponseSelectionExpression)
	}
	if r.ko.Spec.Target != nil {
		res.SetTarget(*r.ko.Spec.Target)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_DynamoDB_Table_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "dynamodb")

	crd := testutil.GetCRDByName(t, g, "Table")
	require.NotNil(crd)

	expected := `
	if r.ko.Spec.AttributeDefinitions != nil {
		f0 := []*svcsdk.AttributeDefinition{}
		for _, f0iter := range r.ko.Spec.AttributeDefinitions {
			f0elem := &svcsdk.AttributeDefinition{}
			if f0iter.AttributeName != nil {
				f0elem.SetAttributeName(*f0iter.AttributeName)
			}
			if f0iter.AttributeType != nil {
				f0elem.SetAttributeType(*f0iter.AttributeType)
			}
			f0 = append(f0, f0elem)
		}
		res.SetAttributeDefinitions(f0)
	}
	if r.ko.Spec.BillingMode != nil {
		res.SetBillingMode(*r.ko.Spec.BillingMode)
	}
	if r.ko.Spec.GlobalSecondaryIndexes != nil {
		f2 := []*svcsdk.GlobalSecondaryIndex{}
		for _, f2iter := range r.ko.Spec.GlobalSecondaryIndexes {
			f2elem := &svcsdk.GlobalSecondaryIndex{}
			if f2iter.IndexName != nil {
				f2elem.SetIndexName(*f2iter.IndexName)
			}
			if f2iter.KeySchema != nil {
				f2elemf1 := []*svcsdk.KeySchemaElement{}
				for _, f2elemf1iter := range f2iter.KeySchema {
					f2elemf1elem := &svcsdk.KeySchemaElement{}
					if f2elemf1iter.AttributeName != nil {
						f2elemf1elem.SetAttributeName(*f2elemf1iter.AttributeName)
					}
					if f2elemf1iter.KeyType != nil {
						f2elemf1elem.SetKeyType(*f2elemf1iter.KeyType)
					}
					f2elemf1 = append(f2elemf1, f2elemf1elem)
				}
				f2elem.SetKeySchema(f2elemf1)
			}
			if f2iter.Projection != nil {
				f2elemf2 := &svcsdk.Projection{}
				if f2iter.Projection.NonKeyAttributes != nil {
					f2elemf2f0 := []*string{}
					for _, f2elemf2f0iter := range f2iter.Projection.NonKeyAttributes {
						var f2elemf2f0elem string
						f2elemf2f0elem = *f2elemf2f0iter
						f2elemf2f0 = append(f2elemf2f0, &f2elemf2f0elem)
					}
					f2elemf2.SetNonKeyAttributes(f2elemf2f0)
				}
				if f2iter.Projection.ProjectionType != nil {
					f2elemf2.SetProjectionType(*f2iter.Projection.ProjectionType)
				}
				f2elem.SetProjection(f2elemf2)
			}
			if f2iter.ProvisionedThroughput != nil {
				f2elemf3 := &svcsdk.ProvisionedThroughput{}
				if f2iter.ProvisionedThroughput.ReadCapacityUnits != nil {
					f2elemf3.SetReadCapacityUnits(*f2iter.ProvisionedThroughput.ReadCapacityUnits)
				}
				if f2iter.ProvisionedThroughput.WriteCapacityUnits != nil {
					f2elemf3.SetWriteCapacityUnits(*f2iter.ProvisionedThroughput.WriteCapacityUnits)
				}
				f2elem.SetProvisionedThroughput(f2elemf3)
			}
			f2 = append(f2, f2elem)
		}
		res.SetGlobalSecondaryIndexes(f2)
	}
	if r.ko.Spec.KeySchema != nil {
		f3 := []*svcsdk.KeySchemaElement{}
		for _, f3iter := range r.ko.Spec.KeySchema {
			f3elem := &svcsdk.KeySchemaElement{}
			if f3iter.AttributeName != nil {
				f3elem.SetAttributeName(*f3iter.AttributeName)
			}
			if f3iter.KeyType != nil {
				f3elem.SetKeyType(*f3iter.KeyType)
			}
			f3 = append(f3, f3elem)
		}
		res.SetKeySchema(f3)
	}
	if r.ko.Spec.LocalSecondaryIndexes != nil {
		f4 := []*svcsdk.LocalSecondaryIndex{}
		for _, f4iter := range r.ko.Spec.LocalSecondaryIndexes {
			f4elem := &svcsdk.LocalSecondaryIndex{}
			if f4iter.IndexName != nil {
				f4elem.SetIndexName(*f4iter.IndexName)
			}
			if f4iter.KeySchema != nil {
				f4elemf1 := []*svcsdk.KeySchemaElement{}
				for _, f4elemf1iter := range f4iter.KeySchema {
					f4elemf1elem := &svcsdk.KeySchemaElement{}
					if f4elemf1iter.AttributeName != nil {
						f4elemf1elem.SetAttributeName(*f4elemf1iter.AttributeName)
					}
					if f4elemf1iter.KeyType != nil {
						f4elemf1elem.SetKeyType(*f4elemf1iter.KeyType)
					}
					f4elemf1 = append(f4elemf1, f4elemf1elem)
				}
				f4elem.SetKeySchema(f4elemf1)
			}
			if f4iter.Projection != nil {
				f4elemf2 := &svcsdk.Projection{}
				if f4iter.Projection.NonKeyAttributes != nil {
					f4elemf2f0 := []*string{}
					for _, f4elemf2f0iter := range f4iter.Projection.NonKeyAttributes {
						var f4elemf2f0elem string
						f4elemf2f0elem = *f4elemf2f0iter
						f4elemf2f0 = append(f4elemf2f0, &f4elemf2f0elem)
					}
					f4elemf2.SetNonKeyAttributes(f4elemf2f0)
				}
				if f4iter.Projection.ProjectionType != nil {
					f4elemf2.SetProjectionType(*f4iter.Projection.ProjectionType)
				}
				f4elem.SetProjection(f4elemf2)
			}
			f4 = append(f4, f4elem)
		}
		res.SetLocalSecondaryIndexes(f4)
	}
	if r.ko.Spec.ProvisionedThroughput != nil {
		f5 := &svcsdk.ProvisionedThroughput{}
		if r.ko.Spec.ProvisionedThroughput.ReadCapacityUnits != nil {
			f5.SetReadCapacityUnits(*r.ko.Spec.ProvisionedThroughput.ReadCapacityUnits)
		}
		if r.ko.Spec.ProvisionedThroughput.WriteCapacityUnits != nil {
			f5.SetWriteCapacityUnits(*r.ko.Spec.ProvisionedThroughput.WriteCapacityUnits)
		}
		res.SetProvisionedThroughput(f5)
	}
	if r.ko.Spec.SSESpecification != nil {
		f6 := &svcsdk.SSESpecification{}
		if r.ko.Spec.SSESpecification.Enabled != nil {
			f6.SetEnabled(*r.ko.Spec.SSESpecification.Enabled)
		}
		if r.ko.Spec.SSESpecification.KMSMasterKeyID != nil {
			f6.SetKMSMasterKeyId(*r.ko.Spec.SSESpecification.KMSMasterKeyID)
		}
		if r.ko.Spec.SSESpecification.SSEType != nil {
			f6.SetSSEType(*r.ko.Spec.SSESpecification.SSEType)
		}
		res.SetSSESpecification(f6)
	}
	if r.ko.Spec.StreamSpecification != nil {
		f7 := &svcsdk.StreamSpecification{}
		if r.ko.Spec.StreamSpecification.StreamEnabled != nil {
			f7.SetStreamEnabled(*r.ko.Spec.StreamSpecification.StreamEnabled)
		}
		if r.ko.Spec.StreamSpecification.StreamViewType != nil {
			f7.SetStreamViewType(*r.ko.Spec.StreamSpecification.StreamViewType)
		}
		res.SetStreamSpecification(f7)
	}
	if r.ko.Spec.TableName != nil {
		res.SetTableName(*r.ko.Spec.TableName)
	}
	if r.ko.Spec.Tags != nil {
		f9 := []*svcsdk.Tag{}
		for _, f9iter := range r.ko.Spec.Tags {
			f9elem := &svcsdk.Tag{}
			if f9iter.Key != nil {
				f9elem.SetKey(*f9iter.Key)
			}
			if f9iter.Value != nil {
				f9elem.SetValue(*f9iter.Value)
			}
			f9 = append(f9, f9elem)
		}
		res.SetTags(f9)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_EC2_LaunchTemplate_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "LaunchTemplate")
	require.NotNil(crd)

	// LaunchTemplateName is in the LaunchTemplate resource's CreateTopicInput shape and also
	// returned in the CreateLaunchTemplateResult shape, so it should have
	// Go code to set the Input Shape member from the Spec field but not set a
	// Status field from the Create Output Shape member
	expected := `
	if r.ko.Spec.ClientToken != nil {
		res.SetClientToken(*r.ko.Spec.ClientToken)
	}
	if r.ko.Spec.DryRun != nil {
		res.SetDryRun(*r.ko.Spec.DryRun)
	}
	if r.ko.Spec.LaunchTemplateData != nil {
		f2 := &svcsdk.RequestLaunchTemplateData{}
		if r.ko.Spec.LaunchTemplateData.BlockDeviceMappings != nil {
			f2f0 := []*svcsdk.LaunchTemplateBlockDeviceMappingRequest{}
			for _, f2f0iter := range r.ko.Spec.LaunchTemplateData.BlockDeviceMappings {
				f2f0elem := &svcsdk.LaunchTemplateBlockDeviceMappingRequest{}
				if f2f0iter.DeviceName != nil {
					f2f0elem.SetDeviceName(*f2f0iter.DeviceName)
				}
				if f2f0iter.EBS != nil {
					f2f0elemf1 := &svcsdk.LaunchTemplateEbsBlockDeviceRequest{}
					if f2f0iter.EBS.DeleteOnTermination != nil {
						f2f0elemf1.SetDeleteOnTermination(*f2f0iter.EBS.DeleteOnTermination)
					}
					if f2f0iter.EBS.Encrypted != nil {
						f2f0elemf1.SetEncrypted(*f2f0iter.EBS.Encrypted)
					}
					if f2f0iter.EBS.IOPS != nil {
						f2f0elemf1.SetIops(*f2f0iter.EBS.IOPS)
					}
					if f2f0iter.EBS.KMSKeyID != nil {
						f2f0elemf1.SetKmsKeyId(*f2f0iter.EBS.KMSKeyID)
					}
					if f2f0iter.EBS.SnapshotID != nil {
						f2f0elemf1.SetSnapshotId(*f2f0iter.EBS.SnapshotID)
					}
					if f2f0iter.EBS.Throughput != nil {
						f2f0elemf1.SetThroughput(*f2f0iter.EBS.Throughput)
					}
					if f2f0iter.EBS.VolumeSize != nil {
						f2f0elemf1.SetVolumeSize(*f2f0iter.EBS.VolumeSize)
					}
					if f2f0iter.EBS.VolumeType != nil {
						f2f0elemf1.SetVolumeType(*f2f0iter.EBS.VolumeType)
					}
					f2f0elem.SetEbs(f2f0elemf1)
				}
				if f2f0iter.NoDevice != nil {
					f2f0elem.SetNoDevice(*f2f0iter.NoDevice)
				}
				if f2f0iter.VirtualName != nil {
					f2f0elem.SetVirtualName(*f2f0iter.VirtualName)
				}
				f2f0 = append(f2f0, f2f0elem)
			}
			f2.SetBlockDeviceMappings(f2f0)
		}
		if r.ko.Spec.LaunchTemplateData.CapacityReservationSpecification != nil {
			f2f1 := &svcsdk.LaunchTemplateCapacityReservationSpecificationRequest{}
			if r.ko.Spec.LaunchTemplateData.CapacityReservationSpecification.CapacityReservationPreference != nil {
				f2f1.SetCapacityReservationPreference(*r.ko.Spec.LaunchTemplateData.CapacityReservationSpecification.CapacityReservationPreference)
			}
			if r.ko.Spec.LaunchTemplateData.CapacityReservationSpecification.CapacityReservationTarget != nil {
				f2f1f1 := &svcsdk.CapacityReservationTarget{}
				if r.ko.Spec.LaunchTemplateData.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationID != nil {
					f2f1f1.SetCapacityReservationId(*r.ko.Spec.LaunchTemplateData.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationID)
				}
				if r.ko.Spec.LaunchTemplateData.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationResourceGroupARN != nil {
					f2f1f1.SetCapacityReservationResourceGroupArn(*r.ko.Spec.LaunchTemplateData.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationResourceGroupARN)
				}
				f2f1.SetCapacityReservationTarget(f2f1f1)
			}
			f2.SetCapacityReservationSpecification(f2f1)
		}
		if r.ko.Spec.LaunchTemplateData.CPUOptions != nil {
			f2f2 := &svcsdk.LaunchTemplateCpuOptionsRequest{}
			if r.ko.Spec.LaunchTemplateData.CPUOptions.CoreCount != nil {
				f2f2.SetCoreCount(*r.ko.Spec.LaunchTemplateData.CPUOptions.CoreCount)
			}
			if r.ko.Spec.LaunchTemplateData.CPUOptions.ThreadsPerCore != nil {
				f2f2.SetThreadsPerCore(*r.ko.Spec.LaunchTemplateData.CPUOptions.ThreadsPerCore)
			}
			f2.SetCpuOptions(f2f2)
		}
		if r.ko.Spec.LaunchTemplateData.CreditSpecification != nil {
			f2f3 := &svcsdk.CreditSpecificationRequest{}
			if r.ko.Spec.LaunchTemplateData.CreditSpecification.CPUCredits != nil {
				f2f3.SetCpuCredits(*r.ko.Spec.LaunchTemplateData.CreditSpecification.CPUCredits)
			}
			f2.SetCreditSpecification(f2f3)
		}
		if r.ko.Spec.LaunchTemplateData.DisableAPITermination != nil {
			f2.SetDisableApiTermination(*r.ko.Spec.LaunchTemplateData.DisableAPITermination)
		}
		if r.ko.Spec.LaunchTemplateData.EBSOptimized != nil {
			f2.SetEbsOptimized(*r.ko.Spec.LaunchTemplateData.EBSOptimized)
		}
		if r.ko.Spec.LaunchTemplateData.ElasticGPUSpecifications != nil {
			f2f6 := []*svcsdk.ElasticGpuSpecification{}
			for _, f2f6iter := range r.ko.Spec.LaunchTemplateData.ElasticGPUSpecifications {
				f2f6elem := &svcsdk.ElasticGpuSpecification{}
				if f2f6iter.Type != nil {
					f2f6elem.SetType(*f2f6iter.Type)
				}
				f2f6 = append(f2f6, f2f6elem)
			}
			f2.SetElasticGpuSpecifications(f2f6)
		}
		if r.ko.Spec.LaunchTemplateData.ElasticInferenceAccelerators != nil {
			f2f7 := []*svcsdk.LaunchTemplateElasticInferenceAccelerator{}
			for _, f2f7iter := range r.ko.Spec.LaunchTemplateData.ElasticInferenceAccelerators {
				f2f7elem := &svcsdk.LaunchTemplateElasticInferenceAccelerator{}
				if f2f7iter.Count != nil {
					f2f7elem.SetCount(*f2f7iter.Count)
				}
				if f2f7iter.Type != nil {
					f2f7elem.SetType(*f2f7iter.Type)
				}
				f2f7 = append(f2f7, f2f7elem)
			}
			f2.SetElasticInferenceAccelerators(f2f7)
		}
		if r.ko.Spec.LaunchTemplateData.EnclaveOptions != nil {
			f2f8 := &svcsdk.LaunchTemplateEnclaveOptionsRequest{}
			if r.ko.Spec.LaunchTemplateData.EnclaveOptions.Enabled != nil {
				f2f8.SetEnabled(*r.ko.Spec.LaunchTemplateData.EnclaveOptions.Enabled)
			}
			f2.SetEnclaveOptions(f2f8)
		}
		if r.ko.Spec.LaunchTemplateData.HibernationOptions != nil {
			f2f9 := &svcsdk.LaunchTemplateHibernationOptionsRequest{}
			if r.ko.Spec.LaunchTemplateData.HibernationOptions.Configured != nil {
				f2f9.SetConfigured(*r.ko.Spec.LaunchTemplateData.HibernationOptions.Configured)
			}
			f2.SetHibernationOptions(f2f9)
		}
		if r.ko.Spec.LaunchTemplateData.IAMInstanceProfile != nil {
			f2f10 := &svcsdk.LaunchTemplateIamInstanceProfileSpecificationRequest{}
			if r.ko.Spec.LaunchTemplateData.IAMInstanceProfile.ARN != nil {
				f2f10.SetArn(*r.ko.Spec.LaunchTemplateData.IAMInstanceProfile.ARN)
			}
			if r.ko.Spec.LaunchTemplateData.IAMInstanceProfile.Name != nil {
				f2f10.SetName(*r.ko.Spec.LaunchTemplateData.IAMInstanceProfile.Name)
			}
			f2.SetIamInstanceProfile(f2f10)
		}
		if r.ko.Spec.LaunchTemplateData.ImageID != nil {
			f2.SetImageId(*r.ko.Spec.LaunchTemplateData.ImageID)
		}
		if r.ko.Spec.LaunchTemplateData.InstanceInitiatedShutdownBehavior != nil {
			f2.SetInstanceInitiatedShutdownBehavior(*r.ko.Spec.LaunchTemplateData.InstanceInitiatedShutdownBehavior)
		}
		if r.ko.Spec.LaunchTemplateData.InstanceMarketOptions != nil {
			f2f13 := &svcsdk.LaunchTemplateInstanceMarketOptionsRequest{}
			if r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.MarketType != nil {
				f2f13.SetMarketType(*r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.MarketType)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions != nil {
				f2f13f1 := &svcsdk.LaunchTemplateSpotMarketOptionsRequest{}
				if r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions.BlockDurationMinutes != nil {
					f2f13f1.SetBlockDurationMinutes(*r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions.BlockDurationMinutes)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions.InstanceInterruptionBehavior != nil {
					f2f13f1.SetInstanceInterruptionBehavior(*r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions.InstanceInterruptionBehavior)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions.MaxPrice != nil {
					f2f13f1.SetMaxPrice(*r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions.MaxPrice)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions.SpotInstanceType != nil {
					f2f13f1.SetSpotInstanceType(*r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions.SpotInstanceType)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions.ValidUntil != nil {
					f2f13f1.SetValidUntil(r.ko.Spec.LaunchTemplateData.InstanceMarketOptions.SpotOptions.ValidUntil.Time)
				}
				f2f13.SetSpotOptions(f2f13f1)
			}
			f2.SetInstanceMarketOptions(f2f13)
		}
		if r.ko.Spec.LaunchTemplateData.InstanceRequirements != nil {
			f2f14 := &svcsdk.InstanceRequirementsRequest{}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorCount != nil {
				f2f14f0 := &svcsdk.AcceleratorCountRequest{}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorCount.Max != nil {
					f2f14f0.SetMax(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorCount.Max)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorCount.Min != nil {
					f2f14f0.SetMin(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorCount.Min)
				}
				f2f14.SetAcceleratorCount(f2f14f0)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorManufacturers != nil {
				f2f14f1 := []*string{}
				for _, f2f14f1iter := range r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorManufacturers {
					var f2f14f1elem string
					f2f14f1elem = *f2f14f1iter
					f2f14f1 = append(f2f14f1, &f2f14f1elem)
				}
				f2f14.SetAcceleratorManufacturers(f2f14f1)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorNames != nil {
				f2f14f2 := []*string{}
				for _, f2f14f2iter := range r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorNames {
					var f2f14f2elem string
					f2f14f2elem = *f2f14f2iter
					f2f14f2 = append(f2f14f2, &f2f14f2elem)
				}
				f2f14.SetAcceleratorNames(f2f14f2)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorTotalMemoryMiB != nil {
				f2f14f3 := &svcsdk.AcceleratorTotalMemoryMiBRequest{}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorTotalMemoryMiB.Max != nil {
					f2f14f3.SetMax(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorTotalMemoryMiB.Max)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorTotalMemoryMiB.Min != nil {
					f2f14f3.SetMin(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorTotalMemoryMiB.Min)
				}
				f2f14.SetAcceleratorTotalMemoryMiB(f2f14f3)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorTypes != nil {
				f2f14f4 := []*string{}
				for _, f2f14f4iter := range r.ko.Spec.LaunchTemplateData.InstanceRequirements.AcceleratorTypes {
					var f2f14f4elem string
					f2f14f4elem = *f2f14f4iter
					f2f14f4 = append(f2f14f4, &f2f14f4elem)
				}
				f2f14.SetAcceleratorTypes(f2f14f4)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.BareMetal != nil {
				f2f14.SetBareMetal(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.BareMetal)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.BaselineEBSBandwidthMbps != nil {
				f2f14f6 := &svcsdk.BaselineEbsBandwidthMbpsRequest{}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.BaselineEBSBandwidthMbps.Max != nil {
					f2f14f6.SetMax(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.BaselineEBSBandwidthMbps.Max)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.BaselineEBSBandwidthMbps.Min != nil {
					f2f14f6.SetMin(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.BaselineEBSBandwidthMbps.Min)
				}
				f2f14.SetBaselineEbsBandwidthMbps(f2f14f6)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.BurstablePerformance != nil {
				f2f14.SetBurstablePerformance(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.BurstablePerformance)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.CPUManufacturers != nil {
				f2f14f8 := []*string{}
				for _, f2f14f8iter := range r.ko.Spec.LaunchTemplateData.InstanceRequirements.CPUManufacturers {
					var f2f14f8elem string
					f2f14f8elem = *f2f14f8iter
					f2f14f8 = append(f2f14f8, &f2f14f8elem)
				}
				f2f14.SetCpuManufacturers(f2f14f8)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.ExcludedInstanceTypes != nil {
				f2f14f9 := []*string{}
				for _, f2f14f9iter := range r.ko.Spec.LaunchTemplateData.InstanceRequirements.ExcludedInstanceTypes {
					var f2f14f9elem string
					f2f14f9elem = *f2f14f9iter
					f2f14f9 = append(f2f14f9, &f2f14f9elem)
				}
				f2f14.SetExcludedInstanceTypes(f2f14f9)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.InstanceGenerations != nil {
				f2f14f10 := []*string{}
				for _, f2f14f10iter := range r.ko.Spec.LaunchTemplateData.InstanceRequirements.InstanceGenerations {
					var f2f14f10elem string
					f2f14f10elem = *f2f14f10iter
					f2f14f10 = append(f2f14f10, &f2f14f10elem)
				}
				f2f14.SetInstanceGenerations(f2f14f10)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.LocalStorage != nil {
				f2f14.SetLocalStorage(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.LocalStorage)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.LocalStorageTypes != nil {
				f2f14f12 := []*string{}
				for _, f2f14f12iter := range r.ko.Spec.LaunchTemplateData.InstanceRequirements.LocalStorageTypes {
					var f2f14f12elem string
					f2f14f12elem = *f2f14f12iter
					f2f14f12 = append(f2f14f12, &f2f14f12elem)
				}
				f2f14.SetLocalStorageTypes(f2f14f12)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.MemoryGiBPerVCPU != nil {
				f2f14f13 := &svcsdk.MemoryGiBPerVCpuRequest{}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.MemoryGiBPerVCPU.Max != nil {
					f2f14f13.SetMax(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.MemoryGiBPerVCPU.Max)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.MemoryGiBPerVCPU.Min != nil {
					f2f14f13.SetMin(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.MemoryGiBPerVCPU.Min)
				}
				f2f14.SetMemoryGiBPerVCpu(f2f14f13)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.MemoryMiB != nil {
				f2f14f14 := &svcsdk.MemoryMiBRequest{}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.MemoryMiB.Max != nil {
					f2f14f14.SetMax(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.MemoryMiB.Max)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.MemoryMiB.Min != nil {
					f2f14f14.SetMin(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.MemoryMiB.Min)
				}
				f2f14.SetMemoryMiB(f2f14f14)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.NetworkInterfaceCount != nil {
				f2f14f15 := &svcsdk.NetworkInterfaceCountRequest{}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.NetworkInterfaceCount.Max != nil {
					f2f14f15.SetMax(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.NetworkInterfaceCount.Max)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.NetworkInterfaceCount.Min != nil {
					f2f14f15.SetMin(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.NetworkInterfaceCount.Min)
				}
				f2f14.SetNetworkInterfaceCount(f2f14f15)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.OnDemandMaxPricePercentageOverLowestPrice != nil {
				f2f14.SetOnDemandMaxPricePercentageOverLowestPrice(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.OnDemandMaxPricePercentageOverLowestPrice)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.RequireHibernateSupport != nil {
				f2f14.SetRequireHibernateSupport(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.RequireHibernateSupport)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.SpotMaxPricePercentageOverLowestPrice != nil {
				f2f14.SetSpotMaxPricePercentageOverLowestPrice(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.SpotMaxPricePercentageOverLowestPrice)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.TotalLocalStorageGB != nil {
				f2f14f19 := &svcsdk.TotalLocalStorageGBRequest{}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.TotalLocalStorageGB.Max != nil {
					f2f14f19.SetMax(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.TotalLocalStorageGB.Max)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.TotalLocalStorageGB.Min != nil {
					f2f14f19.SetMin(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.TotalLocalStorageGB.Min)
				}
				f2f14.SetTotalLocalStorageGB(f2f14f19)
			}
			if r.ko.Spec.LaunchTemplateData.InstanceRequirements.VCPUCount != nil {
				f2f14f20 := &svcsdk.VCpuCountRangeRequest{}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.VCPUCount.Max != nil {
					f2f14f20.SetMax(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.VCPUCount.Max)
				}
				if r.ko.Spec.LaunchTemplateData.InstanceRequirements.VCPUCount.Min != nil {
					f2f14f20.SetMin(*r.ko.Spec.LaunchTemplateData.InstanceRequirements.VCPUCount.Min)
				}
				f2f14.SetVCpuCount(f2f14f20)
			}
			f2.SetInstanceRequirements(f2f14)
		}
		if r.ko.Spec.LaunchTemplateData.InstanceType != nil {
			f2.SetInstanceType(*r.ko.Spec.LaunchTemplateData.InstanceType)
		}
		if r.ko.Spec.LaunchTemplateData.KernelID != nil {
			f2.SetKernelId(*r.ko.Spec.LaunchTemplateData.KernelID)
		}
		if r.ko.Spec.LaunchTemplateData.KeyName != nil {
			f2.SetKeyName(*r.ko.Spec.LaunchTemplateData.KeyName)
		}
		if r.ko.Spec.LaunchTemplateData.LicenseSpecifications != nil {
			f2f18 := []*svcsdk.LaunchTemplateLicenseConfigurationRequest{}
			for _, f2f18iter := range r.ko.Spec.LaunchTemplateData.LicenseSpecifications {
				f2f18elem := &svcsdk.LaunchTemplateLicenseConfigurationRequest{}
				if f2f18iter.LicenseConfigurationARN != nil {
					f2f18elem.SetLicenseConfigurationArn(*f2f18iter.LicenseConfigurationARN)
				}
				f2f18 = append(f2f18, f2f18elem)
			}
			f2.SetLicenseSpecifications(f2f18)
		}
		if r.ko.Spec.LaunchTemplateData.MetadataOptions != nil {
			f2f19 := &svcsdk.LaunchTemplateInstanceMetadataOptionsRequest{}
			if r.ko.Spec.LaunchTemplateData.MetadataOptions.HTTPEndpoint != nil {
				f2f19.SetHttpEndpoint(*r.ko.Spec.LaunchTemplateData.MetadataOptions.HTTPEndpoint)
			}
			if r.ko.Spec.LaunchTemplateData.MetadataOptions.HTTPProtocolIPv6 != nil {
				f2f19.SetHttpProtocolIpv6(*r.ko.Spec.LaunchTemplateData.MetadataOptions.HTTPProtocolIPv6)
			}
			if r.ko.Spec.LaunchTemplateData.MetadataOptions.HTTPPutResponseHopLimit != nil {
				f2f19.SetHttpPutResponseHopLimit(*r.ko.Spec.LaunchTemplateData.MetadataOptions.HTTPPutResponseHopLimit)
			}
			if r.ko.Spec.LaunchTemplateData.MetadataOptions.HTTPTokens != nil {
				f2f19.SetHttpTokens(*r.ko.Spec.LaunchTemplateData.MetadataOptions.HTTPTokens)
			}
			f2.SetMetadataOptions(f2f19)
		}
		if r.ko.Spec.LaunchTemplateData.Monitoring != nil {
			f2f20 := &svcsdk.LaunchTemplatesMonitoringRequest{}
			if r.ko.Spec.LaunchTemplateData.Monitoring.Enabled != nil {
				f2f20.SetEnabled(*r.ko.Spec.LaunchTemplateData.Monitoring.Enabled)
			}
			f2.SetMonitoring(f2f20)
		}
		if r.ko.Spec.LaunchTemplateData.NetworkInterfaces != nil {
			f2f21 := []*svcsdk.LaunchTemplateInstanceNetworkInterfaceSpecificationRequest{}
			for _, f2f21iter := range r.ko.Spec.LaunchTemplateData.NetworkInterfaces {
				f2f21elem := &svcsdk.LaunchTemplateInstanceNetworkInterfaceSpecificationRequest{}
				if f2f21iter.AssociateCarrierIPAddress != nil {
					f2f21elem.SetAssociateCarrierIpAddress(*f2f21iter.AssociateCarrierIPAddress)
				}
				if f2f21iter.AssociatePublicIPAddress != nil {
					f2f21elem.SetAssociatePublicIpAddress(*f2f21iter.AssociatePublicIPAddress)
				}
				if f2f21iter.DeleteOnTermination != nil {
					f2f21elem.SetDeleteOnTermination(*f2f21iter.DeleteOnTermination)
				}
				if f2f21iter.Description != nil {
					f2f21elem.SetDescription(*f2f21iter.Description)
				}
				if f2f21iter.DeviceIndex != nil {
					f2f21elem.SetDeviceIndex(*f2f21iter.DeviceIndex)
				}
				if f2f21iter.Groups != nil {
					f2f21elemf5 := []*string{}
					for _, f2f21elemf5iter := range f2f21iter.Groups {
						var f2f21elemf5elem string
						f2f21elemf5elem = *f2f21elemf5iter
						f2f21elemf5 = append(f2f21elemf5, &f2f21elemf5elem)
					}
					f2f21elem.SetGroups(f2f21elemf5)
				}
				if f2f21iter.InterfaceType != nil {
					f2f21elem.SetInterfaceType(*f2f21iter.InterfaceType)
				}
				if f2f21iter.IPv4PrefixCount != nil {
					f2f21elem.SetIpv4PrefixCount(*f2f21iter.IPv4PrefixCount)
				}
				if f2f21iter.IPv4Prefixes != nil {
					f2f21elemf8 := []*svcsdk.Ipv4PrefixSpecificationRequest{}
					for _, f2f21elemf8iter := range f2f21iter.IPv4Prefixes {
						f2f21elemf8elem := &svcsdk.Ipv4PrefixSpecificationRequest{}
						if f2f21elemf8iter.IPv4Prefix != nil {
							f2f21elemf8elem.SetIpv4Prefix(*f2f21elemf8iter.IPv4Prefix)
						}
						f2f21elemf8 = append(f2f21elemf8, f2f21elemf8elem)
					}
					f2f21elem.SetIpv4Prefixes(f2f21elemf8)
				}
				if f2f21iter.IPv6AddressCount != nil {
					f2f21elem.SetIpv6AddressCount(*f2f21iter.IPv6AddressCount)
				}
				if f2f21iter.IPv6Addresses != nil {
					f2f21elemf10 := []*svcsdk.InstanceIpv6AddressRequest{}
					for _, f2f21elemf10iter := range f2f21iter.IPv6Addresses {
						f2f21elemf10elem := &svcsdk.InstanceIpv6AddressRequest{}
						if f2f21elemf10iter.IPv6Address != nil {
							f2f21elemf10elem.SetIpv6Address(*f2f21elemf10iter.IPv6Address)
						}
						f2f21elemf10 = append(f2f21elemf10, f2f21elemf10elem)
					}
					f2f21elem.SetIpv6Addresses(f2f21elemf10)
				}
				if f2f21iter.IPv6PrefixCount != nil {
					f2f21elem.SetIpv6PrefixCount(*f2f21iter.IPv6PrefixCount)
				}
				if f2f21iter.IPv6Prefixes != nil {
					f2f21elemf12 := []*svcsdk.Ipv6PrefixSpecificationRequest{}
					for _, f2f21elemf12iter := range f2f21iter.IPv6Prefixes {
						f2f21elemf12elem := &svcsdk.Ipv6PrefixSpecificationRequest{}
						if f2f21elemf12iter.IPv6Prefix != nil {
							f2f21elemf12elem.SetIpv6Prefix(*f2f21elemf12iter.IPv6Prefix)
						}
						f2f21elemf12 = append(f2f21elemf12, f2f21elemf12elem)
					}
					f2f21elem.SetIpv6Prefixes(f2f21elemf12)
				}
				if f2f21iter.NetworkCardIndex != nil {
					f2f21elem.SetNetworkCardIndex(*f2f21iter.NetworkCardIndex)
				}
				if f2f21iter.NetworkInterfaceID != nil {
					f2f21elem.SetNetworkInterfaceId(*f2f21iter.NetworkInterfaceID)
				}
				if f2f21iter.PrivateIPAddress != nil {
					f2f21elem.SetPrivateIpAddress(*f2f21iter.PrivateIPAddress)
				}
				if f2f21iter.PrivateIPAddresses != nil {
					f2f21elemf16 := []*svcsdk.PrivateIpAddressSpecification{}
					for _, f2f21elemf16iter := range f2f21iter.PrivateIPAddresses {
						f2f21elemf16elem := &svcsdk.PrivateIpAddressSpecification{}
						if f2f21elemf16iter.Primary != nil {
							f2f21elemf16elem.SetPrimary(*f2f21elemf16iter.Primary)
						}
						if f2f21elemf16iter.PrivateIPAddress != nil {
							f2f21elemf16elem.SetPrivateIpAddress(*f2f21elemf16iter.PrivateIPAddress)
						}
						f2f21elemf16 = append(f2f21elemf16, f2f21elemf16elem)
					}
					f2f21elem.SetPrivateIpAddresses(f2f21elemf16)
				}
				if f2f21iter.SecondaryPrivateIPAddressCount != nil {
					f2f21elem.SetSecondaryPrivateIpAddressCount(*f2f21iter.SecondaryPrivateIPAddressCount)
				}
				if f2f21iter.SubnetID != nil {
					f2f21elem.SetSubnetId(*f2f21iter.SubnetID)
				}
				f2f21 = append(f2f21, f2f21elem)
			}
			f2.SetNetworkInterfaces(f2f21)
		}
		if r.ko.Spec.LaunchTemplateData.Placement != nil {
			f2f22 := &svcsdk.LaunchTemplatePlacementRequest{}
			if r.ko.Spec.LaunchTemplateData.Placement.Affinity != nil {
				f2f22.SetAffinity(*r.ko.Spec.LaunchTemplateData.Placement.Affinity)
			}
			if r.ko.Spec.LaunchTemplateData.Placement.AvailabilityZone != nil {
				f2f22.SetAvailabilityZone(*r.ko.Spec.LaunchTemplateData.Placement.AvailabilityZone)
			}
			if r.ko.Spec.LaunchTemplateData.Placement.GroupName != nil {
				f2f22.SetGroupName(*r.ko.Spec.LaunchTemplateData.Placement.GroupName)
			}
			if r.ko.Spec.LaunchTemplateData.Placement.HostID != nil {
				f2f22.SetHostId(*r.ko.Spec.LaunchTemplateData.Placement.HostID)
			}
			if r.ko.Spec.LaunchTemplateData.Placement.HostResourceGroupARN != nil {
				f2f22.SetHostResourceGroupArn(*r.ko.Spec.LaunchTemplateData.Placement.HostResourceGroupARN)
			}
			if r.ko.Spec.LaunchTemplateData.Placement.PartitionNumber != nil {
				f2f22.SetPartitionNumber(*r.ko.Spec.LaunchTemplateData.Placement.PartitionNumber)
			}
			if r.ko.Spec.LaunchTemplateData.Placement.SpreadDomain != nil {
				f2f22.SetSpreadDomain(*r.ko.Spec.LaunchTemplateData.Placement.SpreadDomain)
			}
			if r.ko.Spec.LaunchTemplateData.Placement.Tenancy != nil {
				f2f22.SetTenancy(*r.ko.Spec.LaunchTemplateData.Placement.Tenancy)
			}
			f2.SetPlacement(f2f22)
		}
		if r.ko.Spec.LaunchTemplateData.RAMDiskID != nil {
			f2.SetRamDiskId(*r.ko.Spec.LaunchTemplateData.RAMDiskID)
		}
		if r.ko.Spec.LaunchTemplateData.SecurityGroupIDs != nil {
			f2f24 := []*string{}
			for _, f2f24iter := range r.ko.Spec.LaunchTemplateData.SecurityGroupIDs {
				var f2f24elem string
				f2f24elem = *f2f24iter
				f2f24 = append(f2f24, &f2f24elem)
			}
			f2.SetSecurityGroupIds(f2f24)
		}
		if r.ko.Spec.LaunchTemplateData.SecurityGroups != nil {
			f2f25 := []*string{}
			for _, f2f25iter := range r.ko.Spec.LaunchTemplateData.SecurityGroups {
				var f2f25elem string
				f2f25elem = *f2f25iter
				f2f25 = append(f2f25, &f2f25elem)
			}
			f2.SetSecurityGroups(f2f25)
		}
		if r.ko.Spec.LaunchTemplateData.TagSpecifications != nil {
			f2f26 := []*svcsdk.LaunchTemplateTagSpecificationRequest{}
			for _, f2f26iter := range r.ko.Spec.LaunchTemplateData.TagSpecifications {
				f2f26elem := &svcsdk.LaunchTemplateTagSpecificationRequest{}
				if f2f26iter.ResourceType != nil {
					f2f26elem.SetResourceType(*f2f26iter.ResourceType)
				}
				if f2f26iter.Tags != nil {
					f2f26elemf1 := []*svcsdk.Tag{}
					for _, f2f26elemf1iter := range f2f26iter.Tags {
						f2f26elemf1elem := &svcsdk.Tag{}
						if f2f26elemf1iter.Key != nil {
							f2f26elemf1elem.SetKey(*f2f26elemf1iter.Key)
						}
						if f2f26elemf1iter.Value != nil {
							f2f26elemf1elem.SetValue(*f2f26elemf1iter.Value)
						}
						f2f26elemf1 = append(f2f26elemf1, f2f26elemf1elem)
					}
					f2f26elem.SetTags(f2f26elemf1)
				}
				f2f26 = append(f2f26, f2f26elem)
			}
			f2.SetTagSpecifications(f2f26)
		}
		if r.ko.Spec.LaunchTemplateData.UserData != nil {
			f2.SetUserData(*r.ko.Spec.LaunchTemplateData.UserData)
		}
		res.SetLaunchTemplateData(f2)
	}
	if r.ko.Spec.LaunchTemplateName != nil {
		res.SetLaunchTemplateName(*r.ko.Spec.LaunchTemplateName)
	}
	if r.ko.Spec.TagSpecifications != nil {
		f4 := []*svcsdk.TagSpecification{}
		for _, f4iter := range r.ko.Spec.TagSpecifications {
			f4elem := &svcsdk.TagSpecification{}
			if f4iter.ResourceType != nil {
				f4elem.SetResourceType(*f4iter.ResourceType)
			}
			if f4iter.Tags != nil {
				f4elemf1 := []*svcsdk.Tag{}
				for _, f4elemf1iter := range f4iter.Tags {
					f4elemf1elem := &svcsdk.Tag{}
					if f4elemf1iter.Key != nil {
						f4elemf1elem.SetKey(*f4elemf1iter.Key)
					}
					if f4elemf1iter.Value != nil {
						f4elemf1elem.SetValue(*f4elemf1iter.Value)
					}
					f4elemf1 = append(f4elemf1, f4elemf1elem)
				}
				f4elem.SetTags(f4elemf1)
			}
			f4 = append(f4, f4elem)
		}
		res.SetTagSpecifications(f4)
	}
	if r.ko.Spec.VersionDescription != nil {
		res.SetVersionDescription(*r.ko.Spec.VersionDescription)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_ECR_Repository_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ecr")

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)

	// ImageScanningConfiguration is in the Repository resource's
	// CreateRepositoryInput shape and also returned in the
	// CreateRepositoryOutput shape, so it should produce Go code to set the
	// appropriate input shape member.
	expected := `
	if r.ko.Spec.ImageScanningConfiguration != nil {
		f0 := &svcsdk.ImageScanningConfiguration{}
		if r.ko.Spec.ImageScanningConfiguration.ScanOnPush != nil {
			f0.SetScanOnPush(*r.ko.Spec.ImageScanningConfiguration.ScanOnPush)
		}
		res.SetImageScanningConfiguration(f0)
	}
	if r.ko.Spec.ImageTagMutability != nil {
		res.SetImageTagMutability(*r.ko.Spec.ImageTagMutability)
	}
	if r.ko.Spec.RepositoryName != nil {
		res.SetRepositoryName(*r.ko.Spec.RepositoryName)
	}
	if r.ko.Spec.Tags != nil {
		f3 := []*svcsdk.Tag{}
		for _, f3iter := range r.ko.Spec.Tags {
			f3elem := &svcsdk.Tag{}
			if f3iter.Key != nil {
				f3elem.SetKey(*f3iter.Key)
			}
			if f3iter.Value != nil {
				f3elem.SetValue(*f3iter.Value)
			}
			f3 = append(f3, f3elem)
		}
		res.SetTags(f3)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_Elasticache_ReplicationGroup_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")

	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	expected := `
	if r.ko.Spec.AtRestEncryptionEnabled != nil {
		res.SetAtRestEncryptionEnabled(*r.ko.Spec.AtRestEncryptionEnabled)
	}
	if r.ko.Spec.AuthToken != nil {
		tmpSecret, err := rm.rr.SecretValueFromReference(ctx, r.ko.Spec.AuthToken)
		if err != nil {
			return nil, ackrequeue.Needed(err)
		}
		if tmpSecret != "" {
			res.SetAuthToken(tmpSecret)
		}
	}
	if r.ko.Spec.AutoMinorVersionUpgrade != nil {
		res.SetAutoMinorVersionUpgrade(*r.ko.Spec.AutoMinorVersionUpgrade)
	}
	if r.ko.Spec.AutomaticFailoverEnabled != nil {
		res.SetAutomaticFailoverEnabled(*r.ko.Spec.AutomaticFailoverEnabled)
	}
	if r.ko.Spec.CacheNodeType != nil {
		res.SetCacheNodeType(*r.ko.Spec.CacheNodeType)
	}
	if r.ko.Spec.CacheParameterGroupName != nil {
		res.SetCacheParameterGroupName(*r.ko.Spec.CacheParameterGroupName)
	}
	if r.ko.Spec.CacheSecurityGroupNames != nil {
		f6 := []*string{}
		for _, f6iter := range r.ko.Spec.CacheSecurityGroupNames {
			var f6elem string
			f6elem = *f6iter
			f6 = append(f6, &f6elem)
		}
		res.SetCacheSecurityGroupNames(f6)
	}
	if r.ko.Spec.CacheSubnetGroupName != nil {
		res.SetCacheSubnetGroupName(*r.ko.Spec.CacheSubnetGroupName)
	}
	if r.ko.Spec.Engine != nil {
		res.SetEngine(*r.ko.Spec.Engine)
	}
	if r.ko.Spec.EngineVersion != nil {
		res.SetEngineVersion(*r.ko.Spec.EngineVersion)
	}
	if r.ko.Spec.KMSKeyID != nil {
		res.SetKmsKeyId(*r.ko.Spec.KMSKeyID)
	}
	if r.ko.Spec.LogDeliveryConfigurations != nil {
		f11 := []*svcsdk.LogDeliveryConfigurationRequest{}
		for _, f11iter := range r.ko.Spec.LogDeliveryConfigurations {
			f11elem := &svcsdk.LogDeliveryConfigurationRequest{}
			if f11iter.DestinationDetails != nil {
				f11elemf0 := &svcsdk.DestinationDetails{}
				if f11iter.DestinationDetails.CloudWatchLogsDetails != nil {
					f11elemf0f0 := &svcsdk.CloudWatchLogsDestinationDetails{}
					if f11iter.DestinationDetails.CloudWatchLogsDetails.LogGroup != nil {
						f11elemf0f0.SetLogGroup(*f11iter.DestinationDetails.CloudWatchLogsDetails.LogGroup)
					}
					f11elemf0.SetCloudWatchLogsDetails(f11elemf0f0)
				}
				if f11iter.DestinationDetails.KinesisFirehoseDetails != nil {
					f11elemf0f1 := &svcsdk.KinesisFirehoseDestinationDetails{}
					if f11iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream != nil {
						f11elemf0f1.SetDeliveryStream(*f11iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream)
					}
					f11elemf0.SetKinesisFirehoseDetails(f11elemf0f1)
				}
				f11elem.SetDestinationDetails(f11elemf0)
			}
			if f11iter.DestinationType != nil {
				f11elem.SetDestinationType(*f11iter.DestinationType)
			}
			if f11iter.Enabled != nil {
				f11elem.SetEnabled(*f11iter.Enabled)
			}
			if f11iter.LogFormat != nil {
				f11elem.SetLogFormat(*f11iter.LogFormat)
			}
			if f11iter.LogType != nil {
				f11elem.SetLogType(*f11iter.LogType)
			}
			f11 = append(f11, f11elem)
		}
		res.SetLogDeliveryConfigurations(f11)
	}
	if r.ko.Spec.MultiAZEnabled != nil {
		res.SetMultiAZEnabled(*r.ko.Spec.MultiAZEnabled)
	}
	if r.ko.Spec.NodeGroupConfiguration != nil {
		f13 := []*svcsdk.NodeGroupConfiguration{}
		for _, f13iter := range r.ko.Spec.NodeGroupConfiguration {
			f13elem := &svcsdk.NodeGroupConfiguration{}
			if f13iter.NodeGroupID != nil {
				f13elem.SetNodeGroupId(*f13iter.NodeGroupID)
			}
			if f13iter.PrimaryAvailabilityZone != nil {
				f13elem.SetPrimaryAvailabilityZone(*f13iter.PrimaryAvailabilityZone)
			}
			if f13iter.PrimaryOutpostARN != nil {
				f13elem.SetPrimaryOutpostArn(*f13iter.PrimaryOutpostARN)
			}
			if f13iter.ReplicaAvailabilityZones != nil {
				f13elemf3 := []*string{}
				for _, f13elemf3iter := range f13iter.ReplicaAvailabilityZones {
					var f13elemf3elem string
					f13elemf3elem = *f13elemf3iter
					f13elemf3 = append(f13elemf3, &f13elemf3elem)
				}
				f13elem.SetReplicaAvailabilityZones(f13elemf3)
			}
			if f13iter.ReplicaCount != nil {
				f13elem.SetReplicaCount(*f13iter.ReplicaCount)
			}
			if f13iter.ReplicaOutpostARNs != nil {
				f13elemf5 := []*string{}
				for _, f13elemf5iter := range f13iter.ReplicaOutpostARNs {
					var f13elemf5elem string
					f13elemf5elem = *f13elemf5iter
					f13elemf5 = append(f13elemf5, &f13elemf5elem)
				}
				f13elem.SetReplicaOutpostArns(f13elemf5)
			}
			if f13iter.Slots != nil {
				f13elem.SetSlots(*f13iter.Slots)
			}
			f13 = append(f13, f13elem)
		}
		res.SetNodeGroupConfiguration(f13)
	}
	if r.ko.Spec.NotificationTopicARN != nil {
		res.SetNotificationTopicArn(*r.ko.Spec.NotificationTopicARN)
	}
	if r.ko.Spec.NumCacheClusters != nil {
		res.SetNumCacheClusters(*r.ko.Spec.NumCacheClusters)
	}
	if r.ko.Spec.NumNodeGroups != nil {
		res.SetNumNodeGroups(*r.ko.Spec.NumNodeGroups)
	}
	if r.ko.Spec.Port != nil {
		res.SetPort(*r.ko.Spec.Port)
	}
	if r.ko.Spec.PreferredCacheClusterAZs != nil {
		f18 := []*string{}
		for _, f18iter := range r.ko.Spec.PreferredCacheClusterAZs {
			var f18elem string
			f18elem = *f18iter
			f18 = append(f18, &f18elem)
		}
		res.SetPreferredCacheClusterAZs(f18)
	}
	if r.ko.Spec.PreferredMaintenanceWindow != nil {
		res.SetPreferredMaintenanceWindow(*r.ko.Spec.PreferredMaintenanceWindow)
	}
	if r.ko.Spec.PrimaryClusterID != nil {
		res.SetPrimaryClusterId(*r.ko.Spec.PrimaryClusterID)
	}
	if r.ko.Spec.ReplicasPerNodeGroup != nil {
		res.SetReplicasPerNodeGroup(*r.ko.Spec.ReplicasPerNodeGroup)
	}
	if r.ko.Spec.ReplicationGroupDescription != nil {
		res.SetReplicationGroupDescription(*r.ko.Spec.ReplicationGroupDescription)
	}
	if r.ko.Spec.ReplicationGroupID != nil {
		res.SetReplicationGroupId(*r.ko.Spec.ReplicationGroupID)
	}
	if r.ko.Spec.SecurityGroupIDs != nil {
		f24 := []*string{}
		for _, f24iter := range r.ko.Spec.SecurityGroupIDs {
			var f24elem string
			f24elem = *f24iter
			f24 = append(f24, &f24elem)
		}
		res.SetSecurityGroupIds(f24)
	}
	if r.ko.Spec.SnapshotARNs != nil {
		f25 := []*string{}
		for _, f25iter := range r.ko.Spec.SnapshotARNs {
			var f25elem string
			f25elem = *f25iter
			f25 = append(f25, &f25elem)
		}
		res.SetSnapshotArns(f25)
	}
	if r.ko.Spec.SnapshotName != nil {
		res.SetSnapshotName(*r.ko.Spec.SnapshotName)
	}
	if r.ko.Spec.SnapshotRetentionLimit != nil {
		res.SetSnapshotRetentionLimit(*r.ko.Spec.SnapshotRetentionLimit)
	}
	if r.ko.Spec.SnapshotWindow != nil {
		res.SetSnapshotWindow(*r.ko.Spec.SnapshotWindow)
	}
	if r.ko.Spec.Tags != nil {
		f29 := []*svcsdk.Tag{}
		for _, f29iter := range r.ko.Spec.Tags {
			f29elem := &svcsdk.Tag{}
			if f29iter.Key != nil {
				f29elem.SetKey(*f29iter.Key)
			}
			if f29iter.Value != nil {
				f29elem.SetValue(*f29iter.Value)
			}
			f29 = append(f29, f29elem)
		}
		res.SetTags(f29)
	}
	if r.ko.Spec.TransitEncryptionEnabled != nil {
		res.SetTransitEncryptionEnabled(*r.ko.Spec.TransitEncryptionEnabled)
	}
	if r.ko.Spec.UserGroupIDs != nil {
		f31 := []*string{}
		for _, f31iter := range r.ko.Spec.UserGroupIDs {
			var f31elem string
			f31elem = *f31iter
			f31 = append(f31, &f31elem)
		}
		res.SetUserGroupIds(f31)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_Elasticache_ReplicationGroup_ReadMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")

	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	// Elasticache doesn't have a ReadOne operation; only a List/ReadMany
	// operation. Let's verify that the construction of the
	// DescribeCacheClustersInput and processing of the
	// DescribeCacheClustersOutput shapes is correct.
	expected := `
	if r.ko.Spec.ReplicationGroupID != nil {
		res.SetReplicationGroupId(*r.ko.Spec.ReplicationGroupID)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeList, "r.ko", "res", 1),
	)
}

func TestSetSDK_Elasticache_ReplicationGroup_Update_Override_Values(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")

	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	expected := `
	res.SetApplyImmediately(true)
	if r.ko.Spec.AuthToken != nil {
		tmpSecret, err := rm.rr.SecretValueFromReference(ctx, r.ko.Spec.AuthToken)
		if err != nil {
			return nil, ackrequeue.Needed(err)
		}
		if tmpSecret != "" {
			res.SetAuthToken(tmpSecret)
		}
	}
	if r.ko.Spec.AutoMinorVersionUpgrade != nil {
		res.SetAutoMinorVersionUpgrade(*r.ko.Spec.AutoMinorVersionUpgrade)
	}
	if r.ko.Spec.AutomaticFailoverEnabled != nil {
		res.SetAutomaticFailoverEnabled(*r.ko.Spec.AutomaticFailoverEnabled)
	}
	if r.ko.Spec.CacheNodeType != nil {
		res.SetCacheNodeType(*r.ko.Spec.CacheNodeType)
	}
	if r.ko.Spec.CacheParameterGroupName != nil {
		res.SetCacheParameterGroupName(*r.ko.Spec.CacheParameterGroupName)
	}
	if r.ko.Spec.CacheSecurityGroupNames != nil {
		f7 := []*string{}
		for _, f7iter := range r.ko.Spec.CacheSecurityGroupNames {
			var f7elem string
			f7elem = *f7iter
			f7 = append(f7, &f7elem)
		}
		res.SetCacheSecurityGroupNames(f7)
	}
	if r.ko.Spec.LogDeliveryConfigurations != nil {
		f8 := []*svcsdk.LogDeliveryConfigurationRequest{}
		for _, f8iter := range r.ko.Spec.LogDeliveryConfigurations {
			f8elem := &svcsdk.LogDeliveryConfigurationRequest{}
			if f8iter.DestinationDetails != nil {
				f8elemf0 := &svcsdk.DestinationDetails{}
				if f8iter.DestinationDetails.CloudWatchLogsDetails != nil {
					f8elemf0f0 := &svcsdk.CloudWatchLogsDestinationDetails{}
					if f8iter.DestinationDetails.CloudWatchLogsDetails.LogGroup != nil {
						f8elemf0f0.SetLogGroup(*f8iter.DestinationDetails.CloudWatchLogsDetails.LogGroup)
					}
					f8elemf0.SetCloudWatchLogsDetails(f8elemf0f0)
				}
				if f8iter.DestinationDetails.KinesisFirehoseDetails != nil {
					f8elemf0f1 := &svcsdk.KinesisFirehoseDestinationDetails{}
					if f8iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream != nil {
						f8elemf0f1.SetDeliveryStream(*f8iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream)
					}
					f8elemf0.SetKinesisFirehoseDetails(f8elemf0f1)
				}
				f8elem.SetDestinationDetails(f8elemf0)
			}
			if f8iter.DestinationType != nil {
				f8elem.SetDestinationType(*f8iter.DestinationType)
			}
			if f8iter.Enabled != nil {
				f8elem.SetEnabled(*f8iter.Enabled)
			}
			if f8iter.LogFormat != nil {
				f8elem.SetLogFormat(*f8iter.LogFormat)
			}
			if f8iter.LogType != nil {
				f8elem.SetLogType(*f8iter.LogType)
			}
			f8 = append(f8, f8elem)
		}
		res.SetLogDeliveryConfigurations(f8)
	}
	if r.ko.Spec.MultiAZEnabled != nil {
		res.SetMultiAZEnabled(*r.ko.Spec.MultiAZEnabled)
	}
	if r.ko.Spec.NotificationTopicARN != nil {
		res.SetNotificationTopicArn(*r.ko.Spec.NotificationTopicARN)
	}
	if r.ko.Spec.PreferredMaintenanceWindow != nil {
		res.SetPreferredMaintenanceWindow(*r.ko.Spec.PreferredMaintenanceWindow)
	}
	if r.ko.Spec.PrimaryClusterID != nil {
		res.SetPrimaryClusterId(*r.ko.Spec.PrimaryClusterID)
	}
	if r.ko.Spec.ReplicationGroupDescription != nil {
		res.SetReplicationGroupDescription(*r.ko.Spec.ReplicationGroupDescription)
	}
	if r.ko.Spec.ReplicationGroupID != nil {
		res.SetReplicationGroupId(*r.ko.Spec.ReplicationGroupID)
	}
	if r.ko.Spec.SnapshotRetentionLimit != nil {
		res.SetSnapshotRetentionLimit(*r.ko.Spec.SnapshotRetentionLimit)
	}
	if r.ko.Spec.SnapshotWindow != nil {
		res.SetSnapshotWindow(*r.ko.Spec.SnapshotWindow)
	}
	if r.ko.Status.SnapshottingClusterID != nil {
		res.SetSnapshottingClusterId(*r.ko.Status.SnapshottingClusterID)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeUpdate, "r.ko", "res", 1),
	)
}

func TestSetSDK_Elasticache_User_Create_Override_Values(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")

	crd := testutil.GetCRDByName(t, g, "User")
	require.NotNil(crd)

	expected := `
	if r.ko.Spec.AccessString != nil {
		res.SetAccessString(*r.ko.Spec.AccessString)
	}
	if r.ko.Spec.NoPasswordRequired != nil {
		res.SetNoPasswordRequired(*r.ko.Spec.NoPasswordRequired)
	}
	if r.ko.Spec.Passwords != nil {
		f3 := []*string{}
		for _, f3iter := range r.ko.Spec.Passwords {
			var f3elem string
			if f3iter != nil {
				tmpSecret, err := rm.rr.SecretValueFromReference(ctx, f3iter)
				if err != nil {
					return nil, ackrequeue.Needed(err)
				}
				if tmpSecret != "" {
					f3elem = tmpSecret
				}
			}
			f3 = append(f3, &f3elem)
		}
		res.SetPasswords(f3)
	}
	if r.ko.Spec.UserID != nil {
		res.SetUserId(*r.ko.Spec.UserID)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeUpdate, "r.ko", "res", 1),
	)
}

func TestSetSDK_RDS_DBInstance_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "rds")

	crd := testutil.GetCRDByName(t, g, "DBInstance")
	require.NotNil(crd)

	expected := `
	if r.ko.Spec.AllocatedStorage != nil {
		res.SetAllocatedStorage(*r.ko.Spec.AllocatedStorage)
	}
	if r.ko.Spec.AutoMinorVersionUpgrade != nil {
		res.SetAutoMinorVersionUpgrade(*r.ko.Spec.AutoMinorVersionUpgrade)
	}
	if r.ko.Spec.AvailabilityZone != nil {
		res.SetAvailabilityZone(*r.ko.Spec.AvailabilityZone)
	}
	if r.ko.Spec.BackupRetentionPeriod != nil {
		res.SetBackupRetentionPeriod(*r.ko.Spec.BackupRetentionPeriod)
	}
	if r.ko.Spec.CharacterSetName != nil {
		res.SetCharacterSetName(*r.ko.Spec.CharacterSetName)
	}
	if r.ko.Spec.CopyTagsToSnapshot != nil {
		res.SetCopyTagsToSnapshot(*r.ko.Spec.CopyTagsToSnapshot)
	}
	if r.ko.Spec.DBClusterIdentifier != nil {
		res.SetDBClusterIdentifier(*r.ko.Spec.DBClusterIdentifier)
	}
	if r.ko.Spec.DBInstanceClass != nil {
		res.SetDBInstanceClass(*r.ko.Spec.DBInstanceClass)
	}
	if r.ko.Spec.DBInstanceIdentifier != nil {
		res.SetDBInstanceIdentifier(*r.ko.Spec.DBInstanceIdentifier)
	}
	if r.ko.Spec.DBName != nil {
		res.SetDBName(*r.ko.Spec.DBName)
	}
	if r.ko.Spec.DBParameterGroupName != nil {
		res.SetDBParameterGroupName(*r.ko.Spec.DBParameterGroupName)
	}
	if r.ko.Spec.DBSecurityGroups != nil {
		f11 := []*string{}
		for _, f11iter := range r.ko.Spec.DBSecurityGroups {
			var f11elem string
			f11elem = *f11iter
			f11 = append(f11, &f11elem)
		}
		res.SetDBSecurityGroups(f11)
	}
	if r.ko.Spec.DBSubnetGroupName != nil {
		res.SetDBSubnetGroupName(*r.ko.Spec.DBSubnetGroupName)
	}
	if r.ko.Spec.DeletionProtection != nil {
		res.SetDeletionProtection(*r.ko.Spec.DeletionProtection)
	}
	if r.ko.Spec.Domain != nil {
		res.SetDomain(*r.ko.Spec.Domain)
	}
	if r.ko.Spec.DomainIAMRoleName != nil {
		res.SetDomainIAMRoleName(*r.ko.Spec.DomainIAMRoleName)
	}
	if r.ko.Spec.EnableCloudwatchLogsExports != nil {
		f16 := []*string{}
		for _, f16iter := range r.ko.Spec.EnableCloudwatchLogsExports {
			var f16elem string
			f16elem = *f16iter
			f16 = append(f16, &f16elem)
		}
		res.SetEnableCloudwatchLogsExports(f16)
	}
	if r.ko.Spec.EnableIAMDatabaseAuthentication != nil {
		res.SetEnableIAMDatabaseAuthentication(*r.ko.Spec.EnableIAMDatabaseAuthentication)
	}
	if r.ko.Spec.EnablePerformanceInsights != nil {
		res.SetEnablePerformanceInsights(*r.ko.Spec.EnablePerformanceInsights)
	}
	if r.ko.Spec.Engine != nil {
		res.SetEngine(*r.ko.Spec.Engine)
	}
	if r.ko.Spec.EngineVersion != nil {
		res.SetEngineVersion(*r.ko.Spec.EngineVersion)
	}
	if r.ko.Spec.IOPS != nil {
		res.SetIops(*r.ko.Spec.IOPS)
	}
	if r.ko.Spec.KMSKeyID != nil {
		res.SetKmsKeyId(*r.ko.Spec.KMSKeyID)
	}
	if r.ko.Spec.LicenseModel != nil {
		res.SetLicenseModel(*r.ko.Spec.LicenseModel)
	}
	if r.ko.Spec.MasterUserPassword != nil {
		res.SetMasterUserPassword(*r.ko.Spec.MasterUserPassword)
	}
	if r.ko.Spec.MasterUsername != nil {
		res.SetMasterUsername(*r.ko.Spec.MasterUsername)
	}
	if r.ko.Spec.MaxAllocatedStorage != nil {
		res.SetMaxAllocatedStorage(*r.ko.Spec.MaxAllocatedStorage)
	}
	if r.ko.Spec.MonitoringInterval != nil {
		res.SetMonitoringInterval(*r.ko.Spec.MonitoringInterval)
	}
	if r.ko.Spec.MonitoringRoleARN != nil {
		res.SetMonitoringRoleArn(*r.ko.Spec.MonitoringRoleARN)
	}
	if r.ko.Spec.MultiAZ != nil {
		res.SetMultiAZ(*r.ko.Spec.MultiAZ)
	}
	if r.ko.Spec.OptionGroupName != nil {
		res.SetOptionGroupName(*r.ko.Spec.OptionGroupName)
	}
	if r.ko.Spec.PerformanceInsightsKMSKeyID != nil {
		res.SetPerformanceInsightsKMSKeyId(*r.ko.Spec.PerformanceInsightsKMSKeyID)
	}
	if r.ko.Spec.PerformanceInsightsRetentionPeriod != nil {
		res.SetPerformanceInsightsRetentionPeriod(*r.ko.Spec.PerformanceInsightsRetentionPeriod)
	}
	if r.ko.Spec.Port != nil {
		res.SetPort(*r.ko.Spec.Port)
	}
	if r.ko.Spec.PreferredBackupWindow != nil {
		res.SetPreferredBackupWindow(*r.ko.Spec.PreferredBackupWindow)
	}
	if r.ko.Spec.PreferredMaintenanceWindow != nil {
		res.SetPreferredMaintenanceWindow(*r.ko.Spec.PreferredMaintenanceWindow)
	}
	if r.ko.Spec.ProcessorFeatures != nil {
		f36 := []*svcsdk.ProcessorFeature{}
		for _, f36iter := range r.ko.Spec.ProcessorFeatures {
			f36elem := &svcsdk.ProcessorFeature{}
			if f36iter.Name != nil {
				f36elem.SetName(*f36iter.Name)
			}
			if f36iter.Value != nil {
				f36elem.SetValue(*f36iter.Value)
			}
			f36 = append(f36, f36elem)
		}
		res.SetProcessorFeatures(f36)
	}
	if r.ko.Spec.PromotionTier != nil {
		res.SetPromotionTier(*r.ko.Spec.PromotionTier)
	}
	if r.ko.Spec.PubliclyAccessible != nil {
		res.SetPubliclyAccessible(*r.ko.Spec.PubliclyAccessible)
	}
	if r.ko.Spec.StorageEncrypted != nil {
		res.SetStorageEncrypted(*r.ko.Spec.StorageEncrypted)
	}
	if r.ko.Spec.StorageType != nil {
		res.SetStorageType(*r.ko.Spec.StorageType)
	}
	if r.ko.Spec.Tags != nil {
		f41 := []*svcsdk.Tag{}
		for _, f41iter := range r.ko.Spec.Tags {
			f41elem := &svcsdk.Tag{}
			if f41iter.Key != nil {
				f41elem.SetKey(*f41iter.Key)
			}
			if f41iter.Value != nil {
				f41elem.SetValue(*f41iter.Value)
			}
			f41 = append(f41, f41elem)
		}
		res.SetTags(f41)
	}
	if r.ko.Spec.TDECredentialARN != nil {
		res.SetTdeCredentialArn(*r.ko.Spec.TDECredentialARN)
	}
	if r.ko.Spec.TDECredentialPassword != nil {
		res.SetTdeCredentialPassword(*r.ko.Spec.TDECredentialPassword)
	}
	if r.ko.Spec.Timezone != nil {
		res.SetTimezone(*r.ko.Spec.Timezone)
	}
	if r.ko.Spec.VPCSecurityGroupIDs != nil {
		f45 := []*string{}
		for _, f45iter := range r.ko.Spec.VPCSecurityGroupIDs {
			var f45elem string
			f45elem = *f45iter
			f45 = append(f45, &f45elem)
		}
		res.SetVpcSecurityGroupIds(f45)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_S3_Bucket_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "s3")

	crd := testutil.GetCRDByName(t, g, "Bucket")
	require.NotNil(crd)

	expected := `
	if r.ko.Spec.ACL != nil {
		res.SetACL(*r.ko.Spec.ACL)
	}
	if r.ko.Spec.Name != nil {
		res.SetBucket(*r.ko.Spec.Name)
	}
	if r.ko.Spec.CreateBucketConfiguration != nil {
		f2 := &svcsdk.CreateBucketConfiguration{}
		if r.ko.Spec.CreateBucketConfiguration.LocationConstraint != nil {
			f2.SetLocationConstraint(*r.ko.Spec.CreateBucketConfiguration.LocationConstraint)
		}
		res.SetCreateBucketConfiguration(f2)
	}
	if r.ko.Spec.GrantFullControl != nil {
		res.SetGrantFullControl(*r.ko.Spec.GrantFullControl)
	}
	if r.ko.Spec.GrantRead != nil {
		res.SetGrantRead(*r.ko.Spec.GrantRead)
	}
	if r.ko.Spec.GrantReadACP != nil {
		res.SetGrantReadACP(*r.ko.Spec.GrantReadACP)
	}
	if r.ko.Spec.GrantWrite != nil {
		res.SetGrantWrite(*r.ko.Spec.GrantWrite)
	}
	if r.ko.Spec.GrantWriteACP != nil {
		res.SetGrantWriteACP(*r.ko.Spec.GrantWriteACP)
	}
	if r.ko.Spec.ObjectLockEnabledForBucket != nil {
		res.SetObjectLockEnabledForBucket(*r.ko.Spec.ObjectLockEnabledForBucket)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_S3_Bucket_Delete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "s3")

	crd := testutil.GetCRDByName(t, g, "Bucket")
	require.NotNil(crd)

	expected := `
	if r.ko.Spec.Name != nil {
		res.SetBucket(*r.ko.Spec.Name)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeDelete, "r.ko", "res", 1),
	)
}

func TestSetSDK_SNS_Topic_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sns")

	crd := testutil.GetCRDByName(t, g, "Topic")
	require.NotNil(crd)

	// The input shape for the Create operation is set from a variety of scalar
	// and non-scalar types and the SNS API features an Attributes parameter
	// that is actually a map[string]*string of real field values that are
	// unpacked by the code generator.
	expected := `
	attrMap := map[string]*string{}
	if r.ko.Spec.DeliveryPolicy != nil {
		attrMap["DeliveryPolicy"] = r.ko.Spec.DeliveryPolicy
	}
	if r.ko.Spec.DisplayName != nil {
		attrMap["DisplayName"] = r.ko.Spec.DisplayName
	}
	if r.ko.Spec.KMSMasterKeyID != nil {
		attrMap["KmsMasterKeyId"] = r.ko.Spec.KMSMasterKeyID
	}
	if r.ko.Spec.Policy != nil {
		attrMap["Policy"] = r.ko.Spec.Policy
	}
	res.SetAttributes(attrMap)
	if r.ko.Spec.Name != nil {
		res.SetName(*r.ko.Spec.Name)
	}
	if r.ko.Spec.Tags != nil {
		f2 := []*svcsdk.Tag{}
		for _, f2iter := range r.ko.Spec.Tags {
			f2elem := &svcsdk.Tag{}
			if f2iter.Key != nil {
				f2elem.SetKey(*f2iter.Key)
			}
			if f2iter.Value != nil {
				f2elem.SetValue(*f2iter.Value)
			}
			f2 = append(f2, f2elem)
		}
		res.SetTags(f2)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_SNS_Topic_GetAttributes(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sns")

	crd := testutil.GetCRDByName(t, g, "Topic")
	require.NotNil(crd)

	// The input shape for the GetAttributes operation has a single TopicArn
	// field. This field represents the ARN of the primary resource (the Topic
	// itself) and should be set specially from the ACKResourceMetadata.ARN
	// field in the TopicStatus struct
	expected := `
	if r.ko.Status.ACKResourceMetadata != nil && r.ko.Status.ACKResourceMetadata.ARN != nil {
		res.SetTopicArn(string(*r.ko.Status.ACKResourceMetadata.ARN))
	} else {
		res.SetTopicArn(rm.ARNFromName(*r.ko.Spec.Name))
	}
`
	assert.Equal(
		expected,
		code.SetSDKGetAttributes(crd.Config(), crd, "r.ko", "res", 1),
	)
}

func TestSetSDK_SQS_Queue_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sqs")

	crd := testutil.GetCRDByName(t, g, "Queue")
	require.NotNil(crd)

	expected := `
	attrMap := map[string]*string{}
	if r.ko.Spec.ContentBasedDeduplication != nil {
		attrMap["ContentBasedDeduplication"] = r.ko.Spec.ContentBasedDeduplication
	}
	if r.ko.Spec.DelaySeconds != nil {
		attrMap["DelaySeconds"] = r.ko.Spec.DelaySeconds
	}
	if r.ko.Spec.FifoQueue != nil {
		attrMap["FifoQueue"] = r.ko.Spec.FifoQueue
	}
	if r.ko.Spec.KMSDataKeyReusePeriodSeconds != nil {
		attrMap["KmsDataKeyReusePeriodSeconds"] = r.ko.Spec.KMSDataKeyReusePeriodSeconds
	}
	if r.ko.Spec.KMSMasterKeyID != nil {
		attrMap["KmsMasterKeyId"] = r.ko.Spec.KMSMasterKeyID
	}
	if r.ko.Spec.MaximumMessageSize != nil {
		attrMap["MaximumMessageSize"] = r.ko.Spec.MaximumMessageSize
	}
	if r.ko.Spec.MessageRetentionPeriod != nil {
		attrMap["MessageRetentionPeriod"] = r.ko.Spec.MessageRetentionPeriod
	}
	if r.ko.Spec.Policy != nil {
		attrMap["Policy"] = r.ko.Spec.Policy
	}
	if r.ko.Spec.ReceiveMessageWaitTimeSeconds != nil {
		attrMap["ReceiveMessageWaitTimeSeconds"] = r.ko.Spec.ReceiveMessageWaitTimeSeconds
	}
	if r.ko.Spec.RedrivePolicy != nil {
		attrMap["RedrivePolicy"] = r.ko.Spec.RedrivePolicy
	}
	if r.ko.Spec.VisibilityTimeout != nil {
		attrMap["VisibilityTimeout"] = r.ko.Spec.VisibilityTimeout
	}
	res.SetAttributes(attrMap)
	if r.ko.Spec.QueueName != nil {
		res.SetQueueName(*r.ko.Spec.QueueName)
	}
	if r.ko.Spec.Tags != nil {
		f2 := map[string]*string{}
		for f2key, f2valiter := range r.ko.Spec.Tags {
			var f2val string
			f2val = *f2valiter
			f2[f2key] = &f2val
		}
		res.SetTags(f2)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_SQS_Queue_GetAttributes(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sqs")

	crd := testutil.GetCRDByName(t, g, "Queue")
	require.NotNil(crd)

	// The input shape for the GetAttributes operation technically has two
	// fields in it: an AttributeNames list of attribute keys to file
	// attributes for and a QueueUrl field. We only care about the QueueUrl
	// field, since we look for all attributes for a queue.
	expected := `
	{
		tmpVals := []*string{}
		tmpVal0 := "All"
		tmpVals = append(tmpVals, &tmpVal0)
		res.SetAttributeNames(tmpVals)
	}
	if r.ko.Status.QueueURL != nil {
		res.SetQueueUrl(*r.ko.Status.QueueURL)
	}
`
	assert.Equal(
		expected,
		code.SetSDKGetAttributes(crd.Config(), crd, "r.ko", "res", 1),
	)
}

func TestSetSDK_MQ_Broker_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "mq")

	crd := testutil.GetCRDByName(t, g, "Broker")
	require.NotNil(crd)

	expected := `
	if r.ko.Spec.AuthenticationStrategy != nil {
		res.SetAuthenticationStrategy(*r.ko.Spec.AuthenticationStrategy)
	}
	if r.ko.Spec.AutoMinorVersionUpgrade != nil {
		res.SetAutoMinorVersionUpgrade(*r.ko.Spec.AutoMinorVersionUpgrade)
	}
	if r.ko.Spec.BrokerName != nil {
		res.SetBrokerName(*r.ko.Spec.BrokerName)
	}
	if r.ko.Spec.Configuration != nil {
		f3 := &svcsdk.ConfigurationId{}
		if r.ko.Spec.Configuration.ID != nil {
			f3.SetId(*r.ko.Spec.Configuration.ID)
		}
		if r.ko.Spec.Configuration.Revision != nil {
			f3.SetRevision(*r.ko.Spec.Configuration.Revision)
		}
		res.SetConfiguration(f3)
	}
	if r.ko.Spec.CreatorRequestID != nil {
		res.SetCreatorRequestId(*r.ko.Spec.CreatorRequestID)
	}
	if r.ko.Spec.DeploymentMode != nil {
		res.SetDeploymentMode(*r.ko.Spec.DeploymentMode)
	}
	if r.ko.Spec.EncryptionOptions != nil {
		f6 := &svcsdk.EncryptionOptions{}
		if r.ko.Spec.EncryptionOptions.KMSKeyID != nil {
			f6.SetKmsKeyId(*r.ko.Spec.EncryptionOptions.KMSKeyID)
		}
		if r.ko.Spec.EncryptionOptions.UseAWSOwnedKey != nil {
			f6.SetUseAwsOwnedKey(*r.ko.Spec.EncryptionOptions.UseAWSOwnedKey)
		}
		res.SetEncryptionOptions(f6)
	}
	if r.ko.Spec.EngineType != nil {
		res.SetEngineType(*r.ko.Spec.EngineType)
	}
	if r.ko.Spec.EngineVersion != nil {
		res.SetEngineVersion(*r.ko.Spec.EngineVersion)
	}
	if r.ko.Spec.HostInstanceType != nil {
		res.SetHostInstanceType(*r.ko.Spec.HostInstanceType)
	}
	if r.ko.Spec.LDAPServerMetadata != nil {
		f10 := &svcsdk.LdapServerMetadataInput{}
		if r.ko.Spec.LDAPServerMetadata.Hosts != nil {
			f10f0 := []*string{}
			for _, f10f0iter := range r.ko.Spec.LDAPServerMetadata.Hosts {
				var f10f0elem string
				f10f0elem = *f10f0iter
				f10f0 = append(f10f0, &f10f0elem)
			}
			f10.SetHosts(f10f0)
		}
		if r.ko.Spec.LDAPServerMetadata.RoleBase != nil {
			f10.SetRoleBase(*r.ko.Spec.LDAPServerMetadata.RoleBase)
		}
		if r.ko.Spec.LDAPServerMetadata.RoleName != nil {
			f10.SetRoleName(*r.ko.Spec.LDAPServerMetadata.RoleName)
		}
		if r.ko.Spec.LDAPServerMetadata.RoleSearchMatching != nil {
			f10.SetRoleSearchMatching(*r.ko.Spec.LDAPServerMetadata.RoleSearchMatching)
		}
		if r.ko.Spec.LDAPServerMetadata.RoleSearchSubtree != nil {
			f10.SetRoleSearchSubtree(*r.ko.Spec.LDAPServerMetadata.RoleSearchSubtree)
		}
		if r.ko.Spec.LDAPServerMetadata.ServiceAccountPassword != nil {
			f10.SetServiceAccountPassword(*r.ko.Spec.LDAPServerMetadata.ServiceAccountPassword)
		}
		if r.ko.Spec.LDAPServerMetadata.ServiceAccountUsername != nil {
			f10.SetServiceAccountUsername(*r.ko.Spec.LDAPServerMetadata.ServiceAccountUsername)
		}
		if r.ko.Spec.LDAPServerMetadata.UserBase != nil {
			f10.SetUserBase(*r.ko.Spec.LDAPServerMetadata.UserBase)
		}
		if r.ko.Spec.LDAPServerMetadata.UserRoleName != nil {
			f10.SetUserRoleName(*r.ko.Spec.LDAPServerMetadata.UserRoleName)
		}
		if r.ko.Spec.LDAPServerMetadata.UserSearchMatching != nil {
			f10.SetUserSearchMatching(*r.ko.Spec.LDAPServerMetadata.UserSearchMatching)
		}
		if r.ko.Spec.LDAPServerMetadata.UserSearchSubtree != nil {
			f10.SetUserSearchSubtree(*r.ko.Spec.LDAPServerMetadata.UserSearchSubtree)
		}
		res.SetLdapServerMetadata(f10)
	}
	if r.ko.Spec.Logs != nil {
		f11 := &svcsdk.Logs{}
		if r.ko.Spec.Logs.Audit != nil {
			f11.SetAudit(*r.ko.Spec.Logs.Audit)
		}
		if r.ko.Spec.Logs.General != nil {
			f11.SetGeneral(*r.ko.Spec.Logs.General)
		}
		res.SetLogs(f11)
	}
	if r.ko.Spec.MaintenanceWindowStartTime != nil {
		f12 := &svcsdk.WeeklyStartTime{}
		if r.ko.Spec.MaintenanceWindowStartTime.DayOfWeek != nil {
			f12.SetDayOfWeek(*r.ko.Spec.MaintenanceWindowStartTime.DayOfWeek)
		}
		if r.ko.Spec.MaintenanceWindowStartTime.TimeOfDay != nil {
			f12.SetTimeOfDay(*r.ko.Spec.MaintenanceWindowStartTime.TimeOfDay)
		}
		if r.ko.Spec.MaintenanceWindowStartTime.TimeZone != nil {
			f12.SetTimeZone(*r.ko.Spec.MaintenanceWindowStartTime.TimeZone)
		}
		res.SetMaintenanceWindowStartTime(f12)
	}
	if r.ko.Spec.PubliclyAccessible != nil {
		res.SetPubliclyAccessible(*r.ko.Spec.PubliclyAccessible)
	}
	if r.ko.Spec.SecurityGroups != nil {
		f14 := []*string{}
		for _, f14iter := range r.ko.Spec.SecurityGroups {
			var f14elem string
			f14elem = *f14iter
			f14 = append(f14, &f14elem)
		}
		res.SetSecurityGroups(f14)
	}
	if r.ko.Spec.StorageType != nil {
		res.SetStorageType(*r.ko.Spec.StorageType)
	}
	if r.ko.Spec.SubnetIDs != nil {
		f16 := []*string{}
		for _, f16iter := range r.ko.Spec.SubnetIDs {
			var f16elem string
			f16elem = *f16iter
			f16 = append(f16, &f16elem)
		}
		res.SetSubnetIds(f16)
	}
	if r.ko.Spec.Tags != nil {
		f17 := map[string]*string{}
		for f17key, f17valiter := range r.ko.Spec.Tags {
			var f17val string
			f17val = *f17valiter
			f17[f17key] = &f17val
		}
		res.SetTags(f17)
	}
	if r.ko.Spec.Users != nil {
		f18 := []*svcsdk.User{}
		for _, f18iter := range r.ko.Spec.Users {
			f18elem := &svcsdk.User{}
			if f18iter.ConsoleAccess != nil {
				f18elem.SetConsoleAccess(*f18iter.ConsoleAccess)
			}
			if f18iter.Groups != nil {
				f18elemf1 := []*string{}
				for _, f18elemf1iter := range f18iter.Groups {
					var f18elemf1elem string
					f18elemf1elem = *f18elemf1iter
					f18elemf1 = append(f18elemf1, &f18elemf1elem)
				}
				f18elem.SetGroups(f18elemf1)
			}
			if f18iter.Password != nil {
				tmpSecret, err := rm.rr.SecretValueFromReference(ctx, f18iter.Password)
				if err != nil {
					return nil, ackrequeue.Needed(err)
				}
				if tmpSecret != "" {
					f18elem.SetPassword(tmpSecret)
				}
			}
			if f18iter.Username != nil {
				f18elem.SetUsername(*f18iter.Username)
			}
			f18 = append(f18, f18elem)
		}
		res.SetUsers(f18)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_EC2_VPC_ReadMany(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "Vpc")
	require.NotNil(crd)

	expected := `
	if r.ko.Status.VPCID != nil {
		f4 := []*string{}
		f4 = append(f4, r.ko.Status.VPCID)
		res.SetVpcIds(f4)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeList, "r.ko", "res", 1),
	)
}

func Test_SetSDK_ECR_Repository_newListRequestPayload(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-renamed-identifier-field.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)

	expected := `
	if r.ko.Status.RegistryID != nil {
		res.SetRegistryId(*r.ko.Status.RegistryID)
	}
	if r.ko.Spec.Name != nil {
		f3 := []*string{}
		f3 = append(f3, r.ko.Spec.Name)
		res.SetRepositoryNames(f3)
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeList, "r.ko", "res", 1),
	)
}
