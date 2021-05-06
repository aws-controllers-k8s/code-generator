{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import (
    ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"

    svcapitypes "github.com/aws-controllers-k8s/{{ .ServiceIDClean }}-controller/apis/v1alpha1"
)

var (
    _ = svcapitypes.GroupVersion
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
{{ GoCodeCompareHelpers .CRD 0 -}}