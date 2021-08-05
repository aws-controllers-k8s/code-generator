
---
apiVersion: rbac.authorization.k8s.io/v1
{{ if not .Values.namespacedInstallation }}
kind: ClusterRole
metadata:
  creationTimestamp: null
  name: ack-sagemaker-controller
{{ else }}
kind: Role
metadata:
  creationTimestamp: null
  name: ack-sagemaker-controller
  namespace: {{ .Release.Namespace }}
{{ end }}
