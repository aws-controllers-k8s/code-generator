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
		ko.Spec.AuthorizationScopes = aws.StringSlice(resp.AuthorizationScopes)
	} else {
		ko.Spec.AuthorizationScopes = nil
	}
	if resp.AuthorizationType != "" {
		ko.Spec.AuthorizationType = aws.String(string(resp.AuthorizationType))
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
		ko.Spec.RequestModels = aws.StringMap(resp.RequestModels)
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
		ko.Spec.AuthorizationScopes = aws.StringSlice(resp.AuthorizationScopes)
	} else {
		ko.Spec.AuthorizationScopes = nil
	}
	if resp.AuthorizationType != "" {
		ko.Spec.AuthorizationType = aws.String(string(resp.AuthorizationType))
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
		ko.Spec.RequestModels = aws.StringMap(resp.RequestModels)
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


func TestSetResource_SageMaker_Domain_ReadOne(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sagemaker")

	crd := testutil.GetCRDByName(t, g, "Domain")
	require.NotNil(crd)

	expected := `
	if resp.AppNetworkAccessType != "" {
		ko.Spec.AppNetworkAccessType = aws.String(string(resp.AppNetworkAccessType))
	} else {
		ko.Spec.AppNetworkAccessType = nil
	}
	if resp.AppSecurityGroupManagement != "" {
		ko.Spec.AppSecurityGroupManagement = aws.String(string(resp.AppSecurityGroupManagement))
	} else {
		ko.Spec.AppSecurityGroupManagement = nil
	}
	if resp.AuthMode != "" {
		ko.Spec.AuthMode = aws.String(string(resp.AuthMode))
	} else {
		ko.Spec.AuthMode = nil
	}
	if resp.DefaultUserSettings != nil {
		f4 := &svcapitypes.UserSettings{}
		if resp.DefaultUserSettings.CodeEditorAppSettings != nil {
			f4f0 := &svcapitypes.CodeEditorAppSettings{}
			if resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec != nil {
				f4f0f0 := &svcapitypes.ResourceSpec{}
				if resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec.InstanceType != "" {
					f4f0f0.InstanceType = aws.String(string(resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec.InstanceType))
				}
				if resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec.LifecycleConfigArn != nil {
					f4f0f0.LifecycleConfigARN = resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec.LifecycleConfigArn
				}
				if resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec.SageMakerImageArn != nil {
					f4f0f0.SageMakerImageARN = resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec.SageMakerImageArn
				}
				if resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec.SageMakerImageVersionAlias != nil {
					f4f0f0.SageMakerImageVersionAlias = resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec.SageMakerImageVersionAlias
				}
				if resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec.SageMakerImageVersionArn != nil {
					f4f0f0.SageMakerImageVersionARN = resp.DefaultUserSettings.CodeEditorAppSettings.DefaultResourceSpec.SageMakerImageVersionArn
				}
				f4f0.DefaultResourceSpec = f4f0f0
			}
			if resp.DefaultUserSettings.CodeEditorAppSettings.LifecycleConfigArns != nil {
				f4f0.LifecycleConfigARNs = aws.StringSlice(resp.DefaultUserSettings.CodeEditorAppSettings.LifecycleConfigArns)
			}
			f4.CodeEditorAppSettings = f4f0
		}
		if resp.DefaultUserSettings.CustomFileSystemConfigs != nil {
			f4f1 := []*svcapitypes.CustomFileSystemConfig{}
			for _, f4f1iter := range resp.DefaultUserSettings.CustomFileSystemConfigs {
				f4f1elem := &svcapitypes.CustomFileSystemConfig{}
				switch f4f1iter.(type) {
				case *svcsdktypes.CustomFileSystemConfigMemberEFSFileSystemConfig:
					f4f1elemf0 := f4f1iter.(*svcsdktypes.CustomFileSystemConfigMemberEFSFileSystemConfig)
					if f4f1elemf0 != nil {
						f4f1elemf0f0 := &svcapitypes.EFSFileSystemConfig{}
						if f4f1elemf0.Value.FileSystemId != nil {
							f4f1elemf0f0.FileSystemID = f4f1elemf0.Value.FileSystemId
						}
						if f4f1elemf0.Value.FileSystemPath != nil {
							f4f1elemf0f0.FileSystemPath = f4f1elemf0.Value.FileSystemPath
						}
						f4f1elem.EFSFileSystemConfig = f4f1elemf0f0
					}
				case *svcsdktypes.CustomFileSystemConfigMemberFSxLustreFileSystemConfig:
					f4f1elemf1 := f4f1iter.(*svcsdktypes.CustomFileSystemConfigMemberFSxLustreFileSystemConfig)
					if f4f1elemf1 != nil {
						f4f1elemf1f1 := &svcapitypes.FSxLustreFileSystemConfig{}
						if f4f1elemf1.Value.FileSystemId != nil {
							f4f1elemf1f1.FileSystemID = f4f1elemf1.Value.FileSystemId
						}
						if f4f1elemf1.Value.FileSystemPath != nil {
							f4f1elemf1f1.FileSystemPath = f4f1elemf1.Value.FileSystemPath
						}
						f4f1elem.FSxLustreFileSystemConfig = f4f1elemf1f1
					}
				}
				f4f1 = append(f4f1, f4f1elem)
			}
			f4.CustomFileSystemConfigs = f4f1
		}
		if resp.DefaultUserSettings.CustomPosixUserConfig != nil {
			f4f2 := &svcapitypes.CustomPosixUserConfig{}
			if resp.DefaultUserSettings.CustomPosixUserConfig.Gid != nil {
				f4f2.GID = resp.DefaultUserSettings.CustomPosixUserConfig.Gid
			}
			if resp.DefaultUserSettings.CustomPosixUserConfig.Uid != nil {
				f4f2.UID = resp.DefaultUserSettings.CustomPosixUserConfig.Uid
			}
			f4.CustomPosixUserConfig = f4f2
		}
		if resp.DefaultUserSettings.DefaultLandingUri != nil {
			f4.DefaultLandingURI = resp.DefaultUserSettings.DefaultLandingUri
		}
		if resp.DefaultUserSettings.ExecutionRole != nil {
			f4.ExecutionRole = resp.DefaultUserSettings.ExecutionRole
		}
		if resp.DefaultUserSettings.JupyterLabAppSettings != nil {
			f4f5 := &svcapitypes.JupyterLabAppSettings{}
			if resp.DefaultUserSettings.JupyterLabAppSettings.CustomImages != nil {
				f4f5f0 := []*svcapitypes.CustomImage{}
				for _, f4f5f0iter := range resp.DefaultUserSettings.JupyterLabAppSettings.CustomImages {
					f4f5f0elem := &svcapitypes.CustomImage{}
					if f4f5f0iter.AppImageConfigName != nil {
						f4f5f0elem.AppImageConfigName = f4f5f0iter.AppImageConfigName
					}
					if f4f5f0iter.ImageName != nil {
						f4f5f0elem.ImageName = f4f5f0iter.ImageName
					}
					if f4f5f0iter.ImageVersionNumber != nil {
						imageVersionNumberCopy := int64(*f4f5f0iter.ImageVersionNumber)
						f4f5f0elem.ImageVersionNumber = &imageVersionNumberCopy
					}
					f4f5f0 = append(f4f5f0, f4f5f0elem)
				}
				f4f5.CustomImages = f4f5f0
			}
			if resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec != nil {
				f4f5f1 := &svcapitypes.ResourceSpec{}
				if resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec.InstanceType != "" {
					f4f5f1.InstanceType = aws.String(string(resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec.InstanceType))
				}
				if resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec.LifecycleConfigArn != nil {
					f4f5f1.LifecycleConfigARN = resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec.LifecycleConfigArn
				}
				if resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec.SageMakerImageArn != nil {
					f4f5f1.SageMakerImageARN = resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec.SageMakerImageArn
				}
				if resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec.SageMakerImageVersionAlias != nil {
					f4f5f1.SageMakerImageVersionAlias = resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec.SageMakerImageVersionAlias
				}
				if resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec.SageMakerImageVersionArn != nil {
					f4f5f1.SageMakerImageVersionARN = resp.DefaultUserSettings.JupyterLabAppSettings.DefaultResourceSpec.SageMakerImageVersionArn
				}
				f4f5.DefaultResourceSpec = f4f5f1
			}
			if resp.DefaultUserSettings.JupyterLabAppSettings.LifecycleConfigArns != nil {
				f4f5.LifecycleConfigARNs = aws.StringSlice(resp.DefaultUserSettings.JupyterLabAppSettings.LifecycleConfigArns)
			}
			f4.JupyterLabAppSettings = f4f5
		}
		if resp.DefaultUserSettings.JupyterServerAppSettings != nil {
			f4f6 := &svcapitypes.JupyterServerAppSettings{}
			if resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec != nil {
				f4f6f0 := &svcapitypes.ResourceSpec{}
				if resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType != "" {
					f4f6f0.InstanceType = aws.String(string(resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec.InstanceType))
				}
				if resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigArn != nil {
					f4f6f0.LifecycleConfigARN = resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec.LifecycleConfigArn
				}
				if resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageArn != nil {
					f4f6f0.SageMakerImageARN = resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageArn
				}
				if resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionAlias != nil {
					f4f6f0.SageMakerImageVersionAlias = resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionAlias
				}
				if resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionArn != nil {
					f4f6f0.SageMakerImageVersionARN = resp.DefaultUserSettings.JupyterServerAppSettings.DefaultResourceSpec.SageMakerImageVersionArn
				}
				f4f6.DefaultResourceSpec = f4f6f0
			}
			if resp.DefaultUserSettings.JupyterServerAppSettings.LifecycleConfigArns != nil {
				f4f6.LifecycleConfigARNs = aws.StringSlice(resp.DefaultUserSettings.JupyterServerAppSettings.LifecycleConfigArns)
			}
			f4.JupyterServerAppSettings = f4f6
		}
		if resp.DefaultUserSettings.KernelGatewayAppSettings != nil {
			f4f7 := &svcapitypes.KernelGatewayAppSettings{}
			if resp.DefaultUserSettings.KernelGatewayAppSettings.CustomImages != nil {
				f4f7f0 := []*svcapitypes.CustomImage{}
				for _, f4f7f0iter := range resp.DefaultUserSettings.KernelGatewayAppSettings.CustomImages {
					f4f7f0elem := &svcapitypes.CustomImage{}
					if f4f7f0iter.AppImageConfigName != nil {
						f4f7f0elem.AppImageConfigName = f4f7f0iter.AppImageConfigName
					}
					if f4f7f0iter.ImageName != nil {
						f4f7f0elem.ImageName = f4f7f0iter.ImageName
					}
					if f4f7f0iter.ImageVersionNumber != nil {
						imageVersionNumberCopy := int64(*f4f7f0iter.ImageVersionNumber)
						f4f7f0elem.ImageVersionNumber = &imageVersionNumberCopy
					}
					f4f7f0 = append(f4f7f0, f4f7f0elem)
				}
				f4f7.CustomImages = f4f7f0
			}
			if resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec != nil {
				f4f7f1 := &svcapitypes.ResourceSpec{}
				if resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType != "" {
					f4f7f1.InstanceType = aws.String(string(resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec.InstanceType))
				}
				if resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigArn != nil {
					f4f7f1.LifecycleConfigARN = resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec.LifecycleConfigArn
				}
				if resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageArn != nil {
					f4f7f1.SageMakerImageARN = resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageArn
				}
				if resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionAlias != nil {
					f4f7f1.SageMakerImageVersionAlias = resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionAlias
				}
				if resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionArn != nil {
					f4f7f1.SageMakerImageVersionARN = resp.DefaultUserSettings.KernelGatewayAppSettings.DefaultResourceSpec.SageMakerImageVersionArn
				}
				f4f7.DefaultResourceSpec = f4f7f1
			}
			if resp.DefaultUserSettings.KernelGatewayAppSettings.LifecycleConfigArns != nil {
				f4f7.LifecycleConfigARNs = aws.StringSlice(resp.DefaultUserSettings.KernelGatewayAppSettings.LifecycleConfigArns)
			}
			f4.KernelGatewayAppSettings = f4f7
		}
		if resp.DefaultUserSettings.RStudioServerProAppSettings != nil {
			f4f8 := &svcapitypes.RStudioServerProAppSettings{}
			if resp.DefaultUserSettings.RStudioServerProAppSettings.AccessStatus != "" {
				f4f8.AccessStatus = aws.String(string(resp.DefaultUserSettings.RStudioServerProAppSettings.AccessStatus))
			}
			if resp.DefaultUserSettings.RStudioServerProAppSettings.UserGroup != "" {
				f4f8.UserGroup = aws.String(string(resp.DefaultUserSettings.RStudioServerProAppSettings.UserGroup))
			}
			f4.RStudioServerProAppSettings = f4f8
		}
		if resp.DefaultUserSettings.SecurityGroups != nil {
			f4.SecurityGroups = aws.StringSlice(resp.DefaultUserSettings.SecurityGroups)
		}
		if resp.DefaultUserSettings.SharingSettings != nil {
			f4f10 := &svcapitypes.SharingSettings{}
			if resp.DefaultUserSettings.SharingSettings.NotebookOutputOption != "" {
				f4f10.NotebookOutputOption = aws.String(string(resp.DefaultUserSettings.SharingSettings.NotebookOutputOption))
			}
			if resp.DefaultUserSettings.SharingSettings.S3KmsKeyId != nil {
				f4f10.S3KMSKeyID = resp.DefaultUserSettings.SharingSettings.S3KmsKeyId
			}
			if resp.DefaultUserSettings.SharingSettings.S3OutputPath != nil {
				f4f10.S3OutputPath = resp.DefaultUserSettings.SharingSettings.S3OutputPath
			}
			f4.SharingSettings = f4f10
		}
		if resp.DefaultUserSettings.SpaceStorageSettings != nil {
			f4f11 := &svcapitypes.DefaultSpaceStorageSettings{}
			if resp.DefaultUserSettings.SpaceStorageSettings.DefaultEbsStorageSettings != nil {
				f4f11f0 := &svcapitypes.DefaultEBSStorageSettings{}
				if resp.DefaultUserSettings.SpaceStorageSettings.DefaultEbsStorageSettings.DefaultEbsVolumeSizeInGb != nil {
					defaultEBSVolumeSizeInGbCopy := int64(*resp.DefaultUserSettings.SpaceStorageSettings.DefaultEbsStorageSettings.DefaultEbsVolumeSizeInGb)
					f4f11f0.DefaultEBSVolumeSizeInGb = &defaultEBSVolumeSizeInGbCopy
				}
				if resp.DefaultUserSettings.SpaceStorageSettings.DefaultEbsStorageSettings.MaximumEbsVolumeSizeInGb != nil {
					maximumEBSVolumeSizeInGbCopy := int64(*resp.DefaultUserSettings.SpaceStorageSettings.DefaultEbsStorageSettings.MaximumEbsVolumeSizeInGb)
					f4f11f0.MaximumEBSVolumeSizeInGb = &maximumEBSVolumeSizeInGbCopy
				}
				f4f11.DefaultEBSStorageSettings = f4f11f0
			}
			f4.SpaceStorageSettings = f4f11
		}
		if resp.DefaultUserSettings.StudioWebPortal != "" {
			f4.StudioWebPortal = aws.String(string(resp.DefaultUserSettings.StudioWebPortal))
		}
		if resp.DefaultUserSettings.TensorBoardAppSettings != nil {
			f4f13 := &svcapitypes.TensorBoardAppSettings{}
			if resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec != nil {
				f4f13f0 := &svcapitypes.ResourceSpec{}
				if resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType != "" {
					f4f13f0.InstanceType = aws.String(string(resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec.InstanceType))
				}
				if resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigArn != nil {
					f4f13f0.LifecycleConfigARN = resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec.LifecycleConfigArn
				}
				if resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageArn != nil {
					f4f13f0.SageMakerImageARN = resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageArn
				}
				if resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionAlias != nil {
					f4f13f0.SageMakerImageVersionAlias = resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionAlias
				}
				if resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionArn != nil {
					f4f13f0.SageMakerImageVersionARN = resp.DefaultUserSettings.TensorBoardAppSettings.DefaultResourceSpec.SageMakerImageVersionArn
				}
				f4f13.DefaultResourceSpec = f4f13f0
			}
			f4.TensorBoardAppSettings = f4f13
		}
		ko.Spec.DefaultUserSettings = f4
	} else {
		ko.Spec.DefaultUserSettings = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.DomainArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.DomainArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.DomainId != nil {
		ko.Status.DomainID = resp.DomainId
	} else {
		ko.Status.DomainID = nil
	}
	if resp.DomainName != nil {
		ko.Spec.DomainName = resp.DomainName
	} else {
		ko.Spec.DomainName = nil
	}
	if resp.DomainSettings != nil {
		f8 := &svcapitypes.DomainSettings{}
		if resp.DomainSettings.DockerSettings != nil {
			f8f0 := &svcapitypes.DockerSettings{}
			if resp.DomainSettings.DockerSettings.EnableDockerAccess != "" {
				f8f0.EnableDockerAccess = aws.String(string(resp.DomainSettings.DockerSettings.EnableDockerAccess))
			}
			if resp.DomainSettings.DockerSettings.VpcOnlyTrustedAccounts != nil {
				f8f0.VPCOnlyTrustedAccounts = aws.StringSlice(resp.DomainSettings.DockerSettings.VpcOnlyTrustedAccounts)
			}
			f8.DockerSettings = f8f0
		}
		if resp.DomainSettings.RStudioServerProDomainSettings != nil {
			f8f1 := &svcapitypes.RStudioServerProDomainSettings{}
			if resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec != nil {
				f8f1f0 := &svcapitypes.ResourceSpec{}
				if resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec.InstanceType != "" {
					f8f1f0.InstanceType = aws.String(string(resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec.InstanceType))
				}
				if resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec.LifecycleConfigArn != nil {
					f8f1f0.LifecycleConfigARN = resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec.LifecycleConfigArn
				}
				if resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec.SageMakerImageArn != nil {
					f8f1f0.SageMakerImageARN = resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec.SageMakerImageArn
				}
				if resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec.SageMakerImageVersionAlias != nil {
					f8f1f0.SageMakerImageVersionAlias = resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec.SageMakerImageVersionAlias
				}
				if resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec.SageMakerImageVersionArn != nil {
					f8f1f0.SageMakerImageVersionARN = resp.DomainSettings.RStudioServerProDomainSettings.DefaultResourceSpec.SageMakerImageVersionArn
				}
				f8f1.DefaultResourceSpec = f8f1f0
			}
			if resp.DomainSettings.RStudioServerProDomainSettings.DomainExecutionRoleArn != nil {
				f8f1.DomainExecutionRoleARN = resp.DomainSettings.RStudioServerProDomainSettings.DomainExecutionRoleArn
			}
			if resp.DomainSettings.RStudioServerProDomainSettings.RStudioConnectUrl != nil {
				f8f1.RStudioConnectURL = resp.DomainSettings.RStudioServerProDomainSettings.RStudioConnectUrl
			}
			if resp.DomainSettings.RStudioServerProDomainSettings.RStudioPackageManagerUrl != nil {
				f8f1.RStudioPackageManagerURL = resp.DomainSettings.RStudioServerProDomainSettings.RStudioPackageManagerUrl
			}
			f8.RStudioServerProDomainSettings = f8f1
		}
		if resp.DomainSettings.SecurityGroupIds != nil {
			f8.SecurityGroupIDs = aws.StringSlice(resp.DomainSettings.SecurityGroupIds)
		}
		ko.Spec.DomainSettings = f8
	} else {
		ko.Spec.DomainSettings = nil
	}
	if resp.HomeEfsFileSystemKmsKeyId != nil {
		ko.Spec.HomeEFSFileSystemKMSKeyID = resp.HomeEfsFileSystemKmsKeyId
	} else {
		ko.Spec.HomeEFSFileSystemKMSKeyID = nil
	}
	if resp.KmsKeyId != nil {
		ko.Spec.KMSKeyID = resp.KmsKeyId
	} else {
		ko.Spec.KMSKeyID = nil
	}
	if resp.Status != "" {
		ko.Status.Status = aws.String(string(resp.Status))
	} else {
		ko.Status.Status = nil
	}
	if resp.SubnetIds != nil {
		ko.Spec.SubnetIDs = aws.StringSlice(resp.SubnetIds)
	} else {
		ko.Spec.SubnetIDs = nil
	}
	if resp.Url != nil {
		ko.Status.URL = resp.Url
	} else {
		ko.Status.URL = nil
	}
	if resp.VpcId != nil {
		ko.Spec.VPCID = resp.VpcId
	} else {
		ko.Spec.VPCID = nil
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
	if resp.BackupDescription.BackupDetails.BackupStatus != "" {
		ko.Status.BackupStatus = aws.String(string(resp.BackupDescription.BackupDetails.BackupStatus))
	} else {
		ko.Status.BackupStatus = nil
	}
	if resp.BackupDescription.BackupDetails.BackupType != "" {
		ko.Status.BackupType = aws.String(string(resp.BackupDescription.BackupDetails.BackupType))
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
			if f1iter.AttributeType != "" {
				f1elem.AttributeType = aws.String(string(f1iter.AttributeType))
			}
			f1 = append(f1, f1elem)
		}
		ko.Spec.AttributeDefinitions = f1
	} else {
		ko.Spec.AttributeDefinitions = nil
	}
	if resp.Table.BillingModeSummary != nil {
		f2 := &svcapitypes.BillingModeSummary{}
		if resp.Table.BillingModeSummary.BillingMode != "" {
			f2.BillingMode = aws.String(string(resp.Table.BillingModeSummary.BillingMode))
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
					if f4elemf6iter.KeyType != "" {
						f4elemf6elem.KeyType = aws.String(string(f4elemf6iter.KeyType))
					}
					f4elemf6 = append(f4elemf6, f4elemf6elem)
				}
				f4elem.KeySchema = f4elemf6
			}
			if f4iter.Projection != nil {
				f4elemf8 := &svcapitypes.Projection{}
				if f4iter.Projection.NonKeyAttributes != nil {
					f4elemf8.NonKeyAttributes = aws.StringSlice(f4iter.Projection.NonKeyAttributes)
				}
				if f4iter.Projection.ProjectionType != "" {
					f4elemf8.ProjectionType = aws.String(string(f4iter.Projection.ProjectionType))
				}
				f4elem.Projection = f4elemf8
			}
			if f4iter.ProvisionedThroughput != nil {
				f4elemf9 := &svcapitypes.ProvisionedThroughput{}
				if f4iter.ProvisionedThroughput.ReadCapacityUnits != nil {
					f4elemf9.ReadCapacityUnits = f4iter.ProvisionedThroughput.ReadCapacityUnits
				}
				if f4iter.ProvisionedThroughput.WriteCapacityUnits != nil {
					f4elemf9.WriteCapacityUnits = f4iter.ProvisionedThroughput.WriteCapacityUnits
				}
				f4elem.ProvisionedThroughput = f4elemf9
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
			if f7iter.KeyType != "" {
				f7elem.KeyType = aws.String(string(f7iter.KeyType))
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
					if f10elemf4iter.KeyType != "" {
						f10elemf4elem.KeyType = aws.String(string(f10elemf4iter.KeyType))
					}
					f10elemf4 = append(f10elemf4, f10elemf4elem)
				}
				f10elem.KeySchema = f10elemf4
			}
			if f10iter.Projection != nil {
				f10elemf5 := &svcapitypes.Projection{}
				if f10iter.Projection.NonKeyAttributes != nil {
					f10elemf5.NonKeyAttributes = aws.StringSlice(f10iter.Projection.NonKeyAttributes)
				}
				if f10iter.Projection.ProjectionType != "" {
					f10elemf5.ProjectionType = aws.String(string(f10iter.Projection.ProjectionType))
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
			if f12iter.ReplicaStatus != "" {
				f12elem.ReplicaStatus = aws.String(string(f12iter.ReplicaStatus))
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
		if resp.Table.SSEDescription.SSEType != "" {
			f14.SSEType = aws.String(string(resp.Table.SSEDescription.SSEType))
		}
		if resp.Table.SSEDescription.Status != "" {
			f14.Status = aws.String(string(resp.Table.SSEDescription.Status))
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
		if resp.Table.StreamSpecification.StreamViewType != "" {
			f15.StreamViewType = aws.String(string(resp.Table.StreamSpecification.StreamViewType))
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
	if resp.Table.TableStatus != "" {
		ko.Status.TableStatus = aws.String(string(resp.Table.TableStatus))
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
		f1.ScanOnPush = &resp.Repository.ImageScanningConfiguration.ScanOnPush
		ko.Spec.ImageScanningConfiguration = f1
	} else {
		ko.Spec.ImageScanningConfiguration = nil
	}
	if resp.Repository.ImageTagMutability != "" {
		ko.Spec.ImageTagMutability = aws.String(string(resp.Repository.ImageTagMutability))
	} else {
		ko.Spec.ImageTagMutability = nil
	}
	if resp.Repository.RegistryId != nil {
		ko.Spec.RegistryID = resp.Repository.RegistryId
	} else {
		ko.Spec.RegistryID = nil
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
			f1.ScanOnPush = &elem.ImageScanningConfiguration.ScanOnPush
			ko.Spec.ImageScanningConfiguration = f1
		} else {
			ko.Spec.ImageScanningConfiguration = nil
		}
		if elem.ImageTagMutability != "" {
			ko.Spec.ImageTagMutability = aws.String(string(elem.ImageTagMutability))
		} else {
			ko.Spec.ImageTagMutability = nil
		}
		if elem.RegistryId != nil {
			ko.Spec.RegistryID = elem.RegistryId
		} else {
			ko.Spec.RegistryID = nil
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
	f0, ok := resp.Attributes["DeliveryPolicy"]
	if ok {
		ko.Spec.DeliveryPolicy = &f0
	} else {
		ko.Spec.DeliveryPolicy = nil
	}
	f1, ok := resp.Attributes["DisplayName"]
	if ok {
		ko.Spec.DisplayName = &f1
	} else {
		ko.Spec.DisplayName = nil
	}
	f2, ok := resp.Attributes["EffectiveDeliveryPolicy"]
	if ok {
		ko.Status.EffectiveDeliveryPolicy = &f2
	} else {
		ko.Status.EffectiveDeliveryPolicy = nil
	}
	f3, ok := resp.Attributes["KmsMasterKeyId"]
	if ok {
		ko.Spec.KMSMasterKeyID = &f3
	} else {
		ko.Spec.KMSMasterKeyID = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	tmpOwnerID := ackv1alpha1.AWSAccountID(resp.Attributes["Owner"])
	ko.Status.ACKResourceMetadata.OwnerAccountID = &tmpOwnerID
	f5, ok := resp.Attributes["Policy"]
	if ok {
		ko.Spec.Policy = &f5
	} else {
		ko.Spec.Policy = nil
	}
	tmpARN := ackv1alpha1.AWSResourceName(resp.Attributes["TopicArn"])
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
	f0, ok := resp.Attributes["ContentBasedDeduplication"]
	if ok {
		ko.Spec.ContentBasedDeduplication = &f0
	} else {
		ko.Spec.ContentBasedDeduplication = nil
	}
	f1, ok := resp.Attributes["CreatedTimestamp"]
	if ok {
		ko.Spec.CreatedTimestamp = &f1
	} else {
		ko.Spec.CreatedTimestamp = nil
	}
	f2, ok := resp.Attributes["DelaySeconds"]
	if ok {
		ko.Spec.DelaySeconds = &f2
	} else {
		ko.Spec.DelaySeconds = nil
	}
	f3, ok := resp.Attributes["FifoQueue"]
	if ok {
		ko.Spec.FIFOQueue = &f3
	} else {
		ko.Spec.FIFOQueue = nil
	}
	f4, ok := resp.Attributes["KmsDataKeyReusePeriodSeconds"]
	if ok {
		ko.Spec.KMSDataKeyReusePeriodSeconds = &f4
	} else {
		ko.Spec.KMSDataKeyReusePeriodSeconds = nil
	}
	f5, ok := resp.Attributes["KmsMasterKeyId"]
	if ok {
		ko.Spec.KMSMasterKeyID = &f5
	} else {
		ko.Spec.KMSMasterKeyID = nil
	}
	f6, ok := resp.Attributes["MaximumMessageSize"]
	if ok {
		ko.Spec.MaximumMessageSize = &f6
	} else {
		ko.Spec.MaximumMessageSize = nil
	}
	f7, ok := resp.Attributes["MessageRetentionPeriod"]
	if ok {
		ko.Spec.MessageRetentionPeriod = &f7
	} else {
		ko.Spec.MessageRetentionPeriod = nil
	}
	f8, ok := resp.Attributes["Policy"]
	if ok {
		ko.Spec.Policy = &f8
	} else {
		ko.Spec.Policy = nil
	}
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	tmpARN := ackv1alpha1.AWSResourceName(resp.Attributes["QueueArn"])
	ko.Status.ACKResourceMetadata.ARN = &tmpARN
	f10, ok := resp.Attributes["ReceiveMessageWaitTimeSeconds"]
	if ok {
		ko.Spec.ReceiveMessageWaitTimeSeconds = &f10
	} else {
		ko.Spec.ReceiveMessageWaitTimeSeconds = nil
	}
	f11, ok := resp.Attributes["RedrivePolicy"]
	if ok {
		ko.Spec.RedrivePolicy = &f11
	} else {
		ko.Spec.RedrivePolicy = nil
	}
	f12, ok := resp.Attributes["VisibilityTimeout"]
	if ok {
		ko.Spec.VisibilityTimeout = &f12
	} else {
		ko.Spec.VisibilityTimeout = nil
	}
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
		if elem.SupportedNetworkTypes != nil {
			ko.Status.SupportedNetworkTypes = aws.StringSlice(elem.SupportedNetworkTypes)
		} else {
			ko.Status.SupportedNetworkTypes = nil
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

func TestSetResource_SNS_Topics_SetResourceIdentifiers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sns")

	crd := testutil.GetCRDByName(t, g, "Topic")
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

func TestSetResource_SQS_Queues_SetResourceIdentifiers(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sqs")

	crd := testutil.GetCRDByName(t, g, "Queue")
	require.NotNil(crd)

	expected := `
	if identifier.NameOrID == "" {
		return ackerrors.MissingNameIdentifier
	}
	r.ko.Status.QueueURL = &identifier.NameOrID

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
		r.ko.Spec.DomainName = aws.String(f1)
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
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.SecurityGroupArn != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.SecurityGroupArn)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.Tags != nil {
		f2 := []*svcapitypes.Tag{}
		for _, f2iter := range resp.Tags {
			f2elem := &svcapitypes.Tag{}
			if f2iter.Key != nil {
				f2elem.Key = f2iter.Key
			}
			if f2iter.Value != nil {
				f2elem.Value = f2iter.Value
			}
			f2 = append(f2, f2elem)
		}
		ko.Status.Tags = f2
	} else {
		ko.Status.Tags = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, model.OpTypeCreate, "resp", "ko", 1),
	)
}

