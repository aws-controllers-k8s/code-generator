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
}
