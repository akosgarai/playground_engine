package model

import (
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
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

type StreetLampBuilder struct {
	position      mgl32.Vec3
	wrapper       interfaces.GLWrapper
	rotationX     float32
	rotationY     float32
	rotationZ     float32
	assetsBaseDir string
	poleLength    float32
	bulbMaterial  *material.Material
	constantTerm  float32
	linearTerm    float32
	quadraticTerm float32
	cutoff        float32
	outerCutoff   float32
	lampOn        bool
}

func NewStreetLampBuilder() *StreetLampBuilder {
	return &StreetLampBuilder{
		position:      mgl32.Vec3{0.0, 0.0, 0.0},
		wrapper:       nil,
		rotationX:     0,
		rotationY:     0,
		rotationZ:     0,
		assetsBaseDir: baseDirModel(),
		poleLength:    1,
		bulbMaterial:  bulbMaterial,
		constantTerm:  0.0,
		linearTerm:    0.0,
		quadraticTerm: 0.0,
		cutoff:        0.0,
		outerCutoff:   0.0,
		lampOn:        true,
	}
}

// SetPosition sets the position.
func (b *StreetLampBuilder) SetPosition(p mgl32.Vec3) {
	b.position = p
}

// SetWrapper sets the wrapper.
func (b *StreetLampBuilder) SetWrapper(w interfaces.GLWrapper) {
	b.wrapper = w
}

// SetRotation sets the rotationX, rotationY, rotationZ values. The inputs has to be angles.
func (b *StreetLampBuilder) SetRotation(x, y, z float32) {
	b.rotationX = x
	b.rotationY = y
	b.rotationZ = z
}

// SetAssetsBaseDir sets the base direction path string.
func (b *StreetLampBuilder) SetAssetsBaseDir(path string) {
	b.assetsBaseDir = path
}

// SetPoleLength sets the poleLength.
func (b *StreetLampBuilder) SetPoleLength(x float32) {
	b.poleLength = x
}

// SetBulbMaterial sets the bulbMaterial. This ambient, diffuse, specular color components
// are also used as the components of the lightsource.
func (b *StreetLampBuilder) SetBulbMaterial(mat *material.Material) {
	b.bulbMaterial = mat
}

// SetLightTerms sets the constantTerm, linearTerm, quadraticTerm params of the spot lightsource.
func (b *StreetLampBuilder) SetLightTerms(c, l, q float32) {
	b.constantTerm = c
	b.linearTerm = l
	b.quadraticTerm = q
}

// SetCutoff sets the cutoff, outerCutoff params of the spot lightsource.
func (b *StreetLampBuilder) SetCutoff(cutoff, outerCutoff float32) {
	b.cutoff = cutoff
	b.outerCutoff = outerCutoff
}

// SetLampOn sets the lampOn flag.
func (b *StreetLampBuilder) SetLampOn(v bool) {
	b.lampOn = v
}
func (b *StreetLampBuilder) rotationTransformationMatrix() mgl32.Mat4 {
	return mgl32.HomogRotate3DY(mgl32.DegToRad(b.rotationY)).Mul4(
		mgl32.HomogRotate3DX(mgl32.DegToRad(b.rotationX))).Mul4(
		mgl32.HomogRotate3DZ(mgl32.DegToRad(b.rotationZ)))
}
func (b *StreetLampBuilder) rotationTransformationMatrixTextureTop() mgl32.Mat4 {
	transformedUp := b.transformedUpDirection()
	transformedFront := b.transformedFrontDirection()
	rotationMatrixUp := mgl32.HomogRotate3D(mgl32.DegToRad(90), transformedUp)
	rotationMatrixFront := mgl32.HomogRotate3D(mgl32.DegToRad(90), transformedFront)
	return rotationMatrixUp.Mul4(rotationMatrixFront)
}

func (b *StreetLampBuilder) transformedUpDirection() mgl32.Vec3 {
	up := mgl32.Vec3{0.0, 1.0, 0.0}
	return mgl32.TransformNormal(up, b.rotationTransformationMatrix())
}
func (b *StreetLampBuilder) transformedFrontDirection() mgl32.Vec3 {
	up := mgl32.Vec3{0.0, 0.0, 1.0}
	return mgl32.TransformNormal(up, b.rotationTransformationMatrix())
}

// materialTopPosition returns the position of the top mesh for a material street lamp
func (b *StreetLampBuilder) materialTopPosition() mgl32.Vec3 {
	height, width, length, _ := b.getSizes()
	defaultPos := mgl32.Vec3{(length - width) / 2, 0, (height + width) / 2}
	return mgl32.TransformCoordinate(defaultPos, b.rotationTransformationMatrix())
}

// materialBulbPosition returns the position of the bulb mesh
func (b *StreetLampBuilder) materialBulbPosition() mgl32.Vec3 {
	_, width, length, bulbRadius := b.getSizes()
	defaultPos := mgl32.Vec3{length/2 - 4*bulbRadius, 0, -width / 2}
	return mgl32.TransformCoordinate(defaultPos, b.rotationTransformationMatrix())
}

// textureTopPosition returns the position of the top mesh for a texture street lamp
func (b *StreetLampBuilder) textureTopPosition() mgl32.Vec3 {
	height, width, length, _ := b.getSizes()
	defaultPos := mgl32.Vec3{(length - width) / 2, 0, (height) / 2}
	return mgl32.TransformCoordinate(defaultPos, b.rotationTransformationMatrix())
}

// textureBulbPosition returns the position of the bulb mesh
func (b *StreetLampBuilder) textureBulbPosition() mgl32.Vec3 {
	_, width, length, bulbRadius := b.getSizes()
	defaultPos := mgl32.Vec3{length/2 - 4*bulbRadius, 0, -width / 2}
	return mgl32.TransformCoordinate(defaultPos, b.rotationTransformationMatrix().Mul4(b.rotationTransformationMatrixTextureTop()))
}

// BuildMaterial returns a street lamp like model. The StreetLamp is a mesh system.
// The 'position' input is the bottom center point of the 'pole' of the lamp. The top of the pole
// points to the +Z axis. The 'top' is the head of the lamp. Its position is relative to the pole.
// The 'bulb' is positioned relative to the 'top'.
func (b *StreetLampBuilder) BuildMaterial() *StreetLamp {
	if b.wrapper == nil {
		panic("Wrapper is missing")
	}
	pole := b.materialPole()
	pole.RotateX(b.rotationX)
	pole.RotateY(b.rotationY)
	pole.RotateZ(b.rotationZ)

	top := b.materialTop()
	top.SetParent(pole)
	bulb := b.materialBulb()
	bulb.SetParent(top)

	m := newCDModel()
	m.AddMesh(pole)
	m.AddMesh(top)
	m.AddMesh(bulb)
	ls := light.NewSpotLight([5]mgl32.Vec3{
		mgl32.TransformCoordinate(bulb.GetPosition(), bulb.ModelTransformation()),
		mgl32.TransformNormal(mgl32.Vec3{0.0, 1.0, 0.0}, b.rotationTransformationMatrix()),
		b.bulbMaterial.GetAmbient(),
		b.bulbMaterial.GetDiffuse(),
		b.bulbMaterial.GetSpecular(),
	}, [5]float32{
		b.constantTerm,
		b.linearTerm,
		b.quadraticTerm,
		b.cutoff,
		b.outerCutoff,
	})

	sl := &StreetLamp{BaseCollisionDetectionModel: *m, lightSource: ls}
	if !b.lampOn {
		sl.TurnLampOff()
	}
	return sl
}

// NewTexturedStreetLamp returns a StreetLamp model that uses textured and textured material meshes.
// The 'position' input is the bottom center point of the 'pole' of the lamp. The top of the pole
// points to the +Z axis. The 'top' is the head of the lamp. Its position is relative to the pole.
// The 'bulb' is positioned relative to the 'top'.
func (b *StreetLampBuilder) BuildTexture() *StreetLamp {
	if b.wrapper == nil {
		panic("Wrapper is missing")
	}
	var metalTexture texture.Textures
	metalTexture.AddTexture(b.assetsBaseDir+"/assets/metal.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", b.wrapper)
	metalTexture.AddTexture(b.assetsBaseDir+"/assets/metal.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", b.wrapper)

	pole := b.texturePole(metalTexture)
	pole.RotateX(b.rotationX)
	pole.RotateY(b.rotationY)
	pole.RotateZ(b.rotationZ)

	top := b.textureTop(metalTexture)
	top.SetParent(pole)

	var bulbTexture texture.Textures
	bulbTexture.AddTexture(b.assetsBaseDir+"/assets/crystal-ball.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", b.wrapper)
	bulbTexture.AddTexture(b.assetsBaseDir+"/assets/crystal-ball.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", b.wrapper)
	bulb := b.textureBulb(bulbTexture)
	bulb.SetParent(top)

	m := newCDModel()
	m.AddMesh(pole)
	m.AddMesh(top)
	m.AddMesh(bulb)

	ls := light.NewSpotLight([5]mgl32.Vec3{
		mgl32.TransformCoordinate(bulb.GetPosition(), bulb.ModelTransformation()),
		mgl32.TransformNormal(mgl32.Vec3{0.0, 1.0, 0.0}, b.rotationTransformationMatrix()),
		b.bulbMaterial.GetAmbient(),
		b.bulbMaterial.GetDiffuse(),
		b.bulbMaterial.GetSpecular(),
	}, [5]float32{
		b.constantTerm,
		b.linearTerm,
		b.quadraticTerm,
		b.cutoff,
		b.outerCutoff,
	})

	sl := &StreetLamp{BaseCollisionDetectionModel: *m, lightSource: ls}
	if !b.lampOn {
		sl.TurnLampOff()
	}
	return sl
}

// It returns the size of the streetlamp components. They are calculated from the inputs and a couple of
// constant ratios.
func (b *StreetLampBuilder) getSizes() (float32, float32, float32, float32) {
	height := defaultPoleHeight * b.poleLength
	width := height * widthHeightRatio
	length := height * lengthHeightRatio
	bulbRadius := height * bulbRadiusHeightRatio
	return height, width, length, bulbRadius
}
func (b *StreetLampBuilder) materialPole() *mesh.MaterialMesh {
	height, width, _, _ := b.getSizes()
	poleCuboid := cuboid.New(width, height, width)
	V, I, bo := poleCuboid.MaterialMeshInput()
	pole := mesh.NewMaterialMesh(V, I, material.Chrome, b.wrapper)
	pole.SetPosition(b.position.Add(mgl32.Vec3{0.0, height / 2, 0.0}))
	pole.SetBoundingObject(bo)
	return pole
}
func (b *StreetLampBuilder) materialTop() *mesh.MaterialMesh {
	_, width, length, _ := b.getSizes()
	topCuboid := cuboid.New(length, width, width)
	V, I, bo := topCuboid.MaterialMeshInput()
	top := mesh.NewMaterialMesh(V, I, material.Chrome, b.wrapper)
	top.SetPosition(b.materialTopPosition())
	top.SetBoundingObject(bo)

	return top
}
func (b *StreetLampBuilder) materialBulb() *mesh.MaterialMesh {
	_, _, _, bulbRadius := b.getSizes()
	sph := sphere.New(15)
	V, I, bo := sph.TexturedMeshInput()
	bulb := mesh.NewMaterialMesh(V, I, b.bulbMaterial, b.wrapper)
	bulb.SetPosition(b.materialBulbPosition())
	bulb.SetScale(mgl32.Vec3{1.0, 1.0, 1.0}.Mul(bulbRadius))
	bulb.SetBoundingObject(bo)
	return bulb
}
func (b *StreetLampBuilder) texturePole(tex texture.Textures) *mesh.TexturedMesh {
	height, width, _, _ := b.getSizes()
	poleCylinder := cylinder.New(width/2, 20, height)
	V, I, bo := poleCylinder.TexturedMeshInput()
	pole := mesh.NewTexturedMesh(V, I, tex, b.wrapper)
	pole.SetPosition(b.position.Add(mgl32.Vec3{0.0, height / 2.0, 0.0}))
	pole.SetBoundingObject(bo)
	return pole
}
func (b *StreetLampBuilder) textureTop(tex texture.Textures) *mesh.TexturedMesh {
	_, width, length, _ := b.getSizes()
	topCylinder := cylinder.NewHalfCircleBased(width/2, 20, length)
	V, I, bo := topCylinder.TexturedMeshInput()
	top := mesh.NewTexturedMesh(V, I, tex, b.wrapper)
	top.SetPosition(b.textureTopPosition())
	rX, rY, rZ := matrixToAngles(b.rotationTransformationMatrixTextureTop())
	top.RotateZ(rZ)
	top.RotateX(rX)
	top.RotateY(rY)
	top.SetBoundingObject(bo)
	return top
}
func (b *StreetLampBuilder) textureBulb(tex texture.Textures) *mesh.TexturedMaterialMesh {
	_, _, _, bulbRadius := b.getSizes()
	sph := sphere.New(15)
	V, I, bo := sph.TexturedMeshInput()
	bulb := mesh.NewTexturedMaterialMesh(V, I, tex, b.bulbMaterial, b.wrapper)
	bulb.SetPosition(b.textureBulbPosition())
	bulb.SetScale(mgl32.Vec3{1.0, 1.0, 1.0}.Mul(bulbRadius))
	bulb.SetBoundingObject(bo)
	return bulb
}

type StreetLamp struct {
	BaseCollisionDetectionModel
	lightSource *light.Light
}

// TurnLampOn sets the lisghtsource color to the bulb material color
func (s *StreetLamp) TurnLampOn() {
	bulb := s.getBulb()
	switch bulb.(type) {
	case *mesh.MaterialMesh:
		b := bulb.(*mesh.MaterialMesh)
		s.lightSource.SetAmbient(b.Material.GetAmbient())
		s.lightSource.SetDiffuse(b.Material.GetDiffuse())
		s.lightSource.SetSpecular(b.Material.GetSpecular())
		break
	case *mesh.TexturedMaterialMesh:
		b := bulb.(*mesh.TexturedMaterialMesh)
		s.lightSource.SetAmbient(b.Material.GetAmbient())
		s.lightSource.SetDiffuse(b.Material.GetDiffuse())
		s.lightSource.SetSpecular(b.Material.GetSpecular())
		break
	}
}

// TurnLampOff sets the lisghtsource color to dark color
func (s *StreetLamp) TurnLampOff() {
	dark := mgl32.Vec3{0.0, 0.0, 0.0}
	s.lightSource.SetAmbient(dark)
	s.lightSource.SetDiffuse(dark)
	s.lightSource.SetSpecular(dark)
}

// GetLightSource returns the lightsource of the lamp.
func (s *StreetLamp) GetLightSource() *light.Light {
	return s.lightSource
}

// Update function loops over each of the meshes and calls their Update function.
func (s *StreetLamp) Update(dt float64) {
	for i, _ := range s.meshes {
		s.meshes[i].Update(dt)
	}
}
func (s *StreetLamp) getBulb() interfaces.Mesh {
	return s.meshes[2]
}