func TestSetResource_EKS_Cluster_PopulateResourceFromAnnotation(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "eks")

	crd := testutil.GetCRDByName(t, g, "Cluster")
	require.NotNil(crd)

	expected := `
	tmp, ok := fields["name"]
	if !ok {
		return ackerrors.NewTerminalError(fmt.Errorf("required field missing: name"))
	}
	r.ko.Spec.Name = &tmp

`
	assert.Equal(
		expected,
		code.PopulateResourceFromAnnotation(crd.Config(), crd, "fields", "r.ko", 1),
	)
}

func TestSetResource_SageMaker_ModelPackage_PopulateResourceFromAnnotation(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sagemaker")

	crd := testutil.GetCRDByName(t, g, "ModelPackage")
	require.NotNil(crd)

	expected := `
	tmp, ok := identifier["arn"]
	if !ok {
		return ackerrors.NewTerminalError(fmt.Errorf("required field missing: arn"))
	}

	if r.ko.Status.ACKResourceMetadata == nil {
		r.ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	arn := ackv1alpha1.AWSResourceName(tmp)
	r.ko.Status.ACKResourceMetadata.ARN = &arn
`
	assert.Equal(
		expected,
		code.PopulateResourceFromAnnotation(crd.Config(), crd, "identifier", "r.ko", 1),
	)
}

