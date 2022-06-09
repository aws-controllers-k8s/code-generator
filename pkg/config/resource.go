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
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

// ResourceConfig represents instructions to the ACK code generator
// for a particular CRD/resource on an AWS service API
type ResourceConfig struct {
	// UnpackAttributeMapConfig contains instructions for converting a raw
	// `map[string]*string` into real fields on a CRD's Spec or Status object
	UnpackAttributesMapConfig *UnpackAttributesMapConfig `json:"unpack_attributes_map,omitempty"`
	// Exceptions identifies the exception codes for the resource. Some API
	// model files don't contain the ErrorInfo struct that contains the
	// HTTPStatusCode attribute that we usually look for to identify 404 Not
	// Found and other common error types for primary resources, and thus we
	// need these instructions.
	Exceptions *ExceptionsConfig `json:"exceptions,omitempty"`
	// Hooks is a map, keyed by the hook identifier, of instructions for the
	// the code generator about a custom callback hooks that should be injected
	// into the resource's manager or SDK binding code.
	Hooks map[string]*HooksConfig `json:"hooks"`
	// Synced contains instructions for the code generator to generate Go code
	// that verifies whether a resource is synced or not.
	Synced *SyncedConfig `json:"synced"`
	// Renames identifies fields in Operations that should be renamed.
	Renames *RenamesConfig `json:"renames,omitempty"`
	// ListOperation contains instructions for the code generator to generate
	// Go code that filters the results of a List operation looking for a
	// singular object. Certain AWS services (e.g. S3's ListBuckets API) have
	// absolutely no way to pass a filter to the operation. Instead, the List
	// operation always returns ALL objects of that type.
	//
	// The ListOperationConfig object enables us to inject some custom code to
	// filter the results of these List operations from within the generated
	// code in sdk.go's sdkFind().
	ListOperation *ListOperationConfig `json:"list_operation,omitempty"`
	// UpdateOperation contains instructions for the code generator to generate
	// Go code for the update operation for the resource. For some APIs, the
	// way that a resource's attributes are updated after creation is, well,
	// very odd. Some APIs have separate API calls for each attribute or set of
	// related attributes of the resource. For example, the ECR API has
	// separate API calls for PutImageScanningConfiguration,
	// PutImageTagMutability, PutLifecyclePolicy and SetRepositoryPolicy. FOr
	// these APIs, we basically need to revert to custom code because there's
	// very little consistency to the APIs that we can use to instruct the code
	// generator :(
	UpdateOperation *UpdateOperationConfig `json:"update_operation,omitempty"`
	// Reconcile describes options for controlling the reconciliation
	// logic for a particular resource.
	Reconcile *ReconcileConfig `json:"reconcile,omitempty"`
	// UpdateConditionsCustomMethodName provides the name of the custom method on the
	// `resourceManager` struct that will set Conditions on a `resource` struct
	// depending on the status of the resource.
	UpdateConditionsCustomMethodName string `json:"update_conditions_custom_method_name,omitempty"`
	// Fields is a map, keyed by the field name, of instructions for how the
	// code generator should interpret and handle a particular field in the
	// resource.
	Fields map[string]*FieldConfig `json:"fields"`
	// Compare contains instructions for the code generation to generate custom
	// comparison logic.
	Compare *CompareConfig `json:"compare,omitempty"`
	// ShortNames represent the CRD list of aliases. Short names allow shorter strings to
	// match a CR on the CLI.
	// All ShortNames must be distinct from any other ShortNames installed into the cluster,
	// otherwise the CRD will fail to install.
	ShortNames []string `json:"shortNames,omitempty"`
	// APIVersions represents the API versions defined for the generated CRD.
	// Default version to be used is the one specified via the "--version"
	// command-line argument, if none is specified here.
	// For the Crossplane pipeline, the generation target is the version
	// marked as the storage version.
	APIVersions []APIVersion `json:"api_versions,omitempty"`
	// IsAdoptable determines whether the CRD should be accepted by the adoption reconciler.
	// If set to false, the user will be given an error if they attempt to adopt a resource
	// with this type.
	IsAdoptable *bool `json:"is_adoptable,omitempty"`
	// Print contains instructions for the code generator to generate kubebuilder printcolumns
	// marker comments.
	Print *PrintConfig `json:"print,omitempty"`
	// IsARNPrimaryKey determines whether the CRD uses the ARN as the primary
	// identifier in the ReadOne operations.
	IsARNPrimaryKey bool `json:"is_arn_primary_key"`
	// TagConfig contains instructions for the code generator to generate
	// custom code for ensuring tags
	TagConfig *TagConfig `json:"tags,omitempty"`
}

