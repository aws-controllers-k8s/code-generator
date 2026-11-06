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
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/aws-controllers-k8s/pkg/names"
	"github.com/gertd/go-pluralize"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/code-generator/pkg/generate/code"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// awsAPI is how a resource is read and deleted through the AWS CLI, derived
// from its read and delete operations.
type awsAPI struct {
	// describe is the CLI call that reads the resource, e.g.
	// `aws ec2 describe-vpcs --vpc-ids $AWS_ID`
	describe string
	// root is the JMESPath to the resource in describe's output, e.g. Vpcs[0].
	// Empty if the output is the resource itself.
	root string
	// delete is the CLI call that deletes the resource, or empty if the
	// resource has no delete operation with a usable identifier.
	delete string
	// waitDeleted is the CLI call that waits for the resource to be deleted,
	// from a waiter in the API model, e.g.
	// `aws ec2 wait nat-gateway-deleted --nat-gateway-ids $AWS_ID`, or empty
	// if the model has none.
	waitDeleted string
	// notFoundCode is the error code describe fails with once the resource
	// is gone, or empty if the generator config doesn't declare one.
	notFoundCode string
	// idPath is the kubectl JSONPath to the resource's AWS ID in the CR.
	idPath string
	// idJSONName is the JSON name of the CR field holding the AWS ID.
	idJSONName string
	// members maps the JSON name of each spec field that describe's output
	// returns to the JMESPath of the output member holding it.
	members map[string]string
	// memberRefs are the output members' shapes, by JMESPath.
	memberRefs map[string]*awssdkmodel.ShapeRef
	// reads are the spec fields read back by a call of their own, e.g. VPC's
	// enableDNSSupport by DescribeVpcAttribute, by JSON name.
	reads map[string]fieldRead
	// why says why each spec field that isn't read back can't be, by JSON
	// name.
	why map[string]string
	// refs are the reference fields, e.g. RouteTable's vpcRef, by JSON name.
	// A reference is read back through the field it fills, e.g. vpcID, and
	// is expected to hold the referenced resource's AWS ID.
	refs map[string]refTarget
	// nestedRefs are the references inside spec fields' values, e.g.
	// RouteTable's routes[].gatewayRef, by the JSON name of the spec field
	// holding them.
	nestedRefs map[string][]nestedRef
	// apiVersions are the API versions of the kinds references point to, by
	// kind, where the inputs' dependencies give them.
	apiVersions map[string]string

	service string
	// readOp is the name of the resource's read operation.
	readOp string
	// idMember is the read operation's input member taking the AWS ID.
	idMember string
}

// refTarget is the field a reference field fills, and where the referenced
// resource keeps the AWS ID it fills it with.
type refTarget struct {
	// concrete is the JSON name of the field the reference fills, e.g. vpcID.
	concrete string
	// kind and group are the referenced resource's, e.g. VPC and
	// ec2.services.k8s.aws.
	kind, group string
	// path is the JMESPath to the AWS ID in the referenced CR, e.g.
	// status.vpcID.
	path string
	// list is whether the reference holds several resources, e.g. subnetRefs.
	list bool
}

// nestedRef is a reference inside a spec field's value, e.g. RouteTable's
// routes[].gatewayRef, which fills routes[].gatewayID.
type nestedRef struct {
	// containers are the JSON names of the fields holding it, outermost
	// first, e.g. routes, or ingressRules and userIDGroupPairs.
	containers []string
	// ref is the reference's JSON name, e.g. gatewayRef.
	ref    string
	target refTarget
}

// fieldRead reads a spec field back from AWS with a call of its own.
type fieldRead struct {
	// call is the AWS CLI call, without --query.
	call string
	// path is the JMESPath to the field in call's output, or empty if the
	// field is the whole output.
	path string
	// ref is the field's shape in call's output.
	ref *awssdkmodel.ShapeRef
}

