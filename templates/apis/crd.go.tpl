{{- template "boilerplate" }}

package {{ .APIVersion }}

import (
{{- if .CRD.TypeImports }}
{{- range $packagePath, $alias := .CRD.TypeImports }}
    {{ if $alias }}{{ $alias }} {{ end }}"{{ $packagePath }}"
{{ end }}

{{- end }}
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

{{ .CRD.Documentation }}
type {{ .CRD.Kind }}Spec struct {
{{ range $fieldName, $field := .CRD.SpecFields }}
{{ if $field.GetDocumentation -}}
    {{ $field.GetDocumentation }}
{{ end -}}

{{- if $field.IsImmutable -}}
    // +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable once set"
{{ end -}}

{{- if and ($field.IsRequired) (not $field.HasReference) -}}
    // +kubebuilder:validation:Required
{{ end -}}

    {{ $field.Names.Camel }} {{ $field.GoType }} {{ $field.GetGoTag }}
{{- end }}
}

// {{ .CRD.Kind }}Status defines the observed state of {{ .CRD.Kind }}
type {{ .CRD.Kind }}Status struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRs managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	{{- range $fieldName, $field := .CRD.StatusFields }}
	{{- if $field.GetDocumentation }}
	{{ $field.GetDocumentation }}
	{{- end }}
	// +kubebuilder:validation:Optional
	{{ $field.Names.Camel }} {{ $field.GoType }} {{ $field.GetGoTag }}
{{- end }}
}

// {{ .CRD.Kind }} is the Schema for the {{ .CRD.Plural }} API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:metadata:annotations="controller-version={{ .ReleaseVersion }}"
{{- range $column := .CRD.AdditionalPrinterColumns }}
// +kubebuilder:printcolumn:name="{{$column.Name}}",type={{$column.Type}},priority={{$column.Priority}},JSONPath=`{{$column.JSONPath}}`
{{- end }}
{{- if .CRD.PrintSyncedColumn }}
// +kubebuilder:printcolumn:name="Synced",type="string",priority=0,JSONPath=".status.conditions[?(@.type==\"ACK.ResourceSynced\")].status"
{{- end }}
{{- if .CRD.PrintAgeColumn }}
// +kubebuilder:printcolumn:name="Age",type="date",priority=0,JSONPath=".metadata.creationTimestamp"
{{- end }}
{{- if .CRD.ShortNames }}
// +kubebuilder:resource:shortName={{ Join .CRD.ShortNames ";" }}
{{- end }}
type {{ .CRD.Kind }} struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec   {{ .CRD.Kind }}Spec   `json:"spec,omitempty"`
	Status {{ .CRD.Kind }}Status `json:"status,omitempty"`
}

// {{ .CRD.Kind }}List contains a list of {{ .CRD.Kind }}
// +kubebuilder:object:root=true
type {{ .CRD.Kind }}List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items []{{ .CRD.Kind }} `json:"items"`
}

func init() {
	SchemeBuilder.Register(&{{ .CRD.Kind }}{}, &{{ .CRD.Kind }}List{})
}