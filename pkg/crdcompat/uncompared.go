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
	"bytes"
	"encoding/json"
	"sort"
	"strings"

	kro "github.com/kubernetes-sigs/kro/pkg/graph/crd/compat"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// ChangeUncomparedChanged is a difference kro's compat library does not look
// at: a schema keyword other than the ones kro compares (maxLength, default,
// x-kubernetes-validations, an enum where one side has none, ...), or a CRD
// or version setting (scope, names, printer columns, ...).
const ChangeUncomparedChanged = "UNCOMPARED_CHANGED"

// Uncompared returns what changed between two versions of a CRD file that
// kro's CompareVersions does not report, so a gate never mistakes such a
// change for no change. kro compares a schema node's description, type,
// pattern, properties, required list, items, and enum (only when both sides
// have one), and a version's name, served flag and status subresource.
//
// Every other difference at a schema node is one UNCOMPARED_CHANGED change at
// that node, naming the differing keys in New; it lands in default_changed
// when default is the only one, else in other. kro never enters
// additionalProperties, so where both sides have a map value schema that
// subtree is also run through kro, and its changes are classified as kro's.
// A newly required name that is a property on neither side, which kro does
// not check, is a breaking REQUIRED_ADDED (validation_tightened).
// CRD-level differences are one change per spec or version key (New names
// the key, or for spec.names the members that differ), all in
// other, and breaking when they change the CRD's identity (group, scope,
// names.kind, names.plural).
//
// The CRD's metadata (annotations, labels, e.g. the controller-gen version
// annotation) is deliberately not compared: it changes on every
// regeneration and does not change what the API server accepts.
func Uncompared(file, kind string, old, new *v1.CustomResourceDefinition) []Change {
	u := uncompared{file: file, kind: kind}
	u.crd(old, new)
	if len(old.Spec.Versions) == 1 && len(new.Spec.Versions) == 1 {
		ov, nv := old.Spec.Versions[0], new.Spec.Versions[0]
		// kro reports a renamed version and stops; so does this.
		if ov.Name == nv.Name {
			u.version(ov, nv)
			if ov.Schema != nil && nv.Schema != nil {
				u.schema("", ov.Schema.OpenAPIV3Schema, nv.Schema.OpenAPIV3Schema)
				tightenUnderOpenParents(u.out, ov.Schema.OpenAPIV3Schema)
			}
		}
	}
	return u.out
}

type uncompared struct {
	file, kind string
	out        []Change
}

// emit records a change at path; keys names what differs there.
func (u *uncompared) emit(path, bucket, keys string, breaking bool) {
	u.out = append(u.out, Change{
		File: u.file, Kind: u.kind, Bucket: bucket,
		Field: FieldPath(kro.Change{Path: path}), Path: path,
		ChangeType: ChangeUncomparedChanged, New: keys, Breaking: breaking,
	})
}

// crd compares the CRD spec outside its versions.
func (u *uncompared) crd(old, new *v1.CustomResourceDefinition) {
	oldSpec, newSpec := old.Spec, new.Spec
	oldSpec.Versions, newSpec.Versions = nil, nil
	keys, err := jsonDiff(oldSpec, newSpec)
	if err != nil {
		u.emit("spec", BucketOther, cannotSerialize(err), false)
		return
	}
	for _, key := range keys {
		what := key
		breaking := false
		switch key {
		case "group", "scope":
			breaking = true
		case "names":
			breaking = old.Spec.Names.Kind != new.Spec.Names.Kind || old.Spec.Names.Plural != new.Spec.Names.Plural
			// names serialized above, so it serializes here.
			members, _ := jsonDiff(old.Spec.Names, new.Spec.Names)
			what = strings.Join(members, ",")
		}
		u.emit("spec."+key, BucketOther, what, breaking)
	}
}

// version compares a version's settings other than the ones kro compares.
func (u *uncompared) version(old, new v1.CustomResourceDefinitionVersion) {
	if !old.Served && new.Served {
		// kro reports only served going away.
		u.emit("version.served", BucketOther, "served", false)
	}
	osub, nsub := old.Subresources, new.Subresources
	for _, v := range []*v1.CustomResourceDefinitionVersion{&old, &new} {
		v.Name, v.Served, v.Schema, v.Subresources = "", false, nil, nil
	}
	keys, err := jsonDiff(old, new)
	if err != nil {
		u.emit("version", BucketOther, cannotSerialize(err), false)
	}
	for _, key := range keys {
		u.emit("version."+key, BucketOther, key, false)
	}

	var oldSub, newSub v1.CustomResourceSubresources
	if osub != nil {
		oldSub = *osub
	}
	if nsub != nil {
		newSub = *nsub
	}
	if oldSub.Status != nil {
		// kro reports the status subresource going away, and present on
		// both sides it has no settings.
		oldSub.Status, newSub.Status = nil, nil
	}
	keys, err = jsonDiff(oldSub, newSub)
	if err != nil {
		u.emit("version.subresources", BucketOther, cannotSerialize(err), false)
	}
	for _, key := range keys {
		u.emit("version.subresources."+key, BucketOther, key, false)
	}
}

