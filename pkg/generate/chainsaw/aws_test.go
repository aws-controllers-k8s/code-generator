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

package chainsaw

import (
	"fmt"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
	"github.com/aws-controllers-k8s/pkg/names"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
)

func TestCLIName(t *testing.T) {
	for in, want := range map[string]string{
		"DescribeVpcs":             "describe-vpcs",
		"VpcIds":                   "vpc-ids",
		"InternetGatewayIds":       "internet-gateway-ids",
		"DescribeSecurityGroups":   "describe-security-groups",
		"DeleteSecurityGroup":      "delete-security-group",
		"DescribeDhcpOptions":      "describe-dhcp-options",
		"GetBucketVersioning":      "get-bucket-versioning",
		"ListTagsForResourceInput": "list-tags-for-resource-input",
	} {
		assert.Equal(t, want, cliName(in), in)
	}
}

// ipPermissions is a list of structures with a nested list of structures,
// like EC2's IpPermissions
var ipPermissions = &awssdkmodel.ShapeRef{Shape: &awssdkmodel.Shape{
	Type: "list",
	MemberRef: awssdkmodel.ShapeRef{Shape: &awssdkmodel.Shape{
		Type: "structure",
		MemberRefs: map[string]*awssdkmodel.ShapeRef{
			"FromPort":   {Shape: &awssdkmodel.Shape{Type: "integer"}},
			"IpProtocol": {Shape: &awssdkmodel.Shape{Type: "string"}},
			"IpRanges": {Shape: &awssdkmodel.Shape{
				Type: "list",
				MemberRef: awssdkmodel.ShapeRef{Shape: &awssdkmodel.Shape{
					Type: "structure",
					MemberRefs: map[string]*awssdkmodel.ShapeRef{
						"CidrIp":      {Shape: &awssdkmodel.Shape{Type: "string"}},
						"Description": {Shape: &awssdkmodel.Shape{Type: "string"}},
					},
				}},
			}},
			"ToPort": {Shape: &awssdkmodel.Shape{Type: "integer"}},
		},
	}},
}}

var crRules = []interface{}{map[string]interface{}{
	"ipProtocol": "tcp",
	"fromPort":   443,
	"ipRanges":   []interface{}{map[string]interface{}{"cidrIP": "10.0.0.0/16"}},
}}

func TestToAWS(t *testing.T) {
	assert.Equal(t, []interface{}{map[string]interface{}{
		"IpProtocol": "tcp",
		"FromPort":   443,
		"IpRanges":   []interface{}{map[string]interface{}{"CidrIp": "10.0.0.0/16"}},
	}}, toAWS(crRules, ipPermissions))
}

func TestProjection(t *testing.T) {
	// Only what the CR value sets, since AWS returns more
	assert.Equal(t,
		"[].{FromPort: FromPort, IpProtocol: IpProtocol, IpRanges: IpRanges[].{CidrIp: CidrIp}}",
		projection(crRules, ipPermissions),
	)
	assert.Equal(t, "", projection("default", &awssdkmodel.ShapeRef{Shape: &awssdkmodel.Shape{Type: "string"}}))
}

func TestMerge(t *testing.T) {
	base := map[string]interface{}{
		"cidrBlocks": []interface{}{"10.0.0.0/16", "10.2.0.0/16"},
		"vpcRef":     map[string]interface{}{"from": map[string]interface{}{"name": "a"}},
		"tags":       []interface{}{"x"},
	}
	got := merge(base, map[string]interface{}{
		"cidrBlocks": []interface{}{"10.0.0.0/16"},
		"vpcRef":     map[string]interface{}{"from": map[string]interface{}{"namespace": "b"}},
		"tags":       nil,
	})
	assert.Equal(t, map[string]interface{}{
		// Lists are replaced, maps merged, and null kept so a merge patch
		// removes the field
		"cidrBlocks": []interface{}{"10.0.0.0/16"},
		"vpcRef":     map[string]interface{}{"from": map[string]interface{}{"name": "a", "namespace": "b"}},
		"tags":       nil,
	}, got)
	// The base is unchanged
	assert.Equal(t, []interface{}{"10.0.0.0/16", "10.2.0.0/16"}, base["cidrBlocks"])
}

