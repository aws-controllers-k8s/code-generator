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

package code

import (
	"fmt"
	"sort"
	"strings"

	awssdkmodel "github.com/aws-controllers-k8s/code-generator/pkg/api"
	"github.com/aws-controllers-k8s/pkg/names"

	ackgenconfig "github.com/aws-controllers-k8s/code-generator/pkg/config"
	"github.com/aws-controllers-k8s/code-generator/pkg/model"
)

// CompareResource returns the Go code that traverses a set of two Resources,
// adding differences between the two Resources to an `ackcompare.Delta`
//
// By default, we produce Go code that only looks at the fields in a resource's
// Spec, since those are the fields that represent the desired state of a
// resource. When we make a ReadOne/ReadMany/GetAttributes call to a backend
// AWS API, we construct a Resource and set the Spec fields to values contained
// in the ReadOne/ReadMany/GetAttributes Output shape. This Resource,
// constructed from the Read operation, is compared to the Resource we got from
// the Kubernetes API server's event bus. The code that is returned from this
// function is the code that compares those two Resources.
//
// The Go code we return depends on the Go type of the various fields for the
// resource being compared.
//
// For *scalar* Go types, the output Go code looks like this:
//
//	if ackcompare.HasNilDifference(a.ko.Spec.GrantFullControl, b.ko.Spec.GrantFullControl) {
//	    delta.Add("Spec.GrantFullControl", a.ko.Spec.GrantFullControl, b.ko.Spec.GrantFullControl)
//	} else if a.ko.Spec.GrantFullControl != nil && b.ko.Spec.GrantFullControl != nil {
//
//	    if *a.ko.Spec.GrantFullControl != *b.ko.Spec.GrantFullControl {
//	        delta.Add("Spec.GrantFullControl", a.ko.Spec.GrantFullControl, b.ko.Spec.GrantFullControl)
//	    }
//	}
//
// For *struct* Go types, the output Go code looks like this (note that it is a
// simple recursive-descent output of all the struct's fields...):
//
//	if ackcompare.HasNilDifference(a.ko.Spec.CreateBucketConfiguration, b.ko.Spec.CreateBucketConfiguration) {
//	    delta.Add("Spec.CreateBucketConfiguration", a.ko.Spec.CreateBucketConfiguration, b.ko.Spec.CreateBucketConfiguration)
//	} else if a.ko.Spec.CreateBucketConfiguration != nil && b.ko.Spec.CreateBucketConfiguration != nil {
//
//	    if ackcompare.HasNilDifference(a.ko.Spec.CreateBucketConfiguration.LocationConstraint, b.ko.Spec.CreateBucketConfiguration.LocationConstraint) {
//	        delta.Add("Spec.CreateBucketConfiguration.LocationConstraint", a.ko.Spec.CreateBucketConfiguration.LocationConstraint, b.ko.Spec.CreateBucketConfiguration.LocationConstraint)
//	    } else if a.ko.Spec.CreateBucketConfiguration.LocationConstraint != nil && b.ko.Spec.CreateBucketConfiguration.LocationConstraint != nil {
//	        if *a.ko.Spec.CreateBucketConfiguration.LocationConstraint != *b.ko.Spec.CreateBucketConfiguration.LocationConstraint {
//	            delta.Add("Spec.CreateBucketConfiguration.LocationConstraint", a.ko.Spec.CreateBucketConfiguration.LocationConstraint, b.ko.Spec.CreateBucketConfiguration.LocationConstraint)
//	        }
//	    }
//	}
//
// For *slice of strings* Go types, the output Go code looks like this:
//
//	if ackcompare.HasNilDifference(a.ko.Spec.AllowedPublishers, b.ko.Spec.AllowedPublishers) {
//	    delta.Add("Spec.AllowedPublishers", a.ko.Spec.AllowedPublishers, b.ko.Spec.AllowedPublishers)
//	} else if a.ko.Spec.AllowedPublishers != nil && b.ko.Spec.AllowedPublishers != nil {
//
//	    if !ackcompare.SliceStringPEqual(a.ko.Spec.AllowedPublishers.SigningProfileVersionARNs, b.ko.Spec.AllowedPublishers.SigningProfileVersionARNs) {
//	        delta.Add("Spec.AllowedPublishers.SigningProfileVersionARNs", a.ko.Spec.AllowedPublishers.SigningProfileVersionARNs, b.ko.Spec.AllowedPublishers.SigningProfileVersionARNs)
//	    }
//	}
func CompareResource(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like "a.ko". See
	// `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like "b.ko". See
	// `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := "\n"

	resConfig := cfg.GetResourceConfig(r.Names.Camel)

	tagField, err := r.GetTagField()
	if err != nil {
		panic(err)
	}

	// We need a deterministic order to traverse our top-level fields...
	specFieldNames := []string{}
	for fieldName := range r.SpecFields {
		specFieldNames = append(specFieldNames, fieldName)
	}
	sort.Strings(specFieldNames)

	for _, fieldName := range specFieldNames {
		specField := r.SpecFields[fieldName]
		indent := strings.Repeat("\t", indentLevel)
		firstResAdaptedVarName := firstResVarName + cfg.PrefixConfig.SpecField
		firstResAdaptedVarName += "." + specField.Names.Camel
		secondResAdaptedVarName := secondResVarName + cfg.PrefixConfig.SpecField
		secondResAdaptedVarName += "." + specField.Names.Camel

		var fieldConfig *ackgenconfig.FieldConfig
		var compareConfig *ackgenconfig.CompareFieldConfig

		if resConfig != nil {
			fieldConfig = resConfig.GetFieldConfig(fieldName)
		}
		if fieldConfig != nil {
			compareConfig = fieldConfig.Compare
		}

		if compareConfig != nil && compareConfig.IsIgnored {
			continue
		}

		// this is the "path" to the field within the structs being compared.
		// This is passed down into the compareXXX functions recursively and
		// appended to with each level of nested structs we recurse into.
		fieldPath := strings.TrimPrefix(
			cfg.PrefixConfig.SpecField+"."+specField.Names.Camel, ".",
		)

		// Use equality.Semantic.Equalities.DeepEqual for comparing Reference fields because
		// some of reference fields are list of pointer to structs and
		// DeepEqual is easy way to compare them
		if specField.IsReference() {
			out += fmt.Sprintf("%sif !equality.Semantic.Equalities.DeepEqual(%s, %s) {\n",
				indent, firstResAdaptedVarName, secondResAdaptedVarName)
			out += fmt.Sprintf("%s\t%s.Add(\"%s\", %s, %s)\n", indent,
				deltaVarName, fieldPath, firstResAdaptedVarName,
				secondResAdaptedVarName)
			out += fmt.Sprintf("%s}\n", indent)
			continue
		}

		// Use a special comparison model for tags, since they need to be
		// converted into the common ACK tag type before doing a map delta
		if tagField != nil && specField == tagField {
			out += compareTags(deltaVarName, firstResAdaptedVarName, secondResAdaptedVarName, fieldPath, indentLevel)
			continue
		}

		memberShapeRef := specField.ShapeRef
		memberShape := memberShapeRef.Shape

		// Use len, bytes.Equal and HasNilDifference to fast compare types, and
		// try to avoid deep comparison as much as possible.
		fastComparisonOutput, needToCloseBlock := fastCompareTypes(
			compareConfig,
			memberShape,
			deltaVarName,
			fieldPath,
			firstResAdaptedVarName,
			secondResAdaptedVarName,
			indentLevel,
		)
		out += fastComparisonOutput

		switch memberShape.Type {
		case "blob":
			// We already handled the case of blobs above, so we can skip it here.
		case "structure":
			// Recurse through all the struct's fields and subfields, building
			// nested conditionals and calls to `delta.Add()`...
			out += CompareStruct(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				fieldPath,
				indentLevel+1,
			)
		case "list":
			// Returns Go code that compares all the elements of the slice fields...
			out += compareSlice(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				fieldPath,
				indentLevel+1,
			)
		case "map":
			// Returns Go code that compares all the elements of the map fields...
			out += compareMap(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				fieldPath,
				indentLevel+1,
			)
		default:
			//   if *a.ko.Spec.Name != *b.ko.Spec.Name) {
			//     delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
			//   }
			out += compareScalar(
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				fieldPath,
				indentLevel+1,
			)
		}
		if needToCloseBlock {
			// }
			out += fmt.Sprintf(
				"%s}\n", indent,
			)
		}
	}
	return out
}

