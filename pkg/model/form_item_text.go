package model

import (
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

type FormItemText struct {
	*FormItemCharBase
	value string
}

// GetValue returns the value of the form item.
func (fi *FormItemText) GetValue() string {
	return fi.value
}

// SetValue returns the value of the form item.
func (fi *FormItemText) SetValue(v string) {
	fi.value = v
}
func NewFormItemText(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemText {
	base := NewFormItemCharBase(label, mat, position, wrapper)
	return &FormItemText{
		FormItemCharBase: base,
		value:            "",
	}
}
func (fi *FormItemText) CharCallback(r rune, offsetX float32) {
	if !fi.validRune(r) {
		return
	}
	fi.value = fi.value + string(r)
	fi.cursorOffsetX = fi.cursorOffsetX + offsetX
	fi.charOffsets = append(fi.charOffsets, offsetX)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}
func (fi *FormItemText) validRune(r rune) bool {
	return true
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemText) ValueToString() string {
	return fi.value
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemText) DeleteLastCharacter() {
	if len(fi.charOffsets) == 0 {
		return
	}
	fi.value = fi.value[:len(fi.value)-1]
	offsetX := fi.charOffsets[len(fi.charOffsets)-1]
	fi.cursorOffsetX = fi.cursorOffsetX - offsetX
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
	fi.charOffsets = fi.charOffsets[:len(fi.charOffsets)-1]
}
