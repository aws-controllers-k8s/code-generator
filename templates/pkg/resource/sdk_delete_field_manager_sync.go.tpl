{{- define "sdk_delete_field_manager_sync" -}}
{{- $mfInfos := FieldManagerInfos .CRD -}}
{{- if $mfInfos }}
	// Clean up managed fields before deleting the parent resource.
	// For each managed field, sync with a nil/empty desired state so all
	// items are deleted.
	koCopy := r.ko.DeepCopy()
{{- range $info := $mfInfos }}
	koCopy.{{ $info.FieldPath }} = nil
{{- end }}
{{- range $info := $mfInfos }}
	mgr_{{ $info.PackageName }} := {{ $info.PackageName }}.NewManager(rm.sdkapi, rm.metrics)
	if err = mgr_{{ $info.PackageName }}.Sync(ctx, koCopy, r.ko); err != nil {
		return nil, err
	}
{{- end }}
{{- end }}
{{- end -}}
