package model

import (
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/sphere"

	"github.com/go-gl/mathgl/mgl32"
)

type Bug struct {
	BaseCollisionDetectionModel
}

// NewBug returns a bug instance. A Bug is a sphered mesh system. Its 'Body'
// is the parent mesh. Its position is absolute in the world coordinate system.
// The other 3 mesh, the 'Bottom', the 'Eye1' and 'Eye2' are child meshes, their
// position is relative to the parent.
func NewBug(position, scale mgl32.Vec3, wrapper interfaces.GLWrapper) *Bug {
	sphereBase := sphere.New(20)
	i, v, bo := sphereBase.MaterialMeshInput()
	// Body supposed to be other green. Like green rubber
	Body := mesh.NewMaterialMesh(i, v, material.Greenrubber, wrapper)
	Body.SetScale(scale)
	Body.SetPosition(position)
	Body.SetBoundingObject(bo)

	// Bottom supposed to be greenish color / material like emerald
	Bottom := mesh.NewMaterialMesh(i, v, material.Emerald, wrapper)
	Bottom.SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
	Bottom.SetPosition(mgl32.Vec3{scale.X() * -1, 0, 0})
	Bottom.SetParent(Body)
	Bottom.SetBoundingObject(bo)
	// Eyes are red. (red plastic)
	Eye1 := mesh.NewMaterialMesh(i, v, material.Ruby, wrapper)
	Eye1.SetScale(mgl32.Vec3{0.1, 0.1, 0.1})
	initPosBase := (mgl32.Vec3{1, 1, 1}).Normalize()
	initPosScaled := mgl32.Vec3{initPosBase.X() * scale.X(), initPosBase.Y() * scale.Y(), initPosBase.Z() * scale.Z()}
	Eye1.SetPosition(initPosScaled)
	Eye1.SetParent(Body)
	Eye1.SetBoundingObject(bo)

	Eye2 := mesh.NewMaterialMesh(i, v, material.Ruby, wrapper)
	Eye2.SetScale(mgl32.Vec3{0.1, 0.1, 0.1})
	initPosBase = (mgl32.Vec3{1, 1, -1}).Normalize()
	initPosScaled = mgl32.Vec3{initPosBase.X() * scale.X(), initPosBase.Y() * scale.Y(), initPosBase.Z() * scale.Z()}
	Eye2.SetPosition(initPosScaled)
	Eye2.SetParent(Body)
	Eye2.SetBoundingObject(bo)

	m := newCDModel()
	m.AddMesh(Bottom)
	m.AddMesh(Body)
	m.AddMesh(Eye1)
	m.AddMesh(Eye2)

	return &Bug{BaseCollisionDetectionModel: *m}
}

// GetBottomPosition returns the current position of the bottom mesh.
// Transformations are applied, due to the relative position.
func (b *Bug) GetBottomPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(b.meshes[0].GetPosition(), b.meshes[0].ModelTransformation())
}

// GetBodyPosition returns the current position of the body mesh.
func (b *Bug) GetBodyPosition() mgl32.Vec3 {
	return b.meshes[1].GetPosition()
}

// GetEye1Position returns the current position of the eye1 mesh.
// Transformations are applied, due to the relative position.
func (b *Bug) GetEye1Position() mgl32.Vec3 {
	return mgl32.TransformCoordinate(b.meshes[2].GetPosition(), b.meshes[2].ModelTransformation())
}

// GetEye2Position returns the current position of the eye2 mesh.
// Transformations are applied, due to the relative position.
func (b *Bug) GetEye2Position() mgl32.Vec3 {
	return mgl32.TransformCoordinate(b.meshes[3].GetPosition(), b.meshes[3].ModelTransformation())
}

// Update function loops over each of the meshes and calls their Update function.
func (b *Bug) Update(dt float64) {
	for i, _ := range b.meshes {
		b.meshes[i].Update(dt)
	}
}
