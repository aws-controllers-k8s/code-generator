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

package fieldpath

import (
	"encoding/json"
	"strings"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
)

// Path provides a JSONPath-like struct and field-member "route" to a
// particular field within a resource. Path implements json.Marshaler
// interface.
type Path struct {
	parts []string
}

// String returns the dotted-notation representation of the Path
func (p *Path) String() string {
	return strings.Join(p.parts, ".")
}

// MarshalJSON returns the JSON encoding of a Path object.
func (p *Path) MarshalJSON() ([]byte, error) {
	// Since json.Marshal doesn't encode unexported struct fields we have to
	// copy the Path instance into a new struct object with exported fields.
	// See https://github.com/aws-controllers-k8s/community/issues/772
	return json.Marshal(
		struct {
			Parts []string
		}{
			p.parts,
		},
	)
}

// Pop removes the last part from the Path and returns it.
func (p *Path) Pop() (part string) {
	if len(p.parts) > 0 {
		part = p.parts[len(p.parts)-1]
		p.parts = p.parts[:len(p.parts)-1]
	}
	return part
}

// At returns the part of the Path at the supplied index, or empty string if
// index exceeds boundary.
func (p *Path) At(index int) string {
	if index < 0 || len(p.parts) == 0 || index > len(p.parts)-1 {
		return ""
	}
	return p.parts[index]
}

// Front returns the first part of the Path or empty string if the Path has no
// parts.
func (p *Path) Front() string {
	if len(p.parts) == 0 {
		return ""
	}
	return p.parts[0]
}

// PopFront removes the first part of the Path and returns it.
func (p *Path) PopFront() (part string) {
	if len(p.parts) > 0 {
		part = p.parts[0]
		p.parts = p.parts[1:]
	}
	return part
}

// Back returns the last part of the Path or empty string if the Path has no
// parts.
func (p *Path) Back() string {
	if len(p.parts) == 0 {
		return ""
	}
	return p.parts[len(p.parts)-1]
}

// PushBack adds a new part to the end of the Path.
func (p *Path) PushBack(part string) {
	p.parts = append(p.parts, part)
}

// Copy returns a new Path that is a copy of this Path
func (p *Path) Copy() *Path {
	return &Path{p.parts}
}

// CopyAt returns a new Path that is a copy of this Path up to the supplied
// index.
//
// e.g. given Path $A containing "X.Y", $A.CopyAt(0) would return a new Path
// containing just "X". $A.CopyAt(1) would return a new Path containing "X.Y".
func (p *Path) CopyAt(index int) *Path {
	if index < 0 || len(p.parts) == 0 || index > len(p.parts)-1 {
		return nil
	}
	return &Path{p.parts[0 : index+1]}
}

// Empty returns true if there are no parts to the Path
func (p *Path) Empty() bool {
	return len(p.parts) == 0
}

// Size returns the Path number of parts
func (p *Path) Size() int {
	return len(p.parts)
}

// ShapeRef returns an aws-sdk-go ShapeRef within the supplied ShapeRef that
// matches the Path. Returns nil if no matching ShapeRef could be found.
//
// Assume a ShapeRef that looks like this:
//
//	authShapeRef := &awssdkmodel.ShapeRef{
//	  ShapeName: "Author",
//	  Shape: &awssdkmodel.Shape{
//	    Type: "structure",
//	    MemberRefs: map[string]*awssdkmodel.ShapeRef{
//	      "Name": &awssdkmodel.ShapeRef{
//	        ShapeName: "Name",
//	        Shape: &awssdkmodel.Shape{
//	          Type: "string",
//	        },
//	      },
//	      "Address": &awssdkmodel.ShapeRef{
//	        ShapeName: "Address",
//	        Shape: &awssdkmodel.Shape{
//	          Type: "structure",
//	          MemberRefs: map[string]*awssdkmodel.ShapeRef{
//	            "State": &awssdkmodel.ShapeRef{
//	              ShapeName: "StateCode",
//	              Shape: &awssdkmodel.Shape{
//	                Type: "string",
//	              },
//	            },
//	            "Country": &awssdkmodel.ShapeRef{
//	              ShapeName: "CountryCode",
//	              Shape: &awssdkmodel.Shape{
//	                Type: "string",
//	              },
//	            },
//	          },
//	        },
//	      },
//	    },
//	  },
//	}
//
// If I have the following Path:
//
// p := fieldpath.FromString("Author.Address.Country")
//
// calling p.ShapeRef(authShapeRef) would return the following:
//
//	&awssdkmodel.ShapeRef{
//	  ShapeName: "CountryCode",
//	  Shape: &awssdkmodel.Shape{
//	    Type: "string",
//	  },
//	},
func (p *Path) ShapeRef(
	subject *awssdkmodel.ShapeRef,
) *awssdkmodel.ShapeRef {
	if subject == nil || p == nil || len(p.parts) == 0 {
		return nil
	}

	// We first check that the first part in the path matches the supplied
	// subject shape's name.
	var compare *awssdkmodel.ShapeRef = subject
	cp := p.Copy()
	cur := cp.PopFront()
	if compare.ShapeName != cur {
		return nil
	}
	// And then we walk through the path, searching through the supplied
	// ShapeRef for a member ShapeRef matching each path element.
	for !cp.Empty() {
		cur = cp.PopFront()
		if compare = memberShapeRef(compare, cur); compare == nil {
			return nil
		}
	}
	return compare
}

