{{- define "sdk_update_sub_resource_sync" -}}
{{- $subResInfos := SubResourceManagerInfos .CRD -}}
{{- if $subResInfos }}
	// Sync sub-resource managers for fields managed by separate API operations.
{{- range $info := $subResInfos }}
	if delta.DifferentAt("{{ $info.FieldPath }}") {
		mgr_{{ $info.PackageName }} := {{ $info.PackageName }}.NewManager(rm.sdkapi, rm.metrics)
		if err = mgr_{{ $info.PackageName }}.Sync(ctx, desired.ko, latest.ko); err != nil {
			return nil, err
		}
	}
{{- end }}
	if !delta.DifferentExcept({{ range $i, $info := $subResInfos }}{{ if $i }}, {{ end }}"{{ $info.FieldPath }}"{{ end }}) {
		return desired, nil
	}
{{- end }}
{{- end -}}
