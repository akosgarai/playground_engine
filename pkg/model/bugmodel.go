package model

import (
	"math"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/primitives/sphere"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	_WING_BOTTOM = 0
	_WING_UP     = 1
	_WING_TOP    = 2
	_WING_DOWN   = 3
)

type BugBuilder struct {
	position              mgl32.Vec3
	scale                 mgl32.Vec3
	wrapper               interfaces.GLWrapper
	bodyMaterial          *material.Material
	bottomMaterial        *material.Material
	eyeMaterial           *material.Material
	rotationX             float32
	rotationY             float32
	rotationZ             float32
	spherePrecision       int
	lightAmbient          mgl32.Vec3
	lightDiffuse          mgl32.Vec3
	lightSpecular         mgl32.Vec3
	constantTerm          float32
	linearTerm            float32
	quadraticTerm         float32
	withLight             bool
	velocity              float32
	direction             mgl32.Vec3
	movementRotationAngle float32
	movementRotationAxis  mgl32.Vec3
	sameDirectionTime     float32
	withWings             bool
}

func NewBugBuilder() *BugBuilder {
	return &BugBuilder{
		position:              mgl32.Vec3{0.0, 0.0, 0.0},
		scale:                 mgl32.Vec3{1.0, 1.0, 1.0},
		wrapper:               nil,
		bodyMaterial:          material.Greenrubber,
		bottomMaterial:        material.Emerald,
		eyeMaterial:           material.Ruby,
		rotationX:             0,
		rotationY:             0,
		rotationZ:             0,
		spherePrecision:       20,
		lightAmbient:          mgl32.Vec3{1.0, 1.0, 1.0},
		lightDiffuse:          mgl32.Vec3{1.0, 1.0, 1.0},
		lightSpecular:         mgl32.Vec3{1.0, 1.0, 1.0},
		constantTerm:          1.0,
		linearTerm:            0.14,
		quadraticTerm:         0.07,
		withLight:             true,
		velocity:              0.0,
		direction:             mgl32.Vec3{0, 0, 0},
		movementRotationAngle: 0.0,
		movementRotationAxis:  mgl32.Vec3{0, 0, 0},
		sameDirectionTime:     1000.0,
		withWings:             false,
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

// SetVelocity updates the velocity.
func (b *BugBuilder) SetVelocity(v float32) {
	b.velocity = v
}

// SetMovementRotationAngle updates the movementRotationAngle.
func (b *BugBuilder) SetMovementRotationAngle(v float32) {
	b.movementRotationAngle = v
}

// SetMovementRotationAxis updates the movementRotationAxis.
func (b *BugBuilder) SetMovementRotationAxis(v mgl32.Vec3) {
	b.movementRotationAxis = v
}

// SetDirection updates the direction.
func (b *BugBuilder) SetDirection(v mgl32.Vec3) {
	b.direction = v
}

// SetSameDirectionTime updates the sameDirectionTime.
func (b *BugBuilder) SetSameDirectionTime(v float32) {
	b.sameDirectionTime = v
}

// SetWithWings updates the withWings flag.
func (b *BugBuilder) SetWithWings(w bool) {
	b.withWings = w
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
	m.AddMesh(Body)
	m.AddMesh(Bottom)
	m.AddMesh(Eye1)
	m.AddMesh(Eye2)
	m.SetSpeed(b.velocity)
	m.SetDirection(b.direction)

	wingStrikeTime := float64(0.0)
	// attach point mesh
	attachPointWing1 := mesh.NewPointMesh(b.wrapper)
	attachPointWing1.SetParent(Body)
	attachPointWing2 := mesh.NewPointMesh(b.wrapper)
	attachPointWing2.SetParent(Body)
	m.AddMesh(attachPointWing1)
	m.AddMesh(attachPointWing2)
	if b.withWings {
		wingBase := rectangle.NewExact(1.0, 1.0)
		V, I, bo := wingBase.MeshInput()
		wing1 := mesh.NewMaterialMesh(V, I, b.bottomMaterial, b.wrapper)
		wing1.SetPosition(b.wing1Position())
		wing1.SetParent(attachPointWing1)
		wing1.SetBoundingObject(bo)
		m.AddMesh(wing1)
		wing2 := mesh.NewMaterialMesh(V, I, b.bottomMaterial, b.wrapper)
		wing2.SetPosition(b.wing2Position())
		wing2.SetParent(attachPointWing2)
		wing2.SetBoundingObject(bo)
		m.AddMesh(wing2)
		wingStrikeTime = 3000
		attachPointWing1.SetPosition(b.wing1AttachPointPosition())
		attachPointWing2.SetPosition(b.wing2AttachPointPosition())
	}

	bug := &Bug{
		BaseCollisionDetectionModel:  *m,
		movementRotationAngle:        b.movementRotationAngle,
		currentMovementRotationAngle: float32(0.0),
		movementRotationAxis:         b.movementRotationAxis,
		sinceLastRotate:              0.0,
		sameDirectionTime:            b.sameDirectionTime,
		wingStrikeTime:               wingStrikeTime,
		wingState:                    _WING_BOTTOM,
		currentWingAnimationTime:     0.0,
		maxWingRotationAngle:         75.0,
		currentWingRotationAngle:     0.0,
		wing1AttachPoint:             attachPointWing1,
		wing2AttachPoint:             attachPointWing2,
	}
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
func (b *BugBuilder) wing1AttachPointPosition() mgl32.Vec3 {
	return b.wingAttachPointPosition((mgl32.Vec3{0, -1, 1}).Normalize())
}
func (b *BugBuilder) wing2AttachPointPosition() mgl32.Vec3 {
	return b.wingAttachPointPosition((mgl32.Vec3{0, -1, -1}).Normalize())
}
func (b *BugBuilder) wingAttachPointPosition(basePos mgl32.Vec3) mgl32.Vec3 {
	normalBase := basePos.Normalize()
	baseScaled := mgl32.Vec3{normalBase.X() * b.scale.X(), normalBase.Y() * b.scale.Y(), normalBase.Z() * b.scale.Z()}
	return mgl32.TransformCoordinate(baseScaled, b.rotationTransformationMatrix())
}
func (b *BugBuilder) wing1Position() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{0.0, 0.0, 0.5}, b.rotationTransformationMatrix())
}
func (b *BugBuilder) wing2Position() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{0.0, 0.0, -0.5}, b.rotationTransformationMatrix())
}

