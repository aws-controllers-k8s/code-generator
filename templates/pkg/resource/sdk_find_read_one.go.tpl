{{- define "sdk_find_read_one" -}}
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()

{{- if $hookCode := Hook .CRD "sdk_read_one_pre_build_request" }}
{{ $hookCode }}
{{- end }}
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadOneInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newDescribeRequestPayload(r)
	if err != nil {
		return nil, err
	}
{{- if $hookCode := Hook .CRD "sdk_read_one_post_build_request" }}
{{ $hookCode }}
{{- end }}

	var resp {{ .CRD.GetOutputShapeGoType .CRD.Ops.ReadOne }}
	resp, err = rm.sdkapi.{{ .CRD.Ops.ReadOne.ExportedName }}WithContext(ctx, input)
{{- if $hookCode := Hook .CRD "sdk_read_one_post_request" }}
{{ $hookCode }}
{{- end }}
	rm.metrics.RecordAPICall("READ_ONE", "{{ .CRD.Ops.ReadOne.ExportedName }}", err)
	if err != nil {
		if reqErr, ok := ackerr.AWSRequestFailure(err); ok && reqErr.StatusCode() == 404 {
			return nil, ackerr.NotFound
        }
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "{{ ResourceExceptionCode .CRD 404 }}" {{ GoCodeSetExceptionMessageCheck .CRD 404 }}{
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()
{{- if $hookCode := Hook .CRD "sdk_read_one_pre_set_output" }}
{{ $hookCode }}
{{- end }}
{{ GoCodeSetReadOneOutput .CRD "resp" "ko" 1 }}
	rm.setStatusDefaults(ko)
{{- if $setOutputCustomMethodName := .CRD.SetOutputCustomMethodName .CRD.Ops.ReadOne }}
	// custom set output from response
	ko, err = rm.{{ $setOutputCustomMethodName }}(ctx, r, resp, ko)
	if err != nil {
		return nil, err
	}
{{- end }}
{{- if $hookCode := Hook .CRD "sdk_read_one_post_set_output" }}
{{ $hookCode }}
{{- end }}
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadOneInput returns true if there are any fields
// for the ReadOne Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadOneInput(
	r *resource,
) bool {
{{- if $customCheckMethod := .CRD.GetCustomCheckRequiredFieldsMissingMethod .CRD.Ops.ReadOne }}
return rm.{{ $customCheckMethod }}(r)
{{- else }}
{{ GoCodeRequiredFieldsMissingFromReadOneInput .CRD "r.ko" 1 }}
{{- end }}
}

// newDescribeRequestPayload returns SDK-specific struct for the HTTP request
// payload of the Describe API call for the resource
func (rm *resourceManager) newDescribeRequestPayload(
	r *resource,
) (*svcsdk.{{ .CRD.Ops.ReadOne.InputRef.Shape.ShapeName }}, error) {
	res := &svcsdk.{{ .CRD.Ops.ReadOne.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetReadOneInput .CRD "r.ko" "res" 1 }}
	return res, nil
}
{{- end -}}
