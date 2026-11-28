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
	"sigs.k8s.io/yaml"
)

// crd wraps spec properties in a one-version namespaced CRD.
func crd(spec map[string]v1.JSONSchemaProps) *v1.CustomResourceDefinition {
	c := &v1.CustomResourceDefinition{}
	c.APIVersion = "apiextensions.k8s.io/v1"
	c.Kind = "CustomResourceDefinition"
	c.Name = "widgets.example.services.k8s.aws"
	c.Spec.Group = "example.services.k8s.aws"
	c.Spec.Scope = v1.NamespaceScoped
	c.Spec.Names = v1.CustomResourceDefinitionNames{Kind: "Widget", ListKind: "WidgetList", Plural: "widgets",
		Singular: "widget"}
	s := root(obj(spec), obj(map[string]v1.JSONSchemaProps{"arn": str()}))
	c.Spec.Versions = []v1.CustomResourceDefinitionVersion{{
		Name: "v1alpha1", Served: true, Storage: true,
		Schema:       &v1.CustomResourceValidation{OpenAPIV3Schema: s},
		Subresources: &v1.CustomResourceSubresources{Status: &v1.CustomResourceSubresourceStatus{}},
	}}
	return c
}

func strMap(value v1.JSONSchemaProps) v1.JSONSchemaProps {
	return v1.JSONSchemaProps{Type: "object", AdditionalProperties: &v1.JSONSchemaPropsOrBool{Allows: true,
		Schema: &value}}
}

func uncomparedOnly(t *testing.T, old, new *v1.CustomResourceDefinition) Change {
	t.Helper()
	got := Uncompared("f.yaml", "Widget", old, new)
	require.Len(t, got, 1, "%+v", got)
	assert.Equal(t, "f.yaml", got[0].File)
	assert.Equal(t, "Widget", got[0].Kind)
	return got[0]
}

func TestUncomparedSchemaChanges(t *testing.T) {
	enumOf := func(vals ...string) v1.JSONSchemaProps {
		s := str()
		for _, v := range vals {
			s.Enum = append(s.Enum, v1.JSON{Raw: []byte(`"` + v + `"`)})
		}
		return s
	}
	maxLen := func(n int64) v1.JSONSchemaProps { s := str(); s.MaxLength = &n; return s }
	withDefault := func(s v1.JSONSchemaProps, raw string) v1.JSONSchemaProps {
		s.Default = &v1.JSON{Raw: []byte(raw)}
		return s
	}
	cel := str()
	cel.XValidations = v1.ValidationRules{{Rule: "self == oldSelf", Message: "immutable"}}
	list := func(item v1.JSONSchemaProps) v1.JSONSchemaProps {
		return v1.JSONSchemaProps{Type: "array", Items: &v1.JSONSchemaPropsOrArray{Schema: &item}}
	}

	cases := []struct {
		name          string
		old, new      v1.JSONSchemaProps
		bucket, field string
		path, keys    string
	}{
		{"enum added where none existed", str(), enumOf("a", "b"),
			BucketOther, "spec.x", ".properties.spec.properties.x", "enum"},
		{"x-kubernetes-validations added", str(), cel,
			BucketOther, "spec.x", ".properties.spec.properties.x", "x-kubernetes-validations"},
		{"maxLength changed", maxLen(5), maxLen(10),
			BucketOther, "spec.x", ".properties.spec.properties.x", "maxLength"},
		{"default only", withDefault(str(), `"a"`), withDefault(str(), `"b"`),
			BucketDefaultChanged, "spec.x", ".properties.spec.properties.x", "default"},
		{"default and maxLength", withDefault(maxLen(5), `"a"`), withDefault(maxLen(6), `"b"`),
			BucketOther, "spec.x", ".properties.spec.properties.x", "default,maxLength"},
		{"list item maxLength", list(maxLen(5)), list(maxLen(6)),
			BucketOther, "spec.x[]", ".properties.spec.properties.x.items", "maxLength"},
		{"map value maxLength", strMap(maxLen(5)), strMap(maxLen(6)),
			BucketOther, "spec.x{}", ".properties.spec.properties.x.additionalProperties", "maxLength"},
		{"map value added", obj(nil), strMap(str()),
			BucketOther, "spec.x", ".properties.spec.properties.x", "additionalProperties"},
		{"map value dropped", strMap(str()), obj(nil),
			BucketOther, "spec.x", ".properties.spec.properties.x", "additionalProperties"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			old := crd(map[string]v1.JSONSchemaProps{"x": tc.old})
			new := crd(map[string]v1.JSONSchemaProps{"x": tc.new})
			c := uncomparedOnly(t, old, new)
			assert.Equal(t, tc.bucket, c.Bucket)
			assert.Equal(t, tc.field, c.Field)
			assert.Equal(t, tc.path, c.Path)
			assert.Equal(t, ChangeUncomparedChanged, c.ChangeType)
			assert.Equal(t, tc.keys, c.New)
			assert.False(t, c.Breaking)
		})
	}
}

