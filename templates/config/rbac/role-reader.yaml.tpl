---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  creationTimestamp: null
  name: ack-{{ .ServicePackageName }}-reader
  namespace: default
rules:
- apiGroups:
  - {{ .APIGroup }}
  resources:
{{- range $crdName := .CRDNames }}
  - {{ $crdName }}
{{- end }}
  verbs:
  - get
  - list
  - watch
