package model

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

func testFormItemBool(t *testing.T) *FormItemBool {
	mat := material.Chrome
	pos := mgl32.Vec3{0, 0, 0}
	fi := NewFormItemBool(DefaultMaxWidth, ITEM_WIDTH_HALF, 1.0, DefaultFormItemLabel, DefaultFormItemDescription, mat, pos, wrapperMock)

	if fi.label != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.label)
	}
	if fi.description != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.description)
	}
	return fi
}
func TestNewFormItemBool(t *testing.T) {
	_ = testFormItemBool(t)
}
func TestFormItemBoolGetLabel(t *testing.T) {
	fi := testFormItemBool(t)
	if fi.GetLabel() != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.GetLabel())
	}
}
func TestFormItemBoolGetDescription(t *testing.T) {
	fi := testFormItemBool(t)
	if fi.GetDescription() != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.GetDescription())
	}
}
func TestFormItemBoolGetValue(t *testing.T) {
	fi := testFormItemBool(t)
	val := true
	fi.value = val
	if fi.GetValue() != val {
		t.Errorf("Invalid form item value. Instead of '%v', it is '%v'.", val, fi.GetValue())
	}
}
func TestFormItemBoolSetValue(t *testing.T) {
	fi := testFormItemBool(t)
	val := true
	fi.value = val
	fi.SetValue(!val)
	if fi.GetValue() != !val {
		t.Errorf("Invalid form item value. Instead of '%v', it is '%v'.", !val, fi.GetValue())
	}
}
func TestFormItemBoolGetSurface(t *testing.T) {
	fi := testFormItemBool(t)
	if fi.GetSurface() != fi.meshes[0] {
		t.Error("Invalid surface mesh")
	}
}
func TestFormItemBoolGetLight(t *testing.T) {
	fi := testFormItemBool(t)
	if fi.GetLight() != fi.meshes[1] {
		t.Error("Invalid light mesh")
	}
}
func TestFormItemBoolValueToString(t *testing.T) {
	fi := testFormItemBool(t)
	if fi.ValueToString() != "false" {
		t.Error("Invalid value string.")
	}
	fi.SetValue(true)
	if fi.ValueToString() != "true" {
		t.Error("Invalid value string.")
	}
}
