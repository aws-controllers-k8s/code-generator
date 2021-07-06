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

	"github.com/aws-controllers-k8s/code-generator/pkg/model/multiversion"
)

func Test_computeRenames(t *testing.T) {
	type args struct {
		srcRenames map[string]string
		dstRenames map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "empty source and destination maps",
			args: args{
				srcRenames: nil,
				dstRenames: nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "empty source map",
			args: args{
				srcRenames: nil,
				dstRenames: map[string]string{
					"A": "B",
					"C": "D",
				},
			},
			want: map[string]string{
				"A": "B",
				"C": "D",
			},
			wantErr: false,
		},
		{
			name: "empty destination map",
			args: args{
				srcRenames: map[string]string{
					"A": "B",
					"C": "D",
				},
				dstRenames: nil,
			},
			want: map[string]string{
				"B": "A",
				"D": "C",
			},
			wantErr: false,
		},
		{
			name: "equal non-empty destination and source maps",
			args: args{
				srcRenames: map[string]string{
					"A": "B",
					"C": "D",
				},
				dstRenames: map[string]string{
					"A": "B",
					"C": "D",
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "different renames for the same field",
			args: args{
				srcRenames: map[string]string{
					"A": "B",
				},
				dstRenames: map[string]string{
					"A": "C",
				},
			},
			want: map[string]string{
				"B": "C",
			},
			wantErr: false,
		},
		{
			name: "different renames for the same field",
			args: args{
				srcRenames: map[string]string{
					"A": "B",
				},
				dstRenames: map[string]string{
					"A": "C",
				},
			},
			want: map[string]string{
				"B": "C",
			},
			wantErr: false,
		},
		{
			name: "conflicting renames",
			args: args{
				srcRenames: map[string]string{
					"A": "C",
				},
				dstRenames: map[string]string{
					"B": "C",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := multiversion.ComputeRenamesDelta(tt.args.srcRenames, tt.args.dstRenames)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeRenamesDelta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Since 1.12 formating functions prints maps in key-sorted order.
			// See https://golang.org/doc/go1.12#fmt
			if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", tt.want) {
				t.Errorf("computeRenamesDelta() = %v, want %v", got, tt.want)
			}
		})
	}
}