func TestChangedValue(t *testing.T) {
	for in, want := range map[interface{}]interface{}{"a": "a-changed", true: false, 1: 2} {
		got, ok := changedValue(in)
		assert.True(t, ok)
		assert.Equal(t, want, got)
	}
	_, ok := changedValue([]interface{}{"a"})
	assert.False(t, ok)
}

func TestProjectionEmptyList(t *testing.T) {
	// Nothing to narrow, and {} isn't valid JMESPath
	assert.Equal(t, "", projection([]interface{}{}, ipPermissions))
}

func TestBestMember(t *testing.T) {
	both := []string{"IpPermissions", "IpPermissionsEgress"}
	// Egress shares a word with IpPermissionsEgress only
	assert.Equal(t, "IpPermissionsEgress", bestMember("EgressRules", both, nil))
	// Ingress shares a word with neither, so it's a tie...
	assert.Equal(t, "", bestMember("IngressRules", both, nil))
	// ...until the other is taken
	assert.Equal(t, "IpPermissions", bestMember("IngressRules", both, map[string]bool{"IpPermissionsEgress": true}))
}

func TestFindPath(t *testing.T) {
	boolValue := &awssdkmodel.ShapeRef{Shape: &awssdkmodel.Shape{Type: "boolean"}}
	output := &awssdkmodel.Shape{Type: "structure", MemberRefs: map[string]*awssdkmodel.ShapeRef{
		"EnableDnsSupport": {Shape: &awssdkmodel.Shape{Type: "structure", MemberRefs: map[string]*awssdkmodel.ShapeRef{
			"Value": boolValue,
		}}},
		"DefaultVersionNumber": {Shape: &awssdkmodel.Shape{Type: "long"}},
	}}
	path, ref := findPath(output, []string{"EnableDnsSupport", "Value"})
	assert.Equal(t, "EnableDnsSupport.Value", path)
	assert.Equal(t, boolValue, ref)
	// A write operation's wrapper isn't in the read's output
	path, _ = findPath(output, []string{"LaunchTemplate", "DefaultVersionNumber"})
	assert.Equal(t, "DefaultVersionNumber", path)
	_, ref = findPath(output, []string{"Missing"})
	assert.Nil(t, ref)
}

func TestPairedRead(t *testing.T) {
	ops := map[string]*awssdkmodel.Operation{
		"DescribeVpcAttribute":             {},
		"DescribeInternetGateways":         {},
		"GetFunctionEventInvokeConfig":     {},
		"GetBucketAccelerateConfiguration": {},
	}
	for subject, want := range map[string]string{
		"VpcAttribute":                  "DescribeVpcAttribute",
		"InternetGateway":               "DescribeInternetGateways",
		"FunctionEventInvokeConfig":     "GetFunctionEventInvokeConfig",
		"BucketAccelerateConfiguration": "GetBucketAccelerateConfiguration",
		"VpcPeeringConnectionOptions":   "",
	} {
		assert.Equal(t, want, pairedRead(ops, subject), subject)
	}
}

func TestSortScalars(t *testing.T) {
	assert.Equal(t, []interface{}{"10.0.0.0/16", "10.2.0.0/16"}, sortScalars([]interface{}{"10.2.0.0/16", "10.0.0.0/16"}))
	assert.Equal(t, true, sortScalars(true))
}

