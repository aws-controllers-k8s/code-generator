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

// key returns the primary key value for a resource, used to identify
// individual items within the sub-resource collection.
func key(r *resource) string {
	return *r.ko.{{ ManagerPrimaryKey .CRD }}
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

// sync reconciles the sub-resource collection by comparing the desired and
// latest state, computing a diff, and issuing the necessary create/update/
// delete SDK calls.
func (rm *resourceManager) sync(
	ctx context.Context,
	desired []resource,
	latest []resource,
) error {
	var err error
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sync")
	defer func() {
		exit(err)
	}()

	d := computeDelta(desired, latest)

{{ $srcType := ManagerSourceType .CRD }}
{{ if $srcType.BatchFieldPath }}
	if len(d.toCreate) > 0 {
		merged := mergeResources(d.toCreate)
		if merged != nil {
			_, err = rm.sdkCreate(ctx, merged)
			if err != nil {
				return err
			}
		}
	}
	if len(d.toUpdate) > 0 {
		merged := mergeResources(d.toUpdate)
		if merged != nil {
			_, err = rm.sdkCreate(ctx, merged)
			if err != nil {
				return err
			}
		}
	}
	if len(d.toDelete) > 0 {
		merged := mergeResources(d.toDelete)
		if merged != nil {
			_, err = rm.sdkDelete(ctx, merged)
			if err != nil {
				return err
			}
		}
	}
{{ else }}
	for _, r := range d.toCreate {
		_, err = rm.sdkCreate(ctx, r)
		if err != nil {
			return err
		}
	}
	for _, r := range d.toUpdate {
		_, err = rm.sdkCreate(ctx, r)
		if err != nil {
			return err
		}
	}
	for _, r := range d.toDelete {
		_, err = rm.sdkDelete(ctx, r)
		if err != nil {
			return err
		}
	}
{{ end }}
	return nil
}

{{ if $srcType.BatchFieldPath }}
// mergeResources combines multiple sub-resource items into a single resource
// by concatenating the batch collection field ({{ $srcType.BatchFieldPath }}).
func mergeResources(items []*resource) *resource {
	if len(items) == 0 {
		return nil
	}
	if len(items) == 1 {
		return items[0]
	}
	merged := items[0].ko.DeepCopy()
	merged.{{ $srcType.BatchFieldPath }} = nil
	for _, item := range items {
		merged.{{ $srcType.BatchFieldPath }} = append(
			merged.{{ $srcType.BatchFieldPath }},
			item.ko.{{ $srcType.BatchFieldPath }}...,
		)
	}
	return &resource{ko: merged}
}
{{ end }}

{{ $readFieldPath := ManagerReadFieldPath .CRD }}

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
// resource slices, then delegates to sync.
func (rm *resourceManager) Sync(
	ctx context.Context,
	desired any,
	latest any,
) error {
	desiredResources := convertToResources(desired)
	latestResources := convertToResources(latest)
	return rm.sync(ctx, desiredResources, latestResources)
}

// Get reads the current state of the sub-resources from AWS and writes the
// result back onto the parent resource.
func (rm *resourceManager) Get(
	ctx context.Context,
	parent *svcapitypes.{{ $srcType.ParentKind }},
) error {
	seeds := convertFromParent(parent)
	if len(seeds) == 0 {
		parent.{{ $srcType.FieldPath }} = nil
		return nil
	}
	found, err := rm.sdkFind(ctx, &seeds[0])
	if err != nil {
		// For sub-resources, NotFound from sdkFind means the parent exists
		// but has no items (e.g. no attached policies, no tags). This is
		// valid — clear the parent field and return success.
		if err == ackerr.NotFound {
			parent.{{ $srcType.FieldPath }} = nil
			return nil
		}
		return err
	}
	if found == nil {
		parent.{{ $srcType.FieldPath }} = nil
		return nil
	}
{{- if $readFieldPath }}
	parent.{{ $srcType.FieldPath }} = found.ko.{{ $readFieldPath }}
{{- end }}
	return nil
}

// convertToResources maps fields from the parent resource into sub-resource
// instances using the source-type-specific convertFromParent.
func convertToResources(parent any) []resource {
	if parent == nil {
		return nil
	}
	switch v := parent.(type) {
	case *svcapitypes.{{ $srcType.ParentKind }}:
		if v == nil {
			return nil
		}
		return convertFromParent(v)
	}
	return nil
}

// convertFromParent is source-type-specific: it maps parent fields into
// sub-resource instances.
{{- if $srcType.IsScalar }}
{{ template "sub_resource_manager_scalar" . }}
{{- else if $srcType.IsStruct }}
{{ template "sub_resource_manager_struct" . }}
{{- else if $srcType.IsListScalar }}
{{ template "sub_resource_manager_list_scalar" . }}
{{- else if $srcType.IsListStruct }}
{{ template "sub_resource_manager_list_struct" . }}
{{- else if $srcType.IsMap }}
{{ template "sub_resource_manager_map" . }}
{{- else if $srcType.IsMapScalar }}
{{ template "sub_resource_manager_map_scalar" . }}
{{- else if $srcType.IsMapStruct }}
{{ template "sub_resource_manager_map_struct" . }}
{{- end }}
