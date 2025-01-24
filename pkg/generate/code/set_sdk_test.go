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
		res.ApiId = r.ko.Spec.APIID
	}
	if r.ko.Spec.APIKeyRequired != nil {
		res.ApiKeyRequired = r.ko.Spec.APIKeyRequired
	}
	if r.ko.Spec.AuthorizationScopes != nil {
		res.AuthorizationScopes = aws.ToStringSlice(r.ko.Spec.AuthorizationScopes)
	}
	if r.ko.Spec.AuthorizationType != nil {
		res.AuthorizationType = svcsdktypes.AuthorizationType(*r.ko.Spec.AuthorizationType)
	}
	if r.ko.Spec.AuthorizerID != nil {
		res.AuthorizerId = r.ko.Spec.AuthorizerID
	}
	if r.ko.Spec.ModelSelectionExpression != nil {
		res.ModelSelectionExpression = r.ko.Spec.ModelSelectionExpression
	}
	if r.ko.Spec.OperationName != nil {
		res.OperationName = r.ko.Spec.OperationName
	}
	if r.ko.Spec.RequestModels != nil {
		res.RequestModels = aws.ToStringMap(r.ko.Spec.RequestModels)
	}
	if r.ko.Spec.RequestParameters != nil {
		f8 := map[string]svcsdktypes.ParameterConstraints{}
		for f8key, f8valiter := range r.ko.Spec.RequestParameters {
			f8val := &svcsdktypes.ParameterConstraints{}
			if f8valiter.Required != nil {
				f8val.Required = f8valiter.Required
			}
			f8[f8key] = *f8val
		}
		res.RequestParameters = f8
	}
	if r.ko.Spec.RouteKey != nil {
		res.RouteKey = r.ko.Spec.RouteKey
	}
	if r.ko.Spec.RouteResponseSelectionExpression != nil {
		res.RouteResponseSelectionExpression = r.ko.Spec.RouteResponseSelectionExpression
	}
	if r.ko.Spec.Target != nil {
		res.Target = r.ko.Spec.Target
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

func TestSetSDK_OpenSearch_Domain_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "opensearch")

	crd := testutil.GetCRDByName(t, g, "Domain")
	require.NotNil(crd)

	expected := `
	if r.ko.Spec.AIMLOptions != nil {
		f0 := &svcsdktypes.AIMLOptionsInput{}
		if r.ko.Spec.AIMLOptions.NATuralLanguageQueryGenerationOptions != nil {
			f0f0 := &svcsdktypes.NaturalLanguageQueryGenerationOptionsInput{}
			if r.ko.Spec.AIMLOptions.NATuralLanguageQueryGenerationOptions.DesiredState != nil {
				f0f0.DesiredState = svcsdktypes.NaturalLanguageQueryGenerationDesiredState(*r.ko.Spec.AIMLOptions.NATuralLanguageQueryGenerationOptions.DesiredState)
			}
			f0.NaturalLanguageQueryGenerationOptions = f0f0
		}
		res.AIMLOptions = f0
	}
	if r.ko.Spec.AccessPolicies != nil {
		res.AccessPolicies = r.ko.Spec.AccessPolicies
	}
	if r.ko.Spec.AdvancedOptions != nil {
		res.AdvancedOptions = aws.ToStringMap(r.ko.Spec.AdvancedOptions)
	}
	if r.ko.Spec.AdvancedSecurityOptions != nil {
		f3 := &svcsdktypes.AdvancedSecurityOptionsInput{}
		if r.ko.Spec.AdvancedSecurityOptions.AnonymousAuthEnabled != nil {
			f3.AnonymousAuthEnabled = r.ko.Spec.AdvancedSecurityOptions.AnonymousAuthEnabled
		}
		if r.ko.Spec.AdvancedSecurityOptions.Enabled != nil {
			f3.Enabled = r.ko.Spec.AdvancedSecurityOptions.Enabled
		}
		if r.ko.Spec.AdvancedSecurityOptions.InternalUserDatabaseEnabled != nil {
			f3.InternalUserDatabaseEnabled = r.ko.Spec.AdvancedSecurityOptions.InternalUserDatabaseEnabled
		}
		if r.ko.Spec.AdvancedSecurityOptions.JWTOptions != nil {
			f3f3 := &svcsdktypes.JWTOptionsInput{}
			if r.ko.Spec.AdvancedSecurityOptions.JWTOptions.Enabled != nil {
				f3f3.Enabled = r.ko.Spec.AdvancedSecurityOptions.JWTOptions.Enabled
			}
			if r.ko.Spec.AdvancedSecurityOptions.JWTOptions.PublicKey != nil {
				f3f3.PublicKey = r.ko.Spec.AdvancedSecurityOptions.JWTOptions.PublicKey
			}
			if r.ko.Spec.AdvancedSecurityOptions.JWTOptions.RolesKey != nil {
				f3f3.RolesKey = r.ko.Spec.AdvancedSecurityOptions.JWTOptions.RolesKey
			}
			if r.ko.Spec.AdvancedSecurityOptions.JWTOptions.SubjectKey != nil {
				f3f3.SubjectKey = r.ko.Spec.AdvancedSecurityOptions.JWTOptions.SubjectKey
			}
			f3.JWTOptions = f3f3
		}
		if r.ko.Spec.AdvancedSecurityOptions.MasterUserOptions != nil {
			f3f4 := &svcsdktypes.MasterUserOptions{}
			if r.ko.Spec.AdvancedSecurityOptions.MasterUserOptions.MasterUserARN != nil {
				f3f4.MasterUserARN = r.ko.Spec.AdvancedSecurityOptions.MasterUserOptions.MasterUserARN
			}
			if r.ko.Spec.AdvancedSecurityOptions.MasterUserOptions.MasterUserName != nil {
				f3f4.MasterUserName = r.ko.Spec.AdvancedSecurityOptions.MasterUserOptions.MasterUserName
			}
			if r.ko.Spec.AdvancedSecurityOptions.MasterUserOptions.MasterUserPassword != nil {
				tmpSecret, err := rm.rr.SecretValueFromReference(ctx, r.ko.Spec.AdvancedSecurityOptions.MasterUserOptions.MasterUserPassword)
				if err != nil {
					return nil, ackrequeue.Needed(err)
				}
				if tmpSecret != "" {
					f3f4.MasterUserPassword = aws.String(tmpSecret)
				}
			}
			f3.MasterUserOptions = f3f4
		}
		if r.ko.Spec.AdvancedSecurityOptions.SAMLOptions != nil {
			f3f5 := &svcsdktypes.SAMLOptionsInput{}
			if r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.Enabled != nil {
				f3f5.Enabled = r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.Enabled
			}
			if r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.IDp != nil {
				f3f5f1 := &svcsdktypes.SAMLIdp{}
				if r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.IDp.EntityID != nil {
					f3f5f1.EntityId = r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.IDp.EntityID
				}
				if r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.IDp.MetadataContent != nil {
					f3f5f1.MetadataContent = r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.IDp.MetadataContent
				}
				f3f5.Idp = f3f5f1
			}
			if r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.MasterBackendRole != nil {
				f3f5.MasterBackendRole = r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.MasterBackendRole
			}
			if r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.MasterUserName != nil {
				f3f5.MasterUserName = r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.MasterUserName
			}
			if r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.RolesKey != nil {
				f3f5.RolesKey = r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.RolesKey
			}
			if r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.SessionTimeoutMinutes != nil {
				if *r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.SessionTimeoutMinutes > math.MaxInt32 || *r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.SessionTimeoutMinutes < math.MinInt32 {
					return nil, fmt.Errorf("error: field SessionTimeoutMinutes is of type int32")
				}
				sessionTimeoutMinutesCopy := int32(*r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.SessionTimeoutMinutes)
				f3f5.SessionTimeoutMinutes = &sessionTimeoutMinutesCopy
			}
			if r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.SubjectKey != nil {
				f3f5.SubjectKey = r.ko.Spec.AdvancedSecurityOptions.SAMLOptions.SubjectKey
			}
			f3.SAMLOptions = f3f5
		}
		res.AdvancedSecurityOptions = f3
	}
	if r.ko.Spec.AutoTuneOptions != nil {
		f4 := &svcsdktypes.AutoTuneOptionsInput{}
		if r.ko.Spec.AutoTuneOptions.DesiredState != nil {
			f4.DesiredState = svcsdktypes.AutoTuneDesiredState(*r.ko.Spec.AutoTuneOptions.DesiredState)
		}
		if r.ko.Spec.AutoTuneOptions.MaintenanceSchedules != nil {
			f4f1 := []svcsdktypes.AutoTuneMaintenanceSchedule{}
			for _, f4f1iter := range r.ko.Spec.AutoTuneOptions.MaintenanceSchedules {
				f4f1elem := &svcsdktypes.AutoTuneMaintenanceSchedule{}
				if f4f1iter.CronExpressionForRecurrence != nil {
					f4f1elem.CronExpressionForRecurrence = f4f1iter.CronExpressionForRecurrence
				}
				if f4f1iter.Duration != nil {
					f4f1elemf1 := &svcsdktypes.Duration{}
					if f4f1iter.Duration.Unit != nil {
						f4f1elemf1.Unit = svcsdktypes.TimeUnit(*f4f1iter.Duration.Unit)
					}
					if f4f1iter.Duration.Value != nil {
						f4f1elemf1.Value = f4f1iter.Duration.Value
					}
					f4f1elem.Duration = f4f1elemf1
				}
				if f4f1iter.StartAt != nil {
					f4f1elem.StartAt = &f4f1iter.StartAt.Time
				}
				f4f1 = append(f4f1, *f4f1elem)
			}
			f4.MaintenanceSchedules = f4f1
		}
		if r.ko.Spec.AutoTuneOptions.UseOffPeakWindow != nil {
			f4.UseOffPeakWindow = r.ko.Spec.AutoTuneOptions.UseOffPeakWindow
		}
		res.AutoTuneOptions = f4
	}
	if r.ko.Spec.ClusterConfig != nil {
		f5 := &svcsdktypes.ClusterConfig{}
		if r.ko.Spec.ClusterConfig.ColdStorageOptions != nil {
			f5f0 := &svcsdktypes.ColdStorageOptions{}
			if r.ko.Spec.ClusterConfig.ColdStorageOptions.Enabled != nil {
				f5f0.Enabled = r.ko.Spec.ClusterConfig.ColdStorageOptions.Enabled
			}
			f5.ColdStorageOptions = f5f0
		}
		if r.ko.Spec.ClusterConfig.DedicatedMasterCount != nil {
			if *r.ko.Spec.ClusterConfig.DedicatedMasterCount > math.MaxInt32 || *r.ko.Spec.ClusterConfig.DedicatedMasterCount < math.MinInt32 {
				return nil, fmt.Errorf("error: field DedicatedMasterCount is of type int32")
			}
			dedicatedMasterCountCopy := int32(*r.ko.Spec.ClusterConfig.DedicatedMasterCount)
			f5.DedicatedMasterCount = &dedicatedMasterCountCopy
		}
		if r.ko.Spec.ClusterConfig.DedicatedMasterEnabled != nil {
			f5.DedicatedMasterEnabled = r.ko.Spec.ClusterConfig.DedicatedMasterEnabled
		}
		if r.ko.Spec.ClusterConfig.DedicatedMasterType != nil {
			f5.DedicatedMasterType = svcsdktypes.OpenSearchPartitionInstanceType(*r.ko.Spec.ClusterConfig.DedicatedMasterType)
		}
		if r.ko.Spec.ClusterConfig.InstanceCount != nil {
			if *r.ko.Spec.ClusterConfig.InstanceCount > math.MaxInt32 || *r.ko.Spec.ClusterConfig.InstanceCount < math.MinInt32 {
				return nil, fmt.Errorf("error: field InstanceCount is of type int32")
			}
			instanceCountCopy := int32(*r.ko.Spec.ClusterConfig.InstanceCount)
			f5.InstanceCount = &instanceCountCopy
		}
		if r.ko.Spec.ClusterConfig.InstanceType != nil {
			f5.InstanceType = svcsdktypes.OpenSearchPartitionInstanceType(*r.ko.Spec.ClusterConfig.InstanceType)
		}
		if r.ko.Spec.ClusterConfig.MultiAZWithStandbyEnabled != nil {
			f5.MultiAZWithStandbyEnabled = r.ko.Spec.ClusterConfig.MultiAZWithStandbyEnabled
		}
		if r.ko.Spec.ClusterConfig.NodeOptions != nil {
			f5f7 := []svcsdktypes.NodeOption{}
			for _, f5f7iter := range r.ko.Spec.ClusterConfig.NodeOptions {
				f5f7elem := &svcsdktypes.NodeOption{}
				if f5f7iter.NodeConfig != nil {
					f5f7elemf0 := &svcsdktypes.NodeConfig{}
					if f5f7iter.NodeConfig.Count != nil {
						if *f5f7iter.NodeConfig.Count > math.MaxInt32 || *f5f7iter.NodeConfig.Count < math.MinInt32 {
							return nil, fmt.Errorf("error: field Count is of type int32")
						}
						countCopy := int32(*f5f7iter.NodeConfig.Count)
						f5f7elemf0.Count = &countCopy
					}
					if f5f7iter.NodeConfig.Enabled != nil {
						f5f7elemf0.Enabled = f5f7iter.NodeConfig.Enabled
					}
					if f5f7iter.NodeConfig.Type != nil {
						f5f7elemf0.Type = svcsdktypes.OpenSearchPartitionInstanceType(*f5f7iter.NodeConfig.Type)
					}
					f5f7elem.NodeConfig = f5f7elemf0
				}
				if f5f7iter.NodeType != nil {
					f5f7elem.NodeType = svcsdktypes.NodeOptionsNodeType(*f5f7iter.NodeType)
				}
				f5f7 = append(f5f7, *f5f7elem)
			}
			f5.NodeOptions = f5f7
		}
		if r.ko.Spec.ClusterConfig.WarmCount != nil {
			if *r.ko.Spec.ClusterConfig.WarmCount > math.MaxInt32 || *r.ko.Spec.ClusterConfig.WarmCount < math.MinInt32 {
				return nil, fmt.Errorf("error: field WarmCount is of type int32")
			}
			warmCountCopy := int32(*r.ko.Spec.ClusterConfig.WarmCount)
			f5.WarmCount = &warmCountCopy
		}
		if r.ko.Spec.ClusterConfig.WarmEnabled != nil {
			f5.WarmEnabled = r.ko.Spec.ClusterConfig.WarmEnabled
		}
		if r.ko.Spec.ClusterConfig.WarmType != nil {
			f5.WarmType = svcsdktypes.OpenSearchWarmPartitionInstanceType(*r.ko.Spec.ClusterConfig.WarmType)
		}
		if r.ko.Spec.ClusterConfig.ZoneAwarenessConfig != nil {
			f5f11 := &svcsdktypes.ZoneAwarenessConfig{}
			if r.ko.Spec.ClusterConfig.ZoneAwarenessConfig.AvailabilityZoneCount != nil {
				if *r.ko.Spec.ClusterConfig.ZoneAwarenessConfig.AvailabilityZoneCount > math.MaxInt32 || *r.ko.Spec.ClusterConfig.ZoneAwarenessConfig.AvailabilityZoneCount < math.MinInt32 {
					return nil, fmt.Errorf("error: field AvailabilityZoneCount is of type int32")
				}
				availabilityZoneCountCopy := int32(*r.ko.Spec.ClusterConfig.ZoneAwarenessConfig.AvailabilityZoneCount)
				f5f11.AvailabilityZoneCount = &availabilityZoneCountCopy
			}
			f5.ZoneAwarenessConfig = f5f11
		}
		if r.ko.Spec.ClusterConfig.ZoneAwarenessEnabled != nil {
			f5.ZoneAwarenessEnabled = r.ko.Spec.ClusterConfig.ZoneAwarenessEnabled
		}
		res.ClusterConfig = f5
	}
	if r.ko.Spec.CognitoOptions != nil {
		f6 := &svcsdktypes.CognitoOptions{}
		if r.ko.Spec.CognitoOptions.Enabled != nil {
			f6.Enabled = r.ko.Spec.CognitoOptions.Enabled
		}
		if r.ko.Spec.CognitoOptions.IdentityPoolID != nil {
			f6.IdentityPoolId = r.ko.Spec.CognitoOptions.IdentityPoolID
		}
		if r.ko.Spec.CognitoOptions.RoleARN != nil {
			f6.RoleArn = r.ko.Spec.CognitoOptions.RoleARN
		}
		if r.ko.Spec.CognitoOptions.UserPoolID != nil {
			f6.UserPoolId = r.ko.Spec.CognitoOptions.UserPoolID
		}
		res.CognitoOptions = f6
	}
	if r.ko.Spec.DomainEndpointOptions != nil {
		f7 := &svcsdktypes.DomainEndpointOptions{}
		if r.ko.Spec.DomainEndpointOptions.CustomEndpoint != nil {
			f7.CustomEndpoint = r.ko.Spec.DomainEndpointOptions.CustomEndpoint
		}
		if r.ko.Spec.DomainEndpointOptions.CustomEndpointCertificateARN != nil {
			f7.CustomEndpointCertificateArn = r.ko.Spec.DomainEndpointOptions.CustomEndpointCertificateARN
		}
		if r.ko.Spec.DomainEndpointOptions.CustomEndpointEnabled != nil {
			f7.CustomEndpointEnabled = r.ko.Spec.DomainEndpointOptions.CustomEndpointEnabled
		}
		if r.ko.Spec.DomainEndpointOptions.EnforceHTTPS != nil {
			f7.EnforceHTTPS = r.ko.Spec.DomainEndpointOptions.EnforceHTTPS
		}
		if r.ko.Spec.DomainEndpointOptions.TLSSecurityPolicy != nil {
			f7.TLSSecurityPolicy = svcsdktypes.TLSSecurityPolicy(*r.ko.Spec.DomainEndpointOptions.TLSSecurityPolicy)
		}
		res.DomainEndpointOptions = f7
	}
	if r.ko.Spec.Name != nil {
		res.DomainName = r.ko.Spec.Name
	}
	if r.ko.Spec.EBSOptions != nil {
		f9 := &svcsdktypes.EBSOptions{}
		if r.ko.Spec.EBSOptions.EBSEnabled != nil {
			f9.EBSEnabled = r.ko.Spec.EBSOptions.EBSEnabled
		}
		if r.ko.Spec.EBSOptions.IOPS != nil {
			if *r.ko.Spec.EBSOptions.IOPS > math.MaxInt32 || *r.ko.Spec.EBSOptions.IOPS < math.MinInt32 {
				return nil, fmt.Errorf("error: field Iops is of type int32")
			}
			iopsCopy := int32(*r.ko.Spec.EBSOptions.IOPS)
			f9.Iops = &iopsCopy
		}
		if r.ko.Spec.EBSOptions.Throughput != nil {
			if *r.ko.Spec.EBSOptions.Throughput > math.MaxInt32 || *r.ko.Spec.EBSOptions.Throughput < math.MinInt32 {
				return nil, fmt.Errorf("error: field Throughput is of type int32")
			}
			throughputCopy := int32(*r.ko.Spec.EBSOptions.Throughput)
			f9.Throughput = &throughputCopy
		}
		if r.ko.Spec.EBSOptions.VolumeSize != nil {
			if *r.ko.Spec.EBSOptions.VolumeSize > math.MaxInt32 || *r.ko.Spec.EBSOptions.VolumeSize < math.MinInt32 {
				return nil, fmt.Errorf("error: field VolumeSize is of type int32")
			}
			volumeSizeCopy := int32(*r.ko.Spec.EBSOptions.VolumeSize)
			f9.VolumeSize = &volumeSizeCopy
		}
		if r.ko.Spec.EBSOptions.VolumeType != nil {
			f9.VolumeType = svcsdktypes.VolumeType(*r.ko.Spec.EBSOptions.VolumeType)
		}
		res.EBSOptions = f9
	}
	if r.ko.Spec.EncryptionAtRestOptions != nil {
		f10 := &svcsdktypes.EncryptionAtRestOptions{}
		if r.ko.Spec.EncryptionAtRestOptions.Enabled != nil {
			f10.Enabled = r.ko.Spec.EncryptionAtRestOptions.Enabled
		}
		if r.ko.Spec.EncryptionAtRestOptions.KMSKeyID != nil {
			f10.KmsKeyId = r.ko.Spec.EncryptionAtRestOptions.KMSKeyID
		}
		res.EncryptionAtRestOptions = f10
	}
	if r.ko.Spec.EngineVersion != nil {
		res.EngineVersion = r.ko.Spec.EngineVersion
	}
	if r.ko.Spec.IPAddressType != nil {
		res.IPAddressType = svcsdktypes.IPAddressType(*r.ko.Spec.IPAddressType)
	}
	if r.ko.Spec.IdentityCenterOptions != nil {
		f13 := &svcsdktypes.IdentityCenterOptionsInput{}
		if r.ko.Spec.IdentityCenterOptions.EnabledAPIAccess != nil {
			f13.EnabledAPIAccess = r.ko.Spec.IdentityCenterOptions.EnabledAPIAccess
		}
		if r.ko.Spec.IdentityCenterOptions.IdentityCenterInstanceARN != nil {
			f13.IdentityCenterInstanceARN = r.ko.Spec.IdentityCenterOptions.IdentityCenterInstanceARN
		}
		if r.ko.Spec.IdentityCenterOptions.RolesKey != nil {
			f13.RolesKey = svcsdktypes.RolesKeyIdCOption(*r.ko.Spec.IdentityCenterOptions.RolesKey)
		}
		if r.ko.Spec.IdentityCenterOptions.SubjectKey != nil {
			f13.SubjectKey = svcsdktypes.SubjectKeyIdCOption(*r.ko.Spec.IdentityCenterOptions.SubjectKey)
		}
		res.IdentityCenterOptions = f13
	}
	if r.ko.Spec.LogPublishingOptions != nil {
		f14 := map[string]svcsdktypes.LogPublishingOption{}
		for f14key, f14valiter := range r.ko.Spec.LogPublishingOptions {
			f14val := &svcsdktypes.LogPublishingOption{}
			if f14valiter.CloudWatchLogsLogGroupARN != nil {
				f14val.CloudWatchLogsLogGroupArn = f14valiter.CloudWatchLogsLogGroupARN
			}
			if f14valiter.Enabled != nil {
				f14val.Enabled = f14valiter.Enabled
			}
			f14[f14key] = *f14val
		}
		res.LogPublishingOptions = f14
	}
	if r.ko.Spec.NodeToNodeEncryptionOptions != nil {
		f15 := &svcsdktypes.NodeToNodeEncryptionOptions{}
		if r.ko.Spec.NodeToNodeEncryptionOptions.Enabled != nil {
			f15.Enabled = r.ko.Spec.NodeToNodeEncryptionOptions.Enabled
		}
		res.NodeToNodeEncryptionOptions = f15
	}
	if r.ko.Spec.OffPeakWindowOptions != nil {
		f16 := &svcsdktypes.OffPeakWindowOptions{}
		if r.ko.Spec.OffPeakWindowOptions.Enabled != nil {
			f16.Enabled = r.ko.Spec.OffPeakWindowOptions.Enabled
		}
		if r.ko.Spec.OffPeakWindowOptions.OffPeakWindow != nil {
			f16f1 := &svcsdktypes.OffPeakWindow{}
			if r.ko.Spec.OffPeakWindowOptions.OffPeakWindow.WindowStartTime != nil {
				f16f1f0 := &svcsdktypes.WindowStartTime{}
				if r.ko.Spec.OffPeakWindowOptions.OffPeakWindow.WindowStartTime.Hours != nil {
					f16f1f0.Hours = *r.ko.Spec.OffPeakWindowOptions.OffPeakWindow.WindowStartTime.Hours
				}
				if r.ko.Spec.OffPeakWindowOptions.OffPeakWindow.WindowStartTime.Minutes != nil {
					f16f1f0.Minutes = *r.ko.Spec.OffPeakWindowOptions.OffPeakWindow.WindowStartTime.Minutes
				}
				f16f1.WindowStartTime = f16f1f0
			}
			f16.OffPeakWindow = f16f1
		}
		res.OffPeakWindowOptions = f16
	}
	if r.ko.Spec.SoftwareUpdateOptions != nil {
		f17 := &svcsdktypes.SoftwareUpdateOptions{}
		if r.ko.Spec.SoftwareUpdateOptions.AutoSoftwareUpdateEnabled != nil {
			f17.AutoSoftwareUpdateEnabled = r.ko.Spec.SoftwareUpdateOptions.AutoSoftwareUpdateEnabled
		}
		res.SoftwareUpdateOptions = f17
	}
	if r.ko.Spec.Tags != nil {
		f18 := []svcsdktypes.Tag{}
		for _, f18iter := range r.ko.Spec.Tags {
			f18elem := &svcsdktypes.Tag{}
			if f18iter.Key != nil {
				f18elem.Key = f18iter.Key
			}
			if f18iter.Value != nil {
				f18elem.Value = f18iter.Value
			}
			f18 = append(f18, *f18elem)
		}
		res.TagList = f18
	}
	if r.ko.Spec.VPCOptions != nil {
		f19 := &svcsdktypes.VPCOptions{}
		if r.ko.Spec.VPCOptions.SecurityGroupIDs != nil {
			f19.SecurityGroupIds = aws.ToStringSlice(r.ko.Spec.VPCOptions.SecurityGroupIDs)
		}
		if r.ko.Spec.VPCOptions.SubnetIDs != nil {
			f19.SubnetIds = aws.ToStringSlice(r.ko.Spec.VPCOptions.SubnetIDs)
		}
		res.VPCOptions = f19
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
		f0 := []svcsdktypes.AttributeDefinition{}
		for _, f0iter := range r.ko.Spec.AttributeDefinitions {
			f0elem := &svcsdktypes.AttributeDefinition{}
			if f0iter.AttributeName != nil {
				f0elem.AttributeName = f0iter.AttributeName
			}
			if f0iter.AttributeType != nil {
				f0elem.AttributeType = svcsdktypes.ScalarAttributeType(*f0iter.AttributeType)
			}
			f0 = append(f0, *f0elem)
		}
		res.AttributeDefinitions = f0
	}
	if r.ko.Spec.BillingMode != nil {
		res.BillingMode = svcsdktypes.BillingMode(*r.ko.Spec.BillingMode)
	}
	if r.ko.Spec.DeletionProtectionEnabled != nil {
		res.DeletionProtectionEnabled = r.ko.Spec.DeletionProtectionEnabled
	}
	if r.ko.Spec.GlobalSecondaryIndexes != nil {
		f3 := []svcsdktypes.GlobalSecondaryIndex{}
		for _, f3iter := range r.ko.Spec.GlobalSecondaryIndexes {
			f3elem := &svcsdktypes.GlobalSecondaryIndex{}
			if f3iter.IndexName != nil {
				f3elem.IndexName = f3iter.IndexName
			}
			if f3iter.KeySchema != nil {
				f3elemf1 := []svcsdktypes.KeySchemaElement{}
				for _, f3elemf1iter := range f3iter.KeySchema {
					f3elemf1elem := &svcsdktypes.KeySchemaElement{}
					if f3elemf1iter.AttributeName != nil {
						f3elemf1elem.AttributeName = f3elemf1iter.AttributeName
					}
					if f3elemf1iter.KeyType != nil {
						f3elemf1elem.KeyType = svcsdktypes.KeyType(*f3elemf1iter.KeyType)
					}
					f3elemf1 = append(f3elemf1, *f3elemf1elem)
				}
				f3elem.KeySchema = f3elemf1
			}
			if f3iter.Projection != nil {
				f3elemf2 := &svcsdktypes.Projection{}
				if f3iter.Projection.NonKeyAttributes != nil {
					f3elemf2.NonKeyAttributes = aws.ToStringSlice(f3iter.Projection.NonKeyAttributes)
				}
				if f3iter.Projection.ProjectionType != nil {
					f3elemf2.ProjectionType = svcsdktypes.ProjectionType(*f3iter.Projection.ProjectionType)
				}
				f3elem.Projection = f3elemf2
			}
			if f3iter.ProvisionedThroughput != nil {
				f3elemf3 := &svcsdktypes.ProvisionedThroughput{}
				if f3iter.ProvisionedThroughput.ReadCapacityUnits != nil {
					f3elemf3.ReadCapacityUnits = f3iter.ProvisionedThroughput.ReadCapacityUnits
				}
				if f3iter.ProvisionedThroughput.WriteCapacityUnits != nil {
					f3elemf3.WriteCapacityUnits = f3iter.ProvisionedThroughput.WriteCapacityUnits
				}
				f3elem.ProvisionedThroughput = f3elemf3
			}
			f3 = append(f3, *f3elem)
		}
		res.GlobalSecondaryIndexes = f3
	}
	if r.ko.Spec.KeySchema != nil {
		f4 := []svcsdktypes.KeySchemaElement{}
		for _, f4iter := range r.ko.Spec.KeySchema {
			f4elem := &svcsdktypes.KeySchemaElement{}
			if f4iter.AttributeName != nil {
				f4elem.AttributeName = f4iter.AttributeName
			}
			if f4iter.KeyType != nil {
				f4elem.KeyType = svcsdktypes.KeyType(*f4iter.KeyType)
			}
			f4 = append(f4, *f4elem)
		}
		res.KeySchema = f4
	}
	if r.ko.Spec.LocalSecondaryIndexes != nil {
		f5 := []svcsdktypes.LocalSecondaryIndex{}
		for _, f5iter := range r.ko.Spec.LocalSecondaryIndexes {
			f5elem := &svcsdktypes.LocalSecondaryIndex{}
			if f5iter.IndexName != nil {
				f5elem.IndexName = f5iter.IndexName
			}
			if f5iter.KeySchema != nil {
				f5elemf1 := []svcsdktypes.KeySchemaElement{}
				for _, f5elemf1iter := range f5iter.KeySchema {
					f5elemf1elem := &svcsdktypes.KeySchemaElement{}
					if f5elemf1iter.AttributeName != nil {
						f5elemf1elem.AttributeName = f5elemf1iter.AttributeName
					}
					if f5elemf1iter.KeyType != nil {
						f5elemf1elem.KeyType = svcsdktypes.KeyType(*f5elemf1iter.KeyType)
					}
					f5elemf1 = append(f5elemf1, *f5elemf1elem)
				}
				f5elem.KeySchema = f5elemf1
			}
			if f5iter.Projection != nil {
				f5elemf2 := &svcsdktypes.Projection{}
				if f5iter.Projection.NonKeyAttributes != nil {
					f5elemf2.NonKeyAttributes = aws.ToStringSlice(f5iter.Projection.NonKeyAttributes)
				}
				if f5iter.Projection.ProjectionType != nil {
					f5elemf2.ProjectionType = svcsdktypes.ProjectionType(*f5iter.Projection.ProjectionType)
				}
				f5elem.Projection = f5elemf2
			}
			f5 = append(f5, *f5elem)
		}
		res.LocalSecondaryIndexes = f5
	}
	if r.ko.Spec.OnDemandThroughput != nil {
		f6 := &svcsdktypes.OnDemandThroughput{}
		if r.ko.Spec.OnDemandThroughput.MaxReadRequestUnits != nil {
			f6.MaxReadRequestUnits = r.ko.Spec.OnDemandThroughput.MaxReadRequestUnits
		}
		if r.ko.Spec.OnDemandThroughput.MaxWriteRequestUnits != nil {
			f6.MaxWriteRequestUnits = r.ko.Spec.OnDemandThroughput.MaxWriteRequestUnits
		}
		res.OnDemandThroughput = f6
	}
	if r.ko.Spec.ProvisionedThroughput != nil {
		f7 := &svcsdktypes.ProvisionedThroughput{}
		if r.ko.Spec.ProvisionedThroughput.ReadCapacityUnits != nil {
			f7.ReadCapacityUnits = r.ko.Spec.ProvisionedThroughput.ReadCapacityUnits
		}
		if r.ko.Spec.ProvisionedThroughput.WriteCapacityUnits != nil {
			f7.WriteCapacityUnits = r.ko.Spec.ProvisionedThroughput.WriteCapacityUnits
		}
		res.ProvisionedThroughput = f7
	}
	if r.ko.Spec.ResourcePolicy != nil {
		res.ResourcePolicy = r.ko.Spec.ResourcePolicy
	}
	if r.ko.Spec.SSESpecification != nil {
		f9 := &svcsdktypes.SSESpecification{}
		if r.ko.Spec.SSESpecification.Enabled != nil {
			f9.Enabled = r.ko.Spec.SSESpecification.Enabled
		}
		if r.ko.Spec.SSESpecification.KMSMasterKeyID != nil {
			f9.KMSMasterKeyId = r.ko.Spec.SSESpecification.KMSMasterKeyID
		}
		if r.ko.Spec.SSESpecification.SSEType != nil {
			f9.SSEType = svcsdktypes.SSEType(*r.ko.Spec.SSESpecification.SSEType)
		}
		res.SSESpecification = f9
	}
	if r.ko.Spec.StreamSpecification != nil {
		f10 := &svcsdktypes.StreamSpecification{}
		if r.ko.Spec.StreamSpecification.StreamEnabled != nil {
			f10.StreamEnabled = r.ko.Spec.StreamSpecification.StreamEnabled
		}
		if r.ko.Spec.StreamSpecification.StreamViewType != nil {
			f10.StreamViewType = svcsdktypes.StreamViewType(*r.ko.Spec.StreamSpecification.StreamViewType)
		}
		res.StreamSpecification = f10
	}
	if r.ko.Spec.TableClass != nil {
		res.TableClass = svcsdktypes.TableClass(*r.ko.Spec.TableClass)
	}
	if r.ko.Spec.TableName != nil {
		res.TableName = r.ko.Spec.TableName
	}
	if r.ko.Spec.Tags != nil {
		f13 := []svcsdktypes.Tag{}
		for _, f13iter := range r.ko.Spec.Tags {
			f13elem := &svcsdktypes.Tag{}
			if f13iter.Key != nil {
				f13elem.Key = f13iter.Key
			}
			if f13iter.Value != nil {
				f13elem.Value = f13iter.Value
			}
			f13 = append(f13, *f13elem)
		}
		res.Tags = f13
	}
	if r.ko.Spec.WarmThroughput != nil {
		f14 := &svcsdktypes.WarmThroughput{}
		if r.ko.Spec.WarmThroughput.ReadUnitsPerSecond != nil {
			f14.ReadUnitsPerSecond = r.ko.Spec.WarmThroughput.ReadUnitsPerSecond
		}
		if r.ko.Spec.WarmThroughput.WriteUnitsPerSecond != nil {
			f14.WriteUnitsPerSecond = r.ko.Spec.WarmThroughput.WriteUnitsPerSecond
		}
		res.WarmThroughput = f14
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
		f0 := &svcsdktypes.ImageScanningConfiguration{}
		if r.ko.Spec.ImageScanningConfiguration.ScanOnPush != nil {
			f0.ScanOnPush = *r.ko.Spec.ImageScanningConfiguration.ScanOnPush
		}
		res.ImageScanningConfiguration = f0
	}
	if r.ko.Spec.ImageTagMutability != nil {
		res.ImageTagMutability = svcsdktypes.ImageTagMutability(*r.ko.Spec.ImageTagMutability)
	}
	if r.ko.Spec.RegistryID != nil {
		res.RegistryId = r.ko.Spec.RegistryID
	}
	if r.ko.Spec.RepositoryName != nil {
		res.RepositoryName = r.ko.Spec.RepositoryName
	}
	if r.ko.Spec.Tags != nil {
		f4 := []svcsdktypes.Tag{}
		for _, f4iter := range r.ko.Spec.Tags {
			f4elem := &svcsdktypes.Tag{}
			if f4iter.Key != nil {
				f4elem.Key = f4iter.Key
			}
			if f4iter.Value != nil {
				f4elem.Value = f4iter.Value
			}
			f4 = append(f4, *f4elem)
		}
		res.Tags = f4
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}

// func TestSetSDK_Elasticache_ReplicationGroup_Create(t *testing.T) {
// 	assert := assert.New(t)
// 	require := require.New(t)

// 	g := testutil.NewModelForService(t, "elasticache")

// 	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
// 	require.NotNil(crd)

// 	expected := `
// 	if r.ko.Spec.AtRestEncryptionEnabled != nil {
// 		res.SetAtRestEncryptionEnabled(*r.ko.Spec.AtRestEncryptionEnabled)
// 	}
// 	if r.ko.Spec.AuthToken != nil {
// 		tmpSecret, err := rm.rr.SecretValueFromReference(ctx, r.ko.Spec.AuthToken)
// 		if err != nil {
// 			return nil, ackrequeue.Needed(err)
// 		}
// 		if tmpSecret != "" {
// 			res.SetAuthToken(tmpSecret)
// 		}
// 	}
// 	if r.ko.Spec.AutoMinorVersionUpgrade != nil {
// 		res.SetAutoMinorVersionUpgrade(*r.ko.Spec.AutoMinorVersionUpgrade)
// 	}
// 	if r.ko.Spec.AutomaticFailoverEnabled != nil {
// 		res.SetAutomaticFailoverEnabled(*r.ko.Spec.AutomaticFailoverEnabled)
// 	}
// 	if r.ko.Spec.CacheNodeType != nil {
// 		res.SetCacheNodeType(*r.ko.Spec.CacheNodeType)
// 	}
// 	if r.ko.Spec.CacheParameterGroupName != nil {
// 		res.SetCacheParameterGroupName(*r.ko.Spec.CacheParameterGroupName)
// 	}
// 	if r.ko.Spec.CacheSecurityGroupNames != nil {
// 		f6 := []*string{}
// 		for _, f6iter := range r.ko.Spec.CacheSecurityGroupNames {
// 			var f6elem string
// 			f6elem = *f6iter
// 			f6 = append(f6, &f6elem)
// 		}
// 		res.SetCacheSecurityGroupNames(f6)
// 	}
// 	if r.ko.Spec.CacheSubnetGroupName != nil {
// 		res.SetCacheSubnetGroupName(*r.ko.Spec.CacheSubnetGroupName)
// 	}
// 	if r.ko.Spec.Engine != nil {
// 		res.SetEngine(*r.ko.Spec.Engine)
// 	}
// 	if r.ko.Spec.EngineVersion != nil {
// 		res.SetEngineVersion(*r.ko.Spec.EngineVersion)
// 	}
// 	if r.ko.Spec.KMSKeyID != nil {
// 		res.SetKmsKeyId(*r.ko.Spec.KMSKeyID)
// 	}
// 	if r.ko.Spec.LogDeliveryConfigurations != nil {
// 		f11 := []*svcsdk.LogDeliveryConfigurationRequest{}
// 		for _, f11iter := range r.ko.Spec.LogDeliveryConfigurations {
// 			f11elem := &svcsdk.LogDeliveryConfigurationRequest{}
// 			if f11iter.DestinationDetails != nil {
// 				f11elemf0 := &svcsdk.DestinationDetails{}
// 				if f11iter.DestinationDetails.CloudWatchLogsDetails != nil {
// 					f11elemf0f0 := &svcsdk.CloudWatchLogsDestinationDetails{}
// 					if f11iter.DestinationDetails.CloudWatchLogsDetails.LogGroup != nil {
// 						f11elemf0f0.SetLogGroup(*f11iter.DestinationDetails.CloudWatchLogsDetails.LogGroup)
// 					}
// 					f11elemf0.SetCloudWatchLogsDetails(f11elemf0f0)
// 				}
// 				if f11iter.DestinationDetails.KinesisFirehoseDetails != nil {
// 					f11elemf0f1 := &svcsdk.KinesisFirehoseDestinationDetails{}
// 					if f11iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream != nil {
// 						f11elemf0f1.SetDeliveryStream(*f11iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream)
// 					}
// 					f11elemf0.SetKinesisFirehoseDetails(f11elemf0f1)
// 				}
// 				f11elem.SetDestinationDetails(f11elemf0)
// 			}
// 			if f11iter.DestinationType != nil {
// 				f11elem.SetDestinationType(*f11iter.DestinationType)
// 			}
// 			if f11iter.Enabled != nil {
// 				f11elem.SetEnabled(*f11iter.Enabled)
// 			}
// 			if f11iter.LogFormat != nil {
// 				f11elem.SetLogFormat(*f11iter.LogFormat)
// 			}
// 			if f11iter.LogType != nil {
// 				f11elem.SetLogType(*f11iter.LogType)
// 			}
// 			f11 = append(f11, f11elem)
// 		}
// 		res.SetLogDeliveryConfigurations(f11)
// 	}
// 	if r.ko.Spec.MultiAZEnabled != nil {
// 		res.SetMultiAZEnabled(*r.ko.Spec.MultiAZEnabled)
// 	}
// 	if r.ko.Spec.NodeGroupConfiguration != nil {
// 		f13 := []*svcsdk.NodeGroupConfiguration{}
// 		for _, f13iter := range r.ko.Spec.NodeGroupConfiguration {
// 			f13elem := &svcsdk.NodeGroupConfiguration{}
// 			if f13iter.NodeGroupID != nil {
// 				f13elem.SetNodeGroupId(*f13iter.NodeGroupID)
// 			}
// 			if f13iter.PrimaryAvailabilityZone != nil {
// 				f13elem.SetPrimaryAvailabilityZone(*f13iter.PrimaryAvailabilityZone)
// 			}
// 			if f13iter.PrimaryOutpostARN != nil {
// 				f13elem.SetPrimaryOutpostArn(*f13iter.PrimaryOutpostARN)
// 			}
// 			if f13iter.ReplicaAvailabilityZones != nil {
// 				f13elemf3 := []*string{}
// 				for _, f13elemf3iter := range f13iter.ReplicaAvailabilityZones {
// 					var f13elemf3elem string
// 					f13elemf3elem = *f13elemf3iter
// 					f13elemf3 = append(f13elemf3, &f13elemf3elem)
// 				}
// 				f13elem.SetReplicaAvailabilityZones(f13elemf3)
// 			}
// 			if f13iter.ReplicaCount != nil {
// 				f13elem.SetReplicaCount(*f13iter.ReplicaCount)
// 			}
// 			if f13iter.ReplicaOutpostARNs != nil {
// 				f13elemf5 := []*string{}
// 				for _, f13elemf5iter := range f13iter.ReplicaOutpostARNs {
// 					var f13elemf5elem string
// 					f13elemf5elem = *f13elemf5iter
// 					f13elemf5 = append(f13elemf5, &f13elemf5elem)
// 				}
// 				f13elem.SetReplicaOutpostArns(f13elemf5)
// 			}
// 			if f13iter.Slots != nil {
// 				f13elem.SetSlots(*f13iter.Slots)
// 			}
// 			f13 = append(f13, f13elem)
// 		}
// 		res.SetNodeGroupConfiguration(f13)
// 	}
// 	if r.ko.Spec.NotificationTopicARN != nil {
// 		res.SetNotificationTopicArn(*r.ko.Spec.NotificationTopicARN)
// 	}
// 	if r.ko.Spec.NumCacheClusters != nil {
// 		res.SetNumCacheClusters(*r.ko.Spec.NumCacheClusters)
// 	}
// 	if r.ko.Spec.NumNodeGroups != nil {
// 		res.SetNumNodeGroups(*r.ko.Spec.NumNodeGroups)
// 	}
// 	if r.ko.Spec.Port != nil {
// 		res.SetPort(*r.ko.Spec.Port)
// 	}
// 	if r.ko.Spec.PreferredCacheClusterAZs != nil {
// 		f18 := []*string{}
// 		for _, f18iter := range r.ko.Spec.PreferredCacheClusterAZs {
// 			var f18elem string
// 			f18elem = *f18iter
// 			f18 = append(f18, &f18elem)
// 		}
// 		res.SetPreferredCacheClusterAZs(f18)
// 	}
// 	if r.ko.Spec.PreferredMaintenanceWindow != nil {
// 		res.SetPreferredMaintenanceWindow(*r.ko.Spec.PreferredMaintenanceWindow)
// 	}
// 	if r.ko.Spec.PrimaryClusterID != nil {
// 		res.SetPrimaryClusterId(*r.ko.Spec.PrimaryClusterID)
// 	}
// 	if r.ko.Spec.ReplicasPerNodeGroup != nil {
// 		res.SetReplicasPerNodeGroup(*r.ko.Spec.ReplicasPerNodeGroup)
// 	}
// 	if r.ko.Spec.ReplicationGroupDescription != nil {
// 		res.SetReplicationGroupDescription(*r.ko.Spec.ReplicationGroupDescription)
// 	}
// 	if r.ko.Spec.ReplicationGroupID != nil {
// 		res.SetReplicationGroupId(*r.ko.Spec.ReplicationGroupID)
// 	}
// 	if r.ko.Spec.SecurityGroupIDs != nil {
// 		f24 := []*string{}
// 		for _, f24iter := range r.ko.Spec.SecurityGroupIDs {
// 			var f24elem string
// 			f24elem = *f24iter
// 			f24 = append(f24, &f24elem)
// 		}
// 		res.SetSecurityGroupIds(f24)
// 	}
// 	if r.ko.Spec.SnapshotARNs != nil {
// 		f25 := []*string{}
// 		for _, f25iter := range r.ko.Spec.SnapshotARNs {
// 			var f25elem string
// 			f25elem = *f25iter
// 			f25 = append(f25, &f25elem)
// 		}
// 		res.SetSnapshotArns(f25)
// 	}
// 	if r.ko.Spec.SnapshotName != nil {
// 		res.SetSnapshotName(*r.ko.Spec.SnapshotName)
// 	}
// 	if r.ko.Spec.SnapshotRetentionLimit != nil {
// 		res.SetSnapshotRetentionLimit(*r.ko.Spec.SnapshotRetentionLimit)
// 	}
// 	if r.ko.Spec.SnapshotWindow != nil {
// 		res.SetSnapshotWindow(*r.ko.Spec.SnapshotWindow)
// 	}
// 	if r.ko.Spec.Tags != nil {
// 		f29 := []*svcsdk.Tag{}
// 		for _, f29iter := range r.ko.Spec.Tags {
// 			f29elem := &svcsdk.Tag{}
// 			if f29iter.Key != nil {
// 				f29elem.SetKey(*f29iter.Key)
// 			}
// 			if f29iter.Value != nil {
// 				f29elem.SetValue(*f29iter.Value)
// 			}
// 			f29 = append(f29, f29elem)
// 		}
// 		res.SetTags(f29)
// 	}
// 	if r.ko.Spec.TransitEncryptionEnabled != nil {
// 		res.SetTransitEncryptionEnabled(*r.ko.Spec.TransitEncryptionEnabled)
// 	}
// 	if r.ko.Spec.UserGroupIDs != nil {
// 		f31 := []*string{}
// 		for _, f31iter := range r.ko.Spec.UserGroupIDs {
// 			var f31elem string
// 			f31elem = *f31iter
// 			f31 = append(f31, &f31elem)
// 		}
// 		res.SetUserGroupIds(f31)
// 	}
// `
// 	assert.Equal(
// 		expected,
// 		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
// 	)
// }

