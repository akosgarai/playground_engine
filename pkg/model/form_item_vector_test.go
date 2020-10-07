package model

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

func testFormItemVector(t *testing.T) *FormItemVector {
	mat := material.Chrome
	pos := mgl32.Vec3{0, 0, 0}
	fi := NewFormItemVector(DefaultMaxWidth, ITEM_WIDTH_FULL, 1.0, DefaultFormItemLabel, DefaultFormItemDescription, 10, mat, pos, wrapperMock)

	if fi.label != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.label)
	}
	if fi.description != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.description)
	}
	return fi
}
func TestNewFormItemVector(t *testing.T) {
	_ = testFormItemVector(t)
}
func TestFormItemVectorGetLabel(t *testing.T) {
	fi := testFormItemVector(t)
	if fi.GetLabel() != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.GetLabel())
	}
}
func TestFormItemVectorGetDescription(t *testing.T) {
	fi := testFormItemVector(t)
	if fi.GetDescription() != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.GetDescription())
	}
}
func TestFormItemVectorGetValue(t *testing.T) {
	fi := testFormItemVector(t)
	nullVec := mgl32.Vec3{0, 0, 0}
	if fi.GetValue() != nullVec {
		t.Errorf("Invalid form item value. Instead of '%v', it is '%v'.", nullVec, fi.GetValue())
	}
	fi.values[0] = "3.2"
	fi.values[1] = "4.1"
	fi.values[2] = "5.0"
	vec := mgl32.Vec3{3.2, 4.1, 5.0}
	if fi.GetValue() != vec {
		t.Errorf("Invalid form item value. Instead of '%v', it is '%v'.", vec, fi.GetValue())
	}
}
func TestFormItemVectorGetSurface(t *testing.T) {
	fi := testFormItemVector(t)
	if fi.GetSurface() != fi.meshes[0] {
		t.Error("Invalid surface mesh")
	}
}
func TestFormItemVectorGetTarget(t *testing.T) {
	fi := testFormItemVector(t)
	if fi.GetTarget() != fi.meshes[1] {
		t.Error("Invalid target mesh")
	}
}
func TestFormItemVectorAddCursor(t *testing.T) {
	fi := testFormItemVector(t)
	if len(fi.meshes) != 4 {
		t.Error("Invalid number of target mesh")
	}
	fi.AddCursor()
	if len(fi.meshes) != 5 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemVectorDeleteCursor(t *testing.T) {
	fi := testFormItemVector(t)
	fi.AddCursor()
	if len(fi.meshes) != 5 {
		t.Error("Invalid number of target mesh")
	}
	fi.DeleteCursor()
	if len(fi.meshes) != 4 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemVectorCharCallback(t *testing.T) {
	fi := testFormItemVector(t)
	fi.AddCursor()
	// start with 0
	fi.CharCallback('0', 0.1)
	if fi.values[0] != "0" {
		t.Errorf("Invalid value. Instead of '0', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "P0" {
		t.Errorf("Invalid typeState. Instead of 'P0', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('1', 0.1)
	if fi.values[0] != "0" {
		t.Errorf("Invalid value. Instead of '0', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "P0" {
		t.Errorf("Invalid typeState. Instead of 'P0', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('.', 0.1)
	if fi.values[0] != "0." {
		t.Errorf("Invalid value. Instead of '0.', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('2', 0.1)
	if fi.values[0] != "0.2" {
		t.Errorf("Invalid value. Instead of '0.2', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('8', 0.1)
	if fi.values[0] != "0.28" {
		t.Errorf("Invalid value. Instead of '0.28', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeStates[0])
	}
	// start with .
	fi = testFormItemVector(t)
	fi.AddCursor()
	fi.CharCallback('.', 0.1)
	if fi.values[0] != "" {
		t.Errorf("Invalid value. Instead of '', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('1', 0.1)
	if fi.values[0] != "1" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "PI" {
		t.Errorf("Invalid typeState. Instead of 'PI', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('2', 0.1)
	if fi.values[0] != "12" {
		t.Errorf("Invalid value. Instead of '12', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "PI" {
		t.Errorf("Invalid typeState. Instead of 'PI', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('.', 0.1)
	if fi.values[0] != "12." {
		t.Errorf("Invalid value. Instead of '12.', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('8', 0.1)
	if fi.values[0] != "12.8" {
		t.Errorf("Invalid value. Instead of '12.8', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeStates[0])
	}
	// start with -
	fi = testFormItemVector(t)
	fi.AddCursor()
	fi.CharCallback('-', 0.1)
	if fi.values[0] != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('.', 0.1)
	if fi.values[0] != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('1', 0.1)
	if fi.values[0] != "-1" {
		t.Errorf("Invalid value. Instead of '-1', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "NI" {
		t.Errorf("Invalid typeState. Instead of 'NI', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('2', 0.1)
	if fi.values[0] != "-12" {
		t.Errorf("Invalid value. Instead of '-12', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "NI" {
		t.Errorf("Invalid typeState. Instead of 'NI', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('.', 0.1)
	if fi.values[0] != "-12." {
		t.Errorf("Invalid value. Instead of '-12.', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('8', 0.1)
	if fi.values[0] != "-12.8" {
		t.Errorf("Invalid value. Instead of '-12.8', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeStates[0])
	}
	// start with -0.
	fi = testFormItemVector(t)
	fi.AddCursor()
	fi.CharCallback('-', 0.1)
	if fi.values[0] != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('0', 0.1)
	if fi.values[0] != "-0" {
		t.Errorf("Invalid value. Instead of '-0', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "N0" {
		t.Errorf("Invalid typeState. Instead of 'N0', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('.', 0.1)
	if fi.values[0] != "-0." {
		t.Errorf("Invalid value. Instead of '-0.', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('8', 0.1)
	if fi.values[0] != "-0.8" {
		t.Errorf("Invalid value. Instead of '-0.8', we have '%s'.", fi.values[0])
	}
	if fi.typeStates[0] != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeStates[0])
	}
	// validator
	fi = testFormItemVector(t)
	fi.AddCursor()
	validator := func(i float32) bool {
		return i <= 1.0 && i >= 0.0
	}
	fi.SetValidator(validator)
	fi.CharCallback('-', 0.1)
	if fi.values[0] != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.values[0])
	}
	fi.CharCallback('0', 0.1)
	if fi.values[0] != "-0" {
		t.Errorf("Invalid value. Instead of '-0', we have '%s'.", fi.values[0])
	}
	fi.CharCallback('.', 0.1)
	if fi.values[0] != "-0." {
		t.Errorf("Invalid value. Instead of '-0.', we have '%s'.", fi.values[0])
	}
	fi.CharCallback('3', 0.1)
	if fi.values[0] != "-0." {
		t.Errorf("Invalid value. Instead of '-0.', we have '%s'.", fi.values[0])
	}
	fi.SetTarget(1)
	fi.CharCallback('1', 0.1)
	if fi.values[1] != "1" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.values[0])
	}
	fi.CharCallback('.', 0.1)
	if fi.values[1] != "1." {
		t.Errorf("Invalid value. Instead of '1.', we have '%s'.", fi.values[0])
	}
	fi.CharCallback('3', 0.1)
	if fi.values[1] != "1." {
		t.Errorf("Invalid value. Instead of '1.', we have '%s'.", fi.values[1])
	}
	fi.CharCallback('0', 0.1)
	if fi.values[1] != "1.0" {
		t.Errorf("Invalid value. Instead of '1.0', we have '%s'.", fi.values[1])
	}
}
func TestFormItemVectorValueToString(t *testing.T) {
	fi := testFormItemVector(t)
	fi.AddCursor()
	fi.values[0] = "-"
	if fi.ValueToString() != "-" {
		t.Errorf("Invalid valuestring. instead of '-', we have '%s'.", fi.ValueToString())
	}
	fi.values[0] = "-3"
	if fi.ValueToString() != "-3" {
		t.Errorf("Invalid valuestring. instead of '-3', we have '%s'.", fi.ValueToString())
	}
	fi.values[0] = "-3."
	if fi.ValueToString() != "-3." {
		t.Errorf("Invalid valuestring. instead of '-3.', we have '%s'.", fi.ValueToString())
	}
	fi.values[0] = "-3.3"
	if fi.ValueToString() != "-3.3" {
		t.Errorf("Invalid valuestring. instead of '-3.3', we have '%s'.", fi.ValueToString())
	}
}
func TestFormItemVectorDeleteLastCharacter(t *testing.T) {
	fi := testFormItemVector(t)
	fi.AddCursor()
	fi.CharCallback('-', 0.1)
	fi.CharCallback('3', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('3', 0.1)
	if fi.ValueToString() != "-3.3" {
		t.Errorf("Invalid valuestring. instead of '-3.3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-3." {
		t.Errorf("Invalid valuestring. instead of '-3.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-3" {
		t.Errorf("Invalid valuestring. instead of '-3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "NI" {
		t.Errorf("Invalid typeState. Instead of 'NI', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-" {
		t.Errorf("Invalid valuestring. instead of '-', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('-', 0.1)
	fi.CharCallback('0', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('2', 0.1)
	if fi.ValueToString() != "-0.2" {
		t.Errorf("Invalid valuestring. instead of '-0.2', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-0." {
		t.Errorf("Invalid valuestring. instead of '-0.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-0" {
		t.Errorf("Invalid valuestring. instead of '-0', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "N0" {
		t.Errorf("Invalid typeState. Instead of 'N0', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-" {
		t.Errorf("Invalid valuestring. instead of '-', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('3', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('3', 0.1)
	if fi.ValueToString() != "3.3" {
		t.Errorf("Invalid valuestring. instead of '3.3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "3." {
		t.Errorf("Invalid valuestring. instead of '3.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "3" {
		t.Errorf("Invalid valuestring. instead of '3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "PI" {
		t.Errorf("Invalid typeState. Instead of 'PI', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeStates[0])
	}
	fi.CharCallback('0', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('3', 0.1)
	if fi.ValueToString() != "0.3" {
		t.Errorf("Invalid valuestring. instead of '0.3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "0." {
		t.Errorf("Invalid valuestring. instead of '0.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "0" {
		t.Errorf("Invalid valuestring. instead of '0', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "P0" {
		t.Errorf("Invalid typeState. Instead of 'P0', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeStates[0])
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeStates[0] != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeStates[0])
	}
}
