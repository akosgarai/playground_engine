package model

import (
	"fmt"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/modelexport"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"
)

type Model struct {
	meshes []interfaces.Mesh
}
type BaseModel struct {
	Model
}

// Update function loops over each of the meshes and calls their Update function.
func (m *BaseModel) Update(dt float64) {
	for i, _ := range m.meshes {
		m.meshes[i].Update(dt)
	}
}
func newModel() *Model {
	return &Model{
		meshes: []interfaces.Mesh{},
	}
}

func New() *BaseModel {
	m := newModel()
	return &BaseModel{
		*m,
	}
}

// AddMesh function adds a mesh to the meshes.
func (m *Model) AddMesh(msh interfaces.Mesh) {
	m.meshes = append(m.meshes, msh)
}

// Draw function loops over each of the meshes and calls their Draw function.
func (m *Model) Draw(s interfaces.Shader) {
	for i, _ := range m.meshes {
		m.meshes[i].Draw(s)
	}
}

// Export function exports the meshes to a file
func (m *Model) Export(path string) {
	exporter := modelexport.New(m.meshes)
	err := exporter.Export(path)
	if err != nil {
		fmt.Printf("Export failed. '%s'\n", err.Error())
	}
}

// SetSpeed function loops over each of the parent meshes and calls their SetSpeed function.
func (m *Model) SetSpeed(s float32) {
	for i, _ := range m.meshes {
		if m.meshes[i].IsParentMesh() {
			m.meshes[i].SetSpeed(s)
		}
	}
}

// SetDirection function loops over each of the parent meshes and calls their SetDirection function.
func (m *Model) SetDirection(p mgl32.Vec3) {
	for i, _ := range m.meshes {
		if m.meshes[i].IsParentMesh() {
			m.meshes[i].SetDirection(p)
		}
	}
}

// RotateX function rotates the model with the given angle (has to be degree).
// It calls the RotateX function of each mesh.
func (m *Model) RotateX(angleDeg float32) {
	for i, _ := range m.meshes {
		if m.meshes[i].IsParentMesh() {
			m.meshes[i].RotateX(angleDeg)
		} else {
			m.meshes[i].RotatePosition(angleDeg, mgl32.Vec3{1.0, 0.0, 0.0})
		}
	}
}

// RotateY function rotates the model with the given angle (has to be degree).
// It calls the RotateY function of each mesh.
func (m *Model) RotateY(angleDeg float32) {
	for i, _ := range m.meshes {
		if m.meshes[i].IsParentMesh() {
			m.meshes[i].RotateY(angleDeg)
		} else {
			m.meshes[i].RotatePosition(angleDeg, mgl32.Vec3{0.0, 1.0, 0.0})
		}
	}
}

// RotateZ function rotates the model with the given angle (has to be degree).
// It calls the RotateZ function of each mesh.
func (m *Model) RotateZ(angleDeg float32) {
	for i, _ := range m.meshes {
		if m.meshes[i].IsParentMesh() {
			m.meshes[i].RotateZ(angleDeg)
		} else {
			m.meshes[i].RotatePosition(angleDeg, mgl32.Vec3{0.0, 0.0, 1.0})
		}
	}
}

// CollideTestWithSphere is the collision detection function for items in this mesh vs sphere.
func (m *Model) CollideTestWithSphere(boundingSphere *coldet.Sphere) bool {
	for i, _ := range m.meshes {
		if m.meshes[i].IsBoundingObjectSet() {
			meshBo := m.meshes[i].GetBoundingObject()
			meshInWorld := m.meshes[i].GetPosition()
			if !m.meshes[i].IsParentMesh() {
				meshTransTransform := m.meshes[i].GetParentTranslationTransformation()
				meshInWorld = mgl32.TransformCoordinate(meshInWorld, meshTransTransform)
			}
			pos := [3]float32{meshInWorld.X(), meshInWorld.Y(), meshInWorld.Z()}
			params := meshBo.Params()
			if meshBo.Type() == "AABB" {
				aabb := coldet.NewBoundingBox(pos, params["width"], params["height"], params["length"])

				if coldet.CheckSphereVsAabb(*boundingSphere, *aabb) {
					//fmt.Printf("BoundingSphere: %v\naabb: %v\nparams: %v\n", boundingSphere, aabb, params)
					return true
				}
			} else if meshBo.Type() == "Sphere" {
				params := meshBo.Params()
				bs := coldet.NewBoundingSphere(pos, params["radius"])

				if coldet.CheckSphereVsSphere(*boundingSphere, *bs) {
					//fmt.Printf("BoundingSphere: %v\nbs: %v\nparams: %v\n", boundingSphere, bs, params)
					return true
				}
			}
		}
	}
	return false
}
