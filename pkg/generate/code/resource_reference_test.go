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

func Test_ResolveReferences_NoReferenceConfig(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
		GeneratorConfigFile: "generator-with-reference.yaml",
	})

	crd := testutil.GetCRDByName(t, g, "Api")
	require.NotNil(crd)
	expected :=
		`	ko := rm.concreteResource(res).ko.DeepCopy()
	referencePresent := false
	if referencePresent {
		return ackcondition.WithReferencesResolvedCondition(&resource{ko}, nil)
	}
	return &resource{ko}, nil`
	assert.Equal(expected, code.ResolveReferences(crd, "ctx", "apiReader", "res", 1))
}

func Test_ResolveReferences_SingleReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Integration")
	require.NotNil(crd)
	expected :=
		`	ko := rm.concreteResource(res).ko.DeepCopy()
	referencePresent := false
	if ko.Spec.APIIDRef != nil && ko.Spec.APIID != nil {
		return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("'APIID' field should not be present when using reference field 'APIIDRef'"))
	}
	if ko.Spec.APIIDRef == nil && ko.Spec.APIID == nil {
		return &resource{ko}, fmt.Errorf("At least one of 'APIID' or 'APIIDRef' field should be present")
	}
	// Checking Referenced Field APIIDRef
	if ko.Spec.APIIDRef != nil && ko.Spec.APIIDRef.From != nil {
		referencePresent = true
		arr := ko.Spec.APIIDRef.From
		if arr == nil || arr.Name == nil || *arr.Name == "" {
			return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("provided resource reference is nil or empty"))
		}
		namespacedName := types.NamespacedName{Namespace: res.MetaObject().GetNamespace(), Name: *arr.Name}
		obj := acksvcv1alpha1.API{} 
		err := apiReader.Get(ctx, namespacedName, &obj)
		if err != nil {
			return ackcondition.WithReferencesResolvedCondition(&resource{ko}, err)
		}
		var refResourceSynced bool
		for _, cond := range obj.Status.Conditions {
			if cond.Type == ackv1alpha1.ConditionTypeResourceSynced && cond.Status == corev1.ConditionTrue {
				refResourceSynced = true
				break
			}
		}
		if !refResourceSynced {
			//TODO(vijtrip2) Uncomment below return statment once ConditionTypeResourceSynced(True/False) is set for all resources
			//return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("referenced 'API' resource " + *arr.Name + " does not have 'ACK.ResourceSynced' condition status 'True'"))
		}
		if obj.Status.APIID == nil {
			return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("'Status.APIID' is not yet present for referenced 'API' resource " + *arr.Name))
		}
		ko.Spec.APIID = obj.Status.APIID
	}
	if referencePresent {
		return ackcondition.WithReferencesResolvedCondition(&resource{ko}, nil)
	}
	return &resource{ko}, nil`
	assert.Equal(expected, code.ResolveReferences(crd, "ctx", "apiReader", "res", 1))
}

func Test_ResolveReferences_SliceOfReferences(t *testing.T) {
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
		`	ko := rm.concreteResource(res).ko.DeepCopy()
	referencePresent := false
	if ko.Spec.SecurityGroupIDsRef != nil && ko.Spec.SecurityGroupIDs != nil {
		return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("'SecurityGroupIDs' field should not be present when using reference field 'SecurityGroupIDsRef'"))
	}
	// Checking Referenced Field SecurityGroupIDsRef
	if ko.Spec.SecurityGroupIDsRef != nil && len(ko.Spec.SecurityGroupIDsRef) > 0 {
		referencePresent = true
		resolvedReferences := []*string{}
		for _, arrw := range ko.Spec.SecurityGroupIDsRef {
			arr := arrw.From
			if arr == nil || arr.Name == nil || *arr.Name == "" {
				return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("provided resource reference is nil or empty"))
			}
			namespacedName := types.NamespacedName{Namespace: res.MetaObject().GetNamespace(), Name: *arr.Name}
			obj := acksvcv1alpha1.API{} 
			err := apiReader.Get(ctx, namespacedName, &obj)
			if err != nil {
				return ackcondition.WithReferencesResolvedCondition(&resource{ko}, err)
			}
			var refResourceSynced bool
			for _, cond := range obj.Status.Conditions {
				if cond.Type == ackv1alpha1.ConditionTypeResourceSynced && cond.Status == corev1.ConditionTrue {
					refResourceSynced = true
					break
				}
			}
			if !refResourceSynced {
				//TODO(vijtrip2) Uncomment below return statment once ConditionTypeResourceSynced(True/False) is set for all resources
				//return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("referenced 'API' resource " + *arr.Name + " does not have 'ACK.ResourceSynced' condition status 'True'"))
			}
			if obj.Status.APIID == nil {
				return ackcondition.WithReferencesResolvedCondition(&resource{ko}, fmt.Errorf("'Status.APIID' is not yet present for referenced 'API' resource " + *arr.Name))
			}
			resolvedReferences = append(resolvedReferences, obj.Status.APIID)
		}
		ko.Spec.SecurityGroupIDs = resolvedReferences
	}
	if referencePresent {
		return ackcondition.WithReferencesResolvedCondition(&resource{ko}, nil)
	}
	return &resource{ko}, nil`
	assert.Equal(expected, code.ResolveReferences(crd, "ctx", "apiReader", "res", 1))
}