// func TestSetSDK_Elasticache_ReplicationGroup_ReadMany(t *testing.T) {
// 	assert := assert.New(t)
// 	require := require.New(t)

// 	g := testutil.NewModelForService(t, "elasticache")

// 	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
// 	require.NotNil(crd)

// 	// Elasticache doesn't have a ReadOne operation; only a List/ReadMany
// 	// operation. Let's verify that the construction of the
// 	// DescribeCacheClustersInput and processing of the
// 	// DescribeCacheClustersOutput shapes is correct.
// 	expected := `
// 	if r.ko.Spec.ReplicationGroupID != nil {
// 		res.SetReplicationGroupId(*r.ko.Spec.ReplicationGroupID)
// 	}
// `
// 	assert.Equal(
// 		expected,
// 		code.SetSDK(crd.Config(), crd, model.OpTypeList, "r.ko", "res", 1),
// 	)
// }

// func TestSetSDK_Elasticache_ReplicationGroup_Update_Override_Values(t *testing.T) {
// 	assert := assert.New(t)
// 	require := require.New(t)

// 	g := testutil.NewModelForService(t, "elasticache")

// 	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
// 	require.NotNil(crd)

// 	expected := `
// 	res.SetApplyImmediately(true)
// 	if r.ko.Spec.AuthToken != nil {
// 		tmpSecret, err := rm.rr.SecretValueFromReference(ctx, r.ko.Spec.AuthToken)
// 		if err != nil {
// 			return nil, ackrequeue.Needed(err)
// 		}
// 		if tmpSecret != "" {
// 			res.SetAuthToken(tmpSecret)
// 		}
// 	}
// 	if r.ko.Spec.AutoMinorVersionUpgrade != nil {
// 		res.SetAutoMinorVersionUpgrade(*r.ko.Spec.AutoMinorVersionUpgrade)
// 	}
// 	if r.ko.Spec.AutomaticFailoverEnabled != nil {
// 		res.SetAutomaticFailoverEnabled(*r.ko.Spec.AutomaticFailoverEnabled)
// 	}
// 	if r.ko.Spec.CacheNodeType != nil {
// 		res.SetCacheNodeType(*r.ko.Spec.CacheNodeType)
// 	}
// 	if r.ko.Spec.CacheParameterGroupName != nil {
// 		res.SetCacheParameterGroupName(*r.ko.Spec.CacheParameterGroupName)
// 	}
// 	if r.ko.Spec.CacheSecurityGroupNames != nil {
// 		f7 := []*string{}
// 		for _, f7iter := range r.ko.Spec.CacheSecurityGroupNames {
// 			var f7elem string
// 			f7elem = *f7iter
// 			f7 = append(f7, &f7elem)
// 		}
// 		res.SetCacheSecurityGroupNames(f7)
// 	}
// 	if r.ko.Spec.LogDeliveryConfigurations != nil {
// 		f8 := []*svcsdk.LogDeliveryConfigurationRequest{}
// 		for _, f8iter := range r.ko.Spec.LogDeliveryConfigurations {
// 			f8elem := &svcsdk.LogDeliveryConfigurationRequest{}
// 			if f8iter.DestinationDetails != nil {
// 				f8elemf0 := &svcsdk.DestinationDetails{}
// 				if f8iter.DestinationDetails.CloudWatchLogsDetails != nil {
// 					f8elemf0f0 := &svcsdk.CloudWatchLogsDestinationDetails{}
// 					if f8iter.DestinationDetails.CloudWatchLogsDetails.LogGroup != nil {
// 						f8elemf0f0.SetLogGroup(*f8iter.DestinationDetails.CloudWatchLogsDetails.LogGroup)
// 					}
// 					f8elemf0.SetCloudWatchLogsDetails(f8elemf0f0)
// 				}
// 				if f8iter.DestinationDetails.KinesisFirehoseDetails != nil {
// 					f8elemf0f1 := &svcsdk.KinesisFirehoseDestinationDetails{}
// 					if f8iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream != nil {
// 						f8elemf0f1.SetDeliveryStream(*f8iter.DestinationDetails.KinesisFirehoseDetails.DeliveryStream)
// 					}
// 					f8elemf0.SetKinesisFirehoseDetails(f8elemf0f1)
// 				}
// 				f8elem.SetDestinationDetails(f8elemf0)
// 			}
// 			if f8iter.DestinationType != nil {
// 				f8elem.SetDestinationType(*f8iter.DestinationType)
// 			}
// 			if f8iter.Enabled != nil {
// 				f8elem.SetEnabled(*f8iter.Enabled)
// 			}
// 			if f8iter.LogFormat != nil {
// 				f8elem.SetLogFormat(*f8iter.LogFormat)
// 			}
// 			if f8iter.LogType != nil {
// 				f8elem.SetLogType(*f8iter.LogType)
// 			}
// 			f8 = append(f8, f8elem)
// 		}
// 		res.SetLogDeliveryConfigurations(f8)
// 	}
// 	if r.ko.Spec.MultiAZEnabled != nil {
// 		res.SetMultiAZEnabled(*r.ko.Spec.MultiAZEnabled)
// 	}
// 	if r.ko.Spec.NotificationTopicARN != nil {
// 		res.SetNotificationTopicArn(*r.ko.Spec.NotificationTopicARN)
// 	}
// 	if r.ko.Spec.PreferredMaintenanceWindow != nil {
// 		res.SetPreferredMaintenanceWindow(*r.ko.Spec.PreferredMaintenanceWindow)
// 	}
// 	if r.ko.Spec.PrimaryClusterID != nil {
// 		res.SetPrimaryClusterId(*r.ko.Spec.PrimaryClusterID)
// 	}
// 	if r.ko.Spec.ReplicationGroupDescription != nil {
// 		res.SetReplicationGroupDescription(*r.ko.Spec.ReplicationGroupDescription)
// 	}
// 	if r.ko.Spec.ReplicationGroupID != nil {
// 		res.SetReplicationGroupId(*r.ko.Spec.ReplicationGroupID)
// 	}
// 	if r.ko.Spec.SnapshotRetentionLimit != nil {
// 		res.SetSnapshotRetentionLimit(*r.ko.Spec.SnapshotRetentionLimit)
// 	}
// 	if r.ko.Spec.SnapshotWindow != nil {
// 		res.SetSnapshotWindow(*r.ko.Spec.SnapshotWindow)
// 	}
// 	if r.ko.Status.SnapshottingClusterID != nil {
// 		res.SetSnapshottingClusterId(*r.ko.Status.SnapshottingClusterID)
// 	}
// `
// 	assert.Equal(
// 		expected,
// 		code.SetSDK(crd.Config(), crd, model.OpTypeUpdate, "r.ko", "res", 1),
// 	)
// }

