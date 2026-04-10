{{- define "sub_resource_manager_scalar" -}}
{{ $srcType := ManagerSourceType .CRD }}
{{ $mapper := ManagerMapper .CRD }}
func convertFromParent(parent *svcapitypes.{{ $srcType.ParentKind }}) []resource {
	if parent.{{ $srcType.FieldPath }} == nil {
		return nil
	}
	ko := &svcapitypes.{{ .CRD.Kind }}{}
{{- range $m := $mapper }}
	ko.{{ $m.To }} = parent.{{ $m.From }}
{{- end }}
	return []resource{ {ko: ko} }
}
{{- end }}
