package model

import (
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/cuboid"
	"github.com/akosgarai/playground_engine/pkg/primitives/cylinder"
	"github.com/akosgarai/playground_engine/pkg/primitives/sphere"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	defaultPoleHeight     = float32(1.0)
	widthHeightRatio      = float32(1.0) / float32(15.0)
	lengthHeightRatio     = float32(1.0) / float32(4.0)
	bulbRadiusHeightRatio = float32(1.0) / float32(60.0)
)

var (
	bulbMaterial = material.New(mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, mgl32.Vec3{1.0, 1.0, 1.0}, 256.0)
)

type StreetLamp struct {
	Model
}

// NewStreetLamp returns a street lamp like model. The StreetLamp is a mesh system.
// The 'position' input is the bottom center point of the 'pole' of the lamp. The top of the pole
// points to the +Z axis. The 'top' is the head of the lamp. Its position is relative to the pole.
// The 'bulb' is positioned relative to the 'top'.
func NewMaterialStreetLamp(position mgl32.Vec3, scale float32) *StreetLamp {
	height := defaultPoleHeight * scale
	width := height * widthHeightRatio
	length := height * lengthHeightRatio
	bulbRadius := height * bulbRadiusHeightRatio

	// pole
	poleCuboid := cuboid.New(width, height, width)
	poleV, poleI, bo := poleCuboid.MaterialMeshInput()
	pole := mesh.NewMaterialMesh(poleV, poleI, material.Chrome, glWrapper)
	pole.SetPosition(mgl32.Vec3{position.X(), position.Y() + height/2, position.Z()})
	pole.SetBoundingObject(bo)
	// top
	topCuboid := cuboid.New(length, width, width)
	topV, topI, bo := topCuboid.MaterialMeshInput()
	top := mesh.NewMaterialMesh(topV, topI, material.Chrome, glWrapper)
	top.SetPosition(mgl32.Vec3{(length - width) / 2, 0, (height + width) / 2})
	top.SetParent(pole)
	top.SetBoundingObject(bo)
	// bulb
	sph := sphere.New(15)
	sphereV, sphereI, bo := sph.TexturedMeshInput()
	bulb := mesh.NewMaterialMesh(sphereV, sphereI, bulbMaterial, glWrapper)
	bulb.SetPosition(mgl32.Vec3{length/2 - 4*bulbRadius, 0, -width / 2})
	bulb.SetScale(mgl32.Vec3{1.0, 1.0, 1.0}.Mul(bulbRadius))
	bulb.SetParent(top)
	bulb.SetBoundingObject(bo)

	m := newModel()
	m.AddMesh(bulb)
	m.AddMesh(top)
	m.AddMesh(pole)

	return &StreetLamp{Model: *m}
}

// NewTexturedStreetLamp returns a StreetLamp model that uses textured and textured material meshes.
// The 'position' input is the bottom center point of the 'pole' of the lamp. The top of the pole
// points to the +Z axis. The 'top' is the head of the lamp. Its position is relative to the pole.
// The 'bulb' is positioned relative to the 'top'.
func NewTexturedStreetLamp(position mgl32.Vec3, scale float32) *StreetLamp {
	// Setup the variables based on the given scale.
	height := defaultPoleHeight * scale
	width := height * widthHeightRatio
	length := height * lengthHeightRatio
	bulbRadius := height * bulbRadiusHeightRatio

	var metalTexture texture.Textures
	metalTexture.AddTexture("pkg/model/assets/metal.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	metalTexture.AddTexture("pkg/model/assets/metal.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	var bulbTexture texture.Textures
	bulbTexture.AddTexture("pkg/model/assets/crystal-ball.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	bulbTexture.AddTexture("pkg/model/assets/crystal-ball.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	// pole
	poleCylinder := cylinder.New(width/2, 20, height)
	poleV, poleI, bo := poleCylinder.TexturedMeshInput()
	pole := mesh.NewTexturedMesh(poleV, poleI, metalTexture, glWrapper)
	pole.SetPosition(mgl32.Vec3{position.X(), position.Y() + height/2, position.Z()})
	pole.SetBoundingObject(bo)

	// top
	topCylinder := cylinder.NewHalfCircleBased(width/2, 20, length)
	topV, topI, bo := topCylinder.TexturedMeshInput()
	top := mesh.NewTexturedMesh(topV, topI, metalTexture, glWrapper)
	top.SetPosition(mgl32.Vec3{(length - width) / 2, 0, (height) / 2})
	top.SetParent(pole)
	top.RotateZ(90)
	top.RotateY(90)
	top.SetBoundingObject(bo)

	// bulb
	sph := sphere.New(15)
	sphereV, sphereI, bo := sph.TexturedMeshInput()
	bulb := mesh.NewTexturedMaterialMesh(sphereV, sphereI, bulbTexture, bulbMaterial, glWrapper)
	bulb.SetPosition(mgl32.Vec3{length/2 - 4*bulbRadius, 0, 0})
	bulb.SetScale(mgl32.Vec3{1.0, 1.0, 1.0}.Mul(bulbRadius))
	bulb.SetParent(top)
	bulb.SetBoundingObject(bo)

	m := newModel()
	m.AddMesh(pole)
	m.AddMesh(top)
	m.AddMesh(bulb)

	return &StreetLamp{Model: *m}
}

// GetPolePosition returns the current position of the pole mesh.
func (s *StreetLamp) GetPolePosition() mgl32.Vec3 {
	return s.meshes[0].GetPosition()
}

// GetTopPosition returns the current position of the top mesh.
// Transformations are applied, due to the relative position.
func (s *StreetLamp) GetTopPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(s.meshes[1].GetPosition(), s.meshes[1].ModelTransformation())
}

// GetBulbPosition returns the current position of the bulb mesh.
// Transformations are applied, due to the relative position.
func (s *StreetLamp) GetBulbPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(s.meshes[2].GetPosition(), s.meshes[2].ModelTransformation())
}

// Update function loops over each of the meshes and calls their Update function.
func (s *StreetLamp) Update(dt float64) {
	for i, _ := range s.meshes {
		s.meshes[i].Update(dt)
	}
}
