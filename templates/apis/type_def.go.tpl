{{- define "type_def" -}}
{{- if .Shape.Documentation }}
{{ .Shape.Documentation }}
{{- end }}
{{- /*
  Immutability is enforced from the containing struct, not from the member.
  Kubernetes skips a transition rule when the old value is absent, so a
  member-level "self == oldSelf" rule never fires for an optional member that
  was unset at creation. Evaluating from here lets the rule see the member's
  absence and freeze presence as well as value.
*/ -}}
{{- range $attrName := .SortedAttrNames }}
{{- $attr := (index $.Attrs $attrName) }}
{{- if $attr.IsImmutable }}
// +kubebuilder:validation:XValidation:rule="{{ $attr.ImmutabilityCELRule }}",message="Value is immutable once set",fieldPath="{{ $attr.ImmutabilityCELFieldPath }}"
{{- end }}
{{- end }}
type {{ .Names.Camel }} struct {
{{- range $attrName := .SortedAttrNames }}
{{- $attr := (index $.Attrs $attrName) }}
	{{- if $attr.Shape.Documentation }}
	{{ $attr.Shape.Documentation }}
	{{- end }}
	{{ $attr.Names.Camel }} {{ $attr.GoType }} {{ $attr.GetGoTag }}
{{- end }}
}
{{- end -}}
