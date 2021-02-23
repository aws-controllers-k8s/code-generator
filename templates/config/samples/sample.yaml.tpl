{{- range .Samples -}}
---
apiVersion: s3.services.k8s.aws/v1alpha1
kind: {{.Kind}}
metadata:
  name: example
spec:
  {{.Spec}}
{{ end }}