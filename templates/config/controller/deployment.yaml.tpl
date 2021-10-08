apiVersion: v1
kind: Namespace
metadata:
  labels:
    control-plane: controller
  name: ack-system
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ack-{{ .ServiceAlias }}-controller
  namespace: ack-system
  labels:
    control-plane: controller
spec:
  selector:
    matchLabels:
      control-plane: controller
  replicas: 1
  template:
    metadata:
      labels:
        control-plane: controller
    spec:
      containers:
      - command:
        - ./bin/controller
        args:
        - --aws-region
        - "$(AWS_REGION)"
        - --enable-development-logging
        - "$(ACK_ENABLE_DEVELOPMENT_LOGGING)"
        - --log-level
        - "$(ACK_LOG_LEVEL)"
        - --resource-tags
        - "$(ACK_RESOURCE_TAGS)"
        - --watch-namespace
        - "$(ACK_WATCH_NAMESPACE)"
        image: controller:latest
        name: controller
        ports:
          - name: http
            containerPort: 8080
        resources:
          limits:
            cpu: 100m
            memory: 300Mi
          requests:
            cpu: 100m
            memory: 200Mi
        env:
        - name: K8S_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
      terminationGracePeriodSeconds: 10
