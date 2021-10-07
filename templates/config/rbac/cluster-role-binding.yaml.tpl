apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ack-{{ .ServicePackageName }}-controller-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ack-{{ .ServicePackageName }}-controller
subjects:
- kind: ServiceAccount
  name: default
  namespace: ack-system
