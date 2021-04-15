apiVersion: v1
name: ack-{{ .ServiceIDClean }}-controller
description: A Helm chart for the ACK service controller for {{ .ServiceIDClean }}
version: {{ .ReleaseVersion }}
appVersion: {{ .ReleaseVersion }}
home: https://github.com/aws-controllers-k8s/{{ .ServiceIDClean }}-controller
icon: https://raw.githubusercontent.com/aws/eks-charts/master/docs/logo/aws.png
sources:
  - https://github.com/aws-controllers-k8s/{{ .ServiceIDClean }}-controller
maintainers:
  - name: ACK Admins
    url: https://github.com/orgs/aws-controllers-k8s/teams/ack-admin
  - name: {{ .ServiceIDClean }} Admins
    url: https://github.com/orgs/aws-controllers-k8s/teams/{{ .ServiceIDClean }}-maintainer
keywords:
  - aws
  - kubernetes
  - {{ .ServiceIDClean }}