func newAWSAPI(m *model.Model, crd *model.CRD, service string, notFoundCode string) (*awsAPI, error) {
	cfg := crd.Config()
	op := crd.Ops.ReadOne
	if op == nil {
		op = crd.Ops.ReadMany
	}
	if op == nil {
		return nil, fmt.Errorf("%s has no read operation", crd.Kind)
	}
	inputShape, err := crd.GetInputShape(op)
	if err != nil {
		return nil, err
	}

	// The CR field holding the AWS ID, and the read operation's input member
	// it goes in (possibly a list, e.g. DescribeVpcs's VpcIds)
	var crID, shapeID string
	// A field the generator config marks is_primary_key goes in the input
	// member of the same name, or its plural, e.g. ElasticIPAddress's
	// AllocationID in DescribeAddresses's AllocationIds
	if pk, err := crd.GetPrimaryKeyField(); err == nil && pk != nil {
		p := pluralize.NewClient()
		for _, name := range inputShape.MemberNames() {
			if strings.EqualFold(name, pk.Names.Original) || strings.EqualFold(name, p.Plural(pk.Names.Original)) {
				crID, shapeID = pk.Names.Original, name
				break
			}
		}
	}
	if shapeID == "" && op == crd.Ops.ReadMany {
		crID, shapeID = code.FindPluralizedIdentifiersInShape(crd, inputShape, op)
	}
	if shapeID == "" {
		if crID, shapeID, err = code.FindPrimaryIdentifierFieldNames(cfg, crd, op); err != nil {
			return nil, err
		}
	}
	if shapeID == "" {
		return nil, fmt.Errorf("%s: no identifier in %s's input", crd.Kind, op.ExportedName)
	}
	// The identifier can be the renamed one, e.g. Ids for DescribeSecurityGroups's
	// GroupIds, but the CLI takes the API's name
	shapeID = originalMemberName(crd, op, inputShape, shapeID)
	idField, idInStatus := findField(crd, crID)
	if idField == nil {
		return nil, fmt.Errorf("%s: no CR field for identifier %q", crd.Kind, crID)
	}
	a := &awsAPI{
		describe: fmt.Sprintf(
			"aws %s %s --%s $AWS_ID", service, cliName(op.ExportedName), cliName(shapeID),
		),
		idJSONName:  idField.GetCRDJSONFieldName(),
		members:     map[string]string{},
		memberRefs:  map[string]*awssdkmodel.ShapeRef{},
		reads:       map[string]fieldRead{},
		why:         map[string]string{},
		refs:        map[string]refTarget{},
		nestedRefs:  map[string][]nestedRef{},
		apiVersions: map[string]string{},
		service:     service,
		readOp:      op.ExportedName,
		idMember:    shapeID,
	}
	if idInStatus {
		a.idPath = fmt.Sprintf("{.status.%s}", a.idJSONName)
	} else {
		a.idPath = fmt.Sprintf("{.spec.%s}", a.idJSONName)
	}

	elemShape, root, err := resourceInOutput(crd, op)
	if err != nil {
		return nil, err
	}
	a.root = root
	// Output members are usually named as in the create operation, so its
	// renames apply too, e.g. SecurityGroup's GroupName is the Name field
	renameOps := []string{op.ExportedName}
	if crd.Ops.Create != nil {
		renameOps = append(renameOps, crd.Ops.Create.ExportedName)
	}
	for _, memberName := range elemShape.MemberNames() {
		for _, opName := range renameOps {
			fieldName := cfg.GetResourceFieldName(crd.Names.Original, opName, memberName)
			if f, ok := crd.SpecFields[fieldName]; ok {
				a.members[f.GetCRDJSONFieldName()] = memberName
				a.memberRefs[memberName] = elemShape.MemberRefs[memberName]
				break
			}
		}
	}
	a.mapOtherFields(m, crd, elemShape)

	a.notFoundCode = notFoundCode
	if resCfg, ok := cfg.Resources[crd.Names.Original]; ok && resCfg.Exceptions != nil {
		if e, ok := resCfg.Exceptions.Errors[404]; ok {
			a.notFoundCode = e.Code
		}
	}

	a.waitDeleted = deletedWaiter(m, crd, op, service, shapeID)

	// The delete operation takes the singular form of the read operation's
	// identifier, e.g. DeleteVpc's VpcId for DescribeVpcs's VpcIds
	if del := crd.Ops.Delete; del != nil {
		singular := pluralize.NewClient().Singular(shapeID)
		if delInput, err := crd.GetInputShape(del); err == nil {
			for _, memberName := range delInput.MemberNames() {
				if memberName == shapeID || memberName == singular {
					a.delete = fmt.Sprintf(
						"aws %s %s --%s $AWS_ID", service, cliName(del.ExportedName), cliName(memberName),
					)
					break
				}
			}
		}
	}
	return a, nil
}

// deletedWaiter returns the CLI call of the API model's waiter for the
// resource being deleted, e.g. NatGatewayDeleted or InstanceTerminated: a
// waiter on the read operation, named for the resource and ending in Deleted
// or Terminated. It's "" if there's none.
func deletedWaiter(m *model.Model, crd *model.CRD, op *awssdkmodel.Operation, service, idMember string) string {
	if m == nil || m.SDKAPI == nil || m.SDKAPI.API == nil {
		return ""
	}
	for _, w := range m.SDKAPI.API.Waiters {
		if w.OperationName != op.ExportedName {
			continue
		}
		subject := strings.TrimSuffix(strings.TrimSuffix(w.Name, "Deleted"), "Terminated")
		if subject == w.Name || !strings.EqualFold(subject, crd.Names.Original) {
			continue
		}
		return fmt.Sprintf("aws %s wait %s --%s $AWS_ID", service, cliName(w.Name), cliName(idMember))
	}
	return ""
}

// originalMemberName returns the input member of op that the possibly renamed
// name refers to.
func originalMemberName(
	crd *model.CRD,
	op *awssdkmodel.Operation,
	inputShape *awssdkmodel.Shape,
	name string,
) string {
	if _, ok := inputShape.MemberRefs[name]; ok {
		return name
	}
	for _, memberName := range inputShape.MemberNames() {
		if crd.Config().GetResourceFieldName(crd.Names.Original, op.ExportedName, memberName) == name {
			return memberName
		}
	}
	return name
}

// mapOtherFields works out how to read back the spec fields describe's output
// doesn't return under their own names: fields whose type comes from another
// operation (`from` in the generator config), and custom fields holding a list
// of a shape describe's output also holds. It records why for those it can't.
func (a *awsAPI) mapOtherFields(m *model.Model, crd *model.CRD, elemShape *awssdkmodel.Shape) {
	var custom []*model.Field
	for _, f := range crd.SpecFields {
		k := f.GetCRDJSONFieldName()
		if _, ok := a.members[k]; ok {
			continue
		}
		switch {
		case f.IsReference():
			a.why[k] = "a reference, only resolved when the test runs"
		case f.FieldConfig != nil && f.FieldConfig.From != nil:
			if why := a.mapFromField(m, f, elemShape); why != "" {
				a.why[k] = why
			}
		case f.FieldConfig != nil && f.FieldConfig.CustomField != nil:
			custom = append(custom, f)
		case f.FieldConfig != nil && f.FieldConfig.Type != nil:
			a.why[k] = "not in the API (a type in the generator config), so hooks handle it"
		default:
			a.why[k] = a.readOp + " doesn't return it"
		}
	}
	sort.Slice(custom, func(i, j int) bool { return custom[i].Names.Camel < custom[j].Names.Camel })
	a.mapCustomFields(custom, elemShape)
	a.mapReferences(crd)
}

