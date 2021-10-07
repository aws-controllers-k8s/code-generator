apiVersion: v1
kind: Service
metadata:
  name: ack-{{ .ServiceAlias }}-metrics-service
  namespace: ack-system
spec:
  selector:
    control-plane: controller
  ports:
    - name: metricsport
      port: 8080
      targetPort: http
      protocol: TCP
  type: NodePort
