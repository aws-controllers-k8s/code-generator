package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestMemberFields(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("LaunchTemplate", crds)
	require.NotNil(crd)

	specFields := crd.SpecFields

	// LaunchTemplate.Spec.LaunchTemplateData field is itself a struct field
	// with a number of member fields. We are checking here to ensure that the
	// LaunchTemplate.Spec.LaunchTemplateData Field definition properly
	// gathered the member field information.
	ltdField := specFields["LaunchTemplateData"]
	require.NotNil(ltdField)
	require.Equal(ltdField.Path, "LaunchTemplateData")
	require.NotNil(ltdField.ShapeRef)
	require.Equal(ltdField.ShapeRef.Shape.Type, "structure")

	assert.NotNil(ltdField.MemberFields)

	// HibernationOptions is a member of the LaunchTemplateData structure and
	// is itself another structure field. Make sure that the Field definition
	// for this nested structure member field exists.
	hoField := ltdField.MemberFields["HibernationOptions"]
	assert.NotNil(hoField)
	assert.Equal(hoField.Path, "LaunchTemplateData.HibernationOptions")
	assert.NotNil(hoField.MemberFields)

	hocField := hoField.MemberFields["Configured"]
	assert.NotNil(hocField)
	assert.Equal(hocField.Path, "LaunchTemplateData.HibernationOptions.Configured")
}
