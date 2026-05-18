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

package ack

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	rbacv1 "k8s.io/api/rbac/v1"

	acksdk "github.com/aws-controllers-k8s/code-generator/pkg/sdk"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

const (
	goimportsPackage = "golang.org/x/tools/cmd/goimports@latest"
)

// BuildControllerOptions contains all the options needed to run the full
// controller generation pipeline.
type BuildControllerOptions struct {
	// SvcAlias is the AWS service alias (e.g., "s3", "ec2")
	SvcAlias string
	// ControllerSourcePath is the root of the service controller repo
	ControllerSourcePath string
	// APIVersion is the Kubernetes API version (e.g., "v1alpha1")
	APIVersion string
	// RBACRoleName is the name of the RBAC ClusterRole
	RBACRoleName string
	// RuntimeVersion is the runtime module version for CRD generation
	RuntimeVersion string
	// ControllerGenVersion is the required version of controller-gen
	ControllerGenVersion string
	// CacheDir is the path to the ack-generate cache directory
	CacheDir string
	// BoilerplatePath is the path to the boilerplate header file
	BoilerplatePath string
	// BoilerplateFS is the embedded filesystem containing boilerplate files
	BoilerplateFS fs.FS
	// TemplatesFS is the embedded filesystem containing templates
	TemplatesFS fs.FS
}

// BuildControllerPreCodegen runs the pipeline steps that must happen before
// code generation: runtime CRDs, deepcopy, and CRD manifests.
func BuildControllerPreCodegen(ctx context.Context, opts BuildControllerOptions) error {
	if err := checkControllerGen(opts.ControllerGenVersion); err != nil {
		return err
	}

	configDir := filepath.Join(opts.ControllerSourcePath, "config")
	apisDir := filepath.Join(opts.ControllerSourcePath, "apis", opts.APIVersion)

	boilerplatePath := opts.BoilerplatePath
	if boilerplatePath == "" {
		svcBoilerplate := filepath.Join(opts.ControllerSourcePath, "templates", "boilerplate.txt")
		if util.FileExists(svcBoilerplate) {
			boilerplatePath = svcBoilerplate
		} else if opts.TemplatesFS != nil {
			data, err := fs.ReadFile(opts.TemplatesFS, "boilerplate.txt")
			if err != nil {
				return fmt.Errorf("reading embedded boilerplate.txt: %w", err)
			}
			tmpFile, err := os.CreateTemp("", "boilerplate-*.txt")
			if err != nil {
				return fmt.Errorf("creating temp boilerplate file: %w", err)
			}
			defer os.Remove(tmpFile.Name())
			if _, err := tmpFile.Write(data); err != nil {
				tmpFile.Close()
				return fmt.Errorf("writing temp boilerplate file: %w", err)
			}
			tmpFile.Close()
			boilerplatePath = tmpFile.Name()
		}
	}

	// Step 1: Copy runtime CRD config (preserves bases/ + kustomization.yaml)
	fmt.Println("Copying common custom resource definitions")
	commonCRDDir := filepath.Join(configDir, "crd", "common")
	if err := os.MkdirAll(commonCRDDir, 0755); err != nil {
		return fmt.Errorf("creating common CRD directory: %w", err)
	}
	if err := copyRuntimeCRDConfig(ctx, opts.CacheDir, opts.RuntimeVersion, commonCRDDir); err != nil {
		return fmt.Errorf("copying runtime CRD config: %w", err)
	}

	// Step 2: controller-gen object (deepcopy)
	fmt.Printf("Generating deepcopy code for %s\n", opts.SvcAlias)
	objectArgs := []string{"object:headerFile=" + boilerplatePath, "paths=./..."}
	if err := runControllerGen(objectArgs, apisDir); err != nil {
		return fmt.Errorf("generating deepcopy code: %w", err)
	}

	// Step 3: controller-gen crd
	fmt.Printf("Generating custom resource definitions for %s\n", opts.SvcAlias)
	crdArgs := []string{
		"crd:allowDangerousTypes=true",
		"paths=./...",
		"output:crd:artifacts:config=" + filepath.Join(configDir, "crd", "bases"),
	}
	if err := runControllerGen(crdArgs, apisDir); err != nil {
		return fmt.Errorf("generating CRDs: %w", err)
	}

	return nil
}

