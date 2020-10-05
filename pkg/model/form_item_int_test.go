package model

import (
	"math"
	"strconv"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

func testFormItemInt(t *testing.T) *FormItemInt {
	mat := material.Chrome
	pos := mgl32.Vec3{0, 0, 0}
	fi := NewFormItemInt(DefaultMaxWidth, ITEM_WIDTH_HALF, 1.0, DefaultFormItemLabel, DefaultFormItemDescription, mat, pos, wrapperMock)

	if fi.label != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.label)
	}
	if fi.description != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.description)
	}
	return fi
}
func TestNewFormItemInt(t *testing.T) {
	_ = testFormItemInt(t)
}
func TestFormItemIntGetLabel(t *testing.T) {
	fi := testFormItemInt(t)
	if fi.GetLabel() != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.GetLabel())
	}
}
func TestFormItemIntGetDescription(t *testing.T) {
	fi := testFormItemInt(t)
	if fi.GetDescription() != DefaultFormItemDescription {
		t.Errorf("Invalid form item description. Instead of '%s', we have '%s'.", DefaultFormItemDescription, fi.GetDescription())
	}
}
func TestFormItemIntGetValue(t *testing.T) {
	fi := testFormItemInt(t)
	val := 3
	fi.value = "3"
	if fi.GetValue() != val {
		t.Errorf("Invalid form item value. Instead of '%d', it is '%d'.", val, fi.GetValue())
	}
}
func TestFormItemIntGetSurface(t *testing.T) {
	fi := testFormItemInt(t)
	if fi.GetSurface() != fi.meshes[0] {
		t.Error("Invalid surface mesh")
	}
}
func TestFormItemIntGetTarget(t *testing.T) {
	fi := testFormItemInt(t)
	if fi.GetTarget() != fi.meshes[1] {
		t.Error("Invalid target mesh")
	}
}
func TestFormItemIntAddCursor(t *testing.T) {
	fi := testFormItemInt(t)
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemIntDeleteCursor(t *testing.T) {
	fi := testFormItemInt(t)
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
	fi.DeleteCursor()
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemIntCharCallback(t *testing.T) {
	fi := testFormItemInt(t)
	fi.AddCursor()
	fi.DeleteLastCharacter()
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
	fi = testFormItemInt(t)
	fi.AddCursor()
	fi.CharCallback('1', 0.1)
	if fi.value != "1" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.value)
	}
	fi.value = strconv.Itoa(math.MaxInt32 - 10)
	fi.CharCallback('1', 0.1)
	if fi.value != strconv.Itoa(math.MaxInt32-10) {
		t.Errorf("Invalid value. Instead of '%d', we have '%s'.", math.MaxInt32-10, fi.value)
	}
	fi = testFormItemInt(t)
	fi.AddCursor()
	validator := func(i int) bool {
		return i <= 100 && i >= -10
	}
	fi.SetValidator(validator)
	fi.CharCallback('1', 0.1)
	if fi.value != "1" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.value)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "11" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.value)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "11" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	fi.DeleteLastCharacter()
	fi.CharCallback('-', 0.1)
	fi.CharCallback('1', 0.1)
	if fi.value != "-1" {
		t.Errorf("Invalid value. Instead of '-1', we have '%s'.", fi.value)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "-1" {
		t.Errorf("Invalid value. Instead of '-1', we have '%s'.", fi.value)
	}
	fi.CharCallback('0', 0.1)
	if fi.value != "-10" {
		t.Errorf("Invalid value. Instead of '-10', we have '%s'.", fi.value)
	}
	fi.CharCallback('0', 0.1)
	if fi.value != "-10" {
		t.Errorf("Invalid value. Instead of '-10', we have '%s'.", fi.value)
	}
}
func TestFormItemIntValueToString(t *testing.T) {
	fi := testFormItemInt(t)
	val := fi.ValueToString()
	if val != "" {
		t.Errorf("Invalid value. Instead of '', we have '%s'.", val)
	}
	fi.value = strconv.Itoa(3)
	val = fi.ValueToString()
	if val != "3" {
		t.Errorf("Invalid value. Instead of '3', we have '%s'.", val)
	}
}
func TestFormItemIntDeleteLastCharacter(t *testing.T) {
	fi := testFormItemInt(t)
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
	fi.value = "-1234"
	fi.typeState = "NI"
	fi.charOffsets = []float32{0.1, 0.1, 0.1, 0.1, 0.1}
	fi.DeleteLastCharacter()
	if fi.value != "-123" {
		t.Errorf("Invalid value. Instead of '-123', we have '%s'.", fi.value)
	}
	fi.DeleteLastCharacter()
	if fi.value != "-12" {
		t.Errorf("Invalid value. Instead of '-12', we have '%s'.", fi.value)
	}
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
