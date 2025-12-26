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

func TestCompareResource_S3_Bucket(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "s3")

	crd := testutil.GetCRDByName(t, g, "Bucket")
	require.NotNil(crd)

	// The ACL field is ignored in the S3 generator config and therefore should
	// not appear in this output.
	expected := `
	if ackcompare.HasNilDifference(a.ko.Spec.CreateBucketConfiguration, b.ko.Spec.CreateBucketConfiguration) {
		delta.Add("Spec.CreateBucketConfiguration", a.ko.Spec.CreateBucketConfiguration, b.ko.Spec.CreateBucketConfiguration)
	} else if a.ko.Spec.CreateBucketConfiguration != nil && b.ko.Spec.CreateBucketConfiguration != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.CreateBucketConfiguration.LocationConstraint, b.ko.Spec.CreateBucketConfiguration.LocationConstraint) {
			delta.Add("Spec.CreateBucketConfiguration.LocationConstraint", a.ko.Spec.CreateBucketConfiguration.LocationConstraint, b.ko.Spec.CreateBucketConfiguration.LocationConstraint)
		} else if a.ko.Spec.CreateBucketConfiguration.LocationConstraint != nil && b.ko.Spec.CreateBucketConfiguration.LocationConstraint != nil {
			if *a.ko.Spec.CreateBucketConfiguration.LocationConstraint != *b.ko.Spec.CreateBucketConfiguration.LocationConstraint {
				delta.Add("Spec.CreateBucketConfiguration.LocationConstraint", a.ko.Spec.CreateBucketConfiguration.LocationConstraint, b.ko.Spec.CreateBucketConfiguration.LocationConstraint)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.GrantFullControl, b.ko.Spec.GrantFullControl) {
		delta.Add("Spec.GrantFullControl", a.ko.Spec.GrantFullControl, b.ko.Spec.GrantFullControl)
	} else if a.ko.Spec.GrantFullControl != nil && b.ko.Spec.GrantFullControl != nil {
		if *a.ko.Spec.GrantFullControl != *b.ko.Spec.GrantFullControl {
			delta.Add("Spec.GrantFullControl", a.ko.Spec.GrantFullControl, b.ko.Spec.GrantFullControl)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.GrantRead, b.ko.Spec.GrantRead) {
		delta.Add("Spec.GrantRead", a.ko.Spec.GrantRead, b.ko.Spec.GrantRead)
	} else if a.ko.Spec.GrantRead != nil && b.ko.Spec.GrantRead != nil {
		if *a.ko.Spec.GrantRead != *b.ko.Spec.GrantRead {
			delta.Add("Spec.GrantRead", a.ko.Spec.GrantRead, b.ko.Spec.GrantRead)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.GrantReadACP, b.ko.Spec.GrantReadACP) {
		delta.Add("Spec.GrantReadACP", a.ko.Spec.GrantReadACP, b.ko.Spec.GrantReadACP)
	} else if a.ko.Spec.GrantReadACP != nil && b.ko.Spec.GrantReadACP != nil {
		if *a.ko.Spec.GrantReadACP != *b.ko.Spec.GrantReadACP {
			delta.Add("Spec.GrantReadACP", a.ko.Spec.GrantReadACP, b.ko.Spec.GrantReadACP)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.GrantWrite, b.ko.Spec.GrantWrite) {
		delta.Add("Spec.GrantWrite", a.ko.Spec.GrantWrite, b.ko.Spec.GrantWrite)
	} else if a.ko.Spec.GrantWrite != nil && b.ko.Spec.GrantWrite != nil {
		if *a.ko.Spec.GrantWrite != *b.ko.Spec.GrantWrite {
			delta.Add("Spec.GrantWrite", a.ko.Spec.GrantWrite, b.ko.Spec.GrantWrite)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.GrantWriteACP, b.ko.Spec.GrantWriteACP) {
		delta.Add("Spec.GrantWriteACP", a.ko.Spec.GrantWriteACP, b.ko.Spec.GrantWriteACP)
	} else if a.ko.Spec.GrantWriteACP != nil && b.ko.Spec.GrantWriteACP != nil {
		if *a.ko.Spec.GrantWriteACP != *b.ko.Spec.GrantWriteACP {
			delta.Add("Spec.GrantWriteACP", a.ko.Spec.GrantWriteACP, b.ko.Spec.GrantWriteACP)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Logging, b.ko.Spec.Logging) {
		delta.Add("Spec.Logging", a.ko.Spec.Logging, b.ko.Spec.Logging)
	} else if a.ko.Spec.Logging != nil && b.ko.Spec.Logging != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.Logging.LoggingEnabled, b.ko.Spec.Logging.LoggingEnabled) {
			delta.Add("Spec.Logging.LoggingEnabled", a.ko.Spec.Logging.LoggingEnabled, b.ko.Spec.Logging.LoggingEnabled)
		} else if a.ko.Spec.Logging.LoggingEnabled != nil && b.ko.Spec.Logging.LoggingEnabled != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.Logging.LoggingEnabled.TargetBucket, b.ko.Spec.Logging.LoggingEnabled.TargetBucket) {
				delta.Add("Spec.Logging.LoggingEnabled.TargetBucket", a.ko.Spec.Logging.LoggingEnabled.TargetBucket, b.ko.Spec.Logging.LoggingEnabled.TargetBucket)
			} else if a.ko.Spec.Logging.LoggingEnabled.TargetBucket != nil && b.ko.Spec.Logging.LoggingEnabled.TargetBucket != nil {
				if *a.ko.Spec.Logging.LoggingEnabled.TargetBucket != *b.ko.Spec.Logging.LoggingEnabled.TargetBucket {
					delta.Add("Spec.Logging.LoggingEnabled.TargetBucket", a.ko.Spec.Logging.LoggingEnabled.TargetBucket, b.ko.Spec.Logging.LoggingEnabled.TargetBucket)
				}
			}
			if len(a.ko.Spec.Logging.LoggingEnabled.TargetGrants) != len(b.ko.Spec.Logging.LoggingEnabled.TargetGrants) {
				delta.Add("Spec.Logging.LoggingEnabled.TargetGrants", a.ko.Spec.Logging.LoggingEnabled.TargetGrants, b.ko.Spec.Logging.LoggingEnabled.TargetGrants)
			} else if len(a.ko.Spec.Logging.LoggingEnabled.TargetGrants) > 0 {
				if !equality.Semantic.Equalities.DeepEqual(a.ko.Spec.Logging.LoggingEnabled.TargetGrants, b.ko.Spec.Logging.LoggingEnabled.TargetGrants) {
					delta.Add("Spec.Logging.LoggingEnabled.TargetGrants", a.ko.Spec.Logging.LoggingEnabled.TargetGrants, b.ko.Spec.Logging.LoggingEnabled.TargetGrants)
				}
			}
			if ackcompare.HasNilDifference(a.ko.Spec.Logging.LoggingEnabled.TargetPrefix, b.ko.Spec.Logging.LoggingEnabled.TargetPrefix) {
				delta.Add("Spec.Logging.LoggingEnabled.TargetPrefix", a.ko.Spec.Logging.LoggingEnabled.TargetPrefix, b.ko.Spec.Logging.LoggingEnabled.TargetPrefix)
			} else if a.ko.Spec.Logging.LoggingEnabled.TargetPrefix != nil && b.ko.Spec.Logging.LoggingEnabled.TargetPrefix != nil {
				if *a.ko.Spec.Logging.LoggingEnabled.TargetPrefix != *b.ko.Spec.Logging.LoggingEnabled.TargetPrefix {
					delta.Add("Spec.Logging.LoggingEnabled.TargetPrefix", a.ko.Spec.Logging.LoggingEnabled.TargetPrefix, b.ko.Spec.Logging.LoggingEnabled.TargetPrefix)
				}
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name) {
		delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	} else if a.ko.Spec.Name != nil && b.ko.Spec.Name != nil {
		if *a.ko.Spec.Name != *b.ko.Spec.Name {
			delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ObjectLockEnabledForBucket, b.ko.Spec.ObjectLockEnabledForBucket) {
		delta.Add("Spec.ObjectLockEnabledForBucket", a.ko.Spec.ObjectLockEnabledForBucket, b.ko.Spec.ObjectLockEnabledForBucket)
	} else if a.ko.Spec.ObjectLockEnabledForBucket != nil && b.ko.Spec.ObjectLockEnabledForBucket != nil {
		if *a.ko.Spec.ObjectLockEnabledForBucket != *b.ko.Spec.ObjectLockEnabledForBucket {
			delta.Add("Spec.ObjectLockEnabledForBucket", a.ko.Spec.ObjectLockEnabledForBucket, b.ko.Spec.ObjectLockEnabledForBucket)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Tagging, b.ko.Spec.Tagging) {
		delta.Add("Spec.Tagging", a.ko.Spec.Tagging, b.ko.Spec.Tagging)
	} else if a.ko.Spec.Tagging != nil && b.ko.Spec.Tagging != nil {
		desiredACKTags, _ := convertToOrderedACKTags(a.ko.Spec.Tagging.TagSet)
		latestACKTags, _ := convertToOrderedACKTags(b.ko.Spec.Tagging.TagSet)
		if !ackcompare.MapStringStringEqual(desiredACKTags, latestACKTags) {
			delta.Add("Spec.Tagging", a.ko.Spec.Tagging.TagSet, b.ko.Spec.Tagging.TagSet)
		}
	}
`
	assert.Equal(
		expected,
		code.CompareResource(
			crd.Config(), crd, "delta", "a.ko", "b.ko", 1,
		),
	)
}

