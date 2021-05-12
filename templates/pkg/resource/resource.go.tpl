{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
{{- if $idField := .CRD.SpecIdentifierField }}
	ackerrors "github.com/aws-controllers-k8s/runtime/pkg/errors"
{{- end }}
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

// SetObjectMeta sets the ObjectMeta field for the resource
func (r *resource) SetObjectMeta(meta metav1.ObjectMeta) {
	r.ko.ObjectMeta = meta;
}

// SetIdentifiers sets the Spec or Status field that is referenced as the unique
// resource identifier and any additional spec fields that may be required for
// describing the resource.
func (r *resource) SetIdentifiers(identifier *ackv1alpha1.AWSIdentifiers) error {
{{- if $idField := .CRD.SpecIdentifierField }}
	if identifier.NameOrID == nil {
		return ackerrors.MissingNameIdentifier
	}
	r.ko.Spec.{{ $idField }} = identifier.NameOrID
{{- else }}
	if r.ko.Status.ACKResourceMetadata == nil {
		r.ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	r.ko.Status.ACKResourceMetadata.ARN = identifier.ARN
{{- end }}

	if len(identifier.AdditionalKeys) == 0 {
		return nil
	}

	specRef := reflect.Indirect(reflect.ValueOf(&r.ko.Spec))
	specType := specRef.Type()

	// Iterate over spec fields and associate corresponding json tags
	for i := 0; i < specRef.NumField(); i++ {
		// Get only the first field in the json tag (the field name)
		jsonTag := strings.Split(specType.Field(i).Tag.Get("json"), ",")[0]
		val, ok := identifier.AdditionalKeys[jsonTag]
		if ok {
			// Set the corresponding field with the value from the mapping
			specRef.FieldByName(specType.Field(i).Name).Set(reflect.ValueOf(val))
		}
	}

	return nil
}