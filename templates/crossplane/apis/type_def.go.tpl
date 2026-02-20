{{- define "type_def" -}}
// +kubebuilder:skipversion
type {{ .Names.Camel }} struct {
{{- range $attrName := .SortedAttrNames }}
{{- $attr := (index $.Attrs $attrName) }}
	{{- if $attr.Shape }}
	{{ $attr.Shape.Documentation }}
	{{- end }}
	{{ $attr.Names.Camel }} {{ $attr.GoType }} `json:"{{ $attr.Names.CamelLower }},omitempty"`
{{- end }}
}
{{- end -}}
