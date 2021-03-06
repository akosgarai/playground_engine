package model

import (
	"math"

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
	CURSOR_WIDTH           = float32(0.015)
	CHAR_NUM_INT           = 10
	CHAR_NUM_FLOAT         = 10
	CHAR_NUM_INT64         = 20
	CHAR_NUM_TEXT          = 35
)

type FloatValidator func(f float32) bool
type Int64Validator func(i int64) bool
type IntValidator func(i int) bool
type StringValidator func(s string) bool

type FormItemBase struct {
	*BaseModel
	width       float32
	size        float32
	label       string
	description string
	// This value is used for scaling the items.
	aspect float32
}

// NewFormItemBase returns a FormItemBase. Its input is the width of the screen,
// the size (scale) of the item, (some constants are provided)
// the label, the material of the surface and a gl wrapper.
// In case of invalid input enum, it panics.
// It creates the surface mesh.
func NewFormItemBase(w, size, aspect float32, label, description string, mat *material.Material, wrapper interfaces.GLWrapper) *FormItemBase {
	m := New()
	fi := &FormItemBase{
		BaseModel:   m,
		width:       w,
		size:        size,
		label:       label,
		description: description,
		aspect:      aspect,
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
	return fi.width * ITEM_HEIGHT_MULTIPLIER * fi.aspect
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

// GetDescription returns the description string of the item.
func (fi *FormItemBase) GetDescription() string {
	return fi.description
}

// GetTargetHeight returns the height size of the target mesh
// (text area, button)
func (fi *FormItemBase) GetTargetHeight() float32 {
	return fi.GetFormItemHeight() * 0.9
}

// GetTargetWidth returns the width size of the target mesh
// (text area, button)
func (fi *FormItemBase) GetTargetWidth() float32 {
	return fi.getWidthWithoutLabel() * 0.95
}
func (fi *FormItemBase) getWidthWithoutLabel() float32 {
	return fi.GetFormItemWidth() - fi.GetLabelAreaWidth()
}

// GetTargetPosition returns the position vector of the target mesh.
func (fi *FormItemBase) GetTargetPosition() mgl32.Vec3 {
	pX := -fi.GetFormItemWidth()/2 + fi.GetLabelAreaWidth() + (fi.GetTargetWidth() / 2)
	return mgl32.Vec3{pX, -0.01, 0.0}
}

// GetVectorTargetPosition returns the position vector of the target mesh.
func (fi *FormItemBase) GetVectorTargetPosition(index int) mgl32.Vec3 {
	containerWidth := fi.GetVectorTargetWidth()
	firstMiddlePosition := -fi.GetFormItemWidth()/2 + fi.GetLabelAreaWidth() + containerWidth/2
	pX := firstMiddlePosition + float32(index)*containerWidth
	return mgl32.Vec3{pX, -0.01, 0.0}
}

// GetVectorTargetWidth returns the width of a vector target item.
func (fi *FormItemBase) GetVectorTargetWidth() float32 {
	return fi.getWidthWithoutLabel() / 3.0 * 0.95
}

// GetCursorHeight returns the height size of the cursor.
func (fi *FormItemBase) GetCursorHeight() float32 {
	return fi.GetFormItemHeight() * 0.7
}

// GetCursorWidth returns the width size of the cursor.
func (fi *FormItemBase) GetCursorWidth() float32 {
	return CURSOR_WIDTH * float32(math.Min(1, float64(fi.aspect)))
}

// GetCursorInitialPosition returns the initial position vector of the cursor.
func (fi *FormItemBase) GetCursorInitialPosition() mgl32.Vec3 {
	return mgl32.Vec3{(fi.GetTargetWidth()*0.85 - fi.GetCursorWidth()) / 2, -0.01, 0.0}
}

// GetVectorCursorInitialPosition returns the initial position vector of the cursor.
func (fi *FormItemBase) GetVectorCursorInitialPosition() mgl32.Vec3 {
	return mgl32.Vec3{(fi.GetVectorTargetWidth()*0.85 - fi.GetCursorWidth()) / 2, -0.01, 0.0}
}
