{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"context"
	"fmt"
	"time"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcfg "github.com/aws-controllers-k8s/runtime/pkg/config"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrt "github.com/aws-controllers-k8s/runtime/pkg/runtime"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	ackutil "github.com/aws-controllers-k8s/runtime/pkg/util"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	svcsdk "github.com/aws/aws-sdk-go/service/{{ .ServicePackageName }}"
	svcsdkV2{{ .ServicePackageName }} "github.com/aws/aws-sdk-go-v2/service/{{ .ServicePackageName }}"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdkapi "github.com/aws/aws-sdk-go/service/{{ .ServicePackageName }}/{{ .ServicePackageName }}iface"

	svcapitypes "github.com/aws-controllers-k8s/{{ .ControllerName }}-controller/apis/{{ .APIVersion }}"
)

var (
	_ = ackutil.InStrings
	_ = acktags.NewTags()
	_ = ackrt.MissingImageTagValue
	_ = svcapitypes.{{ .CRD.Kind }}{}
)

// +kubebuilder:rbac:groups={{ .APIGroup }},resources={{ ToLower .CRD.Plural }},verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups={{ .APIGroup }},resources={{ ToLower .CRD.Plural }}/status,verbs=get;update;patch

{{ GoCodeFindLateInitializedFieldNames .CRD "lateInitializeFieldNames" 1 }}

// resourceManager is responsible for providing a consistent way to perform
// CRUD operations in a backend AWS service API for Book custom resources.
type resourceManager struct {
	// cfg is a copy of the ackcfg.Config object passed on start of the service
	// controller
	cfg ackcfg.Config
	// log refers to the logr.Logger object handling logging for the service
	// controller
	log logr.Logger
	// metrics contains a collection of Prometheus metric objects that the
	// service controller and its reconcilers track
	metrics *ackmetrics.Metrics
	// rr is the Reconciler which can be used for various utility
	// functions such as querying for Secret values given a SecretReference
	rr acktypes.Reconciler
	// awsAccountID is the AWS account identifier that contains the resources
	// managed by this resource manager
	awsAccountID ackv1alpha1.AWSAccountID
	// The AWS Region that this resource manager targets
	awsRegion ackv1alpha1.AWSRegion
	// sess is the AWS SDK Session object used to communicate with the backend
	// AWS service API
	sess *session.Session
	// sdk is a pointer to the AWS service API interface exposed by the
	// aws-sdk-go/services/{alias}/{alias}iface package.
	sdkapi svcsdkapi.{{ .ClientInterfaceTypeName }}

	// This is for AWS-SDK-GO-V2
	config aws.Config

	// clientV2 is the AWS SDK V2 Client object used to communicate with the backend AWS Service API
	// This is for AWS-SDK-GO-V2
	clientV2 *svcsdkV2{{ .ServicePackageName }}.Client
}

// concreteResource returns a pointer to a resource from the supplied
// generic AWSResource interface
func (rm *resourceManager) concreteResource(
	res acktypes.AWSResource,
) *resource {
	// cast the generic interface into a pointer type specific to the concrete
	// implementing resource type managed by this resource manager
	return res.(*resource)
}

// ReadOne returns the currently-observed state of the supplied AWSResource in
// the backend AWS service API.
func (rm *resourceManager) ReadOne(
	ctx context.Context,
	res acktypes.AWSResource,
) (acktypes.AWSResource, error) {
	r := rm.concreteResource(res)
	if r.ko == nil {
		// Should never happen... if it does, it's buggy code.
		panic("resource manager's ReadOne() method received resource with nil CR object")
	}
	observed, err := rm.sdkFind(ctx, r)
	if err != nil {
		if observed != nil {
			return rm.onError(observed, err)
		}
		return rm.onError(r, err)
	}
	return rm.onSuccess(observed)
}

// Create attempts to create the supplied AWSResource in the backend AWS
// service API, returning an AWSResource representing the newly-created
// resource
func (rm *resourceManager) Create(
	ctx context.Context,
	res acktypes.AWSResource,
) (acktypes.AWSResource, error) {
	r := rm.concreteResource(res)
	if r.ko == nil {
		// Should never happen... if it does, it's buggy code.
		panic("resource manager's Create() method received resource with nil CR object")
	}
	created, err := rm.sdkCreate(ctx, r)
	if err != nil {
	    if created != nil {
	        return rm.onError(created, err)
	    }
		return rm.onError(r, err)
	}
	return rm.onSuccess(created)
}

