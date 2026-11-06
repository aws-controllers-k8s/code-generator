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

// Package chainsaw generates Chainsaw e2e tests for a controller's resources
// from their model, plus a small inputs file per resource with the values the
// model can't supply.
//
// A generated test runs a resource through its lifecycle: create, check it in
// AWS, try to change each immutable field, apply each update and check it in
// AWS, delete and check it's gone. If the inputs give an invalid spec, the
// test also checks it goes ACK.Terminal. Each create variant in the inputs
// gets a lifecycle test of its own. If the resource is adoptable, another
// test adopts one created through a retained CR.
//
// The AWS-side checks read each field back the way the model allows: through
// the read operation if it returns the field, through the operation matching
// the one that sets it (e.g. DescribeVpcAttribute for a field set by
// ModifyVpcAttribute), or for a custom field holding a list of some structure,
// through the read operation's member holding a list of the same structure.
// Where the model doesn't say, the inputs can declare a check. The report
// lists the fields that are left, and why.
package chainsaw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// Options configure the generated tests.
type Options struct {
	// Service is the AWS CLI name of the service, e.g. ec2.
	Service string
	// APIVersion is the CRDs' API version, e.g. v1alpha1.
	APIVersion string
	// TemplatesDir is the path of the shared step templates, relative to the
	// generated test's directory.
	TemplatesDir string
	// NamePrefix is prepended to the names of the generated Chainsaw tests.
	NamePrefix string
}

// Report describes what a generated test covers, and what it can't.
type Report struct {
	Kind string
	// Variants are the names of the variants' tests.
	Variants []string
	// Checked are the spec fields the test checks in AWS.
	Checked []string
	// Unchecked are the spec fields the inputs set but the test can't check in
	// AWS, each with why.
	Unchecked []string
	// NotUpdated are spec fields that aren't immutable but no update changes.
	NotUpdated []string
	// ImmutableTested are the immutable fields the test tries to change.
	ImmutableTested []string
	// ImmutableSkipped are immutable fields the test can't try to change.
	ImmutableSkipped []string
	// Notes are other gaps, e.g. no adoption test.
	Notes []string
}

// String returns the report as YAML comment lines.
func (r *Report) String() string {
	var b strings.Builder
	line := func(title string, items []string) {
		if len(items) == 0 {
			return
		}
		fmt.Fprintf(&b, "# %s:\n", title)
		for _, item := range items {
			fmt.Fprintf(&b, "#   - %s\n", item)
		}
	}
	line("Create variants, each with its own test", r.Variants)
	line("Checked in AWS", r.Checked)
	line("Not checked in AWS (declare a check in inputs.yaml)", r.Unchecked)
	line("Not immutable, but no update changes them (mark them is_immutable if the controller can't update them)", r.NotUpdated)
	line("Immutable, tried to change", r.ImmutableTested)
	line("Immutable, not tried (no automatic changed value)", r.ImmutableSkipped)
	line("Notes", r.Notes)
	return b.String()
}

