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
	"regexp"
	"sort"
	"strings"

	kro "github.com/kubernetes-sigs/kro/pkg/graph/crd/compat"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// sensitiveName matches field names that may hold a secret.
var sensitiveName = regexp.MustCompile(`(?i)password|secret|token|privatekey|credential`)

// SensitiveCandidate is an added field whose name suggests a secret. It needs
// an explicit is_secret decision. The name alone decides nothing: a CRD does
// not carry Smithy's @sensitive trait, so the match is a tripwire only.
type SensitiveCandidate struct {
	File  string `json:"file"`
	Kind  string `json:"kind,omitempty"`
	Field string `json:"field"`
}

// SensitiveFields returns each added field, and every field beneath it
// (list elements and map values included), whose name matches. kro reports an
// added object once, at the object, so its members are found by walking root,
// the new version's openAPIV3Schema. changes may include Uncompared's, which
// report fields added inside map values.
func SensitiveFields(file, kind string, root *v1.JSONSchemaProps, changes []Change) []SensitiveCandidate {
	var out []SensitiveCandidate
	emit := func(field string) { out = append(out, SensitiveCandidate{File: file, Kind: kind, Field: field}) }
	for _, c := range changes {
		if c.ChangeType != string(kro.PropertyAdded) {
			continue
		}
		if s := schemaAt(root, c.Path); s != nil {
			walk(c.Field, s, emit)
		}
	}
	return out
}

// NewCRDSensitiveFields is SensitiveFields for a CRD the base does not have,
// where every Spec and Status field is new.
func NewCRDSensitiveFields(file, kind string, root *v1.JSONSchemaProps) []SensitiveCandidate {
	if root == nil {
		return nil
	}
	var out []SensitiveCandidate
	emit := func(field string) { out = append(out, SensitiveCandidate{File: file, Kind: kind, Field: field}) }
	for _, top := range []string{"spec", "status"} {
		if s, ok := root.Properties[top]; ok {
			walk(top, &s, emit)
		}
	}
	return out
}

// walk emits field when its name matches, then visits its members. A
// SecretKeyReference is already a secret and is neither emitted nor entered.
func walk(field string, s *v1.JSONSchemaProps, emit func(string)) {
	if isSecretKeyReference(s) {
		return
	}
	if sensitiveName.MatchString(lastName(field)) {
		emit(field)
	}
	walkMembers(field, s, emit)
}

// walkMembers visits s's properties, list elements ("[]") and map values
// ("{}"). An element or value has no name of its own, so it is not
// re-matched against the list's or map's name, which walk already checked.
func walkMembers(field string, s *v1.JSONSchemaProps, emit func(string)) {
	if isSecretKeyReference(s) {
		return
	}
	names := make([]string, 0, len(s.Properties))
	for name := range s.Properties {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		p := s.Properties[name]
		child := name
		if field != "" {
			child = field + "." + name
		}
		walk(child, &p, emit)
	}
	if s.Items != nil && s.Items.Schema != nil {
		walkMembers(field+"[]", s.Items.Schema, emit)
	}
	if s.AdditionalProperties != nil && s.AdditionalProperties.Schema != nil {
		walkMembers(field+"{}", s.AdditionalProperties.Schema, emit)
	}
}

// isSecretKeyReference reports whether s has exactly the members of ACK's
// SecretKeyReference, which code-generator emits for an is_secret field:
// key, name and namespace, all strings.
func isSecretKeyReference(s *v1.JSONSchemaProps) bool {
	if len(s.Properties) != 3 {
		return false
	}
	for _, name := range []string{"key", "name", "namespace"} {
		if p, ok := s.Properties[name]; !ok || p.Type != "string" {
			return false
		}
	}
	return true
}

// lastName is the last property name in a field path, without list or map
// markers.
func lastName(field string) string {
	field = strings.TrimRight(field, "[]{}")
	if i := strings.LastIndex(field, "."); i >= 0 {
		return field[i+1:]
	}
	return field
}

// schemaAt resolves kro's schema path in root ("" is root itself), or
// returns nil when the path does not lead to a schema node.
func schemaAt(root *v1.JSONSchemaProps, path string) *v1.JSONSchemaProps {
	if root == nil || path == "" {
		return root
	}
	if !strings.HasPrefix(path, ".") {
		return nil
	}
	s := root
	segs := strings.Split(strings.TrimPrefix(path, "."), ".")
	for i := 0; i < len(segs); i++ {
		switch segs[i] {
		case "properties":
			if i+1 >= len(segs) {
				return nil
			}
			i++
			p, ok := s.Properties[segs[i]]
			if !ok {
				return nil
			}
			s = &p
		case "items":
			if s.Items == nil || s.Items.Schema == nil {
				return nil
			}
			s = s.Items.Schema
		case "additionalProperties":
			if s.AdditionalProperties == nil || s.AdditionalProperties.Schema == nil {
				return nil
			}
			s = s.AdditionalProperties.Schema
		default:
			return nil
		}
	}
	return s
}
