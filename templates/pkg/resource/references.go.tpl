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
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
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

// ResolveReferences finds if there are any Reference field(s) present
// inside AWSResource passed in the parameter and attempts to resolve
// those reference field(s) into target field(s).
// It returns an AWSResource with resolved reference(s), and an error if the
// passed AWSResource's reference field(s) cannot be resolved.
// This method also adds/updates the ConditionTypeReferencesResolved for the
// AWSResource.
func (rm *resourceManager) ResolveReferences(
	ctx context.Context,
	apiReader client.Reader,
	res acktypes.AWSResource,
) (acktypes.AWSResource, error) {
{{ if not .CRD.HasReferenceFields -}}
	return res, nil
{{ else -}}
	namespace := res.MetaObject().GetNamespace()
	ko := rm.concreteResource(res).ko.DeepCopy()
	err := validateReferenceFields(ko)
{{- if $hookCode := Hook .CRD "references_pre_resolve" }}
{{ $hookCode }}
{{- end }}
	{{ range $fieldName, $field := .CRD.Fields -}}
	{{ if $field.HasReference -}}
	if err == nil {
		err = resolveReferenceFor{{ $field.FieldPathWithUnderscore }}(ctx, apiReader, namespace, ko)
	}
	{{ end -}}
	{{ end -}}
{{- if $hookCode := Hook .CRD "references_post_resolve" }}
{{ $hookCode }}
{{- end }}
	// If there was an error while resolving any reference, reset all the
	// resolved values so that they do not get persisted inside etcd
	if err != nil {
		ko = rm.concreteResource(res).ko.DeepCopy()
	}
	if hasNonNilReferences(ko) {
		return ackcondition.WithReferencesResolvedCondition(&resource{ko}, err)
	}
	return &resource{ko}, err
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

{{ range $fieldName, $field := .CRD.Fields }}
{{ if $field.HasReference }}
// resolveReferenceFor{{ $field.FieldPathWithUnderscore }} reads the resource referenced
// from {{ $field.ReferenceFieldPath }} field and sets the {{ $field.Path }}
// from referenced resource
func resolveReferenceFor{{ $field.FieldPathWithUnderscore }}(
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

{{ GoCodeResolveReference .CRD $field "ko" 1 }}
	return nil
}
{{ end -}}
{{ end -}}