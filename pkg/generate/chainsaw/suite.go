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

package chainsaw

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"
)

// suiteFS holds the files every controller's generated tests share: the step
// templates they use, and the suite's Chainsaw configuration.
//
//go:embed suite/templates/*.yaml suite/chainsaw.yaml.tpl
var suiteFS embed.FS

const (
	// TemplatesDirName is the directory of the shared step templates, in the
	// tests directory.
	TemplatesDirName = "_templates"
	// ConfigFileName is the suite's Chainsaw configuration, in the tests
	// directory. test-infra passes it to `chainsaw test`.
	ConfigFileName = ".chainsaw.yaml"
)

// SuiteFiles returns the files a service's generated tests share, by path
// relative to the tests directory: the step templates in _templates/ and the
// suite's .chainsaw.yaml. They're the same for every service, apart from the
// service's name in the configuration. Each starts with GeneratedHeader.
func SuiteFiles(service string) (map[string][]byte, error) {
	files := map[string][]byte{}
	templates, err := fs.Glob(suiteFS, "suite/templates/*.yaml")
	if err != nil {
		return nil, err
	}
	for _, p := range templates {
		b, err := suiteFS.ReadFile(p)
		if err != nil {
			return nil, err
		}
		files[path.Join(TemplatesDirName, path.Base(p))] = append([]byte(GeneratedHeader), b...)
	}

	config, err := template.ParseFS(suiteFS, "suite/chainsaw.yaml.tpl")
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	b.WriteString(GeneratedHeader)
	if err := config.Execute(&b, struct{ Service string }{service}); err != nil {
		return nil, fmt.Errorf("rendering %s: %w", ConfigFileName, err)
	}
	files[ConfigFileName] = b.Bytes()
	return files, nil
}

// WriteGeneratedFile writes content to path, unless path holds a file that
// doesn't start with GeneratedHeader: a controller's hand-written replacement
// for a generated file is kept. It reports whether it wrote the file.
func WriteGeneratedFile(path string, content []byte) (bool, error) {
	existing, err := os.ReadFile(path)
	switch {
	case err == nil && !strings.HasPrefix(string(existing), GeneratedHeader):
		return false, nil
	case err != nil && !os.IsNotExist(err):
		return false, err
	}
	return true, os.WriteFile(path, content, 0o644)
}