// compareNil outputs Go code that compares two field values for nullability
// and, if there is a nil difference, adds the difference to a variable
// representing the `ackcompare.Delta`
//
// Output code will look something like this:
//
//	if ackcompare.HasNilDifferenceStringP(a.ko.Spec.Name, b.ko.Spec.Name == nil) {
//	  delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
//	}
func compareNil(
	// struct informing code generator how to compare the field values
	compareConfig *ackgenconfig.CompareFieldConfig,
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	switch shape.Type {
	case "boolean", "string", "character", "byte", "short", "integer", "long",
		"float", "double", "timestamp", "structure", "jsonvalue":
		// if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name) {
		out += fmt.Sprintf(
			"%sif ackcompare.HasNilDifference(%s, %s) {\n",
			indent, firstResVarName, secondResVarName,
		)
	default:
		panic("Unsupported shape type in generate.code.compareNil: " + shape.Type)
	}
	//   delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	out += fmt.Sprintf(
		"%s\t%s.Add(\"%s\", %s, %s)\n",
		indent, deltaVarName, fieldPath, firstResVarName, secondResVarName,
	)
	// }
	out += fmt.Sprintf(
		"%s}", indent,
	)

	return out
}

// compareScalar outputs Go code that compares two scalar values from two
// resource fields and, if there is a difference, adds the difference to a
// variable representing an `ackcompare.Delta`.
//
// Output code will look something like this:
//
//	if *a.ko.Spec.Name != *b.ko.Spec.Name) {
//	  delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
//	}
func compareScalar(
	// struct informing code generator how to compare the field values
	compareConfig *ackgenconfig.CompareFieldConfig,
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	switch shape.Type {
	case "boolean", "string", "character", "byte", "short", "integer", "long", "float", "double":
		// if *a.ko.Spec.Name != *b.ko.Spec.Name {
		out += fmt.Sprintf(
			"%sif *%s != *%s {\n",
			indent, firstResVarName, secondResVarName,
		)
	case "timestamp":
		// if !a.ko.Spec.CreatedAt.Equal(b.ko.Spec.CreatedAt) {
		out += fmt.Sprintf(
			"%sif !%s.Equal(%s) {\n",
			indent, firstResVarName, secondResVarName,
		)
	default:
		panic("Unsupported shape type in generate.code.compareScalar: " + shape.Type)
	}
	//   delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	out += fmt.Sprintf(
		"%s\t%s.Add(\"%s\", %s, %s)\n",
		indent, deltaVarName, fieldPath, firstResVarName, secondResVarName,
	)
	// }
	out += fmt.Sprintf(
		"%s}\n", indent,
	)

	return out
}

