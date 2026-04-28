{{- define "sdk_update_not_implemented" -}}
{{- $subResInfos := SubResourceManagerInfos .CRD -}}
{{- if $subResInfos -}}
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
{{ template "sdk_update_sub_resource_sync" . }}
	// No Update API for this resource — any non-sub-resource delta is a
	// terminal error.
	return nil, ackerr.NewTerminalError(ackerr.NotImplemented)
}
{{- else -}}
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	return nil, ackerr.NewTerminalError(ackerr.NotImplemented)
}
{{- end -}}
{{- end -}}
