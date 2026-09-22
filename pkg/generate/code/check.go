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

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/fieldpath"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// CheckExceptionMessage returns Go code that contains a condition to
// check if the message_prefix/message_suffix specified for a particular HTTP status code in
// generator config is a prefix for the exception message returned by AWS API.
// If message_prefix/message_suffix field was not specified for this HTTP code in generator
// config, we return an empty string
//
// Sample Output:
//
// && strings.HasPrefix(awsErr.Message(), "Could not find model")
// && strings.HasSuffix(awsErr.Message(), "does not exist.")
func CheckExceptionMessage(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	httpStatusCode int,
) string {
	rConfig := cfg.GetResourceConfig(r.Names.Original)
	if rConfig != nil && rConfig.Exceptions != nil {
		excConfig, ok := rConfig.Exceptions.Errors[httpStatusCode]
		if !ok {
			return ""
		}
		if excConfig.MessagePrefix != nil {
			return fmt.Sprintf("&& strings.HasPrefix(awsErr.ErrorMessage(), \"%s\") ",
				*excConfig.MessagePrefix)
		}
		if excConfig.MessageSuffix != nil {
			return fmt.Sprintf("&& strings.HasSuffix(awsErr.ErrorMessage(), \"%s\") ",
				*excConfig.MessageSuffix)
		}
	}
	return ""
}

// CheckRequiredFieldsMissingFromShape returns Go code that contains a
// condition checking that the required fields in the supplied Shape have a
// non-nil value in the corresponding CR's Spec or Status substruct.
//
// Sample Output:
//
// return r.ko.Spec.APIID == nil || r.ko.Status.RouteID == nil
func CheckRequiredFieldsMissingFromShape(
	r *model.CRD,
	opType model.OpType,
	koVarName string,
	indentLevel int,
) (string, error) {
	var op *awssdkmodel.Operation
	switch opType {
	case model.OpTypeGet:
		op = r.Ops.ReadOne
	case model.OpTypeList:
		op = r.Ops.ReadMany
		// Use the wrapper-aware input shape so a ReadMany op configured with
		// input_wrapper_field_path surfaces the unwrapped identifier members
		// to the required-fields check. No-op when no wrapper is configured.
		readManyShape, err := r.GetInputShape(op)
		if err != nil {
			return "", err
		}
		return checkRequiredFieldsMissingFromShapeReadMany(
			r, koVarName, indentLevel, op, readManyShape)
	case model.OpTypeGetAttributes:
		op = r.Ops.GetAttributes
	case model.OpTypeSetAttributes:
		op = r.Ops.SetAttributes
	default:
		return "", nil
	}

	shape := op.InputRef.Shape
	return checkRequiredFieldsMissingFromShape(
		r,
		koVarName,
		indentLevel,
		op,
		shape,
	)
}

// mutuallyExclusiveIdentifierNilConditions returns, for each of the resource's
// configured mutually-exclusive identifier fields, a "<path> == nil" condition
// string, along with the set of CR paths so callers can exclude them from
// per-field required checks. The resource is uniquely identified by exactly one
// of these fields, so it is considered incomplete only when all of them are
// nil. Returns nil slices when the resource has no mutually-exclusive
// identifiers.
func mutuallyExclusiveIdentifierNilConditions(
	r *model.CRD,
	koVarName string,
) ([]string, map[string]bool, error) {
	if !r.HasMutuallyExclusiveIdentifiers() {
		return nil, nil, nil
	}
	identifierFields, err := r.GetMutuallyExclusiveIdentifierFields()
	if err != nil {
		return nil, nil, err
	}
	cfg := r.Config()
	conditions := make([]string, 0, len(identifierFields))
	paths := make(map[string]bool, len(identifierFields))
	for _, identifierField := range identifierFields {
		memberPath, targetField := findFieldInCR(cfg, r, identifierField.Names.Original)
		if targetField == nil {
			return nil, nil, fmt.Errorf(
				"resource %q: mutually_exclusive_identifiers field %q is not in the CR's Spec or Status",
				r.Names.Original, identifierField.Names.Original,
			)
		}
		path := fmt.Sprintf("%s%s.%s", koVarName, memberPath, targetField.Path)
		conditions = append(conditions, fmt.Sprintf("%s == nil", path))
		paths[path] = true
	}
	return conditions, paths, nil
}

