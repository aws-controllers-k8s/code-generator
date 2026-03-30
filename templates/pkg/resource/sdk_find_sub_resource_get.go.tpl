{{- define "sdk_find_sub_resource_get" -}}
{{- $subResInfos := SubResourceManagerInfos .CRD -}}
{{- range $info := $subResInfos }}
{{- if $info.ReadFieldPath }}
	{
		mgr := {{ $info.PackageName }}.NewManager(rm.sdkapi, rm.metrics)
		getResult, err := mgr.Get(ctx, ko)
		if err != nil {
			return nil, err
		}
		if getResult != nil {
			ko.{{ $info.FieldPath }}, _ = getResult.({{ $info.FieldGoType }})
		}
	}
{{- end }}
{{- end }}
{{- end -}}