type Bug struct {
	BaseCollisionDetectionModel
	lightSource                  *light.Light
	movementRotationAngle        float32
	currentMovementRotationAngle float32
	movementRotationAxis         mgl32.Vec3
	sinceLastRotate              float32
	sameDirectionTime            float32
	wingStrikeTime               float64
	wing1AttachPoint             interfaces.Mesh
	wing2AttachPoint             interfaces.Mesh
	// current wing animation time
	currentWingAnimationTime float64
	// wing state (going up or down or edge position)
	wingState int
	// maximum rotation angle of the wings
	maxWingRotationAngle float32
	// current rotation angle of the wings
	currentWingRotationAngle float32
}

// Bottom returns the bottom mesh.
func (b *Bug) Bottom() interfaces.Mesh {
	return b.meshes[1]
}

// Body returns the body mesh.
func (b *Bug) Body() interfaces.Mesh {
	return b.meshes[0]
}

// Eye1 returns the eye1 mesh.
func (b *Bug) Eye1() interfaces.Mesh {
	return b.meshes[2]
}

// Eye2 returns the eye2 mesh.
func (b *Bug) Eye2() interfaces.Mesh {
	return b.meshes[3]
}

// GetLightSource returns the lightsource of the lamp.
func (b *Bug) GetLightSource() *light.Light {
	return b.lightSource
}

