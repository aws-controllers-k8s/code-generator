{{ "{{ $labels := .Values.role.labels }}" }}
{{ "{{ $appVersion := .Chart.AppVersion | quote }}" }}
{{ VarIncludeTemplate "rbacRules" "rbac-rules" }}
{{ VarIncludeTemplate "fullname" "app.fullname" }}
{{ VarIncludeTemplate "chartVersion" "chart.name-version" }}
{{ "{{ if eq .Values.installScope \"cluster\" }}" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{ IncludeTemplate "app.fullname" }}
  labels:
    {{ IncludeTemplate "app.labels" | nindent 4 }}
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
    {{ IncludeTemplate "app.labels" | nindent 4 }}
  {{ "{{- range $key, $value := $labels }}" }}
    {{ "{{ $key }}: {{ $value | quote }}" }}
  {{ "{{- end }}" }}
{{ "{{ $rbacRules }}" }}
{{ "{{ end }}" }}
{{ "{{ end }}" }}