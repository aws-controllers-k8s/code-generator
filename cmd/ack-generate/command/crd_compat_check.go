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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"

	"github.com/aws-controllers-k8s/code-generator/pkg/crdcompat"
	crdcompatkro "github.com/kubernetes-sigs/kro/pkg/graph/crd/compat"
)

var (
	optBaseRef  string
	optCRDPaths []string
	optOutput   string
)

var crdCompatCheckCmd = &cobra.Command{
	Use:   "crd-compat-check",
	Short: "Check CRDs for breaking changes against a git base ref",
	Long: `Compares CRD files on disk against the same files at a given git ref
and reports any breaking schema changes. A CRD file that exists at the
ref but not on disk, or no longer parses as a CRD, is a breaking change.
Exits 1 when breaking changes are detected or on an error.

By default it checks config/crd/bases and helm/crds against the
main branch. It must run from the repository root; --crd-paths must be
inside the repository (absolute paths are made root-relative).

Besides what kro's compat library compares (description, type, pattern,
properties, required, enum, items, version name/served/status), any other
difference in a CRD present on both sides is reported as
UNCOMPARED_CHANGED: other schema keywords (maxLength, default, CEL rules,
an enum added where there was none, map value schemas), and CRD or version
settings (group, scope, names, conversion, printer columns, a version
starting to be served, ...). A change to group, scope, names.kind or
names.plural is breaking.

CRD metadata (annotations, labels) is not compared: it changes with every
regeneration (the controller-gen version annotation) and does not affect
what the API server accepts.

--output json prints one ack.crd-compat/v1 document to stdout: every
change sorted into buckets (added_optional_spec, added_status,
added_required, type_changed, removed, enum_changed,
validation_tightened, validation_relaxed, default_changed,
description_only, crd_added, crd_removed, other), and the added fields
whose names suggest a secret (sensitive_candidates). The document is
printed if and only if the analysis completed: on exit 1 with a document
the changes are breaking; on exit 1 without one there was an error
(reported on stderr) and nothing was decided.`,
	RunE: checkCRDCompat,
}

func init() {
	crdCompatCheckCmd.Flags().StringVar(
		&optBaseRef, "base-ref", "main",
		"Git ref to compare against (branch, tag, or SHA)",
	)
	crdCompatCheckCmd.Flags().StringSliceVar(
		&optCRDPaths, "crd-paths",
		[]string{"config/crd/bases", "helm/crds"},
		"Paths to directories containing CRD YAML files",
	)
	crdCompatCheckCmd.Flags().StringVar(
		&optOutput, "output", "text",
		`Output format: "text" or "json"`,
	)
	rootCmd.AddCommand(crdCompatCheckCmd)
}

func checkCRDCompat(cmd *cobra.Command, args []string) error {
	breaking, err := runCRDCompat(cmd.OutOrStdout(), optBaseRef, optCRDPaths, optOutput)
	if err != nil {
		return err
	}
	if breaking {
		return fmt.Errorf("breaking CRD changes detected")
	}
	return nil
}

