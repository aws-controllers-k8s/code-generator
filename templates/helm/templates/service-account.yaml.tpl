{{ "{{- if .Values.serviceAccount.create }}" }}
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    app.kubernetes.io/name: {{ IncludeTemplate "app.name" }}
    app.kubernetes.io/instance: {{ "{{ .Release.Name }}" }}
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/version: {{ "{{ .Chart.AppVersion | quote }}" }}
    k8s-app: {{ IncludeTemplate "app.name" }}
    helm.sh/chart: {{ IncludeTemplate "chart.name-version" }}
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
  annotations:
  {{ "{{- range $key, $value := .Values.serviceAccount.annotations }}" }}
    {{ "{{ $key }}: {{ $value | quote }}" }}
  {{ "{{- end }}" }}
{{ "{{- end }}" }}
