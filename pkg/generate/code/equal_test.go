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

func TestCompareResource_ApigatewayV2_RouteSettings(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "apigatewayv2")

	tdef := testutil.GetTypeDefByName(t, g, "RouteSettings")
	require.NotNil(tdef)

	expected := `	if ackcompare.HasNilDifference(a.DataTraceEnabled, b.DataTraceEnabled) {
		return false
	} else if a.DataTraceEnabled != nil && b.DataTraceEnabled != nil {
		if *a.DataTraceEnabled != *b.DataTraceEnabled {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.DetailedMetricsEnabled, b.DetailedMetricsEnabled) {
		return false
	} else if a.DetailedMetricsEnabled != nil && b.DetailedMetricsEnabled != nil {
		if *a.DetailedMetricsEnabled != *b.DetailedMetricsEnabled {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.LoggingLevel, b.LoggingLevel) {
		return false
	} else if a.LoggingLevel != nil && b.LoggingLevel != nil {
		if *a.LoggingLevel != *b.LoggingLevel {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.ThrottlingBurstLimit, b.ThrottlingBurstLimit) {
		return false
	} else if a.ThrottlingBurstLimit != nil && b.ThrottlingBurstLimit != nil {
		if *a.ThrottlingBurstLimit != *b.ThrottlingBurstLimit {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.ThrottlingRateLimit, b.ThrottlingRateLimit) {
		return false
	} else if a.ThrottlingRateLimit != nil && b.ThrottlingRateLimit != nil {
		if *a.ThrottlingRateLimit != *b.ThrottlingRateLimit {
			return false
		}
	}
`
	assert.Equal(
		expected,
		code.IsEqualTypeDef(
			tdef, "a", "b", 1,
		),
	)
}

func TestCompareResource_ECR_Repository_SDK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "ecr")

	tdef := testutil.GetTypeDefByName(t, g, "Repository")
	require.NotNil(tdef)

	expected := `	if ackcompare.HasNilDifference(a.CreatedAt, b.CreatedAt) {
		return false
	} else if a.CreatedAt != nil && b.CreatedAt != nil {
		if *a.CreatedAt != *b.CreatedAt {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.ImageScanningConfiguration, b.ImageScanningConfiguration) {
		return false
	} else if a.ImageScanningConfiguration != nil && b.ImageScanningConfiguration != nil {
		if !IsEqualImageScanningConfiguration(a.ImageScanningConfiguration, b.ImageScanningConfiguration) {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.ImageTagMutability, b.ImageTagMutability) {
		return false
	} else if a.ImageTagMutability != nil && b.ImageTagMutability != nil {
		if *a.ImageTagMutability != *b.ImageTagMutability {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.RegistryID, b.RegistryID) {
		return false
	} else if a.RegistryID != nil && b.RegistryID != nil {
		if *a.RegistryID != *b.RegistryID {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.RepositoryARN, b.RepositoryARN) {
		return false
	} else if a.RepositoryARN != nil && b.RepositoryARN != nil {
		if *a.RepositoryARN != *b.RepositoryARN {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.RepositoryName, b.RepositoryName) {
		return false
	} else if a.RepositoryName != nil && b.RepositoryName != nil {
		if *a.RepositoryName != *b.RepositoryName {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.RepositoryURI, b.RepositoryURI) {
		return false
	} else if a.RepositoryURI != nil && b.RepositoryURI != nil {
		if *a.RepositoryURI != *b.RepositoryURI {
			return false
		}
	}
`

	assert.Equal(
		expected,
		code.IsEqualTypeDef(
			tdef, "a", "b", 1,
		),
	)
}

func TestCompareResource_ApigatewayV2_JWTConfiguration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "apigatewayv2")

	tdef := testutil.GetTypeDefByName(t, g, "JWTConfiguration")
	require.NotNil(tdef)

	expected := `
	//TODO(a-hilaly): equality check for slices
	if ackcompare.HasNilDifference(a.Issuer, b.Issuer) {
		return false
	} else if a.Issuer != nil && b.Issuer != nil {
		if *a.Issuer != *b.Issuer {
			return false
		}
	}
`

	assert.Equal(
		expected,
		code.IsEqualTypeDef(
			tdef, "a", "b", 1,
		),
	)
}
