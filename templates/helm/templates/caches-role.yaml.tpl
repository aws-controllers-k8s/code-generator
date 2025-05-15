apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-namespaces-cache
rules:
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-configmaps-cache
  namespace: {{ "{{ .Release.Namespace }}" }}
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - watch