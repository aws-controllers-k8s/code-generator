module github.com/aws-controllers-k8s/{{ .ServiceIDClean }}-controller

go {{ .GoModule.CompilerVersion }}

require (
	github.com/aws/aws-controllers-k8s {{ .GoModule.RequiredModulesVersions.ACKCore }}
	github.com/aws/aws-sdk-go {{ .GoModule.RequiredModulesVersions.AWSSDKGo }}
	github.com/go-logr/logr v0.1.0
	github.com/google/go-cmp v0.3.1
	github.com/spf13/pflag v1.0.5
	k8s.io/api {{ .GoModule.RequiredModulesVersions.K8sAPI }}
	k8s.io/apimachinery {{ .GoModule.RequiredModulesVersions.K8sAPIMachinery }}
	k8s.io/client-go {{ .GoModule.RequiredModulesVersions.K8sClientGo }}
	sigs.k8s.io/controller-runtime {{ .GoModule.RequiredModulesVersions.K8sControllerRuntime }}
)