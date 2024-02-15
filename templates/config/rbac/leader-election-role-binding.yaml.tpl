---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: ack-system
  name: {{.ControllerName}}-leader-election-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{.ControllerName}}-leader-election-role
subjects:
- kind: ServiceAccount
  name: {{.ServiceAccountName}}
  namespace: ack-system
