// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestSkipRepeatedUpdates_Configured(t *testing.T) {
	g := testutil.NewModelForServiceWithOptions(t, "ec2",
		&testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-skip-repeated-updates.yaml"})
	crds, err := g.GetCRDs()
	require.NoError(t, err)
	require.NotEmpty(t, crds)
	for _, crd := range crds {
		assert.True(t, crd.SkipRepeatedUpdates(), crd.Names.Original)
	}
}

func TestSkipRepeatedUpdates_NotConfigured(t *testing.T) {
	g := testutil.NewModelForService(t, "ec2")
	crds, err := g.GetCRDs()
	require.NoError(t, err)
	crd := getCRDByName("DhcpOptions", crds)
	require.NotNil(t, crd)
	assert.False(t, crd.SkipRepeatedUpdates(), "off unless generator.yaml enables it")
}
