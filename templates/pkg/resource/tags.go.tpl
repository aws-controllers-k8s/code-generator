{{ template "boilerplate" }}

package {{ .CRD.Names.Snake }}

import(
    acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"

    svcapitypes "github.com/aws-controllers-k8s/{{ .ControllerName }}-controller/apis/{{ .APIVersion }}"
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
        if t.{{ $keyMemberName}} != nil {
            if t.{{ $valueMemberName }} == nil {
                result[*t.{{ $keyMemberName}}] = ""
            } else {
                result[*t.{{ $keyMemberName }}] = *t.{{ $valueMemberName }}
            }
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

// IgnoreAWSTags ignores tags that have keys that start with "aws:"
// is needed to ensure the controller does not attempt to remove
// tags set by AWS
// Eg. resources created with cloudformation have tags that cannot be
// removed by an ACK controller
func IgnoreAWSTags(tags acktags.Tags) {
	for k := range tags {
		if strings.HasPrefix(k, "aws:") ||
			strings.HasPrefix(k, "services.k8s.aws/") {
			delete(tags, k)
		}
	}
}

// SyncAWSTags ensures AWS-managed tags (prefixed with "aws:") from the latest resource state
// are preserved in the desired state. This prevents the controller from attempting to
// modify AWS-managed tags, which would result in an error.
//
// AWS-managed tags are automatically added by AWS services (e.g., CloudFormation, Service Catalog)
// and cannot be modified or deleted through normal tag operations. Common examples include:
// - aws:cloudformation:stack-name
// - aws:servicecatalog:productArn
//
// Parameters:
//   - a: The target Tags map to be updated (typically desired state)
//   - b: The source Tags map containing AWS-managed tags (typically latest state)
//
// Example:
//
//	latest := Tags{"aws:cloudformation:stack-name": "my-stack", "environment": "prod"}
//	desired := Tags{"environment": "dev"}
//	SyncAWSTags(desired, latest)
//	desired now contains {"aws:cloudformation:stack-name": "my-stack", "environment": "dev"}
func SyncAWSTags(a acktags.Tags, b acktags.Tags) {
	for k := range b {
		if strings.HasPrefix(k, "aws:") {
			a[k] = b[k]
		}
	}
}
{{ end }}