// Update attempts to mutate the supplied desired AWSResource in the backend AWS
// service API, returning an AWSResource representing the newly-mutated
// resource.
// Note for specialized logic implementers can check to see how the latest
// observed resource differs from the supplied desired state. The
// higher-level reonciler determines whether or not the desired differs
// from the latest observed and decides whether to call the resource
// manager's Update method
func (rm *resourceManager) Update(
	ctx context.Context,
	resDesired acktypes.AWSResource,
	resLatest acktypes.AWSResource,
	delta *ackcompare.Delta,
) (acktypes.AWSResource, error) {
	desired := rm.concreteResource(resDesired)
	latest := rm.concreteResource(resLatest)
	if desired.ko == nil || latest.ko == nil {
		// Should never happen... if it does, it's buggy code.
		panic("resource manager's Update() method received resource with nil CR object")
	}
	updated, err := rm.sdkUpdate(ctx, desired, latest, delta)
	if err != nil {
	    if updated != nil {
	        return rm.onError(updated, err)
	    }
		return rm.onError(latest, err)
	}
	return rm.onSuccess(updated)
}

// Delete attempts to destroy the supplied AWSResource in the backend AWS
// service API, returning an AWSResource representing the
// resource being deleted (if delete is asynchronous and takes time)
func (rm *resourceManager) Delete(
	ctx context.Context,
	res acktypes.AWSResource,
) (acktypes.AWSResource, error) {
	r := rm.concreteResource(res)
	if r.ko == nil {
		// Should never happen... if it does, it's buggy code.
		panic("resource manager's Update() method received resource with nil CR object")
	}
	observed, err := rm.sdkDelete(ctx, r)
	if err != nil {
		if observed != nil {
			return rm.onError(observed, err)
		}
		return rm.onError(r, err)
	}

	return rm.onSuccess(observed)
}

// ARNFromName returns an AWS Resource Name from a given string name. This
// is useful for constructing ARNs for APIs that require ARNs in their
// GetAttributes operations but all we have (for new CRs at least) is a
// name for the resource
func (rm *resourceManager) ARNFromName(name string) string {
	return fmt.Sprintf(
		"arn:aws:{{ .ControllerName }}:%s:%s:%s",
		rm.awsRegion,
		rm.awsAccountID,
		name,
	)
}

// LateInitialize returns an acktypes.AWSResource after setting the late initialized
// fields from the readOne call. This method will initialize the optional fields
// which were not provided by the k8s user but were defaulted by the AWS service.
// If there are no such fields to be initialized, the returned object is similar to
// object passed in the parameter.
func (rm *resourceManager) LateInitialize(
	ctx context.Context,
	latest acktypes.AWSResource,
) (acktypes.AWSResource, error) {
	rlog := ackrtlog.FromContext(ctx)
	// If there are no fields to late initialize, do nothing
	if len(lateInitializeFieldNames) == 0 {
		rlog.Debug("no late initialization required.")
		return latest, nil
	}
	latestCopy := latest.DeepCopy()
	lateInitConditionReason := ""
	lateInitConditionMessage := ""
{{- if $hookCode := Hook .CRD "late_initialize_pre_read_one" }}
{{ $hookCode }}
{{- end }}
	observed, err := rm.ReadOne(ctx, latestCopy)
	if err != nil {
		lateInitConditionMessage = "Unable to complete Read operation required for late initialization"
		lateInitConditionReason = "Late Initialization Failure"
		ackcondition.SetLateInitialized(latestCopy, corev1.ConditionFalse, &lateInitConditionMessage, &lateInitConditionReason)
		ackcondition.SetSynced(latestCopy, corev1.ConditionFalse, nil, nil)
		return latestCopy, err
	}
{{- if $hookCode := Hook .CRD "late_initialize_post_read_one" }}
{{ $hookCode }}
{{- end }}
	lateInitializedRes := rm.lateInitializeFromReadOneOutput(observed, latestCopy)
	incompleteInitialization := rm.incompleteLateInitialization(lateInitializedRes)
	if incompleteInitialization {
		// Add the condition with LateInitialized=False
		lateInitConditionMessage = "Late initialization did not complete, requeuing with delay of 5 seconds"
		lateInitConditionReason = "Delayed Late Initialization"
		ackcondition.SetLateInitialized(lateInitializedRes, corev1.ConditionFalse, &lateInitConditionMessage, &lateInitConditionReason)
		ackcondition.SetSynced(lateInitializedRes, corev1.ConditionFalse, nil, nil)
		return lateInitializedRes, ackrequeue.NeededAfter(nil, time.Duration(5)*time.Second)
	}
	// Set LateInitialized condition to True
	lateInitConditionMessage = "Late initialization successful"
	lateInitConditionReason = "Late initialization successful"
	ackcondition.SetLateInitialized(lateInitializedRes, corev1.ConditionTrue, &lateInitConditionMessage, &lateInitConditionReason)
	return lateInitializedRes, nil
}

// incompleteLateInitialization return true if there are fields which were supposed to be
// late initialized but are not. If all the fields are late initialized, false is returned
func (rm *resourceManager) incompleteLateInitialization(
	res acktypes.AWSResource,
) bool {
{{ GoCodeIncompleteLateInitialization .CRD "res" 1 }}
}

