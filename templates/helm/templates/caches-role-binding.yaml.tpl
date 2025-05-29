apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ IncludeTemplate "app.fullname" }}-namespace-caches
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
  name: {{ IncludeTemplate "app.fullname" }}-namespace-caches
subjects:
- kind: ServiceAccount
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
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
roleRef:
  kind: Role
  apiGroup: rbac.authorization.k8s.io
  name: {{ IncludeTemplate "app.fullname" }}-configmaps-cache
subjects:
- kind: ServiceAccount
  name: {{ IncludeTemplate "service-account.name" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
