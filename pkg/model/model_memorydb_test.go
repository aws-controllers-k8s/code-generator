package model_test

import (
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserPasswords_IsSecret(t *testing.T) {
	require := require.New(t)

	g := testutil.NewModelForService(t, "memorydb")
	crds, err := g.GetCRDs()

	require.Nil(err)

	crd := getCRDByName("User", crds)
	require.NotNil(crd)
	assert := assert.New(t)
	assert.Equal("[]*ackv1alpha1.SecretKeyReference", crd.SpecFields["AuthenticationMode"].MemberFields["Passwords"].GoType)
	assert.Equal("SecretKeyReference", crd.SpecFields["AuthenticationMode"].MemberFields["Passwords"].GoTypeElem)
	assert.Equal("[]*ackv1alpha1.SecretKeyReference", crd.SpecFields["AuthenticationMode"].MemberFields["Passwords"].GoTypeWithPkgName)
}
