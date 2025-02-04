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
	"sort"
	"strings"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/pkg/names"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// SetSDK returns the Go code that sets an SDK input shape's member fields from
// a CRD's fields.
//
// Assume a CRD called Repository that looks like this pseudo-schema:
//
// .Status
//
//	.Authors ([]*string)
//	.ImageData
//	  .Location (*string)
//	  .Tag (*string)
//	.Name (*string)
//
// And assume an SDK Shape CreateRepositoryInput that looks like this
// pseudo-schema:
//
// .Repository
//
//	.Authors ([]*string)
//	.ImageData
//	  .Location (*string)
//	  .Tag (*string)
//	.Name
//
// This function is called from a template that generates the Go code that
// represents linkage between the Kubernetes objects (CRs) and the aws-sdk-go
// (SDK) objects. If we call this function with the following parameters:
//
//	opType:			OpTypeCreate
//	sourceVarName:	ko
//	targetVarName:	res
//	indentLevel:	1
//
// Then this function should output something like this:
//
//	  field1 := []*string{}
//	  for _, elem0 := range r.ko.Spec.Authors {
//	      elem0 := &string{*elem0}
//	      field0 = append(field0, elem0)
//	  }
//	  res.Authors = field1
//	  field1 := &svcsdk.ImageData{}
//	  field1.SetLocation(*r.ko.Spec.ImageData.Location)
//	  field1.SetTag(*r.ko.Spec.ImageData.Tag)
//	  res.ImageData = field1
//		 res.SetName(*r.ko.Spec.Name)
//
// Note that for scalar fields, we use the SetXXX methods that are on all
// aws-sdk-go SDK structs
func SetSDK(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The type of operation to look for the Input shape
	opType model.OpType,
	// String representing the name of the variable that we will grab the Input
	// shape from. This will likely be "r.ko" since in the templates that call
	// this method, the "source variable" is the CRD struct which is used to
	// populate the target variable, which is the Input shape
	sourceVarName string,
	// String representing the name of the variable that we will be **setting**
	// with values we get from the Output shape. This will likely be
	// "res" since that is the name of the "target variable" that the
	// templates that call this method use for the Input shape.
	targetVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	var op *awssdkmodel.Operation
	switch opType {
	case model.OpTypeCreate:
		op = r.Ops.Create
	case model.OpTypeGet:
		op = r.Ops.ReadOne
	case model.OpTypeList:
		op = r.Ops.ReadMany
		return setSDKReadMany(cfg, r, op,
			sourceVarName, targetVarName, indentLevel)
	case model.OpTypeUpdate:
		op = r.Ops.Update
	case model.OpTypeDelete:
		op = r.Ops.Delete
	default:
		return ""
	}
	if op == nil {
		return ""
	}
	inputShape := op.InputRef.Shape
	if inputShape == nil {
		return ""
	}

	out := "\n"
	indent := strings.Repeat("\t", indentLevel)

	// Some input shapes for APIs that use GetAttributes API calls don't have
	// an Attributes member (example: all the Delete shapes...)
	_, foundAttrs := inputShape.MemberRefs["Attributes"]
	if r.UnpacksAttributesMap() && foundAttrs {
		// For APIs that use a pattern of a parameter called "Attributes" that
		// is of type `map[string]*string` to represent real, schema'd fields,
		// we need to set the input shape's "Attributes" member field to the
		// re-constructed, packed set of fields.
		//
		// Therefore, we output here something like this (example from SNS
		// Topic's Attributes map):
		//
		// attrMap := map[string]*string{}
		// attrMap["DeliveryPolicy"] = r.ko.Spec.DeliveryPolicy
		// attrMap["DisplayName"} = r.ko.Spec.DisplayName
		// attrMap["KmsMasterKeyId"] = r.ko.Spec.KMSMasterKeyID
		// attrMap["Policy"] = r.ko.Spec.Policy
		// if len(attrMap) > 0 {
		//     res.SetAttributes(attrMap)
		// }
		fieldConfigs := cfg.GetFieldConfigs(r.Names.Original)
		out += fmt.Sprintf("%sattrMap := map[string]*string{}\n", indent)
		sortedAttrFieldNames := []string{}
		for fName, fConfig := range fieldConfigs {
			if fConfig.IsAttribute {
				sortedAttrFieldNames = append(sortedAttrFieldNames, fName)
			}
		}
		sort.Strings(sortedAttrFieldNames)
		for _, fieldName := range sortedAttrFieldNames {
			fieldConfig := fieldConfigs[fieldName]
			fieldNames := names.New(fieldName)
			if !fieldConfig.IsReadOnly {
				sourceAdaptedVarName := sourceVarName + cfg.PrefixConfig.SpecField + "." + fieldNames.Camel
				out += fmt.Sprintf(
					"%sif %s != nil {\n",
					indent, sourceAdaptedVarName,
				)
				out += fmt.Sprintf(
					"%s\tattrMap[\"%s\"] = %s\n",
					indent, fieldName, sourceAdaptedVarName,
				)
				out += fmt.Sprintf(
					"%s}\n", indent,
				)
			}
		}
		out += fmt.Sprintf(
			"%sif len(attrMap) > 0 {\n", indent,
		)
		out += fmt.Sprintf("\t%s%s.SetAttributes(attrMap)\n", indent, targetVarName)
		out += fmt.Sprintf(
			"%s}\n", indent,
		)
	}

	opConfig, override := cfg.GetOverrideValues(op.ExportedName)
	for memberIndex, memberName := range inputShape.MemberNames() {
		if r.UnpacksAttributesMap() && memberName == "Attributes" {
			continue
		}

		if override {
			value, ok := opConfig[memberName]
			memberShapeRef, _ := inputShape.MemberRefs[memberName]
			memberShape := memberShapeRef.Shape

			if ok {
				switch memberShape.Type {
				case "boolean", "integer":
				case "string":
					value = "\"" + value + "\""
				default:
					panic("Member type not handled")
				}

				out += fmt.Sprintf("%s%s.Set%s(%s)\n", indent, targetVarName, memberName, value)
				continue
			}
		}

		if r.IsPrimaryARNField(memberName) {
			// if ko.Status.ACKResourceMetadata != nil && ko.Status.ACKResourceMetadata.ARN != nil {
			//     res.SetTopicArn(string(*ko.Status.ACKResourceMetadata.ARN))
			// }
			out += fmt.Sprintf(
				"%sif %s.Status.ACKResourceMetadata != nil && %s.Status.ACKResourceMetadata.ARN != nil {\n",
				indent, sourceVarName, sourceVarName,
			)
			out += fmt.Sprintf(
				"%s\t%s.Set%s(string(*%s.Status.ACKResourceMetadata.ARN))\n",
				indent, targetVarName, memberName, sourceVarName,
			)
			out += fmt.Sprintf(
				"%s}\n", indent,
			)
			continue
		}

		// Determine whether the input shape's field is in the Spec or the
		// Status struct and set the source variable appropriately.
		var f *model.Field
		sourceAdaptedVarName := sourceVarName

		// Handles field renames, if applicable
		fieldName := cfg.GetResourceFieldName(
			r.Names.Original,
			op.ExportedName,
			memberName,
		)

		// Check if we have any configurations instructing the code
		// generator to set an SDK input field from this specific
		// field path.
		fallbackFieldName := r.GetMatchingInputShapeFieldName(opType, fieldName)
		if fallbackFieldName != "" {
			fieldName = fallbackFieldName
		}

		inSpec, inStatus := r.HasMember(fieldName, op.ExportedName)
		if inSpec {
			sourceAdaptedVarName += cfg.PrefixConfig.SpecField
			f = r.SpecFields[fieldName]
		} else if inStatus {
			sourceAdaptedVarName += cfg.PrefixConfig.StatusField
			f = r.StatusFields[fieldName]
		} else {
			// TODO(jaypipes): check generator config for exceptions?
			continue
		}

		sourceAdaptedVarName += "." + f.Names.Camel
		sourceFieldPath := f.Names.Camel

		memberShapeRef, _ := inputShape.MemberRefs[memberName]
		memberShape := memberShapeRef.Shape

		// we construct variables containing temporary storage for sub-elements
		// and sub-fields that are structs. Names of fields are "f" appended by
		// the 0-based index of the field within the set of the target struct's
		// set of fields. Nested structs simply append another "f" and the
		// field index to the variable name.
		//
		// This means you can tell what field a temporary fields variable
		// represents by the name.
		//
		// For example, the field variable name "f0f5f2", it contains the third
		// field of the sixth field of the first field of the input shape being
		// constructed.
		//
		// If we have two levels of nested struct fields, we will end
		// up with a targetVarName of "field0f0f0" and the generated code
		// might look something like this:
		//
		// res := &sdkapi.CreateBookInput{}
		// f0 := &sdkapi.BookData{}
		// if ko.Spec.Author != nil {
		//     f0f0 := &sdkapi.Author{}
		//     if ko.Spec.Author.Address != nil {
		//         f0f0f0 := &sdkapi.Address{}
		//         f0f0f0.SetStreet(*ko.Spec.Author.Address.Street)
		//         f0f0f0.SetCity(*ko.Spec.Author.Address.City)
		//         f0f0f0.SetState(*ko.Spec.Author.Address.State)
		//         f0f0.Address = f0f0f0
		//     }
		//     if ko.Spec.Author.Name != nil {
		//         f0f0.SetName(*r.ko.Author.Name)
		//         f0.Author = f0f0
		//     }
		//     res.Book = f0
		// }
		//
		// It's ugly but at least consistent and mostly readable...
		//
		// For populating list fields, we need an iterator and a temporary
		// element variable. We name these "{fieldName}iter" and
		// "{fieldName}elem" respectively. For nested levels, the names will be
		// progressively longer.
		//
		// For list fields, we want to end up with something like this:
		//
		// res := &sdkapi.CreateCustomAvailabilityZoneInput{}
		// if ko.Spec.VPNGroupsMemberships != nil {
		//     f0 := []*sdkapi.VpnGroupMembership{}
		//     for _, f0iter := ko.Spec.VPNGroupMemberships {
		//         f0elem := &sdkapi.VpnGroupMembership{}
		//         f0elem.SetVpnId(f0elem.VPNID)
		//         f0 := append(f0, f0elem)
		//     }
		//     res.VpnMemberships = f0
		// }

		omitUnchangedFieldsOnUpdate := op == r.Ops.Update && r.OmitUnchangedFieldsOnUpdate()
		if omitUnchangedFieldsOnUpdate && inSpec {
			fieldJSONPath := fmt.Sprintf("%s.%s", cfg.PrefixConfig.SpecField[1:], f.Names.Camel)
			out += fmt.Sprintf(
				"%sif delta.DifferentAt(%q) {\n", indent, fieldJSONPath,
			)

			// increase indentation level
			indentLevel++
			indent = "\t" + indent
		}

		out += fmt.Sprintf(
			"%sif %s != nil {\n", indent, sourceAdaptedVarName,
		)

		switch memberShape.Type {
		case "list", "structure", "map":
			{
				memberVarName := fmt.Sprintf("f%d", memberIndex)
				out += varEmptyConstructorSDKType(
					cfg, r,
					memberVarName,
					memberShape,
					indentLevel+1,
				)
				out += setSDKForContainer(
					cfg, r,
					memberName,
					memberVarName,
					sourceFieldPath,
					sourceAdaptedVarName,
					memberShapeRef,
					opType,
					indentLevel+1,
				)
				out += setSDKForScalar(
					cfg, r,
					memberName,
					targetVarName,
					inputShape.Type,
					sourceFieldPath,
					memberVarName,
					memberShapeRef,
					indentLevel+1,
				)
			}
		default:
			if r.IsSecretField(memberName) {
				out += setSDKForSecret(
					cfg, r,
					memberName,
					targetVarName,
					sourceAdaptedVarName,
					indentLevel,
				)
			} else {
				out += setSDKForScalar(
					cfg, r,
					memberName,
					targetVarName,
					inputShape.Type,
					sourceFieldPath,
					sourceAdaptedVarName,
					memberShapeRef,
					indentLevel+1,
				)
			}
		}
		out += fmt.Sprintf(
			"%s}\n", indent,
		)

		if omitUnchangedFieldsOnUpdate && inSpec {
			// decrease indentation level
			indentLevel--
			indent = indent[1:]

			out += fmt.Sprintf(
				"%s}\n", indent,
			)
		}
	}
	return out
}

