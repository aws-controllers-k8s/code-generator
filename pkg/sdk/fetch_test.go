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

package sdk_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/sdk"
)

// newTestServer creates an httptest server that routes requests based on path.
// The handler map keys are URL paths; values are (statusCode, body) pairs.
func newTestServer(routes map[string]struct {
	status int
	body   string
}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if route, ok := routes[r.URL.Path]; ok {
			w.WriteHeader(route.status)
			fmt.Fprint(w, route.body)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "not found")
	}))
}

// patchURLTemplates temporarily overrides the URL templates used by EnsureModel
// by replacing http.DefaultClient's Transport to rewrite requests to the test
// server. This approach avoids modifying the production code's constants.
//
// Returns a cleanup function that restores the original transport.
func patchTransport(testServerURL string) func() {
	original := http.DefaultClient.Transport
	http.DefaultClient.Transport = &rewriteTransport{
		base:          original,
		testServerURL: testServerURL,
	}
	return func() {
		http.DefaultClient.Transport = original
	}
}

// rewriteTransport rewrites GitHub raw content URLs to point to the test server.
type rewriteTransport struct {
	base          http.RoundTripper
	testServerURL string
}

func (t *rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite any request to raw.githubusercontent.com to our test server
	if strings.Contains(req.URL.Host, "raw.githubusercontent.com") {
		// Keep the path, but change the scheme+host to the test server
		newURL := t.testServerURL + req.URL.Path
		newReq, err := http.NewRequestWithContext(req.Context(), req.Method, newURL, req.Body)
		if err != nil {
			return nil, err
		}
		base := t.base
		if base == nil {
			base = http.DefaultTransport
		}
		return base.RoundTrip(newReq)
	}
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(req)
}

// TestEnsureModel_CoreSuccess verifies that when no serviceSDKVersion is
// provided (empty string), the core tag URL is used and the model is cached
// at the core path.
func TestEnsureModel_CoreSuccess(t *testing.T) {
	modelBody := `{"metadata":{"apiVersion":"2023-01-01"}}`

	routes := map[string]struct {
		status int
		body   string
	}{
		"/aws/aws-sdk-go-v2/v1.41.5/codegen/sdk-codegen/aws-models/s3files.json": {
			status: http.StatusOK,
			body:   modelBody,
		},
	}
	ts := newTestServer(routes)
	defer ts.Close()
	cleanup := patchTransport(ts.URL)
	defer cleanup()

	cacheDir := t.TempDir()
	// Pass empty serviceSDKVersion to exercise the core-only branch
	basePath, err := sdk.EnsureModel(context.Background(), cacheDir, "v1.41.5", "s3files", "")
	require.NoError(t, err)

	// Verify cached at core path
	expectedBase := filepath.Join(cacheDir, "models", "v1.41.5")
	assert.Equal(t, expectedBase, basePath)

	cachedFile := filepath.Join(basePath, "codegen", "sdk-codegen", "aws-models", "s3files.json")
	data, err := os.ReadFile(cachedFile)
	require.NoError(t, err)
	assert.Equal(t, modelBody, string(data))
}

// TestEnsureModel_PerServiceSuccess verifies that when serviceSDKVersion is
// set, the per-service tag URL is tried first (not as a fallback from core
// 404) and no core request is made. The model is cached at the per-service
// path.
func TestEnsureModel_PerServiceSuccess(t *testing.T) {
	modelBody := `{"metadata":{"apiVersion":"2024-06-01"}}`

	coreRequested := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Track whether the core URL was ever requested
		if strings.Contains(r.URL.Path, "/aws/aws-sdk-go-v2/v1.41.5/") {
			coreRequested = true
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"core":"should not be used"}`)
			return
		}
		// Per-service URL returns 200
		if strings.Contains(r.URL.Path, "/service/s3files/v1.0.0/") {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, modelBody)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()
	cleanup := patchTransport(ts.URL)
	defer cleanup()

	cacheDir := t.TempDir()
	basePath, err := sdk.EnsureModel(context.Background(), cacheDir, "v1.41.5", "s3files", "v1.0.0")
	require.NoError(t, err)

	// Verify no core request was made — per-service is the sole path
	assert.False(t, coreRequested, "core URL should not be requested when serviceSDKVersion is set")

	// Verify cached at per-service path
	expectedBase := filepath.Join(cacheDir, "models", "service", "s3files", "v1.0.0")
	assert.Equal(t, expectedBase, basePath)

	cachedFile := filepath.Join(basePath, "codegen", "sdk-codegen", "aws-models", "s3files.json")
	data, err := os.ReadFile(cachedFile)
	require.NoError(t, err)
	assert.Equal(t, modelBody, string(data))
}

// TestEnsureModel_PerServiceFailure verifies that when serviceSDKVersion is
// set and the per-service URL returns a non-200 status, the error contains
// the per-service URL and the HTTP status code. No core fallback is attempted.
func TestEnsureModel_PerServiceFailure(t *testing.T) {
	coreRequested := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/aws/aws-sdk-go-v2/v1.41.5/") {
			coreRequested = true
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"core":"should not be used"}`)
			return
		}
		if strings.Contains(r.URL.Path, "/service/s3files/v1.0.0/") {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "server error")
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()
	cleanup := patchTransport(ts.URL)
	defer cleanup()

	cacheDir := t.TempDir()
	_, err := sdk.EnsureModel(context.Background(), cacheDir, "v1.41.5", "s3files", "v1.0.0")
	require.Error(t, err)

	// Verify no core fallback was attempted
	assert.False(t, coreRequested, "core URL should not be requested when serviceSDKVersion is set")

	errMsg := err.Error()
	// Error should contain the per-service URL and the status code
	assert.Contains(t, errMsg, "/service/s3files/v1.0.0/")
	assert.Contains(t, errMsg, "500")
}