// compareMap outputs Go code that compares two map values from two resource
// fields and, if there is a difference, adds the difference to a variable
// representing an `ackcompare.Delta`.
//
// Output code will look something like this:
//
//	if !ackcompare.MapStringStringPEqual(a.ko.Spec.Tags, b.ko.Spec.Tags) {
//	  delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
//	}
func compareMap(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// struct informing code generator how to compare the field values
	compareConfig *ackgenconfig.CompareFieldConfig,
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	keyType := shape.KeyRef.Shape.Type

	if keyType != "string" {
		panic("generate.code.compareMap cannot deal with non-string key types: " + keyType)
	}

	valType := shape.ValueRef.Shape.Type

	switch valType {
	case "string":
		// if !ackcompare.MapStringStringPEqual(a.ko.Spec.Tags, b.ko.Spec.Tags) {
		out += fmt.Sprintf(
			"%sif !ackcompare.MapStringStringPEqual(%s, %s) {\n",
			indent, firstResVarName, secondResVarName,
		)
	default:
		// NOTE(jaypipes): Using reflect here is really punting. We should
		// implement this in a cleaner, more efficient fashion by walking the
		// keys and struct values and comparing each struct individually,
		// building up the fieldPath appropriately and calling into a
		// struct-specific comparator function...
		out += fmt.Sprintf(
			"%sif !equality.Semantic.Equalities.DeepEqual(%s, %s) {\n",
			indent, firstResVarName, secondResVarName,
		)
	}
	//   delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	out += fmt.Sprintf(
		"%s\t%s.Add(\"%s\", %s, %s)\n",
		indent, deltaVarName, fieldPath, firstResVarName, secondResVarName,
	)
	// }
	out += fmt.Sprintf(
		"%s}\n", indent,
	)

	return out
}

