apiVersion: kustomize.config.k8s.io/v1
kind: Kustomization
resources:
{{- range .CRDNames }}
- {{ $.APIGroup }}_{{ . }}.yaml
{{- end }}
- services.k8s.aws_fieldexports.yaml
- services.k8s.aws_iamroleselectors.yaml
