package code_test

import(
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FindLateInitializedFieldsWithDelay_EmptyFieldConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ecr")

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	// NO fieldConfig
	assert.Empty(crd.Config().ResourceFields(crd.Names.Original))
	assert.Empty(code.FindLateInitializedFieldsWithDelay(crd.Config(), crd, 1))
}

func Test_FindLateInitializedFieldsWithDelay_NoLateInitializations(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-field-config.yaml"})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	// FieldConfig without lateInitialize
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["Name"])
	assert.Nil(crd.Config().ResourceFields(crd.Names.Original)["Name"].LateInitialize)
	assert.Empty(code.FindLateInitializedFieldsWithDelay(crd.Config(), crd, 1))
}

func Test_FindLateInitializedFieldsWithDelay(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-late-initialize.yaml"})

	crd := testutil.GetCRDByName(t, g, "Repository")
	require.NotNil(crd)
	// FieldConfig without lateInitialize
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["Name"])
	assert.NotEmpty(crd.Config().ResourceFields(crd.Names.Original)["ImageTagMutability"])
	assert.NotNil(crd.Config().ResourceFields(crd.Names.Original)["Name"].LateInitialize)
	assert.NotNil(crd.Config().ResourceFields(crd.Names.Original)["ImageTagMutability"].LateInitialize)
	expected := "\tvar lateInitializeFieldToDelaySeconds = map[string]int{\"ImageTagMutability\":5,\"Name\":0,}\n"
	assert.Equal(expected, code.FindLateInitializedFieldsWithDelay(crd.Config(), crd, 1))
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
	// FieldConfig without lateInitialize
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
	// FieldConfig without lateInitialize
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
`
	assert.Equal(expected, code.LateInitializeFromReadOne(crd.Config(), crd, "observed", "koWithDefaults", 1))
}
