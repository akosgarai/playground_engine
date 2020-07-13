package model

import (
	"strconv"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

type FormItemInt64 struct {
	*FormItemCharBase
	typeState string
}

// GetValue returns the value of the form item.
func (fi *FormItemInt64) GetValue() int64 {
	val, _ := strconv.Atoi(fi.value)
	return int64(val)
}

// SetValue returns the value of the form item.
func (fi *FormItemInt64) SetValue(v int64) {
	fi.value = strconv.Itoa(int(v))
}

// NewFormItemInt64 returns a form item that maintains an int64 value.
func NewFormItemInt64(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemInt64 {
	base := NewFormItemCharBase(label, mat, position, wrapper)
	return &FormItemInt64{
		FormItemCharBase: base,
		typeState:        "P",
	}
}

func (fi *FormItemInt64) validRune(r rune) bool {
	var validRunes []rune
	switch fi.typeState {
	case "P":
		validRunes = []rune("0123456789-")
		break

	case "N":
		validRunes = []rune("123456789")
		break
	case "PI", "NI":
		validRunes = []rune("0123456789")
		break
	}
	for _, v := range validRunes {
		if v == r {
			return true
		}
	}
	return false
}
func (fi *FormItemInt64) popState(r rune) {
	switch fi.typeState {
	case "N":
		fi.typeState = "P"
		break
	case "NI":
		if len(fi.value) == 1 {
			fi.typeState = "N"
		}
		break
	case "PI":
		if len(fi.value) == 0 {
			fi.typeState = "P"
		}
		break
	case "P0":
		fi.typeState = "P"
		break
	}
}
func (fi *FormItemInt64) pushState(r rune) {
	switch fi.typeState {
	case "P":
		if r == rune('-') {
			fi.typeState = "N"
		} else if r == rune('0') {
			fi.typeState = "P0"
		} else {
			fi.typeState = "PI"
		}
		break
	case "N":
		fi.typeState = "NI"
		break
	}
}

// CharCallback validates the input character and appends it to the value if valid.
func (fi *FormItemInt64) CharCallback(r rune, offsetX float32) {
	if !fi.validRune(r) || len(fi.value) > fi.maxLen {
		return
	}
	fi.value = fi.value + string(r)
	fi.MoveCursorWithOffset(offsetX)
	fi.pushState(r)
}

// DeleteLastCharacter removes the last typed character from the form item.
func (fi *FormItemInt64) DeleteLastCharacter() {
	if len(fi.charOffsets) == 0 {
		return
	}
	fi.value = fi.value[:len(fi.value)-1]
	fi.StepBackCursor()
	if len(fi.value) > 0 {
		fi.popState(rune(fi.value[len(fi.value)-1]))
	} else {
		// dummy value for pop state.
		fi.popState('.')
	}
}