// compareSlice outputs Go code that compares two slice values from two
// resource fields and, if there is a difference, adds the difference to a
// variable representing an `ackcompare.Delta`.
//
// Output code will look something like this:
//
//	if !ackcompare.SliceStringPEqual(a.ko.Spec.SecurityGroupIDs, b.ko.Spec.SecurityGroupIDs) {
//	  delta.Add("Spec.SecurityGroupIDs", a.ko.Spec.SecurityGroupIDs, b.ko.Spec.SecurityGroupIDs)
//	}
func compareSlice(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// struct informing code generator how to compare the field values
	compareConfig *ackgenconfig.CompareFieldConfig,
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	elemType := shape.MemberRef.Shape.Type

	switch elemType {
	case "string":
		// if !ackcompare.SliceStringPEqual(a.ko.Spec.SecurityGroupIDs, b.ko.Spec.SecurityGroupIDs) {
		out += fmt.Sprintf(
			"%sif !ackcompare.SliceStringPEqual(%s, %s) {\n",
			indent, firstResVarName, secondResVarName,
		)
	case "structure", "union":
		// NOTE(jaypipes): Using reflect here is really punting. We should
		// implement this in a cleaner, more efficient fashion by walking the
		// struct values and comparing each struct individually, building up
		// the fieldPath appropriately and calling into a struct-specific
		// comparator function...the tricky part of this is figuring out how to
		// sort the slice of structs...
		out += fmt.Sprintf(
			"%sif !equality.Semantic.Equalities.DeepEqual(%s, %s) {\n",
			indent, firstResVarName, secondResVarName,
		)
	default:
		panic("Unsupported shape type in generate.code.compareSlice: " + shape.Type)
	}
	//   delta.Add("Spec.SecurityGroupIDs", a.ko.Spec.SecurityGroupIDs, b.ko.Spec.SecurityGroupIDs)
	out += fmt.Sprintf(
		"%s\t%s.Add(\"%s\", %s, %s)\n",
		indent, deltaVarName, fieldPath, firstResVarName, secondResVarName,
	)
	// }
	out += fmt.Sprintf(
		"%s}\n", indent,
	)

	return out
}

// compareTags outputs Go code that compares two slices of tags from two
// resource fields by first converting them to the common ACK tag type and then
// using a map comparison. If there is a difference, adds the difference to a
// variable representing an `ackcompare.Delta`.
//
// Output code will look something like this:
//
//	 desiredACKTags, _ := convertToOrderedACKTags(a.ko.Spec.Tags)
//	 latestACKTags, _ := convertToOrderedACKTags(b.ko.Spec.Tags)
//		if !ackcompare.MapStringStringEqual(desiredACKTags, latestACKTags) {
//		  delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
//		}
func compareTags(
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""
	indent := strings.Repeat("\t", indentLevel)

	out += fmt.Sprintf("%sdesiredACKTags, _ := convertToOrderedACKTags(%s)\n", indent, firstResVarName)
	out += fmt.Sprintf("%slatestACKTags, _ := convertToOrderedACKTags(%s)\n", indent, secondResVarName)
	out += fmt.Sprintf("%sif !ackcompare.MapStringStringEqual(desiredACKTags, latestACKTags) {\n", indent)
	out += fmt.Sprintf("%s\t%s.Add(\"%s\", %s, %s)\n", indent, deltaVarName, fieldPath, firstResVarName, secondResVarName)
	out += fmt.Sprintf("%s}\n", indent)

	return out
}

