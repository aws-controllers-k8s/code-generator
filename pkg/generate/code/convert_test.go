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
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestConvertResource_ECR_Repository_v1alpha1_v1alpha2(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	srcModel := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-v1alpha1.yaml",
	})
	dstModel := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-v1alpha2.yaml",
	})
	hubCRDs, err := srcModel.GetCRDs()
	require.Nil(err)
	spokeCRDs, err := dstModel.GetCRDs()
	require.Nil(err)
	require.Len(hubCRDs, 1)
	require.Len(spokeCRDs, 1)

	hubCRD := hubCRDs[0]
	spokeCRD := spokeCRDs[0]

	expectedConvertTo := `
	dst := dstRaw.(*v1alpha2.Repository)
	if src.Spec.ScanConfig != nil {
		imageScanningConfigurationCopy := &v1alpha2.ImageScanningConfiguration{}
		imageScanningConfigurationCopy.ScanOnPush = src.Spec.ScanConfig.ScanOnPush
		dst.Spec.ImageScanningConfiguration = imageScanningConfigurationCopy
	}

	dst.Spec.ImageTagMutability = src.Spec.ImageTagMutability
	dst.Spec.RepositoryName = src.Spec.Name
	if src.Spec.Tags != nil {
		tagListCopy := make([]*v1alpha2.Tag, 0, len(src.Spec.Tags))
		for i, element := range src.Spec.Tags {
			_ = i // non-used value guard.
			elementCopy := &v1alpha2.Tag{}
			if element != nil {
				tagCopy := &v1alpha2.Tag{}
				tagCopy.Key = element.Key
				tagCopy.Value = element.Value
				elementCopy = tagCopy
			}

			tagListCopy = append(tagListCopy, elementCopy)
		}
		dst.Spec.Tags = tagListCopy
	}

	dst.Status.CreatedAt = src.Status.CreatedAt
	dst.Status.RegistryID = src.Status.RegistryID
	dst.Status.RepositoryURI = src.Status.RepositoryURI
	dst.Status.ACKResourceMetadata = src.Status.ACKResourceMetadata
	dst.Status.Conditions = src.Status.Conditions

	dst.ObjectMeta = src.ObjectMeta
	return nil`

	assert.Equal(
		expectedConvertTo,
		code.Convert(
			spokeCRD, hubCRD, true, "v1alpha2", "src", "dstRaw", 1,
		),
	)

	expectedConvertFrom := `
	src := srcRaw.(*v1alpha2.Repository)
	if src.Spec.ImageScanningConfiguration != nil {
		imageScanningConfigurationCopy := &ImageScanningConfiguration{}
		imageScanningConfigurationCopy.ScanOnPush = src.Spec.ImageScanningConfiguration.ScanOnPush
		dst.Spec.ScanConfig = imageScanningConfigurationCopy
	}

	dst.Spec.ImageTagMutability = src.Spec.ImageTagMutability
	dst.Spec.Name = src.Spec.RepositoryName
	if src.Spec.Tags != nil {
		tagListCopy := make([]*Tag, 0, len(src.Spec.Tags))
		for i, element := range src.Spec.Tags {
			_ = i // non-used value guard.
			elementCopy := &Tag{}
			if element != nil {
				tagCopy := &Tag{}
				tagCopy.Key = element.Key
				tagCopy.Value = element.Value
				elementCopy = tagCopy
			}

			tagListCopy = append(tagListCopy, elementCopy)
		}
		dst.Spec.Tags = tagListCopy
	}

	dst.Status.CreatedAt = src.Status.CreatedAt
	dst.Status.RegistryID = src.Status.RegistryID
	dst.Status.RepositoryURI = src.Status.RepositoryURI
	dst.Status.ACKResourceMetadata = src.Status.ACKResourceMetadata
	dst.Status.Conditions = src.Status.Conditions

	dst.ObjectMeta = src.ObjectMeta
	return nil`

	assert.Equal(
		expectedConvertFrom,
		code.Convert(
			spokeCRD, hubCRD, false, "v1alpha2", "src", "dstRaw", 1,
		),
	)
}

