{{- define "sdk_find_read_many" -}}
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer exit(err)

{{- if $hookCode := Hook .CRD "sdk_read_many_pre_build_request" }}
{{ $hookCode }}
{{- end }}
	input, err := rm.newListRequestPayload(r)
	if err != nil {
		return nil, err
	}
{{- if $hookCode := Hook .CRD "sdk_read_many_post_build_request" }}
{{ $hookCode }}
{{- end }}
	var resp {{ .CRD.GetOutputShapeGoType .CRD.Ops.ReadMany }}
	resp, err = rm.sdkapi.{{ .CRD.Ops.ReadMany.ExportedName }}WithContext(ctx, input)
{{- if $hookCode := Hook .CRD "sdk_read_many_post_request" }}
{{ $hookCode }}
{{- end }}
	rm.metrics.RecordAPICall("READ_MANY", "{{ .CRD.Ops.ReadMany.ExportedName }}", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "{{ ResourceExceptionCode .CRD 404 }}" {{ GoCodeSetExceptionMessageCheck .CRD 404 }}{
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()
{{- if $hookCode := Hook .CRD "sdk_read_many_pre_set_output" }}
{{ $hookCode }}
{{- end }}
{{ GoCodeSetReadManyOutput .CRD "resp" "ko" 1 }}
	rm.setStatusDefaults(ko)
{{- if $setOutputCustomMethodName := .CRD.SetOutputCustomMethodName .CRD.Ops.ReadMany }}
	// custom set output from response
	ko, err = rm.{{ $setOutputCustomMethodName }}(ctx, r, resp, ko)
	if err != nil {
		return nil, err
	}
{{- end }}
{{- if $hookCode := Hook .CRD "sdk_read_many_post_set_output" }}
{{ $hookCode }}
{{- end }}
	return &resource{ko}, nil
}

// newListRequestPayload returns SDK-specific struct for the HTTP request
// payload of the List API call for the resource
func (rm *resourceManager) newListRequestPayload(
	r *resource,
) (*svcsdk.{{ .CRD.Ops.ReadMany.InputRef.Shape.ShapeName }}, error) {
	res := &svcsdk.{{ .CRD.Ops.ReadMany.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetReadManyInput .CRD "r.ko" "res" 1 }}
	return res, nil
}
{{- end -}}
