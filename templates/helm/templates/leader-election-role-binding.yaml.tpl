{{- "{{ if .Values.leaderElection.enabled }}" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-leader-election-rolebinding
{{ "{{ if .Values.leaderElection.namespace }}" }}
  namespace: {{ "{{ .Values.leaderElection.namespace }}" }}
{{ "{{ else }}" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
{{ "{{ end }}" }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ IncludeTemplate "app.fullname" }}-leader-election-role
subjects:
- kind: ServiceAccount
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
{{- "{{- end }}" }}
