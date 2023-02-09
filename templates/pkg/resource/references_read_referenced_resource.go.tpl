{{/*
"read_referenced_resource_and_validate" template should be invoked with a field as
parameter
Ex: {{ template "read_referenced_resource_and_validate" $field }}
Where field is of type 'Field' from aws-controllers-k8s/code-generator/pkg/model
 */}}
{{- define "read_referenced_resource_and_validate" -}}
{{- $objType := ( printf "%sapitypes.%s" .ReferencedServiceName .FieldConfig.References.Resource ) -}}
{{- if eq .FieldConfig.References.ServiceName "" -}}
{{- $objType = ( printf "svcapitypes.%s" .FieldConfig.References.Resource) -}}
{{ end -}}
// getReferencedResourceState_{{ .FieldConfig.References.Resource }} looks up whether a referenced resource
// exists and is in a ACK.ResourceSynced=True state. If the referenced resource does exist and is
// in a Synced state, returns nil, otherwise returns `ackerr.ResourceReferenceTerminalFor` or
// `ResourceReferenceNotSyncedFor` depending on if the resource is in a Terminal state.
func getReferencedResourceState_{{ .FieldConfig.References.Resource }}(
	ctx context.Context,
	apiReader client.Reader,
	obj *{{ $objType }},
	name string, // the Kubernetes name of the referenced resource
	namespace string, // the Kubernetes namespace of the referenced resource
) error {
	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name: name,
	}
	err := apiReader.Get(ctx, namespacedName, obj)
	if err != nil {
		return err
	}
	var refResourceSynced, refResourceTerminal bool
	for _, cond := range obj.Status.Conditions {
		if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
			cond.Status == corev1.ConditionTrue {
			refResourceSynced = true
		}
		if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
			cond.Status == corev1.ConditionTrue {
			return ackerr.ResourceReferenceTerminalFor(
				"{{ .FieldConfig.References.Resource }}",
				namespace, name)
		}
	}
	if refResourceTerminal {
		return ackerr.ResourceReferenceTerminalFor(
			"{{ .FieldConfig.References.Resource }}",
			namespace, name)
	}
	if !refResourceSynced {
		return ackerr.ResourceReferenceNotSyncedFor(
			"{{ .FieldConfig.References.Resource }}",
			namespace, name)
	}
	if {{ CheckNilReferencesPath . "obj" }} {
		return ackerr.ResourceReferenceMissingTargetFieldFor(
			"{{ .FieldConfig.References.Resource }}",
			namespace, name,
			"{{ .FieldConfig.References.Path }}")
	}
	return nil
}
{{- end -}}