// mapReferences reads each reference field back through the field it fills,
// e.g. RouteTable's vpcRef through vpcID, or records why it can't.
func (a *awsAPI) mapReferences(crd *model.CRD) {
	for _, c := range crd.SpecFields {
		if !c.HasReference() {
			continue
		}
		refNames, err := c.GetReferenceFieldName()
		if err != nil {
			continue
		}
		// Spec fields are keyed by their API names, e.g. VpcRef for VPCRef
		var ref *model.Field
		for key, f := range crd.SpecFields {
			if f.IsReference() && (strings.EqualFold(key, refNames.Camel) || f.Names.Camel == refNames.Camel) {
				ref = f
				break
			}
		}
		if ref == nil {
			continue
		}
		rk, ck := ref.GetCRDJSONFieldName(), c.GetCRDJSONFieldName()
		cfg := c.FieldConfig.References
		a.refs[rk] = refTarget{
			concrete: ck,
			kind:     cfg.Resource,
			group:    c.ReferencedServiceName() + ".services.k8s.aws",
			path:     refPath(cfg.Path),
			list:     strings.HasPrefix(c.GoType, "[]"),
		}
		if member, ok := a.members[ck]; ok {
			a.members[rk] = member
		} else if r, ok := a.reads[ck]; ok {
			a.reads[rk] = r
		} else {
			why := a.why[ck]
			if why == "" {
				why = a.readOp + " doesn't return it"
			}
			a.why[rk] = fmt.Sprintf("fills %s, which can't be read back: %s", ck, why)
			continue
		}
		delete(a.why, rk)
	}
	a.mapNestedReferences(crd)
}

// mapNestedReferences records the references inside spec fields' values,
// from the generator config's references on nested fields, e.g. RouteTable's
// Routes.GatewayId.
func (a *awsAPI) mapNestedReferences(crd *model.CRD) {
	for _, c := range crd.Fields {
		segments := strings.Split(c.Path, ".")
		if !c.HasReference() || len(segments) < 2 {
			continue
		}
		refNames, err := c.GetReferenceFieldName()
		if err != nil {
			continue
		}
		var containers []string
		for i := 1; i < len(segments); i++ {
			f, ok := crd.Fields[strings.Join(segments[:i], ".")]
			if !ok {
				containers = nil
				break
			}
			containers = append(containers, f.GetCRDJSONFieldName())
		}
		if len(containers) == 0 {
			continue
		}
		cfg := c.FieldConfig.References
		a.nestedRefs[containers[0]] = append(a.nestedRefs[containers[0]], nestedRef{
			containers: containers,
			ref:        refNames.CamelLower,
			target: refTarget{
				concrete: c.GetCRDJSONFieldName(),
				kind:     cfg.Resource,
				group:    c.ReferencedServiceName() + ".services.k8s.aws",
				path:     refPath(cfg.Path),
				list:     strings.HasPrefix(c.GoType, "[]"),
			},
		})
	}
	for k := range a.nestedRefs {
		refs := a.nestedRefs[k]
		sort.Slice(refs, func(i, j int) bool { return refs[i].ref < refs[j].ref })
	}
}

// resolveNested returns v, the value of the spec field k, with each nested
// reference to a CR by name replaced by the field it fills, holding the
// referenced resource's AWS ID, read when the test runs. For example,
// routes[].gatewayRef {from: {name: rt-igw}} becomes routes[].gatewayID
// (x_k8s_get(... 'rt-igw').status.internetGatewayID).
func (a *awsAPI) resolveNested(k string, v interface{}) interface{} {
	for _, n := range a.nestedRefs[k] {
		v = a.resolveAt(v, n, n.containers[1:])
	}
	return v
}

func (a *awsAPI) resolveAt(v interface{}, n nestedRef, rest []string) interface{} {
	switch t := v.(type) {
	case []interface{}:
		out := make([]interface{}, len(t))
		for i, e := range t {
			out[i] = a.resolveAt(e, n, rest)
		}
		return out
	case map[string]interface{}:
		out := make(map[string]interface{}, len(t))
		for mk, mv := range t {
			out[mk] = mv
		}
		if len(rest) != 0 {
			if c, ok := out[rest[0]]; ok {
				out[rest[0]] = a.resolveAt(c, n, rest[1:])
			}
			return out
		}
		if r, ok := out[n.ref]; ok {
			if rv, ok := a.refValue(n.target, r); ok {
				delete(out, n.ref)
				out[n.target.concrete] = rv
			}
		}
		return out
	}
	return v
}

// refPath returns the JMESPath to a field of a CR from the generator config's
// path to it, e.g. status.vpcID for Status.VPCID.
func refPath(path string) string {
	segments := strings.Split(path, ".")
	for i, seg := range segments {
		if seg == "ACKResourceMetadata" {
			segments[i] = "ackResourceMetadata"
		} else {
			segments[i] = names.New(seg).CamelLower
		}
	}
	return strings.Join(segments, ".")
}

// refGet returns the Chainsaw expression reading the AWS ID out of the
// referenced CR of the supplied name, in the test's namespace.
func (a *awsAPI) refGet(t refTarget, name string) string {
	apiVersion, ok := a.apiVersions[t.kind]
	if !ok {
		apiVersion = t.group + "/v1alpha1"
	}
	return fmt.Sprintf(
		"x_k8s_get($client, '%s', '%s', $namespace, '%s').%s", apiVersion, t.kind, name, t.path,
	)
}

