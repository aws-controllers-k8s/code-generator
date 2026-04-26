{{- define "sub_resource_manager_map" -}}
{{ $srcType := ManagerSourceType .CRD }}
// TODO: implement map source type convertFromParent
func convertFromParent(parent *svcapitypes.{{ $srcType.ParentKind }}) []resource {
	return nil
}
{{- end }}
