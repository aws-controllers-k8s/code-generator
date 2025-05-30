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
    app.kubernetes.io/name: {{ IncludeTemplate "app.name" }}
    app.kubernetes.io/instance: {{ "{{ .Release.Name }}" }}
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/version: {{ "{{ .Chart.AppVersion | quote }}" }}
    k8s-app: {{ IncludeTemplate "app.name" }}
    helm.sh/chart: {{ IncludeTemplate "chart.name-version" }}
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
    app.kubernetes.io/name: {{ "{{ $fullname }}" }}
    app.kubernetes.io/instance: {{ "{{ $.Release.Name }}" }}
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/version: {{ "{{ $appVersion }}" }}
    k8s-app: {{ "{{ $fullname }}" }}
    helm.sh/chart: {{ "{{ $chartVersion }}" }}
  {{ "{{- range $key, $value := $labels }}" }}
    {{ "{{ $key }}: {{ $value | quote }}" }}
  {{ "{{- end }}" }}
{{ "{{ $rbacRules }}" }}
{{ "{{ end }}" }}
{{ "{{ end }}" }}