// refNames returns the names of the CRs a reference field's value points to,
// or false if it doesn't point to CRs in the test's namespace by name.
func refNames(t refTarget, v interface{}) ([]string, bool) {
	one := func(w interface{}) (string, bool) {
		m, _ := w.(map[string]interface{})
		from, _ := m["from"].(map[string]interface{})
		name, _ := from["name"].(string)
		_, other := from["namespace"]
		return name, name != "" && !other
	}
	if !t.list {
		name, ok := one(v)
		return []string{name}, ok
	}
	l, ok := v.([]interface{})
	if !ok {
		return nil, false
	}
	out := []string{}
	for _, e := range l {
		name, ok := one(e)
		if !ok {
			return nil, false
		}
		out = append(out, name)
	}
	return out, true
}

// refValue returns what AWS holds for a reference field's value: a Chainsaw
// expression reading the referenced resource's AWS ID from its CR when the
// check runs, or for a list, the sorted list of their IDs. It returns false if
// the value doesn't point to CRs by name.
func (a *awsAPI) refValue(t refTarget, v interface{}) (interface{}, bool) {
	crs, ok := refNames(t, v)
	if !ok {
		return nil, false
	}
	if t.list && len(crs) == 0 {
		return []interface{}{}, true
	}
	gets := make([]string, len(crs))
	for i, name := range crs {
		gets[i] = a.refGet(t, name)
	}
	if t.list {
		return fmt.Sprintf("(sort([%s]))", strings.Join(gets, ", ")), true
	}
	return "(" + gets[0] + ")", true
}

// mapFromField works out how to read back a field whose type comes from
// another operation, or returns why it can't. A field set by a write
// operation is read back by the matching read operation, e.g. ModifyVpcAttribute's
// EnableDnsSupport by DescribeVpcAttribute.
func (a *awsAPI) mapFromField(m *model.Model, f *model.Field, elemShape *awssdkmodel.Shape) string {
	k := f.GetCRDJSONFieldName()
	from := f.FieldConfig.From
	ops := m.SDKAPI.API.Operations
	verb := leadingWordRegexp.FindString(from.Operation)
	var readName string
	switch verb {
	case "Describe", "Get", "List":
		readName = from.Operation
	case "Create", "Restore", "Import", "Copy":
		return fmt.Sprintf("an input of %s, which AWS doesn't return", from.Operation)
	default:
		if readName = pairedRead(ops, from.Operation[len(verb):]); readName == "" {
			return fmt.Sprintf("set by %s, which has no matching read operation", from.Operation)
		}
	}
	segments := strings.Split(from.Path, ".")
	notReturned := fmt.Sprintf("set by %s, but %s doesn't return %s", from.Operation, readName, from.Path)

	// Describe's output holds it
	if readName == a.readOp {
		path, ref := findPath(elemShape, segments)
		if ref == nil {
			return notReturned
		}
		a.members[k] = path
		a.memberRefs[path] = ref
		return ""
	}

	read, ok := ops[readName]
	if !ok {
		return notReturned
	}
	args, why := a.readArgs(read, segments[0])
	if why != "" {
		return fmt.Sprintf("set by %s, but %s %s", from.Operation, readName, why)
	}
	r := fieldRead{call: fmt.Sprintf("aws %s %s%s", a.service, cliName(readName), args), ref: &read.OutputRef}
	if from.Path != "." {
		path, ref := findPath(read.OutputRef.Shape, segments)
		switch {
		case ref != nil:
			r.path, r.ref = path, ref
		// The output can be the structure the write operation takes, e.g.
		// GetBucketCors returns the members of PutBucketCors's CORSConfiguration
		case !sharesMembers(f.ShapeRef, read.OutputRef.Shape):
			return notReturned
		}
	}
	a.reads[k] = r
	return ""
}

// readArgs returns the CLI arguments of a call to read the resource with the
// supplied operation, or why there aren't any. The AWS ID goes in the member
// named like describe's identifier, and a required enum gets the value named
// like the field, e.g. DescribeVpcAttribute's Attribute gets enableDnsSupport
// for the field at EnableDnsSupport.Value.
func (a *awsAPI) readArgs(read *awssdkmodel.Operation, fieldName string) (string, string) {
	input := read.InputRef.Shape
	if input == nil {
		return "", "takes no input"
	}
	p := pluralize.NewClient()
	var args strings.Builder
	hasID := false
	for _, name := range input.MemberNames() {
		ref := input.MemberRefs[name]
		isID := strings.EqualFold(name, a.idMember) ||
			strings.EqualFold(name, p.Singular(a.idMember)) ||
			strings.EqualFold(p.Plural(name), a.idMember)
		switch {
		case isID && !hasID:
			fmt.Fprintf(&args, " --%s $AWS_ID", cliName(name))
			hasID = true
		case !input.IsRequired(name):
		case ref.Shape != nil && ref.Shape.IsEnum():
			v := ""
			for _, e := range ref.Shape.Enum {
				if strings.EqualFold(e, fieldName) {
					v = e
				}
			}
			if v == "" {
				return "", fmt.Sprintf("needs a value for %s", name)
			}
			fmt.Fprintf(&args, " --%s %s", cliName(name), v)
		default:
			return "", fmt.Sprintf("needs a value for %s", name)
		}
	}
	if !hasID {
		return "", fmt.Sprintf("has no input for the AWS ID (%s)", a.idMember)
	}
	return args.String(), ""
}

