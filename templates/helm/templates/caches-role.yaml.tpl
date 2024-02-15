apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: ack-namespaces-cache-{{ .ControllerName }}-controller
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
  name: ack-configmaps-cache-{{ .ControllerName }}-controller
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