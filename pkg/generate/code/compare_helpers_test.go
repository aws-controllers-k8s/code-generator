package code_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestCompareResourceHelpers_ApiGatewayV2_Stage(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "Stage")
	require.NotNil(crd)

	expected := `
func equalRouteSettingsMap(a, b map[string]*svcapitypes.RouteSettings) bool {
	if len(a) != len(b) {
		return false
	}
	for ka := range a {
		_, ok := b[ka]
		if !ok {
			return false
		}
	}
	for ka, aX := range a {
		bY := b[ka]
		if ackcompare.HasNilDifference(aX.DataTraceEnabled, bY.DataTraceEnabled) {
			return false
		} else if aX.DataTraceEnabled != nil && bY.DataTraceEnabled != nil {
			if *aX.DataTraceEnabled != *bY.DataTraceEnabled {
				return false
			}
		}
		if ackcompare.HasNilDifference(aX.DetailedMetricsEnabled, bY.DetailedMetricsEnabled) {
			return false
		} else if aX.DetailedMetricsEnabled != nil && bY.DetailedMetricsEnabled != nil {
			if *aX.DetailedMetricsEnabled != *bY.DetailedMetricsEnabled {
				return false
			}
		}
		if ackcompare.HasNilDifference(aX.LoggingLevel, bY.LoggingLevel) {
			return false
		} else if aX.LoggingLevel != nil && bY.LoggingLevel != nil {
			if *aX.LoggingLevel != *bY.LoggingLevel {
				return false
			}
		}
		if ackcompare.HasNilDifference(aX.ThrottlingBurstLimit, bY.ThrottlingBurstLimit) {
			return false
		} else if aX.ThrottlingBurstLimit != nil && bY.ThrottlingBurstLimit != nil {
			if *aX.ThrottlingBurstLimit != *bY.ThrottlingBurstLimit {
				return false
			}
		}
		if ackcompare.HasNilDifference(aX.ThrottlingRateLimit, bY.ThrottlingRateLimit) {
			return false
		} else if aX.ThrottlingRateLimit != nil && bY.ThrottlingRateLimit != nil {
			if *aX.ThrottlingRateLimit != *bY.ThrottlingRateLimit {
				return false
			}
		}
	}
	return true
}
`
	assert.Equal(
		expected,
		code.CompareResourceHelpers(
			crd.Config(), crd, 1,
		),
	)
}

func TestCompareResourceHelpers_ApiGatewayV2_DomainNameConfiguration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "DomainName")
	require.NotNil(crd)

	expected := `
func equalDomainNameConfigurations(a, b []*svcapitypes.DomainNameConfiguration) bool {
	//TODO(a-hilaly) implement this function
	return true
}
`
	assert.Equal(
		expected,
		code.CompareResourceHelpers(
			crd.Config(), crd, 1,
		),
	)
}
