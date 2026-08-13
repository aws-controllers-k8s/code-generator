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

	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// CustomSyncUpdate returns Go code that invokes the hand-written sync function
// for each of the resource's `custom_sync` Spec fields, and then short-circuits
// out of sdkUpdate when those fields are the only ones that differ.
//
// Fields configured with `custom_sync` are managed by a separate AWS API rather
// than by the resource's Update operation, so when nothing else has changed
// there is no Update call left to make. The generated code returns the desired
// state with the observed status grafted on, which is what the reconciler
// expects back from Update.
//
// All `custom_sync` fields on the resource are collected into a single
// DifferentExcept call. That is the main reason this is generated rather than
// hand-written: adding a second out-of-band field to a hand-written hook
// requires widening the existing DifferentExcept, and forgetting to do so makes
// the resource silently short-circuit out of legitimate updates.
//
// Placement in sdk_update.go.tpl is load-bearing. This block must come AFTER
// the `updateable.when` guard, so that a resource in a state where mutations
// are not allowed is requeued before any out-of-band API call is made, and so
// that the short-circuit below cannot return success past the guard. It must
// come BEFORE the Update operation's `custom_implementation`, which returns
// from sdkUpdate directly and would otherwise skip the sync entirely.
//
// The empty string is returned when the resource has no `custom_sync` fields.
//
// Sample output:
//
//	updatedDesired := desired.DeepCopy()
//	updatedDesired.SetStatus(latest)
//	if delta.DifferentAt("Spec.Tags") {
//	    err = rm.syncTags(ctx, desired, latest)
//	    if err != nil {
//	        return nil, err
//	    }
//	}
//	if !delta.DifferentExcept("Spec.Tags") {
//	    return rm.concreteResource(updatedDesired), nil
//	}
func CustomSyncUpdate(
	r *model.CRD,
	// desired resource variable name — "desired" for sdkUpdate
	desiredVarName string,
	// latest resource variable name — "latest" for sdkUpdate
	latestVarName string,
	// delta variable name — "delta" for sdkUpdate
	deltaVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	fields := r.CustomSyncFields()
	if len(fields) == 0 {
		return ""
	}
	indent := strings.Repeat("\t", indentLevel)

	fieldPaths := customSyncFieldPaths(r, fields)

	out := "\n"
	// The reconciler expects Update to hand back the desired state carrying the
	// observed status. Build it before syncing so the short-circuit below has
	// something to return.
	out += fmt.Sprintf(
		"%supdatedDesired := %s.DeepCopy()\n", indent, desiredVarName,
	)
	out += fmt.Sprintf(
		"%supdatedDesired.SetStatus(%s)\n", indent, latestVarName,
	)
	for i, f := range fields {
		out += fmt.Sprintf(
			"%sif %s.DifferentAt(%q) {\n", indent, deltaVarName, fieldPaths[i],
		)
		out += fmt.Sprintf(
			"%s\terr = rm.%s(ctx, %s, %s)\n",
			indent, f.CustomSyncMethodName(), desiredVarName, latestVarName,
		)
		out += fmt.Sprintf("%s\tif err != nil {\n", indent)
		out += fmt.Sprintf("%s\t\treturn nil, err\n", indent)
		out += fmt.Sprintf("%s\t}\n", indent)
		out += fmt.Sprintf("%s}\n", indent)
	}
	// Every custom_sync field has now been reconciled through its own API. If
	// nothing outside that set differs, there is no Update call to make.
	quoted := make([]string, 0, len(fieldPaths))
	for _, fp := range fieldPaths {
		quoted = append(quoted, fmt.Sprintf("%q", fp))
	}
	out += fmt.Sprintf(
		"%sif !%s.DifferentExcept(%s) {\n",
		indent, deltaVarName, strings.Join(quoted, ", "),
	)
	out += fmt.Sprintf(
		"%s\treturn rm.concreteResource(updatedDesired), nil\n", indent,
	)
	out += fmt.Sprintf("%s}\n", indent)
	return out
}

// CustomSyncCreate returns Go code that marks the resource unsynced after a
// successful create when any of its `custom_sync` fields is set.
//
// A `custom_sync` field is only ever applied in the update path, so immediately
// after create the field is present in the resource's Spec but has not been
// pushed to AWS. Setting the Synced condition to false makes the runtime requeue
// after requeue.DefaultRequeueAfterDuration (30 seconds), which lands in
// sdkUpdate and runs the sync. Without this, the resource would report itself
// synced while the field was still unapplied, and the correction would wait for
// the full resync period.
//
// The Synced condition carries a message naming the fields still to be applied,
// so that a user running `kubectl describe` sees why the resource is not synced
// yet and that the controller intends to sync again on its own.
//
// The empty string is returned when the resource has no `custom_sync` fields.
//
// Sample output:
//
//	if ko.Spec.Tags != nil {
//	    msg := "Sync pending for Spec.Tags; resource will be requeued"
//	    ackcondition.SetSynced(&resource{ko}, corev1.ConditionFalse, &msg, nil)
//	}
func CustomSyncCreate(
	r *model.CRD,
	// the variable name of the resource's Kubernetes object — "ko" for sdkCreate
	koVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	fields := r.CustomSyncFields()
	if len(fields) == 0 {
		return ""
	}
	indent := strings.Repeat("\t", indentLevel)
	specPrefix := customSyncSpecPrefix(r)

	// Only mark unsynced when there is actually something to sync. A resource
	// created without any of these fields set has nothing for the follow-up
	// reconcile to do.
	conditions := make([]string, 0, len(fields))
	for _, f := range fields {
		conditions = append(conditions, fmt.Sprintf(
			"%s.%s.%s != nil", koVarName, specPrefix, f.Names.Camel,
		))
	}

	// Name the pending fields in the condition message so `kubectl describe`
	// explains both why the resource is not synced and that the controller will
	// sync again without the user doing anything. Phrased to read correctly for
	// any number of fields.
	msg := fmt.Sprintf(
		"Sync pending for %s; resource will be requeued",
		strings.Join(customSyncFieldPaths(r, fields), ", "),
	)

	out := "\n"
	out += fmt.Sprintf(
		"%sif %s {\n", indent, strings.Join(conditions, " || "),
	)
	out += fmt.Sprintf("%s\tmsg := %q\n", indent, msg)
	out += fmt.Sprintf(
		"%s\tackcondition.SetSynced(&resource{%s}, corev1.ConditionFalse, &msg, nil)\n",
		indent, koVarName,
	)
	out += fmt.Sprintf("%s}\n", indent)
	return out
}

// customSyncFieldPaths returns the delta field path for each supplied field,
// e.g. "Spec.Tags". These are the paths that the generated delta.go passes to
// delta.Add, so they must be constructed the same way — see the fieldPath
// calculation in CompareResource.
func customSyncFieldPaths(r *model.CRD, fields []*model.Field) []string {
	specPrefix := customSyncSpecPrefix(r)
	paths := make([]string, 0, len(fields))
	for _, f := range fields {
		paths = append(paths, fmt.Sprintf("%s.%s", specPrefix, f.Names.Camel))
	}
	return paths
}

// customSyncSpecPrefix returns the Spec prefix without its leading dot, i.e.
// "Spec" for the default prefix config.
func customSyncSpecPrefix(r *model.CRD) string {
	return strings.TrimPrefix(r.Config().PrefixConfig.SpecField, ".")
}
