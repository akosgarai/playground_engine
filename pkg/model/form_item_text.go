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

type FormItemText struct {
	*BaseModel
	cursor        interfaces.Mesh
	cursorOffsetX float32
	charOffsets   []float32
	label         string
	value         string
}

// GetLabel returns the label string of the item.
func (fi *FormItemText) GetLabel() string {
	return fi.label
}

// GetValue returns the value of the form item.
func (fi *FormItemText) GetValue() string {
	return fi.value
}

// SetValue returns the value of the form item.
func (fi *FormItemText) SetValue(v string) {
	fi.value = v
}
func NewFormItemText(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemText {
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
	return &FormItemText{
		BaseModel:     m,
		label:         label,
		cursor:        cursor,
		cursorOffsetX: 0.0,
		charOffsets:   []float32{},
		value:         "",
	}
}

// GetSurface returns the formItemMesh
func (fi *FormItemText) GetSurface() interfaces.Mesh {
	return fi.meshes[0]
}

// GetTarget returns the input target Mesh
func (fi *FormItemText) GetTarget() interfaces.Mesh {
	return fi.meshes[1]
}

// AddCursor displays a cursor on the target surface.
func (fi *FormItemText) AddCursor() {
	fi.AddMesh(fi.cursor)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}
func (fi *FormItemText) DeleteCursor() {
	if len(fi.meshes) == 3 {
		fi.meshes = fi.meshes[:len(fi.meshes)-1]
	}
}
func (fi *FormItemText) CharCallback(r rune, offsetX float32) {
	if !fi.validRune(r) {
		return
	}
	fi.value = fi.value + string(r)
	fi.cursorOffsetX = fi.cursorOffsetX + offsetX
	fi.charOffsets = append(fi.charOffsets, offsetX)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}
func (fi *FormItemText) validRune(r rune) bool {
	return true
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemText) ValueToString() string {
	return fi.value
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemText) DeleteLastCharacter() {
	if len(fi.charOffsets) == 0 {
		return
	}
	fi.value = fi.value[:len(fi.value)-1]
	offsetX := fi.charOffsets[len(fi.charOffsets)-1]
	fi.cursorOffsetX = fi.cursorOffsetX - offsetX
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
	fi.charOffsets = fi.charOffsets[:len(fi.charOffsets)-1]
}