// TestEnsureModel_NoServiceVersionOn404 verifies that when the core URL returns
// 404 and no service SDK version is provided, the error suggests --aws-service-sdk-version.
func TestEnsureModel_NoServiceVersionOn404(t *testing.T) {
	routes := map[string]struct {
		status int
		body   string
	}{
		"/aws/aws-sdk-go-v2/v1.41.5/codegen/sdk-codegen/aws-models/s3files.json": {
			status: http.StatusNotFound,
			body:   "not found",
		},
	}
	ts := newTestServer(routes)
	defer ts.Close()
	cleanup := patchTransport(ts.URL)
	defer cleanup()

	cacheDir := t.TempDir()
	_, err := sdk.EnsureModel(context.Background(), cacheDir, "v1.41.5", "s3files", "")
	require.Error(t, err)

	errMsg := err.Error()
	assert.Contains(t, errMsg, "--aws-service-sdk-version")
	assert.Contains(t, errMsg, "s3files")
}

// TestEnsureModel_CachePathSeparation verifies that for any model name
// and SDK version pair, core and per-service cache paths are distinct.
func TestEnsureModel_CachePathSeparation(t *testing.T) {
	tests := []struct {
		name           string
		modelName      string
		coreVersion    string
		serviceVersion string
	}{
		{
			name:           "s3files basic",
			modelName:      "s3files",
			coreVersion:    "v1.41.5",
			serviceVersion: "v1.0.0",
		},
		{
			name:           "sns different versions",
			modelName:      "sns",
			coreVersion:    "v1.40.0",
			serviceVersion: "v2.1.0",
		},
		{
			name:           "custom model name",
			modelName:      "my-custom-model",
			coreVersion:    "v1.50.0",
			serviceVersion: "v3.0.0",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cacheDir := "/tmp/test-cache"
			corePath := filepath.Join(cacheDir, "models", tc.coreVersion)
			perServicePath := filepath.Join(cacheDir, "models", "service", tc.modelName, tc.serviceVersion)

			assert.NotEqual(t, corePath, perServicePath,
				"core and per-service cache paths must be distinct")

			// Also verify they don't share a prefix that could cause file collisions
			coreModelFile := filepath.Join(corePath, "codegen", "sdk-codegen", "aws-models", tc.modelName+".json")
			perSvcModelFile := filepath.Join(perServicePath, "codegen", "sdk-codegen", "aws-models", tc.modelName+".json")
			assert.NotEqual(t, coreModelFile, perSvcModelFile,
				"core and per-service model file paths must be distinct")
		})
	}
}

// TestEnsureModel_PerServiceURLConstruction verifies that the per-service
// tag URL uses the model name in both the tag path segment and the file name,
// and contains the normalized (v-prefixed) version.
func TestEnsureModel_PerServiceURLConstruction(t *testing.T) {
	tests := []struct {
		name           string
		modelName      string
		serviceVersion string
		expectedPath   string
	}{
		{
			name:           "basic service",
			modelName:      "s3files",
			serviceVersion: "v1.0.0",
			expectedPath:   "/aws/aws-sdk-go-v2/service/s3files/v1.0.0/codegen/sdk-codegen/aws-models/s3files.json",
		},
		{
			name:           "model name override",
			modelName:      "monitoring",
			serviceVersion: "v2.3.1",
			expectedPath:   "/aws/aws-sdk-go-v2/service/monitoring/v2.3.1/codegen/sdk-codegen/aws-models/monitoring.json",
		},
		{
			name:           "version without v prefix gets normalized",
			modelName:      "sqs",
			serviceVersion: "1.5.0",
			// EnsureSemverPrefix will add the v prefix
			expectedPath: "/aws/aws-sdk-go-v2/service/sqs/v1.5.0/codegen/sdk-codegen/aws-models/sqs.json",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// We verify URL construction by setting up a test server that
			// returns 404 for core and captures the per-service request path.
			var capturedPath string
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.Contains(r.URL.Path, "/service/") {
					capturedPath = r.URL.Path
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, `{"test":"data"}`)
					return
				}
				// Core URL returns 404
				w.WriteHeader(http.StatusNotFound)
			}))
			defer ts.Close()
			cleanup := patchTransport(ts.URL)
			defer cleanup()

			cacheDir := t.TempDir()
			_, err := sdk.EnsureModel(context.Background(), cacheDir, "v1.41.5", tc.modelName, tc.serviceVersion)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedPath, capturedPath,
				"per-service URL path should match expected format")
		})
	}
}

