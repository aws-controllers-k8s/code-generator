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
