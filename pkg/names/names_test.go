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

package names_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws-controllers-k8s/code-generator/pkg/names"
)

func TestNames(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		original         string
		expectCamel      string
		expectCamelLower string
		expectSnake      string
	}{
		{"Ami", "AMI", "ami", "ami"},
		{"Identifier", "Identifier", "identifier", "identifier"},
		{"Id", "ID", "id", "id"},
		{"ID", "ID", "id", "id"},
		{"Family", "Family", "family", "family"},
		{"Examine", "Examine", "examine", "examine"},
		{"Dynamic", "Dynamic", "dynamic", "dynamic"},
		{"Camila", "Camila", "camila", "camila"},
		{"Amis", "AMIs", "amis", "amis"},
		{"AmiType", "AMIType", "amiType", "ami_type"},
		{"AmiLaunchIndex", "AMILaunchIndex", "amiLaunchIndex", "ami_launch_index"},
		{"NumberOfAmiToKeep", "NumberOfAMIToKeep", "numberOfAMIToKeep", "number_of_ami_to_keep"},
		{"RequestedAmiVersion", "RequestedAMIVersion", "requestedAMIVersion", "requested_ami_version"},
		{"KeyIdentifier", "KeyIdentifier", "keyIdentifier", "key_identifier"},
		{"KeyId", "KeyID", "keyID", "key_id"},
		{"KeyID", "KeyID", "keyID", "key_id"},
		{"SSEKMSKeyID", "SSEKMSKeyID", "sseKMSKeyID", "sse_kms_key_id"},
		{"DbiResourceId", "DBIResourceID", "dbiResourceID", "dbi_resource_id"},
		{"DbInstanceId", "DBInstanceID", "dbInstanceID", "db_instance_id"},
		{"DBInstanceId", "DBInstanceID", "dbInstanceID", "db_instance_id"},
		{"DBInstanceID", "DBInstanceID", "dbInstanceID", "db_instance_id"},
		{"DBInstanceIdentifier", "DBInstanceIdentifier", "dbInstanceIdentifier", "db_instance_identifier"},
		{"MaxIdleConnectionsPercent", "MaxIdleConnectionsPercent", "maxIdleConnectionsPercent", "max_idle_connections_percent"},
		{"CacheSecurityGroup", "CacheSecurityGroup", "cacheSecurityGroup", "cache_security_group"},
		{"RepositoryUriTest", "RepositoryURITest", "repositoryURITest", "repository_uri_test"},
		{"Ip", "IP", "ip", "ip"},
		{"MultipartUpload", "MultipartUpload", "multipartUpload", "multipart_upload"},
		{"Package", "Package", "package_", "package_"},
		{"LdapServerMetadata", "LDAPServerMetadata", "ldapServerMetadata", "ldap_server_metadata"},
		{"DpdTimeoutAction", "DPDTimeoutAction", "dpdTimeoutAction", "dpd_timeout_action"},
		{"Iops", "IOPS", "iops", "iops"},
		{"IoPerformance", "IOPerformance", "ioPerformance", "io_performance"},
		{"Vlan", "VLAN", "vlan", "vlan"},
		{"Ecmp", "ECMP", "ecmp", "ecmp"},
		{"Ena", "ENA", "ena", "ena"},
		{"UUID", "UUID", "uuid", "uuid"},
		{"Nat", "NAT", "nat", "nat"},
		{"NatGateway", "NATGateway", "natGateway", "nat_gateway"},
		{"NativeAuditFieldsIncluded", "NativeAuditFieldsIncluded", "nativeAuditFieldsIncluded", "native_audit_fields_included"},
	}
	for _, tc := range testCases {
		n := names.New(tc.original)
		msg := fmt.Sprintf("for original %s expected camel name of %s but got %s", tc.original, tc.expectCamel, n.Camel)
		assert.Equal(tc.expectCamel, n.Camel, msg)
		msg = fmt.Sprintf("for original %s expected lowercase camel name of %s but got %s", tc.original, tc.expectCamelLower, n.CamelLower)
		assert.Equal(tc.expectCamelLower, n.CamelLower, msg)
		msg = fmt.Sprintf("for original %s expected snake name of %s but got %s", tc.original, tc.expectSnake, n.Snake)
		assert.Equal(tc.expectSnake, n.Snake, msg)
	}
}
