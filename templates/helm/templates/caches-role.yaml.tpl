{{ VarIncludeTemplate "featuregates" "feature-gates" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-namespaces-cache
  labels:
    app.kubernetes.io/name: {{ IncludeTemplate "app.name" }}
    app.kubernetes.io/instance: {{ "{{ .Release.Name }}" }}
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/version: {{ "{{ .Chart.AppVersion | quote }}" }}
    k8s-app: {{ IncludeTemplate "app.name" }}
    helm.sh/chart: {{ IncludeTemplate "chart.name-version" }}
rules:
{{ "{{ if contains \"IAMRoleSelector=true\" $featuregates }}" }}
- apiGroups:
  - services.k8s.aws
  resources:
  - iamroleselectors
  verbs:
  - get
  - list
  - watch
{{ "{{ end }}" }}
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - get
  - list
  - watch
---
{{ "{{ if eq .Values.enableCARM true }}" }}
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-configmaps-cache
  namespace: {{ "{{ .Release.Namespace }}" }}
  labels:
    app.kubernetes.io/name: {{ IncludeTemplate "app.name" }}
    app.kubernetes.io/instance: {{ "{{ .Release.Name }}" }}
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/version: {{ "{{ .Chart.AppVersion | quote }}" }}
    k8s-app: {{ IncludeTemplate "app.name" }}
    helm.sh/chart: {{ IncludeTemplate "chart.name-version" }}
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - watch
{{ "{{ end }}" }}