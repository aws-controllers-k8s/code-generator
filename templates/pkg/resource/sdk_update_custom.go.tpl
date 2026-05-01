{{- define "sdk_update_custom" -}}
{{- $mfInfos := FieldManagerInfos .CRD -}}
{{- if $mfInfos -}}
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() {
		exit(err)
	}()
{{ template "sdk_update_field_manager_sync" . }}
	return rm.{{ .CRD.CustomUpdateMethodName }}(ctx, desired, latest, delta)
}
{{- else -}}
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	return rm.{{ .CRD.CustomUpdateMethodName }}(ctx, desired, latest, delta)
}
{{- end -}}
{{- end -}}
