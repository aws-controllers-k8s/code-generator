{{ template "boilerplate" }}

package main

import (
	"os"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcfg "github.com/aws-controllers-k8s/runtime/pkg/config"
	ackrt "github.com/aws-controllers-k8s/runtime/pkg/runtime"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	ackrtutil "github.com/aws-controllers-k8s/runtime/pkg/util"
	ackrtwebhook "github.com/aws-controllers-k8s/runtime/pkg/webhook"
	flag "github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrlrt "sigs.k8s.io/controller-runtime"
	ctrlrtcache "sigs.k8s.io/controller-runtime/pkg/cache"
	ctrlrtmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	ctrlrtwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"

{{- /* Import the go types from service controllers whose resources are referenced in this service controller.
If these referenced types are not added to scheme, this service controller will not be able to read
resources across service controller. */ -}}
{{- $servicePackageName := .ServicePackageName }}
{{- $apiVersion := .APIVersion }}
{{- range $referencedServiceName := .ReferencedServiceNames }}
{{- if not (eq $referencedServiceName $servicePackageName) }}
	{{ $referencedServiceName }}apitypes "github.com/aws-controllers-k8s/{{ $referencedServiceName }}-controller/apis/{{ $apiVersion }}"
{{- end }}
{{- end }}

	svcresource "github.com/aws-controllers-k8s/{{ .ServicePackageName }}-controller/pkg/resource"
	svcsdk "github.com/aws/aws-sdk-go/service/{{ .ServicePackageName }}"
	svctypes "github.com/aws-controllers-k8s/{{ .ServicePackageName }}-controller/apis/{{ .APIVersion }}"

	{{/* TODO(a-hilaly): import apis/* packages to register webhooks */}}
	{{range $crdName := .SnakeCasedCRDNames }}_ "github.com/aws-controllers-k8s/{{ $servicePackageName }}-controller/pkg/resource/{{ $crdName }}"
	{{end}}
	"github.com/aws-controllers-k8s/{{ .ServicePackageName }}-controller/pkg/version"
)

var (
	awsServiceAPIGroup      = "{{ .APIGroup }}"
	awsServiceAlias	        = "{{ .ServicePackageName }}"
	awsServiceEndpointsID   = svcsdk.EndpointsID
	scheme			        = runtime.NewScheme()
	setupLog		        = ctrlrt.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	{{/* TODO(a-hilaly): register all the apis/* schemes */}}
	_ = svctypes.AddToScheme(scheme)
	_ = ackv1alpha1.AddToScheme(scheme)
{{- range $referencedServiceName := .ReferencedServiceNames }}
{{- if not (eq $referencedServiceName $servicePackageName) }}
	_ = {{ $referencedServiceName }}apitypes.AddToScheme(scheme)
{{- end }}
{{- end }}
}

func main() {
	var ackCfg ackcfg.Config
	ackCfg.BindFlags()
	flag.Parse()
	ackCfg.SetupLogger()

	managerFactories := svcresource.GetManagerFactories()
	resourceGVKs := make([]schema.GroupVersionKind, 0, len(managerFactories))
	for _, mf := range managerFactories {
		resourceGVKs = append(resourceGVKs, mf.ResourceDescriptor().GroupVersionKind())
	}

	if err := ackCfg.Validate(ackcfg.WithGVKs(resourceGVKs)); err != nil {
		setupLog.Error(
			err, "Unable to create controller manager",
			"aws.service", awsServiceAlias,
		)
		os.Exit(1)
	}

	host, port, err := ackrtutil.GetHostPort(ackCfg.WebhookServerAddr)
	if err != nil {
		setupLog.Error(
			err, "Unable to parse webhook server address.",
			"aws.service", awsServiceAlias,
		)
		os.Exit(1)
	}

	watchNamespaces := make(map[string]ctrlrtcache.Config, 0)
	namespaces, err := ackCfg.GetWatchNamespaces()
	if err != nil {
		setupLog.Error(
			err, "Unable to parse watch namespaces.",
			"aws.service", ackCfg.WatchNamespace,
		)
		os.Exit(1)
	}
	{{/* If namespaces is an empty slice, then we watch all namespaces */}}
	{{/* If namespaces is a slice with multiple elements, then we watch only those namespaces */}}
	for _, namespace := range namespaces {
		watchNamespaces[namespace] = ctrlrtcache.Config{}
	}
	mgr, err := ctrlrt.NewManager(ctrlrt.GetConfigOrDie(), ctrlrt.Options{
		Scheme: scheme,
		Cache: ctrlrtcache.Options{
			Scheme:            scheme,
			DefaultNamespaces: watchNamespaces,
		},
		WebhookServer: &ctrlrtwebhook.DefaultServer{
			Options: ctrlrtwebhook.Options{
				Port: port,
				Host: host,
			},
		},
		Metrics:                 metricsserver.Options{BindAddress: ackCfg.MetricsAddr},
		LeaderElection:          ackCfg.EnableLeaderElection,
		LeaderElectionID:        "ack-" + awsServiceAPIGroup,
		LeaderElectionNamespace: ackCfg.LeaderElectionNamespace,
	})
	if err != nil {
		setupLog.Error(
			err, "unable to create controller manager",
			"aws.service", awsServiceAlias,
		)
		os.Exit(1)
	}

	stopChan := ctrlrt.SetupSignalHandler()

	setupLog.Info(
		"initializing service controller",
		"aws.service", awsServiceAlias,
	)
	sc := ackrt.NewServiceController(
		awsServiceAlias, awsServiceAPIGroup, awsServiceEndpointsID,
		acktypes.VersionInfo{
			version.GitCommit,
			version.GitVersion,
			version.BuildDate,
		},
	).WithLogger(
		ctrlrt.Log,
	).WithResourceManagerFactories(
		svcresource.GetManagerFactories(),
	).WithPrometheusRegistry(
		ctrlrtmetrics.Registry,
	)

	if ackCfg.EnableWebhookServer {
		webhooks := ackrtwebhook.GetWebhooks()
		for _, webhook := range webhooks {
			if err := webhook.Setup(mgr); err != nil {
				setupLog.Error(
					err, "unable to register webhook "+webhook.UID(),
					"aws.service", awsServiceAlias,
				)
			}
		}
	}

	if err = sc.BindControllerManager(mgr, ackCfg); err != nil {
		setupLog.Error(
			err, "unable bind to controller manager to service controller",
			"aws.service", awsServiceAlias,
		)
		os.Exit(1)
	}

	setupLog.Info(
		"starting manager",
		"aws.service", awsServiceAlias,
	)
	if err := mgr.Start(stopChan); err != nil {
		setupLog.Error(
			err, "unable to start controller manager",
			"aws.service", awsServiceAlias,
		)
		os.Exit(1)
	}
}
