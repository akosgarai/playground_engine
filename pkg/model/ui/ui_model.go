package ui

import (
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"

	"github.com/go-gl/mathgl/mgl32"
)

type UIModel struct {
	*model.BaseModel
}

func newModel() *UIModel {
	return &UIModel{
		model.New(),
	}
}

// AttachToScreen gets the screen mesh and the position vector (relative to screen) as input.
// It sets the parent mesh of the frame to the screen, and sets the frame position to the
// new one.
func (m *UIModel) AttachToScreen(screen interfaces.Mesh, position mgl32.Vec3) {
	f, _ := m.GetMeshByIndex(0)
	frame := f.(*mesh.TexturedMaterialMesh)
	frame.SetParent(screen)
	frame.SetPosition(position)
}
