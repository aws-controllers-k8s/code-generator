# Default values for ack-{{ .ServicePackageName }}-controller.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

image:
  repository: {{ .ImageRepository }}
  tag: {{ .ReleaseVersion }}
  pullPolicy: IfNotPresent
  pullSecrets: []

nameOverride: ""
fullnameOverride: ""

deployment:
  annotations: {}
  labels: {}
  containerPort: 8080
  # Which nodeSelector to set?
  # See: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector
  nodeSelector:
    kubernetes.io/os: linux
  # Which tolerations to set?
  # See: https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/
  tolerations: {}
  # What affinity to set?
  # See: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity
  affinity: {}
  # Which priorityClassName to set?
  # See: https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/#pod-priority
  priorityClassName: ""
  
metrics:
  service:
    # Set to true to automatically create a Kubernetes Service resource for the
    # Prometheus metrics server endpoint in controller
    create: false
    # Which Type to use for the Kubernetes Service?
    # See: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
    type: "ClusterIP"

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
  endpoint_url: ""

# log level for the controller
log:
  enable_development_logging: false
  level: info

# Set to "namespace" to install the controller in a namespaced scope, will only
# watch for object creation in the namespace. By default installScope is
# cluster wide.
installScope: cluster

resourceTags:
  # Configures the ACK service controller to always set key/value pairs tags on
  # resources that it manages.
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
