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

// SourceFieldConfig instructs the code generator how to handle a field in the
// Resource's SpecFields/StatusFields collection that takes its value from an
// abnormal source -- in other words, not the Create operation's Input or
// Output shape.
//
// This additional field can source its value from a shape in a different API
// Operation entirely.
//
// The data type (Go type) that a field is assigned during code generation
// depends on whether the field is part of the Create Operation's Input shape
// which go into the Resource's Spec fields collection, or the Create
// Operation's Output shape which, if not present in the Input shape, means the
// field goes into the Resource's Status fields collection).
//
// Each Resource typically also has a ReadOne Operation. The ACK service
// controller will call this ReadOne Operation to get the latest observed state
// of a particular resource in the backend AWS API service. The service
// controller sets the observed Resource's Spec and Status fields from the
// Output shape of the ReadOne Operation. The code generator is responsible for
// producing the Go code that performs these "setter" methods on the Resource.
// The way the code generator determines how to set the Spec or Status fields
// from the Output shape's member fields is by looking at the data type of the
// Spec or Status field with the same name as the Output shape's member field.
//
// Importantly, in producing this "setter" Go code the code generator **assumes
// that the data types (Go types) in the source (the Output shape's member
// field) and target (the Spec or Status field) are the same**.
//
// There are some APIs, however, where the Go type of the field in the Create
// Operation's Input shape is actually different from the same-named field in
// the ReadOne Operation's Output shape. A good example of this is the Lambda
// CreateFunction API call, which has a `Code` member of its Input shape that
// looks like this:
//
// "Code": {
//   "ImageUri": "string",
//   "S3Bucket": "string",
//   "S3Key": "string",
//   "S3ObjectVersion": "string",
//   "ZipFile": blob
// },
//
// The GetFunction API call's Output shape has a same-named field called
// `Code` in it, but this field looks like this:
//
// "Code": {
//   "ImageUri": "string",
//   "Location": "string",
//   "RepositoryType": "string",
//   "ResolvedImageUri": "string"
// },
//
// This presents a conundrum to the ACK code generator, which, as noted above,
// assumes the data types of same-named fields in the Create Operation's Input
// shape and ReadOne Operation's Output shape are the same.
//
// The SourceFieldConfig struct allows us to explain to the code generator
// how to handle situations like this.
//
// For the Lambda Function Resource's `Code` field, we can inform the code
// generator to create three new Status fields (readonly) from the `Location`,
// `RepositoryType` and `ResolvedImageUri` fields in the `Code` member of the
// ReadOne Operation's Output shape:
//
// resources:
//   Function:
//     fields:
//       CodeLocation:
//         is_read_only: true
//         from:
//           operation: GetFunction
//           path: Code.Location
//       CodeRepositoryType:
//         is_read_only: true
//         from:
//           operation: GetFunction
//           path: Code.RepositoryType
//       CodeRegisteredImageURI:
//         is_read_only: true
//         from:
//           operation: GetFunction
//           path: Code.RegisteredImageUri
type SourceFieldConfig struct {
	// Operation refers to the ID of the API Operation where we will
	// determine the field's Go type.
	Operation string `json:"operation"`
	// Path refers to the field path of the member of the Input or Output
	// shape in the Operation identified by OperationID that we will take as
	// our additional spec/status field's value.
	Path string `json:"path"`
}

