{{ template "boilerplate" }}

package version

// GitVersion, GitCommit, and BuildDate describe the service controller build
// and are injected at controller build time via -ldflags.
var (
	GitVersion string
	GitCommit  string
	BuildDate  string
)

// ACKGenerateVersion and ACKGenerateBuildDate identify the ack-generate build
// that generated this controller's code. They are baked in at code-generation
// time and are immutable for the life of the generated code.
const (
	ACKGenerateVersion   = "{{ .Version }}"
	ACKGenerateBuildDate = "{{ .BuildDate }}"
)
