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

package command

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"

	"github.com/aws-controllers-k8s/code-generator/pkg/crdcompat"
)

const crdDir = "config/crd/bases"

// newRepo makes a git repository in a temp dir and changes into it.
func newRepo(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	t.Chdir(dir)
	git(t, "init", "-q", "-b", "main")
	git(t, "config", "user.email", "test@example.com")
	git(t, "config", "user.name", "test")
	git(t, "config", "commit.gpgsign", "false")
}

func git(t *testing.T, args ...string) {
	t.Helper()
	out, err := exec.Command("git", args...).CombinedOutput()
	require.NoError(t, err, "git %v: %s", args, out)
}

func commitAll(t *testing.T) {
	t.Helper()
	git(t, "add", "-A")
	git(t, "commit", "-q", "-m", "base")
}

func str() v1.JSONSchemaProps { return v1.JSONSchemaProps{Type: "string"} }

// writeCRD writes a one-version CRD for kind with the given spec properties.
func writeCRD(t *testing.T, file, kind string, spec map[string]v1.JSONSchemaProps) {
	t.Helper()
	crd := v1.CustomResourceDefinition{}
	crd.APIVersion = "apiextensions.k8s.io/v1"
	crd.Kind = "CustomResourceDefinition"
	crd.Name = kind + "s.example.services.k8s.aws"
	crd.Spec.Names.Kind = kind
	crd.Spec.Versions = []v1.CustomResourceDefinitionVersion{{
		Name: "v1alpha1", Served: true, Storage: true,
		Schema: &v1.CustomResourceValidation{OpenAPIV3Schema: &v1.JSONSchemaProps{
			Type: "object",
			Properties: map[string]v1.JSONSchemaProps{
				"spec":   {Type: "object", Properties: spec},
				"status": {Type: "object", Properties: map[string]v1.JSONSchemaProps{"arn": str()}},
			},
		}},
	}}
	data, err := yaml.Marshal(crd)
	require.NoError(t, err)
	require.NoError(t, os.MkdirAll(filepath.Dir(file), 0o755))
	require.NoError(t, os.WriteFile(file, data, 0o644))
}

func runJSON(t *testing.T, ref string) (crdcompat.Result, bool) {
	t.Helper()
	var buf bytes.Buffer
	breaking, err := runCRDCompat(&buf, ref, []string{crdDir}, "json")
	require.NoError(t, err)
	var res crdcompat.Result
	require.NoError(t, json.Unmarshal(buf.Bytes(), &res), buf.String())
	return res, breaking
}

func TestCRDCompatJSONBucketsAndSensitive(t *testing.T) {
	newRepo(t)
	file := crdDir + "/widgets.yaml"
	writeCRD(t, file, "Widget", map[string]v1.JSONSchemaProps{"name": str()})
	commitAll(t)
	writeCRD(t, file, "Widget", map[string]v1.JSONSchemaProps{
		"name":     str(),
		"apiToken": str(),
		"config":   {Type: "object", Properties: map[string]v1.JSONSchemaProps{"password": str()}},
	})

	res, breaking := runJSON(t, "main")
	assert.False(t, breaking)
	assert.False(t, res.Breaking)
	assert.Equal(t, crdcompat.Schema, res.Schema)
	assert.Equal(t, []crdcompat.FileResult{{File: file, Kind: "Widget", Status: crdcompat.FileChanged}}, res.Files)
	added := res.Buckets[crdcompat.BucketAddedOptionalSpec]
	require.Len(t, added, 2)
	assert.Equal(t, "spec.apiToken", added[0].Field)
	assert.Equal(t, "spec.config", added[1].Field)
	assert.Equal(t, []crdcompat.SensitiveCandidate{
		{File: file, Kind: "Widget", Field: "spec.apiToken"},
		{File: file, Kind: "Widget", Field: "spec.config.password"},
	}, res.SensitiveCandidates)
}

func TestCRDCompatNewAndDeletedCRDs(t *testing.T) {
	newRepo(t)
	writeCRD(t, crdDir+"/widgets.yaml", "Widget", map[string]v1.JSONSchemaProps{"name": str()})
	writeCRD(t, crdDir+"/gadgets.yaml", "Gadget", map[string]v1.JSONSchemaProps{"name": str()})
	commitAll(t)
	require.NoError(t, os.Remove(crdDir+"/gadgets.yaml"))
	writeCRD(t, crdDir+"/gizmos.yaml", "Gizmo", map[string]v1.JSONSchemaProps{"clientSecret": str()})
	require.NoError(t, os.WriteFile(crdDir+"/kustomization.yaml", []byte("resources: []\n"), 0o644))

	res, breaking := runJSON(t, "main")
	assert.True(t, breaking, "a deleted CRD is breaking")
	assert.Equal(t, []crdcompat.FileResult{
		{File: crdDir + "/gadgets.yaml", Kind: "Gadget", Status: crdcompat.FileDeleted},
		{File: crdDir + "/gizmos.yaml", Kind: "Gizmo", Status: crdcompat.FileNew},
		{File: crdDir + "/kustomization.yaml", Status: crdcompat.FileSkipped, Reason: "not a CRD (kind=)"},
		{File: crdDir + "/widgets.yaml", Kind: "Widget", Status: crdcompat.FileUnchanged},
	}, res.Files)
	assert.Equal(t, []crdcompat.Change{crdcompat.CRDRemoved(crdDir+"/gadgets.yaml", "Gadget")},
		res.Buckets[crdcompat.BucketCRDRemoved])
	assert.Equal(t, []crdcompat.Change{crdcompat.CRDAdded(crdDir+"/gizmos.yaml", "Gizmo")},
		res.Buckets[crdcompat.BucketCRDAdded])
	assert.Equal(t, []crdcompat.SensitiveCandidate{
		{File: crdDir + "/gizmos.yaml", Kind: "Gizmo", Field: "spec.clientSecret"},
	}, res.SensitiveCandidates)

	var text bytes.Buffer
	breaking, err := runCRDCompat(&text, "main", []string{crdDir}, "text")
	require.NoError(t, err)
	assert.True(t, breaking)
	assert.Contains(t, text.String(), "FAIL "+crdDir+"/gadgets.yaml\n  BREAKING: CRD removed (in main, not on disk)\n")
	assert.Contains(t, text.String(), "NEW  "+crdDir+"/gizmos.yaml (not in main)\n")
}

