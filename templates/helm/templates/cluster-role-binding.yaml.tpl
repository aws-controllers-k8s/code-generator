apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ "{{ include \"app.fullname\" . }}" }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ack-{{ .ServiceIDClean }}-controller
subjects:
- kind: ServiceAccount
  name: {{ "{{ include \"service-account.name\" . }}" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
