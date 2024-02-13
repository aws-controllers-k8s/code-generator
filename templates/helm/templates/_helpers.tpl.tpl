{{ "{{/* The name of the application this chart installs */}}" }}
{{ DefineTemplate "app.name" }}
{{ "{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix \"-\" -}}" }}
{{ "{{- end -}}" }}

{{ "{{/*" }}
{{ "Create a default fully qualified app name." }}
{{ "We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec)." }}
{{ "If release name contains chart name it will be used as a full name." }}
{{ "*/}}" }}
{{ DefineTemplate "app.fullname" }}
{{ "{{- if .Values.fullnameOverride -}}" }}
{{ "{{- .Values.fullnameOverride | trunc 63 | trimSuffix \"-\" -}}" }}
{{ "{{- else -}}" }}
{{ "{{- $name := default .Chart.Name .Values.nameOverride -}}" }}
{{ "{{- if contains $name .Release.Name -}}" }}
{{ "{{- .Release.Name | trunc 63 | trimSuffix \"-\" -}}" }}
{{ "{{- else -}}" }}
{{ "{{- printf \"%s-%s\" .Release.Name $name | trunc 63 | trimSuffix \"-\" -}}" }}
{{ "{{- end -}}" }}
{{ "{{- end -}}" }}
{{ "{{- end -}}" }}

{{ "{{/* The name and version as used by the chart label */}}" }}
{{ DefineTemplate "chart.name-version" }}
{{ "{{- printf \"%s-%s\" .Chart.Name .Chart.Version | replace \"+\" \"_\" | trunc 63 | trimSuffix \"-\" -}}" }}
{{ "{{- end -}}" }}

{{ "{{/* The name of the service account to use */}}" }}
{{ DefineTemplate "service-account.name" }}
    {{ "{{ default \"default\" .Values.serviceAccount.name }}" }}
{{ "{{- end -}}" }}

{{ DefineTemplate "watch-namespace" }}
{{ "{{- if eq .Values.installScope \"namespace\" -}}" }}
{{ "{{ .Values.watchNamespace | default .Release.Namespace }}" }}
{{ "{{- end -}}" }}
{{ "{{- end -}}" }}

{{ "{{/* The mount path for the shared credentials file */}}" }}
{{ DefineTemplate "aws.credentials.secret_mount_path" }}
{{ "{{- \"/var/run/secrets/aws\" -}}" }}
{{ "{{- end -}}" }}

{{ "{{/* The path the shared credentials file is mounted */}}" }}
{{ DefineTemplate "aws.credentials.path" }}
{{ "{{- printf \"%s/%s\" (include \"aws.credentials.secret_mount_path\" .) .Values.aws.credentials.secretKey -}}" }}
{{ "{{- end -}}" }}

{{ "{{/* The rules a of ClusterRole or Role */}}" }}
{{ DefineTemplate "rbac-rules" }}
SEDREPLACERULES
{{ "{{- end }}" }}