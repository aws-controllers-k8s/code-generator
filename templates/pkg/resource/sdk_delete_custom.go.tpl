{{- define "sdk_delete_custom" -}}
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	return rm.{{ .CRD.CustomDeleteMethodName }}(ctx, r)
}
{{- end -}}