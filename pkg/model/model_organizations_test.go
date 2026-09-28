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

package model

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/pkg/names"
)

func TestOrganizations_AccountOperationSpecificRenames(t *testing.T) {
	cfg, err := ackgenconfig.New(
		filepath.Join("..", "testdata", "models", "apis", "organizations", "0000-00-00", "generator.yaml"),
		ackgenconfig.Config{},
	)
	require.NoError(t, err)
	create := &awssdkmodel.Operation{ExportedName: "CreateAccount"}
	describe := &awssdkmodel.Operation{ExportedName: "DescribeAccount"}
	crd := &CRD{
		cfg:          &cfg,
		Names:        names.New("Account"),
		Ops:          Ops{Create: create, ReadOne: describe},
		SpecFields:   make(map[string]*Field),
		StatusFields: make(map[string]*Field),
		Fields:       make(map[string]*Field),
	}

	// Building the Account model must allow the same top-level SDK field to
	// have different renames in different API operations: CreateAccount's Id
	// becomes CreateAccountRequestId, while DescribeAccount's Id becomes AccountID.
	err = crd.AddSpecField(names.New("Name"), &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "string"},
	})
	require.NoError(t, err)
	assert.Contains(t, crd.SpecFields, "Name")
	err = crd.AddStatusField(names.New("AccountID"), &awssdkmodel.ShapeRef{
		Shape: &awssdkmodel.Shape{Type: "string"},
	})
	require.NoError(t, err)
	assert.Contains(t, crd.StatusFields, "AccountID")

	assert.Equal(t, map[string]string{
		"AccountName": "Name",
		"Id":          "CreateAccountRequestId",
	}, cfg.GetAllRenames("Account", map[string]*awssdkmodel.Operation{
		create.ExportedName: create,
	}))
	assert.Equal(t, map[string]string{
		"Id":     "AccountID",
		"Status": "State",
	}, cfg.GetAllRenames("Account", map[string]*awssdkmodel.Operation{
		describe.ExportedName: describe,
	}))
	assert.Empty(t, mergedFieldRenames(crd))
}
