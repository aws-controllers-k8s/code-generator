{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"context"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcfg "github.com/aws-controllers-k8s/runtime/pkg/config"
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

{{ $batch := ManagerBatch .CRD }}

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

{{ if $batch }}
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

{{ $conversion := ManagerConversion .CRD }}
{{ if $conversion }}
{{ $parentFieldPath := ManagerParentFieldPath .CRD }}
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

// Sync converts the desired and latest parent specs into internal resource
// slices using the configured conversion logic, then delegates to sync.
func (rm *resourceManager) Sync(
	ctx context.Context,
	desiredSpec any,
	latestSpec any,
) error {
	desired := ConvertToResources(desiredSpec)
	latest := ConvertToResources(latestSpec)
	return rm.sync(ctx, desired, latest)
}

{{ $readFieldPath := ManagerReadFieldPath .CRD }}
{{- if $readFieldPath }}
// Get reads the current state of the sub-resource from AWS by building a
// seed resource from the parent spec, delegating to sdkFind, and returning
// the configured field.
func (rm *resourceManager) Get(
	ctx context.Context,
	parentSpec any,
) (any, error) {
	seeds := ConvertToResources(parentSpec)
	if len(seeds) == 0 {
		return nil, nil
	}
	found, err := rm.sdkFind(ctx, &seeds[0])
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, nil
	}
	return found.ko.{{ $readFieldPath }}, nil
}
{{- end }}

// ConvertToResources converts a parent resource into a slice of internal
// resource objects.
func ConvertToResources(
	spec any,
) []resource {
	if spec == nil {
		return nil
	}

	switch v := spec.(type) {
{{- range $typeName, $src := $conversion }}
	case svcapitypes.{{ $typeName }}:
		return convertFrom{{ $typeName }}(v)
	case *svcapitypes.{{ $typeName }}:
		if v == nil {
			return nil
		}
		return convertFrom{{ $typeName }}(*v)
{{- end }}
	}
	return nil
}

{{- range $typeName, $src := $conversion }}

func convertFrom{{ $typeName }}(parent svcapitypes.{{ $typeName }}) []resource {
{{- if $src.IsSingleton }}
	src := parent.{{ $parentFieldPath }}
	if src == nil {
		return nil
	}
	ko := &svcapitypes.{{ $.CRD.Kind }}{}
{{- range $specField, $elemPath := $src.Element }}
{{- if eq $elemPath "." }}
	ko.Spec.{{ $specField }} = src
{{- else }}
	ko.Spec.{{ $specField }} = src.{{ $elemPath }}
{{- end }}
{{- end }}
{{- range $specField, $parentField := $src.Fields }}
{{- if HasPrefix $parentField "Status.ACKResourceMetadata.ARN" }}
	if parent.Status.ACKResourceMetadata != nil && parent.{{ $parentField }} != nil {
		v := string(*parent.{{ $parentField }})
		ko.Spec.{{ $specField }} = &v
	}
{{- else }}
	ko.Spec.{{ $specField }} = parent.{{ $parentField }}
{{- end }}
{{- end }}
	return []resource{ {ko: ko} }
{{- else }}
	collection := parent.{{ $parentFieldPath }}
	if collection == nil {
		return nil
	}
	var resources []resource
	for {{ if $src.HasMapTokens }}mapKey, mapValue{{ else }}_, elem{{ end }} := range collection {
		ko := &svcapitypes.{{ $.CRD.Kind }}{}
{{- range $specField, $elemPath := $src.Element }}
{{- if eq $elemPath "." }}
		ko.Spec.{{ $specField }} = elem
{{- else if eq $elemPath "[.]" }}
		ko.Spec.{{ $specField }} = append(ko.Spec.{{ $specField }}, elem)
{{- else if eq $elemPath "$key" }}
		{
			k := mapKey
			ko.Spec.{{ $specField }} = &k
		}
{{- else if eq $elemPath "$value" }}
		ko.Spec.{{ $specField }} = mapValue
{{- else }}
		if elem != nil {
			ko.Spec.{{ $specField }} = elem.{{ $elemPath }}
		}
{{- end }}
{{- end }}
{{- range $specField, $parentField := $src.Fields }}
{{- if HasPrefix $parentField "Status.ACKResourceMetadata.ARN" }}
		if parent.Status.ACKResourceMetadata != nil && parent.{{ $parentField }} != nil {
			v := string(*parent.{{ $parentField }})
			ko.Spec.{{ $specField }} = &v
		}
{{- else }}
		ko.Spec.{{ $specField }} = parent.{{ $parentField }}
{{- end }}
{{- end }}
		resources = append(resources, resource{ko: ko})
	}
	return resources
{{- end }}
}
{{- end }}
{{ end }}

{{ if $batch }}
func mergeResources(items []*resource) *resource {
	if len(items) == 0 {
		return nil
	}
	if len(items) == 1 {
		return items[0]
	}
	merged := items[0].ko.DeepCopy()
	merged.Spec.{{ $batch.CollectionFieldPath }} = nil
	for _, item := range items {
		merged.Spec.{{ $batch.CollectionFieldPath }} = append(
			merged.Spec.{{ $batch.CollectionFieldPath }},
			item.ko.Spec.{{ $batch.CollectionFieldPath }}...,
		)
	}
	last := items[len(items)-1]
	lastCopy := last.ko.DeepCopy()
	lastCopy.Spec.{{ $batch.CollectionFieldPath }} = merged.Spec.{{ $batch.CollectionFieldPath }}
	return &resource{ko: lastCopy}
}
{{ end }}