func TestCompareResource_Lambda_CodeSigningConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "lambda")

	crd := testutil.GetCRDByName(t, g, "CodeSigningConfig")
	require.NotNil(crd)

	expected := `
	if ackcompare.HasNilDifference(a.ko.Spec.AllowedPublishers, b.ko.Spec.AllowedPublishers) {
		delta.Add("Spec.AllowedPublishers", a.ko.Spec.AllowedPublishers, b.ko.Spec.AllowedPublishers)
	} else if a.ko.Spec.AllowedPublishers != nil && b.ko.Spec.AllowedPublishers != nil {
		if len(a.ko.Spec.AllowedPublishers.SigningProfileVersionARNs) != len(b.ko.Spec.AllowedPublishers.SigningProfileVersionARNs) {
			delta.Add("Spec.AllowedPublishers.SigningProfileVersionARNs", a.ko.Spec.AllowedPublishers.SigningProfileVersionARNs, b.ko.Spec.AllowedPublishers.SigningProfileVersionARNs)
		} else if len(a.ko.Spec.AllowedPublishers.SigningProfileVersionARNs) > 0 {
			if !ackcompare.SliceStringPEqual(a.ko.Spec.AllowedPublishers.SigningProfileVersionARNs, b.ko.Spec.AllowedPublishers.SigningProfileVersionARNs) {
				delta.Add("Spec.AllowedPublishers.SigningProfileVersionARNs", a.ko.Spec.AllowedPublishers.SigningProfileVersionARNs, b.ko.Spec.AllowedPublishers.SigningProfileVersionARNs)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.CodeSigningPolicies, b.ko.Spec.CodeSigningPolicies) {
		delta.Add("Spec.CodeSigningPolicies", a.ko.Spec.CodeSigningPolicies, b.ko.Spec.CodeSigningPolicies)
	} else if a.ko.Spec.CodeSigningPolicies != nil && b.ko.Spec.CodeSigningPolicies != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment, b.ko.Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment) {
			delta.Add("Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment", a.ko.Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment, b.ko.Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment)
		} else if a.ko.Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment != nil && b.ko.Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment != nil {
			if *a.ko.Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment != *b.ko.Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment {
				delta.Add("Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment", a.ko.Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment, b.ko.Spec.CodeSigningPolicies.UntrustedArtifactOnDeployment)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Description, b.ko.Spec.Description) {
		delta.Add("Spec.Description", a.ko.Spec.Description, b.ko.Spec.Description)
	} else if a.ko.Spec.Description != nil && b.ko.Spec.Description != nil {
		if *a.ko.Spec.Description != *b.ko.Spec.Description {
			delta.Add("Spec.Description", a.ko.Spec.Description, b.ko.Spec.Description)
		}
	}
`
	assert.Equal(
		expected,
		code.CompareResource(
			crd.Config(), crd, "delta", "a.ko", "b.ko", 1,
		),
	)
}

