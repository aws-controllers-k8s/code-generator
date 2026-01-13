{{- define "type_def" -}}
{{- if .Shape.Documentation }}
{{ .Shape.Documentation }}
{{- end }}
type {{ .Names.Camel }} struct {
{{- range $attrName, $attr := .Attrs }}
	{{- if $attr.Shape.Documentation }}
	{{ $attr.Shape.Documentation }}
	{{- end }}
	{{- if $attr.IsImmutable }}
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Value is immutable once set"
	{{- end }}
	{{ $attr.Names.Camel }} {{ $attr.GoType }} {{ $attr.GetGoTag }}
{{- end }}
}
{{- end -}}
