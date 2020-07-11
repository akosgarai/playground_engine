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
	valueInt      int
	valueFloat    int
	floatPosition int
	isNegative    bool
}

// GetLabel returns the label string of the item.
func (fi *FormItemFloat) GetLabel() string {
	return fi.label
}

// GetValue returns the value of the form item.
func (fi *FormItemFloat) GetValue() float32 {
	valueString := fmt.Sprintf("%d.%d", fi.valueInt, fi.valueFloat)
	valueFloat, err := strconv.ParseFloat(valueString, 32)
	if err != nil {
		fmt.Printf("Can't format to float: '%s', err: '%s'\n", valueString, err.Error())
		return 0.0
	}
	return float32(valueFloat)
}

// SetValue returns the value of the form item.
func (fi *FormItemFloat) SetValue(v float32) {
	valueString := fmt.Sprintf("%f", v)
	parts := strings.Split(valueString, ".")
	partInt, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Printf("Can't format to int: '%s', err: '%s'\n", parts[0], err.Error())
		return
	}
	partFloatString := strings.TrimSuffix(parts[1], "0")
	if partFloatString == "" {
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
	fi.valueInt = partInt
	fi.floatPosition = len(parts[0])
	to := 9 - (len(parts[0]) + 1)
	partFloat, err := strconv.Atoi(partFloatString[0:to])
	if err != nil {
		fmt.Printf("Can't format to int: '%s' (orig: '%s'), err: '%s'\n", partFloatString[0:to], parts[1], err.Error())
		return
	}
	fi.valueFloat = partFloat
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
		valueInt:      0,
		valueFloat:    0,
		floatPosition: -1,
		isNegative:    false,
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
	// if the first character is '-', mark the form item as negative.
	if fi.valueInt == 0 && fi.floatPosition == -1 && r == rune('-') {
		fi.isNegative = true
		fi.cursorOffsetX = fi.cursorOffsetX + offsetX
		fi.charOffsets = append(fi.charOffsets, offsetX)
		fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
		return
	}
	// handle '.' here
	if r == rune('.') {
		fi.floatPosition = len(fmt.Sprintf("%d", fi.valueInt)) + 1
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
	if fi.floatPosition > -1 {
		fi.valueFloat = fi.valueFloat*10 + val
	} else {
		fi.valueInt = fi.valueInt*10 + val
	}
	fi.cursorOffsetX = fi.cursorOffsetX + offsetX
	fi.charOffsets = append(fi.charOffsets, offsetX)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}
func (fi *FormItemFloat) validRune(r rune) bool {
	// integer number isn't allowed to start with 0.
	if fi.valueInt == 0 && fi.floatPosition == -1 && r == rune('0') {
		return false
	}
	// if the next iteration the value will be grater than max or less than min
	// return false
	if fi.floatPosition == -1 {
		if fi.valueInt > math.MaxInt32/10 || fi.valueInt < math.MinInt32/10 {
			return false
		}
	} else {
		if fi.valueFloat > math.MaxInt32/10 || fi.valueFloat < math.MinInt32/10 {
			return false
		}
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
func (fi *FormItemFloat) ValueToString() string {
	if fi.isNegative && fi.valueInt == 0 && fi.valueFloat == 0 {
		return "-"
	}
	if fi.valueFloat > 0 {
		return fmt.Sprintf("%d.%d", fi.valueInt, fi.valueFloat)
	}
	return fmt.Sprintf("%d", fi.valueInt)
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemFloat) DeleteLastCharacter() {
	if fi.valueInt == 0 && fi.valueFloat == 0 {
		if fi.isNegative {
			fi.isNegative = false
		} else {
			return
		}
	} else {
		if fi.valueFloat > 0 {
			mod := fi.valueFloat % 10
			fi.valueFloat = (fi.valueFloat - mod) / 10
		} else if fi.floatPosition > -1 {
			fi.floatPosition = -1
		} else {
			mod := fi.valueInt % 10
			fi.valueInt = (fi.valueInt - mod) / 10
		}
	}
	offsetX := fi.charOffsets[len(fi.charOffsets)-1]
	fi.cursorOffsetX = fi.cursorOffsetX - offsetX
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
	fi.charOffsets = fi.charOffsets[:len(fi.charOffsets)-1]
}