func TestCompareResource_Lambda_Function(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "lambda")

	crd := testutil.GetCRDByName(t, g, "Function")
	require.NotNil(crd)

	expected := `
	if ackcompare.HasNilDifference(a.ko.Spec.Code, b.ko.Spec.Code) {
		delta.Add("Spec.Code", a.ko.Spec.Code, b.ko.Spec.Code)
	} else if a.ko.Spec.Code != nil && b.ko.Spec.Code != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.Code.ImageURI, b.ko.Spec.Code.ImageURI) {
			delta.Add("Spec.Code.ImageURI", a.ko.Spec.Code.ImageURI, b.ko.Spec.Code.ImageURI)
		} else if a.ko.Spec.Code.ImageURI != nil && b.ko.Spec.Code.ImageURI != nil {
			if *a.ko.Spec.Code.ImageURI != *b.ko.Spec.Code.ImageURI {
				delta.Add("Spec.Code.ImageURI", a.ko.Spec.Code.ImageURI, b.ko.Spec.Code.ImageURI)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.Code.S3Bucket, b.ko.Spec.Code.S3Bucket) {
			delta.Add("Spec.Code.S3Bucket", a.ko.Spec.Code.S3Bucket, b.ko.Spec.Code.S3Bucket)
		} else if a.ko.Spec.Code.S3Bucket != nil && b.ko.Spec.Code.S3Bucket != nil {
			if *a.ko.Spec.Code.S3Bucket != *b.ko.Spec.Code.S3Bucket {
				delta.Add("Spec.Code.S3Bucket", a.ko.Spec.Code.S3Bucket, b.ko.Spec.Code.S3Bucket)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.Code.S3Key, b.ko.Spec.Code.S3Key) {
			delta.Add("Spec.Code.S3Key", a.ko.Spec.Code.S3Key, b.ko.Spec.Code.S3Key)
		} else if a.ko.Spec.Code.S3Key != nil && b.ko.Spec.Code.S3Key != nil {
			if *a.ko.Spec.Code.S3Key != *b.ko.Spec.Code.S3Key {
				delta.Add("Spec.Code.S3Key", a.ko.Spec.Code.S3Key, b.ko.Spec.Code.S3Key)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.Code.S3ObjectVersion, b.ko.Spec.Code.S3ObjectVersion) {
			delta.Add("Spec.Code.S3ObjectVersion", a.ko.Spec.Code.S3ObjectVersion, b.ko.Spec.Code.S3ObjectVersion)
		} else if a.ko.Spec.Code.S3ObjectVersion != nil && b.ko.Spec.Code.S3ObjectVersion != nil {
			if *a.ko.Spec.Code.S3ObjectVersion != *b.ko.Spec.Code.S3ObjectVersion {
				delta.Add("Spec.Code.S3ObjectVersion", a.ko.Spec.Code.S3ObjectVersion, b.ko.Spec.Code.S3ObjectVersion)
			}
		}
		if !bytes.Equal(a.ko.Spec.Code.ZipFile, b.ko.Spec.Code.ZipFile) {
			delta.Add("Spec.Code.ZipFile", a.ko.Spec.Code.ZipFile, b.ko.Spec.Code.ZipFile)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.CodeSigningConfigARN, b.ko.Spec.CodeSigningConfigARN) {
		delta.Add("Spec.CodeSigningConfigARN", a.ko.Spec.CodeSigningConfigARN, b.ko.Spec.CodeSigningConfigARN)
	} else if a.ko.Spec.CodeSigningConfigARN != nil && b.ko.Spec.CodeSigningConfigARN != nil {
		if *a.ko.Spec.CodeSigningConfigARN != *b.ko.Spec.CodeSigningConfigARN {
			delta.Add("Spec.CodeSigningConfigARN", a.ko.Spec.CodeSigningConfigARN, b.ko.Spec.CodeSigningConfigARN)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.DeadLetterConfig, b.ko.Spec.DeadLetterConfig) {
		delta.Add("Spec.DeadLetterConfig", a.ko.Spec.DeadLetterConfig, b.ko.Spec.DeadLetterConfig)
	} else if a.ko.Spec.DeadLetterConfig != nil && b.ko.Spec.DeadLetterConfig != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.DeadLetterConfig.TargetARN, b.ko.Spec.DeadLetterConfig.TargetARN) {
			delta.Add("Spec.DeadLetterConfig.TargetARN", a.ko.Spec.DeadLetterConfig.TargetARN, b.ko.Spec.DeadLetterConfig.TargetARN)
		} else if a.ko.Spec.DeadLetterConfig.TargetARN != nil && b.ko.Spec.DeadLetterConfig.TargetARN != nil {
			if *a.ko.Spec.DeadLetterConfig.TargetARN != *b.ko.Spec.DeadLetterConfig.TargetARN {
				delta.Add("Spec.DeadLetterConfig.TargetARN", a.ko.Spec.DeadLetterConfig.TargetARN, b.ko.Spec.DeadLetterConfig.TargetARN)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Description, b.ko.Spec.Description) {
		delta.Add("Spec.Description", a.ko.Spec.Description, b.ko.Spec.Description)
	} else if a.ko.Spec.Description != nil && b.ko.Spec.Description != nil {
		if *a.ko.Spec.Description != *b.ko.Spec.Description {
			delta.Add("Spec.Description", a.ko.Spec.Description, b.ko.Spec.Description)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Environment, b.ko.Spec.Environment) {
		delta.Add("Spec.Environment", a.ko.Spec.Environment, b.ko.Spec.Environment)
	} else if a.ko.Spec.Environment != nil && b.ko.Spec.Environment != nil {
		if len(a.ko.Spec.Environment.Variables) != len(b.ko.Spec.Environment.Variables) {
			delta.Add("Spec.Environment.Variables", a.ko.Spec.Environment.Variables, b.ko.Spec.Environment.Variables)
		} else if len(a.ko.Spec.Environment.Variables) > 0 {
			if !ackcompare.MapStringStringPEqual(a.ko.Spec.Environment.Variables, b.ko.Spec.Environment.Variables) {
				delta.Add("Spec.Environment.Variables", a.ko.Spec.Environment.Variables, b.ko.Spec.Environment.Variables)
			}
		}
	}
	if len(a.ko.Spec.FileSystemConfigs) != len(b.ko.Spec.FileSystemConfigs) {
		delta.Add("Spec.FileSystemConfigs", a.ko.Spec.FileSystemConfigs, b.ko.Spec.FileSystemConfigs)
	} else if len(a.ko.Spec.FileSystemConfigs) > 0 {
		if !equality.Semantic.Equalities.DeepEqual(a.ko.Spec.FileSystemConfigs, b.ko.Spec.FileSystemConfigs) {
			delta.Add("Spec.FileSystemConfigs", a.ko.Spec.FileSystemConfigs, b.ko.Spec.FileSystemConfigs)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.FunctionName, b.ko.Spec.FunctionName) {
		delta.Add("Spec.FunctionName", a.ko.Spec.FunctionName, b.ko.Spec.FunctionName)
	} else if a.ko.Spec.FunctionName != nil && b.ko.Spec.FunctionName != nil {
		if *a.ko.Spec.FunctionName != *b.ko.Spec.FunctionName {
			delta.Add("Spec.FunctionName", a.ko.Spec.FunctionName, b.ko.Spec.FunctionName)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Handler, b.ko.Spec.Handler) {
		delta.Add("Spec.Handler", a.ko.Spec.Handler, b.ko.Spec.Handler)
	} else if a.ko.Spec.Handler != nil && b.ko.Spec.Handler != nil {
		if *a.ko.Spec.Handler != *b.ko.Spec.Handler {
			delta.Add("Spec.Handler", a.ko.Spec.Handler, b.ko.Spec.Handler)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ImageConfig, b.ko.Spec.ImageConfig) {
		delta.Add("Spec.ImageConfig", a.ko.Spec.ImageConfig, b.ko.Spec.ImageConfig)
	} else if a.ko.Spec.ImageConfig != nil && b.ko.Spec.ImageConfig != nil {
		if len(a.ko.Spec.ImageConfig.Command) != len(b.ko.Spec.ImageConfig.Command) {
			delta.Add("Spec.ImageConfig.Command", a.ko.Spec.ImageConfig.Command, b.ko.Spec.ImageConfig.Command)
		} else if len(a.ko.Spec.ImageConfig.Command) > 0 {
			if !ackcompare.SliceStringPEqual(a.ko.Spec.ImageConfig.Command, b.ko.Spec.ImageConfig.Command) {
				delta.Add("Spec.ImageConfig.Command", a.ko.Spec.ImageConfig.Command, b.ko.Spec.ImageConfig.Command)
			}
		}
		if len(a.ko.Spec.ImageConfig.EntryPoint) != len(b.ko.Spec.ImageConfig.EntryPoint) {
			delta.Add("Spec.ImageConfig.EntryPoint", a.ko.Spec.ImageConfig.EntryPoint, b.ko.Spec.ImageConfig.EntryPoint)
		} else if len(a.ko.Spec.ImageConfig.EntryPoint) > 0 {
			if !ackcompare.SliceStringPEqual(a.ko.Spec.ImageConfig.EntryPoint, b.ko.Spec.ImageConfig.EntryPoint) {
				delta.Add("Spec.ImageConfig.EntryPoint", a.ko.Spec.ImageConfig.EntryPoint, b.ko.Spec.ImageConfig.EntryPoint)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ImageConfig.WorkingDirectory, b.ko.Spec.ImageConfig.WorkingDirectory) {
			delta.Add("Spec.ImageConfig.WorkingDirectory", a.ko.Spec.ImageConfig.WorkingDirectory, b.ko.Spec.ImageConfig.WorkingDirectory)
		} else if a.ko.Spec.ImageConfig.WorkingDirectory != nil && b.ko.Spec.ImageConfig.WorkingDirectory != nil {
			if *a.ko.Spec.ImageConfig.WorkingDirectory != *b.ko.Spec.ImageConfig.WorkingDirectory {
				delta.Add("Spec.ImageConfig.WorkingDirectory", a.ko.Spec.ImageConfig.WorkingDirectory, b.ko.Spec.ImageConfig.WorkingDirectory)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.KMSKeyARN, b.ko.Spec.KMSKeyARN) {
		delta.Add("Spec.KMSKeyARN", a.ko.Spec.KMSKeyARN, b.ko.Spec.KMSKeyARN)
	} else if a.ko.Spec.KMSKeyARN != nil && b.ko.Spec.KMSKeyARN != nil {
		if *a.ko.Spec.KMSKeyARN != *b.ko.Spec.KMSKeyARN {
			delta.Add("Spec.KMSKeyARN", a.ko.Spec.KMSKeyARN, b.ko.Spec.KMSKeyARN)
		}
	}
	if len(a.ko.Spec.Layers) != len(b.ko.Spec.Layers) {
		delta.Add("Spec.Layers", a.ko.Spec.Layers, b.ko.Spec.Layers)
	} else if len(a.ko.Spec.Layers) > 0 {
		if !ackcompare.SliceStringPEqual(a.ko.Spec.Layers, b.ko.Spec.Layers) {
			delta.Add("Spec.Layers", a.ko.Spec.Layers, b.ko.Spec.Layers)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.MemorySize, b.ko.Spec.MemorySize) {
		delta.Add("Spec.MemorySize", a.ko.Spec.MemorySize, b.ko.Spec.MemorySize)
	} else if a.ko.Spec.MemorySize != nil && b.ko.Spec.MemorySize != nil {
		if *a.ko.Spec.MemorySize != *b.ko.Spec.MemorySize {
			delta.Add("Spec.MemorySize", a.ko.Spec.MemorySize, b.ko.Spec.MemorySize)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.PackageType, b.ko.Spec.PackageType) {
		delta.Add("Spec.PackageType", a.ko.Spec.PackageType, b.ko.Spec.PackageType)
	} else if a.ko.Spec.PackageType != nil && b.ko.Spec.PackageType != nil {
		if *a.ko.Spec.PackageType != *b.ko.Spec.PackageType {
			delta.Add("Spec.PackageType", a.ko.Spec.PackageType, b.ko.Spec.PackageType)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Publish, b.ko.Spec.Publish) {
		delta.Add("Spec.Publish", a.ko.Spec.Publish, b.ko.Spec.Publish)
	} else if a.ko.Spec.Publish != nil && b.ko.Spec.Publish != nil {
		if *a.ko.Spec.Publish != *b.ko.Spec.Publish {
			delta.Add("Spec.Publish", a.ko.Spec.Publish, b.ko.Spec.Publish)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Role, b.ko.Spec.Role) {
		delta.Add("Spec.Role", a.ko.Spec.Role, b.ko.Spec.Role)
	} else if a.ko.Spec.Role != nil && b.ko.Spec.Role != nil {
		if *a.ko.Spec.Role != *b.ko.Spec.Role {
			delta.Add("Spec.Role", a.ko.Spec.Role, b.ko.Spec.Role)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Runtime, b.ko.Spec.Runtime) {
		delta.Add("Spec.Runtime", a.ko.Spec.Runtime, b.ko.Spec.Runtime)
	} else if a.ko.Spec.Runtime != nil && b.ko.Spec.Runtime != nil {
		if *a.ko.Spec.Runtime != *b.ko.Spec.Runtime {
			delta.Add("Spec.Runtime", a.ko.Spec.Runtime, b.ko.Spec.Runtime)
		}
	}
	desiredACKTags, _ := convertToOrderedACKTags(a.ko.Spec.Tags)
	latestACKTags, _ := convertToOrderedACKTags(b.ko.Spec.Tags)
	if !ackcompare.MapStringStringEqual(desiredACKTags, latestACKTags) {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Timeout, b.ko.Spec.Timeout) {
		delta.Add("Spec.Timeout", a.ko.Spec.Timeout, b.ko.Spec.Timeout)
	} else if a.ko.Spec.Timeout != nil && b.ko.Spec.Timeout != nil {
		if *a.ko.Spec.Timeout != *b.ko.Spec.Timeout {
			delta.Add("Spec.Timeout", a.ko.Spec.Timeout, b.ko.Spec.Timeout)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.TracingConfig, b.ko.Spec.TracingConfig) {
		delta.Add("Spec.TracingConfig", a.ko.Spec.TracingConfig, b.ko.Spec.TracingConfig)
	} else if a.ko.Spec.TracingConfig != nil && b.ko.Spec.TracingConfig != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.TracingConfig.Mode, b.ko.Spec.TracingConfig.Mode) {
			delta.Add("Spec.TracingConfig.Mode", a.ko.Spec.TracingConfig.Mode, b.ko.Spec.TracingConfig.Mode)
		} else if a.ko.Spec.TracingConfig.Mode != nil && b.ko.Spec.TracingConfig.Mode != nil {
			if *a.ko.Spec.TracingConfig.Mode != *b.ko.Spec.TracingConfig.Mode {
				delta.Add("Spec.TracingConfig.Mode", a.ko.Spec.TracingConfig.Mode, b.ko.Spec.TracingConfig.Mode)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.VPCConfig, b.ko.Spec.VPCConfig) {
		delta.Add("Spec.VPCConfig", a.ko.Spec.VPCConfig, b.ko.Spec.VPCConfig)
	} else if a.ko.Spec.VPCConfig != nil && b.ko.Spec.VPCConfig != nil {
		if len(a.ko.Spec.VPCConfig.SecurityGroupIDs) != len(b.ko.Spec.VPCConfig.SecurityGroupIDs) {
			delta.Add("Spec.VPCConfig.SecurityGroupIDs", a.ko.Spec.VPCConfig.SecurityGroupIDs, b.ko.Spec.VPCConfig.SecurityGroupIDs)
		} else if len(a.ko.Spec.VPCConfig.SecurityGroupIDs) > 0 {
			if !ackcompare.SliceStringPEqual(a.ko.Spec.VPCConfig.SecurityGroupIDs, b.ko.Spec.VPCConfig.SecurityGroupIDs) {
				delta.Add("Spec.VPCConfig.SecurityGroupIDs", a.ko.Spec.VPCConfig.SecurityGroupIDs, b.ko.Spec.VPCConfig.SecurityGroupIDs)
			}
		}
		if len(a.ko.Spec.VPCConfig.SubnetIDs) != len(b.ko.Spec.VPCConfig.SubnetIDs) {
			delta.Add("Spec.VPCConfig.SubnetIDs", a.ko.Spec.VPCConfig.SubnetIDs, b.ko.Spec.VPCConfig.SubnetIDs)
		} else if len(a.ko.Spec.VPCConfig.SubnetIDs) > 0 {
			if !ackcompare.SliceStringPEqual(a.ko.Spec.VPCConfig.SubnetIDs, b.ko.Spec.VPCConfig.SubnetIDs) {
				delta.Add("Spec.VPCConfig.SubnetIDs", a.ko.Spec.VPCConfig.SubnetIDs, b.ko.Spec.VPCConfig.SubnetIDs)
			}
		}
	}
`
	assert.Equal(
		expected,
		code.CompareResource(
			crd.Config(), crd, "delta", "a.ko", "b.ko", 1,
		),
	)
}

