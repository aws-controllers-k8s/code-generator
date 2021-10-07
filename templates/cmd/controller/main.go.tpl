{{ template "boilerplate" }}

package main

import (
	"os"

	ackcfg "github.com/aws-controllers-k8s/runtime/pkg/config"
	ackrt "github.com/aws-controllers-k8s/runtime/pkg/runtime"
	ackrtutil "github.com/aws-controllers-k8s/runtime/pkg/util"
	ackrtwebhook "github.com/aws-controllers-k8s/runtime/pkg/webhook"
	flag "github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrlrt "sigs.k8s.io/controller-runtime"
	ctrlrtmetrics "sigs.k8s.io/controller-runtime/pkg/metrics"
	svcsdk "github.com/aws/aws-sdk-go/service/{{ .ServicePackageName }}"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	svcresource "github.com/aws-controllers-k8s/{{ .ServicePackageName }}-controller/pkg/resource"
	svctypes "github.com/aws-controllers-k8s/{{ .ServicePackageName }}-controller/apis/{{ .APIVersion }}"
	{{/* TODO(a-hilaly): import apis/* packages to register webhooks */}}
	{{ $serviceAlias := .ServicePackageName }} {{range $crdName := .SnakeCasedCRDNames }}_ "github.com/aws-controllers-k8s/{{ $serviceAlias }}-controller/pkg/resource/{{ $crdName }}"
	{{end}}
)

var (
	awsServiceAPIGroup = "{{ .APIGroup }}"
	awsServiceAlias	= "{{ .ServicePackageName }}"
	awsServiceEndpointsID = svcsdk.EndpointsID
	scheme			 = runtime.NewScheme()
	setupLog		   = ctrlrt.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	{{/* TODO(a-hilaly): register all the apis/* schemes */}}
	_ = svctypes.AddToScheme(scheme)
	_ = ackv1alpha1.AddToScheme(scheme)
}

func main() {
	var ackCfg ackcfg.Config
	ackCfg.BindFlags()
	flag.Parse()
	ackCfg.SetupLogger()

	if err := ackCfg.Validate(); err != nil {
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

	mgr, err := ctrlrt.NewManager(ctrlrt.GetConfigOrDie(), ctrlrt.Options{
		Scheme:             scheme,
		Port:               port,
		Host:               host,
		MetricsBindAddress: ackCfg.MetricsAddr,
		LeaderElection:	    ackCfg.EnableLeaderElection,
		LeaderElectionID:   awsServiceAPIGroup,
		Namespace:          ackCfg.WatchNamespace,
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
		ackrt.VersionInfo{},	// TODO: populate version info
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
