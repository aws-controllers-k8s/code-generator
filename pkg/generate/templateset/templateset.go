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

package templateset

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	ttpl "text/template"

	"github.com/pkg/errors"

	ackutil "github.com/aws-controllers-k8s/code-generator/pkg/util"
)

var (
	// ErrTemplateNotFound is returned when a template path couldn't be located
	// in any base search path
	ErrTemplateNotFound    = errors.New("Template Not Found")
	errMsgTemplateNotFound = "template %s not found in any base search path"
)

func errTemplateNotFound(templatePath string) error {
	return errors.WithMessage(
		ErrTemplateNotFound,
		fmt.Sprintf(errMsgTemplateNotFound, templatePath),
	)
}

// templateWithVars contains a template and the variables injected during execution
type templateWithVars struct {
	t *ttpl.Template
	v interface{}
}

// TemplateSet contains a set of templates and copy files for a particular
// target
type TemplateSet struct {
	// baseSearchPaths is a set of directories from which templates will be
	// searched
	baseSearchPaths []string
	includePaths    []string
	copyPaths       []string
	templates       map[string]templateWithVars
	funcMap         ttpl.FuncMap
	executed        map[string]*bytes.Buffer
}

// New returns a pointer to a TemplateSet
func New(
	baseSearchPaths []string,
	includePaths []string,
	copyPaths []string,
	funcMap ttpl.FuncMap,
) *TemplateSet {
	return &TemplateSet{
		baseSearchPaths: baseSearchPaths,
		includePaths:    includePaths,
		copyPaths:       copyPaths,
		funcMap:         funcMap,
		templates:       map[string]templateWithVars{},
		executed:        map[string]*bytes.Buffer{},
	}
}

// Add constructs a named template from a path and variables
func (ts *TemplateSet) Add(
	outPath string,
	templatePath string,
	vars interface{},
) error {
	var foundPath string
	for _, basePath := range ts.baseSearchPaths {
		path := filepath.Join(basePath, templatePath)
		if ackutil.FileExists(path) {
			foundPath = path
			break
		}
	}

	if foundPath == "" {
		return errTemplateNotFound(templatePath)
	}

	tplContents, err := ioutil.ReadFile(foundPath)
	if err != nil {
		return err
	}
	t := ttpl.New(foundPath)
	t = t.Funcs(ts.funcMap)
	t, err = t.Parse(string(tplContents))
	if err != nil {
		return err
	}
	if err = ts.joinIncludes(t); err != nil {
		return err
	}
	ts.templates[outPath] = templateWithVars{t, vars}
	return nil
}

// joinIncludes adds all include templates to the supplied template
func (ts *TemplateSet) joinIncludes(t *ttpl.Template) error {
	var err error
	for _, basePath := range ts.baseSearchPaths {
		for _, includePath := range ts.includePaths {
			tplPath := filepath.Join(basePath, includePath)
			if !ackutil.FileExists(tplPath) {
				continue
			}
			if t, err = includeTemplate(t, tplPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// Execute runs all of the template and copy files in our TemplateSet and
// returns whether any error occurred executing any of the templates. Once
// Execute() is run, `TemplateSet.Executed()` can be used to iterate over a set
// of byte buffers containing the output of executed templates
func (ts *TemplateSet) Execute() error {
	for path, tv := range ts.templates {
		var b bytes.Buffer
		if err := tv.t.Execute(&b, tv.v); err != nil {
			return err
		}
		ts.executed[path] = &b
	}
	for _, basePath := range ts.baseSearchPaths {
		for _, path := range ts.copyPaths {
			copyPath := filepath.Join(basePath, path)
			if !ackutil.FileExists(copyPath) {
				continue
			}
			b, err := byteBufferFromFile(copyPath)
			if err != nil {
				return err
			}
			ts.executed[path] = b
		}
	}
	return nil
}

// Executed returns a map, keyed by the template or copy file path, of
// *bytes.Buffer objects containing executed template or copied file contents
func (ts *TemplateSet) Executed() map[string]*bytes.Buffer {
	return ts.executed
}

func byteBufferFromFile(path string) (*bytes.Buffer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	fsize := fi.Size()
	b := make([]byte, fsize)

	_, err = f.Read(b)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), nil
}

// includeTemplate includes a template into a supplied Template struct
func includeTemplate(t *ttpl.Template, tplPath string) (*ttpl.Template, error) {
	tplContents, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return nil, err
	}
	if t, err = t.Parse(string(tplContents)); err != nil {
		return nil, err
	}
	return t, nil
}
