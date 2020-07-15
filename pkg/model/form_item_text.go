package model

import (
	"strings"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

type FormItemText struct {
	*FormItemCharBase
}

// GetValue returns the value of the form item.
func (fi *FormItemText) GetValue() string {
	return fi.value
}

// NewFormItemText returns a form item that maintains a string value.
func NewFormItemText(maxWidth float32, label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemText {
	base := NewFormItemCharBase(maxWidth, ITEM_WIDTH_HALF, label, mat, position, wrapper)
	return &FormItemText{
		FormItemCharBase: base,
	}
}
func (fi *FormItemText) validRune(r rune) bool {
	return true
}

// CharCallback validates the input character and appends it to the value if valid.
func (fi *FormItemText) CharCallback(r rune, offsetX float32) {
	if !fi.validRune(r) || len(fi.value)+strings.Count(fi.value, " ") > fi.maxLen {
		return
	}
	fi.value = fi.value + string(r)
	fi.MoveCursorWithOffset(offsetX)
}

// DeleteLastCharacter removes the last typed character from the form item.
func (fi *FormItemText) DeleteLastCharacter() {
	if len(fi.charOffsets) == 0 {
		return
	}
	fi.value = fi.value[:len(fi.value)-1]
	fi.StepBackCursor()
}
