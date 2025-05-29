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
- apiGroups:
  - ""
  resources:
  - namespaces
  verbs:
  - get
  - list
  - watch
---
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