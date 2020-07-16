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
	ITEM_WIDTH_FULL        = float32(1.0)
	ITEM_WIDTH_HALF        = float32(0.5)
	ITEM_WIDTH_LONG        = float32(2.0 / 3.0)
	ITEM_WIDTH_SHORT       = float32(1.0 / 3.0)
	ITEM_HEIGHT_MULTIPLIER = float32(0.1 / 1.96)
)

type FormItemBase struct {
	*BaseModel
	width float32
	size  float32
	label string
}

// NewFormItemBase returns a FormItemBase. Its input is the width of the screen,
// the size (scale) of the item, (some constants are provided)
// the label, the material of the surface and a gl wrapper.
// In case of invalid input enum, it panics.
// It creates the surface mesh.
func NewFormItemBase(w, size float32, label string, mat *material.Material, wrapper interfaces.GLWrapper) *FormItemBase {
	m := New()
	fi := &FormItemBase{
		BaseModel: m,
		width:     w,
		size:      size,
		label:     label,
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
	return fi.size
}

// It returns the width of the form item.
func (fi *FormItemBase) GetFormItemWidth() float32 {
	return fi.width * fi.widthMultiplier()
}

// It returns the height of the form item.
func (fi *FormItemBase) GetFormItemHeight() float32 {
	return fi.width * ITEM_HEIGHT_MULTIPLIER
}

// It returns the width of the label area. (55% of the halfwidth)
func (fi *FormItemBase) GetLabelAreaWidth() float32 {
	return fi.width * 0.275
}

// GetSurface returns the formItemMesh
func (fi *FormItemBase) GetSurface() interfaces.Mesh {
	return fi.meshes[0]
}

// GetLabel returns the label string of the item.
func (fi *FormItemBase) GetLabel() string {
	return fi.label
}

// GetTargetHeight returns the height size of the target mesh
// (text area, button)
func (fi *FormItemBase) GetTargetHeight() float32 {
	return fi.GetFormItemHeight() * 0.9
}

// GetTargetWidth returns the width size of the target mesh
// (text area, button)
func (fi *FormItemBase) GetTargetWidth() float32 {
	return fi.getWidthWithoutLabel() * 0.9
}
func (fi *FormItemBase) getWidthWithoutLabel() float32 {
	return fi.GetFormItemWidth() - fi.GetLabelAreaWidth()
}

// GetTargetPosition returns the position vector of the target mesh.
func (fi *FormItemBase) GetTargetPosition() mgl32.Vec3 {
	pX := -fi.GetFormItemWidth()/2 + fi.GetLabelAreaWidth() + (fi.getWidthWithoutLabel() / 2)
	return mgl32.Vec3{pX, -0.01, 0.0}
}

// GetCursorHeight returns the height size of the cursor.
func (fi *FormItemBase) GetCursorHeight() float32 {
	return fi.GetFormItemHeight() * 0.7
}

// GetCursorWidth returns the width size of the cursor.
func (fi *FormItemBase) GetCursorWidth() float32 {
	return fi.GetFormItemHeight() * 0.15
}

// GetCursorInitialPosition returns the initial position vector of the cursor.
func (fi *FormItemBase) GetCursorInitialPosition() mgl32.Vec3 {
	return mgl32.Vec3{(fi.GetTargetWidth()*0.9 - fi.GetCursorWidth()) / 2, -0.01, 0.0}
}

type FormItemCharBase struct {
	*FormItemBase
	cursor        interfaces.Mesh
	cursorOffsetX float32
	charOffsets   []float32
	value         string
	maxLen        int
}

// NewFormItemCharBase returns a FormItemCharBase that could be the base of text based form items.
func NewFormItemCharBase(maxWidth, widthRatio float32, label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemCharBase {
	m := NewFormItemBase(maxWidth, widthRatio, label, mat, wrapper)
	m.GetSurface().SetPosition(position)
	var writableTexture texture.Textures
	writableTexture.AddTexture(baseDirModel()+"/assets/paper.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", wrapper)
	writableTexture.AddTexture(baseDirModel()+"/assets/paper.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", wrapper)
	writablePrimitive := rectangle.NewExact(m.GetTargetWidth(), m.GetTargetHeight())
	v, i, bo := writablePrimitive.MeshInput()
	writableMesh := mesh.NewTexturedMaterialMesh(v, i, writableTexture, mat, wrapper)
	writableMesh.SetParent(m.GetSurface())
	writableMesh.SetPosition(m.GetTargetPosition())
	writableMesh.SetBoundingObject(bo)
	m.AddMesh(writableMesh)
	cursorPrimitive := rectangle.NewExact(m.GetCursorWidth(), m.GetCursorHeight())
	var ctex texture.Textures
	ctex.TransparentTexture(1, 1, 255, "tex.diffuse", wrapper)
	ctex.TransparentTexture(1, 1, 255, "tex.specular", wrapper)
	v, i, _ = cursorPrimitive.MeshInput()
	cursor := mesh.NewTexturedMaterialMesh(v, i, ctex, material.Greenplastic, wrapper)
	cursor.SetPosition(m.GetCursorInitialPosition())
	cursor.SetParent(writableMesh)
	return &FormItemCharBase{
		FormItemBase:  m,
		cursor:        cursor,
		cursorOffsetX: 0.0,
		charOffsets:   []float32{},
		value:         "",
		maxLen:        9,
	}
}

// GetTarget returns the input target Mesh
func (fi *FormItemCharBase) GetTarget() interfaces.Mesh {
	return fi.meshes[1]
}

// AddCursor displays a cursor on the target surface.
func (fi *FormItemCharBase) AddCursor() {
	fi.AddMesh(fi.cursor)
	fi.cursor.SetPosition(mgl32.Vec3{fi.GetCursorInitialPosition().X() - fi.cursorOffsetX, 0.0, -0.01})
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
	fi.cursor.SetPosition(mgl32.Vec3{fi.GetCursorInitialPosition().X() - fi.cursorOffsetX, 0.0, -0.01})
}

// StepBackCursor moves the cursor back after a character deletion.
func (fi *FormItemCharBase) StepBackCursor() {
	offsetX := fi.charOffsets[len(fi.charOffsets)-1]
	fi.cursorOffsetX = fi.cursorOffsetX - offsetX
	fi.cursor.SetPosition(mgl32.Vec3{fi.GetCursorInitialPosition().X() - fi.cursorOffsetX, 0.0, -0.01})
	fi.charOffsets = fi.charOffsets[:len(fi.charOffsets)-1]
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemCharBase) ValueToString() string {
	return fi.value
}
