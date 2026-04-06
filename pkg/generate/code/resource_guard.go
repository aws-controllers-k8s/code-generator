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

package code

import (
	"fmt"
	"strings"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// defaultRequeueAfterSeconds is the default delay in seconds before a requeue
// when a resource is not in an allowed state for update or delete.
const defaultRequeueAfterSeconds = 30

// ResourceIsUpdateable returns Go code that checks whether a resource can be
// updated based on its current status. If the resource is NOT updateable, the
// generated code returns ackrequeue.NeededAfter.
//
// This follows the same pattern as ResourceIsSynced in synced.go.
//
// Sample output:
//
//	if latest.ko.Status.Status != nil {
//	    if !ackutil.InStrings(*latest.ko.Status.Status, []string{"ACTIVE", "AVAILABLE"}) {
//	        return nil, ackrequeue.NeededAfter(
//	            fmt.Errorf("resource is in %s state, cannot be updated",
//	                *latest.ko.Status.Status),
//	            time.Duration(30)*time.Second,
//	        )
//	    }
//	}
func ResourceIsUpdateable(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// resource variable name — "latest" for sdkUpdate
	resVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) (string, error) {
	return resourceIsGuarded(cfg, r, resVarName, indentLevel, "updateable", "updated")
}

// ResourceIsDeletable returns Go code that checks whether a resource can be
// deleted based on its current status. If the resource is NOT deletable, the
// generated code returns ackrequeue.NeededAfter.
//
// Sample output:
//
//	if r.ko.Status.Status != nil {
//	    if !ackutil.InStrings(*r.ko.Status.Status, []string{"ACTIVE", "AVAILABLE", "FAILED"}) {
//	        return nil, ackrequeue.NeededAfter(
//	            fmt.Errorf("resource is in %s state, cannot be deleted",
//	                *r.ko.Status.Status),
//	            time.Duration(30)*time.Second,
//	        )
//	    }
//	}
func ResourceIsDeletable(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// resource variable name — "r" for sdkDelete
	resVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) (string, error) {
	return resourceIsGuarded(cfg, r, resVarName, indentLevel, "deletable", "deleted")
}

// resourceIsGuarded is the shared implementation for ResourceIsUpdateable and
// ResourceIsDeletable. It reads the appropriate config block and generates
// guard code for each condition.
func resourceIsGuarded(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	resVarName string,
	indentLevel int,
	// configKey is "updateable" or "deletable"
	configKey string,
	// opVerb is "updated" or "deleted" — used in the error message
	opVerb string,
) (string, error) {
	out := ""
	resConfig := cfg.GetResourceConfig(r.Names.Original)
	if resConfig == nil {
		return out, nil
	}

	var conditions []ackgenconfig.StatusCondition
	var requeueSeconds int

	switch configKey {
	case "updateable":
		if resConfig.Updateable == nil || len(resConfig.Updateable.When) == 0 {
			return out, nil
		}
		conditions = resConfig.Updateable.When
		requeueSeconds = defaultRequeueAfterSeconds
		if resConfig.Updateable.RequeueAfterSeconds != nil &&
			*resConfig.Updateable.RequeueAfterSeconds > 0 {
			requeueSeconds = *resConfig.Updateable.RequeueAfterSeconds
		}
	case "deletable":
		if resConfig.Deletable == nil || len(resConfig.Deletable.When) == 0 {
			return out, nil
		}
		conditions = resConfig.Deletable.When
		requeueSeconds = defaultRequeueAfterSeconds
		if resConfig.Deletable.RequeueAfterSeconds != nil &&
			*resConfig.Deletable.RequeueAfterSeconds > 0 {
			requeueSeconds = *resConfig.Deletable.RequeueAfterSeconds
		}
	default:
		return "", fmt.Errorf("unknown config key %q", configKey)
	}

	for _, condCfg := range conditions {
		if condCfg.Path == nil || *condCfg.Path == "" {
			return "", fmt.Errorf(
				"resource %q: %s.when condition has empty path",
				r.Names.Original, configKey,
			)
		}
		if len(condCfg.In) == 0 {
			return "", fmt.Errorf(
				"resource %q, path %q: %s.when condition 'in' must not be empty",
				r.Names.Original, *condCfg.Path, configKey,
			)
		}

		_, err := getTopLevelField(r, *condCfg.Path)
		if err != nil {
			return "", fmt.Errorf(
				"resource %q: cannot find field for path %q: %w",
				r.Names.Original, *condCfg.Path, err,
			)
		}

		out += renderGuardBlock(
			resVarName, *condCfg.Path, condCfg.In,
			requeueSeconds, opVerb, indentLevel,
		)
	}

	return out, nil
}

// renderGuardBlock produces the Go source code for a single condition check.
// It generates a nil check on the field pointer, then an InStrings check
// against the allowed values, returning ackrequeue.NeededAfter if the value
// is not in the allowed set.
func renderGuardBlock(
	resVarName string,
	fieldPath string,
	allowedValues []string,
	requeueSeconds int,
	opVerb string,
	indentLevel int,
) string {
	indent := strings.Repeat("\t", indentLevel)
	fullPath := fmt.Sprintf("%s.ko.%s", resVarName, fieldPath)

	valuesSlice := fmt.Sprintf(`[]string{"%s"}`, strings.Join(allowedValues, `", "`))

	out := ""
	out += fmt.Sprintf("%sif %s != nil {\n", indent, fullPath)
	out += fmt.Sprintf("%s\tif !ackutil.InStrings(*%s, %s) {\n", indent, fullPath, valuesSlice)
	out += fmt.Sprintf("%s\t\treturn nil, ackrequeue.NeededAfter(\n", indent)
	out += fmt.Sprintf("%s\t\t\tfmt.Errorf(\"resource is in %%s state, cannot be %s\",\n", indent, opVerb)
	out += fmt.Sprintf("%s\t\t\t\t*%s),\n", indent, fullPath)
	out += fmt.Sprintf("%s\t\t\ttime.Duration(%d)*time.Second,\n", indent, requeueSeconds)
	out += fmt.Sprintf("%s\t\t)\n", indent)
	out += fmt.Sprintf("%s\t}\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	return out
}