// SetFieldConfig instructs the code generator how to handle setting the value
// of the field from an Output shape in the API response.
//
// These instructions are necessary when the Go type for Input and Output
// fields is different.
//
// For example, consider the `DBSecurityGroups` field for an RDS `DBInstance`
// resource.
//
// The Create operation's Input shape's `DBSecurityGroups` field has a Go type
// of `[]*string` [0] because the user is expected to provide a list of
// DBSecurityGroup identifiers when creating the DBInstance.
//
// However, the Output shape for both the Create and ReadOne operation uses a
// different Go type for the `DBSecurityGroups` field. For these Output shapes,
// the Go type is `[]*DBSecurityGroupMembership` [1]. The
// `DBSecurityGroupMembership` struct contains the DBSecurityGroup's name and
// "status".
//
// The challenge that the ACK code generator has is to figure out how to take
// the Go type of the Output shape and process it into the Go type of the Input
// shape.
//
// In other words, for the `DBInstance.Spec.DBSecurityGroups` field discussed
// above, we need to have the code generator produce the following Go code in
// the SetResource generator:
//
// ```go
// if resp.DBInstance.DBSecurityGroups != nil {
//     f17 := []*string{}
//     for _, f17iter := range resp.DBInstance.DBSecurityGroups {
//         var f17elem string
//         f17elem = *f17iter.DBSecurityGroupName
//         f17 = append(f17, &f17elem)
//     }
//     ko.Spec.DBSecurityGroupNames = f17
// } else {
//     ko.Spec.DBSecurityGroupNames = nil
// }
// ```
//
// [0] https://github.com/aws/aws-sdk-go/blob/0a01aef9caf16d869c7340e729080205760dc2a2/models/apis/rds/2014-10-31/api-2.json#L2985
// [1] https://github.com/aws/aws-sdk-go/blob/0a01aef9caf16d869c7340e729080205760dc2a2/models/apis/rds/2014-10-31/api-2.json#L3815
type SetFieldConfig struct {
	// Method is the resource manager method name whose Output shape will be
	// transformed by this config. If empty, this set field config applies to
	// all resource manager methods.
	//
	// Options: Create, Update, Delete or ReadOne
	Method *string `json:"method,omitempty"`
	// From tells the code generator to output Go code that sets the value of a
	// variable containing the target resource field with the contents of a
	// member field in the source struct.
	//
	// Consider the `DBInstance.Spec.DBSecurityGroups` field discussed above.
	//
	// If we have the following generator.yaml config:
	//
	// ```yaml
	// resources:
	//   DBInstance:
	//     fields:
	//       DBSecurityGroups:
	//         set:
	//           - method: Create
	//			   from: DBSecurityGroupName
	//           - method: ReadOne
	//			   from: DBSecurityGroupName
	// ```
	//
	// That will instruct the code generator to output this Go code when
	// processing the `*DBSecurityGroupMembership` struct elements of the
	// Output shape's DBSecurityGroups field:
	//
	// ```go
	// f17elem = *f17iter.DBSecurityGroupName
	// ```
	From *string `json:"from,omitempty"`
	// Ignore instructs the code generator to ignore this field in the Output
	// shape when setting the value of the resource's field in the Spec. This
	// is useful when we know that, for example, the returned value of field in
	// an Output shape contains stale data, such as when the ElastiCache
	// ModifyReplicationGroup API's Output response shape contains the
	// originally-set value for the LogDeliveryConfiguration field that was
	// updated in the Input shape.
	//
	// See: https://github.com/aws-controllers-k8s/elasticache-controller/pull/59/
	//
	// In the case of ElastiCache, we might have the following generator.yaml
	// config:
	//
	// ```yaml
	// resources:
	//   ReplicationGroup:
	//     fields:
	//       LogDeliveryConfiguration:
	//         set:
	//           - method: Update
	//			   ignore: true
	// ```
	Ignore bool `json:"ignore,omitempty"`
}

// CompareFieldConfig informs the code generator how to compare two values of a
// field
type CompareFieldConfig struct {
	// IsIgnored indicates the field should be ignored when comparing a
	// resource
	IsIgnored bool `json:"is_ignored"`
	// NilEqualsZeroValue indicates a nil pointer and zero-value pointed-to
	// value should be considered equal for the purposes of comparison
	NilEqualsZeroValue bool `json:"nil_equals_zero_value"`
}

// PrintFieldConfig instructs the code generator how to handle kubebuilder:printcolumn
// comment marker generation. If this struct is not nil, the field will be added to the
// columns of `kubectl get` response.
type PrintFieldConfig struct {
	// Name instructs the code generator to override the column name used to
	// include the field in `kubectl get` response. This field is generally used
	// to override very long and redundant columns names.
	Name string `json:"name"`
	// Priority differentiates between fields/columns shown in standard view or wide
	// view (using the -o wide flag). Fields with priority 0 are shown in standard view.
	// Fields with priority greater than 0 are only shown in wide view. Default is 0
	Priority int `json:"priority"`
	// Index informs the code generator about the position/order of a specific field/column in
	// `kubectl get` response. To enable ordering by index, `$resource.print.orderBy` must be set
	// to `index`
	// The field with the smallest index will be right next to the first column (NAME).
	// The field with the biggest index will be positioned right before the last column (AGE).
	Index int `json:"index"`
}

// CustomField instructs the code generator to create a new list or map field
// type using a shape that exists in the SDK.
type CustomFieldConfig struct {
	// ListOf provides the name of the SDK shape which will become the
	// member of a custom slice field.
	ListOf string `json:"list_of,omitempty"`
	// MapOf provides the name of the SDK shape which will become the value
	// shape for a custom map field. All maps will have `string` as their key
	// type.
	MapOf string `json:"map_of,omitempty"`
}

// LateInitializeConfig contains instructions for how to handle the
// retrieval and setting of server-side defaulted fields.
// NOTE: Currently the members of this have no effect on late initialization of fields.
// Currently the late initialization is requeued with static delay of 5 second.
// TODO: (vijat@) Add support of retry/backoff for late initialization.
type LateInitializeConfig struct {
	// MinBackoffSeconds provides the minimum backoff to attempt late initialization again after an unsuccessful
	// attempt to late initialized fields from ReadOne output
	// For every attempt, the reconciler will calculate the delay between MinBackoffSeconds and MaxBackoffSeconds
	// using exponential backoff and retry strategy
	MinBackoffSeconds int `json:"min_backoff_seconds,omitempty"`
	// MaxBackoffSeconds provide the maximum allowed backoff when retrying late initialization after an
	// unsuccessful attempt.
	MaxBackoffSeconds int `json:"max_backoff_seconds"`
}

