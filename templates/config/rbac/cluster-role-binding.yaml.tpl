apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ack-{{ .ServiceAlias }}-controller-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ack-{{ .ServiceAlias }}-controller
subjects:
- kind: ServiceAccount
  name: default
  namespace: ack-system
