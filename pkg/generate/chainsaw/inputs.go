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
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// InputsFileName is the name of the file, in each resource's test directory,
// holding the values the generator can't derive from the resource's model.
const InputsFileName = "inputs.yaml"

// Inputs are the per-resource test values the generator can't derive from the
// resource's model: valid field values, and the other resources it needs.
type Inputs struct {
	// Dependencies are manifests of other resources the resource needs, e.g.
	// the VPC a SecurityGroup belongs to. They are created alongside it.
	Dependencies []map[string]interface{} `yaml:"dependencies,omitempty"`
	// Create is the spec of the resource to create.
	Create map[string]interface{} `yaml:"create"`
	// Updates are applied in order. Each is merged into the spec before it,
	// the way a merge patch is: maps merge, lists and values are replaced.
	Updates []map[string]interface{} `yaml:"updates,omitempty"`
	// Variants are other create specs, each with a lifecycle test of its own,
	// for fields that conflict with Create's, or that only matter at create.
	// They aren't merged with Create.
	Variants []Variant `yaml:"variants,omitempty"`
	// Invalid, if set, is a spec that must make the resource ACK.Terminal.
	Invalid *InvalidInput `yaml:"invalid,omitempty"`
	// Values are worked out when the test runs, after the dependencies are
	// Synced and before the resource is created, by AWS CLI or kubectl
	// commands: lookups such as an AMI ID, a dependency's AWS ID, or AWS
	// resources outside the controller, such as an S3 bucket. A value is $NAME
	// in any string of the inputs.
	Values []Value `yaml:"values,omitempty"`
	// NotFoundCode is the error code the read operation fails with once the
	// resource is gone, for when the generator config doesn't declare one.
	NotFoundCode string `yaml:"notFoundCode,omitempty"`
	// Deleted is a JMESPath expression over the read operation's output that
	// is true once the resource is gone, for resources AWS keeps describing
	// after deletion, e.g. "Reservations[0].Instances[0].State.Name ==
	// 'terminated'".
	Deleted string `yaml:"deleted,omitempty"`
	// AdoptionFields are spec fields, by JSON name, that adopting the
	// resource needs besides its AWS ID, for a controller whose adoption code
	// requires them, e.g. FlowLog's resourceID and resourceType. The adoption
	// test copies their values from the CR it adopts from. Each must be a
	// string field that create sets.
	AdoptionFields stringList `yaml:"adoptionFields,omitempty"`
	// Checks read spec fields back from AWS where the generator can't work
	// out how, by the fields' JSON names. A check here replaces any the
	// generator derives.
	Checks map[string]Check `yaml:"checks,omitempty"`
}

// Check reads a spec field back from AWS.
type Check struct {
	// Path is where the field is in the resource, as the read operation
	// returns it, in AWS's names, e.g. Attachments[0].VpcId. The expected
	// value is the spec's, converted to AWS's names and shape. Use it, or
	// Command.
	Path string `yaml:"path,omitempty"`
	// Where keeps the items of the list at Path whose members have these
	// values, e.g. {CidrBlockState.State: associated}.
	Where map[string]interface{} `yaml:"where,omitempty"`
	// Exclude leaves out the items of the list at Path whose members have
	// these values, e.g. AWS's own {GatewayId: local} route.
	Exclude map[string]interface{} `yaml:"exclude,omitempty"`
	// SortBy sorts the list at Path by these members, e.g. DestinationCidrBlock,
	// or [Egress, RuleNumber]. The expected list is sorted the same way.
	SortBy stringList `yaml:"sortBy,omitempty"`
	// Select reads one member of each item of the list at Path, e.g.
	// CidrBlock.
	Select string `yaml:"select,omitempty"`
	// Command is an AWS CLI call printing the field's value as JSON, in the
	// form the CR holds it, for a field the resource's read operation doesn't
	// return. It can use $AWS_ID, e.g.
	// `aws ec2 describe-vpcs --vpc-ids $AWS_ID --query 'Vpcs[0].CidrBlock'`
	Command string `yaml:"command,omitempty"`
	// Unordered is for a list of plain values whose order AWS doesn't keep.
	// The expected list is sorted, so a Command must print it sorted too.
	Unordered bool `yaml:"unordered,omitempty"`
	// When, if set, limits the check to specs setting the field to this
	// value. It's for fields that only act one way, e.g. VPC's
	// disallowSecurityGroupDefaultRules, whose false leaves AWS as it is.
	When interface{} `yaml:"when,omitempty"`
}

// allSpecs returns every spec the inputs set: the creates, the updates, the
// invalid spec, and the dependencies'.
func (in *Inputs) allSpecs() []map[string]interface{} {
	specs := append([]map[string]interface{}{in.Create}, in.Updates...)
	for _, v := range in.Variants {
		specs = append(specs, v.Create)
		specs = append(specs, v.Updates...)
	}
	if in.Invalid != nil {
		specs = append(specs, in.Invalid.Spec)
	}
	for _, dep := range in.Dependencies {
		if spec, ok := dep["spec"].(map[string]interface{}); ok {
			specs = append(specs, spec)
		}
	}
	return specs
}

