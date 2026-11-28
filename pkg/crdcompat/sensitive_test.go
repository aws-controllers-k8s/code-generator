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

package crdcompat

import (
	"testing"

	kro "github.com/kubernetes-sigs/kro/pkg/graph/crd/compat"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func secretRef() v1.JSONSchemaProps {
	return obj(map[string]v1.JSONSchemaProps{"key": str(), "name": str(), "namespace": str()}, "key")
}

func fields(cs []SensitiveCandidate) []string {
	out := []string{}
	for _, c := range cs {
		out = append(out, c.Field)
	}
	return out
}

func TestSensitiveFieldsWalksAddedSubtrees(t *testing.T) {
	old := root(obj(map[string]v1.JSONSchemaProps{"name": str()}), obj(nil))
	new := root(obj(map[string]v1.JSONSchemaProps{
		"name":     str(),
		"apiToken": str(),
		"auth": obj(map[string]v1.JSONSchemaProps{
			"userName": str(),
			"Password": str(),
			"keys": {Type: "array", Items: &v1.JSONSchemaPropsOrArray{
				Schema: &v1.JSONSchemaProps{Type: "object", Properties: map[string]v1.JSONSchemaProps{
					"privateKey": str(), "label": str()}}}},
		}),
		"masterUserPassword": secretRef(),
		"plain":              str(),
	}), obj(map[string]v1.JSONSchemaProps{"credentialsArn": str()}))

	changes := Classify("f.yaml", "Widget", kro.Compare(old, new))
	got := SensitiveFields("f.yaml", "Widget", new, changes)

	assert.ElementsMatch(t, []string{
		"spec.apiToken", "spec.auth.Password", "spec.auth.keys[].privateKey", "status.credentialsArn",
	}, fields(got), "a SecretKeyReference already has its is_secret decision")
	for _, c := range got {
		assert.Equal(t, "f.yaml", c.File)
		assert.Equal(t, "Widget", c.Kind)
	}
}

func TestSensitiveFieldsIgnoresExistingFields(t *testing.T) {
	both := root(obj(map[string]v1.JSONSchemaProps{"secretName": str()}), obj(nil))
	changed := root(obj(map[string]v1.JSONSchemaProps{"secretName": {Type: "string", Description: "d"}}), obj(nil))
	changes := Classify("f.yaml", "Widget", kro.Compare(both, changed))
	assert.Empty(t, SensitiveFields("f.yaml", "Widget", changed, changes))
}

func TestNewCRDSensitiveFields(t *testing.T) {
	schema := root(obj(map[string]v1.JSONSchemaProps{"clientSecret": str(), "name": str()}),
		obj(map[string]v1.JSONSchemaProps{"sessionToken": str()}))
	assert.ElementsMatch(t, []string{"spec.clientSecret", "status.sessionToken"},
		fields(NewCRDSensitiveFields("n.yaml", "Gizmo", schema)))
	assert.Empty(t, NewCRDSensitiveFields("n.yaml", "Gizmo", nil))
}

func mapOf(value v1.JSONSchemaProps) v1.JSONSchemaProps {
	return v1.JSONSchemaProps{Type: "object", AdditionalProperties: &v1.JSONSchemaPropsOrBool{Allows: true,
		Schema: &value}}
}

func TestSensitiveFieldsWalksMapValues(t *testing.T) {
	old := root(obj(map[string]v1.JSONSchemaProps{"name": str()}), obj(nil))
	new := root(obj(map[string]v1.JSONSchemaProps{
		"name":   str(),
		"labels": mapOf(obj(map[string]v1.JSONSchemaProps{"secretValue": str(), "plain": str()})),
		"tokens": mapOf(str()),
	}), obj(nil))
	changes := Classify("f.yaml", "Widget", kro.Compare(old, new))
	assert.ElementsMatch(t, []string{"spec.labels{}.secretValue", "spec.tokens"},
		fields(SensitiveFields("f.yaml", "Widget", new, changes)))
	assert.ElementsMatch(t, []string{"spec.labels{}.secretValue", "spec.tokens"},
		fields(NewCRDSensitiveFields("f.yaml", "Widget", new)))
}

// A member added inside an existing map's values comes from Uncompared.
func TestSensitiveFieldsInExistingMapValues(t *testing.T) {
	old := crd(map[string]v1.JSONSchemaProps{"m": mapOf(obj(map[string]v1.JSONSchemaProps{"a": str()}))})
	new := crd(map[string]v1.JSONSchemaProps{"m": mapOf(obj(map[string]v1.JSONSchemaProps{
		"a": str(), "password": str()}))})
	changes := Uncompared("f.yaml", "Widget", old, new)
	assert.Equal(t, []string{"spec.m{}.password"},
		fields(SensitiveFields("f.yaml", "Widget", new.Spec.Versions[0].Schema.OpenAPIV3Schema, changes)))
}

// Only key, name and namespace all strings make a SecretKeyReference.
func TestSensitiveSecretKeyReferenceNeedsStrings(t *testing.T) {
	lookalike := obj(map[string]v1.JSONSchemaProps{"key": {Type: "integer"}, "name": str(), "namespace": str()})
	schema := root(obj(map[string]v1.JSONSchemaProps{"dbPassword": lookalike, "adminPassword": secretRef()}), obj(nil))
	assert.Equal(t, []string{"spec.dbPassword"}, fields(NewCRDSensitiveFields("n.yaml", "Gizmo", schema)))
}
