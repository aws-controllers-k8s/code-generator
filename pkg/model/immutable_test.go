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
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/pkg/names"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

// immutabilityRule builds the strict CEL expression we expect for a member with
// the supplied JSON name: presence and value are both frozen. Kept separate from
// the production helper so a change to the rule has to be made deliberately in
// both places.
func immutabilityRule(jsonName string) string {
	return "has(self." + jsonName + ") == has(oldSelf." + jsonName + ") && " +
		"(!has(self." + jsonName + ") || self." + jsonName + " == oldSelf." + jsonName + ")"
}

// onceSetRule builds the weaker CEL expression expected for a member the
// controller populates itself: the first write is permitted, but once set the
// value may not change and may not be removed.
func onceSetRule(jsonName string) string {
	return "!has(oldSelf." + jsonName + ") || " +
		"(has(self." + jsonName + ") && self." + jsonName + " == oldSelf." + jsonName + ")"
}

// TestRoute53_RecordSet_ImmutableFields covers the three shapes of is_immutable
// on the route53 RecordSet fixture: a top-level optional field, a top-level
// required field, and an optional field nested inside an optional struct.
func TestRoute53_RecordSet_ImmutableFields(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m := testutil.NewModelForService(t, "route53")

	crds, err := m.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("RecordSet", crds)
	require.NotNil(crd)

	// --- top-level, optional ---
	// Name is is_immutable but not is_required, so under a field-level
	// `self == oldSelf` rule it was unenforced when unset at creation.
	nameField := crd.SpecFields["Name"]
	require.NotNil(nameField)
	assert.True(nameField.IsImmutable())
	assert.False(nameField.IsRequired())
	assert.Equal("name", nameField.GetCRDJSONFieldName())
	assert.Equal(immutabilityRule("name"), nameField.ImmutabilityCELRule())

	setIdentifierField := crd.SpecFields["SetIdentifier"]
	require.NotNil(setIdentifierField)
	assert.True(setIdentifierField.IsImmutable())
	assert.False(setIdentifierField.IsRequired())
	assert.Equal(immutabilityRule("setIdentifier"), setIdentifierField.ImmutabilityCELRule())

	// --- top-level, required: no regression, still guarded ---
	recordTypeField := crd.SpecFields["RecordType"]
	require.NotNil(recordTypeField)
	assert.True(recordTypeField.IsImmutable())
	assert.True(recordTypeField.IsRequired())
	assert.Equal(immutabilityRule("recordType"), recordTypeField.ImmutabilityCELRule())

	// --- a non-immutable field must not advertise a rule ---
	ttlField := crd.SpecFields["TTL"]
	require.NotNil(ttlField)
	assert.False(ttlField.IsImmutable())

	// --- fieldPath attributes the rejection to the member, not to the struct ---
	assert.Equal(".name", nameField.ImmutabilityCELFieldPath())
	assert.Equal(".recordType", recordTypeField.ImmutabilityCELFieldPath())

	// --- nested, optional ---
	// AliasTarget.DNSName is recorded on the AliasTarget TypeDef, which is what
	// renders the marker, so the rule can observe the member being absent.
	tds, err := m.GetTypeDefs()
	require.Nil(err)

	var aliasTargetTD *model.TypeDef
	for _, td := range tds {
		if td != nil && strings.EqualFold(td.Names.Original, "AliasTarget") {
			aliasTargetTD = td
			break
		}
	}
	require.NotNil(aliasTargetTD)

	dnsNameAttr := aliasTargetTD.GetAttribute("DNSName")
	require.NotNil(dnsNameAttr)
	assert.True(dnsNameAttr.IsImmutable)
	assert.Equal("dnsName", dnsNameAttr.GetCRDJSONFieldName())
	assert.Equal(immutabilityRule("dnsName"), dnsNameAttr.ImmutabilityCELRule())

	// Sibling members of the same struct are untouched.
	hostedZoneIDAttr := aliasTargetTD.GetAttribute("HostedZoneId")
	require.NotNil(hostedZoneIDAttr)
	assert.False(hostedZoneIDAttr.IsImmutable)

	// The nested field is not a Spec field, so no Spec-level rule is produced
	// for it.
	assert.Nil(crd.SpecFields["DNSName"])
}