// schema walks old and new in parallel and reports, at each node, what
// differs besides what kro compares.
func (u *uncompared) schema(path string, old, new *v1.JSONSchemaProps) {
	if old == nil || new == nil || old.Type != new.Type {
		// kro reports an added or removed node, or a type change, and does
		// not look further.
		return
	}
	o, n := *old, *new
	for _, s := range []*v1.JSONSchemaProps{&o, &n} {
		s.Description, s.Type, s.Pattern, s.Properties, s.Required = "", "", "", nil, nil
	}
	if len(old.Enum) > 0 && len(new.Enum) > 0 {
		o.Enum, n.Enum = nil, nil
	}
	// kro compares items, or reports them added or removed, only for arrays.
	arrays := old.Type == "array" && new.Type == "array"
	if arrays {
		o.Items, n.Items = itemsWithoutSchema(old.Items), itemsWithoutSchema(new.Items)
	}
	intoMap := old.AdditionalProperties != nil && new.AdditionalProperties != nil &&
		old.AdditionalProperties.Allows == new.AdditionalProperties.Allows &&
		old.AdditionalProperties.Schema != nil && new.AdditionalProperties.Schema != nil
	if intoMap {
		o.AdditionalProperties, n.AdditionalProperties = nil, nil
	}
	keys, err := jsonDiff(o, n)
	switch {
	case err != nil:
		u.emit(path, BucketOther, cannotSerialize(err), false)
	case len(keys) == 1 && keys[0] == "default":
		u.emit(path, BucketDefaultChanged, "default", false)
	case len(keys) > 0:
		u.emit(path, BucketOther, strings.Join(keys, ","), false)
	}

	// kro checks a newly required name only when the property existed, and
	// an added property only when it is a property now. A name that is a
	// property on neither side is still required of every object.
	oldRequired := map[string]bool{}
	for _, name := range old.Required {
		oldRequired[name] = true
	}
	ghosts := map[string]bool{}
	for _, name := range new.Required {
		_, inOld := old.Properties[name]
		_, inNew := new.Properties[name]
		if !oldRequired[name] && !inOld && !inNew && !ghosts[name] {
			ghosts[name] = true
			u.out = append(u.out, newChange(u.file, u.kind,
				kro.Change{Path: path + ".required", ChangeType: kro.RequiredAdded, NewValue: name}, true))
		}
	}

	names := make([]string, 0, len(old.Properties))
	for name := range old.Properties {
		if _, ok := new.Properties[name]; ok {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	for _, name := range names {
		op, np := old.Properties[name], new.Properties[name]
		u.schema(path+".properties."+name, &op, &np)
	}
	if arrays && old.Items != nil && new.Items != nil {
		u.schema(path+".items", old.Items.Schema, new.Items.Schema)
	}
	if intoMap {
		mapPath := path + ".additionalProperties"
		oldValue, newValue := old.AdditionalProperties.Schema, new.AdditionalProperties.Schema
		r := kro.Compare(oldValue, newValue)
		for _, c := range r.BreakingChanges {
			c.Path = mapPath + c.Path
			u.out = append(u.out, newChange(u.file, u.kind, c, true))
		}
		for _, c := range r.NonBreakingChanges {
			c.Path = mapPath + c.Path
			u.out = append(u.out, newChange(u.file, u.kind, c, false))
		}
		u.schema(mapPath, oldValue, newValue)
	}
}

// itemsWithoutSchema is items minus the schema kro compares, or nil.
func itemsWithoutSchema(items *v1.JSONSchemaPropsOrArray) *v1.JSONSchemaPropsOrArray {
	if items == nil || items.JSONSchemas == nil {
		return nil
	}
	return &v1.JSONSchemaPropsOrArray{JSONSchemas: items.JSONSchemas}
}

// jsonDiff compares a and b as serialized and returns the sorted top-level
// keys whose JSON differs, or none when the bytes are identical. Comparing
// bytes, not decoded values, keeps nil and empty alike (both omitted) and
// integers beyond float64 precision apart. An error means a or b cannot be
// serialized, so they cannot be shown equal.
func jsonDiff(a, b any) ([]string, error) {
	aj, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	bj, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(aj, bj) {
		return nil, nil
	}
	var am, bm map[string]json.RawMessage
	if err := json.Unmarshal(aj, &am); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bj, &bm); err != nil {
		return nil, err
	}
	var keys []string
	for k, av := range am {
		if bv, ok := bm[k]; !ok || !bytes.Equal(av, bv) {
			keys = append(keys, k)
		}
	}
	for k := range bm {
		if _, ok := am[k]; !ok {
			keys = append(keys, k)
		}
	}
	if len(keys) == 0 {
		// Unreachable for objects; never call different bytes equal.
		keys = []string{"(serialized form differs)"}
	}
	sort.Strings(keys)
	return keys, nil
}

func cannotSerialize(err error) string { return "(cannot serialize: " + err.Error() + ")" }

// Description is the change in words, as kro words its own.
func (c Change) Description() string {
	switch c.ChangeType {
	case ChangeUncomparedChanged:
		return "Uncompared change: " + c.New
	case ChangeCRDAdded:
		return "CRD added"
	case ChangeCRDRemoved:
		return "CRD removed"
	}
	return kro.Change{Path: c.Path, ChangeType: kro.ChangeType(c.ChangeType), OldValue: c.Old,
		NewValue: c.New}.Description()
}
