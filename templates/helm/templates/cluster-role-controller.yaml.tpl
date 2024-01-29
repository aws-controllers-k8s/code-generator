{{ "{{ $labels := .Values.role.labels }}" }}
{{ "{{ $rules := include \"controller-role-rules\" . }}" }}
{{ "{{ if eq .Values.installScope \"cluster\" }}" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ack-{{ .ServicePackageName }}-controller
  labels:
  {{ "{{- range $key, $value := $labels }}" }}
    {{ "{{ $key }}: {{ $value | quote }}" }}
  {{ "{{- end }}" }}
{{ "{{- $rules }}" }}
{{ "{{ else if .Values.watchNamespace }}" }}
{{ "{{ $namespaces := split \",\" .Values.watchNamespace }}" }}
{{ "{{ range $namespaces }}" }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: ack-{{ .ServicePackageName }}-controller
  namespace: {{ "{{ . }}" }}
  labels:
  {{ "{{- range $key, $value := $labels }}" }}
    {{ "{{ $key }}: {{ $value | quote }}" }}
  {{ "{{- end }}" }}
{{ "{{- $rules }}" }}
{{ "{{ end }}" }}
{{ "{{ end }}" }}