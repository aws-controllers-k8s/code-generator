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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ackgenerate "github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

// TestCELOnTypeDef verifies that custom_cel_rules configured on a
// nested field are rendered as +kubebuilder:validation:XValidation markers in
// the generated types.go (via the type_def template include).
func TestCELOnTypeDef(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-cel-rules.yaml",
	})

	ts, err := ackgenerate.APIs(g, []string{testutil.TemplatesBasePath(t)})
	require.NoError(err)
	require.NoError(ts.Execute())

	typesGo, ok := ts.Executed()["types.go"]
	require.True(ok, "types.go not found in executed templates")

	output := typesGo.String()
	assert.True(
		strings.Contains(output, `// +kubebuilder:validation:XValidation:rule="self.startsWith('https://')",message="Issuer must be an HTTPS URL"`),
		"expected XValidation marker for Issuer CEL rule in types.go",
	)
}

// TestCELOnSpec verifies that custom_cel_rules configured at the
// resource level are rendered as +kubebuilder:validation:XValidation markers
// on the generated CRD Spec struct (via the crd.go template).
func TestCELOnSpec(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "route53")

	ts, err := ackgenerate.APIs(g, []string{testutil.TemplatesBasePath(t)})
	require.NoError(err)
	require.NoError(ts.Execute())

	hostedZoneGo, ok := ts.Executed()["hosted_zone.go"]
	require.True(ok, "hosted_zone.go not found in executed templates")

	output := hostedZoneGo.String()
	assert.True(
		strings.Contains(output, `// +kubebuilder:validation:XValidation:rule="!has(self.hostedZoneConfig) || !self.hostedZoneConfig.privateZone || has(self.vpc)",message="spec.vpc is required for private hosted zones"`),
		"expected XValidation marker for first HostedZone CEL rule",
	)
	assert.True(
		strings.Contains(output, `// +kubebuilder:validation:XValidation:rule="size(self.name) > 0"`),
		"expected XValidation marker for second HostedZone CEL rule (no message)",
	)
}

// TestCELOnField verifies that custom_cel_rules configured on a
// top-level spec field are rendered as +kubebuilder:validation:XValidation
// markers on the individual field in the generated CRD Spec struct.
func TestCELOnField(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "route53")

	ts, err := ackgenerate.APIs(g, []string{testutil.TemplatesBasePath(t)})
	require.NoError(err)
	require.NoError(ts.Execute())

	recordSetGo, ok := ts.Executed()["record_set.go"]
	require.True(ok, "record_set.go not found in executed templates")

	output := recordSetGo.String()
	assert.True(
		strings.Contains(output, `// +kubebuilder:validation:XValidation:rule="self.endsWith('.')",message="DNS name must end with a dot"`),
		"expected XValidation marker for Name field CEL rule in record_set.go",
	)
}
