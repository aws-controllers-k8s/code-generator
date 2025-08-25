{{ "{{ if eq .Values.installScope \"cluster\" }}" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-rolebinding
  labels:
    {{ IncludeTemplate "app.labels" | nindent 4 }}
roleRef:
  kind: ClusterRole
  apiGroup: rbac.authorization.k8s.io
  name: {{ IncludeTemplate "app.fullname" }}
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
{{ VarIncludeTemplate "chartVersion" "chart.name-version" }}
{{ "{{ $appVersion := .Chart.AppVersion | quote }}" }}
{{ "{{ range $namespaces }}" }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ "{{ $fullname }}" }}-{{ "{{ . }}" }}
  namespace: {{ "{{ . }}" }}
  labels:
    {{ IncludeTemplate "app.labels" | nindent 4 }}
roleRef:
  kind: Role
  apiGroup: rbac.authorization.k8s.io
  name: {{ "{{ $fullname }}" }}-{{ "{{ . }}" }}
subjects:
- kind: ServiceAccount
  name: {{ "{{ $serviceAccountName }}" }}
  namespace: {{ "{{ $releaseNamespace }}" }}
{{ "{{ end }}" }}
{{ "{{ end }}" }}