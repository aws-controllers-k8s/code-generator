apiVersion: apps/v1
kind: Deployment
metadata:
  name: ack-{{.ControllerName}}-controller
  namespace: {{.Annotations.SuggestedNamespace}}
spec:
  template:
    spec:
      containers:
      - name: controller
        envFrom:
          - configMapRef:
              name: ack-{{.ControllerName}}-user-config
              optional: false
          - secretRef:
              name: ack-{{.ControllerName}}-user-secrets
              optional: true