func TestSetResource_APIGWV2_ApiMapping_PopulateResourceFromAnnotation(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "ApiMapping")
	require.NotNil(crd)

	expected := `
	tmp, ok := fields["apiMappingID"]
	if !ok {
		return ackerrors.NewTerminalError(fmt.Errorf("required field missing: apiMappingID"))
	}
	r.ko.Status.APIMappingID = &tmp

	f1, f1ok := fields["domainName"]
	if f1ok {
		r.ko.Spec.DomainName = aws.String(f1)
	}
`
	assert.Equal(
		expected,
		code.PopulateResourceFromAnnotation(crd.Config(), crd, "fields", "r.ko", 1),
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
		maxSessionDurationCopy := int64(*resp.Role.MaxSessionDuration)
		ko.Spec.MaxSessionDuration = &maxSessionDurationCopy
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
					var f0elemf1elem *string
					if f0elemf1iter.Value != nil {
						f0elemf1elem = f0elemf1iter.Value
					}
					f0elemf1 = append(f0elemf1, f0elemf1elem)
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

func TestSetResource_EC2_Instance_Create(t *testing.T) {
	// Check that the RunInstances output (Reservation)
	// uses the first element of the returned list of Instances
	// to populate Instance CR
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")
	op := model.OpTypeCreate

	crd := testutil.GetCRDByName(t, g, "Instance")
	require.NotNil(crd)

	expected := `
	found := false
	for _, elem := range resp.Instances {
		if elem.AmiLaunchIndex != nil {
			amiLaunchIndexCopy := int64(*elem.AmiLaunchIndex)
			ko.Status.AMILaunchIndex = &amiLaunchIndexCopy
		} else {
			ko.Status.AMILaunchIndex = nil
		}
		if elem.Architecture != "" {
			ko.Status.Architecture = aws.String(string(elem.Architecture))
		} else {
			ko.Status.Architecture = nil
		}
		if elem.BlockDeviceMappings != nil {
			f2 := []*svcapitypes.BlockDeviceMapping{}
			for _, f2iter := range elem.BlockDeviceMappings {
				f2elem := &svcapitypes.BlockDeviceMapping{}
				if f2iter.DeviceName != nil {
					f2elem.DeviceName = f2iter.DeviceName
				}
				if f2iter.Ebs != nil {
					f2elemf1 := &svcapitypes.EBSBlockDevice{}
					if f2iter.Ebs.DeleteOnTermination != nil {
						f2elemf1.DeleteOnTermination = f2iter.Ebs.DeleteOnTermination
					}
					f2elem.EBS = f2elemf1
				}
				f2 = append(f2, f2elem)
			}
			ko.Spec.BlockDeviceMappings = f2
		} else {
			ko.Spec.BlockDeviceMappings = nil
		}
		if elem.BootMode != "" {
			ko.Status.BootMode = aws.String(string(elem.BootMode))
		} else {
			ko.Status.BootMode = nil
		}
		if elem.CapacityReservationId != nil {
			ko.Status.CapacityReservationID = elem.CapacityReservationId
		} else {
			ko.Status.CapacityReservationID = nil
		}
		if elem.CapacityReservationSpecification != nil {
			f5 := &svcapitypes.CapacityReservationSpecification{}
			if elem.CapacityReservationSpecification.CapacityReservationPreference != "" {
				f5.CapacityReservationPreference = aws.String(string(elem.CapacityReservationSpecification.CapacityReservationPreference))
			}
			if elem.CapacityReservationSpecification.CapacityReservationTarget != nil {
				f5f1 := &svcapitypes.CapacityReservationTarget{}
				if elem.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationId != nil {
					f5f1.CapacityReservationID = elem.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationId
				}
				if elem.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationResourceGroupArn != nil {
					f5f1.CapacityReservationResourceGroupARN = elem.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationResourceGroupArn
				}
				f5.CapacityReservationTarget = f5f1
			}
			ko.Spec.CapacityReservationSpecification = f5
		} else {
			ko.Spec.CapacityReservationSpecification = nil
		}
		if elem.CpuOptions != nil {
			f6 := &svcapitypes.CPUOptionsRequest{}
			if elem.CpuOptions.CoreCount != nil {
				coreCountCopy := int64(*elem.CpuOptions.CoreCount)
				f6.CoreCount = &coreCountCopy
			}
			if elem.CpuOptions.ThreadsPerCore != nil {
				threadsPerCoreCopy := int64(*elem.CpuOptions.ThreadsPerCore)
				f6.ThreadsPerCore = &threadsPerCoreCopy
			}
			ko.Spec.CPUOptions = f6
		} else {
			ko.Spec.CPUOptions = nil
		}
		if elem.EbsOptimized != nil {
			ko.Spec.EBSOptimized = elem.EbsOptimized
		} else {
			ko.Spec.EBSOptimized = nil
		}
		if elem.ElasticGpuAssociations != nil {
			f8 := []*svcapitypes.ElasticGPUAssociation{}
			for _, f8iter := range elem.ElasticGpuAssociations {
				f8elem := &svcapitypes.ElasticGPUAssociation{}
				if f8iter.ElasticGpuAssociationId != nil {
					f8elem.ElasticGPUAssociationID = f8iter.ElasticGpuAssociationId
				}
				if f8iter.ElasticGpuAssociationState != nil {
					f8elem.ElasticGPUAssociationState = f8iter.ElasticGpuAssociationState
				}
				if f8iter.ElasticGpuAssociationTime != nil {
					f8elem.ElasticGPUAssociationTime = f8iter.ElasticGpuAssociationTime
				}
				if f8iter.ElasticGpuId != nil {
					f8elem.ElasticGPUID = f8iter.ElasticGpuId
				}
				f8 = append(f8, f8elem)
			}
			ko.Status.ElasticGPUAssociations = f8
		} else {
			ko.Status.ElasticGPUAssociations = nil
		}
		if elem.ElasticInferenceAcceleratorAssociations != nil {
			f9 := []*svcapitypes.ElasticInferenceAcceleratorAssociation{}
			for _, f9iter := range elem.ElasticInferenceAcceleratorAssociations {
				f9elem := &svcapitypes.ElasticInferenceAcceleratorAssociation{}
				if f9iter.ElasticInferenceAcceleratorArn != nil {
					f9elem.ElasticInferenceAcceleratorARN = f9iter.ElasticInferenceAcceleratorArn
				}
				if f9iter.ElasticInferenceAcceleratorAssociationId != nil {
					f9elem.ElasticInferenceAcceleratorAssociationID = f9iter.ElasticInferenceAcceleratorAssociationId
				}
				if f9iter.ElasticInferenceAcceleratorAssociationState != nil {
					f9elem.ElasticInferenceAcceleratorAssociationState = f9iter.ElasticInferenceAcceleratorAssociationState
				}
				if f9iter.ElasticInferenceAcceleratorAssociationTime != nil {
					f9elem.ElasticInferenceAcceleratorAssociationTime = &metav1.Time{*f9iter.ElasticInferenceAcceleratorAssociationTime}
				}
				f9 = append(f9, f9elem)
			}
			ko.Status.ElasticInferenceAcceleratorAssociations = f9
		} else {
			ko.Status.ElasticInferenceAcceleratorAssociations = nil
		}
		if elem.EnaSupport != nil {
			ko.Status.ENASupport = elem.EnaSupport
		} else {
			ko.Status.ENASupport = nil
		}
		if elem.EnclaveOptions != nil {
			f11 := &svcapitypes.EnclaveOptionsRequest{}
			if elem.EnclaveOptions.Enabled != nil {
				f11.Enabled = elem.EnclaveOptions.Enabled
			}
			ko.Spec.EnclaveOptions = f11
		} else {
			ko.Spec.EnclaveOptions = nil
		}
		if elem.HibernationOptions != nil {
			f12 := &svcapitypes.HibernationOptionsRequest{}
			if elem.HibernationOptions.Configured != nil {
				f12.Configured = elem.HibernationOptions.Configured
			}
			ko.Spec.HibernationOptions = f12
		} else {
			ko.Spec.HibernationOptions = nil
		}
		if elem.Hypervisor != "" {
			ko.Status.Hypervisor = aws.String(string(elem.Hypervisor))
		} else {
			ko.Status.Hypervisor = nil
		}
		if elem.IamInstanceProfile != nil {
			f14 := &svcapitypes.IAMInstanceProfileSpecification{}
			if elem.IamInstanceProfile.Arn != nil {
				f14.ARN = elem.IamInstanceProfile.Arn
			}
			ko.Spec.IAMInstanceProfile = f14
		} else {
			ko.Spec.IAMInstanceProfile = nil
		}
		if elem.ImageId != nil {
			ko.Spec.ImageID = elem.ImageId
		} else {
			ko.Spec.ImageID = nil
		}
		if elem.InstanceId != nil {
			ko.Status.InstanceID = elem.InstanceId
		} else {
			ko.Status.InstanceID = nil
		}
		if elem.InstanceLifecycle != "" {
			ko.Status.InstanceLifecycle = aws.String(string(elem.InstanceLifecycle))
		} else {
			ko.Status.InstanceLifecycle = nil
		}
		if elem.InstanceType != "" {
			ko.Spec.InstanceType = aws.String(string(elem.InstanceType))
		} else {
			ko.Spec.InstanceType = nil
		}
		if elem.KernelId != nil {
			ko.Spec.KernelID = elem.KernelId
		} else {
			ko.Spec.KernelID = nil
		}
		if elem.KeyName != nil {
			ko.Spec.KeyName = elem.KeyName
		} else {
			ko.Spec.KeyName = nil
		}
		if elem.LaunchTime != nil {
			ko.Status.LaunchTime = &metav1.Time{*elem.LaunchTime}
		} else {
			ko.Status.LaunchTime = nil
		}
		if elem.Licenses != nil {
			f22 := []*svcapitypes.LicenseConfiguration{}
			for _, f22iter := range elem.Licenses {
				f22elem := &svcapitypes.LicenseConfiguration{}
				if f22iter.LicenseConfigurationArn != nil {
					f22elem.LicenseConfigurationARN = f22iter.LicenseConfigurationArn
				}
				f22 = append(f22, f22elem)
			}
			ko.Status.Licenses = f22
		} else {
			ko.Status.Licenses = nil
		}
		if elem.MetadataOptions != nil {
			f23 := &svcapitypes.InstanceMetadataOptionsRequest{}
			if elem.MetadataOptions.HttpEndpoint != "" {
				f23.HTTPEndpoint = aws.String(string(elem.MetadataOptions.HttpEndpoint))
			}
			if elem.MetadataOptions.HttpProtocolIpv6 != "" {
				f23.HTTPProtocolIPv6 = aws.String(string(elem.MetadataOptions.HttpProtocolIpv6))
			}
			if elem.MetadataOptions.HttpPutResponseHopLimit != nil {
				httpPutResponseHopLimitCopy := int64(*elem.MetadataOptions.HttpPutResponseHopLimit)
				f23.HTTPPutResponseHopLimit = &httpPutResponseHopLimitCopy
			}
			if elem.MetadataOptions.HttpTokens != "" {
				f23.HTTPTokens = aws.String(string(elem.MetadataOptions.HttpTokens))
			}
			ko.Spec.MetadataOptions = f23
		} else {
			ko.Spec.MetadataOptions = nil
		}
		if elem.Monitoring != nil {
			f24 := &svcapitypes.RunInstancesMonitoringEnabled{}
			ko.Spec.Monitoring = f24
		} else {
			ko.Spec.Monitoring = nil
		}
		if elem.NetworkInterfaces != nil {
			f25 := []*svcapitypes.InstanceNetworkInterfaceSpecification{}
			for _, f25iter := range elem.NetworkInterfaces {
				f25elem := &svcapitypes.InstanceNetworkInterfaceSpecification{}
				if f25iter.Description != nil {
					f25elem.Description = f25iter.Description
				}
				if f25iter.InterfaceType != nil {
					f25elem.InterfaceType = f25iter.InterfaceType
				}
				if f25iter.Ipv4Prefixes != nil {
					f25elemf6 := []*svcapitypes.IPv4PrefixSpecificationRequest{}
					for _, f25elemf6iter := range f25iter.Ipv4Prefixes {
						f25elemf6elem := &svcapitypes.IPv4PrefixSpecificationRequest{}
						if f25elemf6iter.Ipv4Prefix != nil {
							f25elemf6elem.IPv4Prefix = f25elemf6iter.Ipv4Prefix
						}
						f25elemf6 = append(f25elemf6, f25elemf6elem)
					}
					f25elem.IPv4Prefixes = f25elemf6
				}
				if f25iter.Ipv6Addresses != nil {
					f25elemf7 := []*svcapitypes.InstanceIPv6Address{}
					for _, f25elemf7iter := range f25iter.Ipv6Addresses {
						f25elemf7elem := &svcapitypes.InstanceIPv6Address{}
						if f25elemf7iter.Ipv6Address != nil {
							f25elemf7elem.IPv6Address = f25elemf7iter.Ipv6Address
						}
						f25elemf7 = append(f25elemf7, f25elemf7elem)
					}
					f25elem.IPv6Addresses = f25elemf7
				}
				if f25iter.Ipv6Prefixes != nil {
					f25elemf8 := []*svcapitypes.IPv6PrefixSpecificationRequest{}
					for _, f25elemf8iter := range f25iter.Ipv6Prefixes {
						f25elemf8elem := &svcapitypes.IPv6PrefixSpecificationRequest{}
						if f25elemf8iter.Ipv6Prefix != nil {
							f25elemf8elem.IPv6Prefix = f25elemf8iter.Ipv6Prefix
						}
						f25elemf8 = append(f25elemf8, f25elemf8elem)
					}
					f25elem.IPv6Prefixes = f25elemf8
				}
				if f25iter.NetworkInterfaceId != nil {
					f25elem.NetworkInterfaceID = f25iter.NetworkInterfaceId
				}
				if f25iter.PrivateIpAddress != nil {
					f25elem.PrivateIPAddress = f25iter.PrivateIpAddress
				}
				if f25iter.PrivateIpAddresses != nil {
					f25elemf15 := []*svcapitypes.PrivateIPAddressSpecification{}
					for _, f25elemf15iter := range f25iter.PrivateIpAddresses {
						f25elemf15elem := &svcapitypes.PrivateIPAddressSpecification{}
						if f25elemf15iter.Primary != nil {
							f25elemf15elem.Primary = f25elemf15iter.Primary
						}
						if f25elemf15iter.PrivateIpAddress != nil {
							f25elemf15elem.PrivateIPAddress = f25elemf15iter.PrivateIpAddress
						}
						f25elemf15 = append(f25elemf15, f25elemf15elem)
					}
					f25elem.PrivateIPAddresses = f25elemf15
				}
				if f25iter.SubnetId != nil {
					f25elem.SubnetID = f25iter.SubnetId
				}
				f25 = append(f25, f25elem)
			}
			ko.Spec.NetworkInterfaces = f25
		} else {
			ko.Spec.NetworkInterfaces = nil
		}
		if elem.OutpostArn != nil {
			ko.Status.OutpostARN = elem.OutpostArn
		} else {
			ko.Status.OutpostARN = nil
		}
		if elem.Placement != nil {
			f27 := &svcapitypes.Placement{}
			if elem.Placement.Affinity != nil {
				f27.Affinity = elem.Placement.Affinity
			}
			if elem.Placement.AvailabilityZone != nil {
				f27.AvailabilityZone = elem.Placement.AvailabilityZone
			}
			if elem.Placement.GroupName != nil {
				f27.GroupName = elem.Placement.GroupName
			}
			if elem.Placement.HostId != nil {
				f27.HostID = elem.Placement.HostId
			}
			if elem.Placement.HostResourceGroupArn != nil {
				f27.HostResourceGroupARN = elem.Placement.HostResourceGroupArn
			}
			if elem.Placement.PartitionNumber != nil {
				partitionNumberCopy := int64(*elem.Placement.PartitionNumber)
				f27.PartitionNumber = &partitionNumberCopy
			}
			if elem.Placement.SpreadDomain != nil {
				f27.SpreadDomain = elem.Placement.SpreadDomain
			}
			if elem.Placement.Tenancy != "" {
				f27.Tenancy = aws.String(string(elem.Placement.Tenancy))
			}
			ko.Spec.Placement = f27
		} else {
			ko.Spec.Placement = nil
		}
		if elem.Platform != "" {
			ko.Status.Platform = aws.String(string(elem.Platform))
		} else {
			ko.Status.Platform = nil
		}
		if elem.PlatformDetails != nil {
			ko.Status.PlatformDetails = elem.PlatformDetails
		} else {
			ko.Status.PlatformDetails = nil
		}
		if elem.PrivateDnsName != nil {
			ko.Status.PrivateDNSName = elem.PrivateDnsName
		} else {
			ko.Status.PrivateDNSName = nil
		}
		if elem.PrivateIpAddress != nil {
			ko.Spec.PrivateIPAddress = elem.PrivateIpAddress
		} else {
			ko.Spec.PrivateIPAddress = nil
		}
		if elem.ProductCodes != nil {
			f32 := []*svcapitypes.ProductCode{}
			for _, f32iter := range elem.ProductCodes {
				f32elem := &svcapitypes.ProductCode{}
				if f32iter.ProductCodeId != nil {
					f32elem.ProductCodeID = f32iter.ProductCodeId
				}
				if f32iter.ProductCodeType != "" {
					f32elem.ProductCodeType = aws.String(string(f32iter.ProductCodeType))
				}
				f32 = append(f32, f32elem)
			}
			ko.Status.ProductCodes = f32
		} else {
			ko.Status.ProductCodes = nil
		}
		if elem.PublicDnsName != nil {
			ko.Status.PublicDNSName = elem.PublicDnsName
		} else {
			ko.Status.PublicDNSName = nil
		}
		if elem.PublicIpAddress != nil {
			ko.Status.PublicIPAddress = elem.PublicIpAddress
		} else {
			ko.Status.PublicIPAddress = nil
		}
		if elem.RamdiskId != nil {
			ko.Spec.RAMDiskID = elem.RamdiskId
		} else {
			ko.Spec.RAMDiskID = nil
		}
		if elem.RootDeviceName != nil {
			ko.Status.RootDeviceName = elem.RootDeviceName
		} else {
			ko.Status.RootDeviceName = nil
		}
		if elem.RootDeviceType != "" {
			ko.Status.RootDeviceType = aws.String(string(elem.RootDeviceType))
		} else {
			ko.Status.RootDeviceType = nil
		}
		if elem.SecurityGroups != nil {
			f38 := []*string{}
			for _, f38iter := range elem.SecurityGroups {
				var f38elem *string
				f38elem = f38iter.GroupName
				f38 = append(f38, f38elem)
			}
			ko.Spec.SecurityGroups = f38
		} else {
			ko.Spec.SecurityGroups = nil
		}
		if elem.SourceDestCheck != nil {
			ko.Status.SourceDestCheck = elem.SourceDestCheck
		} else {
			ko.Status.SourceDestCheck = nil
		}
		if elem.SpotInstanceRequestId != nil {
			ko.Status.SpotInstanceRequestID = elem.SpotInstanceRequestId
		} else {
			ko.Status.SpotInstanceRequestID = nil
		}
		if elem.SriovNetSupport != nil {
			ko.Status.SRIOVNetSupport = elem.SriovNetSupport
		} else {
			ko.Status.SRIOVNetSupport = nil
		}
		if elem.State != nil {
			f42 := &svcapitypes.InstanceState{}
			if elem.State.Code != nil {
				codeCopy := int64(*elem.State.Code)
				f42.Code = &codeCopy
			}
			if elem.State.Name != "" {
				f42.Name = aws.String(string(elem.State.Name))
			}
			ko.Status.State = f42
		} else {
			ko.Status.State = nil
		}
		if elem.StateReason != nil {
			f43 := &svcapitypes.StateReason{}
			if elem.StateReason.Code != nil {
				f43.Code = elem.StateReason.Code
			}
			if elem.StateReason.Message != nil {
				f43.Message = elem.StateReason.Message
			}
			ko.Status.StateReason = f43
		} else {
			ko.Status.StateReason = nil
		}
		if elem.StateTransitionReason != nil {
			ko.Status.StateTransitionReason = elem.StateTransitionReason
		} else {
			ko.Status.StateTransitionReason = nil
		}
		if elem.SubnetId != nil {
			ko.Spec.SubnetID = elem.SubnetId
		} else {
			ko.Spec.SubnetID = nil
		}
		if elem.Tags != nil {
			f46 := []*svcapitypes.Tag{}
			for _, f46iter := range elem.Tags {
				f46elem := &svcapitypes.Tag{}
				if f46iter.Key != nil {
					f46elem.Key = f46iter.Key
				}
				if f46iter.Value != nil {
					f46elem.Value = f46iter.Value
				}
				f46 = append(f46, f46elem)
			}
			ko.Status.Tags = f46
		} else {
			ko.Status.Tags = nil
		}
		if elem.UsageOperation != nil {
			ko.Status.UsageOperation = elem.UsageOperation
		} else {
			ko.Status.UsageOperation = nil
		}
		if elem.UsageOperationUpdateTime != nil {
			ko.Status.UsageOperationUpdateTime = &metav1.Time{*elem.UsageOperationUpdateTime}
		} else {
			ko.Status.UsageOperationUpdateTime = nil
		}
		if elem.VirtualizationType != "" {
			ko.Status.VirtualizationType = aws.String(string(elem.VirtualizationType))
		} else {
			ko.Status.VirtualizationType = nil
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
		code.SetResource(crd.Config(), crd, op, "resp", "ko", 1),
	)
}

func TestSetResource_EC2_Instance_ReadMany(t *testing.T) {
	// DescribeInstances returns a list of Reservations
	// containing a list of Instances. The first Instance
	// in the first Reservation list should be used to populate CR.
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")
	op := model.OpTypeList

	crd := testutil.GetCRDByName(t, g, "Instance")
	require.NotNil(crd)

	expected := `
	found := false
	for _, iter0 := range resp.Reservations {
		for _, elem := range iter0.Instances {
			if elem.AmiLaunchIndex != nil {
				amiLaunchIndexCopy := int64(*elem.AmiLaunchIndex)
				ko.Status.AMILaunchIndex = &amiLaunchIndexCopy
			} else {
				ko.Status.AMILaunchIndex = nil
			}
			if elem.Architecture != "" {
				ko.Status.Architecture = aws.String(string(elem.Architecture))
			} else {
				ko.Status.Architecture = nil
			}
			if elem.BlockDeviceMappings != nil {
				f2 := []*svcapitypes.BlockDeviceMapping{}
				for _, f2iter := range elem.BlockDeviceMappings {
					f2elem := &svcapitypes.BlockDeviceMapping{}
					if f2iter.DeviceName != nil {
						f2elem.DeviceName = f2iter.DeviceName
					}
					if f2iter.Ebs != nil {
						f2elemf1 := &svcapitypes.EBSBlockDevice{}
						if f2iter.Ebs.DeleteOnTermination != nil {
							f2elemf1.DeleteOnTermination = f2iter.Ebs.DeleteOnTermination
						}
						f2elem.EBS = f2elemf1
					}
					f2 = append(f2, f2elem)
				}
				ko.Spec.BlockDeviceMappings = f2
			} else {
				ko.Spec.BlockDeviceMappings = nil
			}
			if elem.BootMode != "" {
				ko.Status.BootMode = aws.String(string(elem.BootMode))
			} else {
				ko.Status.BootMode = nil
			}
			if elem.CapacityReservationId != nil {
				ko.Status.CapacityReservationID = elem.CapacityReservationId
			} else {
				ko.Status.CapacityReservationID = nil
			}
			if elem.CapacityReservationSpecification != nil {
				f5 := &svcapitypes.CapacityReservationSpecification{}
				if elem.CapacityReservationSpecification.CapacityReservationPreference != "" {
					f5.CapacityReservationPreference = aws.String(string(elem.CapacityReservationSpecification.CapacityReservationPreference))
				}
				if elem.CapacityReservationSpecification.CapacityReservationTarget != nil {
					f5f1 := &svcapitypes.CapacityReservationTarget{}
					if elem.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationId != nil {
						f5f1.CapacityReservationID = elem.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationId
					}
					if elem.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationResourceGroupArn != nil {
						f5f1.CapacityReservationResourceGroupARN = elem.CapacityReservationSpecification.CapacityReservationTarget.CapacityReservationResourceGroupArn
					}
					f5.CapacityReservationTarget = f5f1
				}
				ko.Spec.CapacityReservationSpecification = f5
			} else {
				ko.Spec.CapacityReservationSpecification = nil
			}
			if elem.CpuOptions != nil {
				f6 := &svcapitypes.CPUOptionsRequest{}
				if elem.CpuOptions.CoreCount != nil {
					coreCountCopy := int64(*elem.CpuOptions.CoreCount)
					f6.CoreCount = &coreCountCopy
				}
				if elem.CpuOptions.ThreadsPerCore != nil {
					threadsPerCoreCopy := int64(*elem.CpuOptions.ThreadsPerCore)
					f6.ThreadsPerCore = &threadsPerCoreCopy
				}
				ko.Spec.CPUOptions = f6
			} else {
				ko.Spec.CPUOptions = nil
			}
			if elem.EbsOptimized != nil {
				ko.Spec.EBSOptimized = elem.EbsOptimized
			} else {
				ko.Spec.EBSOptimized = nil
			}
			if elem.ElasticGpuAssociations != nil {
				f8 := []*svcapitypes.ElasticGPUAssociation{}
				for _, f8iter := range elem.ElasticGpuAssociations {
					f8elem := &svcapitypes.ElasticGPUAssociation{}
					if f8iter.ElasticGpuAssociationId != nil {
						f8elem.ElasticGPUAssociationID = f8iter.ElasticGpuAssociationId
					}
					if f8iter.ElasticGpuAssociationState != nil {
						f8elem.ElasticGPUAssociationState = f8iter.ElasticGpuAssociationState
					}
					if f8iter.ElasticGpuAssociationTime != nil {
						f8elem.ElasticGPUAssociationTime = f8iter.ElasticGpuAssociationTime
					}
					if f8iter.ElasticGpuId != nil {
						f8elem.ElasticGPUID = f8iter.ElasticGpuId
					}
					f8 = append(f8, f8elem)
				}
				ko.Status.ElasticGPUAssociations = f8
			} else {
				ko.Status.ElasticGPUAssociations = nil
			}
			if elem.ElasticInferenceAcceleratorAssociations != nil {
				f9 := []*svcapitypes.ElasticInferenceAcceleratorAssociation{}
				for _, f9iter := range elem.ElasticInferenceAcceleratorAssociations {
					f9elem := &svcapitypes.ElasticInferenceAcceleratorAssociation{}
					if f9iter.ElasticInferenceAcceleratorArn != nil {
						f9elem.ElasticInferenceAcceleratorARN = f9iter.ElasticInferenceAcceleratorArn
					}
					if f9iter.ElasticInferenceAcceleratorAssociationId != nil {
						f9elem.ElasticInferenceAcceleratorAssociationID = f9iter.ElasticInferenceAcceleratorAssociationId
					}
					if f9iter.ElasticInferenceAcceleratorAssociationState != nil {
						f9elem.ElasticInferenceAcceleratorAssociationState = f9iter.ElasticInferenceAcceleratorAssociationState
					}
					if f9iter.ElasticInferenceAcceleratorAssociationTime != nil {
						f9elem.ElasticInferenceAcceleratorAssociationTime = &metav1.Time{*f9iter.ElasticInferenceAcceleratorAssociationTime}
					}
					f9 = append(f9, f9elem)
				}
				ko.Status.ElasticInferenceAcceleratorAssociations = f9
			} else {
				ko.Status.ElasticInferenceAcceleratorAssociations = nil
			}
			if elem.EnaSupport != nil {
				ko.Status.ENASupport = elem.EnaSupport
			} else {
				ko.Status.ENASupport = nil
			}
			if elem.EnclaveOptions != nil {
				f11 := &svcapitypes.EnclaveOptionsRequest{}
				if elem.EnclaveOptions.Enabled != nil {
					f11.Enabled = elem.EnclaveOptions.Enabled
				}
				ko.Spec.EnclaveOptions = f11
			} else {
				ko.Spec.EnclaveOptions = nil
			}
			if elem.HibernationOptions != nil {
				f12 := &svcapitypes.HibernationOptionsRequest{}
				if elem.HibernationOptions.Configured != nil {
					f12.Configured = elem.HibernationOptions.Configured
				}
				ko.Spec.HibernationOptions = f12
			} else {
				ko.Spec.HibernationOptions = nil
			}
			if elem.Hypervisor != "" {
				ko.Status.Hypervisor = aws.String(string(elem.Hypervisor))
			} else {
				ko.Status.Hypervisor = nil
			}
			if elem.IamInstanceProfile != nil {
				f14 := &svcapitypes.IAMInstanceProfileSpecification{}
				if elem.IamInstanceProfile.Arn != nil {
					f14.ARN = elem.IamInstanceProfile.Arn
				}
				ko.Spec.IAMInstanceProfile = f14
			} else {
				ko.Spec.IAMInstanceProfile = nil
			}
			if elem.ImageId != nil {
				ko.Spec.ImageID = elem.ImageId
			} else {
				ko.Spec.ImageID = nil
			}
			if elem.InstanceId != nil {
				ko.Status.InstanceID = elem.InstanceId
			} else {
				ko.Status.InstanceID = nil
			}
			if elem.InstanceLifecycle != "" {
				ko.Status.InstanceLifecycle = aws.String(string(elem.InstanceLifecycle))
			} else {
				ko.Status.InstanceLifecycle = nil
			}
			if elem.InstanceType != "" {
				ko.Spec.InstanceType = aws.String(string(elem.InstanceType))
			} else {
				ko.Spec.InstanceType = nil
			}
			if elem.KernelId != nil {
				ko.Spec.KernelID = elem.KernelId
			} else {
				ko.Spec.KernelID = nil
			}
			if elem.KeyName != nil {
				ko.Spec.KeyName = elem.KeyName
			} else {
				ko.Spec.KeyName = nil
			}
			if elem.LaunchTime != nil {
				ko.Status.LaunchTime = &metav1.Time{*elem.LaunchTime}
			} else {
				ko.Status.LaunchTime = nil
			}
			if elem.Licenses != nil {
				f22 := []*svcapitypes.LicenseConfiguration{}
				for _, f22iter := range elem.Licenses {
					f22elem := &svcapitypes.LicenseConfiguration{}
					if f22iter.LicenseConfigurationArn != nil {
						f22elem.LicenseConfigurationARN = f22iter.LicenseConfigurationArn
					}
					f22 = append(f22, f22elem)
				}
				ko.Status.Licenses = f22
			} else {
				ko.Status.Licenses = nil
			}
			if elem.MetadataOptions != nil {
				f23 := &svcapitypes.InstanceMetadataOptionsRequest{}
				if elem.MetadataOptions.HttpEndpoint != "" {
					f23.HTTPEndpoint = aws.String(string(elem.MetadataOptions.HttpEndpoint))
				}
				if elem.MetadataOptions.HttpProtocolIpv6 != "" {
					f23.HTTPProtocolIPv6 = aws.String(string(elem.MetadataOptions.HttpProtocolIpv6))
				}
				if elem.MetadataOptions.HttpPutResponseHopLimit != nil {
					httpPutResponseHopLimitCopy := int64(*elem.MetadataOptions.HttpPutResponseHopLimit)
					f23.HTTPPutResponseHopLimit = &httpPutResponseHopLimitCopy
				}
				if elem.MetadataOptions.HttpTokens != "" {
					f23.HTTPTokens = aws.String(string(elem.MetadataOptions.HttpTokens))
				}
				ko.Spec.MetadataOptions = f23
			} else {
				ko.Spec.MetadataOptions = nil
			}
			if elem.Monitoring != nil {
				f24 := &svcapitypes.RunInstancesMonitoringEnabled{}
				ko.Spec.Monitoring = f24
			} else {
				ko.Spec.Monitoring = nil
			}
			if elem.NetworkInterfaces != nil {
				f25 := []*svcapitypes.InstanceNetworkInterfaceSpecification{}
				for _, f25iter := range elem.NetworkInterfaces {
					f25elem := &svcapitypes.InstanceNetworkInterfaceSpecification{}
					if f25iter.Description != nil {
						f25elem.Description = f25iter.Description
					}
					if f25iter.InterfaceType != nil {
						f25elem.InterfaceType = f25iter.InterfaceType
					}
					if f25iter.Ipv4Prefixes != nil {
						f25elemf6 := []*svcapitypes.IPv4PrefixSpecificationRequest{}
						for _, f25elemf6iter := range f25iter.Ipv4Prefixes {
							f25elemf6elem := &svcapitypes.IPv4PrefixSpecificationRequest{}
							if f25elemf6iter.Ipv4Prefix != nil {
								f25elemf6elem.IPv4Prefix = f25elemf6iter.Ipv4Prefix
							}
							f25elemf6 = append(f25elemf6, f25elemf6elem)
						}
						f25elem.IPv4Prefixes = f25elemf6
					}
					if f25iter.Ipv6Addresses != nil {
						f25elemf7 := []*svcapitypes.InstanceIPv6Address{}
						for _, f25elemf7iter := range f25iter.Ipv6Addresses {
							f25elemf7elem := &svcapitypes.InstanceIPv6Address{}
							if f25elemf7iter.Ipv6Address != nil {
								f25elemf7elem.IPv6Address = f25elemf7iter.Ipv6Address
							}
							f25elemf7 = append(f25elemf7, f25elemf7elem)
						}
						f25elem.IPv6Addresses = f25elemf7
					}
					if f25iter.Ipv6Prefixes != nil {
						f25elemf8 := []*svcapitypes.IPv6PrefixSpecificationRequest{}
						for _, f25elemf8iter := range f25iter.Ipv6Prefixes {
							f25elemf8elem := &svcapitypes.IPv6PrefixSpecificationRequest{}
							if f25elemf8iter.Ipv6Prefix != nil {
								f25elemf8elem.IPv6Prefix = f25elemf8iter.Ipv6Prefix
							}
							f25elemf8 = append(f25elemf8, f25elemf8elem)
						}
						f25elem.IPv6Prefixes = f25elemf8
					}
					if f25iter.NetworkInterfaceId != nil {
						f25elem.NetworkInterfaceID = f25iter.NetworkInterfaceId
					}
					if f25iter.PrivateIpAddress != nil {
						f25elem.PrivateIPAddress = f25iter.PrivateIpAddress
					}
					if f25iter.PrivateIpAddresses != nil {
						f25elemf15 := []*svcapitypes.PrivateIPAddressSpecification{}
						for _, f25elemf15iter := range f25iter.PrivateIpAddresses {
							f25elemf15elem := &svcapitypes.PrivateIPAddressSpecification{}
							if f25elemf15iter.Primary != nil {
								f25elemf15elem.Primary = f25elemf15iter.Primary
							}
							if f25elemf15iter.PrivateIpAddress != nil {
								f25elemf15elem.PrivateIPAddress = f25elemf15iter.PrivateIpAddress
							}
							f25elemf15 = append(f25elemf15, f25elemf15elem)
						}
						f25elem.PrivateIPAddresses = f25elemf15
					}
					if f25iter.SubnetId != nil {
						f25elem.SubnetID = f25iter.SubnetId
					}
					f25 = append(f25, f25elem)
				}
				ko.Spec.NetworkInterfaces = f25
			} else {
				ko.Spec.NetworkInterfaces = nil
			}
			if elem.OutpostArn != nil {
				ko.Status.OutpostARN = elem.OutpostArn
			} else {
				ko.Status.OutpostARN = nil
			}
			if elem.Placement != nil {
				f27 := &svcapitypes.Placement{}
				if elem.Placement.Affinity != nil {
					f27.Affinity = elem.Placement.Affinity
				}
				if elem.Placement.AvailabilityZone != nil {
					f27.AvailabilityZone = elem.Placement.AvailabilityZone
				}
				if elem.Placement.GroupName != nil {
					f27.GroupName = elem.Placement.GroupName
				}
				if elem.Placement.HostId != nil {
					f27.HostID = elem.Placement.HostId
				}
				if elem.Placement.HostResourceGroupArn != nil {
					f27.HostResourceGroupARN = elem.Placement.HostResourceGroupArn
				}
				if elem.Placement.PartitionNumber != nil {
					partitionNumberCopy := int64(*elem.Placement.PartitionNumber)
					f27.PartitionNumber = &partitionNumberCopy
				}
				if elem.Placement.SpreadDomain != nil {
					f27.SpreadDomain = elem.Placement.SpreadDomain
				}
				if elem.Placement.Tenancy != "" {
					f27.Tenancy = aws.String(string(elem.Placement.Tenancy))
				}
				ko.Spec.Placement = f27
			} else {
				ko.Spec.Placement = nil
			}
			if elem.Platform != "" {
				ko.Status.Platform = aws.String(string(elem.Platform))
			} else {
				ko.Status.Platform = nil
			}
			if elem.PlatformDetails != nil {
				ko.Status.PlatformDetails = elem.PlatformDetails
			} else {
				ko.Status.PlatformDetails = nil
			}
			if elem.PrivateDnsName != nil {
				ko.Status.PrivateDNSName = elem.PrivateDnsName
			} else {
				ko.Status.PrivateDNSName = nil
			}
			if elem.PrivateIpAddress != nil {
				ko.Spec.PrivateIPAddress = elem.PrivateIpAddress
			} else {
				ko.Spec.PrivateIPAddress = nil
			}
			if elem.ProductCodes != nil {
				f32 := []*svcapitypes.ProductCode{}
				for _, f32iter := range elem.ProductCodes {
					f32elem := &svcapitypes.ProductCode{}
					if f32iter.ProductCodeId != nil {
						f32elem.ProductCodeID = f32iter.ProductCodeId
					}
					if f32iter.ProductCodeType != "" {
						f32elem.ProductCodeType = aws.String(string(f32iter.ProductCodeType))
					}
					f32 = append(f32, f32elem)
				}
				ko.Status.ProductCodes = f32
			} else {
				ko.Status.ProductCodes = nil
			}
			if elem.PublicDnsName != nil {
				ko.Status.PublicDNSName = elem.PublicDnsName
			} else {
				ko.Status.PublicDNSName = nil
			}
			if elem.PublicIpAddress != nil {
				ko.Status.PublicIPAddress = elem.PublicIpAddress
			} else {
				ko.Status.PublicIPAddress = nil
			}
			if elem.RamdiskId != nil {
				ko.Spec.RAMDiskID = elem.RamdiskId
			} else {
				ko.Spec.RAMDiskID = nil
			}
			if elem.RootDeviceName != nil {
				ko.Status.RootDeviceName = elem.RootDeviceName
			} else {
				ko.Status.RootDeviceName = nil
			}
			if elem.RootDeviceType != "" {
				ko.Status.RootDeviceType = aws.String(string(elem.RootDeviceType))
			} else {
				ko.Status.RootDeviceType = nil
			}
			if elem.SecurityGroups != nil {
				f38 := []*string{}
				for _, f38iter := range elem.SecurityGroups {
					var f38elem *string
					f38elem = f38iter.GroupName
					f38 = append(f38, f38elem)
				}
				ko.Spec.SecurityGroups = f38
			} else {
				ko.Spec.SecurityGroups = nil
			}
			if elem.SourceDestCheck != nil {
				ko.Status.SourceDestCheck = elem.SourceDestCheck
			} else {
				ko.Status.SourceDestCheck = nil
			}
			if elem.SpotInstanceRequestId != nil {
				ko.Status.SpotInstanceRequestID = elem.SpotInstanceRequestId
			} else {
				ko.Status.SpotInstanceRequestID = nil
			}
			if elem.SriovNetSupport != nil {
				ko.Status.SRIOVNetSupport = elem.SriovNetSupport
			} else {
				ko.Status.SRIOVNetSupport = nil
			}
			if elem.State != nil {
				f42 := &svcapitypes.InstanceState{}
				if elem.State.Code != nil {
					codeCopy := int64(*elem.State.Code)
					f42.Code = &codeCopy
				}
				if elem.State.Name != "" {
					f42.Name = aws.String(string(elem.State.Name))
				}
				ko.Status.State = f42
			} else {
				ko.Status.State = nil
			}
			if elem.StateReason != nil {
				f43 := &svcapitypes.StateReason{}
				if elem.StateReason.Code != nil {
					f43.Code = elem.StateReason.Code
				}
				if elem.StateReason.Message != nil {
					f43.Message = elem.StateReason.Message
				}
				ko.Status.StateReason = f43
			} else {
				ko.Status.StateReason = nil
			}
			if elem.StateTransitionReason != nil {
				ko.Status.StateTransitionReason = elem.StateTransitionReason
			} else {
				ko.Status.StateTransitionReason = nil
			}
			if elem.SubnetId != nil {
				ko.Spec.SubnetID = elem.SubnetId
			} else {
				ko.Spec.SubnetID = nil
			}
			if elem.Tags != nil {
				f46 := []*svcapitypes.Tag{}
				for _, f46iter := range elem.Tags {
					f46elem := &svcapitypes.Tag{}
					if f46iter.Key != nil {
						f46elem.Key = f46iter.Key
					}
					if f46iter.Value != nil {
						f46elem.Value = f46iter.Value
					}
					f46 = append(f46, f46elem)
				}
				ko.Status.Tags = f46
			} else {
				ko.Status.Tags = nil
			}
			if elem.UsageOperation != nil {
				ko.Status.UsageOperation = elem.UsageOperation
			} else {
				ko.Status.UsageOperation = nil
			}
			if elem.UsageOperationUpdateTime != nil {
				ko.Status.UsageOperationUpdateTime = &metav1.Time{*elem.UsageOperationUpdateTime}
			} else {
				ko.Status.UsageOperationUpdateTime = nil
			}
			if elem.VirtualizationType != "" {
				ko.Status.VirtualizationType = aws.String(string(elem.VirtualizationType))
			} else {
				ko.Status.VirtualizationType = nil
			}
			if elem.VpcId != nil {
				ko.Status.VPCID = elem.VpcId
			} else {
				ko.Status.VPCID = nil
			}
			found = true
			break
		}
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, op, "resp", "ko", 1),
	)
}

