package model

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

type FormItemVector struct {
	*FormItemBase
	cursors       []interfaces.Mesh
	charOffsets   [][]float32
	values        []string
	maxLen        int
	currentTarget int // the index of the currently edited mesh
	validator     FloatValidator
	typeStates    []string
}

func NewFormItemVector(maxWidth, widthRatio, aspect float32, label, description string, inputMaxLen int, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemVector {
	m := NewFormItemBase(maxWidth, widthRatio, aspect, label, description, mat, wrapper)
	m.GetSurface().SetPosition(position)

	var writableTexture texture.Textures
	writableTexture.AddTexture(baseDirModel()+"/assets/paper.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", wrapper)
	writableTexture.AddTexture(baseDirModel()+"/assets/paper.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", wrapper)
	writablePrimitive := rectangle.NewExact(m.GetVectorTargetWidth(), m.GetTargetHeight())

	cursorPrimitive := rectangle.NewExact(m.GetCursorWidth(), m.GetCursorHeight())
	var cursorTex texture.Textures
	cursorTex.TransparentTexture(1, 1, 255, "tex.diffuse", wrapper)
	cursorTex.TransparentTexture(1, 1, 255, "tex.specular", wrapper)

	var cursors []interfaces.Mesh
	var offsets [][]float32
	var values []string
	var typeState []string
	for index := 0; index < 3; index++ {
		v, i, bo := writablePrimitive.MeshInput()
		writableMesh := mesh.NewTexturedMaterialMesh(v, i, writableTexture, mat, wrapper)
		writableMesh.SetParent(m.GetSurface())
		writableMesh.SetPosition(m.GetVectorTargetPosition(index))
		writableMesh.SetBoundingObject(bo)
		m.AddMesh(writableMesh)

		v, i, _ = cursorPrimitive.MeshInput()
		cursor := mesh.NewTexturedMaterialMesh(v, i, cursorTex, material.Greenplastic, wrapper)
		cursor.SetPosition(m.GetVectorCursorInitialPosition())
		cursor.SetParent(writableMesh)
		cursors = append(cursors, cursor)
		offsets = append(offsets, []float32{})
		values = append(values, "")
		typeState = append(typeState, "P")
	}
	return &FormItemVector{
		FormItemBase: m,
		cursors:      cursors,
		charOffsets:  offsets,
		values:       values,
		maxLen:       inputMaxLen,
		validator:    nil,
		typeStates:   typeState,
	}
}
func (fi *FormItemVector) cursorOffsetX() float32 {
	result := float32(0.0)
	if len(fi.charOffsets[fi.currentTarget]) > 0 {
		for i := 0; i < len(fi.charOffsets[fi.currentTarget]); i++ {
			result = result + fi.charOffsets[fi.currentTarget][i]
		}
	}
	return result
}
func (fi *FormItemVector) GetIndex(m interfaces.Mesh) int {
	for i, _ := range fi.meshes {
		if reflect.DeepEqual(m, reflect.ValueOf(fi.meshes[i]).Interface()) {
			return i
		}
	}
	return -1
}

// GetTarget returns the input target Mesh
func (fi *FormItemVector) GetTarget() interfaces.Mesh {
	return fi.meshes[fi.currentTarget+1]
}

// SetTarget updates the currentTarget index
func (fi *FormItemVector) SetTarget(i int) {
	fi.currentTarget = i
}

// AddCursor displays a cursor on the target surface.
func (fi *FormItemVector) AddCursor() {
	fi.AddMesh(fi.cursors[fi.currentTarget])
	fi.cursors[fi.currentTarget].SetPosition(mgl32.Vec3{fi.GetVectorCursorInitialPosition().X() - fi.cursorOffsetX(), 0.0, -0.01})
}

// DeleteCursor removes the cursor from the meshes.
func (fi *FormItemVector) DeleteCursor() {
	if len(fi.meshes) == 5 {
		fi.meshes = fi.meshes[:len(fi.meshes)-1]
	}
}

// MoveCursorWithOffset moves to cursor with the offset.
// It adds the new offset to the offsets, increments the sum offset
// and sets the cursor position.
func (fi *FormItemVector) MoveCursorWithOffset(offsetX float32) {
	fi.charOffsets[fi.currentTarget] = append(fi.charOffsets[fi.currentTarget], offsetX)
	fi.cursors[fi.currentTarget].SetPosition(mgl32.Vec3{fi.GetVectorCursorInitialPosition().X() - fi.cursorOffsetX(), 0.0, -0.01})
}

// StepBackCursor moves the cursor back after a character deletion.
func (fi *FormItemVector) StepBackCursor() {
	fi.charOffsets[fi.currentTarget] = fi.charOffsets[fi.currentTarget][:len(fi.charOffsets[fi.currentTarget])-1]
	fi.cursors[fi.currentTarget].SetPosition(mgl32.Vec3{fi.GetVectorCursorInitialPosition().X() - fi.cursorOffsetX(), 0.0, -0.01})
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemVector) ValueToString() string {
	return fi.values[fi.currentTarget]
}

// SetValidator sets the validator function
func (fi *FormItemVector) SetValidator(validator FloatValidator) {
	fi.validator = validator
}

// getValue returns the value of the current target form item. If the value can't parse as float32,
// some error message is printed out to the console, and the returned value is 0.0
func (fi *FormItemVector) getValue() float32 {
	valueFloat, err := strconv.ParseFloat(fi.values[fi.currentTarget], 32)
	if err != nil {
		fmt.Printf("Can't format to float: '%s', err: '%s'\n", fi.values[fi.currentTarget], err.Error())
		return 0.0
	}
	return float32(valueFloat)
}

// GetValue returns the value of the form item. If a component can't be parsed as float32,
// some error message is printed out to the console, and '0.0' is used instead of the
// wrong component.
func (fi *FormItemVector) GetValue() mgl32.Vec3 {
	var components [3]float32
	for i := 0; i < 3; i++ {
		valueFloat, err := strconv.ParseFloat(fi.values[i], 32)
		if err != nil {
			fmt.Printf("Can't format to float: '%s', err: '%s'\n", fi.values[i], err.Error())
			valueFloat = 0.0
		}
		components[i] = float32(valueFloat)
	}
	return mgl32.Vec3{components[0], components[1], components[2]}
}

// CharCallback validates the input character and appends it to the value if valid.
func (fi *FormItemVector) CharCallback(r rune, offsetX float32) {
	if !fi.validRune(r) || len(fi.values[fi.currentTarget]) >= fi.maxLen {
		return
	}
	fi.values[fi.currentTarget] = fi.values[fi.currentTarget] + string(r)
	fi.MoveCursorWithOffset(offsetX)
	fi.pushState(r)
	if fi.validator != nil {
		if !fi.validator(fi.getValue()) {
			fi.DeleteLastCharacter()
		}
	}
}

// DeleteLastCharacter removes the last typed character from the form item.
func (fi *FormItemVector) DeleteLastCharacter() {
	if len(fi.charOffsets[fi.currentTarget]) == 0 {
		return
	}
	fi.values[fi.currentTarget] = fi.values[fi.currentTarget][:len(fi.values[fi.currentTarget])-1]
	fi.StepBackCursor()
	if len(fi.values[fi.currentTarget]) > 0 {
		fi.popState(rune(fi.values[fi.currentTarget][len(fi.values[fi.currentTarget])-1]))
	} else {
		// dummy value for pop state.
		fi.popState('.')
	}
}
func (fi *FormItemVector) pushState(r rune) {
	switch fi.typeStates[fi.currentTarget] {
	case "P":
		if r == rune('-') {
			fi.typeStates[fi.currentTarget] = "N"
		} else if r == rune('0') {
			fi.typeStates[fi.currentTarget] = "P0"
		} else {
			fi.typeStates[fi.currentTarget] = "PI"
		}
		break
	case "P0":
		fi.typeStates[fi.currentTarget] = "P."
		break
	case "P.":
		fi.typeStates[fi.currentTarget] = "PF"
		break
	case "PI":
		if r == rune('.') {
			fi.typeStates[fi.currentTarget] = "P."
		}
		break
	case "N":
		if r == rune('0') {
			fi.typeStates[fi.currentTarget] = "N0"
		} else {
			fi.typeStates[fi.currentTarget] = "NI"
		}
		break
	case "N0":
		fi.typeStates[fi.currentTarget] = "N."
		break
	case "N.":
		fi.typeStates[fi.currentTarget] = "NF"
		break
	case "NI":
		if r == rune('.') {
			fi.typeStates[fi.currentTarget] = "N."
		}
		break
	}
}
func (fi *FormItemVector) popState(r rune) {
	switch fi.typeStates[fi.currentTarget] {
	case "NF":
		if r == rune('.') {
			fi.typeStates[fi.currentTarget] = "N."
		}
		break
	case "N.":
		if len(fi.values[fi.currentTarget]) == 2 && fi.values[fi.currentTarget][1] == '0' {
			fi.typeStates[fi.currentTarget] = "N0"
		} else {
			fi.typeStates[fi.currentTarget] = "NI"
		}
		break
	case "N0":
		fi.typeStates[fi.currentTarget] = "N"
		break
	case "NI":
		if len(fi.values[fi.currentTarget]) == 1 {
			fi.typeStates[fi.currentTarget] = "N"
		}
		break
	case "N":
		fi.typeStates[fi.currentTarget] = "P"
		break
	case "PF":
		if r == rune('.') {
			fi.typeStates[fi.currentTarget] = "P."
		}
		break
	case "P.":
		if len(fi.values[fi.currentTarget]) == 1 && fi.values[fi.currentTarget][0] == '0' {
			fi.typeStates[fi.currentTarget] = "P0"
		} else {
			fi.typeStates[fi.currentTarget] = "PI"
		}
		break
	case "P0":
		fi.typeStates[fi.currentTarget] = "P"
		break
	case "PI":
		if len(fi.values[fi.currentTarget]) == 0 {
			fi.typeStates[fi.currentTarget] = "P"
		}
		break
	}
}
func (fi *FormItemVector) validRune(r rune) bool {
	var validRunes []rune
	switch fi.typeStates[fi.currentTarget] {
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
