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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptyResultHasEveryBucket(t *testing.T) {
	out, err := NewResult("main").JSON()
	require.NoError(t, err)
	assert.Equal(t, byte('\n'), out[len(out)-1])

	var doc map[string]any
	require.NoError(t, json.Unmarshal(out, &doc))
	assert.Equal(t, Schema, doc["schema"])
	assert.Equal(t, "main", doc["base_ref"])
	assert.Equal(t, false, doc["breaking"])
	assert.Equal(t, []any{}, doc["files"])
	assert.Equal(t, []any{}, doc["sensitive_candidates"])
	buckets := doc["buckets"].(map[string]any)
	assert.Len(t, buckets, len(Buckets))
	for _, b := range Buckets {
		assert.Equal(t, []any{}, buckets[b], b)
	}
}

func TestResultBreakingAndOrder(t *testing.T) {
	r := NewResult("abc123")
	r.AddFile(FileResult{File: "b.yaml", Kind: "B", Status: FileChanged})
	r.AddFile(FileResult{File: "a.yaml", Kind: "A", Status: FileUnchanged})
	r.AddChanges(
		Change{File: "b.yaml", Bucket: BucketAddedOptionalSpec, Field: "spec.z", ChangeType: "PROPERTY_ADDED"},
		Change{File: "a.yaml", Bucket: BucketAddedOptionalSpec, Field: "spec.y", ChangeType: "PROPERTY_ADDED"},
		Change{File: "b.yaml", Bucket: BucketAddedOptionalSpec, Field: "spec.a", ChangeType: "PROPERTY_ADDED"},
	)
	r.AddSensitive(SensitiveCandidate{File: "b.yaml", Field: "spec.token"},
		SensitiveCandidate{File: "a.yaml", Field: "spec.secret"})
	assert.False(t, r.Breaking)

	r.AddChanges(Change{File: "c.yaml", Bucket: "not-a-bucket", ChangeType: "X", Breaking: true})
	assert.True(t, r.Breaking)

	_, err := r.JSON()
	require.NoError(t, err)
	assert.Equal(t, []string{"a.yaml", "b.yaml"}, []string{r.Files[0].File, r.Files[1].File})
	got := r.Buckets[BucketAddedOptionalSpec]
	assert.Equal(t, []string{"spec.y", "spec.a", "spec.z"}, []string{got[0].Field, got[1].Field, got[2].Field})
	assert.Equal(t, "a.yaml", r.SensitiveCandidates[0].File)
	require.Len(t, r.Buckets[BucketOther], 1, "an unknown bucket fails closed into other")
	assert.Equal(t, BucketOther, r.Buckets[BucketOther][0].Bucket)
}
