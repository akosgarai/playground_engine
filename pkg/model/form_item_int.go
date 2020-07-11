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

type FormItemInt struct {
	*BaseModel
	cursor        interfaces.Mesh
	cursorOffsetX float32
	label         string
	value         int
	isNegative    bool
}

// GetLabel returns the label string of the item.
func (fi *FormItemInt) GetLabel() string {
	return fi.label
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
	return &FormItemInt{
		BaseModel:     m,
		label:         label,
		cursor:        cursor,
		cursorOffsetX: 0.0,
		value:         0,
		isNegative:    false,
	}
}

// GetSurface returns the formItemMesh
func (fi *FormItemInt) GetSurface() interfaces.Mesh {
	return fi.meshes[0]
}

// GetTarget returns the input target Mesh
func (fi *FormItemInt) GetTarget() interfaces.Mesh {
	return fi.meshes[1]
}

// AddCursor displays a cursor on the target surface.
func (fi *FormItemInt) AddCursor() {
	fi.AddMesh(fi.cursor)
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}
func (fi *FormItemInt) DeleteCursor() {
	if len(fi.meshes) == 3 {
		fi.meshes = fi.meshes[:len(fi.meshes)-1]
	}
}
func (fi *FormItemInt) CharCallback(r rune, offsetX float32) {
	// if the first character is '-', mark the form item as negative.
	if fi.value == 0 && r == rune('-') {
		fi.isNegative = true
		fi.cursorOffsetX = fi.cursorOffsetX + offsetX
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
	fi.cursor.SetPosition(mgl32.Vec3{CursorInitX - fi.cursorOffsetX, 0.0, -0.01})
}
func (fi *FormItemInt) validRune(r rune) bool {
	// integer number isn't allowed to start with 0.
	if fi.value == 0 && r == rune('0') {
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
