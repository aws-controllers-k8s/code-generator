// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

var (
	// BuildDate of application at compile time (-X 'main.buildDate=$(BUILDDATE)').
	BuildDate string = "No Build Date Provided."
	// Version of application at compile time (-X 'main.version=$(VERSION)').
	Version string = "(Unknown Version)"
	// BuildHash is the GIT hash of application at compile time (-X 'main.buildHash=$(GITCOMMIT)').
	BuildHash string = "No Git-hash Provided."
	// GoVersion is the Go compiler version used to compile this binary
	GoVersion string
)

func init() {
	GoVersion = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	var revision string
	var dirty bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			dirty = s.Value == "true"
		}
	}

	if BuildHash == "No Git-hash Provided." && revision != "" {
		BuildHash = revision
	}

	// When installed via `go install ...@version`, ldflags aren't set.
	// Build a version string from the embedded module/VCS info.
	if Version == "(Unknown Version)" {
		Version = versionFromBuildInfo(info.Main.Version, revision, dirty)
	}
}

func versionFromBuildInfo(moduleVersion string, revision string, dirty bool) string {
	// Remote install (go install @v0.58.1): clean tag, no VCS info.
	if revision == "" {
		if moduleVersion != "" && moduleVersion != "(devel)" {
			return moduleVersion
		}
		return "v0.0.0-dev"
	}

	// Local build: construct git-describe-style from VCS info.
	short := revision
	if len(short) > 7 {
		short = short[:7]
	}

	// Extract the semver base from a pseudo-version.
	// e.g. "v0.58.2-0.20260515001229-e6c811834c1d+dirty" → "v0.58.2"
	base := moduleVersion
	if i := strings.Index(base, "+"); i >= 0 {
		base = base[:i]
	}
	if parts := strings.SplitN(base, "-", 2); len(parts) > 0 && strings.HasPrefix(parts[0], "v") {
		base = parts[0]
	}
	if base == "" || base == "(devel)" {
		base = "v0.0.0"
	}

	v := fmt.Sprintf("%s-g%s", base, short)
	if dirty {
		v += "+dirty"
	}
	return v
}
