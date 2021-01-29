package config

// GoModule contains controller Go module informations
type GoModule struct {
	// Go compiler version. Defaults to 1.15
	CompilerVersion string `json:"compiler_version"`
	// Required modules versions
	RequiredModulesVersions RequiredModulesVersions `json:"required_modules_versions"`
}

// RequiredModulesVersions contains required modules versions
type RequiredModulesVersions struct {
	// Version of github.com/aws/aws-sdk-go. Defaults to v1.35.5
	AWSSDKGo string `json:"aws_sdk_go"`
	// Version of github.com/aws/aws-controllers-k8s. Defaults to v0.0.2
	ACKCore string `json:"ack_core"`
	// Version of k8s.io/api. Defaults to v0.18.2
	K8sAPI string `json:"k8s_api"`
	// Version of k8s.io/apimachinery. Defaults to v0.18.6
	K8sAPIMachinery string `json:"k8s_api_machinery"`
	// Version of k8s.io/client-go. Defaults to v0.18.2
	K8sClientGo string `json:"k8s_client_go"`
	// Version of sigs.k8s.io/controller-runtime version. Defaults to v0.6.0
	K8sControllerRuntime string `json:"k8s_controller_runtime"`
}
