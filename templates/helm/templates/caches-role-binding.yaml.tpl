apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ IncludeTemplate "app.fullname" }}
roleRef:
  kind: ClusterRole
  apiGroup: rbac.authorization.k8s.io
  name: {{ IncludeTemplate "app.fullname" }}
subjects:
- kind: ServiceAccount
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ IncludeTemplate "app.fullname" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
roleRef:
  kind: Role
  apiGroup: rbac.authorization.k8s.io
  name: {{ IncludeTemplate "app.fullname" }}
subjects:
- kind: ServiceAccount
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
