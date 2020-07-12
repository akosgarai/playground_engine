package model

import (
	"strconv"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

type FormItemInt struct {
	*FormItemCharBase
	typeState string
}

// GetValue returns the value of the form item.
func (fi *FormItemInt) GetValue() int {
	val, _ := strconv.Atoi(fi.value)
	return val
}

// SetValue returns the value of the form item.
func (fi *FormItemInt) SetValue(v int) {
	fi.value = strconv.Itoa(v)
}
func NewFormItemInt(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemInt {
	base := NewFormItemCharBase(label, mat, position, wrapper)
	return &FormItemInt{
		FormItemCharBase: base,
		typeState:        "P",
	}
}

func (fi *FormItemInt) validRune(r rune) bool {
	var validRunes []rune
	switch fi.typeState {
	case "P":
		validRunes = []rune("123456789-")
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
func (fi *FormItemInt) popState(r rune) {
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
	}
}
func (fi *FormItemInt) pushState(r rune) {
	switch fi.typeState {
	case "P":
		if r == rune('-') {
			fi.typeState = "N"
		} else {
			fi.typeState = "PI"
		}
		break
	case "N":
		fi.typeState = "NI"
		break
	}
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemInt) ValueToString() string {
	return fi.value
}

// CharCallback validates the input character and appends it to the value if valid.
func (fi *FormItemInt) CharCallback(r rune, offsetX float32) {
	if !fi.validRune(r) || len(fi.value) > fi.maxLen {
		return
	}
	fi.value = fi.value + string(r)
	fi.MoveCursorWithOffset(offsetX)
	fi.pushState(r)
}

// DeleteLastCharacter removes the last typed character from the form item.
func (fi *FormItemInt) DeleteLastCharacter() {
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
