// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	 http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package ack

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8srt "k8s.io/apimachinery/pkg/runtime"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/require"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcfg "github.com/aws-controllers-k8s/runtime/pkg/config"
	ackmetrics "github.com/aws-controllers-k8s/runtime/pkg/metrics"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
)

type fakeIdentifiers struct{}

func (ids *fakeIdentifiers) ARN() *ackv1alpha1.AWSResourceName {
	arn := ackv1alpha1.AWSResourceName("fake-arn")
	return &arn
}

func (ids *fakeIdentifiers) OwnerAccountID() *ackv1alpha1.AWSAccountID {
	owner := ackv1alpha1.AWSAccountID("fake-owner-account-id")
	return &owner
}

type fakeDescriptor struct{}

func (fd *fakeDescriptor) GroupKind() *metav1.GroupKind {
	return nil
}

func (fd *fakeDescriptor) EmptyRuntimeObject() k8srt.Object {
	return nil
}

func (fd *fakeDescriptor) ResourceFromRuntimeObject(o k8srt.Object) acktypes.AWSResource {
	return nil
}

func (fd *fakeDescriptor) Delta(a, b acktypes.AWSResource) *ackcompare.Delta {
	return nil
}

func (fd *fakeDescriptor) UpdateCRStatus(acktypes.AWSResource) (bool, error) {
	return false, nil
}

func (fd *fakeDescriptor) IsManaged(acktypes.AWSResource) bool {
	return false
}

func (fd *fakeDescriptor) MarkManaged(acktypes.AWSResource) {
}

func (fd *fakeDescriptor) MarkUnmanaged(acktypes.AWSResource) {
}

func (fd *fakeDescriptor) MarkAdopted(acktypes.AWSResource) {
}

type fakeRMF struct{}

func (rmf *fakeRMF) ResourceDescriptor() acktypes.AWSResourceDescriptor {
	return &fakeDescriptor{}
}

func (rmf *fakeRMF) ManagerFor(
	ackcfg.Config, // passed by-value to avoid mutation by consumers
	logr.Logger,
	*ackmetrics.Metrics,
	acktypes.Reconciler,
	*session.Session,
	ackv1alpha1.AWSAccountID,
	ackv1alpha1.AWSRegion,
) (acktypes.AWSResourceManager, error) {
	return nil, nil
}
func (rmf *fakeRMF) IsAdoptable() bool            { return false }
func (rmf *fakeRMF) RequeueOnSuccessSeconds() int { return 10 }

// This test is mostly just a hack to introduce a Go module dependency between
// the ACK runtime library and the code generator. The code generator doesn't
// actually depend on Go code in the ACK runtime, but it *produces* templated
// Go code that itself depends on the ACK runtime's types and interfaces.
func TestRuntimeDependency(t *testing.T) {
	require := require.New(t)

	require.Implements((*acktypes.AWSResourceIdentifiers)(nil), new(fakeIdentifiers))
	require.Implements((*acktypes.AWSResourceDescriptor)(nil), new(fakeDescriptor))

	// ACK runtime 0.2.3 introduced a new logger that is now passed into the
	// Context and retrievable using the `pkg/runtime/log.FromContext`
	// function.  This function returns NoopLogger if no such logger is found
	// in the context, but this check here is mostly to ensure that the new
	// function used in ACK runtime 0.2.3 and templates in code-generator
	// consuming 0.2.3 are properly pinned.
	require.Implements((*acktypes.Logger)(nil), ackrtlog.FromContext(context.TODO()))

	// ACK runtime 0.3.0 introduced a new RequeueOnSuccessSeconds method to the
	// resource manager factory
	require.Implements((*acktypes.AWSResourceManagerFactory)(nil), new(fakeRMF))
}
