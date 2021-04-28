{{- range .Samples -}}
---
apiVersion: {{$.APIGroup}}/{{$.APIVersion}}
kind: {{.Kind}}
metadata:
  name: example
spec:
  {{.Spec}}
{{ end }}