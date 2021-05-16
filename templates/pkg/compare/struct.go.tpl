{{ template "boilerplate" }}

package compare

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"

	svcapitypes "github.com/aws-controllers-k8s/{{ .ServiceIDClean }}-controller/apis/v1alpha1"
)

{{- range $typeDef := .TypeDefs }}

// IsEqual{{ .Names.Camel }} compares two {{ .Names.Camel }} object pointers. 
func IsEqual{{ .Names.Camel }}(a, b *svcapitypes.{{ .Names.Camel -}}) bool {
{{ GoCodeIsEqual $typeDef "a" "b" 1 }}
	return true
}

{{- end -}}