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
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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
ref but not on disk is a breaking change. Exits non-zero if breaking
changes are detected.

By default it checks config/crd/bases and helm/crds against the
main branch. Run it from the repository root.

--output json prints one ack.crd-compat/v1 document to stdout: every
change sorted into buckets (added_optional_spec, added_status,
added_required, type_changed, removed, enum_changed,
validation_tightened, validation_relaxed, default_changed,
description_only, crd_added, crd_removed, other), and the added fields
whose names suggest a secret (sensitive_candidates).`,
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
		var newCRD v1.CustomResourceDefinition
		if err := yaml.Unmarshal(newData, &newCRD); err != nil {
			reason := fmt.Sprintf("not a valid CRD (%v)", err)
			res.AddFile(crdcompat.FileResult{File: file, Status: crdcompat.FileSkipped, Reason: reason})
			if text {
				fmt.Fprintf(w, "SKIP %s: %s\n", file, reason)
			}
			continue
		}
		if newCRD.Kind != "CustomResourceDefinition" {
			reason := fmt.Sprintf("not a CRD (kind=%s)", newCRD.Kind)
			res.AddFile(crdcompat.FileResult{File: file, Status: crdcompat.FileSkipped, Reason: reason})
			if text {
				fmt.Fprintf(w, "SKIP %s: %s\n", file, reason)
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
		var oldCRD v1.CustomResourceDefinition
		if err := yaml.Unmarshal(oldData, &oldCRD); err != nil {
			return false, fmt.Errorf("parsing old %s at %s: %w", file, baseRef, err)
		}
		report, err := crdcompatkro.CompareVersions(oldCRD.Spec.Versions, newCRD.Spec.Versions)
		if err != nil {
			return false, fmt.Errorf("comparing %s: %w", file, err)
		}

		checked++
		changes := crdcompat.Classify(file, kind, report)
		res.AddChanges(changes...)
		res.AddSensitive(crdcompat.SensitiveFields(file, kind, firstSchema(newCRD), changes)...)
		status := crdcompat.FileUnchanged
		if report.HasChanges() {
			status = crdcompat.FileChanged
		}
		res.AddFile(crdcompat.FileResult{File: file, Kind: kind, Status: status})
		if !text {
			continue
		}
		switch {
		case report.HasBreakingChanges():
			fmt.Fprintf(w, "FAIL %s\n", file)
			for _, c := range report.BreakingChanges {
				fmt.Fprintf(w, "  BREAKING: [%s] %s\n", c.Path, c.Description())
			}
			for _, c := range report.NonBreakingChanges {
				fmt.Fprintf(w, "  ok:       [%s] %s\n", c.Path, c.Description())
			}
		case report.HasChanges():
			fmt.Fprintf(w, "OK   %s (non-breaking changes only)\n", file)
		default:
			fmt.Fprintf(w, "OK   %s (no changes)\n", file)
		}
	}

	for _, file := range sortedKeys(base) {
		if onDisk[file] {
			continue
		}
		oldData, err := gitShow(baseRef, file)
		if err != nil {
			return false, err
		}
		var oldCRD v1.CustomResourceDefinition
		if err := yaml.Unmarshal(oldData, &oldCRD); err != nil || oldCRD.Kind != "CustomResourceDefinition" {
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

// headCRDFiles returns the YAML files directly in dirs, sorted.
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
	sort.Strings(files)
	return files, nil
}

// baseCRDFiles returns the YAML files directly in dirs at ref. A directory
// missing at ref has none.
func baseCRDFiles(ref string, dirs []string) (map[string]bool, error) {
	files := map[string]bool{}
	for _, dir := range dirs {
		dir = filepath.ToSlash(filepath.Clean(dir))
		out, err := gitOutput("ls-tree", "--name-only", ref, "--", dir+"/")
		if err != nil {
			return nil, fmt.Errorf("listing %s at %s: %w", dir, ref, err)
		}
		for _, name := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if ext := filepath.Ext(name); (ext == ".yaml" || ext == ".yml") && filepath.ToSlash(filepath.Dir(name)) == dir {
				files[name] = true
			}
		}
	}
	return files, nil
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
