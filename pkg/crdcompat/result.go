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
	"cmp"
	"encoding/json"
	"slices"
)

// Schema names the JSON document's format. A reader rejects any other value.
const Schema = "ack.crd-compat/v1"

// File statuses.
const (
	FileUnchanged = "unchanged"
	FileChanged   = "changed"
	FileNew       = "new"
	FileDeleted   = "deleted"
	FileSkipped   = "skipped"
)

// FileResult is what happened to one CRD file.
type FileResult struct {
	File   string `json:"file"`
	Kind   string `json:"kind,omitempty"`
	Status string `json:"status"`
	// Reason says why a file was skipped, or why a file still on disk is
	// deleted: it was a CRD at the base and is not one now.
	Reason string `json:"reason,omitempty"`
}

// Result is the `crd-compat-check --output json` document.
type Result struct {
	Schema  string `json:"schema"`
	BaseRef string `json:"base_ref"`
	// Breaking is true when any change is breaking; the command then exits
	// non-zero.
	Breaking bool         `json:"breaking"`
	Files    []FileResult `json:"files"`
	// Buckets holds every bucket in Buckets, empty ones as [].
	Buckets             map[string][]Change  `json:"buckets"`
	SensitiveCandidates []SensitiveCandidate `json:"sensitive_candidates"`
}

// NewResult returns an empty result against baseRef.
func NewResult(baseRef string) *Result {
	r := &Result{
		Schema:              Schema,
		BaseRef:             baseRef,
		Files:               []FileResult{},
		Buckets:             make(map[string][]Change, len(Buckets)),
		SensitiveCandidates: []SensitiveCandidate{},
	}
	for _, b := range Buckets {
		r.Buckets[b] = []Change{}
	}
	return r
}

// AddFile records one file's status.
func (r *Result) AddFile(f FileResult) { r.Files = append(r.Files, f) }

// AddChanges files each change under its bucket. A bucket name the result
// does not have becomes BucketOther, so it cannot pass as additive.
func (r *Result) AddChanges(cs ...Change) {
	for _, c := range cs {
		if _, ok := r.Buckets[c.Bucket]; !ok {
			c.Bucket = BucketOther
		}
		r.Buckets[c.Bucket] = append(r.Buckets[c.Bucket], c)
		r.Breaking = r.Breaking || c.Breaking
	}
}

// AddSensitive records sensitive candidates.
func (r *Result) AddSensitive(cs ...SensitiveCandidate) {
	r.SensitiveCandidates = append(r.SensitiveCandidates, cs...)
}

// sort orders everything, since kro walks maps in random order.
func (r *Result) sort() {
	slices.SortFunc(r.Files, func(a, b FileResult) int { return cmp.Compare(a.File, b.File) })
	for _, cs := range r.Buckets {
		slices.SortFunc(cs, func(a, b Change) int {
			return cmp.Or(cmp.Compare(a.File, b.File), cmp.Compare(a.Field, b.Field),
				cmp.Compare(a.Path, b.Path), cmp.Compare(a.ChangeType, b.ChangeType),
				cmp.Compare(a.Old, b.Old), cmp.Compare(a.New, b.New))
		})
	}
	slices.SortFunc(r.SensitiveCandidates, func(a, b SensitiveCandidate) int {
		return cmp.Or(cmp.Compare(a.File, b.File), cmp.Compare(a.Field, b.Field))
	})
}

// JSON sorts the result and returns it as indented JSON with a trailing newline.
func (r *Result) JSON() ([]byte, error) {
	r.sort()
	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}