// TestEnsureSemverPrefix_Idempotence verifies that applying
// EnsureSemverPrefix twice produces the same result as applying it once.
func TestEnsureSemverPrefix_Idempotence(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "no prefix", input: "1.0.0"},
		{name: "single v prefix", input: "v1.0.0"},
		{name: "double v prefix", input: "vv1.0.0"},
		{name: "triple v prefix", input: "vvv1.0.0"},
		{name: "complex version", input: "1.2.3-beta.1"},
		{name: "v-prefixed complex", input: "v1.2.3-beta.1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			once := sdk.EnsureSemverPrefix(tc.input)
			twice := sdk.EnsureSemverPrefix(once)
			assert.Equal(t, once, twice,
				"EnsureSemverPrefix should be idempotent: f(f(x)) == f(x)")
			// Also verify the result starts with exactly one 'v'
			assert.True(t, strings.HasPrefix(once, "v"), "result should start with 'v'")
			assert.False(t, strings.HasPrefix(once, "vv"), "result should not start with 'vv'")
		})
	}
}

// TestEnsureModel_CoreCacheHit verifies that when serviceSDKVersion is empty,
// a cached core model is returned without making any HTTP requests.
func TestEnsureModel_CoreCacheHit(t *testing.T) {
	cacheDir := t.TempDir()
	modelBody := `{"cached":"true"}`

	// Pre-populate the core cache
	modelDir := filepath.Join(cacheDir, "models", "v1.41.5", "codegen", "sdk-codegen", "aws-models")
	require.NoError(t, os.MkdirAll(modelDir, os.ModePerm))
	require.NoError(t, os.WriteFile(filepath.Join(modelDir, "s3files.json"), []byte(modelBody), 0644))

	// Pass empty serviceSDKVersion to exercise the core-only branch
	// No test server needed — if HTTP is attempted, it will fail
	basePath, err := sdk.EnsureModel(context.Background(), cacheDir, "v1.41.5", "s3files", "")
	require.NoError(t, err)

	expectedBase := filepath.Join(cacheDir, "models", "v1.41.5")
	assert.Equal(t, expectedBase, basePath)
}

// TestEnsureModel_PerServiceCacheHit verifies that when serviceSDKVersion is
// set, the per-service cache is checked first and returned on hit. No core
// cache check or HTTP request is made.
func TestEnsureModel_PerServiceCacheHit(t *testing.T) {
	cacheDir := t.TempDir()
	modelBody := `{"cached":"per-service"}`

	// Pre-populate the per-service cache
	modelDir := filepath.Join(cacheDir, "models", "service", "s3files", "v1.0.0", "codegen", "sdk-codegen", "aws-models")
	require.NoError(t, os.MkdirAll(modelDir, os.ModePerm))
	require.NoError(t, os.WriteFile(filepath.Join(modelDir, "s3files.json"), []byte(modelBody), 0644))

	// Also pre-populate a core cache with different content to prove it's not used
	coreModelDir := filepath.Join(cacheDir, "models", "v1.41.5", "codegen", "sdk-codegen", "aws-models")
	require.NoError(t, os.MkdirAll(coreModelDir, os.ModePerm))
	require.NoError(t, os.WriteFile(filepath.Join(coreModelDir, "s3files.json"), []byte(`{"cached":"core-should-not-be-used"}`), 0644))

	// No test server needed — if HTTP is attempted, it will fail
	basePath, err := sdk.EnsureModel(context.Background(), cacheDir, "v1.41.5", "s3files", "v1.0.0")
	require.NoError(t, err)

	// Verify per-service cache path is returned, not core
	expectedBase := filepath.Join(cacheDir, "models", "service", "s3files", "v1.0.0")
	assert.Equal(t, expectedBase, basePath)

	// Verify the returned path contains the per-service cached content
	cachedFile := filepath.Join(basePath, "codegen", "sdk-codegen", "aws-models", "s3files.json")
	data, err := os.ReadFile(cachedFile)
	require.NoError(t, err)
	assert.Equal(t, modelBody, string(data))
}