// lateInitializeFromReadOneOutput late initializes the 'latest' resource from the 'observed'
// resource and returns 'latest' resource
func (rm *resourceManager) lateInitializeFromReadOneOutput(
	observed acktypes.AWSResource,
	latest acktypes.AWSResource,
) acktypes.AWSResource {
{{ GoCodeLateInitializeFromReadOne .CRD "observed" "latest" 1 }}
}

// IsSynced returns true if the resource is synced.
func (rm *resourceManager) IsSynced(ctx context.Context, res acktypes.AWSResource) (bool, error) {
	r := rm.concreteResource(res)
	if r.ko == nil {
		// Should never happen... if it does, it's buggy code.
		panic("resource manager's IsSynced() method received resource with nil CR object")
	}
{{ GoCodeIsSynced .CRD "r.ko" 1}}
	return true, nil
}

// EnsureTags ensures that tags are present inside the AWSResource.
// If the AWSResource does not have any existing resource tags, the 'tags'
// field is initialized and the controller tags are added.
// If the AWSResource has existing resource tags, then controller tags are
// added to the existing resource tags without overriding them.
// If the AWSResource does not support tags, only then the controller tags
// will not be added to the AWSResource.
func (rm *resourceManager) EnsureTags(
    ctx context.Context,
    res acktypes.AWSResource,
    md acktypes.ServiceControllerMetadata,
) error {
{{- if $hookCode := Hook .CRD "ensure_tags" }}
{{ $hookCode }}
{{ else }}
{{ $tagField := .CRD.GetTagField -}}
{{ if $tagField -}}
{{ $tagFieldShapeType := $tagField.ShapeRef.Shape.Type -}}
{{ $tagFieldGoType := $tagField.GoType -}}
{{ if eq "list" $tagFieldShapeType -}}
{{ $tagFieldGoType = (print "[]*svcapitypes." $tagField.GoTypeElem) -}}
{{ end -}}
	r := rm.concreteResource(res)
	if r.ko == nil {
		// Should never happen... if it does, it's buggy code.
		panic("resource manager's EnsureTags method received resource with nil CR object")
	}
	defaultTags := ackrt.GetDefaultTags(&rm.cfg, r.ko, md)
	var existingTags {{ $tagFieldGoType }}
{{ $nilCheck := CheckNilFieldPath $tagField "r.ko.Spec" -}}
{{ if not (eq $nilCheck "") -}}
    if {{ $nilCheck }} {
        existingTags = nil
    } else {
        existingTags = r.ko.Spec.{{ $tagField.Path }}
    }
{{ else -}}
    existingTags = r.ko.Spec.{{ $tagField.Path }}
{{ end -}}
	resourceTags := ToACKTags(existingTags)
	tags := acktags.Merge(resourceTags, defaultTags)
{{ GoCodeInitializeNestedStructField .CRD "r.ko" $tagField "svcapitypes" 1 -}}
	r.ko.Spec.{{ $tagField.Path }} = FromACKTags(tags)
{{- end }}
    return nil
{{- end }}
}

// newResourceManager returns a new struct implementing
// acktypes.AWSResourceManager
// This is for AWS-SDK-GO-V2 - Created newResourceManager With AWS sdk-Go-ClientV2
func newResourceManager(
	cfg ackcfg.Config,
	log logr.Logger,
	metrics *ackmetrics.Metrics,
	rr acktypes.Reconciler,
	sess *session.Session,
	id ackv1alpha1.AWSAccountID,
	region ackv1alpha1.AWSRegion,
	config aws.Config,
	clientV2 *svcsdkV2{{ .ServicePackageName }}.Client,
) (*resourceManager, error) {
	return &resourceManager{
		cfg: cfg,
		log: log,
		metrics: metrics,
		rr: rr,
		awsAccountID: id,
		awsRegion: region,
		sess:		 sess,
		sdkapi:	   svcsdk.New(sess),
		config: 	config,
		clientV2:	clientV2, 
	}, nil
}

// onError updates resource conditions and returns updated resource
// it returns nil if no condition is updated.
func (rm *resourceManager) onError(
	r *resource,
	err error,
) (acktypes.AWSResource, error) {
	if r == nil {
		return nil, err
	}
	r1, updated := rm.updateConditions(r, false, err)
	if !updated {
		return r, err
	}
	for _, condition := range r1.Conditions() {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal &&
			condition.Status == corev1.ConditionTrue {
			// resource is in Terminal condition
			// return Terminal error
			return r1, ackerr.Terminal
		}
	}
	return r1, err
}

// onSuccess updates resource conditions and returns updated resource
// it returns the supplied resource if no condition is updated.
func (rm *resourceManager) onSuccess(
	r *resource,
) (acktypes.AWSResource, error) {
	if r == nil  {
		return nil, nil
	}
	r1, updated := rm.updateConditions(r, true, nil)
	if !updated {
		return r, nil
	}
	return r1, nil
}