func TestCompareResource_APIGatewayv2_Route(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "Route")
	require.NotNil(crd)

	expected := `
	if ackcompare.HasNilDifference(a.ko.Spec.APIID, b.ko.Spec.APIID) {
		delta.Add("Spec.APIID", a.ko.Spec.APIID, b.ko.Spec.APIID)
	} else if a.ko.Spec.APIID != nil && b.ko.Spec.APIID != nil {
		if *a.ko.Spec.APIID != *b.ko.Spec.APIID {
			delta.Add("Spec.APIID", a.ko.Spec.APIID, b.ko.Spec.APIID)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.APIKeyRequired, b.ko.Spec.APIKeyRequired) {
		delta.Add("Spec.APIKeyRequired", a.ko.Spec.APIKeyRequired, b.ko.Spec.APIKeyRequired)
	} else if a.ko.Spec.APIKeyRequired != nil && b.ko.Spec.APIKeyRequired != nil {
		if *a.ko.Spec.APIKeyRequired != *b.ko.Spec.APIKeyRequired {
			delta.Add("Spec.APIKeyRequired", a.ko.Spec.APIKeyRequired, b.ko.Spec.APIKeyRequired)
		}
	}
	if len(a.ko.Spec.AuthorizationScopes) != len(b.ko.Spec.AuthorizationScopes) {
		delta.Add("Spec.AuthorizationScopes", a.ko.Spec.AuthorizationScopes, b.ko.Spec.AuthorizationScopes)
	} else if len(a.ko.Spec.AuthorizationScopes) > 0 {
		if !ackcompare.SliceStringPEqual(a.ko.Spec.AuthorizationScopes, b.ko.Spec.AuthorizationScopes) {
			delta.Add("Spec.AuthorizationScopes", a.ko.Spec.AuthorizationScopes, b.ko.Spec.AuthorizationScopes)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.AuthorizationType, b.ko.Spec.AuthorizationType) {
		delta.Add("Spec.AuthorizationType", a.ko.Spec.AuthorizationType, b.ko.Spec.AuthorizationType)
	} else if a.ko.Spec.AuthorizationType != nil && b.ko.Spec.AuthorizationType != nil {
		if *a.ko.Spec.AuthorizationType != *b.ko.Spec.AuthorizationType {
			delta.Add("Spec.AuthorizationType", a.ko.Spec.AuthorizationType, b.ko.Spec.AuthorizationType)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.AuthorizerID, b.ko.Spec.AuthorizerID) {
		delta.Add("Spec.AuthorizerID", a.ko.Spec.AuthorizerID, b.ko.Spec.AuthorizerID)
	} else if a.ko.Spec.AuthorizerID != nil && b.ko.Spec.AuthorizerID != nil {
		if *a.ko.Spec.AuthorizerID != *b.ko.Spec.AuthorizerID {
			delta.Add("Spec.AuthorizerID", a.ko.Spec.AuthorizerID, b.ko.Spec.AuthorizerID)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ModelSelectionExpression, b.ko.Spec.ModelSelectionExpression) {
		delta.Add("Spec.ModelSelectionExpression", a.ko.Spec.ModelSelectionExpression, b.ko.Spec.ModelSelectionExpression)
	} else if a.ko.Spec.ModelSelectionExpression != nil && b.ko.Spec.ModelSelectionExpression != nil {
		if *a.ko.Spec.ModelSelectionExpression != *b.ko.Spec.ModelSelectionExpression {
			delta.Add("Spec.ModelSelectionExpression", a.ko.Spec.ModelSelectionExpression, b.ko.Spec.ModelSelectionExpression)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.OperationName, b.ko.Spec.OperationName) {
		delta.Add("Spec.OperationName", a.ko.Spec.OperationName, b.ko.Spec.OperationName)
	} else if a.ko.Spec.OperationName != nil && b.ko.Spec.OperationName != nil {
		if *a.ko.Spec.OperationName != *b.ko.Spec.OperationName {
			delta.Add("Spec.OperationName", a.ko.Spec.OperationName, b.ko.Spec.OperationName)
		}
	}
	if len(a.ko.Spec.RequestModels) != len(b.ko.Spec.RequestModels) {
		delta.Add("Spec.RequestModels", a.ko.Spec.RequestModels, b.ko.Spec.RequestModels)
	} else if len(a.ko.Spec.RequestModels) > 0 {
		if !ackcompare.MapStringStringPEqual(a.ko.Spec.RequestModels, b.ko.Spec.RequestModels) {
			delta.Add("Spec.RequestModels", a.ko.Spec.RequestModels, b.ko.Spec.RequestModels)
		}
	}
	if len(a.ko.Spec.RequestParameters) != len(b.ko.Spec.RequestParameters) {
		delta.Add("Spec.RequestParameters", a.ko.Spec.RequestParameters, b.ko.Spec.RequestParameters)
	} else if len(a.ko.Spec.RequestParameters) > 0 {
		if !equality.Semantic.Equalities.DeepEqual(a.ko.Spec.RequestParameters, b.ko.Spec.RequestParameters) {
			delta.Add("Spec.RequestParameters", a.ko.Spec.RequestParameters, b.ko.Spec.RequestParameters)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.RouteKey, b.ko.Spec.RouteKey) {
		delta.Add("Spec.RouteKey", a.ko.Spec.RouteKey, b.ko.Spec.RouteKey)
	} else if a.ko.Spec.RouteKey != nil && b.ko.Spec.RouteKey != nil {
		if *a.ko.Spec.RouteKey != *b.ko.Spec.RouteKey {
			delta.Add("Spec.RouteKey", a.ko.Spec.RouteKey, b.ko.Spec.RouteKey)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.RouteResponseSelectionExpression, b.ko.Spec.RouteResponseSelectionExpression) {
		delta.Add("Spec.RouteResponseSelectionExpression", a.ko.Spec.RouteResponseSelectionExpression, b.ko.Spec.RouteResponseSelectionExpression)
	} else if a.ko.Spec.RouteResponseSelectionExpression != nil && b.ko.Spec.RouteResponseSelectionExpression != nil {
		if *a.ko.Spec.RouteResponseSelectionExpression != *b.ko.Spec.RouteResponseSelectionExpression {
			delta.Add("Spec.RouteResponseSelectionExpression", a.ko.Spec.RouteResponseSelectionExpression, b.ko.Spec.RouteResponseSelectionExpression)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Target, b.ko.Spec.Target) {
		delta.Add("Spec.Target", a.ko.Spec.Target, b.ko.Spec.Target)
	} else if a.ko.Spec.Target != nil && b.ko.Spec.Target != nil {
		if *a.ko.Spec.Target != *b.ko.Spec.Target {
			delta.Add("Spec.Target", a.ko.Spec.Target, b.ko.Spec.Target)
		}
	}
`
	assert.Equal(
		expected,
		code.CompareResource(
			crd.Config(), crd, "delta", "a.ko", "b.ko", 1,
		),
	)
}

