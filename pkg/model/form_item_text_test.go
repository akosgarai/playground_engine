package model

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

func testFormItemText(t *testing.T) *FormItemText {
	mat := material.Chrome
	pos := mgl32.Vec3{0, 0, 0}
	fi := NewFormItemText(DefaultMaxWidth, ITEM_WIDTH_HALF, 1.0, DefaultFormItemLabel, DefaultFormItemDescription, mat, pos, wrapperMock)

	if fi.label != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.label)
	}
	if fi.description != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.description)
	}
	return fi
}
func TestNewFormItemText(t *testing.T) {
	_ = testFormItemText(t)
}
func TestFormItemTextGetLabel(t *testing.T) {
	fi := testFormItemText(t)
	if fi.GetLabel() != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.GetLabel())
	}
}
func TestFormItemTextGetDescription(t *testing.T) {
	fi := testFormItemText(t)
	if fi.GetDescription() != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.GetDescription())
	}
}
func TestFormItemTextGetValue(t *testing.T) {
	fi := testFormItemText(t)
	val := "test value"
	fi.value = val
	if fi.GetValue() != val {
		t.Errorf("Invalid form item value. Instead of '%s', it is '%s'.", val, fi.GetValue())
	}
}
func TestFormItemTextGetSurface(t *testing.T) {
	fi := testFormItemText(t)
	if fi.GetSurface() != fi.meshes[0] {
		t.Error("Invalid surface mesh")
	}
}
func TestFormItemTextGetTarget(t *testing.T) {
	fi := testFormItemText(t)
	if fi.GetTarget() != fi.meshes[1] {
		t.Error("Invalid target mesh")
	}
}
func TestFormItemTextAddCursor(t *testing.T) {
	fi := testFormItemText(t)
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemTextDeleteCursor(t *testing.T) {
	fi := testFormItemText(t)
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
	fi.DeleteCursor()
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemTextCharCallback(t *testing.T) {
	fi := testFormItemText(t)
	fi.AddCursor()
	fi.CharCallback('b', 0.1)
	if fi.value != "b" {
		t.Errorf("Invalid value. Instead of 'b', we have '%s'.", fi.value)
	}
	fi.CharCallback('a', 0.1)
	if fi.value != "ba" {
		t.Errorf("Invalid value. Instead of 'ba', we have '%s'.", fi.value)
	}
	fi.CharCallback(' ', 0.1)
	if fi.value != "ba " {
		t.Errorf("Invalid value. Instead of 'ba ', we have '%s'.", fi.value)
	}
	fi.CharCallback('b', 0.1)
	if fi.value != "ba b" {
		t.Errorf("Invalid value. Instead of 'ba b', we have '%s'.", fi.value)
	}
	fi.CharCallback('b', 0.1)
	if fi.value != "ba bb" {
		t.Errorf("Invalid value. Instead of 'ba bb', we have '%s'.", fi.value)
	}
	fi.CharCallback('b', 0.1)
	if fi.value != "ba bbb" {
		t.Errorf("Invalid value. Instead of 'ba bbb', we have '%s'.", fi.value)
	}
	fi.CharCallback('b', 0.1)
	if fi.value != "ba bbbb" {
		t.Errorf("Invalid value. Instead of 'ba bbbb', we have '%s'.", fi.value)
	}
	fi.CharCallback('b', 0.1)
	if fi.value != "ba bbbbb" {
		t.Errorf("Invalid value. Instead of 'ba bbbbb', we have '%s'.", fi.value)
	}
	fi.CharCallback('b', 0.1)
	if fi.value != "ba bbbbbb" {
		t.Errorf("Invalid value. Instead of 'ba bbbbbb', we have '%s'.", fi.value)
	}
	valueBase := "ba bbbbbb"
	for i := 0; i < 25; i++ {
		valueBase = valueBase + "b"
		fi.CharCallback('b', 0.1)
		if fi.value != valueBase {
			t.Errorf("Invalid value (%d). Instead of '%s', we have '%s'.", i, valueBase, fi.value)
		}
	}
	for i := 0; i < 10; i++ {
		fi.CharCallback('b', 0.1)
		if fi.value != valueBase {
			t.Errorf("Invalid value (%d). Instead of '%s', we have '%s'.", i, valueBase, fi.value)
		}
	}
	fi = testFormItemText(t)
	fi.AddCursor()
	validator := func(s string) bool {
		return s != "not"
	}
	fi.SetValidator(validator)
	fi.CharCallback('n', 0.1)
	if fi.value != "n" {
		t.Errorf("Invalid value. Instead of 'n', we have '%s'.", fi.value)
	}
	fi.CharCallback('o', 0.1)
	if fi.value != "no" {
		t.Errorf("Invalid value. Instead of 'no', we have '%s'.", fi.value)
	}
	fi.CharCallback('t', 0.1)
	if fi.value != "no" {
		t.Errorf("Invalid value. Instead of 'no', we have '%s'.", fi.value)
	}
	fi.CharCallback('r', 0.1)
	if fi.value != "nor" {
		t.Errorf("Invalid value. Instead of 'nor', we have '%s'.", fi.value)
	}
}
func TestFormItemTextValueToString(t *testing.T) {
	fi := testFormItemText(t)
	val := fi.ValueToString()
	if val != "" {
		t.Errorf("Invalid value. Instead of '', we have '%s'.", val)
	}
	fi.value = "test value"
	val = fi.ValueToString()
	if val != "test value" {
		t.Errorf("Invalid value. Instead of 'test value', we have '%s'.", val)
	}
}
func TestFormItemTextDeleteLastCharacter(t *testing.T) {
	fi := testFormItemText(t)
	fi.value = "12345"
	fi.charOffsets = []float32{0.1, 0.1, 0.1, 0.1, 0.1}
	fi.DeleteLastCharacter()
	if fi.value != "1234" {
		t.Errorf("Invalid value. Instead of '1234', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	if fi.value != "123" {
		t.Errorf("Invalid value. Instead of '123', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	if fi.value != "12" {
		t.Errorf("Invalid value. Instead of '12', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	if fi.value != "1" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	if fi.value != "" {
		t.Errorf("Invalid value. Instead of '', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	if fi.value != "" {
		t.Errorf("Invalid value. Instead of '', we have '%s'.", fi.value)
	}
}
