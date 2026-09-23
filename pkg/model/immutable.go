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

package model

import (
	"fmt"
	"strings"
)

// celReservedSymbols is the set of words CEL reserves and therefore cannot be
// used as a bare identifier in an expression. A CRD property whose JSON name is
// one of these must be referenced through Kubernetes' property escaping scheme
// (`__<word>__`) instead of plain field selection.
//
// This mirrors celReservedSymbols in k8s.io/apiserver/pkg/cel. It is duplicated
// rather than imported because the code-generator does not otherwise depend on
// k8s.io/apiserver.
var celReservedSymbols = map[string]struct{}{
	"true": {}, "false": {}, "null": {}, "in": {},
	"as": {}, "break": {}, "const": {}, "continue": {}, "else": {},
	"for": {}, "function": {}, "if": {}, "import": {}, "let": {},
	"loop": {}, "package": {}, "namespace": {}, "return": {},
	"var": {}, "void": {}, "while": {},
}

// celPropertyEscapes are the character sequences Kubernetes escapes when
// turning an OpenAPI property name into a CEL identifier. `__` is listed first
// because it is the escape marker itself. Mirrors celEscapeSet in
// k8s.io/apiserver/pkg/cel.
var celPropertyEscapes = strings.NewReplacer(
	"__", "__underscores__",
	".", "__dot__",
	"-", "__dash__",
	"/", "__slash__",
)

// escapeCELPropertyName returns the CEL identifier Kubernetes uses to address
// the CRD property named jsonName.
//
// Most ACK field names are plain lower-camel identifiers and pass through
// untouched, but a handful collide with CEL reserved words -- `spec.namespace`
// on the s3 and s3tables Bucket/Table resources, for example -- and would
// otherwise produce a rule that fails to compile when the CRD is created.
func escapeCELPropertyName(jsonName string) string {
	if _, found := celReservedSymbols[jsonName]; found {
		return "__" + jsonName + "__"
	}
	return celPropertyEscapes.Replace(jsonName)
}

// immutabilityCELRule returns the CEL expression that freezes the member named
// jsonName of the struct the rule is attached to.
//
// The rule is deliberately attached to the *containing* struct rather than to
// the member itself. Kubernetes skips a transition rule whenever the old value
// is absent, so a field-level `self == oldSelf` rule never fires for an
// optional field that was unset at creation -- the field could be added later
// and silently accepted. Evaluating from the parent lets the rule observe the
// member's absence, so it can freeze presence as well as value:
//
//	has(self.x) == has(oldSelf.x) && (!has(self.x) || self.x == oldSelf.x)
//
// The first conjunct rejects both absent->present and present->absent. The
// second compares values only when the member is present, so it never
// dereferences an absent field.
//
// Note that this only holds while the containing struct itself is present in
// both the old and the new object: if the parent struct goes absent->present,
// Kubernetes skips this rule too, exactly as it skipped the field-level one.
// See the package documentation on is_immutable for that remaining gap.
func immutabilityCELRule(jsonName string) string {
	n := escapeCELPropertyName(jsonName)
	return fmt.Sprintf(
		"has(self.%[1]s) == has(oldSelf.%[1]s) && (!has(self.%[1]s) || self.%[1]s == oldSelf.%[1]s)",
		n,
	)
}
