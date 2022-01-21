{{ "{{ .Chart.Name }}" }} has been installed. Check its status by running:
  kubectl --namespace {{ "{{ .Release.Namespace }}" }} get pods -l "app.kubernetes.io/instance={{ "{{ .Release.Name }}" }}"

You are now able to create {{ .Metadata.Service.FullName }} ({{ .Metadata.Service.ShortName }})
resources in your cluster!

The controller is running in "{{ "{{ .Values.installScope }}" }}" mode.
The controller is configured to run in the region: "{{ "{{ .Values.aws.region }}" }}"

Visit https://aws-controllers-k8s.github.io/community/reference/ for an API 
reference of all the resources that can be created using this controller.

For more information on the AWS Controller for Kubernetes (ACK) project, visit:
https://aws-controllers-k8s.github.io/community/