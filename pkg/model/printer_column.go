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

package model

import (
	"fmt"
	"sort"
	"strings"
)

// PrinterColumn represents a single field in the CRD's Spec or Status objects
type PrinterColumn struct {
	CRD      *CRD
	Name     string
	Type     string
	JSONPath string
}

// By can sort two PrinterColumns
type By func(a, b *PrinterColumn) bool

// Sort does an in-place sort of the supplied printer columns
func (by By) Sort(subject []*PrinterColumn) {
	pcs := printerColumnSorter{
		cols: subject,
		by:   by,
	}
	sort.Sort(pcs)
}

// printerColumnSorter sorts printer columns by name
type printerColumnSorter struct {
	cols []*PrinterColumn
	by   By
}

// Len implements sort.Interface.Len
func (pcs printerColumnSorter) Len() int {
	return len(pcs.cols)
}

// Swap implements sort.Interface.Swap
func (pcs printerColumnSorter) Swap(i, j int) {
	pcs.cols[i], pcs.cols[j] = pcs.cols[j], pcs.cols[i]
}

// Less implements sort.Interface.Less
func (pcs printerColumnSorter) Less(i, j int) bool {
	return pcs.by(pcs.cols[i], pcs.cols[j])
}

// sortFunction returns a Go function used the sort the printer columns.
func sortFunction(sortByField string) func(a, b *PrinterColumn) bool {
	switch strings.ToLower(sortByField) {
	//TODO(a-hially): add Priority and Order sort functions
	case "name":
		return func(a, b *PrinterColumn) bool {
			return a.Name < b.Name
		}
	case "type":
		return func(a, b *PrinterColumn) bool {
			return a.Type < b.Type
		}
	case "jsonpath":
		return func(a, b *PrinterColumn) bool {
			return a.JSONPath < b.JSONPath
		}
	default:
		msg := fmt.Sprintf("unknown sort-by field: '%s'. must be one of 'Name', 'Type' and 'JSONPath'", sortByField)
		panic(msg)
	}
}

// AdditionalPrinterColumns returns a sorted list of PrinterColumn structs for
// the resource
func (r *CRD) AdditionalPrinterColumns() []*PrinterColumn {
	orderByFieldName := r.GetResourcePrintOrderByName()
	sortFn := sortFunction(orderByFieldName)
	By(sortFn).Sort(r.additionalPrinterColumns)
	return r.additionalPrinterColumns
}

// addPrintableColumn adds an entry to the list of additional printer columns
// using the given path and field types.
func (r *CRD) addPrintableColumn(
	field *Field,
	jsonPath string,
) {
	fieldColumnType := field.GoTypeElem

	// Printable columns must be primitives supported by the OpenAPI list of data
	// types as defined by
	// https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md#data-types
	// This maps Go type to OpenAPI type.
	acceptableColumnMaps := map[string]string{
		"string":      "string",
		"boolean":     "boolean",
		"int":         "integer",
		"int8":        "integer",
		"int16":       "integer",
		"int32":       "integer",
		"int64":       "integer",
		"uint":        "integer",
		"uint8":       "integer",
		"uint16":      "integer",
		"uint32":      "integer",
		"uint64":      "integer",
		"uintptr":     "integer",
		"float32":     "number",
		"float64":     "number",
		"metav1.Time": "date",
	}
	printColumnType, exists := acceptableColumnMaps[fieldColumnType]

	if !exists {
		msg := fmt.Sprintf(
			"GENERATION FAILURE! Unable to generate a printer column for the field %s that has type %s.",
			field.Names.Camel, fieldColumnType,
		)
		panic(msg)
	}

	name := field.Names.Camel
	if field.FieldConfig.Print.Name != "" {
		name = field.FieldConfig.Print.Name
	}

	column := &PrinterColumn{
		CRD:      r,
		Name:     name,
		Type:     printColumnType,
		JSONPath: jsonPath,
	}
	r.additionalPrinterColumns = append(r.additionalPrinterColumns, column)
}

// addSpecPrintableColumn adds an entry to the list of additional printer columns
// using the path of the given spec field.
func (r *CRD) addSpecPrintableColumn(
	field *Field,
) {
	r.addPrintableColumn(
		field,
		//TODO(nithomso): Ideally we'd use `r.cfg.PrefixConfig.SpecField` but it uses uppercase
		fmt.Sprintf("%s.%s", ".spec", field.Names.CamelLower),
	)
}

// addStatusPrintableColumn adds an entry to the list of additional printer columns
// using the path of the given status field.
func (r *CRD) addStatusPrintableColumn(
	field *Field,
) {
	r.addPrintableColumn(
		field,
		//TODO(nithomso): Ideally we'd use `r.cfg.PrefixConfig.StatusField` but it uses uppercase
		fmt.Sprintf("%s.%s", ".status", field.Names.CamelLower),
	)
}