// runCRDCompat compares the CRD files in dirs against baseRef, writes the
// report to w, and returns whether any change is breaking.
func runCRDCompat(w io.Writer, baseRef string, dirs []string, output string) (bool, error) {
	if output != "text" && output != "json" {
		return false, fmt.Errorf(`--output must be "text" or "json", got %q`, output)
	}
	text := output == "text"
	dirs, err := repoDirs(dirs)
	if err != nil {
		return false, err
	}
	if err := verifyCommit(baseRef); err != nil {
		return false, err
	}
	head, err := headCRDFiles(dirs)
	if err != nil {
		return false, err
	}
	if len(head) == 0 {
		if text {
			fmt.Fprintln(w, "No CRD files found, nothing to check.")
		}
		return false, fmt.Errorf("no CRD files found in %v", dirs)
	}
	base, err := baseCRDFiles(baseRef, dirs)
	if err != nil {
		return false, err
	}

	res := crdcompat.NewResult(baseRef)
	checked := 0
	onDisk := map[string]bool{}
	for _, file := range head {
		onDisk[file] = true
		newData, err := os.ReadFile(file)
		if err != nil {
			return false, fmt.Errorf("reading %s: %w", file, err)
		}
		// skip records a file that is not a CRD on disk. If it was a CRD at
		// the base, the kind is gone: that is a removed CRD, not a skip.
		skip := func(reason string) error {
			if base[file] {
				oldCRD, ok, err := baseCRD(baseRef, file)
				if err != nil {
					return err
				}
				if ok {
					kind := oldCRD.Spec.Names.Kind
					res.AddFile(crdcompat.FileResult{File: file, Kind: kind, Status: crdcompat.FileDeleted, Reason: reason})
					res.AddChanges(crdcompat.CRDRemoved(file, kind))
					if text {
						fmt.Fprintf(w, "FAIL %s\n  BREAKING: CRD removed (%s; a CRD in %s)\n", file, reason, baseRef)
					}
					return nil
				}
			}
			res.AddFile(crdcompat.FileResult{File: file, Status: crdcompat.FileSkipped, Reason: reason})
			if text {
				fmt.Fprintf(w, "SKIP %s: %s\n", file, reason)
			}
			return nil
		}
		newData, err = onlyDocument(file, newData)
		if err != nil {
			return false, err
		}
		var newCRD v1.CustomResourceDefinition
		if err := yaml.Unmarshal(newData, &newCRD); err != nil {
			if err := skip(fmt.Sprintf("not a valid CRD (%v)", err)); err != nil {
				return false, err
			}
			continue
		}
		if newCRD.Kind != "CustomResourceDefinition" {
			if err := skip(fmt.Sprintf("not a CRD (kind=%s)", newCRD.Kind)); err != nil {
				return false, err
			}
			continue
		}
		kind := newCRD.Spec.Names.Kind

		if !base[file] {
			res.AddFile(crdcompat.FileResult{File: file, Kind: kind, Status: crdcompat.FileNew})
			res.AddChanges(crdcompat.CRDAdded(file, kind))
			res.AddSensitive(crdcompat.NewCRDSensitiveFields(file, kind, firstSchema(newCRD))...)
			if text {
				fmt.Fprintf(w, "NEW  %s (not in %s)\n", file, baseRef)
			}
			continue
		}

		oldData, err := gitShow(baseRef, file)
		if err != nil {
			return false, err
		}
		if oldData, err = onlyDocument(file, oldData); err != nil {
			return false, fmt.Errorf("%w at %s", err, baseRef)
		}
		var oldCRD v1.CustomResourceDefinition
		if err := yaml.Unmarshal(oldData, &oldCRD); err != nil {
			return false, fmt.Errorf("parsing old %s at %s: %w", file, baseRef, err)
		}
		report, err := crdcompatkro.CompareVersions(oldCRD.Spec.Versions, newCRD.Spec.Versions)
		if err != nil {
			return false, fmt.Errorf("comparing %s: %w", file, err)
		}

		checked++
		// kro compares only some of a CRD; residual is everything else.
		residual := crdcompat.Uncompared(file, kind, &oldCRD, &newCRD)
		changes := append(crdcompat.ClassifyWithBase(file, kind, report, firstSchema(oldCRD)), residual...)
		res.AddChanges(changes...)
		res.AddSensitive(crdcompat.SensitiveFields(file, kind, firstSchema(newCRD), changes)...)
		status := crdcompat.FileUnchanged
		if len(changes) > 0 {
			status = crdcompat.FileChanged
		}
		res.AddFile(crdcompat.FileResult{File: file, Kind: kind, Status: status})
		if !text {
			continue
		}
		residualBreaking := slices.ContainsFunc(residual, func(c crdcompat.Change) bool { return c.Breaking })
		switch {
		case report.HasBreakingChanges() || residualBreaking:
			fmt.Fprintf(w, "FAIL %s\n", file)
			for _, c := range report.BreakingChanges {
				fmt.Fprintf(w, "  BREAKING: [%s] %s\n", c.Path, c.Description())
			}
			for _, c := range residual {
				if c.Breaking {
					fmt.Fprintf(w, "  BREAKING: [%s] %s\n", c.Path, c.Description())
				}
			}
			for _, c := range report.NonBreakingChanges {
				fmt.Fprintf(w, "  ok:       [%s] %s\n", c.Path, c.Description())
			}
			for _, c := range residual {
				if !c.Breaking {
					fmt.Fprintf(w, "  ok:       [%s] %s\n", c.Path, c.Description())
				}
			}
		case len(changes) > 0:
			fmt.Fprintf(w, "OK   %s (non-breaking changes only)\n", file)
		default:
			fmt.Fprintf(w, "OK   %s (no changes)\n", file)
		}
	}

	for _, file := range sortedKeys(base) {
		if onDisk[file] {
			continue
		}
		oldCRD, ok, err := baseCRD(baseRef, file)
		if err != nil {
			return false, err
		}
		if !ok {
			// It was not a CRD at the base either, so nothing was lost.
			continue
		}
		kind := oldCRD.Spec.Names.Kind
		res.AddFile(crdcompat.FileResult{File: file, Kind: kind, Status: crdcompat.FileDeleted})
		res.AddChanges(crdcompat.CRDRemoved(file, kind))
		if text {
			fmt.Fprintf(w, "FAIL %s\n  BREAKING: CRD removed (in %s, not on disk)\n", file, baseRef)
		}
	}

	if text {
		fmt.Fprintf(w, "\nChecked %d CRD(s) across %d file(s).\n", checked, len(head))
		return res.Breaking, nil
	}
	out, err := res.JSON()
	if err != nil {
		return false, err
	}
	if _, err := w.Write(out); err != nil {
		return false, err
	}
	return res.Breaking, nil
}

