apiVersion: v1
kind: Service
metadata:
  name: ack-{{ .ControllerName }}-metrics-service
  namespace: ack-system
spec:
  selector:
    app.kubernetes.io/name: ack-{{ .ControllerName }}-controller
  ports:
    - name: metricsport
      port: 8080
      targetPort: http
      protocol: TCP
  type: NodePort
