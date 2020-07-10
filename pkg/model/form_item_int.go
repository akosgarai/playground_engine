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
)

type FormItemInt struct {
	*BaseModel
	label string
	value int
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
	writableMesh.SetPosition(mgl32.Vec3{0.29, -0.01, 0.0})
	writableMesh.SetBoundingObject(bo)
	m.AddMesh(writableMesh)
	return &FormItemInt{
		BaseModel: m,
		label:     label,
		value:     0,
	}
}

// GetSurface returns the formItemMesh
func (fi *FormItemInt) GetSurface() interfaces.Mesh {
	return fi.meshes[0]
}

// GetTarget returns the ledMesh
func (fi *FormItemInt) GetTarget() interfaces.Mesh {
	return fi.meshes[1]
}