// Generate returns the files of a resource's Chainsaw test, by file name, and
// a report of what the test covers.
func Generate(
	m *model.Model,
	crd *model.CRD,
	in *Inputs,
	opts Options,
) (map[string][]byte, *Report, error) {
	aws, err := newAWSAPI(m, crd, opts.Service, in.NotFoundCode)
	if err != nil {
		return nil, nil, err
	}
	// References point to the dependencies, so their kinds' API versions are
	// the dependencies'
	for _, dep := range in.Dependencies {
		kind, _ := dep["kind"].(string)
		apiVersion, _ := dep["apiVersion"].(string)
		if kind != "" && apiVersion != "" {
			aws.apiVersions[kind] = apiVersion
		}
	}
	specFields := map[string]bool{}
	for _, f := range crd.SpecFields {
		specFields[f.GetCRDJSONFieldName()] = true
	}
	for k := range in.Checks {
		if !specFields[k] {
			return nil, nil, fmt.Errorf("inputs check %q isn't a spec field", k)
		}
	}
	g := &generator{
		crd:        crd,
		in:         in,
		opts:       opts,
		aws:        aws,
		apiVersion: m.APIGroup() + "/" + opts.APIVersion,
		name:       strings.ReplaceAll(crd.Names.Snake, "_", "-"),
		files:      map[string][]byte{},
		report:     &Report{Kind: crd.Kind},
		checked:    map[string]bool{},
		unchecked:  map[string]bool{},
		immutable:  map[string]bool{},
		skipped:    map[string]bool{},
	}
	tests := []test{g.lifecycleTest(lifecycle{
		name: g.name, create: in.Create, updates: in.Updates, invalid: in.Invalid,
	})}
	if in.Invalid == nil {
		g.note("inputs give no invalid spec, so there's no ACK.Terminal check")
	}
	for _, v := range in.Variants {
		lc := lifecycle{
			variant: v.Name, description: v.Description,
			name: g.name + "-" + v.Name, create: v.Create, updates: v.Updates,
		}
		tests = append(tests, g.lifecycleTest(lc))
		g.report.Variants = append(g.report.Variants, g.opts.NamePrefix+lc.name)
	}
	if names := unsuffixedNames(in.allSpecs()...); len(names) != 0 {
		g.note(fmt.Sprintf(
			"%s set without %s, so the AWS names can clash with the same test running at the same time",
			strings.Join(names, ", "), SuffixToken,
		))
	}
	if crd.IsAdoptable() {
		if adoption := g.adoptionTest(); adoption != nil {
			tests = append(tests, *adoption)
		}
	} else {
		g.note("not adoptable, so there's no adoption test")
	}
	g.finishReport()

	var buf bytes.Buffer
	buf.WriteString(GeneratedHeader)
	buf.WriteString("# Edit inputs.yaml and regenerate instead.\n#\n")
	buf.WriteString(g.report.String())
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	for _, t := range tests {
		if err := enc.Encode(t); err != nil {
			return nil, nil, err
		}
	}
	if err := enc.Close(); err != nil {
		return nil, nil, err
	}
	g.files["chainsaw-test.yaml"] = buf.Bytes()
	return g.files, g.report, nil
}

type generator struct {
	crd        *model.CRD
	in         *Inputs
	opts       Options
	aws        *awsAPI
	apiVersion string
	// name is the CR's name, and the base name of its manifests
	name      string
	files     map[string][]byte
	report    *Report
	checked   map[string]bool
	unchecked map[string]bool
	// immutable and skipped are the immutable fields the tests try to
	// change, and those they can't
	immutable map[string]bool
	skipped   map[string]bool
}

// lifecycle is a create spec and the updates to it: the inputs' main one, or
// a variant.
type lifecycle struct {
	// variant is the variant's name, or empty for the main create spec.
	variant     string
	description string
	// name is the CR's name and the test's name, and the base name of the
	// manifests.
	name    string
	create  map[string]interface{}
	updates []map[string]interface{}
	invalid *InvalidInput
}

func (g *generator) lifecycleTest(lc lifecycle) test {
	setup, applyDeps := g.setupSteps()
	create := step{Name: "create"}
	if applyDeps {
		create.Try = append(create.Try, g.applyDependencies()...)
	}
	create.Try = append(create.Try, operation{
		Description: "Create the " + g.crd.Kind,
		Apply:       &fileRef{File: g.manifest(lc.name+".yaml", lc.name, lc.create, nil)},
	})
	if inv := lc.invalid; inv != nil {
		invalidName := g.name + "-invalid"
		create.Bindings = append(create.Bindings, binding{Name: "terminalMessage", Value: g.templated(inv.Message)})
		create.Try = append(create.Try, operation{
			Description: "Create a " + g.crd.Kind + " with an invalid spec",
			Apply:       &fileRef{File: g.manifest(invalidName+".yaml", invalidName, inv.Spec, nil)},
		}, operation{
			Description: "Check the invalid " + g.crd.Kind + " is Terminal",
			Assert: &resourceRef{Resource: map[string]interface{}{
				"apiVersion": g.apiVersion,
				"kind":       g.crd.Kind,
				"metadata":   map[string]interface{}{"name": invalidName},
				"status": map[string]interface{}{
					"(conditions[?type == 'ACK.Terminal'])": []interface{}{map[string]interface{}{
						"status":                                "True",
						"(contains(message, $terminalMessage))": true,
					}},
				},
			}},
		})
	}

	steps := append(setup, create, g.useStep("wait-synced", "", "assert-synced.yaml"))
	if s, ok := g.awsCheckStep("check-created", lc.create); ok {
		steps = append(steps, s)
	}
	steps = append(steps, g.immutableSteps(lc.name, lc.create)...)

	where := ""
	if lc.variant != "" {
		where = "variant " + lc.variant + "'s "
	}
	state := lc.create
	for i, upd := range lc.updates {
		state = merge(state, upd)
		n := i + 1
		file := g.manifest(fmt.Sprintf("%s-update-%d.yaml", lc.name, n), lc.name, state, nil)
		steps = append(steps, g.useStep(
			fmt.Sprintf("update-%d", n), "Change "+strings.Join(sortedKeys(upd), ", "),
			"apply-and-wait-reconciled.yaml", binding{Name: "file", Value: file},
		))
		if s, ok := g.awsCheckStep(fmt.Sprintf("check-update-%d", n), state); ok {
			steps = append(steps, s)
		}
		if changed := sortedKeys(upd); len(g.readableFields(upd)) == 0 {
			g.note(fmt.Sprintf(
				"%supdate-%d only changes fields that can't be checked in AWS (%s), so its check can't tell whether it was applied",
				where, n, strings.Join(changed, ", "),
			))
		}
	}
	steps = append(steps, g.deleteStep())

	description := fmt.Sprintf(
		"%s lifecycle: create, try to change each immutable field, apply each update, delete.",
		g.crd.Kind,
	)
	if lc.variant != "" {
		description = fmt.Sprintf("%s lifecycle, variant %s", g.crd.Kind, lc.variant)
		if lc.description != "" {
			description += ": " + strings.TrimSuffix(lc.description, ".")
		}
		description += "."
	}
	return test{
		APIVersion: "chainsaw.kyverno.io/v1alpha1",
		Kind:       "Test",
		Metadata:   objectMeta{Name: g.opts.NamePrefix + lc.name},
		Spec: testSpec{
			Description: description,
			Bindings:    g.bindings(lc.name),
			Steps:       steps,
		},
	}
}