// SetSDKGetAttributes returns the Go code that sets the Input shape for a
// resource's GetAttributes operation.
//
// As an example, for the GetTopicAttributes SNS API call, the returned code
// looks like this:
//
// res.SetTopicArn(string(*r.ko.Status.ACKResourceMetadata.ARN))
//
// For the SQS API's GetQueueAttributes call, the returned code looks like this:
//
// res.SetQueueUrl(*r.ko.Status.QueueURL)
//
// You will note the difference due to the special handling of the ARN fields.
func SetSDKGetAttributes(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// String representing the name of the variable that we will grab the
	// Input shape from. This will likely be "r.ko.Spec" since in the templates
	// that call this method, the "source variable" is the CRD struct's Spec
	// field which is used to populate the target variable, which is the Input
	// shape
	sourceVarName string,
	// String representing the name of the variable that we will be **setting**
	// with values we get from the Output shape. This will likely be
	// "res" since that is the name of the "target variable" that the
	// templates that call this method use for the Input shape.
	targetVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	op := r.Ops.GetAttributes
	if op == nil {
		return ""
	}
	inputShape := op.InputRef.Shape
	if inputShape == nil {
		return ""
	}
	if !r.UnpacksAttributesMap() {
		// This is a bug in the code generation if this occurs...
		msg := fmt.Sprintf(
			"called SetSDKGetAttributes for a resource '%s' that doesn't unpack attributes map",
			r.Names.Original,
		)
		panic(msg)
	}

	out := "\n"
	indent := strings.Repeat("\t", indentLevel)

	inputFieldOverrides := map[string][]string{}
	rConfig := cfg.GetResourceConfig(r.Names.Original)
	if rConfig == nil {
		// This is a bug in the code generation if this occurs...
		msg := fmt.Sprintf(
			"called SetSDKGetAttributes for a resource '%s' that doesn't have a ResourceConfig",
			r.Names.Original,
		)
		panic(msg)
	}
	attrCfg := rConfig.UnpackAttributesMapConfig
	if attrCfg != nil && attrCfg.GetAttributesInput != nil {
		for memberName, override := range attrCfg.GetAttributesInput.Overrides {
			inputFieldOverrides[memberName] = override.Values
		}
	}
	for _, memberName := range inputShape.MemberNames() {
		if r.IsPrimaryARNField(memberName) {
			// if ko.Status.ACKResourceMetadata != nil && ko.Status.ACKResourceMetadata.ARN != nil {
			//     res.SetTopicArn(string(*ko.Status.ACKResourceMetadata.ARN))
			// } else {
			//     res.SetTopicArn(rm.ARNFromName(*ko.Spec.Name))
			// }
			out += fmt.Sprintf(
				"%sif %s.Status.ACKResourceMetadata != nil && %s.Status.ACKResourceMetadata.ARN != nil {\n",
				indent, sourceVarName, sourceVarName,
			)
			out += fmt.Sprintf(
				"%s\t%s.Set%s(string(*%s.Status.ACKResourceMetadata.ARN))\n",
				indent, targetVarName, memberName, sourceVarName,
			)
			nameField := r.SpecIdentifierField()
			if nameField != nil {
				// There is no name or ID field for the resource, so don't try
				// to set an ARN from a name. Example: Subscription from SNS...
				out += fmt.Sprintf(
					"%s} else {\n", indent,
				)
				out += fmt.Sprintf(
					"%s\t%s.Set%s(rm.ARNFromName(*%s.Spec.%s))\n",
					indent, targetVarName, memberName, sourceVarName, *nameField,
				)
			}
			out += fmt.Sprintf(
				"%s}\n", indent,
			)
			continue
		}

		// Some APIs to retrieve the attributes for a resource require passing
		// specific fields and field values. For example, in order to get all
		// of an SQS Queue's attributes, the SQS GetQueueAttributes API call's
		// Input shape's AttributeNames member needs to be set to
		// []string{"All"}...
		//
		// Go code output in this section will look something like this:
		//
		// {
		//     tmpVals := []*string{}
		//     tmpVal0 := "All"
		//     tmpVals = append(tmpVals, &tmpVal0)
		//     res.SetAttributeNames(tmpVals)
		// }
		if overrideValues, ok := inputFieldOverrides[memberName]; ok {
			memberShapeRef := inputShape.MemberRefs[memberName]
			out += fmt.Sprintf("%s{\n", indent)
			// We need to output a set of temporary strings that we will take a
			// reference to when constructing the values of the []*string or
			// *string members.
			if memberShapeRef.Shape.Type == "list" {
				out += fmt.Sprintf("%s\ttmpVals := []*string{}\n", indent)
				for x, overrideValue := range overrideValues {
					out += fmt.Sprintf("%s\ttmpVal%d := \"%s\"\n", indent, x, overrideValue)
					out += fmt.Sprintf("%s\ttmpVals = append(tmpVals, &tmpVal%d)\n", indent, x)
				}
				out += fmt.Sprintf("%s\t%s.Set%s(tmpVals)\n", indent, targetVarName, memberName)
			} else {
				out += fmt.Sprintf("%s\ttmpVal := \"%s\"\n", indent, overrideValues[0])
				out += fmt.Sprintf("%s\t%s.Set%s(&tmpVal)\n", indent, targetVarName, memberName)
			}
			out += fmt.Sprintf("%s}\n", indent)
			continue
		}

		cleanMemberNames := names.New(memberName)
		cleanMemberName := cleanMemberNames.Camel

		sourceVarPath := sourceVarName
		field, found := r.SpecFields[memberName]
		if found {
			sourceVarPath = sourceVarName + cfg.PrefixConfig.SpecField + "." + cleanMemberName
		} else {
			field, found = r.StatusFields[memberName]
			if !found {
				// If it isn't in our spec/status fields, just ignore it
				continue
			}
			sourceVarPath = sourceVarPath + cfg.PrefixConfig.StatusField + "." + cleanMemberName
		}
		out += fmt.Sprintf(
			"%sif %s != nil {\n",
			indent, sourceVarPath,
		)
		out += setSDKForScalar(
			cfg, r,
			memberName,
			targetVarName,
			inputShape.Type,
			cleanMemberName,
			sourceVarPath,
			field.ShapeRef,
			indentLevel+1,
		)
		out += fmt.Sprintf(
			"%s}\n", indent,
		)
	}
	return out
}

