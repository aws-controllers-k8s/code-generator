{{- define "sdk_find_field_groups" -}}
func (rm *resourceManager) sdkFindFieldGroups(
	ctx context.Context,
	r *resource,
	observed *resource,
) (*resource, error) {
	ko := observed.ko.DeepCopy()
	var err error
{{ range $fg := .CRD.ReadFieldGroups }}
	ko, err = rm.sdkFind{{ $fg.OperationID }}(ctx, r, ko)
	if err != nil {
		return &resource{ko}, err
	}
{{ end }}
	return &resource{ko}, nil
}
{{ range $fg := .CRD.ReadFieldGroups }}
func (rm *resourceManager) sdkFind{{ $fg.OperationID }}(
	ctx context.Context,
	r *resource,
	ko *svcapitypes.{{ $.CRD.Names.Camel }},
) (*svcapitypes.{{ $.CRD.Names.Camel }}, error) {
{{- if $hookCode := Hook $.CRD (print "sdk_read_" $fg.Names.Snake "_pre_build_request") }}
{{ $hookCode }}
{{- end }}
	input, err := rm.new{{ $fg.OperationID }}Input(r)
	if err != nil {
		return ko, err
	}
{{- if $hookCode := Hook $.CRD (print "sdk_read_" $fg.Names.Snake "_post_build_request") }}
{{ $hookCode }}
{{- end }}

	var resp {{ $.CRD.GetOutputShapeGoType $fg.Operation }}; _ = resp
	resp, err = rm.sdkapi.{{ $fg.OperationID }}(ctx, input)
{{- if $hookCode := Hook $.CRD (print "sdk_read_" $fg.Names.Snake "_post_request") }}
{{ $hookCode }}
{{- end }}
	rm.metrics.RecordAPICall("READ_ONE", "{{ $fg.OperationID }}", err)
	if err != nil {
		var awsErr smithy.APIError
{{- if $fg.ExceptionCode404 }}
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "{{ $fg.ExceptionCode404 }}" {
			return ko, nil
		}
{{- else }}
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "{{ ResourceExceptionCode $.CRD 404 }}" {
			return ko, nil
		}
{{- end }}
		return ko, err
	}
{{- if $hookCode := Hook $.CRD (print "sdk_read_" $fg.Names.Snake "_pre_set_output") }}
{{ $hookCode }}
{{- end }}
{{ GoCodeSetFieldGroupOutput $.CRD $fg "resp" "ko" 1 }}
{{- if $hookCode := Hook $.CRD (print "sdk_read_" $fg.Names.Snake "_post_set_output") }}
{{ $hookCode }}
{{- end }}
	return ko, nil
}

func (rm *resourceManager) new{{ $fg.OperationID }}Input(
	r *resource,
) (*svcsdk.{{ $fg.Operation.InputRef.Shape.ShapeName }}, error) {
	res := &svcsdk.{{ $fg.Operation.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetFieldGroupInput $.CRD $fg "r.ko" "res" 1 }}
	return res, nil
}
{{ end }}
{{- end -}}
