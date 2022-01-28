{{/*
read_referenced_resource_and_validate" template should be invoked with a field as
parameter
Ex: {{ template "read_referenced_resource_and_validate" $field }}
Where field is of type 'Field' from aws-controllers-k8s/code-generator/pkg/model
 */}}
{{- define "read_referenced_resource_and_validate" -}}
			if arr == nil || arr.Name == nil || *arr.Name == "" {
				return fmt.Errorf("provided resource reference is nil or empty")
			}
			namespacedName := types.NamespacedName{
				Namespace: namespace,
				Name: *arr.Name,
			}
			obj := svcapitypes.{{ .FieldConfig.References.Resource }}{}
			err := apiReader.Get(ctx, namespacedName, &obj)
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
					refResourceTerminal = true
				}
			}
			if refResourceTerminal {
				return ackerr.ResourceReferenceTerminalFor(
					"{{ .FieldConfig.References.Resource }}",
					namespace, *arr.Name)
			}
			if !refResourceSynced {
				return ackerr.ResourceReferenceNotSyncedFor(
					"{{ .FieldConfig.References.Resource }}",
					namespace, *arr.Name)
			}
			if obj.{{ .FieldConfig.References.Path }} == nil {
				return ackerr.ResourceReferenceMissingTargetFieldFor(
					"{{ .FieldConfig.References.Resource }}",
					namespace, *arr.Name,
					"{{ .FieldConfig.References.Path }}")
			}
{{- end -}}