// SetSDKSetAttributes returns the Go code that sets the Input shape for a
// resource's SetAttributes operation.
//
// Unfortunately, the AWS SetAttributes API operations (even within the *same*
// API) are inconsistent regarding whether the SetAttributes sets a batch of
// attributes or a single attribute. We need to construct the method
// differently depending on this behaviour. For example, the SNS
// SetTopicAttributes API call actually only allows the caller to set a single
// attribute, which needs to be specified in an AttributeName and
// AttributeValue field in the Input shape. On the other hand, the SNS
// SetPlatformApplicationAttributes API call's Input shape has an Attributes
// field which is a map[string]string containing all the attribute key/value
// pairs to replace. Your guess is as good as mine as to why these APIs are
// different.
//
// The returned code looks something like this:
//
// attrMap := map[string]*string{}
//
//	if r.ko.Spec.DeliveryPolicy != nil {
//	    attrMap["DeliveryPolicy"] = r.ko.Spec.DeliveryPolicy
//	}
//
//	if r.ko.Spec.DisplayName != nil {
//	    attrMap["DisplayName"} = r.ko.Spec.DisplayName
//	}
//
//	if r.ko.Spec.KMSMasterKeyID != nil {
//	    attrMap["KmsMasterKeyId"] = r.ko.Spec.KMSMasterKeyID
//	}
//
//	if r.ko.Spec.Policy != nil {
//	    attrMap["Policy"] = r.ko.Spec.Policy
//	}
//
// res.SetAttributes(attrMap)
func SetSDKSetAttributes(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// String representing the name of the variable that we will grab the Input
	// shape from. This will likely be "r.ko" since in the templates that call
	// this method, the "source variable" is the CRD struct which is used to
	// populate the target variable, which is the Input shape
	sourceVarName string,
	// String representing the name of the variable that we will be **setting**
	// with values we get from the Output shape. This will likely be
	// "res" since that is the name of the "target variable" that the
	// templates that call this method use for the Input shape.
	targetVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	op := r.Ops.SetAttributes
	if op == nil {
		return ""
	}
	inputShape := op.InputRef.Shape
	if inputShape == nil {
		return ""
	}
	if !r.UnpacksAttributesMap() {
		// This is a bug in the code generation if this occurs...
		msg := fmt.Sprintf(
			"called SetSDKSetAttributes for a resource '%s' that doesn't unpack attributes map",
			r.Names.Original,
		)
		panic(msg)
	}

	if r.SetAttributesSingleAttribute() {
		// TODO(jaypipes): For now, because these APIs require *multiple* calls
		// to the backend, one for each attribute being set, we'll go ahead and
		// rely on the CustomOperation functionality to write code for these...
		return ""
	}

	out := "\n"
	indent := strings.Repeat("\t", indentLevel)

	for _, memberName := range inputShape.MemberNames() {
		if r.IsPrimaryARNField(memberName) {
			// if ko.Status.ACKResourceMetadata != nil && ko.Status.ACKResourceMetadata.ARN != nil {
			//     res.SetTopicArn(string(*ko.Status.ACKResourceMetadata.ARN))
			// } else {
			//     res.SetTopicArn(rm.ARNFromName(*ko.Spec.Name))
			// }
			out += fmt.Sprintf(
				"%sif %s.Status.ACKResourceMetadata != nil && %s.Status.ACKResourceMetadata.ARN != nil {\n",
				indent, sourceVarName, sourceVarName,
			)
			out += fmt.Sprintf(
				"%s\t%s.Set%s(string(*%s.Status.ACKResourceMetadata.ARN))\n",
				indent, targetVarName, memberName, sourceVarName,
			)
			nameField := r.SpecIdentifierField()
			if nameField != nil {
				// There is no name or ID field for the resource, so don't try
				// to set an ARN from a name. Example: Subscription from SNS...
				out += fmt.Sprintf(
					"%s} else {\n", indent,
				)
				out += fmt.Sprintf(
					"%s\t%s.Set%s(rm.ARNFromName(*%s.Spec.%s))\n",
					indent, targetVarName, memberName, sourceVarName, *nameField,
				)
			}
			out += fmt.Sprintf(
				"%s}\n", indent,
			)
			continue
		}
		if memberName == "Attributes" {
			// For APIs that use a pattern of a parameter called "Attributes" that
			// is of type `map[string]*string` to represent real, schema'd fields,
			// we need to set the input shape's "Attributes" member field to the
			// re-constructed, packed set of fields.
			//
			// Therefore, we output here something like this (example from SNS
			// Topic's Attributes map):
			//
			// attrMap := map[string]*string{}
			// if r.ko.Spec.DeliveryPolicy != nil {
			//     attrMap["DeliveryPolicy"] = r.ko.Spec.DeliveryPolicy
			// }
			// if r.ko.Spec.DisplayName != nil {
			//     attrMap["DisplayName"} = r.ko.Spec.DisplayName
			// }
			// if r.ko.Spec.KMSMasterKeyID != nil {
			//     attrMap["KmsMasterKeyId"] = r.ko.Spec.KMSMasterKeyID
			// }
			// if r.ko.Spec.Policy != nil {
			//     attrMap["Policy"] = r.ko.Spec.Policy
			// }
			// res.SetAttributes(attrMap)
			fieldConfigs := cfg.GetFieldConfigs(r.Names.Original)
			out += fmt.Sprintf("%sattrMap := map[string]*string{}\n", indent)
			sortedAttrFieldNames := []string{}
			for fName, fConfig := range fieldConfigs {
				if fConfig.IsAttribute {
					sortedAttrFieldNames = append(sortedAttrFieldNames, fName)
				}
			}
			sort.Strings(sortedAttrFieldNames)
			for _, fieldName := range sortedAttrFieldNames {
				fieldConfig := fieldConfigs[fieldName]
				fieldNames := names.New(fieldName)
				if !fieldConfig.IsReadOnly {
					sourceAdaptedVarName := sourceVarName + cfg.PrefixConfig.SpecField + "." + fieldNames.Camel
					out += fmt.Sprintf(
						"%sif %s != nil {\n",
						indent, sourceAdaptedVarName,
					)
					out += fmt.Sprintf(
						"%s\tattrMap[\"%s\"] = %s\n",
						indent, fieldName, sourceAdaptedVarName,
					)
					out += fmt.Sprintf(
						"%s}\n", indent,
					)
				}
			}
			out += fmt.Sprintf(
				"%sif len(attrMap) > 0 {\n", indent,
			)
			out += fmt.Sprintf("\t%s%s.SetAttributes(attrMap)\n", indent, targetVarName)
			out += fmt.Sprintf(
				"%s}\n", indent,
			)
			continue
		}

		// Handle setting any other Input shape fields that are not the ARN
		// field or the Attributes unpacked map. The field value may come from
		// either the Spec or the Status fields.
		cleanMemberNames := names.New(memberName)
		cleanMemberName := cleanMemberNames.Camel

		sourceVarPath := sourceVarName
		field, found := r.SpecFields[memberName]
		if found {
			sourceVarPath = sourceVarName + cfg.PrefixConfig.SpecField + "." + cleanMemberName
		} else {
			field, found = r.StatusFields[memberName]
			if !found {
				// If it isn't in our spec/status fields, just ignore it
				continue
			}
			sourceVarPath = sourceVarPath + cfg.PrefixConfig.StatusField + "." + cleanMemberName
		}
		out += fmt.Sprintf(
			"%sif %s != nil {\n",
			indent, sourceVarPath,
		)
		out += setSDKForScalar(
			cfg, r,
			memberName,
			targetVarName,
			inputShape.Type,
			cleanMemberName,
			sourceVarPath,
			field.ShapeRef,
			indentLevel+1,
		)
		out += fmt.Sprintf(
			"%s}\n", indent,
		)
	}
	return out
}