// BuildControllerPostCodegen runs the pipeline steps that happen after code
// generation: go mod tidy, RBAC, formatting, and boilerplate copying.
func BuildControllerPostCodegen(opts BuildControllerOptions) error {
	configDir := filepath.Join(opts.ControllerSourcePath, "config")
	resourceDir := filepath.Join(opts.ControllerSourcePath, "pkg", "resource")

	// Step 1: go mod tidy
	fmt.Println("Running go mod tidy")
	if err := runCommand("go", []string{"mod", "tidy"}, opts.ControllerSourcePath); err != nil {
		return fmt.Errorf("running go mod tidy: %w", err)
	}

	// Step 2: controller-gen rbac
	fmt.Printf("Generating RBAC manifests for %s\n", opts.SvcAlias)
	rbacArgs := []string{
		"rbac:roleName=" + opts.RBACRoleName,
		"paths=./...",
		"output:rbac:artifacts:config=" + filepath.Join(configDir, "rbac"),
	}
	if err := runControllerGen(rbacArgs, resourceDir); err != nil {
		return fmt.Errorf("generating RBAC manifests: %w", err)
	}

	// Step 3: rename role.yaml → cluster-role-controller.yaml
	roleYAML := filepath.Join(configDir, "rbac", "role.yaml")
	clusterRoleYAML := filepath.Join(configDir, "rbac", "cluster-role-controller.yaml")
	if err := os.Rename(roleYAML, clusterRoleYAML); err != nil {
		return fmt.Errorf("renaming role.yaml: %w", err)
	}

	// Step 4: copy namespaced overlay patches from embedded FS
	if opts.TemplatesFS != nil {
		fmt.Println("Copying namespaced overlay patches")
		if err := copyNamespacedOverlays(opts.TemplatesFS, configDir); err != nil {
			return fmt.Errorf("copying namespaced overlays: %w", err)
		}
	}

	// Step 5: gofmt + goimports
	fmt.Printf("Running formatters against generated code for %s\n", opts.SvcAlias)
	if err := runFormatters(opts.ControllerSourcePath); err != nil {
		return fmt.Errorf("running formatters: %w", err)
	}

	// Step 6: copy boilerplate files
	if opts.BoilerplateFS != nil {
		fmt.Println("Updating repository maintenance files")
		if err := copyBoilerplate(opts.BoilerplateFS, opts.ControllerSourcePath); err != nil {
			return fmt.Errorf("copying boilerplate files: %w", err)
		}
	}

	return nil
}

// BuildReleaseOptions contains all the options needed to run the full
// release artifact generation pipeline.
type BuildReleaseOptions struct {
	// SvcAlias is the AWS service alias (e.g., "s3", "ec2")
	SvcAlias string
	// ControllerSourcePath is the root of the service controller repo
	ControllerSourcePath string
	// APIVersion is the Kubernetes API version (e.g., "v1alpha1")
	APIVersion string
	// RBACRoleName is the name of the RBAC ClusterRole
	RBACRoleName string
	// RuntimeVersion is the runtime module version for CRD generation
	RuntimeVersion string
	// ControllerGenVersion is the required version of controller-gen
	ControllerGenVersion string
	// CacheDir is the path to the ack-generate cache directory
	CacheDir string
}

// BuildRelease runs the post-processing steps for release artifact generation:
//  1. Copy runtime CRDs → helm/crds
//  2. controller-gen crd (service CRDs → helm/crds)
//  3. controller-gen rbac (→ helm/templates)
//  4. RBAC rules injection into _helpers.tpl
func BuildRelease(ctx context.Context, opts BuildReleaseOptions) error {
	if err := checkControllerGen(opts.ControllerGenVersion); err != nil {
		return err
	}

	helmDir := filepath.Join(opts.ControllerSourcePath, "helm")
	apisDir := filepath.Join(opts.ControllerSourcePath, "apis", opts.APIVersion)
	resourceDir := filepath.Join(opts.ControllerSourcePath, "pkg", "resource")

	// Step 1: Copy runtime CRDs → helm/crds
	fmt.Println("Copying common custom resource definitions")
	crdsDir := filepath.Join(helmDir, "crds")
	if err := copyRuntimeCRDBases(ctx, opts.CacheDir, opts.RuntimeVersion, crdsDir); err != nil {
		return fmt.Errorf("copying runtime CRDs for helm: %w", err)
	}

	// Step 2: controller-gen crd (service CRDs → helm/crds)
	fmt.Printf("Generating custom resource definitions for %s\n", opts.SvcAlias)
	crdArgs := []string{
		"crd:allowDangerousTypes=true",
		"paths=./...",
		"output:crd:artifacts:config=" + crdsDir,
	}
	if err := runControllerGen(crdArgs, apisDir); err != nil {
		return fmt.Errorf("generating service CRDs for helm: %w", err)
	}

	// Step 3: controller-gen rbac (→ helm/templates)
	fmt.Printf("Generating RBAC manifests for %s\n", opts.SvcAlias)
	rbacArgs := []string{
		"rbac:roleName=" + opts.RBACRoleName,
		"paths=./...",
		"output:rbac:artifacts:config=" + filepath.Join(helmDir, "templates"),
	}
	if err := runControllerGen(rbacArgs, resourceDir); err != nil {
		return fmt.Errorf("generating RBAC manifests for helm: %w", err)
	}

	// Step 4: Inject RBAC rules into _helpers.tpl
	fmt.Println("Injecting RBAC rules into Helm templates")
	if err := injectRBACRules(helmDir); err != nil {
		return fmt.Errorf("injecting RBAC rules: %w", err)
	}

	return nil
}

