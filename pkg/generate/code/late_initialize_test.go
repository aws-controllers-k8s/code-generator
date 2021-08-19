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

package code_test

import (
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FindLateInitializedFieldNames_EmptyFieldConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ecr")

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	// NO fieldConfig
	assert.Empty(crd.Config().ResourceFields(crd.Names.Original))
	expected :=
		`	lateInitializeFieldNames = []string{}
`
	assert.Equal(expected, code.FindLateInitializedFieldNames(crd.Config(), crd, "lateInitializeFieldNames", 1))
}

func Test_FindLateInitializedFieldNames_NoLateInitializations(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-field-config.yaml"})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	// FieldConfig without lateInitialize
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["Name"])
	assert.Nil(crd.Config().ResourceFields(crd.Names.Original)["Name"].LateInitialize)
	expected :=
		`	lateInitializeFieldNames = []string{}
`
	assert.Equal(expected, code.FindLateInitializedFieldNames(crd.Config(), crd, "lateInitializeFieldNames", 1))
}

func Test_FindLateInitializedFieldNames(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-late-initialize.yaml"})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["Name"])
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["ImageTagMutability"])
	assert.NotNil(crd.Config().ResourceFields(crd.Names.Original)["Name"].LateInitialize)
	assert.NotNil(crd.Config().ResourceFields(crd.Names.Original)["ImageTagMutability"].LateInitialize)
	expected :=
		`	lateInitializeFieldNames = []string{"ImageTagMutability","Name",}
`
	assert.Equal(expected, code.FindLateInitializedFieldNames(crd.Config(), crd, "lateInitializeFieldNames", 1))
}

func Test_LateInitializeFromReadOne_NoFieldsToLateInitialize(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ecr")

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	// NO fieldConfig
	assert.Empty(crd.Config().ResourceFields(crd.Names.Original))
	expected := "	return latest"
	assert.Equal(expected, code.LateInitializeFromReadOne(crd.Config(), crd, "observed", "latest", 1))
}

func Test_LateInitializeFromReadOne_NonNestedPath(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-late-initialize.yaml"})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["Name"])
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["ImageTagMutability"])
	assert.NotNil(crd.Config().ResourceFields(crd.Names.Original)["Name"].LateInitialize)
	assert.NotNil(crd.Config().ResourceFields(crd.Names.Original)["ImageTagMutability"].LateInitialize)
	expected :=
		`	observedKo := rm.concreteResource(observed).ko
	latestKo := rm.concreteResource(latest).ko
	if observedKo.Spec.ImageTagMutability != nil && latestKo.Spec.ImageTagMutability == nil {
		latestKo.Spec.ImageTagMutability = observedKo.Spec.ImageTagMutability
	}
	if observedKo.Spec.Name != nil && latestKo.Spec.Name == nil {
		latestKo.Spec.Name = observedKo.Spec.Name
	}
	return latest`
	assert.Equal(expected, code.LateInitializeFromReadOne(crd.Config(), crd, "observed", "latest", 1))
}

func Test_LateInitializeFromReadOne_NestedPath(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-nested-path-late-initialize.yaml"})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["Name"])
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["ImageScanningConfiguration.ScanOnPush"])
	assert.NotNil(crd.Config().ResourceFields(crd.Names.Original)["Name"].LateInitialize)
	assert.NotNil(crd.Config().ResourceFields(crd.Names.Original)["ImageScanningConfiguration.ScanOnPush"].LateInitialize)
	expected :=
		`	observedKo := rm.concreteResource(observed).ko
	latestKo := rm.concreteResource(latest).ko
	if observedKo.Spec.ImageScanningConfiguration != nil && latestKo.Spec.ImageScanningConfiguration != nil {
		if observedKo.Spec.ImageScanningConfiguration.ScanOnPush != nil && latestKo.Spec.ImageScanningConfiguration.ScanOnPush == nil {
			latestKo.Spec.ImageScanningConfiguration.ScanOnPush = observedKo.Spec.ImageScanningConfiguration.ScanOnPush
		}
	}
	if observedKo.Spec.Name != nil && latestKo.Spec.Name == nil {
		latestKo.Spec.Name = observedKo.Spec.Name
	}
	if observedKo.Spec.another != nil && latestKo.Spec.another != nil {
		if observedKo.Spec.another.map != nil && latestKo.Spec.another.map != nil {
			if observedKo.Spec.another.map["lastfield"] != nil && latestKo.Spec.another.map["lastfield"] == nil {
				latestKo.Spec.another.map["lastfield"] = observedKo.Spec.another.map["lastfield"]
			}
		}
	}
	if observedKo.Spec.map != nil && latestKo.Spec.map != nil {
		if observedKo.Spec.map["subfield"] != nil && latestKo.Spec.map["subfield"] != nil {
			if observedKo.Spec.map["subfield"].x != nil && latestKo.Spec.map["subfield"].x == nil {
				latestKo.Spec.map["subfield"].x = observedKo.Spec.map["subfield"].x
			}
		}
	}
	if observedKo.Spec.some != nil && latestKo.Spec.some != nil {
		if observedKo.Spec.some.list != nil && latestKo.Spec.some.list == nil {
			latestKo.Spec.some.list = observedKo.Spec.some.list
		}
	}
	if observedKo.Spec.structA != nil && latestKo.Spec.structA != nil {
		if observedKo.Spec.structA.mapB != nil && latestKo.Spec.structA.mapB != nil {
			if observedKo.Spec.structA.mapB["structC"] != nil && latestKo.Spec.structA.mapB["structC"] != nil {
				if observedKo.Spec.structA.mapB["structC"].valueD != nil && latestKo.Spec.structA.mapB["structC"].valueD == nil {
					latestKo.Spec.structA.mapB["structC"].valueD = observedKo.Spec.structA.mapB["structC"].valueD
				}
			}
		}
	}
	return latest`
	assert.Equal(expected, code.LateInitializeFromReadOne(crd.Config(), crd, "observed", "latest", 1))
}