// adoptionTest creates a resource through a CR with deletion-policy retain,
// deletes the CR, and adopts the resource it leaves behind.
func (g *generator) adoptionTest() *test {
	if g.aws.delete == "" {
		g.note("no delete operation with a usable identifier, so there's no adoption test")
		return nil
	}
	source, adopted := g.name+"-adopt-source", g.name+"-adopted"
	setup, applyDeps := g.setupSteps()
	create := step{Name: "create-source"}
	if applyDeps {
		create.Try = append(create.Try, g.applyDependencies()...)
	}
	create.Try = append(create.Try, operation{
		Description: "Create a " + g.crd.Kind + " to adopt, through a CR with deletion-policy retain",
		Apply: &fileRef{File: g.manifest(source+".yaml", source, g.in.Create, map[string]interface{}{
			"services.k8s.aws/deletion-policy": "retain",
		})},
	})
	// The adoption fields are the AWS ID, plus any the inputs say the
	// controller's adoption code needs, copied from the source CR
	fields := []string{g.aws.idJSONName + ": $sourceID"}
	record := []operation{{
		Description: "Record the AWS ID",
		Script: &script{
			Content: literal(fmt.Sprintf(
				"kubectl get %s %s -n \"$NAMESPACE\" -o jsonpath='%s'",
				g.kubectlResource(), source, g.aws.idPath,
			)),
			Outputs: []binding{{Name: "sourceID", Value: "(trim_space($stdout))"}},
		},
	}}
	if len(g.in.AdoptionFields) != 0 {
		var outputs []binding
		for _, f := range g.in.AdoptionFields {
			name := "source" + strings.ToUpper(f[:1]) + f[1:]
			outputs = append(outputs, binding{Name: name, Value: fmt.Sprintf("(json_parse($stdout).spec.%s)", f)})
			fields = append(fields, fmt.Sprintf("%s: $%s", f, name))
		}
		record = append(record, operation{
			Description: "Record the spec fields adopting it needs: " + strings.Join(g.in.AdoptionFields, ", "),
			Script: &script{
				Content: literal(fmt.Sprintf("kubectl get %s %s -n \"$NAMESPACE\" -o json", g.kubectlResource(), source)),
				Outputs: outputs,
			},
		})
	}
	adoptedFile := g.manifest(adopted+".yaml", adopted, nil, map[string]interface{}{
		"services.k8s.aws/adoption-policy": "adopt",
		"services.k8s.aws/adoption-fields": fmt.Sprintf("(to_string({%s}))", strings.Join(fields, ", ")),
	})
	adopt := step{
		Name: "adopt",
		Try: append(record, operation{
			Description: "Delete the CR, leaving the " + g.crd.Kind + " in AWS",
			Delete:      &deleteOp{Ref: objectRef{APIVersion: g.apiVersion, Kind: g.crd.Kind, Name: source}},
		}, operation{
			Description: "Adopt the " + g.crd.Kind,
			Apply:       &fileRef{File: adoptedFile},
		}),
		// The resource is only deleted with its adopting CR if adoption worked
		Cleanup: []operation{{
			Description: "Delete the " + g.crd.Kind + " if it's still in AWS",
			Script: &script{
				Env:     []binding{{Name: "AWS_ID", Value: "($sourceID)"}},
				Content: literal(g.aws.delete + " >/dev/null 2>&1 || true"),
			},
		}},
	}
	steps := append(setup,
		create,
		g.useStep("wait-synced-source", "", "assert-synced.yaml", binding{Name: "name", Value: source}),
		adopt,
		g.useStep("wait-synced", "", "assert-synced.yaml"),
	)
	if s, ok := g.awsCheckStep("check-adopted", g.in.Create); ok {
		steps = append(steps, s)
	}
	// Adoption fills the spec in from AWS before the CR is Synced, so it's
	// checked once
	if spec := g.readableFields(g.in.Create); len(spec) != 0 {
		steps = append(steps, step{
			Name: "check-adopted-spec",
			Try: []operation{{
				Description: "Check the CR's spec was filled in from AWS",
				Script: &script{
					Content: literal(fmt.Sprintf(`kubectl get %s %s -n "$NAMESPACE" -o json`, g.kubectlResource(), adopted)),
					Check: map[string]interface{}{
						"($error)": nil,
						"(json_parse($stdout))": map[string]interface{}{
							"metadata": map[string]interface{}{
								"annotations": map[string]interface{}{"services.k8s.aws/adopted": "true"},
							},
							"spec": g.templated(jsonParsed(g.adoptedRefs(spec))),
						},
					},
				},
			}},
		})
	}
	steps = append(steps, g.deleteStep())
	return &test{
		APIVersion: "chainsaw.kyverno.io/v1alpha1",
		Kind:       "Test",
		Metadata:   objectMeta{Name: g.opts.NamePrefix + g.name + "-adoption"},
		Spec: testSpec{
			Description: fmt.Sprintf(
				"Adopts a %s left behind by a CR with deletion-policy retain, then deletes it.",
				g.crd.Kind,
			),
			Bindings: g.bindings(adopted),
			Steps:    steps,
		},
	}
}