// kro never enters additionalProperties, so changes kro would report there
// come from Uncompared, classified as kro would.
func TestUncomparedMapValueTypeChange(t *testing.T) {
	old := crd(map[string]v1.JSONSchemaProps{"tags": strMap(str())})
	new := crd(map[string]v1.JSONSchemaProps{"tags": strMap(v1.JSONSchemaProps{Type: "integer"})})
	report, err := kro.CompareVersions(old.Spec.Versions, new.Spec.Versions)
	require.NoError(t, err)
	require.False(t, report.HasChanges(), "kro misses it")

	c := uncomparedOnly(t, old, new)
	assert.Equal(t, BucketTypeChanged, c.Bucket)
	assert.Equal(t, "spec.tags{}", c.Field)
	assert.Equal(t, ".properties.spec.properties.tags.additionalProperties.type", c.Path)
	assert.True(t, c.Breaking)
}

func TestUncomparedMapValueGainsMember(t *testing.T) {
	old := crd(map[string]v1.JSONSchemaProps{"m": strMap(obj(map[string]v1.JSONSchemaProps{"a": str()}))})
	new := crd(map[string]v1.JSONSchemaProps{"m": strMap(obj(map[string]v1.JSONSchemaProps{"a": str(), "b": str()}))})
	c := uncomparedOnly(t, old, new)
	assert.Equal(t, BucketAddedOptionalSpec, c.Bucket)
	assert.Equal(t, "spec.m{}.b", c.Field)
	assert.Equal(t, string(kro.PropertyAdded), c.ChangeType)
}

// kro reports a type change and stops there; Uncompared adds nothing.
func TestUncomparedSkipsTypeChangedNode(t *testing.T) {
	n := int64(5)
	old := crd(map[string]v1.JSONSchemaProps{"x": {Type: "string", MaxLength: &n}})
	new := crd(map[string]v1.JSONSchemaProps{"x": {Type: "integer"}})
	assert.Empty(t, Uncompared("f.yaml", "Widget", old, new))
}

func TestUncomparedCRDLevelChanges(t *testing.T) {
	cases := []struct {
		name     string
		mutate   func(c *v1.CustomResourceDefinition)
		path     string
		breaking bool
		keys     string
	}{
		{"scope", func(c *v1.CustomResourceDefinition) { c.Spec.Scope = v1.ClusterScoped }, "spec.scope", true, "scope"},
		{"group", func(c *v1.CustomResourceDefinition) { c.Spec.Group = "other.k8s.aws" }, "spec.group", true, "group"},
		{"kind", func(c *v1.CustomResourceDefinition) { c.Spec.Names.Kind = "Gadget" }, "spec.names", true, "kind"},
		{"plural", func(c *v1.CustomResourceDefinition) { c.Spec.Names.Plural = "widgetz" }, "spec.names", true, "plural"},
		{"short names", func(c *v1.CustomResourceDefinition) { c.Spec.Names.ShortNames = []string{"wd"} },
			"spec.names", false, "shortNames"},
		{"conversion", func(c *v1.CustomResourceDefinition) {
			c.Spec.Conversion = &v1.CustomResourceConversion{Strategy: v1.NoneConverter}
		}, "spec.conversion", false, "conversion"},
		{"preserveUnknownFields", func(c *v1.CustomResourceDefinition) { c.Spec.PreserveUnknownFields = true },
			"spec.preserveUnknownFields", false, "preserveUnknownFields"},
		{"printer column", func(c *v1.CustomResourceDefinition) {
			c.Spec.Versions[0].AdditionalPrinterColumns = []v1.CustomResourceColumnDefinition{
				{Name: "ARN", Type: "string", JSONPath: ".status.arn"}}
		}, "version.additionalPrinterColumns", false, "additionalPrinterColumns"},
		{"storage", func(c *v1.CustomResourceDefinition) { c.Spec.Versions[0].Storage = false },
			"version.storage", false, "storage"},
		{"deprecated", func(c *v1.CustomResourceDefinition) { c.Spec.Versions[0].Deprecated = true },
			"version.deprecated", false, "deprecated"},
		{"deprecation warning", func(c *v1.CustomResourceDefinition) {
			w := "gone soon"
			c.Spec.Versions[0].DeprecationWarning = &w
		}, "version.deprecationWarning", false, "deprecationWarning"},
		{"selectable fields", func(c *v1.CustomResourceDefinition) {
			c.Spec.Versions[0].SelectableFields = []v1.SelectableField{{JSONPath: ".spec.x"}}
		}, "version.selectableFields", false, "selectableFields"},
		{"scale", func(c *v1.CustomResourceDefinition) {
			c.Spec.Versions[0].Subresources.Scale = &v1.CustomResourceSubresourceScale{
				SpecReplicasPath: ".spec.replicas", StatusReplicasPath: ".status.replicas"}
		}, "version.subresources.scale", false, "scale"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			old := crd(map[string]v1.JSONSchemaProps{"x": str()})
			new := old.DeepCopy()
			tc.mutate(new)
			c := uncomparedOnly(t, old, new)
			assert.Equal(t, tc.path, c.Path)
			assert.Equal(t, tc.path, c.Field)
			assert.Equal(t, BucketOther, c.Bucket)
			assert.Equal(t, ChangeUncomparedChanged, c.ChangeType)
			assert.Equal(t, tc.breaking, c.Breaking)
			assert.Equal(t, tc.keys, c.New)
		})
	}
}

