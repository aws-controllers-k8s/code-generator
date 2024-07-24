apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ IncludeTemplate "app.fullname" }}
  namespace: {{ "{{ .Release.Namespace }}" }}
  labels:
    app.kubernetes.io/name: {{ IncludeTemplate "app.name" }}
    app.kubernetes.io/instance: {{ "{{ .Release.Name }}" }}
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/version: {{ "{{ .Chart.AppVersion | quote }}" }}
    k8s-app: {{ IncludeTemplate "app.name" }}
    helm.sh/chart: {{ IncludeTemplate "chart.name-version" }}
spec:
  replicas: {{ "{{ .Values.deployment.replicas }}" }}
  selector:
    matchLabels:
      app.kubernetes.io/name: {{ IncludeTemplate "app.name" }}
      app.kubernetes.io/instance: {{ "{{ .Release.Name }}" }}
  template:
    metadata:
{{ "{{- if .Values.deployment.annotations }}" }}
      annotations:
      {{ "{{- range $key, $value := .Values.deployment.annotations }}" }}
        {{ "{{ $key }}: {{ $value | quote }}" }}
      {{ "{{- end }}" }}
{{ "{{- end }}" }}
      labels:
        app.kubernetes.io/name: {{ IncludeTemplate "app.name" }}
        app.kubernetes.io/instance: {{ "{{ .Release.Name }}" }}
        app.kubernetes.io/managed-by: Helm
        k8s-app: {{ IncludeTemplate "app.name" }}
{{ "{{- range $key, $value := .Values.deployment.labels }}" }}
        {{ "{{ $key }}: {{ $value | quote }}" }}
{{ "{{- end }}" }}
    spec:
      serviceAccountName: {{ IncludeTemplate "service-account.name" }}
      {{ "{{- if .Values.image.pullSecrets }}" }}
      imagePullSecrets:
      {{ "{{- range .Values.image.pullSecrets }}" }}
        - name: {{ "{{ . }}" }}
      {{ "{{- end }}" }}
      {{ "{{- end }}" }}
      containers:
      - command:
        - ./bin/controller
        args:
        - --aws-region
        - "$(AWS_REGION)"
        - --aws-endpoint-url
        - "$(AWS_ENDPOINT_URL)"
{{ "{{- if .Values.log.enable_development_logging }}" }}
        - --enable-development-logging
{{ "{{- end }}" }}
        - --log-level
        - "$(ACK_LOG_LEVEL)"
        - --resource-tags
        - "$(ACK_RESOURCE_TAGS)"
        - --watch-namespace
        - "$(ACK_WATCH_NAMESPACE)"
        - --deletion-policy
        - "$(DELETION_POLICY)"
{{ "{{- if .Values.leaderElection.enabled }}" }}
        - --enable-leader-election
        - --leader-election-namespace
        - "$(LEADER_ELECTION_NAMESPACE)"
{{ "{{- end }}" }}
{{ "{{- if gt (int .Values.reconcile.defaultResyncPeriod) 0 }}" }}
        - --reconcile-default-resync-seconds
        - "$(RECONCILE_DEFAULT_RESYNC_SECONDS)"
{{ "{{- end }}" }}
{{ "{{- range $key, $value := .Values.reconcile.resourceResyncPeriods }}" }}
        - --reconcile-resource-resync-seconds
        - {{ "\"$(RECONCILE_RESOURCE_RESYNC_SECONDS_{{ $key | upper }})\"" }}
{{ "{{- end }}" }}
{{ "{{- if gt (int .Values.reconcile.defaultMaxConcurrentSyncs) 0 }}" }}
        - --reconcile-default-max-concurrent-syncs
        - "$(RECONCILE_DEFAULT_MAX_CONCURRENT_SYNCS)"
{{ "{{- end }}" }}
{{ "{{- range $key, $value := .Values.reconcile.resourceMaxConcurrentSyncs }}" }}
        - --reconcile-resource-max-concurrent-syncs
        - {{ "\"$(RECONCILE_RESOURCE_MAX_CONCURRENT_SYNCS_{{ $key | upper }})\"" }}
{{ "{{- end }}" }}
{{ "{{- if .Values.featureGates}}" }}
        - --feature-gates
        - "$(FEATURE_GATES)"
{{ "{{- end }}" }}
        image: {{ "{{ .Values.image.repository }}:{{ .Values.image.tag }}" }}
        imagePullPolicy: {{ "{{ .Values.image.pullPolicy }}" }}
        name: controller
        ports:
          - name: http
            containerPort: {{ "{{ .Values.deployment.containerPort }}" }}
        resources:
          {{ "{{- toYaml .Values.resources | nindent 10 }}" }}
        env:
        - name: ACK_SYSTEM_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        - name: AWS_REGION
          value: {{ "{{ .Values.aws.region }}" }}
        - name: AWS_ENDPOINT_URL
          value: {{ "{{ .Values.aws.endpoint_url | quote }}" }}
        - name: ACK_WATCH_NAMESPACE
          value: {{ IncludeTemplate "watch-namespace" }}
        - name: DELETION_POLICY
          value: {{ "{{ .Values.deletionPolicy }}" }}
        - name: LEADER_ELECTION_NAMESPACE
          value: {{ "{{ .Values.leaderElection.namespace | quote }}" }}
        - name: ACK_LOG_LEVEL
          value: {{ "{{ .Values.log.level | quote }}" }}
        - name: ACK_RESOURCE_TAGS
          value: {{ "{{ join \",\" .Values.resourceTags | quote }}" }}
{{ "{{- if gt (int .Values.reconcile.defaultResyncPeriod) 0 }}" }}
        - name: RECONCILE_DEFAULT_RESYNC_SECONDS
          value: {{ "{{ .Values.reconcile.defaultResyncPeriod | quote }}" }}
{{ "{{- end }}" }}
{{ "{{- range $key, $value := .Values.reconcile.resourceResyncPeriods }}" }}
        - name: {{ "RECONCILE_RESOURCE_RESYNC_SECONDS_{{ $key | upper }}" }}
          value: {{ "{{ $key }}={{ $value }}" }}
{{ "{{- end }}" }}
{{ "{{- if gt (int .Values.reconcile.defaultMaxConcurrentSyncs) 0 }}" }}
        - name: RECONCILE_DEFAULT_MAX_CONCURRENT_SYNCS
          value: {{ "{{ .Values.reconcile.defaultMaxConcurrentSyncs | quote }}" }}
{{ "{{- end }}" }}
{{ "{{- range $key, $value := .Values.reconcile.resourceMaxConcurrentSyncs }}" }}
        - name: {{ "RECONCILE_RESOURCE_MAX_CONCURRENT_SYNCS_{{ $key | upper }}" }}
          value: {{ "{{ $key }}={{ $value }}" }}
{{ "{{- end }}" }}
{{ "{{- if .Values.featureGates}}" }}
        - name: FEATURE_GATES
          value: {{ IncludeTemplate "feature-gates" }}
{{ "{{- end }}" }}
        {{ "{{- if .Values.aws.credentials.secretName }}" }}
        - name: AWS_SHARED_CREDENTIALS_FILE
          value: {{ IncludeTemplate "aws.credentials.path" }}
        - name: AWS_PROFILE
          value: {{ "{{ .Values.aws.credentials.profile }}" }}
        {{ "{{- end }}" }}
        {{ "{{- if .Values.deployment.extraEnvVars -}}" }}
          {{ "{{ toYaml .Values.deployment.extraEnvVars | nindent 8 }}" }}
        {{ "{{- end }}" }}
        volumeMounts:
        {{ "{{- if .Values.aws.credentials.secretName }}" }}
          - name: {{ "{{ .Values.aws.credentials.secretName }}" }}
            mountPath: {{ IncludeTemplate "aws.credentials.secret_mount_path" }}
            readOnly: true
        {{ "{{- end }}" }}
        {{ "{{- if .Values.deployment.extraVolumeMounts -}}" }}
          {{ "{{ toYaml .Values.deployment.extraVolumeMounts | nindent 10 }}" }}
        {{ "{{- end }}" }}
        securityContext:
          allowPrivilegeEscalation: false
          privileged: false
          runAsNonRoot: true
          capabilities:
            drop:
              - ALL
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8081
          initialDelaySeconds: 15
          periodSeconds: 20
        readinessProbe:
          httpGet:
            path: /readyz
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 10
      securityContext:
        seccompProfile:
          type: RuntimeDefault
      terminationGracePeriodSeconds: 10
      nodeSelector: {{ "{{ toYaml .Values.deployment.nodeSelector | nindent 8 }}" }}
      {{ "{{ if .Values.deployment.tolerations -}}" }}
      tolerations: {{ "{{ toYaml .Values.deployment.tolerations | nindent 8 }}" }}
      {{ "{{ end -}}" }}
      {{ "{{ if .Values.deployment.affinity -}}" }}
      affinity: {{ "{{ toYaml .Values.deployment.affinity | nindent 8 }}" }}
      {{ "{{ end -}}" }}
      {{ "{{ if .Values.deployment.priorityClassName -}}" }}
      priorityClassName: {{ "{{ .Values.deployment.priorityClassName }}" }}
      {{ "{{ end -}}" }}
      hostIPC: false
      hostPID: false
      hostNetwork: {{ "{{ .Values.deployment.hostNetwork }}" }}
      dnsPolicy: {{ "{{ .Values.deployment.dnsPolicy }}" }}
      volumes:
      {{ "{{- if .Values.aws.credentials.secretName }}" }}
        - name: {{ "{{ .Values.aws.credentials.secretName }}" }}
          secret:
            secretName: {{ "{{ .Values.aws.credentials.secretName }}" }}
      {{ "{{- end }}" }}
{{ "{{- if .Values.deployment.extraVolumes }}" }}
{{ "{{ toYaml .Values.deployment.extraVolumes | indent 8}}" }}
{{ "{{- end }}" }}
