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

package model

import (
	"fmt"
	"sort"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
)

// FieldGroupOpType indicates whether a field group operation is for update or
// read.
type FieldGroupOpType int

const (
	FieldGroupOpTypeUpdate FieldGroupOpType = iota
	FieldGroupOpTypeRead
)

// FieldGroupOperation represents a resolved field-group operation. It ties a
// specific SDK API operation to the subset of CRD fields it manages, with
// fields partitioned into identifiers (locators) and payload (data).
type FieldGroupOperation struct {
	// OpType indicates whether this is an update or read field group.
	OpType FieldGroupOpType
	// OperationID is the SDK operation's exported name
	// (e.g., "PutImageScanningConfiguration").
	OperationID string
	// Operation is the resolved SDK operation pointer.
	Operation *awssdkmodel.Operation
	// IdentifierFields are CRD fields shared with ReadOne/Delete Input
	// shapes — these locate the resource and are always set in the API call.
	IdentifierFields []*Field
	// PayloadFields are the remaining CRD Spec fields from the operation's
	// Input shape — these are the data fields managed by this group.
	PayloadFields []*Field
	// Config holds the original generator.yaml configuration for this field
	// group operation (RequeueOnSuccess, Fields override, etc.).
	Config ackgenconfig.FieldGroupOperationConfig
}

// resolveFieldGroupOperations resolves all configured field-group operations
// (update and read) for a CRD, looking up SDK operations and partitioning
// Input shape members into identifier and payload fields.
//
// This must be called AFTER all CRD fields have been discovered (i.e., after
// processFields in GetCRDs), because it references crd.SpecFields.
func (r *CRD) resolveFieldGroupOperations() error {
	updateCfgs := r.cfg.GetUpdateFieldGroupOperations(r.Names.Original)
	for _, fgCfg := range updateCfgs {
		fg, err := r.resolveOneFieldGroupOperation(FieldGroupOpTypeUpdate, fgCfg)
		if err != nil {
			return fmt.Errorf(
				"resource %s, update_operations: %w",
				r.Names.Original, err,
			)
		}
		r.UpdateFieldGroups = append(r.UpdateFieldGroups, fg)
	}

	readCfgs := r.cfg.GetReadFieldGroupOperations(r.Names.Original)
	for _, fgCfg := range readCfgs {
		fg, err := r.resolveOneFieldGroupOperation(FieldGroupOpTypeRead, fgCfg)
		if err != nil {
			return fmt.Errorf(
				"resource %s, read_operations: %w",
				r.Names.Original, err,
			)
		}
		r.ReadFieldGroups = append(r.ReadFieldGroups, fg)
	}
	return nil
}

// resolveOneFieldGroupOperation resolves a single field-group operation config
// into a FieldGroupOperation with partitioned identifier and payload fields.
//
// For update operations, both identifier and payload fields come from the
// Input shape. For read operations, identifier fields come from the Input
// shape and payload fields come from the Output shape (the data we read back).
func (r *CRD) resolveOneFieldGroupOperation(
	opType FieldGroupOpType,
	fgCfg ackgenconfig.FieldGroupOperationConfig,
) (*FieldGroupOperation, error) {
	opID := fgCfg.OperationID

	op, found := r.sdkAPI.API.Operations[opID]
	if !found {
		return nil, fmt.Errorf("operation %q not found in SDK", opID)
	}

	inputShape := op.InputRef.Shape
	if inputShape == nil {
		return nil, fmt.Errorf("operation %q has nil Input shape", opID)
	}

	// Build the set of identifier member names from ReadOne and Delete Input
	// shapes. These are the fields that locate the resource.
	identifierMemberNames := r.buildIdentifierMemberSet()

	var identifierFields []*Field
	var payloadFields []*Field

	if len(fgCfg.Fields) > 0 {
		// Explicit Fields override: use those as payload
		identifierFields, payloadFields = r.resolveExplicitFields(
			opID, inputShape, identifierMemberNames, fgCfg.Fields,
		)
	} else if opType == FieldGroupOpTypeRead {
		// Read operations: identifiers from Input, payload from Output
		identifierFields = r.resolveIdentifierFieldsFromInput(
			opID, inputShape, identifierMemberNames,
		)
		outputShape := op.OutputRef.Shape
		if outputShape != nil {
			payloadFields = r.resolvePayloadFieldsFromShape(
				opID, outputShape, identifierMemberNames,
			)
		}
	} else {
		// Update operations: both from Input
		identifierFields, payloadFields = r.resolveAutoDetectedFields(
			opID, inputShape, identifierMemberNames,
		)
	}

	return &FieldGroupOperation{
		OpType:           opType,
		OperationID:      opID,
		Operation:        op,
		IdentifierFields: identifierFields,
		PayloadFields:    payloadFields,
		Config:           fgCfg,
	}, nil
}

// buildIdentifierMemberSet returns a set of SDK member names that appear in
// the Input shapes of the ReadOne and/or Delete operations. These represent
// the "identifier" fields used to locate the resource.
func (r *CRD) buildIdentifierMemberSet() map[string]bool {
	idMembers := map[string]bool{}

	addFromOp := func(op *awssdkmodel.Operation) {
		if op == nil || op.InputRef.Shape == nil {
			return
		}
		for _, memberName := range op.InputRef.Shape.MemberNames() {
			idMembers[memberName] = true
		}
	}

	addFromOp(r.Ops.ReadOne)
	addFromOp(r.Ops.Delete)

	// If neither ReadOne nor Delete exist (unusual), fall back to ReadMany
	if r.Ops.ReadOne == nil && r.Ops.Delete == nil {
		addFromOp(r.Ops.ReadMany)
	}

	return idMembers
}

