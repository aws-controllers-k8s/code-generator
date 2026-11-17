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
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func str() v1.JSONSchemaProps { return v1.JSONSchemaProps{Type: "string"} }

func obj(props map[string]v1.JSONSchemaProps, required ...string) v1.JSONSchemaProps {
	return v1.JSONSchemaProps{Type: "object", Properties: props, Required: required}
}

// root builds a CRD openAPIV3Schema with the given spec and status.
func root(spec, status v1.JSONSchemaProps) *v1.JSONSchemaProps {
	r := obj(map[string]v1.JSONSchemaProps{"spec": spec, "status": status})
	return &r
}

// only classifies old→new and requires exactly one change.
func only(t *testing.T, old, new *v1.JSONSchemaProps) Change {
	t.Helper()
	got := Classify("f.yaml", "Widget", kro.Compare(old, new))
	require.Len(t, got, 1, "%+v", got)
	return got[0]
}

func TestClassifyBuckets(t *testing.T) {
	base := func() v1.JSONSchemaProps { return obj(map[string]v1.JSONSchemaProps{"name": str()}) }
	status := obj(map[string]v1.JSONSchemaProps{"arn": str()})

	withSpec := func(extra string, s v1.JSONSchemaProps, required ...string) v1.JSONSchemaProps {
		o := base()
		o.Properties[extra] = s
		o.Required = required
		return o
	}

	cases := []struct {
		name     string
		old, new *v1.JSONSchemaProps
		bucket   string
		field    string
		breaking bool
	}{
		{"optional spec add", root(base(), status), root(withSpec("size", str()), status),
			BucketAddedOptionalSpec, "spec.size", false},
		{"status add", root(base(), status),
			root(base(), obj(map[string]v1.JSONSchemaProps{"arn": str(), "state": str()})),
			BucketAddedStatus, "status.state", false},
		{"required add", root(base(), status), root(withSpec("size", str(), "size"), status),
			BucketAddedRequired, "spec.size", true},
		{"type change", root(base(), status),
			root(obj(map[string]v1.JSONSchemaProps{"name": {Type: "integer"}}), status),
			BucketTypeChanged, "spec.name", true},
		{"removed", root(withSpec("size", str()), status), root(base(), status),
			BucketRemoved, "spec.size", true},
		{"existing field made required", root(base(), status), root(withSpec("name", str(), "name"), status),
			BucketValidationTightened, "spec.name", true},
		{"required dropped", root(withSpec("name", str(), "name"), status), root(base(), status),
			BucketValidationRelaxed, "spec.name", false},
		{"pattern added", root(base(), status),
			root(obj(map[string]v1.JSONSchemaProps{"name": {Type: "string", Pattern: "^a"}}), status),
			BucketValidationTightened, "spec.name", true},
		{"pattern removed",
			root(obj(map[string]v1.JSONSchemaProps{"name": {Type: "string", Pattern: "^a"}}), status),
			root(base(), status), BucketValidationRelaxed, "spec.name", false},
		{"description", root(base(), status),
			root(obj(map[string]v1.JSONSchemaProps{"name": {Type: "string", Description: "new"}}), status),
			BucketDescriptionOnly, "spec.name", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := only(t, tc.old, tc.new)
			assert.Equal(t, tc.bucket, c.Bucket)
			assert.Equal(t, tc.field, c.Field)
			assert.Equal(t, tc.breaking, c.Breaking)
			assert.Equal(t, "f.yaml", c.File)
			assert.Equal(t, "Widget", c.Kind)
		})
	}
}

func TestClassifyEnumKeepsDirection(t *testing.T) {
	enum := func(vals ...string) *v1.JSONSchemaProps {
		s := str()
		for _, v := range vals {
			s.Enum = append(s.Enum, v1.JSON{Raw: []byte(`"` + v + `"`)})
		}
		return root(obj(map[string]v1.JSONSchemaProps{"mode": *s.DeepCopy()}), obj(nil))
	}
	expanded := only(t, enum("a"), enum("a", "b"))
	assert.Equal(t, BucketEnumChanged, expanded.Bucket)
	assert.False(t, expanded.Breaking)
	assert.Equal(t, "spec.mode", expanded.Field)

	restricted := only(t, enum("a", "b"), enum("a"))
	assert.Equal(t, BucketEnumChanged, restricted.Bucket)
	assert.True(t, restricted.Breaking)
}

func TestFieldPath(t *testing.T) {
	cases := []struct {
		change kro.Change
		want   string
	}{
		{kro.Change{Path: ".properties.spec.properties.tags.items.properties.key", ChangeType: kro.PropertyAdded},
			"spec.tags[].key"},
		{kro.Change{Path: ".properties.spec.properties.tags.items", ChangeType: kro.PropertyAdded}, "spec.tags[]"},
		{kro.Change{Path: ".properties.spec.properties.name.type", ChangeType: kro.TypeChanged}, "spec.name"},
		{kro.Change{Path: ".properties.spec.required", ChangeType: kro.RequiredAdded, NewValue: "name"},
			"spec.name"},
		{kro.Change{Path: ".properties.spec.required", ChangeType: kro.RequiredRemoved, OldValue: "name"},
			"spec.name"},
		// A property literally named like a schema keyword is still a name.
		{kro.Change{Path: ".properties.spec.properties.items.properties.required", ChangeType: kro.PropertyAdded},
			"spec.items.required"},
		{kro.Change{Path: ".description", ChangeType: kro.DescriptionChanged}, ""},
		{kro.Change{Path: "version.served", ChangeType: kro.TypeChanged}, "version.served"},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.want, FieldPath(tc.change), tc.change.Path)
	}
}

func TestClassifyVersionRenameIsTypeChanged(t *testing.T) {
	r, err := kro.CompareVersions(
		[]v1.CustomResourceDefinitionVersion{{Name: "v1alpha1"}},
		[]v1.CustomResourceDefinitionVersion{{Name: "v1beta1"}},
	)
	require.NoError(t, err)
	got := Classify("f.yaml", "Widget", r)
	require.Len(t, got, 1)
	assert.Equal(t, BucketTypeChanged, got[0].Bucket)
	assert.Equal(t, "version", got[0].Field)
	assert.True(t, got[0].Breaking)
}

func TestClassifyUnknownChangeTypeIsOther(t *testing.T) {
	got := Classify("f.yaml", "Widget", &kro.Report{NonBreakingChanges: []kro.Change{
		{Path: ".properties.spec.properties.name.format", ChangeType: "FORMAT_CHANGED"},
	}})
	require.Len(t, got, 1)
	assert.Equal(t, BucketOther, got[0].Bucket)
	assert.False(t, got[0].Breaking)
}

func TestCRDFileChanges(t *testing.T) {
	added := CRDAdded("a.yaml", "Gizmo")
	assert.Equal(t, Change{File: "a.yaml", Kind: "Gizmo", Bucket: BucketCRDAdded, ChangeType: ChangeCRDAdded}, added)
	removed := CRDRemoved("r.yaml", "Gadget")
	assert.Equal(t, Change{File: "r.yaml", Kind: "Gadget", Bucket: BucketCRDRemoved, ChangeType: ChangeCRDRemoved,
		Breaking: true}, removed)
}
