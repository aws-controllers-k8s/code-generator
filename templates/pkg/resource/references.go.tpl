{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"context"
{{ if .CRD.HasReferenceFields -}}
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
{{ end -}}
	"sigs.k8s.io/controller-runtime/pkg/client"

{{ if .CRD.HasReferenceFields -}}
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
{{ end -}}
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
{{ $servicePackageName := .ServicePackageName -}}
{{ $apiVersion := .APIVersion -}}
{{ if .CRD.HasReferenceFields -}}
{{ range $referencedServiceName := .CRD.ReferencedServiceNames -}}
{{ if not (eq $referencedServiceName $servicePackageName) -}}
    {{ $referencedServiceName }}apitypes "github.com/aws-controllers-k8s/{{ $referencedServiceName }}-controller/apis/{{ $apiVersion }}"
{{ end }}
{{- end }}
{{- end }}

	svcapitypes "github.com/aws-controllers-k8s/{{ .ServicePackageName }}-controller/apis/{{ .APIVersion }}"
)

{{ if .CRD.HasReferenceFields -}}
{{ range $fieldName, $field := .CRD.Fields -}}
{{ if and $field.HasReference (not (eq $field.ReferencedServiceName $servicePackageName)) -}}
// +kubebuilder:rbac:groups={{ $field.ReferencedServiceName -}}.services.k8s.aws,resources={{ ToLower $field.ReferencedResourceNamePlural }},verbs=get;list
// +kubebuilder:rbac:groups={{ $field.ReferencedServiceName -}}.services.k8s.aws,resources={{ ToLower $field.ReferencedResourceNamePlural }}/status,verbs=get;list

{{ end -}}
{{ end -}}
{{ end -}}

type resourceReferenceMap struct {
	// All of the references that have been resolved for the resource. The field
	// path of the reference field is the key and the resolved value is the
	// value.
	// TODO(redbackthomson): Sync MUTEX
	resolvedReferences map[string]any
}

// setReferencedValue sets the resolved reference value for the resource at a
// given field path
func (rm *resourceManager) setReferencedValue(ctx context.Context, path string, refValue any) {
	rlog := ackrtlog.FromContext(ctx)
	rlog.Info("setting resolved reference", "path", path, "referencedValue", refValue)
	rm.refMap.resolvedReferences[path] = refValue
}

// getReferencedValue will attempt to return the value of a resolved
// reference for a resource, or nil if it does not exist
func (rm *resourceManager) getReferencedValue(path string) (any, bool) {
	val, ok := rm.refMap.resolvedReferences[path]
	return val, ok
}

func (rm *resourceManager) CopyWithResolvedReferences(res acktypes.AWSResource) (acktypes.AWSResource, error) {
	ko := rm.concreteResource(res).ko.DeepCopy()

{{ range $fieldName, $field := .CRD.Fields -}}
{{ if $field.HasReference -}}
{{ GoCodeCopyWithResolvedReferences $field "ko" 1 }}
{{ end -}}
{{ end -}}

	return &resource{ko}, nil
}

func (rm *resourceManager) ClearResolvedReferences(res acktypes.AWSResource) (acktypes.AWSResource, error) {
	ko := rm.concreteResource(res).ko.DeepCopy()

{{ range $fieldName, $field := .CRD.Fields }}
{{ if $field.HasReference }}
{{ GoCodeClearResolvedReferences $field "ko" 1 }}
{{ end -}}
{{ end -}}
	return &resource{ko}, nil
}

func NewReferenceMap() *resourceReferenceMap {
	return &resourceReferenceMap{
		resolvedReferences: make(map[string]any),
	}
}

// ResolveReferences finds if there are any Reference field(s) present
// inside AWSResource passed in the parameter and attempts to resolve
// those reference field(s) into the resolved reference cache within the
// resource manager. No fields are modified within the resource itself.
func (rm *resourceManager) ResolveReferences(
	ctx context.Context,
	apiReader client.Reader,
	res acktypes.AWSResource,
) (bool, error) {
{{ if not .CRD.HasReferenceFields -}}
	return res, nil
{{ else -}}
	namespace := res.MetaObject().GetNamespace()
	ko := rm.concreteResource(res).ko
	err := validateReferenceFields(ko)
{{- if $hookCode := Hook .CRD "references_pre_resolve" }}
{{ $hookCode }}
{{- end }}
	{{ range $fieldName, $field := .CRD.Fields -}}
	{{ if $field.HasReference -}}
	if err == nil {
		err = rm.resolveReferenceFor{{ $field.FieldPathWithUnderscore }}(ctx, apiReader, namespace, ko)
	}
	{{ end -}}
	{{ end -}}
{{- if $hookCode := Hook .CRD "references_post_resolve" }}
{{ $hookCode }}
{{- end }}
	return hasNonNilReferences(ko), err
{{ end -}}
}

// validateReferenceFields validates the reference field and corresponding
// identifier field.
func validateReferenceFields(ko *svcapitypes.{{ .CRD.Names.Camel }}) error {
{{ GoCodeReferencesValidation .CRD "ko" 1 -}}
	return nil
}

// hasNonNilReferences returns true if resource contains a reference to another
// resource
func hasNonNilReferences(ko *svcapitypes.{{ .CRD.Names.Camel }}) bool {
	{{ GoCodeContainsReferences .CRD "ko"}}
}

{{- $getReferencedResourceStateResources := (Nil) -}}

{{ range $fieldName, $field := .CRD.Fields }}
{{ if $field.HasReference }}
// resolveReferenceFor{{ $field.FieldPathWithUnderscore }} reads the resource referenced
// from {{ $field.ReferenceFieldPath }} field and sets the {{ $field.Path }}
// from referenced resource
func (rm *resourceManager) resolveReferenceFor{{ $field.FieldPathWithUnderscore }}(
	ctx context.Context,
	apiReader client.Reader,
	namespace string,
	ko *svcapitypes.{{ .CRD.Names.Camel }},
) error {
{{ $nilCheck := CheckNilFieldPath $field "ko.Spec" -}}
{{ if not (eq $nilCheck "") -}}
    if {{ $nilCheck }} {
        return nil
    }
{{ end -}}

{{ GoCodeResolveReference $field "ko" 1 }}
	return nil
}

{{- if not (and $getReferencedResourceStateResources (eq (index $getReferencedResourceStateResources .FieldConfig.References.Resource) "true" )) }}
{{- $getReferencedResourceStateResources = AddToMap $getReferencedResourceStateResources .FieldConfig.References.Resource "true" }}
{{ template "read_referenced_resource_and_validate" $field }}
{{ end -}}
{{ end -}}
{{ end -}}

