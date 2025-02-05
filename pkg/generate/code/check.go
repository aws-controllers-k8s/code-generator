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
) string {
	var op *awssdkmodel.Operation
	switch opType {
	case model.OpTypeGet:
		op = r.Ops.ReadOne
	case model.OpTypeList:
		op = r.Ops.ReadMany
		return checkRequiredFieldsMissingFromShapeReadMany(
			r, koVarName, indentLevel, op, op.InputRef.Shape)
	case model.OpTypeGetAttributes:
		op = r.Ops.GetAttributes
	case model.OpTypeSetAttributes:
		op = r.Ops.SetAttributes
	default:
		return ""
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

func checkRequiredFieldsMissingFromShape(
	r *model.CRD,
	koVarName string,
	indentLevel int,
	op *awssdkmodel.Operation,
	shape *awssdkmodel.Shape,
) string {
	indent := strings.Repeat("\t", indentLevel)
	if shape == nil || len(shape.Required) == 0 {
		return fmt.Sprintf("%sreturn false", indent)
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
			// If it isn't in our spec/status fields, we have a problem!
			msg := fmt.Sprintf(
				"GENERATION FAILURE! there's a required field %s in "+
					"Shape %s that isn't in either the CR's Spec or "+
					"Status structs!",
				memberName, shape.ShapeName,
			)
			panic(msg)
		}
		missing = append(missing, fmt.Sprintf("%s == nil", resVarPath))
	}
	// Use '||' because if any of the required fields are missing the object
	// is not created yet
	missingCondition := strings.Join(missing, " || ")
	return fmt.Sprintf("%sreturn %s\n", indent, missingCondition)
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
) string {
	indent := strings.Repeat("\t", indentLevel)
	result := fmt.Sprintf("%sreturn false", indent)

	reqIdentifier, _ := FindPluralizedIdentifiersInShape(r, shape, op)

	resVarPath, err := r.GetSanitizedMemberPath(reqIdentifier, op, koVarName)
	if err != nil {
		return result
	}

	result = fmt.Sprintf("%s == nil", resVarPath)
	return fmt.Sprintf("%sreturn %s\n", indent, result)
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