// kro reports a dropped status subresource, but not a new one.
func TestUncomparedStatusSubresource(t *testing.T) {
	without := crd(map[string]v1.JSONSchemaProps{"x": str()})
	without.Spec.Versions[0].Subresources = nil
	with := crd(map[string]v1.JSONSchemaProps{"x": str()})

	c := uncomparedOnly(t, without, with)
	assert.Equal(t, "version.subresources.status", c.Path)
	assert.Equal(t, "status", c.New)
	assert.Equal(t, BucketOther, c.Bucket)
	assert.False(t, c.Breaking)

	assert.Empty(t, Uncompared("f.yaml", "Widget", with, without), "kro reports the removal")
}

func TestUncomparedNoOp(t *testing.T) {
	n := int64(3)
	spec := map[string]v1.JSONSchemaProps{
		"name": {Type: "string", MaxLength: &n, Default: &v1.JSON{Raw: []byte(`"a"`)}},
		"tags": strMap(str()),
		"list": {Type: "array", Items: &v1.JSONSchemaPropsOrArray{Schema: &v1.JSONSchemaProps{Type: "string"}}},
	}
	old := crd(spec)
	old.Spec.Versions[0].AdditionalPrinterColumns = []v1.CustomResourceColumnDefinition{
		{Name: "ARN", Type: "string", JSONPath: ".status.arn"}}
	data, err := yaml.Marshal(old)
	require.NoError(t, err)
	var trip v1.CustomResourceDefinition
	require.NoError(t, yaml.Unmarshal(data, &trip))

	assert.Empty(t, Uncompared("f.yaml", "Widget", old, &trip))
	assert.Empty(t, Uncompared("f.yaml", "Widget", old, old.DeepCopy()))
}

// Changes kro compares are not reported again.
func TestUncomparedIgnoresWhatKroCompares(t *testing.T) {
	old := crd(map[string]v1.JSONSchemaProps{"x": str(), "y": str()})
	new := crd(map[string]v1.JSONSchemaProps{
		"x": {Type: "string", Description: "d", Pattern: "^a$"},
		"z": str(),
	})
	new.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"] = func() v1.JSONSchemaProps {
		s := new.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"]
		s.Required = []string{"x"}
		return s
	}()
	assert.Empty(t, Uncompared("f.yaml", "Widget", old, new))
}

func TestFieldPathMapValues(t *testing.T) {
	assert.Equal(t, "spec.tags{}", FieldPath(kro.Change{Path: ".properties.spec.properties.tags.additionalProperties"}))
	assert.Equal(t, "spec.tags{}.key",
		FieldPath(kro.Change{Path: ".properties.spec.properties.tags.additionalProperties.properties.key.type"}))
	assert.Equal(t, "spec.x[]{}",
		FieldPath(kro.Change{Path: ".properties.spec.properties.x.items.additionalProperties"}))
}

