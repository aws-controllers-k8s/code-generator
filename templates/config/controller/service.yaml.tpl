apiVersion: v1
kind: Service
metadata:
  name: ack-{{ .ServicePackageName }}-metrics-service
  namespace: ack-system
spec:
  selector:
    app.kubernetes.io/name: ack-{{ .ServicePackageName }}-controller
  ports:
    - name: metricsport
      port: 8080
      targetPort: http
      protocol: TCP
  type: NodePort
