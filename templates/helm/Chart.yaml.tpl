apiVersion: v1
name: {{ .ServiceAlias }}-chart
description: A Helm chart for the ACK service controller for {{ .Metadata.Service.FullName }} ({{ .Metadata.Service.ShortName }})
version: {{ .ReleaseVersion }}
appVersion: {{ .ReleaseVersion }}
home: https://github.com/aws-controllers-k8s/{{ .ServiceAlias }}-controller
icon: https://raw.githubusercontent.com/aws/eks-charts/master/docs/logo/aws.png
sources:
  - https://github.com/aws-controllers-k8s/{{ .ServiceAlias }}-controller
maintainers:
  - name: ACK Admins
    url: https://github.com/orgs/aws-controllers-k8s/teams/ack-admin
  - name: {{ .Metadata.Service.ShortName }} Admins
    url: https://github.com/orgs/aws-controllers-k8s/teams/{{ .ServiceAlias }}-maintainer
keywords:
  - aws
  - kubernetes
  - {{ .ServiceAlias }}
