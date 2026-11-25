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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// skipTests are a resource's generated tests: its own, a variant's and the
// adoption test, named as Generate names them.
func skipTests(prefix string) []test {
	return []test{
		{Metadata: objectMeta{Name: prefix + "launch-template"},
			Spec: testSpec{Description: "LaunchTemplate lifecycle."}},
		{Metadata: objectMeta{Name: prefix + "launch-template-minimal"},
			Spec: testSpec{Description: "LaunchTemplate lifecycle, variant minimal."}},
		{Metadata: objectMeta{Name: prefix + "launch-template-adoption"}},
	}
}

func TestApplySkips(t *testing.T) {
	g := &generator{
		in: &Inputs{Skip: Skip{Tests: map[string]string{
			"launch-template-minimal":  "the delete never sees the template gone, issue #123",
			"launch-template-adoption": "the read never sets spec.name, issue #124",
		}}},
		report: &Report{},
	}
	tests := skipTests("")
	require.NoError(t, g.applySkips(tests))

	// The test that isn't skipped keeps its description and runs
	assert.Nil(t, tests[0].Spec.Skip)
	assert.Equal(t, "LaunchTemplate lifecycle.", tests[0].Spec.Description)

	require.NotNil(t, tests[1].Spec.Skip)
	assert.True(t, *tests[1].Spec.Skip)
	assert.Equal(t,
		"LaunchTemplate lifecycle, variant minimal. Skipped: the delete never sees the template gone, issue #123",
		tests[1].Spec.Description)

	// A test with no description of its own
	require.NotNil(t, tests[2].Spec.Skip)
	assert.Equal(t, "Skipped: the read never sets spec.name, issue #124", tests[2].Spec.Description)

	// The report lists them in a stable order, so regenerating is a no-op
	assert.Equal(t, []string{
		"launch-template-adoption: the read never sets spec.name, issue #124",
		"launch-template-minimal: the delete never sees the template gone, issue #123",
	}, g.report.SkippedTests)
}

func TestApplySkipsNamePrefix(t *testing.T) {
	g := &generator{
		in:     &Inputs{Skip: Skip{Tests: map[string]string{"launch-template": "issue #123"}}},
		opts:   Options{NamePrefix: "ec2-"},
		report: &Report{},
	}
	tests := skipTests("ec2-")
	require.NoError(t, g.applySkips(tests))
	require.NotNil(t, tests[0].Spec.Skip)
}

func TestApplySkipsUnknownTest(t *testing.T) {
	g := &generator{
		in:     &Inputs{Skip: Skip{Tests: map[string]string{"launch-template-typo": "issue #123"}}},
		report: &Report{},
	}
	err := g.applySkips(skipTests(""))
	require.Error(t, err)
	assert.ErrorContains(t, err, `skip "launch-template-typo" isn't one of this resource's tests`)
	// The message lists the names to pick from
	assert.ErrorContains(t, err, "launch-template, launch-template-adoption, launch-template-minimal")
}

func TestWithoutSkippedUpdates(t *testing.T) {
	g := &generator{
		in: &Inputs{Skip: Skip{Updates: map[string]string{
			"instanceInitiatedShutdownBehavior": "AWS never returns it, issue #125",
			"disableAPIStop":                    "AWS never returns it, issue #125",
		}}},
		report: &Report{},
	}
	updates := []map[string]interface{}{
		// Mixed: the skipped fields come out, the rest stays
		{"instanceInitiatedShutdownBehavior": "terminate", "sourceDestCheckEnabled": false, "disableAPIStop": true},
		// Only skipped fields: the whole update is dropped
		{"disableAPIStop": true},
		// Nothing skipped: unchanged
		{"tags": []interface{}{"x"}},
	}
	got := g.withoutSkippedUpdates(updates)
	require.Len(t, got, 2)
	assert.Equal(t, map[string]interface{}{"sourceDestCheckEnabled": false}, got[0])
	assert.Equal(t, map[string]interface{}{"tags": []interface{}{"x"}}, got[1])

	// Each skipped field is recorded once, for the report, even one that no
	// update sets
	assert.Equal(t, []string{
		"disableAPIStop: AWS never returns it, issue #125",
		"instanceInitiatedShutdownBehavior: AWS never returns it, issue #125",
	}, g.report.SkippedUpdates)
}

func TestWithoutSkippedUpdatesNoSkips(t *testing.T) {
	g := &generator{in: &Inputs{}, report: &Report{}}
	in := []map[string]interface{}{{"a": 1}}
	assert.Equal(t, in, g.withoutSkippedUpdates(in))
	assert.Empty(t, g.report.SkippedUpdates)
}
