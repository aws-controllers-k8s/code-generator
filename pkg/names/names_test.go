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
		{"AmiLaunchIndex", "AMILaunchIndex", "amiLaunchIndex", "ami_launch_index"},
		{"Amis", "AMIs", "amis", "amis"},
		{"AmiType", "AMIType", "amiType", "ami_type"},
		{"CacheSecurityGroup", "CacheSecurityGroup", "cacheSecurityGroup", "cache_security_group"},
		{"Camila", "Camila", "camila", "camila"},
		{"DbInstanceId", "DBInstanceID", "dbInstanceID", "db_instance_id"},
		{"DBInstanceId", "DBInstanceID", "dbInstanceID", "db_instance_id"},
		{"DBInstanceID", "DBInstanceID", "dbInstanceID", "db_instance_id"},
		{"DBInstanceIdentifier", "DBInstanceIdentifier", "dbInstanceIdentifier", "db_instance_identifier"},
		{"DbiResourceId", "DBIResourceID", "dbiResourceID", "dbi_resource_id"},
		{"DpdTimeoutAction", "DPDTimeoutAction", "dpdTimeoutAction", "dpd_timeout_action"},
		{"Dynamic", "Dynamic", "dynamic", "dynamic"},
		{"Ecmp", "ECMP", "ecmp", "ecmp"},
		{"EdiPartyName", "EDIPartyName", "ediPartyName", "edi_party_name"},
		{"Editable", "Editable", "editable", "editable"},
		{"Ena", "ENA", "ena", "ena"},
		{"Examine", "Examine", "examine", "examine"},
		{"Family", "Family", "family", "family"},
		{"Id", "ID", "id", "id"},
		{"ID", "ID", "id", "id"},
		{"Identifier", "Identifier", "identifier", "identifier"},
		{"IoPerformance", "IOPerformance", "ioPerformance", "io_performance"},
		{"Iops", "IOPS", "iops", "iops"},
		{"Ip", "IP", "ip", "ip"},
		{"Frame", "Frame", "frame", "frame"},
		{"KeyId", "KeyID", "keyID", "key_id"},
		{"KeyID", "KeyID", "keyID", "key_id"},
		{"KeyIdentifier", "KeyIdentifier", "keyIdentifier", "key_identifier"},
		{"LdapServerMetadata", "LDAPServerMetadata", "ldapServerMetadata", "ldap_server_metadata"},
		{"MaxIdleConnectionsPercent", "MaxIdleConnectionsPercent", "maxIdleConnectionsPercent", "max_idle_connections_percent"},
		{"MultipartUpload", "MultipartUpload", "multipartUpload", "multipart_upload"},
		{"Nat", "NAT", "nat", "nat"},
		{"NatGateway", "NATGateway", "natGateway", "nat_gateway"},
		{"NativeAuditFieldsIncluded", "NativeAuditFieldsIncluded", "nativeAuditFieldsIncluded", "native_audit_fields_included"},
		{"NumberOfAmiToKeep", "NumberOfAMIToKeep", "numberOfAMIToKeep", "number_of_ami_to_keep"},
		{"Package", "Package", "package_", "package_"},
		{"Param", "Param", "param", "param"},
		{"Ram", "RAM", "ram", "ram"},
		{"RamDiskId", "RAMDiskID", "ramDiskID", "ram_disk_id"},
		{"RepositoryUriTest", "RepositoryURITest", "repositoryURITest", "repository_uri_test"},
		{"RequestedAmiVersion", "RequestedAMIVersion", "requestedAMIVersion", "requested_ami_version"},
		{"SriovNetSupport", "SRIOVNetSupport", "sriovNetSupport", "sriov_net_support"},
		{"SSEKMSKeyID", "SSEKMSKeyID", "sseKMSKeyID", "sse_kms_key_id"},
		{"UUID", "UUID", "uuid", "uuid"},
		{"Vlan", "VLAN", "vlan", "vlan"},
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
