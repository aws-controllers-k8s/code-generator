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
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// runValue runs a value's command the way the generated values step does.
func runValue(t *testing.T, command string) (string, error) {
	t.Helper()
	script := "set -eu\nV=$(sh -euc " + shQuote(command) + ") || { echo failed; exit 1; }\necho \"$V\"\n"
	out, err := exec.Command("sh", "-c", script).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func TestValueCommandStopsAtFirstFailure(t *testing.T) {
	// A failed step before the final echo, e.g. a create-bucket the role
	// can't do, must fail the value rather than print a name
	out, err := runValue(t, "false\necho bucket-name")
	assert.Error(t, err)
	assert.Equal(t, "failed", out)
}

func TestValueCommandKeepsQuotes(t *testing.T) {
	out, err := runValue(t, `echo '{"a": 1}' | tr -d ' '`)
	assert.NoError(t, err)
	assert.Equal(t, `{"a":1}`, out)
}
