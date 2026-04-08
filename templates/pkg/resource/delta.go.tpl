{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	"bytes"

	"k8s.io/apimachinery/pkg/api/equality"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"
)

// Hack to avoid import errors during build...
var (
	_ = &bytes.Buffer{}
	_ = &acktags.Tags{}
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

// newResourceDeltaForPreDelete returns a new `ackcompare.Delta` that includes
// fields normally excluded by compare.is_ignored. Used for pre-delete sync.
func newResourceDeltaForPreDelete(
	a *resource,
	b *resource,
) *ackcompare.Delta {
	delta := ackcompare.NewDelta()
	if ((a == nil && b != nil) ||
			(a != nil && b == nil)) {
		delta.Add("", a, b)
		return delta
	}

{{ GoCodeCompareForPreDelete .CRD "delta" "a.ko" "b.ko" 1}}
	return delta
}