func TestTemplated(t *testing.T) {
	g := &generator{in: &Inputs{Values: []Value{{Name: "AMI_ID"}, {Name: "AMI"}}}}
	assert.Equal(t, "plain", g.templated("plain"))
	assert.Equal(t, "(join('', ['ack-test-', $suffix]))", g.templated("ack-test-$SUFFIX"))
	assert.Equal(t, "(join('', [$suffix, '-it\\'s-', $suffix]))", g.templated("$SUFFIX-it's-$SUFFIX"))
	// A value, not mistaken for a shorter one, and other $NAMEs left alone
	cm := "x_k8s_get($client, 'v1', 'ConfigMap', $namespace, 'ack-test-values').data."
	assert.Equal(t, "(join('', ["+cm+"AMI_ID, '/', "+cm+"AMI, '-$AWS_ID']))", g.templated("$AMI_ID/$AMI-$AWS_ID"))
	assert.Equal(t, "$AWS_ID", g.templated("$AWS_ID"))
	assert.Equal(t, map[string]interface{}{
		"name": "(join('', ['sg-', $suffix]))",
		"tags": []interface{}{map[string]interface{}{"key": "k", "value": "(join('', ['v-', $suffix]))"}},
		"port": 443,
	}, g.templated(map[string]interface{}{
		"name": "sg-$SUFFIX",
		"tags": []interface{}{map[string]interface{}{"key": "k", "value": "v-$SUFFIX"}},
		"port": 443,
	}))
}

func TestJSONParsed(t *testing.T) {
	assert.Equal(t, map[string]interface{}{
		"(json_parse(policyDocument))": map[string]interface{}{
			"Version":   "2012-10-17",
			"Statement": []interface{}{map[string]interface{}{"Effect": "Allow", "Resource": "arn:aws:s3:::b-$SUFFIX"}},
		},
		"(json_parse(list))": []interface{}{float64(1), "a"},
		"name":               "{not json",
		"cidrBlocks":         []interface{}{"10.0.0.0/16"},
		"enabled":            true,
	}, jsonParsed(map[string]interface{}{
		"policyDocument": ` {"Version":"2012-10-17","Statement":[{"Effect":"Allow","Resource":"arn:aws:s3:::b-$SUFFIX"}]}`,
		"list":           `[1, "a"]`,
		"name":           "{not json",
		"cidrBlocks":     []interface{}{"10.0.0.0/16"},
		"enabled":        true,
	}))
}

func TestUnsuffixedNames(t *testing.T) {
	assert.Equal(t, []string{"groupName", "name"}, unsuffixedNames(
		map[string]interface{}{"name": "fixed", "description": "not a name"},
		map[string]interface{}{"groupName": "fixed", "roleName": "r-$SUFFIX", "name": "n-$SUFFIX", "bucketName": "$BUCKET"},
	))
}

func TestWrapperPath(t *testing.T) {
	instance := &awssdkmodel.Shape{Type: "structure", ShapeName: "Instance"}
	output := &awssdkmodel.Shape{Type: "structure", MemberRefs: map[string]*awssdkmodel.ShapeRef{
		"Reservations": {Shape: &awssdkmodel.Shape{Type: "list", MemberRef: awssdkmodel.ShapeRef{Shape: &awssdkmodel.Shape{
			Type: "structure",
			MemberRefs: map[string]*awssdkmodel.ShapeRef{
				"Instances": {Shape: &awssdkmodel.Shape{Type: "list", MemberRef: awssdkmodel.ShapeRef{Shape: instance}}},
			},
		}}}},
		"LaunchTemplate": {Shape: &awssdkmodel.Shape{Type: "structure"}},
	}}
	// Each list before the last member is read at its first element
	assert.Equal(t, "Reservations[0].Instances", wrapperPath(output, "Reservations.Instances"))
	assert.Equal(t, "LaunchTemplate", wrapperPath(output, "LaunchTemplate"))
}

func TestRefPath(t *testing.T) {
	assert.Equal(t, "status.vpcID", refPath("Status.VPCID"))
	assert.Equal(t, "status.ackResourceMetadata.arn", refPath("Status.ACKResourceMetadata.ARN"))
	assert.Equal(t, "spec.name", refPath("Spec.Name"))
}

