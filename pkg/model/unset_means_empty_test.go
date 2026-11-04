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

func TestUnsetMeansEmptyPaths_Configured(t *testing.T) {
	g := testutil.NewModelForServiceWithOptions(t, "ec2",
		&testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-unset-means-empty.yaml"})
	crds, err := g.GetCRDs()
	require.NoError(t, err)
	crd := getCRDByName("SecurityGroup", crds)
	require.NotNil(t, crd)

	paths, err := crd.UnsetMeansEmptyPaths()
	require.NoError(t, err)
	assert.Equal(t, []string{"Spec.Description"}, paths)
}

func TestUnsetMeansEmptyPaths_IncludesTagField(t *testing.T) {
	g := testutil.NewModelForService(t, "ecr")
	crds, err := g.GetCRDs()
	require.NoError(t, err)
	crd := getCRDByName("Repository", crds)
	require.NotNil(t, crd)

	paths, err := crd.UnsetMeansEmptyPaths()
	require.NoError(t, err)
	assert.Equal(t, []string{"Spec.Tags"}, paths)
}

func TestUnsetMeansEmptyPaths_NotConfigured(t *testing.T) {
	g := testutil.NewModelForService(t, "ec2")
	crds, err := g.GetCRDs()
	require.NoError(t, err)
	crd := getCRDByName("DhcpOptions", crds)
	require.NotNil(t, crd)

	paths, err := crd.UnsetMeansEmptyPaths()
	require.NoError(t, err)
	assert.Empty(t, paths, "tags are ignored for this resource in the test model")
}

func TestUnsetMeansEmptyPaths_RejectsPathThroughList(t *testing.T) {
	g := testutil.NewModelForServiceWithOptions(t, "ec2",
		&testutil.TestingModelOptions{GeneratorConfigFile: "generator-unset-means-empty-in-list.yaml"})
	crds, err := g.GetCRDs()
	require.NoError(t, err)
	crd := getCRDByName("DhcpOptions", crds)
	require.NotNil(t, crd)

	_, err = crd.UnsetMeansEmptyPaths()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "goes through the list")
}