// immutableSteps try to change each immutable field the create spec sets,
// which the CRD's validation rules must reject.
func (g *generator) immutableSteps(name string, create map[string]interface{}) []step {
	var steps []step
	specFields := g.specFieldsByJSON()
	for _, k := range sortedKeys(create) {
		f, ok := specFields[k]
		if !ok || !f.IsImmutable() {
			continue
		}
		changed, ok := changedValue(create[k])
		if !ok {
			g.skipped[k] = true
			continue
		}
		g.immutable[k] = true
		steps = append(steps, step{
			Name: "immutable-" + strings.ToLower(k),
			Try: []operation{{
				Description: fmt.Sprintf("Check changing %s is rejected", k),
				Patch: &patch{
					Resource: map[string]interface{}{
						"apiVersion": g.apiVersion,
						"kind":       g.crd.Kind,
						"metadata":   map[string]interface{}{"name": name},
						"spec":       map[string]interface{}{k: g.templated(changed)},
					},
					Expect: []expectation{{Check: map[string]interface{}{"($error != null)": true}}},
				},
			}},
		})
	}
	return steps
}

// awsCheckStep returns a step checking in AWS the fields of spec that can be
// read back, or false if none can.
func (g *generator) awsCheckStep(name string, spec map[string]interface{}) (step, bool) {
	checks, checked, unchecked := g.aws.checks(g.crd, spec, g.in.Checks)
	for _, k := range unchecked {
		g.unchecked[k] = true
	}
	for _, k := range checked {
		g.checked[k] = true
	}
	if len(checks) == 0 {
		return step{}, false
	}
	return g.checkStep(name, checks), true
}