// suffixNames ends every spec's name with -$SUFFIX, unless it already has
// $SUFFIX or a value in it. By ACK convention a spec's name is the resource's
// AWS name, and some must be unique in the account, e.g. a LaunchTemplate's,
// so the same test running twice at once, e.g. for two PRs, can't clash.
func (in *Inputs) suffixNames() {
	for _, spec := range in.allSpecs() {
		if s, ok := spec["name"].(string); ok && s != "" && !tokenRegexp.MatchString(s) {
			spec["name"] = s + "-" + SuffixToken
		}
	}
}

// Value is worked out by a command when the test runs.
type Value struct {
	// Name is what the inputs use the value as, $NAME. Uppercase letters,
	// digits and _.
	Name string `yaml:"name"`
	// Command prints the value. It runs with sh -c, and can use $SUFFIX,
	// $NAMESPACE, $AWS_REGION and the values before it.
	Command string `yaml:"command"`
	// Cleanup, if set, runs at the end of the test, even after a failure,
	// before the dependencies are deleted. It's for a value that's an AWS
	// resource the command created, e.g. deleting a bucket. Cleanups run in
	// the reverse order of the values.
	Cleanup string `yaml:"cleanup,omitempty"`
}

// valueNameRegexp matches names that are valid environment variables, and
// aren't mistaken for a longer name in a string.
var valueNameRegexp = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// reservedValueNames are set by the tests or the templates.
var reservedValueNames = map[string]bool{
	"SUFFIX": true, "NAMESPACE": true, "AWS_ID": true, "AWS_REGION": true,
}

// Variant is a create spec other than the main one, e.g. a VPC created with
// disallowSecurityGroupDefaultRules, or with only its required fields.
type Variant struct {
	// Name is the variant's test name and CR name after the resource's, e.g.
	// disallowed for vpc-disallowed.
	Name string `yaml:"name"`
	// Description says what the variant covers.
	Description string `yaml:"description,omitempty"`
	// Create is the spec of the resource to create.
	Create map[string]interface{} `yaml:"create"`
	// Updates are applied in order, as the main ones are.
	Updates []map[string]interface{} `yaml:"updates,omitempty"`
}

// variantNameRegexp matches names that are valid in a Kubernetes object name.
var variantNameRegexp = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

// reservedVariantNames are used by the other generated tests and manifests.
var reservedVariantNames = map[string]bool{
	"adoption": true, "adopted": true, "adopt-source": true, "invalid": true,
}

// InvalidInput is a spec that must make the resource ACK.Terminal.
type InvalidInput struct {
	Spec map[string]interface{} `yaml:"spec"`
	// Message is a substring of the expected ACK.Terminal message.
	Message string `yaml:"message"`
}

// LoadInputs reads the inputs file at path.
func LoadInputs(path string) (*Inputs, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var in Inputs
	if err := yaml.Unmarshal(b, &in); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	if len(in.Create) == 0 {
		return nil, fmt.Errorf("%s: create must set at least one field", path)
	}
	seen := map[string]bool{}
	for i, v := range in.Variants {
		switch {
		case !variantNameRegexp.MatchString(v.Name):
			return nil, fmt.Errorf("%s: variants[%d] needs a name of lowercase letters, digits and -", path, i)
		case reservedVariantNames[v.Name] || strings.HasPrefix(v.Name, "update-"):
			return nil, fmt.Errorf("%s: variant name %q clashes with a generated test or manifest", path, v.Name)
		case seen[v.Name]:
			return nil, fmt.Errorf("%s: more than one variant is named %q", path, v.Name)
		case len(v.Create) == 0:
			return nil, fmt.Errorf("%s: variant %q's create must set at least one field", path, v.Name)
		}
		seen[v.Name] = true
	}
	seenValues := map[string]bool{}
	for i, v := range in.Values {
		switch {
		case !valueNameRegexp.MatchString(v.Name):
			return nil, fmt.Errorf("%s: values[%d] needs a name of uppercase letters, digits and _", path, i)
		case reservedValueNames[v.Name]:
			return nil, fmt.Errorf("%s: value name %s is reserved", path, v.Name)
		case seenValues[v.Name]:
			return nil, fmt.Errorf("%s: more than one value is named %s", path, v.Name)
		case v.Command == "":
			return nil, fmt.Errorf("%s: value %s has no command", path, v.Name)
		}
		seenValues[v.Name] = true
	}
	for field, c := range in.Checks {
		switch {
		case (c.Command == "") == (c.Path == ""):
			return nil, fmt.Errorf("%s: checks.%s needs one of path and command", path, field)
		case c.Path == "" && (len(c.Where) != 0 || len(c.Exclude) != 0 || len(c.SortBy) != 0 || c.Select != ""):
			return nil, fmt.Errorf("%s: checks.%s: where, exclude, sortBy and select only go with path", path, field)
		}
	}
	for _, field := range in.AdoptionFields {
		if _, ok := in.Create[field].(string); !ok {
			return nil, fmt.Errorf("%s: adoptionFields: create must set %s to a string", path, field)
		}
	}
	in.suffixNames()
	return &in, nil
}

// stringList is a list of strings that can also be written as one string.
type stringList []string

func (l *stringList) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		*l = stringList{value.Value}
		return nil
	}
	var many []string
	if err := value.Decode(&many); err != nil {
		return err
	}
	*l = many
	return nil
}