// TagConfig instructs the code  generator on how to generate functions that
// ensure that controller tags are added to the AWS Resource
type TagConfig struct {
	// Ignore is a boolean that indicates whether ensuring controller tags
	// should be ignored for a resource. For AWS resources which do not
	// support tagging, this should be set to True
	Ignore bool `json:"ignore,omitempty"`
	// Path represents the field path for the member which contains the tags
	Path *string `json:"path,omitempty"`
	// KeyMemberName is the name of field which represents AWS tag key inside tag
	// struct. This is only used for tag fields with shape as list of struct,
	// where the struct represents a single tag.
	KeyMemberName *string `json:"key_name,omitempty"`
	// ValueMemberName is the name of field which represents AWS tag value inside
	// tag struct. This is only used for tag fields with shape as list of struct,
	// where the struct represents a single tag.
	ValueMemberName *string `json:"value_name,omitempty"`
}

// SyncedConfig instructs the code generator on how to generate functions that checks
// whether a resource was synced or not.
type SyncedConfig struct {
	// When is a list of conditions that should be satisfied in order to tell whether a
	// a resource was synced or not.
	When []SyncedCondition `json:"when"`
}

// SyncedCondition represent one of the unique condition that should be fullfiled in
// order to assert whether a resource is synced.
type SyncedCondition struct {
	// Path of the field. e.g Status.Processing
	Path *string `json:"path"`
	// In contains a list of possible values `Path` should be equal to.
	In []string `json:"in"`
}

// HooksConfig instructs the code generator how to inject custom callback hooks
// at various places in the resource manager and SDK linkage code.
//
// Example usage from the AmazonMQ generator config:
//
// resources:
//   Broker:
//     hooks:
//       sdk_update_pre_build_request:
//        code: if err := rm.requeueIfNotRunning(latest); err != nil { return nil, err }
//
// Note that the implementor of the AmazonMQ service controller for ACK should
// ensure that there is a `requeueIfNotRunning()` method implementation in
// `pkg/resource/broker`
//
// Instead of placing Go code directly into the generator.yaml file using the
// `code` field, you can reference a template file containing Go code with the
// `template_path` field:
//
// resources:
//   Broker:
//     hooks:
//       sdk_update_pre_build_update_request:
//        template_path: templates/sdk_update_pre_build_request.go.tpl
type HooksConfig struct {
	// Code is the Go code to be injected at the hook point
	Code *string `json:"code,omitempty"`
	// TemplatePath is a path to the template containing the hook code
	TemplatePath *string `json:"template_path,omitempty"`
}

// CompareConfig informs instruct the code generator on how to compare two different
// two objects of the same type
type CompareConfig struct {
	// Ignore is a list of field paths to ignore when comparing two objects
	Ignore []string `json:"ignore"`
}