// checkStep returns a step reading the resource back from AWS once. One
// script runs the step's AWS CLI calls and prints their JSON output as one
// object, by each call's key, and Chainsaw checks it against the expected
// values, as an assertion tree. A referenced resource's AWS ID, a value and
// the suffix are expressions in the tree, read when the step runs. The step
// runs once the CR is Synced, when the controller has just read the resource
// from AWS, so the read isn't polled. The AWS CLI retries throttling and
// transient errors itself.
func (g *generator) checkStep(name string, checks []awsCheck) step {
	var content strings.Builder
	content.WriteString("export AWS_DEFAULT_OUTPUT=json AWS_RETRY_MODE=adaptive\n")
	content.WriteString("# A call that fails prints its error, and leaves jq no JSON to read\n")
	content.WriteString("jq -n \\\n")
	expected := map[string]interface{}{}
	for _, c := range checks {
		fmt.Fprintf(&content, "  --argjson %s \"$(%s)\" \\\n", c.Key, c.Command)
		expected[c.Key] = g.templated(c.Expected)
	}
	content.WriteString("  '$ARGS.named'")
	return step{
		Name: name,
		Try: []operation{{
			Description: "Read back from AWS",
			Script: &script{
				Env:     g.awsEnv(),
				Content: literal(content.String()),
				Check: map[string]interface{}{
					"($error)":              nil,
					"(json_parse($stdout))": expected,
				},
			},
		}},
		Catch: []operation{{
			Describe: &objectRef{APIVersion: g.apiVersion, Kind: g.crd.Kind, Name: "($name)"},
		}},
	}
}

// awsEnv returns the environment of an AWS CLI call reading the resource back:
// its AWS ID, the suffix, and the inputs' values.
func (g *generator) awsEnv() []binding {
	env := []binding{
		{Name: "AWS_ID", Value: fmt.Sprintf(
			"(x_k8s_get($client, '%s', '%s', $namespace, $name).%s)",
			g.apiVersion, g.crd.Kind, strings.Trim(g.aws.idPath, "{}."),
		)},
		{Name: "SUFFIX", Value: "($suffix)"},
	}
	for _, value := range g.in.Values {
		env = append(env, binding{Name: value.Name, Value: fmt.Sprintf(
			"(x_k8s_get($client, 'v1', 'ConfigMap', $namespace, '%s').data.%s)", valuesConfigMap, value.Name,
		)})
	}
	return env
}

// adoptedRefs returns spec with each reference to a single resource replaced
// by the field it fills, holding the referenced resource's AWS ID: adoption
// fills that field in from AWS, not the reference, e.g. vpcID, not vpcRef.
// A list's order in AWS isn't known, so references to several resources are
// left out.
func (g *generator) adoptedRefs(spec map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(spec))
	for k, v := range spec {
		t, ok := g.aws.refs[k]
		if !ok {
			// Adoption fills nested references in as the IDs too
			out[k] = g.aws.resolveNested(k, v)
			continue
		}
		if crs, ok := refNames(t, v); ok && !t.list {
			out[t.concrete] = "(" + g.aws.refGet(t, crs[0]) + ")"
		}
	}
	return out
}

func (g *generator) deleteStep() step {
	if g.aws.notFoundCode == "" && g.in.Deleted == "" && g.aws.waitDeleted == "" {
		g.note("neither the generator config, the API model's waiters nor the inputs give a not-found error code, a deleted waiter or a deleted query, so deletion isn't checked in AWS")
		return step{
			Name: "delete",
			Try: []operation{{
				Delete: &deleteOp{Ref: objectRef{APIVersion: g.apiVersion, Kind: g.crd.Kind, Name: "($name)"}},
			}},
		}
	}
	bindings := []binding{
		{Name: "describeCommand", Value: g.aws.describe},
		{Name: "notFoundCode", Value: g.aws.notFoundCode},
	}
	if g.in.Deleted != "" {
		bindings = append(bindings, binding{Name: "deletedQuery", Value: g.in.Deleted})
	}
	// AWS's own polling, when the API model defines it
	if g.aws.waitDeleted != "" {
		bindings = append(bindings, binding{Name: "waitCommand", Value: g.aws.waitDeleted})
	}
	return g.useStep("delete", "", "delete.yaml", bindings...)
}

// valuesConfigMap holds a test's values, in its namespace, so that every
// step can read them: Chainsaw's script outputs only last for their step.
const valuesConfigMap = "ack-test-values"

