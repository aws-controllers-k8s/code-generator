{{- define "sdk_update" -}}
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
{{- if .CRD.HasImmutableFieldChanges }}
    if immutableFieldChanges := rm.getImmutableFieldChanges(delta); len(immutableFieldChanges) > 0 {
        msg := fmt.Sprintf("Immutable Spec fields have been modified: %s", strings.Join(immutableFieldChanges, ","))
        return nil, ackerr.NewTerminalError(fmt.Errorf(msg))
    }
{{- end }}
{{- if $hookCode := Hook .CRD "sdk_update_pre_build_request" }}
{{ $hookCode }}
{{- end }}
{{- if $customMethod := .CRD.GetCustomImplementation .CRD.Ops.Update }}
	updated, err = rm.{{ $customMethod }}(ctx, desired, latest, delta)
	if updated != nil || err != nil {
		return updated, err
	}
{{- end }}
	input, err := rm.newUpdateRequestPayload(ctx, desired, delta)
	if err != nil {
		return nil, err
	}
{{- if $hookCode := Hook .CRD "sdk_update_post_build_request" }}
{{ $hookCode }} 
{{- end }}

	var resp {{ .CRD.GetOutputShapeGoType .CRD.Ops.Update }}; _ = resp;
	resp, err = rm.sdkapi.{{ .CRD.Ops.Update.ExportedName }}(ctx, input)
{{- if $hookCode := Hook .CRD "sdk_update_post_request" }}
{{ $hookCode }}
{{- end }}
	rm.metrics.RecordAPICall("UPDATE", "{{ .CRD.Ops.Update.ExportedName }}", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()
{{- if $hookCode := Hook .CRD "sdk_update_pre_set_output" }}
{{ $hookCode }}
{{- end }}
{{ GoCodeSetUpdateOutput .CRD "resp" "ko" 1 }}
	rm.setStatusDefaults(ko)
{{- if $setOutputCustomMethodName := .CRD.SetOutputCustomMethodName .CRD.Ops.Update }}
	// custom set output from response
	ko, err = rm.{{ $setOutputCustomMethodName }}(ctx, desired, resp, ko)
	if err != nil {
		return nil, err
	}
{{- end }}
{{- if $hookCode := Hook .CRD "sdk_update_post_set_output" }}
{{ $hookCode }}
{{- end }}
	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	ctx context.Context,
	r *resource,
	delta *ackcompare.Delta,
) (*svcsdk.{{ .CRD.Ops.Update.InputRef.Shape.ShapeName }}, error) {
	res := &svcsdk.{{ .CRD.Ops.Update.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetUpdateInput .CRD "r.ko" "res" 1 }}
	return res, nil
}
{{- end -}}