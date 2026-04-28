{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"context"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcfg "github.com/aws-controllers-k8s/runtime/pkg/config"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/{{ .ServicePackageName }}"
	"github.com/go-logr/logr"

	svcapitypes "github.com/aws-controllers-k8s/{{ .ControllerName }}-controller/apis/{{ .APIVersion }}"
)

// Hack to avoid import errors during build...
var (
	_ = &ackv1alpha1.AWSIdentifiers{}
)

type resource struct {
	ko *svcapitypes.{{ .CRD.Kind }}
}

type resourceManager struct {
	cfg          ackcfg.Config
	clientcfg    aws.Config
	log          logr.Logger
	metrics      *ackmetrics.Metrics
	rr           acktypes.Reconciler
	awsAccountID ackv1alpha1.AWSAccountID
	awsRegion    ackv1alpha1.AWSRegion
	sdkapi       *svcsdk.Client
}

// delta holds the result of comparing desired vs. latest sub-resource items.
type delta struct {
	toCreate []*resource
	toUpdate []*resource
	toDelete []*resource
}

// computeDelta performs a key-based diff between the desired and latest
// sub-resource slices.
func computeDelta(desired, latest []resource) delta {
	d := delta{}
	latestMap := make(map[string]*resource)
	for i := range latest {
		latestMap[key(&latest[i])] = &latest[i]
	}

	desiredMap := make(map[string]bool)
	for i := range desired {
		k := key(&desired[i])
		desiredMap[k] = true
		if lat, exists := latestMap[k]; !exists {
			d.toCreate = append(d.toCreate, &desired[i])
		} else {
			rd := newResourceDelta(&desired[i], lat)
			if len(rd.Differences) > 0 {
				d.toUpdate = append(d.toUpdate, &desired[i])
			}
		}
	}

	for k, lat := range latestMap {
		if !desiredMap[k] {
			d.toDelete = append(d.toDelete, lat)
		}
	}
	return d
}

// NewManager creates a resourceManager using the provided SDK client and
// metrics recorder.
func NewManager(
	sdkapi *svcsdk.Client,
	metrics *ackmetrics.Metrics,
) *resourceManager {
	return &resourceManager{
		sdkapi:  sdkapi,
		metrics: metrics,
	}
}

// Sync converts the desired and latest parent resources into internal
// resource slices, computes the diff, and delegates to sync.
func (rm *resourceManager) Sync(
	ctx context.Context,
	desired any,
	latest any,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.Sync")
	defer func() {
		exit(err)
	}()
	desiredResources := convertToResources(desired)
	latestResources := convertToResources(latest)
	d := computeDelta(desiredResources, latestResources)
	return rm.sync(ctx, d)
}

// convertToResources casts the parent into its concrete type and delegates
// to convertFromParent. Returns nil when the parent has no value at the
// sub-resource's source field.
func convertToResources(parent any) []resource {
	if parent == nil {
		return nil
	}
	v, ok := parent.(*svcapitypes.{{ ParentKind .CRD }})
	if !ok || v == nil {
		return nil
	}
	if v.{{ SubResFieldPath .CRD }} == nil {
		return nil
	}
	return convertFromParent(v)
}

// The singleton implementation (Get, convertFromParent, sync, key) follows.
{{ template "sub_resource_manager_singleton" . }}

// NewSpecDelta returns a delta comparing two {{ .CRD.Kind }}Spec values.
// The parent resource's delta_pre_compare hook can call this to detect
// changes in the sub-resource's spec fields without duplicating the
// comparison logic.
func NewSpecDelta(
	a *svcapitypes.{{ .CRD.Kind }}Spec,
	b *svcapitypes.{{ .CRD.Kind }}Spec,
) *ackcompare.Delta {
	var ra, rb *resource
	if a != nil {
		ra = &resource{ko: &svcapitypes.{{ .CRD.Kind }}{Spec: *a}}
	}
	if b != nil {
		rb = &resource{ko: &svcapitypes.{{ .CRD.Kind }}{Spec: *b}}
	}
	return newResourceDelta(ra, rb)
}
