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

	g := testutil.NewGeneratorForService(t, "s3")

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

	g := testutil.NewGeneratorForService(t, "lambda")

	crd := testutil.GetCRDByName(t, g, "CodeSigningConfig")
	require.NotNil(crd)

	expected := `
	if ackcompare.HasNilDifference(a.ko.Spec.AllowedPublishers, b.ko.Spec.AllowedPublishers) {
		delta.Add("Spec.AllowedPublishers", a.ko.Spec.AllowedPublishers, b.ko.Spec.AllowedPublishers)
	} else if a.ko.Spec.AllowedPublishers != nil && b.ko.Spec.AllowedPublishers != nil {

		if !ackcompare.SliceStringPEqual(a.ko.Spec.AllowedPublishers.SigningProfileVersionARNs, b.ko.Spec.AllowedPublishers.SigningProfileVersionARNs) {
			delta.Add("Spec.AllowedPublishers.SigningProfileVersionARNs", a.ko.Spec.AllowedPublishers.SigningProfileVersionARNs, b.ko.Spec.AllowedPublishers.SigningProfileVersionARNs)
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

func TestCompareResource_ApiGatewayv2_DomainName(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "RouteResponse")
	require.NotNil(crd)

	expected := `
	if ackcompare.HasNilDifference(a.ko.Spec.APIID, b.ko.Spec.APIID) {
		delta.Add("Spec.APIID", a.ko.Spec.APIID, b.ko.Spec.APIID)
	} else if a.ko.Spec.APIID != nil && b.ko.Spec.APIID != nil {
		if *a.ko.Spec.APIID != *b.ko.Spec.APIID {
			delta.Add("Spec.APIID", a.ko.Spec.APIID, b.ko.Spec.APIID)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ModelSelectionExpression, b.ko.Spec.ModelSelectionExpression) {
		delta.Add("Spec.ModelSelectionExpression", a.ko.Spec.ModelSelectionExpression, b.ko.Spec.ModelSelectionExpression)
	} else if a.ko.Spec.ModelSelectionExpression != nil && b.ko.Spec.ModelSelectionExpression != nil {
		if *a.ko.Spec.ModelSelectionExpression != *b.ko.Spec.ModelSelectionExpression {
			delta.Add("Spec.ModelSelectionExpression", a.ko.Spec.ModelSelectionExpression, b.ko.Spec.ModelSelectionExpression)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ResponseModels, b.ko.Spec.ResponseModels) {
		delta.Add("Spec.ResponseModels", a.ko.Spec.ResponseModels, b.ko.Spec.ResponseModels)
	} else if a.ko.Spec.ResponseModels != nil && b.ko.Spec.ResponseModels != nil {
		if !ackcompare.MapStringStringPEqual(a.ko.Spec.ResponseModels, b.ko.Spec.ResponseModels) {
			delta.Add("Spec.ResponseModels", a.ko.Spec.ResponseModels, b.ko.Spec.ResponseModels)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ResponseParameters, b.ko.Spec.ResponseParameters) {
		delta.Add("Spec.ResponseParameters", a.ko.Spec.ResponseParameters, b.ko.Spec.ResponseParameters)
	} else if a.ko.Spec.ResponseParameters != nil && b.ko.Spec.ResponseParameters != nil {
		if !equalRouteParameters(a.ko.Spec.ResponseParameters, b.ko.Spec.ResponseParameters) {
			delta.Add("Spec.ResponseParameters", a.ko.Spec.ResponseParameters, b.ko.Spec.ResponseParameters)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.RouteID, b.ko.Spec.RouteID) {
		delta.Add("Spec.RouteID", a.ko.Spec.RouteID, b.ko.Spec.RouteID)
	} else if a.ko.Spec.RouteID != nil && b.ko.Spec.RouteID != nil {
		if *a.ko.Spec.RouteID != *b.ko.Spec.RouteID {
			delta.Add("Spec.RouteID", a.ko.Spec.RouteID, b.ko.Spec.RouteID)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.RouteResponseKey, b.ko.Spec.RouteResponseKey) {
		delta.Add("Spec.RouteResponseKey", a.ko.Spec.RouteResponseKey, b.ko.Spec.RouteResponseKey)
	} else if a.ko.Spec.RouteResponseKey != nil && b.ko.Spec.RouteResponseKey != nil {
		if *a.ko.Spec.RouteResponseKey != *b.ko.Spec.RouteResponseKey {
			delta.Add("Spec.RouteResponseKey", a.ko.Spec.RouteResponseKey, b.ko.Spec.RouteResponseKey)
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
