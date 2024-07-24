# Default values for ack-{{ .ControllerName }}-controller.
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
  # Number of Deployment replicas
  # This determines how many instances of the controller will be running. It's recommended
  # to enable leader election if you need to increase the number of replicas > 1
  replicas: 1
  # Which nodeSelector to set?
  # See: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodeselector
  nodeSelector:
    kubernetes.io/os: linux
  # Which tolerations to set?
  # See: https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/
  tolerations: []
  # What affinity to set?
  # See: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity
  affinity: {}
  # Which priorityClassName to set?
  # See: https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/#pod-priority
  priorityClassName: ""
  # Specifies the hostname of the Pod.
  # If not specified, the pod's hostname will be set to a system-defined value.
  hostNetwork: false
  # Set DNS policy for the pod.
  # Defaults to "ClusterFirst".
  # Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'.
  # To have DNS options set along with hostNetwork, you have to specify DNS policy
  # explicitly to 'ClusterFirstWithHostNet'.
  dnsPolicy: ClusterFirst
  extraVolumes: []
  extraVolumeMounts: []

  # Additional server container environment variables
  #
  # You specify this manually like you would a raw deployment manifest.
  # This means you can bind in environment variables from secrets.
  #
  # e.g. static environment variable:
  #  - name: DEMO_GREETING
  #    value: "Hello from the environment"
  #
  # e.g. secret environment variable:
  # - name: USERNAME
  #   valueFrom:
  #     secretKeyRef:
  #       name: mysecret
  #       key: username
  extraEnvVars: []


# If "installScope: cluster" then these labels will be applied to ClusterRole
role:
  labels: {}

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
  credentials:
    # If specified, Secret with shared credentials file to use.
    secretName: ""
    # Secret stringData key that contains the credentials
    secretKey: "credentials"
    # Profile used for AWS credentials
    profile: "default"

# log level for the controller
log:
  enable_development_logging: false
  level: info

# Set to "namespace" to install the controller in a namespaced scope, will only
# watch for object creation in the namespace. By default installScope is
# cluster wide.
installScope: cluster

# Set the value of the "namespace" to be watched by the controller
# This value is only used when the `installScope` is set to "namespace". If left empty, the default value is the release namespace for the chart.
# You can set multiple namespaces by providing a comma separated list of namespaces. e.g "namespace1,namespace2"
watchNamespace: ""

resourceTags:
  # Configures the ACK service controller to always set key/value pairs tags on
  # resources that it manages.
  - services.k8s.aws/controller-version=%CONTROLLER_SERVICE%-%CONTROLLER_VERSION%
  - services.k8s.aws/namespace=%K8S_NAMESPACE%

# Set to "retain" to keep all AWS resources intact even after the K8s resources
# have been deleted. By default, the ACK controller will delete the AWS resource
# before the K8s resource is removed.
deletionPolicy: delete

# controller reconciliation configurations
reconcile:
  # The default duration, in seconds, to wait before resyncing desired state of custom resources.
  defaultResyncPeriod: 36000 # 10 Hours
  # An object representing the reconcile resync configuration for each specific resource.
  resourceResyncPeriods: {}

  # The default number of concurrent syncs that a reconciler can perform.
  defaultMaxConcurrentSyncs: 1
  # An object representing the reconcile max concurrent syncs configuration for each specific
  # resource.
  resourceMaxConcurrentSyncs: {}

serviceAccount:
  # Specifies whether a service account should be created
  create: true
  # The name of the service account to use.
  name: {{ .ServiceAccountName }}
  annotations: {}
    # eks.amazonaws.com/role-arn: arn:aws:iam::AWS_ACCOUNT_ID:role/IAM_ROLE_NAME

# Configuration of the leader election. Required for running multiple instances of the
# controller within the same cluster.
# See https://kubernetes.io/docs/concepts/architecture/leases/#leader-election
leaderElection:
  # Enable Controller Leader Election. Set this to true to enable leader election
  # for this controller.
  enabled: false
  # Leader election can be scoped to a specific namespace. By default, the controller
  # will attempt to use the namespace of the service account mounted to the Controller
  # pod.
  namespace: ""

# Configuration for feature gates.  These are optional controller features that
# can be individually enabled ("true") or disabled ("false") by adding key/value
# pairs below.
featureGates: {}
  # featureGate1: true
  # featureGate2: false
