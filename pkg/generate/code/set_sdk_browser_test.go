package code_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

// TestSetSDK_BedrockAgentCoreControl_Browser_DefaultBool verifies that a
// boolean member with @default on the member reference (targeting
// smithy.api#Boolean) is treated as a value type (non-pointer) in generated
// code. Specifically:
//   - SetSDK (Create): should dereference the CRD field with * to assign to the
//     non-pointer SDK field (e.g. f0.Enabled = *r.ko.Spec.BrowserSigning.Enabled)
//   - SetResource (ReadOne): should use & to take the address when assigning to
//     the CRD pointer field (e.g. f2.Enabled = &resp.BrowserSigning.Enabled)
//     and should NOT generate a nil check (the SDK field is not a pointer).
func TestSetSDK_BedrockAgentCoreControl_Browser_DefaultBool(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "bedrock-agentcore-control", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-browser.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "Browser")
	require.NotNil(crd)
	require.NotNil(crd.Ops.Create)

	// --- SetSDK Create ---
	// The SDK's BrowserSigningConfigInput.Enabled is a non-pointer bool because
	// the member targets smithy.api#Boolean with @default(false). The generated
	// code must dereference the CRD's *bool field with *.
	sdkOut, err := code.SetSDK(crd.Config(), crd, model.OpTypeCreate, "r.ko", "res", 1)
	require.NoError(err)

	// Should dereference: f<N>.Enabled = *r.ko.Spec.BrowserSigning.Enabled
	assert.Contains(sdkOut, "= *r.ko.Spec.BrowserSigning.Enabled")
	// Should NOT assign without dereference (pointer to non-pointer mismatch)
	assert.NotContains(sdkOut, "= r.ko.Spec.BrowserSigning.Enabled\n")

	// --- SetResource ReadOne ---
	require.NotNil(crd.Ops.ReadOne)
	resOut, err := code.SetResource(crd.Config(), crd, model.OpTypeGet, "resp", "ko", 1)
	require.NoError(err)

	// Should take address: f<N>.Enabled = &resp.BrowserSigning.Enabled
	assert.Contains(resOut, "= &resp.BrowserSigning.Enabled")
	// Should NOT generate a nil check for the non-pointer bool field
	assert.NotContains(resOut, "if resp.BrowserSigning.Enabled != nil")
}
