resources:
- ../../default
patches:
- path: role.json
  target:
    group: rbac.authorization.k8s.io
    version: v1
    kind: ClusterRole
    name: ack-{{ .ServicePackageName }}-controller
- path: role-binding.json
  target:
    group: rbac.authorization.k8s.io
    version: v1
    kind: ClusterRoleBinding
    name: ack-{{ .ServicePackageName }}-controller-rolebinding