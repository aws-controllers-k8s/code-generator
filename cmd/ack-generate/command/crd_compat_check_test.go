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

// newRepo makes a git repository in a temp dir, changes into it, and
// returns its path as created (on macOS, unresolved: /var, not /private/var).
func newRepo(t *testing.T) string {
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
	return dir
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
	writeCRDWith(t, file, kind, spec, nil)
}

// writeCRDWith is writeCRD with mutate applied to the CRD before writing.
func writeCRDWith(t *testing.T, file, kind string, spec map[string]v1.JSONSchemaProps,
	mutate func(*v1.CustomResourceDefinition)) {
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
	if mutate != nil {
		mutate(&crd)
	}
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

// Changes kro does not compare still make a file changed, and land in a
// bucket a gate cannot pass as additive.
func TestCRDCompatReportsUncomparedChanges(t *testing.T) {
	newRepo(t)
	file := crdDir + "/widgets.yaml"
	max := func(n int64) v1.JSONSchemaProps { s := str(); s.MaxLength = &n; return s }
	writeCRD(t, file, "Widget", map[string]v1.JSONSchemaProps{"name": max(5)})
	commitAll(t)
	writeCRD(t, file, "Widget", map[string]v1.JSONSchemaProps{"name": max(10)})

	res, breaking := runJSON(t, "main")
	assert.False(t, breaking)
	assert.Equal(t, []crdcompat.FileResult{{File: file, Kind: "Widget", Status: crdcompat.FileChanged}}, res.Files)
	require.Len(t, res.Buckets[crdcompat.BucketOther], 1)
	c := res.Buckets[crdcompat.BucketOther][0]
	assert.Equal(t, "spec.name", c.Field)
	assert.Equal(t, crdcompat.ChangeUncomparedChanged, c.ChangeType)
	assert.Equal(t, "maxLength", c.New)

	var text bytes.Buffer
	breaking, err := runCRDCompat(&text, "main", []string{crdDir}, "text")
	require.NoError(t, err)
	assert.False(t, breaking)
	assert.Contains(t, text.String(), "OK   "+file+" (non-breaking changes only)\n")
}

func TestCRDCompatUncomparedBreakingScope(t *testing.T) {
	newRepo(t)
	file := crdDir + "/widgets.yaml"
	spec := map[string]v1.JSONSchemaProps{"name": str()}
	writeCRDWith(t, file, "Widget", spec, func(c *v1.CustomResourceDefinition) { c.Spec.Scope = v1.NamespaceScoped })
	commitAll(t)
	writeCRDWith(t, file, "Widget", spec, func(c *v1.CustomResourceDefinition) {
		c.Spec.Scope = v1.ClusterScoped
		c.Spec.Versions[0].AdditionalPrinterColumns = []v1.CustomResourceColumnDefinition{
			{Name: "ARN", Type: "string", JSONPath: ".status.arn"}}
	})

	res, breaking := runJSON(t, "main")
	assert.True(t, breaking)
	assert.True(t, res.Breaking)
	assert.Equal(t, crdcompat.FileChanged, res.Files[0].Status)
	require.Len(t, res.Buckets[crdcompat.BucketOther], 2)

	var text bytes.Buffer
	breaking, err := runCRDCompat(&text, "main", []string{crdDir}, "text")
	require.NoError(t, err)
	assert.True(t, breaking)
	assert.Contains(t, text.String(), "FAIL "+file+"\n")
	assert.Contains(t, text.String(), "  BREAKING: [spec.scope] Uncompared change: scope\n")
	assert.Contains(t, text.String(),
		"  ok:       [version.additionalPrinterColumns] Uncompared change: additionalPrinterColumns\n")
}

func runJSONDirs(t *testing.T, dirs ...string) crdcompat.Result {
	t.Helper()
	var buf bytes.Buffer
	_, err := runCRDCompat(&buf, "main", dirs, "json")
	require.NoError(t, err)
	var res crdcompat.Result
	require.NoError(t, json.Unmarshal(buf.Bytes(), &res), buf.String())
	return res
}

func TestCRDCompatAbsoluteCRDPath(t *testing.T) {
	dir := newRepo(t)
	file := crdDir + "/widgets.yaml"
	writeCRD(t, file, "Widget", map[string]v1.JSONSchemaProps{"name": str()})
	commitAll(t)

	res := runJSONDirs(t, filepath.Join(dir, crdDir))
	assert.Equal(t, []crdcompat.FileResult{{File: file, Kind: "Widget", Status: crdcompat.FileUnchanged}}, res.Files)
	assert.Empty(t, res.Buckets[crdcompat.BucketCRDAdded])
}

func TestCRDCompatRejectsPathsOutsideRepo(t *testing.T) {
	newRepo(t)
	writeCRD(t, crdDir+"/widgets.yaml", "Widget", map[string]v1.JSONSchemaProps{"name": str()})
	commitAll(t)
	outside := t.TempDir()
	for _, dir := range []string{outside, "../elsewhere", crdDir + "/../../../.."} {
		_, err := runCRDCompat(&bytes.Buffer{}, "main", []string{dir}, "json")
		assert.ErrorContains(t, err, "outside the repository", dir)
	}
}

func TestCRDCompatDedupesDirs(t *testing.T) {
	newRepo(t)
	writeCRD(t, crdDir+"/widgets.yaml", "Widget", map[string]v1.JSONSchemaProps{"name": str()})
	writeCRD(t, crdDir+"/gadgets.yaml", "Gadget", map[string]v1.JSONSchemaProps{"name": str()})
	commitAll(t)
	require.NoError(t, os.Remove(crdDir+"/gadgets.yaml"))
	writeCRD(t, crdDir+"/gizmos.yaml", "Gizmo", map[string]v1.JSONSchemaProps{"name": str()})

	res := runJSONDirs(t, crdDir, crdDir+"/", "./"+crdDir)
	assert.Len(t, res.Files, 3)
	assert.Len(t, res.Buckets[crdcompat.BucketCRDAdded], 1)
	assert.Len(t, res.Buckets[crdcompat.BucketCRDRemoved], 1)

	var text bytes.Buffer
	_, err := runCRDCompat(&text, "main", []string{crdDir, crdDir + "/"}, "text")
	require.NoError(t, err)
	assert.Contains(t, text.String(), "\nChecked 1 CRD(s) across 2 file(s).\n")
}

func TestCRDCompatMustRunFromRepoRoot(t *testing.T) {
	newRepo(t)
	writeCRD(t, crdDir+"/widgets.yaml", "Widget", map[string]v1.JSONSchemaProps{"name": str()})
	commitAll(t)
	t.Chdir(crdDir)
	_, err := runCRDCompat(&bytes.Buffer{}, "main", []string{"."}, "json")
	assert.ErrorContains(t, err, "run crd-compat-check from the repository root")
}

// git quotes non-ASCII names in ls-tree output unless asked for raw names.
func TestCRDCompatNonASCIIFileName(t *testing.T) {
	newRepo(t)
	file := crdDir + "/wïdgets.yaml"
	writeCRD(t, file, "Widget", map[string]v1.JSONSchemaProps{"name": str()})
	commitAll(t)

	res := runJSONDirs(t, crdDir)
	assert.Equal(t, []crdcompat.FileResult{{File: file, Kind: "Widget", Status: crdcompat.FileUnchanged}}, res.Files)

	require.NoError(t, os.Remove(file))
	writeCRD(t, crdDir+"/other.yaml", "Other", map[string]v1.JSONSchemaProps{"name": str()})
	res = runJSONDirs(t, crdDir)
	assert.Equal(t, []crdcompat.Change{crdcompat.CRDRemoved(file, "Widget")}, res.Buckets[crdcompat.BucketCRDRemoved])
}

// Text mode keeps upstream's order: per --crd-paths dir, *.yaml then *.yml.
func TestCRDCompatTextOrder(t *testing.T) {
	newRepo(t)
	spec := map[string]v1.JSONSchemaProps{"name": str()}
	writeCRD(t, "config/crd/bases/b.yaml", "B", spec)
	writeCRD(t, "helm/crds/a.yml", "A", spec)
	writeCRD(t, "helm/crds/z.yaml", "Z", spec)
	commitAll(t)

	var text bytes.Buffer
	_, err := runCRDCompat(&text, "main", []string{"helm/crds", "config/crd/bases"}, "text")
	require.NoError(t, err)
	assert.Equal(t, "OK   helm/crds/z.yaml (no changes)\n"+
		"OK   helm/crds/a.yml (no changes)\n"+
		"OK   config/crd/bases/b.yaml (no changes)\n"+
		"\nChecked 3 CRD(s) across 3 file(s).\n", text.String())
}

func TestCRDCompatHelpStatesExitAndJSONContract(t *testing.T) {
	assert.Contains(t, crdCompatCheckCmd.Long, "Exits 1 when breaking changes are detected or on an error.")
	assert.Contains(t, crdCompatCheckCmd.Long, "printed if and only if the analysis completed")
	assert.Contains(t, crdCompatCheckCmd.Long, "CRD metadata (annotations, labels) is not compared")
}

// yaml.Unmarshal reads only a file's first document, so anything after it
// would go unchecked. One CRD per file is required.
func TestCRDCompatMultipleDocuments(t *testing.T) {
	spec := map[string]v1.JSONSchemaProps{"name": str()}
	appendDoc := func(t *testing.T, file string) {
		t.Helper()
		data, err := os.ReadFile(file)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(file, append(data, []byte("---\nkind: ConfigMap\n")...), 0o644))
	}
	t.Run("head", func(t *testing.T) {
		newRepo(t)
		file := crdDir + "/widgets.yaml"
		writeCRD(t, file, "Widget", spec)
		commitAll(t)
		appendDoc(t, file)
		_, err := runCRDCompat(&bytes.Buffer{}, "main", []string{crdDir}, "json")
		assert.ErrorContains(t, err, "multiple YAML documents in "+file+"; one CRD per file is required")
	})
	t.Run("base", func(t *testing.T) {
		newRepo(t)
		file := crdDir + "/widgets.yaml"
		writeCRD(t, file, "Widget", spec)
		appendDoc(t, file)
		commitAll(t)
		writeCRD(t, file, "Widget", spec)
		_, err := runCRDCompat(&bytes.Buffer{}, "main", []string{crdDir}, "json")
		assert.ErrorContains(t, err, "multiple YAML documents in "+file+"; one CRD per file is required")
	})
	t.Run("base, deleted on disk", func(t *testing.T) {
		newRepo(t)
		file := crdDir + "/widgets.yaml"
		writeCRD(t, file, "Widget", spec)
		appendDoc(t, file)
		writeCRD(t, crdDir+"/other.yaml", "Other", spec)
		commitAll(t)
		require.NoError(t, os.Remove(file))
		_, err := runCRDCompat(&bytes.Buffer{}, "main", []string{crdDir}, "json")
		assert.ErrorContains(t, err, "multiple YAML documents in "+file)
	})
	t.Run("empty documents are fine", func(t *testing.T) {
		newRepo(t)
		file := crdDir + "/widgets.yaml"
		writeCRD(t, file, "Widget", spec)
		data, err := os.ReadFile(file)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(file, append([]byte("---\n# comment\n---\n"), append(data, "---\n"...)...), 0o644))
		commitAll(t)
		res := runJSONDirs(t, crdDir)
		assert.Equal(t, []crdcompat.FileResult{{File: file, Kind: "Widget", Status: crdcompat.FileUnchanged}}, res.Files)
	})
}

func TestCRDCompatAddedUnderPreserveUnknownFields(t *testing.T) {
	newRepo(t)
	file := crdDir + "/widgets.yaml"
	yes := true
	open := func(c *v1.CustomResourceDefinition) {
		s := c.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"]
		s.XPreserveUnknownFields = &yes
		c.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties["spec"] = s
	}
	writeCRDWith(t, file, "Widget", nil, open)
	commitAll(t)
	writeCRDWith(t, file, "Widget", map[string]v1.JSONSchemaProps{"x": str()}, open)

	res, _ := runJSON(t, "main")
	require.Len(t, res.Buckets[crdcompat.BucketValidationTightened], 1)
	assert.Equal(t, "spec.x", res.Buckets[crdcompat.BucketValidationTightened][0].Field)
	assert.Empty(t, res.Buckets[crdcompat.BucketAddedOptionalSpec])
}
