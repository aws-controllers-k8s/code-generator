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
	"testing"

	"github.com/stretchr/testify/require"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
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

// This test is mostly just a hack to introduce a Go module dependency between
// the ACK runtime library and the code generator. The code generator doesn't
// actually depend on Go code in the ACK runtime, but it *produces* templated
// Go code that itself depends on the ACK runtime's types and interfaces.
func TestRuntimeDependency(t *testing.T) {
	require := require.New(t)

	require.Implements((*acktypes.AWSResourceIdentifiers)(nil), new(fakeIdentifiers))
}