// repoDirs checks that the command runs from the repository root, where
// paths on disk and paths at the base ref agree, and returns dirs relative
// to the root, cleaned and without duplicates. An absolute dir inside the
// repository is made relative; a dir outside it is an error.
func repoDirs(dirs []string) ([]string, error) {
	prefix, err := gitOutput("rev-parse", "--show-prefix")
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(string(prefix)) != "" {
		return nil, fmt.Errorf("run crd-compat-check from the repository root")
	}
	top, err := gitOutput("rev-parse", "--show-toplevel")
	if err != nil {
		return nil, err
	}
	root := resolvePath(strings.TrimSpace(string(top)))
	var out []string
	seen := map[string]bool{}
	for _, dir := range dirs {
		rel := filepath.Clean(dir)
		if filepath.IsAbs(rel) {
			// Compare with symlinks resolved: on macOS git reports
			// /private/var for a directory under /var.
			if rel, err = filepath.Rel(root, resolvePath(rel)); err != nil {
				return nil, fmt.Errorf("--crd-paths %q is outside the repository", dir)
			}
		}
		rel = filepath.ToSlash(rel)
		if rel == ".." || strings.HasPrefix(rel, "../") {
			return nil, fmt.Errorf("--crd-paths %q is outside the repository", dir)
		}
		if !seen[rel] {
			seen[rel] = true
			out = append(out, rel)
		}
	}
	return out, nil
}

// resolvePath resolves symlinks in the longest prefix of path that exists.
func resolvePath(path string) string {
	rest := ""
	for p := path; ; p = filepath.Dir(p) {
		if r, err := filepath.EvalSymlinks(p); err == nil {
			return filepath.Join(r, rest)
		}
		if filepath.Dir(p) == p {
			return path
		}
		rest = filepath.Join(filepath.Base(p), rest)
	}
}

// headCRDFiles returns the YAML files directly in dirs: per dir, in the
// order given, *.yaml then *.yml, each sorted. Text output keeps this order;
// the JSON document is sorted on its own.
func headCRDFiles(dirs []string) ([]string, error) {
	var files []string
	for _, dir := range dirs {
		for _, ext := range []string{"*.yaml", "*.yml"} {
			matches, err := filepath.Glob(filepath.Join(dir, ext))
			if err != nil {
				return nil, fmt.Errorf("globbing %s: %w", dir, err)
			}
			for _, m := range matches {
				files = append(files, filepath.ToSlash(m))
			}
		}
	}
	return files, nil
}