func checkRequiredFieldsMissingFromShape(
	r *model.CRD,
	koVarName string,
	indentLevel int,
	op *awssdkmodel.Operation,
	shape *awssdkmodel.Shape,
) (string, error) {
	indent := strings.Repeat("\t", indentLevel)

	// When the resource declares mutually-exclusive identifiers, the resource is
	// uniquely identified by exactly one of the declared fields. Build a single
	// grouped condition that is true only when none of them are set, and collect
	// their CR paths so they are not required individually below. This makes the
	// generated check treat the input as incomplete unless at least one
	// identifier is present, mirroring the ReadMany handling.
	exclusiveConditions, exclusivePaths, err := mutuallyExclusiveIdentifierNilConditions(r, koVarName)
	if err != nil {
		return "", err
	}
	exclusiveGroupCondition := ""
	if len(exclusiveConditions) > 0 {
		exclusiveGroupCondition = fmt.Sprintf("(%s)", strings.Join(exclusiveConditions, " && "))
	}

	if shape == nil || len(shape.Required) == 0 {
		if exclusiveGroupCondition != "" {
			return fmt.Sprintf("%sreturn %s\n", indent, exclusiveGroupCondition), nil
		}
		return fmt.Sprintf("%sreturn false", indent), nil
	}

	// Loop over the required member fields in the shape and identify whether
	// the field exists in either the Status or the Spec of the resource and
	// generate an if condition checking for all required fields having non-nil
	// corresponding resource Spec/Status values
	missing := []string{}
	for _, memberName := range shape.Required {
		if r.UnpacksAttributesMap() {
			// We set the Attributes field specially... depending on whether
			// the SetAttributes API call uses the batch or single attribute
			// flavor
			if r.SetAttributesSingleAttribute() {
				if memberName == "AttributeName" || memberName == "AttributeValue" {
					continue
				}
			} else {
				if memberName == "Attributes" {
					continue
				}
			}
		}
		if r.IsPrimaryARNField(memberName) {
			primaryARNCondition := fmt.Sprintf(
				"(%s.Status.ACKResourceMetadata == nil || %s.Status.ACKResourceMetadata.ARN == nil)",
				koVarName, koVarName,
			)
			missing = append(missing, primaryARNCondition)
			continue
		}

		resVarPath, err := r.GetSanitizedMemberPath(memberName, op, koVarName)
		if err != nil {
			return "", fmt.Errorf(
				"resource %q: required field %q in shape %q is not in the CR's Spec or Status structs",
				r.Names.Original, memberName, shape.ShapeName,
			)
		}
		// Mutually-exclusive identifiers are not required individually; they are
		// covered by the grouped condition appended below.
		if exclusivePaths[resVarPath] {
			continue
		}
		missing = append(missing, fmt.Sprintf("%s == nil", resVarPath))
	}
	if exclusiveGroupCondition != "" {
		missing = append(missing, exclusiveGroupCondition)
	}
	// Use '||' because if any of the required fields are missing the object
	// is not created yet
	missingCondition := strings.Join(missing, " || ")
	return fmt.Sprintf("%sreturn %s\n", indent, missingCondition), nil
}