// resolveAutoDetectedFields partitions the field-group operation's Input shape
// members into identifier and payload fields using auto-detection.
func (r *CRD) resolveAutoDetectedFields(
	opID string,
	inputShape *awssdkmodel.Shape,
	identifierMemberNames map[string]bool,
) ([]*Field, []*Field) {
	var identifierFields []*Field
	var payloadFields []*Field

	for _, memberName := range inputShape.MemberNames() {
		fieldName := r.cfg.GetResourceFieldName(
			r.Names.Original, opID, memberName,
		)
		field, ok := r.SpecFields[fieldName]
		if !ok {
			// Not a CRD Spec field (e.g., request metadata). Skip.
			continue
		}

		if identifierMemberNames[memberName] {
			identifierFields = append(identifierFields, field)
		} else {
			payloadFields = append(payloadFields, field)
		}
	}
	return identifierFields, payloadFields
}

// resolveIdentifierFieldsFromInput returns CRD fields from the Input shape
// that are identifiers (shared with ReadOne/Delete).
func (r *CRD) resolveIdentifierFieldsFromInput(
	opID string,
	inputShape *awssdkmodel.Shape,
	identifierMemberNames map[string]bool,
) []*Field {
	var identifierFields []*Field
	for _, memberName := range inputShape.MemberNames() {
		fieldName := r.cfg.GetResourceFieldName(
			r.Names.Original, opID, memberName,
		)
		field, ok := r.SpecFields[fieldName]
		if !ok {
			continue
		}
		if identifierMemberNames[memberName] {
			identifierFields = append(identifierFields, field)
		}
	}
	return identifierFields
}

// resolvePayloadFieldsFromShape returns CRD Spec fields from a shape's
// members that are NOT identifiers. Used for read operations where the
// payload comes from the Output shape.
func (r *CRD) resolvePayloadFieldsFromShape(
	opID string,
	shape *awssdkmodel.Shape,
	identifierMemberNames map[string]bool,
) []*Field {
	var payloadFields []*Field
	for _, memberName := range shape.MemberNames() {
		if identifierMemberNames[memberName] {
			continue
		}
		fieldName := r.cfg.GetResourceFieldName(
			r.Names.Original, opID, memberName,
		)
		// Check both Spec and Status fields for reads — the output may
		// map to either.
		if field, ok := r.SpecFields[fieldName]; ok {
			payloadFields = append(payloadFields, field)
		} else if field, ok := r.StatusFields[fieldName]; ok {
			payloadFields = append(payloadFields, field)
		}
	}
	return payloadFields
}

// resolveExplicitFields uses the explicitly configured field names to build
// the payload field list, and identifies remaining Input members that are
// identifiers.
func (r *CRD) resolveExplicitFields(
	opID string,
	inputShape *awssdkmodel.Shape,
	identifierMemberNames map[string]bool,
	explicitFields []string,
) ([]*Field, []*Field) {
	// Build a set of explicit payload field names for fast lookup
	explicitSet := make(map[string]bool, len(explicitFields))
	for _, fn := range explicitFields {
		explicitSet[fn] = true
	}

	var identifierFields []*Field
	var payloadFields []*Field

	// First, resolve identifier fields from the Input shape
	for _, memberName := range inputShape.MemberNames() {
		fieldName := r.cfg.GetResourceFieldName(
			r.Names.Original, opID, memberName,
		)
		field, ok := r.SpecFields[fieldName]
		if !ok {
			continue
		}
		if identifierMemberNames[memberName] && !explicitSet[field.Names.Camel] {
			identifierFields = append(identifierFields, field)
		}
	}

	// Then, resolve payload fields from explicit config (in order)
	for _, fn := range explicitFields {
		if field, ok := r.SpecFields[fn]; ok {
			payloadFields = append(payloadFields, field)
		}
	}

	return identifierFields, payloadFields
}

// HasFieldGroupUpdates returns true if this CRD has field-group update
// operations configured and resolved.
func (r *CRD) HasFieldGroupUpdates() bool {
	return len(r.UpdateFieldGroups) > 0
}

// HasFieldGroupReads returns true if this CRD has field-group read
// operations configured and resolved.
func (r *CRD) HasFieldGroupReads() bool {
	return len(r.ReadFieldGroups) > 0
}

// FieldGroupPayloadFieldNames returns a sorted, deduplicated list of all
// payload field names across all field-group operations of the given type.
func (r *CRD) FieldGroupPayloadFieldNames(opType FieldGroupOpType) []string {
	seen := map[string]bool{}
	var groups []*FieldGroupOperation
	if opType == FieldGroupOpTypeUpdate {
		groups = r.UpdateFieldGroups
	} else {
		groups = r.ReadFieldGroups
	}
	for _, fg := range groups {
		for _, f := range fg.PayloadFields {
			seen[f.Names.Camel] = true
		}
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
