{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
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
{{ GoCodeCompare .CRD "delta" "a.ko" "b.ko" 1}}
    return delta
}