func TestRefValue(t *testing.T) {
	a := &awsAPI{apiVersions: map[string]string{"VPC": "ec2.services.k8s.aws/v1alpha1"}}
	vpc := refTarget{concrete: "vpcID", kind: "VPC", group: "ec2.services.k8s.aws", path: "status.vpcID"}
	from := func(name string) map[string]interface{} {
		return map[string]interface{}{"from": map[string]interface{}{"name": name}}
	}

	v, ok := a.refValue(vpc, from("rt-vpc"))
	assert.True(t, ok)
	assert.Equal(t, "(x_k8s_get($client, 'ec2.services.k8s.aws/v1alpha1', 'VPC', $namespace, 'rt-vpc').status.vpcID)", v)

	// A list: sorted, since the IDs' order isn't known until the test runs.
	// A kind without a dependency gets its group's API version.
	subnets := refTarget{concrete: "subnetIDs", kind: "Subnet", group: "ec2.services.k8s.aws", path: "status.subnetID", list: true}
	v, ok = a.refValue(subnets, []interface{}{from("a"), from("b")})
	assert.True(t, ok)
	assert.Equal(t,
		"(sort([x_k8s_get($client, 'ec2.services.k8s.aws/v1alpha1', 'Subnet', $namespace, 'a').status.subnetID, "+
			"x_k8s_get($client, 'ec2.services.k8s.aws/v1alpha1', 'Subnet', $namespace, 'b').status.subnetID]))", v)
	v, ok = a.refValue(subnets, []interface{}{})
	assert.True(t, ok)
	assert.Equal(t, []interface{}{}, v)

	// Not by name in the test's namespace
	_, ok = a.refValue(vpc, map[string]interface{}{"from": map[string]interface{}{"name": "x", "namespace": "other"}})
	assert.False(t, ok)
	_, ok = a.refValue(vpc, "vpc-123")
	assert.False(t, ok)
}

func TestTemplateWithEscapes(t *testing.T) {
	assert.Equal(t,
		`(join('', ['a\\b \'c\' ', x]))`,
		templateWith(`a\b 'c' $X`, map[string]string{"$X": "x"}))
}

func TestCheckStep(t *testing.T) {
	g := &generator{
		in:         &Inputs{Values: []Value{{Name: "AMI_ID"}}},
		aws:        &awsAPI{idPath: "{.status.vpcID}"},
		crd:        &model.CRD{Kind: "VPC"},
		apiVersion: "ec2.services.k8s.aws/v1alpha1",
	}
	s := g.checkStep("check-created", []awsCheck{
		{Key: "describeVpcs", Command: "aws ec2 describe-vpcs --vpc-ids $AWS_ID --query 'Vpcs[0]'", Expected: map[string]interface{}{"name": "n-$SUFFIX"}},
		{Key: "enableDNSSupport", Command: "aws ec2 describe-vpc-attribute --vpc-id $AWS_ID", Expected: true},
	})
	assert.Equal(t, "check-created", s.Name)
	// One script for every AWS call, checked once by Chainsaw
	require.Len(t, s.Try, 1)
	sc := s.Try[0].Script
	env := map[string]interface{}{}
	for _, b := range sc.Env {
		env[b.Name] = b.Value
	}
	assert.Equal(t, "(x_k8s_get($client, 'ec2.services.k8s.aws/v1alpha1', 'VPC', $namespace, $name).status.vpcID)", env["AWS_ID"])
	assert.Equal(t, "(x_k8s_get($client, 'v1', 'ConfigMap', $namespace, 'ack-test-values').data.AMI_ID)", env["AMI_ID"])
	assert.Equal(t, literal(`export AWS_DEFAULT_OUTPUT=json AWS_RETRY_MODE=adaptive
# A call that fails prints its error, and leaves jq no JSON to read
jq -n \
  --argjson describeVpcs "$(aws ec2 describe-vpcs --vpc-ids $AWS_ID --query 'Vpcs[0]')" \
  --argjson enableDNSSupport "$(aws ec2 describe-vpc-attribute --vpc-id $AWS_ID)" \
  '$ARGS.named'`), sc.Content)
	assert.Equal(t, map[string]interface{}{
		"($error)": nil,
		"(json_parse($stdout))": map[string]interface{}{
			"describeVpcs":     map[string]interface{}{"name": "(join('', ['n-', $suffix]))"},
			"enableDNSSupport": true,
		},
	}, sc.Check)
	assert.Equal(t, "VPC", s.Catch[0].Describe.Kind)
}