// setupSteps returns the steps before the create step, and whether the
// create step applies the dependencies itself. With values, the dependencies
// come first, in a step of their own that waits for them to be Synced,
// since a value can read their AWS IDs.
func (g *generator) setupSteps() ([]step, bool) {
	if len(g.in.Values) == 0 {
		return nil, true
	}
	var steps []step
	if len(g.in.Dependencies) != 0 {
		deps := step{Name: "dependencies", Try: g.applyDependencies()}
		for _, dep := range g.in.Dependencies {
			meta, _ := dep["metadata"].(map[string]interface{})
			deps.Try = append(deps.Try, operation{
				Description: fmt.Sprintf("Wait for %s %v to be Synced", dep["kind"], meta["name"]),
				Assert: &resourceRef{Resource: map[string]interface{}{
					"apiVersion": dep["apiVersion"],
					"kind":       dep["kind"],
					"metadata":   map[string]interface{}{"name": meta["name"]},
					"status": map[string]interface{}{
						"(conditions[?type == 'ACK.ResourceSynced'])": []interface{}{
							map[string]interface{}{"status": "True"},
						},
					},
				}},
			})
		}
		steps = append(steps, deps)
	}
	return append(steps, g.valuesStep()), false
}

// shQuote returns s as one single-quoted shell word.
func shQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// valuesStep works out the inputs' values, keeping each in valuesConfigMap
// as soon as it's known. Its cleanup runs the values' cleanups at the end of
// the test, in reverse order, for the values that were set.
func (g *generator) valuesStep() step {
	env := `export SUFFIX NAMESPACE AWS_REGION="${AWS_REGION:-${AWS_DEFAULT_REGION:-}}"` + "\n"
	var run strings.Builder
	run.WriteString("set -eu\n" + env)
	fmt.Fprintf(&run, "kubectl create configmap %s -n \"$NAMESPACE\" >/dev/null\n", valuesConfigMap)
	for _, v := range g.in.Values {
		// Each command runs in a shell of its own, so its set -e holds: a shell
		// ignores set -e in a function called before ||, and a failed step of
		// the command would go unnoticed
		fmt.Fprintf(&run, "\n%[1]s=$(sh -euc %[2]s) || { echo \"value %[1]s: command failed\" >&2; exit 1; }\n",
			v.Name, shQuote(strings.TrimSpace(v.Command)))
		fmt.Fprintf(&run, "[ -n \"$%[1]s\" ] || { echo \"value %[1]s: command printed nothing\" >&2; exit 1; }\n", v.Name)
		fmt.Fprintf(&run, "export %s\n", v.Name)
		fmt.Fprintf(&run,
			"kubectl patch configmap %s -n \"$NAMESPACE\" --type merge -p \"$(jq -n --arg v \"$%s\" '{data: {%s: $v}}')\" >/dev/null\n",
			valuesConfigMap, v.Name, v.Name,
		)
	}

	var cleanup strings.Builder
	cleanup.WriteString("set -u\n" + env)
	fmt.Fprintf(&cleanup,
		"eval \"$(kubectl get configmap %s -n \"$NAMESPACE\" -o json 2>/dev/null | jq -r '.data // {} | to_entries[] | \"export \\(.key)=\\(.value | @sh)\"')\"\n",
		valuesConfigMap,
	)
	cleanup.WriteString("failed=0\n")
	hasCleanup := false
	for i := len(g.in.Values) - 1; i >= 0; i-- {
		v := g.in.Values[i]
		if v.Cleanup == "" {
			continue
		}
		hasCleanup = true
		fmt.Fprintf(&cleanup,
			"\nif [ -n \"${%[1]s:-}\" ]; then sh -euc %[2]s || { echo \"cleanup of %[1]s failed\" >&2; failed=1; }; fi\n",
			v.Name, shQuote(strings.TrimSpace(v.Cleanup)),
		)
	}
	cleanup.WriteString("exit $failed\n")

	suffix := []binding{{Name: "SUFFIX", Value: "($suffix)"}}
	s := step{
		Name: "values",
		Try: []operation{{
			Description: "Work out the values the specs use",
			Script:      &script{Env: suffix, Content: literal(run.String())},
		}},
	}
	if hasCleanup {
		s.Cleanup = []operation{{
			Description: "Clean up what the values' commands created",
			Script:      &script{Env: suffix, Content: literal(cleanup.String())},
		}}
	}
	return s
}

func (g *generator) applyDependencies() []operation {
	var ops []operation
	for _, dep := range g.in.Dependencies {
		meta, _ := dep["metadata"].(map[string]interface{})
		depName, _ := meta["name"].(string)
		kind, _ := dep["kind"].(string)
		file := "dependency-" + depName + ".yaml"
		g.files[file] = append([]byte(GeneratedHeader), marshal(g.templated(dep))...)
		ops = append(ops, operation{
			Description: fmt.Sprintf("Create the %s %s it depends on", kind, depName),
			Apply:       &fileRef{File: file},
		})
	}
	return ops
}

