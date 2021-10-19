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

package multiversion_test

import (
	"fmt"
	"testing"

	awssdkmodel "github.com/aws/aws-sdk-go/private/model/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ackmodel "github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/code-generator/pkg/model/multiversion"
	"github.com/aws-controllers-k8s/code-generator/pkg/testutil"
)

func TestComputeFieldsDiff_ECR(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	type expectDelta struct {
		fieldName  string
		changeType multiversion.FieldChangeType
	}
	type testAPIInfo struct {
		apiVersion    string // e.g v1alpha1
		sdkAPIVersion string // e.g 00-00-0000
		// if given tests will change the shape.Type to `type struct{}`
		mutateShape string
	}
	type args struct {
		dst testAPIInfo
		src testAPIInfo
	}
	tests := []struct {
		name    string
		args    args
		want    []expectDelta
		wantErr bool
	}{
		{
			name: "v1alpha1-v1alpha1: no changes.",
			args: args{
				src: testAPIInfo{
					apiVersion: "v1alpha1",
				},
				dst: testAPIInfo{
					apiVersion: "v1alpha1",
				},
			},
			wantErr: false,
			want: []expectDelta{
				{
					fieldName:  "ImageScanningConfiguration",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "ImageTagMutability",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "RepositoryName",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Tags",
					changeType: multiversion.FieldChangeTypeNone,
				},
			},
		},
		{
			name: "v1alpha1-v1alpha2: renamed fields.",
			args: args{
				src: testAPIInfo{
					apiVersion: "v1alpha1",
				},
				dst: testAPIInfo{
					apiVersion: "v1alpha2",
				},
			},
			wantErr: false,
			want: []expectDelta{
				{
					fieldName:  "ImageTagMutability",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Name",
					changeType: multiversion.FieldChangeTypeRenamed,
				},
				{
					fieldName:  "ScanConfig",
					changeType: multiversion.FieldChangeTypeRenamed,
				},
				{
					fieldName:  "Tags",
					changeType: multiversion.FieldChangeTypeNone,
				},
			},
		},
		{
			name: "v1alpha2-v1beta1: added fields from aws-sdk-go",
			args: args{
				src: testAPIInfo{
					apiVersion: "v1alpha2",
				},
				dst: testAPIInfo{
					apiVersion:    "v1beta1",
					sdkAPIVersion: "0000-00-01",
				},
			},
			wantErr: false,
			want: []expectDelta{
				{
					fieldName:  "EncryptionConfiguration",
					changeType: multiversion.FieldChangeTypeAdded,
				},
				{
					fieldName:  "ImageTagMutability",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Name",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "ScanConfig",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Tags",
					changeType: multiversion.FieldChangeTypeNone,
				},
			},
		},
		{
			name: "v1beta1-v1beta2: added fields from generator configuration",
			args: args{
				src: testAPIInfo{
					apiVersion:    "v1beta1",
					sdkAPIVersion: "0000-00-01",
				},
				dst: testAPIInfo{
					apiVersion:    "v1beta2",
					sdkAPIVersion: "0000-00-01",
				},
			},
			wantErr: false,
			want: []expectDelta{
				{
					fieldName:  "AnotherNameField",
					changeType: multiversion.FieldChangeTypeAdded,
				},
				{
					fieldName:  "EncryptionConfiguration",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "ImageTagMutability",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Name",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "ScanConfig",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Tags",
					changeType: multiversion.FieldChangeTypeNone,
				},
			},
		},
		{
			name: "v1beta2-v1beta1: removed fields from generator configuration",
			args: args{
				src: testAPIInfo{
					apiVersion:    "v1beta2",
					sdkAPIVersion: "0000-00-01",
				},
				dst: testAPIInfo{
					apiVersion:    "v1beta1",
					sdkAPIVersion: "0000-00-01",
				},
			},
			wantErr: false,
			want: []expectDelta{
				{
					fieldName:  "AnotherNameField",
					changeType: multiversion.FieldChangeTypeRemoved,
				},
				{
					fieldName:  "EncryptionConfiguration",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "ImageTagMutability",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Name",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "ScanConfig",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Tags",
					changeType: multiversion.FieldChangeTypeNone,
				},
			},
		},
		{
			name: "v1alpha2-v1alpha3: field change from string to secret",
			args: args{
				src: testAPIInfo{
					apiVersion: "v1alpha2",
				},
				dst: testAPIInfo{
					apiVersion: "v1alpha3",
				},
			},
			wantErr: false,
			want: []expectDelta{
				{
					fieldName:  "ImageTagMutability",
					changeType: multiversion.FieldChangeTypeShapeChangedFromStringToSecret,
				},
				{
					fieldName:  "Name",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "ScanConfig",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Tags",
					changeType: multiversion.FieldChangeTypeNone,
				},
			},
		},
		{
			name: "v1alpha3-v1alpha2: field change from secret to string",
			args: args{
				src: testAPIInfo{
					apiVersion: "v1alpha3",
				},
				dst: testAPIInfo{
					apiVersion: "v1alpha2",
				},
			},
			wantErr: false,
			want: []expectDelta{
				{
					fieldName:  "ImageTagMutability",
					changeType: multiversion.FieldChangeTypeShapeChangedFromSecretToString,
				},
				{
					fieldName:  "Name",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "ScanConfig",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Tags",
					changeType: multiversion.FieldChangeTypeNone,
				},
			},
		},
		{
			name: "v1alpha2-v1alpha3: field type changed",
			args: args{
				src: testAPIInfo{
					apiVersion: "v1alpha2",
				},
				dst: testAPIInfo{
					apiVersion:  "v1alpha3",
					mutateShape: "Tags",
				},
			},
			wantErr: false,
			want: []expectDelta{
				{
					fieldName:  "ImageTagMutability",
					changeType: multiversion.FieldChangeTypeShapeChangedFromStringToSecret,
				},
				{
					fieldName:  "Name",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "ScanConfig",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Tags",
					changeType: multiversion.FieldChangeTypeShapeChanged,
				},
			},
		},
		{
			name: "v1alpha3-v1alpha2: field type changed",
			args: args{
				src: testAPIInfo{
					apiVersion: "v1alpha3",
				},
				dst: testAPIInfo{
					apiVersion:  "v1alpha2",
					mutateShape: "Tags",
				},
			},
			wantErr: false,
			want: []expectDelta{
				{
					fieldName:  "ImageTagMutability",
					changeType: multiversion.FieldChangeTypeShapeChangedFromSecretToString,
				},
				{
					fieldName:  "Name",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "ScanConfig",
					changeType: multiversion.FieldChangeTypeNone,
				},
				{
					fieldName:  "Tags",
					changeType: multiversion.FieldChangeTypeShapeChanged,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dstGeneratorFile := fmt.Sprintf("generator-%s.yaml", tt.args.dst.apiVersion)
			dstModel := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
				GeneratorConfigFile: dstGeneratorFile,
				APIVersion:          tt.args.dst.apiVersion,
				ServiceAPIVersion:   tt.args.dst.sdkAPIVersion,
			})

			srcGeneratorFile := fmt.Sprintf("generator-%s.yaml", tt.args.src.apiVersion)
			srcModel := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{
				GeneratorConfigFile: srcGeneratorFile,
				APIVersion:          tt.args.src.apiVersion,
				ServiceAPIVersion:   tt.args.src.sdkAPIVersion,
			})

			dstCRDs, err := dstModel.GetCRDs()
			require.Nil(err)
			srcCRDs, err := srcModel.GetCRDs()
			require.Nil(err)
			require.Len(dstCRDs, 1)
			require.Len(srcCRDs, 1)

			dstCRD := dstCRDs[0]
			if tt.args.dst.mutateShape != "" {
				field, ok := dstCRD.Fields[tt.args.dst.mutateShape]
				require.True(ok)
				field.ShapeRef.Shape = newEmptyStructShape(tt.args.dst.mutateShape)
			}
			srcCRD := srcCRDs[0]
			if tt.args.src.mutateShape != "" {
				field, ok := srcCRD.Fields[tt.args.src.mutateShape]
				require.True(ok)
				field.ShapeRef.Shape = newEmptyStructShape(tt.args.src.mutateShape)
			}

			srcRenames, err := srcCRD.GetAllRenames(ackmodel.OpTypeCreate)
			require.Nil(err)
			dstRenames, err := dstCRD.GetAllRenames(ackmodel.OpTypeCreate)
			require.Nil(err)
			renames, err := multiversion.ComputeRenamesDelta(srcRenames, dstRenames)
			require.Nil(err)

			require.Nil(err)
			deltas, err := multiversion.ComputeFieldDeltas(
				srcCRD.SpecFields, dstCRD.SpecFields, renames,
			)
			if (err != nil) != tt.wantErr {
				require.Fail(fmt.Sprintf("multiversion.ComputeFieldsDelta() error = %v, wantErr %v", err, tt.wantErr))
			}

			require.Equal(len(deltas), len(tt.want))
			for i, delta := range deltas {
				require.Equal(tt.want[i].changeType, delta.ChangeType)

				// Additional checks
				switch delta.ChangeType {
				case multiversion.FieldChangeTypeNone:
					assert.Equal(tt.want[i].fieldName, delta.Source.Names.Camel)
					assert.Equal(tt.want[i].fieldName, delta.Destination.Names.Camel)
				case multiversion.FieldChangeTypeRemoved:
					assert.Equal(tt.want[i].fieldName, delta.Source.Names.Camel)
				case multiversion.FieldChangeTypeAdded:
					assert.Equal(tt.want[i].fieldName, delta.Destination.Names.Camel)
				case multiversion.FieldChangeTypeRenamed:
					assert.Equal(tt.want[i].fieldName, delta.Destination.Names.Camel)
				}
			}
		})
	}
}