// Update function loops over each of the meshes and calls their Update function.
// It also updates the direction of the bug, if necessary. The lighstource position
// also sync'd with the bottom mesh.
func (b *Bug) Update(dt float64) {
	b.sinceLastRotate = b.sinceLastRotate + float32(dt)
	if b.sinceLastRotate >= b.sameDirectionTime {
		b.sinceLastRotate = 0.0
		b.currentMovementRotationAngle = float32(math.Mod(float64(b.currentMovementRotationAngle+b.movementRotationAngle), 360))
		// current values
		cX, cY, cZ := matrixToAngles(b.Body().RotationTransformation())
		// expected values
		tX, tY, tZ := matrixToAngles(mgl32.HomogRotate3D(mgl32.DegToRad(b.currentMovementRotationAngle), mgl32.TransformNormal(b.movementRotationAxis, b.Body().RotationTransformation())))
		// rotate with the diff
		b.RotateY(tY - cY)
		b.RotateX(tX - cX)
		b.RotateZ(tZ - cZ)
	}
	if b.lightSource != nil {
		b.lightSource.SetPosition(mgl32.TransformCoordinate(mgl32.Vec3{0, 0, 0}, b.Bottom().ModelTransformation()))
	}
	if b.wingStrikeTime > 0.0 {
		b.animateWings(dt)
	}
	for i, _ := range b.meshes {
		b.meshes[i].Update(dt)
	}
}
func (b *Bug) pushState() {
	b.wingState = (b.wingState + 1) % 4
	b.currentWingAnimationTime = 0.0
}

// keep the wings in the edge states for a while, the rest of the time is for the movement.
// Movement starts bottom, then it goes up until the top position. after it goes down until the bottom position.
func (b *Bug) animateWings(dt float64) {
	// calculate the current delta time. If dt is gt than the remaining
	// animation time, it is decresed.
	maxDelta := math.Min(dt, b.wingStrikeTime-b.currentWingAnimationTime+dt)
	b.currentWingAnimationTime += maxDelta

	// calculate the rotation angle. It depends on the wingState.
	if b.wingState == _WING_BOTTOM || b.wingState == _WING_TOP {
		b.pushState()
		return
	}
	currentRotationAngle := float32(b.wingState-2) * b.maxWingRotationAngle / float32(b.wingStrikeTime) * float32(maxDelta)
	b.currentWingRotationAngle = b.currentWingRotationAngle - currentRotationAngle
	// sin, cos of the current angle.
	cosDeg := float32(math.Cos(float64(mgl32.DegToRad(b.currentWingRotationAngle))))
	sinDeg := float32(math.Sin(float64(mgl32.DegToRad(b.currentWingRotationAngle))))

	// rotation matrix of the base mesh.
	rotationMatrix := b.meshes[4].RotationTransformation()
	// current rotation angles of the w1:
	w1X, w1Y, w1Z := matrixToAngles(b.meshes[6].RotationTransformation())
	// current rotation angles of the w1:
	//w2X, w2Y, w2Z := matrixToAngles(b.meshes[5].RotationTransformation())
	// calculate the rotation vector of the door.
	rotatedOrigoBasedVector := mgl32.Vec3{0.0, -sinDeg, cosDeg}
	transformedVectorW1 := mgl32.TransformCoordinate(rotatedOrigoBasedVector, rotationMatrix)
	//transformedVectorW2 := mgl32.TransformNormal(rotatedOrigoBasedVector.Mul(-1), rotationMatrix)
	b.meshes[6].SetPosition(transformedVectorW1.Mul(0.5))
	//b.meshes[5].SetPosition(transformedVectorW2)

	// the rotation angles for the given full angle:
	transformedForward := mgl32.TransformNormal(mgl32.Vec3{1.0, 0.0, 0.0}, rotationMatrix)
	eX, eY, eZ := matrixToAngles(mgl32.HomogRotate3D(mgl32.DegToRad(b.currentWingRotationAngle), transformedForward).Mul4(rotationMatrix))

	b.meshes[6].RotateZ(eZ - w1Z)
	b.meshes[6].RotateX(eX - w1X)
	b.meshes[6].RotateY(eY - w1Y)
	/*

		b.meshes[5].RotateZ(-eZ - w2Z)
		b.meshes[5].RotateX(-eX - w2X)
		b.meshes[5].RotateY(-eY - w2Y)
	*/
	if b.currentWingAnimationTime >= b.wingStrikeTime {
		b.pushState()
	}
}
