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

	require.NotNil(ltdField.MemberFields)

	// HibernationOptions is a member of the LaunchTemplateData structure and
	// is itself another structure field. Make sure that the Field definition
	// for this nested structure member field exists.
	hoField := ltdField.MemberFields["HibernationOptions"]
	assert.NotNil(hoField)
	assert.Equal(hoField.Path, "LaunchTemplateData.HibernationOptions")
	assert.NotNil(hoField.MemberFields)

	hocField := hoField.MemberFields["Configured"]
	require.NotNil(hocField)
	assert.Equal(hocField.Path, "LaunchTemplateData.HibernationOptions.Configured")
}

func TestMemberFields_Containers_ListOfStruct(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("DHCPOptions", crds)
	require.NotNil(crd)

	// The DHCPOptions resource has a DHCPConfigurations field that is of type
	// []*NewDHCPConfiguration. Here, we test to make sure that the inference
	// process deduced the DHCPOptions.DHCPConfigurations Field object's Values
	// Field properly.
	dhcpConfsField := crd.Fields["DHCPConfigurations"]
	require.NotNil(dhcpConfsField)

	require.NotNil(dhcpConfsField.MemberFields)
	valuesField := dhcpConfsField.MemberFields["Values"]
	require.NotNil(valuesField)
	require.NotNil(valuesField.ShapeRef)
	require.NotNil(valuesField.ShapeRef.Shape)
	assert.Equal(valuesField.ShapeRef.Shape.Type, "list")
	assert.Equal(valuesField.Path, "DHCPConfigurations.Values")
}

func TestMemberFields_Containers_MapOfStruct(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "sagemaker")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("TrialComponent", crds)
	require.NotNil(crd)

	// The TrialComponent resource has an InputArtifacts field that is of type
	// map[string]*TrialComponentArtifact. TrialComponentArtifact is a struct
	// with two member fields, Value and MediaType. Here, we test to make sure
	// that the inference process deduced the TrialComponent.InputArtifacts
	// Field object's Value Field properly.
	inArtifactsField := crd.Fields["InputArtifacts"]
	require.NotNil(inArtifactsField)

	require.NotNil(inArtifactsField.MemberFields)
	valueField := inArtifactsField.MemberFields["Value"]
	require.NotNil(valueField)
	require.NotNil(valueField.ShapeRef)
	require.NotNil(valueField.ShapeRef.Shape)
	assert.Equal(valueField.ShapeRef.Shape.Type, "string")
	assert.Equal(valueField.Path, "InputArtifacts.Value")
}