// UnpackAttributesMapConfig informs the code generator that the API follows a
// pattern or using an "Attributes" `map[string]*string` that contains real,
// schema'd fields of the primary resource, and that those fields should be
// "unpacked" from the raw map and into CRD's Spec and Status struct fields.
//
// AWS Simple Notification Service (SNS) and AWS Simple Queue Service (SQS) are
// examples of APIs that use this pattern. For instance, the SNS CreateTopic
// API accepts a parameter called "Attributes" that can contain one of four
// keys:
//
// * DeliveryPolicy – The policy that defines how Amazon SNS retries failed
//   deliveries to HTTP/S endpoints.
// * DisplayName – The display name to use for a topic with SMS subscriptions
// * Policy – The policy that defines who can access your topic.
// * KmsMasterKeyId - The ID of an AWS-managed customer master key (CMK) for
//   Amazon SNS or a custom CMK.
//
// The `CreateTopic` API call **returns** only a single field: the TopicARN.
// But there is a separate `GetTopicAttributes` call that needs to be made that
// returns the above attributes (that are ReadWrite) along with a set of
// key/values that are ReadOnly:
//
// * Owner – The AWS account ID of the topic's owner.
// * SubscriptionsConfirmed – The number of confirmed subscriptions for the
//   topic.
// * SubscriptionsDeleted – The number of deleted subscriptions for the topic.
// * SubscriptionsPending – The number of subscriptions pending confirmation
//   for the topic.
// * TopicArn – The topic's ARN.
// * EffectiveDeliveryPolicy – The JSON serialization of the effective delivery
//   policy, taking system defaults into account.
//
// This structure instructs the code generator about the above real, schema'd
// fields that are masquerading as raw key/value pairs.
type UnpackAttributesMapConfig struct {
	// SetAttributesSingleAttribute indicates that the SetAttributes API call
	// doesn't actually set multiple attributes but rather must be called
	// multiple times, once for each attribute that needs to change. See SNS
	// SetTopicAttributes API call, which can be compared to the "normal" SNS
	// SetPlatformApplicationAttributes API call which accepts multiple
	// attributes and replaces the supplied attributes map key/values...
	SetAttributesSingleAttribute bool `json:"set_attributes_single_attribute"`
	// GetAttributesInput instructs the code generator how to handle the
	// GetAttributes input shape
	GetAttributesInput *GetAttributesInputConfig `json:"get_attributes_input,omitempty"`
}

// GetAttributesInputConfig is used to instruct the code generator how to
// handle the GetAttributes API operation's Input shape.
type GetAttributesInputConfig struct {
	// Overrides is a map of structures instructing the code generator how to
	// handle the override of a particular field in the Input shape for the
	// GetAttributes operation. The map keys are the names of the field in the
	// Input shape to override.
	Overrides map[string]*MemberConstructorConfig `json:"overrides"`
}

// MemberConstructorConfig contains override instructions for how to handle the
// construction of a particular member for a Shape in the API.
type MemberConstructorConfig struct {
	// Values contains the value or values of the member to always set the
	// member to. If the member's type is a []string, the member is set to the
	// Values list. If the type is a string, the member's value is set to the
	// first list element in the Values list.
	Values []string `json:"values"`
}

// ExceptionsConfig contains instructions to the code generator about how to
// handle the exceptions for the operations on a resource. These instructions
// are necessary for those APIs where the API models do not contain any
// information about the HTTP status codes a particular exception has (or, like
// the EC2 API, where the API model has no information at all about error
// responses for any operation)
type ExceptionsConfig struct {
	// Errors is a map of HTTP status code to information about the Exception
	// that corresponds to that HTTP status code for this resource
	Errors map[int]ErrorConfig `json:"errors"`
	// Set of aws exception codes that are terminal exceptions for this resource
	TerminalCodes []string `json:"terminal_codes"`
}

// ErrorConfig contains instructions to the code generator about the exception
// corresponding to a HTTP status code
type ErrorConfig struct {
	// Code corresponds to name of Exception returned by AWS API.
	// In AWS Go SDK terms - awsErr.Code()
	Code string `json:"code"`
	// MessagePrefix is an optional string field to be checked as prefix of the
	// exception message in addition to exception name. This is needed for HTTP codes
	// where the exception name alone is not sufficient to determine the type of error.
	// Example: SageMaker service throws ValidationException if job does not exist
	// as well as if IAM role does not have sufficient permission to fetch the dataset
	// For the former controller should proceed with creation of job whereas the
	// later is a terminal state.
	// In Go SDK terms - awsErr.Message()
	MessagePrefix *string `json:"message_prefix,omitempty"`
	// MessageSuffix is an optional string field to be checked as suffix of the
	// exception message in addition to exception name. This is needed for HTTP codes
	// where the exception name alone is not sufficient to determine the type of error.
	// Example: SageMaker service throws ValidationException if job does not exist
	// as well as if IAM role does not have sufficient permission to fetch the dataset
	// For the former controller should proceed with creation of job whereas the
	// later is a terminal state.
	// In Go SDK terms - awsErr.Message()
	MessageSuffix *string `json:"message_suffix,omitempty"`
}