func TestConvertResource_ECR_Repository_v1beta1_v1beta2(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	srcModel := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-v1beta1.yaml",
		ServiceAPIVersion:   "0000-00-01",
	})
	dstModel := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-v1beta2.yaml",
		ServiceAPIVersion:   "0000-00-01",
	})

	hubCRDs, err := srcModel.GetCRDs()
	require.Nil(err)
	spokeCRDs, err := dstModel.GetCRDs()
	require.Nil(err)
	require.Len(hubCRDs, 1)
	require.Len(spokeCRDs, 1)

	hubCRD := hubCRDs[0]
	spokeCRD := spokeCRDs[0]

	expectedConvertTo := `
	dst := dstRaw.(*v1beta2.Repository)
	if src.Spec.EncryptionConfiguration != nil {
		encryptionConfigurationCopy := &v1beta2.EncryptionConfiguration{}
		encryptionConfigurationCopy.EncryptionType = src.Spec.EncryptionConfiguration.EncryptionType
		encryptionConfigurationCopy.KMSKey = src.Spec.EncryptionConfiguration.KMSKey
		dst.Spec.EncryptionConfiguration = encryptionConfigurationCopy
	}

	dst.Spec.ImageTagMutability = src.Spec.ImageTagMutability
	dst.Spec.Name = src.Spec.Name
	if src.Spec.ScanConfig != nil {
		imageScanningConfigurationCopy := &v1beta2.ImageScanningConfiguration{}
		imageScanningConfigurationCopy.ScanOnPush = src.Spec.ScanConfig.ScanOnPush
		dst.Spec.ScanConfig = imageScanningConfigurationCopy
	}

	if src.Spec.Tags != nil {
		tagListCopy := make([]*v1beta2.Tag, 0, len(src.Spec.Tags))
		for i, element := range src.Spec.Tags {
			_ = i // non-used value guard.
			elementCopy := &v1beta2.Tag{}
			if element != nil {
				tagCopy := &v1beta2.Tag{}
				tagCopy.Key = element.Key
				tagCopy.Value = element.Value
				elementCopy = tagCopy
			}

			tagListCopy = append(tagListCopy, elementCopy)
		}
		dst.Spec.Tags = tagListCopy
	}

	dst.Status.CreatedAt = src.Status.CreatedAt
	dst.Status.RegistryID = src.Status.RegistryID
	dst.Status.RepositoryURI = src.Status.RepositoryURI
	dst.Status.ACKResourceMetadata = src.Status.ACKResourceMetadata
	dst.Status.Conditions = src.Status.Conditions

	dst.ObjectMeta = src.ObjectMeta
	return nil`

	assert.Equal(
		expectedConvertTo,
		code.Convert(
			spokeCRD, hubCRD, true, "v1beta2", "src", "dstRaw", 1,
		),
	)

	expectedConvertFrom := `
	src := srcRaw.(*v1beta2.Repository)
	if src.Spec.EncryptionConfiguration != nil {
		encryptionConfigurationCopy := &EncryptionConfiguration{}
		encryptionConfigurationCopy.EncryptionType = src.Spec.EncryptionConfiguration.EncryptionType
		encryptionConfigurationCopy.KMSKey = src.Spec.EncryptionConfiguration.KMSKey
		dst.Spec.EncryptionConfiguration = encryptionConfigurationCopy
	}

	dst.Spec.ImageTagMutability = src.Spec.ImageTagMutability
	dst.Spec.Name = src.Spec.Name
	if src.Spec.ScanConfig != nil {
		imageScanningConfigurationCopy := &ImageScanningConfiguration{}
		imageScanningConfigurationCopy.ScanOnPush = src.Spec.ScanConfig.ScanOnPush
		dst.Spec.ScanConfig = imageScanningConfigurationCopy
	}

	if src.Spec.Tags != nil {
		tagListCopy := make([]*Tag, 0, len(src.Spec.Tags))
		for i, element := range src.Spec.Tags {
			_ = i // non-used value guard.
			elementCopy := &Tag{}
			if element != nil {
				tagCopy := &Tag{}
				tagCopy.Key = element.Key
				tagCopy.Value = element.Value
				elementCopy = tagCopy
			}

			tagListCopy = append(tagListCopy, elementCopy)
		}
		dst.Spec.Tags = tagListCopy
	}

	dst.Status.CreatedAt = src.Status.CreatedAt
	dst.Status.RegistryID = src.Status.RegistryID
	dst.Status.RepositoryURI = src.Status.RepositoryURI
	dst.Status.ACKResourceMetadata = src.Status.ACKResourceMetadata
	dst.Status.Conditions = src.Status.Conditions

	dst.ObjectMeta = src.ObjectMeta
	return nil`

	assert.Equal(
		expectedConvertFrom,
		code.Convert(
			spokeCRD, hubCRD, false, "v1beta2", "src", "dstRaw", 1,
		),
	)
}