// checkRequiredFieldsMissingFromShapeReadMany is a special-case handling
// of those APIs where there is no ReadOne operation and instead the only way to
// grab information for a single object is to call the ReadMany/List operation
// with one of more filtering fields-- specifically identifier(s). This method
// locates an identifier field in the shape that can be populated with an
// identifier value from the CR.
//
// As an example, DescribeVpcs EC2 API call doesn't have a ReadOne operation or
// required fields. However, the input shape has a VpcIds field which can be
// populated using a VpcId, a field in the VPC CR's Status. Therefore, require
// the VpcId field to be present to ensure the returned array from the API call
// consists only of the desired Vpc.
//
// Sample Output:
//
// return r.ko.Status.VPCID == nil
func checkRequiredFieldsMissingFromShapeReadMany(
	r *model.CRD,
	koVarName string,
	indentLevel int,
	op *awssdkmodel.Operation,
	shape *awssdkmodel.Shape,
) (string, error) {
	indent := strings.Repeat("\t", indentLevel)
	result := fmt.Sprintf("%sreturn false", indent)

	// When the resource declares mutually-exclusive identifiers, the ReadMany
	// input typically has no required members, so the default `return false`
	// would let sdkFind list every resource and match an arbitrary one. Instead,
	// treat the read input as incomplete (returning true so sdkFind bails out
	// with NotFound) unless at least one of the declared identifiers is set on
	// the resource.
	if r.HasMutuallyExclusiveIdentifiers() {
		exclusiveConditions, _, err := mutuallyExclusiveIdentifierNilConditions(r, koVarName)
		if err != nil {
			return "", err
		}
		// Parenthesize the grouped condition to mirror the ReadOne handling and
		// stay correct if a future term is ever joined here with `||`.
		return fmt.Sprintf("%sreturn (%s)\n", indent, strings.Join(exclusiveConditions, " && ")), nil
	}

	reqIdentifier, _ := FindPluralizedIdentifiersInShape(r, shape, op)
	resVarPath, err := r.GetSanitizedMemberPath(reqIdentifier, op, koVarName)
	if err != nil {
		return result, nil
	}

	// Start the incomplete-input condition with the primary pluralized
	// identifier. For a composite key the read input has additional REQUIRED
	// members beyond the primary identifier (e.g. a BatchGet* op unwrapped via
	// input_wrapper_field_path whose element structure requires both Name and
	// Type). AND each such required member that resolves to a CR field into the
	// condition, so sdkFind returns NotFound unless EVERY required identifier
	// component is present. Members are deduplicated against the primary
	// identifier and skipped when they have no corresponding Spec/Status field,
	// preserving the single-identifier output for existing ReadMany resources
	// whose input shape has no required members.
	conditions := []string{fmt.Sprintf("%s == nil", resVarPath)}
	seen := map[string]bool{resVarPath: true}
	for _, memberName := range shape.Required {
		fieldName := r.Config().GetResourceFieldName(
			r.Names.Original, op.ExportedName, memberName,
		)
		memberPath, findErr := r.GetSanitizedMemberPath(fieldName, op, koVarName)
		if findErr != nil {
			// Required member has no corresponding CR field; the payload/read
			// path handles it, so it is not part of the identity gate.
			continue
		}
		if seen[memberPath] {
			continue
		}
		seen[memberPath] = true
		conditions = append(conditions, fmt.Sprintf("%s == nil", memberPath))
	}

	return fmt.Sprintf("%sreturn %s\n", indent, strings.Join(conditions, " || ")), nil
}

// CheckNilFieldPath returns the condition statement for Nil check
// on a field path. This nil check on field path is useful to avoid
// nil pointer panics when accessing a field value.
//
// This function only outputs the logical condition and not the "if" block
// so that the output can be reused in many templates, where
// logic inside "if" block can be different.
//
// Example Output for fieldpath "JWTConfiguration.Issuer.SomeField" is
// "ko.Spec.JWTConfiguration == nil || ko.Spec.JWTConfiguration.Issuer == nil"
func CheckNilFieldPath(field *model.Field, sourceVarName string) string {
	out := ""
	fp := fieldpath.FromString(field.Path)
	// remove fieldName from fieldPath before adding nil checks
	fp.Pop()
	fieldNamePrefix := ""
	for fp.Size() > 0 {
		fieldNamePrefix = fmt.Sprintf("%s.%s", fieldNamePrefix, fp.PopFront())
		out += fmt.Sprintf(" || %s%s == nil", sourceVarName, fieldNamePrefix)
	}
	return strings.TrimPrefix(out, " || ")
}

// CheckNilReferencesPath returns the condition statement for Nil check
// on the path in ReferencesConfig. This nil check on the reference path is
// useful to avoid nil pointer panics when accessing the referenced value.
//
// This function only outputs the logical condition and not the "if" block
// so that the output can be reused in many templates, where
// logic inside "if" block can be different.
//
// Example Output for ReferencesConfig path "Status.ACKResourceMetadata.ARN",
// and sourceVarName "obj" is
// "obj.Status.ACKResourceMetadata == nil || obj.Status.ACKResourceMetadata.ARN == nil"
func CheckNilReferencesPath(field *model.Field, sourceVarName string) string {
	out := ""
	if field.HasReference() {
		refPath := fieldpath.FromString(field.FieldConfig.References.Path)
		// Remove the front from reference path because "Spec" or "Status" being
		// an struct cannot be added in nil check
		fieldNamePrefix := "." + refPath.PopFront()
		for refPath.Size() > 0 {
			fieldNamePrefix = fmt.Sprintf("%s.%s", fieldNamePrefix, refPath.PopFront())
			out += fmt.Sprintf(" || %s%s == nil", sourceVarName, fieldNamePrefix)
		}
	}
	return strings.TrimPrefix(out, " || ")
}