// copyRuntimeCRDConfig fetches the runtime CRD config files from GitHub
// (cached locally) and copies them to the controller's config/crd/common/
// directory, preserving the directory structure (bases/ + kustomization.yaml).
func copyRuntimeCRDConfig(ctx context.Context, cacheDir string, runtimeVersion string, destDir string) error {
	cachedDir, err := acksdk.EnsureRuntimeCRDs(ctx, cacheDir, runtimeVersion)
	if err != nil {
		return err
	}

	for _, relPath := range acksdk.RuntimeCRDFiles() {
		srcPath := filepath.Join(cachedDir, relPath)
		destPath := filepath.Join(destDir, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(destPath, data, 0666); err != nil {
			return err
		}
	}
	return nil
}

// copyRuntimeCRDBases fetches the runtime CRD files and copies only the
// bases/*.yaml files to the destination (used for helm/crds/ which needs
// flat CRD files without kustomization).
func copyRuntimeCRDBases(ctx context.Context, cacheDir string, runtimeVersion string, destDir string) error {
	cachedDir, err := acksdk.EnsureRuntimeCRDs(ctx, cacheDir, runtimeVersion)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	for _, relPath := range acksdk.RuntimeCRDFiles() {
		if !strings.HasPrefix(relPath, "bases/") {
			continue
		}
		srcPath := filepath.Join(cachedDir, relPath)
		destPath := filepath.Join(destDir, filepath.Base(relPath))
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return err
		}
		if err := os.WriteFile(destPath, data, 0666); err != nil {
			return err
		}
	}
	return nil
}

// checkControllerGen verifies that controller-gen is installed and at the
// required version.
func checkControllerGen(requiredVersion string) error {
	path, err := exec.LookPath("controller-gen")
	if err != nil {
		return fmt.Errorf(
			"controller-gen not found in PATH. Install it with:\n"+
				"  go install sigs.k8s.io/controller-tools/cmd/controller-gen@%s",
			requiredVersion,
		)
	}

	out, err := exec.Command(path, "--version").Output()
	if err != nil {
		return fmt.Errorf("checking controller-gen version: %w", err)
	}
	version := strings.TrimSpace(string(out))
	if !strings.Contains(version, requiredVersion) {
		return fmt.Errorf(
			"controller-gen version mismatch: have %q, need %s.\n"+
				"Install the correct version with:\n"+
				"  go install sigs.k8s.io/controller-tools/cmd/controller-gen@%s",
			version, requiredVersion, requiredVersion,
		)
	}
	return nil
}

