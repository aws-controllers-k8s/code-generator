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

// ResolveReferences produces the go code to read all the referenced resource(s)
// inside a resource and populate the target field(s) from referenced resource(s)
//
// Sample code:
//	ko := rm.concreteResource(res).ko.DeepCopy()
//	referencePresent := false
//	if ko.Spec.APIIDRef != nil && ko.Spec.APIID != nil {
//		return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("'APIID' field should not be present when using reference field 'APIIDRef'"))
//	}
//	if ko.Spec.APIIDRef == nil && ko.Spec.APIID == nil {
//		return &resource{ko}, fmt.Errorf("At least one of 'APIID' or 'APIIDRef' field should be present")
//	}
//	// Checking Referenced Field APIIDRef
//	if ko.Spec.APIIDRef != nil && ko.Spec.APIIDRef.From != nil {
//		referencePresent = true
//		arr := ko.Spec.APIIDRef.From
//		if arr == nil || arr.Name == nil || *arr.Name == "" {
//			return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("provided resource reference is nil or empty"))
//		}
//		namespacedName := types.NamespacedName{Namespace: res.MetaObject().GetNamespace(), Name: *arr.Name}
//		obj := acksvcv1alpha1.API{}
//		err := apiReader.Get(ctx, namespacedName, &obj)
//		if err != nil {
//			return ackcondition.WithReferencesResolvedCondition(&resource{ko}, err)
//		}
//		var refResourceSynced bool
//		for _, cond := range obj.Status.Conditions {
//			if cond.Type == ackv1alpha1.ConditionTypeResourceSynced && cond.Status == corev1.ConditionTrue {
//				refResourceSynced = true
//				break
//			}
//		}
//		if !refResourceSynced {
//			//TODO(vijtrip2) Uncomment below return statment once ConditionTypeResourceSynced(True/False) is set for all resources
//			//return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("referenced 'API' resource " + *arr.Name + " does not have 'ACK.ResourceSynced' condition status 'True'"))
//		}
//		if obj.Status.APIID == nil {
//			return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("'Status.APIID' is not yet present for referenced 'API' resource " + *arr.Name))
//		}
//		ko.Spec.APIID = obj.Status.APIID
//	}
//	if referencePresent {
//		return ackcondition.WithReferencesResolvedCondition(&resource{ko}, nil)
//	}
//	return &resource{ko}, nil
func ResolveReferences(
	r *model.CRD,
	contextVarName string,
	apiReaderVarName string,
	resVarName string,
// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)
	out += fmt.Sprintf("%sko := rm.concreteResource(%s).ko.DeepCopy()\n",
		indent, resVarName)
	out += fmt.Sprintf("%sreferencePresent := false\n", indent)
	// Iterate over all fields
	if r.Fields != nil {
		for _, field := range r.Fields {
			if field.HasReference() {
				// Validation to make sure both target field and reference are
				// not present at the same time in desired resource
				out += fmt.Sprintf("%sif ko.Spec.%sRef != nil"+
					" && ko.Spec.%s != nil {\n", indent, field.Names.Camel,
					field.Names.Camel)
				out += fmt.Sprintf("%s\treturn ackcondition."+
					"WithReferencesResolvedCondition(&resource{ko}, "+
					"fmt.Errorf(\"'%s' field should not be present when using"+
					" reference field '%sRef'\"))\n", indent,field.Names.Camel,
					field.Names.Camel)
				out += fmt.Sprintf("%s}\n", indent)

				// If the field is required, make sure either Ref or original
				// field is present in the resource
				if field.IsRequired() {
					out += fmt.Sprintf("%sif ko.Spec.%sRef == nil &&"+
						" ko.Spec.%s == nil {\n", indent, field.Names.Camel,
						field.Names.Camel)
					fmt.Sprintf("%s// No need to add reference resolved"+
						" condition because reference is not present\t", indent)
					out += fmt.Sprintf("%s\treturn &resource{ko},"+
						" fmt.Errorf(\"At least one of '%s' or '%sRef' field"+
						" should be present\")\n", indent, field.Names.Camel,
						field.Names.Camel)
					out += fmt.Sprintf("%s}\n", indent)
				}

				fIndentLevel := indentLevel
				fIndent := indent
				out += fmt.Sprintf("%s// Checking Referenced Field"+
					" %sRef\n", fIndent, field.Names.Camel)
				if field.ShapeRef.Shape.Type == "list" {
					out += fmt.Sprintf("%sif ko.Spec.%sRef != nil &&"+
						" len(ko.Spec.%sRef) > 0 {\n", fIndent,
						field.Names.Camel, field.Names.Camel)
					out += fmt.Sprintf("%s\treferencePresent = true\n",
						fIndent)
					out += fmt.Sprintf("%s\tresolvedReferences :="+
						" []*string{}\n", fIndent)
					// arrw stands for AWSResourceReferenceWrapper
					out += fmt.Sprintf("%s\tfor _, arrw := range"+
						" ko.Spec.%sRef {\n", fIndent, field.Names.Camel)
					out += fmt.Sprintf("%s\t\tarr := arrw.From\n",
						fIndent)
					out += readReferencedResource(field, strings.Repeat("\t",
						fIndentLevel+2), apiReaderVarName, contextVarName,
						resVarName)
					out += checkReferencedResource(field, strings.Repeat("\t",
						fIndentLevel+2))
					out += fmt.Sprintf("%s\t\tresolvedReferences ="+
						" append(resolvedReferences, obj.%s)\n", fIndent,
						field.FieldConfig.References.Path)
					out += fmt.Sprintf("%s\t}\n", fIndent)
					out += fmt.Sprintf("%s\tko.Spec.%s ="+
						" resolvedReferences\n", fIndent, field.Names.Camel)
					out += fmt.Sprintf("%s}\n", fIndent)
				} else {
					out += fmt.Sprintf("%sif ko.Spec.%sRef != nil &&"+
						" ko.Spec.%sRef.From != nil {\n", fIndent,
						field.Names.Camel, field.Names.Camel)
					out += fmt.Sprintf("%s\treferencePresent = true\n",
						fIndent)
					out += fmt.Sprintf("%s\tarr := ko.Spec.%sRef.From\n",
						fIndent, field.Names.Camel)
					out += readReferencedResource(field, strings.Repeat("\t",
						fIndentLevel+1), apiReaderVarName, contextVarName,
						resVarName)
					out += checkReferencedResource(field, strings.Repeat("\t",
						fIndentLevel+1))
					out += fmt.Sprintf("%s\tko.Spec.%s = obj.%s\n",
						fIndent, field.Names.Camel,
						field.FieldConfig.References.Path)
					out += fmt.Sprintf("%s}\n", fIndent)
				}
			}
		}
	}
	out += fmt.Sprintf("%sif referencePresent {\n", indent)
	out += fmt.Sprintf("%s\treturn ackcondition."+
		"WithReferencesResolvedCondition(&resource{ko}, nil)\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%sreturn &resource{ko}, nil", indent)
	return out
}

