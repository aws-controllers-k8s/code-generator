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

package config

import (
	"testing"
)

// opBinding describes an expected operation binding after inference.
type opBinding struct {
	opID    string
	resName string
	opType  string
	isBound bool // true = should exist, false = should NOT exist
}

func TestInferManagedFieldOps(t *testing.T) {
	tests := []struct {
		name       string
		parentName string
		fieldName  string
		sdkOps     []string
		// Pre-existing operation bindings (simulates explicit user config)
		existing map[string]OperationConfig
		expect   []opBinding
	}{
		{
			name:       "standard {Verb}{Parent}{Field} pattern",
			parentName: "BackupVault",
			fieldName:  "LockConfiguration",
			sdkOps: []string{
				"PutBackupVaultLockConfiguration",
				"DeleteBackupVaultLockConfiguration",
				"GetBackupVaultLockConfiguration",
			},
			expect: []opBinding{
				{"PutBackupVaultLockConfiguration", "LockConfiguration", "create", true},
				{"DeleteBackupVaultLockConfiguration", "LockConfiguration", "delete", true},
				{"GetBackupVaultLockConfiguration", "LockConfiguration", "get", true},
			},
		},
		{
			name:       "Put preferred over Create",
			parentName: "BackupVault",
			fieldName:  "AccessPolicy",
			sdkOps: []string{
				"PutBackupVaultAccessPolicy",
				"CreateBackupVaultAccessPolicy",
				"DeleteBackupVaultAccessPolicy",
			},
			expect: []opBinding{
				{"PutBackupVaultAccessPolicy", "AccessPolicy", "create", true},
				{"CreateBackupVaultAccessPolicy", "", "", false},
			},
		},
		{
			name:       "fallback to {Verb}{Field} without parent",
			parentName: "BackupVault",
			fieldName:  "Notifications",
			sdkOps: []string{
				"PutNotifications",
				"DeleteNotifications",
				"GetNotifications",
			},
			expect: []opBinding{
				{"PutNotifications", "Notifications", "create", true},
				{"DeleteNotifications", "Notifications", "delete", true},
				{"GetNotifications", "Notifications", "get", true},
			},
		},
		{
			name:       "singular form: Tags -> Tag for TagResource/UntagResource",
			parentName: "BackupVault",
			fieldName:  "Tags",
			sdkOps: []string{
				"TagResource",
				"UntagResource",
				"ListTags",
			},
			expect: []opBinding{
				{"TagResource", "Tags", "create", true},
				{"UntagResource", "Tags", "delete", true},
				{"ListTags", "Tags", "get", true},
			},
		},
		{
			name:       "generic Resource suffix as last resort",
			parentName: "MyParent",
			fieldName:  "SomeField",
			sdkOps: []string{
				"CreateResource",
				"DeleteResource",
				"GetResource",
			},
			expect: []opBinding{
				{"CreateResource", "SomeField", "create", true},
				{"DeleteResource", "SomeField", "delete", true},
				{"GetResource", "SomeField", "get", true},
			},
		},
		{
			name:       "{Parent}{Field} preferred over {Field} only",
			parentName: "BackupVault",
			fieldName:  "Notifications",
			sdkOps: []string{
				"PutBackupVaultNotifications",
				"PutNotifications",
				"DeleteBackupVaultNotifications",
				"DeleteNotifications",
			},
			expect: []opBinding{
				{"PutBackupVaultNotifications", "Notifications", "create", true},
				{"PutNotifications", "", "", false},
				{"DeleteBackupVaultNotifications", "Notifications", "delete", true},
				{"DeleteNotifications", "", "", false},
			},
		},
		{
			name:       "explicit binding not overwritten",
			parentName: "BackupVault",
			fieldName:  "AccessPolicy",
			sdkOps: []string{
				"PutBackupVaultAccessPolicy",
				"DeleteBackupVaultAccessPolicy",
			},
			existing: map[string]OperationConfig{
				"PutBackupVaultAccessPolicy": {
					ResourceName:  StringArray{"CustomName"},
					OperationType: StringArray{"create"},
				},
			},
			expect: []opBinding{
				{"PutBackupVaultAccessPolicy", "CustomName", "create", true},
				{"DeleteBackupVaultAccessPolicy", "AccessPolicy", "delete", true},
			},
		},
		{
			name:       "no matching ops produces no bindings",
			parentName: "BackupVault",
			fieldName:  "LockConfiguration",
			sdkOps: []string{
				"CreateBucket",
				"DeleteBucket",
			},
			expect: []opBinding{},
		},
		{
			name:       "partial match infers what it can",
			parentName: "BackupVault",
			fieldName:  "LockConfiguration",
			sdkOps: []string{
				"PutBackupVaultLockConfiguration",
				"DeleteBackupVaultLockConfiguration",
			},
			expect: []opBinding{
				{"PutBackupVaultLockConfiguration", "LockConfiguration", "create", true},
				{"DeleteBackupVaultLockConfiguration", "LockConfiguration", "delete", true},
			},
		},
		{
			name:       "Get preferred over Describe",
			parentName: "BackupVault",
			fieldName:  "Policy",
			sdkOps: []string{
				"GetBackupVaultPolicy",
				"DescribeBackupVaultPolicy",
			},
			expect: []opBinding{
				{"GetBackupVaultPolicy", "Policy", "get", true},
				{"DescribeBackupVaultPolicy", "", "", false},
			},
		},
		{
			name:       "List works as get fallback",
			parentName: "MyResource",
			fieldName:  "Tags",
			sdkOps: []string{
				"ListTags",
			},
			expect: []opBinding{
				{"ListTags", "Tags", "get", true},
			},
		},
		{
			name:       "non-plural field skips singular candidates",
			parentName: "BackupVault",
			fieldName:  "AccessPolicy",
			sdkOps: []string{
				"PutBackupVaultAccessPolicy",
				"DeleteBackupVaultAccessPolicy",
			},
			expect: []opBinding{
				{"PutBackupVaultAccessPolicy", "AccessPolicy", "create", true},
				{"DeleteBackupVaultAccessPolicy", "AccessPolicy", "delete", true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Operations: make(map[string]OperationConfig),
			}
			for k, v := range tt.existing {
				cfg.Operations[k] = v
			}

			sdkOps := make(map[string]struct{}, len(tt.sdkOps))
			for _, op := range tt.sdkOps {
				sdkOps[op] = struct{}{}
			}

			cfg.inferManagedFieldOps(tt.parentName, tt.fieldName, sdkOps)

			for _, exp := range tt.expect {
				op, exists := cfg.Operations[exp.opID]
				if exp.isBound {
					if !exists {
						t.Errorf("expected %q to be bound, but it was not", exp.opID)
						continue
					}
					if len(op.ResourceName) != 1 || op.ResourceName[0] != exp.resName {
						t.Errorf("%q: resource_name = %v, want [%s]", exp.opID, op.ResourceName, exp.resName)
					}
					if len(op.OperationType) != 1 || op.OperationType[0] != exp.opType {
						t.Errorf("%q: operation_type = %v, want [%s]", exp.opID, op.OperationType, exp.opType)
					}
				} else {
					if exists && len(op.ResourceName) > 0 && len(op.OperationType) > 0 {
						t.Errorf("expected %q to NOT be bound, but got resource_name=%v operation_type=%v",
							exp.opID, op.ResourceName, op.OperationType)
					}
				}
			}
		})
	}
}
