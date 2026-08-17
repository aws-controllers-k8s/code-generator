{{- define "sdk_update_custom" -}}
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
{{- GoCodeResourceIsUpdateable .CRD "latest" 1 }}
	return rm.{{ .CRD.CustomUpdateMethodName }}(ctx, desired, latest, delta)
}
{{- end -}}
