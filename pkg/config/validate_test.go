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

package config

import (
	"strings"
	"testing"
)

func TestValidateRenameOperations(t *testing.T) {
	sdkOps := map[string]struct{}{
		"CreateBucket": {},
		"DeleteBucket": {},
		"ListBuckets":  {},
	}

	tests := []struct {
		name            string
		resources       map[string]ResourceConfig
		wantErrCount    int
		wantErrContains string
	}{
		{
			name: "valid rename operation",
			resources: map[string]ResourceConfig{
				"Bucket": {
					Renames: &RenamesConfig{
						Operations: map[string]*OperationRenamesConfig{
							"CreateBucket": {InputFields: map[string]string{"foo": "bar"}},
						},
					},
				},
			},
			wantErrCount: 0,
		},
		{
			name: "invalid rename operation",
			resources: map[string]ResourceConfig{
				"Bucket": {
					Renames: &RenamesConfig{
						Operations: map[string]*OperationRenamesConfig{
							"CreateBuckett": {InputFields: map[string]string{"foo": "bar"}},
						},
					},
				},
			},
			wantErrCount:    1,
			wantErrContains: "CreateBuckett",
		},
		{
			name: "no renames config",
			resources: map[string]ResourceConfig{
				"Bucket": {},
			},
			wantErrCount: 0,
		},
		{
			name: "multiple invalid renames across resources",
			resources: map[string]ResourceConfig{
				"Bucket": {
					Renames: &RenamesConfig{
						Operations: map[string]*OperationRenamesConfig{
							"CreateBuckett": {},
						},
					},
				},
				"Object": {
					Renames: &RenamesConfig{
						Operations: map[string]*OperationRenamesConfig{
							"DeleteObjekt": {},
						},
					},
				},
			},
			wantErrCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Resources: tt.resources}
			errs := validateRenameOperations(cfg, sdkOps)
			if len(errs) != tt.wantErrCount {
				t.Errorf("got %d errors, want %d: %v", len(errs), tt.wantErrCount, errs)
			}
			if tt.wantErrContains != "" && len(errs) > 0 {
				if !strings.Contains(errs[0].Error(), tt.wantErrContains) {
					t.Errorf("error %q does not contain %q", errs[0].Error(), tt.wantErrContains)
				}
			}
		})
	}
}

func TestValidateIgnoredOperations(t *testing.T) {
	sdkOps := map[string]struct{}{
		"CreateBucket": {},
		"DeleteBucket": {},
	}

	tests := []struct {
		name            string
		ignored         []string
		wantErrCount    int
		wantErrContains string
	}{
		{
			name:         "valid ignored operation",
			ignored:      []string{"CreateBucket"},
			wantErrCount: 0,
		},
		{
			name:            "invalid ignored operation",
			ignored:         []string{"CreateBuckett"},
			wantErrCount:    1,
			wantErrContains: "CreateBuckett",
		},
		{
			name:         "empty ignored",
			ignored:      nil,
			wantErrCount: 0,
		},
		{
			name:         "multiple invalid ignored operations",
			ignored:      []string{"BadOp1", "BadOp2"},
			wantErrCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Ignore: IgnoreSpec{Operations: tt.ignored}}
			errs := validateIgnoredOperations(cfg, sdkOps)
			if len(errs) != tt.wantErrCount {
				t.Errorf("got %d errors, want %d: %v", len(errs), tt.wantErrCount, errs)
			}
			if tt.wantErrContains != "" && len(errs) > 0 {
				if !strings.Contains(errs[0].Error(), tt.wantErrContains) {
					t.Errorf("error %q does not contain %q", errs[0].Error(), tt.wantErrContains)
				}
			}
		})
	}
}

func TestValidateConfig_AllErrors(t *testing.T) {
	sdkOps := map[string]struct{}{
		"CreateBucket": {},
		"DeleteBucket": {},
	}

	cfg := &Config{
		Resources: map[string]ResourceConfig{
			"Bucket": {
				Renames: &RenamesConfig{
					Operations: map[string]*OperationRenamesConfig{
						"CreateBuckett": {},
					},
				},
			},
		},
		Ignore: IgnoreSpec{
			Operations: []string{"BadOp"},
		},
	}

	errs := ValidateConfig(cfg, sdkOps)
	// 1: CreateBuckett rename, 2: BadOp ignored op
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d: %v", len(errs), errs)
	}
}

func TestValidateConfig_NilConfig(t *testing.T) {
	errs := ValidateConfig(nil, nil)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors for nil config, got %d", len(errs))
	}
}

func TestValidateConfig_ValidConfig(t *testing.T) {
	sdkOps := map[string]struct{}{
		"CreateBucket": {},
		"DeleteBucket": {},
	}

	cfg := &Config{
		Resources: map[string]ResourceConfig{
			"Bucket": {
				Renames: &RenamesConfig{
					Operations: map[string]*OperationRenamesConfig{
						"CreateBucket": {InputFields: map[string]string{"foo": "bar"}},
					},
				},
			},
		},
		Ignore: IgnoreSpec{
			Operations: []string{"DeleteBucket"},
		},
	}

	errs := ValidateConfig(cfg, sdkOps)
	if len(errs) != 0 {
		t.Errorf("expected 0 errors for valid config, got %d: %v", len(errs), errs)
	}
}

func TestValidateConfig_ErrorMessageIncludesAvailable(t *testing.T) {
	sdkOps := map[string]struct{}{
		"CreateBucket": {},
		"DeleteBucket": {},
	}

	cfg := &Config{
		Ignore: IgnoreSpec{
			Operations: []string{"Typo"},
		},
	}

	errs := ValidateConfig(cfg, sdkOps)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	msg := errs[0].Error()
	if !strings.Contains(msg, "CreateBucket") || !strings.Contains(msg, "DeleteBucket") {
		t.Errorf("error should list available operations, got: %s", msg)
	}
}

func TestFormatAvailableTruncated(t *testing.T) {
	items := []string{"A", "B", "C", "D", "E"}
	got := formatAvailableTruncated(items, 3)
	if !strings.Contains(got, "5 total") {
		t.Errorf("expected truncation message, got: %s", got)
	}

	got = formatAvailableTruncated(items, 10)
	if strings.Contains(got, "total") {
		t.Errorf("should not truncate when maxItems > len: %s", got)
	}
}