// mapCustomFields matches custom fields holding a list of some structure to
// the members of describe's output holding a list of the same structure, e.g.
// SecurityGroup's IngressRules and EgressRules, both lists of IpPermission, to
// IpPermissions and IpPermissionsEgress. Where several members could hold a
// field, the one sharing the most words with the field's name wins, and a
// member another field took is ruled out.
func (a *awsAPI) mapCustomFields(fields []*model.Field, elemShape *awssdkmodel.Shape) {
	taken := map[string]bool{}
	for _, path := range a.members {
		taken[path] = true
	}
	candidates := map[string][]string{}
	for _, f := range fields {
		listOf := f.FieldConfig.CustomField.ListOf
		for _, name := range elemShape.MemberNames() {
			if listOf != "" && structureListOf(elemShape.MemberRefs[name]) == listOf {
				candidates[f.Names.Camel] = append(candidates[f.Names.Camel], name)
			}
		}
	}
	for progress := true; progress; {
		progress = false
		for _, f := range fields {
			k := f.GetCRDJSONFieldName()
			if _, done := a.members[k]; done {
				continue
			}
			if best := bestMember(f.Names.Camel, candidates[f.Names.Camel], taken); best != "" {
				a.members[k] = best
				a.memberRefs[best] = elemShape.MemberRefs[best]
				taken[best] = true
				progress = true
			}
		}
	}
	for _, f := range fields {
		k := f.GetCRDJSONFieldName()
		if _, done := a.members[k]; done {
			continue
		}
		cf := f.FieldConfig.CustomField
		switch {
		case cf.ListOf == "" || structureListOf(f.ShapeRef) == "":
			a.why[k] = "a custom field of plain values, with no shape to match in " + a.readOp
		case len(candidates[f.Names.Camel]) == 0:
			a.why[k] = fmt.Sprintf("a custom field of %s, which %s doesn't return", cf.ListOf, a.readOp)
		default:
			a.why[k] = fmt.Sprintf(
				"a custom field of %s, which several members of %s could hold (%s)",
				cf.ListOf, a.readOp, strings.Join(candidates[f.Names.Camel], ", "),
			)
		}
	}
}

// bestMember returns the one candidate not taken, or failing that the one
// sharing the most words with fieldName, or "" if there's no single best.
func bestMember(fieldName string, candidates []string, taken map[string]bool) string {
	var free []string
	for _, c := range candidates {
		if !taken[c] {
			free = append(free, c)
		}
	}
	if len(free) == 1 {
		return free[0]
	}
	best, bestScore, tie := "", 0, false
	for _, c := range free {
		score := 0
		for _, w := range camelWordRegexp.FindAllString(fieldName, -1) {
			if strings.Contains(c, w) {
				score++
			}
		}
		switch {
		case score > bestScore:
			best, bestScore, tie = c, score, false
		case score == bestScore:
			tie = true
		}
	}
	if tie || bestScore == 0 {
		return ""
	}
	return best
}

// pairedRead returns the operation that reads what a write operation on
// subject sets, e.g. DescribeVpcAttribute for ModifyVpcAttribute's
// VpcAttribute, or "" if there's none.
func pairedRead(ops map[string]*awssdkmodel.Operation, subject string) string {
	plural := pluralize.NewClient().Plural(subject)
	for _, name := range []string{
		"Get" + subject,
		"Describe" + subject,
		"Describe" + plural,
		"Get" + strings.Replace(subject, "Configuration", "Config", 1),
		"List" + plural,
	} {
		if _, ok := ops[name]; ok {
			return name
		}
	}
	return ""
}

// findPath returns the JMESPath to the member of shape at the dotted path
// segments, and its shape, trying ever shorter tails of the path, since a
// write operation can wrap the member, e.g. ModifyLaunchTemplate's
// LaunchTemplate.DefaultVersionNumber is a LaunchTemplate's
// DefaultVersionNumber. A list on the way is read at its first element.
func findPath(shape *awssdkmodel.Shape, segments []string) (string, *awssdkmodel.ShapeRef) {
	for start := range segments {
		var path []string
		cur := shape
		var ref *awssdkmodel.ShapeRef
		for i, seg := range segments[start:] {
			ref = nil
			if cur == nil {
				break
			}
			for _, name := range cur.MemberNames() {
				if strings.EqualFold(name, seg) {
					ref = cur.MemberRefs[name]
					path = append(path, name)
					break
				}
			}
			if ref == nil {
				break
			}
			cur = ref.Shape
			if cur != nil && cur.Type == "list" && start+i < len(segments)-1 {
				path[len(path)-1] += "[0]"
				cur = cur.MemberRef.Shape
			}
		}
		if ref != nil {
			return strings.Join(path, "."), ref
		}
	}
	return "", nil
}

// sharesMembers returns whether ref is a structure some of whose members
// shape also has.
func sharesMembers(ref *awssdkmodel.ShapeRef, shape *awssdkmodel.Shape) bool {
	if ref == nil || ref.Shape == nil || ref.Shape.Type != "structure" || shape == nil {
		return false
	}
	for name := range ref.Shape.MemberRefs {
		if _, ok := shape.MemberRefs[name]; ok {
			return true
		}
	}
	return false
}

// structureListOf returns the name of the structure ref is a list of, or ""
// if it isn't a list of structures.
func structureListOf(ref *awssdkmodel.ShapeRef) string {
	if ref == nil || ref.Shape == nil || ref.Shape.Type != "list" {
		return ""
	}
	elem := ref.Shape.MemberRef.Shape
	if elem == nil || elem.Type != "structure" {
		return ""
	}
	return elem.ShapeName
}

// findField returns the CR field with the supplied name, and whether it's in
// the status.
func findField(crd *model.CRD, name string) (*model.Field, bool) {
	camel := names.New(name).Camel
	if f, ok := crd.StatusFields[camel]; ok {
		return f, true
	}
	if f, ok := crd.SpecFields[camel]; ok {
		return f, false
	}
	// The identifier can be an original API name whose CR field name differs
	// in case, e.g. InternetGatewayId for InternetGatewayID
	for key, f := range crd.StatusFields {
		if strings.EqualFold(key, camel) || strings.EqualFold(f.Names.Original, name) {
			return f, true
		}
	}
	for key, f := range crd.SpecFields {
		if strings.EqualFold(key, camel) || strings.EqualFold(f.Names.Original, name) {
			return f, false
		}
	}
	return nil, false
}