func TestCRDCompatTextReportsBreakingTypeChange(t *testing.T) {
	newRepo(t)
	file := crdDir + "/widgets.yaml"
	writeCRD(t, file, "Widget", map[string]v1.JSONSchemaProps{"name": str()})
	commitAll(t)
	writeCRD(t, file, "Widget", map[string]v1.JSONSchemaProps{"name": {Type: "integer"}})

	var text bytes.Buffer
	breaking, err := runCRDCompat(&text, "main", []string{crdDir}, "text")
	require.NoError(t, err)
	assert.True(t, breaking)
	assert.Contains(t, text.String(), "FAIL "+file+"\n")
	assert.Contains(t, text.String(), "  BREAKING: [.properties.spec.properties.name.type] Type changed from string to integer\n")
	assert.Contains(t, text.String(), "\nChecked 1 CRD(s) across 1 file(s).\n")

	res, _ := runJSON(t, "main")
	require.Len(t, res.Buckets[crdcompat.BucketTypeChanged], 1)
	assert.Equal(t, "spec.name", res.Buckets[crdcompat.BucketTypeChanged][0].Field)
}

func TestCRDCompatRejectsBadRefAndOutput(t *testing.T) {
	newRepo(t)
	writeCRD(t, crdDir+"/widgets.yaml", "Widget", map[string]v1.JSONSchemaProps{"name": str()})
	commitAll(t)

	_, err := runCRDCompat(&bytes.Buffer{}, "no-such-ref", []string{crdDir}, "text")
	assert.ErrorContains(t, err, `base ref "no-such-ref" is not a commit`)
	_, err = runCRDCompat(&bytes.Buffer{}, "-x", []string{crdDir}, "text")
	assert.ErrorContains(t, err, "is not a commit")
	_, err = runCRDCompat(&bytes.Buffer{}, "main", []string{crdDir}, "yaml")
	assert.ErrorContains(t, err, `--output must be "text" or "json"`)
}

func TestCRDCompatNoFilesIsAnError(t *testing.T) {
	newRepo(t)
	require.NoError(t, os.WriteFile("README", []byte("x"), 0o644))
	commitAll(t)
	_, err := runCRDCompat(&bytes.Buffer{}, "main", []string{crdDir}, "json")
	assert.ErrorContains(t, err, "no CRD files found")
}

// A head file that no longer parses as a CRD, but was one at the base, is a
// removed CRD: the cluster loses the kind.
func TestCRDCompatHeadFileNoLongerACRDIsRemoved(t *testing.T) {
	cases := []struct {
		name, content, reason string
	}{
		{"emptied", "", "not a CRD (kind=)"},
		{"kustomization", "resources: []\n", "not a CRD (kind=)"},
		{"unparsable", "a: [\n", "not a valid CRD"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			newRepo(t)
			file := crdDir + "/widgets.yaml"
			writeCRD(t, file, "Widget", map[string]v1.JSONSchemaProps{"name": str()})
			commitAll(t)
			require.NoError(t, os.WriteFile(file, []byte(tc.content), 0o644))

			res, breaking := runJSON(t, "main")
			assert.True(t, breaking)
			require.Len(t, res.Files, 1)
			assert.Equal(t, file, res.Files[0].File)
			assert.Equal(t, "Widget", res.Files[0].Kind)
			assert.Equal(t, crdcompat.FileDeleted, res.Files[0].Status)
			assert.Contains(t, res.Files[0].Reason, tc.reason)
			assert.Equal(t, []crdcompat.Change{crdcompat.CRDRemoved(file, "Widget")},
				res.Buckets[crdcompat.BucketCRDRemoved])

			var text bytes.Buffer
			breaking, err := runCRDCompat(&text, "main", []string{crdDir}, "text")
			require.NoError(t, err)
			assert.True(t, breaking)
			assert.Contains(t, text.String(), "FAIL "+file+"\n  BREAKING: CRD removed ("+tc.reason)
		})
	}
}
