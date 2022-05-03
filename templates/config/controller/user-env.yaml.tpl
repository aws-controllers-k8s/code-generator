apiVersion: apps/v1
kind: Deployment
metadata:
  name: ack-{{.ServicePackageName}}-controller
  namespace: {{.Annotations.SuggestedNamespace}}
spec:
  template:
    spec:
      containers:
      - name: controller
        envFrom:
          - configMapRef:
              name: ack-{{.ServicePackageName}}-user-config
              optional: false
          - secretRef:
              name: ack-{{.ServicePackageName}}-user-secrets
              optional: false
