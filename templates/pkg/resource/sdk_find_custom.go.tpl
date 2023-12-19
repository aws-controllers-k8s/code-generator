{{- define "sdk_find_custom" -}}
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (*resource, error) {
	return rm.{{ .CRD.CustomFindMethodName }}(ctx, r)
}
{{- end -}}
