apiVersion: rbac.authorization.k8s.io/v1
{{ "{{ if not .Values.namespacedInstallation }}" }}
kind: ClusterRoleBinding
metadata:
  name: {{ "{{ include \"app.fullname\" . }}" }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ack-{{ .ServiceIDClean }}-controller
{{ "{{ else }}" }}
kind: RoleBinding
metadata:
  name: {{ "{{ include \"app.fullname\" . }}" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: ack-{{ .ServiceIDClean }}-controller
{{ "{{ end }}" }}
subjects:
- kind: ServiceAccount
  name: {{ "{{ include \"service-account.name\" . }}" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