// RenamesConfig contains instructions to the code generator how to rename
// fields in various Operation payloads
type RenamesConfig struct {
	// Operations is a map, keyed by Operation ID, of instructions on how to
	// handle renamed fields in Input and Output shapes.
	Operations map[string]*OperationRenamesConfig `json:"operations"`
}

// OperationRenamesConfig contains instructions to the code generator on how to
// rename fields in an Operation's input and output payload shapes
type OperationRenamesConfig struct {
	// InputFields is a map of Input shape fields to renamed field name.
	InputFields map[string]string `json:"input_fields"`
	// OutputFields is a map of Output shape fields to renamed field name.
	OutputFields map[string]string `json:"output_fields"`
}

// ListOperationConfig contains instructions for the code generator to handle
// List operations for service APIs that have no built-in filtering ability and
// whose List Operation always returns all objects.
type ListOperationConfig struct {
	// MatchFields lists the names of fields in the Shape of the
	// list element in the List Operation's Output shape.
	MatchFields []string `json:"match_fields"`
}

// UpdateOperationConfig contains instructions for the code generator to handle
// Update operations for service APIs that have resources that have
// difficult-to-standardize update operations.
type UpdateOperationConfig struct {
	// CustomMethodName is a string for the method name to replace the
	// sdkUpdate() method implementation for this resource
	CustomMethodName string `json:"custom_method_name"`
}

// PrintConfig informs instruct the code generator on how to sort kubebuilder
// printcolumn marker coments.
type PrintConfig struct {
	// AddAgeColumn a boolean informing the code generator whether to append a kubebuilder
	// marker comment to show a resource Age (created since date) in `kubectl get` response.
	// The Age value is parsed from '.metadata.creationTimestamp'.
	//
	// NOTE: this is the Kubernetes resource Age (creation time at the api-server/etcd)
	// and not the AWS resource Age.
	AddAgeColumn bool `json:"add_age_column"`
	// AddSyncedColumn is used to append a kubebuilder marker comment to show the status of a
	// resource in `kubectl get` response.
	//
	// Default value is true.
	AddSyncedColumn *bool `json:"add_synced_column"`
	// OrderBy is the field used to sort the list of PrinterColumn options.
	OrderBy string `json:"order_by"`
}

// ReconcileConfig describes options for controlling the reconciliation
// logic for a particular resource.
type ReconcileConfig struct {
	// RequeueOnSuccessSeconds indicates the number of seconds after which to requeue a
	// resource that has been successfully reconciled (i.e. ConditionTypeResourceSynced=true)
	// This is useful for resources that are long-lived and may have observable status fields
	// change over time that would be useful to refresh those field values for users.
	// This field is optional and the default behaviour of the ACK runtime is to not requeue
	// resources that have been successfully reconciled. Note that all ACK controllers will
	// *flush and resync their watch caches* every 10 hours by default, which will end up
	// causing ACK controllers to refresh the status views of all watched resources, but this
	// behaviour is expensive and may be turned off in future ACK runtime options.
	RequeueOnSuccessSeconds int `json:"requeue_on_success_seconds,omitempty"`
}

// ResourceIsIgnored returns true if resource name is configured to be ignored
// in generator config for the AWS service
func (c *Config) ResourceIsIgnored(resourceName string) bool {
	if resourceName == "" {
		return true
	}
	if c == nil {
		return false
	}
	return util.InStrings(resourceName, c.Ignore.ResourceNames)
}