// ShapeRefAt returns an aws-sdk-go ShapeRef within the supplied ShapeRef that
// matches the Path at the supplied index. Returns nil if no matching ShapeRef
// could be found or index out of bounds.
func (p *Path) ShapeRefAt(
	subject *awssdkmodel.ShapeRef,
	index int,
) *awssdkmodel.ShapeRef {
	if subject == nil || p == nil || len(p.parts) == 0 {
		return nil
	}

	cp := p.CopyAt(index)
	if cp == nil {
		return nil
	}
	return cp.ShapeRef(subject)
}

// IterShapeRefs returns a slice of ShapeRef pointers representing each part of
// the path
func (p *Path) IterShapeRefs(
	subject *awssdkmodel.ShapeRef,
) []*awssdkmodel.ShapeRef {
	res := make([]*awssdkmodel.ShapeRef, len(p.parts))
	for idx, _ := range p.parts {
		res[idx] = p.ShapeRefAt(subject, idx)
	}
	return res
}

// memberShapeRef returns the named member ShapeRef of the supplied
// ShapeRef
func memberShapeRef(
	shapeRef *awssdkmodel.ShapeRef,
	memberName string,
) *awssdkmodel.ShapeRef {
	if shapeRef.ShapeName == memberName {
		return shapeRef
	}
	switch shapeRef.Shape.Type {
	case "structure":
		// We are looking for a member of a structure. Since the ACK fields and
		// the AWS SDK fields may have different casing (e.g AWSVPCConfiguration
		// and AwsVpcConfiguration) we need to perform a case insensitive
		// comparison to find the correct member reference.
		for memberRefName, memberRefShape := range shapeRef.Shape.MemberRefs {
			if strings.EqualFold(memberRefName, memberName) {
				return memberRefShape
			}
		}
		// If no matching member is found, return nil.
		return nil
	case "list":
		return memberShapeRef(&shapeRef.Shape.MemberRef, memberName)
	case "map":
		return memberShapeRef(&shapeRef.Shape.ValueRef, memberName)
	}
	return nil
}

// HasPrefix returns true if the supplied string, delimited on ".", matches
// p.parts up to the length of the supplied string.
// e.g. if the Path p represents "A.B":
//
//	subject "A" -> true
//	subject "A.B" -> true
//	subject "A.B.C" -> false
//	subject "B" -> false
//	subject "A.C" -> false
func (p *Path) HasPrefix(subject string) bool {
	subjectSplit := strings.Split(subject, ".")

	if len(subjectSplit) > len(p.parts) {
		return false
	}

	for i, s := range subjectSplit {
		if p.parts[i] != s {
			return false
		}
	}

	return true
}

// HasPrefixFold is the same as HasPrefix but uses case-insensitive comparisons
func (p *Path) HasPrefixFold(subject string) bool {
	subjectSplit := strings.Split(subject, ".")

	if len(subjectSplit) > len(p.parts) {
		return false
	}

	for i, s := range subjectSplit {
		if !strings.EqualFold(p.parts[i], s) {
			return false
		}
	}

	return true
}

// FromString returns a new Path from a dotted-notation string, e.g.
// "Author.Name".
func FromString(dotted string) *Path {
	return &Path{strings.Split(dotted, ".")}
}
