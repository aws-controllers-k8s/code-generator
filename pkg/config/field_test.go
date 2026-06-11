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

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestCELRule_Parsing(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	yamlStr := `
resources:
  MyResource:
    custom_cel_rules:
      - rule: "has(self.foo)"
        message: "foo is required"
      - rule: "self.bar > 0"
    fields:
      MyField:
        custom_cel_rules:
          - rule: "self.matches('^[a-z]+')"
            message: "must be lowercase"
`
	cfg, err := New("", Config{})
	require.Nil(err)
	err = yaml.UnmarshalStrict([]byte(yamlStr), &cfg)
	require.Nil(err)

	resConfig, ok := cfg.Resources["MyResource"]
	require.True(ok)

	// Resource-level rules
	require.Len(resConfig.CustomCELRules, 2)
	assert.Equal("has(self.foo)", resConfig.CustomCELRules[0].Rule)
	require.NotNil(resConfig.CustomCELRules[0].Message)
	assert.Equal("foo is required", *resConfig.CustomCELRules[0].Message)
	assert.Equal("self.bar > 0", resConfig.CustomCELRules[1].Rule)
	assert.Nil(resConfig.CustomCELRules[1].Message) // no message key

	// Field-level rules
	fieldConfig, ok := resConfig.Fields["MyField"]
	require.True(ok)
	require.Len(fieldConfig.CustomCELRules, 1)
	assert.Equal("self.matches('^[a-z]+')", fieldConfig.CustomCELRules[0].Rule)
	require.NotNil(fieldConfig.CustomCELRules[0].Message)
	assert.Equal("must be lowercase", *fieldConfig.CustomCELRules[0].Message)
}
