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

package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestEKS_IgnoreFieldPaths(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForService(t, "eks")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Nodegroup", crds)
	require.NotNil(crd)

	// fields we have ignored via generator.yaml
	ignoredFieldPaths := []string{
		"Version",
		"ScalingConfig.DesiredSize",
	}
	for _, ignoredFieldPath := range ignoredFieldPaths {
		_, found := crd.Fields[ignoredFieldPath]
		require.False(found, "field should be ignored: %v", ignoredFieldPath)
	}

	// Ensure we do not accidentally remove sibling fields
	_, found := crd.Fields["ScalingConfig.MaxSize"]
	require.True(found, "Sibling should be kept when ScalingConfig.DesiredSize is ignored")
}
