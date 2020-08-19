package model

import (
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/sphere"

	"github.com/go-gl/mathgl/mgl32"
)

type BugBuilder struct {
	position        mgl32.Vec3
	scale           mgl32.Vec3
	wrapper         interfaces.GLWrapper
	bodyMaterial    *material.Material
	bottomMaterial  *material.Material
	eyeMaterial     *material.Material
	rotationX       float32
	rotationY       float32
	rotationZ       float32
	spherePrecision int
	lightAmbient    mgl32.Vec3
	lightDiffuse    mgl32.Vec3
	lightSpecular   mgl32.Vec3
	constantTerm    float32
	linearTerm      float32
	quadraticTerm   float32
	withLight       bool
}

func NewBugBuilder() *BugBuilder {
	return &BugBuilder{
		position:        mgl32.Vec3{0.0, 0.0, 0.0},
		scale:           mgl32.Vec3{1.0, 1.0, 1.0},
		wrapper:         nil,
		bodyMaterial:    material.Greenrubber,
		bottomMaterial:  material.Emerald,
		eyeMaterial:     material.Ruby,
		rotationX:       0,
		rotationY:       0,
		rotationZ:       0,
		spherePrecision: 20,
		lightAmbient:    mgl32.Vec3{1.0, 1.0, 1.0},
		lightDiffuse:    mgl32.Vec3{1.0, 1.0, 1.0},
		lightSpecular:   mgl32.Vec3{1.0, 1.0, 1.0},
		constantTerm:    1.0,
		linearTerm:      0.14,
		quadraticTerm:   0.07,
		withLight:       true,
	}
}

// SetPosition updates the position of the builder.
func (b *BugBuilder) SetPosition(p mgl32.Vec3) {
	b.position = p
}

// SetScale updates the scale of the builder.
func (b *BugBuilder) SetScale(s mgl32.Vec3) {
	b.scale = s
}

// SetWrapper sets the wrapper.
func (b *BugBuilder) SetWrapper(w interfaces.GLWrapper) {
	b.wrapper = w
}

// SetBodyMaterial updates the material of the body.
func (b *BugBuilder) SetBodyMaterial(m *material.Material) {
	b.bodyMaterial = m
}

// SetBottomMaterial updates the material of the bottom.
func (b *BugBuilder) SetBottomMaterial(m *material.Material) {
	b.bottomMaterial = m
}

// SetEyeMaterial updates the material of the eye.
func (b *BugBuilder) SetEyeMaterial(m *material.Material) {
	b.eyeMaterial = m
}

// SetRotation sets the rotationX, rotationY, rotationZ values. The inputs has to be angles.
func (b *BugBuilder) SetRotation(x, y, z float32) {
	b.rotationX = x
	b.rotationY = y
	b.rotationZ = z
}

// SetSpherePrecision sets the precision of the spheres.
func (b *BugBuilder) SetSpherePrecision(p int) {
	b.spherePrecision = p
}

// SetLightAmbient updates the ambient light component.
func (b *BugBuilder) SetLightAmbient(a mgl32.Vec3) {
	b.lightAmbient = a
}

// SetLightDiffuse updates the diffuse light component.
func (b *BugBuilder) SetLightDiffuse(d mgl32.Vec3) {
	b.lightDiffuse = d
}

// SetLightSpecular updates the specular light component.
func (b *BugBuilder) SetLightSpecular(s mgl32.Vec3) {
	b.lightSpecular = s
}

// SetLightTerms sets the constantTerm, linearTerm, quadraticTerm params of the spot lightsource.
func (b *BugBuilder) SetLightTerms(c, l, q float32) {
	b.constantTerm = c
	b.linearTerm = l
	b.quadraticTerm = q
}

// SetWithLight updates the withLight flag. If it is false, the lightsource of the bug will be nil.
func (b *BugBuilder) SetWithLight(l bool) {
	b.withLight = l
}

