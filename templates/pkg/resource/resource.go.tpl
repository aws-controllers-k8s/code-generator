{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	ackerrors "github.com/aws-controllers-k8s/runtime/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"

	svcapitypes "github.com/aws-controllers-k8s/{{ .ControllerName }}-controller/apis/{{ .APIVersion}}"
)

// Hack to avoid import errors during build...
var (
	_ = &ackerrors.MissingNameIdentifier
)

// resource implements the `aws-controller-k8s/runtime/pkg/types.AWSResource`
// interface
type resource struct {
	// The Kubernetes-native CR representing the resource
	ko *svcapitypes.{{ .CRD.Names.Camel }}
}

// Identifiers returns an AWSResourceIdentifiers object containing various
// identifying information, including the AWS account ID that owns the
// resource, the resource's AWS Resource Name (ARN)
func (r *resource) Identifiers() acktypes.AWSResourceIdentifiers {
	return &resourceIdentifiers{r.ko.Status.ACKResourceMetadata}
}

// IsBeingDeleted returns true if the Kubernetes resource has a non-zero
// deletion timestamp
func (r *resource) IsBeingDeleted() bool {
	return !r.ko.DeletionTimestamp.IsZero()
}

// RuntimeObject returns the Kubernetes apimachinery/runtime representation of
// the AWSResource
func (r *resource) RuntimeObject() rtclient.Object {
	return r.ko
}

// MetaObject returns the Kubernetes apimachinery/apis/meta/v1.Object
// representation of the AWSResource
func (r *resource) MetaObject() metav1.Object {
	return r.ko.GetObjectMeta()
}

// Conditions returns the ACK Conditions collection for the AWSResource
func (r *resource) Conditions() []*ackv1alpha1.Condition {
	return r.ko.Status.Conditions
}

// ReplaceConditions sets the Conditions status field for the resource
func (r *resource) ReplaceConditions(conditions []*ackv1alpha1.Condition) {
	r.ko.Status.Conditions = conditions
}

// SetObjectMeta sets the ObjectMeta field for the resource
func (r *resource) SetObjectMeta(meta metav1.ObjectMeta) {
	r.ko.ObjectMeta = meta;
}

// SetStatus will set the Status field for the resource
func (r *resource) SetStatus(desired acktypes.AWSResource) {
	r.ko.Status = desired.(*resource).ko.Status
}

// SetIdentifiers sets the Spec or Status field that is referenced as the unique
// resource identifier
func (r *resource) SetIdentifiers(identifier *ackv1alpha1.AWSIdentifiers) error {
{{- if $hookCode := Hook .CRD "pre_set_resource_identifiers" }}
{{ $hookCode }}
{{- end }}
{{- GoCodeSetResourceIdentifiers .CRD "identifier" "r.ko" 1}}
{{- if $hookCode := Hook .CRD "post_set_resource_identifiers" }}
{{ $hookCode }}
{{- end }}
	return nil
}

// PopulateResourceFromAnnotation populates the fields passed from adoption annotation 
// 
func (r *resource) PopulateResourceFromAnnotation(fields map[string]string) error {
{{- if $hookCode := Hook .CRD "pre_populate_resource_from_annotation" }}
{{ $hookCode }}
{{- end }}
{{- GoCodePopulateResourceFromAnnotation .CRD "fields" "r.ko" 1}}
{{- if $hookCode := Hook .CRD "post_populate_resource_from_annotation" }}
{{ $hookCode }}
{{- end }}
	return nil
}


// DeepCopy will return a copy of the resource
func (r *resource) DeepCopy() acktypes.AWSResource {
	koCopy := r.ko.DeepCopy()
	return &resource{koCopy}
}
