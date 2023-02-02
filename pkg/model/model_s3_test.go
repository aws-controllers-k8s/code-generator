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

func TestS3_Bucket(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "s3")

	crds, err := g.GetCRDs()
	require.Nil(err)

	// Pronounced "Boo-Kay".
	crd := getCRDByName("Bucket", crds)
	require.NotNil(crd)

	assert.Equal("Bucket", crd.Names.Camel)
	assert.Equal("bucket", crd.Names.CamelLower)
	assert.Equal("bucket", crd.Names.Snake)

	// The ListBucketsResult shape has no defined error codes (in fact, none of
	// the S3 API shapes do). We will need to create exceptions config in the
	// generate.yaml for S3, but this will take quite some manual work. For
	// now, return UNKNOWN
	assert.Equal("UNKNOWN", crd.ExceptionCode(404))

	// The S3 Bucket API is a whole lot of weird. There are Create and Delete
	// operations ("CreateBucket", "DeleteBucket") but there is no ReadOne
	// operation (there are separate API calls for each and every attribute of
	// a Bucket. For instance, there is a GetBucketCord API call, a
	// GetBucketAnalyticsConfiguration API call, a GetBucketLocation call,
	// etc...) or Update operation (there are separate API calls for each and
	// every attribute of a Bucket, though, for instance PutBucketAcl). There
	// is a ReadMany operation (ListBuckets)
	require.NotNil(crd.Ops)

	assert.NotNil(crd.Ops.Create)
	assert.NotNil(crd.Ops.Delete)
	assert.NotNil(crd.Ops.ReadMany)

	assert.Nil(crd.Ops.GetAttributes)
	assert.Nil(crd.Ops.SetAttributes)
	assert.Nil(crd.Ops.ReadOne)
	assert.Nil(crd.Ops.Update)

	specFields := crd.SpecFields
	statusFields := crd.StatusFields

	expSpecFieldCamel := []string{
		"ACL",
		"CreateBucketConfiguration",
		"GrantFullControl",
		"GrantRead",
		"GrantReadACP",
		"GrantWrite",
		"GrantWriteACP",
		"Logging",
		// NOTE(jaypipes): Original field name in CreateBucket input is
		// "Bucket" but should be renamed to "Name" from the generator.yaml (in
		// order to match with the name of the field in the Output shape for a
		// ListBuckets API call...
		"Name",
		"ObjectLockEnabledForBucket",
		"Tagging",
	}
	assert.Equal(expSpecFieldCamel, attrCamelNames(specFields))

	expStatusFieldCamel := []string{
		"Location",
	}
	assert.Equal(expStatusFieldCamel, attrCamelNames(statusFields))

	expTypeDefCamel := []string{
		"BucketLoggingStatus",
		"LoggingEnabled",
		"TargetGrant",
	}
	for _, typeDef := range expTypeDefCamel {
		assert.NotNil(testutil.GetTypeDefByName(t, g, typeDef))
	}
}