// A required name that is a property on neither side: kro checks only
// existing properties, and every object lacking the name stops validating.
func TestUncomparedRequiredNameWithoutProperty(t *testing.T) {
	old := crd(map[string]v1.JSONSchemaProps{"x": str()})
	new := crd(map[string]v1.JSONSchemaProps{"x": str()})
	spec := new.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"]
	spec.Required = []string{"ghost"}
	new.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"] = spec
	report, err := kro.CompareVersions(old.Spec.Versions, new.Spec.Versions)
	require.NoError(t, err)
	require.False(t, report.HasChanges(), "kro misses it")

	c := uncomparedOnly(t, old, new)
	assert.Equal(t, BucketValidationTightened, c.Bucket)
	assert.Equal(t, string(kro.RequiredAdded), c.ChangeType)
	assert.Equal(t, ".properties.spec.required", c.Path)
	assert.Equal(t, "spec.ghost", c.Field)
	assert.Equal(t, "ghost", c.New)
	assert.True(t, c.Breaking)

	assert.Empty(t, Uncompared("f.yaml", "Widget", new, new.DeepCopy()), "already required")
}

func TestUncomparedMapValueOpenParent(t *testing.T) {
	yes := true
	value := func(props map[string]v1.JSONSchemaProps) v1.JSONSchemaProps {
		s := obj(props)
		s.XPreserveUnknownFields = &yes
		return s
	}
	old := crd(map[string]v1.JSONSchemaProps{"m": strMap(value(nil))})
	new := crd(map[string]v1.JSONSchemaProps{"m": strMap(value(map[string]v1.JSONSchemaProps{"b": str()}))})
	c := uncomparedOnly(t, old, new)
	assert.Equal(t, BucketValidationTightened, c.Bucket)
	assert.Equal(t, "spec.m{}.b", c.Field)
}

func TestUncomparedComparesSerializedBytes(t *testing.T) {
	// enum: [] serializes like no enum.
	empty := str()
	empty.Enum = []v1.JSON{}
	assert.Empty(t, Uncompared("f.yaml", "Widget",
		crd(map[string]v1.JSONSchemaProps{"x": str()}), crd(map[string]v1.JSONSchemaProps{"x": empty})))

	// Integers beyond float64 precision still differ, under the right key.
	big := func(n int64, def string) v1.JSONSchemaProps {
		s := v1.JSONSchemaProps{Type: "integer", Maximum: nil}
		s.MaxLength = &n
		s.Default = &v1.JSON{Raw: []byte(def)}
		return s
	}
	c := uncomparedOnly(t, crd(map[string]v1.JSONSchemaProps{"x": big(9007199254740992, "9007199254740992")}),
		crd(map[string]v1.JSONSchemaProps{"x": big(9007199254740993, "9007199254740992")}))
	assert.Equal(t, "maxLength", c.New)
	c = uncomparedOnly(t, crd(map[string]v1.JSONSchemaProps{"x": big(1, "9007199254740992")}),
		crd(map[string]v1.JSONSchemaProps{"x": big(1, "9007199254740993")}))
	assert.Equal(t, "default", c.New)
	assert.Equal(t, BucketDefaultChanged, c.Bucket)
}

// What cannot be serialized cannot be shown equal.
func TestUncomparedMarshalFailureIsAChange(t *testing.T) {
	bad := str()
	bad.Default = &v1.JSON{Raw: []byte("{not json")}
	old := crd(map[string]v1.JSONSchemaProps{"x": bad})
	c := uncomparedOnly(t, old, old.DeepCopy())
	assert.Equal(t, BucketOther, c.Bucket)
	assert.Equal(t, "spec.x", c.Field)
	assert.Contains(t, c.New, "cannot serialize")
}

// kro reports a version no longer served, not one newly served.
func TestUncomparedServedTurnedOn(t *testing.T) {
	off := crd(map[string]v1.JSONSchemaProps{"x": str()})
	off.Spec.Versions[0].Served = false
	on := crd(map[string]v1.JSONSchemaProps{"x": str()})
	c := uncomparedOnly(t, off, on)
	assert.Equal(t, "version.served", c.Path)
	assert.Equal(t, BucketOther, c.Bucket)
	assert.Equal(t, ChangeUncomparedChanged, c.ChangeType)
	assert.False(t, c.Breaking)
	assert.Empty(t, Uncompared("f.yaml", "Widget", on, off), "kro reports it")
}
