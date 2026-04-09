{{- define "sdk_delete_sub_resource_sync" -}}
{{- $subResInfos := SubResourceManagerInfos .CRD -}}
{{- if $subResInfos }}
	// Clean up sub-resources before deleting the parent resource.
	// For each sub-resource, sync with a nil/empty desired state so all
	// items are deleted.
	koCopy := r.ko.DeepCopy()
{{- range $info := $subResInfos }}
	koCopy.{{ $info.FieldPath }} = nil
{{- end }}
{{- range $info := $subResInfos }}
	mgr_{{ $info.PackageName }} := {{ $info.PackageName }}.NewManager(rm.sdkapi, rm.metrics)
	if err = mgr_{{ $info.PackageName }}.Sync(ctx, koCopy, r.ko); err != nil {
		return nil, err
	}
{{- end }}
{{- end }}
{{- end -}}
