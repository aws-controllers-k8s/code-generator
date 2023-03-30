{{- define "sdk_update_set_attributes" -}}
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	var err error
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer func() {
		exit(err)
	}()

{{- if $hookCode := Hook .CRD "sdk_update_pre_build_request" }}
{{ $hookCode }}
{{- end }}
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. And sdkUpdate should never be called if this is the
	// case, and it's an error in the generated code if it is...
	if rm.requiredFieldsMissingFromSetAttributesInput(desired) {
		panic("Required field in SetAttributes input shape missing!")
	}

	input, err := rm.newSetAttributesRequestPayload(desired)
	if err != nil {
		return nil, err
	}
{{- if $hookCode := Hook .CRD "sdk_update_post_build_request" }}
{{ $hookCode }} 
{{- end }}

	// NOTE(jaypipes): SetAttributes calls return a response but they don't
	// contain any useful information. Instead, below, we'll be returning a
	// DeepCopy of the supplied desired state, which should be fine because
	// that desired state has been constructed from a call to GetAttributes...
	_, respErr := rm.sdkapi.{{ .CRD.Ops.SetAttributes.ExportedName }}WithContext(ctx, input)
{{- if $hookCode := Hook .CRD "sdk_update_post_request" }}
{{ $hookCode }}
{{- end }}
	rm.metrics.RecordAPICall("SET_ATTRIBUTES", "{{ .CRD.Ops.SetAttributes.ExportedName }}", respErr)
	if respErr != nil {
		if awsErr, ok := ackerr.AWSError(respErr); ok && awsErr.Code() == "{{ ResourceExceptionCode .CRD 404 }}" {{ GoCodeSetExceptionMessageCheck .CRD 404 }}{
			// Technically, this means someone deleted the backend resource in
			// between the time we got a result back from sdkFind() and here...
			return nil, ackerr.NotFound
		}
		return nil, respErr
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()
{{- if $hookCode := Hook .CRD "sdk_update_pre_set_output" }}
{{ $hookCode }}
{{- end }}
	rm.setStatusDefaults(ko)
{{- if $hookCode := Hook .CRD "sdk_update_post_set_output" }}
{{ $hookCode }}
{{- end }}
	return &resource{ko}, nil
}

// requiredFieldsMissingFromSetAtttributesInput returns true if there are any
// fields for the SetAttributes Input shape that are required by not present in
// the resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromSetAttributesInput(
	r *resource,
) bool {
{{- if $customCheckMethod := .CRD.GetCustomCheckRequiredFieldsMissingMethod .CRD.Ops.SetAttributes }}
return rm.{{ $customCheckMethod }}(r)
{{- else }}
{{ GoCodeRequiredFieldsMissingFromSetAttributesInput .CRD "r.ko" 1 }}
{{- end }}
}

// newSetAttributesRequestPayload returns SDK-specific struct for the HTTP
// request payload of the SetAttributes API call for the resource
func (rm *resourceManager) newSetAttributesRequestPayload(
	r *resource,
) (*svcsdk.{{ .CRD.Ops.SetAttributes.InputRef.Shape.ShapeName }}, error) {
	res := &svcsdk.{{ .CRD.Ops.SetAttributes.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetAttributesSetInput .CRD "r.ko" "res" 1 }}
	return res, nil
}
{{- end -}}
