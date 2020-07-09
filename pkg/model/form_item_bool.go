package model

import (
	"path"
	"runtime"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	width  = float32(0.98)
	length = float32(0.1)
)

func baseDirModel() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

type FormItemBool struct {
	*BaseModel
	label string
}

// GetLabel returns the label string of the item.
func (fi *FormItemBool) GetLabel() string {
	return fi.label
}

func NewFormItemBool(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemBool {
	labelPrimitive := rectangle.NewExact(width, length)
	v, i, bo := labelPrimitive.MeshInput()
	formItemMesh := mesh.NewMaterialMesh(v, i, mat, wrapper)
	formItemMesh.SetBoundingObject(bo)
	formItemMesh.SetPosition(position)
	m := New()
	m.AddMesh(formItemMesh)
	var ledTexture texture.Textures
	ledTexture.AddTexture(baseDirModel()+"/assets/led-button.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", wrapper)
	ledTexture.AddTexture(baseDirModel()+"/assets/led-button.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", wrapper)
	ledPrimitive := rectangle.NewExact(0.2, 0.1)
	v, i, bo = ledPrimitive.MeshInput()
	ledMesh := mesh.NewTexturedMesh(v, i, ledTexture, wrapper)
	ledMesh.SetBoundingObject(bo)
	ledMesh.SetParent(formItemMesh)
	ledMesh.SetPosition(mgl32.Vec3{0.29, 0, 0})
	m.AddMesh(ledMesh)
	return &FormItemBool{
		BaseModel: m,
		label:     label,
	}
}

// GetSurface returns the formItemMesh
func (fi *FormItemBool) GetSurface() interfaces.Mesh {
	return fi.meshes[0]
}
