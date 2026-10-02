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

package code

import (
	"testing"

	"github.com/stretchr/testify/assert"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
)

// TestClearedDefaultBooleanPointerGeneration locks down the code-emission
// behavior selected by the <nil> sentinel that the SDK RemoveDefaults cleanup
// assigns to modeled default values whose generated SDK fields remain
// pointers. NetworkManager's VpcOptions booleans use this path.
func TestClearedDefaultBooleanPointerGeneration(t *testing.T) {
	shapeRef := &awssdkmodel.ShapeRef{
		DefaultValue: "<nil>",
		Shape: &awssdkmodel.Shape{
			Type: "boolean",
		},
	}

	t.Run("SDK to Kubernetes uses direct pointer assignment", func(t *testing.T) {
		got := setResourceForScalar(
			"ko.Spec.Options.DNSSupport",
			"resp.VpcAttachment.Options.DnsSupport",
			shapeRef,
			1,
			false,
			false,
		)

		assert.Equal(t,
			"\tko.Spec.Options.DNSSupport = resp.VpcAttachment.Options.DnsSupport\n",
			got,
		)
		assert.NotContains(t, got, "&resp.VpcAttachment.Options.DnsSupport")
	})

	t.Run("Kubernetes to SDK uses direct pointer assignment", func(t *testing.T) {
		got := setSDKForScalar(
			"DnsSupport",
			"f2",
			"structure",
			"Spec.Options.DNSSupport",
			"r.ko.Spec.Options.DNSSupport",
			false,
			shapeRef,
			1,
		)

		assert.Equal(t,
			"\tf2.DnsSupport = r.ko.Spec.Options.DNSSupport\n",
			got,
		)
		assert.NotContains(t, got, "*r.ko.Spec.Options.DNSSupport")
	})
}
