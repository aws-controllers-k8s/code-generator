{{- define "type_def" -}}
{{- if .Shape.Documentation }}
{{ .Shape.Documentation }}
{{- end }}
type {{ .Names.Camel }} struct {
{{- range $attrName, $attr := .Attrs }}
	{{- if $attr.Shape.Documentation }}
	{{ $attr.Shape.Documentation }}
	{{- end }}
	{{ $attr.Names.Camel }} {{ $attr.GoType }} `json:"{{ $attr.Names.CamelLower }},omitempty"`
{{- end }}
}
{{- end -}}
