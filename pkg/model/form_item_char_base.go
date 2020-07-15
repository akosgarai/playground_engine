package model

import (
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	writableWidth  = float32(0.40)
	writableHeight = float32(0.09)

	cursorHeight = float32(0.07)
	cursorWidth  = float32(0.015)
	CursorInitX  = float32(0.155)
)

var (
	formItemEnums = []string{"Full", "Half", "Long", "Short"}
)

func validFormItemEnum(name string) bool {
	for _, i := range formItemEnums {
		if i == name {
			return true
		}
	}
	return false
}

type FormItemBase struct {
	*BaseModel
	width    float32
	sizeEnum string
}

// NewFormItemBase returns a FormItemBase. Its input is the width of the screen,
// the size enum of the item. Possible values: 'Full', 'Half', 'Long', 'Short'
// the material of the surface and a gl wrapper.
// In case of invalid input enum, it panics.
// It creates the surface mesh.
func NewFormItemBase(w float32, sizeEnum string, mat *material.Material, wrapper interfaces.GLWrapper) *FormItemBase {
	if !validFormItemEnum(sizeEnum) {
		panic("Invalid size enum.")
	}
	m := New()
	fi := &FormItemBase{
		BaseModel: m,
		width:     w,
		sizeEnum:  sizeEnum,
	}
	labelPrimitive := rectangle.NewExact(fi.GetFormItemWidth(), fi.GetFormItemHeight())
	v, i, bo := labelPrimitive.MeshInput()
	var tex texture.Textures
	tex.TransparentTexture(1, 1, 128, "tex.diffuse", wrapper)
	tex.TransparentTexture(1, 1, 128, "tex.specular", wrapper)
	formItemMesh := mesh.NewTexturedMaterialMesh(v, i, tex, mat, wrapper)
	formItemMesh.SetBoundingObject(bo)
	m.AddMesh(formItemMesh)
	return fi
}
func (fi *FormItemBase) widthMultiplier() float32 {
	var result float32
	switch fi.sizeEnum {
	case "Full":
		result = 1.0
		break
	case "Half":
		result = 0.5
		break
	case "Long":
		result = 2.0 / 3.0
		break
	case "Short":
		result = 1.0 / 3.0
		break
	}
	return result
}
func (fi *FormItemBase) heightMultiplier() float32 {
	return 1.0 / 1.96
}

// It returns the width of the form item.
func (fi *FormItemBase) GetFormItemWidth() float32 {
	return fi.width * fi.widthMultiplier()
}

// It returns the height of the form item.
func (fi *FormItemBase) GetFormItemHeight() float32 {
	return fi.width * fi.heightMultiplier()
}

// GetSurface returns the formItemMesh
func (fi *FormItemBase) GetSurface() interfaces.Mesh {
	return fi.meshes[0]
}

type FormItemCharBase struct {
	*FormItemBase
	cursor        interfaces.Mesh
	cursorOffsetX float32
	charOffsets   []float32
	label         string
	value         string
	maxLen        int
}

// NewFormItemCharBase returns a FormItemCharBase that could be the base of text based form items.
func NewFormItemCharBase(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemCharBase {
	m := NewFormItemBase(1.96, "Half", mat, wrapper)
	m.GetSurface().SetPosition(position)
	var writableTexture texture.Textures
	writableTexture.AddTexture(baseDirModel()+"/assets/paper.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", wrapper)
	writableTexture.AddTexture(baseDirModel()+"/assets/paper.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", wrapper)
	writablePrimitive := rectangle.NewExact(writableWidth, writableHeight)
	v, i, bo := writablePrimitive.MeshInput()
	writableMesh := mesh.NewTexturedMaterialMesh(v, i, writableTexture, mat, wrapper)
	writableMesh.SetParent(m.GetSurface())
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
	return &FormItemCharBase{
		FormItemBase:  m,
		label:         label,
		cursor:        cursor,
		cursorOffsetX: 0.0,
		charOffsets:   []float32{},
		value:         "",
		maxLen:        9,
	}
}

// GetLabel returns the label string of the item.
func (fi *FormItemCharBase) GetLabel() string {
	return fi.label
}

// GetTarget returns the input target Mesh
func (fi *FormItemCharBase) GetTarget() interfaces.Mesh {
	return fi.meshes[1]
}

// AddCursor displays a cursor on the target surface.
func (fi *FormItemCharBase) AddCursor() {
	fi.AddMesh(fi.cursor)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}

// DeleteCursor removes the cursor from the meshes.
func (fi *FormItemCharBase) DeleteCursor() {
	if len(fi.meshes) == 3 {
		fi.meshes = fi.meshes[:len(fi.meshes)-1]
	}
}

// MoveCursorWithOffset moves to cursor with the offset.
// It adds the new offset to the offsets, increments the sum offset
// and sets the cursor position.
func (fi *FormItemCharBase) MoveCursorWithOffset(offsetX float32) {
	fi.cursorOffsetX = fi.cursorOffsetX + offsetX
	fi.charOffsets = append(fi.charOffsets, offsetX)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}

// StepBackCursor moves the cursor back after a character deletion.
func (fi *FormItemCharBase) StepBackCursor() {
	offsetX := fi.charOffsets[len(fi.charOffsets)-1]
	fi.cursorOffsetX = fi.cursorOffsetX - offsetX
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
	fi.charOffsets = fi.charOffsets[:len(fi.charOffsets)-1]
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemCharBase) ValueToString() string {
	return fi.value
}
