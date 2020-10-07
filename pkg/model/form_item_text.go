package model

import (
	"strings"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

type FormItemText struct {
	*FormItemCharBase
	validator StringValidator
}

// GetValue returns the value of the form item.
func (fi *FormItemText) GetValue() string {
	return fi.value
}

// NewFormItemText returns a form item that maintains a string value.
func NewFormItemText(maxWidth, itemWidth, aspect float32, label, description string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemText {
	base := NewFormItemCharBase(maxWidth, itemWidth, aspect, label, description, CHAR_NUM_TEXT, mat, position, wrapper)
	return &FormItemText{
		FormItemCharBase: base,
		validator:        nil,
	}
}

// SetValidator sets the validator function
func (fi *FormItemText) SetValidator(validator StringValidator) {
	fi.validator = validator
}
func (fi *FormItemText) validRune(r rune) bool {
	return true
}

// CharCallback validates the input character and appends it to the value if valid.
func (fi *FormItemText) CharCallback(r rune, offsetX float32) {
	if !fi.validRune(r) || len(fi.value)+strings.Count(fi.value, " ") >= fi.maxLen {
		return
	}
	fi.value = fi.value + string(r)
	fi.MoveCursorWithOffset(offsetX)
	if fi.validator != nil {
		if !fi.validator(fi.GetValue()) {
			fi.DeleteLastCharacter()
		}
	}
}

// DeleteLastCharacter removes the last typed character from the form item.
func (fi *FormItemText) DeleteLastCharacter() {
	if len(fi.charOffsets) == 0 {
		return
	}
	fi.value = fi.value[:len(fi.value)-1]
	fi.StepBackCursor()
}
