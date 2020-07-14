package model

import (
	"fmt"
	"strconv"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

type FormItemFloat struct {
	*FormItemCharBase
	typeState string
}

// GetValue returns the value of the form item. If the value can't parse as float32,
// some error message is printed out to the console, and the returned value is 0.0
func (fi *FormItemFloat) GetValue() float32 {
	valueFloat, err := strconv.ParseFloat(fi.value, 32)
	if err != nil {
		fmt.Printf("Can't format to float: '%s', err: '%s'\n", fi.value, err.Error())
		return 0.0
	}
	return float32(valueFloat)
}

// NewFormItemFloat returns a form item that maintains a float32 value.
func NewFormItemFloat(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemFloat {
	base := NewFormItemCharBase(label, mat, position, wrapper)
	return &FormItemFloat{
		FormItemCharBase: base,
		typeState:        "P",
	}
}

func (fi *FormItemFloat) pushState(r rune) {
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
	case "P0":
		fi.typeState = "P."
		break
	case "P.":
		fi.typeState = "PF"
		break
	case "PI":
		if r == rune('.') {
			fi.typeState = "P."
		}
		break
	case "N":
		if r == rune('0') {
			fi.typeState = "N0"
		} else {
			fi.typeState = "NI"
		}
		break
	case "N0":
		fi.typeState = "N."
		break
	case "N.":
		fi.typeState = "NF"
		break
	case "NI":
		if r == rune('.') {
			fi.typeState = "N."
		}
		break
	}
}
func (fi *FormItemFloat) popState(r rune) {
	switch fi.typeState {
	case "NF":
		if r == rune('.') {
			fi.typeState = "N."
		}
		break
	case "N.":
		if len(fi.value) == 2 && fi.value[1] == '0' {
			fi.typeState = "N0"
		} else {
			fi.typeState = "NI"
		}
		break
	case "N0":
		fi.typeState = "N"
		break
	case "NI":
		if len(fi.value) == 1 {
			fi.typeState = "N"
		}
		break
	case "N":
		fi.typeState = "P"
		break
	case "PF":
		if r == rune('.') {
			fi.typeState = "P."
		}
		break
	case "P.":
		if len(fi.value) == 1 && fi.value[0] == '0' {
			fi.typeState = "P0"
		} else {
			fi.typeState = "PI"
		}
		break
	case "P0":
		fi.typeState = "P"
		break
	case "PI":
		if len(fi.value) == 0 {
			fi.typeState = "P"
		}
		break
	}
}
func (fi *FormItemFloat) validRune(r rune) bool {
	var validRunes []rune
	switch fi.typeState {
	case "P":
		validRunes = []rune("0123456789-")
		break
	case "P0", "N0":
		validRunes = []rune(".")
		break
	case "P.", "N.", "PF", "N", "NF":
		validRunes = []rune("0123456789")
		break
	case "PI", "NI":
		validRunes = []rune("0123456789.")
		break
	}
	for _, v := range validRunes {
		if v == r {
			return true
		}
	}
	return false
}

// CharCallback validates the input character and appends it to the value if valid.
func (fi *FormItemFloat) CharCallback(r rune, offsetX float32) {
	if !fi.validRune(r) || len(fi.value) > fi.maxLen {
		return
	}
	fi.value = fi.value + string(r)
	fi.MoveCursorWithOffset(offsetX)
	fi.pushState(r)
}

// DeleteLastCharacter removes the last typed character from the form item.
func (fi *FormItemFloat) DeleteLastCharacter() {
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
