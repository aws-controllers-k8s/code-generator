{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"fmt"
	"sync"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcfg "github.com/aws-controllers-k8s/runtime/pkg/config"
	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-logr/logr"

	svcresource "github.com/aws-controllers-k8s/{{ .ControllerName }}-controller/pkg/resource"
)

// resourceManagerFactory produces resourceManager objects. It implements the
// `types.AWSResourceManagerFactory` interface.
type resourceManagerFactory struct {
	sync.RWMutex
	// rmCache contains resource managers for a particular AWS account ID
	rmCache map[string]*resourceManager
}

// ResourcePrototype returns an AWSResource that resource managers produced by
// this factory will handle
func (f *resourceManagerFactory) ResourceDescriptor() acktypes.AWSResourceDescriptor {
	return &resourceDescriptor{}
}

// GetCachedManager returns a manager object that can manage resources for a
// supplied AWS account if it was already created and cached, or nil if not
func (f *resourceManagerFactory) GetCachedManager(
	id ackv1alpha1.AWSAccountID,
	region ackv1alpha1.AWSRegion,
	roleARN ackv1alpha1.AWSResourceName,
) acktypes.AWSResourceManager {
	// We use the account ID, region, and role ARN to uniquely identify a
	// resource manager. This helps us to avoid creating multiple resource
	// managers for the same account/region/roleARN combination.
	rmId := fmt.Sprintf("%s/%s/%s", id, region, roleARN)
	f.RLock()
	rm, found := f.rmCache[rmId]
	f.RUnlock()
	if !found {
		return nil
	}

	return rm
}

// ManagerFor returns a resource manager object that can manage resources for a
// supplied AWS account
func (f *resourceManagerFactory) ManagerFor(
	cfg ackcfg.Config,
	clientcfg aws.Config,
	log logr.Logger,
	metrics *ackmetrics.Metrics,
	rr acktypes.Reconciler,
	id ackv1alpha1.AWSAccountID,
	region ackv1alpha1.AWSRegion,
	partition ackv1alpha1.AWSPartition,
	roleARN ackv1alpha1.AWSResourceName,
) (acktypes.AWSResourceManager, error) {
	f.Lock()
	defer f.Unlock()

	// We use the account ID, region, partition, and role ARN to uniquely identify a
	// resource manager. This helps us to avoid creating multiple resource
	// managers for the same account/region/roleARN combination.
	rmId := fmt.Sprintf("%s/%s/%s", id, region, roleARN)
	rm, err := newResourceManager(cfg, clientcfg, log, metrics, rr, id, region, partition)
	if err != nil {
		return nil, err
	}
	f.rmCache[rmId] = rm
	return rm, nil
}

// IsAdoptable returns true if the resource is able to be adopted
func (f *resourceManagerFactory) IsAdoptable() bool {
	return {{ .CRD.IsAdoptable }}
}

// RequeueOnSuccessSeconds returns true if the resource should be requeued after specified seconds
// Default is false which means resource will not be requeued after success. 
func (f *resourceManagerFactory) RequeueOnSuccessSeconds() int {
{{- if $reconcileRequeuOnSuccessSeconds := .CRD.ReconcileRequeuOnSuccessSeconds }}
	return {{ $reconcileRequeuOnSuccessSeconds }}
{{- else }}
	return 0
{{- end }}
}

func newResourceManagerFactory() *resourceManagerFactory {
	return &resourceManagerFactory{
		rmCache: map[string]*resourceManager{},
	}
}

func init() {
	svcresource.RegisterManagerFactory(newResourceManagerFactory())
}
