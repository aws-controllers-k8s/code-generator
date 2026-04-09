{{- define "sdk_find_sub_resource_get" -}}
{{- $subResInfos := SubResourceManagerInfos .CRD -}}
{{- range $info := $subResInfos }}
	mgr_{{ $info.PackageName }} := {{ $info.PackageName }}.NewManager(rm.sdkapi, rm.metrics)
	if err := mgr_{{ $info.PackageName }}.Get(ctx, ko); err != nil {
		return nil, err
	}
{{- end }}
{{- end -}}
