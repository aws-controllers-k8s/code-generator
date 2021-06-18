{{ template "boilerplate" }}

package resource

import (
	ackrt "github.com/aws-controllers-k8s/runtime/pkg/runtime"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
)

// +kubebuilder:rbac:groups=services.k8s.aws,resources=adoptedresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=services.k8s.aws,resources=adoptedresources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
{{ if .GeneratorConfig.ResourceContainsSecret -}}
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
{{- end }}

var (
	reg = ackrt.NewRegistry()
)

// GetManagerFactories returns a slice of resource manager factories that are
// registered with this package
func GetManagerFactories() []acktypes.AWSResourceManagerFactory {
	return reg.GetResourceManagerFactories()
}

// RegisterManagerFactory registers a resource manager factory with the
// package's registry
func RegisterManagerFactory(f acktypes.AWSResourceManagerFactory) {
	reg.RegisterResourceManagerFactory(f)
}