// setSDKReadMany is a special-case handling of those APIs where there is no
// ReadOne operation and instead the only way to grab information for a single
// object is to call the ReadMany/List operation with one of more filtering
// fields-- specifically identifier(s). This method populates this identifier
// field with the identifier shared between the shape and the CR. Note, in the
// case of multiple matching identifiers, the identifier field containing 'Id'
// will be the only field populated.
//
// As an example, DescribeVpcs EC2 API call doesn't have a ReadOne operation or
// required fields. However, the input shape VpcIds field can be populated using
// a VpcId, a field in the VPC CR's Status. Therefore, populate VpcIds field
// with the *single* VpcId value to ensure the returned array from the API call
// consists only of the desired Vpc.
//
// Sample Output:
//
//	if r.ko.Status.VPCID != nil {
//		f4 := []*string{}
//		f4 = append(f4, r.ko.Status.VPCID)
//		res.SetVpcIds(f4)
//	}
func setSDKReadMany(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	op *awssdkmodel.Operation,
	sourceVarName string,
	targetVarName string,
	indentLevel int,
) string {
	inputShape := op.InputRef.Shape
	if inputShape == nil {
		return ""
	}

	out := "\n"
	indent := strings.Repeat("\t", indentLevel)

	resVarPath := ""
	opConfig, override := cfg.GetOverrideValues(op.ExportedName)
	var err error
	for memberIndex, memberName := range inputShape.MemberNames() {
		if override {
			value, ok := opConfig[memberName]
			memberShapeRef, _ := inputShape.MemberRefs[memberName]
			memberShape := memberShapeRef.Shape
			if ok {
				switch memberShape.Type {
				case "boolean", "integer":
				case "string":
					value = "\"" + value + "\""
				default:
					panic(fmt.Sprintf("Unsupported shape type %s in "+
						"generate.code.setSDKReadMany", memberShape.Type))
				}

				out += fmt.Sprintf("%s%s.Set%s(%s)\n", indent, targetVarName, memberName, value)
				continue
			}
		}

		// Handles field renames, if applicable
		fieldName := cfg.GetResourceFieldName(
			r.Names.Original,
			op.ExportedName,
			memberName,
		)
		resVarPath, err = r.GetSanitizedMemberPath(memberName, op, sourceVarName)
		if err != nil {
			// memberName could be a plural identifier field, so check for
			// corresponding singular model identifier
			crIdentifier, shapeIdentifier := FindPluralizedIdentifiersInShape(r,
				inputShape, op)
			if strings.EqualFold(fieldName, shapeIdentifier) {
				resVarPath, err = r.GetSanitizedMemberPath(crIdentifier, op, sourceVarName)
				if err != nil {
					panic(fmt.Sprintf(
						"Unable to locate identifier field %s in "+
							"%s Spec/Status in generate.code.setSDKReadMany", crIdentifier, r.Kind))
				}
			} else {
				// TODO(jaypipes): check generator config for exceptions?
				continue
			}
		}

		memberShapeRef, _ := inputShape.MemberRefs[memberName]
		memberShape := memberShapeRef.Shape
		out += fmt.Sprintf(
			"%sif %s != nil {\n", indent, resVarPath,
		)

		switch memberShape.Type {
		case "list":
			// Expecting slice of identifiers
			memberVarName := fmt.Sprintf("f%d", memberIndex)
			// f0 := []*string{}
			out += varEmptyConstructorSDKType(
				cfg, r,
				memberVarName,
				memberShape,
				indentLevel+1,
			)

			//  f0 = append(f0, sourceVarName)
			out += fmt.Sprintf("%s\t%s = append(%s, %s)\n", indent,
				memberVarName, memberVarName, resVarPath)

			// res.SetIds(f0)
			out += setSDKForScalar(
				cfg, r,
				memberName,
				targetVarName,
				inputShape.Type,
				sourceVarName,
				memberVarName,
				memberShapeRef,
				indentLevel+1,
			)
		default:
			// For ReadMany that have a singular identifier field.
			// ex: DescribeReplicationGroups
			out += setSDKForScalar(
				cfg, r,
				memberName,
				targetVarName,
				inputShape.Type,
				sourceVarName,
				resVarPath,
				memberShapeRef,
				indentLevel+1,
			)
		}
		out += fmt.Sprintf(
			"%s}\n", indent,
		)
	}

	return out
}

