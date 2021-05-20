{{- define "sdk_update" -}}
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
{{- if $hookCode := Hook .CRD "sdk_update_pre_build_request" }}
{{ $hookCode }}
{{- end }}
{{ $customMethod := .CRD.GetCustomImplementation .CRD.Ops.Update }}
{{ if $customMethod }}
	customResp, customRespErr := rm.{{ $customMethod }}(ctx, desired, latest, delta)
	if customResp != nil || customRespErr != nil {
		return customResp, customRespErr
	}
{{ end }}

	input, err := rm.newUpdateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}
{{- if $hookCode := Hook .CRD "sdk_update_post_build_request" }}
{{ $hookCode }} 
{{- end }}
{{ $setCode := GoCodeSetUpdateOutput .CRD "resp" "ko" 1 false }}
	{{ if not ( Empty $setCode ) }}resp{{ else }}_{{ end }}, respErr := rm.sdkapi.{{ .CRD.Ops.Update.ExportedName }}WithContext(ctx, input)
{{- if $hookCode := Hook .CRD "sdk_update_post_request" }}
{{ $hookCode }}
{{- end }}
	rm.metrics.RecordAPICall("UPDATE", "{{ .CRD.Ops.Update.ExportedName }}", respErr)
	if respErr != nil {
		return nil, respErr
	}
{{- if .CRD.HasImmutableFieldChanges }}
	desired = rm.handleImmutableFieldsChangedCondition(desired, delta)
{{- end }}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()
{{- if $hookCode := Hook .CRD "sdk_update_pre_set_output" }}
{{ $hookCode }}
{{- end }}
{{ $setCode }}
	rm.setStatusDefaults(ko)
{{ if $setOutputCustomMethodName := .CRD.SetOutputCustomMethodName .CRD.Ops.Update }}
	// custom set output from response
	ko, err = rm.{{ $setOutputCustomMethodName }}(ctx, desired, resp, ko)
	if err != nil {
		return nil, err
	}
{{ end }}
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
) (*svcsdk.{{ .CRD.Ops.Update.InputRef.Shape.ShapeName }}, error) {
	res := &svcsdk.{{ .CRD.Ops.Update.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetUpdateInput .CRD "r.ko" "res" 1 }}
	return res, nil
}
{{- end -}}
