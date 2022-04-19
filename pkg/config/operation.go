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
	"encoding/json"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

// StringArray is a type that can be represented in JSON as *either* a string
// *or* an array of strings
type StringArray []string

// OperationConfig represents instructions to the ACK code generator to
// specify the overriding values for API operation parameters and its custom implementation.
type OperationConfig struct {
	CustomImplementation                   string            `json:"custom_implementation,omitempty"`
	CustomCheckRequiredFieldsMissingMethod string            `json:"custom_check_required_fields_missing_method,omitempty"`
	OverrideValues                         map[string]string `json:"override_values"`
	// SetOutputCustomMethodName provides the name of the custom method on the
	// `resourceManager` struct that will set fields on a `resource` struct
	// depending on the output of the operation.
	SetOutputCustomMethodName string `json:"set_output_custom_method_name,omitempty"`
	// OutputWrapperFieldPath provides the JSON-Path like to the struct field containing
	// information that will be merged into a `resource` object.
	OutputWrapperFieldPath string `json:"output_wrapper_field_path,omitempty"`
	// Override for resource name in case of heuristic failure
	// An example of this is correcting stutter when the resource logic doesn't properly determine the resource name
	ResourceName string `json:"resource_name"`
	// Override for operation type in case of heuristic failure
	// An example of this is `Put...` or `Register...` API operations not being correctly classified as `Create` op type
	// OperationType []string `json:"operation_type"`
	OperationType StringArray `json:"operation_type"`
}

// OperationIsIgnored returns true if Operation Name is configured to be ignored
// in generator config for the AWS service
func (c *Config) OperationIsIgnored(operation *awssdkmodel.Operation) bool {
	if c == nil {
		return false
	}
	if operation == nil {
		return true
	}
	return util.InStrings(operation.Name, c.Ignore.Operations)
}

// UnmarshalJSON parses input for a either a string or
// or a list and returns a StringArray.
func (a *StringArray) UnmarshalJSON(b []byte) error {
	var multi []string
	err := json.Unmarshal(b, &multi)
	if err != nil {
		var single string
		err := json.Unmarshal(b, &single)
		if err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

// GetOutputWrapperFieldPath returns the JSON-Path of the output wrapper field
// as *string for a given operation, if specified in generator config.
func (c *Config) GetOutputWrapperFieldPath(
	op *awssdkmodel.Operation,
) *string {
	if op == nil {
		return nil
	}
	if c == nil {
		return nil
	}
	opConfig, found := c.Operations[op.Name]
	if !found {
		return nil
	}

	if opConfig.OutputWrapperFieldPath == "" {
		return nil
	}
	return &opConfig.OutputWrapperFieldPath
}

// GetSetOutputCustomMethodName returns custom set output operation as *string for
// given operation on custom resource, if specified in generator config
func (c *Config) GetSetOutputCustomMethodName(
	// The operation to look for the Output shape
	op *awssdkmodel.Operation,
) *string {
	if op == nil {
		return nil
	}
	if c == nil {
		return nil
	}
	opConfig, found := c.Operations[op.Name]
	if !found {
		return nil
	}

	if opConfig.SetOutputCustomMethodName == "" {
		return nil
	}
	return &opConfig.SetOutputCustomMethodName
}

// GetCustomImplementation returns custom implementation method name for the
// supplied operation as specified in generator config
func (c *Config) GetCustomImplementation(
	// The type of operation
	op *awssdkmodel.Operation,
) string {
	if op == nil || c == nil {
		return ""
	}

	operationConfig, found := c.Operations[op.Name]
	if !found {
		return ""
	}

	return operationConfig.CustomImplementation
}

// GetCustomCheckRequiredFieldsMissingMethod returns custom check required fields missing method
// as string for custom resource, if specified in generator config
func (c *Config) GetCustomCheckRequiredFieldsMissingMethod(
	// The type of operation
	op *awssdkmodel.Operation,
) string {
	if op == nil || c == nil {
		return ""
	}

	operationConfig, found := c.Operations[op.Name]
	if !found {
		return ""
	}

	return operationConfig.CustomCheckRequiredFieldsMissingMethod
}

// OverrideValues returns a list of member values to override for a given operation
func (c *Config) GetOverrideValues(operationName string) (map[string]string, bool) {
	if c == nil {
		return nil, false
	}
	oConfig, ok := c.Operations[operationName]
	if !ok {
		return nil, false
	}
	return oConfig.OverrideValues, ok
}

// OperationConfig returns the OperationConfig for a given operation
func (c *Config) GetOperationConfig(opID string) (*OperationConfig, bool) {
	if c == nil {
		return nil, false
	}
	opConfig, ok := c.Operations[opID]
	return &opConfig, ok
}
