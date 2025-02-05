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

package fieldpath_test

import (
	"testing"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/stretchr/testify/require"

	"github.com/aws-controllers-k8s/code-generator/pkg/fieldpath"
)

func TestBasics(t *testing.T) {
	require := require.New(t)

	pname := fieldpath.FromString("Author.Name")
	require.Equal("Author.Name", pname.String())

	pstate := fieldpath.FromString("Author.Address.State")
	require.Equal("Author.Address.State", pstate.String())

	require.Equal("Author", pstate.Front())
	require.Equal("State", pstate.Back())

	require.Equal("Author", pstate.At(0))
	require.Equal("Address", pstate.At(1))
	require.Equal("State", pstate.At(2))
	require.Equal("", pstate.At(3))

	pauth := pstate.CopyAt(0)
	require.Equal("Author", pauth.String())

	last := pstate.Pop()
	require.Equal("State", last)
	require.Equal("Address", pstate.Back())

	pstate.PushBack("Country")
	require.Equal("Country", pstate.Back())

	front := pstate.PopFront()
	require.Equal("Author", front)
	require.Equal("Address", pstate.Front())
	require.False(pstate.Empty())
	pstate.Pop()
	require.False(pstate.Empty())
	pstate.Pop()
	require.True(pstate.Empty())
}

func TestHasPrefix(t *testing.T) {
	require := require.New(t)

	p := fieldpath.FromString("Author.Name")
	require.True(p.HasPrefix("Author.Name"))
	require.True(p.HasPrefix("Author"))
	require.False(p.HasPrefix("Name"))
	require.False(p.HasPrefix("Author.Address"))
	// Case-insensitive comparisons...
	require.False(p.HasPrefix("author"))
	require.True(p.HasPrefixFold("author"))
}

func TestShapeRef(t *testing.T) {
	require := require.New(t)

	p := fieldpath.FromString("Author.Name")
	emptyShapeRef := &awssdkmodel.ShapeRef{}
	require.Nil(p.ShapeRef(emptyShapeRef))

	authShapeRef := &awssdkmodel.ShapeRef{
		ShapeName: "Author",
		Shape: &awssdkmodel.Shape{
			Type: "structure",
			MemberRefs: map[string]*awssdkmodel.ShapeRef{
				"Name": &awssdkmodel.ShapeRef{
					ShapeName: "Name",
					Shape: &awssdkmodel.Shape{
						Type: "string",
					},
				},
				"Address": &awssdkmodel.ShapeRef{
					ShapeName: "Address",
					Shape: &awssdkmodel.Shape{
						Type: "structure",
						MemberRefs: map[string]*awssdkmodel.ShapeRef{
							"State": &awssdkmodel.ShapeRef{
								ShapeName: "StateCode",
								Shape: &awssdkmodel.Shape{
									Type: "string",
								},
							},
							"Country": &awssdkmodel.ShapeRef{
								ShapeName: "CountryCode",
								Shape: &awssdkmodel.Shape{
									Type: "string",
								},
							},
						},
					},
				},
				"Books": &awssdkmodel.ShapeRef{
					ShapeName: "BookList",
					Shape: &awssdkmodel.Shape{
						Type: "list",
						MemberRef: awssdkmodel.ShapeRef{
							ShapeName: "Book",
							Shape: &awssdkmodel.Shape{
								Type: "structure",
								MemberRefs: map[string]*awssdkmodel.ShapeRef{
									"Title": &awssdkmodel.ShapeRef{
										ShapeName: "Title",
										Shape: &awssdkmodel.Shape{
											Type: "string",
										},
									},
									"ChapterPageCounts": &awssdkmodel.ShapeRef{
										ShapeName: "ChapterPageCounts",
										Shape: &awssdkmodel.Shape{
											Type: "map",
											KeyRef: awssdkmodel.ShapeRef{
												ShapeName: "ChapterTitle",
												Shape: &awssdkmodel.Shape{
													Type: "string",
												},
											},
											ValueRef: awssdkmodel.ShapeRef{
												ShapeName: "PageCount",
												Shape: &awssdkmodel.Shape{
													Type: "integer",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				"WeirdlycasEdType": &awssdkmodel.ShapeRef{
					ShapeName: "WeirdlycasEdType",
					Shape: &awssdkmodel.Shape{
						Type: "string",
					},
				},
			},
		},
	}
	ref := p.ShapeRef(authShapeRef)
	require.NotNil(ref)
	require.Equal("Name", ref.ShapeName)
	require.Equal("string", ref.Shape.Type)

	p = fieldpath.FromString("Author")
	ref = p.ShapeRefAt(authShapeRef, 0)
	require.NotNil(ref)
	require.Equal("Author", ref.ShapeName)
	ref = p.ShapeRefAt(authShapeRef, 1)
	require.Nil(ref)

	p = fieldpath.FromString("Author.Address")
	ref = p.ShapeRefAt(authShapeRef, 0)
	require.NotNil(ref)
	require.Equal("Author", ref.ShapeName)
	ref = p.ShapeRefAt(authShapeRef, 1)
	require.NotNil(ref)
	require.Equal("Address", ref.ShapeName)
	ref = p.ShapeRefAt(authShapeRef, 2)
	require.Nil(ref)

	for idx, shapeRef := range p.IterShapeRefs(authShapeRef) {
		ref = p.ShapeRefAt(authShapeRef, idx)
		require.NotNil(shapeRef)
		require.NotNil(ref)
		require.Equal(ref.ShapeName, shapeRef.ShapeName)
	}

	// Path needs to match on outer-most shape first before going into member
	// refs
	p = fieldpath.FromString("Address")
	ref = p.ShapeRef(authShapeRef)
	require.Nil(ref)

	// More than a single level of nesting is possible...
	p = fieldpath.FromString("Author.Address.Country")
	ref = p.ShapeRef(authShapeRef)
	require.NotNil(ref)
	// Note that the ShapeName is actually different from the MemberName (which
	// is the string key in the MemberRefs map of the parent shape)
	require.Equal("CountryCode", ref.ShapeName)

	// Can't skip through a member...
	p = fieldpath.FromString("Author.Country")
	ref = p.ShapeRef(authShapeRef)
	require.Nil(ref)

	// Single dot-notation access of list member types...
	p = fieldpath.FromString("Author.Books")
	ref = p.ShapeRef(authShapeRef)
	require.NotNil(ref)
	require.Equal("list", ref.Shape.Type)

	p = fieldpath.FromString("Author.Books.Title")
	ref = p.ShapeRef(authShapeRef)
	require.NotNil(ref)
	require.Equal("Title", ref.ShapeName)
	require.Equal("string", ref.Shape.Type)

	// We support single dot notation even for deeply-nested map types
	p = fieldpath.FromString("Author.Books.ChapterPageCounts.PageCount")
	ref = p.ShapeRef(authShapeRef)
	require.NotNil(ref)
	require.Equal("PageCount", ref.ShapeName)
	require.Equal("integer", ref.Shape.Type)

	// Calling ShapeRef should not modify the original Path
	require.Equal("Author.Books.ChapterPageCounts.PageCount", p.String())

	// Case-insensitive comparisons...
	p = fieldpath.FromString("Author.WeirdlyCasedType")
	ref = p.ShapeRef(authShapeRef)
	require.NotNil(ref)
	require.Equal("WeirdlycasEdType", ref.ShapeName)
	require.Equal("string", ref.Shape.Type)
}
