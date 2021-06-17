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

package ack_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/ack"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestResourceHookCodeInline(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	basePaths := []string{}
	hookID := "sdk_update_pre_build_request"

	g := testutil.NewGeneratorForService(t, "mq")

	crd := testutil.GetCRDByName(t, g, "Broker")
	require.NotNil(crd)

	// The Broker's update operation has a special hook callback configured
	expected := `if err := rm.requeueIfNotRunning(latest); err != nil { return nil, err }`
	got, err := ack.ResourceHookCode(basePaths, crd, hookID, nil, nil)
	assert.Nil(err)
	assert.Equal(expected, got)
}

func TestResourceHookCodeTemplatePath(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	wd, _ := os.Getwd()
	basePaths := []string{
		filepath.Join(wd, "testdata", "templates"),
	}
	hookID := "sdk_delete_pre_build_request"

	g := testutil.NewGeneratorForService(t, "mq")

	crd := testutil.GetCRDByName(t, g, "Broker")
	require.NotNil(crd)

	// The Broker's delete operation has a special hook configured to point to a template.
	expected := "// this is my template.\n"
	got, err := ack.ResourceHookCode(basePaths, crd, hookID, nil, nil)
	assert.Nil(err)
	assert.Equal(expected, got)
}
