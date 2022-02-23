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
{{ range $fieldName, $field := .CRD.Fields -}}
{{ if and $field.HasReference (not (eq $field.ReferencedServiceName $servicePackageName)) -}}
    {{ $field.ReferencedServiceName }}apitypes "github.com/aws-controllers-k8s/{{ $field.ReferencedServiceName }}-controller/apis/{{ $apiVersion }}"
{{ end -}}
{{ end -}}
{{ end -}}

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
		err = resolveReferenceFor{{ $field.Names.Camel }}(ctx, apiReader, namespace, ko)
	}
	{{ end -}}
	{{ end -}}
{{- if $hookCode := Hook .CRD "references_post_resolve" }}
{{ $hookCode }}
{{- end }}
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
	return {{ GoCodeContainsReferences .CRD "ko"}}
}

{{ range $fieldName, $field := .CRD.Fields }}
{{ if $field.HasReference }}
{{ $refField := index .CRD.Fields $field.GetReferenceFieldName.Camel }}
// resolveReferenceFor{{ $field.Names.Camel }} reads the resource referenced
// from {{ $refField.Names.Camel }} field and sets the {{ $field.Names.Camel }}
// from referenced resource
func resolveReferenceFor{{ $field.Names.Camel }}(
	ctx context.Context,
	apiReader client.Reader,
	namespace string,
	ko *svcapitypes.{{ .CRD.Names.Camel }},
) error {
{{ if eq $field.ShapeRef.Shape.Type "list" -}}
	if ko.Spec.{{ $refField.Names.Camel }} != nil &&
	   len(ko.Spec.{{ $refField.Names.Camel }}) > 0 {
		resolvedReferences := []*string{}
		for _, arrw := range ko.Spec.{{ $refField.Names.Camel }} {
			arr := arrw.From
{{ template "read_referenced_resource_and_validate" $field }}
			resolvedReferences = append(resolvedReferences,
								   obj.{{ $field.FieldConfig.References.Path }})
		}
		ko.Spec.{{ $field.Names.Camel }} = resolvedReferences
	}
	return nil
}
{{ else -}}
	if ko.Spec.{{ $refField.Names.Camel }} != nil &&
		ko.Spec.{{ $refField.Names.Camel}}.From != nil {
			arr := ko.Spec.{{ $refField.Names.Camel }}.From
{{ template "read_referenced_resource_and_validate" $field }}
			ko.Spec.{{ $field.Names.Camel }} = obj.{{ $field.FieldConfig.References.Path }}
	}
	return nil
}
{{ end -}}
{{ end -}}
{{ end -}}
