---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  namespace: ack-system
  name: {{.ServicePackageName}}-leader-election-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{.ServicePackageName}}-leader-election-role
subjects:
- kind: ServiceAccount
  name: {{.ServiceAccountName}}
  namespace: ack-system
