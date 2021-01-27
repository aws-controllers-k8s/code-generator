module github.com/aws-controllers-k8s/code-generator

go 1.14

require (
	github.com/aws/aws-sdk-go v1.35.5
	github.com/dlclark/regexp2 v1.4.0
	// pin to v0.1.1 due to release problem with v0.1.2
	github.com/gertd/go-pluralize v0.1.1
	github.com/ghodss/yaml v1.0.0
	github.com/iancoleman/strcase v0.1.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/mod v0.4.1
	k8s.io/apimachinery v0.20.2
)
