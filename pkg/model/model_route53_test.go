// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	 http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestRoute53_RecordSet(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "route53")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("RecordSet", crds)
	require.NotNil(crd)

	assert.Equal("RecordSet", crd.Names.Camel)
	assert.Equal("recordSet", crd.Names.CamelLower)
	assert.Equal("record_set", crd.Names.Snake)

	// The Route53 API has CD as one operation +L:
	//
	// * ChangeResourceRecordSets
	// * ChangeResourceRecordSets
	// * ListResourceRecordSets
	require.NotNil(crd.Ops)

	assert.NotNil(crd.Ops.Create)
	assert.NotNil(crd.Ops.ReadMany)
	assert.NotNil(crd.Ops.Delete)

	// But sadly, has no Update or ReadOne operation :(
	// for update we still use ChangeResourceRecordSets
	assert.Nil(crd.Ops.ReadOne)
	assert.Nil(crd.Ops.Update)

	specFields := crd.SpecFields
	statusFields := crd.StatusFields

	expSpecFieldCamel := []string{
		"AliasTarget",
		// "ChangeBatch", <= Testing that this is removed from spec
		"CIDRRoutingConfig",
		"Failover",
		"GeoLocation",
		"HealthCheckID",
		"HostedZoneID",
		"HostedZoneRef",
		"MultiValueAnswer",
		"Name",
		"RecordType",
		"Region",
		"ResourceRecords",
		"SetIdentifier",
		"TTL",
		"Weight",
	}
	assert.Equal(expSpecFieldCamel, attrCamelNames(specFields))

	expStatusFieldCamel := []string{
		"ID",
		"Status",
		"SubmittedAt",
	}
	assert.Equal(expStatusFieldCamel, attrCamelNames(statusFields))
}