// setSDKForContainer returns a string of Go code that sets the value of a
// target variable to that of a source variable. When the source variable type
// is a map, struct or slice type, then this function is called recursively on
// the elements or members of the source variable.
func setSDKForContainer(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The name of the SDK Input shape member we're outputting for
	targetFieldName string,
	// The variable name that we want to set a value to
	targetVarName string,
	// The path to the field that we access our source value from
	sourceFieldPath string,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	// ShapeRef of the target struct field
	targetShapeRef *awssdkmodel.ShapeRef,
	op model.OpType,
	indentLevel int,
) string {
	switch targetShapeRef.Shape.Type {
	case "structure":
		return SetSDKForStruct(
			cfg, r,
			targetFieldName,
			targetVarName,
			targetShapeRef,
			sourceFieldPath,
			sourceVarName,
			op,
			indentLevel,
		)
	case "list":
		return setSDKForSlice(
			cfg, r,
			targetFieldName,
			targetVarName,
			targetShapeRef,
			sourceFieldPath,
			sourceVarName,
			op,
			indentLevel,
		)
	case "map":
		return setSDKForMap(
			cfg, r,
			targetFieldName,
			targetVarName,
			targetShapeRef,
			sourceFieldPath,
			sourceVarName,
			op,
			indentLevel,
		)
	default:
		if r.IsSecretField(sourceFieldPath) {
			indent := strings.Repeat("\t", indentLevel)
			// if ko.Spec.MasterUserPassword != nil {
			out := fmt.Sprintf(
				"%sif %s != nil {\n",
				indent, sourceVarName,
			)
			out += setSDKForSecret(
				cfg, r,
				"",
				targetVarName,
				sourceVarName,
				indentLevel,
			)
			// }
			out += fmt.Sprintf("%s}\n", indent)
			return out
		}

		return setSDKForScalar(
			cfg, r,
			targetFieldName,
			targetVarName,
			targetShapeRef.Shape.Type,
			sourceFieldPath,
			sourceVarName,
			targetShapeRef,
			indentLevel,
		)
	}
}

