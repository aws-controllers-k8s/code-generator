{{ template "controller_kustomization" . }}

patchesStrategicMerge:
- user-env.yaml