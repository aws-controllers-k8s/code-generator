{{- "{{ if .Values.leaderElection.enabled }}" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{.ServicePackageName}}-leader-election-rolebinding
{{ "{{ if .Values.leaderElection.namespace }}" }}
  namespace: {{ "{{ .Values.leaderElection.namespace }}" }}
{{ "{{ else }}" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
{{ "{{ end }}" }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{.ServicePackageName}}-leader-election-role
subjects:
- kind: ServiceAccount
  name: {{ "{{ include \"service-account.name\" . }}" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
{{- "{{- end }}" }}