// func TestSetSDK_Elasticache_User_Create_Override_Values(t *testing.T) {
// 	assert := assert.New(t)
// 	require := require.New(t)

// 	g := testutil.NewModelForService(t, "elasticache")

// 	crd := testutil.GetCRDByName(t, g, "User")
// 	require.NotNil(crd)

// 	expected := `
// 	if r.ko.Spec.AccessString != nil {
// 		res.SetAccessString(*r.ko.Spec.AccessString)
// 	}
// 	if r.ko.Spec.NoPasswordRequired != nil {
// 		res.SetNoPasswordRequired(*r.ko.Spec.NoPasswordRequired)
// 	}
// 	if r.ko.Spec.Passwords != nil {
// 		f3 := []*string{}
// 		for _, f3iter := range r.ko.Spec.Passwords {
// 			var f3elem string
// 			if f3iter != nil {
// 				tmpSecret, err := rm.rr.SecretValueFromReference(ctx, f3iter)
// 				if err != nil {
// 					return nil, ackrequeue.Needed(err)
// 				}
// 				if tmpSecret != "" {
// 					f3elem = tmpSecret
// 				}
// 			}
// 			f3 = append(f3, &f3elem)
// 		}
// 		res.SetPasswords(f3)
// 	}
// 	if r.ko.Spec.UserID != nil {
// 		res.SetUserId(*r.ko.Spec.UserID)
// 	}
// `
// 	assert.Equal(
// 		expected,
// 		code.SetSDK(crd.Config(), crd, model.OpTypeUpdate, "r.ko", "res", 1),
// 	)
// }

func TestSetSDK_MQ_Broker_newUpdateRequest_OmitUnchangedValues(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "mq")

	crd := testutil.GetCRDByName(t, g, "Broker")
	require.NotNil(crd)

	expected := `
	if delta.DifferentAt("Spec.AuthenticationStrategy") {
		if r.ko.Spec.AuthenticationStrategy != nil {
			res.AuthenticationStrategy = svcsdktypes.AuthenticationStrategy(*r.ko.Spec.AuthenticationStrategy)
		}
	}
	if delta.DifferentAt("Spec.AutoMinorVersionUpgrade") {
		if r.ko.Spec.AutoMinorVersionUpgrade != nil {
			res.AutoMinorVersionUpgrade = r.ko.Spec.AutoMinorVersionUpgrade
		}
	}
	if r.ko.Status.BrokerID != nil {
		res.BrokerId = r.ko.Status.BrokerID
	}
	if delta.DifferentAt("Spec.Configuration") {
		if r.ko.Spec.Configuration != nil {
			f3 := &svcsdktypes.ConfigurationId{}
			if r.ko.Spec.Configuration.ID != nil {
				f3.Id = r.ko.Spec.Configuration.ID
			}
			if r.ko.Spec.Configuration.Revision != nil {
				if *r.ko.Spec.Configuration.Revision > math.MaxInt32 || *r.ko.Spec.Configuration.Revision < math.MinInt32 {
					return nil, fmt.Errorf("error: field Revision is of type int32")
				}
				revisionCopy := int32(*r.ko.Spec.Configuration.Revision)
				f3.Revision = &revisionCopy
			}
			res.Configuration = f3
		}
	}
	if delta.DifferentAt("Spec.EngineVersion") {
		if r.ko.Spec.EngineVersion != nil {
			res.EngineVersion = r.ko.Spec.EngineVersion
		}
	}
	if delta.DifferentAt("Spec.HostInstanceType") {
		if r.ko.Spec.HostInstanceType != nil {
			res.HostInstanceType = r.ko.Spec.HostInstanceType
		}
	}
	if delta.DifferentAt("Spec.LDAPServerMetadata") {
		if r.ko.Spec.LDAPServerMetadata != nil {
			f7 := &svcsdktypes.LdapServerMetadataInput{}
			if r.ko.Spec.LDAPServerMetadata.Hosts != nil {
				f7.Hosts = aws.ToStringSlice(r.ko.Spec.LDAPServerMetadata.Hosts)
			}
			if r.ko.Spec.LDAPServerMetadata.RoleBase != nil {
				f7.RoleBase = r.ko.Spec.LDAPServerMetadata.RoleBase
			}
			if r.ko.Spec.LDAPServerMetadata.RoleName != nil {
				f7.RoleName = r.ko.Spec.LDAPServerMetadata.RoleName
			}
			if r.ko.Spec.LDAPServerMetadata.RoleSearchMatching != nil {
				f7.RoleSearchMatching = r.ko.Spec.LDAPServerMetadata.RoleSearchMatching
			}
			if r.ko.Spec.LDAPServerMetadata.RoleSearchSubtree != nil {
				f7.RoleSearchSubtree = r.ko.Spec.LDAPServerMetadata.RoleSearchSubtree
			}
			if r.ko.Spec.LDAPServerMetadata.ServiceAccountPassword != nil {
				f7.ServiceAccountPassword = r.ko.Spec.LDAPServerMetadata.ServiceAccountPassword
			}
			if r.ko.Spec.LDAPServerMetadata.ServiceAccountUsername != nil {
				f7.ServiceAccountUsername = r.ko.Spec.LDAPServerMetadata.ServiceAccountUsername
			}
			if r.ko.Spec.LDAPServerMetadata.UserBase != nil {
				f7.UserBase = r.ko.Spec.LDAPServerMetadata.UserBase
			}
			if r.ko.Spec.LDAPServerMetadata.UserRoleName != nil {
				f7.UserRoleName = r.ko.Spec.LDAPServerMetadata.UserRoleName
			}
			if r.ko.Spec.LDAPServerMetadata.UserSearchMatching != nil {
				f7.UserSearchMatching = r.ko.Spec.LDAPServerMetadata.UserSearchMatching
			}
			if r.ko.Spec.LDAPServerMetadata.UserSearchSubtree != nil {
				f7.UserSearchSubtree = r.ko.Spec.LDAPServerMetadata.UserSearchSubtree
			}
			res.LdapServerMetadata = f7
		}
	}
	if delta.DifferentAt("Spec.Logs") {
		if r.ko.Spec.Logs != nil {
			f8 := &svcsdktypes.Logs{}
			if r.ko.Spec.Logs.Audit != nil {
				f8.Audit = r.ko.Spec.Logs.Audit
			}
			if r.ko.Spec.Logs.General != nil {
				f8.General = r.ko.Spec.Logs.General
			}
			res.Logs = f8
		}
	}
	if delta.DifferentAt("Spec.MaintenanceWindowStartTime") {
		if r.ko.Spec.MaintenanceWindowStartTime != nil {
			f9 := &svcsdktypes.WeeklyStartTime{}
			if r.ko.Spec.MaintenanceWindowStartTime.DayOfWeek != nil {
				f9.DayOfWeek = svcsdktypes.DayOfWeek(*r.ko.Spec.MaintenanceWindowStartTime.DayOfWeek)
			}
			if r.ko.Spec.MaintenanceWindowStartTime.TimeOfDay != nil {
				f9.TimeOfDay = r.ko.Spec.MaintenanceWindowStartTime.TimeOfDay
			}
			if r.ko.Spec.MaintenanceWindowStartTime.TimeZone != nil {
				f9.TimeZone = r.ko.Spec.MaintenanceWindowStartTime.TimeZone
			}
			res.MaintenanceWindowStartTime = f9
		}
	}
	if delta.DifferentAt("Spec.SecurityGroups") {
		if r.ko.Spec.SecurityGroups != nil {
			res.SecurityGroups = aws.ToStringSlice(r.ko.Spec.SecurityGroups)
		}
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
		if *r.ko.Spec.AllocatedStorage > math.MaxInt32 || *r.ko.Spec.AllocatedStorage < math.MinInt32 {
			return nil, fmt.Errorf("error: field AllocatedStorage is of type int32")
		}
		allocatedStorageCopy := int32(*r.ko.Spec.AllocatedStorage)
		res.AllocatedStorage = &allocatedStorageCopy
	}
	if r.ko.Spec.AutoMinorVersionUpgrade != nil {
		res.AutoMinorVersionUpgrade = r.ko.Spec.AutoMinorVersionUpgrade
	}
	if r.ko.Spec.AvailabilityZone != nil {
		res.AvailabilityZone = r.ko.Spec.AvailabilityZone
	}
	if r.ko.Spec.BackupRetentionPeriod != nil {
		if *r.ko.Spec.BackupRetentionPeriod > math.MaxInt32 || *r.ko.Spec.BackupRetentionPeriod < math.MinInt32 {
			return nil, fmt.Errorf("error: field BackupRetentionPeriod is of type int32")
		}
		backupRetentionPeriodCopy := int32(*r.ko.Spec.BackupRetentionPeriod)
		res.BackupRetentionPeriod = &backupRetentionPeriodCopy
	}
	if r.ko.Spec.CharacterSetName != nil {
		res.CharacterSetName = r.ko.Spec.CharacterSetName
	}
	if r.ko.Spec.CopyTagsToSnapshot != nil {
		res.CopyTagsToSnapshot = r.ko.Spec.CopyTagsToSnapshot
	}
	if r.ko.Spec.DBClusterIdentifier != nil {
		res.DBClusterIdentifier = r.ko.Spec.DBClusterIdentifier
	}
	if r.ko.Spec.DBInstanceClass != nil {
		res.DBInstanceClass = r.ko.Spec.DBInstanceClass
	}
	if r.ko.Spec.DBInstanceIdentifier != nil {
		res.DBInstanceIdentifier = r.ko.Spec.DBInstanceIdentifier
	}
	if r.ko.Spec.DBName != nil {
		res.DBName = r.ko.Spec.DBName
	}
	if r.ko.Spec.DBParameterGroupName != nil {
		res.DBParameterGroupName = r.ko.Spec.DBParameterGroupName
	}
	if r.ko.Spec.DBSecurityGroups != nil {
		res.DBSecurityGroups = aws.ToStringSlice(r.ko.Spec.DBSecurityGroups)
	}
	if r.ko.Spec.DBSubnetGroupName != nil {
		res.DBSubnetGroupName = r.ko.Spec.DBSubnetGroupName
	}
	if r.ko.Spec.DeletionProtection != nil {
		res.DeletionProtection = r.ko.Spec.DeletionProtection
	}
	if r.ko.Spec.Domain != nil {
		res.Domain = r.ko.Spec.Domain
	}
	if r.ko.Spec.DomainIAMRoleName != nil {
		res.DomainIAMRoleName = r.ko.Spec.DomainIAMRoleName
	}
	if r.ko.Spec.EnableCloudwatchLogsExports != nil {
		res.EnableCloudwatchLogsExports = aws.ToStringSlice(r.ko.Spec.EnableCloudwatchLogsExports)
	}
	if r.ko.Spec.EnableIAMDatabaseAuthentication != nil {
		res.EnableIAMDatabaseAuthentication = r.ko.Spec.EnableIAMDatabaseAuthentication
	}
	if r.ko.Spec.EnablePerformanceInsights != nil {
		res.EnablePerformanceInsights = r.ko.Spec.EnablePerformanceInsights
	}
	if r.ko.Spec.Engine != nil {
		res.Engine = r.ko.Spec.Engine
	}
	if r.ko.Spec.EngineVersion != nil {
		res.EngineVersion = r.ko.Spec.EngineVersion
	}
	if r.ko.Spec.IOPS != nil {
		if *r.ko.Spec.IOPS > math.MaxInt32 || *r.ko.Spec.IOPS < math.MinInt32 {
			return nil, fmt.Errorf("error: field Iops is of type int32")
		}
		iopsCopy := int32(*r.ko.Spec.IOPS)
		res.Iops = &iopsCopy
	}
	if r.ko.Spec.KMSKeyID != nil {
		res.KmsKeyId = r.ko.Spec.KMSKeyID
	}
	if r.ko.Spec.LicenseModel != nil {
		res.LicenseModel = r.ko.Spec.LicenseModel
	}
	if r.ko.Spec.MasterUserPassword != nil {
		res.MasterUserPassword = r.ko.Spec.MasterUserPassword
	}
	if r.ko.Spec.MasterUsername != nil {
		res.MasterUsername = r.ko.Spec.MasterUsername
	}
	if r.ko.Spec.MaxAllocatedStorage != nil {
		if *r.ko.Spec.MaxAllocatedStorage > math.MaxInt32 || *r.ko.Spec.MaxAllocatedStorage < math.MinInt32 {
			return nil, fmt.Errorf("error: field MaxAllocatedStorage is of type int32")
		}
		maxAllocatedStorageCopy := int32(*r.ko.Spec.MaxAllocatedStorage)
		res.MaxAllocatedStorage = &maxAllocatedStorageCopy
	}
	if r.ko.Spec.MonitoringInterval != nil {
		if *r.ko.Spec.MonitoringInterval > math.MaxInt32 || *r.ko.Spec.MonitoringInterval < math.MinInt32 {
			return nil, fmt.Errorf("error: field MonitoringInterval is of type int32")
		}
		monitoringIntervalCopy := int32(*r.ko.Spec.MonitoringInterval)
		res.MonitoringInterval = &monitoringIntervalCopy
	}
	if r.ko.Spec.MonitoringRoleARN != nil {
		res.MonitoringRoleArn = r.ko.Spec.MonitoringRoleARN
	}
	if r.ko.Spec.MultiAZ != nil {
		res.MultiAZ = r.ko.Spec.MultiAZ
	}
	if r.ko.Spec.OptionGroupName != nil {
		res.OptionGroupName = r.ko.Spec.OptionGroupName
	}
	if r.ko.Spec.PerformanceInsightsKMSKeyID != nil {
		res.PerformanceInsightsKMSKeyId = r.ko.Spec.PerformanceInsightsKMSKeyID
	}
	if r.ko.Spec.PerformanceInsightsRetentionPeriod != nil {
		if *r.ko.Spec.PerformanceInsightsRetentionPeriod > math.MaxInt32 || *r.ko.Spec.PerformanceInsightsRetentionPeriod < math.MinInt32 {
			return nil, fmt.Errorf("error: field PerformanceInsightsRetentionPeriod is of type int32")
		}
		performanceInsightsRetentionPeriodCopy := int32(*r.ko.Spec.PerformanceInsightsRetentionPeriod)
		res.PerformanceInsightsRetentionPeriod = &performanceInsightsRetentionPeriodCopy
	}
	if r.ko.Spec.Port != nil {
		if *r.ko.Spec.Port > math.MaxInt32 || *r.ko.Spec.Port < math.MinInt32 {
			return nil, fmt.Errorf("error: field Port is of type int32")
		}
		portCopy := int32(*r.ko.Spec.Port)
		res.Port = &portCopy
	}
	if r.ko.Spec.PreferredBackupWindow != nil {
		res.PreferredBackupWindow = r.ko.Spec.PreferredBackupWindow
	}
	if r.ko.Spec.PreferredMaintenanceWindow != nil {
		res.PreferredMaintenanceWindow = r.ko.Spec.PreferredMaintenanceWindow
	}
	if r.ko.Spec.ProcessorFeatures != nil {
		f36 := []svcsdktypes.ProcessorFeature{}
		for _, f36iter := range r.ko.Spec.ProcessorFeatures {
			f36elem := &svcsdktypes.ProcessorFeature{}
			if f36iter.Name != nil {
				f36elem.Name = f36iter.Name
			}
			if f36iter.Value != nil {
				f36elem.Value = f36iter.Value
			}
			f36 = append(f36, *f36elem)
		}
		res.ProcessorFeatures = f36
	}
	if r.ko.Spec.PromotionTier != nil {
		if *r.ko.Spec.PromotionTier > math.MaxInt32 || *r.ko.Spec.PromotionTier < math.MinInt32 {
			return nil, fmt.Errorf("error: field PromotionTier is of type int32")
		}
		promotionTierCopy := int32(*r.ko.Spec.PromotionTier)
		res.PromotionTier = &promotionTierCopy
	}
	if r.ko.Spec.PubliclyAccessible != nil {
		res.PubliclyAccessible = r.ko.Spec.PubliclyAccessible
	}
	if r.ko.Spec.StorageEncrypted != nil {
		res.StorageEncrypted = r.ko.Spec.StorageEncrypted
	}
	if r.ko.Spec.StorageType != nil {
		res.StorageType = r.ko.Spec.StorageType
	}
	if r.ko.Spec.Tags != nil {
		f41 := []svcsdktypes.Tag{}
		for _, f41iter := range r.ko.Spec.Tags {
			f41elem := &svcsdktypes.Tag{}
			if f41iter.Key != nil {
				f41elem.Key = f41iter.Key
			}
			if f41iter.Value != nil {
				f41elem.Value = f41iter.Value
			}
			f41 = append(f41, *f41elem)
		}
		res.Tags = f41
	}
	if r.ko.Spec.TDECredentialARN != nil {
		res.TdeCredentialArn = r.ko.Spec.TDECredentialARN
	}
	if r.ko.Spec.TDECredentialPassword != nil {
		res.TdeCredentialPassword = r.ko.Spec.TDECredentialPassword
	}
	if r.ko.Spec.Timezone != nil {
		res.Timezone = r.ko.Spec.Timezone
	}
	if r.ko.Spec.VPCSecurityGroupIDs != nil {
		res.VpcSecurityGroupIds = aws.ToStringSlice(r.ko.Spec.VPCSecurityGroupIDs)
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
		res.ACL = svcsdktypes.BucketCannedACL(*r.ko.Spec.ACL)
	}
	if r.ko.Spec.Name != nil {
		res.Bucket = r.ko.Spec.Name
	}
	if r.ko.Spec.CreateBucketConfiguration != nil {
		f2 := &svcsdktypes.CreateBucketConfiguration{}
		if r.ko.Spec.CreateBucketConfiguration.LocationConstraint != nil {
			f2.LocationConstraint = svcsdktypes.BucketLocationConstraint(*r.ko.Spec.CreateBucketConfiguration.LocationConstraint)
		}
		res.CreateBucketConfiguration = f2
	}
	if r.ko.Spec.GrantFullControl != nil {
		res.GrantFullControl = r.ko.Spec.GrantFullControl
	}
	if r.ko.Spec.GrantRead != nil {
		res.GrantRead = r.ko.Spec.GrantRead
	}
	if r.ko.Spec.GrantReadACP != nil {
		res.GrantReadACP = r.ko.Spec.GrantReadACP
	}
	if r.ko.Spec.GrantWrite != nil {
		res.GrantWrite = r.ko.Spec.GrantWrite
	}
	if r.ko.Spec.GrantWriteACP != nil {
		res.GrantWriteACP = r.ko.Spec.GrantWriteACP
	}
	if r.ko.Spec.ObjectLockEnabledForBucket != nil {
		res.ObjectLockEnabledForBucket = r.ko.Spec.ObjectLockEnabledForBucket
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
		res.Bucket = r.ko.Spec.Name
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
	attrMap := map[string]string{}
	if r.ko.Spec.DeliveryPolicy != nil {
		attrMap["DeliveryPolicy"] = *r.ko.Spec.DeliveryPolicy
	}
	if r.ko.Spec.DisplayName != nil {
		attrMap["DisplayName"] = *r.ko.Spec.DisplayName
	}
	if r.ko.Spec.KMSMasterKeyID != nil {
		attrMap["KmsMasterKeyId"] = *r.ko.Spec.KMSMasterKeyID
	}
	if r.ko.Spec.Policy != nil {
		attrMap["Policy"] = *r.ko.Spec.Policy
	}
	if len(attrMap) > 0 {
		res.Attributes = attrMap
	}
	if r.ko.Spec.Name != nil {
		res.Name = r.ko.Spec.Name
	}
	if r.ko.Spec.Tags != nil {
		f2 := []svcsdktypes.Tag{}
		for _, f2iter := range r.ko.Spec.Tags {
			f2elem := &svcsdktypes.Tag{}
			if f2iter.Key != nil {
				f2elem.Key = f2iter.Key
			}
			if f2iter.Value != nil {
				f2elem.Value = f2iter.Value
			}
			f2 = append(f2, *f2elem)
		}
		res.Tags = f2
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
		res.TopicArn = aws.String(string(*r.ko.Status.ACKResourceMetadata.ARN))
	} else {
		res.TopicArn = aws.String(rm.ARNFromName(*r.ko.Spec.Name))
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
	attrMap := map[string]string{}
	if r.ko.Spec.ContentBasedDeduplication != nil {
		attrMap["ContentBasedDeduplication"] = *r.ko.Spec.ContentBasedDeduplication
	}
	if r.ko.Spec.CreatedTimestamp != nil {
		attrMap["CreatedTimestamp"] = *r.ko.Spec.CreatedTimestamp
	}
	if r.ko.Spec.DelaySeconds != nil {
		attrMap["DelaySeconds"] = *r.ko.Spec.DelaySeconds
	}
	if r.ko.Spec.FIFOQueue != nil {
		attrMap["FifoQueue"] = *r.ko.Spec.FIFOQueue
	}
	if r.ko.Spec.KMSDataKeyReusePeriodSeconds != nil {
		attrMap["KmsDataKeyReusePeriodSeconds"] = *r.ko.Spec.KMSDataKeyReusePeriodSeconds
	}
	if r.ko.Spec.KMSMasterKeyID != nil {
		attrMap["KmsMasterKeyId"] = *r.ko.Spec.KMSMasterKeyID
	}
	if r.ko.Spec.MaximumMessageSize != nil {
		attrMap["MaximumMessageSize"] = *r.ko.Spec.MaximumMessageSize
	}
	if r.ko.Spec.MessageRetentionPeriod != nil {
		attrMap["MessageRetentionPeriod"] = *r.ko.Spec.MessageRetentionPeriod
	}
	if r.ko.Spec.Policy != nil {
		attrMap["Policy"] = *r.ko.Spec.Policy
	}
	if r.ko.Spec.ReceiveMessageWaitTimeSeconds != nil {
		attrMap["ReceiveMessageWaitTimeSeconds"] = *r.ko.Spec.ReceiveMessageWaitTimeSeconds
	}
	if r.ko.Spec.RedrivePolicy != nil {
		attrMap["RedrivePolicy"] = *r.ko.Spec.RedrivePolicy
	}
	if r.ko.Spec.VisibilityTimeout != nil {
		attrMap["VisibilityTimeout"] = *r.ko.Spec.VisibilityTimeout
	}
	if len(attrMap) > 0 {
		res.Attributes = attrMap
	}
	if r.ko.Spec.QueueName != nil {
		res.QueueName = r.ko.Spec.QueueName
	}
	if r.ko.Spec.Tags != nil {
		res.Tags = aws.ToStringMap(r.ko.Spec.Tags)
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
		tmpVals := []svcsdktypes.QueueAttributeName{}
		tmpVal0 := svcsdktypes.QueueAttributeNameAll
		tmpVals = append(tmpVals, tmpVal0)
		res.AttributeNames = tmpVals
	}
	if r.ko.Status.QueueURL != nil {
		res.QueueUrl = r.ko.Status.QueueURL
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
		res.AuthenticationStrategy = svcsdktypes.AuthenticationStrategy(*r.ko.Spec.AuthenticationStrategy)
	}
	if r.ko.Spec.AutoMinorVersionUpgrade != nil {
		res.AutoMinorVersionUpgrade = r.ko.Spec.AutoMinorVersionUpgrade
	}
	if r.ko.Spec.BrokerName != nil {
		res.BrokerName = r.ko.Spec.BrokerName
	}
	if r.ko.Spec.Configuration != nil {
		f3 := &svcsdktypes.ConfigurationId{}
		if r.ko.Spec.Configuration.ID != nil {
			f3.Id = r.ko.Spec.Configuration.ID
		}
		if r.ko.Spec.Configuration.Revision != nil {
			if *r.ko.Spec.Configuration.Revision > math.MaxInt32 || *r.ko.Spec.Configuration.Revision < math.MinInt32 {
				return nil, fmt.Errorf("error: field Revision is of type int32")
			}
			revisionCopy := int32(*r.ko.Spec.Configuration.Revision)
			f3.Revision = &revisionCopy
		}
		res.Configuration = f3
	}
	if r.ko.Spec.CreatorRequestID != nil {
		res.CreatorRequestId = r.ko.Spec.CreatorRequestID
	}
	if r.ko.Spec.DeploymentMode != nil {
		res.DeploymentMode = svcsdktypes.DeploymentMode(*r.ko.Spec.DeploymentMode)
	}
	if r.ko.Spec.EncryptionOptions != nil {
		f6 := &svcsdktypes.EncryptionOptions{}
		if r.ko.Spec.EncryptionOptions.KMSKeyID != nil {
			f6.KmsKeyId = r.ko.Spec.EncryptionOptions.KMSKeyID
		}
		if r.ko.Spec.EncryptionOptions.UseAWSOwnedKey != nil {
			f6.UseAwsOwnedKey = r.ko.Spec.EncryptionOptions.UseAWSOwnedKey
		}
		res.EncryptionOptions = f6
	}
	if r.ko.Spec.EngineType != nil {
		res.EngineType = svcsdktypes.EngineType(*r.ko.Spec.EngineType)
	}
	if r.ko.Spec.EngineVersion != nil {
		res.EngineVersion = r.ko.Spec.EngineVersion
	}
	if r.ko.Spec.HostInstanceType != nil {
		res.HostInstanceType = r.ko.Spec.HostInstanceType
	}
	if r.ko.Spec.LDAPServerMetadata != nil {
		f10 := &svcsdktypes.LdapServerMetadataInput{}
		if r.ko.Spec.LDAPServerMetadata.Hosts != nil {
			f10.Hosts = aws.ToStringSlice(r.ko.Spec.LDAPServerMetadata.Hosts)
		}
		if r.ko.Spec.LDAPServerMetadata.RoleBase != nil {
			f10.RoleBase = r.ko.Spec.LDAPServerMetadata.RoleBase
		}
		if r.ko.Spec.LDAPServerMetadata.RoleName != nil {
			f10.RoleName = r.ko.Spec.LDAPServerMetadata.RoleName
		}
		if r.ko.Spec.LDAPServerMetadata.RoleSearchMatching != nil {
			f10.RoleSearchMatching = r.ko.Spec.LDAPServerMetadata.RoleSearchMatching
		}
		if r.ko.Spec.LDAPServerMetadata.RoleSearchSubtree != nil {
			f10.RoleSearchSubtree = r.ko.Spec.LDAPServerMetadata.RoleSearchSubtree
		}
		if r.ko.Spec.LDAPServerMetadata.ServiceAccountPassword != nil {
			f10.ServiceAccountPassword = r.ko.Spec.LDAPServerMetadata.ServiceAccountPassword
		}
		if r.ko.Spec.LDAPServerMetadata.ServiceAccountUsername != nil {
			f10.ServiceAccountUsername = r.ko.Spec.LDAPServerMetadata.ServiceAccountUsername
		}
		if r.ko.Spec.LDAPServerMetadata.UserBase != nil {
			f10.UserBase = r.ko.Spec.LDAPServerMetadata.UserBase
		}
		if r.ko.Spec.LDAPServerMetadata.UserRoleName != nil {
			f10.UserRoleName = r.ko.Spec.LDAPServerMetadata.UserRoleName
		}
		if r.ko.Spec.LDAPServerMetadata.UserSearchMatching != nil {
			f10.UserSearchMatching = r.ko.Spec.LDAPServerMetadata.UserSearchMatching
		}
		if r.ko.Spec.LDAPServerMetadata.UserSearchSubtree != nil {
			f10.UserSearchSubtree = r.ko.Spec.LDAPServerMetadata.UserSearchSubtree
		}
		res.LdapServerMetadata = f10
	}
	if r.ko.Spec.Logs != nil {
		f11 := &svcsdktypes.Logs{}
		if r.ko.Spec.Logs.Audit != nil {
			f11.Audit = r.ko.Spec.Logs.Audit
		}
		if r.ko.Spec.Logs.General != nil {
			f11.General = r.ko.Spec.Logs.General
		}
		res.Logs = f11
	}
	if r.ko.Spec.MaintenanceWindowStartTime != nil {
		f12 := &svcsdktypes.WeeklyStartTime{}
		if r.ko.Spec.MaintenanceWindowStartTime.DayOfWeek != nil {
			f12.DayOfWeek = svcsdktypes.DayOfWeek(*r.ko.Spec.MaintenanceWindowStartTime.DayOfWeek)
		}
		if r.ko.Spec.MaintenanceWindowStartTime.TimeOfDay != nil {
			f12.TimeOfDay = r.ko.Spec.MaintenanceWindowStartTime.TimeOfDay
		}
		if r.ko.Spec.MaintenanceWindowStartTime.TimeZone != nil {
			f12.TimeZone = r.ko.Spec.MaintenanceWindowStartTime.TimeZone
		}
		res.MaintenanceWindowStartTime = f12
	}
	if r.ko.Spec.PubliclyAccessible != nil {
		res.PubliclyAccessible = r.ko.Spec.PubliclyAccessible
	}
	if r.ko.Spec.SecurityGroups != nil {
		res.SecurityGroups = aws.ToStringSlice(r.ko.Spec.SecurityGroups)
	}
	if r.ko.Spec.StorageType != nil {
		res.StorageType = svcsdktypes.BrokerStorageType(*r.ko.Spec.StorageType)
	}
	if r.ko.Spec.SubnetIDs != nil {
		res.SubnetIds = aws.ToStringSlice(r.ko.Spec.SubnetIDs)
	}
	if r.ko.Spec.Tags != nil {
		res.Tags = aws.ToStringMap(r.ko.Spec.Tags)
	}
	if r.ko.Spec.Users != nil {
		f18 := []svcsdktypes.User{}
		for _, f18iter := range r.ko.Spec.Users {
			f18elem := &svcsdktypes.User{}
			if f18iter.ConsoleAccess != nil {
				f18elem.ConsoleAccess = f18iter.ConsoleAccess
			}
			if f18iter.Groups != nil {
				f18elem.Groups = aws.ToStringSlice(f18iter.Groups)
			}
			if f18iter.Password != nil {
				tmpSecret, err := rm.rr.SecretValueFromReference(ctx, f18iter.Password)
				if err != nil {
					return nil, ackrequeue.Needed(err)
				}
				if tmpSecret != "" {
					f18elem.Password = aws.String(tmpSecret)
				}
			}
			if f18iter.Username != nil {
				f18elem.Username = f18iter.Username
			}
			f18 = append(f18, *f18elem)
		}
		res.Users = f18
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
		f4 := []string{}
		f4 = append(f4, *r.ko.Status.VPCID)
		res.VpcIds = f4
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
	if r.ko.Spec.RegistryID != nil {
		res.RegistryId = r.ko.Spec.RegistryID
	}
	if r.ko.Spec.Name != nil {
		f3 := []string{}
		f3 = append(f3, *r.ko.Spec.Name)
		res.RepositoryNames = f3
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeList, "r.ko", "res", 1),
	)
}

func TestSetSDK_IAM_User_NewPath(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "iam", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-user-newpath.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "User")
	require.NotNil(crd)
	expected := `
	if r.ko.Spec.Path != nil {
		res.NewPath = r.ko.Spec.Path
	}
	if r.ko.Spec.UserName != nil {
		res.UserName = r.ko.Spec.UserName
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeUpdate, "r.ko", "res", 1),
	)
}

