{{- template "boilerplate" }}

{{ $hubImportAlias := .HubVersion }}
package {{ .APIVersion }}

import (
    "encoding/json"
    "fmt"

    ctrlrtconversion "sigs.k8s.io/controller-runtime/pkg/conversion"
{{- if not .IsHub }}
    ctrlrt "sigs.k8s.io/controller-runtime"
    ackrtwh "github.com/aws-controllers-k8s/runtime/pkg/webhook"

    {{ $hubImportAlias }} "github.com/aws-controllers-k8s/{{ .ServiceIDClean }}-controller/apis/{{ .HubVersion }}"
{{- end }}
)

var (
    _ = fmt.Printf
    _ = json.Marshal
)

{{- if .IsHub }}
// Assert hub interface implementation {{ .SourceCRD.Names.Camel }}
var _ ctrlrtconversion.Hub = &{{ .SourceCRD.Names.Camel }}{}

// Hub marks this type as conversion hub.
func (*{{ .SourceCRD.Kind }}) Hub() {}
{{ else }}

func init() {
    webhook := ackrtwh.New(
        "conversion",
        "{{ .SourceCRD.Names.Camel }}",
        "{{ .APIVersion }}",
		func(mgr ctrlrt.Manager) error {
			return ctrlrt.NewWebhookManagedBy(mgr).
				For(&{{ .SourceCRD.Names.Camel }}{}).
				Complete()
		},
    )
    if err := ackrtwh.RegisterWebhook(webhook); err != nil {
        msg := fmt.Sprintf("cannot register webhook: %v", err)
        panic(msg)
    }
}

// Assert convertible interface implementation {{ .SourceCRD.Names.Camel }}
var _ ctrlrtconversion.Convertible = &{{ .SourceCRD.Names.Camel }}{}

// ConvertTo converts this {{ .SourceCRD.Kind }} to the Hub version ({{ .HubVersion }}).
func (src *{{ .SourceCRD.Kind }}) ConvertTo(dstRaw ctrlrtconversion.Hub) error {
{{- GoCodeConvert .SourceCRD .DestCRD true $hubImportAlias "src" "dstRaw" 1}}
}

// ConvertFrom converts the Hub version ({{ .HubVersion }}) to this {{ .SourceCRD.Kind }}.
func (dst *{{ .SourceCRD.Kind }}) ConvertFrom(srcRaw ctrlrtconversion.Hub) error {
{{- GoCodeConvert .SourceCRD .DestCRD false $hubImportAlias "dst" "srcRaw" 1}}
}

{{- end }}