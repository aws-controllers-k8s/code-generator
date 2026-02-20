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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"

	crdcompat "github.com/kubernetes-sigs/kro/pkg/graph/crd/compat"
)

var (
	optBaseRef  string
	optCRDPaths []string
)

var crdCompatCheckCmd = &cobra.Command{
	Use:   "crd-compat-check",
	Short: "Check CRDs for breaking changes against a git base ref",
	Long: `Compares CRD files on disk against the same files at a given git ref
and reports any breaking schema changes. Exits non-zero if breaking
changes are detected.

By default it checks config/crd/bases and helm/crds against the
main branch.`,
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
	rootCmd.AddCommand(crdCompatCheckCmd)
}

func checkCRDCompat(cmd *cobra.Command, args []string) error {
	var crdFiles []string
	for _, dir := range optCRDPaths {
		for _, ext := range []string{"*.yaml", "*.yml"} {
			matches, err := filepath.Glob(filepath.Join(dir, ext))
			if err != nil {
				return fmt.Errorf("globbing %s: %w", dir, err)
			}
			crdFiles = append(crdFiles, matches...)
		}
	}

	if len(crdFiles) == 0 {
		fmt.Println("No CRD files found, nothing to check.")
		return nil
	}

	hasBreaking := false
	checked := 0

	for _, crdFile := range crdFiles {
		newData, err := os.ReadFile(crdFile)
		if err != nil {
			return fmt.Errorf("reading %s: %w", crdFile, err)
		}

		// Verify it's actually a CRD before proceeding
		var newCRD v1.CustomResourceDefinition
		if err := yaml.Unmarshal(newData, &newCRD); err != nil {
			fmt.Printf("SKIP %s: not a valid CRD (%v)\n", crdFile, err)
			continue
		}
		if newCRD.Kind != "CustomResourceDefinition" {
			fmt.Printf("SKIP %s: not a CRD (kind=%s)\n", crdFile, newCRD.Kind)
			continue
		}

		oldData, err := gitShow(optBaseRef, crdFile)
		if err != nil {
			fmt.Printf("NEW  %s (not in %s)\n", crdFile, optBaseRef)
			continue
		}

		var oldCRD v1.CustomResourceDefinition
		if err := yaml.Unmarshal(oldData, &oldCRD); err != nil {
			return fmt.Errorf("parsing old %s at %s: %w", crdFile, optBaseRef, err)
		}

		report, err := crdcompat.CompareVersions(
			oldCRD.Spec.Versions, newCRD.Spec.Versions,
		)
		if err != nil {
			return fmt.Errorf("comparing %s: %w", crdFile, err)
		}

		checked++
		if report.HasBreakingChanges() {
			hasBreaking = true
			fmt.Printf("FAIL %s\n", crdFile)
			for _, c := range report.BreakingChanges {
				fmt.Printf("  BREAKING: [%s] %s\n", c.Path, c.Description())
			}
			for _, c := range report.NonBreakingChanges {
				fmt.Printf("  ok:       [%s] %s\n", c.Path, c.Description())
			}
		} else if report.HasChanges() {
			fmt.Printf("OK   %s (non-breaking changes only)\n", crdFile)
		} else {
			fmt.Printf("OK   %s (no changes)\n", crdFile)
		}
	}

	fmt.Printf("\nChecked %d CRD(s) across %d file(s).\n", checked, len(crdFiles))
	if hasBreaking {
		return fmt.Errorf("breaking CRD changes detected")
	}
	return nil
}

func gitShow(ref, path string) ([]byte, error) {
	out, err := exec.Command("git", "show", ref+":"+path).Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}