func Test_CalculateRequeueDelay_NoFieldsToLateInitialization(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-field-config.yaml"})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	expected :=
		`	return time.Duration(0), false`
	assert.Equal(expected, code.CalculateRequeueDelay(crd.Config(), crd, "latest", 1))
}

func Test_CalculateRequeueDelay(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-nested-path-late-initialize.yaml"})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["Name"])
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["ImageScanningConfiguration.ScanOnPush"])
	assert.NotNil(crd.Config().ResourceFields(crd.Names.Original)["Name"].LateInitialize)
	assert.NotNil(crd.Config().ResourceFields(crd.Names.Original)["ImageScanningConfiguration.ScanOnPush"].LateInitialize)
	expected :=
		`	ko := rm.concreteResource(latest).ko
	numLateInitializationAttempt := ackannotation.GetNumLateInitializationAttempt(latest.MetaObject())
	requeueDelay := time.Duration(0)*time.Second
	incompleteInitialization := false
	if ko.Spec.ImageScanningConfiguration != nil {
		if ko.Spec.ImageScanningConfiguration.ScanOnPush == nil {
			fDelay := (&acktypes.Exponential{Initial:time.Duration(5)*time.Second, Factor: 2, MaxDelay: time.Duration(15)*time.Second,}).GetBackoff(numLateInitializationAttempt)
			requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
			incompleteInitialization= true
		}
	}
	if ko.Spec.Name == nil {
		fDelay := (&acktypes.Exponential{Initial:time.Duration(0)*time.Second, Factor: 2, MaxDelay: time.Duration(0)*time.Second,}).GetBackoff(numLateInitializationAttempt)
		requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
		incompleteInitialization= true
	}
	if ko.Spec.another != nil {
		if ko.Spec.another.map != nil {
			if ko.Spec.another.map["lastfield"] == nil {
				fDelay := (&acktypes.Exponential{Initial:time.Duration(5)*time.Second, Factor: 2, MaxDelay: time.Duration(0)*time.Second,}).GetBackoff(numLateInitializationAttempt)
				requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
				incompleteInitialization= true
			}
		}
	}
	if ko.Spec.map != nil {
		if ko.Spec.map["subfield"] != nil {
			if ko.Spec.map["subfield"].x == nil {
				fDelay := (&acktypes.Exponential{Initial:time.Duration(5)*time.Second, Factor: 2, MaxDelay: time.Duration(0)*time.Second,}).GetBackoff(numLateInitializationAttempt)
				requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
				incompleteInitialization= true
			}
		}
	}
	if ko.Spec.some != nil {
		if ko.Spec.some.list == nil {
			fDelay := (&acktypes.Exponential{Initial:time.Duration(10)*time.Second, Factor: 2, MaxDelay: time.Duration(0)*time.Second,}).GetBackoff(numLateInitializationAttempt)
			requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
			incompleteInitialization= true
		}
	}
	if ko.Spec.structA != nil {
		if ko.Spec.structA.mapB != nil {
			if ko.Spec.structA.mapB["structC"] != nil {
				if ko.Spec.structA.mapB["structC"].valueD == nil {
					fDelay := (&acktypes.Exponential{Initial:time.Duration(20)*time.Second, Factor: 2, MaxDelay: time.Duration(0)*time.Second,}).GetBackoff(numLateInitializationAttempt)
					requeueDelay = time.Duration(math.Max(requeueDelay.Seconds(), fDelay.Seconds()))*time.Second
					incompleteInitialization= true
				}
			}
		}
	}
	return requeueDelay, incompleteInitialization`
	assert.Equal(expected, code.CalculateRequeueDelay(crd.Config(), crd, "latest", 1))
}
