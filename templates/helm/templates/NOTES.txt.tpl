{{ "{{ .Chart.Name }}" }} has been installed.
This chart deploys "{{ .ImageRepository }}:{{ .ReleaseVersion }}".

Check its status by running:
  kubectl --namespace {{ "{{ .Release.Namespace }}" }} get pods -l "app.kubernetes.io/instance={{ "{{ .Release.Name }}" }}"

You are now able to create {{ .Metadata.Service.FullName }} ({{ .Metadata.Service.ShortName }}) resources!

The controller is running in "{{ "{{ .Values.installScope }}" }}" mode.
The controller is configured to manage AWS resources in region: "{{ "{{ .Values.aws.region }}" }}"

Visit https://aws-controllers-k8s.github.io/community/reference/ for an API 
reference of all the resources that can be created using this controller.

For more information on the AWS Controllers for Kubernetes (ACK) project, visit:
https://aws-controllers-k8s.github.io/community/
