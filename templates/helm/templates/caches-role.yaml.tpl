apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-clusterrole
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
  name: {{ IncludeTemplate "app.fullname" }}-role
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