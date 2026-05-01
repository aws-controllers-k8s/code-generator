{{- define "sdk_update_field_manager_sync" -}}
{{- $mfInfos := FieldManagerInfos .CRD -}}
{{- if $mfInfos }}
	// Sync field managers for fields managed by separate API operations.
{{- range $info := $mfInfos }}
	if delta.DifferentAt("{{ $info.FieldPath }}") {
		mgr_{{ $info.PackageName }} := {{ $info.PackageName }}.NewManager(rm.sdkapi, rm.metrics)
		if err = mgr_{{ $info.PackageName }}.Sync(ctx, desired.ko, latest.ko); err != nil {
			return nil, err
		}
	}
{{- end }}
	if !delta.DifferentExcept({{ range $i, $info := $mfInfos }}{{ if $i }}, {{ end }}"{{ $info.FieldPath }}"{{ end }}) {
		return desired, nil
	}
{{- end }}
{{- end -}}
