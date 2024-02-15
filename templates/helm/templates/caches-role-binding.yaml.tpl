apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ack-namespaces-cache-{{ .ControllerName }}-controller
roleRef:
  kind: ClusterRole
  apiGroup: rbac.authorization.k8s.io
  name: ack-namespaces-cache-{{ .ControllerName }}-controller
subjects:
- kind: ServiceAccount
  name: ack-{{ .ControllerName }}-controller
  namespace: {{ "{{ .Release.Namespace }}" }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: ack-configmaps-cache-{{ .ControllerName }}-controller
  namespace: {{ "{{ .Release.Namespace }}" }}
roleRef:
  kind: Role
  apiGroup: rbac.authorization.k8s.io
  name: ack-configmaps-cache-{{ .ControllerName }}-controller
subjects:
- kind: ServiceAccount
  name: ack-{{ .ControllerName }}-controller
  namespace: {{ "{{ .Release.Namespace }}" }}