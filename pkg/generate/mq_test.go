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

package generate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestMQ_Broker(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "mq")

	crd := testutil.GetCRDByName(t, g, "Broker")
	require.NotNil(crd)

	// We want to verify that the `Password` field of the `Spec.Users` field
	// (which is a `[]*User` type) is findable in the CRD's Fields collection
	// by the path `Spec.Users..Password` and that the FieldConfig associated
	// with this Field is marked as a SecretKeyReference.
	passFieldPath := "Users..Password"
	passField, found := crd.Fields[passFieldPath]
	require.True(found)
	require.NotNil(passField.FieldConfig)
	assert.True(passField.FieldConfig.IsSecret)

	// We now verify that the User TypeDef that is inferred by the code
	// generator has had its Go type changed from `string` to
	// `*ackv1alpha1.SecretKeyReference` as part of the nested fields
	// post-processing of type defs in the GetTypeDefs() method.
	tdef := testutil.GetTypeDefByName(t, g, "User")
	require.NotNil(tdef)
	passAttr, found := tdef.Attrs["Password"]
	require.True(found)
	assert.Equal("*ackv1alpha1.SecretKeyReference", passAttr.GoType)
}

func TestMQ_GetOutputShapeGoType(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewGeneratorForService(t, "mq")

	crd := testutil.GetCRDByName(t, g, "Broker")
	require.NotNil(crd)

	exp := "*svcsdkapi.CreateBrokerResponse"
	otype := crd.GetOutputShapeGoType(crd.Ops.Create)
	assert.Equal(exp, otype)
}
