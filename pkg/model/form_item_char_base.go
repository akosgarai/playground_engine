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

type FormItemCharBase struct {
	*FormItemBase
	cursor      interfaces.Mesh
	charOffsets []float32
	value       string
	maxLen      int
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
		FormItemBase: m,
		cursor:       cursor,
		charOffsets:  []float32{},
		value:        "",
		maxLen:       9,
	}
}
func (fi *FormItemCharBase) cursorOffsetX() float32 {
	result := float32(0.0)
	if len(fi.charOffsets) > 0 {
		for i := 0; i < len(fi.charOffsets); i++ {
			result = result + fi.charOffsets[i]
		}
	}
	return result
}

// GetTarget returns the input target Mesh
func (fi *FormItemCharBase) GetTarget() interfaces.Mesh {
	return fi.meshes[1]
}

// AddCursor displays a cursor on the target surface.
func (fi *FormItemCharBase) AddCursor() {
	fi.AddMesh(fi.cursor)
	fi.cursor.SetPosition(mgl32.Vec3{fi.GetCursorInitialPosition().X() - fi.cursorOffsetX(), 0.0, -0.01})
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
	fi.charOffsets = append(fi.charOffsets, offsetX)
	fi.cursor.SetPosition(mgl32.Vec3{fi.GetCursorInitialPosition().X() - fi.cursorOffsetX(), 0.0, -0.01})
}

// StepBackCursor moves the cursor back after a character deletion.
func (fi *FormItemCharBase) StepBackCursor() {
	fi.charOffsets = fi.charOffsets[:len(fi.charOffsets)-1]
	fi.cursor.SetPosition(mgl32.Vec3{fi.GetCursorInitialPosition().X() - fi.cursorOffsetX(), 0.0, -0.01})
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemCharBase) ValueToString() string {
	return fi.value
}
