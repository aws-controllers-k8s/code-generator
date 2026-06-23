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

package sdk

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"

	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

const (
	runtimeCRDURLTemplate = "https://raw.githubusercontent.com/aws-controllers-k8s/runtime/%s/config/crd/%s"
)

// runtimeCRDFiles is the set of files we need from the runtime config/crd/
// directory. The keys are relative paths within config/crd/.
var runtimeCRDFiles = []string{
	"kustomization.yaml",
	"bases/services.k8s.aws_fieldexports.yaml",
	"bases/services.k8s.aws_iamroleselectors.yaml",
}

// EnsureRuntimeCRDs ensures that we have a locally-cached copy of the runtime
// CRD config files for a given runtime version. If the files are already
// cached, it returns immediately. Otherwise, it fetches them from GitHub.
//
// The returned string is the path to the cached config/crd/ directory.
func EnsureRuntimeCRDs(
	ctx context.Context,
	cacheDir string,
	runtimeVersion string,
) (string, error) {
	crdDir := filepath.Join(cacheDir, "runtime", runtimeVersion, "config", "crd")

	// Check if all files already exist (cache hit)
	allCached := true
	for _, relPath := range runtimeCRDFiles {
		if _, err := os.Stat(filepath.Join(crdDir, relPath)); err != nil {
			allCached = false
			break
		}
	}
	if allCached {
		util.Tracef("EnsureRuntimeCRDs: cache hit for runtime@%s\n", runtimeVersion)
		return crdDir, nil
	}

	// Fetch missing files from GitHub
	for _, relPath := range runtimeCRDFiles {
		destPath := filepath.Join(crdDir, relPath)
		if _, err := os.Stat(destPath); err == nil {
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return "", fmt.Errorf("cannot create runtime CRD cache directory: %v", err)
		}

		url := fmt.Sprintf(runtimeCRDURLTemplate, runtimeVersion, relPath)
		util.Tracef("EnsureRuntimeCRDs: fetching %s\n", url)

		if err := fetchFile(ctx, url, destPath); err != nil {
			return "", fmt.Errorf("fetching runtime CRD file %s: %w", relPath, err)
		}
	}

	util.Tracef("EnsureRuntimeCRDs: cached runtime@%s CRDs\n", runtimeVersion)
	return crdDir, nil
}

// fetchFile downloads a URL and writes it atomically to destPath.
func fetchFile(ctx context.Context, url string, destPath string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultHTTPFetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("cannot create HTTP request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot fetch %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch %s: HTTP %d", url, resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(destPath), ".tmp-*")
	if err != nil {
		return fmt.Errorf("cannot create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()

	_, err = io.Copy(tmpFile, resp.Body)
	closeErr := tmpFile.Close()
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("cannot write file: %v", err)
	}
	if closeErr != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("cannot close temp file: %v", closeErr)
	}

	if err := os.Rename(tmpPath, destPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("cannot rename temp file to cache path: %v", err)
	}

	return nil
}

// GetRuntimeVersion returns the runtime module version from a controller's
// go.mod file.
func GetRuntimeVersion(controllerRepoPath string) (string, error) {
	goModPath := filepath.Join(controllerRepoPath, "go.mod")
	b, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("reading go.mod: %w", err)
	}

	goMod, err := modfile.Parse("", b, nil)
	if err != nil {
		return "", fmt.Errorf("parsing go.mod: %w", err)
	}

	const runtimeModule = "github.com/aws-controllers-k8s/runtime"
	for _, require := range goMod.Require {
		if require.Mod.Path == runtimeModule {
			return require.Mod.Version, nil
		}
	}
	return "", fmt.Errorf("runtime module not found in %s", goModPath)
}

// RuntimeCRDFiles returns the list of runtime CRD files that are cached.
// This is exported for use in the pipeline logic that copies these files
// to the controller's config/crd/common/ directory.
func RuntimeCRDFiles() []string {
	return runtimeCRDFiles
}
