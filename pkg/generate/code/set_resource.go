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

	"github.com/aws-controllers-k8s/pkg/names"
	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/fieldpath"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/util"
)

// SetResource returns the Go code that sets a CRD's field value to the value
// of an output shape's member fields.  Status fields are always updated.
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
// And assume an SDK Shape CreateRepositoryOutput that looks like this
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
//	sourceVarName:	resp
//	targetVarName:	ko.Status
//	indentLevel:	1
//
// Then this function should output something like this:
//
//	field0 := []*string{}
//	for _, iter0 := range resp.Authors {
//	    var elem0 string
//	    elem0 = *iter
//	    field0 = append(field0, &elem0)
//	}
//	ko.Status.Authors = field0
//	field1 := &svcapitypes.ImageData{}
//	field1.Location = resp.ImageData.Location
//	field1.Tag = resp.ImageData.Tag
//	ko.Status.ImageData = field1
//	ko.Status.Name = resp.Name
func SetResource(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The type of operation to look for the Output shape
	opType model.OpType,
	// String representing the name of the variable that we will grab the
	// Output shape from. This will likely be "resp" since in the templates
	// that call this method, the "source variable" is the response struct
	// returned by the aws-sdk-go's SDK API call corresponding to the Operation
	sourceVarName string,
	// String representing the name of the variable that we will be **setting**
	// with values we get from the Output shape. This will likely be
	// "ko.Status" since that is the name of the "target variable" that the
	// templates that call this method use.
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
	outputShape, _ := r.GetOutputShape(op)
	if outputShape == nil {
		return ""
	}

	// If the output shape has a list containing the resource,
	// then call setResourceReadMany to generate for-range loops.
	// Output shape will be a list for ReadMany operations or if
	// designated via output wrapper config.
	wrapperFieldPath := r.GetOutputWrapperFieldPath(op)
	if op == r.Ops.ReadMany {
		return setResourceReadMany(
			cfg, r,
			op, sourceVarName, targetVarName, indentLevel,
		)
	} else if wrapperFieldPath != nil {
		// fieldpath api requires fully-qualified path
		qwfp := fieldpath.FromString(op.OutputRef.ShapeName + "." + *wrapperFieldPath)
		for _, sref := range qwfp.IterShapeRefs(&op.OutputRef) {
			// if there's at least 1 list to unpack, call setResourceReadMany
			if sref.Shape.Type == "list" {
				return setResourceReadMany(
					cfg, r,
					op, sourceVarName, targetVarName, indentLevel,
				)
			}
		}
		sourceVarName += "." + *wrapperFieldPath
	} else {
		// If the wrapper field path is not specified in the config file and if
		// there is a single member shape and that member shape is a structure,
		// unwrap it.
		if outputShape.UsedAsOutput && len(outputShape.MemberRefs) == 1 {
			for memberName, memberRef := range outputShape.MemberRefs {
				if memberRef.Shape.Type == "structure" {
					sourceVarName += "." + memberName
					outputShape = memberRef.Shape
				}
			}
		}
	}

	out := "\n"
	indent := strings.Repeat("\t", indentLevel)

	// Recursively descend through the set of fields on the Output shape,
	// creating temporary variables, populating those temporary variables'
	// fields with further-nested fields as needed
	for memberIndex, memberName := range outputShape.MemberNames() {
		//TODO: (vijat@) should these field be renamed before looking them up in spec?
		sourceAdaptedVarName := sourceVarName + "." + memberName

		// Handle the special case of ARN for primary resource identifier
		if r.IsPrimaryARNField(memberName) {
			// if ko.Status.ACKResourceMetadata == nil {
			//     ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			// }
			out += fmt.Sprintf(
				"%sif %s.Status.ACKResourceMetadata == nil {\n",
				indent,
				targetVarName,
			)
			out += fmt.Sprintf(
				"%s\t%s.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}\n",
				indent,
				targetVarName,
			)
			out += fmt.Sprintf("%s}\n", indent)

			// if resp.BookArn != nil {
			//     ko.Status.ACKResourceMetadata.ARN = resp.BookArn
			// }
			out += fmt.Sprintf(
				"%sif %s != nil {\n",
				indent,
				sourceAdaptedVarName,
			)
			out += fmt.Sprintf(
				"%s\tarn := ackv1alpha1.AWSResourceName(*%s)\n",
				indent,
				sourceAdaptedVarName,
			)
			out += fmt.Sprintf(
				"%s\t%s.Status.ACKResourceMetadata.ARN = &arn\n",
				indent,
				targetVarName,
			)
			out += fmt.Sprintf("%s}\n", indent)
			continue
		}

		// Determine whether the input shape's field is in the Spec or the
		// Status struct and set the source variable appropriately.
		var f *model.Field
		var targetMemberShapeRef *awssdkmodel.ShapeRef
		targetAdaptedVarName := targetVarName

		// Handles field renames, if applicable
		fieldName := cfg.GetResourceFieldName(
			r.Names.Original,
			op.ExportedName,
			memberName,
		)
		inSpec, inStatus := r.HasMember(fieldName, op.ExportedName)
		if inSpec {
			targetAdaptedVarName += cfg.PrefixConfig.SpecField
			f = r.SpecFields[fieldName]
		} else if inStatus {
			targetAdaptedVarName += cfg.PrefixConfig.StatusField
			f = r.StatusFields[fieldName]
		} else {
			// TODO(jaypipes): check generator config for exceptions?
			continue
		}

		targetMemberShapeRef = f.ShapeRef

		// We may have some special instructions for how to handle setting the
		// field value...
		setCfg := f.GetSetterConfig(opType)

		if setCfg != nil && setCfg.Ignore {
			continue
		}

		sourceMemberShapeRef := outputShape.MemberRefs[memberName]
		if sourceMemberShapeRef.Shape == nil {
			// We may have some instructions to specially handle this field by
			// accessing a different Output shape member.
			//
			// As an example, let's say we're dealing with a resource called
			// Bucket that has a "Name" field. In the Output shape of the
			// ReadOne resource manager method, the field name is "Bucket",
			// though. There may be a generator.yaml file that looks like this:
			//
			// resources:
			//   Bucket:
			//     fields:
			//       Name:
			//         set:
			//          - method: ReadOne
			//            from: Bucket
			//
			// And this would indicate to the code generator that the Spec.Name
			// field should be set from the value of the Output shape's Bucket
			// field.
			if setCfg != nil && setCfg.From != nil {
				fp := fieldpath.FromString(*setCfg.From)
				sourceMemberShapeRef = fp.ShapeRef(sourceMemberShapeRef)
			}
			if sourceMemberShapeRef == nil || sourceMemberShapeRef.Shape == nil {
				// Technically this should not happen, so let's bail here if it
				// does...
				msg := fmt.Sprintf(
					"expected .Shape to not be nil for ShapeRef of memberName %s",
					memberName,
				)
				panic(msg)
			}
		}

		targetMemberShape := targetMemberShapeRef.Shape

		// fieldVarName is the name of the variable that is used for temporary
		// storage of complex member field values
		//
		// For struct fields, we want to output code sort of like this:
		//
		//   field0 := &svapitypes.ImageData{}
		//   if resp.ImageData.Location != nil {
		//	     field0.Location = resp.ImageData.Location
		//   }
		//   if resp.ImageData.Tag != nil {
		//       field0.Tag = resp.ImageData.Tag
		//   }
		//   r.ko.Status.ImageData = field0
		//   if resp.Name != nil {
		//	     r.ko.Status.Name = resp.Name
		//   }
		//
		// For list fields, we want to end up with something like this:
		//
		// field0 := []*svcapitypes.VpnGroupMembership{}
		// for _, iter0 := resp.CustomAvailabilityZone.VpnGroupMemberships {
		//     elem0 := &svcapitypes.VPNGroupMembership{}
		//     if iter0.VPNID != nil {
		//         elem0.VPNID = iter0.VPNID
		//     }
		//     field0 := append(field0, elem0)
		// }
		// ko.Status.VpnMemberships = field0

		out += fmt.Sprintf(
			"%sif %s != nil {\n", indent, sourceAdaptedVarName,
		)
		qualifiedTargetVar := fmt.Sprintf(
			"%s.%s", targetAdaptedVarName, f.Names.Camel,
		)

		switch targetMemberShape.Type {
		case "list", "structure", "map":
			{
				memberVarName := fmt.Sprintf("f%d", memberIndex)
				out += varEmptyConstructorK8sType(
					cfg, r,
					memberVarName,
					targetMemberShapeRef.Shape,
					indentLevel+1,
				)
				out += setResourceForContainer(
					cfg, r,
					f.Names.Camel,
					memberVarName,
					targetMemberShapeRef,
					setCfg,
					sourceAdaptedVarName,
					sourceMemberShapeRef,
					f.Names.Camel,
					opType,
					indentLevel+1,
				)
				out += setResourceForScalar(
					qualifiedTargetVar,
					memberVarName,
					sourceMemberShapeRef,
					indentLevel+1,
				)
			}
		default:
			if setCfg != nil && setCfg.From != nil {
				// We have some special instructions to pull the value from a
				// different field or member...
				sourceAdaptedVarName = sourceVarName + "." + *setCfg.From
			}
			out += setResourceForScalar(
				qualifiedTargetVar,
				sourceAdaptedVarName,
				sourceMemberShapeRef,
				indentLevel+1,
			)
		}
		out += fmt.Sprintf(
			"%s} else {\n", indent,
		)
		out += fmt.Sprintf(
			"%s%s%s.%s = nil\n", indent, indent,
			targetAdaptedVarName, f.Names.Camel,
		)
		out += fmt.Sprintf(
			"%s}\n", indent,
		)
	}
	return out
}

