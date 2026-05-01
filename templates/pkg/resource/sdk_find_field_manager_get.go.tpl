{{- define "sdk_find_field_manager_get" -}}
{{- $mfInfos := FieldManagerInfos .CRD -}}
{{- range $info := $mfInfos }}
	mgr_{{ $info.PackageName }} := {{ $info.PackageName }}.NewManager(rm.sdkapi, rm.metrics)
	if err := mgr_{{ $info.PackageName }}.Get(ctx, ko); err != nil {
		return nil, err
	}
{{- end }}
{{- end -}}
