package model

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

func testFormItemFloat(t *testing.T) *FormItemFloat {
	mat := material.Chrome
	pos := mgl32.Vec3{0, 0, 0}
	fi := NewFormItemFloat(DefaultMaxWidth, ITEM_WIDTH_HALF, 1.0, DefaultFormItemLabel, DefaultFormItemDescription, mat, pos, wrapperMock)

	if fi.label != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.label)
	}
	if fi.description != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.description)
	}
	return fi
}
func TestNewFormItemFloat(t *testing.T) {
	_ = testFormItemFloat(t)
}
func TestFormItemFloatGetLabel(t *testing.T) {
	fi := testFormItemFloat(t)
	if fi.GetLabel() != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.GetLabel())
	}
}
func TestFormItemFloatGetDescription(t *testing.T) {
	fi := testFormItemFloat(t)
	if fi.GetDescription() != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.GetDescription())
	}
}
func TestFormItemFloatGetValue(t *testing.T) {
	fi := testFormItemFloat(t)
	valI := 3
	fi.value = "3"
	if fi.GetValue() != float32(valI) {
		t.Errorf("Invalid form item value. Instead of '%d', it is '%f'.", valI, fi.GetValue())
	}
	fi.value = "3.2"
	if fi.GetValue() != 3.2 {
		t.Errorf("Invalid form item value. Instead of '3.2', it is '%f'.", fi.GetValue())
	}
}
func TestFormItemFloatGetSurface(t *testing.T) {
	fi := testFormItemFloat(t)
	if fi.GetSurface() != fi.meshes[0] {
		t.Error("Invalid surface mesh")
	}
}
func TestFormItemFloatGetTarget(t *testing.T) {
	fi := testFormItemFloat(t)
	if fi.GetTarget() != fi.meshes[1] {
		t.Error("Invalid target mesh")
	}
}
func TestFormItemFloatAddCursor(t *testing.T) {
	fi := testFormItemFloat(t)
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemFloatDeleteCursor(t *testing.T) {
	fi := testFormItemFloat(t)
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
	fi.DeleteCursor()
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemFloatCharCallback(t *testing.T) {
	fi := testFormItemFloat(t)
	fi.AddCursor()
	// start with 0
	fi.CharCallback('0', 0.1)
	if fi.value != "0" {
		t.Errorf("Invalid value. Instead of '0', we have '%s'.", fi.value)
	}
	if fi.typeState != "P0" {
		t.Errorf("Invalid typeState. Instead of 'P0', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "0" {
		t.Errorf("Invalid value. Instead of '0', we have '%s'.", fi.value)
	}
	if fi.typeState != "P0" {
		t.Errorf("Invalid typeState. Instead of 'P0', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "0." {
		t.Errorf("Invalid value. Instead of '0.', we have '%s'.", fi.value)
	}
	if fi.typeState != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('2', 0.1)
	if fi.value != "0.2" {
		t.Errorf("Invalid value. Instead of '0.2', we have '%s'.", fi.value)
	}
	if fi.typeState != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('8', 0.1)
	if fi.value != "0.28" {
		t.Errorf("Invalid value. Instead of '0.28', we have '%s'.", fi.value)
	}
	if fi.typeState != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeState)
	}
	// start with .
	fi = testFormItemFloat(t)
	fi.AddCursor()
	fi.CharCallback('.', 0.1)
	if fi.value != "" {
		t.Errorf("Invalid value. Instead of '', we have '%s'.", fi.value)
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "1" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.value)
	}
	if fi.typeState != "PI" {
		t.Errorf("Invalid typeState. Instead of 'PI', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('2', 0.1)
	if fi.value != "12" {
		t.Errorf("Invalid value. Instead of '12', we have '%s'.", fi.value)
	}
	if fi.typeState != "PI" {
		t.Errorf("Invalid typeState. Instead of 'PI', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "12." {
		t.Errorf("Invalid value. Instead of '12.', we have '%s'.", fi.value)
	}
	if fi.typeState != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('8', 0.1)
	if fi.value != "12.8" {
		t.Errorf("Invalid value. Instead of '12.8', we have '%s'.", fi.value)
	}
	if fi.typeState != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeState)
	}
	// start with -
	fi = testFormItemFloat(t)
	fi.AddCursor()
	fi.CharCallback('-', 0.1)
	if fi.value != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.value)
	}
	if fi.typeState != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.value)
	}
	if fi.typeState != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "-1" {
		t.Errorf("Invalid value. Instead of '-1', we have '%s'.", fi.value)
	}
	if fi.typeState != "NI" {
		t.Errorf("Invalid typeState. Instead of 'NI', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('2', 0.1)
	if fi.value != "-12" {
		t.Errorf("Invalid value. Instead of '-12', we have '%s'.", fi.value)
	}
	if fi.typeState != "NI" {
		t.Errorf("Invalid typeState. Instead of 'NI', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "-12." {
		t.Errorf("Invalid value. Instead of '-12.', we have '%s'.", fi.value)
	}
	if fi.typeState != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('8', 0.1)
	if fi.value != "-12.8" {
		t.Errorf("Invalid value. Instead of '-12.8', we have '%s'.", fi.value)
	}
	if fi.typeState != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeState)
	}
	// start with -0.
	fi = testFormItemFloat(t)
	fi.AddCursor()
	fi.CharCallback('-', 0.1)
	if fi.value != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.value)
	}
	if fi.typeState != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('0', 0.1)
	if fi.value != "-0" {
		t.Errorf("Invalid value. Instead of '-0', we have '%s'.", fi.value)
	}
	if fi.typeState != "N0" {
		t.Errorf("Invalid typeState. Instead of 'N0', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "-0." {
		t.Errorf("Invalid value. Instead of '-0.', we have '%s'.", fi.value)
	}
	if fi.typeState != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('8', 0.1)
	if fi.value != "-0.8" {
		t.Errorf("Invalid value. Instead of '-0.8', we have '%s'.", fi.value)
	}
	if fi.typeState != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeState)
	}
	// validator
	fi = testFormItemFloat(t)
	fi.AddCursor()
	validator := func(i float32) bool {
		return i <= 1.0 && i >= 0.0
	}
	fi.SetValidator(validator)
	fi.CharCallback('-', 0.1)
	if fi.value != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.value)
	}
	fi.CharCallback('0', 0.1)
	if fi.value != "-0" {
		t.Errorf("Invalid value. Instead of '-0', we have '%s'.", fi.value)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "-0." {
		t.Errorf("Invalid value. Instead of '-0.', we have '%s'.", fi.value)
	}
	fi.CharCallback('3', 0.1)
	if fi.value != "-0." {
		t.Errorf("Invalid value. Instead of '-0.', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	fi.DeleteLastCharacter()
	fi.DeleteLastCharacter()
	fi.CharCallback('1', 0.1)
	if fi.value != "1" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.value)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "1." {
		t.Errorf("Invalid value. Instead of '1.', we have '%s'.", fi.value)
	}
	fi.CharCallback('3', 0.1)
	if fi.value != "1." {
		t.Errorf("Invalid value. Instead of '1.', we have '%s'.", fi.value)
	}
	fi.CharCallback('0', 0.1)
	if fi.value != "1.0" {
		t.Errorf("Invalid value. Instead of '1.0', we have '%s'.", fi.value)
	}
}
func TestFormItemFloatValueToString(t *testing.T) {
	fi := testFormItemFloat(t)
	fi.AddCursor()
	fi.value = "-"
	if fi.ValueToString() != "-" {
		t.Errorf("Invalid valuestring. instead of '-', we have '%s'.", fi.ValueToString())
	}
	fi.value = "-3"
	if fi.ValueToString() != "-3" {
		t.Errorf("Invalid valuestring. instead of '-3', we have '%s'.", fi.ValueToString())
	}
	fi.value = "-3."
	if fi.ValueToString() != "-3." {
		t.Errorf("Invalid valuestring. instead of '-3.', we have '%s'.", fi.ValueToString())
	}
	fi.value = "-3.3"
	if fi.ValueToString() != "-3.3" {
		t.Errorf("Invalid valuestring. instead of '-3.3', we have '%s'.", fi.ValueToString())
	}

}
func TestFormItemFloatDeleteLastCharacter(t *testing.T) {
	fi := testFormItemFloat(t)
	fi.AddCursor()
	fi.CharCallback('-', 0.1)
	fi.CharCallback('3', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('3', 0.1)
	if fi.ValueToString() != "-3.3" {
		t.Errorf("Invalid valuestring. instead of '-3.3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-3." {
		t.Errorf("Invalid valuestring. instead of '-3.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-3" {
		t.Errorf("Invalid valuestring. instead of '-3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "NI" {
		t.Errorf("Invalid typeState. Instead of 'NI', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-" {
		t.Errorf("Invalid valuestring. instead of '-', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('-', 0.1)
	fi.CharCallback('0', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('2', 0.1)
	if fi.ValueToString() != "-0.2" {
		t.Errorf("Invalid valuestring. instead of '-0.2', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-0." {
		t.Errorf("Invalid valuestring. instead of '-0.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-0" {
		t.Errorf("Invalid valuestring. instead of '-0', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "N0" {
		t.Errorf("Invalid typeState. Instead of 'N0', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-" {
		t.Errorf("Invalid valuestring. instead of '-', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('3', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('3', 0.1)
	if fi.ValueToString() != "3.3" {
		t.Errorf("Invalid valuestring. instead of '3.3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "3." {
		t.Errorf("Invalid valuestring. instead of '3.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "3" {
		t.Errorf("Invalid valuestring. instead of '3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "PI" {
		t.Errorf("Invalid typeState. Instead of 'PI', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('0', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('3', 0.1)
	if fi.ValueToString() != "0.3" {
		t.Errorf("Invalid valuestring. instead of '0.3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "0." {
		t.Errorf("Invalid valuestring. instead of '0.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "0" {
		t.Errorf("Invalid valuestring. instead of '0', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P0" {
		t.Errorf("Invalid typeState. Instead of 'P0', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
}