func ListMemberNameInReadManyOutput(
	r *model.CRD,
) string {
	// Find the element in the output shape that contains the list of
	// resources. This heuristic is simplistic (just look for the field with a
	// list type) but seems to be followed consistently by the aws-sdk-go for
	// List operations.
	for memberName, memberShapeRef := range r.Ops.ReadMany.OutputRef.Shape.MemberRefs {
		if memberShapeRef.Shape.Type == "list" {
			return memberName
		}
	}
	panic("List output shape had no field of type 'list'")
}

// setResourceReadMany sets the supplied target variable from the results of a
// List operation. This is a special-case handling of those APIs where there is
// no ReadOne operation and instead the only way to grab information for a
// single object is to call the ReadMany/List operation with one of more
// filtering fields and then look for one element in the returned array of
// results and unpack that into the target variable.
//
// As an example, for the DescribeCacheClusters Elasticache API call, the
// returned code looks like this:
//
// Note: "resp" is the source variable and represents the
//
//	DescribeCacheClustersOutput shape/struct in the aws-sdk-go API for
//	Elasticache
//
// Note: "ko" is the target variable and represents the thing we'll be
//
//			 setting fields on
//
//	 if len(resp.CacheClusters) == 0 {
//	     return nil, ackerr.NotFound
//	 }
//	 found := false
//	 for _, elem := range resp.CacheClusters {
//	     if elem.ARN != nil {
//	         if ko.Status.ACKResourceMetadata == nil {
//	             ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
//	         }
//	         tmpARN := ackv1alpha1.AWSResourceName(*elemARN)
//	         ko.Status.ACKResourceMetadata.ARN = &tmpARN
//	     }
//	     if elem.AtRestEncryptionEnabled != nil {
//	         ko.Status.AtRestEncryptionEnabled = elem.AtRestEncryptionEnabled
//	     } else {
//	         ko.Status.AtRestEncryptionEnabled = nil
//	     }
//	     ...
//	     if elem.CacheClusterId != nil {
//	         if ko.Spec.CacheClusterID != nil {
//	             if *elem.CacheClusterId != *ko.Spec.CacheClusterID {
//	                 continue
//	             }
//	         }
//	         r.ko.Spec.CacheClusterID = elem.CacheClusterId
//	     } else {
//	         r.ko.Spec.CacheClusterID = nil
//	     }
//	     found = true
//	 }
//	 if !found {
//	     return nil, ackerr.NotFound
//	 }
func setResourceReadMany(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The ReadMany operation descriptor
	op *awssdkmodel.Operation,
	// String representing the name of the variable that we will grab the
	// Output shape from. This will likely be "resp" since in the templates
	// that call this method, the "source variable" is the response struct
	// returned by the aws-sdk-go's SDK API call corresponding to the Operation
	sourceVarName string,
	// String representing the name of the variable that we will be **setting**
	// with values we get from the Output shape. This will likely be
	// "ko" since that is the name of the "target variable" that the
	// templates that call this method use.
	targetVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	outputShape := op.OutputRef.Shape
	if outputShape == nil {
		return ""
	}

	out := "\n"
	indent := strings.Repeat("\t", indentLevel)

	listShapeName := ""
	var sourceElemShape *awssdkmodel.Shape

	// Find the element in the output shape that contains the list of
	// resources:
	// Check if there's a wrapperFieldPath, which will
	// point directly to the shape.
	wrapperFieldPath := r.GetOutputWrapperFieldPath(op)
	if wrapperFieldPath != nil {
		// fieldpath API needs fully qualified name
		wfp := fieldpath.FromString(outputShape.ShapeName + "." + *wrapperFieldPath)
		wfpShapeRef := wfp.ShapeRef(&op.OutputRef)
		if wfpShapeRef != nil {
			listShapeName = wfpShapeRef.ShapeName
			sourceElemShape = wfpShapeRef.Shape.MemberRef.Shape
		}
	}

	// If listShape can't be found using wrapperFieldPath,
	// then fall back to looking for the first field with a list type;
	// this heuristic seems to work for most list operations in aws-sdk-go.
	if listShapeName == "" {
		for memberName, memberShapeRef := range outputShape.MemberRefs {
			if memberShapeRef.Shape.Type == "list" {
				listShapeName = memberName
				sourceElemShape = memberShapeRef.Shape.MemberRef.Shape
				break
			}
		}
	}

	if listShapeName == "" {
		panic("List output shape had no field of type 'list'")
	}

	// Set of field names in the element shape that, if the generator config
	// instructs us to, we will write Go code to filter results of the List
	// operation by checking for matching values in these fields.
	matchFieldNames := r.ListOpMatchFieldNames()

	for _, mfName := range matchFieldNames {
		if inSpec, inStat := r.HasMember(mfName, op.ExportedName); !inSpec && !inStat {
			msg := fmt.Sprintf(
				"Match field name %s is not in %s Spec or Status fields",
				mfName, r.Names.Camel,
			)
			panic(msg)
		}
	}

	// found := false
	out += fmt.Sprintf("%sfound := false\n", indent)
	elemVarName := "elem"
	pathToShape := listShapeName
	if wrapperFieldPath != nil {
		pathToShape = *wrapperFieldPath
	}

	// for _, elem := range resp.CacheClusters {
	opening, closing, flIndentLvl := generateForRangeLoops(&op.OutputRef, pathToShape, sourceVarName, elemVarName, indentLevel)
	innerForIndent := strings.Repeat("\t", flIndentLvl)
	out += opening

	for memberIndex, memberName := range sourceElemShape.MemberNames() {
		sourceMemberShapeRef := sourceElemShape.MemberRefs[memberName]
		sourceMemberShape := sourceMemberShapeRef.Shape
		sourceAdaptedVarName := elemVarName + "." + memberName
		if r.IsPrimaryARNField(memberName) {
			out += fmt.Sprintf(
				"%sif %s != nil {\n", innerForIndent, sourceAdaptedVarName,
			)
			//     if ko.Status.ACKResourceMetadata == nil {
			//  	   ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			//     }
			out += fmt.Sprintf(
				"%s\tif %s.Status.ACKResourceMetadata == nil {\n",
				innerForIndent, targetVarName,
			)
			out += fmt.Sprintf(
				"%s\t\t%s.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}\n",
				innerForIndent, targetVarName,
			)
			out += fmt.Sprintf(
				"\t%s}\n", innerForIndent,
			)
			//          tmpARN := ackv1alpha1.AWSResourceName(*elemARN)
			//  		ko.Status.ACKResourceMetadata.ARN = &tmpARN
			out += fmt.Sprintf(
				"%s\ttmpARN := ackv1alpha1.AWSResourceName(*%s)\n",
				innerForIndent,
				sourceAdaptedVarName,
			)
			out += fmt.Sprintf(
				"%s\t%s.Status.ACKResourceMetadata.ARN = &tmpARN\n",
				innerForIndent,
				targetVarName,
			)
			out += fmt.Sprintf(
				"%s}\n", innerForIndent,
			)
			continue
		}
		// Determine whether the input shape's field is in the Spec or the
		// Status struct and set the source variable appropriately.
		var f *model.Field
		var targetMemberShapeRef *awssdkmodel.ShapeRef
		targetAdaptedVarName := targetVarName

		// Handles field renames, if applicable
		fieldName := cfg.GetResourceFieldName(
			r.Names.Original,
			op.ExportedName,
			memberName,
		)
		inSpec, inStatus := r.HasMember(fieldName, op.ExportedName)
		if inSpec {
			targetAdaptedVarName += cfg.PrefixConfig.SpecField
			f = r.SpecFields[fieldName]
		} else if inStatus {
			targetAdaptedVarName += cfg.PrefixConfig.StatusField
			f = r.StatusFields[fieldName]
		} else {
			// field not found in Spec or Status
			continue
		}

		// We may have some special instructions for how to handle setting the
		// field value...
		setCfg := f.GetSetterConfig(model.OpTypeList)

		if setCfg != nil && setCfg.Ignore {
			continue
		}

		targetMemberShapeRef = f.ShapeRef
		out += fmt.Sprintf(
			"%sif %s != nil {\n", innerForIndent, sourceAdaptedVarName,
		)

		//ex: r.ko.Spec.CacheClusterID
		qualifiedTargetVar := fmt.Sprintf(
			"%s.%s", targetAdaptedVarName, f.Names.Camel,
		)
		switch sourceMemberShape.Type {
		case "list", "structure", "map":
			{
				memberVarName := fmt.Sprintf("f%d", memberIndex)
				out += varEmptyConstructorK8sType(
					cfg, r,
					memberVarName,
					targetMemberShapeRef.Shape,
					flIndentLvl+1,
				)
				out += setResourceForContainer(
					cfg, r,
					f.Names.Camel,
					memberVarName,
					targetMemberShapeRef,
					setCfg,
					sourceAdaptedVarName,
					sourceMemberShapeRef,
					f.Names.Camel,
					model.OpTypeList,
					flIndentLvl+1,
				)
				out += setResourceForScalar(
					qualifiedTargetVar,
					memberVarName,
					sourceMemberShapeRef,
					flIndentLvl+1,
				)
			}
		default:
			//          if ko.Spec.CacheClusterID != nil {
			//              if *elem.CacheClusterId != *ko.Spec.CacheClusterID {
			//                  continue
			//              }
			//          }
			if util.InStrings(fieldName, matchFieldNames) {
				out += fmt.Sprintf(
					"%s\tif %s.%s != nil {\n",
					innerForIndent,
					targetAdaptedVarName,
					f.Names.Camel,
				)
				out += fmt.Sprintf(
					"%s\t\tif *%s != *%s.%s {\n",
					innerForIndent,
					sourceAdaptedVarName,
					targetAdaptedVarName,
					f.Names.Camel,
				)
				out += fmt.Sprintf(
					"%s\t\t\tcontinue\n", innerForIndent,
				)
				out += fmt.Sprintf(
					"%s\t\t}\n", innerForIndent,
				)
				out += fmt.Sprintf(
					"%s\t}\n", innerForIndent,
				)
			}
			//          r.ko.Spec.CacheClusterID = elem.CacheClusterId
			out += setResourceForScalar(
				qualifiedTargetVar,
				sourceAdaptedVarName,
				sourceMemberShapeRef,
				flIndentLvl+1,
			)
		}
		out += fmt.Sprintf(
			"%s} else {\n", innerForIndent,
		)
		out += fmt.Sprintf(
			"%s\t%s.%s = nil\n", innerForIndent,
			targetAdaptedVarName, f.Names.Camel,
		)
		out += fmt.Sprintf(
			"%s}\n", innerForIndent,
		)
	}
	// When we don't have custom matching/filtering logic for the list
	// operation, we just take the first element in the returned slice
	// of objects. When we DO have match fields, the generated Go code
	// above will output a `continue` when the required fields don't
	// match. Thus, we will break here only when getting a record where
	// all match fields have matched.
	out += fmt.Sprintf(
		"%sfound = true\n", innerForIndent,
	)

	// End of for-range loops
	out += closing

	//  if !found {
	//      return nil, ackerr.NotFound
	//  }
	out += fmt.Sprintf("%sif !found {\n", indent)
	out += fmt.Sprintf("%s\t%s\n", indent, cfg.SetManyOutputNotFoundErrReturn)
	out += fmt.Sprintf("%s}\n", indent)
	return out
}

