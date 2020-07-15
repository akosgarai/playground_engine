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
	FormItemWidth  = float32(0.98)
	FormItemLength = float32(0.1)
	ledWidth       = float32(0.2)
	ledHeight      = float32(0.09)
)

type FormItemBool struct {
	*FormItemBase
	label string
	value bool
}

// GetLabel returns the label string of the item.
func (fi *FormItemBool) GetLabel() string {
	return fi.label
}

// GetValue returns the value of the form item.
func (fi *FormItemBool) GetValue() bool {
	return fi.value
}

// SetValue returns the value of the form item.
func (fi *FormItemBool) SetValue(v bool) {
	fi.value = v
}

func NewFormItemBool(label string, mat *material.Material, position mgl32.Vec3, wrapper interfaces.GLWrapper) *FormItemBool {
	m := NewFormItemBase(1.96, "Half", mat, wrapper)
	m.GetSurface().SetPosition(position)
	var ledTexture texture.Textures
	ledTexture.AddTexture(baseDirModel()+"/assets/led-button.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", wrapper)
	ledTexture.AddTexture(baseDirModel()+"/assets/led-button.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", wrapper)
	ledPrimitive := rectangle.NewExact(ledWidth, ledHeight)
	v, i, _ := ledPrimitive.MeshInput()
	ledMesh := mesh.NewTexturedMaterialMesh(v, i, ledTexture, mat, wrapper)
	ledMesh.SetParent(m.GetSurface())
	ledMesh.SetPosition(mgl32.Vec3{0.29, -0.01, 0.0})
	m.AddMesh(ledMesh)
	return &FormItemBool{
		FormItemBase: m,
		label:        label,
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
