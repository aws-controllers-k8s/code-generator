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
		namespace := ko.ObjectMeta.GetNamespace()
		if arr.Namespace != nil && *arr.Namespace != "" {
			namespace = *arr.Namespace
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
		namespace := ko.ObjectMeta.GetNamespace()
		if arr.Namespace != nil && *arr.Namespace != "" {
			namespace = *arr.Namespace
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
			namespace := ko.ObjectMeta.GetNamespace()
			if arr.Namespace != nil && *arr.Namespace != "" {
				namespace = *arr.Namespace
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
			namespace := ko.ObjectMeta.GetNamespace()
			if arr.Namespace != nil && *arr.Namespace != "" {
				namespace = *arr.Namespace
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
				namespace := ko.ObjectMeta.GetNamespace()
				if arr.Namespace != nil && *arr.Namespace != "" {
					namespace = *arr.Namespace
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
			namespace := ko.ObjectMeta.GetNamespace()
			if arr.Namespace != nil && *arr.Namespace != "" {
				namespace = *arr.Namespace
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
							namespace := ko.ObjectMeta.GetNamespace()
							if arr.Namespace != nil && *arr.Namespace != "" {
								namespace = *arr.Namespace
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
