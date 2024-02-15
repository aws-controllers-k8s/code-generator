apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: ack-{{ .ControllerName }}-controller-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ack-{{ .ControllerName }}-controller
subjects:
- kind: ServiceAccount
  name: {{ .ServiceAccountName }}
  namespace: ack-system
