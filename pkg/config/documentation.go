package config

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// DocumentationConfig represents the configuration of the documentation file,
// used to override or append documentation to any of the resource fields
type DocumentationConfig struct {
	Resources map[string]*ResourceDocsConfig `json:"resources,omitempty"`
}

// ResourceDocsConfig represents the configuration for the documentation
// overrides of a single resource
type ResourceDocsConfig struct {
	Fields map[string]*FieldDocsConfig `json:"fields,omitempty"`
}

// FieldDocsConfig represents the configuration for the documentation overrides
// of a single field
type FieldDocsConfig struct {
	// Append specifies a string that will be added to the end of the existing
	// GoDoc comment for the field
	Append *string `json:"append,omitempty"`
	// Prepend specifies a string that will be added before the existing
	// GoDoc comment for the field
	Prepend *string `json:"prepend,omitempty"`
	// Override will entirely replace the GoDoc comment for the field
	Override *string `json:"override,omitempty"`
}

// NewDocumentationConfig returns a new DocumentationConfig object given a
// supplied path to a config file
func NewDocumentationConfig(
	configPath string,
) (DocumentationConfig, error) {
	if configPath == "" {
		return DocumentationConfig{}, nil
	}
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return DocumentationConfig{}, err
	}
	gc := DocumentationConfig{}
	if err = yaml.Unmarshal(content, &gc); err != nil {
		return DocumentationConfig{}, err
	}
	return gc, nil
}
