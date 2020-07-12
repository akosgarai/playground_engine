package model

import (
	"math"
	"strconv"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

type FormItemInt struct {
	*FormItemCharBase
	value      int
	isNegative bool
}

// GetValue returns the value of the form item.
func (fi *FormItemInt) GetValue() int {
	return fi.value
}

// SetValue returns the value of the form item.
func (fi *FormItemInt) SetValue(v int) {
	fi.value = v
}
func NewFormItemInt(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemInt {
	base := NewFormItemCharBase(label, mat, position, wrapper)
	return &FormItemInt{
		FormItemCharBase: base,
		value:            0,
		isNegative:       false,
	}
}

func (fi *FormItemInt) CharCallback(r rune, offsetX float32) {
	// if the first character is '-', mark the form item as negative.
	if fi.value == 0 && r == rune('-') {
		fi.isNegative = true
		fi.cursorOffsetX = fi.cursorOffsetX + offsetX
		fi.charOffsets = append(fi.charOffsets, offsetX)
		fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
		return
	}
	if !fi.validRune(r) {
		return
	}
	val := int(r - '0')
	if fi.isNegative {
		val = -val
	}
	fi.value = fi.value*10 + val
	fi.cursorOffsetX = fi.cursorOffsetX + offsetX
	fi.charOffsets = append(fi.charOffsets, offsetX)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}
func (fi *FormItemInt) validRune(r rune) bool {
	// integer number isn't allowed to start with 0.
	if fi.value == 0 && r == rune('0') {
		return false
	}
	// if the next iteration the value will be grater than max or less than min
	// return false
	if fi.value > math.MaxInt32/10 || fi.value < math.MinInt32/10 {
		return false
	}
	validRunes := []rune("0123456789")
	for _, v := range validRunes {
		if v == r {
			return true
		}
	}
	return false
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemInt) ValueToString() string {
	if fi.isNegative && fi.value == 0 {
		return "-"
	}
	return strconv.Itoa(fi.value)
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemInt) DeleteLastCharacter() {
	if fi.value == 0 {
		if fi.isNegative {
			fi.isNegative = false
		} else {
			return
		}
	} else {
		mod := fi.value % 10
		fi.value = (fi.value - mod) / 10
	}
	offsetX := fi.charOffsets[len(fi.charOffsets)-1]
	fi.cursorOffsetX = fi.cursorOffsetX - offsetX
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
	fi.charOffsets = fi.charOffsets[:len(fi.charOffsets)-1]
}
