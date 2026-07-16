package model_test

import (
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestReplacePkgName(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		subject         string
		pkgName         string
		replacePkgAlias string
		keepPointer     bool
		want            string
	}{
		{ // most frequent case
			"*ecr.Repository",
			"ecr",
			"svcsdk",
			true,
			"*svcsdk.Repository",
		},
		{ // don't keep pointer
			"*ecr.Repository",
			"ecr",
			"svcsdk",
			false,
			"svcsdk.Repository",
		},
		{ // non sdk type
			"*time.Time",
			"ecr",
			"svcsdk",
			true,
			"*time.Time",
		},
		{ // map type
			"map[string]*ecr.Repository",
			"ecr",
			"svcsdk",
			true,
			"map[string]*svcsdk.Repository",
		},
		{ // nested map type
			"map[string]map[string]uint8",
			"ec2",
			"svcsdk",
			true,
			"map[string]map[string]uint8",
		},
		{ // slice type
			"[]ecr.Repository",
			"ecr",
			"svcsdk",
			true,
			"[]svcsdk.Repository",
		},
		{ // nested slices type
			"[][]*codedeploy.EC2TagFilter",
			"codedeploy",
			"svcsdk",
			true,
			"[][]*svcsdk.EC2TagFilter",
		},
	}

	for _, tc := range testCases {
		result := model.ReplacePkgName(
			tc.subject,
			tc.pkgName,
			tc.replacePkgAlias,
			tc.keepPointer,
		)
		assert.Equal(tc.want, result)
	}
}

func TestCleanGoTypeSecretReference(t *testing.T) {
	gte, gt, gtwp := model.CleanGoType(
		nil,
		nil,
		&api.Shape{Type: "blob"},
		&config.FieldConfig{IsSecretReference: true},
	)

	assert.Equal(t, "SecretReference", gte)
	assert.Equal(t, "*ackv1alpha1.SecretReference", gt)
	assert.Equal(t, "*ackv1alpha1.SecretReference", gtwp)
}