// resourceInOutput returns the shape of the resource in the read operation's
// output, and the JMESPath to it. A list operation's output holds a list of
// resources, of which the read returns just one.
func resourceInOutput(crd *model.CRD, op *awssdkmodel.Operation) (*awssdkmodel.Shape, string, error) {
	outputShape, err := crd.GetOutputShape(op)
	if err != nil {
		return nil, "", err
	}
	wrapper := crd.GetOutputWrapperFieldPath(op)
	if op != crd.Ops.ReadMany {
		if wrapper != nil {
			return outputShape, wrapperPath(op.OutputRef.Shape, *wrapper), nil
		}
		return outputShape, "", nil
	}
	if wrapper != nil && outputShape.Type == "list" {
		return outputShape.MemberRef.Shape, wrapperPath(op.OutputRef.Shape, *wrapper) + "[0]", nil
	}
	// The wrapper can end in the list, which the model already reads the
	// resource from, e.g. DescribeInstances's Reservations.Instances
	if wrapper != nil && outputShape.Type == "structure" && wrapperEndsInList(op.OutputRef.Shape, *wrapper) {
		return outputShape, wrapperPath(op.OutputRef.Shape, *wrapper) + "[0]", nil
	}
	for _, memberName := range outputShape.MemberNames() {
		ref := outputShape.MemberRefs[memberName]
		if ref.Shape.Type == "list" && ref.Shape.MemberRef.Shape.Type == "structure" {
			return ref.Shape.MemberRef.Shape, memberName + "[0]", nil
		}
	}
	return nil, "", fmt.Errorf("%s: %s's output has no list of resources", crd.Kind, op.ExportedName)
}

// readable returns whether a spec field can be read back from AWS, by describe,
// a call of its own, or a check the inputs declare. A reference is read back
// through the field it fills.
func (a *awsAPI) readable(k string, declared map[string]Check) bool {
	_, member := a.members[k]
	_, read := a.reads[k]
	_, check := declared[k]
	return member || read || check
}

// checks returns the AWS CLI calls reading back the fields of spec, with the
// JSON each must print, the fields they check, and the fields they can't.
// Fields describe's output returns share one call.
func (a *awsAPI) checks(
	crd *model.CRD,
	spec map[string]interface{},
	declared map[string]Check,
) (checks []awsCheck, checked, unchecked []string) {
	described := map[string]interface{}{}
	// Declared checks reading a path in the resource, as the read operation
	// returns it, by the expression Chainsaw evaluates
	resource := map[string]interface{}{}
	fields := map[string]*model.Field{}
	for _, f := range crd.SpecFields {
		fields[f.GetCRDJSONFieldName()] = f
	}
	for _, k := range sortedKeys(spec) {
		v := a.resolveNested(k, spec[k])
		t, isRef := a.refs[k]
		if isRef {
			// The ID of the resource it points to, known when the test runs
			rv, ok := a.refValue(t, v)
			if !ok {
				unchecked = append(unchecked, k)
				continue
			}
			v = rv
		}
		if c, ok := declared[k]; ok {
			checked = append(checked, k)
			if c.When != nil && !reflect.DeepEqual(c.When, spec[k]) {
				continue
			}
			if c.Path != "" {
				var ref *awssdkmodel.ShapeRef
				if f, ok := fields[k]; ok && !isRef {
					ref = f.ShapeRef
				}
				expr, expected := pathCheck(c, toAWS(v, ref))
				// An expression, which Chainsaw evaluates over the resource
				resource["("+expr+")"] = expected
				continue
			}
			expected := v
			if c.Unordered && !isRef {
				expected = sortScalars(expected)
			}
			checks = append(checks, awsCheck{Key: k, Command: c.Command, Expected: expected})
			continue
		}
		if r, ok := a.reads[k]; ok {
			query := strings.TrimPrefix(r.path+projection(v, r.ref), ".")
			if isRef && t.list {
				query = fmt.Sprintf("sort(%s)", query)
			}
			command := r.call
			if query != "" {
				command += fmt.Sprintf(" --query '%s'", query)
			}
			checks = append(checks, awsCheck{Key: k, Command: command, Expected: toAWS(v, r.ref)})
			checked = append(checked, k)
			continue
		}
		if _, ok := a.members[k]; ok {
			described[k] = v
			checked = append(checked, k)
			continue
		}
		unchecked = append(unchecked, k)
	}
	if len(resource) != 0 {
		command := a.describe
		if a.root != "" {
			command += fmt.Sprintf(" --query '%s'", a.root)
		}
		checks = append(checks, awsCheck{Key: "resource", Command: command, Expected: resource})
	}
	if query, expected := a.query(crd, described); query != "" {
		checks = append([]awsCheck{{
			Key:      names.New(a.readOp).CamelLower,
			Command:  fmt.Sprintf("%s --query '%s'", a.describe, query),
			Expected: expected,
		}}, checks...)
	}
	return checks, checked, unchecked
}

