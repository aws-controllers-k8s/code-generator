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

package code_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestResourceIsUpdateable_SingleCondition(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "lambda")
	crd := testutil.GetCRDByName(t, g, "Function")
	require.NotNil(crd)

	expected := `	if latest.ko.Status.State != nil {
		if !ackutil.InStrings(*latest.ko.Status.State, []string{"Active"}) {
			return nil, ackrequeue.NeededAfter(
				fmt.Errorf("resource is in %s state, cannot be updated",
					*latest.ko.Status.State),
				time.Duration(30)*time.Second,
			)
		}
	}
`
	got, err := code.ResourceIsUpdateable(
		crd.Config(), crd, "latest", 1,
	)
	require.NoError(err)
	assert.Equal(expected, got)
}
func TestResourceIsDeletable_SingleCondition(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "lambda")
	crd := testutil.GetCRDByName(t, g, "Function")
	require.NotNil(crd)

	expected := `	if r.ko.Status.State != nil {
		if !ackutil.InStrings(*r.ko.Status.State, []string{"Active", "Failed"}) {
			return nil, ackrequeue.NeededAfter(
				fmt.Errorf("resource is in %s state, cannot be deleted",
					*r.ko.Status.State),
				time.Duration(30)*time.Second,
			)
		}
	}
`
	got, err := code.ResourceIsDeletable(
		crd.Config(), crd, "r", 1,
	)
	require.NoError(err)
	assert.Equal(expected, got)
}

func TestResourceIsUpdateable_MultipleConditions_CustomRequeue(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "lambda", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-updateable-multi-condition.yaml",
	})
	crd := testutil.GetCRDByName(t, g, "Function")
	require.NotNil(crd)

	expected := `	if latest.ko.Status.State != nil {
		if !ackutil.InStrings(*latest.ko.Status.State, []string{"Active"}) {
			return nil, ackrequeue.NeededAfter(
				fmt.Errorf("resource is in %s state, cannot be updated",
					*latest.ko.Status.State),
				time.Duration(15)*time.Second,
			)
		}
	}
	if latest.ko.Status.LastUpdateStatus != nil {
		if !ackutil.InStrings(*latest.ko.Status.LastUpdateStatus, []string{"Successful"}) {
			return nil, ackrequeue.NeededAfter(
				fmt.Errorf("resource is in %s state, cannot be updated",
					*latest.ko.Status.LastUpdateStatus),
				time.Duration(15)*time.Second,
			)
		}
	}
`
	got, err := code.ResourceIsUpdateable(
		crd.Config(), crd, "latest", 1,
	)
	require.NoError(err)
	assert.Equal(expected, got)
}

func TestResourceIsUpdateable_NoConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "lambda")
	crd := testutil.GetCRDByName(t, g, "CodeSigningConfig")
	require.NotNil(crd)

	got, err := code.ResourceIsUpdateable(
		crd.Config(), crd, "latest", 1,
	)
	require.NoError(err)
	assert.Equal("", got)
}

func TestResourceIsDeletable_NoConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "lambda")
	crd := testutil.GetCRDByName(t, g, "CodeSigningConfig")
	require.NotNil(crd)

	got, err := code.ResourceIsDeletable(
		crd.Config(), crd, "r", 1,
	)
	require.NoError(err)
	assert.Equal("", got)
}

func TestResourceIsUpdateable_EmptyPath(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "lambda", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-updateable-empty-path.yaml",
	})
	crd := testutil.GetCRDByName(t, g, "Function")
	require.NotNil(crd)

	_, err := code.ResourceIsUpdateable(
		crd.Config(), crd, "latest", 1,
	)
	require.Error(err)
	require.Contains(err.Error(), "empty path")
}

func TestResourceIsUpdateable_EmptyIn(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "lambda", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-updateable-empty-in.yaml",
	})
	crd := testutil.GetCRDByName(t, g, "Function")
	require.NotNil(crd)

	_, err := code.ResourceIsUpdateable(
		crd.Config(), crd, "latest", 1,
	)
	require.Error(err)
	require.Contains(err.Error(), "must not be empty")
}

func TestResourceIsDeletable_MultipleConditions_CustomRequeue(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "lambda", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-deletable-multi-condition.yaml",
	})
	crd := testutil.GetCRDByName(t, g, "Function")
	require.NotNil(crd)

	expected := `	if r.ko.Status.State != nil {
		if !ackutil.InStrings(*r.ko.Status.State, []string{"Active"}) {
			return nil, ackrequeue.NeededAfter(
				fmt.Errorf("resource is in %s state, cannot be deleted",
					*r.ko.Status.State),
				time.Duration(20)*time.Second,
			)
		}
	}
	if r.ko.Status.LastUpdateStatus != nil {
		if !ackutil.InStrings(*r.ko.Status.LastUpdateStatus, []string{"Successful"}) {
			return nil, ackrequeue.NeededAfter(
				fmt.Errorf("resource is in %s state, cannot be deleted",
					*r.ko.Status.LastUpdateStatus),
				time.Duration(20)*time.Second,
			)
		}
	}
`
	got, err := code.ResourceIsDeletable(
		crd.Config(), crd, "r", 1,
	)
	require.NoError(err)
	assert.Equal(expected, got)
}

func TestResourceIsUpdateable_InvalidFieldPath(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "lambda", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-updateable-invalid-path.yaml",
	})
	crd := testutil.GetCRDByName(t, g, "Function")
	require.NotNil(crd)

	_, err := code.ResourceIsUpdateable(
		crd.Config(), crd, "latest", 1,
	)
	require.Error(err)
	require.Contains(err.Error(), "cannot find field")
}
