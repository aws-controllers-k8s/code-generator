{{ "{{ if eq .Values.installScope \"cluster\" }}" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-rolebinding
  labels:
    app.kubernetes.io/name: {{ IncludeTemplate "app.name" }}
    app.kubernetes.io/instance: {{ "{{ .Release.Name }}" }}
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/version: {{ "{{ .Chart.AppVersion | quote }}" }}
    k8s-app: {{ IncludeTemplate "app.name" }}
    helm.sh/chart: {{ IncludeTemplate "chart.name-version" }}
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
    app.kubernetes.io/name: {{ "{{ $fullname }}" }}
    app.kubernetes.io/instance: {{ "{{ $.Release.Name }}" }}
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/version: {{ "{{ $appVersion }}" }}
    k8s-app: {{ "{{ $fullname }}" }}
    helm.sh/chart: {{ "{{ $chartVersion }}" }}
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