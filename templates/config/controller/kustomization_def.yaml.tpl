{{- define "controller_kustomization" -}}
resources:
- deployment.yaml
- service.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
images:
- name: controller
  newName: ack-{{ .ServicePackageName }}-controller
  newTag: latest
{{end}}