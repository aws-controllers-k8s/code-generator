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
	errs = append(errs, validateFieldGroupOperations(cfg, sdkOperations)...)

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

// validateFieldGroupOperations checks that field group operation configs
// (update_operations and read_operations) reference valid SDK operations and
// are not misconfigured.
func validateFieldGroupOperations(
	cfg *Config,
	sdkOperations map[string]struct{},
) []error {
	var errs []error
	for resName, resCfg := range cfg.Resources {
		// Check update_operations
		seen := make(map[string]bool)
		for i, fgCfg := range resCfg.UpdateOperations {
			if fgCfg.OperationID == "" {
				errs = append(errs, fmt.Errorf(
					"resources.%s.update_operations[%d]: operation_id must not be empty",
					resName, i,
				))
				continue
			}
			if _, ok := sdkOperations[fgCfg.OperationID]; !ok {
				errs = append(errs, fmt.Errorf(
					"resources.%s.update_operations[%d]: operation %q not found in SDK. available: %s",
					resName, i, fgCfg.OperationID, formatAvailableTruncated(sortedKeys(sdkOperations), 10),
				))
			}
			if seen[fgCfg.OperationID] {
				errs = append(errs, fmt.Errorf(
					"resources.%s.update_operations[%d]: duplicate operation_id %q",
					resName, i, fgCfg.OperationID,
				))
			}
			seen[fgCfg.OperationID] = true
		}

		// Mutual exclusivity: custom_method_name + update_operations
		if len(resCfg.UpdateOperations) > 0 && resCfg.UpdateOperation != nil && resCfg.UpdateOperation.CustomMethodName != "" {
			errs = append(errs, fmt.Errorf(
				"resources.%s: update_operations cannot be used together with update_operation.custom_method_name",
				resName,
			))
		}

		// Check read_operations
		seen = make(map[string]bool)
		for i, fgCfg := range resCfg.ReadOperations {
			if fgCfg.OperationID == "" {
				errs = append(errs, fmt.Errorf(
					"resources.%s.read_operations[%d]: operation_id must not be empty",
					resName, i,
				))
				continue
			}
			if _, ok := sdkOperations[fgCfg.OperationID]; !ok {
				errs = append(errs, fmt.Errorf(
					"resources.%s.read_operations[%d]: operation %q not found in SDK. available: %s",
					resName, i, fgCfg.OperationID, formatAvailableTruncated(sortedKeys(sdkOperations), 10),
				))
			}
			if seen[fgCfg.OperationID] {
				errs = append(errs, fmt.Errorf(
					"resources.%s.read_operations[%d]: duplicate operation_id %q",
					resName, i, fgCfg.OperationID,
				))
			}
			seen[fgCfg.OperationID] = true
		}
	}
	return errs
}

// formatAvailableTruncated formats a sorted slice, showing at most maxItems
// entries with an ellipsis if truncated.
func formatAvailableTruncated(items []string, maxItems int) string {
	if len(items) <= maxItems {
		return strings.Join(items, ", ")
	}
	return strings.Join(items[:maxItems], ", ") + fmt.Sprintf(", ... (%d total)", len(items))
}
