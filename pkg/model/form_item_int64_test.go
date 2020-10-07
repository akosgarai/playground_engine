package model

import (
	"math"
	"strconv"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

func testFormItemInt64(t *testing.T) *FormItemInt64 {
	mat := material.Chrome
	pos := mgl32.Vec3{0, 0, 0}
	fi := NewFormItemInt64(DefaultMaxWidth, ITEM_WIDTH_HALF, 1.0, DefaultFormItemLabel, DefaultFormItemDescription, mat, pos, wrapperMock)

	if fi.label != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.label)
	}
	if fi.description != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.description)
	}
	return fi
}
func TestNewFormItemInt64(t *testing.T) {
	_ = testFormItemInt64(t)
}
func TestFormItemInt64GetLabel(t *testing.T) {
	fi := testFormItemInt64(t)
	if fi.GetLabel() != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.GetLabel())
	}
}
func TestFormItemInt64GetDescription(t *testing.T) {
	fi := testFormItemInt64(t)
	if fi.GetDescription() != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.GetDescription())
	}
}
func TestFormItemInt64GetValue(t *testing.T) {
	fi := testFormItemInt64(t)
	val := int64(3)
	fi.value = "3"
	if fi.GetValue() != val {
		t.Errorf("Invalid form item value. Instead of '%d', it is '%d'.", val, fi.GetValue())
	}
}
func TestFormItemInt64GetSurface(t *testing.T) {
	fi := testFormItemInt64(t)
	if fi.GetSurface() != fi.meshes[0] {
		t.Error("Invalid surface mesh")
	}
}
func TestFormItemInt64GetTarget(t *testing.T) {
	fi := testFormItemInt64(t)
	if fi.GetTarget() != fi.meshes[1] {
		t.Error("Invalid target mesh")
	}
}
func TestFormItemInt64AddCursor(t *testing.T) {
	fi := testFormItemInt64(t)
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemInt64DeleteCursor(t *testing.T) {
	fi := testFormItemInt64(t)
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
	fi.DeleteCursor()
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemInt64CharCallback(t *testing.T) {
	fi := testFormItemInt64(t)
	fi.AddCursor()
	fi.DeleteLastCharacter()
	if fi.value != "" {
		t.Errorf("Invalid value. Instead of '', we have '%s'.", fi.value)
	}
	fi.CharCallback('b', 0.1)
	if fi.value != "" {
		t.Errorf("Invalid value. Instead of '', we have '%s'.", fi.value)
	}
	fi.CharCallback('0', 0.1)
	if fi.value != "0" {
		t.Errorf("Invalid value. Instead of '0', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	fi.CharCallback('-', 0.1)
	if fi.value != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.value)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "-1" {
		t.Errorf("Invalid value. Instead of '-1', we have '%s'.", fi.value)
	}
	fi = testFormItemInt64(t)
	fi.AddCursor()
	fi.CharCallback('1', 0.1)
	if fi.value != "1" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.value)
	}
	fi.value = strconv.Itoa(math.MaxInt64 - 10)
	maxVal := strconv.Itoa(math.MaxInt64-10) + "1"
	fi.CharCallback('1', 0.1)
	if fi.value != maxVal {
		t.Errorf("Invalid value. Instead of '%s', we have '%s'.", maxVal, fi.value)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != maxVal {
		t.Errorf("Invalid value. Instead of '%s', we have '%s'.", maxVal, fi.value)
	}
	fi = testFormItemInt64(t)
	fi.AddCursor()
	validator := func(i int64) bool {
		return i >= -300 && i <= 200
	}
	fi.SetValidator(validator)
	fi.CharCallback('1', 0.1)
	if fi.value != "1" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.value)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "11" {
		t.Errorf("Invalid value. Instead of '11', we have '%s'.", fi.value)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "111" {
		t.Errorf("Invalid value. Instead of '111', we have '%s'.", fi.value)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "111" {
		t.Errorf("Invalid value. Instead of '111', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	fi.DeleteLastCharacter()
	fi.DeleteLastCharacter()
	fi.CharCallback('-', 0.1)
	fi.CharCallback('3', 0.1)
	if fi.value != "-3" {
		t.Errorf("Invalid value. Instead of '-3', we have '%s'.", fi.value)
	}
	fi.CharCallback('3', 0.1)
	if fi.value != "-33" {
		t.Errorf("Invalid value. Instead of '-33', we have '%s'.", fi.value)
	}
	fi.CharCallback('3', 0.1)
	if fi.value != "-33" {
		t.Errorf("Invalid value. Instead of '-33', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	fi.CharCallback('0', 0.1)
	if fi.value != "-30" {
		t.Errorf("Invalid value. Instead of '-30', we have '%s'.", fi.value)
	}
	fi.CharCallback('0', 0.1)
	if fi.value != "-300" {
		t.Errorf("Invalid value. Instead of '-300', we have '%s'.", fi.value)
	}
}
func TestFormItemInt64ValueToString(t *testing.T) {
	fi := testFormItemInt64(t)
	val := fi.ValueToString()
	if val != "" {
		t.Errorf("Invalid value. Instead of '', we have '%s'.", val)
	}
	fi.value = "3"
	val = fi.ValueToString()
	if val != "3" {
		t.Errorf("Invalid value. Instead of '3', we have '%s'.", val)
	}
}
func TestFormItemInt64DeleteLastCharacter(t *testing.T) {
	fi := testFormItemInt64(t)
	fi.value = "12345"
	fi.typeState = "PI"
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
	fi.value = "-12"
	fi.typeState = "NI"
	fi.charOffsets = []float32{0.1, 0.1, 0.1}
	fi.DeleteLastCharacter()
	if fi.value != "-1" {
		t.Errorf("Invalid value. Instead of '-1', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	if fi.value != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.value)
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
