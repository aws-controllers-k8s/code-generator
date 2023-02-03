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

package code_test

import (
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReferenceFieldsValidation_NoReferenceConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Api")
	require.NotNil(crd)
	expected := ""
	assert.Equal(expected, code.ReferenceFieldsValidation(crd, "ko", 1))
}

func Test_ReferenceFieldsValidation_SingleReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Integration")
	require.NotNil(crd)
	expected :=
		`	if ko.Spec.APIRef != nil && ko.Spec.APIID != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("APIID", "APIRef")
	}
	if ko.Spec.APIRef == nil && ko.Spec.APIID == nil {
		return ackerr.ResourceReferenceOrIDRequiredFor("APIID", "APIRef")
	}
`
	assert.Equal(expected, code.ReferenceFieldsValidation(crd, "ko", 1))
}

func Test_ReferenceFieldsValidation_SliceOfReferences(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	//NOTE: For the moment, we are substituting SecurityGroupId with ApiId
	// just to test code generation for slices of reference
	crd := testutil.GetCRDByName(t, g, "VpcLink")
	require.NotNil(crd)
	expected :=
		`	if ko.Spec.SecurityGroupRefs != nil && ko.Spec.SecurityGroupIDs != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("SecurityGroupIDs", "SecurityGroupRefs")
	}
	if ko.Spec.SubnetRefs != nil && ko.Spec.SubnetIDs != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("SubnetIDs", "SubnetRefs")
	}
	if ko.Spec.SubnetRefs == nil && ko.Spec.SubnetIDs == nil {
		return ackerr.ResourceReferenceOrIDRequiredFor("SubnetIDs", "SubnetRefs")
	}
`
	assert.Equal(expected, code.ReferenceFieldsValidation(crd, "ko", 1))
}

func Test_ReferenceFieldsValidation_NestedReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Authorizer")
	require.NotNil(crd)
	expected :=
		`	if ko.Spec.JWTConfiguration != nil {
		if ko.Spec.JWTConfiguration.IssuerRef != nil && ko.Spec.JWTConfiguration.Issuer != nil {
			return ackerr.ResourceReferenceAndIDNotSupportedFor("JWTConfiguration.Issuer", "JWTConfiguration.IssuerRef")
		}
	}
`
	assert.Equal(expected, code.ReferenceFieldsValidation(crd, "ko", 1))
}

func Test_ReferenceFieldsPresent_NoReferenceConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Api")
	require.NotNil(crd)
	expected := "return false"
	assert.Equal(expected, code.ReferenceFieldsPresent(crd, "ko"))
}

func Test_ReferenceFieldsPresent_SingleReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Integration")
	require.NotNil(crd)
	expected := "return false || (ko.Spec.APIRef != nil)"
	assert.Equal(expected, code.ReferenceFieldsPresent(crd, "ko"))
}

func Test_ReferenceFieldsPresent_SliceOfReferences(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	//NOTE: For the moment, we are substituting SecurityGroupId with ApiId
	// just to test code generation for slices of reference
	crd := testutil.GetCRDByName(t, g, "VpcLink")
	require.NotNil(crd)
	expected := "return false || (ko.Spec.SecurityGroupRefs != nil) || (ko.Spec.SubnetRefs != nil)"
	assert.Equal(expected, code.ReferenceFieldsPresent(crd, "ko"))
}

func Test_ReferenceFieldsPresent_NestedReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Authorizer")
	require.NotNil(crd)
	expected := "return false || (ko.Spec.JWTConfiguration != nil && ko.Spec.JWTConfiguration.IssuerRef != nil)"
	assert.Equal(expected, code.ReferenceFieldsPresent(crd, "ko"))
}

func Test_ReferenceFieldsPresent_NestedSliceOfStructsReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ec2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "RouteTable")
	require.NotNil(crd)
	expected :=
		`if ko.Spec.Routes != nil {
	for _, iter35 := range ko.Spec.Routes {
		if iter35.GatewayRef != nil {
			return true
		}
	}
}
if ko.Spec.Routes != nil {
	for _, iter38 := range ko.Spec.Routes {
		if iter38.NATGatewayRef != nil {
			return true
		}
	}
}
return false || (ko.Spec.VPCRef != nil)`
	assert.Equal(expected, code.ReferenceFieldsPresent(crd, "ko"))
}

