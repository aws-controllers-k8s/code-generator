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
	"strconv"
	"time"

	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

const (
	sdkModelURLTemplate        = "https://raw.githubusercontent.com/aws/aws-sdk-go-v2/%s/codegen/sdk-codegen/aws-models/%s.json"
	perServiceModelURLTemplate = "https://raw.githubusercontent.com/aws/aws-sdk-go-v2/service/%s/%s/codegen/sdk-codegen/aws-models/%s.json"
	defaultHTTPFetchTimeout    = 60 * time.Second
	// maxFetchAttempts bounds how many times a model file is requested when
	// GitHub answers with a transient status (429 or 5xx).
	maxFetchAttempts = 6
)

// Backoff between model fetch attempts. The delay doubles after each attempt,
// starting at fetchRetryBaseDelay, and never exceeds fetchRetryMaxDelay, even
// when the server's Retry-After asks for longer. These are variables so tests
// can shorten them.
var (
	fetchRetryBaseDelay = 2 * time.Second
	fetchRetryMaxDelay  = 60 * time.Second
)

// EnsureModel ensures that we have a locally-cached copy of the AWS SDK model
// JSON file for a given service and SDK version. If the file is already cached,
// it returns immediately. Otherwise, it fetches the file from GitHub.
//
// The fetch strategy is determined by serviceSDKVersion:
//   - When serviceSDKVersion is non-empty: check per-service cache, then fetch
//     per-service tag URL. On non-200, return an error — no core fallback.
//   - When serviceSDKVersion is empty: check core cache, then fetch core tag
//     URL. On 404, suggest --aws-service-sdk-version. On other non-200, return
//     an error with the URL and status code.
//
// svcPackageName is the aws-sdk-go-v2 service package directory name (e.g.
// "route53"), used as the per-service tag path segment. It can differ from
// modelName, which identifies the model JSON file (e.g. "route-53"). It is only
// consulted on the per-service path.
//
// The returned string is the base path to use with NewHelper — it mirrors the
// SDK repo directory structure so that ModelAndDocsPath works unchanged.
func EnsureModel(
	ctx context.Context,
	cacheDir string,
	sdkVersion string,
	svcPackageName string,
	modelName string,
	serviceSDKVersion string,
) (string, error) {
	totalStart := time.Now()

	if serviceSDKVersion != "" {
		return ensureModelPerService(ctx, cacheDir, svcPackageName, modelName, serviceSDKVersion, totalStart)
	}
	return ensureModelCore(ctx, cacheDir, sdkVersion, modelName, totalStart)
}

// ensureModelPerService handles the per-service-only fetch path.
// It checks the per-service cache first, then fetches from the per-service tag
// URL. It never falls back to the core SDK tag.
//
// The per-service tag is named after the aws-sdk-go-v2 service package
// (svcPackageName, e.g. "route53"), which can differ from the model file name
// (modelName, e.g. "route-53"). When svcPackageName is empty we fall back to
// modelName so existing callers keep working.
func ensureModelPerService(
	ctx context.Context,
	cacheDir string,
	svcPackageName string,
	modelName string,
	serviceSDKVersion string,
	totalStart time.Time,
) (string, error) {
	if svcPackageName == "" {
		svcPackageName = modelName
	}
	normalizedSvcVer := EnsureSemverPrefix(serviceSDKVersion)
	basePath := filepath.Join(cacheDir, "models", "service", svcPackageName, normalizedSvcVer)
	modelDir := filepath.Join(basePath, "codegen", "sdk-codegen", "aws-models")
	modelPath := filepath.Join(modelDir, fmt.Sprintf("%s.json", modelName))

	// Check per-service cache
	if _, err := os.Stat(modelPath); err == nil {
		util.Tracef("EnsureModel: per-service cache hit for %s@%s\n", svcPackageName, normalizedSvcVer)
		return basePath, nil
	}

	// Create the per-service cache directory
	if err := os.MkdirAll(modelDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("cannot create per-service model cache directory %s: %v", modelDir, err)
	}

	// Fetch from GitHub (per-service tag)
	perSvcURL := fmt.Sprintf(perServiceModelURLTemplate, svcPackageName, normalizedSvcVer, modelName)
	util.Tracef("EnsureModel: fetching %s\n", perSvcURL)

	fetchStart := time.Now()
	status, body, err := httpGetWithRetry(ctx, perSvcURL)
	if err != nil {
		return "", fmt.Errorf("cannot fetch model file from %s: %v", perSvcURL, err)
	}

	if status == http.StatusOK {
		if err := writeModelToCache(modelDir, modelName, body); err != nil {
			return "", err
		}
		util.Tracef("EnsureModel: fetched %s from per-service tag (%s)\n", modelName, time.Since(fetchStart))
		util.Tracef("EnsureModel total: %s\n", time.Since(totalStart))
		return basePath, nil
	}

	// Non-200 — return error with URL and status code, no core fallback
	return "", fmt.Errorf("failed to fetch model file from %s: HTTP %d", perSvcURL, status)
}

