apiVersion: apps/v1
kind: Deployment
metadata:
  name: ack-{{.ServiceIDClean}}-controller
  namespace: {{.Annotations.SuggestedNamespace}}
spec:
  template:
    spec:
      containers:
      - name: controller
        envFrom:
          - configMapRef:
              name: ack-user-config
              optional: true
          - secretRef:
              name: ack-user-secrets
              optional: true