func TestDeletedWaiter(t *testing.T) {
	describe := &awssdkmodel.Operation{ExportedName: "DescribeNatGateways"}
	m := &model.Model{SDKAPI: &model.SDKAPI{API: &awssdkmodel.API{Waiters: []awssdkmodel.Waiter{
		{Name: "NatGatewayAvailable", OperationName: "DescribeNatGateways"},
		{Name: "NatGatewayDeleted", OperationName: "DescribeNatGateways"},
		{Name: "InstanceTerminated", OperationName: "DescribeInstances"},
	}}}}
	nat := &model.CRD{Names: names.New("NatGateway")}
	assert.Equal(t, "aws ec2 wait nat-gateway-deleted --nat-gateway-ids $AWS_ID",
		deletedWaiter(m, nat, describe, "ec2", "NatGatewayIds"))
	instance := &model.CRD{Names: names.New("Instance")}
	assert.Equal(t, "aws ec2 wait instance-terminated --instance-ids $AWS_ID",
		deletedWaiter(m, instance, &awssdkmodel.Operation{ExportedName: "DescribeInstances"}, "ec2", "InstanceIds"))
	// No waiter named for the resource
	vpc := &model.CRD{Names: names.New("Vpc")}
	assert.Equal(t, "", deletedWaiter(m, vpc, &awssdkmodel.Operation{ExportedName: "DescribeVpcs"}, "ec2", "VpcIds"))
}

func TestResolveNested(t *testing.T) {
	igw := refTarget{concrete: "gatewayID", kind: "InternetGateway", group: "ec2.services.k8s.aws", path: "status.internetGatewayID"}
	group := refTarget{concrete: "groupID", kind: "SecurityGroup", group: "ec2.services.k8s.aws", path: "status.id"}
	a := &awsAPI{
		apiVersions: map[string]string{},
		nestedRefs: map[string][]nestedRef{
			"routes":       {{containers: []string{"routes"}, ref: "gatewayRef", target: igw}},
			"ingressRules": {{containers: []string{"ingressRules", "userIDGroupPairs"}, ref: "groupRef", target: group}},
		},
	}
	from := func(name string) map[string]interface{} {
		return map[string]interface{}{"from": map[string]interface{}{"name": name}}
	}
	get := func(kind, name, path string) string {
		return fmt.Sprintf("(x_k8s_get($client, 'ec2.services.k8s.aws/v1alpha1', '%s', $namespace, '%s').%s)", kind, name, path)
	}

	// A reference in each list item
	assert.Equal(t, []interface{}{
		map[string]interface{}{"destinationCIDRBlock": "0.0.0.0/0", "gatewayID": get("InternetGateway", "rt-igw", "status.internetGatewayID")},
		map[string]interface{}{"destinationCIDRBlock": "10.1.0.0/16", "gatewayID": "igw-literal"},
	}, a.resolveNested("routes", []interface{}{
		map[string]interface{}{"destinationCIDRBlock": "0.0.0.0/0", "gatewayRef": from("rt-igw")},
		map[string]interface{}{"destinationCIDRBlock": "10.1.0.0/16", "gatewayID": "igw-literal"},
	}))

	// Two levels down
	assert.Equal(t, []interface{}{
		map[string]interface{}{"fromPort": 443, "userIDGroupPairs": []interface{}{
			map[string]interface{}{"groupID": get("SecurityGroup", "other-sg", "status.id")},
		}},
	}, a.resolveNested("ingressRules", []interface{}{
		map[string]interface{}{"fromPort": 443, "userIDGroupPairs": []interface{}{
			map[string]interface{}{"groupRef": from("other-sg")},
		}},
	}))

	// Fields without nested references, and an empty list, are left as they are
	assert.Equal(t, "x", a.resolveNested("name", "x"))
	assert.Equal(t, []interface{}{}, a.resolveNested("routes", []interface{}{}))
}

