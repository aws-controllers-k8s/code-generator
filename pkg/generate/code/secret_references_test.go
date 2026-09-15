// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package code_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestSecretReferences_RDS_DBInstance(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "rds")
	crd := testutil.GetCRDByName(t, g, "DBInstance")
	require.NotNil(crd)

	expected := `	var refs []*ackv1alpha1.SecretKeyReference
	if r.ko.Spec.MasterUserPassword != nil {
		refs = append(refs, r.ko.Spec.MasterUserPassword)
	}
	return refs
`
	got := code.SecretReferences(crd, "r.ko", 1)
	assert.Equal(expected, got)
}

func TestSecretReferences_Lambda_Function_MapOfSecrets(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "lambda", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-environment-vars-secret.yaml",
	})
	crd := testutil.GetCRDByName(t, g, "Function")
	require.NotNil(crd)

	expected := `	var refs []*ackv1alpha1.SecretKeyReference
	if r.ko.Spec.Environment != nil {
		for _, v := range r.ko.Spec.Environment.Variables {
			if v != nil {
				refs = append(refs, v)
			}
		}
	}
	return refs
`
	got := code.SecretReferences(crd, "r.ko", 1)
	assert.Equal(expected, got)
}

func TestSecretReferences_MQ_Broker_SliceOfStructsWithSecret(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "mq")
	crd := testutil.GetCRDByName(t, g, "Broker")
	require.NotNil(crd)

	expected := `	var refs []*ackv1alpha1.SecretKeyReference
	if r.ko.Spec.Users != nil {
		for _, elem := range r.ko.Spec.Users {
			if elem == nil {
				continue
			}
			if elem.Password != nil {
				refs = append(refs, elem.Password)
			}
		}
	}
	return refs
`
	got := code.SecretReferences(crd, "r.ko", 1)
	assert.Equal(expected, got)
}

func TestSecretReferences_MemoryDB_User_SliceSecretInStruct(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "memorydb")
	crd := testutil.GetCRDByName(t, g, "User")
	require.NotNil(crd)

	expected := `	var refs []*ackv1alpha1.SecretKeyReference
	if r.ko.Spec.AuthenticationMode != nil {
		if r.ko.Spec.AuthenticationMode.Passwords != nil {
			refs = append(refs, r.ko.Spec.AuthenticationMode.Passwords...)
		}
	}
	return refs
`
	got := code.SecretReferences(crd, "r.ko", 1)
	assert.Equal(expected, got)
}

func TestSecretReferences_Elasticache_ReplicationGroup_TopLevelSecret(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")
	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	got := code.SecretReferences(crd, "r.ko", 1)
	assert.Contains(got, "r.ko.Spec.AuthToken")
	assert.Contains(got, "refs = append(refs, r.ko.Spec.AuthToken)")
}

func TestCompareResource_RDS_DBInstance_Secret(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "rds")
	crd := testutil.GetCRDByName(t, g, "DBInstance")
	require.NotNil(crd)

	got, err := code.CompareResource(
		crd.Config(), crd, "delta", "a.ko", "b.ko", 1,
	)
	require.NoError(err)

	assert.Contains(got, "ackcompare.HasNilDifference(a.ko.Spec.MasterUserPassword, b.ko.Spec.MasterUserPassword)")
	assert.Contains(got, "ackcompare.SecretKeyReferenceEqual(a.ko.Spec.MasterUserPassword, b.ko.Spec.MasterUserPassword)")
	assert.Contains(got, "ackcompare.SecretDataChanged(a.ko.GetAnnotations(), b.ko.GetAnnotations(), acksecret.IndexKey(a.ko.Spec.MasterUserPassword, a.ko.GetNamespace()))")
}

func TestCompareResource_Elasticache_ReplicationGroup_Secret(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "elasticache")
	crd := testutil.GetCRDByName(t, g, "ReplicationGroup")
	require.NotNil(crd)

	got, err := code.CompareResource(
		crd.Config(), crd, "delta", "a.ko", "b.ko", 1,
	)
	require.NoError(err)

	assert.Contains(got, "ackcompare.SecretKeyReferenceEqual(a.ko.Spec.AuthToken, b.ko.Spec.AuthToken)")
	assert.Contains(got, "ackcompare.SecretDataChanged(a.ko.GetAnnotations(), b.ko.GetAnnotations(), acksecret.IndexKey(a.ko.Spec.AuthToken, a.ko.GetNamespace()))")
}

func TestCompareResource_MemoryDB_User_SecretIsCompareIgnored(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "memorydb")
	crd := testutil.GetCRDByName(t, g, "User")
	require.NotNil(crd)

	got, err := code.CompareResource(
		crd.Config(), crd, "delta", "a.ko", "b.ko", 1,
	)
	require.NoError(err)

	// Passwords is compare.is_ignored so it should NOT appear in compare output
	assert.NotContains(got, "Passwords")
	assert.NotContains(got, "SecretDataChanged")
}

func TestHasSecretFields(t *testing.T) {
	require := require.New(t)

	tests := []struct {
		service    string
		resource   string
		options    *testutil.TestingModelOptions
		hasSecret  bool
	}{
		{"rds", "DBInstance", nil, true},
		{"memorydb", "User", nil, true},
		{"s3", "Bucket", nil, false},
		{"lambda", "Function", &testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-environment-vars-secret.yaml",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.service+"_"+tt.resource, func(t *testing.T) {
			if tt.options != nil {
				g := testutil.NewModelForServiceWithOptions(t, tt.service, tt.options)
				crd := testutil.GetCRDByName(t, g, tt.resource)
				require.NotNil(crd)
				assert.Equal(t, tt.hasSecret, crd.HasSecretFields())
			} else {
				g := testutil.NewModelForService(t, tt.service)
				crd := testutil.GetCRDByName(t, g, tt.resource)
				require.NotNil(crd)
				assert.Equal(t, tt.hasSecret, crd.HasSecretFields())
			}
		})
	}
}

func TestSdkFind_RDS_DBInstance_SetsSecretResourceVersions(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "rds")
	crd := testutil.GetCRDByName(t, g, "DBInstance")
	require.NotNil(crd)

	// The sdk_find_read_one template includes SetResourceVersionsAnnotation
	// when HasSecretFields is true. Verify that the CRD reports having secrets
	// and the SecretReferences output is non-empty.
	assert.True(crd.HasSecretFields())
	refs := code.SecretReferences(crd, "r.ko", 1)
	assert.True(strings.Contains(refs, "MasterUserPassword"))
}
