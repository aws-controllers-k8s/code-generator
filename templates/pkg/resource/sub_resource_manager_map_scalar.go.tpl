{{- define "sub_resource_manager_map_scalar" -}}
{{ $srcType := ManagerSourceType .CRD }}
{{ $mapper := ManagerMapper .CRD }}
func convertFromParent(parent *svcapitypes.{{ $srcType.ParentKind }}) []resource {
	items := parent.{{ $srcType.FieldPath }}
	if len(items) == 0 {
		return nil
	}
	var resources []resource
	for mapKey, mapValue := range items {
		ko := &svcapitypes.{{ .CRD.Kind }}{}
{{- range $m := $mapper }}
{{- if eq $m.From "$key" }}
		{
			k := mapKey
			ko.{{ $m.To }} = &k
		}
{{- else if eq $m.From "$value" }}
		ko.{{ $m.To }} = mapValue
{{- else }}
		ko.{{ $m.To }} = parent.{{ $m.From }}
{{- end }}
{{- end }}
		resources = append(resources, resource{ko: ko})
	}
	return resources
}
{{- end }}
