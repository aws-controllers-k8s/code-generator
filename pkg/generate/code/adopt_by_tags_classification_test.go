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

package code_test

import (
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// adoptionClass is the tag-based-adoption classification of a resource, derived
// from the code IdentifierFieldsFromARN generates for it.
type adoptionClass string

const (
	// classPositional: single identifier field derived from the ARN's resource
	// segment (no config needed).
	classPositional adoptionClass = "positional"
	// classARNPrimary: the whole ARN is the identifier (is_arn_primary_key).
	classARNPrimary adoptionClass = "arn-primary"
	// classRequiresTemplate: supports tags but has multiple identifier fields,
	// so it needs an adoption.arn_identifier_template override; without one it
	// emits a terminal error at runtime.
	classRequiresTemplate adoptionClass = "requires-template"
	// classUnsupportedNoTags: the resource has no tag field, so it cannot be
	// matched by a tag selector; tag-based adoption is unsupported (terminal).
	classUnsupportedNoTags adoptionClass = "unsupported-no-tags"
	// classUnsupportedNoID: the resource supports tags, but its ReadOne
	// identifier cannot be derived from the model, so it is unsupported until a
	// template/hook is added (terminal).
	classUnsupportedNoID adoptionClass = "unsupported-no-id"
)

// classifyAdoptionByTags returns the adoptionClass for a CRD by inspecting the
// generated IdentifierFieldsFromARN body. It fails the test if code generation
// errors - generation must never fail for any resource.
func classifyAdoptionByTags(t *testing.T, crd *model.CRD) adoptionClass {
	t.Helper()
	body, err := code.IdentifierFieldsFromARN(crd.Config(), crd, "arn", 1)
	require.NoError(t, err, "IdentifierFieldsFromARN must never fail codegen for %s", crd.Names.Original)
	require.NotEmpty(t, body, crd.Names.Original)
	switch {
	case strings.Contains(body, "IdentifierFieldsFromARNPositional"):
		return classPositional
	case strings.Contains(body, `map[string]string{"arn"`):
		return classARNPrimary
	case strings.Contains(body, "requires an adoption.arn_identifier_template"):
		return classRequiresTemplate
	case strings.Contains(body, "does not support tag-based adoption"):
		return classUnsupportedNoTags
	case strings.Contains(body, "no derivable identifier"):
		return classUnsupportedNoID
	default:
		t.Fatalf("%s: unrecognized IdentifierFieldsFromARN body:\n%s", crd.Names.Original, body)
		return ""
	}
}

// TestAdoptByTags_Classification documents and asserts, per resource, how each
// is handled by tag-based adoption. It is deliberately explicit (rather than a
// pass/fail counter) so it is clear which services/resources support the
// feature, which need an arn_identifier_template, and which are unsupported and
// therefore return a terminal error.
//
// The expected ResourceTypeFilter is also asserted: it is empty exactly for the
// "unsupported-no-tags" class and non-empty otherwise.
func TestAdoptByTags_Classification(t *testing.T) {
	cases := []struct {
		service  string
		resource string
		expected adoptionClass
	}{
		// Single identifier derived positionally from the ARN - zero config.
		{"eks", "Cluster", classPositional},
		{"sqs", "Queue", classPositional},
		{"lambda", "Function", classPositional},
		{"iam", "Role", classPositional},
		{"rds", "DBInstance", classPositional},
		{"ecr", "Repository", classPositional},
		{"dynamodb", "Table", classPositional},

		// ARN is itself the identifier (is_arn_primary_key) + tags supported.
		{"sns", "Topic", classARNPrimary},

		// Multiple identifier fields: require an arn_identifier_template because
		// ARN segment order is not knowable from the model.
		{"eks", "Nodegroup", classRequiresTemplate},
		{"eks", "Addon", classRequiresTemplate},
		{"eks", "FargateProfile", classRequiresTemplate},
		{"apigatewayv2", "Stage", classRequiresTemplate},
		{"wafv2", "WebACL", classRequiresTemplate},

		// No tag field in the CRD -> cannot be matched by tags -> unsupported.
		{"sns", "PlatformApplication", classUnsupportedNoTags},
		{"lambda", "Alias", classUnsupportedNoTags},
		{"s3", "Session", classUnsupportedNoTags},
		{"sagemaker", "ModelPackage", classUnsupportedNoTags},

		// Tags supported but ReadOne identifier not derivable from the model ->
		// unsupported until a template/hook is added.
		{"iam", "User", classUnsupportedNoID},
	}

	for _, tc := range cases {
		t.Run(tc.service+"/"+tc.resource, func(t *testing.T) {
			g := testutil.NewModelForService(t, tc.service)
			crd := testutil.GetCRDByName(t, g, tc.resource)
			require.NotNil(t, crd, "%s/%s not found", tc.service, tc.resource)

			got := classifyAdoptionByTags(t, crd)
			assert.Equal(t, tc.expected, got, "%s/%s classification", tc.service, tc.resource)

			// The type filter is empty exactly when the kind has no tag support.
			if tc.expected == classUnsupportedNoTags {
				assert.Equal(t, "", crd.ResourceTypeFilter(),
					"%s/%s must have empty ResourceTypeFilter", tc.service, tc.resource)
			} else {
				assert.NotEqual(t, "", crd.ResourceTypeFilter(),
					"%s/%s must have a ResourceTypeFilter", tc.service, tc.resource)
			}
		})
	}
}

// TestAdoptByTags_CodegenNeverFails is the safety-net invariant behind the
// classification above: across a broad set of services, IdentifierFieldsFromARN
// must never fail code generation for ANY resource - unsupported resources emit
// a terminal error body instead. This guards against a resource shape that
// would otherwise break `make build-controller`.
func TestAdoptByTags_CodegenNeverFails(t *testing.T) {
	services := []string{
		"eks", "ec2", "s3", "iam", "rds", "lambda", "dynamodb", "sns", "sqs",
		"ecr", "elasticache", "sagemaker", "apigatewayv2", "mq", "memorydb",
		"opensearchserverless", "wafv2", "eventbridge", "route53",
	}
	for _, svc := range services {
		g := testutil.NewModelForService(t, svc)
		crds, err := g.GetCRDs()
		require.NoError(t, err, svc)
		for _, crd := range crds {
			// classifyAdoptionByTags requires codegen to succeed and the body to
			// match one of the known classes; that is the invariant we want.
			_ = classifyAdoptionByTags(t, crd)
		}
	}
}