// pathCheck returns the JMESPath expression a declared check with a path reads
// from the resource, and the value it must hold: expected, in AWS's names,
// sorted as the expression sorts it. A JSON document, e.g. a policy, is
// compared by content.
func pathCheck(c Check, expected interface{}) (string, interface{}) {
	expr := c.Path
	var conds []string
	for _, k := range sortedKeys(c.Where) {
		conds = append(conds, pathCond(k, "==", c.Where[k]))
	}
	for _, k := range sortedKeys(c.Exclude) {
		conds = append(conds, pathCond(k, "!=", c.Exclude[k]))
	}
	if len(conds) != 0 {
		expr = fmt.Sprintf("%s[?%s]", expr, strings.Join(conds, " && "))
	}
	// sort_by is stable, so the first key is sorted on last
	for i := len(c.SortBy) - 1; i >= 0; i-- {
		expr = fmt.Sprintf("sort_by(%s, &%s)", expr, sortKey(c.SortBy[i], expected))
	}
	if len(c.SortBy) != 0 {
		expected = sortByMembers(expected, c.SortBy)
	}
	if c.Select != "" {
		expr += "[]." + c.Select
	}
	if c.Unordered {
		expr = fmt.Sprintf("sort(%s)", expr)
		expected = sortScalars(expected)
	}
	if s, ok := expected.(string); ok {
		t := strings.TrimSpace(s)
		var parsed interface{}
		if (strings.HasPrefix(t, "{") || strings.HasPrefix(t, "[")) && json.Unmarshal([]byte(t), &parsed) == nil {
			return fmt.Sprintf("json_parse(%s)", expr), parsed
		}
	}
	return expr, expected
}

// pathCond returns a JMESPath filter condition comparing an item's member to
// a value. Numbers that json_parse returns only compare to JMESPath numbers
// through to_number.
func pathCond(member, op string, value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%s %s '%s'", member, op, strings.ReplaceAll(v, "'", `\'`))
	case int, int64, float64:
		return fmt.Sprintf("to_number(%s) %s `%v`", member, op, v)
	default:
		b, _ := json.Marshal(v)
		return fmt.Sprintf("%s %s `%s`", member, op, b)
	}
}

// sortKey returns the JMESPath sort_by key for a member, by the type of its
// values in expected: sort_by only takes numbers and strings, and only sorts
// json_parse's numbers through to_number.
func sortKey(member string, expected interface{}) string {
	l, _ := expected.([]interface{})
	for _, e := range l {
		switch memberValue(e, member).(type) {
		case int, int64, float64:
			return fmt.Sprintf("to_number(%s)", member)
		case bool:
			return fmt.Sprintf("to_string(%s)", member)
		case string:
			return member
		}
	}
	return member
}

// sortByMembers returns a list of structures sorted by the supplied members,
// the first first, as JMESPath's sort_by sorts them.
func sortByMembers(v interface{}, members []string) interface{} {
	l, ok := v.([]interface{})
	if !ok {
		return v
	}
	sorted := append([]interface{}{}, l...)
	sort.SliceStable(sorted, func(i, j int) bool {
		for _, m := range members {
			a, b := memberValue(sorted[i], m), memberValue(sorted[j], m)
			if less, more := compareValues(a, b), compareValues(b, a); less || more {
				return less
			}
		}
		return false
	})
	return sorted
}

func compareValues(a, b interface{}) bool {
	switch x := a.(type) {
	case bool:
		y, _ := b.(bool)
		return fmt.Sprint(x) < fmt.Sprint(y)
	case string:
		y, _ := b.(string)
		return x < y
	}
	x, xok := toFloat(a)
	y, yok := toFloat(b)
	return xok && yok && x < y
}

func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case float64:
		return n, true
	}
	return 0, false
}

// memberValue returns the value at a dotted member path in a structure.
func memberValue(v interface{}, path string) interface{} {
	for _, seg := range strings.Split(path, ".") {
		m, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}
		v = m[seg]
	}
	return v
}

// wrapperEndsInList returns whether the member at an output wrapper field path
// is a list.
func wrapperEndsInList(output *awssdkmodel.Shape, wrapper string) bool {
	cur := output
	for _, seg := range strings.Split(wrapper, ".") {
		if cur == nil {
			return false
		}
		if cur.Type == "list" {
			cur = cur.MemberRef.Shape
		}
		ref, ok := cur.MemberRefs[seg]
		if !ok || ref.Shape == nil {
			return false
		}
		cur = ref.Shape
	}
	return cur != nil && cur.Type == "list"
}

// wrapperPath returns the JMESPath to the member at an output wrapper field
// path, reading each list before the last member at its first element, e.g.
// Reservations[0].Instances for DescribeInstances's Reservations.Instances.
func wrapperPath(output *awssdkmodel.Shape, wrapper string) string {
	segments := strings.Split(wrapper, ".")
	cur := output
	for i, seg := range segments {
		if cur == nil {
			break
		}
		ref, ok := cur.MemberRefs[seg]
		if !ok || ref.Shape == nil {
			break
		}
		cur = ref.Shape
		if cur.Type == "list" && i < len(segments)-1 {
			segments[i] += "[0]"
			cur = cur.MemberRef.Shape
		}
	}
	return strings.Join(segments, ".")
}

