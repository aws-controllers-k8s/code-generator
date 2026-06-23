package api

import (
	"testing"
)

func TestIsNonPointerInSDK_AddedDefault(t *testing.T) {
	// A member with @default(false) AND @addedDefault should remain a pointer
	ref := &ShapeRef{
		DefaultValue: "false",
		AddedDefault: true,
		Shape: &Shape{
			Type: "boolean",
		},
	}
	if ref.IsNonPointerInSDK() {
		t.Error("expected IsNonPointerInSDK() == false when AddedDefault is set, got true")
	}

	// Without AddedDefault, @default(false) on a boolean should be non-pointer
	ref2 := &ShapeRef{
		DefaultValue: "false",
		Shape: &Shape{
			Type: "boolean",
		},
	}
	if !ref2.IsNonPointerInSDK() {
		t.Error("expected IsNonPointerInSDK() == true for @default(false) boolean without AddedDefault, got false")
	}
}

func TestIsNonPointerInSDK_ClientOptional(t *testing.T) {
	// A member with @default(0) AND @clientOptional should remain a pointer
	ref := &ShapeRef{
		DefaultValue:   "0",
		ClientOptional: true,
		Shape: &Shape{
			Type: "integer",
		},
	}
	if ref.IsNonPointerInSDK() {
		t.Error("expected IsNonPointerInSDK() == false when ClientOptional is set, got true")
	}

	// Without ClientOptional, @default(0) on an integer should be non-pointer
	ref2 := &ShapeRef{
		DefaultValue: "0",
		Shape: &Shape{
			Type: "integer",
		},
	}
	if !ref2.IsNonPointerInSDK() {
		t.Error("expected IsNonPointerInSDK() == true for @default(0) integer without ClientOptional, got false")
	}
}

func TestIsNonPointerInSDK_NilDefault(t *testing.T) {
	// <nil> sentinel value means the default was explicitly cleared
	ref := &ShapeRef{
		DefaultValue: "<nil>",
		Shape: &Shape{
			Type: "boolean",
		},
	}
	if ref.IsNonPointerInSDK() {
		t.Error("expected IsNonPointerInSDK() == false for <nil> DefaultValue, got true")
	}
}

func TestIsNonPointerInSDK_NonZeroDefault(t *testing.T) {
	// @default(true) on a boolean should still be a pointer (non-zero default)
	ref := &ShapeRef{
		DefaultValue: "true",
		Shape: &Shape{
			Type: "boolean",
		},
	}
	if ref.IsNonPointerInSDK() {
		t.Error("expected IsNonPointerInSDK() == false for @default(true) boolean, got true")
	}
}

func TestIsNonPointerInSDK_ShapeLevelDefault(t *testing.T) {
	// Default value from shape level (e.g. PrimitiveBoolean)
	ref := &ShapeRef{
		Shape: &Shape{
			Type:         "boolean",
			DefaultValue: "false",
		},
	}
	if !ref.IsNonPointerInSDK() {
		t.Error("expected IsNonPointerInSDK() == true for shape-level @default(false) boolean, got false")
	}

	// Shape-level default with AddedDefault on the ref should still be pointer
	ref2 := &ShapeRef{
		AddedDefault: true,
		Shape: &Shape{
			Type:         "boolean",
			DefaultValue: "false",
		},
	}
	if ref2.IsNonPointerInSDK() {
		t.Error("expected IsNonPointerInSDK() == false when AddedDefault is set with shape-level default, got true")
	}
}
