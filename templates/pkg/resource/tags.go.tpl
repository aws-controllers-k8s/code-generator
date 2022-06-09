{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import(
    acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"

    svcapitypes "github.com/aws-controllers-k8s/{{ .ServicePackageName }}-controller/apis/{{ .APIVersion }}"
)

var (
	_ = svcapitypes.{{ .CRD.Kind }}{}
	_ = acktags.NewTags()
)

{{- if $hookCode := Hook .CRD "convert_tags" }}
{{ $hookCode }}
{{ else -}}
{{- $tagField := .CRD.GetTagField }}
{{- if $tagField }}
{{- $tagFieldShapeType := $tagField.ShapeRef.Shape.Type }}
{{- $tagFieldGoType := $tagField.GoType }}
{{- $keyMemberName := .CRD.GetTagKeyMemberName }}
{{- $valueMemberName := .CRD.GetTagValueMemberName }}
{{- if eq "list" $tagFieldShapeType }}
{{- $tagFieldGoType = (print "[]*svcapitypes." $tagField.GoTypeElem) }}
{{- end }}
// ToACKTags converts the tags parameter into 'acktags.Tags' shape.
// This method helps in creating the hub(acktags.Tags) for merging
// default controller tags with existing resource tags.
func ToACKTags(tags {{ $tagFieldGoType }}) acktags.Tags {
    result := acktags.NewTags()
{{- if $hookCode := Hook .CRD "pre_convert_to_ack_tags" }}
{{ $hookCode }}
{{ end }}
    if tags == nil || len(tags) == 0 {
        return result
    }
{{ if eq "map" $tagFieldShapeType }}
    for k, v := range tags {
        if v == nil {
            result[k] = ""
        } else {
            result[k] = *v
        }
    }
{{ else if eq "list" $tagFieldShapeType }}
    for _, t := range tags {
        if t.{{ $valueMemberName }} == nil {
            result[*t.{{ $keyMemberName}}] = ""
        } else {
            result[*t.{{ $keyMemberName }}] = *t.{{ $valueMemberName }}
        }
    }
{{ end }}
{{- if $hookCode := Hook .CRD "post_convert_to_ack_tags" }}
{{ $hookCode }}
{{ end }}
    return result
}

// FromACKTags converts the tags parameter into {{ $tagFieldGoType }} shape.
// This method helps in setting the tags back inside AWSResource after merging
// default controller tags with existing resource tags.
func FromACKTags(tags acktags.Tags) {{ $tagFieldGoType }} {
    result := {{ $tagFieldGoType }}{}
{{- if $hookCode := Hook .CRD "pre_convert_from_ack_tags" }}
{{ $hookCode }}
{{ end }}
    for k, v := range tags {
{{- if eq "map" $tagFieldShapeType }}
        vCopy := v
        result[k] = &vCopy
{{- else if eq "list" $tagFieldShapeType }}
        kCopy := k
        vCopy := v
        tag := svcapitypes.{{ $tagField.GoTypeElem }}{ {{ $keyMemberName }}: &kCopy, {{ $valueMemberName }} : &vCopy}
        result = append(result, &tag)
{{- end }}
    }
{{- if $hookCode := Hook .CRD "post_convert_from_ack_tags" }}
{{ $hookCode }}
{{ end }}
    return result
}
{{ end }}
{{ end }}