// manifest writes a manifest of the resource and returns its file name.
func (g *generator) manifest(
	file, name string,
	spec map[string]interface{},
	annotations map[string]interface{},
) string {
	meta := map[string]interface{}{"name": name}
	if len(annotations) != 0 {
		meta["annotations"] = annotations
	}
	obj := map[string]interface{}{
		"apiVersion": g.apiVersion,
		"kind":       g.crd.Kind,
		"metadata":   meta,
	}
	if len(spec) != 0 {
		obj["spec"] = g.templated(spec)
	}
	g.files[file] = append([]byte(GeneratedHeader), marshal(obj)...)
	return file
}

func (g *generator) useStep(name, description, template string, bindings ...binding) step {
	s := step{
		Name:        name,
		Description: description,
		Use:         &use{Template: path.Join(g.opts.TemplatesDir, template)},
	}
	if len(bindings) != 0 {
		s.Use.With = &with{Bindings: bindings}
	}
	return s
}

func (g *generator) bindings(name string) []binding {
	return []binding{
		// A test's AWS names end in it, so tests running at the same time
		// in one account, e.g. for two PRs, don't clash
		{Name: "suffix", Value: "(random('[a-z0-9]{8}'))"},
		{Name: "apiVersion", Value: g.apiVersion},
		{Name: "kind", Value: g.crd.Kind},
		{Name: "name", Value: name},
		{Name: "idPath", Value: g.aws.idPath},
	}
}

// kubectlResource returns the CR's resource name for kubectl, e.g.
// vpc.ec2.services.k8s.aws.
func (g *generator) kubectlResource() string {
	group := strings.SplitN(g.apiVersion, "/", 2)[0]
	return strings.ToLower(g.crd.Kind) + "." + group
}

// readableFields returns the fields of spec that can be read back from AWS.
func (g *generator) readableFields(spec map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range spec {
		if g.aws.readable(k, g.in.Checks) {
			out[k] = v
		}
	}
	return out
}

// jsonParsed returns spec with each string field that holds a JSON object or
// array asserted through json_parse, by content. AWS can return a document
// such as a policy with its keys in another order, or reformatted, so the
// string an adopted CR gets isn't the one the inputs gave.
func jsonParsed(spec map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(spec))
	for k, v := range spec {
		if s, ok := v.(string); ok {
			t := strings.TrimSpace(s)
			var parsed interface{}
			if (strings.HasPrefix(t, "{") || strings.HasPrefix(t, "[")) &&
				json.Unmarshal([]byte(t), &parsed) == nil {
				out["(json_parse("+k+"))"] = parsed
				continue
			}
		}
		out[k] = v
	}
	return out
}

func (g *generator) specFieldsByJSON() map[string]*model.Field {
	out := map[string]*model.Field{}
	for _, f := range g.crd.SpecFields {
		out[f.GetCRDJSONFieldName()] = f
	}
	return out
}

func (g *generator) note(n string) {
	for _, existing := range g.report.Notes {
		if existing == n {
			return
		}
	}
	g.report.Notes = append(g.report.Notes, n)
}

func (g *generator) finishReport() {
	g.report.Checked = sortedSet(g.checked)
	for k := range g.checked {
		delete(g.unchecked, k)
	}
	for _, k := range sortedSet(g.unchecked) {
		why, ok := g.aws.why[k]
		if !ok {
			why = g.aws.readOp + " doesn't return it"
		}
		g.report.Unchecked = append(g.report.Unchecked, k+": "+why)
	}

	g.report.ImmutableTested = sortedSet(g.immutable)
	for k := range g.immutable {
		delete(g.skipped, k)
	}
	g.report.ImmutableSkipped = sortedSet(g.skipped)

	updated := map[string]bool{}
	allUpdates := append([]map[string]interface{}{}, g.in.Updates...)
	for _, v := range g.in.Variants {
		allUpdates = append(allUpdates, v.Updates...)
	}
	for _, upd := range allUpdates {
		for k := range upd {
			updated[k] = true
		}
	}
	for k, f := range g.specFieldsByJSON() {
		if f.IsImmutable() || f.IsReference() || updated[k] {
			continue
		}
		// A field set through its reference counts as updated if its
		// reference field is
		if f.HasReference() {
			if refNames, err := f.GetReferenceFieldName(); err == nil && updated[refNames.CamelLower] {
				continue
			}
		}
		g.report.NotUpdated = append(g.report.NotUpdated, k)
	}
	sort.Strings(g.report.NotUpdated)
}

