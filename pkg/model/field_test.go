package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

func TestParentFieldPath(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		subject string
		want    string
	}{
		{
			"Repository.Name",
			"Repository",
		},
		{
			"Users..Password",
			"Users",
		},
		{
			"User.Credentials..Password",
			"User.Credentials",
		},
	}

	for _, tc := range testCases {
		result := model.ParentFieldPath(
			tc.subject,
		)
		assert.Equal(tc.want, result)
	}
}