// ResourceIsAdoptable returns true if resource name is configured to be adoptable.
// Default behavior is every resource can be adopted using AdoptionReconciler
func (c *Config) ResourceIsAdoptable(resourceName string) bool {
	if c == nil {
		return true
	}
	rConfig, ok := c.Resources[resourceName]
	if !ok {
		return true
	}
	if rConfig.IsAdoptable == nil {
		return true
	}
	return *rConfig.IsAdoptable
}

// ResourceContainsAttributesMap returns true if the underlying API has
// Get{Resource}Attributes/Set{Resource}Attributes API calls that map real,
// schema'd fields to a raw `map[string]*string` for a given resource name (see SNS and
// SQS APIs)
func (c *Config) ResourceContainsAttributesMap(resourceName string) bool {
	if c == nil {
		return false
	}
	resGenConfig, found := c.Resources[resourceName]
	if found {
		if resGenConfig.UnpackAttributesMapConfig != nil {
			return true
		}
		for _, fConfig := range resGenConfig.Fields {
			if fConfig.IsAttribute {
				return true
			}
		}
	}
	return false
}

// ResourceDisplaysAgeColumn returns true if the resource is
// configured to display resource age (created since date)
func (c *Config) ResourceDisplaysAgeColumn(resourceName string) bool {
	if c == nil {
		return false
	}
	rConfig, ok := c.Resources[resourceName]
	if !ok {
		return false
	}
	if rConfig.Print != nil {
		return rConfig.Print.AddAgeColumn
	}
	return false
}

// ResourceDisplaysSyncedColumn returns true if the resource is
// configured to display the synced status.
func (c *Config) ResourceDisplaysSyncedColumn(resourceName string) bool {
	if c == nil {
		return false
	}
	rConfig, ok := c.Resources[resourceName]
	if !ok {
		return false
	}
	if rConfig.Print != nil {
		// default value should be true.
		return rConfig.Print.AddSyncedColumn == nil || *rConfig.Print.AddSyncedColumn
	}
	return false
}

// ResourceSetsSingleAttribute returns true if the supplied resource name has
// a SetAttributes operation that only actually changes a single attribute at a
// time. See: SNS SetTopicAttributes API call, which is entirely different from
// the SNS SetPlatformApplicationAttributes API call, which sets multiple
// attributes at once. :shrug:
func (c *Config) ResourceSetsSingleAttribute(resourceName string) bool {
	if c == nil {
		return false
	}
	resGenConfig, found := c.Resources[resourceName]
	if !found || resGenConfig.UnpackAttributesMapConfig == nil {
		return false
	}
	return resGenConfig.UnpackAttributesMapConfig.SetAttributesSingleAttribute
}

// GetResourceConfig returns the ResourceConfig for a given resource name
func (c *Config) GetResourceConfig(resourceName string) *ResourceConfig {
	if c == nil {
		return nil
	}
	rc, _ := c.Resources[resourceName]
	return &rc
}

// GetCompareIgnoredFieldPaths returns the list of field paths to ignore when
// comparing two different objects
func (c *Config) GetCompareIgnoredFieldPaths(resourceName string) []string {
	if c == nil {
		return nil
	}
	rConfig, ok := c.Resources[resourceName]
	if !ok {
		return nil
	}
	if rConfig.Compare == nil {
		return nil
	}
	return rConfig.Compare.Ignore
}

// GetResourceFieldName returns a resource field name
// after applying rename overrides, if configured
func (c *Config) GetResourceFieldName(
	resourceName string,
	opID string,
	origFieldName string,
) string {
	if c == nil {
		return origFieldName
	}
	rConfig, ok := c.Resources[resourceName]
	if !ok {
		return origFieldName
	}
	if rConfig.Renames == nil {
		return origFieldName
	}
	oRenames, ok := rConfig.Renames.Operations[opID]
	if !ok {
		return origFieldName
	}
	renamed, ok := oRenames.InputFields[origFieldName]
	if !ok {
		renamed, ok = oRenames.OutputFields[origFieldName]
		if !ok {
			return origFieldName
		}
	}
	return renamed
}

