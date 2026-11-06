// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package chainsaw

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadInputs(t *testing.T, content string) (*Inputs, error) {
	t.Helper()
	path := filepath.Join(t.TempDir(), InputsFileName)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return LoadInputs(path)
}

func TestLoadInputsVariants(t *testing.T) {
	in, err := loadInputs(t, `
create: {cidrBlocks: [10.0.0.0/16]}
variants:
- name: disallowed
  description: Default security group rules removed at create
  create: {cidrBlocks: [10.0.0.0/16], disallowSecurityGroupDefaultRules: true}
  updates:
  - {enableDNSSupport: false}
checks:
  disallowSecurityGroupDefaultRules: {command: echo true, when: true}
`)
	require.NoError(t, err)
	require.Len(t, in.Variants, 1)
	assert.Equal(t, "disallowed", in.Variants[0].Name)
	assert.Equal(t, true, in.Variants[0].Create["disallowSecurityGroupDefaultRules"])
	assert.Len(t, in.Variants[0].Updates, 1)
	assert.Equal(t, true, in.Checks["disallowSecurityGroupDefaultRules"].When)
}

func TestLoadInputsBadVariants(t *testing.T) {
	for name, variants := range map[string]string{
		"no name":      `[{create: {a: 1}}]`,
		"bad name":     `[{name: Disallowed, create: {a: 1}}]`,
		"reserved":     `[{name: adoption, create: {a: 1}}]`,
		"update clash": `[{name: update-1, create: {a: 1}}]`,
		"duplicate":    `[{name: a, create: {a: 1}}, {name: a, create: {a: 2}}]`,
		"no create":    `[{name: a}]`,
	} {
		_, err := loadInputs(t, "create: {a: 1}\nvariants: "+variants+"\n")
		assert.Error(t, err, name)
	}
}

func TestLoadInputsValues(t *testing.T) {
	in, err := loadInputs(t, `
create: {imageID: $AMI_ID}
values:
- name: AMI_ID
  command: echo ami-1
- name: BUCKET
  command: echo b
  cleanup: echo rm
deleted: "Reservations[0].Instances[0].State.Name == 'terminated'"
`)
	require.NoError(t, err)
	require.Len(t, in.Values, 2)
	assert.Equal(t, "echo rm", in.Values[1].Cleanup)
	assert.Equal(t, "Reservations[0].Instances[0].State.Name == 'terminated'", in.Deleted)

	for name, values := range map[string]string{
		"lowercase": `[{name: ami, command: x}]`,
		"reserved":  `[{name: SUFFIX, command: x}]`,
		"duplicate": `[{name: A, command: x}, {name: A, command: y}]`,
		"no cmd":    `[{name: A}]`,
	} {
		_, err := loadInputs(t, "create: {a: 1}\nvalues: "+values+"\n")
		assert.Error(t, err, name)
	}
}

func TestLoadInputsPathChecks(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "inputs.yaml")
	require.NoError(t, os.WriteFile(p, []byte(`create: {name: x}
checks:
  routes:
    path: Routes
    exclude: {GatewayId: local}
    sortBy: DestinationCidrBlock
  entries:
    path: Entries
    sortBy: [Egress, RuleNumber]
`), 0o644))
	in, err := LoadInputs(p)
	require.NoError(t, err)
	assert.Equal(t, stringList{"DestinationCidrBlock"}, in.Checks["routes"].SortBy)
	assert.Equal(t, stringList{"Egress", "RuleNumber"}, in.Checks["entries"].SortBy)

	for _, bad := range []string{
		"checks: {f: {}}",
		"checks: {f: {path: A, command: b}}",
		"checks: {f: {command: b, sortBy: A}}",
	} {
		require.NoError(t, os.WriteFile(p, []byte("create: {name: x}\n"+bad+"\n"), 0o644))
		_, err := LoadInputs(p)
		assert.Error(t, err, bad)
	}
}

func TestLoadInputsAdoptionFields(t *testing.T) {
	in, err := loadInputs(t, `
create: {resourceID: $VPC_ID, resourceType: VPC}
adoptionFields: [resourceID, resourceType]
`)
	require.NoError(t, err)
	assert.Equal(t, stringList{"resourceID", "resourceType"}, in.AdoptionFields)

	_, err = loadInputs(t, `
create: {resourceID: $VPC_ID}
adoptionFields: resourceType
`)
	assert.ErrorContains(t, err, "create must set resourceType to a string")
}

func TestLoadInputsSuffixesNames(t *testing.T) {
	in, err := loadInputs(t, `
dependencies:
- apiVersion: ec2.services.k8s.aws/v1alpha1
  kind: SecurityGroup
  metadata: {name: dep-sg}
  spec: {name: dep-sg, description: d}
create: {name: ack-test, description: kept}
updates:
- {name: ack-renamed}
- {name: null}
variants:
- name: minimal
  create: {name: ack-minimal}
invalid:
  spec: {name: ack-invalid}
  message: m
values:
- {name: BUCKET, command: echo b}
`)
	require.NoError(t, err)
	assert.Equal(t, "ack-test-$SUFFIX", in.Create["name"])
	assert.Equal(t, "kept", in.Create["description"])
	assert.Equal(t, "ack-renamed-$SUFFIX", in.Updates[0]["name"])
	assert.Nil(t, in.Updates[1]["name"])
	assert.Equal(t, "ack-minimal-$SUFFIX", in.Variants[0].Create["name"])
	assert.Equal(t, "ack-invalid-$SUFFIX", in.Invalid.Spec["name"])
	assert.Equal(t, "dep-sg-$SUFFIX", in.Dependencies[0]["spec"].(map[string]interface{})["name"])
	// The CR name isn't an AWS name
	assert.Equal(t, "dep-sg", in.Dependencies[0]["metadata"].(map[string]interface{})["name"])

	in, err = loadInputs(t, `
create: {name: ack-$SUFFIX-x}
updates:
- {name: $BUCKET}
values:
- {name: BUCKET, command: echo b}
`)
	require.NoError(t, err)
	assert.Equal(t, "ack-$SUFFIX-x", in.Create["name"])
	assert.Equal(t, "$BUCKET", in.Updates[0]["name"])
}
