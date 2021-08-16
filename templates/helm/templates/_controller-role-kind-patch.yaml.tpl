apiVersion: rbac.authorization.k8s.io/v1
{{ "{{ if eq .Values.installScope \"cluster\" }}" }}
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: ack-{{ .ServiceIDClean }}-controller
{{ "{{ else }}" }}
kind: Role
metadata:
  creationTimestamp: null
  name: ack-{{ .ServiceIDClean }}-controller
  namespace: {{ "{{ .Release.Namespace }}" }}
{{ "{{ end }}" }}
