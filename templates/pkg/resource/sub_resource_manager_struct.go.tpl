{{- define "sub_resource_manager_struct" -}}
{{ $srcType := ManagerSourceType .CRD }}
// TODO: implement struct source type convertFromParent
func convertFromParent(parent *svcapitypes.{{ $srcType.ParentKind }}) []resource {
	return nil
}
{{- end }}