func TestConvertResource_APIGateway_Route_v1alpha1_v1alpha1(t *testing.T) {
	assert := assert.New(t)

	srcModel := testutil.NewModelForServiceWithOptions(t, "apigatewayv2", &testutil.TestingModelOptions{})
	crd := testutil.GetCRDByName(t, srcModel, "Route")

	expectedConvertTo := `
	dst := dstRaw.(*v1alpha2.Route)
	dst.Spec.APIID = src.Spec.APIID
	dst.Spec.APIKeyRequired = src.Spec.APIKeyRequired
	src.Spec.AuthorizationScopes = dst.Spec.AuthorizationScopes
	dst.Spec.AuthorizationType = src.Spec.AuthorizationType
	dst.Spec.AuthorizerID = src.Spec.AuthorizerID
	dst.Spec.ModelSelectionExpression = src.Spec.ModelSelectionExpression
	dst.Spec.OperationName = src.Spec.OperationName
	dst.Spec.RequestModels = src.Spec.RequestModels
	if src.Spec.RequestParameters != nil {
		routeParametersCopy := make(map[string]*v1alpha2.ParameterConstraints, len(src.Spec.RequestParameters))
		for k, v := range src.Spec.RequestParameters {
			elementCopy := &v1alpha2.ParameterConstraints{}
			if v != nil {
				parameterConstraintsCopy := &v1alpha2.ParameterConstraints{}
				parameterConstraintsCopy.Required = v.Required
				elementCopy = parameterConstraintsCopy
			}

			routeParametersCopy[k] = elementCopy
		}
		dst.Spec.RequestParameters = routeParametersCopy
	}

	dst.Spec.RouteKey = src.Spec.RouteKey
	dst.Spec.RouteResponseSelectionExpression = src.Spec.RouteResponseSelectionExpression
	dst.Spec.Target = src.Spec.Target
	dst.Status.APIGatewayManaged = src.Status.APIGatewayManaged
	dst.Status.RouteID = src.Status.RouteID
	dst.Status.ACKResourceMetadata = src.Status.ACKResourceMetadata
	dst.Status.Conditions = src.Status.Conditions

	dst.ObjectMeta = src.ObjectMeta
	return nil`

	assert.Equal(
		expectedConvertTo,
		code.Convert(
			crd, crd, true, "v1alpha2", "src", "dstRaw", 1,
		),
	)

	expectedConvertFrom := `
	src := srcRaw.(*v1alpha2.Route)
	dst.Spec.APIID = src.Spec.APIID
	dst.Spec.APIKeyRequired = src.Spec.APIKeyRequired
	src.Spec.AuthorizationScopes = dst.Spec.AuthorizationScopes
	dst.Spec.AuthorizationType = src.Spec.AuthorizationType
	dst.Spec.AuthorizerID = src.Spec.AuthorizerID
	dst.Spec.ModelSelectionExpression = src.Spec.ModelSelectionExpression
	dst.Spec.OperationName = src.Spec.OperationName
	dst.Spec.RequestModels = src.Spec.RequestModels
	if src.Spec.RequestParameters != nil {
		routeParametersCopy := make(map[string]*ParameterConstraints, len(src.Spec.RequestParameters))
		for k, v := range src.Spec.RequestParameters {
			elementCopy := &ParameterConstraints{}
			if v != nil {
				parameterConstraintsCopy := &ParameterConstraints{}
				parameterConstraintsCopy.Required = v.Required
				elementCopy = parameterConstraintsCopy
			}

			routeParametersCopy[k] = elementCopy
		}
		dst.Spec.RequestParameters = routeParametersCopy
	}

	dst.Spec.RouteKey = src.Spec.RouteKey
	dst.Spec.RouteResponseSelectionExpression = src.Spec.RouteResponseSelectionExpression
	dst.Spec.Target = src.Spec.Target
	dst.Status.APIGatewayManaged = src.Status.APIGatewayManaged
	dst.Status.RouteID = src.Status.RouteID
	dst.Status.ACKResourceMetadata = src.Status.ACKResourceMetadata
	dst.Status.Conditions = src.Status.Conditions

	dst.ObjectMeta = src.ObjectMeta
	return nil`

	assert.Equal(
		expectedConvertFrom,
		code.Convert(
			crd, crd, false, "v1alpha2", "src", "dstRaw", 1,
		),
	)
}

