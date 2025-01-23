{{- define "sdk_find_get_attributes" -}}
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()

{{- if $hookCode := Hook .CRD "sdk_get_attributes_pre_build_request" }}
{{ $hookCode }}
{{- end }}
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromGetAttributesInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newGetAttributesRequestPayload(r)
	if err != nil {
		return nil, err
	}
{{- if $hookCode := Hook .CRD "sdk_get_attributes_post_build_request" }}
{{ $hookCode }}
{{- end }}
	var resp {{ .CRD.GetOutputShapeGoType .CRD.Ops.GetAttributes }}
	resp, err = rm.sdkapi.{{ .CRD.Ops.GetAttributes.ExportedName }}(ctx, input)
{{- if $hookCode := Hook .CRD "sdk_get_attributes_post_request" }}
{{ $hookCode }}
{{- end }}
	rm.metrics.RecordAPICall("GET_ATTRIBUTES", "{{ .CRD.Ops.GetAttributes.ExportedName }}", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.ErrorCode() == "{{ ResourceExceptionCode .CRD 404 }}" {{ GoCodeSetExceptionMessageCheck .CRD 404 }} {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()
{{ GoCodeGetAttributesSetOutput .CRD "resp" "ko" 1 }}
{{- if $hookCode := Hook .CRD "sdk_get_attributes_pre_set_output" }}
{{ $hookCode }}
{{- end }}
	rm.setStatusDefaults(ko)
{{- if $hookCode := Hook .CRD "sdk_get_attributes_post_set_output" }}
{{ $hookCode }}
{{- end }}
	return &resource{ko}, nil
}

// requiredFieldsMissingFromGetAtttributesInput returns true if there are any
// fields for the GetAttributes Input shape that are required by not present in
// the resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromGetAttributesInput(
	r *resource,
) bool {
{{- if $customCheckMethod := .CRD.GetCustomCheckRequiredFieldsMissingMethod .CRD.Ops.GetAttributes }}
return rm.{{ $customCheckMethod }}(r)
{{- else }}
{{ GoCodeRequiredFieldsMissingFromGetAttributesInput .CRD "r.ko" 1 }}
{{- end }}
}

// newGetAttributesRequestPayload returns SDK-specific struct for the HTTP
// request payload of the GetAttributes API call for the resource
func (rm *resourceManager) newGetAttributesRequestPayload(
	r *resource,
) (*svcsdk.{{ .CRD.Ops.GetAttributes.InputRef.Shape.ShapeName }}, error) {
	res := &svcsdk.{{ .CRD.Ops.GetAttributes.InputRef.Shape.ShapeName }}{}
{{ GoCodeGetAttributesSetInput .CRD "r.ko" "res" 1 }}
	return res, nil
}
{{- end -}}