func TestSetResource_WAFv2_RuleGroup_ReadOne(t *testing.T) {
	// DescribeInstances returns a list of Reservations
	// containing a list of Instances. The first Instance
	// in the first Reservation list should be used to populate CR.
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "wafv2")
	op := model.OpTypeGet

	crd := testutil.GetCRDByName(t, g, "RuleGroup")
	require.NotNil(crd)

	expected := `
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.RuleGroup.ARN != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.RuleGroup.ARN)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.RuleGroup.Capacity != nil {
		ko.Spec.Capacity = resp.RuleGroup.Capacity
	} else {
		ko.Spec.Capacity = nil
	}
	if resp.RuleGroup.CustomResponseBodies != nil {
		f4 := map[string]*svcapitypes.CustomResponseBody{}
		for f4key, f4valiter := range resp.RuleGroup.CustomResponseBodies {
			f4val := &svcapitypes.CustomResponseBody{}
			if f4valiter.Content != nil {
				f4val.Content = f4valiter.Content
			}
			if f4valiter.ContentType != "" {
				f4val.ContentType = aws.String(string(f4valiter.ContentType))
			}
			f4[f4key] = f4val
		}
		ko.Spec.CustomResponseBodies = f4
	} else {
		ko.Spec.CustomResponseBodies = nil
	}
	if resp.RuleGroup.Description != nil {
		ko.Spec.Description = resp.RuleGroup.Description
	} else {
		ko.Spec.Description = nil
	}
	if resp.RuleGroup.Id != nil {
		ko.Status.ID = resp.RuleGroup.Id
	} else {
		ko.Status.ID = nil
	}
	if resp.RuleGroup.Name != nil {
		ko.Spec.Name = resp.RuleGroup.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.RuleGroup.Rules != nil {
		f9 := []*svcapitypes.Rule{}
		for _, f9iter := range resp.RuleGroup.Rules {
			f9elem := &svcapitypes.Rule{}
			if f9iter.Action != nil {
				f9elemf0 := &svcapitypes.RuleAction{}
				if f9iter.Action.Allow != nil {
					f9elemf0f0 := &svcapitypes.AllowAction{}
					if f9iter.Action.Allow.CustomRequestHandling != nil {
						f9elemf0f0f0 := &svcapitypes.CustomRequestHandling{}
						if f9iter.Action.Allow.CustomRequestHandling.InsertHeaders != nil {
							f9elemf0f0f0f0 := []*svcapitypes.CustomHTTPHeader{}
							for _, f9elemf0f0f0f0iter := range f9iter.Action.Allow.CustomRequestHandling.InsertHeaders {
								f9elemf0f0f0f0elem := &svcapitypes.CustomHTTPHeader{}
								if f9elemf0f0f0f0iter.Name != nil {
									f9elemf0f0f0f0elem.Name = f9elemf0f0f0f0iter.Name
								}
								if f9elemf0f0f0f0iter.Value != nil {
									f9elemf0f0f0f0elem.Value = f9elemf0f0f0f0iter.Value
								}
								f9elemf0f0f0f0 = append(f9elemf0f0f0f0, f9elemf0f0f0f0elem)
							}
							f9elemf0f0f0.InsertHeaders = f9elemf0f0f0f0
						}
						f9elemf0f0.CustomRequestHandling = f9elemf0f0f0
					}
					f9elemf0.Allow = f9elemf0f0
				}
				if f9iter.Action.Block != nil {
					f9elemf0f1 := &svcapitypes.BlockAction{}
					if f9iter.Action.Block.CustomResponse != nil {
						f9elemf0f1f0 := &svcapitypes.CustomResponse{}
						if f9iter.Action.Block.CustomResponse.CustomResponseBodyKey != nil {
							f9elemf0f1f0.CustomResponseBodyKey = f9iter.Action.Block.CustomResponse.CustomResponseBodyKey
						}
						if f9iter.Action.Block.CustomResponse.ResponseCode != nil {
							responseCodeCopy := int64(*f9iter.Action.Block.CustomResponse.ResponseCode)
							f9elemf0f1f0.ResponseCode = &responseCodeCopy
						}
						if f9iter.Action.Block.CustomResponse.ResponseHeaders != nil {
							f9elemf0f1f0f2 := []*svcapitypes.CustomHTTPHeader{}
							for _, f9elemf0f1f0f2iter := range f9iter.Action.Block.CustomResponse.ResponseHeaders {
								f9elemf0f1f0f2elem := &svcapitypes.CustomHTTPHeader{}
								if f9elemf0f1f0f2iter.Name != nil {
									f9elemf0f1f0f2elem.Name = f9elemf0f1f0f2iter.Name
								}
								if f9elemf0f1f0f2iter.Value != nil {
									f9elemf0f1f0f2elem.Value = f9elemf0f1f0f2iter.Value
								}
								f9elemf0f1f0f2 = append(f9elemf0f1f0f2, f9elemf0f1f0f2elem)
							}
							f9elemf0f1f0.ResponseHeaders = f9elemf0f1f0f2
						}
						f9elemf0f1.CustomResponse = f9elemf0f1f0
					}
					f9elemf0.Block = f9elemf0f1
				}
				if f9iter.Action.Captcha != nil {
					f9elemf0f2 := &svcapitypes.CaptchaAction{}
					if f9iter.Action.Captcha.CustomRequestHandling != nil {
						f9elemf0f2f0 := &svcapitypes.CustomRequestHandling{}
						if f9iter.Action.Captcha.CustomRequestHandling.InsertHeaders != nil {
							f9elemf0f2f0f0 := []*svcapitypes.CustomHTTPHeader{}
							for _, f9elemf0f2f0f0iter := range f9iter.Action.Captcha.CustomRequestHandling.InsertHeaders {
								f9elemf0f2f0f0elem := &svcapitypes.CustomHTTPHeader{}
								if f9elemf0f2f0f0iter.Name != nil {
									f9elemf0f2f0f0elem.Name = f9elemf0f2f0f0iter.Name
								}
								if f9elemf0f2f0f0iter.Value != nil {
									f9elemf0f2f0f0elem.Value = f9elemf0f2f0f0iter.Value
								}
								f9elemf0f2f0f0 = append(f9elemf0f2f0f0, f9elemf0f2f0f0elem)
							}
							f9elemf0f2f0.InsertHeaders = f9elemf0f2f0f0
						}
						f9elemf0f2.CustomRequestHandling = f9elemf0f2f0
					}
					f9elemf0.Captcha = f9elemf0f2
				}
				if f9iter.Action.Challenge != nil {
					f9elemf0f3 := &svcapitypes.ChallengeAction{}
					if f9iter.Action.Challenge.CustomRequestHandling != nil {
						f9elemf0f3f0 := &svcapitypes.CustomRequestHandling{}
						if f9iter.Action.Challenge.CustomRequestHandling.InsertHeaders != nil {
							f9elemf0f3f0f0 := []*svcapitypes.CustomHTTPHeader{}
							for _, f9elemf0f3f0f0iter := range f9iter.Action.Challenge.CustomRequestHandling.InsertHeaders {
								f9elemf0f3f0f0elem := &svcapitypes.CustomHTTPHeader{}
								if f9elemf0f3f0f0iter.Name != nil {
									f9elemf0f3f0f0elem.Name = f9elemf0f3f0f0iter.Name
								}
								if f9elemf0f3f0f0iter.Value != nil {
									f9elemf0f3f0f0elem.Value = f9elemf0f3f0f0iter.Value
								}
								f9elemf0f3f0f0 = append(f9elemf0f3f0f0, f9elemf0f3f0f0elem)
							}
							f9elemf0f3f0.InsertHeaders = f9elemf0f3f0f0
						}
						f9elemf0f3.CustomRequestHandling = f9elemf0f3f0
					}
					f9elemf0.Challenge = f9elemf0f3
				}
				if f9iter.Action.Count != nil {
					f9elemf0f4 := &svcapitypes.CountAction{}
					if f9iter.Action.Count.CustomRequestHandling != nil {
						f9elemf0f4f0 := &svcapitypes.CustomRequestHandling{}
						if f9iter.Action.Count.CustomRequestHandling.InsertHeaders != nil {
							f9elemf0f4f0f0 := []*svcapitypes.CustomHTTPHeader{}
							for _, f9elemf0f4f0f0iter := range f9iter.Action.Count.CustomRequestHandling.InsertHeaders {
								f9elemf0f4f0f0elem := &svcapitypes.CustomHTTPHeader{}
								if f9elemf0f4f0f0iter.Name != nil {
									f9elemf0f4f0f0elem.Name = f9elemf0f4f0f0iter.Name
								}
								if f9elemf0f4f0f0iter.Value != nil {
									f9elemf0f4f0f0elem.Value = f9elemf0f4f0f0iter.Value
								}
								f9elemf0f4f0f0 = append(f9elemf0f4f0f0, f9elemf0f4f0f0elem)
							}
							f9elemf0f4f0.InsertHeaders = f9elemf0f4f0f0
						}
						f9elemf0f4.CustomRequestHandling = f9elemf0f4f0
					}
					f9elemf0.Count = f9elemf0f4
				}
				f9elem.Action = f9elemf0
			}
			if f9iter.CaptchaConfig != nil {
				f9elemf1 := &svcapitypes.CaptchaConfig{}
				if f9iter.CaptchaConfig.ImmunityTimeProperty != nil {
					f9elemf1f0 := &svcapitypes.ImmunityTimeProperty{}
					if f9iter.CaptchaConfig.ImmunityTimeProperty.ImmunityTime != nil {
						f9elemf1f0.ImmunityTime = f9iter.CaptchaConfig.ImmunityTimeProperty.ImmunityTime
					}
					f9elemf1.ImmunityTimeProperty = f9elemf1f0
				}
				f9elem.CaptchaConfig = f9elemf1
			}
			if f9iter.ChallengeConfig != nil {
				f9elemf2 := &svcapitypes.ChallengeConfig{}
				if f9iter.ChallengeConfig.ImmunityTimeProperty != nil {
					f9elemf2f0 := &svcapitypes.ImmunityTimeProperty{}
					if f9iter.ChallengeConfig.ImmunityTimeProperty.ImmunityTime != nil {
						f9elemf2f0.ImmunityTime = f9iter.ChallengeConfig.ImmunityTimeProperty.ImmunityTime
					}
					f9elemf2.ImmunityTimeProperty = f9elemf2f0
				}
				f9elem.ChallengeConfig = f9elemf2
			}
			if f9iter.Name != nil {
				f9elem.Name = f9iter.Name
			}
			if f9iter.OverrideAction != nil {
				f9elemf4 := &svcapitypes.OverrideAction{}
				if f9iter.OverrideAction.Count != nil {
					f9elemf4f0 := &svcapitypes.CountAction{}
					if f9iter.OverrideAction.Count.CustomRequestHandling != nil {
						f9elemf4f0f0 := &svcapitypes.CustomRequestHandling{}
						if f9iter.OverrideAction.Count.CustomRequestHandling.InsertHeaders != nil {
							f9elemf4f0f0f0 := []*svcapitypes.CustomHTTPHeader{}
							for _, f9elemf4f0f0f0iter := range f9iter.OverrideAction.Count.CustomRequestHandling.InsertHeaders {
								f9elemf4f0f0f0elem := &svcapitypes.CustomHTTPHeader{}
								if f9elemf4f0f0f0iter.Name != nil {
									f9elemf4f0f0f0elem.Name = f9elemf4f0f0f0iter.Name
								}
								if f9elemf4f0f0f0iter.Value != nil {
									f9elemf4f0f0f0elem.Value = f9elemf4f0f0f0iter.Value
								}
								f9elemf4f0f0f0 = append(f9elemf4f0f0f0, f9elemf4f0f0f0elem)
							}
							f9elemf4f0f0.InsertHeaders = f9elemf4f0f0f0
						}
						f9elemf4f0.CustomRequestHandling = f9elemf4f0f0
					}
					f9elemf4.Count = f9elemf4f0
				}
				if f9iter.OverrideAction.None != nil {
					f9elemf4f1 := map[string]*string{}
					f9elemf4.None = f9elemf4f1
				}
				f9elem.OverrideAction = f9elemf4
			}
			priorityCopy := int64(f9iter.Priority)
			f9elem.Priority = &priorityCopy
			if f9iter.RuleLabels != nil {
				f9elemf6 := []*svcapitypes.Label{}
				for _, f9elemf6iter := range f9iter.RuleLabels {
					f9elemf6elem := &svcapitypes.Label{}
					if f9elemf6iter.Name != nil {
						f9elemf6elem.Name = f9elemf6iter.Name
					}
					f9elemf6 = append(f9elemf6, f9elemf6elem)
				}
				f9elem.RuleLabels = f9elemf6
			}
			if f9iter.Statement != nil {
				f9elemf7 := &svcapitypes.Statement{}
				if f9iter.Statement.ByteMatchStatement != nil {
					f9elemf7f1 := &svcapitypes.ByteMatchStatement{}
					if f9iter.Statement.ByteMatchStatement.FieldToMatch != nil {
						f9elemf7f1f0 := &svcapitypes.FieldToMatch{}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.AllQueryArguments != nil {
							f9elemf7f1f0f0 := map[string]*string{}
							f9elemf7f1f0.AllQueryArguments = f9elemf7f1f0f0
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.Body != nil {
							f9elemf7f1f0f1 := &svcapitypes.Body{}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.Body.OversizeHandling != "" {
								f9elemf7f1f0f1.OversizeHandling = aws.String(string(f9iter.Statement.ByteMatchStatement.FieldToMatch.Body.OversizeHandling))
							}
							f9elemf7f1f0.Body = f9elemf7f1f0f1
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies != nil {
							f9elemf7f1f0f2 := &svcapitypes.Cookies{}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f9elemf7f1f0f2f0 := &svcapitypes.CookieMatchPattern{}
								if f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f9elemf7f1f0f2f0f0 := map[string]*string{}
									f9elemf7f1f0f2f0.All = f9elemf7f1f0f2f0f0
								}
								if f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f9elemf7f1f0f2f0.ExcludedCookies = aws.StringSlice(f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f9elemf7f1f0f2f0.IncludedCookies = aws.StringSlice(f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f9elemf7f1f0f2.MatchPattern = f9elemf7f1f0f2f0
							}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchScope != "" {
								f9elemf7f1f0f2.MatchScope = aws.String(string(f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.MatchScope))
							}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.OversizeHandling != "" {
								f9elemf7f1f0f2.OversizeHandling = aws.String(string(f9iter.Statement.ByteMatchStatement.FieldToMatch.Cookies.OversizeHandling))
							}
							f9elemf7f1f0.Cookies = f9elemf7f1f0f2
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.HeaderOrder != nil {
							f9elemf7f1f0f3 := &svcapitypes.HeaderOrder{}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling != "" {
								f9elemf7f1f0f3.OversizeHandling = aws.String(string(f9iter.Statement.ByteMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling))
							}
							f9elemf7f1f0.HeaderOrder = f9elemf7f1f0f3
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers != nil {
							f9elemf7f1f0f4 := &svcapitypes.Headers{}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern != nil {
								f9elemf7f1f0f4f0 := &svcapitypes.HeaderMatchPattern{}
								if f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f9elemf7f1f0f4f0f0 := map[string]*string{}
									f9elemf7f1f0f4f0.All = f9elemf7f1f0f4f0f0
								}
								if f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f9elemf7f1f0f4f0.ExcludedHeaders = aws.StringSlice(f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f9elemf7f1f0f4f0.IncludedHeaders = aws.StringSlice(f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f9elemf7f1f0f4.MatchPattern = f9elemf7f1f0f4f0
							}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchScope != "" {
								f9elemf7f1f0f4.MatchScope = aws.String(string(f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers.MatchScope))
							}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers.OversizeHandling != "" {
								f9elemf7f1f0f4.OversizeHandling = aws.String(string(f9iter.Statement.ByteMatchStatement.FieldToMatch.Headers.OversizeHandling))
							}
							f9elemf7f1f0.Headers = f9elemf7f1f0f4
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.JA3Fingerprint != nil {
							f9elemf7f1f0f5 := &svcapitypes.JA3Fingerprint{}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != "" {
								f9elemf7f1f0f5.FallbackBehavior = aws.String(string(f9iter.Statement.ByteMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior))
							}
							f9elemf7f1f0.JA3Fingerprint = f9elemf7f1f0f5
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody != nil {
							f9elemf7f1f0f6 := &svcapitypes.JSONBody{}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior != "" {
								f9elemf7f1f0f6.InvalidFallbackBehavior = aws.String(string(f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior))
							}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody.MatchPattern != nil {
								f9elemf7f1f0f6f1 := &svcapitypes.JSONMatchPattern{}
								if f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody.MatchPattern.All != nil {
									f9elemf7f1f0f6f1f0 := map[string]*string{}
									f9elemf7f1f0f6f1.All = f9elemf7f1f0f6f1f0
								}
								if f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths != nil {
									f9elemf7f1f0f6f1.IncludedPaths = aws.StringSlice(f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths)
								}
								f9elemf7f1f0f6.MatchPattern = f9elemf7f1f0f6f1
							}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody.MatchScope != "" {
								f9elemf7f1f0f6.MatchScope = aws.String(string(f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody.MatchScope))
							}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody.OversizeHandling != "" {
								f9elemf7f1f0f6.OversizeHandling = aws.String(string(f9iter.Statement.ByteMatchStatement.FieldToMatch.JsonBody.OversizeHandling))
							}
							f9elemf7f1f0.JSONBody = f9elemf7f1f0f6
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.Method != nil {
							f9elemf7f1f0f7 := map[string]*string{}
							f9elemf7f1f0.Method = f9elemf7f1f0f7
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.QueryString != nil {
							f9elemf7f1f0f8 := map[string]*string{}
							f9elemf7f1f0.QueryString = f9elemf7f1f0f8
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.SingleHeader != nil {
							f9elemf7f1f0f9 := &svcapitypes.SingleHeader{}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.SingleHeader.Name != nil {
								f9elemf7f1f0f9.Name = f9iter.Statement.ByteMatchStatement.FieldToMatch.SingleHeader.Name
							}
							f9elemf7f1f0.SingleHeader = f9elemf7f1f0f9
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.SingleQueryArgument != nil {
							f9elemf7f1f0f10 := &svcapitypes.SingleQueryArgument{}
							if f9iter.Statement.ByteMatchStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f9elemf7f1f0f10.Name = f9iter.Statement.ByteMatchStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f9elemf7f1f0.SingleQueryArgument = f9elemf7f1f0f10
						}
						if f9iter.Statement.ByteMatchStatement.FieldToMatch.UriPath != nil {
							f9elemf7f1f0f11 := map[string]*string{}
							f9elemf7f1f0.URIPath = f9elemf7f1f0f11
						}
						f9elemf7f1.FieldToMatch = f9elemf7f1f0
					}
					if f9iter.Statement.ByteMatchStatement.PositionalConstraint != "" {
						f9elemf7f1.PositionalConstraint = aws.String(string(f9iter.Statement.ByteMatchStatement.PositionalConstraint))
					}
					if f9iter.Statement.ByteMatchStatement.SearchString != nil {
						f9elemf7f1.SearchString = f9iter.Statement.ByteMatchStatement.SearchString
					}
					if f9iter.Statement.ByteMatchStatement.TextTransformations != nil {
						f9elemf7f1f3 := []*svcapitypes.TextTransformation{}
						for _, f9elemf7f1f3iter := range f9iter.Statement.ByteMatchStatement.TextTransformations {
							f9elemf7f1f3elem := &svcapitypes.TextTransformation{}
							priorityCopy := int64(f9elemf7f1f3iter.Priority)
							f9elemf7f1f3elem.Priority = &priorityCopy
							if f9elemf7f1f3iter.Type != "" {
								f9elemf7f1f3elem.Type = aws.String(string(f9elemf7f1f3iter.Type))
							}
							f9elemf7f1f3 = append(f9elemf7f1f3, f9elemf7f1f3elem)
						}
						f9elemf7f1.TextTransformations = f9elemf7f1f3
					}
					f9elemf7.ByteMatchStatement = f9elemf7f1
				}
				if f9iter.Statement.GeoMatchStatement != nil {
					f9elemf7f2 := &svcapitypes.GeoMatchStatement{}
					if f9iter.Statement.GeoMatchStatement.CountryCodes != nil {
						f9elemf7f2f0 := []*string{}
						for _, f9elemf7f2f0iter := range f9iter.Statement.GeoMatchStatement.CountryCodes {
							var f9elemf7f2f0elem *string
							f9elemf7f2f0elem = aws.String(string(f9elemf7f2f0iter))
							f9elemf7f2f0 = append(f9elemf7f2f0, f9elemf7f2f0elem)
						}
						f9elemf7f2.CountryCodes = f9elemf7f2f0
					}
					if f9iter.Statement.GeoMatchStatement.ForwardedIPConfig != nil {
						f9elemf7f2f1 := &svcapitypes.ForwardedIPConfig{}
						if f9iter.Statement.GeoMatchStatement.ForwardedIPConfig.FallbackBehavior != "" {
							f9elemf7f2f1.FallbackBehavior = aws.String(string(f9iter.Statement.GeoMatchStatement.ForwardedIPConfig.FallbackBehavior))
						}
						if f9iter.Statement.GeoMatchStatement.ForwardedIPConfig.HeaderName != nil {
							f9elemf7f2f1.HeaderName = f9iter.Statement.GeoMatchStatement.ForwardedIPConfig.HeaderName
						}
						f9elemf7f2.ForwardedIPConfig = f9elemf7f2f1
					}
					f9elemf7.GeoMatchStatement = f9elemf7f2
				}
				if f9iter.Statement.IPSetReferenceStatement != nil {
					f9elemf7f3 := &svcapitypes.IPSetReferenceStatement{}
					if f9iter.Statement.IPSetReferenceStatement.ARN != nil {
						f9elemf7f3.ARN = f9iter.Statement.IPSetReferenceStatement.ARN
					}
					if f9iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig != nil {
						f9elemf7f3f1 := &svcapitypes.IPSetForwardedIPConfig{}
						if f9iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.FallbackBehavior != "" {
							f9elemf7f3f1.FallbackBehavior = aws.String(string(f9iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.FallbackBehavior))
						}
						if f9iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.HeaderName != nil {
							f9elemf7f3f1.HeaderName = f9iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.HeaderName
						}
						if f9iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.Position != "" {
							f9elemf7f3f1.Position = aws.String(string(f9iter.Statement.IPSetReferenceStatement.IPSetForwardedIPConfig.Position))
						}
						f9elemf7f3.IPSetForwardedIPConfig = f9elemf7f3f1
					}
					f9elemf7.IPSetReferenceStatement = f9elemf7f3
				}
				if f9iter.Statement.LabelMatchStatement != nil {
					f9elemf7f4 := &svcapitypes.LabelMatchStatement{}
					if f9iter.Statement.LabelMatchStatement.Key != nil {
						f9elemf7f4.Key = f9iter.Statement.LabelMatchStatement.Key
					}
					if f9iter.Statement.LabelMatchStatement.Scope != "" {
						f9elemf7f4.Scope = aws.String(string(f9iter.Statement.LabelMatchStatement.Scope))
					}
					f9elemf7.LabelMatchStatement = f9elemf7f4
				}
				if f9iter.Statement.ManagedRuleGroupStatement != nil {
					f9elemf7f5 := &svcapitypes.ManagedRuleGroupStatement{}
					if f9iter.Statement.ManagedRuleGroupStatement.ExcludedRules != nil {
						f9elemf7f5f0 := []*svcapitypes.ExcludedRule{}
						for _, f9elemf7f5f0iter := range f9iter.Statement.ManagedRuleGroupStatement.ExcludedRules {
							f9elemf7f5f0elem := &svcapitypes.ExcludedRule{}
							if f9elemf7f5f0iter.Name != nil {
								f9elemf7f5f0elem.Name = f9elemf7f5f0iter.Name
							}
							f9elemf7f5f0 = append(f9elemf7f5f0, f9elemf7f5f0elem)
						}
						f9elemf7f5.ExcludedRules = f9elemf7f5f0
					}
					if f9iter.Statement.ManagedRuleGroupStatement.ManagedRuleGroupConfigs != nil {
						f9elemf7f5f1 := []*svcapitypes.ManagedRuleGroupConfig{}
						for _, f9elemf7f5f1iter := range f9iter.Statement.ManagedRuleGroupStatement.ManagedRuleGroupConfigs {
							f9elemf7f5f1elem := &svcapitypes.ManagedRuleGroupConfig{}
							if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet != nil {
								f9elemf7f5f1elemf0 := &svcapitypes.AWSManagedRulesACFPRuleSet{}
								if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.CreationPath != nil {
									f9elemf7f5f1elemf0.CreationPath = f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.CreationPath
								}
								f9elemf7f5f1elemf0.EnableRegexInPath = &f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.EnableRegexInPath
								if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RegistrationPagePath != nil {
									f9elemf7f5f1elemf0.RegistrationPagePath = f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RegistrationPagePath
								}
								if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection != nil {
									f9elemf7f5f1elemf0f3 := &svcapitypes.RequestInspectionACFP{}
									if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.AddressFields != nil {
										f9elemf7f5f1elemf0f3f0 := []*svcapitypes.AddressField{}
										for _, f9elemf7f5f1elemf0f3f0iter := range f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.AddressFields {
											f9elemf7f5f1elemf0f3f0elem := &svcapitypes.AddressField{}
											if f9elemf7f5f1elemf0f3f0iter.Identifier != nil {
												f9elemf7f5f1elemf0f3f0elem.Identifier = f9elemf7f5f1elemf0f3f0iter.Identifier
											}
											f9elemf7f5f1elemf0f3f0 = append(f9elemf7f5f1elemf0f3f0, f9elemf7f5f1elemf0f3f0elem)
										}
										f9elemf7f5f1elemf0f3.AddressFields = f9elemf7f5f1elemf0f3f0
									}
									if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.EmailField != nil {
										f9elemf7f5f1elemf0f3f1 := &svcapitypes.EmailField{}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.EmailField.Identifier != nil {
											f9elemf7f5f1elemf0f3f1.Identifier = f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.EmailField.Identifier
										}
										f9elemf7f5f1elemf0f3.EmailField = f9elemf7f5f1elemf0f3f1
									}
									if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PasswordField != nil {
										f9elemf7f5f1elemf0f3f2 := &svcapitypes.PasswordField{}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PasswordField.Identifier != nil {
											f9elemf7f5f1elemf0f3f2.Identifier = f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PasswordField.Identifier
										}
										f9elemf7f5f1elemf0f3.PasswordField = f9elemf7f5f1elemf0f3f2
									}
									if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PayloadType != "" {
										f9elemf7f5f1elemf0f3.PayloadType = aws.String(string(f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PayloadType))
									}
									if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PhoneNumberFields != nil {
										f9elemf7f5f1elemf0f3f4 := []*svcapitypes.PhoneNumberField{}
										for _, f9elemf7f5f1elemf0f3f4iter := range f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.PhoneNumberFields {
											f9elemf7f5f1elemf0f3f4elem := &svcapitypes.PhoneNumberField{}
											if f9elemf7f5f1elemf0f3f4iter.Identifier != nil {
												f9elemf7f5f1elemf0f3f4elem.Identifier = f9elemf7f5f1elemf0f3f4iter.Identifier
											}
											f9elemf7f5f1elemf0f3f4 = append(f9elemf7f5f1elemf0f3f4, f9elemf7f5f1elemf0f3f4elem)
										}
										f9elemf7f5f1elemf0f3.PhoneNumberFields = f9elemf7f5f1elemf0f3f4
									}
									if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.UsernameField != nil {
										f9elemf7f5f1elemf0f3f5 := &svcapitypes.UsernameField{}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.UsernameField.Identifier != nil {
											f9elemf7f5f1elemf0f3f5.Identifier = f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.RequestInspection.UsernameField.Identifier
										}
										f9elemf7f5f1elemf0f3.UsernameField = f9elemf7f5f1elemf0f3f5
									}
									f9elemf7f5f1elemf0.RequestInspection = f9elemf7f5f1elemf0f3
								}
								if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection != nil {
									f9elemf7f5f1elemf0f4 := &svcapitypes.ResponseInspection{}
									if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.BodyContains != nil {
										f9elemf7f5f1elemf0f4f0 := &svcapitypes.ResponseInspectionBodyContains{}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.BodyContains.FailureStrings != nil {
											f9elemf7f5f1elemf0f4f0.FailureStrings = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.BodyContains.FailureStrings)
										}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.BodyContains.SuccessStrings != nil {
											f9elemf7f5f1elemf0f4f0.SuccessStrings = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.BodyContains.SuccessStrings)
										}
										f9elemf7f5f1elemf0f4.BodyContains = f9elemf7f5f1elemf0f4f0
									}
									if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header != nil {
										f9elemf7f5f1elemf0f4f1 := &svcapitypes.ResponseInspectionHeader{}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.FailureValues != nil {
											f9elemf7f5f1elemf0f4f1.FailureValues = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.FailureValues)
										}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.Name != nil {
											f9elemf7f5f1elemf0f4f1.Name = f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.Name
										}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.SuccessValues != nil {
											f9elemf7f5f1elemf0f4f1.SuccessValues = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Header.SuccessValues)
										}
										f9elemf7f5f1elemf0f4.Header = f9elemf7f5f1elemf0f4f1
									}
									if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Json != nil {
										f9elemf7f5f1elemf0f4f2 := &svcapitypes.ResponseInspectionJSON{}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Json.FailureValues != nil {
											f9elemf7f5f1elemf0f4f2.FailureValues = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Json.FailureValues)
										}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Json.Identifier != nil {
											f9elemf7f5f1elemf0f4f2.Identifier = f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Json.Identifier
										}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Json.SuccessValues != nil {
											f9elemf7f5f1elemf0f4f2.SuccessValues = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.Json.SuccessValues)
										}
										f9elemf7f5f1elemf0f4.JSON = f9elemf7f5f1elemf0f4f2
									}
									if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.StatusCode != nil {
										f9elemf7f5f1elemf0f4f3 := &svcapitypes.ResponseInspectionStatusCode{}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.StatusCode.FailureCodes != nil {
											f9elemf7f5f1elemf0f4f3f0 := []*int64{}
											for _, f9elemf7f5f1elemf0f4f3f0iter := range f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.StatusCode.FailureCodes {
												var f9elemf7f5f1elemf0f4f3f0elem *int64
												failureCodeCopy := int64(f9elemf7f5f1elemf0f4f3f0iter)
												f9elemf7f5f1elemf0f4f3f0elem = &failureCodeCopy
												f9elemf7f5f1elemf0f4f3f0 = append(f9elemf7f5f1elemf0f4f3f0, f9elemf7f5f1elemf0f4f3f0elem)
											}
											f9elemf7f5f1elemf0f4f3.FailureCodes = f9elemf7f5f1elemf0f4f3f0
										}
										if f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.StatusCode.SuccessCodes != nil {
											f9elemf7f5f1elemf0f4f3f1 := []*int64{}
											for _, f9elemf7f5f1elemf0f4f3f1iter := range f9elemf7f5f1iter.AWSManagedRulesACFPRuleSet.ResponseInspection.StatusCode.SuccessCodes {
												var f9elemf7f5f1elemf0f4f3f1elem *int64
												successCodeCopy := int64(f9elemf7f5f1elemf0f4f3f1iter)
												f9elemf7f5f1elemf0f4f3f1elem = &successCodeCopy
												f9elemf7f5f1elemf0f4f3f1 = append(f9elemf7f5f1elemf0f4f3f1, f9elemf7f5f1elemf0f4f3f1elem)
											}
											f9elemf7f5f1elemf0f4f3.SuccessCodes = f9elemf7f5f1elemf0f4f3f1
										}
										f9elemf7f5f1elemf0f4.StatusCode = f9elemf7f5f1elemf0f4f3
									}
									f9elemf7f5f1elemf0.ResponseInspection = f9elemf7f5f1elemf0f4
								}
								f9elemf7f5f1elem.AWSManagedRulesACFPRuleSet = f9elemf7f5f1elemf0
							}
							if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet != nil {
								f9elemf7f5f1elemf1 := &svcapitypes.AWSManagedRulesATPRuleSet{}
								f9elemf7f5f1elemf1.EnableRegexInPath = &f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.EnableRegexInPath
								if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.LoginPath != nil {
									f9elemf7f5f1elemf1.LoginPath = f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.LoginPath
								}
								if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection != nil {
									f9elemf7f5f1elemf1f2 := &svcapitypes.RequestInspection{}
									if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.PasswordField != nil {
										f9elemf7f5f1elemf1f2f0 := &svcapitypes.PasswordField{}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.PasswordField.Identifier != nil {
											f9elemf7f5f1elemf1f2f0.Identifier = f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.PasswordField.Identifier
										}
										f9elemf7f5f1elemf1f2.PasswordField = f9elemf7f5f1elemf1f2f0
									}
									if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.PayloadType != "" {
										f9elemf7f5f1elemf1f2.PayloadType = aws.String(string(f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.PayloadType))
									}
									if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.UsernameField != nil {
										f9elemf7f5f1elemf1f2f2 := &svcapitypes.UsernameField{}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.UsernameField.Identifier != nil {
											f9elemf7f5f1elemf1f2f2.Identifier = f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.RequestInspection.UsernameField.Identifier
										}
										f9elemf7f5f1elemf1f2.UsernameField = f9elemf7f5f1elemf1f2f2
									}
									f9elemf7f5f1elemf1.RequestInspection = f9elemf7f5f1elemf1f2
								}
								if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection != nil {
									f9elemf7f5f1elemf1f3 := &svcapitypes.ResponseInspection{}
									if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.BodyContains != nil {
										f9elemf7f5f1elemf1f3f0 := &svcapitypes.ResponseInspectionBodyContains{}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.BodyContains.FailureStrings != nil {
											f9elemf7f5f1elemf1f3f0.FailureStrings = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.BodyContains.FailureStrings)
										}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.BodyContains.SuccessStrings != nil {
											f9elemf7f5f1elemf1f3f0.SuccessStrings = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.BodyContains.SuccessStrings)
										}
										f9elemf7f5f1elemf1f3.BodyContains = f9elemf7f5f1elemf1f3f0
									}
									if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header != nil {
										f9elemf7f5f1elemf1f3f1 := &svcapitypes.ResponseInspectionHeader{}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.FailureValues != nil {
											f9elemf7f5f1elemf1f3f1.FailureValues = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.FailureValues)
										}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.Name != nil {
											f9elemf7f5f1elemf1f3f1.Name = f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.Name
										}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.SuccessValues != nil {
											f9elemf7f5f1elemf1f3f1.SuccessValues = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Header.SuccessValues)
										}
										f9elemf7f5f1elemf1f3.Header = f9elemf7f5f1elemf1f3f1
									}
									if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Json != nil {
										f9elemf7f5f1elemf1f3f2 := &svcapitypes.ResponseInspectionJSON{}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Json.FailureValues != nil {
											f9elemf7f5f1elemf1f3f2.FailureValues = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Json.FailureValues)
										}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Json.Identifier != nil {
											f9elemf7f5f1elemf1f3f2.Identifier = f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Json.Identifier
										}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Json.SuccessValues != nil {
											f9elemf7f5f1elemf1f3f2.SuccessValues = aws.StringSlice(f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.Json.SuccessValues)
										}
										f9elemf7f5f1elemf1f3.JSON = f9elemf7f5f1elemf1f3f2
									}
									if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.StatusCode != nil {
										f9elemf7f5f1elemf1f3f3 := &svcapitypes.ResponseInspectionStatusCode{}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.StatusCode.FailureCodes != nil {
											f9elemf7f5f1elemf1f3f3f0 := []*int64{}
											for _, f9elemf7f5f1elemf1f3f3f0iter := range f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.StatusCode.FailureCodes {
												var f9elemf7f5f1elemf1f3f3f0elem *int64
												failureCodeCopy := int64(f9elemf7f5f1elemf1f3f3f0iter)
												f9elemf7f5f1elemf1f3f3f0elem = &failureCodeCopy
												f9elemf7f5f1elemf1f3f3f0 = append(f9elemf7f5f1elemf1f3f3f0, f9elemf7f5f1elemf1f3f3f0elem)
											}
											f9elemf7f5f1elemf1f3f3.FailureCodes = f9elemf7f5f1elemf1f3f3f0
										}
										if f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.StatusCode.SuccessCodes != nil {
											f9elemf7f5f1elemf1f3f3f1 := []*int64{}
											for _, f9elemf7f5f1elemf1f3f3f1iter := range f9elemf7f5f1iter.AWSManagedRulesATPRuleSet.ResponseInspection.StatusCode.SuccessCodes {
												var f9elemf7f5f1elemf1f3f3f1elem *int64
												successCodeCopy := int64(f9elemf7f5f1elemf1f3f3f1iter)
												f9elemf7f5f1elemf1f3f3f1elem = &successCodeCopy
												f9elemf7f5f1elemf1f3f3f1 = append(f9elemf7f5f1elemf1f3f3f1, f9elemf7f5f1elemf1f3f3f1elem)
											}
											f9elemf7f5f1elemf1f3f3.SuccessCodes = f9elemf7f5f1elemf1f3f3f1
										}
										f9elemf7f5f1elemf1f3.StatusCode = f9elemf7f5f1elemf1f3f3
									}
									f9elemf7f5f1elemf1.ResponseInspection = f9elemf7f5f1elemf1f3
								}
								f9elemf7f5f1elem.AWSManagedRulesATPRuleSet = f9elemf7f5f1elemf1
							}
							if f9elemf7f5f1iter.AWSManagedRulesBotControlRuleSet != nil {
								f9elemf7f5f1elemf2 := &svcapitypes.AWSManagedRulesBotControlRuleSet{}
								if f9elemf7f5f1iter.AWSManagedRulesBotControlRuleSet.EnableMachineLearning != nil {
									f9elemf7f5f1elemf2.EnableMachineLearning = f9elemf7f5f1iter.AWSManagedRulesBotControlRuleSet.EnableMachineLearning
								}
								if f9elemf7f5f1iter.AWSManagedRulesBotControlRuleSet.InspectionLevel != "" {
									f9elemf7f5f1elemf2.InspectionLevel = aws.String(string(f9elemf7f5f1iter.AWSManagedRulesBotControlRuleSet.InspectionLevel))
								}
								f9elemf7f5f1elem.AWSManagedRulesBotControlRuleSet = f9elemf7f5f1elemf2
							}
							if f9elemf7f5f1iter.LoginPath != nil {
								f9elemf7f5f1elem.LoginPath = f9elemf7f5f1iter.LoginPath
							}
							if f9elemf7f5f1iter.PasswordField != nil {
								f9elemf7f5f1elemf4 := &svcapitypes.PasswordField{}
								if f9elemf7f5f1iter.PasswordField.Identifier != nil {
									f9elemf7f5f1elemf4.Identifier = f9elemf7f5f1iter.PasswordField.Identifier
								}
								f9elemf7f5f1elem.PasswordField = f9elemf7f5f1elemf4
							}
							if f9elemf7f5f1iter.PayloadType != "" {
								f9elemf7f5f1elem.PayloadType = aws.String(string(f9elemf7f5f1iter.PayloadType))
							}
							if f9elemf7f5f1iter.UsernameField != nil {
								f9elemf7f5f1elemf6 := &svcapitypes.UsernameField{}
								if f9elemf7f5f1iter.UsernameField.Identifier != nil {
									f9elemf7f5f1elemf6.Identifier = f9elemf7f5f1iter.UsernameField.Identifier
								}
								f9elemf7f5f1elem.UsernameField = f9elemf7f5f1elemf6
							}
							f9elemf7f5f1 = append(f9elemf7f5f1, f9elemf7f5f1elem)
						}
						f9elemf7f5.ManagedRuleGroupConfigs = f9elemf7f5f1
					}
					if f9iter.Statement.ManagedRuleGroupStatement.Name != nil {
						f9elemf7f5.Name = f9iter.Statement.ManagedRuleGroupStatement.Name
					}
					if f9iter.Statement.ManagedRuleGroupStatement.RuleActionOverrides != nil {
						f9elemf7f5f3 := []*svcapitypes.RuleActionOverride{}
						for _, f9elemf7f5f3iter := range f9iter.Statement.ManagedRuleGroupStatement.RuleActionOverrides {
							f9elemf7f5f3elem := &svcapitypes.RuleActionOverride{}
							if f9elemf7f5f3iter.ActionToUse != nil {
								f9elemf7f5f3elemf0 := &svcapitypes.RuleAction{}
								if f9elemf7f5f3iter.ActionToUse.Allow != nil {
									f9elemf7f5f3elemf0f0 := &svcapitypes.AllowAction{}
									if f9elemf7f5f3iter.ActionToUse.Allow.CustomRequestHandling != nil {
										f9elemf7f5f3elemf0f0f0 := &svcapitypes.CustomRequestHandling{}
										if f9elemf7f5f3iter.ActionToUse.Allow.CustomRequestHandling.InsertHeaders != nil {
											f9elemf7f5f3elemf0f0f0f0 := []*svcapitypes.CustomHTTPHeader{}
											for _, f9elemf7f5f3elemf0f0f0f0iter := range f9elemf7f5f3iter.ActionToUse.Allow.CustomRequestHandling.InsertHeaders {
												f9elemf7f5f3elemf0f0f0f0elem := &svcapitypes.CustomHTTPHeader{}
												if f9elemf7f5f3elemf0f0f0f0iter.Name != nil {
													f9elemf7f5f3elemf0f0f0f0elem.Name = f9elemf7f5f3elemf0f0f0f0iter.Name
												}
												if f9elemf7f5f3elemf0f0f0f0iter.Value != nil {
													f9elemf7f5f3elemf0f0f0f0elem.Value = f9elemf7f5f3elemf0f0f0f0iter.Value
												}
												f9elemf7f5f3elemf0f0f0f0 = append(f9elemf7f5f3elemf0f0f0f0, f9elemf7f5f3elemf0f0f0f0elem)
											}
											f9elemf7f5f3elemf0f0f0.InsertHeaders = f9elemf7f5f3elemf0f0f0f0
										}
										f9elemf7f5f3elemf0f0.CustomRequestHandling = f9elemf7f5f3elemf0f0f0
									}
									f9elemf7f5f3elemf0.Allow = f9elemf7f5f3elemf0f0
								}
								if f9elemf7f5f3iter.ActionToUse.Block != nil {
									f9elemf7f5f3elemf0f1 := &svcapitypes.BlockAction{}
									if f9elemf7f5f3iter.ActionToUse.Block.CustomResponse != nil {
										f9elemf7f5f3elemf0f1f0 := &svcapitypes.CustomResponse{}
										if f9elemf7f5f3iter.ActionToUse.Block.CustomResponse.CustomResponseBodyKey != nil {
											f9elemf7f5f3elemf0f1f0.CustomResponseBodyKey = f9elemf7f5f3iter.ActionToUse.Block.CustomResponse.CustomResponseBodyKey
										}
										if f9elemf7f5f3iter.ActionToUse.Block.CustomResponse.ResponseCode != nil {
											responseCodeCopy := int64(*f9elemf7f5f3iter.ActionToUse.Block.CustomResponse.ResponseCode)
											f9elemf7f5f3elemf0f1f0.ResponseCode = &responseCodeCopy
										}
										if f9elemf7f5f3iter.ActionToUse.Block.CustomResponse.ResponseHeaders != nil {
											f9elemf7f5f3elemf0f1f0f2 := []*svcapitypes.CustomHTTPHeader{}
											for _, f9elemf7f5f3elemf0f1f0f2iter := range f9elemf7f5f3iter.ActionToUse.Block.CustomResponse.ResponseHeaders {
												f9elemf7f5f3elemf0f1f0f2elem := &svcapitypes.CustomHTTPHeader{}
												if f9elemf7f5f3elemf0f1f0f2iter.Name != nil {
													f9elemf7f5f3elemf0f1f0f2elem.Name = f9elemf7f5f3elemf0f1f0f2iter.Name
												}
												if f9elemf7f5f3elemf0f1f0f2iter.Value != nil {
													f9elemf7f5f3elemf0f1f0f2elem.Value = f9elemf7f5f3elemf0f1f0f2iter.Value
												}
												f9elemf7f5f3elemf0f1f0f2 = append(f9elemf7f5f3elemf0f1f0f2, f9elemf7f5f3elemf0f1f0f2elem)
											}
											f9elemf7f5f3elemf0f1f0.ResponseHeaders = f9elemf7f5f3elemf0f1f0f2
										}
										f9elemf7f5f3elemf0f1.CustomResponse = f9elemf7f5f3elemf0f1f0
									}
									f9elemf7f5f3elemf0.Block = f9elemf7f5f3elemf0f1
								}
								if f9elemf7f5f3iter.ActionToUse.Captcha != nil {
									f9elemf7f5f3elemf0f2 := &svcapitypes.CaptchaAction{}
									if f9elemf7f5f3iter.ActionToUse.Captcha.CustomRequestHandling != nil {
										f9elemf7f5f3elemf0f2f0 := &svcapitypes.CustomRequestHandling{}
										if f9elemf7f5f3iter.ActionToUse.Captcha.CustomRequestHandling.InsertHeaders != nil {
											f9elemf7f5f3elemf0f2f0f0 := []*svcapitypes.CustomHTTPHeader{}
											for _, f9elemf7f5f3elemf0f2f0f0iter := range f9elemf7f5f3iter.ActionToUse.Captcha.CustomRequestHandling.InsertHeaders {
												f9elemf7f5f3elemf0f2f0f0elem := &svcapitypes.CustomHTTPHeader{}
												if f9elemf7f5f3elemf0f2f0f0iter.Name != nil {
													f9elemf7f5f3elemf0f2f0f0elem.Name = f9elemf7f5f3elemf0f2f0f0iter.Name
												}
												if f9elemf7f5f3elemf0f2f0f0iter.Value != nil {
													f9elemf7f5f3elemf0f2f0f0elem.Value = f9elemf7f5f3elemf0f2f0f0iter.Value
												}
												f9elemf7f5f3elemf0f2f0f0 = append(f9elemf7f5f3elemf0f2f0f0, f9elemf7f5f3elemf0f2f0f0elem)
											}
											f9elemf7f5f3elemf0f2f0.InsertHeaders = f9elemf7f5f3elemf0f2f0f0
										}
										f9elemf7f5f3elemf0f2.CustomRequestHandling = f9elemf7f5f3elemf0f2f0
									}
									f9elemf7f5f3elemf0.Captcha = f9elemf7f5f3elemf0f2
								}
								if f9elemf7f5f3iter.ActionToUse.Challenge != nil {
									f9elemf7f5f3elemf0f3 := &svcapitypes.ChallengeAction{}
									if f9elemf7f5f3iter.ActionToUse.Challenge.CustomRequestHandling != nil {
										f9elemf7f5f3elemf0f3f0 := &svcapitypes.CustomRequestHandling{}
										if f9elemf7f5f3iter.ActionToUse.Challenge.CustomRequestHandling.InsertHeaders != nil {
											f9elemf7f5f3elemf0f3f0f0 := []*svcapitypes.CustomHTTPHeader{}
											for _, f9elemf7f5f3elemf0f3f0f0iter := range f9elemf7f5f3iter.ActionToUse.Challenge.CustomRequestHandling.InsertHeaders {
												f9elemf7f5f3elemf0f3f0f0elem := &svcapitypes.CustomHTTPHeader{}
												if f9elemf7f5f3elemf0f3f0f0iter.Name != nil {
													f9elemf7f5f3elemf0f3f0f0elem.Name = f9elemf7f5f3elemf0f3f0f0iter.Name
												}
												if f9elemf7f5f3elemf0f3f0f0iter.Value != nil {
													f9elemf7f5f3elemf0f3f0f0elem.Value = f9elemf7f5f3elemf0f3f0f0iter.Value
												}
												f9elemf7f5f3elemf0f3f0f0 = append(f9elemf7f5f3elemf0f3f0f0, f9elemf7f5f3elemf0f3f0f0elem)
											}
											f9elemf7f5f3elemf0f3f0.InsertHeaders = f9elemf7f5f3elemf0f3f0f0
										}
										f9elemf7f5f3elemf0f3.CustomRequestHandling = f9elemf7f5f3elemf0f3f0
									}
									f9elemf7f5f3elemf0.Challenge = f9elemf7f5f3elemf0f3
								}
								if f9elemf7f5f3iter.ActionToUse.Count != nil {
									f9elemf7f5f3elemf0f4 := &svcapitypes.CountAction{}
									if f9elemf7f5f3iter.ActionToUse.Count.CustomRequestHandling != nil {
										f9elemf7f5f3elemf0f4f0 := &svcapitypes.CustomRequestHandling{}
										if f9elemf7f5f3iter.ActionToUse.Count.CustomRequestHandling.InsertHeaders != nil {
											f9elemf7f5f3elemf0f4f0f0 := []*svcapitypes.CustomHTTPHeader{}
											for _, f9elemf7f5f3elemf0f4f0f0iter := range f9elemf7f5f3iter.ActionToUse.Count.CustomRequestHandling.InsertHeaders {
												f9elemf7f5f3elemf0f4f0f0elem := &svcapitypes.CustomHTTPHeader{}
												if f9elemf7f5f3elemf0f4f0f0iter.Name != nil {
													f9elemf7f5f3elemf0f4f0f0elem.Name = f9elemf7f5f3elemf0f4f0f0iter.Name
												}
												if f9elemf7f5f3elemf0f4f0f0iter.Value != nil {
													f9elemf7f5f3elemf0f4f0f0elem.Value = f9elemf7f5f3elemf0f4f0f0iter.Value
												}
												f9elemf7f5f3elemf0f4f0f0 = append(f9elemf7f5f3elemf0f4f0f0, f9elemf7f5f3elemf0f4f0f0elem)
											}
											f9elemf7f5f3elemf0f4f0.InsertHeaders = f9elemf7f5f3elemf0f4f0f0
										}
										f9elemf7f5f3elemf0f4.CustomRequestHandling = f9elemf7f5f3elemf0f4f0
									}
									f9elemf7f5f3elemf0.Count = f9elemf7f5f3elemf0f4
								}
								f9elemf7f5f3elem.ActionToUse = f9elemf7f5f3elemf0
							}
							if f9elemf7f5f3iter.Name != nil {
								f9elemf7f5f3elem.Name = f9elemf7f5f3iter.Name
							}
							f9elemf7f5f3 = append(f9elemf7f5f3, f9elemf7f5f3elem)
						}
						f9elemf7f5.RuleActionOverrides = f9elemf7f5f3
					}
					if f9iter.Statement.ManagedRuleGroupStatement.VendorName != nil {
						f9elemf7f5.VendorName = f9iter.Statement.ManagedRuleGroupStatement.VendorName
					}
					if f9iter.Statement.ManagedRuleGroupStatement.Version != nil {
						f9elemf7f5.Version = f9iter.Statement.ManagedRuleGroupStatement.Version
					}
					f9elemf7.ManagedRuleGroupStatement = f9elemf7f5
				}
				if f9iter.Statement.RateBasedStatement != nil {
					f9elemf7f8 := &svcapitypes.RateBasedStatement{}
					if f9iter.Statement.RateBasedStatement.AggregateKeyType != "" {
						f9elemf7f8.AggregateKeyType = aws.String(string(f9iter.Statement.RateBasedStatement.AggregateKeyType))
					}
					if f9iter.Statement.RateBasedStatement.CustomKeys != nil {
						f9elemf7f8f1 := []*svcapitypes.RateBasedStatementCustomKey{}
						for _, f9elemf7f8f1iter := range f9iter.Statement.RateBasedStatement.CustomKeys {
							f9elemf7f8f1elem := &svcapitypes.RateBasedStatementCustomKey{}
							if f9elemf7f8f1iter.Cookie != nil {
								f9elemf7f8f1elemf0 := &svcapitypes.RateLimitCookie{}
								if f9elemf7f8f1iter.Cookie.Name != nil {
									f9elemf7f8f1elemf0.Name = f9elemf7f8f1iter.Cookie.Name
								}
								if f9elemf7f8f1iter.Cookie.TextTransformations != nil {
									f9elemf7f8f1elemf0f1 := []*svcapitypes.TextTransformation{}
									for _, f9elemf7f8f1elemf0f1iter := range f9elemf7f8f1iter.Cookie.TextTransformations {
										f9elemf7f8f1elemf0f1elem := &svcapitypes.TextTransformation{}
										priorityCopy := int64(f9elemf7f8f1elemf0f1iter.Priority)
										f9elemf7f8f1elemf0f1elem.Priority = &priorityCopy
										if f9elemf7f8f1elemf0f1iter.Type != "" {
											f9elemf7f8f1elemf0f1elem.Type = aws.String(string(f9elemf7f8f1elemf0f1iter.Type))
										}
										f9elemf7f8f1elemf0f1 = append(f9elemf7f8f1elemf0f1, f9elemf7f8f1elemf0f1elem)
									}
									f9elemf7f8f1elemf0.TextTransformations = f9elemf7f8f1elemf0f1
								}
								f9elemf7f8f1elem.Cookie = f9elemf7f8f1elemf0
							}
							if f9elemf7f8f1iter.ForwardedIP != nil {
								f9elemf7f8f1elemf1 := map[string]*string{}
								f9elemf7f8f1elem.ForwardedIP = f9elemf7f8f1elemf1
							}
							if f9elemf7f8f1iter.HTTPMethod != nil {
								f9elemf7f8f1elemf2 := map[string]*string{}
								f9elemf7f8f1elem.HTTPMethod = f9elemf7f8f1elemf2
							}
							if f9elemf7f8f1iter.Header != nil {
								f9elemf7f8f1elemf3 := &svcapitypes.RateLimitHeader{}
								if f9elemf7f8f1iter.Header.Name != nil {
									f9elemf7f8f1elemf3.Name = f9elemf7f8f1iter.Header.Name
								}
								if f9elemf7f8f1iter.Header.TextTransformations != nil {
									f9elemf7f8f1elemf3f1 := []*svcapitypes.TextTransformation{}
									for _, f9elemf7f8f1elemf3f1iter := range f9elemf7f8f1iter.Header.TextTransformations {
										f9elemf7f8f1elemf3f1elem := &svcapitypes.TextTransformation{}
										priorityCopy := int64(f9elemf7f8f1elemf3f1iter.Priority)
										f9elemf7f8f1elemf3f1elem.Priority = &priorityCopy
										if f9elemf7f8f1elemf3f1iter.Type != "" {
											f9elemf7f8f1elemf3f1elem.Type = aws.String(string(f9elemf7f8f1elemf3f1iter.Type))
										}
										f9elemf7f8f1elemf3f1 = append(f9elemf7f8f1elemf3f1, f9elemf7f8f1elemf3f1elem)
									}
									f9elemf7f8f1elemf3.TextTransformations = f9elemf7f8f1elemf3f1
								}
								f9elemf7f8f1elem.Header = f9elemf7f8f1elemf3
							}
							if f9elemf7f8f1iter.IP != nil {
								f9elemf7f8f1elemf4 := map[string]*string{}
								f9elemf7f8f1elem.IP = f9elemf7f8f1elemf4
							}
							if f9elemf7f8f1iter.LabelNamespace != nil {
								f9elemf7f8f1elemf5 := &svcapitypes.RateLimitLabelNamespace{}
								if f9elemf7f8f1iter.LabelNamespace.Namespace != nil {
									f9elemf7f8f1elemf5.Namespace = f9elemf7f8f1iter.LabelNamespace.Namespace
								}
								f9elemf7f8f1elem.LabelNamespace = f9elemf7f8f1elemf5
							}
							if f9elemf7f8f1iter.QueryArgument != nil {
								f9elemf7f8f1elemf6 := &svcapitypes.RateLimitQueryArgument{}
								if f9elemf7f8f1iter.QueryArgument.Name != nil {
									f9elemf7f8f1elemf6.Name = f9elemf7f8f1iter.QueryArgument.Name
								}
								if f9elemf7f8f1iter.QueryArgument.TextTransformations != nil {
									f9elemf7f8f1elemf6f1 := []*svcapitypes.TextTransformation{}
									for _, f9elemf7f8f1elemf6f1iter := range f9elemf7f8f1iter.QueryArgument.TextTransformations {
										f9elemf7f8f1elemf6f1elem := &svcapitypes.TextTransformation{}
										priorityCopy := int64(f9elemf7f8f1elemf6f1iter.Priority)
										f9elemf7f8f1elemf6f1elem.Priority = &priorityCopy
										if f9elemf7f8f1elemf6f1iter.Type != "" {
											f9elemf7f8f1elemf6f1elem.Type = aws.String(string(f9elemf7f8f1elemf6f1iter.Type))
										}
										f9elemf7f8f1elemf6f1 = append(f9elemf7f8f1elemf6f1, f9elemf7f8f1elemf6f1elem)
									}
									f9elemf7f8f1elemf6.TextTransformations = f9elemf7f8f1elemf6f1
								}
								f9elemf7f8f1elem.QueryArgument = f9elemf7f8f1elemf6
							}
							if f9elemf7f8f1iter.QueryString != nil {
								f9elemf7f8f1elemf7 := &svcapitypes.RateLimitQueryString{}
								if f9elemf7f8f1iter.QueryString.TextTransformations != nil {
									f9elemf7f8f1elemf7f0 := []*svcapitypes.TextTransformation{}
									for _, f9elemf7f8f1elemf7f0iter := range f9elemf7f8f1iter.QueryString.TextTransformations {
										f9elemf7f8f1elemf7f0elem := &svcapitypes.TextTransformation{}
										priorityCopy := int64(f9elemf7f8f1elemf7f0iter.Priority)
										f9elemf7f8f1elemf7f0elem.Priority = &priorityCopy
										if f9elemf7f8f1elemf7f0iter.Type != "" {
											f9elemf7f8f1elemf7f0elem.Type = aws.String(string(f9elemf7f8f1elemf7f0iter.Type))
										}
										f9elemf7f8f1elemf7f0 = append(f9elemf7f8f1elemf7f0, f9elemf7f8f1elemf7f0elem)
									}
									f9elemf7f8f1elemf7.TextTransformations = f9elemf7f8f1elemf7f0
								}
								f9elemf7f8f1elem.QueryString = f9elemf7f8f1elemf7
							}
							if f9elemf7f8f1iter.UriPath != nil {
								f9elemf7f8f1elemf8 := &svcapitypes.RateLimitURIPath{}
								if f9elemf7f8f1iter.UriPath.TextTransformations != nil {
									f9elemf7f8f1elemf8f0 := []*svcapitypes.TextTransformation{}
									for _, f9elemf7f8f1elemf8f0iter := range f9elemf7f8f1iter.UriPath.TextTransformations {
										f9elemf7f8f1elemf8f0elem := &svcapitypes.TextTransformation{}
										priorityCopy := int64(f9elemf7f8f1elemf8f0iter.Priority)
										f9elemf7f8f1elemf8f0elem.Priority = &priorityCopy
										if f9elemf7f8f1elemf8f0iter.Type != "" {
											f9elemf7f8f1elemf8f0elem.Type = aws.String(string(f9elemf7f8f1elemf8f0iter.Type))
										}
										f9elemf7f8f1elemf8f0 = append(f9elemf7f8f1elemf8f0, f9elemf7f8f1elemf8f0elem)
									}
									f9elemf7f8f1elemf8.TextTransformations = f9elemf7f8f1elemf8f0
								}
								f9elemf7f8f1elem.URIPath = f9elemf7f8f1elemf8
							}
							f9elemf7f8f1 = append(f9elemf7f8f1, f9elemf7f8f1elem)
						}
						f9elemf7f8.CustomKeys = f9elemf7f8f1
					}
					f9elemf7f8.EvaluationWindowSec = &f9iter.Statement.RateBasedStatement.EvaluationWindowSec
					if f9iter.Statement.RateBasedStatement.ForwardedIPConfig != nil {
						f9elemf7f8f3 := &svcapitypes.ForwardedIPConfig{}
						if f9iter.Statement.RateBasedStatement.ForwardedIPConfig.FallbackBehavior != "" {
							f9elemf7f8f3.FallbackBehavior = aws.String(string(f9iter.Statement.RateBasedStatement.ForwardedIPConfig.FallbackBehavior))
						}
						if f9iter.Statement.RateBasedStatement.ForwardedIPConfig.HeaderName != nil {
							f9elemf7f8f3.HeaderName = f9iter.Statement.RateBasedStatement.ForwardedIPConfig.HeaderName
						}
						f9elemf7f8.ForwardedIPConfig = f9elemf7f8f3
					}
					if f9iter.Statement.RateBasedStatement.Limit != nil {
						f9elemf7f8.Limit = f9iter.Statement.RateBasedStatement.Limit
					}
					f9elemf7.RateBasedStatement = f9elemf7f8
				}
				if f9iter.Statement.RegexMatchStatement != nil {
					f9elemf7f9 := &svcapitypes.RegexMatchStatement{}
					if f9iter.Statement.RegexMatchStatement.FieldToMatch != nil {
						f9elemf7f9f0 := &svcapitypes.FieldToMatch{}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.AllQueryArguments != nil {
							f9elemf7f9f0f0 := map[string]*string{}
							f9elemf7f9f0.AllQueryArguments = f9elemf7f9f0f0
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.Body != nil {
							f9elemf7f9f0f1 := &svcapitypes.Body{}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.Body.OversizeHandling != "" {
								f9elemf7f9f0f1.OversizeHandling = aws.String(string(f9iter.Statement.RegexMatchStatement.FieldToMatch.Body.OversizeHandling))
							}
							f9elemf7f9f0.Body = f9elemf7f9f0f1
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies != nil {
							f9elemf7f9f0f2 := &svcapitypes.Cookies{}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f9elemf7f9f0f2f0 := &svcapitypes.CookieMatchPattern{}
								if f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f9elemf7f9f0f2f0f0 := map[string]*string{}
									f9elemf7f9f0f2f0.All = f9elemf7f9f0f2f0f0
								}
								if f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f9elemf7f9f0f2f0.ExcludedCookies = aws.StringSlice(f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f9elemf7f9f0f2f0.IncludedCookies = aws.StringSlice(f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f9elemf7f9f0f2.MatchPattern = f9elemf7f9f0f2f0
							}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchScope != "" {
								f9elemf7f9f0f2.MatchScope = aws.String(string(f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.MatchScope))
							}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.OversizeHandling != "" {
								f9elemf7f9f0f2.OversizeHandling = aws.String(string(f9iter.Statement.RegexMatchStatement.FieldToMatch.Cookies.OversizeHandling))
							}
							f9elemf7f9f0.Cookies = f9elemf7f9f0f2
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.HeaderOrder != nil {
							f9elemf7f9f0f3 := &svcapitypes.HeaderOrder{}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling != "" {
								f9elemf7f9f0f3.OversizeHandling = aws.String(string(f9iter.Statement.RegexMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling))
							}
							f9elemf7f9f0.HeaderOrder = f9elemf7f9f0f3
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers != nil {
							f9elemf7f9f0f4 := &svcapitypes.Headers{}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern != nil {
								f9elemf7f9f0f4f0 := &svcapitypes.HeaderMatchPattern{}
								if f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f9elemf7f9f0f4f0f0 := map[string]*string{}
									f9elemf7f9f0f4f0.All = f9elemf7f9f0f4f0f0
								}
								if f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f9elemf7f9f0f4f0.ExcludedHeaders = aws.StringSlice(f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f9elemf7f9f0f4f0.IncludedHeaders = aws.StringSlice(f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f9elemf7f9f0f4.MatchPattern = f9elemf7f9f0f4f0
							}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchScope != "" {
								f9elemf7f9f0f4.MatchScope = aws.String(string(f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers.MatchScope))
							}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers.OversizeHandling != "" {
								f9elemf7f9f0f4.OversizeHandling = aws.String(string(f9iter.Statement.RegexMatchStatement.FieldToMatch.Headers.OversizeHandling))
							}
							f9elemf7f9f0.Headers = f9elemf7f9f0f4
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.JA3Fingerprint != nil {
							f9elemf7f9f0f5 := &svcapitypes.JA3Fingerprint{}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != "" {
								f9elemf7f9f0f5.FallbackBehavior = aws.String(string(f9iter.Statement.RegexMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior))
							}
							f9elemf7f9f0.JA3Fingerprint = f9elemf7f9f0f5
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody != nil {
							f9elemf7f9f0f6 := &svcapitypes.JSONBody{}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior != "" {
								f9elemf7f9f0f6.InvalidFallbackBehavior = aws.String(string(f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior))
							}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody.MatchPattern != nil {
								f9elemf7f9f0f6f1 := &svcapitypes.JSONMatchPattern{}
								if f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody.MatchPattern.All != nil {
									f9elemf7f9f0f6f1f0 := map[string]*string{}
									f9elemf7f9f0f6f1.All = f9elemf7f9f0f6f1f0
								}
								if f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths != nil {
									f9elemf7f9f0f6f1.IncludedPaths = aws.StringSlice(f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths)
								}
								f9elemf7f9f0f6.MatchPattern = f9elemf7f9f0f6f1
							}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody.MatchScope != "" {
								f9elemf7f9f0f6.MatchScope = aws.String(string(f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody.MatchScope))
							}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody.OversizeHandling != "" {
								f9elemf7f9f0f6.OversizeHandling = aws.String(string(f9iter.Statement.RegexMatchStatement.FieldToMatch.JsonBody.OversizeHandling))
							}
							f9elemf7f9f0.JSONBody = f9elemf7f9f0f6
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.Method != nil {
							f9elemf7f9f0f7 := map[string]*string{}
							f9elemf7f9f0.Method = f9elemf7f9f0f7
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.QueryString != nil {
							f9elemf7f9f0f8 := map[string]*string{}
							f9elemf7f9f0.QueryString = f9elemf7f9f0f8
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.SingleHeader != nil {
							f9elemf7f9f0f9 := &svcapitypes.SingleHeader{}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.SingleHeader.Name != nil {
								f9elemf7f9f0f9.Name = f9iter.Statement.RegexMatchStatement.FieldToMatch.SingleHeader.Name
							}
							f9elemf7f9f0.SingleHeader = f9elemf7f9f0f9
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.SingleQueryArgument != nil {
							f9elemf7f9f0f10 := &svcapitypes.SingleQueryArgument{}
							if f9iter.Statement.RegexMatchStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f9elemf7f9f0f10.Name = f9iter.Statement.RegexMatchStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f9elemf7f9f0.SingleQueryArgument = f9elemf7f9f0f10
						}
						if f9iter.Statement.RegexMatchStatement.FieldToMatch.UriPath != nil {
							f9elemf7f9f0f11 := map[string]*string{}
							f9elemf7f9f0.URIPath = f9elemf7f9f0f11
						}
						f9elemf7f9.FieldToMatch = f9elemf7f9f0
					}
					if f9iter.Statement.RegexMatchStatement.RegexString != nil {
						f9elemf7f9.RegexString = f9iter.Statement.RegexMatchStatement.RegexString
					}
					if f9iter.Statement.RegexMatchStatement.TextTransformations != nil {
						f9elemf7f9f2 := []*svcapitypes.TextTransformation{}
						for _, f9elemf7f9f2iter := range f9iter.Statement.RegexMatchStatement.TextTransformations {
							f9elemf7f9f2elem := &svcapitypes.TextTransformation{}
							priorityCopy := int64(f9elemf7f9f2iter.Priority)
							f9elemf7f9f2elem.Priority = &priorityCopy
							if f9elemf7f9f2iter.Type != "" {
								f9elemf7f9f2elem.Type = aws.String(string(f9elemf7f9f2iter.Type))
							}
							f9elemf7f9f2 = append(f9elemf7f9f2, f9elemf7f9f2elem)
						}
						f9elemf7f9.TextTransformations = f9elemf7f9f2
					}
					f9elemf7.RegexMatchStatement = f9elemf7f9
				}
				if f9iter.Statement.RegexPatternSetReferenceStatement != nil {
					f9elemf7f10 := &svcapitypes.RegexPatternSetReferenceStatement{}
					if f9iter.Statement.RegexPatternSetReferenceStatement.ARN != nil {
						f9elemf7f10.ARN = f9iter.Statement.RegexPatternSetReferenceStatement.ARN
					}
					if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch != nil {
						f9elemf7f10f1 := &svcapitypes.FieldToMatch{}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.AllQueryArguments != nil {
							f9elemf7f10f1f0 := map[string]*string{}
							f9elemf7f10f1.AllQueryArguments = f9elemf7f10f1f0
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Body != nil {
							f9elemf7f10f1f1 := &svcapitypes.Body{}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Body.OversizeHandling != "" {
								f9elemf7f10f1f1.OversizeHandling = aws.String(string(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Body.OversizeHandling))
							}
							f9elemf7f10f1.Body = f9elemf7f10f1f1
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies != nil {
							f9elemf7f10f1f2 := &svcapitypes.Cookies{}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f9elemf7f10f1f2f0 := &svcapitypes.CookieMatchPattern{}
								if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f9elemf7f10f1f2f0f0 := map[string]*string{}
									f9elemf7f10f1f2f0.All = f9elemf7f10f1f2f0f0
								}
								if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f9elemf7f10f1f2f0.ExcludedCookies = aws.StringSlice(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f9elemf7f10f1f2f0.IncludedCookies = aws.StringSlice(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f9elemf7f10f1f2.MatchPattern = f9elemf7f10f1f2f0
							}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchScope != "" {
								f9elemf7f10f1f2.MatchScope = aws.String(string(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.MatchScope))
							}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.OversizeHandling != "" {
								f9elemf7f10f1f2.OversizeHandling = aws.String(string(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Cookies.OversizeHandling))
							}
							f9elemf7f10f1.Cookies = f9elemf7f10f1f2
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.HeaderOrder != nil {
							f9elemf7f10f1f3 := &svcapitypes.HeaderOrder{}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.HeaderOrder.OversizeHandling != "" {
								f9elemf7f10f1f3.OversizeHandling = aws.String(string(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.HeaderOrder.OversizeHandling))
							}
							f9elemf7f10f1.HeaderOrder = f9elemf7f10f1f3
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers != nil {
							f9elemf7f10f1f4 := &svcapitypes.Headers{}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern != nil {
								f9elemf7f10f1f4f0 := &svcapitypes.HeaderMatchPattern{}
								if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f9elemf7f10f1f4f0f0 := map[string]*string{}
									f9elemf7f10f1f4f0.All = f9elemf7f10f1f4f0f0
								}
								if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f9elemf7f10f1f4f0.ExcludedHeaders = aws.StringSlice(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f9elemf7f10f1f4f0.IncludedHeaders = aws.StringSlice(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f9elemf7f10f1f4.MatchPattern = f9elemf7f10f1f4f0
							}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchScope != "" {
								f9elemf7f10f1f4.MatchScope = aws.String(string(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.MatchScope))
							}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.OversizeHandling != "" {
								f9elemf7f10f1f4.OversizeHandling = aws.String(string(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Headers.OversizeHandling))
							}
							f9elemf7f10f1.Headers = f9elemf7f10f1f4
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JA3Fingerprint != nil {
							f9elemf7f10f1f5 := &svcapitypes.JA3Fingerprint{}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != "" {
								f9elemf7f10f1f5.FallbackBehavior = aws.String(string(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior))
							}
							f9elemf7f10f1.JA3Fingerprint = f9elemf7f10f1f5
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody != nil {
							f9elemf7f10f1f6 := &svcapitypes.JSONBody{}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior != "" {
								f9elemf7f10f1f6.InvalidFallbackBehavior = aws.String(string(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior))
							}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody.MatchPattern != nil {
								f9elemf7f10f1f6f1 := &svcapitypes.JSONMatchPattern{}
								if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody.MatchPattern.All != nil {
									f9elemf7f10f1f6f1f0 := map[string]*string{}
									f9elemf7f10f1f6f1.All = f9elemf7f10f1f6f1f0
								}
								if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths != nil {
									f9elemf7f10f1f6f1.IncludedPaths = aws.StringSlice(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths)
								}
								f9elemf7f10f1f6.MatchPattern = f9elemf7f10f1f6f1
							}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody.MatchScope != "" {
								f9elemf7f10f1f6.MatchScope = aws.String(string(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody.MatchScope))
							}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody.OversizeHandling != "" {
								f9elemf7f10f1f6.OversizeHandling = aws.String(string(f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.JsonBody.OversizeHandling))
							}
							f9elemf7f10f1.JSONBody = f9elemf7f10f1f6
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.Method != nil {
							f9elemf7f10f1f7 := map[string]*string{}
							f9elemf7f10f1.Method = f9elemf7f10f1f7
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.QueryString != nil {
							f9elemf7f10f1f8 := map[string]*string{}
							f9elemf7f10f1.QueryString = f9elemf7f10f1f8
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleHeader != nil {
							f9elemf7f10f1f9 := &svcapitypes.SingleHeader{}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleHeader.Name != nil {
								f9elemf7f10f1f9.Name = f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleHeader.Name
							}
							f9elemf7f10f1.SingleHeader = f9elemf7f10f1f9
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleQueryArgument != nil {
							f9elemf7f10f1f10 := &svcapitypes.SingleQueryArgument{}
							if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f9elemf7f10f1f10.Name = f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f9elemf7f10f1.SingleQueryArgument = f9elemf7f10f1f10
						}
						if f9iter.Statement.RegexPatternSetReferenceStatement.FieldToMatch.UriPath != nil {
							f9elemf7f10f1f11 := map[string]*string{}
							f9elemf7f10f1.URIPath = f9elemf7f10f1f11
						}
						f9elemf7f10.FieldToMatch = f9elemf7f10f1
					}
					if f9iter.Statement.RegexPatternSetReferenceStatement.TextTransformations != nil {
						f9elemf7f10f2 := []*svcapitypes.TextTransformation{}
						for _, f9elemf7f10f2iter := range f9iter.Statement.RegexPatternSetReferenceStatement.TextTransformations {
							f9elemf7f10f2elem := &svcapitypes.TextTransformation{}
							priorityCopy := int64(f9elemf7f10f2iter.Priority)
							f9elemf7f10f2elem.Priority = &priorityCopy
							if f9elemf7f10f2iter.Type != "" {
								f9elemf7f10f2elem.Type = aws.String(string(f9elemf7f10f2iter.Type))
							}
							f9elemf7f10f2 = append(f9elemf7f10f2, f9elemf7f10f2elem)
						}
						f9elemf7f10.TextTransformations = f9elemf7f10f2
					}
					f9elemf7.RegexPatternSetReferenceStatement = f9elemf7f10
				}
				if f9iter.Statement.RuleGroupReferenceStatement != nil {
					f9elemf7f11 := &svcapitypes.RuleGroupReferenceStatement{}
					if f9iter.Statement.RuleGroupReferenceStatement.ARN != nil {
						f9elemf7f11.ARN = f9iter.Statement.RuleGroupReferenceStatement.ARN
					}
					if f9iter.Statement.RuleGroupReferenceStatement.ExcludedRules != nil {
						f9elemf7f11f1 := []*svcapitypes.ExcludedRule{}
						for _, f9elemf7f11f1iter := range f9iter.Statement.RuleGroupReferenceStatement.ExcludedRules {
							f9elemf7f11f1elem := &svcapitypes.ExcludedRule{}
							if f9elemf7f11f1iter.Name != nil {
								f9elemf7f11f1elem.Name = f9elemf7f11f1iter.Name
							}
							f9elemf7f11f1 = append(f9elemf7f11f1, f9elemf7f11f1elem)
						}
						f9elemf7f11.ExcludedRules = f9elemf7f11f1
					}
					if f9iter.Statement.RuleGroupReferenceStatement.RuleActionOverrides != nil {
						f9elemf7f11f2 := []*svcapitypes.RuleActionOverride{}
						for _, f9elemf7f11f2iter := range f9iter.Statement.RuleGroupReferenceStatement.RuleActionOverrides {
							f9elemf7f11f2elem := &svcapitypes.RuleActionOverride{}
							if f9elemf7f11f2iter.ActionToUse != nil {
								f9elemf7f11f2elemf0 := &svcapitypes.RuleAction{}
								if f9elemf7f11f2iter.ActionToUse.Allow != nil {
									f9elemf7f11f2elemf0f0 := &svcapitypes.AllowAction{}
									if f9elemf7f11f2iter.ActionToUse.Allow.CustomRequestHandling != nil {
										f9elemf7f11f2elemf0f0f0 := &svcapitypes.CustomRequestHandling{}
										if f9elemf7f11f2iter.ActionToUse.Allow.CustomRequestHandling.InsertHeaders != nil {
											f9elemf7f11f2elemf0f0f0f0 := []*svcapitypes.CustomHTTPHeader{}
											for _, f9elemf7f11f2elemf0f0f0f0iter := range f9elemf7f11f2iter.ActionToUse.Allow.CustomRequestHandling.InsertHeaders {
												f9elemf7f11f2elemf0f0f0f0elem := &svcapitypes.CustomHTTPHeader{}
												if f9elemf7f11f2elemf0f0f0f0iter.Name != nil {
													f9elemf7f11f2elemf0f0f0f0elem.Name = f9elemf7f11f2elemf0f0f0f0iter.Name
												}
												if f9elemf7f11f2elemf0f0f0f0iter.Value != nil {
													f9elemf7f11f2elemf0f0f0f0elem.Value = f9elemf7f11f2elemf0f0f0f0iter.Value
												}
												f9elemf7f11f2elemf0f0f0f0 = append(f9elemf7f11f2elemf0f0f0f0, f9elemf7f11f2elemf0f0f0f0elem)
											}
											f9elemf7f11f2elemf0f0f0.InsertHeaders = f9elemf7f11f2elemf0f0f0f0
										}
										f9elemf7f11f2elemf0f0.CustomRequestHandling = f9elemf7f11f2elemf0f0f0
									}
									f9elemf7f11f2elemf0.Allow = f9elemf7f11f2elemf0f0
								}
								if f9elemf7f11f2iter.ActionToUse.Block != nil {
									f9elemf7f11f2elemf0f1 := &svcapitypes.BlockAction{}
									if f9elemf7f11f2iter.ActionToUse.Block.CustomResponse != nil {
										f9elemf7f11f2elemf0f1f0 := &svcapitypes.CustomResponse{}
										if f9elemf7f11f2iter.ActionToUse.Block.CustomResponse.CustomResponseBodyKey != nil {
											f9elemf7f11f2elemf0f1f0.CustomResponseBodyKey = f9elemf7f11f2iter.ActionToUse.Block.CustomResponse.CustomResponseBodyKey
										}
										if f9elemf7f11f2iter.ActionToUse.Block.CustomResponse.ResponseCode != nil {
											responseCodeCopy := int64(*f9elemf7f11f2iter.ActionToUse.Block.CustomResponse.ResponseCode)
											f9elemf7f11f2elemf0f1f0.ResponseCode = &responseCodeCopy
										}
										if f9elemf7f11f2iter.ActionToUse.Block.CustomResponse.ResponseHeaders != nil {
											f9elemf7f11f2elemf0f1f0f2 := []*svcapitypes.CustomHTTPHeader{}
											for _, f9elemf7f11f2elemf0f1f0f2iter := range f9elemf7f11f2iter.ActionToUse.Block.CustomResponse.ResponseHeaders {
												f9elemf7f11f2elemf0f1f0f2elem := &svcapitypes.CustomHTTPHeader{}
												if f9elemf7f11f2elemf0f1f0f2iter.Name != nil {
													f9elemf7f11f2elemf0f1f0f2elem.Name = f9elemf7f11f2elemf0f1f0f2iter.Name
												}
												if f9elemf7f11f2elemf0f1f0f2iter.Value != nil {
													f9elemf7f11f2elemf0f1f0f2elem.Value = f9elemf7f11f2elemf0f1f0f2iter.Value
												}
												f9elemf7f11f2elemf0f1f0f2 = append(f9elemf7f11f2elemf0f1f0f2, f9elemf7f11f2elemf0f1f0f2elem)
											}
											f9elemf7f11f2elemf0f1f0.ResponseHeaders = f9elemf7f11f2elemf0f1f0f2
										}
										f9elemf7f11f2elemf0f1.CustomResponse = f9elemf7f11f2elemf0f1f0
									}
									f9elemf7f11f2elemf0.Block = f9elemf7f11f2elemf0f1
								}
								if f9elemf7f11f2iter.ActionToUse.Captcha != nil {
									f9elemf7f11f2elemf0f2 := &svcapitypes.CaptchaAction{}
									if f9elemf7f11f2iter.ActionToUse.Captcha.CustomRequestHandling != nil {
										f9elemf7f11f2elemf0f2f0 := &svcapitypes.CustomRequestHandling{}
										if f9elemf7f11f2iter.ActionToUse.Captcha.CustomRequestHandling.InsertHeaders != nil {
											f9elemf7f11f2elemf0f2f0f0 := []*svcapitypes.CustomHTTPHeader{}
											for _, f9elemf7f11f2elemf0f2f0f0iter := range f9elemf7f11f2iter.ActionToUse.Captcha.CustomRequestHandling.InsertHeaders {
												f9elemf7f11f2elemf0f2f0f0elem := &svcapitypes.CustomHTTPHeader{}
												if f9elemf7f11f2elemf0f2f0f0iter.Name != nil {
													f9elemf7f11f2elemf0f2f0f0elem.Name = f9elemf7f11f2elemf0f2f0f0iter.Name
												}
												if f9elemf7f11f2elemf0f2f0f0iter.Value != nil {
													f9elemf7f11f2elemf0f2f0f0elem.Value = f9elemf7f11f2elemf0f2f0f0iter.Value
												}
												f9elemf7f11f2elemf0f2f0f0 = append(f9elemf7f11f2elemf0f2f0f0, f9elemf7f11f2elemf0f2f0f0elem)
											}
											f9elemf7f11f2elemf0f2f0.InsertHeaders = f9elemf7f11f2elemf0f2f0f0
										}
										f9elemf7f11f2elemf0f2.CustomRequestHandling = f9elemf7f11f2elemf0f2f0
									}
									f9elemf7f11f2elemf0.Captcha = f9elemf7f11f2elemf0f2
								}
								if f9elemf7f11f2iter.ActionToUse.Challenge != nil {
									f9elemf7f11f2elemf0f3 := &svcapitypes.ChallengeAction{}
									if f9elemf7f11f2iter.ActionToUse.Challenge.CustomRequestHandling != nil {
										f9elemf7f11f2elemf0f3f0 := &svcapitypes.CustomRequestHandling{}
										if f9elemf7f11f2iter.ActionToUse.Challenge.CustomRequestHandling.InsertHeaders != nil {
											f9elemf7f11f2elemf0f3f0f0 := []*svcapitypes.CustomHTTPHeader{}
											for _, f9elemf7f11f2elemf0f3f0f0iter := range f9elemf7f11f2iter.ActionToUse.Challenge.CustomRequestHandling.InsertHeaders {
												f9elemf7f11f2elemf0f3f0f0elem := &svcapitypes.CustomHTTPHeader{}
												if f9elemf7f11f2elemf0f3f0f0iter.Name != nil {
													f9elemf7f11f2elemf0f3f0f0elem.Name = f9elemf7f11f2elemf0f3f0f0iter.Name
												}
												if f9elemf7f11f2elemf0f3f0f0iter.Value != nil {
													f9elemf7f11f2elemf0f3f0f0elem.Value = f9elemf7f11f2elemf0f3f0f0iter.Value
												}
												f9elemf7f11f2elemf0f3f0f0 = append(f9elemf7f11f2elemf0f3f0f0, f9elemf7f11f2elemf0f3f0f0elem)
											}
											f9elemf7f11f2elemf0f3f0.InsertHeaders = f9elemf7f11f2elemf0f3f0f0
										}
										f9elemf7f11f2elemf0f3.CustomRequestHandling = f9elemf7f11f2elemf0f3f0
									}
									f9elemf7f11f2elemf0.Challenge = f9elemf7f11f2elemf0f3
								}
								if f9elemf7f11f2iter.ActionToUse.Count != nil {
									f9elemf7f11f2elemf0f4 := &svcapitypes.CountAction{}
									if f9elemf7f11f2iter.ActionToUse.Count.CustomRequestHandling != nil {
										f9elemf7f11f2elemf0f4f0 := &svcapitypes.CustomRequestHandling{}
										if f9elemf7f11f2iter.ActionToUse.Count.CustomRequestHandling.InsertHeaders != nil {
											f9elemf7f11f2elemf0f4f0f0 := []*svcapitypes.CustomHTTPHeader{}
											for _, f9elemf7f11f2elemf0f4f0f0iter := range f9elemf7f11f2iter.ActionToUse.Count.CustomRequestHandling.InsertHeaders {
												f9elemf7f11f2elemf0f4f0f0elem := &svcapitypes.CustomHTTPHeader{}
												if f9elemf7f11f2elemf0f4f0f0iter.Name != nil {
													f9elemf7f11f2elemf0f4f0f0elem.Name = f9elemf7f11f2elemf0f4f0f0iter.Name
												}
												if f9elemf7f11f2elemf0f4f0f0iter.Value != nil {
													f9elemf7f11f2elemf0f4f0f0elem.Value = f9elemf7f11f2elemf0f4f0f0iter.Value
												}
												f9elemf7f11f2elemf0f4f0f0 = append(f9elemf7f11f2elemf0f4f0f0, f9elemf7f11f2elemf0f4f0f0elem)
											}
											f9elemf7f11f2elemf0f4f0.InsertHeaders = f9elemf7f11f2elemf0f4f0f0
										}
										f9elemf7f11f2elemf0f4.CustomRequestHandling = f9elemf7f11f2elemf0f4f0
									}
									f9elemf7f11f2elemf0.Count = f9elemf7f11f2elemf0f4
								}
								f9elemf7f11f2elem.ActionToUse = f9elemf7f11f2elemf0
							}
							if f9elemf7f11f2iter.Name != nil {
								f9elemf7f11f2elem.Name = f9elemf7f11f2iter.Name
							}
							f9elemf7f11f2 = append(f9elemf7f11f2, f9elemf7f11f2elem)
						}
						f9elemf7f11.RuleActionOverrides = f9elemf7f11f2
					}
					f9elemf7.RuleGroupReferenceStatement = f9elemf7f11
				}
				if f9iter.Statement.SizeConstraintStatement != nil {
					f9elemf7f12 := &svcapitypes.SizeConstraintStatement{}
					if f9iter.Statement.SizeConstraintStatement.ComparisonOperator != "" {
						f9elemf7f12.ComparisonOperator = aws.String(string(f9iter.Statement.SizeConstraintStatement.ComparisonOperator))
					}
					if f9iter.Statement.SizeConstraintStatement.FieldToMatch != nil {
						f9elemf7f12f1 := &svcapitypes.FieldToMatch{}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.AllQueryArguments != nil {
							f9elemf7f12f1f0 := map[string]*string{}
							f9elemf7f12f1.AllQueryArguments = f9elemf7f12f1f0
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Body != nil {
							f9elemf7f12f1f1 := &svcapitypes.Body{}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Body.OversizeHandling != "" {
								f9elemf7f12f1f1.OversizeHandling = aws.String(string(f9iter.Statement.SizeConstraintStatement.FieldToMatch.Body.OversizeHandling))
							}
							f9elemf7f12f1.Body = f9elemf7f12f1f1
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies != nil {
							f9elemf7f12f1f2 := &svcapitypes.Cookies{}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f9elemf7f12f1f2f0 := &svcapitypes.CookieMatchPattern{}
								if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f9elemf7f12f1f2f0f0 := map[string]*string{}
									f9elemf7f12f1f2f0.All = f9elemf7f12f1f2f0f0
								}
								if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f9elemf7f12f1f2f0.ExcludedCookies = aws.StringSlice(f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f9elemf7f12f1f2f0.IncludedCookies = aws.StringSlice(f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f9elemf7f12f1f2.MatchPattern = f9elemf7f12f1f2f0
							}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchScope != "" {
								f9elemf7f12f1f2.MatchScope = aws.String(string(f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.MatchScope))
							}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.OversizeHandling != "" {
								f9elemf7f12f1f2.OversizeHandling = aws.String(string(f9iter.Statement.SizeConstraintStatement.FieldToMatch.Cookies.OversizeHandling))
							}
							f9elemf7f12f1.Cookies = f9elemf7f12f1f2
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.HeaderOrder != nil {
							f9elemf7f12f1f3 := &svcapitypes.HeaderOrder{}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.HeaderOrder.OversizeHandling != "" {
								f9elemf7f12f1f3.OversizeHandling = aws.String(string(f9iter.Statement.SizeConstraintStatement.FieldToMatch.HeaderOrder.OversizeHandling))
							}
							f9elemf7f12f1.HeaderOrder = f9elemf7f12f1f3
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers != nil {
							f9elemf7f12f1f4 := &svcapitypes.Headers{}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern != nil {
								f9elemf7f12f1f4f0 := &svcapitypes.HeaderMatchPattern{}
								if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f9elemf7f12f1f4f0f0 := map[string]*string{}
									f9elemf7f12f1f4f0.All = f9elemf7f12f1f4f0f0
								}
								if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f9elemf7f12f1f4f0.ExcludedHeaders = aws.StringSlice(f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f9elemf7f12f1f4f0.IncludedHeaders = aws.StringSlice(f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f9elemf7f12f1f4.MatchPattern = f9elemf7f12f1f4f0
							}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchScope != "" {
								f9elemf7f12f1f4.MatchScope = aws.String(string(f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.MatchScope))
							}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.OversizeHandling != "" {
								f9elemf7f12f1f4.OversizeHandling = aws.String(string(f9iter.Statement.SizeConstraintStatement.FieldToMatch.Headers.OversizeHandling))
							}
							f9elemf7f12f1.Headers = f9elemf7f12f1f4
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.JA3Fingerprint != nil {
							f9elemf7f12f1f5 := &svcapitypes.JA3Fingerprint{}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != "" {
								f9elemf7f12f1f5.FallbackBehavior = aws.String(string(f9iter.Statement.SizeConstraintStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior))
							}
							f9elemf7f12f1.JA3Fingerprint = f9elemf7f12f1f5
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody != nil {
							f9elemf7f12f1f6 := &svcapitypes.JSONBody{}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior != "" {
								f9elemf7f12f1f6.InvalidFallbackBehavior = aws.String(string(f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior))
							}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody.MatchPattern != nil {
								f9elemf7f12f1f6f1 := &svcapitypes.JSONMatchPattern{}
								if f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody.MatchPattern.All != nil {
									f9elemf7f12f1f6f1f0 := map[string]*string{}
									f9elemf7f12f1f6f1.All = f9elemf7f12f1f6f1f0
								}
								if f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths != nil {
									f9elemf7f12f1f6f1.IncludedPaths = aws.StringSlice(f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths)
								}
								f9elemf7f12f1f6.MatchPattern = f9elemf7f12f1f6f1
							}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody.MatchScope != "" {
								f9elemf7f12f1f6.MatchScope = aws.String(string(f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody.MatchScope))
							}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody.OversizeHandling != "" {
								f9elemf7f12f1f6.OversizeHandling = aws.String(string(f9iter.Statement.SizeConstraintStatement.FieldToMatch.JsonBody.OversizeHandling))
							}
							f9elemf7f12f1.JSONBody = f9elemf7f12f1f6
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.Method != nil {
							f9elemf7f12f1f7 := map[string]*string{}
							f9elemf7f12f1.Method = f9elemf7f12f1f7
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.QueryString != nil {
							f9elemf7f12f1f8 := map[string]*string{}
							f9elemf7f12f1.QueryString = f9elemf7f12f1f8
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.SingleHeader != nil {
							f9elemf7f12f1f9 := &svcapitypes.SingleHeader{}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.SingleHeader.Name != nil {
								f9elemf7f12f1f9.Name = f9iter.Statement.SizeConstraintStatement.FieldToMatch.SingleHeader.Name
							}
							f9elemf7f12f1.SingleHeader = f9elemf7f12f1f9
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.SingleQueryArgument != nil {
							f9elemf7f12f1f10 := &svcapitypes.SingleQueryArgument{}
							if f9iter.Statement.SizeConstraintStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f9elemf7f12f1f10.Name = f9iter.Statement.SizeConstraintStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f9elemf7f12f1.SingleQueryArgument = f9elemf7f12f1f10
						}
						if f9iter.Statement.SizeConstraintStatement.FieldToMatch.UriPath != nil {
							f9elemf7f12f1f11 := map[string]*string{}
							f9elemf7f12f1.URIPath = f9elemf7f12f1f11
						}
						f9elemf7f12.FieldToMatch = f9elemf7f12f1
					}
					f9elemf7f12.Size = &f9iter.Statement.SizeConstraintStatement.Size
					if f9iter.Statement.SizeConstraintStatement.TextTransformations != nil {
						f9elemf7f12f3 := []*svcapitypes.TextTransformation{}
						for _, f9elemf7f12f3iter := range f9iter.Statement.SizeConstraintStatement.TextTransformations {
							f9elemf7f12f3elem := &svcapitypes.TextTransformation{}
							priorityCopy := int64(f9elemf7f12f3iter.Priority)
							f9elemf7f12f3elem.Priority = &priorityCopy
							if f9elemf7f12f3iter.Type != "" {
								f9elemf7f12f3elem.Type = aws.String(string(f9elemf7f12f3iter.Type))
							}
							f9elemf7f12f3 = append(f9elemf7f12f3, f9elemf7f12f3elem)
						}
						f9elemf7f12.TextTransformations = f9elemf7f12f3
					}
					f9elemf7.SizeConstraintStatement = f9elemf7f12
				}
				if f9iter.Statement.SqliMatchStatement != nil {
					f9elemf7f13 := &svcapitypes.SQLIMatchStatement{}
					if f9iter.Statement.SqliMatchStatement.FieldToMatch != nil {
						f9elemf7f13f0 := &svcapitypes.FieldToMatch{}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.AllQueryArguments != nil {
							f9elemf7f13f0f0 := map[string]*string{}
							f9elemf7f13f0.AllQueryArguments = f9elemf7f13f0f0
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.Body != nil {
							f9elemf7f13f0f1 := &svcapitypes.Body{}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.Body.OversizeHandling != "" {
								f9elemf7f13f0f1.OversizeHandling = aws.String(string(f9iter.Statement.SqliMatchStatement.FieldToMatch.Body.OversizeHandling))
							}
							f9elemf7f13f0.Body = f9elemf7f13f0f1
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies != nil {
							f9elemf7f13f0f2 := &svcapitypes.Cookies{}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f9elemf7f13f0f2f0 := &svcapitypes.CookieMatchPattern{}
								if f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f9elemf7f13f0f2f0f0 := map[string]*string{}
									f9elemf7f13f0f2f0.All = f9elemf7f13f0f2f0f0
								}
								if f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f9elemf7f13f0f2f0.ExcludedCookies = aws.StringSlice(f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f9elemf7f13f0f2f0.IncludedCookies = aws.StringSlice(f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f9elemf7f13f0f2.MatchPattern = f9elemf7f13f0f2f0
							}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies.MatchScope != "" {
								f9elemf7f13f0f2.MatchScope = aws.String(string(f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies.MatchScope))
							}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies.OversizeHandling != "" {
								f9elemf7f13f0f2.OversizeHandling = aws.String(string(f9iter.Statement.SqliMatchStatement.FieldToMatch.Cookies.OversizeHandling))
							}
							f9elemf7f13f0.Cookies = f9elemf7f13f0f2
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.HeaderOrder != nil {
							f9elemf7f13f0f3 := &svcapitypes.HeaderOrder{}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling != "" {
								f9elemf7f13f0f3.OversizeHandling = aws.String(string(f9iter.Statement.SqliMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling))
							}
							f9elemf7f13f0.HeaderOrder = f9elemf7f13f0f3
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers != nil {
							f9elemf7f13f0f4 := &svcapitypes.Headers{}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers.MatchPattern != nil {
								f9elemf7f13f0f4f0 := &svcapitypes.HeaderMatchPattern{}
								if f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f9elemf7f13f0f4f0f0 := map[string]*string{}
									f9elemf7f13f0f4f0.All = f9elemf7f13f0f4f0f0
								}
								if f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f9elemf7f13f0f4f0.ExcludedHeaders = aws.StringSlice(f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f9elemf7f13f0f4f0.IncludedHeaders = aws.StringSlice(f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f9elemf7f13f0f4.MatchPattern = f9elemf7f13f0f4f0
							}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers.MatchScope != "" {
								f9elemf7f13f0f4.MatchScope = aws.String(string(f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers.MatchScope))
							}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers.OversizeHandling != "" {
								f9elemf7f13f0f4.OversizeHandling = aws.String(string(f9iter.Statement.SqliMatchStatement.FieldToMatch.Headers.OversizeHandling))
							}
							f9elemf7f13f0.Headers = f9elemf7f13f0f4
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.JA3Fingerprint != nil {
							f9elemf7f13f0f5 := &svcapitypes.JA3Fingerprint{}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != "" {
								f9elemf7f13f0f5.FallbackBehavior = aws.String(string(f9iter.Statement.SqliMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior))
							}
							f9elemf7f13f0.JA3Fingerprint = f9elemf7f13f0f5
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody != nil {
							f9elemf7f13f0f6 := &svcapitypes.JSONBody{}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior != "" {
								f9elemf7f13f0f6.InvalidFallbackBehavior = aws.String(string(f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior))
							}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody.MatchPattern != nil {
								f9elemf7f13f0f6f1 := &svcapitypes.JSONMatchPattern{}
								if f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody.MatchPattern.All != nil {
									f9elemf7f13f0f6f1f0 := map[string]*string{}
									f9elemf7f13f0f6f1.All = f9elemf7f13f0f6f1f0
								}
								if f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths != nil {
									f9elemf7f13f0f6f1.IncludedPaths = aws.StringSlice(f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths)
								}
								f9elemf7f13f0f6.MatchPattern = f9elemf7f13f0f6f1
							}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody.MatchScope != "" {
								f9elemf7f13f0f6.MatchScope = aws.String(string(f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody.MatchScope))
							}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody.OversizeHandling != "" {
								f9elemf7f13f0f6.OversizeHandling = aws.String(string(f9iter.Statement.SqliMatchStatement.FieldToMatch.JsonBody.OversizeHandling))
							}
							f9elemf7f13f0.JSONBody = f9elemf7f13f0f6
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.Method != nil {
							f9elemf7f13f0f7 := map[string]*string{}
							f9elemf7f13f0.Method = f9elemf7f13f0f7
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.QueryString != nil {
							f9elemf7f13f0f8 := map[string]*string{}
							f9elemf7f13f0.QueryString = f9elemf7f13f0f8
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.SingleHeader != nil {
							f9elemf7f13f0f9 := &svcapitypes.SingleHeader{}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.SingleHeader.Name != nil {
								f9elemf7f13f0f9.Name = f9iter.Statement.SqliMatchStatement.FieldToMatch.SingleHeader.Name
							}
							f9elemf7f13f0.SingleHeader = f9elemf7f13f0f9
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.SingleQueryArgument != nil {
							f9elemf7f13f0f10 := &svcapitypes.SingleQueryArgument{}
							if f9iter.Statement.SqliMatchStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f9elemf7f13f0f10.Name = f9iter.Statement.SqliMatchStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f9elemf7f13f0.SingleQueryArgument = f9elemf7f13f0f10
						}
						if f9iter.Statement.SqliMatchStatement.FieldToMatch.UriPath != nil {
							f9elemf7f13f0f11 := map[string]*string{}
							f9elemf7f13f0.URIPath = f9elemf7f13f0f11
						}
						f9elemf7f13.FieldToMatch = f9elemf7f13f0
					}
					if f9iter.Statement.SqliMatchStatement.SensitivityLevel != "" {
						f9elemf7f13.SensitivityLevel = aws.String(string(f9iter.Statement.SqliMatchStatement.SensitivityLevel))
					}
					if f9iter.Statement.SqliMatchStatement.TextTransformations != nil {
						f9elemf7f13f2 := []*svcapitypes.TextTransformation{}
						for _, f9elemf7f13f2iter := range f9iter.Statement.SqliMatchStatement.TextTransformations {
							f9elemf7f13f2elem := &svcapitypes.TextTransformation{}
							priorityCopy := int64(f9elemf7f13f2iter.Priority)
							f9elemf7f13f2elem.Priority = &priorityCopy
							if f9elemf7f13f2iter.Type != "" {
								f9elemf7f13f2elem.Type = aws.String(string(f9elemf7f13f2iter.Type))
							}
							f9elemf7f13f2 = append(f9elemf7f13f2, f9elemf7f13f2elem)
						}
						f9elemf7f13.TextTransformations = f9elemf7f13f2
					}
					f9elemf7.SQLIMatchStatement = f9elemf7f13
				}
				if f9iter.Statement.XssMatchStatement != nil {
					f9elemf7f14 := &svcapitypes.XSSMatchStatement{}
					if f9iter.Statement.XssMatchStatement.FieldToMatch != nil {
						f9elemf7f14f0 := &svcapitypes.FieldToMatch{}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.AllQueryArguments != nil {
							f9elemf7f14f0f0 := map[string]*string{}
							f9elemf7f14f0.AllQueryArguments = f9elemf7f14f0f0
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.Body != nil {
							f9elemf7f14f0f1 := &svcapitypes.Body{}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.Body.OversizeHandling != "" {
								f9elemf7f14f0f1.OversizeHandling = aws.String(string(f9iter.Statement.XssMatchStatement.FieldToMatch.Body.OversizeHandling))
							}
							f9elemf7f14f0.Body = f9elemf7f14f0f1
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies != nil {
							f9elemf7f14f0f2 := &svcapitypes.Cookies{}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies.MatchPattern != nil {
								f9elemf7f14f0f2f0 := &svcapitypes.CookieMatchPattern{}
								if f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies.MatchPattern.All != nil {
									f9elemf7f14f0f2f0f0 := map[string]*string{}
									f9elemf7f14f0f2f0.All = f9elemf7f14f0f2f0f0
								}
								if f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies != nil {
									f9elemf7f14f0f2f0.ExcludedCookies = aws.StringSlice(f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies.MatchPattern.ExcludedCookies)
								}
								if f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies != nil {
									f9elemf7f14f0f2f0.IncludedCookies = aws.StringSlice(f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies.MatchPattern.IncludedCookies)
								}
								f9elemf7f14f0f2.MatchPattern = f9elemf7f14f0f2f0
							}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies.MatchScope != "" {
								f9elemf7f14f0f2.MatchScope = aws.String(string(f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies.MatchScope))
							}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies.OversizeHandling != "" {
								f9elemf7f14f0f2.OversizeHandling = aws.String(string(f9iter.Statement.XssMatchStatement.FieldToMatch.Cookies.OversizeHandling))
							}
							f9elemf7f14f0.Cookies = f9elemf7f14f0f2
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.HeaderOrder != nil {
							f9elemf7f14f0f3 := &svcapitypes.HeaderOrder{}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling != "" {
								f9elemf7f14f0f3.OversizeHandling = aws.String(string(f9iter.Statement.XssMatchStatement.FieldToMatch.HeaderOrder.OversizeHandling))
							}
							f9elemf7f14f0.HeaderOrder = f9elemf7f14f0f3
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.Headers != nil {
							f9elemf7f14f0f4 := &svcapitypes.Headers{}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.Headers.MatchPattern != nil {
								f9elemf7f14f0f4f0 := &svcapitypes.HeaderMatchPattern{}
								if f9iter.Statement.XssMatchStatement.FieldToMatch.Headers.MatchPattern.All != nil {
									f9elemf7f14f0f4f0f0 := map[string]*string{}
									f9elemf7f14f0f4f0.All = f9elemf7f14f0f4f0f0
								}
								if f9iter.Statement.XssMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders != nil {
									f9elemf7f14f0f4f0.ExcludedHeaders = aws.StringSlice(f9iter.Statement.XssMatchStatement.FieldToMatch.Headers.MatchPattern.ExcludedHeaders)
								}
								if f9iter.Statement.XssMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders != nil {
									f9elemf7f14f0f4f0.IncludedHeaders = aws.StringSlice(f9iter.Statement.XssMatchStatement.FieldToMatch.Headers.MatchPattern.IncludedHeaders)
								}
								f9elemf7f14f0f4.MatchPattern = f9elemf7f14f0f4f0
							}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.Headers.MatchScope != "" {
								f9elemf7f14f0f4.MatchScope = aws.String(string(f9iter.Statement.XssMatchStatement.FieldToMatch.Headers.MatchScope))
							}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.Headers.OversizeHandling != "" {
								f9elemf7f14f0f4.OversizeHandling = aws.String(string(f9iter.Statement.XssMatchStatement.FieldToMatch.Headers.OversizeHandling))
							}
							f9elemf7f14f0.Headers = f9elemf7f14f0f4
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.JA3Fingerprint != nil {
							f9elemf7f14f0f5 := &svcapitypes.JA3Fingerprint{}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior != "" {
								f9elemf7f14f0f5.FallbackBehavior = aws.String(string(f9iter.Statement.XssMatchStatement.FieldToMatch.JA3Fingerprint.FallbackBehavior))
							}
							f9elemf7f14f0.JA3Fingerprint = f9elemf7f14f0f5
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody != nil {
							f9elemf7f14f0f6 := &svcapitypes.JSONBody{}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior != "" {
								f9elemf7f14f0f6.InvalidFallbackBehavior = aws.String(string(f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody.InvalidFallbackBehavior))
							}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody.MatchPattern != nil {
								f9elemf7f14f0f6f1 := &svcapitypes.JSONMatchPattern{}
								if f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody.MatchPattern.All != nil {
									f9elemf7f14f0f6f1f0 := map[string]*string{}
									f9elemf7f14f0f6f1.All = f9elemf7f14f0f6f1f0
								}
								if f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths != nil {
									f9elemf7f14f0f6f1.IncludedPaths = aws.StringSlice(f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody.MatchPattern.IncludedPaths)
								}
								f9elemf7f14f0f6.MatchPattern = f9elemf7f14f0f6f1
							}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody.MatchScope != "" {
								f9elemf7f14f0f6.MatchScope = aws.String(string(f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody.MatchScope))
							}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody.OversizeHandling != "" {
								f9elemf7f14f0f6.OversizeHandling = aws.String(string(f9iter.Statement.XssMatchStatement.FieldToMatch.JsonBody.OversizeHandling))
							}
							f9elemf7f14f0.JSONBody = f9elemf7f14f0f6
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.Method != nil {
							f9elemf7f14f0f7 := map[string]*string{}
							f9elemf7f14f0.Method = f9elemf7f14f0f7
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.QueryString != nil {
							f9elemf7f14f0f8 := map[string]*string{}
							f9elemf7f14f0.QueryString = f9elemf7f14f0f8
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.SingleHeader != nil {
							f9elemf7f14f0f9 := &svcapitypes.SingleHeader{}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.SingleHeader.Name != nil {
								f9elemf7f14f0f9.Name = f9iter.Statement.XssMatchStatement.FieldToMatch.SingleHeader.Name
							}
							f9elemf7f14f0.SingleHeader = f9elemf7f14f0f9
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.SingleQueryArgument != nil {
							f9elemf7f14f0f10 := &svcapitypes.SingleQueryArgument{}
							if f9iter.Statement.XssMatchStatement.FieldToMatch.SingleQueryArgument.Name != nil {
								f9elemf7f14f0f10.Name = f9iter.Statement.XssMatchStatement.FieldToMatch.SingleQueryArgument.Name
							}
							f9elemf7f14f0.SingleQueryArgument = f9elemf7f14f0f10
						}
						if f9iter.Statement.XssMatchStatement.FieldToMatch.UriPath != nil {
							f9elemf7f14f0f11 := map[string]*string{}
							f9elemf7f14f0.URIPath = f9elemf7f14f0f11
						}
						f9elemf7f14.FieldToMatch = f9elemf7f14f0
					}
					if f9iter.Statement.XssMatchStatement.TextTransformations != nil {
						f9elemf7f14f1 := []*svcapitypes.TextTransformation{}
						for _, f9elemf7f14f1iter := range f9iter.Statement.XssMatchStatement.TextTransformations {
							f9elemf7f14f1elem := &svcapitypes.TextTransformation{}
							priorityCopy := int64(f9elemf7f14f1iter.Priority)
							f9elemf7f14f1elem.Priority = &priorityCopy
							if f9elemf7f14f1iter.Type != "" {
								f9elemf7f14f1elem.Type = aws.String(string(f9elemf7f14f1iter.Type))
							}
							f9elemf7f14f1 = append(f9elemf7f14f1, f9elemf7f14f1elem)
						}
						f9elemf7f14.TextTransformations = f9elemf7f14f1
					}
					f9elemf7.XSSMatchStatement = f9elemf7f14
				}
				f9elem.Statement = f9elemf7
			}
			if f9iter.VisibilityConfig != nil {
				f9elemf8 := &svcapitypes.VisibilityConfig{}
				f9elemf8.CloudWatchMetricsEnabled = &f9iter.VisibilityConfig.CloudWatchMetricsEnabled
				if f9iter.VisibilityConfig.MetricName != nil {
					f9elemf8.MetricName = f9iter.VisibilityConfig.MetricName
				}
				f9elemf8.SampledRequestsEnabled = &f9iter.VisibilityConfig.SampledRequestsEnabled
				f9elem.VisibilityConfig = f9elemf8
			}
			f9 = append(f9, f9elem)
		}
		ko.Spec.Rules = f9
	} else {
		ko.Spec.Rules = nil
	}
	if resp.RuleGroup.VisibilityConfig != nil {
		f10 := &svcapitypes.VisibilityConfig{}
		f10.CloudWatchMetricsEnabled = &resp.RuleGroup.VisibilityConfig.CloudWatchMetricsEnabled
		if resp.RuleGroup.VisibilityConfig.MetricName != nil {
			f10.MetricName = resp.RuleGroup.VisibilityConfig.MetricName
		}
		f10.SampledRequestsEnabled = &resp.RuleGroup.VisibilityConfig.SampledRequestsEnabled
		ko.Spec.VisibilityConfig = f10
	} else {
		ko.Spec.VisibilityConfig = nil
	}
`
	assert.Equal(
		expected,
		code.SetResource(crd.Config(), crd, op, "resp", "ko", 1),
	)
}
