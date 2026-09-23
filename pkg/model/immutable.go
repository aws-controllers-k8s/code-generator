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
// member's absence.
//
// There are two forms, because "immutable" means something weaker for a field
// the controller itself fills in.
//
// lateInitialized == false -- freeze presence and value:
//
//	has(self.x) == has(oldSelf.x) && (!has(self.x) || self.x == oldSelf.x)
//
// The first conjunct rejects both absent->present and present->absent. The
// second compares values only when the member is present, so it never
// dereferences an absent field.
//
// lateInitialized == true -- freeze value once set, but permit the first write:
//
//	!has(oldSelf.x) || (has(self.x) && self.x == oldSelf.x)
//
// A late-initialized field is populated by the controller after creation, from
// whatever the AWS API reports. That patch is an absent->present transition and
// the strict form rejects it, leaving the resource stuck in a reconcile error.
// CEL cannot tell a controller late-init write from a user edit -- both arrive
// as an UPDATE adding the field -- so absent->present has to stay permitted for
// these fields. Change-after-set and remove-after-set are still rejected, which
// is strictly more than the old field-level rule caught (it skipped removals
// entirely).
//
// Neither form holds while the containing struct itself is absent: if the parent
// struct goes absent->present, Kubernetes skips this rule too, exactly as it
// skipped the field-level one. See the is_immutable documentation for that
// remaining gap.
//
// KNOWN LIMITATION -- the strict form is not safe for every is_immutable field,
// because the ACK reconciler itself moves spec fields in and out of the stored
// object and its patches go through admission like any other update:
//
//   - Adoption materializes the whole spec from the ReadOne response
//     (runtime reconciler, the AdoptionPolicy_Adopt branch calls
//     patchResourceMetadataAndSpec with the sdkFind output as the target), so
//     any immutable non-required field the Describe returns is an
//     absent->present transition the strict form rejects. acm Certificate
//     domainName is exactly this shape.
//   - Generated SetResource emits an unconditional `else { ko.Spec.X = nil }`
//     for read/create/update output members, so a user-set field the AWS
//     response omits is deleted from the stored spec -- a present->absent
//     transition that BOTH forms reject.
//
// late_initialize is the one cause of this that is visible in generator.yaml and
// so can be handled here. The others depend on per-field AWS response behaviour
// and cannot be determined statically, which is why freezing presence by default
// for every is_immutable field is not obviously correct. Treat this as the open
// design question it is rather than as settled.
// immutabilityCELFieldPath returns the value for the XValidation marker's
// fieldPath, which is what the API server attributes the rejection to.
//
// Moving the rule from the member to the containing struct would otherwise lose
// that attribution: with several members of one struct frozen, every rejection
// would report the struct and share an identical message, and the user would
// have no way to tell which member they touched. fieldPath restores it without
// changing the message text.
//
// Unlike the rule, this is a JSONPath and not CEL, so it uses the raw JSON name.
// Names that are not simple identifiers take the documented bracket form,
// e.g. `.['deployment-type']`.
func immutabilityCELFieldPath(jsonName string) string {
	if isSimpleJSONPathIdent(jsonName) {
		return "." + jsonName
	}
	return fmt.Sprintf(".['%s']", jsonName)
}

// isSimpleJSONPathIdent reports whether name can be used with dot notation in a
// fieldPath.
func isSimpleJSONPathIdent(name string) bool {
	if name == "" {
		return false
	}
	for i, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r == '_':
		case i > 0 && r >= '0' && r <= '9':
		default:
			return false
		}
	}
	return true
}

func immutabilityCELRule(jsonName string, lateInitialized bool) string {
	n := escapeCELPropertyName(jsonName)
	if lateInitialized {
		return fmt.Sprintf(
			"!has(oldSelf.%[1]s) || (has(self.%[1]s) && self.%[1]s == oldSelf.%[1]s)",
			n,
		)
	}
	return fmt.Sprintf(
		"has(self.%[1]s) == has(oldSelf.%[1]s) && (!has(self.%[1]s) || self.%[1]s == oldSelf.%[1]s)",
		n,
	)
}
