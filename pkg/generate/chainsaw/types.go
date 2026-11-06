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

import "gopkg.in/yaml.v3"

// The subset of Chainsaw's v1alpha1 Test API that generated tests use. Field
// order here is the order in the generated YAML.

type test struct {
	APIVersion string     `yaml:"apiVersion"`
	Kind       string     `yaml:"kind"`
	Metadata   objectMeta `yaml:"metadata"`
	Spec       testSpec   `yaml:"spec"`
}

type objectMeta struct {
	Name string `yaml:"name"`
}

type testSpec struct {
	Description string    `yaml:"description,omitempty"`
	Bindings    []binding `yaml:"bindings,omitempty"`
	Steps       []step    `yaml:"steps"`
}

type binding struct {
	Name  string      `yaml:"name"`
	Value interface{} `yaml:"value"`
}

type step struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description,omitempty"`
	Bindings    []binding   `yaml:"bindings,omitempty"`
	Use         *use        `yaml:"use,omitempty"`
	Try         []operation `yaml:"try,omitempty"`
	Catch       []operation `yaml:"catch,omitempty"`
	Finally     []operation `yaml:"finally,omitempty"`
	Cleanup     []operation `yaml:"cleanup,omitempty"`
}

type use struct {
	Template string `yaml:"template"`
	With     *with  `yaml:"with,omitempty"`
}

type with struct {
	Bindings []binding `yaml:"bindings"`
}

type operation struct {
	Description string       `yaml:"description,omitempty"`
	Apply       *fileRef     `yaml:"apply,omitempty"`
	Assert      *resourceRef `yaml:"assert,omitempty"`
	Patch       *patch       `yaml:"patch,omitempty"`
	Delete      *deleteOp    `yaml:"delete,omitempty"`
	Describe    *objectRef   `yaml:"describe,omitempty"`
	Script      *script      `yaml:"script,omitempty"`
}

type fileRef struct {
	File string `yaml:"file"`
}

type resourceRef struct {
	Resource interface{} `yaml:"resource"`
}

type patch struct {
	Resource map[string]interface{} `yaml:"resource"`
	Expect   []expectation          `yaml:"expect,omitempty"`
}

type expectation struct {
	Check map[string]interface{} `yaml:"check"`
}

type deleteOp struct {
	Ref objectRef `yaml:"ref"`
}

type objectRef struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Name       string `yaml:"name"`
}

type script struct {
	Env     []binding              `yaml:"env,omitempty"`
	Content literal                `yaml:"content"`
	Check   map[string]interface{} `yaml:"check,omitempty"`
	Outputs []binding              `yaml:"outputs,omitempty"`
}

// awsCheck is an AWS CLI call reading fields back, and the JSON it must
// print, as a Chainsaw assertion tree.
type awsCheck struct {
	// Key names the call's output in a check step's combined output: the
	// field it reads, or the read operation for the fields it returns, e.g.
	// enableDNSSupport or describeVpcs.
	Key      string
	Command  string
	Expected interface{}
}

// literal is a string written as a YAML literal block, so shell commands
// don't need YAML quoting on top of their own.
type literal string

func (l literal) MarshalYAML() (interface{}, error) {
	return &yaml.Node{Kind: yaml.ScalarNode, Style: yaml.LiteralStyle, Value: string(l)}, nil
}