func newEmptyStructShape(name string) *awssdkmodel.Shape {
	return &awssdkmodel.Shape{
		Type:      "structure",
		ShapeName: name,
	}
}

func TestAreEqualShapes_ECR_Repository(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	model := testutil.NewModelForServiceWithOptions(t, "ecr", &testutil.TestingModelOptions{})
	crds, err := model.GetCRDs()
	require.Nil(err)
	require.Len(crds, 1)
	repositoryCRD := crds[0]

	delete(repositoryCRD.SpecFields, "ImageTagMutability")

	for _, fieldNameX := range repositoryCRD.SpecFieldNames() {
		for _, fieldNameY := range repositoryCRD.SpecFieldNames() {
			testName := fmt.Sprintf("comparing %s with %s", fieldNameX, fieldNameY)
			t.Run(testName, func(t *testing.T) {

				fieldX := repositoryCRD.SpecFields[fieldNameX]
				fieldY := repositoryCRD.SpecFields[fieldNameY]
				equal, _ := multiversion.AreEqualShapes(fieldX.ShapeRef.Shape, fieldY.ShapeRef.Shape, true)
				if fieldNameY == fieldNameX {
					assert.True(equal)
				} else {
					assert.False(equal)
				}
			})
		}
	}
}

func TestAreEqualShapes_APIGatewayV2_DomainName(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	model := testutil.NewModelForServiceWithOptions(t, "apigatewayv2", &testutil.TestingModelOptions{})
	crds, err := model.GetCRDs()
	require.Nil(err)
	require.Len(crds, 12)
	domainNameCRD := crds[4]

	for _, fieldNameX := range domainNameCRD.SpecFieldNames() {
		for _, fieldNameY := range domainNameCRD.SpecFieldNames() {
			testName := fmt.Sprintf("comparing %s with %s", fieldNameX, fieldNameY)
			t.Run(testName, func(t *testing.T) {

				fieldX := domainNameCRD.SpecFields[fieldNameX]
				fieldY := domainNameCRD.SpecFields[fieldNameY]
				equal, _ := multiversion.AreEqualShapes(fieldX.ShapeRef.Shape, fieldY.ShapeRef.Shape, false)
				if fieldNameY == fieldNameX {
					assert.True(equal)
				} else {
					assert.False(equal)
				}
			})
		}
	}
}

func TestComputeCRDDeltas_APIGatewayV2_DomainName(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	model := testutil.NewModelForServiceWithOptions(t, "apigatewayv2", &testutil.TestingModelOptions{})
	crds, err := model.GetCRDs()
	require.Nil(err)
	require.Len(crds, 12)
	domainNameCRD := crds[4]

	deltas, err := multiversion.ComputeCRDFieldDeltas(domainNameCRD, domainNameCRD)
	require.Nil(err)
	assert.Len(deltas.SpecDeltas, len(domainNameCRD.SpecFields))
	assert.Len(deltas.StatusDeltas, len(domainNameCRD.StatusFields))

	for _, delta := range deltas.SpecDeltas {
		assert.Equal(delta.ChangeType, multiversion.FieldChangeTypeNone)
	}
	for _, delta := range deltas.StatusDeltas {
		assert.Equal(delta.ChangeType, multiversion.FieldChangeTypeNone)
	}
}