// GetResourceShortNames returns the CRD list of aliases
func (c *Config) GetResourceShortNames(resourceName string) []string {
	if c == nil {
		return nil
	}
	rConfig, ok := c.Resources[resourceName]
	if !ok {
		return nil
	}
	return rConfig.ShortNames
}

// GetResourcePrintOrderByName returns the Printer Column order-by field name
func (c *Config) GetResourcePrintOrderByName(resourceName string) string {
	if c == nil {
		return ""
	}
	rConfig, ok := c.Resources[resourceName]
	if !ok {
		return ""
	}
	if rConfig.Print != nil {
		return rConfig.Print.OrderBy
	}
	return ""
}

// GetUpdateConditionsCustomMethodName returns custom update conditions operation
// as *string for custom resource, if specified in generator config
func (c *Config) GetUpdateConditionsCustomMethodName(resourceName string) string {
	if c == nil {
		return ""
	}
	resGenConfig, found := c.Resources[resourceName]
	if !found {
		return ""
	}
	return resGenConfig.UpdateConditionsCustomMethodName
}

// GetReconcileRequeueOnSuccessSeconds returns the duration after which to requeue
// the custom resource as int, if specified in generator config.
func (c *Config) GetReconcileRequeueOnSuccessSeconds(resourceName string) int {
	if c == nil {
		return 0
	}
	resGenConfig, found := c.Resources[resourceName]
	if !found {
		return 0
	}
	reconcile := resGenConfig.Reconcile
	if reconcile != nil {
		return reconcile.RequeueOnSuccessSeconds
	}
	// handles the default case
	return 0
}

// GetCustomUpdateMethodName returns the name of the custom resourceManager method
// for updating the resource state, if any has been specified in the generator
// config
func (c *Config) GetCustomUpdateMethodName(resourceName string) string {
	if c == nil {
		return ""
	}
	rConfig, found := c.Resources[resourceName]
	if found {
		if rConfig.UpdateOperation != nil {
			return rConfig.UpdateOperation.CustomMethodName
		}
	}
	return ""
}

// GetAllRenames returns all of the CRD's field renames observed in the generator config
// for a given map of operations.
func (c *Config) GetAllRenames(
	resourceName string,
	operations map[string]*awssdkmodel.Operation,
) map[string]string {
	renames := make(map[string]string)
	if c == nil {
		return renames
	}
	resourceConfig, ok := c.Resources[resourceName]
	if !ok {
		return renames
	}
	if resourceConfig.Renames == nil || resourceConfig.Renames.Operations == nil {
		return renames
	}

	opRenameConfigs := resourceConfig.Renames.Operations
	for opName, opRenameConfigs := range opRenameConfigs {
		for _, op := range operations {
			if opName != op.Name {
				continue
			}
			for old, new := range opRenameConfigs.InputFields {
				renames[old] = new
			}
			for old, new := range opRenameConfigs.OutputFields {
				renames[old] = new
			}
		}
	}
	return renames
}

// GetTerminalExceptionCodes returns terminal exception codes as
// []string for custom resource, if specified in generator config
func (c *Config) GetTerminalExceptionCodes(resourceName string) []string {
	if c == nil {
		return nil
	}
	resGenConfig, found := c.Resources[resourceName]
	if found && resGenConfig.Exceptions != nil {
		return resGenConfig.Exceptions.TerminalCodes
	}
	return nil
}

// GetListOpMatchFieldNames returns a slice of strings representing the field
// names in the List operation's Output shape's element Shape that we should
// check a corresponding value in the target Spec exists.
func (c *Config) GetListOpMatchFieldNames(
	resName string,
) []string {
	res := []string{}
	if c == nil {
		return res
	}
	rConfig, found := c.Resources[resName]
	if !found {
		return res
	}
	if rConfig.ListOperation == nil {
		return res
	}
	return rConfig.ListOperation.MatchFields
}