// CompareStruct outputs Go code that compares two struct values from two
// resource fields and, if there is a difference, adds the difference to a
// variable representing an `ackcompare.Delta`.
func CompareStruct(
	cfg *ackgenconfig.Config,
	r *model.CRD,
	// struct informing code generator how to compare the field values
	compareConfig *ackgenconfig.CompareFieldConfig,
	// struct describing the SDK type of the field being compared
	shape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like
	// "a.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like
	// "b.ko.Spec.Name". See `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// String indicating the current field path being evaluated, e.g.
	// "Author.Name". This does not include the top-level Spec or Status
	// struct.
	fieldPath string,
	// Number of levels of indentation to use
	indentLevel int,
) string {
	out := ""

	tagField, err := r.GetTagField()
	if err != nil {
		panic(err)
	}

	fieldConfigs := cfg.GetFieldConfigs(r.Names.Original)

	for _, memberName := range shape.MemberNames() {
		memberShapeRef := shape.MemberRefs[memberName]
		// TODO(jaypipes): This is fragile. We actually need to have a way of
		// normalizing names in a lossless fashion...
		memberNames := names.New(memberName)
		memberNameClean := memberNames.Camel

		// this is the "path" to the field within the structs being compared.
		// This is passed down into the compareXXX functions recursively and
		// appended to with each level of nested structs we recurse into.
		memberFieldPath := fieldPath + "." + memberNameClean
		indent := strings.Repeat("\t", indentLevel)
		firstResAdaptedVarName := firstResVarName + "." + memberNameClean
		secondResAdaptedVarName := secondResVarName + "." + memberNameClean

		var compareConfig *ackgenconfig.CompareFieldConfig
		// memberFieldPath contains the field path along with the prefix cfg.PrefixConfig.SpecField + "." hence we
		// would need to substring to exclude cfg.PrefixConfig.SpecField + "." to get correct field config.
		specFieldLen := len(strings.TrimPrefix(cfg.PrefixConfig.SpecField, "."))
		trimmedFieldPath := memberFieldPath[specFieldLen+1:]

		fieldConfig := fieldConfigs[trimmedFieldPath]
		if fieldConfig != nil {
			compareConfig = fieldConfig.Compare
		}

		if compareConfig != nil && compareConfig.IsIgnored {
			continue
		}

		memberShape := memberShapeRef.Shape

		// Use a special comparison model for tags, since they need to be
		// converted into the common ACK tag type before doing a map delta
		if tagField != nil && tagField.Path == trimmedFieldPath {
			out += compareTags(deltaVarName, firstResAdaptedVarName, secondResAdaptedVarName, fieldPath, indentLevel)
			continue
		}

		fastComparisonOutput, needToCloseBlock := fastCompareTypes(
			compareConfig,
			memberShape,
			deltaVarName,
			memberFieldPath,
			firstResAdaptedVarName,
			secondResAdaptedVarName,
			indentLevel,
		)
		out += fastComparisonOutput

		switch memberShape.Type {
		case "blob":
			// We already handled the case of blobs above, so we can skip it here.
		case "structure":
			// Recurse through all the struct's fields and subfields, building
			// nested conditionals and calls to `delta.Add()`...
			out += CompareStruct(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				memberFieldPath,
				indentLevel+1,
			)
		case "list":
			// Returns Go code that compares all the elements of the slice fields...
			out += compareSlice(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				memberFieldPath,
				indentLevel+1,
			)
		case "map":
			// Returns Go code that compares all the elements of the map fields...
			out += compareMap(
				cfg, r,
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				memberFieldPath,
				indentLevel+1,
			)
		default:
			//   if *a.ko.Spec.Name != *b.ko.Spec.Name {
			//     delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
			//   }
			out += compareScalar(
				compareConfig,
				memberShape,
				deltaVarName,
				firstResAdaptedVarName,
				secondResAdaptedVarName,
				memberFieldPath,
				indentLevel+1,
			)
		}
		if needToCloseBlock {
			// }
			out += fmt.Sprintf(
				"%s}\n", indent,
			)
		}
	}
	return out
}