// TestRDS_DBInstance_LateInitializedImmutableField asserts that a field
// configured with BOTH late_initialize and is_immutable gets the weaker once-set
// rule rather than the strict presence-freezing one.
//
// The strict rule rejects the controller's own late-initialization patch: that
// patch is an absent->present transition, which is exactly what the strict rule
// forbids. Five fields across acm, iam, mwaa and rds are configured this way, and
// freezing presence for them leaves the resource wedged in a reconcile error.
func TestRDS_DBInstance_LateInitializedImmutableField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	m := testutil.NewModelForService(t, "rds")

	crds, err := m.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("DBInstance", crds)
	require.NotNil(crd)

	azField := crd.SpecFields["AvailabilityZone"]
	require.NotNil(azField)
	require.True(azField.IsImmutable())
	require.True(azField.IsLateInitialized(),
		"fixture is expected to configure late_initialize on AvailabilityZone")

	// The once-set form permits absent->present (the late-init write) but still
	// rejects change-after-set and remove-after-set.
	assert.Equal(onceSetRule("availabilityZone"), azField.ImmutabilityCELRule())
	assert.NotEqual(immutabilityRule("availabilityZone"), azField.ImmutabilityCELRule())

	// A field that is immutable but not late-initialized keeps the strict form.
	idField := crd.SpecFields["DBInstanceIdentifier"]
	require.NotNil(idField)
	if idField.IsImmutable() {
		assert.False(idField.IsLateInitialized())
		assert.Equal(immutabilityRule("dbInstanceIdentifier"), idField.ImmutabilityCELRule())
	}
}

// TestAttr_ImmutabilityCELRule_EscapesReservedWords asserts that a member whose
// JSON name collides with a CEL reserved word is addressed through Kubernetes'
// `__word__` property escaping. Without this the generated rule fails to
// compile when the CRD is created -- `spec.namespace` is a real immutable field
// on the s3 Bucket and s3tables Table resources.
func TestAttr_ImmutabilityCELRule_EscapesReservedWords(t *testing.T) {
	testCases := []struct {
		name     string
		attrName string
		expected string
	}{
		{"plain lower-camel name", "DeploymentType", immutabilityRule("deploymentType")},
		{"CEL reserved word", "Namespace", immutabilityRule("__namespace__")},
		{"CEL reserved word", "Function", immutabilityRule("__function__")},
		// names.New already escapes Go keywords with a trailing underscore, so
		// the handful of words that are reserved in both Go and CEL never reach
		// the CEL escaping. A single underscore is a legal CEL identifier
		// character and is left alone.
		{"Go and CEL reserved word", "Package", immutabilityRule("package_")},
		{"Go keyword name", "Type", immutabilityRule("type_")},
	}
	for _, tc := range testCases {
		t.Run(tc.name+"/"+tc.attrName, func(t *testing.T) {
			attr := model.NewAttr(names.New(tc.attrName), "*string", nil)
			assert.Equal(t, tc.expected, attr.ImmutabilityCELRule())
		})
	}
}

// TestAttr_ImmutabilityCELRule_UsesGoTagOverride asserts the rule follows a
// `go_tag` override, so it always names the property the CRD actually
// serializes.
func TestAttr_ImmutabilityCELRule_UsesGoTagOverride(t *testing.T) {
	attr := model.NewAttr(names.New("DeploymentType"), "*string", nil)
	attr.GoTag = "`json:\"deployment-type,omitempty\"`"

	assert.Equal(t, "deployment-type", attr.GetCRDJSONFieldName())
	// '-' is illegal in a CEL identifier and is escaped by Kubernetes.
	assert.Equal(t, immutabilityRule("deployment__dash__type"), attr.ImmutabilityCELRule())
}
