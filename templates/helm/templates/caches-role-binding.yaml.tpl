apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-namespaces-cache
  labels:
    {{ IncludeTemplate "app.labels" | nindent 4 }}
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
  labels:
    {{ IncludeTemplate "app.labels" | nindent 4 }}
roleRef:
  kind: Role
  apiGroup: rbac.authorization.k8s.io
  name: {{ IncludeTemplate "app.fullname" }}-configmaps-cache
subjects:
- kind: ServiceAccount
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
