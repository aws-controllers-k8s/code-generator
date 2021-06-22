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

package ack

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	ttpl "text/template"

	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	ackutil "github.com/aws-controllers-k8s/code-generator/pkg/util"
)

/*
The following hook points are supported in the ACK controller resource manager
code paths:

* sdk_read_one_pre_build_request
* sdk_read_many_pre_build_request
* sdk_get_attributes_pre_build_request
* sdk_create_pre_build_request
* sdk_update_pre_build_request
* sdk_delete_pre_build_request
* sdk_read_one_post_build_request
* sdk_read_many_post_build_request
* sdk_get_attributes_post_build_request
* sdk_create_post_build_request
* sdk_update_post_build_request
* sdk_delete_post_build_request
* sdk_read_one_post_request
* sdk_read_many_post_request
* sdk_get_attributes_post_request
* sdk_create_post_request
* sdk_update_post_request
* sdk_delete_post_request
* sdk_read_one_pre_set_output
* sdk_read_many_pre_set_output
* sdk_get_attributes_pre_set_output
* sdk_create_pre_set_output
* sdk_update_pre_set_output
* sdk_read_one_post_set_output
* sdk_read_many_post_set_output
* sdk_get_attributes_post_set_output
* sdk_create_post_set_output
* sdk_update_post_set_output
* delta_pre_compare
* delta_post_compare

The "pre_build_request" hooks are called BEFORE the call to construct
the Input shape that is used in the API operation and therefore BEFORE
any call to validate that Input shape.

The "post_build_request" hooks are called AFTER the call to construct
the Input shape but BEFORE the API operation.

The "post_request" hooks are called IMMEDIATELY AFTER the API operation
aws-sdk-go client call.  These hooks will have access to a Go variable
named `resp` that refers to the aws-sdk-go client response and a Go
variable named `respErr` that refers to any error returned from the
aws-sdk-go client call.

The "pre_set_output" hooks are called BEFORE the code that processes the
Outputshape (the pkg/generate/code.SetOutput function). These hooks will
have access to a Go variable named `ko` that represents the concrete
Kubernetes CR object that will be returned from the main method
(sdkFind, sdkCreate, etc). This `ko` variable will have been defined
immediately before the "pre_set_output" hooks as a copy of the resource
that is supplied to the main method, like so:

```go
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()
```

The "post_set_output" hooks are called AFTER the the information from the API call
is merged with the copy of the original Kubernetes object. These hooks will
have access to the updated Kubernetes object `ko`, the response of the API call
(and the original Kubernetes CR object if its sdkUpdate)

The "delta_pre_compare" hooks are called BEFORE the generated code that
compares two resources.

The "delta_post_compare" hooks are called AFTER the generated code that
compares two resources.

*/

// ResourceHookCode returns a string with custom callback code for a resource
// and hook identifier
func ResourceHookCode(
	templateBasePaths []string,
	r *ackmodel.CRD,
	hookID string,
	vars interface{},
	funcMap ttpl.FuncMap,
) (string, error) {
	resourceName := r.Names.Original
	if resourceName == "" || hookID == "" {
		return "", nil
	}
	c := r.Config()
	if c == nil {
		return "", nil
	}
	rConfig, ok := c.Resources[resourceName]
	if !ok {
		return "", nil
	}
	hook, ok := rConfig.Hooks[hookID]
	if !ok {
		return "", nil
	}
	if hook.Code != nil {
		return *hook.Code, nil
	}
	if hook.TemplatePath == nil {
		err := fmt.Errorf(
			"resource %s hook config for %s is invalid. Need either code or template_path",
			resourceName, hookID,
		)
		return "", err
	}
	for _, basePath := range templateBasePaths {
		tplPath := filepath.Join(basePath, *hook.TemplatePath)
		if !ackutil.FileExists(tplPath) {
			continue
		}
		tplContents, err := ioutil.ReadFile(tplPath)
		if err != nil {
			err := fmt.Errorf(
				"resource %s hook config for %s is invalid: error reading %s: %s",
				resourceName, hookID, tplPath, err,
			)
			return "", err
		}
		t := ttpl.New(tplPath)
		t = t.Funcs(funcMap)
		if t, err = t.Parse(string(tplContents)); err != nil {
			err := fmt.Errorf(
				"resource %s hook config for %s is invalid: error parsing %s: %s",
				resourceName, hookID, tplPath, err,
			)
			return "", err
		}
		var b bytes.Buffer
		// TODO(jaypipes): Instead of nil for template vars here, maybe pass in
		// a struct of variables?
		if err := t.Execute(&b, vars); err != nil {
			err := fmt.Errorf(
				"resource %s hook config for %s is invalid: error executing %s: %s",
				resourceName, hookID, tplPath, err,
			)
			return "", err
		}
		return b.String(), nil
	}
	err := fmt.Errorf(
		"resource %s hook config for %s is invalid: template_path %s not found",
		resourceName, hookID, *hook.TemplatePath,
	)
	return "", err
}
