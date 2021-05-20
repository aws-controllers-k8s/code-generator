# Default values for ack-{{ .ServiceIDClean }}-controller.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

image:
  repository: {{ .ImageRepository }}
  tag: {{ .ServiceIDClean }}-{{ .ReleaseVersion }}
  pullPolicy: IfNotPresent
  pullSecrets: []

nameOverride: ""
fullnameOverride: ""

deployment:
  annotations: {}
  labels: {}
  containerPort: 8080

resources:
  requests:
    memory: "64Mi"
    cpu: "50m"
  limits:
    memory: "128Mi"
    cpu: "100m"

aws:
  # If specified, use the AWS region for AWS API calls
  region: ""
  account_id: ""

# log level for the controller
log:
  enable_development_logging: false
  level: info

# If specified, the service controller will watch for object creation only in the provided namespace
watchNamespace: ""

resourceTags:
  # Configures the ACK service controller to always set key/value pairs tags on resources that it manages.
  - services.k8s.aws/managed=true
  - services.k8s.aws/created=%UTCNOW%
  - services.k8s.aws/namespace=%KUBERNETES_NAMESPACE%

serviceAccount:
  # Specifies whether a service account should be created
  create: true
  # The name of the service account to use.
  name: {{ .ServiceAccountName }}
  annotations: {}
    # eks.amazonaws.com/role-arn: arn:aws:iam::AWS_ACCOUNT_ID:role/IAM_ROLE_NAME
