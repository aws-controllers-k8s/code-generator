{{ "{{- if .Values.metrics.service.create }}" }}
apiVersion: v1
kind: Service
metadata:
  name: {{ "{{ .Chart.Name | trimSuffix \"-chart\" | trunc 44 }}-controller-metrics" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
  labels:
    {{ IncludeTemplate "app.labels" | nindent 4 }}
spec:
  selector:
    {{ IncludeTemplate "app.selectorLabels" | nindent 4 }}
{{ "{{- range $key, $value := .Values.deployment.labels }}" }}
    {{ "{{ $key }}: {{ $value | quote }}" }}
{{ "{{- end }}" }}
  type: {{ "{{ .Values.metrics.service.type }}" }}
  ports:
  - name: metricsport
    port: 8080
    targetPort: http
    protocol: TCP
{{ "{{- end }}" }}
