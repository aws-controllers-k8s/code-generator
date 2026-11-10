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

package chainsaw

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestSuiteFiles(t *testing.T) {
	files, err := SuiteFiles("s3")
	require.NoError(t, err)

	// The step templates the generated tests use, and the configuration
	var names []string
	for name := range files {
		names = append(names, name)
	}
	assert.ElementsMatch(t, []string{
		"_templates/apply-and-wait-reconciled.yaml",
		"_templates/assert-synced.yaml",
		"_templates/delete.yaml",
		".chainsaw.yaml",
	}, names)

	for name, content := range files {
		assert.True(t, strings.HasPrefix(string(content), GeneratedHeader), name)
		var doc map[string]interface{}
		require.NoError(t, yaml.Unmarshal(content, &doc), name)
		if strings.HasPrefix(name, TemplatesDirName+"/") {
			assert.Equal(t, "StepTemplate", doc["kind"], name)
		}
	}

	config := string(files[ConfigFileName])
	assert.Contains(t, config, "name: s3-controller")
	assert.Contains(t, config, "selector: app.kubernetes.io/name=ack-s3-controller")
	assert.NotContains(t, config, "{{")
	assert.NotContains(t, config, "ec2")
}

func TestWriteGeneratedFile(t *testing.T) {
	dir := t.TempDir()
	generated := []byte(GeneratedHeader + "kind: StepTemplate\n")

	// A new file is written, and so is a generated one
	path := filepath.Join(dir, "delete.yaml")
	written, err := WriteGeneratedFile(path, generated)
	require.NoError(t, err)
	assert.True(t, written)
	written, err = WriteGeneratedFile(path, append(generated, "# changed\n"...))
	require.NoError(t, err)
	assert.True(t, written)

	// A hand-written replacement is kept
	custom := filepath.Join(dir, ".chainsaw.yaml")
	require.NoError(t, os.WriteFile(custom, []byte("# our own timeouts\n"), 0o644))
	written, err = WriteGeneratedFile(custom, generated)
	require.NoError(t, err)
	assert.False(t, written)
	b, err := os.ReadFile(custom)
	require.NoError(t, err)
	assert.Equal(t, "# our own timeouts\n", string(b))
}
