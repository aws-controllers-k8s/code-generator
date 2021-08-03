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
	assert.Empty(code.FindLateInitializedFieldNames(crd.Config(), crd, "lateInitializeFieldNames", 1))
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
	assert.Empty(code.FindLateInitializedFieldNames(crd.Config(), crd, "lateInitializeFieldNames", 1))
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
		`	var lateInitializeFieldNames = []string{"ImageTagMutability","Name",}
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
	assert.Empty(code.LateInitializeFromReadOne(crd.Config(), crd, "observed", "koWithDefaults", 1))
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
		`	if observed.Spec.ImageTagMutability != nil && koWithDefaults.Spec.ImageTagMutability == nil {
		koWithDefaults.Spec.ImageTagMutability = observed.Spec.ImageTagMutability
	}
	if observed.Spec.Name != nil && koWithDefaults.Spec.Name == nil {
		koWithDefaults.Spec.Name = observed.Spec.Name
	}
`
	assert.Equal(expected, code.LateInitializeFromReadOne(crd.Config(), crd, "observed", "koWithDefaults", 1))
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
		`	if observed.Spec.ImageScanningConfiguration != nil && koWithDefaults.Spec.ImageScanningConfiguration != nil {
		if observed.Spec.ImageScanningConfiguration.ScanOnPush != nil && koWithDefaults.Spec.ImageScanningConfiguration.ScanOnPush == nil {
			koWithDefaults.Spec.ImageScanningConfiguration.ScanOnPush = observed.Spec.ImageScanningConfiguration.ScanOnPush
		}
	}
	if observed.Spec.Name != nil && koWithDefaults.Spec.Name == nil {
		koWithDefaults.Spec.Name = observed.Spec.Name
	}
	if observed.Spec.another != nil && koWithDefaults.Spec.another != nil {
		if observed.Spec.another.map != nil && koWithDefaults.Spec.another.map != nil {
			if observed.Spec.another.map["lastfield"] != nil && koWithDefaults.Spec.another.map["lastfield"] == nil {
				koWithDefaults.Spec.another.map["lastfield"] = observed.Spec.another.map["lastfield"]
			}
		}
	}
	if observed.Spec.map != nil && koWithDefaults.Spec.map != nil {
		if observed.Spec.map["subfield"] != nil && koWithDefaults.Spec.map["subfield"] != nil {
			if observed.Spec.map["subfield"].x != nil && koWithDefaults.Spec.map["subfield"].x == nil {
				koWithDefaults.Spec.map["subfield"].x = observed.Spec.map["subfield"].x
			}
		}
	}
	if observed.Spec.some != nil && koWithDefaults.Spec.some != nil {
		if observed.Spec.some.list != nil && koWithDefaults.Spec.some.list == nil {
			koWithDefaults.Spec.some.list = observed.Spec.some.list
		}
	}
	if observed.Spec.structA != nil && koWithDefaults.Spec.structA != nil {
		if observed.Spec.structA.mapB != nil && koWithDefaults.Spec.structA.mapB != nil {
			if observed.Spec.structA.mapB["structC"] != nil && koWithDefaults.Spec.structA.mapB["structC"] != nil {
				if observed.Spec.structA.mapB["structC"].valueD != nil && koWithDefaults.Spec.structA.mapB["structC"].valueD == nil {
					koWithDefaults.Spec.structA.mapB["structC"].valueD = observed.Spec.structA.mapB["structC"].valueD
				}
			}
		}
	}
`
	assert.Equal(expected, code.LateInitializeFromReadOne(crd.Config(), crd, "observed", "koWithDefaults", 1))
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
		`	ko := latestWithDefaults.ko
	numLateInitializationAttempt := ackannotation.GetNumLateInitializationAttempt(ko)
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
	return requeueDelay, incompleteInitialization
`
	assert.Equal(expected, code.CalculateRequeueDelay(crd.Config(), crd, "latestWithDefaults", 1))
}
