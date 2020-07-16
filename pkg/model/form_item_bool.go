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

type FormItemBool struct {
	*FormItemBase
	value bool
}

// GetValue returns the value of the form item.
func (fi *FormItemBool) GetValue() bool {
	return fi.value
}

// SetValue returns the value of the form item.
func (fi *FormItemBool) SetValue(v bool) {
	fi.value = v
}

func NewFormItemBool(maxWidth, itemWidth float32, label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemBool {
	m := NewFormItemBase(maxWidth, itemWidth, label, mat, wrapper)
	m.GetSurface().SetPosition(position)
	var ledTexture texture.Textures
	ledTexture.AddTexture(baseDirModel()+"/assets/led-button.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", wrapper)
	ledTexture.AddTexture(baseDirModel()+"/assets/led-button.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", wrapper)
	ledPrimitive := rectangle.NewExact(m.GetTargetWidth(), m.GetTargetHeight())
	v, i, _ := ledPrimitive.MeshInput()
	ledMesh := mesh.NewTexturedMaterialMesh(v, i, ledTexture, mat, wrapper)
	ledMesh.SetParent(m.GetSurface())
	ledMesh.SetPosition(m.GetTargetPosition())
	m.AddMesh(ledMesh)
	return &FormItemBool{
		FormItemBase: m,
		value:        false,
	}
}

// GetLight returns the ledMesh
func (fi *FormItemBool) GetLight() interfaces.Mesh {
	return fi.meshes[1]
}

// ValueToString returns the string representation of the value of the form item.
func (fi *FormItemBool) ValueToString() string {
	if fi.value {
		return "true"
	}
	return "false"
}