// setSDKForSecret returns a string of Go code that sets a target variable to
// the value of a Secret when the type of the source variable is a
// SecretKeyReference.
//
// The Go code output from this function looks like this:
//
//     tmpSecret, err := rm.rr.SecretValueFromReference(ctx, ko.Spec.MasterUserPassword)
//     if err != nil {
//         return nil, ackrequeue.Needed(err)
//     }
//     if tmpSecret != "" {
//         res.SetMasterUserPassword(tmpSecret)
//     }
//
//     or:
//
//     tmpSecret, err := rm.rr.SecretValueFromReference(ctx, f3iter)
//     if err != nil {
//         return nil, ackrequeue.Needed(err)
//     }
//     if tmpSecret != "" {
//         f3elem = tmpSecret
//     }
//
// The second case is used when the SecretKeyReference field
// is a slice of `[]*string` in the original AWS API Input shape.

func setSDKForSecret(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The name of the SDK Shape field we're setting
	targetFieldName string,
	// The variable name that we want to set a value on
	targetVarName string,
	// The CR field that we access our source value from
	sourceVarName string,
	indentLevel int,
) string {

	out := ""
	indent := strings.Repeat("\t", indentLevel)
	secVar := "tmpSecret"

	//     tmpSecret, err := rm.rr.SecretValueFromReference(ctx, ko.Spec.MasterUserPassword)
	out += fmt.Sprintf(
		"%s\t%s, err := rm.rr.SecretValueFromReference(ctx, %s)\n",
		indent, secVar, sourceVarName,
	)
	//     if err != nil {
	//         return nil, ackrequeue.Needed(err)
	//     }
	out += fmt.Sprintf("%s\tif err != nil {\n", indent)
	out += fmt.Sprintf("%s\t\treturn nil, ackrequeue.Needed(err)\n", indent)
	out += fmt.Sprintf("%s\t}\n", indent)
	//     if tmpSecret != "" {
	//         res.SetMasterUserPassword(tmpSecret)
	//     }
	out += fmt.Sprintf("%s\tif tmpSecret != \"\" {\n", indent)
	if targetFieldName == "" {
		out += fmt.Sprintf(
			"%s\t\t%s = %s\n",
			indent, targetVarName, secVar,
		)
	} else {
		out += fmt.Sprintf(
			"%s\t\t%s.Set%s(%s)\n",
			indent, targetVarName, targetFieldName, secVar,
		)
	}
	out += fmt.Sprintf("%s\t}\n", indent)
	return out
}

