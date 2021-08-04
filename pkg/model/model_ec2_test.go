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

func TestEC2_LaunchTemplate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("LaunchTemplate", crds)
	require.NotNil(crd)

	assert.Equal("LaunchTemplate", crd.Names.Camel)
	assert.Equal("launchTemplate", crd.Names.CamelLower)
	assert.Equal("launch_template", crd.Names.Snake)

	// The DescribeLaunchTemplatesResult shape has no defined error codes (in
	// fact, none of the EC2 API shapes do). We will need to create exceptions
	// config in the generate.yaml for EC2, but this will take quite some
	// manual work. For now, return UNKNOWN
	assert.Equal("UNKNOWN", crd.ExceptionCode(404))

	specFields := crd.SpecFields
	statusFields := crd.StatusFields

	expSpecFieldCamel := []string{
		// TODO(jaypipes): DryRun and ClientToken are examples of two fields in
		// the resource input shape that need to be stripped out of the CRD. We
		// need to instruct the code generator that these types of fields are
		// not germane to the resource itself...
		"ClientToken",
		"DryRun",
		"LaunchTemplateData",
		"LaunchTemplateName",
		// TODO(jaypipes): Here's an example of where we need to instruct the
		// code generator to rename the "TagSpecifications" field to simply
		// "Tags" and place it into the common Spec.Tags field.
		"TagSpecifications",
		"VersionDescription",
	}
	assert.Equal(expSpecFieldCamel, attrCamelNames(specFields))

	expStatusFieldCamel := []string{
		"CreateTime",
		"CreatedBy",
		"DefaultVersionNumber",
		"LatestVersionNumber",
		// TODO(jaypipes): Handle "Id" Fields like "LaunchTemplateId" in the
		// same way as we handle ARN-ified modern service APIs and use the
		// SDKMapper to instruct the code generator that this field represents
		// the primary resource object's identifier field.
		"LaunchTemplateID",
		// LaunchTemplateName excluded because it matches input shape.,
		// TODO(jaypipes): Tags field should be excluded because it is the same
		// as the input shape's "TagSpecifications" field...
		"Tags",
	}
	assert.Equal(expStatusFieldCamel, attrCamelNames(statusFields))

	// The EC2 LaunchTemplate API has a "normal" set of CUD operations:
	//
	// * CreateLaunchTemplate
	// * ModifyLaunchTemplate
	// * DeleteLaunchTemplate
	require.NotNil(crd.Ops)

	assert.NotNil(crd.Ops.Create)
	assert.NotNil(crd.Ops.Delete)
	assert.NotNil(crd.Ops.Update)

	// However, oddly, there is no ReadOne operation. There is only a
	// ReadMany/List operation "DescribeLaunchTemplates" :(
	//
	// TODO(jaypipes): Develop strategy for informing the code generator via
	// the SDKMapper that certain APIs don't have ReadOne but only ReadMany
	// APIs...
	assert.Nil(crd.Ops.ReadOne)
	assert.NotNil(crd.Ops.ReadMany)
}

func TestEC2_Volume(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForService(t, "ec2")

	crds, err := g.GetCRDs()
	require.Nil(err)

	crd := getCRDByName("Volume", crds)
	require.NotNil(crd)

	assert.Equal("Volume", crd.Names.Camel)
	assert.Equal("volume", crd.Names.CamelLower)
	assert.Equal("volume", crd.Names.Snake)

	specFields := crd.SpecFields
	statusFields := crd.StatusFields

	expSpecFieldCamel := []string{
		"AvailabilityZone",
		"DryRun",
		"Encrypted",
		"IOPS",
		"KMSKeyID",
		"MultiAttachEnabled",
		"OutpostARN",
		"Size",
		"SnapshotID",
		"TagSpecifications",
		"VolumeType",
	}
	assert.Equal(expSpecFieldCamel, attrCamelNames(specFields))

	expStatusFieldCamel := []string{
		"Attachments",
		"CreateTime",
		"FastRestored",
		"State",
		"Tags",
		"VolumeID",
	}
	assert.Equal(expStatusFieldCamel, attrCamelNames(statusFields))

	// Ensure that we generate TypeDefs for the VolumeAttachment field.
	// This field is the payload of the `AttachVolume` payload, but should
	// be included because it is the field value for the `attachments` status
	// field
	assert.NotNil(testutil.GetTypeDefByName(t, g, "VolumeAttachment"))
}