// fastCompareTypes outputs Go code that fast-compares two objects of the same
// type, by leveraging nil check comparison and length checks. This is used
// when we want to quickly determine if two objects are different, but don't
// need to know the specific differences. For example, when determining if a
// that an array of structs has changed, we don't need to know which structs
// have changed, if the size of the array has changed.
//
// Generally, we can distinguish between the following cases:
// 1. Both objects are structures or scalars type pointers.
// 2. Both objects are collections (slices or  maps).
// 3. Both objects are blobs.
//
// In the case of 1, we can use the HasNilDifference function to
// eliminate early on the case where one of the objects is nil and other
// is not.
//
// In the case of 2, we can use the built-in len function to eliminate
// early on the case where the collections have different lengths.
// The trick here is that len(nil) is 0, so we don't need to check for
// nils explicitly.
//
// In the case of 3, we can use the bytes.Equal function to compare the
// byte arrays. bytes.Equal works well with nil arrays too.
// Output code will look something like this:
//
//	if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name) {
//	  delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
//	} else if a.ko.Spec.Name != nil && b.ko.Spec.Name != nil {
//	  if *a.ko.Spec.Name != *b.ko.Spec.Name) {
//	    delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
//	  }
//	}
//
//	if len(a.ko.Spec.Tags) != len(b.ko.Spec.Tags) {
//	  delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
//	} else if len(a.ko.Spec.Tags) > 0 {
//	  if !ackcompare.SliceStringPEqual(a.ko.Spec.Tags, b.ko.Spec.Tags) {
//	    delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
//	  }
//	}
func fastCompareTypes(
	// struct informing code generator how to compare the field values.
	compareConfig *ackgenconfig.CompareFieldConfig,
	// AWSSDK Shape describing the type of the field being compared.
	memberShape *awssdkmodel.Shape,
	// String representing the name of the variable that is of type
	// `*ackcompare.Delta`. We will generate Go code that calls the `Add()`
	// method of this variable when differences between fields are detected.
	deltaVarName string,
	// String representing the json path of the field being compared, e.g.
	// "Spec.Name".
	fieldPath string,
	// String representing the name of the variable that represents the first
	// CR under comparison. This will typically be something like "a.ko". See
	// `templates/pkg/resource/delta.go.tpl`.
	firstResVarName string,
	// String representing the name of the variable that represents the second
	// CR under comparison. This will typically be something like "b.ko". See
	// `templates/pkg/resource/delta.go.tpl`.
	secondResVarName string,
	// Number of levels of indentation to use
	indentLevel int,
) (string, bool) {
	out := ""
	needToCloseBlock := false
	indent := strings.Repeat("\t", indentLevel)

	switch memberShape.Type {
	case "list", "map":
		// Note once we eliminate the case of non equal sizes, we can
		// safely assume that the objects have the same length and we can
		// we can only iterate over them if they have more than 0 elements.
		//
		// if len(a.ko.Spec.Tags) != len(b.ko.Spec.Tags) {
		//   delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
		// } else len(a.ko.Spec.Tags) > 0 {
		//    ...
		// }
		out += fmt.Sprintf(
			"%sif len(%s) != len(%s) {\n",
			indent,
			firstResVarName,
			secondResVarName,
		)
		out += fmt.Sprintf(
			"%s\t%s.Add(\"%s\", %s, %s)\n",
			indent,
			deltaVarName,
			fieldPath,
			firstResVarName,
			secondResVarName,
		)
		// For sure we are inside the else block of the nil/size check, so we need
		// to increase the indentation level.
		out += fmt.Sprintf(
			"%s} else if len(%s) > 0 {\n",
			indent,
			firstResVarName,
		)
		// For sure we are inside the else block of the nil/size check, so we need
		// to increase the indentation level and ask the caller to close the block.
		needToCloseBlock = true
	case "blob":
		// Blob is a special case because we need to compare the byte arrays
		// using the bytes.Equal function.
		//
		// if !bytes.Equal(a.ko.Spec.Certificate, b.ko.Spec.Certificate) {
		//   delta.Add("Spec.Certificate", a.ko.Spec.Certificate, b.ko.Spec.Certificate)
		// }
		out += fmt.Sprintf(
			"%sif !bytes.Equal(%s, %s) {\n",
			indent,
			firstResVarName,
			secondResVarName,
		)
		out += fmt.Sprintf(
			"%s\t%s.Add(\"%s\", %s, %s)\n",
			indent,
			deltaVarName,
			fieldPath,
			firstResVarName,
			secondResVarName,
		)
		out += fmt.Sprintf(
			"%s}\n", indent,
		)
	default:
		// For any other type, we can use the HasNilDifference function to
		// eliminate early on the case where one of the objects is nil and
		// other is not.
		out += compareNil(
			compareConfig,
			memberShape,
			deltaVarName,
			firstResVarName,
			secondResVarName,
			fieldPath,
			indentLevel,
		)

		// if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name) {
		//   delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
		// } else if a.ko.Spec.Name != nil && b.ko.Spec.Name != nil {
		out += fmt.Sprintf(
			" else if %s != nil && %s != nil {\n",
			firstResVarName, secondResVarName,
		)
		// For sure we are inside the else block of the nil/size check, so we need
		// to increase the indentation level and ask the caller to close the block.
		needToCloseBlock = true
	}
	return out, needToCloseBlock
}
