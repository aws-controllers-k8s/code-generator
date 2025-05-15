apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-namespaces-cache
roleRef:
  kind: ClusterRole
  apiGroup: rbac.authorization.k8s.io
  name: {{ IncludeTemplate "app.fullname" }}-namespaces-cache
subjects:
- kind: ServiceAccount
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-configmaps-cache
  namespace: {{ "{{ .Release.Namespace }}" }}
roleRef:
  kind: Role
  apiGroup: rbac.authorization.k8s.io
  name: {{ IncludeTemplate "app.fullname" }}-configmaps-cache
subjects:
- kind: ServiceAccount
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
