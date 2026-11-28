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

// Package crdcompat sorts the schema changes kro's compat library reports
// between two versions of a CRD into buckets that review gates can decide on,
// for `ack-generate crd-compat-check --output json`.
package crdcompat

import (
	"strings"

	kro "github.com/kubernetes-sigs/kro/pkg/graph/crd/compat"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Change buckets. A change type this package does not know goes to
// BucketOther, which a gate must not treat as additive.
//
// BucketDefaultChanged holds schema nodes whose only uncompared difference is
// default (kro declares DEFAULT_CHANGED but never emits it). BucketOther also
// holds every other UNCOMPARED_CHANGED change (see Uncompared): schema
// keywords kro does not compare, and CRD or version settings.
const (
	BucketAddedOptionalSpec   = "added_optional_spec"
	BucketAddedStatus         = "added_status"
	BucketAddedRequired       = "added_required"
	BucketTypeChanged         = "type_changed"
	BucketRemoved             = "removed"
	BucketEnumChanged         = "enum_changed"
	BucketValidationTightened = "validation_tightened" // also: list gains items; field added under an open parent
	BucketValidationRelaxed   = "validation_relaxed"
	BucketDefaultChanged      = "default_changed"
	BucketDescriptionOnly     = "description_only"
	BucketCRDAdded            = "crd_added"
	BucketCRDRemoved          = "crd_removed"
	BucketOther               = "other"
)

// Buckets lists every bucket.
var Buckets = []string{
	BucketAddedOptionalSpec, BucketAddedStatus, BucketAddedRequired, BucketTypeChanged,
	BucketRemoved, BucketEnumChanged, BucketValidationTightened, BucketValidationRelaxed,
	BucketDefaultChanged, BucketDescriptionOnly, BucketCRDAdded, BucketCRDRemoved, BucketOther,
}

// Change types for whole CRD files. kro compares only files present on both
// sides, so these come from the file listing.
const (
	ChangeCRDAdded   = "CRD_ADDED"
	ChangeCRDRemoved = "CRD_REMOVED"
)

// Change is one bucketed schema change in one CRD file.
type Change struct {
	// File is the CRD file, relative to the repository root.
	File string `json:"file"`
	// Kind is the CRD's spec.names.kind.
	Kind   string `json:"kind,omitempty"`
	Bucket string `json:"bucket"`
	// Field is the changed field as a dotted path ("spec.tags[].key"), or ""
	// for the schema root or a whole file.
	Field string `json:"field,omitempty"`
	// Path is kro's schema path, kept so a reader can locate the exact node.
	Path       string `json:"path,omitempty"`
	ChangeType string `json:"change_type"`
	Old        string `json:"old,omitempty"`
	New        string `json:"new,omitempty"`
	// Breaking is kro's verdict; a removed CRD is breaking.
	Breaking bool `json:"breaking"`
}

// Classify converts kro's report for one CRD file into bucketed changes. kro
// does not compare everything; Uncompared reports the rest. Without the old
// schema it cannot see open parents; prefer ClassifyWithBase.
func Classify(file, kind string, r *kro.Report) []Change {
	return ClassifyWithBase(file, kind, r, nil)
}

// ClassifyWithBase is Classify knowing oldRoot, the base version's
// openAPIV3Schema. A property added under a parent that had
// x-kubernetes-preserve-unknown-fields or x-kubernetes-embedded-resource is
// validation_tightened: objects already stored may hold that name with any
// value, which the new schema can now reject.
func ClassifyWithBase(file, kind string, r *kro.Report, oldRoot *v1.JSONSchemaProps) []Change {
	out := make([]Change, 0, len(r.BreakingChanges)+len(r.NonBreakingChanges))
	for _, c := range r.BreakingChanges {
		out = append(out, newChange(file, kind, c, true))
	}
	for _, c := range r.NonBreakingChanges {
		out = append(out, newChange(file, kind, c, false))
	}
	tightenUnderOpenParents(out, oldRoot)
	return out
}

// tightenUnderOpenParents rebuckets additive PROPERTY_ADDED changes whose
// parent in oldRoot kept unknown fields or was an embedded resource.
func tightenUnderOpenParents(cs []Change, oldRoot *v1.JSONSchemaProps) {
	if oldRoot == nil {
		return
	}
	for i, c := range cs {
		if c.ChangeType != string(kro.PropertyAdded) {
			continue
		}
		switch c.Bucket {
		case BucketAddedOptionalSpec, BucketAddedStatus, BucketOther:
		default:
			continue
		}
		at := strings.LastIndex(c.Path, ".properties.")
		if at < 0 {
			continue
		}
		parent := schemaAt(oldRoot, c.Path[:at])
		if parent == nil {
			continue
		}
		if (parent.XPreserveUnknownFields != nil && *parent.XPreserveUnknownFields) || parent.XEmbeddedResource {
			cs[i].Bucket = BucketValidationTightened
		}
	}
}

// CRDAdded is the change for a CRD file the base ref does not have.
func CRDAdded(file, kind string) Change {
	return Change{File: file, Kind: kind, Bucket: BucketCRDAdded, ChangeType: ChangeCRDAdded}
}

// CRDRemoved is the change for a CRD file the base ref has and the working
// tree does not. Every user with resources of that kind loses them.
func CRDRemoved(file, kind string) Change {
	return Change{File: file, Kind: kind, Bucket: BucketCRDRemoved, ChangeType: ChangeCRDRemoved, Breaking: true}
}

func newChange(file, kind string, c kro.Change, breaking bool) Change {
	field := FieldPath(c)
	return Change{
		File: file, Kind: kind,
		Bucket:     bucketFor(c.ChangeType, field, breaking),
		Field:      field,
		Path:       c.Path,
		ChangeType: string(c.ChangeType),
		Old:        c.OldValue,
		New:        c.NewValue,
		Breaking:   breaking,
	}
}

func bucketFor(t kro.ChangeType, field string, breaking bool) string {
	switch t {
	case kro.PropertyAdded:
		switch {
		case strings.HasSuffix(field, "[]"):
			// An untyped list gained an items schema: values that were
			// accepted may no longer be. kro calls it non-breaking.
			return BucketValidationTightened
		case breaking:
			// Required with no default: existing manifests stop validating.
			return BucketAddedRequired
		case under(field, "spec"):
			return BucketAddedOptionalSpec
		case under(field, "status"):
			return BucketAddedStatus
		}
		return BucketOther
	case kro.PropertyRemoved:
		return BucketRemoved
	case kro.TypeChanged:
		return BucketTypeChanged
	case kro.EnumRestricted, kro.EnumExpanded:
		return BucketEnumChanged
	case kro.RequiredAdded, kro.PatternAdded, kro.PatternChanged:
		return BucketValidationTightened
	case kro.RequiredRemoved, kro.PatternRemoved:
		return BucketValidationRelaxed
	case kro.DefaultChanged:
		return BucketDefaultChanged
	case kro.DescriptionChanged:
		return BucketDescriptionOnly
	}
	return BucketOther
}

// under reports whether field is top or lies beneath it.
func under(field, top string) bool {
	return field == top || strings.HasPrefix(field, top+".") || strings.HasPrefix(field, top+"[")
}

// FieldPath turns kro's schema path into a field path:
// ".properties.spec.properties.tags.items.properties.key" is "spec.tags[].key",
// and a map's values, ".additionalProperties", are "{}" ("spec.labels{}").
// A change to an attribute (description, type, pattern, enum) names its
// property, and a required-list change names the member that changed. Paths
// kro builds outside the schema ("version.served") are returned unchanged.
func FieldPath(c kro.Change) string {
	if c.Path == "" {
		return ""
	}
	if !strings.HasPrefix(c.Path, ".") {
		return c.Path
	}
	segs := strings.Split(strings.TrimPrefix(c.Path, "."), ".")
	var parts []string
	for i := 0; i < len(segs); i++ {
		switch segs[i] {
		case "properties":
			// The next segment is a property name, whatever it spells.
			if i+1 < len(segs) {
				i++
				parts = append(parts, segs[i])
			}
		case "items":
			if len(parts) > 0 {
				parts[len(parts)-1] += "[]"
			} else {
				parts = append(parts, "[]")
			}
		case "additionalProperties":
			// kro never builds this; Uncompared does, for a map's values.
			if len(parts) > 0 {
				parts[len(parts)-1] += "{}"
			} else {
				parts = append(parts, "{}")
			}
		case "required":
			member := c.NewValue
			if c.ChangeType == kro.RequiredRemoved {
				member = c.OldValue
			}
			parts = append(parts, member)
		}
		// Other segments (description, type, pattern, enum) are attributes of
		// the property already named.
	}
	return strings.Join(parts, ".")
}
