{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
	k8sctrlutil "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	svcapitypes "github.com/aws-controllers-k8s/{{ .ServicePackageName }}-controller/apis/{{ .APIVersion }}"
)

const (
	finalizerString = "finalizers.{{ .APIGroup }}/{{ .CRD.Kind }}"
)

var (
	GroupVersionResource = svcapitypes.GroupVersion.WithResource("{{ ToLower .CRD.Plural }}")
	GroupKind = metav1.GroupKind{
		Group: "{{ .APIGroup }}",
		Kind:  "{{ .CRD.Kind }}",
	}
)

// resourceDescriptor implements the
// `aws-service-operator-k8s/pkg/types.AWSResourceDescriptor` interface
type resourceDescriptor struct {
}

// GroupVersionKind returns a Kubernetes schema.GroupVersionKind struct that
// describes the API Group, Version and Kind of CRs described by the descriptor
func (d *resourceDescriptor) GroupVersionKind() schema.GroupVersionKind {
	return svcapitypes.GroupVersion.WithKind(GroupKind.Kind)
}

// EmptyRuntimeObject returns an empty object prototype that may be used in
// apimachinery and k8s client operations
func (d *resourceDescriptor) EmptyRuntimeObject() rtclient.Object {
	return &svcapitypes.{{ .CRD.Kind }}{}
}

// ResourceFromRuntimeObject returns an AWSResource that has been initialized
// with the supplied runtime.Object
func (d *resourceDescriptor) ResourceFromRuntimeObject(
	obj rtclient.Object,
) acktypes.AWSResource {
	return &resource{
		ko: obj.(*svcapitypes.{{ .CRD.Kind }}),
	}
}

// Delta returns an `ackcompare.Delta` object containing the difference between
// one `AWSResource` and another.
func (d *resourceDescriptor) Delta(a, b acktypes.AWSResource) *ackcompare.Delta {
    return newResourceDelta(a.(*resource), b.(*resource))
}

// IsManaged returns true if the supplied AWSResource is under the management
// of an ACK service controller. What this means in practice is that the
// underlying custom resource (CR) in the AWSResource has had a
// resource-specific finalizer associated with it.
func (d *resourceDescriptor) IsManaged(
	res acktypes.AWSResource,
) bool {
	obj := res.RuntimeObject()
	if obj == nil {
		// Should not happen. If it does, there is a bug in the code
		panic("nil RuntimeMetaObject in AWSResource")
	}
	// Remove use of custom code once
	// https://github.com/kubernetes-sigs/controller-runtime/issues/994 is
	// fixed. This should be able to be:
	//
	// return k8sctrlutil.ContainsFinalizer(obj, finalizerString)
	return containsFinalizer(obj, finalizerString)
}

// Remove once https://github.com/kubernetes-sigs/controller-runtime/issues/994
// is fixed.
func containsFinalizer(obj rtclient.Object, finalizer string) bool {
	f := obj.GetFinalizers()
	for _, e := range f {
		if e == finalizer {
			return true
		}
	}
	return false
}

// MarkManaged places the supplied resource under the management of ACK.  What
// this typically means is that the resource manager will decorate the
// underlying custom resource (CR) with a finalizer that indicates ACK is
// managing the resource and the underlying CR may not be deleted until ACK is
// finished cleaning up any backend AWS service resources associated with the
// CR.
func (d *resourceDescriptor) MarkManaged(
	res acktypes.AWSResource,
) {
	obj := res.RuntimeObject()
	if obj == nil {
		// Should not happen. If it does, there is a bug in the code
		panic("nil RuntimeMetaObject in AWSResource")
	}
	k8sctrlutil.AddFinalizer(obj, finalizerString)
}

// MarkUnmanaged removes the supplied resource from management by ACK.  What
// this typically means is that the resource manager will remove a finalizer
// underlying custom resource (CR) that indicates ACK is managing the resource.
// This will allow the Kubernetes API server to delete the underlying CR.
func (d *resourceDescriptor) MarkUnmanaged(
	res acktypes.AWSResource,
) {
	obj := res.RuntimeObject()
	if obj == nil {
		// Should not happen. If it does, there is a bug in the code
		panic("nil RuntimeMetaObject in AWSResource")
	}
	k8sctrlutil.RemoveFinalizer(obj, finalizerString)
}

// MarkAdopted places descriptors on the custom resource that indicate the
// resource was not created from within ACK.
func (d *resourceDescriptor) MarkAdopted(
	res acktypes.AWSResource,
) {
	obj := res.RuntimeObject()
	if obj == nil {
		// Should not happen. If it does, there is a bug in the code
		panic("nil RuntimeObject in AWSResource")
	}
	curr := obj.GetAnnotations()
	if curr == nil {
		curr = make(map[string]string)
	}
	curr[ackv1alpha1.AnnotationAdopted] = "true"
	obj.SetAnnotations(curr)
}