// query returns a JMESPath query selecting the fields of spec from describe's
// output, which must return them all, and the JSON it must return.
func (a *awsAPI) query(
	crd *model.CRD,
	spec map[string]interface{},
) (query string, expected map[string]interface{}) {
	tagJSON, tagKey := "", crd.GetTagKeyMemberName()
	if f, err := crd.GetTagField(); err == nil && f != nil {
		tagJSON = f.GetCRDJSONFieldName()
	}

	expected = map[string]interface{}{}
	var selects []string
	for _, k := range sortedKeys(spec) {
		member := a.members[k]
		ref := a.memberRefs[member]
		value := toAWS(spec[k], ref)
		if k == tagJSON && ref.Shape.Type == "list" {
			// AWS also returns ACK's system tags, which the spec doesn't hold
			selects = append(selects, fmt.Sprintf(
				"%s: sort_by(%s[?!starts_with(%s, `\"services.k8s.aws/\"`)], &%s)",
				k, member, tagKey, tagKey,
			), fmt.Sprintf(
				"systemTags: sort(%s[?starts_with(%s, `\"services.k8s.aws/\"`)].%s)",
				member, tagKey, tagKey,
			))
			expected[k] = sortByKey(value, tagKey)
			expected["systemTags"] = []interface{}{
				"services.k8s.aws/controller-version", "services.k8s.aws/namespace",
			}
			continue
		}
		if t, ok := a.refs[k]; ok && t.list {
			// The IDs' order isn't known until the test runs, so both sides
			// are sorted
			selects = append(selects, fmt.Sprintf("%s: sort(%s)", k, member))
			expected[k] = value
			continue
		}
		selects = append(selects, fmt.Sprintf("%s: %s%s", k, member, projection(spec[k], ref)))
		expected[k] = value
	}
	if len(selects) == 0 {
		return "", nil
	}
	sel := "{" + strings.Join(selects, ", ") + "}"
	if a.root != "" {
		sel = a.root + "." + sel
	}
	return sel, expected
}

// projection returns the JMESPath that narrows a member of describe's output
// to the parts the CR value sets, since AWS returns more than the spec holds,
// e.g. `[].{Key: Key, Value: Value}` for a list of structures.
func projection(v interface{}, ref *awssdkmodel.ShapeRef) string {
	if ref == nil || ref.Shape == nil {
		return ""
	}
	switch ref.Shape.Type {
	case "structure":
		m, ok := v.(map[string]interface{})
		// An empty list's elements set nothing, and {} isn't valid JMESPath
		if !ok || len(m) == 0 {
			return ""
		}
		var selects []string
		for _, k := range sortedKeys(m) {
			member := memberName(ref.Shape, k)
			selects = append(selects, fmt.Sprintf(
				"%s: %s%s", member, member, projection(m[k], ref.Shape.MemberRefs[member]),
			))
		}
		return ".{" + strings.Join(selects, ", ") + "}"
	case "list":
		l, ok := v.([]interface{})
		if !ok || ref.Shape.MemberRef.Shape == nil || ref.Shape.MemberRef.Shape.Type != "structure" {
			return ""
		}
		// Every key any element sets
		union := map[string]interface{}{}
		for _, e := range l {
			if m, ok := e.(map[string]interface{}); ok {
				for k, ev := range m {
					union[k] = ev
				}
			}
		}
		if len(union) == 0 {
			return ""
		}
		return "[]" + projection(union, &ref.Shape.MemberRef)
	}
	return ""
}

// toAWS converts a CR value to the JSON AWS returns for it: the same values
// under the API's member names, e.g. {key: a} becomes {Key: a}.
func toAWS(v interface{}, ref *awssdkmodel.ShapeRef) interface{} {
	if ref == nil || ref.Shape == nil {
		return v
	}
	switch ref.Shape.Type {
	case "structure":
		m, ok := v.(map[string]interface{})
		if !ok {
			return v
		}
		out := map[string]interface{}{}
		for k, mv := range m {
			member := memberName(ref.Shape, k)
			out[member] = toAWS(mv, ref.Shape.MemberRefs[member])
		}
		return out
	case "list":
		l, ok := v.([]interface{})
		if !ok {
			return v
		}
		out := make([]interface{}, len(l))
		for i, e := range l {
			out[i] = toAWS(e, &ref.Shape.MemberRef)
		}
		return out
	case "map":
		m, ok := v.(map[string]interface{})
		if !ok {
			return v
		}
		out := map[string]interface{}{}
		for k, mv := range m {
			out[k] = toAWS(mv, &ref.Shape.ValueRef)
		}
		return out
	}
	return v
}

// memberName returns the member of a structure that a CR JSON key holds.
func memberName(shape *awssdkmodel.Shape, jsonKey string) string {
	for _, m := range shape.MemberNames() {
		if names.New(m).CamelLower == jsonKey {
			return m
		}
	}
	for _, m := range shape.MemberNames() {
		if strings.EqualFold(m, jsonKey) {
			return m
		}
	}
	return jsonKey
}

func sortByKey(v interface{}, key string) interface{} {
	l, ok := v.([]interface{})
	if !ok {
		return v
	}
	sorted := append([]interface{}{}, l...)
	sort.SliceStable(sorted, func(i, j int) bool {
		mi, _ := sorted[i].(map[string]interface{})
		mj, _ := sorted[j].(map[string]interface{})
		return fmt.Sprint(mi[key]) < fmt.Sprint(mj[key])
	})
	return sorted
}

// sortScalars returns a list of strings or numbers sorted, and anything else
// as is.
func sortScalars(v interface{}) interface{} {
	l, ok := v.([]interface{})
	if !ok {
		return v
	}
	sorted := append([]interface{}{}, l...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return fmt.Sprint(sorted[i]) < fmt.Sprint(sorted[j])
	})
	return sorted
}

var (
	cliWordRegexp  = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	cliUpperRegexp = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	// The first word of an operation name, e.g. Modify in ModifyVpcAttribute
	leadingWordRegexp = regexp.MustCompile(`^[A-Z][a-z]+`)
	camelWordRegexp   = regexp.MustCompile(`[A-Z][a-z0-9]*`)
)

// cliName returns the AWS CLI name of an API operation or parameter, e.g.
// describe-vpcs for DescribeVpcs.
func cliName(s string) string {
	s = cliWordRegexp.ReplaceAllString(s, "${1}-${2}")
	s = cliUpperRegexp.ReplaceAllString(s, "${1}-${2}")
	return strings.ToLower(s)
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
