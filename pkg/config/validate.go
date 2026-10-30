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
	"fmt"
	"sort"
	"strings"

	"github.com/aws-controllers-k8s/pkg/names"
)

// ValidateConfig checks that generator.yaml references to SDK operations are
// valid. It returns all validation errors found, not just the first.
//
// Validates:
//   - Operation names in resources[R].renames.operations (must exist in SDK)
//   - Operation names in ignore.operations (must exist in SDK)
//
// Does NOT validate:
//   - Resource names: controllers define resources with custom names
//   - Shape names (ignore/empty): shapes may not exist in every SDK version
//     and ApplyShapeIgnoreRules already tolerates missing shapes
//
// Parameter sdkOperations is a set of API.Operations keys (ExportedNames).
func ValidateConfig(
	cfg *Config,
	sdkOperations map[string]struct{},
) []error {
	if cfg == nil {
		return nil
	}

	var errs []error

	errs = append(errs, validateRenameOperations(cfg, sdkOperations)...)
	errs = append(errs, validateIgnoredOperations(cfg, sdkOperations)...)
	errs = append(errs, validateMutuallyExclusiveIdentifiers(cfg)...)
	errs = append(errs, validateSpecValidations(cfg)...)

	return errs
}

// validateMutuallyExclusiveIdentifiers checks that a resource's
// mutually_exclusive_identifiers configuration is internally consistent: it
// must list at least two fields (a single identifier is not mutually exclusive
// with anything), the fields must be distinct after name normalization (two
// entries that resolve to the same CR field can never be "exactly one"), and it
// cannot be combined with is_arn_primary_key, since an ARN-primary resource
// always requires its ARN.
func validateMutuallyExclusiveIdentifiers(cfg *Config) []error {
	var errs []error
	for resName, resCfg := range cfg.Resources {
		identifiers := resCfg.MutuallyExclusiveIdentifiers
		if len(identifiers) == 0 {
			continue
		}
		if len(identifiers) < 2 {
			errs = append(errs, fmt.Errorf(
				"resources.%s.mutually_exclusive_identifiers: must list at least two fields, got %d",
				resName, len(identifiers),
			))
		}
		// Entries are resolved to CR fields by camel-casing the name (see
		// CRD.GetMutuallyExclusiveIdentifierFields), so e.g. "PolicyName" and
		// "policyName" collapse to the same field. Reject duplicates after
		// normalization: the generated exactly-one guard would otherwise count
		// the same field twice and adoption could never satisfy it.
		seen := make(map[string]bool, len(identifiers))
		for _, id := range identifiers {
			norm := names.New(id).Camel
			if seen[norm] {
				errs = append(errs, fmt.Errorf(
					"resources.%s.mutually_exclusive_identifiers: %q resolves to the same field as another entry",
					resName, id,
				))
			}
			seen[norm] = true
		}
		if resCfg.IsARNPrimaryKey {
			errs = append(errs, fmt.Errorf(
				"resources.%s.mutually_exclusive_identifiers: cannot be combined with is_arn_primary_key",
				resName,
			))
		}
	}
	return errs
}

func validateSpecValidations(cfg *Config) []error {
	var errs []error
	for resourceName, resourceConfig := range cfg.Resources {
		for index, validation := range resourceConfig.SpecValidations {
			if strings.TrimSpace(validation.Rule) == "" {
				errs = append(errs, fmt.Errorf(
					"resources.%s.spec_validations[%d].rule: must not be empty",
					resourceName, index,
				))
			}
			if strings.TrimSpace(validation.Message) == "" {
				errs = append(errs, fmt.Errorf(
					"resources.%s.spec_validations[%d].message: must not be empty",
					resourceName, index,
				))
			}
		}
	}
	return errs
}

// validateRenameOperations checks that operation names referenced in
// resources[R].renames.operations[OpName] exist in the SDK.
func validateRenameOperations(
	cfg *Config,
	sdkOperations map[string]struct{},
) []error {
	var errs []error
	for resName, resCfg := range cfg.Resources {
		if resCfg.Renames == nil || resCfg.Renames.Operations == nil {
			continue
		}
		for opName := range resCfg.Renames.Operations {
			if _, ok := sdkOperations[opName]; !ok {
				errs = append(errs, fmt.Errorf(
					"resources.%s.renames.operations.%s: operation not found in SDK. available: %s",
					resName, opName, formatAvailableTruncated(sortedKeys(sdkOperations), 10),
				))
			}
		}
	}
	return errs
}

// validateIgnoredOperations checks that operation names in
// ignore.operations exist in the SDK.
func validateIgnoredOperations(
	cfg *Config,
	sdkOperations map[string]struct{},
) []error {
	var errs []error
	for _, opName := range cfg.Ignore.Operations {
		if _, ok := sdkOperations[opName]; !ok {
			errs = append(errs, fmt.Errorf(
				"ignore.operations: operation %q not found in SDK. available: %s",
				opName, formatAvailableTruncated(sortedKeys(sdkOperations), 10),
			))
		}
	}
	return errs
}

// sortedKeys returns sorted keys from a map[string]struct{}.
func sortedKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// formatAvailableTruncated formats a sorted slice, showing at most maxItems
// entries with an ellipsis if truncated.
func formatAvailableTruncated(items []string, maxItems int) string {
	if len(items) <= maxItems {
		return strings.Join(items, ", ")
	}
	return strings.Join(items[:maxItems], ", ") + fmt.Sprintf(", ... (%d total)", len(items))
}