// func TestSetSDK_Lambda_Ignore_Code_SHA256(t *testing.T) {
// 	assert := assert.New(t)
// 	require := require.New(t)

// 	g := testutil.NewModelForServiceWithOptions(t, "lambda", &testutil.TestingModelOptions{
// 		GeneratorConfigFile: "generator-lambda-ignore-code-sha256.yaml",
// 	})

// 	crd := testutil.GetCRDByName(t, g, "Function")
// 	require.NotNil(crd)

// 	expected := `
// 	if r.ko.Spec.Code != nil {
// 		f0 := &svcsdk.FunctionCode{}
// 		if r.ko.Spec.Code.ImageURI != nil {
// 			f0.SetImageUri(*r.ko.Spec.Code.ImageURI)
// 		}
// 		if r.ko.Spec.Code.S3Bucket != nil {
// 			f0.SetS3Bucket(*r.ko.Spec.Code.S3Bucket)
// 		}
// 		if r.ko.Spec.Code.S3Key != nil {
// 			f0.SetS3Key(*r.ko.Spec.Code.S3Key)
// 		}
// 		if r.ko.Spec.Code.S3ObjectVersion != nil {
// 			f0.SetS3ObjectVersion(*r.ko.Spec.Code.S3ObjectVersion)
// 		}
// 		if r.ko.Spec.Code.ZipFile != nil {
// 			f0.SetZipFile(r.ko.Spec.Code.ZipFile)
// 		}
// 		res.SetCode(f0)
// 	}
// 	if r.ko.Spec.CodeSigningConfigARN != nil {
// 		res.SetCodeSigningConfigArn(*r.ko.Spec.CodeSigningConfigARN)
// 	}
// 	if r.ko.Spec.DeadLetterConfig != nil {
// 		f2 := &svcsdk.DeadLetterConfig{}
// 		if r.ko.Spec.DeadLetterConfig.TargetARN != nil {
// 			f2.SetTargetArn(*r.ko.Spec.DeadLetterConfig.TargetARN)
// 		}
// 		res.SetDeadLetterConfig(f2)
// 	}
// 	if r.ko.Spec.Description != nil {
// 		res.SetDescription(*r.ko.Spec.Description)
// 	}
// 	if r.ko.Spec.Environment != nil {
// 		f4 := &svcsdk.Environment{}
// 		if r.ko.Spec.Environment.Variables != nil {
// 			f4f0 := map[string]*string{}
// 			for f4f0key, f4f0valiter := range r.ko.Spec.Environment.Variables {
// 				var f4f0val string
// 				f4f0val = *f4f0valiter
// 				f4f0[f4f0key] = &f4f0val
// 			}
// 			f4.SetVariables(f4f0)
// 		}
// 		res.SetEnvironment(f4)
// 	}
// 	if r.ko.Spec.FileSystemConfigs != nil {
// 		f5 := []*svcsdk.FileSystemConfig{}
// 		for _, f5iter := range r.ko.Spec.FileSystemConfigs {
// 			f5elem := &svcsdk.FileSystemConfig{}
// 			if f5iter.ARN != nil {
// 				f5elem.SetArn(*f5iter.ARN)
// 			}
// 			if f5iter.LocalMountPath != nil {
// 				f5elem.SetLocalMountPath(*f5iter.LocalMountPath)
// 			}
// 			f5 = append(f5, f5elem)
// 		}
// 		res.SetFileSystemConfigs(f5)
// 	}
// 	if r.ko.Spec.FunctionName != nil {
// 		res.SetFunctionName(*r.ko.Spec.FunctionName)
// 	}
// 	if r.ko.Spec.Handler != nil {
// 		res.SetHandler(*r.ko.Spec.Handler)
// 	}
// 	if r.ko.Spec.ImageConfig != nil {
// 		f8 := &svcsdk.ImageConfig{}
// 		if r.ko.Spec.ImageConfig.Command != nil {
// 			f8f0 := []*string{}
// 			for _, f8f0iter := range r.ko.Spec.ImageConfig.Command {
// 				var f8f0elem string
// 				f8f0elem = *f8f0iter
// 				f8f0 = append(f8f0, &f8f0elem)
// 			}
// 			f8.SetCommand(f8f0)
// 		}
// 		if r.ko.Spec.ImageConfig.EntryPoint != nil {
// 			f8f1 := []*string{}
// 			for _, f8f1iter := range r.ko.Spec.ImageConfig.EntryPoint {
// 				var f8f1elem string
// 				f8f1elem = *f8f1iter
// 				f8f1 = append(f8f1, &f8f1elem)
// 			}
// 			f8.SetEntryPoint(f8f1)
// 		}
// 		if r.ko.Spec.ImageConfig.WorkingDirectory != nil {
// 			f8.SetWorkingDirectory(*r.ko.Spec.ImageConfig.WorkingDirectory)
// 		}
// 		res.SetImageConfig(f8)
// 	}
// 	if r.ko.Spec.KMSKeyARN != nil {
// 		res.SetKMSKeyArn(*r.ko.Spec.KMSKeyARN)
// 	}
// 	if r.ko.Spec.Layers != nil {
// 		f10 := []*string{}
// 		for _, f10iter := range r.ko.Spec.Layers {
// 			var f10elem string
// 			f10elem = *f10iter
// 			f10 = append(f10, &f10elem)
// 		}
// 		res.SetLayers(f10)
// 	}
// 	if r.ko.Spec.MemorySize != nil {
// 		res.SetMemorySize(*r.ko.Spec.MemorySize)
// 	}
// 	if r.ko.Spec.PackageType != nil {
// 		res.SetPackageType(*r.ko.Spec.PackageType)
// 	}
// 	if r.ko.Spec.Publish != nil {
// 		res.SetPublish(*r.ko.Spec.Publish)
// 	}
// 	if r.ko.Spec.Role != nil {
// 		res.SetRole(*r.ko.Spec.Role)
// 	}
// 	if r.ko.Spec.Runtime != nil {
// 		res.SetRuntime(*r.ko.Spec.Runtime)
// 	}
// 	if r.ko.Spec.Tags != nil {
// 		f16 := map[string]*string{}
// 		for f16key, f16valiter := range r.ko.Spec.Tags {
// 			var f16val string
// 			f16val = *f16valiter
// 			f16[f16key] = &f16val
// 		}
// 		res.SetTags(f16)
// 	}
// 	if r.ko.Spec.Timeout != nil {
// 		res.SetTimeout(*r.ko.Spec.Timeout)
// 	}
// 	if r.ko.Spec.TracingConfig != nil {
// 		f18 := &svcsdk.TracingConfig{}
// 		if r.ko.Spec.TracingConfig.Mode != nil {
// 			f18.SetMode(*r.ko.Spec.TracingConfig.Mode)
// 		}
// 		res.SetTracingConfig(f18)
// 	}
// 	if r.ko.Spec.VPCConfig != nil {
// 		f19 := &svcsdk.VpcConfig{}
// 		if r.ko.Spec.VPCConfig.SecurityGroupIDs != nil {
// 			f19f0 := []*string{}
// 			for _, f19f0iter := range r.ko.Spec.VPCConfig.SecurityGroupIDs {
// 				var f19f0elem string
// 				f19f0elem = *f19f0iter
// 				f19f0 = append(f19f0, &f19f0elem)
// 			}
// 			f19.SetSecurityGroupIds(f19f0)
// 		}
// 		if r.ko.Spec.VPCConfig.SubnetIDs != nil {
// 			f19f1 := []*string{}
// 			for _, f19f1iter := range r.ko.Spec.VPCConfig.SubnetIDs {
// 				var f19f1elem string
// 				f19f1elem = *f19f1iter
// 				f19f1 = append(f19f1, &f19f1elem)
// 			}
// 			f19.SetSubnetIds(f19f1)
// 		}
// 		res.SetVpcConfig(f19)
// 	}
// `
// 	assert.Equal(
// 		expected,
// 		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
// 	)
// }