// ReferencesConfig contains the instructions for how to add the referenced resource
// configuration for a field.
// Example:
// ```
// Integration:
//    fields:
//      ApiId:
//        references:
//          resource: API
//          path: Status.APIID
//```
// The above configuration will result in generation of a new field 'APIRef'
// of type 'AWSResourceReference' for ApiGatewayv2-Integration crd.
// When 'APIRef' field is present in custom resource manifest, reconciler will
// read the referred 'API' resource and copy the value from 'Status.APIID' in
// 'Integration' resource's 'APIID' field
type ReferencesConfig struct {
	// ServiceName mentions the AWS service name where "Resource" exists.
	// This field is used to generate the API Group for the "Resource".
	//
	// ServiceName is the go package name for AWS service in
	// aws-sdk-go/service/<package_name>/api.go
	// Ex: Use "opensearchservice" to refer "Domain" resource from
	// opensearchservice-controller because it is the go package name for
	// "aws-sdk-go/service/opensearchservice/api.go"
	//
	// When not specified, 'ServiceName' defaults to service name of controller
	// which contains generator.yaml
	ServiceName string `json:"service_name,omitempty"`
	// Resource mentions the K8s resource which is read to resolve the
	// reference
	Resource string `json:"resource"`
	// Path refers to the the path of field which should be copied
	// to resolve the reference
	Path string `json:"path"`
}

// FieldConfig contains instructions to the code generator about how
// to interpret the value of an Attribute and how to map it to a CRD's Spec or
// Status field
type FieldConfig struct {
	// IsAttribute informs the code generator that this field is part of an
	// "Attributes Map".
	//
	// Some resources for some service APIs follow a pattern or using an
	// "Attributes" `map[string]*string` that contains real, schema'd fields of
	// the primary resource, and that those fields should be "unpacked" from
	// the raw map and into CRD's Spec and Status struct fields.
	IsAttribute bool `json:"is_attribute"`
	// IsReadOnly indicates the field's value can not be set by a Kubernetes
	// user; in other words, the field should go in the CR's Status struct
	IsReadOnly bool `json:"is_read_only"`
	// Required indicates whether this field is a required member or not.
	// This field is used to configure '+kubebuilder:validation:Required' on API object's members.
	IsRequired *bool `json:"is_required,omitempty"`
	// IsPrimaryKey indicates the field represents the primary name/string
	// identifier field for the resource.  This allows the generator config to
	// override the default behaviour of considering a field called "Name" or
	// "{Resource}Name" or "{Resource}Id" as the "name field" for the resource.
	IsPrimaryKey bool `json:"is_primary_key"`
	// IsOwnerAccountID indicates the field contains the AWS Account ID
	// that owns the resource. This is a special field that we direct to
	// storage in the common `Status.ACKResourceMetadata.OwnerAccountID` field.
	IsOwnerAccountID bool `json:"is_owner_account_id"`
	// IsARN indicates the field represents the ARN for the resource.
	// This allows the generator config to override the
	// default behaviour of considering a field called "Arn" or
	// "{Resource}Arn" (case in-sensitive) as the "ARN field" for the resource.
	IsARN bool `json:"is_arn"`
	// IsSecret instructs the code generator that this field should be a
	// SecretKeyReference.
	IsSecret bool `json:"is_secret"`
	// IsImmutable instructs the code generator to add advisory conditions
	// if user modifies the spec field after resource was created.
	IsImmutable bool `json:"is_immutable"`
	// From instructs the code generator that the value of the field should
	// be retrieved from the specified operation and member path
	From *SourceFieldConfig `json:"from,omitempty"`
	// CustomField instructs the code generator to create a new field that does
	// not exist in the SDK.
	CustomField *CustomFieldConfig `json:"custom_field,omitempty"`
	// Compare instructs the code generator how to produce code that compares
	// the value of the field in two resources
	Compare *CompareFieldConfig `json:"compare,omitempty"`
	// Set contains instructions for the code generator how to deal with
	// fields where the Go type of the same-named fields in an Output shape is
	// different from the Go type of the Input shape.
	Set []*SetFieldConfig `json:"set"`
	// Print instructs the code generator how to generate comment markers that
	// influence hows field are printed in `kubectl get` response. If this field
	// is not nil, it will be added to the columns of `kubectl get`.
	Print *PrintFieldConfig `json:"print,omitempty"`
	// Late Initialize instructs the code generator how to handle the late initialization
	// of the field.
	LateInitialize *LateInitializeConfig `json:"late_initialize,omitempty"`
	// References instructs the code generator how to refer this field from
	// other custom resource
	References *ReferencesConfig `json:"references,omitempty"`
	// Type *overrides* the inferred Go type of the field. This is required for
	// custom fields that are not inferred either as a Create Input/Output
	// shape or via the SourceFieldConfig attribute.
	//
	// As an example, assume you have a Role resource where you want to add a
	// custom spec field called Policies that is a slice of string pointers.
	// The generator.yaml file might look like this:
	//
	// resources:
	//   Role:
	//     fields:
	//       Policies:
	//         type: []*string
	//
	// TODO(jaypipes,crtbry): Figure out if we can roll the CustomShape stuff
	// into this type override...
	Type *string `json:"type,omitempty"`
}