func TestCompareResource_IAM_OIDC_URL(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "iam", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-oidc-url.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "OpenIDConnectProvider")
	require.NotNil(crd)
	expected := `
	if len(a.ko.Spec.ClientIDList) != len(b.ko.Spec.ClientIDList) {
		delta.Add("Spec.ClientIDList", a.ko.Spec.ClientIDList, b.ko.Spec.ClientIDList)
	} else if len(a.ko.Spec.ClientIDList) > 0 {
		if !ackcompare.SliceStringPEqual(a.ko.Spec.ClientIDList, b.ko.Spec.ClientIDList) {
			delta.Add("Spec.ClientIDList", a.ko.Spec.ClientIDList, b.ko.Spec.ClientIDList)
		}
	}
	desiredACKTags, _ := convertToOrderedACKTags(a.ko.Spec.Tags)
	latestACKTags, _ := convertToOrderedACKTags(b.ko.Spec.Tags)
	if !ackcompare.MapStringStringEqual(desiredACKTags, latestACKTags) {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	}
	if len(a.ko.Spec.ThumbprintList) != len(b.ko.Spec.ThumbprintList) {
		delta.Add("Spec.ThumbprintList", a.ko.Spec.ThumbprintList, b.ko.Spec.ThumbprintList)
	} else if len(a.ko.Spec.ThumbprintList) > 0 {
		if !ackcompare.SliceStringPEqual(a.ko.Spec.ThumbprintList, b.ko.Spec.ThumbprintList) {
			delta.Add("Spec.ThumbprintList", a.ko.Spec.ThumbprintList, b.ko.Spec.ThumbprintList)
		}
	}
`
	assert.Equal(
		expected,
		code.CompareResource(
			crd.Config(), crd, "delta", "a.ko", "b.ko", 1,
		),
	)
}

