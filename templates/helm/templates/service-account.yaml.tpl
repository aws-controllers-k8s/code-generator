{{ "{{- if .Values.serviceAccount.create }}" }}
apiVersion: v1
kind: ServiceAccount
metadata:
  labels:
    {{ IncludeTemplate "app.labels" | nindent 4 }}
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
  annotations:
  {{ "{{- range $key, $value := .Values.serviceAccount.annotations }}" }}
    {{ "{{ $key }}: {{ $value | quote }}" }}
  {{ "{{- end }}" }}
{{ "{{- end }}" }}
