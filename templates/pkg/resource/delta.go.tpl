{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

// Hack to avoid import errors during build...
var (
	_ = &reflect.Method{}
)

// newResourceDelta returns a new `ackcompare.Delta` used to compare two
// resources
func newResourceDelta(
	a *resource,
	b *resource,
) *ackcompare.Delta {
	delta := ackcompare.NewDelta()
	if ((a == nil && b != nil) ||
			(a != nil && b == nil)) {
		delta.Add("", a, b)
		return delta
	}

{{- if $hookCode := Hook .CRD "delta_pre_compare" }}
{{ $hookCode }}
{{- end }}
{{ GoCodeCompare .CRD "delta" "a.ko" "b.ko" 1}}
{{- if $hookCode := Hook .CRD "delta_post_compare" }}
{{ $hookCode }}
{{- end }}
	return delta
}
