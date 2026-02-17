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

func TestToACKTagsForListShape(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "Vpc")
	require.NotNil(crd)

	expectedSyncedConditions := `
	if len(tags) == 0 {
		return result, keyOrder
	}
	for _, t := range tags {
		if t.Key != nil {
			keyOrder = append(keyOrder, *t.Key)
			if t.Value != nil {
				result[*t.Key] = *t.Value
			} else {
				result[*t.Key] = ""
			}
		}
	}
`
	got, err := code.GoCodeConvertToACKTags(
		crd, "tags", "result", "keyOrder", 1,
	)
	require.NoError(err)
	assert.Equal(expectedSyncedConditions, got)
}

func TestToACKTagsForMapShape(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "Api")
	require.NotNil(crd)

	expectedSyncedConditions := `
	if len(tags) == 0 {
		return result, keyOrder
	}
	for k, v := range tags {
		if v == nil {
			result[k] = ""
		} else {
			result[k] = *v
		}
	}
`
	got, err := code.GoCodeConvertToACKTags(
		crd, "tags", "result", "keyOrder", 1,
	)
	require.NoError(err)
	assert.Equal(expectedSyncedConditions, got)
}

func TestFromACKTagsForListShape(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crd := testutil.GetCRDByName(t, g, "Vpc")
	require.NotNil(crd)

	expectedSyncedConditions := `
	for _, k := range keyOrder {
		v, ok := tags[k]
		if ok {
			tag := svcapitypes.Tag{Key: &k, Value: &v}
			result = append(result, &tag)
			delete(tags, k)
		}
	}
	for k, v := range tags {
		tag := svcapitypes.Tag{Key: &k, Value: &v}
		result = append(result, &tag)
	}
`
	got, err := code.GoCodeFromACKTags(
		crd, "tags", "keyOrder", "result", 1,
	)
	require.NoError(err)
	assert.Equal(expectedSyncedConditions, got)
}

func TestFromACKTagsForMapShape(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "apigatewayv2")

	crd := testutil.GetCRDByName(t, g, "Api")
	require.NotNil(crd)

	expectedSyncedConditions := `
	_ = keyOrder
	for k, v := range tags {
		result[k] = &v
	}
`
	got, err := code.GoCodeFromACKTags(
		crd, "tags", "keyOrder", "result", 1,
	)
	require.NoError(err)
	assert.Equal(expectedSyncedConditions, got)
}