func TestPathCheck(t *testing.T) {
	// A plain path
	expr, expected := pathCheck(Check{Path: "Attachments[0].VpcId"}, "(x_k8s_get(...))")
	assert.Equal(t, "Attachments[0].VpcId", expr)
	assert.Equal(t, "(x_k8s_get(...))", expected)

	// Filtered and sorted, the expected list sorted the same way
	expr, expected = pathCheck(Check{
		Path:    "Entries",
		Exclude: map[string]interface{}{"RuleNumber": 32767},
		SortBy:  stringList{"Egress", "RuleNumber"},
	}, []interface{}{
		map[string]interface{}{"Egress": true, "RuleNumber": 100},
		map[string]interface{}{"Egress": false, "RuleNumber": 200},
		map[string]interface{}{"Egress": false, "RuleNumber": 100},
	})
	assert.Equal(t, "sort_by(sort_by(Entries[?to_number(RuleNumber) != `32767`], &to_number(RuleNumber)), &to_string(Egress))", expr)
	assert.Equal(t, []interface{}{
		map[string]interface{}{"Egress": false, "RuleNumber": 100},
		map[string]interface{}{"Egress": false, "RuleNumber": 200},
		map[string]interface{}{"Egress": true, "RuleNumber": 100},
	}, expected)

	// Kept by a nested member, one member of each item, sorted
	expr, expected = pathCheck(Check{
		Path:      "CidrBlockAssociationSet",
		Where:     map[string]interface{}{"CidrBlockState.State": "associated"},
		Select:    "CidrBlock",
		Unordered: true,
	}, []interface{}{"10.2.0.0/16", "10.0.0.0/16"})
	assert.Equal(t, "sort(CidrBlockAssociationSet[?CidrBlockState.State == 'associated'][].CidrBlock)", expr)
	assert.Equal(t, []interface{}{"10.0.0.0/16", "10.2.0.0/16"}, expected)

	// A JSON document, compared by content
	expr, expected = pathCheck(Check{Path: "PolicyDocument"}, `{"Version":"2012-10-17"}`)
	assert.Equal(t, "json_parse(PolicyDocument)", expr)
	assert.Equal(t, map[string]interface{}{"Version": "2012-10-17"}, expected)
}

func TestWrapperEndsInList(t *testing.T) {
	instance := &awssdkmodel.Shape{Type: "structure"}
	instances := &awssdkmodel.Shape{Type: "list", MemberRef: awssdkmodel.ShapeRef{Shape: instance}}
	reservation := &awssdkmodel.Shape{Type: "structure", MemberRefs: map[string]*awssdkmodel.ShapeRef{"Instances": {Shape: instances}}}
	reservations := &awssdkmodel.Shape{Type: "list", MemberRef: awssdkmodel.ShapeRef{Shape: reservation}}
	output := &awssdkmodel.Shape{Type: "structure", MemberRefs: map[string]*awssdkmodel.ShapeRef{
		"Reservations": {Shape: reservations},
		"RouteTable":   {Shape: instance},
	}}
	assert.True(t, wrapperEndsInList(output, "Reservations.Instances"))
	assert.Equal(t, "Reservations[0].Instances", wrapperPath(output, "Reservations.Instances"))
	assert.False(t, wrapperEndsInList(output, "RouteTable"))
	assert.False(t, wrapperEndsInList(output, "Missing"))
}