// ackResourceMetadataGuardConstructor returns Go code representing a nil-guard
// and constructor for an ACKResourceMetadata struct:
//
//	if ko.Status.ACKResourceMetadata == nil {
//	    ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
//	}
func ackResourceMetadataGuardConstructor(
	// String representing the name of the variable that we will be **setting**
	// with values we get from the Output shape. This will likely be
	// "ko.Status" since that is the name of the "target variable" that the
	// templates that call this method use.
	targetVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	indent := strings.Repeat("\t", indentLevel)
	out := fmt.Sprintf(
		"%sif %s.ACKResourceMetadata == nil {\n",
		indent,
		targetVarName,
	)
	out += fmt.Sprintf(
		"%s\t%s.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}\n",
		indent,
		targetVarName,
	)
	out += fmt.Sprintf("%s}\n", indent)
	return out
}

// identifierNameOrIDGuardConstructor returns Go code representing a nil-guard
// and returns a `MissingNameIdentifier` error:
//
//	if identifier.NameOrID == "" {
//	 return ackerrors.MissingNameIdentifier
//	}
func identifierNameOrIDGuardConstructor(
	// String representing the name of the identifier that should have the
	// `NameOrID` pointer defined
	sourceVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	indent := strings.Repeat("\t", indentLevel)
	out := fmt.Sprintf("%sif %s.NameOrID == \"\" {\n", indent, sourceVarName)
	out += fmt.Sprintf("%s\treturn ackerrors.MissingNameIdentifier\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	return out
}

// SetResourceGetAttributes returns the Go code that sets the Status fields
// from the Output shape returned from a resource's GetAttributes operation.
//
// As an example, for the GetTopicAttributes SNS API call, the returned code
// looks like this:
//
//	if ko.Status.ACKResourceMetadata == nil {
//	    ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
//	}
//
// ko.Status.EffectiveDeliveryPolicy = resp.Attributes["EffectiveDeliveryPolicy"]
// ko.Status.ACKResourceMetadata.OwnerAccountID = ackv1alpha1.AWSAccountID(resp.Attributes["Owner"])
// ko.Status.ACKResourceMetadata.ARN = ackv1alpha1.AWSResourceName(resp.Attributes["TopicArn"])
func SetResourceGetAttributes(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// String representing the name of the variable that we will grab the
	// Output shape from. This will likely be "resp" since in the templates
	// that call this method, the "source variable" is the response struct
	// returned by the aws-sdk-go's SDK API call corresponding to the Operation
	sourceVarName string,
	// String representing the name of the variable that we will be **setting**
	// with values we get from the Output shape. This will likely be
	// "ko.Status" since that is the name of the "target variable" that the
	// templates that call this method use.
	targetVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	if !r.UnpacksAttributesMap() {
		// This is a bug in the code generation if this occurs...
		msg := fmt.Sprintf(
			"called SetResourceGetAttributes for a resource '%s' that doesn't unpack attributes map",
			r.Ops.GetAttributes.Name,
		)
		panic(msg)
	}
	op := r.Ops.GetAttributes
	if op == nil {
		return ""
	}
	inputShape := op.InputRef.Shape
	if inputShape == nil {
		return ""
	}

	out := "\n"
	indent := strings.Repeat("\t", indentLevel)

	// did we output an ACKResourceMetadata guard and constructor snippet?
	mdGuardOut := false
	fieldConfigs := cfg.GetFieldConfigs(r.Names.Original)
	sortedAttrFieldNames := []string{}
	for fName, fConfig := range fieldConfigs {
		if fConfig.IsAttribute {
			sortedAttrFieldNames = append(sortedAttrFieldNames, fName)
		}
	}
	sort.Strings(sortedAttrFieldNames)
	for _, fieldName := range sortedAttrFieldNames {
		adaptiveTargetVarName := targetVarName + cfg.PrefixConfig.StatusField
		if r.IsPrimaryARNField(fieldName) {
			if !mdGuardOut {
				out += ackResourceMetadataGuardConstructor(
					adaptiveTargetVarName, indentLevel,
				)
				mdGuardOut = true
			}
			out += fmt.Sprintf(
				"%stmpARN := ackv1alpha1.AWSResourceName(*%s.Attributes[\"%s\"])\n",
				indent,
				sourceVarName,
				fieldName,
			)
			out += fmt.Sprintf(
				"%s%s.ACKResourceMetadata.ARN = &tmpARN\n",
				indent,
				adaptiveTargetVarName,
			)
			continue
		}

		fieldConfig := fieldConfigs[fieldName]
		if fieldConfig.IsOwnerAccountID && cfg.IncludeACKMetadata {
			if !mdGuardOut {
				out += ackResourceMetadataGuardConstructor(
					adaptiveTargetVarName, indentLevel,
				)
				mdGuardOut = true
			}
			out += fmt.Sprintf(
				"%stmpOwnerID := ackv1alpha1.AWSAccountID(*%s.Attributes[\"%s\"])\n",
				indent,
				sourceVarName,
				fieldName,
			)
			out += fmt.Sprintf(
				"%s%s.ACKResourceMetadata.OwnerAccountID = &tmpOwnerID\n",
				indent,
				adaptiveTargetVarName,
			)
			continue
		}

		fieldNames := names.New(fieldName)
		if !fieldConfig.IsReadOnly {
			adaptiveTargetVarName = targetVarName + cfg.PrefixConfig.SpecField
		}
		out += fmt.Sprintf(
			"%s%s.%s = %s.Attributes[\"%s\"]\n",
			indent,
			adaptiveTargetVarName,
			fieldNames.Camel,
			sourceVarName,
			fieldName,
		)
	}
	return out
}

// SetResourceIdentifiers returns the Go code that sets an empty CR object with
// Spec and Status field values that correspond to the primary identifier (be
// that an ARN, ID or Name) and any other "additional keys" required for the AWS
// service to uniquely identify the object.
//
// The method will attempt to look for the field denoted with a value of true
// for `is_primary_key`, or will use the ARN if the resource has a value of true
// for `is_arn_primary_key`. Otherwise, the method will attempt to use the
// `ReadOne` operation, if present, falling back to using `ReadMany`.
// If it detects the operation uses an ARN to identify the resource it will read
// it from the metadata status field. Otherwise it will use any field with a
// name that matches the primary identifier from the operation, pulling from
// top-level spec or status fields.
//
// An example of code with no additional keys:
//
// ```
//
//	if identifier.NameOrID == nil {
//		return ackerrors.MissingNameIdentifier
//	}
//	r.ko.Status.BrokerID = identifier.NameOrID
//
// ```
//
// An example of code with additional keys:
//
// ```
//
//	if identifier.NameOrID == nil {
//		  return ackerrors.MissingNameIdentifier
//	}
//
// r.ko.Spec.ResourceID = identifier.NameOrID
//
// f0, f0ok := identifier.AdditionalKeys["scalableDimension"]
//
//	if f0ok {
//		  r.ko.Spec.ScalableDimension = f0
//	}
//
// f1, f1ok := identifier.AdditionalKeys["serviceNamespace"]
//
//	if f1ok {
//		  r.ko.Spec.ServiceNamespace = f1
//	}
//
// ```
// An example of code that uses the ARN:
//
// ```
//
//	if r.ko.Status.ACKResourceMetadata == nil {
//		r.ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
//	}
//
// r.ko.Status.ACKResourceMetadata.ARN = identifier.ARN
//
// f0, f0ok := identifier.AdditionalKeys["modelPackageName"]
//
//	if f0ok {
//		r.ko.Spec.ModelPackageName = &f0
//	}
//
// ```
func SetResourceIdentifiers(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// String representing the name of the variable that we will grab the Input
	// shape from. This will likely be "identifier" since in the templates that
	// call this method, the "source variable" is the CRD struct which is used
	// to populate the target variable, which is the struct of unique
	// identifiers
	sourceVarName string,
	// String representing the name of the variable that we will be **setting**
	// with values we get from the Output shape. This will likely be
	// "r.ko" since that is the name of the "target variable" that the
	// templates that call this method use for the Input shape.
	targetVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	op := r.Ops.ReadOne
	if op == nil {
		switch {
		case r.Ops.GetAttributes != nil:
			// If single lookups can only be done with GetAttributes
			op = r.Ops.GetAttributes
		case r.Ops.ReadMany != nil:
			// If single lookups can only be done using ReadMany
			op = r.Ops.ReadMany
		}
	}
	inputShape := op.InputRef.Shape
	if inputShape == nil {
		return ""
	}

	primaryKeyOut := ""
	additionalKeyOut := "\n"

	indent := strings.Repeat("\t", indentLevel)

	primaryKeyConditionalOut := "\n"
	primaryKeyConditionalOut += identifierNameOrIDGuardConstructor(sourceVarName, indentLevel)

	// if r.ko.Status.ACKResourceMetadata == nil {
	//  r.ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	// }
	// r.ko.Status.ACKResourceMetadata.ARN = identifier.ARN
	arnOut := "\n"
	arnOut += ackResourceMetadataGuardConstructor(fmt.Sprintf("%s.Status", targetVarName), indentLevel)
	arnOut += fmt.Sprintf(
		"%s%s.Status.ACKResourceMetadata.ARN = %s.ARN\n",
		indent, targetVarName, sourceVarName,
	)

	// Check if the CRD defines the primary keys
	if r.IsARNPrimaryKey() {
		return arnOut
	}
	primaryField, err := r.GetPrimaryKeyField()
	if err != nil {
		panic(err)
	}

	var primaryCRField, primaryShapeField string
	isPrimarySet := primaryField != nil
	if isPrimarySet {
		memberPath, _ := findFieldInCR(cfg, r, primaryField.Names.Original)
		targetVarPath := fmt.Sprintf("%s%s", targetVarName, memberPath)
		primaryKeyOut += setResourceIdentifierPrimaryIdentifier(cfg, r,
			primaryField,
			targetVarPath,
			sourceVarName,
			indentLevel)
	} else {
		primaryCRField, primaryShapeField = FindPrimaryIdentifierFieldNames(cfg, r, op)
		if primaryShapeField == PrimaryIdentifierARNOverride {
			return arnOut
		}
	}

	paginatorFieldLookup := []string{
		"NextToken",
		"MaxResults",
	}

	for memberIndex, memberName := range inputShape.MemberNames() {
		if util.InStrings(memberName, paginatorFieldLookup) {
			continue
		}

		inputShapeRef := inputShape.MemberRefs[memberName]
		inputMemberShape := inputShapeRef.Shape

		// Only strings and list of strings are currently accepted as valid
		// inputs for additional key fields
		if inputMemberShape.Type != "string" &&
			(inputMemberShape.Type != "list" ||
				inputMemberShape.MemberRef.Shape.Type != "string") {
			continue
		}

		if r.IsSecretField(memberName) {
			// Secrets cannot be used as fields in identifiers
			continue
		}

		if r.IsPrimaryARNField(memberName) {
			continue
		}

		// Handles field renames, if applicable
		fieldName := cfg.GetResourceFieldName(
			r.Names.Original,
			op.ExportedName,
			memberName,
		)

		// Check to see if we've already set the field as the primary identifier
		if isPrimarySet && fieldName == primaryField.Names.Camel {
			continue
		}

		isPrimaryIdentifier := fieldName == primaryShapeField

		searchField := ""
		if isPrimaryIdentifier {
			searchField = primaryCRField
		} else {
			searchField = fieldName
		}

		memberPath, targetField := findFieldInCR(cfg, r, searchField)
		if targetField == nil || (isPrimarySet && targetField == primaryField) {
			continue
		}

		switch targetField.ShapeRef.Shape.Type {
		case "list", "structure", "map":
			panic("primary identifier '" + targetField.Path + "' must be a scalar type since NameOrID is a string")
		default:
			break
		}

		targetVarPath := fmt.Sprintf("%s%s", targetVarName, memberPath)
		if isPrimaryIdentifier {
			primaryKeyOut += setResourceIdentifierPrimaryIdentifier(cfg, r,
				targetField,
				targetVarPath,
				sourceVarName,
				indentLevel)
		} else {
			additionalKeyOut += setResourceIdentifierAdditionalKey(
				cfg, r,
				memberIndex,
				targetField,
				targetVarPath,
				sourceVarName,
				names.New(fieldName).CamelLower,
				indentLevel)
		}
	}

	return primaryKeyConditionalOut + primaryKeyOut + additionalKeyOut
}

// findFieldInCR will search for a given field, by its name, in a CR and returns
// the member path and Field type if one is found.
func findFieldInCR(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The name of the field to search for
	searchField string,
) (memberPath string, targetField *model.Field) {
	specField, inSpec := r.SpecFields[searchField]
	statusField, inStatus := r.StatusFields[searchField]
	switch {
	case inSpec:
		memberPath = cfg.PrefixConfig.SpecField
		targetField = specField
	case inStatus:
		memberPath = cfg.PrefixConfig.StatusField
		targetField = statusField
	default:
		return "", nil
	}
	return memberPath, targetField
}

// setResourceIdentifierPrimaryIdentifier returns a string of Go code that sets
// the primary identifier Spec or Status field on a given resource to the value
// in the identifier `NameOrID` field:
//
// r.ko.Status.BrokerID = &identifier.NameOrID
func setResourceIdentifierPrimaryIdentifier(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The field that will be set on the target variable
	targetField *model.Field,
	// The variable name that we want to set a value to
	targetVarName string,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	adaptedMemberPath := fmt.Sprintf("&%s.NameOrID", sourceVarName)
	qualifiedTargetVar := fmt.Sprintf("%s.%s", targetVarName, targetField.Path)

	return setResourceForScalar(
		qualifiedTargetVar,
		adaptedMemberPath,
		targetField.ShapeRef,
		indentLevel,
	)
}

// setResourceIdentifierAdditionalKey returns a string of Go code that sets a
// Spec or Status field on a given resource to the value in the identifier's
// `AdditionalKeys` mapping:
//
// f0, f0ok := identifier.AdditionalKeys["scalableDimension"]
//
//	if f0ok {
//		r.ko.Spec.ScalableDimension = f0
//	}
func setResourceIdentifierAdditionalKey(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	fieldIndex int,
	// The field that will be set on the target variable
	targetField *model.Field,
	// The variable name that we want to set a value to
	targetVarName string,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	// The key in the `AdditionalKeys` map storing the source variable
	sourceVarKey string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	indent := strings.Repeat("\t", indentLevel)

	additionalKeyOut := ""

	fieldIndexName := fmt.Sprintf("f%d", fieldIndex)
	sourceAdaptedVarName := fmt.Sprintf("%s.AdditionalKeys[\"%s\"]", sourceVarName, sourceVarKey)

	// TODO(RedbackThomson): If the identifiers don't exist, we should be
	// throwing an error accessible to the user
	additionalKeyOut += fmt.Sprintf("%s%s, %sok := %s\n", indent, fieldIndexName, fieldIndexName, sourceAdaptedVarName)
	additionalKeyOut += fmt.Sprintf("%sif %sok {\n", indent, fieldIndexName)
	qualifiedTargetVar := fmt.Sprintf("%s.%s", targetVarName, targetField.Path)
	additionalKeyOut += setResourceForScalar(
		qualifiedTargetVar,
		fmt.Sprintf("&%s", fieldIndexName),
		targetField.ShapeRef,
		indentLevel+1,
	)
	additionalKeyOut += fmt.Sprintf("%s}\n", indent)

	return additionalKeyOut
}

// setResourceForContainer returns a string of Go code that sets the value of a
// target variable to that of a source variable. When the source variable type
// is a map, struct or slice type, then this function is called recursively on
// the elements or members of the source variable.
func setResourceForContainer(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The name of the CR field we're outputting for
	targetFieldName string,
	// The variable name that we want to set a value to
	targetVarName string,
	// Shape Ref of the target struct field
	targetShapeRef *awssdkmodel.ShapeRef,
	// SetFieldConfig of the *target* field
	targetSetCfg *ackgenconfig.SetFieldConfig,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	// ShapeRef of the source struct field
	sourceShapeRef *awssdkmodel.ShapeRef,
	// targetFieldPath is the field path for the target containing field
	targetFieldPath string,
	op model.OpType,
	indentLevel int,
) string {
	switch sourceShapeRef.Shape.Type {
	case "structure":
		return SetResourceForStruct(
			cfg, r,
			targetVarName,
			targetShapeRef,
			targetSetCfg,
			sourceVarName,
			sourceShapeRef,
			targetFieldPath,
			op,
			indentLevel,
		)
	case "list":
		return setResourceForSlice(
			cfg, r,
			targetFieldName,
			targetVarName,
			targetShapeRef,
			targetSetCfg,
			sourceVarName,
			sourceShapeRef,
			targetFieldPath,
			op,
			indentLevel,
		)
	case "map":
		return setResourceForMap(
			cfg, r,
			targetFieldName,
			targetVarName,
			targetShapeRef,
			targetSetCfg,
			sourceVarName,
			sourceShapeRef,
			targetFieldPath,
			op,
			indentLevel,
		)
	default:
		return setResourceForScalar(
			fmt.Sprintf("%s.%s", targetFieldName, targetVarName),
			sourceVarName,
			sourceShapeRef,
			indentLevel,
		)
	}
}

// SetResourceForStruct returns a string of Go code that sets a target variable
// value to a source variable when the type of the source variable is a struct.
func SetResourceForStruct(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The variable name that we want to set a value to
	targetVarName string,
	// Shape Ref of the target struct field
	targetShapeRef *awssdkmodel.ShapeRef,
	// SetFieldConfig of the *target* field
	targetSetCfg *ackgenconfig.SetFieldConfig,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	// ShapeRef of the source struct field
	sourceShapeRef *awssdkmodel.ShapeRef,
	// targetFieldPath is the field path to targetFieldName
	targetFieldPath string,
	op model.OpType,
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	sourceShape := sourceShapeRef.Shape
	targetShape := targetShapeRef.Shape

	var sourceMemberShapeRef *awssdkmodel.ShapeRef
	var sourceAdaptedVarName, qualifiedTargetVar string

	for _, targetMemberName := range targetShape.MemberNames() {
		sourceMemberShapeRef = sourceShape.MemberRefs[targetMemberName]
		if sourceMemberShapeRef == nil {
			continue
		}
		// Upstream logic iterates over sourceShape members and therefore uses
		// the sourceShape's index; continue using sourceShape's index here for consistency.
		sourceMemberIndex, err := GetMemberIndex(sourceShape, targetMemberName)
		if err != nil {
			msg := fmt.Sprintf(
				"could not determine source shape index: %v", err)
			panic(msg)
		}

		targetMemberShapeRef := targetShape.MemberRefs[targetMemberName]
		indexedVarName := fmt.Sprintf("%sf%d", targetVarName, sourceMemberIndex)
		sourceMemberShape := sourceMemberShapeRef.Shape
		targetMemberCleanNames := names.New(targetMemberName)
		sourceAdaptedVarName = sourceVarName + "." + targetMemberName
		out += fmt.Sprintf(
			"%sif %s != nil {\n", indent, sourceAdaptedVarName,
		)
		qualifiedTargetVar = fmt.Sprintf(
			"%s.%s", targetVarName, targetMemberCleanNames.Camel,
		)
		updatedTargetFieldPath := targetFieldPath + "." + targetMemberCleanNames.Camel

		switch sourceMemberShape.Type {
		case "list", "structure", "map":
			{
				out += varEmptyConstructorK8sType(
					cfg, r,
					indexedVarName,
					targetMemberShapeRef.Shape,
					indentLevel+1,
				)
				out += setResourceForContainer(
					cfg, r,
					targetMemberCleanNames.Camel,
					indexedVarName,
					targetMemberShapeRef,
					nil,
					sourceAdaptedVarName,
					sourceMemberShapeRef,
					updatedTargetFieldPath,
					op,
					indentLevel+1,
				)
				out += setResourceForScalar(
					qualifiedTargetVar,
					indexedVarName,
					sourceMemberShapeRef,
					indentLevel+1,
				)
			}
		default:
			out += setResourceForScalar(
				qualifiedTargetVar,
				sourceAdaptedVarName,
				sourceMemberShapeRef,
				indentLevel+1,
			)
		}
		out += fmt.Sprintf(
			"%s}\n", indent,
		)
	}
	if len(targetShape.MemberNames()) == 0 {
		// This scenario can occur when the targetShape is a primitive, but
		// the sourceShape is a struct. For example, EC2 resource DHCPOptions
		// has a field NewDhcpConfiguration.Values(targetShape = string) whose name
		// aligns with DhcpConfiguration.Values(sourceShape = AttributeValue).
		// Although the names correspond, the shapes/types are different and the intent
		// is to set NewDhcpConfiguration.Values using DhcpConfiguration.Values.Value
		// (AttributeValue.Value) shape instead. This behavior can be configured using
		// SetConfig.

		// Check if target field has a SetConfig, validate SetConfig.From points
		// to a shape within sourceShape, and generate Go code using
		// said shape. Using the example above, SetConfig is set
		// for NewDhcpConfiguration.Values and Setconfig.From points
		// to AttributeValue.Value (string), which leads to generating Go
		// code referencing DhcpConfiguration.Values.Value instead of 'Values'.

		if targetField, ok := r.Fields[targetFieldPath]; ok {
			setCfg := targetField.GetSetterConfig(op)
			if setCfg != nil && setCfg.From != nil {
				fp := fieldpath.FromString(*setCfg.From)
				sourceMemberShapeRef = fp.ShapeRef(sourceShapeRef)
				if sourceMemberShapeRef != nil && sourceMemberShapeRef.Shape != nil {
					names := names.New(sourceMemberShapeRef.LocationName)
					sourceAdaptedVarName = sourceVarName + "." + names.Camel
					out += fmt.Sprintf(
						"%sif %s != nil {\n", indent, sourceAdaptedVarName,
					)
					qualifiedTargetVar = targetVarName

					// Use setResourceForScalar and dereference sourceAdaptedVarName
					// because primitives are being set.
					sourceAdaptedVarName = "*" + sourceAdaptedVarName
					out += setResourceForScalar(
						qualifiedTargetVar,
						sourceAdaptedVarName,
						sourceMemberShapeRef,
						indentLevel+1,
					)
					out += fmt.Sprintf(
						"%s}\n", indent,
					)
				}
			}
		}
	}
	return out
}

// setResourceForSlice returns a string of Go code that sets a target variable
// value to a source variable when the type of the source variable is a slice.
func setResourceForSlice(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The name of the CR field we're outputting for
	targetFieldName string,
	// The variable name that we want to set a value to
	targetVarName string,
	// Shape Ref of the target slice field
	targetShapeRef *awssdkmodel.ShapeRef,
	// SetFieldConfig of the *target* field
	targetSetCfg *ackgenconfig.SetFieldConfig,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	// ShapeRef of the source slice field
	sourceShapeRef *awssdkmodel.ShapeRef,
	// targetFieldPath is the field path to targetFieldName
	targetFieldPath string,
	op model.OpType,
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	sourceShape := sourceShapeRef.Shape
	targetShape := targetShapeRef.Shape
	iterVarName := fmt.Sprintf("%siter", targetVarName)
	elemVarName := fmt.Sprintf("%selem", targetVarName)
	// for _, f0iter0 := range resp.TagSpecifications {
	out += fmt.Sprintf("%sfor _, %s := range %s {\n", indent, iterVarName, sourceVarName)
	//		var f0elem0 string
	out += varEmptyConstructorK8sType(
		cfg, r,
		elemVarName,
		targetShape.MemberRef.Shape,
		indentLevel+1,
	)
	// We may have some instructions to specially handle this field by
	// accessing a different Output shape member.
	//
	// As an example, let's say we're dealing with a resource called
	// DBInstance that has a "DBSecurityGroups" field of type `[]*string`.
	// In the Output shape of the ReadOne resource manager method, the
	// field name is "DBSecurityGroups" of type
	// `[]*DBSecurityGroupMembership`, though.  This
	// `DBSecurityGroupMembership` struct contains two fields,
	// `DBSecurityGroupName` and `DBSecurityGroupStatus` and we wish to
	// grab only the value of the `DBSecurityGroupName` fields and set the
	// `Spec.DBSecurityGroups` field to the set of those field values.
	//
	// There may be a generator.yaml file that looks like this:
	//
	// resources:
	//   DBInstance:
	//     fields:
	//       DBSecurityGroups:
	//         set:
	//          - method: ReadOne
	//            from: DBSecurityGroupName
	//
	// And this would indicate to the code generator that the
	// Spec.DBSecurityGroups field should be set from the value of the
	// Output shape's DBSecurityGroups..DBSecurityGroupName fields.
	if targetSetCfg != nil && targetSetCfg.From != nil {
		if sourceMemberShapeRef, found := sourceShape.MemberRef.Shape.MemberRefs[*targetSetCfg.From]; found {
			out += setResourceForScalar(
				elemVarName,
				fmt.Sprintf("*%s.%s", iterVarName, *targetSetCfg.From),
				sourceMemberShapeRef,
				indentLevel+1,
			)
		} else {
			// This is a bug in the code generation if this occurs...
			msg := fmt.Sprintf(
				"could not find targetSetCfg.From %s in struct member type.",
				*targetSetCfg.From,
			)
			panic(msg)
		}
	} else {
		//  f0elem0 = *f0iter0
		//
		// or
		//
		//  f0elem0.SetMyField(*f0iter0)
		containerFieldName := ""
		if sourceShape.MemberRef.Shape.Type == "structure" {
			containerFieldName = targetFieldName
		}
		out += setResourceForContainer(
			cfg, r,
			containerFieldName,
			elemVarName,
			&targetShape.MemberRef,
			targetSetCfg,
			iterVarName,
			&sourceShape.MemberRef,
			targetFieldPath,
			op,
			indentLevel+1,
		)
	}
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

// setResourceForMap returns a string of Go code that sets a target variable
// value to a source variable when the type of the source variable is a map.
func setResourceForMap(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The name of the CR field we're outputting for
	targetFieldName string,
	// The variable name that we want to set a value to
	targetVarName string,
	// Shape Ref of the target map field
	targetShapeRef *awssdkmodel.ShapeRef,
	// SetFieldConfig of the *target* field
	targetSetCfg *ackgenconfig.SetFieldConfig,
	// The struct or struct field that we access our source value from
	sourceVarName string,
	// ShapeRef of the source map field
	sourceShapeRef *awssdkmodel.ShapeRef,
	// targetFieldPath is the field path to targetFieldName
	targetFieldPath string,
	op model.OpType,
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	sourceShape := sourceShapeRef.Shape
	targetShape := targetShapeRef.Shape

	valIterVarName := fmt.Sprintf("%svaliter", targetVarName)
	keyVarName := fmt.Sprintf("%skey", targetVarName)
	valVarName := fmt.Sprintf("%sval", targetVarName)
	// for f0key, f0valiter := range resp.Tags {
	out += fmt.Sprintf("%sfor %s, %s := range %s {\n", indent, keyVarName, valIterVarName, sourceVarName)
	//		f0elem := string{}
	out += varEmptyConstructorK8sType(
		cfg, r,
		valVarName,
		targetShape.ValueRef.Shape,
		indentLevel+1,
	)
	//  f0val = *f0valiter
	containerFieldName := ""
	if sourceShape.ValueRef.Shape.Type == "structure" {
		containerFieldName = targetFieldName
	}
	out += setResourceForContainer(
		cfg, r,
		containerFieldName,
		valVarName,
		&targetShape.ValueRef,
		nil,
		valIterVarName,
		&sourceShape.ValueRef,
		targetFieldPath,
		op,
		indentLevel+1,
	)
	addressOfVar := ""
	switch sourceShape.ValueRef.Shape.Type {
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

// setResourceForScalar returns a string of Go code that sets a target variable
// value to a source variable when the type of the source variable is a scalar
// type (not a map, slice or struct).
func setResourceForScalar(
	// The fully-qualified variable that will be set to sourceVar
	targetVar string,
	// The struct or struct field that we access our source value from
	sourceVar string,
	shapeRef *awssdkmodel.ShapeRef,
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	setTo := sourceVar
	shape := shapeRef.Shape
	if shape.Type == "timestamp" {
		setTo = "&metav1.Time{*" + sourceVar + "}"
	}
	if strings.HasPrefix(targetVar, ".") {
		targetVar = targetVar[1:]
		setTo = "*" + setTo
	}
	out += fmt.Sprintf("%s%s = %s\n", indent, targetVar, setTo)
	return out
}

// generateForRangeLoops returns strings of Go code and an int
// representing indentLevel of the inner-most for loop + 1.
// This function unpacks a collection from a shapeRef
// using path and builds the opening and closing
// pieces of a for-range loop written in Go in order
// to access the element contained within the collection.
// The name of this element value is designated with outputVarName:
// ex: 'for _, <outputVarName> := ...'
// Limitations: path supports lists and structs only and
// the 'break' is coupled with for-range loop to only take
// first element of a list.
//
// Sample Input:
// - shapeRef: DescribeInstancesOutputShape
// - path: Reservations.Instances
// - sourceVarName: resp
// - outputVarName: elem
// - indentLevel: 1
//
// Sample Output (omit formatting for readability):
//   - opening: "for _, iter0 := range resp.Reservations {
//     for _, elem := range iter0.Instances {"
//   - closing: "break } break }"
//   - updatedIndentLevel: 3
func generateForRangeLoops(
	// shapeRef of the shape containing element
	shapeRef *awssdkmodel.ShapeRef,
	// path is the path to the element relative to shapeRef
	path string,
	// sourceVarName is the name of struct or field used to access source value
	sourceVarName string,
	// outputVarName is the desired name of the element, once unwrapped
	outputVarName string,
	indentLevel int,
) (string, string, int) {
	opening, closing := "", ""
	updatedIndentLevel := indentLevel

	fp := fieldpath.FromString(path)
	unwrapCount := 0
	iterVarName := fmt.Sprintf("iter%d", unwrapCount)
	collectionVarName := sourceVarName
	unpackShape := shapeRef.Shape

	for fp.Size() > 0 {
		pathPart := fp.PopFront()
		partShapeRef, _ := unpackShape.MemberRefs[pathPart]
		unpackShape = partShapeRef.Shape
		indent := strings.Repeat("\t", updatedIndentLevel)
		iterVarName = fmt.Sprintf("iter%d", unwrapCount)
		collectionVarName += "." + pathPart

		// Using the fieldpath as a guide, unwrap the shapeRef
		// to generate for-range loops. If pathPart points
		// to a struct member, then simply append struct name
		// to collectionVarName and move on to unwrap the next pathPart/shape.
		// If pathPart points to a list member, then generate for-range loop
		// code and update collectionVarName, unpackShape, and updatedIndentLevel
		// for processing the next loop, if applicable.
		if partShapeRef.Shape.Type == "list" {
			// ex: for _, iter0 := range resp.Reservations {
			opening += fmt.Sprintf("%sfor _, %s := range %s {\n", indent, iterVarName, collectionVarName)
			// ex:
			//        break
			//    }
			closeLoop := fmt.Sprintf("%s\tbreak\n%s}\n", indent, indent)
			if closing != "" {
				// nested loops need to output inner most closing braces first
				closeLoop += closing
				closing = closeLoop
			} else {
				closing += closeLoop
			}
			// reference iterVarName in subsequent for-loop, if any
			collectionVarName = iterVarName
			unpackShape = partShapeRef.Shape.MemberRef.Shape
			updatedIndentLevel += 1
		}
		unwrapCount += 1
	}
	// replace inner-most range loop value's name with desired outputVarName
	opening = strings.Replace(opening, iterVarName, outputVarName, 1)
	return opening, closing, updatedIndentLevel
}

func FindInArray(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// The type of operation to look for the Output shape
	opType model.OpType,
	// String representing the name of the variable that we will grab the
	// Output shape from. This will likely be "resp" since in the templates
	// that call this method, the "source variable" is the response struct
	// returned by the aws-sdk-go's SDK API call corresponding to the Operation
	sourceVarName string,
	// String representing the name of the variable that we will be **setting**
	// with values we get from the Output shape. This will likely be
	// "ko.Status" since that is the name of the "target variable" that the
	// templates that call this method use.
	targetVarName string,
	matchWith string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := "\n"
	indent := strings.Repeat("\t", indentLevel)

	op := r.Ops.ReadMany
	matchFieldNames := r.ListOpMatchFieldNames()

	outputShape := op.OutputRef.Shape

	listShapeName := ""
	var sourceElemShape *awssdkmodel.Shape
	_ = sourceElemShape

	// Find the element in the output shape that contains the list of
	// resources:
	// Check if there's a wrapperFieldPath, which will
	// point directly to the shape.
	wrapperFieldPath := r.GetOutputWrapperFieldPath(op)
	if wrapperFieldPath != nil {
		// fieldpath API needs fully qualified name
		wfp := fieldpath.FromString(outputShape.ShapeName + "." + *wrapperFieldPath)
		wfpShapeRef := wfp.ShapeRef(&op.OutputRef)
		if wfpShapeRef != nil {
			listShapeName = wfpShapeRef.ShapeName
			sourceElemShape = wfpShapeRef.Shape.MemberRef.Shape
		}
	}

	// If listShape can't be found using wrapperFieldPath,
	// then fall back to looking for the first field with a list type;
	// this heuristic seems to work for most list operations in aws-sdk-go.
	if listShapeName == "" {
		for memberName, memberShapeRef := range outputShape.MemberRefs {
			if memberShapeRef.Shape.Type == "list" {
				listShapeName = memberName
				sourceElemShape = memberShapeRef.Shape.MemberRef.Shape
				break
			}
		}
	}

	if listShapeName == "" {
		panic("List output shape had no field of type 'list'")
	}

	for _, mfName := range matchFieldNames {
		if inSpec, inStat := r.HasMember(mfName, op.ExportedName); !inSpec && !inStat {
			msg := fmt.Sprintf(
				"Match field name %s is not in %s Spec or Status fields",
				mfName, r.Names.Camel,
			)
			panic(msg)
		}
	}

	elemVarName := "elem"
	pathToShape := listShapeName
	if wrapperFieldPath != nil {
		pathToShape = *wrapperFieldPath
	}

	// for _, elem := range resp.InternetGateways {
	opening, _, _ := generateForRangeLoops(&op.OutputRef, pathToShape, sourceVarName, elemVarName, indentLevel)
	out += opening

	for _, memberName := range sourceElemShape.MemberNames() {
		sourceMemberShapeRef := sourceElemShape.MemberRefs[memberName]
		sourceMemberShape := sourceMemberShapeRef.Shape
		sourceAdaptedVarName := elemVarName + "." + memberName
		// Determine whether the input shape's field is in the Spec or the
		// Status struct and set the source variable appropriately.
		var f *model.Field
		// var targetMemberShapeRef *awssdkmodel.ShapeRef
		targetAdaptedVarName := targetVarName

		// Handles field renames, if applicable
		fieldName := cfg.GetResourceFieldName(
			r.Names.Original,
			op.ExportedName,
			memberName,
		)
		inSpec, inStatus := r.HasMember(fieldName, op.ExportedName)
		if inSpec {
			targetAdaptedVarName += cfg.PrefixConfig.SpecField
			f = r.SpecFields[fieldName]
		} else if inStatus {
			targetAdaptedVarName += cfg.PrefixConfig.StatusField
			f = r.StatusFields[fieldName]
		} else {
			// field not found in Spec or Status
			continue
		}

		// We may have some special instructions for how to handle setting the
		// field value...
		setCfg := f.GetSetterConfig(model.OpTypeList)

		if setCfg != nil && setCfg.Ignore {
			continue
		}

		_ = f.ShapeRef

		//ex: r.ko.Spec.CacheClusterID
		_ = fmt.Sprintf(
			"%s.%s", targetAdaptedVarName, f.Names.Camel,
		)
		var matchWithPath string
		if inSpec {
			matchWithPath = matchWith + ".Spec"
		}

		switch sourceMemberShape.Type {
		case "list", "structure", "map":
		default:
			//          if ko.Spec.CacheClusterID != nil {
			//              if *elem.CacheClusterId != *ko.Spec.CacheClusterID {
			//                  continue
			//              }
			//          }

			if util.InStrings(fieldName, matchFieldNames) {
				out += fmt.Sprintf(
					"%s\tif %s != nil {\n", indent, sourceAdaptedVarName,
				)
				out += fmt.Sprintf(
					"%s\t\tif %s.%s != nil {\n",
					indent,
					matchWithPath,
					f.Names.Camel,
				)
				out += fmt.Sprintf(
					"%s\t\t\tif *%s != *%s.%s {\n",
					indent,
					sourceAdaptedVarName,
					matchWithPath,
					f.Names.Camel,
				)
				out += fmt.Sprintf("\t\t\t\t%srespMatch = elem\n", indent)
				out += fmt.Sprintf("\t\t\t\t%sbreak\n", indent)

				out += fmt.Sprintf(
					"%s\t\t\t}\n", indent,
				)
				out += fmt.Sprintf(
					"%s\t\t}\n", indent,
				)
				out += fmt.Sprintf(
					"%s\t}\n", indent,
				)
			}
		}

	}

	// }
	out += fmt.Sprintf("%s}\n", indent)

	return out
}
