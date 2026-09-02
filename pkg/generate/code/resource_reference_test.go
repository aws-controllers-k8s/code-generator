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
	"strings"
	"testing"

	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ReferenceFieldsValidation_SingleReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Integration")
	require.NotNil(crd)

	field := crd.Fields["APIID"]
	expected :=
		`	if ko.Spec.APIRef != nil && ko.Spec.APIID != nil {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("APIID", "APIRef")
	}
	if ko.Spec.APIRef == nil && ko.Spec.APIID == nil {
		return ackerr.ResourceReferenceOrIDRequiredFor("APIID", "APIRef")
	}
`
	got, err := code.ReferenceFieldsValidation(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ReferenceFieldsValidation_WithOptional_SliceOfReferences(t *testing.T) {
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

	field := crd.Fields["SecurityGroupIDs"]
	expected :=
		`	if len(ko.Spec.SecurityGroupRefs) > 0 && len(ko.Spec.SecurityGroupIDs) > 0 {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("SecurityGroupIDs", "SecurityGroupRefs")
	}
`
	got, err := code.ReferenceFieldsValidation(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ReferenceFieldsValidation_WithRequired_SliceOfReferences(t *testing.T) {
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

	field := crd.Fields["SubnetIDs"]
	expected :=
		`	if len(ko.Spec.SubnetRefs) > 0 && len(ko.Spec.SubnetIDs) > 0 {
		return ackerr.ResourceReferenceAndIDNotSupportedFor("SubnetIDs", "SubnetRefs")
	}
	if len(ko.Spec.SubnetRefs) == 0 && len(ko.Spec.SubnetIDs) == 0 {
		return ackerr.ResourceReferenceOrIDRequiredFor("SubnetIDs", "SubnetRefs")
	}
`
	got, err := code.ReferenceFieldsValidation(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
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

	field := crd.Fields["JWTConfiguration.Issuer"]
	expected :=
		`	if ko.Spec.JWTConfiguration != nil {
		if ko.Spec.JWTConfiguration.IssuerRef != nil && ko.Spec.JWTConfiguration.Issuer != nil {
			return ackerr.ResourceReferenceAndIDNotSupportedFor("JWTConfiguration.Issuer", "JWTConfiguration.IssuerRef")
		}
	}
`
	got, err := code.ReferenceFieldsValidation(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ResolveReferencesForField_SingleReference(t *testing.T) {
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
		hasReferences = true
		arr := ko.Spec.APIRef.From
		if arr.Name == nil || *arr.Name == "" {
			return hasReferences, fmt.Errorf("provided resource reference is nil or empty: APIRef")
		}
		namespace, err := ackrt.ResolveCrossNamespaceReference(
			ctx,
			rm.cfg.EnableCrossNamespace,
			&ko.Status.Conditions,
			ackrt.CrossNamespaceRefKindResource,
			ko.ObjectMeta.GetNamespace(),
			arr.Namespace,
			*arr.Name,
		)
		if err != nil {
			return hasReferences, err
		}
		obj := &svcapitypes.API{}
		if err := getReferencedResourceState_API(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
			return hasReferences, err
		}
		ko.Spec.APIID = (*string)(obj.Status.APIID)
	}
`

	field := crd.Fields["APIID"]
	got, err := code.ResolveReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ResolveReferencesForField_ReferencingARN(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "iam",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "User")
	require.NotNil(crd)
	expected :=
		`	if ko.Spec.PermissionsBoundaryRef != nil && ko.Spec.PermissionsBoundaryRef.From != nil {
		hasReferences = true
		arr := ko.Spec.PermissionsBoundaryRef.From
		if arr.Name == nil || *arr.Name == "" {
			return hasReferences, fmt.Errorf("provided resource reference is nil or empty: PermissionsBoundaryRef")
		}
		namespace, err := ackrt.ResolveCrossNamespaceReference(
			ctx,
			rm.cfg.EnableCrossNamespace,
			&ko.Status.Conditions,
			ackrt.CrossNamespaceRefKindResource,
			ko.ObjectMeta.GetNamespace(),
			arr.Namespace,
			*arr.Name,
		)
		if err != nil {
			return hasReferences, err
		}
		obj := &svcapitypes.Policy{}
		if err := getReferencedResourceState_Policy(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
			return hasReferences, err
		}
		ko.Spec.PermissionsBoundary = (*string)(obj.Status.ACKResourceMetadata.ARN)
	}
`

	field := crd.Fields["PermissionsBoundary"]
	got, err := code.ResolveReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ResolveReferencesForField_SliceOfReferences(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "VpcLink")
	require.NotNil(crd)
	expected :=
		`	for _, f0iter := range ko.Spec.SecurityGroupRefs {
		if f0iter != nil && f0iter.From != nil {
			hasReferences = true
			arr := f0iter.From
			if arr.Name == nil || *arr.Name == "" {
				return hasReferences, fmt.Errorf("provided resource reference is nil or empty: SecurityGroupRefs")
			}
			namespace, err := ackrt.ResolveCrossNamespaceReference(
				ctx,
				rm.cfg.EnableCrossNamespace,
				&ko.Status.Conditions,
				ackrt.CrossNamespaceRefKindResource,
				ko.ObjectMeta.GetNamespace(),
				arr.Namespace,
				*arr.Name,
			)
			if err != nil {
				return hasReferences, err
			}
			obj := &ec2apitypes.SecurityGroup{}
			if err := getReferencedResourceState_SecurityGroup(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
				return hasReferences, err
			}
			if ko.Spec.SecurityGroupIDs == nil {
				ko.Spec.SecurityGroupIDs = make([]*string, 0, 1)
			}
			ko.Spec.SecurityGroupIDs = append(ko.Spec.SecurityGroupIDs, (*string)(obj.Status.ID))
		}
	}
`

	field := crd.Fields["SecurityGroupIDs"]
	got, err := code.ResolveReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ResolveReferencesForField_NestedSingleReference(t *testing.T) {
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
		if ko.Spec.JWTConfiguration.IssuerRef != nil && ko.Spec.JWTConfiguration.IssuerRef.From != nil {
			hasReferences = true
			arr := ko.Spec.JWTConfiguration.IssuerRef.From
			if arr.Name == nil || *arr.Name == "" {
				return hasReferences, fmt.Errorf("provided resource reference is nil or empty: JWTConfiguration.IssuerRef")
			}
			namespace, err := ackrt.ResolveCrossNamespaceReference(
				ctx,
				rm.cfg.EnableCrossNamespace,
				&ko.Status.Conditions,
				ackrt.CrossNamespaceRefKindResource,
				ko.ObjectMeta.GetNamespace(),
				arr.Namespace,
				*arr.Name,
			)
			if err != nil {
				return hasReferences, err
			}
			obj := &svcapitypes.API{}
			if err := getReferencedResourceState_API(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
				return hasReferences, err
			}
			ko.Spec.JWTConfiguration.Issuer = (*string)(obj.Status.APIID)
		}
	}
`

	field := crd.Fields["JWTConfiguration.Issuer"]
	got, err := code.ResolveReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ResolveReferencesForField_SingleReference_DeeplyNested(t *testing.T) {
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
		`	if ko.Spec.Logging != nil {
		if ko.Spec.Logging.LoggingEnabled != nil {
			if ko.Spec.Logging.LoggingEnabled.TargetBucketRef != nil && ko.Spec.Logging.LoggingEnabled.TargetBucketRef.From != nil {
				hasReferences = true
				arr := ko.Spec.Logging.LoggingEnabled.TargetBucketRef.From
				if arr.Name == nil || *arr.Name == "" {
					return hasReferences, fmt.Errorf("provided resource reference is nil or empty: Logging.LoggingEnabled.TargetBucketRef")
				}
				namespace, err := ackrt.ResolveCrossNamespaceReference(
					ctx,
					rm.cfg.EnableCrossNamespace,
					&ko.Status.Conditions,
					ackrt.CrossNamespaceRefKindResource,
					ko.ObjectMeta.GetNamespace(),
					arr.Namespace,
					*arr.Name,
				)
				if err != nil {
					return hasReferences, err
				}
				obj := &svcapitypes.Bucket{}
				if err := getReferencedResourceState_Bucket(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
					return hasReferences, err
				}
				ko.Spec.Logging.LoggingEnabled.TargetBucket = (*string)(obj.Spec.Name)
			}
		}
	}
`

	field := crd.Fields["Logging.LoggingEnabled.TargetBucket"]
	got, err := code.ResolveReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ResolveReferencesForField_SingleReference_WithinSlice(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ec2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-references.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "RouteTable")
	require.NotNil(crd)

	// the Go template has the appropriate nil checks to ensure the parent path exists
	expected :=
		`	for f0idx, f0iter := range ko.Spec.Routes {
		if f0iter.GatewayRef != nil && f0iter.GatewayRef.From != nil {
			hasReferences = true
			arr := f0iter.GatewayRef.From
			if arr.Name == nil || *arr.Name == "" {
				return hasReferences, fmt.Errorf("provided resource reference is nil or empty: Routes.GatewayRef")
			}
			namespace, err := ackrt.ResolveCrossNamespaceReference(
				ctx,
				rm.cfg.EnableCrossNamespace,
				&ko.Status.Conditions,
				ackrt.CrossNamespaceRefKindResource,
				ko.ObjectMeta.GetNamespace(),
				arr.Namespace,
				*arr.Name,
			)
			if err != nil {
				return hasReferences, err
			}
			obj := &svcapitypes.InternetGateway{}
			if err := getReferencedResourceState_InternetGateway(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
				return hasReferences, err
			}
			ko.Spec.Routes[f0idx].GatewayID = (*string)(obj.Status.InternetGatewayID)
		}
	}
`

	field := crd.Fields["Routes.GatewayID"]
	got, err := code.ResolveReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ResolveReferencesForField_SingleReference_WithinMultipleSlices(t *testing.T) {
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
		`	if ko.Spec.Notification != nil {
		for f0idx, f0iter := range ko.Spec.Notification.LambdaFunctionConfigurations {
			if f0iter.Filter != nil {
				if f0iter.Filter.Key != nil {
					for f1idx, f1iter := range f0iter.Filter.Key.FilterRules {
						if f1iter.ValueRef != nil && f1iter.ValueRef.From != nil {
							hasReferences = true
							arr := f1iter.ValueRef.From
							if arr.Name == nil || *arr.Name == "" {
								return hasReferences, fmt.Errorf("provided resource reference is nil or empty: Notification.LambdaFunctionConfigurations.Filter.Key.FilterRules.ValueRef")
							}
							namespace, err := ackrt.ResolveCrossNamespaceReference(
								ctx,
								rm.cfg.EnableCrossNamespace,
								&ko.Status.Conditions,
								ackrt.CrossNamespaceRefKindResource,
								ko.ObjectMeta.GetNamespace(),
								arr.Namespace,
								*arr.Name,
							)
							if err != nil {
								return hasReferences, err
							}
							obj := &svcapitypes.Bucket{}
							if err := getReferencedResourceState_Bucket(ctx, apiReader, obj, *arr.Name, namespace); err != nil {
								return hasReferences, err
							}
							ko.Spec.Notification.LambdaFunctionConfigurations[f0idx].Filter.Key.FilterRules[f1idx].Value = (*string)(obj.Spec.Name)
						}
					}
				}
			}
		}
	}
`

	field := crd.Fields["Notification.LambdaFunctionConfigurations.Filter.Key.FilterRules.Value"]
	got, err := code.ResolveReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ClearResolvedReferencesForField_SingleReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Integration")
	require.NotNil(crd)
	expected :=
		`	if ko.Spec.APIRef != nil {
		ko.Spec.APIID = nil
	}
`

	field := crd.Fields["APIID"]
	got, err := code.ClearResolvedReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ClearResolvedReferencesForField_SliceOfReferences(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "VpcLink")
	require.NotNil(crd)
	expected :=
		`	if len(ko.Spec.SecurityGroupRefs) > 0 {
		ko.Spec.SecurityGroupIDs = nil
	}
`

	field := crd.Fields["SecurityGroupIDs"]
	got, err := code.ClearResolvedReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ClearResolvedReferencesForField_NestedSingleReference(t *testing.T) {
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
		if ko.Spec.JWTConfiguration.IssuerRef != nil {
			ko.Spec.JWTConfiguration.Issuer = nil
		}
	}
`

	field := crd.Fields["JWTConfiguration.Issuer"]
	got, err := code.ClearResolvedReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ClearResolvedReferencesForField_SingleReference_DeeplyNested(t *testing.T) {
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
		`	if ko.Spec.Logging != nil {
		if ko.Spec.Logging.LoggingEnabled != nil {
			if ko.Spec.Logging.LoggingEnabled.TargetBucketRef != nil {
				ko.Spec.Logging.LoggingEnabled.TargetBucket = nil
			}
		}
	}
`

	field := crd.Fields["Logging.LoggingEnabled.TargetBucket"]
	got, err := code.ClearResolvedReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ClearResolvedReferencesForField_SingleReference_WithinSlice(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ec2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-references.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "RouteTable")
	require.NotNil(crd)

	// the Go template has the appropriate nil checks to ensure the parent path exists
	expected :=
		`	for f0idx, f0iter := range ko.Spec.Routes {
		if f0iter.GatewayRef != nil {
			ko.Spec.Routes[f0idx].GatewayID = nil
		}
	}
`

	field := crd.Fields["Routes.GatewayID"]
	got, err := code.ClearResolvedReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_ClearResolvedReferencesForField_SingleReference_WithinMultipleSlices(t *testing.T) {
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
		`	if ko.Spec.Notification != nil {
		for f0idx, f0iter := range ko.Spec.Notification.LambdaFunctionConfigurations {
			if f0iter.Filter != nil {
				if f0iter.Filter.Key != nil {
					for f1idx, f1iter := range f0iter.Filter.Key.FilterRules {
						if f1iter.ValueRef != nil {
							ko.Spec.Notification.LambdaFunctionConfigurations[f0idx].Filter.Key.FilterRules[f1idx].Value = nil
						}
					}
				}
			}
		}
	}
`

	field := crd.Fields["Notification.LambdaFunctionConfigurations.Filter.Key.FilterRules.Value"]
	got, err := code.ClearResolvedReferencesForField(field, "ko", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_EnsureReferences_TopLevelReference_EmitsNothing(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference.yaml",
		})

	// Integration's only reference is the top-level APIID/APIRef pair; VpcLink's
	// are top-level lists of references. A top-level *Ref has no parent that could
	// be rebuilt, so it always survives and nothing needs emitting.
	for _, kind := range []string{"Integration", "VpcLink"} {
		crd := testutil.GetCRDByName(t, g, kind)
		require.NotNil(crd)

		got, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 1)
		require.NoError(err)
		assert.Equal("", got, "resource %s", kind)
	}
}

func Test_EnsureReferences_StructPath_AssignsOnlyTheReference(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-reference.yaml",
		})

	// Reached through a struct, so the reference has one fixed address: guard the
	// ancestors on both objects and assign just that field.
	crd := testutil.GetCRDByName(t, g, "Authorizer")
	require.NotNil(crd)
	expected :=
		`	if desiredKO.Spec.JWTConfiguration != nil {
		if desiredKO.Spec.JWTConfiguration.IssuerRef != nil {
			if latestKO.Spec.JWTConfiguration == nil {
				latestKO.Spec.JWTConfiguration = &svcapitypes.JWTConfiguration{}
			}
			if latestKO.Spec.JWTConfiguration.IssuerRef == nil {
				latestKO.Spec.JWTConfiguration.IssuerRef = desiredKO.Spec.JWTConfiguration.IssuerRef
			}
		}
	}
`

	got, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 1)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_EnsureReferences_StructPath_ListOfReferences(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "eks",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-reference.yaml",
		})

	// The reference field is itself a list (*Refs) but sits in a struct at a fixed
	// address, so the list is the leaf rather than part of the path. It is copied
	// whole, guarded on length, as ClearResolvedReferences treats the same shape.
	// This is the shape community#2431 was filed for.
	crd := testutil.GetCRDByName(t, g, "Cluster")
	require.NotNil(crd)
	expected :=
		`	if desiredKO.Spec.ResourcesVPCConfig != nil {
		if len(desiredKO.Spec.ResourcesVPCConfig.SecurityGroupRefs) > 0 {
			if latestKO.Spec.ResourcesVPCConfig == nil {
				latestKO.Spec.ResourcesVPCConfig = &svcapitypes.VPCConfigRequest{}
			}
			if len(latestKO.Spec.ResourcesVPCConfig.SecurityGroupRefs) == 0 {
				latestKO.Spec.ResourcesVPCConfig.SecurityGroupRefs = desiredKO.Spec.ResourcesVPCConfig.SecurityGroupRefs
			}
		}
	}
`

	got, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 1)
	require.NoError(err)
	assert.Equal(expected, got)
	// The concrete sibling is never read or written.
	assert.NotContains(got, "SecurityGroupIDs")
}

func Test_EnsureReferences_ListPath_IsSkipped(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "ec2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-references.yaml",
		})

	// RouteTable's references are two inside spec.Routes plus a top-level VPCID.
	// The top-level one needs no help and the list-nested ones are skipped, so
	// nothing is emitted and the template's `if $ensureReferences` guard leaves
	// RouteTable without the method entirely.
	crd := testutil.GetCRDByName(t, g, "RouteTable")
	require.NotNil(crd)

	got, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 1)
	require.NoError(err)
	assert.Equal("", got)
}

func Test_EnsureReferences_MixedShapes_EmitsOnlyTheStructPath(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "s3",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-references.yaml",
		})

	// Bucket carries one of each shape, pinning that they are treated differently
	// within a single resource: Logging.LoggingEnabled.TargetBucket is reached
	// through structs alone, while
	// Notification.LambdaFunctionConfigurations[].Filter.Key.FilterRules[].Value
	// sits two lists deep.
	crd := testutil.GetCRDByName(t, g, "Bucket")
	require.NotNil(crd)

	got, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 1)
	require.NoError(err)

	// The struct path is restored, writing nothing but the reference itself.
	expected :=
		`	if desiredKO.Spec.Logging != nil {
		if desiredKO.Spec.Logging.LoggingEnabled != nil {
			if desiredKO.Spec.Logging.LoggingEnabled.TargetBucketRef != nil {
				if latestKO.Spec.Logging == nil {
					latestKO.Spec.Logging = &svcapitypes.BucketLoggingStatus{}
				}
				if latestKO.Spec.Logging.LoggingEnabled == nil {
					latestKO.Spec.Logging.LoggingEnabled = &svcapitypes.LoggingEnabled{}
				}
				if latestKO.Spec.Logging.LoggingEnabled.TargetBucketRef == nil {
					latestKO.Spec.Logging.LoggingEnabled.TargetBucketRef = desiredKO.Spec.Logging.LoggingEnabled.TargetBucketRef
				}
			}
		}
	}
`
	assert.Equal(expected, got)

	// The list-nested reference contributes nothing, and in particular the
	// containing list is not assigned.
	assert.NotContains(got, "LambdaFunctionConfigurations")
	assert.NotContains(got, "FilterRules")

	// Nothing is iterated or indexed on the way there.
	assert.NotContains(got, "range")
	assert.NotContains(got, "[f0idx]")
	assert.NotContains(got, "[f1idx]")
}

func Test_EnsureReferences_RespectsIndentLevel(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-reference.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Authorizer")
	require.NotNil(crd)
	expected :=
		`			if desiredKO.Spec.JWTConfiguration != nil {
				if desiredKO.Spec.JWTConfiguration.IssuerRef != nil {
					if latestKO.Spec.JWTConfiguration == nil {
						latestKO.Spec.JWTConfiguration = &svcapitypes.JWTConfiguration{}
					}
					if latestKO.Spec.JWTConfiguration.IssuerRef == nil {
						latestKO.Spec.JWTConfiguration.IssuerRef = desiredKO.Spec.JWTConfiguration.IssuerRef
					}
				}
			}
`

	got, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 3)
	require.NoError(err)
	assert.Equal(expected, got)
}

func Test_EnsureReferences_ReferenceWithinMap_IsRejected(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "apigatewayv2",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-reference-in-map.yaml",
		})

	// Stage's RouteSettings is a RouteSettingsMap, so LoggingLevel is reachable
	// only by inventing a map key -- worse than the list case, which at least has
	// positions. Generation must fail rather than silently miss the reference.
	crd := testutil.GetCRDByName(t, g, "Stage")
	require.NotNil(crd)

	got, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 1)
	require.Error(err)
	assert.Contains(err.Error(), "references cannot be within a map")
	assert.Equal("", got, "nothing may be emitted when generation fails")
}

func Test_EnsureReferences_MissingAncestorField_IsRejected(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "s3",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-references.yaml",
		})

	crd := testutil.GetCRDByName(t, g, "Bucket")
	require.NotNil(crd)

	// With the model intact the struct-nested reference under Logging is emitted.
	// Establishing that first keeps the negative case below from being vacuous.
	before, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 1)
	require.NoError(err)
	require.Contains(before, "latestKO.Spec.Logging.LoggingEnabled.TargetBucketRef")

	// Drop the `Logging` ancestor, leaving the reference field that walks through
	// it. No generator.yaml can produce this -- the model always registers the
	// ancestors of a field it registers -- so reaching the guard means breaking
	// that invariant directly. The guard is what turns an inconsistent model into
	// a build failure naming the path instead of a nil dereference. The model is
	// built fresh per test, so the mutation cannot leak.
	require.Contains(crd.Fields, "Logging")
	delete(crd.Fields, "Logging")

	got, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 1)
	require.Error(err)
	assert.Contains(err.Error(), `unable to find field with path "Logging"`)
	assert.Contains(err.Error(), `resource "Bucket"`)
	assert.Equal("", got, "nothing may be emitted when generation fails")
}

func Test_EnsureReferences_MaterialisesMissingContainerOnTarget(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "s3",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-references.yaml",
		})

	// Generated set-output code rebuilds a struct from the response and nils it
	// when the response omits it, so the target can be missing a container the
	// source has. Restoring only the leaf would then be a no-op and the spec patch
	// would delete the whole declared block, which is the case the hand-maintained
	// hooks in eks/cluster, lambda/function and opensearchservice/domain exist to
	// prevent. Every container on the path must therefore be constructed.
	crd := testutil.GetCRDByName(t, g, "Bucket")
	require.NotNil(crd)

	got, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 1)
	require.NoError(err)

	// Both levels of the path, not just the innermost.
	assert.Contains(got, "if latestKO.Spec.Logging == nil {\n")
	assert.Contains(got, "latestKO.Spec.Logging = &svcapitypes.BucketLoggingStatus{}")
	assert.Contains(got, "if latestKO.Spec.Logging.LoggingEnabled == nil {\n")
	assert.Contains(got, "latestKO.Spec.Logging.LoggingEnabled = &svcapitypes.LoggingEnabled{}")

	// A container is only constructed once it is known the source holds a
	// reference to put in it, so a source that declares the container but no
	// reference leaves the target untouched. That ordering is what keeps an empty
	// container out of the patch.
	srcGuard := strings.Index(got, "if desiredKO.Spec.Logging.LoggingEnabled.TargetBucketRef != nil {")
	materialise := strings.Index(got, "if latestKO.Spec.Logging == nil {")
	require.NotEqual(-1, srcGuard)
	require.NotEqual(-1, materialise)
	assert.Less(srcGuard, materialise,
		"the source-side reference guard must enclose the materialisation")
}

func Test_EnsureReferences_GuardsAreNestedNotConcatenated(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	g := testutil.NewModelForServiceWithOptions(t, "s3",
		&testutil.TestingModelOptions{
			GeneratorConfigFile: "generator-with-nested-references.yaml",
		})

	// A single combined condition reaches several hundred characters on a deep
	// path -- 661 on ecs/CapacityProvider's three-level path -- and gofmt does not
	// wrap it. Nesting keeps every line readable and mirrors
	// ClearResolvedReferences, which walks the same paths.
	crd := testutil.GetCRDByName(t, g, "Bucket")
	require.NotNil(crd)

	got, err := code.EnsureReferences(crd, "desiredKO", "latestKO", 1)
	require.NoError(err)

	for _, line := range strings.Split(got, "\n") {
		assert.LessOrEqual(len(line), 120,
			"generated line should stay readable, got %d chars: %s", len(line), line)
	}
	assert.NotContains(got, " && ", "guards must be nested, not concatenated")
}
