package model

import (
	"errors"
	"fmt"
	"math"
	"path"
	"runtime"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/modelexport"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	emptyMeshesError = errors.New("EMPTY_MESHES")
)

func baseDirModel() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

type Model struct {
	meshes      []interfaces.Mesh
	transparent bool
	// uniforms, that needs to be set for this model.
	uniformFloat  map[string]float32    // map for float32
	uniformVector map[string]mgl32.Vec3 // map for 3 float32
}
type BaseModel struct {
	Model
}
type BaseCollisionDetectionModel struct {
	Model
}

// CollideTestWithSphere is the collision detection function for items in this mesh vs sphere.
func (m *BaseCollisionDetectionModel) CollideTestWithSphere(boundingSphere *coldet.Sphere) bool {
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

// Update function loops over each of the meshes and calls their Update function.
func (m *BaseModel) Update(dt float64) {
	for i, _ := range m.meshes {
		m.meshes[i].Update(dt)
	}
}
func newModel() *Model {
	return &Model{
		meshes:        []interfaces.Mesh{},
		transparent:   false,
		uniformFloat:  make(map[string]float32),
		uniformVector: make(map[string]mgl32.Vec3),
	}
}
func newCDModel() *BaseCollisionDetectionModel {
	return &BaseCollisionDetectionModel{
		Model{
			meshes:        []interfaces.Mesh{},
			transparent:   false,
			uniformFloat:  make(map[string]float32),
			uniformVector: make(map[string]mgl32.Vec3),
		},
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

// GetMeshByIndex function returns the mesh with the given index and nil.
// If the index is greater than the mesh size, or less than 0, it returns error.
// If the meshes is empty, it returns error.
func (m *Model) GetMeshByIndex(index int) (interfaces.Mesh, error) {
	if len(m.meshes) == 0 {
		return nil, emptyMeshesError
	}
	if index < 0 {
		return nil, errors.New("INVALID_INDEX")
	}
	if index+1 > len(m.meshes) {
		return nil, errors.New("INVALID_INDEX")
	}
	return m.meshes[index], nil
}

// Draw function loops over each of the meshes and calls their Draw function.
func (m *Model) Draw(s interfaces.Shader) {
	m.customUniforms(s)
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

// SetTransparent function updates the transparent flag.
func (m *Model) SetTransparent(t bool) {
	m.transparent = t
}

// IsTransparent function returns the transparent flag.
func (m *Model) IsTransparent() bool {
	return m.transparent
}

// SetUniformFloat sets the given float value to the given string key in
// the uniformFloat map.
func (m *Model) SetUniformFloat(key string, value float32) {
	m.uniformFloat[key] = value
}

// SetUniformVector sets the given mgl32.Vec3 value to the given string key in
// the uniformVector map.
func (m *Model) SetUniformVector(key string, value mgl32.Vec3) {
	m.uniformVector[key] = value
}

// Setup custom uniforms for the shader application.
func (m *Model) customUniforms(s interfaces.Shader) {
	for name, value := range m.uniformFloat {
		s.SetUniform1f(name, value)
	}
	for name, value := range m.uniformVector {
		s.SetUniform3f(name, value.X(), value.Y(), value.Z())
	}
}

// ClosestMeshTo returns the closest mesh to the given point.
func (m *Model) ClosestMeshTo(position mgl32.Vec3) (interfaces.Mesh, float32) {
	closest := float32(math.MaxFloat32)
	var closestMesh interfaces.Mesh
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
			var distance float32
			if meshBo.Type() == "AABB" {
				aabb := coldet.NewBoundingBox(pos, params["width"], params["height"], params["length"])
				distance = aabb.Distance([3]float32{position.X(), position.Y(), position.Z()})

			} else if meshBo.Type() == "Sphere" {
				bs := coldet.NewBoundingSphere(pos, params["radius"])
				distance = bs.Distance([3]float32{position.X(), position.Y(), position.Z()})

			}
			if distance < closest {
				closest = distance
				closestMesh = m.meshes[i]
			}
		}
	}
	return closestMesh, closest
}

// Clear function deletes the current meshes.
func (m *Model) Clear() {
	m.meshes = []interfaces.Mesh{}
}

// CollideTestWithSphere is the collision detection function for items in this mesh vs sphere.
func (m *BaseModel) CollideTestWithSphere(boundingSphere *coldet.Sphere) bool {
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
