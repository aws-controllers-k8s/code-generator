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

package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestEMRContainers_JobRun(t *testing.T) {
	assert := assert.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "emrcontainers", &testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-cycle.yaml",
	})

	assert.Panics(func() { g.GetCRDs() })
}
