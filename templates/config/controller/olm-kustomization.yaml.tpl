{{ template "controller_kustomization" . }}

patchesJson6902:
  - target:
      group: apps
      version: v1
      kind: Deployment
      name: ack-{{ .ServicePackageName }}-controller
    patch: |-
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

patchesStrategicMerge:
  - user-env.yaml
