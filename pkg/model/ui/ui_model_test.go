package ui

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNew(t *testing.T) {
	model := newModel()
	if _, err := model.GetMeshByIndex(0); err.Error() != "EMPTY_MESHES" {
		t.Error("Invalid empty meshes error.")
	}
}
func TestAttachToScreen(t *testing.T) {
	model := newModel()
	rect := rectangle.NewSquare()
	V, I, _ := rect.MeshInput()
	var tex texture.Textures
	msh := mesh.NewTexturedMaterialMesh(V, I, tex, material.Jade, wrapperMock)
	model.AddMesh(msh)
	pos := mgl32.Vec3{1, 1, 1}
	target := mesh.NewPointMesh(wrapperMock)
	model.AttachToScreen(target, pos)
	m, err := model.GetMeshByIndex(0)
	if err != nil {
		t.Errorf("Shouldn't have error here. '%#v'.", err)
	}
	if m.GetParent() != target {
		t.Error("Invalid target")
	}
}
