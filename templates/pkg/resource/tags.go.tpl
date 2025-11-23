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
{{- if eq "list" $tagFieldShapeType }}
{{- $tagFieldGoType = (print "[]*svcapitypes." $tagField.GoTypeElem) }}
{{- end }}
// convertToOrderedACKTags converts the tags parameter into 'acktags.Tags' shape.
// This method helps in creating the hub(acktags.Tags) for merging
// default controller tags with existing resource tags. It also returns a slice
// of keys maintaining the original key Order when the tags are a list
func convertToOrderedACKTags(tags {{ $tagFieldGoType }}) (acktags.Tags, []string) {
    result := acktags.NewTags()
    keyOrder := []string{}
{{- if $hookCode := Hook .CRD "pre_convert_to_ack_tags" }}
{{ $hookCode }}
{{ end }}
{{ GoCodeConvertToACKTags .CRD "tags" "result" "keyOrder" 1 }}
{{- if $hookCode := Hook .CRD "post_convert_to_ack_tags" }}
{{ $hookCode }}
{{ end }}
    return result, keyOrder
}

// fromACKTags converts the tags parameter into {{ $tagFieldGoType }} shape.
// This method helps in setting the tags back inside AWSResource after merging
// default controller tags with existing resource tags. When a list, 
// it maintains the order from original 
func fromACKTags(tags acktags.Tags, keyOrder []string) {{ $tagFieldGoType }} {
    result := {{ $tagFieldGoType }}{}
{{- if $hookCode := Hook .CRD "pre_convert_from_ack_tags" }}
{{ $hookCode }}
{{ end }}
{{ GoCodeFromACKTags .CRD "tags" "keyOrder" "result" 1 }}
{{- if $hookCode := Hook .CRD "post_convert_from_ack_tags" }}
{{ $hookCode }}
{{ end }}
    return result
}
{{ end }}

// ignoreSystemTags ignores tags that have keys that start with "aws:"
// and systemTags defined on startup via the --resource-tags flag,
// to avoid patching them to the resourceSpec.
// Eg. resources created with cloudformation have tags that cannot be
// removed by an ACK controller
func ignoreSystemTags(tags acktags.Tags, systemTags []string) {
	for k := range tags {
		if strings.HasPrefix(k, "aws:") ||
			slices.Contains(systemTags, k) {
			delete(tags, k)
		}
	}
}

// syncAWSTags ensures AWS-managed tags (prefixed with "aws:") from the latest resource state
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
func syncAWSTags(a acktags.Tags, b acktags.Tags) {
	for k := range b {
		if strings.HasPrefix(k, "aws:") {
			a[k] = b[k]
		}
	}
}
{{ end }}