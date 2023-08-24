{{ template "controller_kustomization" . }}

patches:
- patch: |-
    - op: replace
      path: '/spec/template/spec/containers/0/env'
      value: []
    - op: add
      path: '/spec/template/spec/containers/0/env/0'
      value:
        name: ACK_SYSTEM_NAMESPACE
        valueFrom:
          fieldRef:
            fieldPath: metadata.namespace
  target:
    group: apps
    kind: Deployment
    name: ack-{{ .ServicePackageName }}-controller
    version: v1
- path: user-env.yaml