func Test_ReferenceForField_SingleReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Integration")
	require.NotNil(crd)
	expected :=
		`	if ko.Spec.APIRef != nil && ko.Spec.APIRef.From != nil {
		arr := ko.Spec.APIRef.From
		if arr == nil || arr.Name == nil || *arr.Name == "" {
			return fmt.Errorf("provided resource reference is nil or empty")
		}
		namespacedName := types.NamespacedName{
			Namespace: namespace,
			Name: *arr.Name,
		}
		obj := svcapitypes.API{}
		err := apiReader.Get(ctx, namespacedName, &obj)
		if err != nil {
			return err
		}
		var refResourceSynced, refResourceTerminal bool
		for _, cond := range obj.Status.Conditions {
			if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
				cond.Status == corev1.ConditionTrue {
				refResourceSynced = true
			}
			if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
				cond.Status == corev1.ConditionTrue {
				refResourceTerminal = true
			}
		}
		if refResourceTerminal {
			return ackerr.ResourceReferenceTerminalFor(
				"API",
				namespace, *arr.Name)
		}
		if !refResourceSynced {
			return ackerr.ResourceReferenceNotSyncedFor(
				"API",
				namespace, *arr.Name)
		}
		if obj.Status.APIID == nil {
			return ackerr.ResourceReferenceMissingTargetFieldFor(
				"API",
				namespace, *arr.Name,
				"Status.APIID")
		}
		referencedValue := string(*obj.Status.APIID)
		ko.Spec.APIID = &referencedValue
	}
	return nil`

	field := crd.Fields["APIID"]
	assert.Equal(expected, code.ResolveReferencesForField(crd, field, 1))
}

func Test_ReferenceForField_SliceOfReferences(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "VpcLink")
	require.NotNil(crd)
	expected :=
		`	if ko.Spec.SecurityGroupRefs != nil &&
		len(ko.Spec.SecurityGroupRefs) > 0 {
		resolvedReferences := []*string{}
		for _, arrw := range ko.Spec.SecurityGroupRefs {
			arr := arrw.From
			if arr == nil || arr.Name == nil || *arr.Name == "" {
				return fmt.Errorf("provided resource reference is nil or empty")
			}
			namespacedName := types.NamespacedName{
				Namespace: namespace,
				Name: *arr.Name,
			}
			obj := ec2apitypes.SecurityGroup{}
			err := apiReader.Get(ctx, namespacedName, &obj)
			if err != nil {
				return err
			}
			var refResourceSynced, refResourceTerminal bool
			for _, cond := range obj.Status.Conditions {
				if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
					cond.Status == corev1.ConditionTrue {
					refResourceSynced = true
				}
				if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
					cond.Status == corev1.ConditionTrue {
					refResourceTerminal = true
				}
			}
			if refResourceTerminal {
				return ackerr.ResourceReferenceTerminalFor(
					"SecurityGroup",
					namespace, *arr.Name)
			}
			if !refResourceSynced {
				return ackerr.ResourceReferenceNotSyncedFor(
					"SecurityGroup",
					namespace, *arr.Name)
			}
			if obj.Status.ID == nil {
				return ackerr.ResourceReferenceMissingTargetFieldFor(
					"SecurityGroup",
					namespace, *arr.Name,
					"Status.ID")
			}
			referencedValue := string(*obj.Status.ID)
			resolvedReferences = append(resolvedReferences, &referencedValue)
		}
		ko.Spec.SecurityGroupIDs = resolvedReferences
	}
	return nil`

	field := crd.Fields["SecurityGroupIDs"]
	assert.Equal(expected, code.ResolveReferencesForField(crd, field, 1))
}

