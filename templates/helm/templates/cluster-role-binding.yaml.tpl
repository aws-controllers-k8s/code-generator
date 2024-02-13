{{ "{{ if eq .Values.installScope \"cluster\" }}" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ IncludeTemplate "app.fullname" }}
roleRef:
  kind: ClusterRole
  apiGroup: rbac.authorization.k8s.io
  name: ack-{{ .ServicePackageName }}-controller
subjects:
- kind: ServiceAccount
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
{{ "{{ else if eq .Values.installScope \"namespace\" }}" }}
{{ VarIncludeTemplate "wn" "watch-namespace" }}
{{ "{{ $namespaces := split \",\" $wn }}" }}
{{ VarIncludeTemplate "fullname" "app.fullname" }}
{{ "{{ $releaseNamespace := .Release.Namespace }}" }}
{{ VarIncludeTemplate "serviceAccountName" "service-account.name" }}
{{ "{{ range $namespaces }}" }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ "{{ $fullname }}" }}
  namespace: {{ "{{ . }}" }}
roleRef:
  kind: Role
  apiGroup: rbac.authorization.k8s.io
  name: ack-{{ .ServicePackageName }}-controller
subjects:
- kind: ServiceAccount
  name: {{ "{{ $serviceAccountName }}" }}
  namespace: {{ "{{ $releaseNamespace }}" }}
{{ "{{ end }}" }}
{{ "{{ end }}" }}