func TestSetSDK_WAFv2_RuleGroup_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "wafv2", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "RuleGroup")
	require.NotNil(crd)
	expected := `
	if r.ko.Spec.Capacity != nil {
		res.Capacity = r.ko.Spec.Capacity
	}
	if r.ko.Spec.CustomResponseBodies != nil {
		f1 := map[string]svcsdktypes.CustomResponseBody{}
		for f1key, f1valiter := range r.ko.Spec.CustomResponseBodies {
			f1val := &svcsdktypes.CustomResponseBody{}
			if f1valiter.Content != nil {
				f1val.Content = f1valiter.Content
			}
			if f1valiter.ContentType != nil {
				f1val.ContentType = svcsdktypes.ResponseContentType(*f1valiter.ContentType)
			}
			f1[f1key] = *f1val
		}
		res.CustomResponseBodies = f1
	}
	if r.ko.Spec.Description != nil {
		res.Description = r.ko.Spec.Description
	}
	if r.ko.Spec.Name != nil {
		res.Name = r.ko.Spec.Name
	}
	if r.ko.Spec.Rules != nil {
		f4 := []svcsdktypes.Rule{}
		for _, f4iter := range r.ko.Spec.Rules {
			f4elem := &svcsdktypes.Rule{}
			if f4iter.Action != nil {
				f4elemf0 := &svcsdktypes.RuleAction{}
				if f4iter.Action.Allow != nil {
					f4elemf0f0 := &svcsdktypes.AllowAction{}
					if f4iter.Action.Allow.CustomRequestHandling != nil {
						f4elemf0f0f0 := &svcsdktypes.CustomRequestHandling{}
						if f4iter.Action.Allow.CustomRequestHandling.InsertHeaders != nil {
							f4elemf0f0f0f0 := []svcsdktypes.CustomHTTPHeader{}
							for _, f4elemf0f0f0f0iter := range f4iter.Action.Allow.CustomRequestHandling.InsertHeaders {
								f4elemf0f0f0f0elem := &svcsdktypes.CustomHTTPHeader{}
								if f4elemf0f0f0f0iter.Name != nil {
									f4elemf0f0f0f0elem.Name = f4elemf0f0f0f0iter.Name
								}
								if f4elemf0f0f0f0iter.Value != nil {
									f4elemf0f0f0f0elem.Value = f4elemf0f0f0f0iter.Value
								}
								f4elemf0f0f0f0 = append(f4elemf0f0f0f0, *f4elemf0f0f0f0elem)
							}
							f4elemf0f0f0.InsertHeaders = f4elemf0f0f0f0
						}
						f4elemf0f0.CustomRequestHandling = f4elemf0f0f0
					}
					f4elemf0.Allow = f4elemf0f0
				}
				if f4iter.Action.Block != nil {
					f4elemf0f1 := &svcsdktypes.BlockAction{}
					if f4iter.Action.Block.CustomResponse != nil {
						f4elemf0f1f0 := &svcsdktypes.CustomResponse{}
						if f4iter.Action.Block.CustomResponse.CustomResponseBodyKey != nil {
							f4elemf0f1f0.CustomResponseBodyKey = f4iter.Action.Block.CustomResponse.CustomResponseBodyKey
						}
						if f4iter.Action.Block.CustomResponse.ResponseCode != nil {
							if *f4iter.Action.Block.CustomResponse.ResponseCode > math.MaxInt32 || *f4iter.Action.Block.CustomResponse.ResponseCode < math.MinInt32 {
								return nil, fmt.Errorf("error: field ResponseCode is of type int32")
							}
							responseCodeCopy := int32(*f4iter.Action.Block.CustomResponse.ResponseCode)
							f4elemf0f1f0.ResponseCode = &responseCodeCopy
						}
						if f4iter.Action.Block.CustomResponse.ResponseHeaders != nil {
							f4elemf0f1f0f2 := []svcsdktypes.CustomHTTPHeader{}
							for _, f4elemf0f1f0f2iter := range f4iter.Action.Block.CustomResponse.ResponseHeaders {
								f4elemf0f1f0f2elem := &svcsdktypes.CustomHTTPHeader{}
								if f4elemf0f1f0f2iter.Name != nil {
									f4elemf0f1f0f2elem.Name = f4elemf0f1f0f2iter.Name
								}
								if f4elemf0f1f0f2iter.Value != nil {
									f4elemf0f1f0f2elem.Value = f4elemf0f1f0f2iter.Value
								}
								f4elemf0f1f0f2 = append(f4elemf0f1f0f2, *f4elemf0f1f0f2elem)
							}
							f4elemf0f1f0.ResponseHeaders = f4elemf0f1f0f2
						}
						f4elemf0f1.CustomResponse = f4elemf0f1f0
					}
					f4elemf0.Block = f4elemf0f1
				}
				if f4iter.Action.Captcha != nil {
					f4elemf0f2 := &svcsdktypes.CaptchaAction{}
					if f4iter.Action.Captcha.CustomRequestHandling != nil {
						f4elemf0f2f0 := &svcsdktypes.CustomRequestHandling{}
						if f4iter.Action.Captcha.CustomRequestHandling.InsertHeaders != nil {
							f4elemf0f2f0f0 := []svcsdktypes.CustomHTTPHeader{}
							for _, f4elemf0f2f0f0iter := range f4iter.Action.Captcha.CustomRequestHandling.InsertHeaders {
								f4elemf0f2f0f0elem := &svcsdktypes.CustomHTTPHeader{}
								if f4elemf0f2f0f0iter.Name != nil {
									f4elemf0f2f0f0elem.Name = f4elemf0f2f0f0iter.Name
								}
								if f4elemf0f2f0f0iter.Value != nil {
									f4elemf0f2f0f0elem.Value = f4elemf0f2f0f0iter.Value
								}
								f4elemf0f2f0f0 = append(f4elemf0f2f0f0, *f4elemf0f2f0f0elem)
							}
							f4elemf0f2f0.InsertHeaders = f4elemf0f2f0f0
						}
						f4elemf0f2.CustomRequestHandling = f4elemf0f2f0
					}
					f4elemf0.Captcha = f4elemf0f2
				}
				if f4iter.Action.Challenge != nil {
					f4elemf0f3 := &svcsdktypes.ChallengeAction{}
					if f4iter.Action.Challenge.CustomRequestHandling != nil {
						f4elemf0f3f0 := &svcsdktypes.CustomRequestHandling{}
						if f4iter.Action.Challenge.CustomRequestHandling.InsertHeaders != nil {
							f4elemf0f3f0f0 := []svcsdktypes.CustomHTTPHeader{}
							for _, f4elemf0f3f0f0iter := range f4iter.Action.Challenge.CustomRequestHandling.InsertHeaders {
								f4elemf0f3f0f0elem := &svcsdktypes.CustomHTTPHeader{}
								if f4elemf0f3f0f0iter.Name != nil {
									f4elemf0f3f0f0elem.Name = f4elemf0f3f0f0iter.Name
								}
								if f4elemf0f3f0f0iter.Value != nil {
									f4elemf0f3f0f0elem.Value = f4elemf0f3f0f0iter.Value
								}
								f4elemf0f3f0f0 = append(f4elemf0f3f0f0, *f4elemf0f3f0f0elem)
							}
							f4elemf0f3f0.InsertHeaders = f4elemf0f3f0f0
						}
						f4elemf0f3.CustomRequestHandling = f4elemf0f3f0
					}
					f4elemf0.Challenge = f4elemf0f3
				}
				if f4iter.Action.Count != nil {
					f4elemf0f4 := &svcsdktypes.CountAction{}
					if f4iter.Action.Count.CustomRequestHandling != nil {
						f4elemf0f4f0 := &svcsdktypes.CustomRequestHandling{}
						if f4iter.Action.Count.CustomRequestHandling.InsertHeaders != nil {
							f4elemf0f4f0f0 := []svcsdktypes.CustomHTTPHeader{}
							for _, f4elemf0f4f0f0iter := range f4iter.Action.Count.CustomRequestHandling.InsertHeaders {
								f4elemf0f4f0f0elem := &svcsdktypes.CustomHTTPHeader{}
								if f4elemf0f4f0f0iter.Name != nil {
									f4elemf0f4f0f0elem.Name = f4elemf0f4f0f0iter.Name
								}
								if f4elemf0f4f0f0iter.Value != nil {
									f4elemf0f4f0f0elem.Value = f4elemf0f4f0f0iter.Value
								}
								f4elemf0f4f0f0 = append(f4elemf0f4f0f0, *f4elemf0f4f0f0elem)
							}
							f4elemf0f4f0.InsertHeaders = f4elemf0f4f0f0
						}
						f4elemf0f4.CustomRequestHandling = f4elemf0f4f0
					}
					f4elemf0.Count = f4elemf0f4
				}
				f4elem.Action = f4elemf0
			}
			if f4iter.CaptchaConfig != nil {
				f4elemf1 := &svcsdktypes.CaptchaConfig{}
				if f4iter.CaptchaConfig.ImmunityTimeProperty != nil {
					f4elemf1f0 := &svcsdktypes.ImmunityTimeProperty{}
					if f4iter.CaptchaConfig.ImmunityTimeProperty.ImmunityTime != nil {
						f4elemf1f0.ImmunityTime = f4iter.CaptchaConfig.ImmunityTimeProperty.ImmunityTime
					}
					f4elemf1.ImmunityTimeProperty = f4elemf1f0
				}
				f4elem.CaptchaConfig = f4elemf1
			}
			if f4iter.ChallengeConfig != nil {
				f4elemf2 := &svcsdktypes.ChallengeConfig{}
				if f4iter.ChallengeConfig.ImmunityTimeProperty != nil {
					f4elemf2f0 := &svcsdktypes.ImmunityTimeProperty{}
					if f4iter.ChallengeConfig.ImmunityTimeProperty.ImmunityTime != nil {
						f4elemf2f0.ImmunityTime = f4iter.ChallengeConfig.ImmunityTimeProperty.ImmunityTime
					}
					f4elemf2.ImmunityTimeProperty = f4elemf2f0
				}
				f4elem.ChallengeConfig = f4elemf2
			}
			if f4iter.Name != nil {
				f4elem.Name = f4iter.Name
			}
			if f4iter.OverrideAction != nil {
				f4elemf4 := &svcsdktypes.OverrideAction{}
				if f4iter.OverrideAction.Count != nil {
					f4elemf4f0 := &svcsdktypes.CountAction{}
					if f4iter.OverrideAction.Count.CustomRequestHandling != nil {
						f4elemf4f0f0 := &svcsdktypes.CustomRequestHandling{}
						if f4iter.OverrideAction.Count.CustomRequestHandling.InsertHeaders != nil {
							f4elemf4f0f0f0 := []svcsdktypes.CustomHTTPHeader{}
							for _, f4elemf4f0f0f0iter := range f4iter.OverrideAction.Count.CustomRequestHandling.InsertHeaders {
								f4elemf4f0f0f0elem := &svcsdktypes.CustomHTTPHeader{}
								if f4elemf4f0f0f0iter.Name != nil {
									f4elemf4f0f0f0elem.Name = f4elemf4f0f0f0iter.Name
								}
								if f4elemf4f0f0f0iter.Value != nil {
									f4elemf4f0f0f0elem.Value = f4elemf4f0f0f0iter.Value
								}
								f4elemf4f0f0f0 = append(f4elemf4f0f0f0, *f4elemf4f0f0f0elem)
							}
							f4elemf4f0f0.InsertHeaders = f4elemf4f0f0f0
						}
						f4elemf4f0.CustomRequestHandling = f4elemf4f0f0
					}
					f4elemf4.Count = f4elemf4f0
				}
				if f4iter.OverrideAction.None != nil {
					f4elemf4f1 := &svcsdktypes.NoneAction{}
					f4elemf4.None = f4elemf4f1
				}
				f4elem.OverrideAction = f4elemf4
			}
			if f4iter.Priority != nil {
				if *f4iter.Priority > math.MaxInt32 || *f4iter.Priority < math.MinInt32 {
					return nil, fmt.Errorf("error: field Priority is of type int32")
				}
				priorityCopy := int32(*f4iter.Priority)
				f4elem.Priority = priorityCopy
			}
			if f4iter.RuleLabels != nil {
				f4elemf6 := []svcsdktypes.Label{}
				for _, f4elemf6iter := range f4iter.RuleLabels {
					f4elemf6elem := &svcsdktypes.Label{}
					if f4elemf6iter.Name != nil {
						f4elemf6elem.Name = f4elemf6iter.Name
					}
					f4elemf6 = append(f4elemf6, *f4elemf6elem)
				}
				f4elem.RuleLabels = f4elemf6
			}
			if f4iter.Statement != nil {
				f4elemf7 := &svcsdktypes.Statement{}
				if f4iter.Statement.ByteMatchStatement != nil {
					f4elemf7f1 := &svcsdktypes.ByteMatchStatement{}
					if f4iter.Statement.ByteMatchStatement.FieldToMatch != nil {
						f4elemf7f1f0 := &svcsdktypes.FieldToMatch{}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.AllQueryArguments != nil {
							f4elemf7f1f0f0 := &svcsdktypes.AllQueryArguments{}
							f4elemf7f1f0.AllQueryArguments = f4elemf7f1f0f0
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.Body != nil {
							f4elemf7f1f0f1 := &svcsdktypes.Body{}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.Body.OversizeHandling != nil {
								f4elemf7f1f0f1.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.ByteMatchStatement.FieldToMatch.Body.OversizeHandling)
							}
							f4elemf7f1f0.Body = f4elemf7f1f0f1
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies != nil {
							f4elemf7f1f0f2 := &svcsdktypes.Cookies{}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f4elemf7f1f0f2f0 := &svcsdktypes.CookieMatchPattern{}
								if f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f4elemf7f1f0f2f0f0 := &svcsdktypes.All{}
									f4elemf7f1f0f2f0.All = f4elemf7f1f0f2f0f0
								}
								if f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f4elemf7f1f0f2f0.ExcludedCookies = aws.ToStringSlice(f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f4elemf7f1f0f2f0.IncludedCookies = aws.ToStringSlice(f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f4elemf7f1f0f2.MatchPattern = f4elemf7f1f0f2f0
							}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchScope != nil {
								f4elemf7f1f0f2.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchScope)
							}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.OversizeHandling != nil {
								f4elemf7f1f0f2.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.OversizeHandling)
							}
							f4elemf7f1f0.Cookies = f4elemf7f1f0f2
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.HeaderOrder != nil {
							f4elemf7f1f0f3 := &svcsdktypes.HeaderOrder{}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling != nil {
								f4elemf7f1f0f3.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.ByteMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling)
							}
							f4elemf7f1f0.HeaderOrder = f4elemf7f1f0f3
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers != nil {
							f4elemf7f1f0f4 := &svcsdktypes.Headers{}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern != nil {
								f4elemf7f1f0f4f0 := &svcsdktypes.HeaderMatchPattern{}
								if f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f4elemf7f1f0f4f0f0 := &svcsdktypes.All{}
									f4elemf7f1f0f4f0.All = f4elemf7f1f0f4f0f0
								}
								if f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f4elemf7f1f0f4f0.ExcludedHeaders = aws.ToStringSlice(f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f4elemf7f1f0f4f0.IncludedHeaders = aws.ToStringSlice(f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f4elemf7f1f0f4.MatchPattern = f4elemf7f1f0f4f0
							}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchScope != nil {
								f4elemf7f1f0f4.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchScope)
							}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers.OversizeHandling != nil {
								f4elemf7f1f0f4.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.ByteMatchStatement.FieldToMatch.Headers.OversizeHandling)
							}
							f4elemf7f1f0.Headers = f4elemf7f1f0f4
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.JA3Fingerprint != nil {
							f4elemf7f1f0f5 := &svcsdktypes.JA3Fingerprint{}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != nil {
								f4elemf7f1f0f5.FallbackBehavior = svcsdktypes.FallbackBehavior(*f4iter.Statement.ByteMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior)
							}
							f4elemf7f1f0.JA3Fingerprint = f4elemf7f1f0f5
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody != nil {
							f4elemf7f1f0f6 := &svcsdktypes.JsonBody{}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior != nil {
								f4elemf7f1f0f6.InvalidFallbackBehavior = svcsdktypes.BodyParsingFallbackBehavior(*f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior)
							}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody.MatchPattern != nil {
								f4elemf7f1f0f6f1 := &svcsdktypes.JsonMatchPattern{}
								if f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody.MatchPattern.All != nil {
									f4elemf7f1f0f6f1f0 := &svcsdktypes.All{}
									f4elemf7f1f0f6f1.All = f4elemf7f1f0f6f1f0
								}
								if f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths != nil {
									f4elemf7f1f0f6f1.IncludedPaths = aws.ToStringSlice(f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths)
								}
								f4elemf7f1f0f6.MatchPattern = f4elemf7f1f0f6f1
							}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody.MatchScope != nil {
								f4elemf7f1f0f6.MatchScope = svcsdktypes.JsonMatchScope(*f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody.MatchScope)
							}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody.OversizeHandling != nil {
								f4elemf7f1f0f6.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.ByteMatchStatement.FieldToMatch.JSONBody.OversizeHandling)
							}
							f4elemf7f1f0.JsonBody = f4elemf7f1f0f6
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.Method != nil {
							f4elemf7f1f0f7 := &svcsdktypes.Method{}
							f4elemf7f1f0.Method = f4elemf7f1f0f7
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.QueryString != nil {
							f4elemf7f1f0f8 := &svcsdktypes.QueryString{}
							f4elemf7f1f0.QueryString = f4elemf7f1f0f8
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.SingleHeader != nil {
							f4elemf7f1f0f9 := &svcsdktypes.SingleHeader{}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.SingleHeader.Name != nil {
								f4elemf7f1f0f9.Name = f4iter.Statement.ByteMatchStatement.FieldToMatch.SingleHeader.Name
							}
							f4elemf7f1f0.SingleHeader = f4elemf7f1f0f9
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.SingleQueryArgument != nil {
							f4elemf7f1f0f10 := &svcsdktypes.SingleQueryArgument{}
							if f4iter.Statement.ByteMatchStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f4elemf7f1f0f10.Name = f4iter.Statement.ByteMatchStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f4elemf7f1f0.SingleQueryArgument = f4elemf7f1f0f10
						}
						if f4iter.Statement.ByteMatchStatement.FieldToMatch.URIPath != nil {
							f4elemf7f1f0f11 := &svcsdktypes.UriPath{}
							f4elemf7f1f0.UriPath = f4elemf7f1f0f11
						}
						f4elemf7f1.FieldToMatch = f4elemf7f1f0
					}
					if f4iter.Statement.ByteMatchStatement.PositionalConstraint != nil {
						f4elemf7f1.PositionalConstraint = svcsdktypes.PositionalConstraint(*f4iter.Statement.ByteMatchStatement.PositionalConstraint)
					}
					if f4iter.Statement.ByteMatchStatement.SearchString != nil {
						f4elemf7f1.SearchString = f4iter.Statement.ByteMatchStatement.SearchString
					}
					if f4iter.Statement.ByteMatchStatement.TextTransformations != nil {
						f4elemf7f1f3 := []svcsdktypes.TextTransformation{}
						for _, f4elemf7f1f3iter := range f4iter.Statement.ByteMatchStatement.TextTransformations {
							f4elemf7f1f3elem := &svcsdktypes.TextTransformation{}
							if f4elemf7f1f3iter.Priority != nil {
								if *f4elemf7f1f3iter.Priority > math.MaxInt32 || *f4elemf7f1f3iter.Priority < math.MinInt32 {
									return nil, fmt.Errorf("error: field Priority is of type int32")
								}
								priorityCopy := int32(*f4elemf7f1f3iter.Priority)
								f4elemf7f1f3elem.Priority = priorityCopy
							}
							if f4elemf7f1f3iter.Type != nil {
								f4elemf7f1f3elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f1f3iter.Type)
							}
							f4elemf7f1f3 = append(f4elemf7f1f3, *f4elemf7f1f3elem)
						}
						f4elemf7f1.TextTransformations = f4elemf7f1f3
					}
					f4elemf7.ByteMatchStatement = f4elemf7f1
				}
				if f4iter.Statement.GeoMatchStatement != nil {
					f4elemf7f2 := &svcsdktypes.GeoMatchStatement{}
					if f4iter.Statement.GeoMatchStatement.CountryCodes != nil {
						f4elemf7f2f0 := []svcsdktypes.CountryCode{}
						for _, f4elemf7f2f0iter := range f4iter.Statement.GeoMatchStatement.CountryCodes {
							var f4elemf7f2f0elem string
							f4elemf7f2f0elem = string(*f4elemf7f2f0iter)
							f4elemf7f2f0 = append(f4elemf7f2f0, svcsdktypes.CountryCode(f4elemf7f2f0elem))
						}
						f4elemf7f2.CountryCodes = f4elemf7f2f0
					}
					if f4iter.Statement.GeoMatchStatement.ForwardedIPConfig != nil {
						f4elemf7f2f1 := &svcsdktypes.ForwardedIPConfig{}
						if f4iter.Statement.GeoMatchStatement.ForwardedIPConfig.FallbackBehavior != nil {
							f4elemf7f2f1.FallbackBehavior = svcsdktypes.FallbackBehavior(*f4iter.Statement.GeoMatchStatement.ForwardedIPConfig.FallbackBehavior)
						}
						if f4iter.Statement.GeoMatchStatement.ForwardedIPConfig.HeaderName != nil {
							f4elemf7f2f1.HeaderName = f4iter.Statement.GeoMatchStatement.ForwardedIPConfig.HeaderName
						}
						f4elemf7f2.ForwardedIPConfig = f4elemf7f2f1
					}
					f4elemf7.GeoMatchStatement = f4elemf7f2
				}
				if f4iter.Statement.IPSetReferenceStatement != nil {
					f4elemf7f3 := &svcsdktypes.IPSetReferenceStatement{}
					if f4iter.Statement.IPSetReferenceStatement.ARN != nil {
						f4elemf7f3.ARN = f4iter.Statement.IPSetReferenceStatement.ARN
					}
					if f4iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig != nil {
						f4elemf7f3f1 := &svcsdktypes.IPSetForwardedIPConfig{}
						if f4iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.FallbackBehavior != nil {
							f4elemf7f3f1.FallbackBehavior = svcsdktypes.FallbackBehavior(*f4iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.FallbackBehavior)
						}
						if f4iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.HeaderName != nil {
							f4elemf7f3f1.HeaderName = f4iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.HeaderName
						}
						if f4iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.Position != nil {
							f4elemf7f3f1.Position = svcsdktypes.ForwardedIPPosition(*f4iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.Position)
						}
						f4elemf7f3.IPSetForwardedIPConfig = f4elemf7f3f1
					}
					f4elemf7.IPSetReferenceStatement = f4elemf7f3
				}
				if f4iter.Statement.LabelMatchStatement != nil {
					f4elemf7f4 := &svcsdktypes.LabelMatchStatement{}
					if f4iter.Statement.LabelMatchStatement.Key != nil {
						f4elemf7f4.Key = f4iter.Statement.LabelMatchStatement.Key
					}
					if f4iter.Statement.LabelMatchStatement.Scope != nil {
						f4elemf7f4.Scope = svcsdktypes.LabelMatchScope(*f4iter.Statement.LabelMatchStatement.Scope)
					}
					f4elemf7.LabelMatchStatement = f4elemf7f4
				}
				if f4iter.Statement.ManagedRuleGroupStatement != nil {
					f4elemf7f5 := &svcsdktypes.ManagedRuleGroupStatement{}
					if f4iter.Statement.ManagedRuleGroupStatement.ExcludedRules != nil {
						f4elemf7f5f0 := []svcsdktypes.ExcludedRule{}
						for _, f4elemf7f5f0iter := range f4iter.Statement.ManagedRuleGroupStatement.ExcludedRules {
							f4elemf7f5f0elem := &svcsdktypes.ExcludedRule{}
							if f4elemf7f5f0iter.Name != nil {
								f4elemf7f5f0elem.Name = f4elemf7f5f0iter.Name
							}
							f4elemf7f5f0 = append(f4elemf7f5f0, *f4elemf7f5f0elem)
						}
						f4elemf7f5.ExcludedRules = f4elemf7f5f0
					}
					if f4iter.Statement.ManagedRuleGroupStatement.ManagedRuleGroupConfigs != nil {
						f4elemf7f5f1 := []svcsdktypes.ManagedRuleGroupConfig{}
						for _, f4elemf7f5f1iter := range f4iter.Statement.ManagedRuleGroupStatement.ManagedRuleGroupConfigs {
							f4elemf7f5f1elem := &svcsdktypes.ManagedRuleGroupConfig{}
							if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet != nil {
								f4elemf7f5f1elemf0 := &svcsdktypes.AWSManagedRulesACFPRuleSet{}
								if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.CreationPath != nil {
									f4elemf7f5f1elemf0.CreationPath = f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.CreationPath
								}
								if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.EnableRegexInPath != nil {
									f4elemf7f5f1elemf0.EnableRegexInPath = *f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.EnableRegexInPath
								}
								if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RegistrationPagePath != nil {
									f4elemf7f5f1elemf0.RegistrationPagePath = f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RegistrationPagePath
								}
								if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection != nil {
									f4elemf7f5f1elemf0f3 := &svcsdktypes.RequestInspectionACFP{}
									if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.AddressFields != nil {
										f4elemf7f5f1elemf0f3f0 := []svcsdktypes.AddressField{}
										for _, f4elemf7f5f1elemf0f3f0iter := range f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.AddressFields {
											f4elemf7f5f1elemf0f3f0elem := &svcsdktypes.AddressField{}
											if f4elemf7f5f1elemf0f3f0iter.Identifier != nil {
												f4elemf7f5f1elemf0f3f0elem.Identifier = f4elemf7f5f1elemf0f3f0iter.Identifier
											}
											f4elemf7f5f1elemf0f3f0 = append(f4elemf7f5f1elemf0f3f0, *f4elemf7f5f1elemf0f3f0elem)
										}
										f4elemf7f5f1elemf0f3.AddressFields = f4elemf7f5f1elemf0f3f0
									}
									if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.EmailField != nil {
										f4elemf7f5f1elemf0f3f1 := &svcsdktypes.EmailField{}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.EmailField.Identifier != nil {
											f4elemf7f5f1elemf0f3f1.Identifier = f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.EmailField.Identifier
										}
										f4elemf7f5f1elemf0f3.EmailField = f4elemf7f5f1elemf0f3f1
									}
									if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PasswordField != nil {
										f4elemf7f5f1elemf0f3f2 := &svcsdktypes.PasswordField{}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PasswordField.Identifier != nil {
											f4elemf7f5f1elemf0f3f2.Identifier = f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PasswordField.Identifier
										}
										f4elemf7f5f1elemf0f3.PasswordField = f4elemf7f5f1elemf0f3f2
									}
									if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PayloadType != nil {
										f4elemf7f5f1elemf0f3.PayloadType = svcsdktypes.PayloadType(*f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PayloadType)
									}
									if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PhoneNumberFields != nil {
										f4elemf7f5f1elemf0f3f4 := []svcsdktypes.PhoneNumberField{}
										for _, f4elemf7f5f1elemf0f3f4iter := range f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PhoneNumberFields {
											f4elemf7f5f1elemf0f3f4elem := &svcsdktypes.PhoneNumberField{}
											if f4elemf7f5f1elemf0f3f4iter.Identifier != nil {
												f4elemf7f5f1elemf0f3f4elem.Identifier = f4elemf7f5f1elemf0f3f4iter.Identifier
											}
											f4elemf7f5f1elemf0f3f4 = append(f4elemf7f5f1elemf0f3f4, *f4elemf7f5f1elemf0f3f4elem)
										}
										f4elemf7f5f1elemf0f3.PhoneNumberFields = f4elemf7f5f1elemf0f3f4
									}
									if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.UsernameField != nil {
										f4elemf7f5f1elemf0f3f5 := &svcsdktypes.UsernameField{}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.UsernameField.Identifier != nil {
											f4elemf7f5f1elemf0f3f5.Identifier = f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.UsernameField.Identifier
										}
										f4elemf7f5f1elemf0f3.UsernameField = f4elemf7f5f1elemf0f3f5
									}
									f4elemf7f5f1elemf0.RequestInspection = f4elemf7f5f1elemf0f3
								}
								if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection != nil {
									f4elemf7f5f1elemf0f4 := &svcsdktypes.ResponseInspection{}
									if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.BodyContains != nil {
										f4elemf7f5f1elemf0f4f0 := &svcsdktypes.ResponseInspectionBodyContains{}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.BodyContains.FailureStrings != nil {
											f4elemf7f5f1elemf0f4f0.FailureStrings = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.BodyContains.FailureStrings)
										}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.BodyContains.SuccessStrings != nil {
											f4elemf7f5f1elemf0f4f0.SuccessStrings = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.BodyContains.SuccessStrings)
										}
										f4elemf7f5f1elemf0f4.BodyContains = f4elemf7f5f1elemf0f4f0
									}
									if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header != nil {
										f4elemf7f5f1elemf0f4f1 := &svcsdktypes.ResponseInspectionHeader{}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.FailureValues != nil {
											f4elemf7f5f1elemf0f4f1.FailureValues = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.FailureValues)
										}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.Name != nil {
											f4elemf7f5f1elemf0f4f1.Name = f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.Name
										}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.SuccessValues != nil {
											f4elemf7f5f1elemf0f4f1.SuccessValues = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.SuccessValues)
										}
										f4elemf7f5f1elemf0f4.Header = f4elemf7f5f1elemf0f4f1
									}
									if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.JSON != nil {
										f4elemf7f5f1elemf0f4f2 := &svcsdktypes.ResponseInspectionJson{}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.JSON.FailureValues != nil {
											f4elemf7f5f1elemf0f4f2.FailureValues = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.JSON.FailureValues)
										}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.JSON.Identifier != nil {
											f4elemf7f5f1elemf0f4f2.Identifier = f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.JSON.Identifier
										}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.JSON.SuccessValues != nil {
											f4elemf7f5f1elemf0f4f2.SuccessValues = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.JSON.SuccessValues)
										}
										f4elemf7f5f1elemf0f4.Json = f4elemf7f5f1elemf0f4f2
									}
									if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.StatusCode != nil {
										f4elemf7f5f1elemf0f4f3 := &svcsdktypes.ResponseInspectionStatusCode{}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.StatusCode.FailureCodes != nil {
											f4elemf7f5f1elemf0f4f3f0 := []int32{}
											for _, f4elemf7f5f1elemf0f4f3f0iter := range f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.StatusCode.FailureCodes {
												var f4elemf7f5f1elemf0f4f3f0elem int32
												if *f4elemf7f5f1elemf0f4f3f0iter > math.MaxInt32 || *f4elemf7f5f1elemf0f4f3f0iter < math.MinInt32 {
													return nil, fmt.Errorf("error: field FailureCode is of type int32")
												}
												failureCodeCopy := int32(*f4elemf7f5f1elemf0f4f3f0iter)
												f4elemf7f5f1elemf0f4f3f0elem = failureCodeCopy
												f4elemf7f5f1elemf0f4f3f0 = append(f4elemf7f5f1elemf0f4f3f0, f4elemf7f5f1elemf0f4f3f0elem)
											}
											f4elemf7f5f1elemf0f4f3.FailureCodes = f4elemf7f5f1elemf0f4f3f0
										}
										if f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.StatusCode.SuccessCodes != nil {
											f4elemf7f5f1elemf0f4f3f1 := []int32{}
											for _, f4elemf7f5f1elemf0f4f3f1iter := range f4elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.StatusCode.SuccessCodes {
												var f4elemf7f5f1elemf0f4f3f1elem int32
												if *f4elemf7f5f1elemf0f4f3f1iter > math.MaxInt32 || *f4elemf7f5f1elemf0f4f3f1iter < math.MinInt32 {
													return nil, fmt.Errorf("error: field SuccessCode is of type int32")
												}
												successCodeCopy := int32(*f4elemf7f5f1elemf0f4f3f1iter)
												f4elemf7f5f1elemf0f4f3f1elem = successCodeCopy
												f4elemf7f5f1elemf0f4f3f1 = append(f4elemf7f5f1elemf0f4f3f1, f4elemf7f5f1elemf0f4f3f1elem)
											}
											f4elemf7f5f1elemf0f4f3.SuccessCodes = f4elemf7f5f1elemf0f4f3f1
										}
										f4elemf7f5f1elemf0f4.StatusCode = f4elemf7f5f1elemf0f4f3
									}
									f4elemf7f5f1elemf0.ResponseInspection = f4elemf7f5f1elemf0f4
								}
								f4elemf7f5f1elem.AWSManagedRulesACFPRuleSet = f4elemf7f5f1elemf0
							}
							if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet != nil {
								f4elemf7f5f1elemf1 := &svcsdktypes.AWSManagedRulesATPRuleSet{}
								if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.EnableRegexInPath != nil {
									f4elemf7f5f1elemf1.EnableRegexInPath = *f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.EnableRegexInPath
								}
								if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.LoginPath != nil {
									f4elemf7f5f1elemf1.LoginPath = f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.LoginPath
								}
								if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection != nil {
									f4elemf7f5f1elemf1f2 := &svcsdktypes.RequestInspection{}
									if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.PasswordField != nil {
										f4elemf7f5f1elemf1f2f0 := &svcsdktypes.PasswordField{}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.PasswordField.Identifier != nil {
											f4elemf7f5f1elemf1f2f0.Identifier = f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.PasswordField.Identifier
										}
										f4elemf7f5f1elemf1f2.PasswordField = f4elemf7f5f1elemf1f2f0
									}
									if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.PayloadType != nil {
										f4elemf7f5f1elemf1f2.PayloadType = svcsdktypes.PayloadType(*f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.PayloadType)
									}
									if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.UsernameField != nil {
										f4elemf7f5f1elemf1f2f2 := &svcsdktypes.UsernameField{}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.UsernameField.Identifier != nil {
											f4elemf7f5f1elemf1f2f2.Identifier = f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.UsernameField.Identifier
										}
										f4elemf7f5f1elemf1f2.UsernameField = f4elemf7f5f1elemf1f2f2
									}
									f4elemf7f5f1elemf1.RequestInspection = f4elemf7f5f1elemf1f2
								}
								if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection != nil {
									f4elemf7f5f1elemf1f3 := &svcsdktypes.ResponseInspection{}
									if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.BodyContains != nil {
										f4elemf7f5f1elemf1f3f0 := &svcsdktypes.ResponseInspectionBodyContains{}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.BodyContains.FailureStrings != nil {
											f4elemf7f5f1elemf1f3f0.FailureStrings = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.BodyContains.FailureStrings)
										}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.BodyContains.SuccessStrings != nil {
											f4elemf7f5f1elemf1f3f0.SuccessStrings = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.BodyContains.SuccessStrings)
										}
										f4elemf7f5f1elemf1f3.BodyContains = f4elemf7f5f1elemf1f3f0
									}
									if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header != nil {
										f4elemf7f5f1elemf1f3f1 := &svcsdktypes.ResponseInspectionHeader{}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.FailureValues != nil {
											f4elemf7f5f1elemf1f3f1.FailureValues = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.FailureValues)
										}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.Name != nil {
											f4elemf7f5f1elemf1f3f1.Name = f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.Name
										}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.SuccessValues != nil {
											f4elemf7f5f1elemf1f3f1.SuccessValues = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.SuccessValues)
										}
										f4elemf7f5f1elemf1f3.Header = f4elemf7f5f1elemf1f3f1
									}
									if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.JSON != nil {
										f4elemf7f5f1elemf1f3f2 := &svcsdktypes.ResponseInspectionJson{}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.JSON.FailureValues != nil {
											f4elemf7f5f1elemf1f3f2.FailureValues = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.JSON.FailureValues)
										}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.JSON.Identifier != nil {
											f4elemf7f5f1elemf1f3f2.Identifier = f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.JSON.Identifier
										}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.JSON.SuccessValues != nil {
											f4elemf7f5f1elemf1f3f2.SuccessValues = aws.ToStringSlice(f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.JSON.SuccessValues)
										}
										f4elemf7f5f1elemf1f3.Json = f4elemf7f5f1elemf1f3f2
									}
									if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.StatusCode != nil {
										f4elemf7f5f1elemf1f3f3 := &svcsdktypes.ResponseInspectionStatusCode{}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.StatusCode.FailureCodes != nil {
											f4elemf7f5f1elemf1f3f3f0 := []int32{}
											for _, f4elemf7f5f1elemf1f3f3f0iter := range f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.StatusCode.FailureCodes {
												var f4elemf7f5f1elemf1f3f3f0elem int32
												if *f4elemf7f5f1elemf1f3f3f0iter > math.MaxInt32 || *f4elemf7f5f1elemf1f3f3f0iter < math.MinInt32 {
													return nil, fmt.Errorf("error: field FailureCode is of type int32")
												}
												failureCodeCopy := int32(*f4elemf7f5f1elemf1f3f3f0iter)
												f4elemf7f5f1elemf1f3f3f0elem = failureCodeCopy
												f4elemf7f5f1elemf1f3f3f0 = append(f4elemf7f5f1elemf1f3f3f0, f4elemf7f5f1elemf1f3f3f0elem)
											}
											f4elemf7f5f1elemf1f3f3.FailureCodes = f4elemf7f5f1elemf1f3f3f0
										}
										if f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.StatusCode.SuccessCodes != nil {
											f4elemf7f5f1elemf1f3f3f1 := []int32{}
											for _, f4elemf7f5f1elemf1f3f3f1iter := range f4elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.StatusCode.SuccessCodes {
												var f4elemf7f5f1elemf1f3f3f1elem int32
												if *f4elemf7f5f1elemf1f3f3f1iter > math.MaxInt32 || *f4elemf7f5f1elemf1f3f3f1iter < math.MinInt32 {
													return nil, fmt.Errorf("error: field SuccessCode is of type int32")
												}
												successCodeCopy := int32(*f4elemf7f5f1elemf1f3f3f1iter)
												f4elemf7f5f1elemf1f3f3f1elem = successCodeCopy
												f4elemf7f5f1elemf1f3f3f1 = append(f4elemf7f5f1elemf1f3f3f1, f4elemf7f5f1elemf1f3f3f1elem)
											}
											f4elemf7f5f1elemf1f3f3.SuccessCodes = f4elemf7f5f1elemf1f3f3f1
										}
										f4elemf7f5f1elemf1f3.StatusCode = f4elemf7f5f1elemf1f3f3
									}
									f4elemf7f5f1elemf1.ResponseInspection = f4elemf7f5f1elemf1f3
								}
								f4elemf7f5f1elem.AWSManagedRulesATPRuleSet = f4elemf7f5f1elemf1
							}
							if f4elemf7f5f1iter.AWSManagedRulesBotControlRuleSet != nil {
								f4elemf7f5f1elemf2 := &svcsdktypes.AWSManagedRulesBotControlRuleSet{}
								if f4elemf7f5f1iter.AWSManagedRulesBotControlRuleSet.EnableMachineLearning != nil {
									f4elemf7f5f1elemf2.EnableMachineLearning = f4elemf7f5f1iter.AWSManagedRulesBotControlRuleSet.EnableMachineLearning
								}
								if f4elemf7f5f1iter.AWSManagedRulesBotControlRuleSet.InspectionLevel != nil {
									f4elemf7f5f1elemf2.InspectionLevel = svcsdktypes.InspectionLevel(*f4elemf7f5f1iter.AWSManagedRulesBotControlRuleSet.InspectionLevel)
								}
								f4elemf7f5f1elem.AWSManagedRulesBotControlRuleSet = f4elemf7f5f1elemf2
							}
							if f4elemf7f5f1iter.LoginPath != nil {
								f4elemf7f5f1elem.LoginPath = f4elemf7f5f1iter.LoginPath
							}
							if f4elemf7f5f1iter.PasswordField != nil {
								f4elemf7f5f1elemf4 := &svcsdktypes.PasswordField{}
								if f4elemf7f5f1iter.PasswordField.Identifier != nil {
									f4elemf7f5f1elemf4.Identifier = f4elemf7f5f1iter.PasswordField.Identifier
								}
								f4elemf7f5f1elem.PasswordField = f4elemf7f5f1elemf4
							}
							if f4elemf7f5f1iter.PayloadType != nil {
								f4elemf7f5f1elem.PayloadType = svcsdktypes.PayloadType(*f4elemf7f5f1iter.PayloadType)
							}
							if f4elemf7f5f1iter.UsernameField != nil {
								f4elemf7f5f1elemf6 := &svcsdktypes.UsernameField{}
								if f4elemf7f5f1iter.UsernameField.Identifier != nil {
									f4elemf7f5f1elemf6.Identifier = f4elemf7f5f1iter.UsernameField.Identifier
								}
								f4elemf7f5f1elem.UsernameField = f4elemf7f5f1elemf6
							}
							f4elemf7f5f1 = append(f4elemf7f5f1, *f4elemf7f5f1elem)
						}
						f4elemf7f5.ManagedRuleGroupConfigs = f4elemf7f5f1
					}
					if f4iter.Statement.ManagedRuleGroupStatement.Name != nil {
						f4elemf7f5.Name = f4iter.Statement.ManagedRuleGroupStatement.Name
					}
					if f4iter.Statement.ManagedRuleGroupStatement.RuleActionOverrides != nil {
						f4elemf7f5f3 := []svcsdktypes.RuleActionOverride{}
						for _, f4elemf7f5f3iter := range f4iter.Statement.ManagedRuleGroupStatement.RuleActionOverrides {
							f4elemf7f5f3elem := &svcsdktypes.RuleActionOverride{}
							if f4elemf7f5f3iter.ActionToUse != nil {
								f4elemf7f5f3elemf0 := &svcsdktypes.RuleAction{}
								if f4elemf7f5f3iter.ActionToUse.Allow != nil {
									f4elemf7f5f3elemf0f0 := &svcsdktypes.AllowAction{}
									if f4elemf7f5f3iter.ActionToUse.Allow.CustomRequestHandling != nil {
										f4elemf7f5f3elemf0f0f0 := &svcsdktypes.CustomRequestHandling{}
										if f4elemf7f5f3iter.ActionToUse.Allow.CustomRequestHandling.InsertHeaders != nil {
											f4elemf7f5f3elemf0f0f0f0 := []svcsdktypes.CustomHTTPHeader{}
											for _, f4elemf7f5f3elemf0f0f0f0iter := range f4elemf7f5f3iter.ActionToUse.Allow.CustomRequestHandling.InsertHeaders {
												f4elemf7f5f3elemf0f0f0f0elem := &svcsdktypes.CustomHTTPHeader{}
												if f4elemf7f5f3elemf0f0f0f0iter.Name != nil {
													f4elemf7f5f3elemf0f0f0f0elem.Name = f4elemf7f5f3elemf0f0f0f0iter.Name
												}
												if f4elemf7f5f3elemf0f0f0f0iter.Value != nil {
													f4elemf7f5f3elemf0f0f0f0elem.Value = f4elemf7f5f3elemf0f0f0f0iter.Value
												}
												f4elemf7f5f3elemf0f0f0f0 = append(f4elemf7f5f3elemf0f0f0f0, *f4elemf7f5f3elemf0f0f0f0elem)
											}
											f4elemf7f5f3elemf0f0f0.InsertHeaders = f4elemf7f5f3elemf0f0f0f0
										}
										f4elemf7f5f3elemf0f0.CustomRequestHandling = f4elemf7f5f3elemf0f0f0
									}
									f4elemf7f5f3elemf0.Allow = f4elemf7f5f3elemf0f0
								}
								if f4elemf7f5f3iter.ActionToUse.Block != nil {
									f4elemf7f5f3elemf0f1 := &svcsdktypes.BlockAction{}
									if f4elemf7f5f3iter.ActionToUse.Block.CustomResponse != nil {
										f4elemf7f5f3elemf0f1f0 := &svcsdktypes.CustomResponse{}
										if f4elemf7f5f3iter.ActionToUse.Block.CustomResponse.CustomResponseBodyKey != nil {
											f4elemf7f5f3elemf0f1f0.CustomResponseBodyKey = f4elemf7f5f3iter.ActionToUse.Block.CustomResponse.CustomResponseBodyKey
										}
										if f4elemf7f5f3iter.ActionToUse.Block.CustomResponse.ResponseCode != nil {
											if *f4elemf7f5f3iter.ActionToUse.Block.CustomResponse.ResponseCode > math.MaxInt32 || *f4elemf7f5f3iter.ActionToUse.Block.CustomResponse.ResponseCode < math.MinInt32 {
												return nil, fmt.Errorf("error: field ResponseCode is of type int32")
											}
											responseCodeCopy := int32(*f4elemf7f5f3iter.ActionToUse.Block.CustomResponse.ResponseCode)
											f4elemf7f5f3elemf0f1f0.ResponseCode = &responseCodeCopy
										}
										if f4elemf7f5f3iter.ActionToUse.Block.CustomResponse.ResponseHeaders != nil {
											f4elemf7f5f3elemf0f1f0f2 := []svcsdktypes.CustomHTTPHeader{}
											for _, f4elemf7f5f3elemf0f1f0f2iter := range f4elemf7f5f3iter.ActionToUse.Block.CustomResponse.ResponseHeaders {
												f4elemf7f5f3elemf0f1f0f2elem := &svcsdktypes.CustomHTTPHeader{}
												if f4elemf7f5f3elemf0f1f0f2iter.Name != nil {
													f4elemf7f5f3elemf0f1f0f2elem.Name = f4elemf7f5f3elemf0f1f0f2iter.Name
												}
												if f4elemf7f5f3elemf0f1f0f2iter.Value != nil {
													f4elemf7f5f3elemf0f1f0f2elem.Value = f4elemf7f5f3elemf0f1f0f2iter.Value
												}
												f4elemf7f5f3elemf0f1f0f2 = append(f4elemf7f5f3elemf0f1f0f2, *f4elemf7f5f3elemf0f1f0f2elem)
											}
											f4elemf7f5f3elemf0f1f0.ResponseHeaders = f4elemf7f5f3elemf0f1f0f2
										}
										f4elemf7f5f3elemf0f1.CustomResponse = f4elemf7f5f3elemf0f1f0
									}
									f4elemf7f5f3elemf0.Block = f4elemf7f5f3elemf0f1
								}
								if f4elemf7f5f3iter.ActionToUse.Captcha != nil {
									f4elemf7f5f3elemf0f2 := &svcsdktypes.CaptchaAction{}
									if f4elemf7f5f3iter.ActionToUse.Captcha.CustomRequestHandling != nil {
										f4elemf7f5f3elemf0f2f0 := &svcsdktypes.CustomRequestHandling{}
										if f4elemf7f5f3iter.ActionToUse.Captcha.CustomRequestHandling.InsertHeaders != nil {
											f4elemf7f5f3elemf0f2f0f0 := []svcsdktypes.CustomHTTPHeader{}
											for _, f4elemf7f5f3elemf0f2f0f0iter := range f4elemf7f5f3iter.ActionToUse.Captcha.CustomRequestHandling.InsertHeaders {
												f4elemf7f5f3elemf0f2f0f0elem := &svcsdktypes.CustomHTTPHeader{}
												if f4elemf7f5f3elemf0f2f0f0iter.Name != nil {
													f4elemf7f5f3elemf0f2f0f0elem.Name = f4elemf7f5f3elemf0f2f0f0iter.Name
												}
												if f4elemf7f5f3elemf0f2f0f0iter.Value != nil {
													f4elemf7f5f3elemf0f2f0f0elem.Value = f4elemf7f5f3elemf0f2f0f0iter.Value
												}
												f4elemf7f5f3elemf0f2f0f0 = append(f4elemf7f5f3elemf0f2f0f0, *f4elemf7f5f3elemf0f2f0f0elem)
											}
											f4elemf7f5f3elemf0f2f0.InsertHeaders = f4elemf7f5f3elemf0f2f0f0
										}
										f4elemf7f5f3elemf0f2.CustomRequestHandling = f4elemf7f5f3elemf0f2f0
									}
									f4elemf7f5f3elemf0.Captcha = f4elemf7f5f3elemf0f2
								}
								if f4elemf7f5f3iter.ActionToUse.Challenge != nil {
									f4elemf7f5f3elemf0f3 := &svcsdktypes.ChallengeAction{}
									if f4elemf7f5f3iter.ActionToUse.Challenge.CustomRequestHandling != nil {
										f4elemf7f5f3elemf0f3f0 := &svcsdktypes.CustomRequestHandling{}
										if f4elemf7f5f3iter.ActionToUse.Challenge.CustomRequestHandling.InsertHeaders != nil {
											f4elemf7f5f3elemf0f3f0f0 := []svcsdktypes.CustomHTTPHeader{}
											for _, f4elemf7f5f3elemf0f3f0f0iter := range f4elemf7f5f3iter.ActionToUse.Challenge.CustomRequestHandling.InsertHeaders {
												f4elemf7f5f3elemf0f3f0f0elem := &svcsdktypes.CustomHTTPHeader{}
												if f4elemf7f5f3elemf0f3f0f0iter.Name != nil {
													f4elemf7f5f3elemf0f3f0f0elem.Name = f4elemf7f5f3elemf0f3f0f0iter.Name
												}
												if f4elemf7f5f3elemf0f3f0f0iter.Value != nil {
													f4elemf7f5f3elemf0f3f0f0elem.Value = f4elemf7f5f3elemf0f3f0f0iter.Value
												}
												f4elemf7f5f3elemf0f3f0f0 = append(f4elemf7f5f3elemf0f3f0f0, *f4elemf7f5f3elemf0f3f0f0elem)
											}
											f4elemf7f5f3elemf0f3f0.InsertHeaders = f4elemf7f5f3elemf0f3f0f0
										}
										f4elemf7f5f3elemf0f3.CustomRequestHandling = f4elemf7f5f3elemf0f3f0
									}
									f4elemf7f5f3elemf0.Challenge = f4elemf7f5f3elemf0f3
								}
								if f4elemf7f5f3iter.ActionToUse.Count != nil {
									f4elemf7f5f3elemf0f4 := &svcsdktypes.CountAction{}
									if f4elemf7f5f3iter.ActionToUse.Count.CustomRequestHandling != nil {
										f4elemf7f5f3elemf0f4f0 := &svcsdktypes.CustomRequestHandling{}
										if f4elemf7f5f3iter.ActionToUse.Count.CustomRequestHandling.InsertHeaders != nil {
											f4elemf7f5f3elemf0f4f0f0 := []svcsdktypes.CustomHTTPHeader{}
											for _, f4elemf7f5f3elemf0f4f0f0iter := range f4elemf7f5f3iter.ActionToUse.Count.CustomRequestHandling.InsertHeaders {
												f4elemf7f5f3elemf0f4f0f0elem := &svcsdktypes.CustomHTTPHeader{}
												if f4elemf7f5f3elemf0f4f0f0iter.Name != nil {
													f4elemf7f5f3elemf0f4f0f0elem.Name = f4elemf7f5f3elemf0f4f0f0iter.Name
												}
												if f4elemf7f5f3elemf0f4f0f0iter.Value != nil {
													f4elemf7f5f3elemf0f4f0f0elem.Value = f4elemf7f5f3elemf0f4f0f0iter.Value
												}
												f4elemf7f5f3elemf0f4f0f0 = append(f4elemf7f5f3elemf0f4f0f0, *f4elemf7f5f3elemf0f4f0f0elem)
											}
											f4elemf7f5f3elemf0f4f0.InsertHeaders = f4elemf7f5f3elemf0f4f0f0
										}
										f4elemf7f5f3elemf0f4.CustomRequestHandling = f4elemf7f5f3elemf0f4f0
									}
									f4elemf7f5f3elemf0.Count = f4elemf7f5f3elemf0f4
								}
								f4elemf7f5f3elem.ActionToUse = f4elemf7f5f3elemf0
							}
							if f4elemf7f5f3iter.Name != nil {
								f4elemf7f5f3elem.Name = f4elemf7f5f3iter.Name
							}
							f4elemf7f5f3 = append(f4elemf7f5f3, *f4elemf7f5f3elem)
						}
						f4elemf7f5.RuleActionOverrides = f4elemf7f5f3
					}
					if f4iter.Statement.ManagedRuleGroupStatement.VendorName != nil {
						f4elemf7f5.VendorName = f4iter.Statement.ManagedRuleGroupStatement.VendorName
					}
					if f4iter.Statement.ManagedRuleGroupStatement.Version != nil {
						f4elemf7f5.Version = f4iter.Statement.ManagedRuleGroupStatement.Version
					}
					f4elemf7.ManagedRuleGroupStatement = f4elemf7f5
				}
				if f4iter.Statement.RateBasedStatement != nil {
					f4elemf7f8 := &svcsdktypes.RateBasedStatement{}
					if f4iter.Statement.RateBasedStatement.AggregateKeyType != nil {
						f4elemf7f8.AggregateKeyType = svcsdktypes.RateBasedStatementAggregateKeyType(*f4iter.Statement.RateBasedStatement.AggregateKeyType)
					}
					if f4iter.Statement.RateBasedStatement.CustomKeys != nil {
						f4elemf7f8f1 := []svcsdktypes.RateBasedStatementCustomKey{}
						for _, f4elemf7f8f1iter := range f4iter.Statement.RateBasedStatement.CustomKeys {
							f4elemf7f8f1elem := &svcsdktypes.RateBasedStatementCustomKey{}
							if f4elemf7f8f1iter.Cookie != nil {
								f4elemf7f8f1elemf0 := &svcsdktypes.RateLimitCookie{}
								if f4elemf7f8f1iter.Cookie.Name != nil {
									f4elemf7f8f1elemf0.Name = f4elemf7f8f1iter.Cookie.Name
								}
								if f4elemf7f8f1iter.Cookie.TextTransformations != nil {
									f4elemf7f8f1elemf0f1 := []svcsdktypes.TextTransformation{}
									for _, f4elemf7f8f1elemf0f1iter := range f4elemf7f8f1iter.Cookie.TextTransformations {
										f4elemf7f8f1elemf0f1elem := &svcsdktypes.TextTransformation{}
										if f4elemf7f8f1elemf0f1iter.Priority != nil {
											if *f4elemf7f8f1elemf0f1iter.Priority > math.MaxInt32 || *f4elemf7f8f1elemf0f1iter.Priority < math.MinInt32 {
												return nil, fmt.Errorf("error: field Priority is of type int32")
											}
											priorityCopy := int32(*f4elemf7f8f1elemf0f1iter.Priority)
											f4elemf7f8f1elemf0f1elem.Priority = priorityCopy
										}
										if f4elemf7f8f1elemf0f1iter.Type != nil {
											f4elemf7f8f1elemf0f1elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f8f1elemf0f1iter.Type)
										}
										f4elemf7f8f1elemf0f1 = append(f4elemf7f8f1elemf0f1, *f4elemf7f8f1elemf0f1elem)
									}
									f4elemf7f8f1elemf0.TextTransformations = f4elemf7f8f1elemf0f1
								}
								f4elemf7f8f1elem.Cookie = f4elemf7f8f1elemf0
							}
							if f4elemf7f8f1iter.ForwardedIP != nil {
								f4elemf7f8f1elemf1 := &svcsdktypes.RateLimitForwardedIP{}
								f4elemf7f8f1elem.ForwardedIP = f4elemf7f8f1elemf1
							}
							if f4elemf7f8f1iter.HTTPMethod != nil {
								f4elemf7f8f1elemf2 := &svcsdktypes.RateLimitHTTPMethod{}
								f4elemf7f8f1elem.HTTPMethod = f4elemf7f8f1elemf2
							}
							if f4elemf7f8f1iter.Header != nil {
								f4elemf7f8f1elemf3 := &svcsdktypes.RateLimitHeader{}
								if f4elemf7f8f1iter.Header.Name != nil {
									f4elemf7f8f1elemf3.Name = f4elemf7f8f1iter.Header.Name
								}
								if f4elemf7f8f1iter.Header.TextTransformations != nil {
									f4elemf7f8f1elemf3f1 := []svcsdktypes.TextTransformation{}
									for _, f4elemf7f8f1elemf3f1iter := range f4elemf7f8f1iter.Header.TextTransformations {
										f4elemf7f8f1elemf3f1elem := &svcsdktypes.TextTransformation{}
										if f4elemf7f8f1elemf3f1iter.Priority != nil {
											if *f4elemf7f8f1elemf3f1iter.Priority > math.MaxInt32 || *f4elemf7f8f1elemf3f1iter.Priority < math.MinInt32 {
												return nil, fmt.Errorf("error: field Priority is of type int32")
											}
											priorityCopy := int32(*f4elemf7f8f1elemf3f1iter.Priority)
											f4elemf7f8f1elemf3f1elem.Priority = priorityCopy
										}
										if f4elemf7f8f1elemf3f1iter.Type != nil {
											f4elemf7f8f1elemf3f1elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f8f1elemf3f1iter.Type)
										}
										f4elemf7f8f1elemf3f1 = append(f4elemf7f8f1elemf3f1, *f4elemf7f8f1elemf3f1elem)
									}
									f4elemf7f8f1elemf3.TextTransformations = f4elemf7f8f1elemf3f1
								}
								f4elemf7f8f1elem.Header = f4elemf7f8f1elemf3
							}
							if f4elemf7f8f1iter.IP != nil {
								f4elemf7f8f1elemf4 := &svcsdktypes.RateLimitIP{}
								f4elemf7f8f1elem.IP = f4elemf7f8f1elemf4
							}
							if f4elemf7f8f1iter.LabelNamespace != nil {
								f4elemf7f8f1elemf5 := &svcsdktypes.RateLimitLabelNamespace{}
								if f4elemf7f8f1iter.LabelNamespace.Namespace != nil {
									f4elemf7f8f1elemf5.Namespace = f4elemf7f8f1iter.LabelNamespace.Namespace
								}
								f4elemf7f8f1elem.LabelNamespace = f4elemf7f8f1elemf5
							}
							if f4elemf7f8f1iter.QueryArgument != nil {
								f4elemf7f8f1elemf6 := &svcsdktypes.RateLimitQueryArgument{}
								if f4elemf7f8f1iter.QueryArgument.Name != nil {
									f4elemf7f8f1elemf6.Name = f4elemf7f8f1iter.QueryArgument.Name
								}
								if f4elemf7f8f1iter.QueryArgument.TextTransformations != nil {
									f4elemf7f8f1elemf6f1 := []svcsdktypes.TextTransformation{}
									for _, f4elemf7f8f1elemf6f1iter := range f4elemf7f8f1iter.QueryArgument.TextTransformations {
										f4elemf7f8f1elemf6f1elem := &svcsdktypes.TextTransformation{}
										if f4elemf7f8f1elemf6f1iter.Priority != nil {
											if *f4elemf7f8f1elemf6f1iter.Priority > math.MaxInt32 || *f4elemf7f8f1elemf6f1iter.Priority < math.MinInt32 {
												return nil, fmt.Errorf("error: field Priority is of type int32")
											}
											priorityCopy := int32(*f4elemf7f8f1elemf6f1iter.Priority)
											f4elemf7f8f1elemf6f1elem.Priority = priorityCopy
										}
										if f4elemf7f8f1elemf6f1iter.Type != nil {
											f4elemf7f8f1elemf6f1elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f8f1elemf6f1iter.Type)
										}
										f4elemf7f8f1elemf6f1 = append(f4elemf7f8f1elemf6f1, *f4elemf7f8f1elemf6f1elem)
									}
									f4elemf7f8f1elemf6.TextTransformations = f4elemf7f8f1elemf6f1
								}
								f4elemf7f8f1elem.QueryArgument = f4elemf7f8f1elemf6
							}
							if f4elemf7f8f1iter.QueryString != nil {
								f4elemf7f8f1elemf7 := &svcsdktypes.RateLimitQueryString{}
								if f4elemf7f8f1iter.QueryString.TextTransformations != nil {
									f4elemf7f8f1elemf7f0 := []svcsdktypes.TextTransformation{}
									for _, f4elemf7f8f1elemf7f0iter := range f4elemf7f8f1iter.QueryString.TextTransformations {
										f4elemf7f8f1elemf7f0elem := &svcsdktypes.TextTransformation{}
										if f4elemf7f8f1elemf7f0iter.Priority != nil {
											if *f4elemf7f8f1elemf7f0iter.Priority > math.MaxInt32 || *f4elemf7f8f1elemf7f0iter.Priority < math.MinInt32 {
												return nil, fmt.Errorf("error: field Priority is of type int32")
											}
											priorityCopy := int32(*f4elemf7f8f1elemf7f0iter.Priority)
											f4elemf7f8f1elemf7f0elem.Priority = priorityCopy
										}
										if f4elemf7f8f1elemf7f0iter.Type != nil {
											f4elemf7f8f1elemf7f0elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f8f1elemf7f0iter.Type)
										}
										f4elemf7f8f1elemf7f0 = append(f4elemf7f8f1elemf7f0, *f4elemf7f8f1elemf7f0elem)
									}
									f4elemf7f8f1elemf7.TextTransformations = f4elemf7f8f1elemf7f0
								}
								f4elemf7f8f1elem.QueryString = f4elemf7f8f1elemf7
							}
							if f4elemf7f8f1iter.URIPath != nil {
								f4elemf7f8f1elemf8 := &svcsdktypes.RateLimitUriPath{}
								if f4elemf7f8f1iter.URIPath.TextTransformations != nil {
									f4elemf7f8f1elemf8f0 := []svcsdktypes.TextTransformation{}
									for _, f4elemf7f8f1elemf8f0iter := range f4elemf7f8f1iter.URIPath.TextTransformations {
										f4elemf7f8f1elemf8f0elem := &svcsdktypes.TextTransformation{}
										if f4elemf7f8f1elemf8f0iter.Priority != nil {
											if *f4elemf7f8f1elemf8f0iter.Priority > math.MaxInt32 || *f4elemf7f8f1elemf8f0iter.Priority < math.MinInt32 {
												return nil, fmt.Errorf("error: field Priority is of type int32")
											}
											priorityCopy := int32(*f4elemf7f8f1elemf8f0iter.Priority)
											f4elemf7f8f1elemf8f0elem.Priority = priorityCopy
										}
										if f4elemf7f8f1elemf8f0iter.Type != nil {
											f4elemf7f8f1elemf8f0elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f8f1elemf8f0iter.Type)
										}
										f4elemf7f8f1elemf8f0 = append(f4elemf7f8f1elemf8f0, *f4elemf7f8f1elemf8f0elem)
									}
									f4elemf7f8f1elemf8.TextTransformations = f4elemf7f8f1elemf8f0
								}
								f4elemf7f8f1elem.UriPath = f4elemf7f8f1elemf8
							}
							f4elemf7f8f1 = append(f4elemf7f8f1, *f4elemf7f8f1elem)
						}
						f4elemf7f8.CustomKeys = f4elemf7f8f1
					}
					if f4iter.Statement.RateBasedStatement.EvaluationWindowSec != nil {
						f4elemf7f8.EvaluationWindowSec = *f4iter.Statement.RateBasedStatement.EvaluationWindowSec
					}
					if f4iter.Statement.RateBasedStatement.ForwardedIPConfig != nil {
						f4elemf7f8f3 := &svcsdktypes.ForwardedIPConfig{}
						if f4iter.Statement.RateBasedStatement.ForwardedIPConfig.FallbackBehavior != nil {
							f4elemf7f8f3.FallbackBehavior = svcsdktypes.FallbackBehavior(*f4iter.Statement.RateBasedStatement.ForwardedIPConfig.FallbackBehavior)
						}
						if f4iter.Statement.RateBasedStatement.ForwardedIPConfig.HeaderName != nil {
							f4elemf7f8f3.HeaderName = f4iter.Statement.RateBasedStatement.ForwardedIPConfig.HeaderName
						}
						f4elemf7f8.ForwardedIPConfig = f4elemf7f8f3
					}
					if f4iter.Statement.RateBasedStatement.Limit != nil {
						f4elemf7f8.Limit = f4iter.Statement.RateBasedStatement.Limit
					}
					f4elemf7.RateBasedStatement = f4elemf7f8
				}
				if f4iter.Statement.RegexMatchStatement != nil {
					f4elemf7f9 := &svcsdktypes.RegexMatchStatement{}
					if f4iter.Statement.RegexMatchStatement.FieldToMatch != nil {
						f4elemf7f9f0 := &svcsdktypes.FieldToMatch{}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.AllQueryArguments != nil {
							f4elemf7f9f0f0 := &svcsdktypes.AllQueryArguments{}
							f4elemf7f9f0.AllQueryArguments = f4elemf7f9f0f0
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.Body != nil {
							f4elemf7f9f0f1 := &svcsdktypes.Body{}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.Body.OversizeHandling != nil {
								f4elemf7f9f0f1.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.RegexMatchStatement.FieldToMatch.Body.OversizeHandling)
							}
							f4elemf7f9f0.Body = f4elemf7f9f0f1
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies != nil {
							f4elemf7f9f0f2 := &svcsdktypes.Cookies{}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f4elemf7f9f0f2f0 := &svcsdktypes.CookieMatchPattern{}
								if f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f4elemf7f9f0f2f0f0 := &svcsdktypes.All{}
									f4elemf7f9f0f2f0.All = f4elemf7f9f0f2f0f0
								}
								if f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f4elemf7f9f0f2f0.ExcludedCookies = aws.ToStringSlice(f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f4elemf7f9f0f2f0.IncludedCookies = aws.ToStringSlice(f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f4elemf7f9f0f2.MatchPattern = f4elemf7f9f0f2f0
							}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchScope != nil {
								f4elemf7f9f0f2.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchScope)
							}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.OversizeHandling != nil {
								f4elemf7f9f0f2.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.OversizeHandling)
							}
							f4elemf7f9f0.Cookies = f4elemf7f9f0f2
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.HeaderOrder != nil {
							f4elemf7f9f0f3 := &svcsdktypes.HeaderOrder{}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling != nil {
								f4elemf7f9f0f3.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.RegexMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling)
							}
							f4elemf7f9f0.HeaderOrder = f4elemf7f9f0f3
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers != nil {
							f4elemf7f9f0f4 := &svcsdktypes.Headers{}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern != nil {
								f4elemf7f9f0f4f0 := &svcsdktypes.HeaderMatchPattern{}
								if f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f4elemf7f9f0f4f0f0 := &svcsdktypes.All{}
									f4elemf7f9f0f4f0.All = f4elemf7f9f0f4f0f0
								}
								if f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f4elemf7f9f0f4f0.ExcludedHeaders = aws.ToStringSlice(f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f4elemf7f9f0f4f0.IncludedHeaders = aws.ToStringSlice(f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f4elemf7f9f0f4.MatchPattern = f4elemf7f9f0f4f0
							}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchScope != nil {
								f4elemf7f9f0f4.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchScope)
							}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers.OversizeHandling != nil {
								f4elemf7f9f0f4.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.RegexMatchStatement.FieldToMatch.Headers.OversizeHandling)
							}
							f4elemf7f9f0.Headers = f4elemf7f9f0f4
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.JA3Fingerprint != nil {
							f4elemf7f9f0f5 := &svcsdktypes.JA3Fingerprint{}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != nil {
								f4elemf7f9f0f5.FallbackBehavior = svcsdktypes.FallbackBehavior(*f4iter.Statement.RegexMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior)
							}
							f4elemf7f9f0.JA3Fingerprint = f4elemf7f9f0f5
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody != nil {
							f4elemf7f9f0f6 := &svcsdktypes.JsonBody{}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior != nil {
								f4elemf7f9f0f6.InvalidFallbackBehavior = svcsdktypes.BodyParsingFallbackBehavior(*f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior)
							}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody.MatchPattern != nil {
								f4elemf7f9f0f6f1 := &svcsdktypes.JsonMatchPattern{}
								if f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody.MatchPattern.All != nil {
									f4elemf7f9f0f6f1f0 := &svcsdktypes.All{}
									f4elemf7f9f0f6f1.All = f4elemf7f9f0f6f1f0
								}
								if f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths != nil {
									f4elemf7f9f0f6f1.IncludedPaths = aws.ToStringSlice(f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths)
								}
								f4elemf7f9f0f6.MatchPattern = f4elemf7f9f0f6f1
							}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody.MatchScope != nil {
								f4elemf7f9f0f6.MatchScope = svcsdktypes.JsonMatchScope(*f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody.MatchScope)
							}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody.OversizeHandling != nil {
								f4elemf7f9f0f6.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.RegexMatchStatement.FieldToMatch.JSONBody.OversizeHandling)
							}
							f4elemf7f9f0.JsonBody = f4elemf7f9f0f6
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.Method != nil {
							f4elemf7f9f0f7 := &svcsdktypes.Method{}
							f4elemf7f9f0.Method = f4elemf7f9f0f7
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.QueryString != nil {
							f4elemf7f9f0f8 := &svcsdktypes.QueryString{}
							f4elemf7f9f0.QueryString = f4elemf7f9f0f8
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.SingleHeader != nil {
							f4elemf7f9f0f9 := &svcsdktypes.SingleHeader{}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.SingleHeader.Name != nil {
								f4elemf7f9f0f9.Name = f4iter.Statement.RegexMatchStatement.FieldToMatch.SingleHeader.Name
							}
							f4elemf7f9f0.SingleHeader = f4elemf7f9f0f9
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.SingleQueryArgument != nil {
							f4elemf7f9f0f10 := &svcsdktypes.SingleQueryArgument{}
							if f4iter.Statement.RegexMatchStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f4elemf7f9f0f10.Name = f4iter.Statement.RegexMatchStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f4elemf7f9f0.SingleQueryArgument = f4elemf7f9f0f10
						}
						if f4iter.Statement.RegexMatchStatement.FieldToMatch.URIPath != nil {
							f4elemf7f9f0f11 := &svcsdktypes.UriPath{}
							f4elemf7f9f0.UriPath = f4elemf7f9f0f11
						}
						f4elemf7f9.FieldToMatch = f4elemf7f9f0
					}
					if f4iter.Statement.RegexMatchStatement.RegexString != nil {
						f4elemf7f9.RegexString = f4iter.Statement.RegexMatchStatement.RegexString
					}
					if f4iter.Statement.RegexMatchStatement.TextTransformations != nil {
						f4elemf7f9f2 := []svcsdktypes.TextTransformation{}
						for _, f4elemf7f9f2iter := range f4iter.Statement.RegexMatchStatement.TextTransformations {
							f4elemf7f9f2elem := &svcsdktypes.TextTransformation{}
							if f4elemf7f9f2iter.Priority != nil {
								if *f4elemf7f9f2iter.Priority > math.MaxInt32 || *f4elemf7f9f2iter.Priority < math.MinInt32 {
									return nil, fmt.Errorf("error: field Priority is of type int32")
								}
								priorityCopy := int32(*f4elemf7f9f2iter.Priority)
								f4elemf7f9f2elem.Priority = priorityCopy
							}
							if f4elemf7f9f2iter.Type != nil {
								f4elemf7f9f2elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f9f2iter.Type)
							}
							f4elemf7f9f2 = append(f4elemf7f9f2, *f4elemf7f9f2elem)
						}
						f4elemf7f9.TextTransformations = f4elemf7f9f2
					}
					f4elemf7.RegexMatchStatement = f4elemf7f9
				}
				if f4iter.Statement.RegexPatternSetReferenceStatement != nil {
					f4elemf7f10 := &svcsdktypes.RegexPatternSetReferenceStatement{}
					if f4iter.Statement.RegexPatternSetReferenceStatement.ARN != nil {
						f4elemf7f10.ARN = f4iter.Statement.RegexPatternSetReferenceStatement.ARN
					}
					if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch != nil {
						f4elemf7f10f1 := &svcsdktypes.FieldToMatch{}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.AllQueryArguments != nil {
							f4elemf7f10f1f0 := &svcsdktypes.AllQueryArguments{}
							f4elemf7f10f1.AllQueryArguments = f4elemf7f10f1f0
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Body != nil {
							f4elemf7f10f1f1 := &svcsdktypes.Body{}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Body.OversizeHandling != nil {
								f4elemf7f10f1f1.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Body.OversizeHandling)
							}
							f4elemf7f10f1.Body = f4elemf7f10f1f1
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies != nil {
							f4elemf7f10f1f2 := &svcsdktypes.Cookies{}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f4elemf7f10f1f2f0 := &svcsdktypes.CookieMatchPattern{}
								if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f4elemf7f10f1f2f0f0 := &svcsdktypes.All{}
									f4elemf7f10f1f2f0.All = f4elemf7f10f1f2f0f0
								}
								if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f4elemf7f10f1f2f0.ExcludedCookies = aws.ToStringSlice(f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f4elemf7f10f1f2f0.IncludedCookies = aws.ToStringSlice(f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f4elemf7f10f1f2.MatchPattern = f4elemf7f10f1f2f0
							}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchScope != nil {
								f4elemf7f10f1f2.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchScope)
							}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.OversizeHandling != nil {
								f4elemf7f10f1f2.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.OversizeHandling)
							}
							f4elemf7f10f1.Cookies = f4elemf7f10f1f2
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.HeaderOrder != nil {
							f4elemf7f10f1f3 := &svcsdktypes.HeaderOrder{}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.HeaderOrder.OversizeHandling != nil {
								f4elemf7f10f1f3.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.HeaderOrder.OversizeHandling)
							}
							f4elemf7f10f1.HeaderOrder = f4elemf7f10f1f3
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers != nil {
							f4elemf7f10f1f4 := &svcsdktypes.Headers{}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern != nil {
								f4elemf7f10f1f4f0 := &svcsdktypes.HeaderMatchPattern{}
								if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f4elemf7f10f1f4f0f0 := &svcsdktypes.All{}
									f4elemf7f10f1f4f0.All = f4elemf7f10f1f4f0f0
								}
								if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f4elemf7f10f1f4f0.ExcludedHeaders = aws.ToStringSlice(f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f4elemf7f10f1f4f0.IncludedHeaders = aws.ToStringSlice(f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f4elemf7f10f1f4.MatchPattern = f4elemf7f10f1f4f0
							}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchScope != nil {
								f4elemf7f10f1f4.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchScope)
							}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.OversizeHandling != nil {
								f4elemf7f10f1f4.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.OversizeHandling)
							}
							f4elemf7f10f1.Headers = f4elemf7f10f1f4
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JA3Fingerprint != nil {
							f4elemf7f10f1f5 := &svcsdktypes.JA3Fingerprint{}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != nil {
								f4elemf7f10f1f5.FallbackBehavior = svcsdktypes.FallbackBehavior(*f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior)
							}
							f4elemf7f10f1.JA3Fingerprint = f4elemf7f10f1f5
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody != nil {
							f4elemf7f10f1f6 := &svcsdktypes.JsonBody{}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior != nil {
								f4elemf7f10f1f6.InvalidFallbackBehavior = svcsdktypes.BodyParsingFallbackBehavior(*f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior)
							}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody.MatchPattern != nil {
								f4elemf7f10f1f6f1 := &svcsdktypes.JsonMatchPattern{}
								if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody.MatchPattern.All != nil {
									f4elemf7f10f1f6f1f0 := &svcsdktypes.All{}
									f4elemf7f10f1f6f1.All = f4elemf7f10f1f6f1f0
								}
								if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths != nil {
									f4elemf7f10f1f6f1.IncludedPaths = aws.ToStringSlice(f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths)
								}
								f4elemf7f10f1f6.MatchPattern = f4elemf7f10f1f6f1
							}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody.MatchScope != nil {
								f4elemf7f10f1f6.MatchScope = svcsdktypes.JsonMatchScope(*f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody.MatchScope)
							}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody.OversizeHandling != nil {
								f4elemf7f10f1f6.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JSONBody.OversizeHandling)
							}
							f4elemf7f10f1.JsonBody = f4elemf7f10f1f6
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Method != nil {
							f4elemf7f10f1f7 := &svcsdktypes.Method{}
							f4elemf7f10f1.Method = f4elemf7f10f1f7
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.QueryString != nil {
							f4elemf7f10f1f8 := &svcsdktypes.QueryString{}
							f4elemf7f10f1.QueryString = f4elemf7f10f1f8
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleHeader != nil {
							f4elemf7f10f1f9 := &svcsdktypes.SingleHeader{}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleHeader.Name != nil {
								f4elemf7f10f1f9.Name = f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleHeader.Name
							}
							f4elemf7f10f1.SingleHeader = f4elemf7f10f1f9
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleQueryArgument != nil {
							f4elemf7f10f1f10 := &svcsdktypes.SingleQueryArgument{}
							if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f4elemf7f10f1f10.Name = f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f4elemf7f10f1.SingleQueryArgument = f4elemf7f10f1f10
						}
						if f4iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.URIPath != nil {
							f4elemf7f10f1f11 := &svcsdktypes.UriPath{}
							f4elemf7f10f1.UriPath = f4elemf7f10f1f11
						}
						f4elemf7f10.FieldToMatch = f4elemf7f10f1
					}
					if f4iter.Statement.RegexPatternSetReferenceStatement.TextTransformations != nil {
						f4elemf7f10f2 := []svcsdktypes.TextTransformation{}
						for _, f4elemf7f10f2iter := range f4iter.Statement.RegexPatternSetReferenceStatement.TextTransformations {
							f4elemf7f10f2elem := &svcsdktypes.TextTransformation{}
							if f4elemf7f10f2iter.Priority != nil {
								if *f4elemf7f10f2iter.Priority > math.MaxInt32 || *f4elemf7f10f2iter.Priority < math.MinInt32 {
									return nil, fmt.Errorf("error: field Priority is of type int32")
								}
								priorityCopy := int32(*f4elemf7f10f2iter.Priority)
								f4elemf7f10f2elem.Priority = priorityCopy
							}
							if f4elemf7f10f2iter.Type != nil {
								f4elemf7f10f2elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f10f2iter.Type)
							}
							f4elemf7f10f2 = append(f4elemf7f10f2, *f4elemf7f10f2elem)
						}
						f4elemf7f10.TextTransformations = f4elemf7f10f2
					}
					f4elemf7.RegexPatternSetReferenceStatement = f4elemf7f10
				}
				if f4iter.Statement.RuleGroupReferenceStatement != nil {
					f4elemf7f11 := &svcsdktypes.RuleGroupReferenceStatement{}
					if f4iter.Statement.RuleGroupReferenceStatement.ARN != nil {
						f4elemf7f11.ARN = f4iter.Statement.RuleGroupReferenceStatement.ARN
					}
					if f4iter.Statement.RuleGroupReferenceStatement.ExcludedRules != nil {
						f4elemf7f11f1 := []svcsdktypes.ExcludedRule{}
						for _, f4elemf7f11f1iter := range f4iter.Statement.RuleGroupReferenceStatement.ExcludedRules {
							f4elemf7f11f1elem := &svcsdktypes.ExcludedRule{}
							if f4elemf7f11f1iter.Name != nil {
								f4elemf7f11f1elem.Name = f4elemf7f11f1iter.Name
							}
							f4elemf7f11f1 = append(f4elemf7f11f1, *f4elemf7f11f1elem)
						}
						f4elemf7f11.ExcludedRules = f4elemf7f11f1
					}
					if f4iter.Statement.RuleGroupReferenceStatement.RuleActionOverrides != nil {
						f4elemf7f11f2 := []svcsdktypes.RuleActionOverride{}
						for _, f4elemf7f11f2iter := range f4iter.Statement.RuleGroupReferenceStatement.RuleActionOverrides {
							f4elemf7f11f2elem := &svcsdktypes.RuleActionOverride{}
							if f4elemf7f11f2iter.ActionToUse != nil {
								f4elemf7f11f2elemf0 := &svcsdktypes.RuleAction{}
								if f4elemf7f11f2iter.ActionToUse.Allow != nil {
									f4elemf7f11f2elemf0f0 := &svcsdktypes.AllowAction{}
									if f4elemf7f11f2iter.ActionToUse.Allow.CustomRequestHandling != nil {
										f4elemf7f11f2elemf0f0f0 := &svcsdktypes.CustomRequestHandling{}
										if f4elemf7f11f2iter.ActionToUse.Allow.CustomRequestHandling.InsertHeaders != nil {
											f4elemf7f11f2elemf0f0f0f0 := []svcsdktypes.CustomHTTPHeader{}
											for _, f4elemf7f11f2elemf0f0f0f0iter := range f4elemf7f11f2iter.ActionToUse.Allow.CustomRequestHandling.InsertHeaders {
												f4elemf7f11f2elemf0f0f0f0elem := &svcsdktypes.CustomHTTPHeader{}
												if f4elemf7f11f2elemf0f0f0f0iter.Name != nil {
													f4elemf7f11f2elemf0f0f0f0elem.Name = f4elemf7f11f2elemf0f0f0f0iter.Name
												}
												if f4elemf7f11f2elemf0f0f0f0iter.Value != nil {
													f4elemf7f11f2elemf0f0f0f0elem.Value = f4elemf7f11f2elemf0f0f0f0iter.Value
												}
												f4elemf7f11f2elemf0f0f0f0 = append(f4elemf7f11f2elemf0f0f0f0, *f4elemf7f11f2elemf0f0f0f0elem)
											}
											f4elemf7f11f2elemf0f0f0.InsertHeaders = f4elemf7f11f2elemf0f0f0f0
										}
										f4elemf7f11f2elemf0f0.CustomRequestHandling = f4elemf7f11f2elemf0f0f0
									}
									f4elemf7f11f2elemf0.Allow = f4elemf7f11f2elemf0f0
								}
								if f4elemf7f11f2iter.ActionToUse.Block != nil {
									f4elemf7f11f2elemf0f1 := &svcsdktypes.BlockAction{}
									if f4elemf7f11f2iter.ActionToUse.Block.CustomResponse != nil {
										f4elemf7f11f2elemf0f1f0 := &svcsdktypes.CustomResponse{}
										if f4elemf7f11f2iter.ActionToUse.Block.CustomResponse.CustomResponseBodyKey != nil {
											f4elemf7f11f2elemf0f1f0.CustomResponseBodyKey = f4elemf7f11f2iter.ActionToUse.Block.CustomResponse.CustomResponseBodyKey
										}
										if f4elemf7f11f2iter.ActionToUse.Block.CustomResponse.ResponseCode != nil {
											if *f4elemf7f11f2iter.ActionToUse.Block.CustomResponse.ResponseCode > math.MaxInt32 || *f4elemf7f11f2iter.ActionToUse.Block.CustomResponse.ResponseCode < math.MinInt32 {
												return nil, fmt.Errorf("error: field ResponseCode is of type int32")
											}
											responseCodeCopy := int32(*f4elemf7f11f2iter.ActionToUse.Block.CustomResponse.ResponseCode)
											f4elemf7f11f2elemf0f1f0.ResponseCode = &responseCodeCopy
										}
										if f4elemf7f11f2iter.ActionToUse.Block.CustomResponse.ResponseHeaders != nil {
											f4elemf7f11f2elemf0f1f0f2 := []svcsdktypes.CustomHTTPHeader{}
											for _, f4elemf7f11f2elemf0f1f0f2iter := range f4elemf7f11f2iter.ActionToUse.Block.CustomResponse.ResponseHeaders {
												f4elemf7f11f2elemf0f1f0f2elem := &svcsdktypes.CustomHTTPHeader{}
												if f4elemf7f11f2elemf0f1f0f2iter.Name != nil {
													f4elemf7f11f2elemf0f1f0f2elem.Name = f4elemf7f11f2elemf0f1f0f2iter.Name
												}
												if f4elemf7f11f2elemf0f1f0f2iter.Value != nil {
													f4elemf7f11f2elemf0f1f0f2elem.Value = f4elemf7f11f2elemf0f1f0f2iter.Value
												}
												f4elemf7f11f2elemf0f1f0f2 = append(f4elemf7f11f2elemf0f1f0f2, *f4elemf7f11f2elemf0f1f0f2elem)
											}
											f4elemf7f11f2elemf0f1f0.ResponseHeaders = f4elemf7f11f2elemf0f1f0f2
										}
										f4elemf7f11f2elemf0f1.CustomResponse = f4elemf7f11f2elemf0f1f0
									}
									f4elemf7f11f2elemf0.Block = f4elemf7f11f2elemf0f1
								}
								if f4elemf7f11f2iter.ActionToUse.Captcha != nil {
									f4elemf7f11f2elemf0f2 := &svcsdktypes.CaptchaAction{}
									if f4elemf7f11f2iter.ActionToUse.Captcha.CustomRequestHandling != nil {
										f4elemf7f11f2elemf0f2f0 := &svcsdktypes.CustomRequestHandling{}
										if f4elemf7f11f2iter.ActionToUse.Captcha.CustomRequestHandling.InsertHeaders != nil {
											f4elemf7f11f2elemf0f2f0f0 := []svcsdktypes.CustomHTTPHeader{}
											for _, f4elemf7f11f2elemf0f2f0f0iter := range f4elemf7f11f2iter.ActionToUse.Captcha.CustomRequestHandling.InsertHeaders {
												f4elemf7f11f2elemf0f2f0f0elem := &svcsdktypes.CustomHTTPHeader{}
												if f4elemf7f11f2elemf0f2f0f0iter.Name != nil {
													f4elemf7f11f2elemf0f2f0f0elem.Name = f4elemf7f11f2elemf0f2f0f0iter.Name
												}
												if f4elemf7f11f2elemf0f2f0f0iter.Value != nil {
													f4elemf7f11f2elemf0f2f0f0elem.Value = f4elemf7f11f2elemf0f2f0f0iter.Value
												}
												f4elemf7f11f2elemf0f2f0f0 = append(f4elemf7f11f2elemf0f2f0f0, *f4elemf7f11f2elemf0f2f0f0elem)
											}
											f4elemf7f11f2elemf0f2f0.InsertHeaders = f4elemf7f11f2elemf0f2f0f0
										}
										f4elemf7f11f2elemf0f2.CustomRequestHandling = f4elemf7f11f2elemf0f2f0
									}
									f4elemf7f11f2elemf0.Captcha = f4elemf7f11f2elemf0f2
								}
								if f4elemf7f11f2iter.ActionToUse.Challenge != nil {
									f4elemf7f11f2elemf0f3 := &svcsdktypes.ChallengeAction{}
									if f4elemf7f11f2iter.ActionToUse.Challenge.CustomRequestHandling != nil {
										f4elemf7f11f2elemf0f3f0 := &svcsdktypes.CustomRequestHandling{}
										if f4elemf7f11f2iter.ActionToUse.Challenge.CustomRequestHandling.InsertHeaders != nil {
											f4elemf7f11f2elemf0f3f0f0 := []svcsdktypes.CustomHTTPHeader{}
											for _, f4elemf7f11f2elemf0f3f0f0iter := range f4elemf7f11f2iter.ActionToUse.Challenge.CustomRequestHandling.InsertHeaders {
												f4elemf7f11f2elemf0f3f0f0elem := &svcsdktypes.CustomHTTPHeader{}
												if f4elemf7f11f2elemf0f3f0f0iter.Name != nil {
													f4elemf7f11f2elemf0f3f0f0elem.Name = f4elemf7f11f2elemf0f3f0f0iter.Name
												}
												if f4elemf7f11f2elemf0f3f0f0iter.Value != nil {
													f4elemf7f11f2elemf0f3f0f0elem.Value = f4elemf7f11f2elemf0f3f0f0iter.Value
												}
												f4elemf7f11f2elemf0f3f0f0 = append(f4elemf7f11f2elemf0f3f0f0, *f4elemf7f11f2elemf0f3f0f0elem)
											}
											f4elemf7f11f2elemf0f3f0.InsertHeaders = f4elemf7f11f2elemf0f3f0f0
										}
										f4elemf7f11f2elemf0f3.CustomRequestHandling = f4elemf7f11f2elemf0f3f0
									}
									f4elemf7f11f2elemf0.Challenge = f4elemf7f11f2elemf0f3
								}
								if f4elemf7f11f2iter.ActionToUse.Count != nil {
									f4elemf7f11f2elemf0f4 := &svcsdktypes.CountAction{}
									if f4elemf7f11f2iter.ActionToUse.Count.CustomRequestHandling != nil {
										f4elemf7f11f2elemf0f4f0 := &svcsdktypes.CustomRequestHandling{}
										if f4elemf7f11f2iter.ActionToUse.Count.CustomRequestHandling.InsertHeaders != nil {
											f4elemf7f11f2elemf0f4f0f0 := []svcsdktypes.CustomHTTPHeader{}
											for _, f4elemf7f11f2elemf0f4f0f0iter := range f4elemf7f11f2iter.ActionToUse.Count.CustomRequestHandling.InsertHeaders {
												f4elemf7f11f2elemf0f4f0f0elem := &svcsdktypes.CustomHTTPHeader{}
												if f4elemf7f11f2elemf0f4f0f0iter.Name != nil {
													f4elemf7f11f2elemf0f4f0f0elem.Name = f4elemf7f11f2elemf0f4f0f0iter.Name
												}
												if f4elemf7f11f2elemf0f4f0f0iter.Value != nil {
													f4elemf7f11f2elemf0f4f0f0elem.Value = f4elemf7f11f2elemf0f4f0f0iter.Value
												}
												f4elemf7f11f2elemf0f4f0f0 = append(f4elemf7f11f2elemf0f4f0f0, *f4elemf7f11f2elemf0f4f0f0elem)
											}
											f4elemf7f11f2elemf0f4f0.InsertHeaders = f4elemf7f11f2elemf0f4f0f0
										}
										f4elemf7f11f2elemf0f4.CustomRequestHandling = f4elemf7f11f2elemf0f4f0
									}
									f4elemf7f11f2elemf0.Count = f4elemf7f11f2elemf0f4
								}
								f4elemf7f11f2elem.ActionToUse = f4elemf7f11f2elemf0
							}
							if f4elemf7f11f2iter.Name != nil {
								f4elemf7f11f2elem.Name = f4elemf7f11f2iter.Name
							}
							f4elemf7f11f2 = append(f4elemf7f11f2, *f4elemf7f11f2elem)
						}
						f4elemf7f11.RuleActionOverrides = f4elemf7f11f2
					}
					f4elemf7.RuleGroupReferenceStatement = f4elemf7f11
				}
				if f4iter.Statement.SizeConstraintStatement != nil {
					f4elemf7f12 := &svcsdktypes.SizeConstraintStatement{}
					if f4iter.Statement.SizeConstraintStatement.ComparisonOperator != nil {
						f4elemf7f12.ComparisonOperator = svcsdktypes.ComparisonOperator(*f4iter.Statement.SizeConstraintStatement.ComparisonOperator)
					}
					if f4iter.Statement.SizeConstraintStatement.FieldToMatch != nil {
						f4elemf7f12f1 := &svcsdktypes.FieldToMatch{}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.AllQueryArguments != nil {
							f4elemf7f12f1f0 := &svcsdktypes.AllQueryArguments{}
							f4elemf7f12f1.AllQueryArguments = f4elemf7f12f1f0
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Body != nil {
							f4elemf7f12f1f1 := &svcsdktypes.Body{}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Body.OversizeHandling != nil {
								f4elemf7f12f1f1.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.SizeConstraintStatement.FieldToMatch.Body.OversizeHandling)
							}
							f4elemf7f12f1.Body = f4elemf7f12f1f1
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies != nil {
							f4elemf7f12f1f2 := &svcsdktypes.Cookies{}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f4elemf7f12f1f2f0 := &svcsdktypes.CookieMatchPattern{}
								if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f4elemf7f12f1f2f0f0 := &svcsdktypes.All{}
									f4elemf7f12f1f2f0.All = f4elemf7f12f1f2f0f0
								}
								if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f4elemf7f12f1f2f0.ExcludedCookies = aws.ToStringSlice(f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f4elemf7f12f1f2f0.IncludedCookies = aws.ToStringSlice(f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f4elemf7f12f1f2.MatchPattern = f4elemf7f12f1f2f0
							}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchScope != nil {
								f4elemf7f12f1f2.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchScope)
							}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.OversizeHandling != nil {
								f4elemf7f12f1f2.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.OversizeHandling)
							}
							f4elemf7f12f1.Cookies = f4elemf7f12f1f2
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.HeaderOrder != nil {
							f4elemf7f12f1f3 := &svcsdktypes.HeaderOrder{}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.HeaderOrder.OversizeHandling != nil {
								f4elemf7f12f1f3.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.SizeConstraintStatement.FieldToMatch.HeaderOrder.OversizeHandling)
							}
							f4elemf7f12f1.HeaderOrder = f4elemf7f12f1f3
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers != nil {
							f4elemf7f12f1f4 := &svcsdktypes.Headers{}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern != nil {
								f4elemf7f12f1f4f0 := &svcsdktypes.HeaderMatchPattern{}
								if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f4elemf7f12f1f4f0f0 := &svcsdktypes.All{}
									f4elemf7f12f1f4f0.All = f4elemf7f12f1f4f0f0
								}
								if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f4elemf7f12f1f4f0.ExcludedHeaders = aws.ToStringSlice(f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f4elemf7f12f1f4f0.IncludedHeaders = aws.ToStringSlice(f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f4elemf7f12f1f4.MatchPattern = f4elemf7f12f1f4f0
							}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchScope != nil {
								f4elemf7f12f1f4.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchScope)
							}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.OversizeHandling != nil {
								f4elemf7f12f1f4.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.OversizeHandling)
							}
							f4elemf7f12f1.Headers = f4elemf7f12f1f4
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.JA3Fingerprint != nil {
							f4elemf7f12f1f5 := &svcsdktypes.JA3Fingerprint{}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != nil {
								f4elemf7f12f1f5.FallbackBehavior = svcsdktypes.FallbackBehavior(*f4iter.Statement.SizeConstraintStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior)
							}
							f4elemf7f12f1.JA3Fingerprint = f4elemf7f12f1f5
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody != nil {
							f4elemf7f12f1f6 := &svcsdktypes.JsonBody{}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior != nil {
								f4elemf7f12f1f6.InvalidFallbackBehavior = svcsdktypes.BodyParsingFallbackBehavior(*f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior)
							}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody.MatchPattern != nil {
								f4elemf7f12f1f6f1 := &svcsdktypes.JsonMatchPattern{}
								if f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody.MatchPattern.All != nil {
									f4elemf7f12f1f6f1f0 := &svcsdktypes.All{}
									f4elemf7f12f1f6f1.All = f4elemf7f12f1f6f1f0
								}
								if f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths != nil {
									f4elemf7f12f1f6f1.IncludedPaths = aws.ToStringSlice(f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths)
								}
								f4elemf7f12f1f6.MatchPattern = f4elemf7f12f1f6f1
							}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody.MatchScope != nil {
								f4elemf7f12f1f6.MatchScope = svcsdktypes.JsonMatchScope(*f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody.MatchScope)
							}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody.OversizeHandling != nil {
								f4elemf7f12f1f6.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.SizeConstraintStatement.FieldToMatch.JSONBody.OversizeHandling)
							}
							f4elemf7f12f1.JsonBody = f4elemf7f12f1f6
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.Method != nil {
							f4elemf7f12f1f7 := &svcsdktypes.Method{}
							f4elemf7f12f1.Method = f4elemf7f12f1f7
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.QueryString != nil {
							f4elemf7f12f1f8 := &svcsdktypes.QueryString{}
							f4elemf7f12f1.QueryString = f4elemf7f12f1f8
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.SingleHeader != nil {
							f4elemf7f12f1f9 := &svcsdktypes.SingleHeader{}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.SingleHeader.Name != nil {
								f4elemf7f12f1f9.Name = f4iter.Statement.SizeConstraintStatement.FieldToMatch.SingleHeader.Name
							}
							f4elemf7f12f1.SingleHeader = f4elemf7f12f1f9
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.SingleQueryArgument != nil {
							f4elemf7f12f1f10 := &svcsdktypes.SingleQueryArgument{}
							if f4iter.Statement.SizeConstraintStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f4elemf7f12f1f10.Name = f4iter.Statement.SizeConstraintStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f4elemf7f12f1.SingleQueryArgument = f4elemf7f12f1f10
						}
						if f4iter.Statement.SizeConstraintStatement.FieldToMatch.URIPath != nil {
							f4elemf7f12f1f11 := &svcsdktypes.UriPath{}
							f4elemf7f12f1.UriPath = f4elemf7f12f1f11
						}
						f4elemf7f12.FieldToMatch = f4elemf7f12f1
					}
					if f4iter.Statement.SizeConstraintStatement.Size != nil {
						f4elemf7f12.Size = *f4iter.Statement.SizeConstraintStatement.Size
					}
					if f4iter.Statement.SizeConstraintStatement.TextTransformations != nil {
						f4elemf7f12f3 := []svcsdktypes.TextTransformation{}
						for _, f4elemf7f12f3iter := range f4iter.Statement.SizeConstraintStatement.TextTransformations {
							f4elemf7f12f3elem := &svcsdktypes.TextTransformation{}
							if f4elemf7f12f3iter.Priority != nil {
								if *f4elemf7f12f3iter.Priority > math.MaxInt32 || *f4elemf7f12f3iter.Priority < math.MinInt32 {
									return nil, fmt.Errorf("error: field Priority is of type int32")
								}
								priorityCopy := int32(*f4elemf7f12f3iter.Priority)
								f4elemf7f12f3elem.Priority = priorityCopy
							}
							if f4elemf7f12f3iter.Type != nil {
								f4elemf7f12f3elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f12f3iter.Type)
							}
							f4elemf7f12f3 = append(f4elemf7f12f3, *f4elemf7f12f3elem)
						}
						f4elemf7f12.TextTransformations = f4elemf7f12f3
					}
					f4elemf7.SizeConstraintStatement = f4elemf7f12
				}
				if f4iter.Statement.SQLIMatchStatement != nil {
					f4elemf7f13 := &svcsdktypes.SqliMatchStatement{}
					if f4iter.Statement.SQLIMatchStatement.FieldToMatch != nil {
						f4elemf7f13f0 := &svcsdktypes.FieldToMatch{}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.AllQueryArguments != nil {
							f4elemf7f13f0f0 := &svcsdktypes.AllQueryArguments{}
							f4elemf7f13f0.AllQueryArguments = f4elemf7f13f0f0
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Body != nil {
							f4elemf7f13f0f1 := &svcsdktypes.Body{}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Body.OversizeHandling != nil {
								f4elemf7f13f0f1.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.SQLIMatchStatement.FieldToMatch.Body.OversizeHandling)
							}
							f4elemf7f13f0.Body = f4elemf7f13f0f1
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies != nil {
							f4elemf7f13f0f2 := &svcsdktypes.Cookies{}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f4elemf7f13f0f2f0 := &svcsdktypes.CookieMatchPattern{}
								if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f4elemf7f13f0f2f0f0 := &svcsdktypes.All{}
									f4elemf7f13f0f2f0.All = f4elemf7f13f0f2f0f0
								}
								if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f4elemf7f13f0f2f0.ExcludedCookies = aws.ToStringSlice(f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f4elemf7f13f0f2f0.IncludedCookies = aws.ToStringSlice(f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f4elemf7f13f0f2.MatchPattern = f4elemf7f13f0f2f0
							}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies.MatchScope != nil {
								f4elemf7f13f0f2.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies.MatchScope)
							}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies.OversizeHandling != nil {
								f4elemf7f13f0f2.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.SQLIMatchStatement.FieldToMatch.Cookies.OversizeHandling)
							}
							f4elemf7f13f0.Cookies = f4elemf7f13f0f2
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.HeaderOrder != nil {
							f4elemf7f13f0f3 := &svcsdktypes.HeaderOrder{}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling != nil {
								f4elemf7f13f0f3.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.SQLIMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling)
							}
							f4elemf7f13f0.HeaderOrder = f4elemf7f13f0f3
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers != nil {
							f4elemf7f13f0f4 := &svcsdktypes.Headers{}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers.MatchPattern != nil {
								f4elemf7f13f0f4f0 := &svcsdktypes.HeaderMatchPattern{}
								if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f4elemf7f13f0f4f0f0 := &svcsdktypes.All{}
									f4elemf7f13f0f4f0.All = f4elemf7f13f0f4f0f0
								}
								if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f4elemf7f13f0f4f0.ExcludedHeaders = aws.ToStringSlice(f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f4elemf7f13f0f4f0.IncludedHeaders = aws.ToStringSlice(f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f4elemf7f13f0f4.MatchPattern = f4elemf7f13f0f4f0
							}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers.MatchScope != nil {
								f4elemf7f13f0f4.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers.MatchScope)
							}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers.OversizeHandling != nil {
								f4elemf7f13f0f4.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.SQLIMatchStatement.FieldToMatch.Headers.OversizeHandling)
							}
							f4elemf7f13f0.Headers = f4elemf7f13f0f4
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.JA3Fingerprint != nil {
							f4elemf7f13f0f5 := &svcsdktypes.JA3Fingerprint{}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != nil {
								f4elemf7f13f0f5.FallbackBehavior = svcsdktypes.FallbackBehavior(*f4iter.Statement.SQLIMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior)
							}
							f4elemf7f13f0.JA3Fingerprint = f4elemf7f13f0f5
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody != nil {
							f4elemf7f13f0f6 := &svcsdktypes.JsonBody{}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior != nil {
								f4elemf7f13f0f6.InvalidFallbackBehavior = svcsdktypes.BodyParsingFallbackBehavior(*f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior)
							}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody.MatchPattern != nil {
								f4elemf7f13f0f6f1 := &svcsdktypes.JsonMatchPattern{}
								if f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody.MatchPattern.All != nil {
									f4elemf7f13f0f6f1f0 := &svcsdktypes.All{}
									f4elemf7f13f0f6f1.All = f4elemf7f13f0f6f1f0
								}
								if f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths != nil {
									f4elemf7f13f0f6f1.IncludedPaths = aws.ToStringSlice(f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths)
								}
								f4elemf7f13f0f6.MatchPattern = f4elemf7f13f0f6f1
							}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody.MatchScope != nil {
								f4elemf7f13f0f6.MatchScope = svcsdktypes.JsonMatchScope(*f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody.MatchScope)
							}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody.OversizeHandling != nil {
								f4elemf7f13f0f6.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.SQLIMatchStatement.FieldToMatch.JSONBody.OversizeHandling)
							}
							f4elemf7f13f0.JsonBody = f4elemf7f13f0f6
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.Method != nil {
							f4elemf7f13f0f7 := &svcsdktypes.Method{}
							f4elemf7f13f0.Method = f4elemf7f13f0f7
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.QueryString != nil {
							f4elemf7f13f0f8 := &svcsdktypes.QueryString{}
							f4elemf7f13f0.QueryString = f4elemf7f13f0f8
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.SingleHeader != nil {
							f4elemf7f13f0f9 := &svcsdktypes.SingleHeader{}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.SingleHeader.Name != nil {
								f4elemf7f13f0f9.Name = f4iter.Statement.SQLIMatchStatement.FieldToMatch.SingleHeader.Name
							}
							f4elemf7f13f0.SingleHeader = f4elemf7f13f0f9
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.SingleQueryArgument != nil {
							f4elemf7f13f0f10 := &svcsdktypes.SingleQueryArgument{}
							if f4iter.Statement.SQLIMatchStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f4elemf7f13f0f10.Name = f4iter.Statement.SQLIMatchStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f4elemf7f13f0.SingleQueryArgument = f4elemf7f13f0f10
						}
						if f4iter.Statement.SQLIMatchStatement.FieldToMatch.URIPath != nil {
							f4elemf7f13f0f11 := &svcsdktypes.UriPath{}
							f4elemf7f13f0.UriPath = f4elemf7f13f0f11
						}
						f4elemf7f13.FieldToMatch = f4elemf7f13f0
					}
					if f4iter.Statement.SQLIMatchStatement.SensitivityLevel != nil {
						f4elemf7f13.SensitivityLevel = svcsdktypes.SensitivityLevel(*f4iter.Statement.SQLIMatchStatement.SensitivityLevel)
					}
					if f4iter.Statement.SQLIMatchStatement.TextTransformations != nil {
						f4elemf7f13f2 := []svcsdktypes.TextTransformation{}
						for _, f4elemf7f13f2iter := range f4iter.Statement.SQLIMatchStatement.TextTransformations {
							f4elemf7f13f2elem := &svcsdktypes.TextTransformation{}
							if f4elemf7f13f2iter.Priority != nil {
								if *f4elemf7f13f2iter.Priority > math.MaxInt32 || *f4elemf7f13f2iter.Priority < math.MinInt32 {
									return nil, fmt.Errorf("error: field Priority is of type int32")
								}
								priorityCopy := int32(*f4elemf7f13f2iter.Priority)
								f4elemf7f13f2elem.Priority = priorityCopy
							}
							if f4elemf7f13f2iter.Type != nil {
								f4elemf7f13f2elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f13f2iter.Type)
							}
							f4elemf7f13f2 = append(f4elemf7f13f2, *f4elemf7f13f2elem)
						}
						f4elemf7f13.TextTransformations = f4elemf7f13f2
					}
					f4elemf7.SqliMatchStatement = f4elemf7f13
				}
				if f4iter.Statement.XSSMatchStatement != nil {
					f4elemf7f14 := &svcsdktypes.XssMatchStatement{}
					if f4iter.Statement.XSSMatchStatement.FieldToMatch != nil {
						f4elemf7f14f0 := &svcsdktypes.FieldToMatch{}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.AllQueryArguments != nil {
							f4elemf7f14f0f0 := &svcsdktypes.AllQueryArguments{}
							f4elemf7f14f0.AllQueryArguments = f4elemf7f14f0f0
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.Body != nil {
							f4elemf7f14f0f1 := &svcsdktypes.Body{}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.Body.OversizeHandling != nil {
								f4elemf7f14f0f1.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.XSSMatchStatement.FieldToMatch.Body.OversizeHandling)
							}
							f4elemf7f14f0.Body = f4elemf7f14f0f1
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies != nil {
							f4elemf7f14f0f2 := &svcsdktypes.Cookies{}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f4elemf7f14f0f2f0 := &svcsdktypes.CookieMatchPattern{}
								if f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f4elemf7f14f0f2f0f0 := &svcsdktypes.All{}
									f4elemf7f14f0f2f0.All = f4elemf7f14f0f2f0f0
								}
								if f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f4elemf7f14f0f2f0.ExcludedCookies = aws.ToStringSlice(f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f4elemf7f14f0f2f0.IncludedCookies = aws.ToStringSlice(f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f4elemf7f14f0f2.MatchPattern = f4elemf7f14f0f2f0
							}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies.MatchScope != nil {
								f4elemf7f14f0f2.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies.MatchScope)
							}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies.OversizeHandling != nil {
								f4elemf7f14f0f2.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.XSSMatchStatement.FieldToMatch.Cookies.OversizeHandling)
							}
							f4elemf7f14f0.Cookies = f4elemf7f14f0f2
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.HeaderOrder != nil {
							f4elemf7f14f0f3 := &svcsdktypes.HeaderOrder{}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling != nil {
								f4elemf7f14f0f3.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.XSSMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling)
							}
							f4elemf7f14f0.HeaderOrder = f4elemf7f14f0f3
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers != nil {
							f4elemf7f14f0f4 := &svcsdktypes.Headers{}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers.MatchPattern != nil {
								f4elemf7f14f0f4f0 := &svcsdktypes.HeaderMatchPattern{}
								if f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f4elemf7f14f0f4f0f0 := &svcsdktypes.All{}
									f4elemf7f14f0f4f0.All = f4elemf7f14f0f4f0f0
								}
								if f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f4elemf7f14f0f4f0.ExcludedHeaders = aws.ToStringSlice(f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f4elemf7f14f0f4f0.IncludedHeaders = aws.ToStringSlice(f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f4elemf7f14f0f4.MatchPattern = f4elemf7f14f0f4f0
							}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers.MatchScope != nil {
								f4elemf7f14f0f4.MatchScope = svcsdktypes.MapMatchScope(*f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers.MatchScope)
							}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers.OversizeHandling != nil {
								f4elemf7f14f0f4.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.XSSMatchStatement.FieldToMatch.Headers.OversizeHandling)
							}
							f4elemf7f14f0.Headers = f4elemf7f14f0f4
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.JA3Fingerprint != nil {
							f4elemf7f14f0f5 := &svcsdktypes.JA3Fingerprint{}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != nil {
								f4elemf7f14f0f5.FallbackBehavior = svcsdktypes.FallbackBehavior(*f4iter.Statement.XSSMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior)
							}
							f4elemf7f14f0.JA3Fingerprint = f4elemf7f14f0f5
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody != nil {
							f4elemf7f14f0f6 := &svcsdktypes.JsonBody{}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior != nil {
								f4elemf7f14f0f6.InvalidFallbackBehavior = svcsdktypes.BodyParsingFallbackBehavior(*f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody.InvalidFallbackBehavior)
							}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody.MatchPattern != nil {
								f4elemf7f14f0f6f1 := &svcsdktypes.JsonMatchPattern{}
								if f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody.MatchPattern.All != nil {
									f4elemf7f14f0f6f1f0 := &svcsdktypes.All{}
									f4elemf7f14f0f6f1.All = f4elemf7f14f0f6f1f0
								}
								if f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths != nil {
									f4elemf7f14f0f6f1.IncludedPaths = aws.ToStringSlice(f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody.MatchPattern.IncludedPaths)
								}
								f4elemf7f14f0f6.MatchPattern = f4elemf7f14f0f6f1
							}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody.MatchScope != nil {
								f4elemf7f14f0f6.MatchScope = svcsdktypes.JsonMatchScope(*f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody.MatchScope)
							}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody.OversizeHandling != nil {
								f4elemf7f14f0f6.OversizeHandling = svcsdktypes.OversizeHandling(*f4iter.Statement.XSSMatchStatement.FieldToMatch.JSONBody.OversizeHandling)
							}
							f4elemf7f14f0.JsonBody = f4elemf7f14f0f6
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.Method != nil {
							f4elemf7f14f0f7 := &svcsdktypes.Method{}
							f4elemf7f14f0.Method = f4elemf7f14f0f7
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.QueryString != nil {
							f4elemf7f14f0f8 := &svcsdktypes.QueryString{}
							f4elemf7f14f0.QueryString = f4elemf7f14f0f8
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.SingleHeader != nil {
							f4elemf7f14f0f9 := &svcsdktypes.SingleHeader{}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.SingleHeader.Name != nil {
								f4elemf7f14f0f9.Name = f4iter.Statement.XSSMatchStatement.FieldToMatch.SingleHeader.Name
							}
							f4elemf7f14f0.SingleHeader = f4elemf7f14f0f9
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.SingleQueryArgument != nil {
							f4elemf7f14f0f10 := &svcsdktypes.SingleQueryArgument{}
							if f4iter.Statement.XSSMatchStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f4elemf7f14f0f10.Name = f4iter.Statement.XSSMatchStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f4elemf7f14f0.SingleQueryArgument = f4elemf7f14f0f10
						}
						if f4iter.Statement.XSSMatchStatement.FieldToMatch.URIPath != nil {
							f4elemf7f14f0f11 := &svcsdktypes.UriPath{}
							f4elemf7f14f0.UriPath = f4elemf7f14f0f11
						}
						f4elemf7f14.FieldToMatch = f4elemf7f14f0
					}
					if f4iter.Statement.XSSMatchStatement.TextTransformations != nil {
						f4elemf7f14f1 := []svcsdktypes.TextTransformation{}
						for _, f4elemf7f14f1iter := range f4iter.Statement.XSSMatchStatement.TextTransformations {
							f4elemf7f14f1elem := &svcsdktypes.TextTransformation{}
							if f4elemf7f14f1iter.Priority != nil {
								if *f4elemf7f14f1iter.Priority > math.MaxInt32 || *f4elemf7f14f1iter.Priority < math.MinInt32 {
									return nil, fmt.Errorf("error: field Priority is of type int32")
								}
								priorityCopy := int32(*f4elemf7f14f1iter.Priority)
								f4elemf7f14f1elem.Priority = priorityCopy
							}
							if f4elemf7f14f1iter.Type != nil {
								f4elemf7f14f1elem.Type = svcsdktypes.TextTransformationType(*f4elemf7f14f1iter.Type)
							}
							f4elemf7f14f1 = append(f4elemf7f14f1, *f4elemf7f14f1elem)
						}
						f4elemf7f14.TextTransformations = f4elemf7f14f1
					}
					f4elemf7.XssMatchStatement = f4elemf7f14
				}
				f4elem.Statement = f4elemf7
			}
			if f4iter.VisibilityConfig != nil {
				f4elemf8 := &svcsdktypes.VisibilityConfig{}
				if f4iter.VisibilityConfig.CloudWatchMetricsEnabled != nil {
					f4elemf8.CloudWatchMetricsEnabled = *f4iter.VisibilityConfig.CloudWatchMetricsEnabled
				}
				if f4iter.VisibilityConfig.MetricName != nil {
					f4elemf8.MetricName = f4iter.VisibilityConfig.MetricName
				}
				if f4iter.VisibilityConfig.SampledRequestsEnabled != nil {
					f4elemf8.SampledRequestsEnabled = *f4iter.VisibilityConfig.SampledRequestsEnabled
				}
				f4elem.VisibilityConfig = f4elemf8
			}
			f4 = append(f4, *f4elem)
		}
		res.Rules = f4
	}
	if r.ko.Spec.Scope != nil {
		res.Scope = svcsdktypes.Scope(*r.ko.Spec.Scope)
	}
	if r.ko.Spec.Tags != nil {
		f6 := []svcsdktypes.Tag{}
		for _, f6iter := range r.ko.Spec.Tags {
			f6elem := &svcsdktypes.Tag{}
			if f6iter.Key != nil {
				f6elem.Key = f6iter.Key
			}
			if f6iter.Value != nil {
				f6elem.Value = f6iter.Value
			}
			f6 = append(f6, *f6elem)
		}
		res.Tags = f6
	}
	if r.ko.Spec.VisibilityConfig != nil {
		f7 := &svcsdktypes.VisibilityConfig{}
		if r.ko.Spec.VisibilityConfig.CloudWatchMetricsEnabled != nil {
			f7.CloudWatchMetricsEnabled = *r.ko.Spec.VisibilityConfig.CloudWatchMetricsEnabled
		}
		if r.ko.Spec.VisibilityConfig.MetricName != nil {
			f7.MetricName = r.ko.Spec.VisibilityConfig.MetricName
		}
		if r.ko.Spec.VisibilityConfig.SampledRequestsEnabled != nil {
			f7.SampledRequestsEnabled = *r.ko.Spec.VisibilityConfig.SampledRequestsEnabled
		}
		res.VisibilityConfig = f7
	}
`
	assert.Equal(
		expected,
		code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1),
	)
}