// baseCRDFiles returns the YAML files directly in dirs (root-relative, from
// repoDirs) at ref. A directory missing at ref has none.
func baseCRDFiles(ref string, dirs []string) (map[string]bool, error) {
	files := map[string]bool{}
	for _, dir := range dirs {
		pathspec := dir + "/"
		if dir == "." {
			pathspec = "."
		}
		// -z: names raw and NUL-terminated (git otherwise quotes non-ASCII
		// names); --full-name: root-relative whatever the cwd.
		out, err := gitOutput("ls-tree", "-z", "--name-only", "--full-name", ref, "--", pathspec)
		if err != nil {
			return nil, fmt.Errorf("listing %s at %s: %w", dir, ref, err)
		}
		for _, name := range strings.Split(string(out), "\x00") {
			if ext := filepath.Ext(name); (ext == ".yaml" || ext == ".yml") && filepath.ToSlash(filepath.Dir(name)) == dir {
				files[name] = true
			}
		}
	}
	return files, nil
}

// baseCRD reads file at ref and reports whether it is a CRD there.
func baseCRD(ref, file string) (v1.CustomResourceDefinition, bool, error) {
	var crd v1.CustomResourceDefinition
	data, err := gitShow(ref, file)
	if err != nil {
		return crd, false, err
	}
	if data, err = onlyDocument(file, data); err != nil {
		return crd, false, fmt.Errorf("%w at %s", err, ref)
	}
	if err := yaml.Unmarshal(data, &crd); err != nil || crd.Kind != "CustomResourceDefinition" {
		return crd, false, nil
	}
	return crd, true, nil
}

// onlyDocument returns the one non-empty YAML document in data (or an empty
// one). yaml.Unmarshal reads only the first document, so a second CRD in a
// file would go unchecked: more than one non-empty document is an error.
// A document holding only comments or whitespace is empty.
func onlyDocument(file string, data []byte) ([]byte, error) {
	r := k8syaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(data)))
	var found []byte
	n := 0
	for {
		doc, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", file, err)
		}
		if emptyDocument(doc) {
			continue
		}
		if n++; n > 1 {
			return nil, fmt.Errorf("multiple YAML documents in %s; one CRD per file is required", file)
		}
		found = doc
	}
	return found, nil
}

func emptyDocument(doc []byte) bool {
	for _, line := range strings.Split(string(doc), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && line != "---" && !strings.HasPrefix(line, "#") {
			return false
		}
	}
	return true
}

// verifyCommit fails unless ref names a commit. Without this check a typo in
// --base-ref made every file NEW, and the check passed.
func verifyCommit(ref string) error {
	if ref == "" || strings.HasPrefix(ref, "-") {
		return fmt.Errorf("base ref %q is not a commit in this repository", ref)
	}
	if _, err := gitOutput("rev-parse", "--verify", "--quiet", ref+"^{commit}"); err != nil {
		return fmt.Errorf("base ref %q is not a commit in this repository", ref)
	}
	return nil
}

// firstSchema returns the first version's openAPIV3Schema, or nil.
func firstSchema(crd v1.CustomResourceDefinition) *v1.JSONSchemaProps {
	if len(crd.Spec.Versions) == 0 || crd.Spec.Versions[0].Schema == nil {
		return nil
	}
	return crd.Spec.Versions[0].Schema.OpenAPIV3Schema
}

func sortedKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func gitShow(ref, path string) ([]byte, error) {
	out, err := gitOutput("show", ref+":"+path)
	if err != nil {
		return nil, fmt.Errorf("reading %s at %s: %w", path, ref, err)
	}
	return out, nil
}

// gitOutput runs git and returns its stdout. A failure includes git's stderr.
func gitOutput(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && stderr.Len() > 0 {
			return nil, fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(stderr.String()))
		}
		return nil, fmt.Errorf("git %s: %w", args[0], err)
	}
	return out, nil
}
