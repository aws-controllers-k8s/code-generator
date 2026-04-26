{{- define "sub_resource_manager_map_struct" -}}
{{ $srcType := ManagerSourceType .CRD }}
// TODO: implement map struct source type convertFromParent
func convertFromParent(parent *svcapitypes.{{ $srcType.ParentKind }}) []resource {
	return nil
}
{{- end }}