func (b *BugBuilder) BuildMaterial() *Bug {
	sphereBase := sphere.New(b.spherePrecision)
	V, I, bo := sphereBase.MaterialMeshInput()

	Body := mesh.NewMaterialMesh(V, I, b.bodyMaterial, b.wrapper)
	Body.SetScale(b.scale)
	Body.SetPosition(b.position)
	Body.RotateX(b.rotationX)
	Body.RotateY(b.rotationY)
	Body.RotateZ(b.rotationZ)

	Body.SetBoundingObject(bo)

	Bottom := mesh.NewMaterialMesh(V, I, b.bottomMaterial, b.wrapper)
	Bottom.SetScale(b.bottomScale())
	Bottom.SetPosition(b.bottomPosition())
	Bottom.SetParent(Body)
	Bottom.SetBoundingObject(bo)

	Eye1 := mesh.NewMaterialMesh(V, I, b.eyeMaterial, b.wrapper)
	Eye1.SetScale(b.eyeScale())
	Eye1.SetPosition(b.eye1Position())
	Eye1.SetParent(Body)
	Eye1.SetBoundingObject(bo)

	Eye2 := mesh.NewMaterialMesh(V, I, b.eyeMaterial, b.wrapper)
	Eye2.SetScale(b.eyeScale())
	Eye2.SetPosition(b.eye2Position())
	Eye2.SetParent(Body)
	Eye2.SetBoundingObject(bo)

	m := newCDModel()
	m.AddMesh(Bottom)
	m.AddMesh(Body)
	m.AddMesh(Eye1)
	m.AddMesh(Eye2)
	bug := &Bug{BaseCollisionDetectionModel: *m}
	if b.withLight {
		l := light.NewPointLight([4]mgl32.Vec3{
			b.bottomPosition(), // position
			b.lightAmbient,     // ambient
			b.lightDiffuse,     // diffuse
			b.lightSpecular,    // specular
		}, [3]float32{
			b.constantTerm,
			b.linearTerm,
			b.quadraticTerm,
		})
		bug.lightSource = l
	} else {
		bug.lightSource = nil
	}

	return bug
}
func (b *BugBuilder) rotationTransformationMatrix() mgl32.Mat4 {
	return mgl32.HomogRotate3DY(mgl32.DegToRad(b.rotationY)).Mul4(
		mgl32.HomogRotate3DX(mgl32.DegToRad(b.rotationX))).Mul4(
		mgl32.HomogRotate3DZ(mgl32.DegToRad(b.rotationZ)))
}
func (b *BugBuilder) bottomScale() mgl32.Vec3 {
	return mgl32.Vec3{0.5, 0.5, 0.5}
}
func (b *BugBuilder) eyeScale() mgl32.Vec3 {
	return mgl32.Vec3{0.1, 0.1, 0.1}
}
func (b *BugBuilder) bottomPosition() mgl32.Vec3 {
	origPos := mgl32.Vec3{b.scale.X() * -1, 0, 0}
	return mgl32.TransformCoordinate(origPos, b.rotationTransformationMatrix())
}
func (b *BugBuilder) eye1Position() mgl32.Vec3 {
	return b.eyePosition((mgl32.Vec3{1, 1, 1}).Normalize())
}
func (b *BugBuilder) eye2Position() mgl32.Vec3 {
	return b.eyePosition((mgl32.Vec3{1, 1, -1}).Normalize())
}
func (b *BugBuilder) eyePosition(basePos mgl32.Vec3) mgl32.Vec3 {
	normalBase := basePos.Normalize()
	baseScaled := mgl32.Vec3{normalBase.X() * b.scale.X(), normalBase.Y() * b.scale.Y(), normalBase.Z() * b.scale.Z()}
	return mgl32.TransformCoordinate(baseScaled, b.rotationTransformationMatrix())
}

type Bug struct {
	BaseCollisionDetectionModel
	lightSource *light.Light
}

// GetBottomPosition returns the current position of the bottom mesh.
// Transformations are applied, due to the relative position.
func (b *Bug) GetBottomPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{0, 0, 0}, b.meshes[0].ModelTransformation())
}

// GetBodyPosition returns the current position of the body mesh.
func (b *Bug) GetBodyPosition() mgl32.Vec3 {
	return b.meshes[1].GetPosition()
}

// GetEye1Position returns the current position of the eye1 mesh.
// Transformations are applied, due to the relative position.
func (b *Bug) GetEye1Position() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{0, 0, 0}, b.meshes[2].ModelTransformation())
}

// GetEye2Position returns the current position of the eye2 mesh.
// Transformations are applied, due to the relative position.
func (b *Bug) GetEye2Position() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{0, 0, 0}, b.meshes[3].ModelTransformation())
}

// GetLightSource returns the lightsource of the lamp.
func (b *Bug) GetLightSource() *light.Light {
	return b.lightSource
}

// Update function loops over each of the meshes and calls their Update function.
func (b *Bug) Update(dt float64) {
	for i, _ := range b.meshes {
		b.meshes[i].Update(dt)
	}
}
