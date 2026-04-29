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

// FieldGroupOperationConfig defines an API operation that manages a subset of
// a resource's fields. Used in both update_operations and read_operations to
// specify per-field-group API calls.
//
// When the generator encounters these configs, it auto-detects which fields
// belong to each operation by inspecting the operation's Input/Output shape
// and cross-referencing with the resource's identifier fields (those shared
// with ReadOne/Delete operations). Explicitly listing fields in the Fields
// slice overrides this auto-detection.
//
// Example generator.yaml:
//
//	resources:
//	  Repository:
//	    update_operations:
//	      - operation_id: PutImageScanningConfiguration
//	        requeue_on_success: true
//	      - operation_id: PutImageTagMutability
//	    read_operations:
//	      - operation_id: GetLifecyclePolicy
//	      - operation_id: GetRepositoryPolicy
type FieldGroupOperationConfig struct {
	// OperationID is the SDK operation's exported name
	// (e.g., "PutImageScanningConfiguration").
	OperationID string `json:"operation_id"`
	// Fields optionally overrides auto-detection of payload fields. When
	// empty, payload fields are auto-detected from the operation's Input
	// shape by excluding identifier fields. When set, only the listed CRD
	// field names are treated as payload fields for this group.
	Fields []string `json:"fields,omitempty"`
	// RequeueOnSuccess, when true, causes the reconciler to requeue after
	// this operation succeeds. This is useful for operations whose response
	// does not contain the updated field values, requiring a subsequent
	// ReadOne to refresh state.
	RequeueOnSuccess bool `json:"requeue_on_success,omitempty"`
}
