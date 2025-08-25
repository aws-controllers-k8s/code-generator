---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  creationTimestamp: null
  name: {{ IncludeTemplate "app.fullname" }}-reader
  namespace: {{ "{{ .Release.Namespace }}" }}
  labels:
    {{ IncludeTemplate "app.labels" | nindent 4 }}
rules:
- apiGroups:
  - {{ .APIGroup }}
  resources:
{{- range $crdName := .CRDNames }}
  - {{ $crdName }}
{{- end }}
  verbs:
  - get
  - list
  - watch
