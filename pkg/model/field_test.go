package model_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/pkg/names"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestFieldDocumentation(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "eks",
		&testutil.TestingModelOptions{
			GeneratorConfigFile:     "generator.yaml",
			DocumentationConfigFile: "documentation.yaml",
		},
	)

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Cluster", crds)
	require.NotNil(crd)

	specFields := crd.SpecFields

	// We have not altered the docstring for Version from the
	// docstring that comes in the doc-2.json file...
	ltdField := specFields["Version"]
	require.NotNil(ltdField)
	require.NotNil(ltdField.ShapeRef)

	require.Equal(
		"// The desired Kubernetes version for your cluster. If you don't specify a value\n// here, the latest version available in Amazon EKS is used.",
		ltdField.GetDocumentation(),
	)

	prependField := crd.Fields["ResourcesVPCConfig.SecurityGroupIDs"]
	require.NotNil(prependField)
	require.NotNil(prependField.ShapeRef)

	require.True(
		strings.HasPrefix(prependField.GetDocumentation(), "// !!! Let's take it from the top"),
	)

	appendField := crd.Fields["RoleARN"]
	require.NotNil(appendField)
	require.NotNil(appendField.ShapeRef)

	require.True(
		strings.HasSuffix(appendField.GetDocumentation(), "// !!! Hello earthlings! I have come to assume your permissions."),
	)

	overrideField := crd.Fields["Identity.OIDC.Issuer"]
	require.NotNil(overrideField)
	require.NotNil(overrideField.ShapeRef)

	require.Equal(
		"// !!! All your docs has become mine\n"+
			"// \n"+
			"// That whitespace is entirely on purpose",
		overrideField.GetDocumentation(),
	)

}

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

func TestCustomFieldType(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "iam")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Role", crds)
	require.NotNil(crd)

	// The Role resource has a custom field called Policies that is of type
	// []*string. This field is custom because it is not inferred via either
	// the Create Input/Output shape or the SourceFieldConfig attribute in the
	// generator.yaml file but rather via a `type` attribute of the
	// FieldConfig, which overrides the Go type of the custom field.
	policiesField := crd.Fields["Policies"]
	require.NotNil(policiesField)

	assert.Equal("[]*string", policiesField.GoType)
	require.NotNil(policiesField.ShapeRef)

	// A map and a scalar custom field are also added in the testdata
	// generator.yaml file.
	logConfigField := crd.Fields["LoggingConfig"]
	require.NotNil(logConfigField)

	assert.Equal("map[string]*bool", logConfigField.GoType)
	require.NotNil(logConfigField.ShapeRef)

	myIntField := crd.Fields["MyCustomInteger"]
	require.NotNil(myIntField)

	assert.Equal("*int64", myIntField.GoType)
	require.NotNil(myIntField.ShapeRef)
}

func TestGetReferenceFieldName(t *testing.T) {
	assert := assert.New(t)

	stringShape := api.ShapeRef{
		Shape: &api.Shape{
			Type: "string",
		},
	}

	listShape := api.ShapeRef{
		Shape: &api.Shape{
			Type: "list",
		},
	}

	testCases := []struct {
		fieldName                  string
		expectedReferenceFieldName string
		shapeRef                   *api.ShapeRef
	}{
		{"ClusterName", "ClusterRef", &stringShape},
		{"ClusterNames", "ClusterRefs", &listShape},
		{"ClusterARN", "ClusterRef", &stringShape},
		{"ClusterARNs", "ClusterRefs", &listShape},
		{"ClusterID", "ClusterRef", &stringShape},
		{"ClusterId", "ClusterRef", &stringShape},
		{"ClusterIds", "ClusterRefs", &listShape},
		{"ClusterIDs", "ClusterRefs", &listShape},
		{"Cluster", "ClusterRef", &stringShape},
		{"Clusters", "ClusterRefs", &listShape},
		// When the resource name indicates plural but it is singular. Ex: DHCPOptions
		{"Clusters", "ClustersRef", &stringShape},
		{"BlueDeploymentId", "BlueDeploymentRef", &stringShape},
		{"GreenDeploymentId", "GreenDeploymentRef", &stringShape},
	}

	for _, tc := range testCases {
		f := model.Field{}
		f.ShapeRef = tc.shapeRef
		f.Names = names.New(tc.fieldName)
		referenceFieldName := f.GetReferenceFieldName().Camel
		msg := fmt.Sprintf("for %s, expected reference field name of %s but got %s",
			tc.fieldName, tc.expectedReferenceFieldName, referenceFieldName)
		assert.Equal(tc.expectedReferenceFieldName, referenceFieldName, msg)
	}
}

func TestGetGoTag(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	g := testutil.NewModelForServiceWithOptions(t, "eks",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-gotag.yaml",
		},
	)

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Cluster", crds)
	require.NotNil(crd)

	testCases := []struct {
		name        string
		field       *model.Field
		expectedTag string
	}{
		{
			name:        "not required spec field",
			field:       crd.SpecFields["Logging"],
			expectedTag: "`json:\"logging,omitempty\"`",
		},
		{
			name:        "required spec field",
			field:       crd.SpecFields["Name"],
			expectedTag: "`json:\"name\"`",
		},
		{
			name:        "spec field with go_tag override",
			field:       crd.SpecFields["Version"],
			expectedTag: "`json:\"myCustomVersionName,omitempty\"`",
		},
		{
			// Status fields are not required, so they should always have omitempty directive
			name:        "status field",
			field:       crd.StatusFields["Endpoint"],
			expectedTag: "`json:\"endpoint,omitempty\"`",
		},
		{
			name:        "status field with go_tag override",
			field:       crd.StatusFields["Status"],
			expectedTag: "`json:\"clusterState,omitempty\" yaml:\"some_extra_tags\"`",
		},
	}
	for _, tc := range testCases {
		tag := tc.field.GetGoTag()
		assert.Equal(tc.expectedTag, tag)
	}
}

func TestReferenceFieldPath(t *testing.T) {
	assert := assert.New(t)

	field := model.Field{}
	//Non nested fieldPath
	field.Path = "MyField"
	field.Names = names.New("MyField")
	assert.Equal("MyFieldRef", field.ReferenceFieldPath())

	// Nested fieldPath
	field.Path = "subPathA.subPathB.MyField"
	field.Names = names.New("MyField")
	assert.Equal("subPathA.subPathB.MyFieldRef", field.ReferenceFieldPath())
}

func TestFieldPathWithUnderscore(t *testing.T) {
	assert := assert.New(t)

	field := model.Field{}
	//Empty fieldPath
	field.Path = ""
	assert.Equal("", field.FieldPathWithUnderscore())

	//Non nested fieldPath
	field.Path = "MyField"
	assert.Equal("MyField", field.FieldPathWithUnderscore())

	// Nested fieldPath
	field.Path = "subPathA.subPathB.MyField"
	assert.Equal("subPathA_subPathB_MyField", field.FieldPathWithUnderscore())
}