func Test_ReferenceForField_NestedSingleReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Authorizer")
	require.NotNil(crd)
	expected :=
		`	if ko.Spec.JWTConfiguration.IssuerRef != nil &&
		len(ko.Spec.JWTConfiguration.IssuerRef) > 0 {
		resolvedReferences := []*string{}
		for _, arrw := range ko.Spec.JWTConfiguration.IssuerRef {
			arr := arrw.From
			if arr == nil || arr.Name == nil || *arr.Name == "" {
				return fmt.Errorf("provided resource reference is nil or empty")
			}
			namespacedName := types.NamespacedName{
				Namespace: namespace,
				Name: *arr.Name,
			}
			obj := svcapitypes.API{}
			err := apiReader.Get(ctx, namespacedName, &obj)
			if err != nil {
				return err
			}
			var refResourceSynced, refResourceTerminal bool
			for _, cond := range obj.Status.Conditions {
				if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
					cond.Status == corev1.ConditionTrue {
					refResourceSynced = true
				}
				if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
					cond.Status == corev1.ConditionTrue {
					refResourceTerminal = true
				}
			}
			if refResourceTerminal {
				return ackerr.ResourceReferenceTerminalFor(
					"API",
					namespace, *arr.Name)
			}
			if !refResourceSynced {
				return ackerr.ResourceReferenceNotSyncedFor(
					"API",
					namespace, *arr.Name)
			}
			if obj.Status.APIID == nil {
				return ackerr.ResourceReferenceMissingTargetFieldFor(
					"API",
					namespace, *arr.Name,
					"Status.APIID")
			}
			referencedValue := string(*obj.Status.APIID)
			resolvedReferences = append(resolvedReferences, &referencedValue)
		}
		ko.Spec.JWTConfiguration.Issuer = resolvedReferences
	}
	return nil`

	field := crd.Fields["JWTConfiguration.Issuer"]
	assert.Equal(expected, code.ResolveReferencesForField(crd, field, 1))
}

func Test_ReferenceForField_SingleReference_DeeplyNested(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "s3",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-references.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Bucket")
	require.NotNil(crd)

	// the Go template has the appropriate nil checks to ensure the parent path exists
	expected :=
		`	if ko.Spec.Analytics.StorageClassAnalysis.DataExport.Destination.S3BucketDestination.BucketRef != nil &&
		len(ko.Spec.Analytics.StorageClassAnalysis.DataExport.Destination.S3BucketDestination.BucketRef) > 0 {
		resolvedReferences := []*string{}
		for _, arrw := range ko.Spec.Analytics.StorageClassAnalysis.DataExport.Destination.S3BucketDestination.BucketRef {
			arr := arrw.From
			if arr == nil || arr.Name == nil || *arr.Name == "" {
				return fmt.Errorf("provided resource reference is nil or empty")
			}
			namespacedName := types.NamespacedName{
				Namespace: namespace,
				Name: *arr.Name,
			}
			obj := svcapitypes.Bucket{}
			err := apiReader.Get(ctx, namespacedName, &obj)
			if err != nil {
				return err
			}
			var refResourceSynced, refResourceTerminal bool
			for _, cond := range obj.Status.Conditions {
				if cond.Type == ackv1alpha1.ConditionTypeResourceSynced &&
					cond.Status == corev1.ConditionTrue {
					refResourceSynced = true
				}
				if cond.Type == ackv1alpha1.ConditionTypeTerminal &&
					cond.Status == corev1.ConditionTrue {
					refResourceTerminal = true
				}
			}
			if refResourceTerminal {
				return ackerr.ResourceReferenceTerminalFor(
					"Bucket",
					namespace, *arr.Name)
			}
			if !refResourceSynced {
				return ackerr.ResourceReferenceNotSyncedFor(
					"Bucket",
					namespace, *arr.Name)
			}
			if obj.Spec.Name == nil {
				return ackerr.ResourceReferenceMissingTargetFieldFor(
					"Bucket",
					namespace, *arr.Name,
					"Spec.Name")
			}
			referencedValue := string(*obj.Spec.Name)
			resolvedReferences = append(resolvedReferences, &referencedValue)
		}
		ko.Spec.Analytics.StorageClassAnalysis.DataExport.Destination.S3BucketDestination.Bucket = resolvedReferences
	}
	return nil`

	field := crd.Fields["Analytics.StorageClassAnalysis.DataExport.Destination.S3BucketDestination.Bucket"]
	assert.Equal(expected, code.ResolveReferencesForField(crd, field, 1))
}
