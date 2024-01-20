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

	svcapitypes "github.com/aws-controllers-k8s/{{ .ControllerName }}-controller/apis/{{ .APIVersion }}"
)

{{ if .CRD.HasReferenceFields -}}
{{ range $fieldName, $field := .CRD.Fields -}}
{{ if and $field.HasReference (not (eq $field.ReferencedServiceName $servicePackageName)) -}}
// +kubebuilder:rbac:groups={{ $field.ReferencedServiceName -}}.services.k8s.aws,resources={{ ToLower $field.ReferencedResourceNamePlural }},verbs=get;list
// +kubebuilder:rbac:groups={{ $field.ReferencedServiceName -}}.services.k8s.aws,resources={{ ToLower $field.ReferencedResourceNamePlural }}/status,verbs=get;list

{{ end -}}
{{ end -}}
{{ end -}}

// ClearResolvedReferences removes any reference values that were made
// concrete in the spec. It returns a copy of the input AWSResource which
// contains the original *Ref values, but none of their respective concrete
// values.
func (rm *resourceManager) ClearResolvedReferences(res acktypes.AWSResource) (acktypes.AWSResource) {
	ko := rm.concreteResource(res).ko.DeepCopy()

{{ range $fieldName, $field := .CRD.Fields -}}
{{ if $field.HasReference -}}
{{ GoCodeClearResolvedReferences $field "ko" 1 }}
{{ end -}}
{{ end -}}
	return &resource{ko}
}

// ResolveReferences finds if there are any Reference field(s) present
// inside AWSResource passed in the parameter and attempts to resolve those
// reference field(s) into their respective target field(s). It returns a
// copy of the input AWSResource with resolved reference(s), a boolean which
// is set to true if the resource contains any references (regardless of if
// they are resolved successfully) and an error if the passed AWSResource's
// reference field(s) could not be resolved.
func (rm *resourceManager) ResolveReferences(
	ctx context.Context,
	apiReader client.Reader,
	res acktypes.AWSResource,
) (acktypes.AWSResource, bool, error) {
{{ if not .CRD.HasReferenceFields -}}
	return res, false, nil
{{ else -}}
	namespace := res.MetaObject().GetNamespace()
	ko := rm.concreteResource(res).ko

	resourceHasReferences := false
	err := validateReferenceFields(ko)
{{- if $hookCode := Hook .CRD "references_pre_resolve" }}
{{ $hookCode }}
{{- end }}
	{{ range $fieldName, $field := .CRD.Fields -}}
	{{ if $field.HasReference -}}
	if fieldHasReferences, err := rm.resolveReferenceFor{{ $field.FieldPathWithUnderscore }}(ctx, apiReader, namespace, ko); err != nil {
		return &resource{ko}, (resourceHasReferences || fieldHasReferences), err
	} else {
		resourceHasReferences = resourceHasReferences || fieldHasReferences
	}
	
	{{ end -}}
	{{ end -}}
{{- if $hookCode := Hook .CRD "references_post_resolve" }}
{{ $hookCode }}
{{- end }}
	return &resource{ko}, resourceHasReferences, err
{{ end -}}
}

// validateReferenceFields validates the reference field and corresponding
// identifier field.
func validateReferenceFields(ko *svcapitypes.{{ .CRD.Names.Camel }}) error {
{{ range $fieldName, $field := .CRD.Fields -}}
{{ if $field.HasReference }}
{{ GoCodeReferencesValidation $field "ko" 1 -}}
{{ end -}}
{{ end -}}
	return nil
}

{{- $getReferencedResourceStateResources := (Nil) -}}

{{ range $fieldName, $field := .CRD.Fields -}}
{{ if $field.HasReference }}
// resolveReferenceFor{{ $field.FieldPathWithUnderscore }} reads the resource referenced
// from {{ $field.ReferenceFieldPath }} field and sets the {{ $field.Path }}
// from referenced resource. Returns a boolean indicating whether a reference
// contains references, or an error 
func (rm *resourceManager) resolveReferenceFor{{ $field.FieldPathWithUnderscore }}(
	ctx context.Context,
	apiReader client.Reader,
	namespace string,
	ko *svcapitypes.{{ .CRD.Names.Camel }},
) (hasReferences bool, err error) {
{{ GoCodeResolveReference $field "ko" 1 }}
	return hasReferences, nil
}

{{- if not (and $getReferencedResourceStateResources (eq (index $getReferencedResourceStateResources .FieldConfig.References.Resource) "true" )) }}
{{- $getReferencedResourceStateResources = AddToMap $getReferencedResourceStateResources .FieldConfig.References.Resource "true" }}
{{ template "read_referenced_resource_and_validate" $field }}
{{ end -}}
{{ end -}}
{{ end -}}

