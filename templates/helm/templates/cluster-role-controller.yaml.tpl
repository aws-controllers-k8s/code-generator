{{ "{{ $labels := .Values.role.labels }}" }}
{{ VarIncludeTemplate "rbacRules" "rbac-rules" }}
{{ VarIncludeTemplate "fullname" "app.fullname" }}
{{ "{{ if eq .Values.installScope \"cluster\" }}" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{ IncludeTemplate "app.fullname" }}
  labels:
  {{ "{{- range $key, $value := $labels }}" }}
    {{ "{{ $key }}: {{ $value | quote }}" }}
  {{ "{{- end }}" }}
{{ "{{$rbacRules }}" }}
{{ "{{ else if eq .Values.installScope \"namespace\" }}" }}
{{ VarIncludeTemplate "wn" "watch-namespace" }}
{{ "{{ $namespaces := split \",\" $wn }}" }}
{{ "{{ range $namespaces }}" }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ "{{ $fullname }}" }}-{{ "{{ . }}" }}
  namespace: {{ "{{ . }}" }}
  labels:
  {{ "{{- range $key, $value := $labels }}" }}
    {{ "{{ $key }}: {{ $value | quote }}" }}
  {{ "{{- end }}" }}
{{ "{{ $rbacRules }}" }}
{{ "{{ end }}" }}
{{ "{{ end }}" }}