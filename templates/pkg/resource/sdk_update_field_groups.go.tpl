{{- define "sdk_update_field_groups" -}}
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
{{- if $hookCode := Hook .CRD "sdk_update_pre_build_request" }}
{{ $hookCode }}
{{- end }}

	ko := desired.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	var requeueNeeded bool
{{ range $fg := .CRD.UpdateFieldGroups }}
	if {{ GoCodeFieldGroupDeltaCheck $.CRD $fg "delta" }} {
{{- if $fg.IsSync }}
		err = rm.sync{{ $fg.OperationID }}(ctx, desired, latest)
{{- else }}
		ko, err = rm.sdkUpdate{{ $fg.OperationID }}(ctx, desired, ko)
{{- end }}
		if err != nil {
			return nil, err
		}
{{- if $fg.Config.RequeueOnSuccess }}
		requeueNeeded = true
{{- end }}
	}
{{ end }}
{{- if $hookCode := Hook .CRD "sdk_update_post_set_output" }}
{{ $hookCode }}
{{- end }}
	if requeueNeeded {
		return &resource{ko}, ackrequeue.NeededAfter(nil, 0)
	}
	return &resource{ko}, nil
}
{{ range $fg := .CRD.UpdateFieldGroups }}
{{- if $fg.IsSyncList }}
// sync{{ $fg.OperationID }} examines the desired and latest values of the
// {{ (index $fg.PayloadFields 0).Names.Camel }} field and calls the {{ $fg.OperationID }} and
// {{ $fg.RemoveOperationID }} APIs to bring the observed state in line with desired.
func (rm *resourceManager) sync{{ $fg.OperationID }}(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sync{{ $fg.OperationID }}")
	defer func() { exit(err) }()

	toAdd, toRemove := ackcompare.SliceStringPDifference(
		desired.ko.Spec.{{ (index $fg.PayloadFields 0).Names.Camel }},
		latest.ko.Spec.{{ (index $fg.PayloadFields 0).Names.Camel }},
	)

	for _, item := range toAdd {
		rlog.Debug("adding item via {{ $fg.OperationID }}", "item", *item)
		input := &svcsdk.{{ $fg.Operation.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetSyncInput $.CRD $fg "desired" "input" "item" 2 }}
		_, err = rm.sdkapi.{{ $fg.OperationID }}(ctx, input)
		rm.metrics.RecordAPICall("UPDATE", "{{ $fg.OperationID }}", err)
		if err != nil {
			return err
		}
	}
	for _, item := range toRemove {
		rlog.Debug("removing item via {{ $fg.RemoveOperationID }}", "item", *item)
		input := &svcsdk.{{ $fg.RemoveOperation.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetSyncInput $.CRD $fg "desired" "input" "item" 2 }}
		_, err = rm.sdkapi.{{ $fg.RemoveOperationID }}(ctx, input)
		rm.metrics.RecordAPICall("UPDATE", "{{ $fg.RemoveOperationID }}", err)
		if err != nil {
			return err
		}
	}
	return nil
}
{{- else if $fg.IsSyncMap }}
// sync{{ $fg.OperationID }} examines the desired and latest values of the
// {{ (index $fg.PayloadFields 0).Names.Camel }} field and calls the {{ $fg.OperationID }} and
// {{ $fg.RemoveOperationID }} APIs to bring the observed state in line with desired.
func (rm *resourceManager) sync{{ $fg.OperationID }}(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sync{{ $fg.OperationID }}")
	defer func() { exit(err) }()

	toAddOrUpdate, toRemove := ackcompare.MapStringStringPDifference(
		desired.ko.Spec.{{ (index $fg.PayloadFields 0).Names.Camel }},
		latest.ko.Spec.{{ (index $fg.PayloadFields 0).Names.Camel }},
	)

	for key, val := range toAddOrUpdate {
		rlog.Debug("adding/updating item via {{ $fg.OperationID }}", "key", key)
		input := &svcsdk.{{ $fg.Operation.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetSyncMapAddInput $.CRD $fg "desired" "input" "key" "val" 2 }}
		_, err = rm.sdkapi.{{ $fg.OperationID }}(ctx, input)
		rm.metrics.RecordAPICall("UPDATE", "{{ $fg.OperationID }}", err)
		if err != nil {
			return err
		}
	}
	for _, key := range toRemove {
		rlog.Debug("removing item via {{ $fg.RemoveOperationID }}", "key", key)
		input := &svcsdk.{{ $fg.RemoveOperation.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetSyncMapRemoveInput $.CRD $fg "desired" "input" "key" 2 }}
		_, err = rm.sdkapi.{{ $fg.RemoveOperationID }}(ctx, input)
		rm.metrics.RecordAPICall("UPDATE", "{{ $fg.RemoveOperationID }}", err)
		if err != nil {
			return err
		}
	}
	return nil
}
{{- else }}
func (rm *resourceManager) sdkUpdate{{ $fg.OperationID }}(
	ctx context.Context,
	desired *resource,
	ko *svcapitypes.{{ $.CRD.Names.Camel }},
) (*svcapitypes.{{ $.CRD.Names.Camel }}, error) {
{{- if $hookCode := Hook $.CRD (print "sdk_update_" $fg.Names.Snake "_pre_build_request") }}
{{ $hookCode }}
{{- end }}
	input, err := rm.new{{ $fg.OperationID }}Payload(desired)
	if err != nil {
		return ko, err
	}
{{- if $hookCode := Hook $.CRD (print "sdk_update_" $fg.Names.Snake "_post_build_request") }}
{{ $hookCode }}
{{- end }}

	var resp {{ $.CRD.GetOutputShapeGoType $fg.Operation }}; _ = resp
	resp, err = rm.sdkapi.{{ $fg.OperationID }}(ctx, input)
{{- if $hookCode := Hook $.CRD (print "sdk_update_" $fg.Names.Snake "_post_request") }}
{{ $hookCode }}
{{- end }}
	rm.metrics.RecordAPICall("UPDATE", "{{ $fg.OperationID }}", err)
	if err != nil {
		return ko, err
	}
{{- if $hookCode := Hook $.CRD (print "sdk_update_" $fg.Names.Snake "_pre_set_output") }}
{{ $hookCode }}
{{- end }}
{{ GoCodeSetFieldGroupOutput $.CRD $fg "resp" "ko" 1 }}
{{- if $hookCode := Hook $.CRD (print "sdk_update_" $fg.Names.Snake "_post_set_output") }}
{{ $hookCode }}
{{- end }}
	return ko, nil
}

func (rm *resourceManager) new{{ $fg.OperationID }}Payload(
	r *resource,
) (*svcsdk.{{ $fg.Operation.InputRef.Shape.ShapeName }}, error) {
	res := &svcsdk.{{ $fg.Operation.InputRef.Shape.ShapeName }}{}
{{ GoCodeSetFieldGroupInput $.CRD $fg "r.ko" "res" 1 }}
	return res, nil
}
{{- end }}
{{ end }}
{{- end -}}
