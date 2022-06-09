package code_test

import (
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInitializeNestedStructField(t *testing.T) {
	assert := assert.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "s3",
		&testutil.TestingModelOptions{GeneratorConfigFile: "generator-with-tags.yaml"})

	crd := testutil.GetCRDByName(t, g, "Bucket")
	assert.NotNil(crd)

	f := crd.Fields["Logging.LoggingEnabled.TargetBucket"]

	s := code.InitializeNestedStructField(crd, "r.ko", f,
		"svcapitypes", 1)
	expected :=
		`	r.ko.Spec.Logging = &svcapitypes.BucketLoggingStatus{}
	r.ko.Spec.Logging.LoggingEnabled = &svcapitypes.LoggingEnabled{}
`
	assert.Equal(expected, s)
}