// ensureModelCore handles the core-only fetch path.
// It checks the core cache first, then fetches from the core tag URL.
// On 404 it suggests --aws-service-sdk-version; on other non-200 it returns
// an error with the URL and status code.
func ensureModelCore(
	ctx context.Context,
	cacheDir string,
	sdkVersion string,
	modelName string,
	totalStart time.Time,
) (string, error) {
	basePath := filepath.Join(cacheDir, "models", sdkVersion)
	modelDir := filepath.Join(basePath, "codegen", "sdk-codegen", "aws-models")
	modelPath := filepath.Join(modelDir, fmt.Sprintf("%s.json", modelName))

	// Check core cache
	if _, err := os.Stat(modelPath); err == nil {
		util.Tracef("EnsureModel: cache hit for %s@%s\n", modelName, sdkVersion)
		return basePath, nil
	}

	// Create the core cache directory
	if err := os.MkdirAll(modelDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("cannot create model cache directory %s: %v", modelDir, err)
	}

	// Fetch from GitHub (core tag)
	coreURL := fmt.Sprintf(sdkModelURLTemplate, sdkVersion, modelName)
	util.Tracef("EnsureModel: fetching %s\n", coreURL)

	fetchStart := time.Now()
	status, body, err := httpGetWithRetry(ctx, coreURL)
	if err != nil {
		return "", fmt.Errorf("cannot fetch model file from %s: %v", coreURL, err)
	}

	if status == http.StatusOK {
		if err := writeModelToCache(modelDir, modelName, body); err != nil {
			return "", err
		}
		util.Tracef("EnsureModel: fetched %s (%s)\n", modelName, time.Since(fetchStart))
		util.Tracef("EnsureModel total: %s\n", time.Since(totalStart))
		return basePath, nil
	}

	if status == http.StatusNotFound {
		return "", fmt.Errorf(
			"model %s not found at core tag %s; to use a per-service tag, provide --aws-service-sdk-version",
			modelName, sdkVersion,
		)
	}

	// Non-404 error from core URL
	return "", fmt.Errorf("failed to fetch model file from %s: HTTP %d", coreURL, status)
}

// httpGetWithRetry performs an HTTP GET request like httpGet, retrying with
// backoff while the server answers 429 or 5xx. raw.githubusercontent.com
// throttles unauthenticated clients, and jobs that regenerate many controllers
// in a row otherwise fail whichever model happened to be throttled.
//
// It returns the status and body of the last attempt, so the caller's handling
// of a non-200 status is unchanged once the attempts run out. Transport errors
// and other statuses are returned straight away.
func httpGetWithRetry(ctx context.Context, url string) (int, []byte, error) {
	delay := fetchRetryBaseDelay
	for attempt := 1; ; attempt++ {
		status, body, header, err := httpGet(ctx, url)
		if err != nil || !isRetryableStatus(status) || attempt == maxFetchAttempts {
			return status, body, err
		}

		wait := delay
		if retryAfter, ok := parseRetryAfter(header.Get("Retry-After"), time.Now()); ok {
			wait = retryAfter
		}
		if wait > fetchRetryMaxDelay {
			wait = fetchRetryMaxDelay
		}
		util.Warnf("EnsureModel: %s returned HTTP %d, retrying in %s (attempt %d of %d)\n",
			url, status, wait, attempt+1, maxFetchAttempts)

		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return status, body, ctx.Err()
		case <-timer.C:
		}
		delay *= 2
	}
}

// isRetryableStatus reports whether an HTTP status is worth retrying: rate
// limiting, or a server-side error.
func isRetryableStatus(status int) bool {
	return status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
}

// parseRetryAfter parses a Retry-After header value, which is either a number
// of seconds or an HTTP date. It returns false when the value is missing,
// malformed, or not in the future.
func parseRetryAfter(value string, now time.Time) (time.Duration, bool) {
	if value == "" {
		return 0, false
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		if seconds <= 0 {
			return 0, false
		}
		return time.Duration(seconds) * time.Second, true
	}
	if at, err := http.ParseTime(value); err == nil {
		if wait := at.Sub(now); wait > 0 {
			return wait, true
		}
	}
	return 0, false
}

// httpGet performs an HTTP GET request and returns the status code, response
// body, response headers, and any error. The caller is responsible for the body
// bytes.
func httpGet(ctx context.Context, url string) (int, []byte, http.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultHTTPFetchTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("cannot create HTTP request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, resp.Header, fmt.Errorf("cannot read response body: %v", err)
	}

	return resp.StatusCode, body, resp.Header, nil
}

// writeModelToCache writes model data to the cache directory atomically.
func writeModelToCache(modelDir string, modelName string, data []byte) error {
	modelPath := filepath.Join(modelDir, fmt.Sprintf("%s.json", modelName))

	tmpFile, err := os.CreateTemp(modelDir, ".tmp-"+modelName+"-*.json")
	if err != nil {
		return fmt.Errorf("cannot create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()

	_, err = tmpFile.Write(data)
	closeErr := tmpFile.Close()
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("cannot write model file: %v", err)
	}
	if closeErr != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("cannot close temp file: %v", closeErr)
	}

	if err := os.Rename(tmpPath, modelPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("cannot rename temp file to cache path: %v", err)
	}

	return nil
}