//readReferencedResource generates the go code for reading a referenced resource
func readReferencedResource(
	field *model.Field,
	indent string,
	apiReaderVarName string,
	contextVarName string,
	resVarName string,
) string {
	out := ""
	out += fmt.Sprintf("%sif arr == nil || arr.Name == nil ||"+
		" *arr.Name == \"\" {\n", indent)
	out += fmt.Sprintf("%s\treturn ackcondition."+
		"WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf(\"provided"+
		" resource reference is nil or empty\"))\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%snamespacedName := types.NamespacedName{"+
		"Namespace: %s.MetaObject().GetNamespace(), Name: *arr.Name}\n", indent,
		resVarName)
	out += fmt.Sprintf("%sobj := acksvcv1alpha1.%s{} \n", indent,
		field.FieldConfig.References.Resource)
	out += fmt.Sprintf("%serr := %s.Get(%s, namespacedName, &obj)\n",
		indent, apiReaderVarName, contextVarName)
	out += fmt.Sprintf("%sif err != nil {\n", indent)
	out += fmt.Sprintf("%s\treturn ackcondition."+
		"WithReferencesResolvedCondition(&resource{ko}, err)\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	return out
}

// checkReferencedResource generates the go code for validating if the
// referenced resource is ready for copying the target field from referenced
// resource
func checkReferencedResource(field *model.Field, indent string) string {
	out := ""
	out += fmt.Sprintf("%svar refResourceSynced bool\n", indent)
	out += fmt.Sprintf("%sfor _, cond := range obj.Status.Conditions {\n",
		indent)
	out += fmt.Sprintf("%s\tif cond.Type == ackv1alpha1."+
		"ConditionTypeResourceSynced && cond.Status == corev1.ConditionTrue {\n",
		indent)
	out += fmt.Sprintf("%s\t\trefResourceSynced = true\n", indent)
	out += fmt.Sprintf("%s\t\tbreak\n", indent)
	out += fmt.Sprintf("%s\t}\n", indent)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%sif !refResourceSynced {\n", indent)
	out += fmt.Sprintf("%s\t//TODO(vijtrip2) Uncomment below return "+
		"statment once ConditionTypeResourceSynced(True/False) is set for all"+
		" resources\n", indent)
	out += fmt.Sprintf("%s\t//return ackcondition."+
		"WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf(\"referenced"+
		" '%s' resource \" + *arr.Name + \" does not have 'ACK.ResourceSynced'"+
		" condition status 'True'\"))\n", indent,
		field.FieldConfig.References.Resource)
	out += fmt.Sprintf("%s}\n", indent)
	out += fmt.Sprintf("%sif obj.%s == nil {\n", indent,
		field.FieldConfig.References.Path)
	out += fmt.Sprintf("%s\treturn ackcondition."+
		"WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf(\"'%s' is"+
		" not yet present for referenced '%s' resource \" + *arr.Name))\n",
		indent, field.FieldConfig.References.Path,
		field.FieldConfig.References.Resource)
	out += fmt.Sprintf("%s}\n", indent)
	return out
}
