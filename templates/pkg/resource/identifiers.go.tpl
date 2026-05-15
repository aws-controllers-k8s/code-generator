{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
)

// resourceIdentifiers implements the
// `aws-service-operator-k8s/pkg/types.AWSResourceIdentifiers` interface
type resourceIdentifiers struct {
	meta *ackv1alpha1.ResourceMetadata
}

// ARN returns the AWS Resource Name for the backend AWS resource. If nil,
// this means the resource has not yet been created in the backend AWS
// service.
func (ri *resourceIdentifiers) ARN() *ackv1alpha1.AWSResourceName {
	if ri.meta != nil {
		return ri.meta.ARN
	}
	return nil
}

// OwnerAccountID returns the AWS account identifier in which the
// backend AWS resource resides, or nil if this information is not known
// for the resource
func (ri *resourceIdentifiers) OwnerAccountID() *ackv1alpha1.AWSAccountID {
	if ri.meta != nil {
		return ri.meta.OwnerAccountID
	}
	return nil
}

// Region returns the AWS region in which the resource exists, or
// nil if this information is not known.
func (ri *resourceIdentifiers) Region() *ackv1alpha1.AWSRegion {
	if ri.meta != nil {
		return ri.meta.Region
	}
	return nil
}

// Partition returns the AWS partition in which the reosurce exists, or
// nil if this information is not known. 
func (ri *resourceIdentifiers) Partition() *ackv1alpha1.AWSPartition {
	if ri.meta != nil {
		return ri.meta.Partition
	}
	return nil
}
