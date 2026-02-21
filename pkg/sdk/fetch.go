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
	"time"

	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

const (
	sdkModelURLTemplate    = "https://raw.githubusercontent.com/aws/aws-sdk-go-v2/%s/codegen/sdk-codegen/aws-models/%s.json"
	defaultHTTPFetchTimeout = 60 * time.Second
)

// EnsureModel ensures that we have a locally-cached copy of the AWS SDK model
// JSON file for a given service and SDK version. If the file is already cached,
// it returns immediately. Otherwise, it fetches the file from GitHub.
//
// The returned string is the base path to use with NewHelper — it mirrors the
// SDK repo directory structure so that ModelAndDocsPath works unchanged.
func EnsureModel(
	ctx context.Context,
	cacheDir string,
	sdkVersion string,
	modelName string,
) (string, error) {
	totalStart := time.Now()

	basePath := filepath.Join(cacheDir, "models", sdkVersion)
	modelDir := filepath.Join(basePath, "codegen", "sdk-codegen", "aws-models")
	modelPath := filepath.Join(modelDir, fmt.Sprintf("%s.json", modelName))

	// Check cache first
	if _, err := os.Stat(modelPath); err == nil {
		util.Tracef("EnsureModel: cache hit for %s@%s\n", modelName, sdkVersion)
		return basePath, nil
	}

	// Create the cache directory
	if err := os.MkdirAll(modelDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("cannot create model cache directory %s: %v", modelDir, err)
	}

	// Fetch from GitHub
	url := fmt.Sprintf(sdkModelURLTemplate, sdkVersion, modelName)
	util.Tracef("EnsureModel: fetching %s\n", url)

	fetchStart := time.Now()
	ctx, cancel := context.WithTimeout(ctx, defaultHTTPFetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("cannot create HTTP request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot fetch model file from %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch model file from %s: HTTP %d", url, resp.StatusCode)
	}

	// Write to a temp file, then rename atomically into the cache path
	tmpFile, err := os.CreateTemp(modelDir, ".tmp-"+modelName+"-*.json")
	if err != nil {
		return "", fmt.Errorf("cannot create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()

	_, err = io.Copy(tmpFile, resp.Body)
	closeErr := tmpFile.Close()
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("cannot write model file: %v", err)
	}
	if closeErr != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("cannot close temp file: %v", closeErr)
	}

	if err := os.Rename(tmpPath, modelPath); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("cannot rename temp file to cache path: %v", err)
	}

	util.Tracef("EnsureModel: fetched %s (%s)\n", modelName, time.Since(fetchStart))
	util.Tracef("EnsureModel total: %s\n", time.Since(totalStart))
	return basePath, nil
}
