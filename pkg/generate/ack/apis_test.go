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

package ack_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

// templateBasePaths returns the repository's templates/ directory, relative to
// this package.
func templateBasePaths() []string {
	return []string{filepath.Join("..", "..", "..", "templates")}
}

// renderAPIs executes the apis/ template set for the supplied model and returns
// the rendered files keyed by their output path.
func renderAPIs(t *testing.T, m *ackmodel.Model) map[string]string {
	t.Helper()

	ts, err := ackgenerate.APIs(m, templateBasePaths())
	require.NoError(t, err)
	require.NoError(t, ts.Execute())

	out := map[string]string{}
	for path, contents := range ts.Executed() {
		out[path] = contents.String()
	}
	return out
}

// TestAPIs_ImmutableFieldsRenderOnContainingStruct asserts that an is_immutable
// field renders its CEL transition rule as an XValidation marker on the struct
// that *contains* the field, referencing the field by its JSON name, rather than
// as a bare `self == oldSelf` marker on the field itself.
//
// The field-level form is silently unenforced for an optional field: Kubernetes
// skips a transition rule whenever the old value is absent, so a field that was
// unset at creation could be added later and the rule would never fire.
func TestAPIs_ImmutableFieldsRenderOnContainingStruct(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// The route53 RecordSet fixture carries, between them, all three shapes of
	// is_immutable we care about:
	//
	//   Name                (optional, top-level)
	//   SetIdentifier       (optional, top-level)
	//   RecordType          (required, top-level)
	//   AliasTarget.DNSName (optional, nested inside an optional struct)
	m := testutil.NewModelForService(t, "route53")
	rendered := renderAPIs(t, m)

	recordSet, found := rendered[filepath.Join("record_set.go")]
	require.True(found, "expected a rendered record_set.go, got %v", keys(rendered))
	types, found := rendered[filepath.Join("types.go")]
	require.True(found, "expected a rendered types.go, got %v", keys(rendered))

	// The old, ineffective field-level form must be gone everywhere.
	assert.NotContains(recordSet, `rule="self == oldSelf"`)
	assert.NotContains(types, `rule="self == oldSelf"`)

	// The message text is deliberately the original "Value is immutable once
	// set": controller e2e suites assert on that literal string. Field
	// attribution is carried by fieldPath instead, which the API server reports
	// with the rejection -- without it, moving the rule to the parent struct
	// would leave several members sharing one indistinguishable message.
	specMarker := func(jsonName string) string {
		return `// +kubebuilder:validation:XValidation:rule="` +
			`has(self.` + jsonName + `) == has(oldSelf.` + jsonName + `) && ` +
			`(!has(self.` + jsonName + `) || self.` + jsonName + ` == oldSelf.` + jsonName + `)",` +
			`message="Value is immutable once set",fieldPath=".` + jsonName + `"`
	}

	// --- optional, top-level: the case the old rule missed entirely ---
	assert.Contains(recordSet, specMarker("name"))
	assert.Contains(recordSet, specMarker("setIdentifier"))

	// --- required, top-level: still guarded, no regression ---
	assert.Contains(recordSet, specMarker("recordType"))

	// Every marker must sit on the Spec struct, i.e. above `type RecordSetSpec
	// struct {` and not inside the struct body.
	specDecl := "type RecordSetSpec struct {"
	specDeclIdx := strings.Index(recordSet, specDecl)
	require.NotEqual(-1, specDeclIdx)
	for _, jsonName := range []string{"name", "setIdentifier", "recordType"} {
		idx := strings.Index(recordSet, specMarker(jsonName))
		require.NotEqual(-1, idx, "marker for %q not rendered", jsonName)
		assert.Less(idx, specDeclIdx,
			"marker for %q must precede the Spec struct declaration", jsonName)
	}

	// --- optional, nested: the rule belongs to the containing AliasTarget
	// struct, not to the Spec and not to the member ---
	assert.NotContains(recordSet, specMarker("dnsName"),
		"nested field rule must not be emitted on the Spec struct")
	aliasTargetDecl := "type AliasTarget struct {"
	aliasTargetIdx := strings.Index(types, aliasTargetDecl)
	require.NotEqual(-1, aliasTargetIdx)
	dnsNameMarkerIdx := strings.Index(types, specMarker("dnsName"))
	require.NotEqual(-1, dnsNameMarkerIdx,
		"expected an immutability marker for AliasTarget.DNSName in types.go")
	assert.Less(dnsNameMarkerIdx, aliasTargetIdx,
		"marker must precede the AliasTarget struct declaration")

	// ...and it must be the AliasTarget struct it precedes, not some earlier
	// type. Nothing but comment lines may separate the two.
	between := types[dnsNameMarkerIdx+len(specMarker("dnsName")) : aliasTargetIdx]
	for _, line := range strings.Split(between, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		assert.True(strings.HasPrefix(line, "//"),
			"unexpected non-comment line %q between the marker and AliasTarget", line)
	}
}