// SetSDKForStruct returns a string of Go code that sets a target variable
// value to a source variable when the type of the source variable is a struct.
func SetSDKForStruct(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The name of the CR field we're outputting for
	targetFieldName string,
	// The variable name that we want to set a value to
	targetVarName string,
	// Shape Ref of the target struct field
	targetShapeRef *awssdkmodel.ShapeRef,
	// The path to the field that we access our source value from
	sourceFieldPath string,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	op model.OpType,
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	targetShape := targetShapeRef.Shape

	for memberIndex, memberName := range targetShape.MemberNames() {
		memberShapeRef := targetShape.MemberRefs[memberName]
		memberShape := memberShapeRef.Shape
		cleanMemberNames := names.New(memberName)
		cleanMemberName := cleanMemberNames.Camel
		sourceAdaptedVarName := sourceVarName + "." + cleanMemberName
		memberFieldPath := sourceFieldPath + "." + cleanMemberName

		// todo: To make `ignore` functionality work for all fields that has `ignore` set to `true`,
		// we need to add the below logic inside `SetSDK` function.

		// To check if the field member has `ignore` set to `true`.
		// This condition currently applies only for members of a field whose shape is `structure`
		var setCfg *ackgenconfig.SetFieldConfig
		f, ok := r.Fields[sourceFieldPath]
		if ok {
			mf, ok := f.MemberFields[memberName]
			if ok {
				setCfg = mf.GetSetterConfig(op)
				if setCfg != nil && setCfg.IgnoreSDKSetter() {
					continue
				}
			}
		}

		fallBackName := r.GetMatchingInputShapeFieldName(op, targetFieldName)
		if fallBackName == memberName {
			// TODO: implement @AmineHilaly
		}

		out += fmt.Sprintf(
			"%sif %s != nil {\n", indent, sourceAdaptedVarName,
		)
		switch memberShape.Type {
		case "list", "structure", "map":
			{
				memberVarName := fmt.Sprintf(
					"%sf%d",
					targetVarName, memberIndex,
				)
				out += varEmptyConstructorSDKType(
					cfg, r,
					memberVarName,
					memberShape,
					indentLevel+1,
				)
				out += setSDKForContainer(
					cfg, r,
					memberName,
					memberVarName,
					memberFieldPath,
					sourceAdaptedVarName,
					memberShapeRef,
					op,
					indentLevel+1,
				)
				out += setSDKForScalar(
					cfg, r,
					memberName,
					targetVarName,
					targetShape.Type,
					memberFieldPath,
					memberVarName,
					memberShapeRef,
					indentLevel+1,
				)
			}
		default:
			if r.IsSecretField(memberFieldPath) {
				out += setSDKForSecret(
					cfg, r,
					memberName,
					targetVarName,
					sourceAdaptedVarName,
					indentLevel,
				)
			} else {
				out += setSDKForScalar(
					cfg, r,
					memberName,
					targetVarName,
					targetShape.Type,
					memberFieldPath,
					sourceAdaptedVarName,
					memberShapeRef,
					indentLevel+1,
				)
			}
		}
		out += fmt.Sprintf(
			"%s}\n", indent,
		)
	}
	return out
}

// setSDKForSlice returns a string of Go code that sets a target variable value
// to a source variable when the type of the source variable is a slice.
func setSDKForSlice(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The name of the CR field we're outputting for
	targetFieldName string,
	// The variable name that we want to set a value to
	targetVarName string,
	// Shape Ref of the target struct field
	targetShapeRef *awssdkmodel.ShapeRef,
	// The path to the field that we access our source value from
	sourceFieldPath string,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	op model.OpType,
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	targetShape := targetShapeRef.Shape

	iterVarName := fmt.Sprintf("%siter", targetVarName)
	elemVarName := fmt.Sprintf("%selem", targetVarName)
	// for _, f0iter := range r.ko.Spec.Tags {
	out += fmt.Sprintf("%sfor _, %s := range %s {\n", indent, iterVarName, sourceVarName)
	//		f0elem := string{}
	out += varEmptyConstructorSDKType(
		cfg, r,
		elemVarName,
		targetShape.MemberRef.Shape,
		indentLevel+1,
	)
	//  f0elem = *f0iter
	//
	// or
	//
	//  f0elem.SetMyField(*f0iter)
	containerFieldName := ""
	if targetShape.MemberRef.Shape.Type == "structure" {
		containerFieldName = targetFieldName
	}
	out += setSDKForContainer(
		cfg, r,
		containerFieldName,
		elemVarName,
		sourceFieldPath,
		iterVarName,
		&targetShape.MemberRef,
		op,
		indentLevel+1,
	)
	addressOfVar := ""
	switch targetShape.MemberRef.Shape.Type {
	case "structure", "list", "map":
		break
	default:
		addressOfVar = "&"
	}
	//  f0 = append(f0, elem0)
	out += fmt.Sprintf("%s\t%s = append(%s, %s%s)\n", indent, targetVarName, targetVarName, addressOfVar, elemVarName)
	out += fmt.Sprintf("%s}\n", indent)
	return out
}

