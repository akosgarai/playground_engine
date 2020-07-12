package model

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

type FormItemFloat struct {
	*BaseModel
	cursor        interfaces.Mesh
	cursorOffsetX float32
	charOffsets   []float32
	label         string
	value         string
	typeState     string
}

// GetLabel returns the label string of the item.
func (fi *FormItemFloat) GetLabel() string {
	return fi.label
}

// GetValue returns the value of the form item.
func (fi *FormItemFloat) GetValue() float32 {
	valueFloat, err := strconv.ParseFloat(fi.value, 32)
	if err != nil {
		fmt.Printf("Can't format to float: '%s', err: '%s'\n", fi.value, err.Error())
		return 0.0
	}
	return float32(valueFloat)
}

// SetValue sets the value of the form item.
func (fi *FormItemFloat) SetValue(v float32) {
	valueString := fmt.Sprintf("%f", v)
	parts := strings.Split(valueString, ".")
	partInt, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("Can't format to int: '%s', err: '%s'\n", parts[0], err.Error())
		return
	}
	partFloatString := strings.TrimSuffix("."+parts[1], "0")
	if partFloatString == "." {
		partFloatString = "0"
	}
	// further validation: max 9 char included the '.'
	if len(parts[0]) > 9 {
		fmt.Printf("Can't set value, int part length '%d' > 9\n", len(parts[0]))
		return
	}
	if len(parts[0])+1+int(math.Min(1, float64(len(partFloatString)))) > 9 {
		fmt.Printf("Can't set this precision, int part length '%d' + 1 + '%d' > 9\n", len(parts[0]), int(math.Min(1, float64(len(partFloatString)))))
		return
	}
	to := 9 - (len(parts[0]) + 1)
	if to > len(partFloatString) {
		to = len(partFloatString)
	}
	partFloat, err := strconv.Atoi(partFloatString[1:to])
	if err != nil {
		fmt.Printf("Can't format to int: '%s' (orig: '%s'), err: '%s'\n", partFloatString[0:to], parts[1], err.Error())
		return
	}
	fi.value = fmt.Sprintf("%d.%d", partInt, partFloat)
}
func NewFormItemFloat(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemFloat {
	labelPrimitive := rectangle.NewExact(FormItemWidth, FormItemLength)
	v, i, bo := labelPrimitive.MeshInput()
	var tex texture.Textures
	tex.TransparentTexture(1, 1, 128, "tex.diffuse", wrapper)
	tex.TransparentTexture(1, 1, 128, "tex.specular", wrapper)
	formItemMesh := mesh.NewTexturedMaterialMesh(v, i, tex, mat, wrapper)
	formItemMesh.SetBoundingObject(bo)
	formItemMesh.SetPosition(position)
	m := New()
	m.AddMesh(formItemMesh)
	var writableTexture texture.Textures
	writableTexture.AddTexture(baseDirModel()+"/assets/paper.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", wrapper)
	writableTexture.AddTexture(baseDirModel()+"/assets/paper.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", wrapper)
	writablePrimitive := rectangle.NewExact(writableWidth, writableHeight)
	v, i, bo = writablePrimitive.MeshInput()
	writableMesh := mesh.NewTexturedMaterialMesh(v, i, writableTexture, mat, wrapper)
	writableMesh.SetParent(formItemMesh)
	writableMesh.SetPosition(mgl32.Vec3{0.24, -0.01, 0.0})
	writableMesh.SetBoundingObject(bo)
	m.AddMesh(writableMesh)
	cursorPrimitive := rectangle.NewExact(cursorWidth, cursorHeight)
	var ctex texture.Textures
	ctex.TransparentTexture(1, 1, 255, "tex.diffuse", wrapper)
	ctex.TransparentTexture(1, 1, 255, "tex.specular", wrapper)
	v, i, _ = cursorPrimitive.MeshInput()
	cursor := mesh.NewTexturedMaterialMesh(v, i, ctex, material.Greenplastic, wrapper)
	cursor.SetPosition(mgl32.Vec3{CursorInitX, 0.0, -0.01})
	cursor.SetParent(writableMesh)
	return &FormItemFloat{
		BaseModel:     m,
		label:         label,
		cursor:        cursor,
		cursorOffsetX: 0.0,
		charOffsets:   []float32{},
		value:         "",
		typeState:     "P",
	}
}

// GetSurface returns the formItemMesh
func (fi *FormItemFloat) GetSurface() interfaces.Mesh {
	return fi.meshes[0]
}

// GetTarget returns the input target Mesh
func (fi *FormItemFloat) GetTarget() interfaces.Mesh {
	return fi.meshes[1]
}

// AddCursor displays a cursor on the target surface.
func (fi *FormItemFloat) AddCursor() {
	fi.AddMesh(fi.cursor)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}

// DeleteCursor removes the cursor from the meshes.
func (fi *FormItemFloat) DeleteCursor() {
	if len(fi.meshes) == 3 {
		fi.meshes = fi.meshes[:len(fi.meshes)-1]
	}
}
func (fi *FormItemFloat) CharCallback(r rune, offsetX float32) {
	if !fi.validRune(r) {
		return
	}
	fi.value = fi.value + string(r)
	fi.cursorOffsetX = fi.cursorOffsetX + offsetX
	fi.charOffsets = append(fi.charOffsets, offsetX)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
	fi.pushState(r)
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
		if len(fi.value) == 2 && fi.value[1] == '0' {
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

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemFloat) ValueToString() string {
	return fi.value
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemFloat) DeleteLastCharacter() {
	if len(fi.charOffsets) == 0 {
		return
	}
	lastChar := fi.value[len(fi.value)-1]
	fi.value = fi.value[:len(fi.value)-1]
	offsetX := fi.charOffsets[len(fi.charOffsets)-1]
	fi.cursorOffsetX = fi.cursorOffsetX - offsetX
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
	fi.charOffsets = fi.charOffsets[:len(fi.charOffsets)-1]
	fi.popState(rune(lastChar))
}