func TestConvertResource_CodeDeploy_Deployment_v1alpha1_v1alpha1(t *testing.T) {
	assert := assert.New(t)

	srcModel := testutil.NewModelForServiceWithOptions(t, "codedeploy", &testutil.TestingModelOptions{})
	crd := testutil.GetCRDByName(t, srcModel, "Deployment")

	expectedConvertTo := `
	dst := dstRaw.(*v1alpha2.Deployment)
	dst.Spec.ApplicationName = src.Spec.ApplicationName
	if src.Spec.AutoRollbackConfiguration != nil {
		autoRollbackConfigurationCopy := &v1alpha2.AutoRollbackConfiguration{}
		autoRollbackConfigurationCopy.Enabled = src.Spec.AutoRollbackConfiguration.Enabled
		src.Spec.AutoRollbackConfiguration.Events = autoRollbackConfigurationCopy.Events
		dst.Spec.AutoRollbackConfiguration = autoRollbackConfigurationCopy
	}

	dst.Spec.DeploymentConfigName = src.Spec.DeploymentConfigName
	dst.Spec.DeploymentGroupName = src.Spec.DeploymentGroupName
	dst.Spec.Description = src.Spec.Description
	dst.Spec.FileExistsBehavior = src.Spec.FileExistsBehavior
	dst.Spec.IgnoreApplicationStopFailures = src.Spec.IgnoreApplicationStopFailures
	if src.Spec.Revision != nil {
		revisionLocationCopy := &v1alpha2.RevisionLocation{}
		if src.Spec.Revision.AppSpecContent != nil {
			appSpecContentCopy := &v1alpha2.AppSpecContent{}
			appSpecContentCopy.Content = src.Spec.Revision.AppSpecContent.Content
			appSpecContentCopy.SHA256 = src.Spec.Revision.AppSpecContent.SHA256
			revisionLocationCopy.AppSpecContent = appSpecContentCopy
		}

		if src.Spec.Revision.GitHubLocation != nil {
			gitHubLocationCopy := &v1alpha2.GitHubLocation{}
			gitHubLocationCopy.CommitID = src.Spec.Revision.GitHubLocation.CommitID
			gitHubLocationCopy.Repository = src.Spec.Revision.GitHubLocation.Repository
			revisionLocationCopy.GitHubLocation = gitHubLocationCopy
		}

		revisionLocationCopy.RevisionType = src.Spec.Revision.RevisionType
		if src.Spec.Revision.S3Location != nil {
			s3LocationCopy := &v1alpha2.S3Location{}
			s3LocationCopy.Bucket = src.Spec.Revision.S3Location.Bucket
			s3LocationCopy.BundleType = src.Spec.Revision.S3Location.BundleType
			s3LocationCopy.ETag = src.Spec.Revision.S3Location.ETag
			s3LocationCopy.Key = src.Spec.Revision.S3Location.Key
			s3LocationCopy.Version = src.Spec.Revision.S3Location.Version
			revisionLocationCopy.S3Location = s3LocationCopy
		}

		if src.Spec.Revision.String != nil {
			rawStringCopy := &v1alpha2.RawString{}
			rawStringCopy.Content = src.Spec.Revision.String.Content
			rawStringCopy.SHA256 = src.Spec.Revision.String.SHA256
			revisionLocationCopy.String = rawStringCopy
		}

		dst.Spec.Revision = revisionLocationCopy
	}

	if src.Spec.TargetInstances != nil {
		targetInstancesCopy := &v1alpha2.TargetInstances{}
		src.Spec.TargetInstances.AutoScalingGroups = targetInstancesCopy.AutoScalingGroups
		if src.Spec.TargetInstances.EC2TagSet != nil {
			ec2TagSetCopy := &v1alpha2.EC2TagSet{}
			if src.Spec.TargetInstances.EC2TagSet.EC2TagSetList != nil {
				ec2TagSetListCopy := make([][]*v1alpha2.EC2TagFilter, 0, len(src.Spec.TargetInstances.EC2TagSet.EC2TagSetList))
				for i, element := range src.Spec.TargetInstances.EC2TagSet.EC2TagSetList {
					_ = i // non-used value guard.
					elementCopy := make([]*v1alpha2.EC2TagFilter, 0, len(element))
					if element != nil {
						ec2TagFilterListCopy := make([]*v1alpha2.EC2TagFilter, 0, len(element))
						for i, element := range element {
							_ = i // non-used value guard.
							elementCopy := &v1alpha2.EC2TagFilter{}
							if element != nil {
								ec2TagFilterCopy := &v1alpha2.EC2TagFilter{}
								ec2TagFilterCopy.Key = element.Key
								ec2TagFilterCopy.Type = element.Type
								ec2TagFilterCopy.Value = element.Value
								elementCopy = ec2TagFilterCopy
							}

							ec2TagFilterListCopy = append(ec2TagFilterListCopy, elementCopy)
						}
						elementCopy = ec2TagFilterListCopy
					}

					ec2TagSetListCopy = append(ec2TagSetListCopy, elementCopy)
				}
				ec2TagSetCopy.EC2TagSetList = ec2TagSetListCopy
			}

			targetInstancesCopy.EC2TagSet = ec2TagSetCopy
		}

		if src.Spec.TargetInstances.TagFilters != nil {
			ec2TagFilterListCopy := make([]*v1alpha2.EC2TagFilter, 0, len(src.Spec.TargetInstances.TagFilters))
			for i, element := range src.Spec.TargetInstances.TagFilters {
				_ = i // non-used value guard.
				elementCopy := &v1alpha2.EC2TagFilter{}
				if element != nil {
					ec2TagFilterCopy := &v1alpha2.EC2TagFilter{}
					ec2TagFilterCopy.Key = element.Key
					ec2TagFilterCopy.Type = element.Type
					ec2TagFilterCopy.Value = element.Value
					elementCopy = ec2TagFilterCopy
				}

				ec2TagFilterListCopy = append(ec2TagFilterListCopy, elementCopy)
			}
			targetInstancesCopy.TagFilters = ec2TagFilterListCopy
		}

		dst.Spec.TargetInstances = targetInstancesCopy
	}

	dst.Spec.UpdateOutdatedInstancesOnly = src.Spec.UpdateOutdatedInstancesOnly
	dst.Status.DeploymentID = src.Status.DeploymentID
	dst.Status.ACKResourceMetadata = src.Status.ACKResourceMetadata
	dst.Status.Conditions = src.Status.Conditions

	dst.ObjectMeta = src.ObjectMeta
	return nil`

	assert.Equal(
		expectedConvertTo,
		code.Convert(
			crd, crd, true, "v1alpha2", "src", "dstRaw", 1,
		),
	)

	expectedConvertFrom := `
	src := srcRaw.(*v1alpha2.Deployment)
	dst.Spec.ApplicationName = src.Spec.ApplicationName
	if src.Spec.AutoRollbackConfiguration != nil {
		autoRollbackConfigurationCopy := &AutoRollbackConfiguration{}
		autoRollbackConfigurationCopy.Enabled = src.Spec.AutoRollbackConfiguration.Enabled
		src.Spec.AutoRollbackConfiguration.Events = autoRollbackConfigurationCopy.Events
		dst.Spec.AutoRollbackConfiguration = autoRollbackConfigurationCopy
	}

	dst.Spec.DeploymentConfigName = src.Spec.DeploymentConfigName
	dst.Spec.DeploymentGroupName = src.Spec.DeploymentGroupName
	dst.Spec.Description = src.Spec.Description
	dst.Spec.FileExistsBehavior = src.Spec.FileExistsBehavior
	dst.Spec.IgnoreApplicationStopFailures = src.Spec.IgnoreApplicationStopFailures
	if src.Spec.Revision != nil {
		revisionLocationCopy := &RevisionLocation{}
		if src.Spec.Revision.AppSpecContent != nil {
			appSpecContentCopy := &AppSpecContent{}
			appSpecContentCopy.Content = src.Spec.Revision.AppSpecContent.Content
			appSpecContentCopy.SHA256 = src.Spec.Revision.AppSpecContent.SHA256
			revisionLocationCopy.AppSpecContent = appSpecContentCopy
		}

		if src.Spec.Revision.GitHubLocation != nil {
			gitHubLocationCopy := &GitHubLocation{}
			gitHubLocationCopy.CommitID = src.Spec.Revision.GitHubLocation.CommitID
			gitHubLocationCopy.Repository = src.Spec.Revision.GitHubLocation.Repository
			revisionLocationCopy.GitHubLocation = gitHubLocationCopy
		}

		revisionLocationCopy.RevisionType = src.Spec.Revision.RevisionType
		if src.Spec.Revision.S3Location != nil {
			s3LocationCopy := &S3Location{}
			s3LocationCopy.Bucket = src.Spec.Revision.S3Location.Bucket
			s3LocationCopy.BundleType = src.Spec.Revision.S3Location.BundleType
			s3LocationCopy.ETag = src.Spec.Revision.S3Location.ETag
			s3LocationCopy.Key = src.Spec.Revision.S3Location.Key
			s3LocationCopy.Version = src.Spec.Revision.S3Location.Version
			revisionLocationCopy.S3Location = s3LocationCopy
		}

		if src.Spec.Revision.String != nil {
			rawStringCopy := &RawString{}
			rawStringCopy.Content = src.Spec.Revision.String.Content
			rawStringCopy.SHA256 = src.Spec.Revision.String.SHA256
			revisionLocationCopy.String = rawStringCopy
		}

		dst.Spec.Revision = revisionLocationCopy
	}

	if src.Spec.TargetInstances != nil {
		targetInstancesCopy := &TargetInstances{}
		src.Spec.TargetInstances.AutoScalingGroups = targetInstancesCopy.AutoScalingGroups
		if src.Spec.TargetInstances.EC2TagSet != nil {
			ec2TagSetCopy := &EC2TagSet{}
			if src.Spec.TargetInstances.EC2TagSet.EC2TagSetList != nil {
				ec2TagSetListCopy := make([][]*EC2TagFilter, 0, len(src.Spec.TargetInstances.EC2TagSet.EC2TagSetList))
				for i, element := range src.Spec.TargetInstances.EC2TagSet.EC2TagSetList {
					_ = i // non-used value guard.
					elementCopy := make([]*EC2TagFilter, 0, len(element))
					if element != nil {
						ec2TagFilterListCopy := make([]*EC2TagFilter, 0, len(element))
						for i, element := range element {
							_ = i // non-used value guard.
							elementCopy := &EC2TagFilter{}
							if element != nil {
								ec2TagFilterCopy := &EC2TagFilter{}
								ec2TagFilterCopy.Key = element.Key
								ec2TagFilterCopy.Type = element.Type
								ec2TagFilterCopy.Value = element.Value
								elementCopy = ec2TagFilterCopy
							}

							ec2TagFilterListCopy = append(ec2TagFilterListCopy, elementCopy)
						}
						elementCopy = ec2TagFilterListCopy
					}

					ec2TagSetListCopy = append(ec2TagSetListCopy, elementCopy)
				}
				ec2TagSetCopy.EC2TagSetList = ec2TagSetListCopy
			}

			targetInstancesCopy.EC2TagSet = ec2TagSetCopy
		}

		if src.Spec.TargetInstances.TagFilters != nil {
			ec2TagFilterListCopy := make([]*EC2TagFilter, 0, len(src.Spec.TargetInstances.TagFilters))
			for i, element := range src.Spec.TargetInstances.TagFilters {
				_ = i // non-used value guard.
				elementCopy := &EC2TagFilter{}
				if element != nil {
					ec2TagFilterCopy := &EC2TagFilter{}
					ec2TagFilterCopy.Key = element.Key
					ec2TagFilterCopy.Type = element.Type
					ec2TagFilterCopy.Value = element.Value
					elementCopy = ec2TagFilterCopy
				}

				ec2TagFilterListCopy = append(ec2TagFilterListCopy, elementCopy)
			}
			targetInstancesCopy.TagFilters = ec2TagFilterListCopy
		}

		dst.Spec.TargetInstances = targetInstancesCopy
	}

	dst.Spec.UpdateOutdatedInstancesOnly = src.Spec.UpdateOutdatedInstancesOnly
	dst.Status.DeploymentID = src.Status.DeploymentID
	dst.Status.ACKResourceMetadata = src.Status.ACKResourceMetadata
	dst.Status.Conditions = src.Status.Conditions

	dst.ObjectMeta = src.ObjectMeta
	return nil`

	assert.Equal(
		expectedConvertFrom,
		code.Convert(
			crd, crd, false, "v1alpha2", "src", "dstRaw", 1,
		),
	)
}
