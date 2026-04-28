{{- /*
  sdk_create_sub_resource_requeue emits a block at the end of sdkCreate
  that flips ResourceSynced=False when any of the parent's sub-resources
  has a desired value in spec. Sub-resource writes happen via sdkUpdate
  (not the parent's create), so a resource with non-empty sub-resource
  spec on first apply is not truly synced until a follow-up reconcile.
  Marking it unsynced forces that follow-up to run promptly.
*/ -}}
{{- define "sdk_create_sub_resource_requeue" -}}
{{- $subResInfos := SubResourceManagerInfos .CRD -}}
{{- if $subResInfos }}
	// Sub-resources (fields managed by separate API operations) are not
	// applied by the parent's create call. If the user set any of them in
	// spec, flip ResourceSynced=False so the reconciler loops back
	// promptly, picks up the delta, and drives the sub-resource writes
	// through sdkUpdate.
{{- range $info := $subResInfos }}
	if desired.ko.{{ $info.FieldPath }} != nil {
		ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, nil, nil)
	}
{{- end }}
{{- end }}
{{- end -}}
