package model

import (
	"math"
	"strconv"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

type FormItemInt64 struct {
	*FormItemCharBase
	value      int64
	isNegative bool
}

// GetValue returns the value of the form item.
func (fi *FormItemInt64) GetValue() int64 {
	return fi.value
}

// SetValue returns the value of the form item.
func (fi *FormItemInt64) SetValue(v int64) {
	fi.value = v
}
func NewFormItemInt64(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemInt64 {
	base := NewFormItemCharBase(label, mat, position, wrapper)
	return &FormItemInt64{
		FormItemCharBase: base,
		value:            0,
		isNegative:       false,
	}
}

func (fi *FormItemInt64) CharCallback(r rune, offsetX float32) {
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
	fi.value = fi.value*10 + int64(val)
	fi.cursorOffsetX = fi.cursorOffsetX + offsetX
	fi.charOffsets = append(fi.charOffsets, offsetX)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}
func (fi *FormItemInt64) validRune(r rune) bool {
	// integer number isn't allowed to start with 0.
	if fi.value == 0 && r == rune('0') {
		return false
	}
	// if the next iteration the value will be grater than max or less than min
	// return false
	if fi.value > math.MaxInt64/10 || fi.value < math.MinInt64/10 {
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
func (fi *FormItemInt64) ValueToString() string {
	if fi.isNegative && fi.value == 0 {
		return "-"
	}
	return strconv.FormatInt(fi.value, 10)
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemInt64) DeleteLastCharacter() {
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
