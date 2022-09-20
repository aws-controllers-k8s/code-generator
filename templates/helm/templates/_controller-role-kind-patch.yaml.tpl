apiVersion: rbac.authorization.k8s.io/v1
{{ "{{ if eq .Values.installScope \"cluster\" }}" }}
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: ack-{{ .ServicePackageName }}-controller
  labels:
  {{ "{{- range $key, $value := .Values.clusterRole.labels }}" }}
    {{ "{{ $key }}: {{ $value | quote }}" }}
  {{ "{{- end }}" }}
{{ "{{ else }}" }}
kind: Role
metadata:
  creationTimestamp: null
  name: ack-{{ .ServicePackageName }}-controller
  namespace: {{ "{{ .Release.Namespace }}" }}
{{ "{{ end }}" }}
