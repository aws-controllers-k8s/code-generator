{{- "{{ if .Values.leaderElection.enabled }}" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-leaderelection
{{ "{{ if .Values.leaderElection.namespace }}" }}
  namespace: {{ "{{ .Values.leaderElection.namespace }}" }}
{{ "{{ else }}" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
{{ "{{ end }}" }}
rules:
- apiGroups:
  - coordination.k8s.io
  resources:
  - leases
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
{{- "{{- end }}" }}