// setSDKForMap returns a string of Go code that sets a target variable value
// to a source variable when the type of the source variable is a map.
func setSDKForMap(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The name of the CR field we're outputting for
	targetFieldName string,
	// The variable name that we want to set a value to
	targetVarName string,
	// Shape Ref of the target struct field
	targetShapeRef *awssdkmodel.ShapeRef,
	// The path to the field that we access our source value from
	sourceFieldPath string,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	op model.OpType,
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	targetShape := targetShapeRef.Shape

	valIterVarName := fmt.Sprintf("%svaliter", targetVarName)
	keyVarName := fmt.Sprintf("%skey", targetVarName)
	valVarName := fmt.Sprintf("%sval", targetVarName)
	// for f0key, f0valiter := range r.ko.Spec.Tags {
	out += fmt.Sprintf("%sfor %s, %s := range %s {\n", indent, keyVarName, valIterVarName, sourceVarName)
	//		f0elem := string{}
	out += varEmptyConstructorSDKType(
		cfg, r,
		valVarName,
		targetShape.ValueRef.Shape,
		indentLevel+1,
	)
	//  f0val = *f0valiter
	//
	// or
	//
	//  f0val.SetMyField(*f0valiter)
	containerFieldName := ""
	if targetShape.ValueRef.Shape.Type == "structure" {
		containerFieldName = targetFieldName
	}
	out += setSDKForContainer(
		cfg, r,
		containerFieldName,
		valVarName,
		sourceFieldPath,
		valIterVarName,
		&targetShape.ValueRef,
		op,
		indentLevel+1,
	)
	addressOfVar := ""
	switch targetShape.ValueRef.Shape.Type {
	case "structure", "list", "map":
		break
	default:
		addressOfVar = "&"
	}
	// f0[f0key] = f0val
	out += fmt.Sprintf("%s\t%s[%s] = %s%s\n", indent, targetVarName, keyVarName, addressOfVar, valVarName)
	out += fmt.Sprintf("%s}\n", indent)
	return out
}

func varEmptyConstructorSDKType(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	varName string,
	// The shape we want to construct a new thing for
	shape *awssdkmodel.Shape,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	goType := shape.GoTypeWithPkgName()
	keepPointer := (shape.Type == "list" || shape.Type == "map")
	goType = model.ReplacePkgName(goType, r.SDKAPIPackageName(), "svcsdk", keepPointer)
	switch shape.Type {
	case "structure":
		// f0 := &svcsdk.BookData{}
		out += fmt.Sprintf("%s%s := &%s{}\n", indent, varName, goType)
	case "list", "map":
		// f0 := []*string{}
		out += fmt.Sprintf("%s%s := %s{}\n", indent, varName, goType)
	default:
		// var f0 string
		out += fmt.Sprintf("%svar %s %s\n", indent, varName, goType)
	}
	return out
}

func varEmptyConstructorK8sType(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	varName string,
	// The shape we want to construct a new thing for
	shape *awssdkmodel.Shape,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	goType := shape.GoTypeWithPkgName()
	keepPointer := (shape.Type == "list" || shape.Type == "map")
	goType = model.ReplacePkgName(goType, r.SDKAPIPackageName(), "svcapitypes", keepPointer)
	goTypeNoPkg := goType
	goPkg := ""
	hadPkg := false
	if strings.Contains(goType, ".") {
		parts := strings.Split(goType, ".")
		goTypeNoPkg = parts[1]
		goPkg = parts[0]
		hadPkg = true
	}
	renames := r.TypeRenames()
	altTypeName, renamed := renames[goTypeNoPkg]
	if renamed {
		goTypeNoPkg = altTypeName
	} else if hadPkg {
		cleanNames := names.New(goTypeNoPkg)
		goTypeNoPkg = cleanNames.Camel
	}
	goType = goTypeNoPkg
	if hadPkg {
		goType = goPkg + "." + goType
	}

	switch shape.Type {
	case "structure":
		if r.Config().HasEmptyShape(shape.ShapeName) {
			// f0 := map[string]*string{}
			out += fmt.Sprintf("%s%s := map[string]*string{}\n", indent, varName)
		} else {
			// f0 := &svcapitypes.BookData{}
			out += fmt.Sprintf("%s%s := &%s{}\n", indent, varName, goType)
		}
	case "list", "map":
		// f0 := []*string{}
		out += fmt.Sprintf("%s%s := %s{}\n", indent, varName, goType)
	default:
		// var f0 string
		out += fmt.Sprintf("%svar %s %s\n", indent, varName, goType)
	}
	return out
}

// setSDKForScalar returns the Go code that sets the value of a target variable
// or field to a scalar value. For target variables that are structs, we output
// the aws-sdk-go's common SetXXX() method. For everything else, we output
// normal assignment operations.
func setSDKForScalar(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The name of the Input SDK Shape member we're outputting for
	targetFieldName string,
	// The variable name that we want to set a value to
	targetVarName string,
	// The type of shape of the target variable
	targetVarType string,
	// The path to the field that we access our source value from
	sourceFieldPath string,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	shapeRef *awssdkmodel.ShapeRef,
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	setTo := sourceVarName
	shape := shapeRef.Shape
	if shape.Type == "timestamp" {
		setTo += ".Time"
	} else if shapeRef.UseIndirection() {
		setTo = "*" + setTo
	}
	if targetVarType == "structure" {
		out += fmt.Sprintf("%s%s.Set%s(%s)\n", indent, targetVarName, targetFieldName, setTo)
	} else {
		targetVarPath := targetVarName
		if targetFieldName != "" {
			targetVarPath += "." + targetFieldName
		}
		out += fmt.Sprintf("%s%s = %s\n", indent, targetVarPath, setTo)
	}
	return out
}