func TestCompareResource_MemoryDB_User(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "memorydb")

	crd := testutil.GetCRDByName(t, g, "User")
	require.NotNil(crd)
	expected := `
	if ackcompare.HasNilDifference(a.ko.Spec.AccessString, b.ko.Spec.AccessString) {
		delta.Add("Spec.AccessString", a.ko.Spec.AccessString, b.ko.Spec.AccessString)
	} else if a.ko.Spec.AccessString != nil && b.ko.Spec.AccessString != nil {
		if *a.ko.Spec.AccessString != *b.ko.Spec.AccessString {
			delta.Add("Spec.AccessString", a.ko.Spec.AccessString, b.ko.Spec.AccessString)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.AuthenticationMode, b.ko.Spec.AuthenticationMode) {
		delta.Add("Spec.AuthenticationMode", a.ko.Spec.AuthenticationMode, b.ko.Spec.AuthenticationMode)
	} else if a.ko.Spec.AuthenticationMode != nil && b.ko.Spec.AuthenticationMode != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.AuthenticationMode.Type, b.ko.Spec.AuthenticationMode.Type) {
			delta.Add("Spec.AuthenticationMode.Type", a.ko.Spec.AuthenticationMode.Type, b.ko.Spec.AuthenticationMode.Type)
		} else if a.ko.Spec.AuthenticationMode.Type != nil && b.ko.Spec.AuthenticationMode.Type != nil {
			if *a.ko.Spec.AuthenticationMode.Type != *b.ko.Spec.AuthenticationMode.Type {
				delta.Add("Spec.AuthenticationMode.Type", a.ko.Spec.AuthenticationMode.Type, b.ko.Spec.AuthenticationMode.Type)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name) {
		delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	} else if a.ko.Spec.Name != nil && b.ko.Spec.Name != nil {
		if *a.ko.Spec.Name != *b.ko.Spec.Name {
			delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
		}
	}
	desiredACKTags, _ := convertToOrderedACKTags(a.ko.Spec.Tags)
	latestACKTags, _ := convertToOrderedACKTags(b.ko.Spec.Tags)
	if !ackcompare.MapStringStringEqual(desiredACKTags, latestACKTags) {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	}
`
	assert.Equal(
		expected,
		code.CompareResource(
			crd.Config(), crd, "delta", "a.ko", "b.ko", 1,
		),
	)
}