// TestAPIs_LateInitializedImmutableFieldRendersOnceSetRule asserts the rendered
// marker for a field configured with both late_initialize and is_immutable uses
// the once-set rule, so the controller's own late-init patch is not rejected.
func TestAPIs_LateInitializedImmutableFieldRendersOnceSetRule(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m := testutil.NewModelForService(t, "rds")
	rendered := renderAPIs(t, m)

	dbInstance, found := rendered["db_instance.go"]
	require.True(found, "expected a rendered db_instance.go, got %v", keys(rendered))

	assert.Contains(dbInstance,
		`rule="!has(oldSelf.availabilityZone) || (has(self.availabilityZone) && self.availabilityZone == oldSelf.availabilityZone)"`,
		"late-initialized immutable field must use the once-set rule")
	assert.NotContains(dbInstance,
		`rule="has(self.availabilityZone) == has(oldSelf.availabilityZone)`,
		"late-initialized immutable field must not freeze presence")
	// Message and attribution are unchanged between the two forms.
	assert.Contains(dbInstance, `message="Value is immutable once set",fieldPath=".availabilityZone"`)
}

// TestAPIs_ImmutableReferenceCompanionIsNotEnforced pins the documented non-goal
// that immutability is NOT propagated to a field's generated reference
// companion.
//
// route53 RecordSet.HostedZoneId is is_immutable and reference-backed, so the
// CRD gets both `hostedZoneID` and a `hostedZoneRef`. Only the concrete field
// carries a rule. Because the runtime strips resolved reference values out of
// the stored spec (ClearResolvedReferences), a CR authored through the reference
// keeps `hostedZoneID` absent, so editing only `hostedZoneRef` is not caught.
//
// This test exists to make that gap visible and to fail loudly if the reference
// companion ever starts carrying a rule without the hazard being reconsidered --
// a rule on a *nested* reference companion would reject the controller's own
// patch, since rebuilding the containing struct from an AWS response drops the
// Ref and materializes the resolved value.
func TestAPIs_ImmutableReferenceCompanionIsNotEnforced(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m := testutil.NewModelForService(t, "route53")
	rendered := renderAPIs(t, m)

	recordSet, found := rendered["record_set.go"]
	require.True(found)

	// The concrete field is frozen.
	assert.Contains(recordSet, `fieldPath=".hostedZoneID"`)
	// Its reference companion is generated...
	require.Contains(recordSet, "HostedZoneRef",
		"fixture is expected to generate a reference companion")
	// ...but deliberately carries no rule.
	assert.NotContains(recordSet, `fieldPath=".hostedZoneRef"`)
	assert.NotContains(recordSet, "oldSelf.hostedZoneRef")
}

// TestAPIs_NestedImmutableParentIsNotGuarded pins the documented gap that a
// nested immutable member is only protected while its containing struct is
// present in both the old and the new object.
//
// The rule for AliasTarget.DNSName lives on the AliasTarget struct, so an
// `aliasTarget` absent->present transition skips it entirely (Kubernetes skips a
// transition rule when oldSelf is absent). Guarding that would need a rule on
// the Spec root naming the full path; this PR does not emit one, and the 22
// nested immutable paths in the fleet -- 20 of them in fsx -- therefore still
// need their controller-side runtime guard.
func TestAPIs_NestedImmutableParentIsNotGuarded(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m := testutil.NewModelForService(t, "route53")
	rendered := renderAPIs(t, m)

	recordSet, found := rendered["record_set.go"]
	require.True(found)
	types, found := rendered["types.go"]
	require.True(found)

	// The rule is on the nested struct...
	assert.Contains(types, `fieldPath=".dnsName"`)
	// ...and there is deliberately no Spec-root rule covering the full path, so
	// adding `aliasTarget` wholesale after creation is not rejected.
	assert.NotContains(recordSet, "aliasTarget.dnsName")
	assert.NotContains(recordSet, `fieldPath=".dnsName"`)
	assert.NotContains(recordSet, "oldSelf.aliasTarget")
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
