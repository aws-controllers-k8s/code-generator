apiVersion: rbac.authorization.k8s.io/v1
{{ "{{ if not .Values.namespacedInstallation }}" }}
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: ack-{{ .ServiceIDClean }}-controller
{{ "{{ else if and .Values.namespacedInstallation (required \"watchNamespace must be set for namespaced installation\"  .Values.watchNamespace) }}" }}
kind: Role
metadata:
  creationTimestamp: null
  name: ack-{{ .ServiceIDClean }}-controller
  namespace: {{ "{{ .Release.Namespace }}" }}
{{ "{{ end }}" }}
