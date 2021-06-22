{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	ackerrors "github.com/aws-controllers-k8s/runtime/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8srt "k8s.io/apimachinery/pkg/runtime"

	svcapitypes "github.com/aws-controllers-k8s/{{ .ServiceIDClean }}-controller/apis/{{ .APIVersion}}"
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
// deletion timestemp
func (r *resource) IsBeingDeleted() bool {
	return !r.ko.DeletionTimestamp.IsZero()
}

// RuntimeObject returns the Kubernetes apimachinery/runtime representation of
// the AWSResource
func (r *resource) RuntimeObject() k8srt.Object {
	return r.ko
}

// MetaObject returns the Kubernetes apimachinery/apis/meta/v1.Object
// representation of the AWSResource
func (r *resource) MetaObject() metav1.Object {
	return r.ko
}

// RuntimeMetaObject returns an object that implements both the Kubernetes
// apimachinery/runtime.Object and the Kubernetes
// apimachinery/apis/meta/v1.Object interfaces
func (r *resource) RuntimeMetaObject() acktypes.RuntimeMetaObject {
	return r.ko
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

// SetIdentifiers sets the Spec or Status field that is referenced as the unique
// resource identifier
func (r *resource) SetIdentifiers(identifier *ackv1alpha1.AWSIdentifiers) error {
{{- if $idField := .CRD.SpecIdentifierField }}
	if identifier.NameOrID == nil {
		return ackerrors.MissingNameIdentifier
	}
	r.ko.Spec.{{ $idField }} = identifier.NameOrID
{{- else }}
	r.ko.Status.ACKResourceMetadata.ARN = identifier.ARN
{{- end }}
{{ GoCodeSetResourceIdentifiers .CRD "r.ko" "res" 1}}
	return nil
}