// changedValue returns a different value of the same type, for checking an
// immutable field can't be changed. The CRD's validation rejects the change
// before the controller sees it, so the value needn't be valid in AWS.
func changedValue(v interface{}) (interface{}, bool) {
	switch t := v.(type) {
	case string:
		return t + "-changed", true
	case bool:
		return !t, true
	case int:
		return t + 1, true
	case float64:
		return t + 1, true
	}
	return nil, false
}

// SuffixToken, in a string in the inputs, stands for the test's random
// suffix, e.g. a SecurityGroup name of ack-test-$SUFFIX. In checks, it's an
// environment variable of the command.
const SuffixToken = "$SUFFIX"

// GeneratedHeader starts every file the generator writes. A file in a test
// directory that starts with it, and that the generator no longer writes, is
// stale.
const GeneratedHeader = "# Code generated by ack-generate chainsaw. DO NOT EDIT.\n"

// tokenRegexp matches a $NAME token: $SUFFIX, or the name of a value.
var tokenRegexp = regexp.MustCompile(`\$[A-Z][A-Z0-9_]*`)

// templated returns v with each string holding $SUFFIX or a value's $NAME
// made a Chainsaw expression putting the test's suffix or the value in its
// place. For example, ack-test-$SUFFIX becomes
//
//	(join('', ['ack-test-', $suffix]))
func (g *generator) templated(v interface{}) interface{} {
	exprs := map[string]string{SuffixToken: "$suffix"}
	for _, value := range g.in.Values {
		exprs["$"+value.Name] = fmt.Sprintf(
			"x_k8s_get($client, 'v1', 'ConfigMap', $namespace, '%s').data.%s", valuesConfigMap, value.Name,
		)
	}
	return templateWith(v, exprs)
}

// templateWith returns v with each string holding one of the tokens of exprs
// made a Chainsaw expression joining the rest of the string and the tokens'
// expressions. Other $NAMEs are left as they are.
func templateWith(v interface{}, exprs map[string]string) interface{} {
	switch t := v.(type) {
	case string:
		var parts []string
		last := 0
		literal := func(s string) {
			if s != "" {
				// In a JMESPath raw string, \\ and \' are escapes
				s = strings.ReplaceAll(s, `\`, `\\`)
				parts = append(parts, "'"+strings.ReplaceAll(s, "'", `\'`)+"'")
			}
		}
		for _, loc := range tokenRegexp.FindAllStringIndex(t, -1) {
			expr, ok := exprs[t[loc[0]:loc[1]]]
			if !ok {
				continue
			}
			literal(t[last:loc[0]])
			parts = append(parts, expr)
			last = loc[1]
		}
		if last == 0 {
			return t
		}
		literal(t[last:])
		return fmt.Sprintf("(join('', [%s]))", strings.Join(parts, ", "))
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for k, mv := range t {
			out[k] = templateWith(mv, exprs)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(t))
		for i, e := range t {
			out[i] = templateWith(e, exprs)
		}
		return out
	}
	return v
}

// unsuffixedNames returns the name-like fields of specs, e.g. name or
// groupName, set to a string without SuffixToken or a value. Their AWS names can clash
// with those of the same test running at the same time.
func unsuffixedNames(specs ...map[string]interface{}) []string {
	found := map[string]bool{}
	for _, spec := range specs {
		for k, v := range spec {
			s, ok := v.(string)
			isName := k == "name" || strings.HasSuffix(k, "Name")
			// A value's $NAME can hold the suffix itself
			if ok && isName && !tokenRegexp.MatchString(s) {
				found[k] = true
			}
		}
	}
	return sortedSet(found)
}

// merge returns base with patch merged in, the way a JSON merge patch is:
// maps merge, everything else is replaced.
func merge(base, patch map[string]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range base {
		out[k] = v
	}
	for k, v := range patch {
		pm, pok := v.(map[string]interface{})
		bm, bok := out[k].(map[string]interface{})
		if pok && bok {
			out[k] = merge(bm, pm)
		} else {
			out[k] = v
		}
	}
	return out
}

func marshal(v interface{}) []byte {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	// Encoding plain maps and slices can't fail
	_ = enc.Encode(v)
	_ = enc.Close()
	return buf.Bytes()
}

func sortedSet(s map[string]bool) []string {
	out := make([]string, 0, len(s))
	for k := range s {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
