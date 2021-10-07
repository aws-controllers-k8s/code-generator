apiVersion: rbac.authorization.k8s.io/v1
{{ "{{ if eq .Values.installScope \"cluster\" }}" }}
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: ack-{{ .ServiceAlias }}-controller
{{ "{{ else }}" }}
kind: Role
metadata:
  creationTimestamp: null
  name: ack-{{ .ServiceAlias }}-controller
  namespace: {{ "{{ .Release.Namespace }}" }}
{{ "{{ end }}" }}