// ensureGoimports checks if goimports is installed and installs it if not.
func ensureGoimports() error {
	if _, err := exec.LookPath("goimports"); err == nil {
		return nil
	}
	fmt.Println("Installing goimports...")
	cmd := exec.Command("go", "install", goimportsPackage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runControllerGen executes controller-gen with the given arguments in the
// specified directory.
func runControllerGen(args []string, dir string) error {
	return runCommand("controller-gen", args, dir)
}

// runCommand executes a command with the given arguments in the specified
// directory, forwarding stdout/stderr.
func runCommand(name string, args []string, dir string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runFormatters runs gofmt and goimports on the given directory.
func runFormatters(dir string) error {
	if err := runCommand("gofmt", []string{"-w", dir}, ""); err != nil {
		return fmt.Errorf("running gofmt: %w", err)
	}

	if err := ensureGoimports(); err != nil {
		return fmt.Errorf("ensuring goimports is installed: %w", err)
	}

	if err := runCommand("goimports", []string{"-w", dir}, ""); err != nil {
		return fmt.Errorf("running goimports: %w", err)
	}
	return nil
}

// injectRBACRules reads the controller-gen-generated role.yaml, extracts the
// RBAC rules, and injects them into _helpers.tpl replacing the
// SEDREPLACERULES marker.
func injectRBACRules(helmDir string) error {
	roleYAMLPath := filepath.Join(helmDir, "templates", "role.yaml")
	helpersPath := filepath.Join(helmDir, "templates", "_helpers.tpl")

	roleData, err := os.ReadFile(roleYAMLPath)
	if err != nil {
		return fmt.Errorf("reading role.yaml: %w", err)
	}

	var role rbacv1.ClusterRole
	if err := yaml.Unmarshal(roleData, &role); err != nil {
		return fmt.Errorf("unmarshaling role.yaml: %w", err)
	}

	// Marshal just the rules back to YAML
	rulesYAML, err := yaml.Marshal(role.Rules)
	if err != nil {
		return fmt.Errorf("marshaling RBAC rules: %w", err)
	}

	// The rules need to be formatted as they appear in a ClusterRole spec,
	// with "rules:" prefix
	var rulesBlock bytes.Buffer
	rulesBlock.WriteString("rules:\n")
	for _, line := range strings.Split(strings.TrimSpace(string(rulesYAML)), "\n") {
		rulesBlock.WriteString(line)
		rulesBlock.WriteString("\n")
	}

	helpersData, err := os.ReadFile(helpersPath)
	if err != nil {
		return fmt.Errorf("reading _helpers.tpl: %w", err)
	}

	helpersStr := string(helpersData)
	result := strings.Replace(
		helpersStr,
		"SEDREPLACERULES",
		strings.TrimSpace(rulesBlock.String()),
		1,
	)
	if result == helpersStr {
		return fmt.Errorf("SEDREPLACERULES marker not found in %s; ensure _helpers.tpl contains the marker", helpersPath)
	}

	if err := os.WriteFile(helpersPath, []byte(result), 0666); err != nil {
		return fmt.Errorf("writing _helpers.tpl: %w", err)
	}

	// Remove the original role.yaml
	return os.Remove(roleYAMLPath)
}

// copyBoilerplate writes embedded boilerplate files (LICENSE, CONTRIBUTING.md,
// etc.) to the controller repo.
func copyBoilerplate(fsys fs.FS, destDir string) error {
	files := []string{
		"CODE_OF_CONDUCT.md",
		"CONTRIBUTING.md",
		"GOVERNANCE.md",
		"LICENSE",
		"NOTICE",
	}
	for _, name := range files {
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return fmt.Errorf("reading embedded %s: %w", name, err)
		}
		destPath := filepath.Join(destDir, name)
		if err := os.WriteFile(destPath, data, 0666); err != nil {
			return fmt.Errorf("writing %s: %w", name, err)
		}
	}
	return nil
}

// copyNamespacedOverlays copies the namespaced overlay JSON patches from
// the embedded templates FS to the controller's config directory.
func copyNamespacedOverlays(fsys fs.FS, configDir string) error {
	overlayDir := filepath.Join(configDir, "overlays", "namespaced")
	if err := os.MkdirAll(overlayDir, 0755); err != nil {
		return err
	}

	files := []string{
		"role-binding.json",
		"role.json",
	}
	for _, name := range files {
		srcPath := filepath.Join("config", "overlays", "namespaced", name)
		data, err := fs.ReadFile(fsys, filepath.ToSlash(srcPath))
		if err != nil {
			return fmt.Errorf("reading embedded %s: %w", srcPath, err)
		}
		destPath := filepath.Join(overlayDir, name)
		if err := os.WriteFile(destPath, data, 0666); err != nil {
			return fmt.Errorf("writing %s: %w", name, err)
		}
	}
	return nil
}

// ResolveConfigPaths checks for config files at the controller source path
// and returns the resolved paths. Only returns paths for files that exist.
type ResolvedConfigPaths struct {
	GeneratorConfigPath     string
	MetadataConfigPath      string
	DocumentationConfigPath string
}

// ResolveConfigPaths checks for config files at the controller source path,
// returning paths only for files that exist on disk.
func ResolveConfigPaths(controllerSourcePath string) ResolvedConfigPaths {
	resolved := ResolvedConfigPaths{}
	candidates := []struct {
		filename string
		target   *string
	}{
		{"generator.yaml", &resolved.GeneratorConfigPath},
		{"metadata.yaml", &resolved.MetadataConfigPath},
		{"documentation.yaml", &resolved.DocumentationConfigPath},
	}
	for _, c := range candidates {
		path := filepath.Join(controllerSourcePath, c.filename)
		if util.FileExists(path) {
			*c.target = path
		}
	}
	return resolved
}
