{{- define "sdk_delete_custom" -}}
	return rm.{{ .CRD.CustomDeleteMethodName }}(ctx, r)
{{- end -}}