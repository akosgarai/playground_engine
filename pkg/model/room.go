package model

import (
	"fmt"
	"math"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/cuboid"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	_DOOR_OPENED      = 0
	_DOOR_CLOSING     = 1
	_DOOR_CLOSED      = 2
	_DOOR_OPENING     = 3
	doorAnimationTime = float64(1000)
)

type RoomBuilder struct {
	position      mgl32.Vec3 // the position of the room (center point of the floor mesh)
	width         float32    // the length of the usable area in the x axis
	height        float32    // the length of the usable area in the y axis
	length        float32    // the length of the usable area in the z axis
	wallWidth     float32    // the width of the walls
	doorWidth     float32    // the width of the door that is on the right side of the front wall.
	doorHeight    float32    // the height of the door that is on the right side of the front wall.
	rotationX     float32
	rotationY     float32
	rotationZ     float32
	assetsBaseDir string // In case of textured room, we have to know where are the assets.
	frontWindow   bool
	backWindow    bool
	leftWindow    bool
	rightWindow   bool
	windowWidth   float32 // the width of the windows that we could set on the textured rooms.
	windowHeight  float32 // the height of the windows that we could set on the textured rooms.
	wrapper       interfaces.GLWrapper
}

func NewRoomBuilder() *RoomBuilder {
	return &RoomBuilder{
		position:      mgl32.Vec3{0.0, 0.0, 0.0},
		width:         1.0,
		height:        1.0,
		length:        1.0,
		wallWidth:     0.005,
		doorWidth:     0.4,
		doorHeight:    0.6,
		rotationX:     0.0,
		rotationY:     0.0,
		rotationZ:     0.0,
		wrapper:       nil,
		frontWindow:   false,
		backWindow:    false,
		leftWindow:    false,
		rightWindow:   false,
		windowWidth:   0.2,
		windowHeight:  0.4,
		assetsBaseDir: baseDirModel(),
	}
}

// WithFrontWindow sets the frontWindow flag
func (b *RoomBuilder) WithFrontWindow(v bool) {
	b.frontWindow = v
}

// WithBackWindow sets the backWindow flag
func (b *RoomBuilder) WithBackWindow(v bool) {
	b.backWindow = v
}

// WithLeftWindow sets the leftWindow flag
func (b *RoomBuilder) WithLeftWindow(v bool) {
	b.leftWindow = v
}

// WithRightWindow sets the rightWindow flag
func (b *RoomBuilder) WithRightWindow(v bool) {
	b.rightWindow = v
}

// SetPosition sets the position.
func (b *RoomBuilder) SetPosition(p mgl32.Vec3) {
	b.position = p
}

// SetWrapper sets the wrapper.
func (b *RoomBuilder) SetWrapper(w interfaces.GLWrapper) {
	b.wrapper = w
}

// SetSize sets the width, height, length values.
func (b *RoomBuilder) SetSize(w, h, l float32) {
	b.width = w
	b.height = h
	b.length = l
}

// SetWallWidth sets the wallWidth.
func (b *RoomBuilder) SetWallWidth(w float32) {
	b.wallWidth = w
}

// SetDoorSize sets the doorWidth, doorHeight values.
func (b *RoomBuilder) SetDoorSize(w, h float32) {
	b.doorWidth = w
	b.doorHeight = h
}

// SetRotation sets the rotationX, rotationY, rotationZ values. The inputs has to be angles.
func (b *RoomBuilder) SetRotation(x, y, z float32) {
	b.rotationX = x
	b.rotationY = y
	b.rotationZ = z
}

// SetWindowSize sets the windowWidth, windowHeight values.
func (b *RoomBuilder) SetWindowSize(w, h float32) {
	b.windowWidth = w
	b.windowHeight = h
}

// SetAssetsBaseDir sets the base direction path string.
func (b *RoomBuilder) SetAssetsBaseDir(path string) {
	b.assetsBaseDir = path
}
func (b *RoomBuilder) rotationTransformationMatrix() mgl32.Mat4 {
	return mgl32.HomogRotate3DY(mgl32.DegToRad(b.rotationY)).Mul4(
		mgl32.HomogRotate3DX(mgl32.DegToRad(b.rotationX))).Mul4(
		mgl32.HomogRotate3DZ(mgl32.DegToRad(b.rotationZ)))
}

// the ceiling is above the floor.
func (b *RoomBuilder) ceilingPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{0.0, b.height, 0.0}, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) fullBackWallPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{0.0, (b.height - b.wallWidth) / 2, -(b.length + b.wallWidth) / 2}, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) fullLeftWallPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{(b.width - b.wallWidth) / 2, b.height / 2, 0.0}, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) fullRightWallPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{-(b.width - b.wallWidth) / 2, b.height / 2, 0.0}, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) frontDoorPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{-b.doorWidth / 2, 0.0, 0.0}, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) frontDoorWallAttachPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{b.width / 2, b.doorHeight / 2, (b.length - b.wallWidth) / 2}, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) frontAboveDoorWallPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{((b.width - b.doorWidth) / 2), (b.height + b.doorHeight) / 2, (b.length - b.wallWidth) / 2}, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) fullFrontWallPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{-b.doorWidth / 2, b.height / 2, (b.length - b.wallWidth) / 2}, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) stripFrontLongWidth() float32 {
	return (b.width - b.doorWidth - b.windowWidth) / 2
}
func (b *RoomBuilder) stripFrontShortHeight() float32 {
	return (b.height - b.windowHeight) / 2
}
func (b *RoomBuilder) stripFrontLeftWallPosition() mgl32.Vec3 {
	return mgl32.TransformCoordinate(mgl32.Vec3{-(b.width - b.stripFrontLongWidth()) / 2, b.height / 2, (b.length - b.wallWidth) / 2}, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) stripFrontRightWallPosition() mgl32.Vec3 {
	origPosition := mgl32.Vec3{-b.width/2 + b.stripFrontLongWidth() + b.windowWidth + b.stripFrontLongWidth()/2, b.height / 2, (b.length - b.wallWidth) / 2}
	return mgl32.TransformCoordinate(origPosition, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) stripFrontTopWallPosition() mgl32.Vec3 {
	origPosition := mgl32.Vec3{
		-(b.width - 2*b.stripFrontLongWidth() - b.windowWidth) / 2,
		b.stripFrontShortHeight() + b.windowHeight + b.stripFrontShortHeight()/2,
		(b.length - b.wallWidth) / 2,
	}
	return mgl32.TransformCoordinate(origPosition, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) stripFrontBottomWallPosition() mgl32.Vec3 {
	origPosition := mgl32.Vec3{
		-(b.width - 2*b.stripFrontLongWidth() - b.windowWidth) / 2,
		b.stripFrontShortHeight() / 2,
		(b.length - b.wallWidth) / 2,
	}
	return mgl32.TransformCoordinate(origPosition, b.rotationTransformationMatrix())
}
func (b *RoomBuilder) frontWindowPosition() mgl32.Vec3 {
	origPosition := mgl32.Vec3{
		-(b.width - 2*b.stripFrontLongWidth() - b.windowWidth) / 2,
		b.stripFrontShortHeight() + b.windowHeight/2,
		(b.length - b.wallWidth) / 2,
	}
	return mgl32.TransformCoordinate(origPosition, b.rotationTransformationMatrix())
}

// BuildTexture returns a textured material room that is constructed from the given setup.
func (b *RoomBuilder) BuildTexture() *Room {
	if b.wrapper == nil {
		panic("Wrapper is missing.")
	}
	var concreteTexture texture.Textures
	concreteTexture.AddTexture(b.assetsBaseDir+"/assets/concrete-wall.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", b.wrapper)
	concreteTexture.AddTexture(b.assetsBaseDir+"/assets/concrete-wall.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", b.wrapper)
	var doorTexture texture.Textures
	doorTexture.AddTexture(b.assetsBaseDir+"/assets/door.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", b.wrapper)
	doorTexture.AddTexture(b.assetsBaseDir+"/assets/door.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", b.wrapper)

	var windowTexture texture.Textures
	windowTexture.AddTexture(b.assetsBaseDir+"/assets/window.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", b.wrapper)
	windowTexture.AddTexture(b.assetsBaseDir+"/assets/window.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", b.wrapper)

	m := newCDModel()

	// floor + ceiling
	basementSizeCuboid := cuboid.New(b.width, b.length, b.wallWidth)
	basementV, basementI, bo := basementSizeCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)

	floor := mesh.NewTexturedMaterialMesh(basementV, basementI, concreteTexture, material.Chrome, b.wrapper)
	floor.SetPosition(b.position)
	floor.SetBoundingObject(bo)
	floor.RotateY(b.rotationY)
	floor.RotateX(b.rotationX)
	floor.RotateZ(b.rotationZ)
	m.AddMesh(floor)

	ceiling := mesh.NewTexturedMaterialMesh(basementV, basementI, concreteTexture, material.Chrome, b.wrapper)
	ceiling.SetPosition(b.ceilingPosition())
	ceiling.SetParent(floor)
	ceiling.SetBoundingObject(bo)
	m.AddMesh(ceiling)

	// attach point mesh
	attachPoint := mesh.NewPointMesh(b.wrapper)
	attachPoint.SetParent(floor)
	attachPoint.SetPosition(b.frontDoorWallAttachPosition())

	// door
	doorCuboid := cuboid.New(b.doorWidth, b.wallWidth, b.doorHeight)
	V, I, bo := doorCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_SAME)

	door := mesh.NewTexturedMesh(V, I, doorTexture, b.wrapper)
	door.SetPosition(b.frontDoorPosition())
	door.SetParent(attachPoint)
	door.SetBoundingObject(bo)
	m.AddMesh(door)

	// front above the door.
	frontTopCuboid := cuboid.New(b.doorWidth, b.wallWidth, b.height-b.doorHeight)
	V, I, bo = frontTopCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)

	frontWallRest := mesh.NewTexturedMesh(V, I, concreteTexture, b.wrapper)
	frontWallRest.SetPosition(b.frontAboveDoorWallPosition())
	frontWallRest.SetParent(floor)
	frontWallRest.SetBoundingObject(bo)
	m.AddMesh(frontWallRest)

	// back wall
	backWallSizeCuboid := cuboid.New(b.width, b.wallWidth, b.height)
	backWallV, backWallI, bo := backWallSizeCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)

	backWall := mesh.NewTexturedMaterialMesh(backWallV, backWallI, concreteTexture, material.Chrome, b.wrapper)
	backWall.SetPosition(b.fullBackWallPosition())
	backWall.SetParent(floor)
	backWall.SetBoundingObject(bo)
	m.AddMesh(backWall)

	// side wall
	sideWallSizeCuboid := cuboid.New(b.wallWidth, b.length, b.height)
	sideWallV, sideWallI, bo := sideWallSizeCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)

	rightWall := mesh.NewTexturedMaterialMesh(sideWallV, sideWallI, concreteTexture, material.Chrome, b.wrapper)
	rightWall.SetPosition(b.fullRightWallPosition())
	rightWall.SetParent(floor)
	rightWall.SetBoundingObject(bo)
	m.AddMesh(rightWall)

	leftWall := mesh.NewTexturedMaterialMesh(sideWallV, sideWallI, concreteTexture, material.Chrome, b.wrapper)
	leftWall.SetPosition(b.fullLeftWallPosition())
	leftWall.SetParent(floor)
	leftWall.SetBoundingObject(bo)
	m.AddMesh(leftWall)

	// front wall parts
	if b.frontWindow {
		frontSideCuboid := cuboid.New(b.stripFrontLongWidth(), b.wallWidth, b.height)
		V, I, bo = frontSideCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)

		frontWallMain1 := mesh.NewTexturedMaterialMesh(V, I, concreteTexture, material.Chrome, b.wrapper)
		frontWallMain1.SetPosition(b.stripFrontLeftWallPosition())
		frontWallMain1.SetParent(floor)
		frontWallMain1.SetBoundingObject(bo)
		m.AddMesh(frontWallMain1)

		frontWallMain2 := mesh.NewTexturedMaterialMesh(V, I, concreteTexture, material.Chrome, b.wrapper)
		frontWallMain2.SetPosition(b.stripFrontRightWallPosition())
		frontWallMain2.SetParent(floor)
		frontWallMain2.SetBoundingObject(bo)
		m.AddMesh(frontWallMain2)

		frontSmallCuboid := cuboid.New(b.windowWidth, b.wallWidth, b.stripFrontShortHeight())
		V, I, bo = frontSmallCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)

		frontWallMain3 := mesh.NewTexturedMaterialMesh(V, I, concreteTexture, material.Chrome, b.wrapper)
		frontWallMain3.SetPosition(b.stripFrontTopWallPosition())
		frontWallMain3.SetParent(floor)
		frontWallMain3.SetBoundingObject(bo)
		m.AddMesh(frontWallMain3)

		frontWallMain4 := mesh.NewTexturedMaterialMesh(V, I, concreteTexture, material.Chrome, b.wrapper)
		frontWallMain4.SetPosition(b.stripFrontBottomWallPosition())
		frontWallMain4.SetParent(floor)
		frontWallMain4.SetBoundingObject(bo)
		m.AddMesh(frontWallMain4)

		windowCuboid := cuboid.New(b.windowWidth, b.wallWidth, b.windowHeight)
		V, I, bo = windowCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)

		window := mesh.NewTexturedMesh(V, I, windowTexture, b.wrapper)
		window.SetPosition(b.frontWindowPosition())
		window.SetParent(floor)
		window.SetBoundingObject(bo)
		m.AddMesh(window)
	} else {
		frontCuboid := cuboid.New(b.width-b.doorWidth, b.wallWidth, b.height)
		V, I, bo := frontCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)

		frontWallMain := mesh.NewTexturedMaterialMesh(V, I, concreteTexture, material.Chrome, b.wrapper)
		frontWallMain.SetPosition(b.fullFrontWallPosition())
		frontWallMain.SetParent(floor)
		frontWallMain.SetBoundingObject(bo)
		m.AddMesh(frontWallMain)
	}

	if b.frontWindow || b.backWindow || b.leftWindow || b.rightWindow {
		m.SetTransparent(true)
	}
	return &Room{
		BaseCollisionDetectionModel: *m,
		doorState:                   _DOOR_CLOSED,
		currentAnimationTime:        0,
		doorAnimationonAngle:        90.0,
		doorWidth:                   b.doorWidth,
		doorWallAttachPoint:         attachPoint,
	}
}

// BuildMaterial returns a material room that is constructed from the given setup.
func (b *RoomBuilder) BuildMaterial() *Room {
	if b.wrapper == nil {
		panic("Wrapper is missing.")
	}
	m := newCDModel()

	// floor + ceiling
	basementSizeCuboid := cuboid.New(b.width, b.length, b.wallWidth)
	basementV, basementI, bo := basementSizeCuboid.MaterialMeshInput()

	floor := mesh.NewMaterialMesh(basementV, basementI, material.Chrome, b.wrapper)
	floor.SetPosition(b.position)
	floor.SetBoundingObject(bo)
	floor.RotateY(b.rotationY)
	floor.RotateX(b.rotationX)
	floor.RotateZ(b.rotationZ)
	m.AddMesh(floor)

	ceiling := mesh.NewMaterialMesh(basementV, basementI, material.Chrome, b.wrapper)
	ceiling.SetPosition(b.ceilingPosition())
	ceiling.SetParent(floor)
	ceiling.SetBoundingObject(bo)
	m.AddMesh(ceiling)

	// attach point mesh
	attachPoint := mesh.NewPointMesh(b.wrapper)
	attachPoint.SetParent(floor)
	attachPoint.SetPosition(b.frontDoorWallAttachPosition())
	// front door
	doorCuboid := cuboid.New(b.doorWidth, b.wallWidth, b.doorHeight)
	V, I, bo := doorCuboid.MaterialMeshInput()

	door := mesh.NewMaterialMesh(V, I, material.Bronze, b.wrapper)
	door.SetPosition(b.frontDoorPosition())
	door.SetParent(attachPoint)
	door.SetBoundingObject(bo)
	m.AddMesh(door)

	// front above the door.
	frontTopCuboid := cuboid.New(b.doorWidth, b.wallWidth, b.height-b.doorHeight)
	V, I, bo = frontTopCuboid.MaterialMeshInput()

	frontWallRest := mesh.NewMaterialMesh(V, I, material.Chrome, b.wrapper)
	frontWallRest.SetPosition(b.frontAboveDoorWallPosition())
	frontWallRest.SetParent(floor)
	frontWallRest.SetBoundingObject(bo)
	m.AddMesh(frontWallRest)

	// back wall
	backWallSizeCuboid := cuboid.New(b.width, b.wallWidth, b.height)
	backWallV, backWallI, bo := backWallSizeCuboid.MaterialMeshInput()

	backWall := mesh.NewMaterialMesh(backWallV, backWallI, material.Chrome, b.wrapper)
	backWall.SetPosition(b.fullBackWallPosition())
	backWall.SetParent(floor)
	backWall.SetBoundingObject(bo)
	m.AddMesh(backWall)

	// side wall
	sideWallSizeCuboid := cuboid.New(b.wallWidth, b.length, b.height)
	sideWallV, sideWallI, bo := sideWallSizeCuboid.MaterialMeshInput()

	rightWall := mesh.NewMaterialMesh(sideWallV, sideWallI, material.Chrome, b.wrapper)
	rightWall.SetPosition(b.fullRightWallPosition())
	rightWall.SetParent(floor)
	rightWall.SetBoundingObject(bo)
	m.AddMesh(rightWall)

	leftWall := mesh.NewMaterialMesh(sideWallV, sideWallI, material.Chrome, b.wrapper)
	leftWall.SetPosition(b.fullLeftWallPosition())
	leftWall.SetParent(floor)
	leftWall.SetBoundingObject(bo)
	m.AddMesh(leftWall)

	// front wall parts
	frontCuboid := cuboid.New(b.width-b.doorWidth, b.wallWidth, b.height)
	V, I, bo = frontCuboid.MaterialMeshInput()

	frontWallMain := mesh.NewMaterialMesh(V, I, material.Chrome, b.wrapper)
	frontWallMain.SetPosition(b.fullFrontWallPosition())
	frontWallMain.SetParent(floor)
	frontWallMain.SetBoundingObject(bo)
	m.AddMesh(frontWallMain)

	return &Room{
		BaseCollisionDetectionModel: *m,
		doorState:                   _DOOR_CLOSED,
		currentAnimationTime:        0,
		doorAnimationonAngle:        90.0,
		doorWidth:                   b.doorWidth,
		doorWallAttachPoint:         attachPoint,
	}
}

type Room struct {
	BaseCollisionDetectionModel
	doorState            int
	currentAnimationTime float64
	doorWallAttachPoint  interfaces.Mesh
	doorAnimationonAngle float32
	doorWidth            float32
}

func (r *Room) PushDoorState() {
	if r.doorState == _DOOR_OPENED || r.doorState == _DOOR_CLOSED {
		r.doorState += 1
		r.currentAnimationTime = 0
	}
}
func (r *Room) animateDoor(dt float64) {
	// early return if possible
	if r.doorState == _DOOR_OPENED || r.doorState == _DOOR_CLOSED {
		return
	}
	// calculate the current delta time. If dt is gt than the remaining
	// animation time, it is decresed.
	maxDelta := math.Min(dt, doorAnimationTime-r.currentAnimationTime+dt)
	r.currentAnimationTime += maxDelta

	// calculate the rotation angle. It depends on the doorState.
	var rotationDegY float32
	if r.doorState == _DOOR_OPENING {
		rotationDegY = float32(90.0 / doorAnimationTime * maxDelta)
	}
	if r.doorState == _DOOR_CLOSING {
		rotationDegY = float32(-90.0 / doorAnimationTime * maxDelta)
	}
	// The current animation angle is increased with the current rotation deg.
	r.doorAnimationonAngle = r.doorAnimationonAngle - rotationDegY

	// sin, cos of the current angle.
	cosDeg := float32(math.Cos(float64(mgl32.DegToRad(r.doorAnimationonAngle))))
	sinDeg := float32(math.Sin(float64(mgl32.DegToRad(r.doorAnimationonAngle))))

	// calculate the rotation vector of the door.
	rotatedOrigoBasedVector := mgl32.Vec3{-sinDeg, 0.0, cosDeg}
	// what if i transform the robv with the rotation of the parent mesh. in this case it will be the fine position of the stuff.
	attachPointRotationMatrix := r.doorWallAttachPoint.RotationTransformation()
	transformedVector := mgl32.TransformNormal(rotatedOrigoBasedVector, attachPointRotationMatrix)
	// rotation weight form the components of the rotated up vector
	transformedUp := mgl32.TransformNormal(mgl32.Vec3{0.0, 1.0, 0.0}, attachPointRotationMatrix)
	transformedUpInvert := mgl32.TransformNormal(mgl32.Vec3{0.0, 1.0, 0.0}, attachPointRotationMatrix.Inv())

	// the new position of the door.
	doorPosFromAttachPoint := transformedVector.Mul(r.doorWidth / 2)
	fmt.Printf("------------\nDoorNewPosition:\t%v\nDoorAttachPoint:\t%v\nRotatedUnitVector:\t%v\nTransformedVector:\t%v\nTransformedUp:\t\t%v\nTransformedUpInv:\t%v\n",
		doorPosFromAttachPoint, r.doorWallAttachPoint.GetPosition(), rotatedOrigoBasedVector, transformedVector, transformedUp, transformedUpInvert)

	// get rotation euler angles. transformedUp is the axis of our transformation. The angle is rotationDegY.
	// From the HomogRotate3D matrix, the euler angle could be computed. https://www.geometrictools.com/Documentation/EulerAngles.pdf (2.3)
	rX, rY, rZ := r.matrixToAngles(mgl32.HomogRotate3D(mgl32.DegToRad(rotationDegY), transformedUp))
	rX2, rY2, rZ2 := r.matrixToAngles(mgl32.HomogRotate3D(mgl32.DegToRad(rotationDegY), transformedUpInvert))
	door := r.GetDoor()
	// translation transformation inv. rotation, translation back:
	translationTransformationMatrix := r.doorWallAttachPoint.TranslationTransformation()
	fullMatrix := translationTransformationMatrix.Inv().Mul4(attachPointRotationMatrix.Inv()).Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(rotationDegY), transformedUp)).Mul4(attachPointRotationMatrix).Mul4(translationTransformationMatrix)
	fullMatrixI := translationTransformationMatrix.Inv().Mul4(attachPointRotationMatrix.Inv()).Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(rotationDegY), transformedUpInvert)).Mul4(attachPointRotationMatrix).Mul4(translationTransformationMatrix)
	rXOA, _, rZOA := r.matrixToAngles(attachPointRotationMatrix)
	apXZ := mgl32.HomogRotate3DX(mgl32.DegToRad(rXOA)).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(rZOA)))
	fullMatrixOrigAxis := translationTransformationMatrix.Inv().Mul4(apXZ.Inv()).Mul4(mgl32.HomogRotate3DY(mgl32.DegToRad(rotationDegY))).Mul4(apXZ).Mul4(translationTransformationMatrix)
	fullMatrixOrigAxisIDontKnow := translationTransformationMatrix.Inv().Mul4(apXZ.Inv()).Mul4(mgl32.HomogRotate3D(mgl32.DegToRad(rotationDegY), transformedUp)).Mul4(apXZ).Mul4(translationTransformationMatrix)
	rX3, rY3, rZ3 := r.matrixToAngles(fullMatrix)
	rX4, rY4, rZ4 := r.matrixToAngles(fullMatrixI)
	rX5, rY5, rZ5 := r.matrixToAngles(fullMatrixOrigAxis)
	rX6, rY6, rZ6 := r.matrixToAngles(fullMatrixOrigAxisIDontKnow)
	fmt.Printf("----------------\nrotationDegY: %f\n", rotationDegY)
	fmt.Printf("---------------\nRx: %f Ry: %f Rz: %f\n", rX, rY, rZ)
	fmt.Printf("---------------\nRx2: %f Ry2: %f Rz2: %f\n", rX2, rY2, rZ2)
	fmt.Printf("---------------\nRx3: %f Ry3: %f Rz3: %f\n", rX3, rY3, rZ3)
	fmt.Printf("---------------\nRx4: %f Ry4: %f Rz4: %f\n", rX4, rY4, rZ4)
	fmt.Printf("---------------\nRx5: %f Ry5: %f Rz5: %f\n", rX5, rY5, rZ5)
	fmt.Printf("---------------\nRx6: %f Ry6: %f Rz6: %f\n", rX6, rY6, rZ6)

	// Update door position to the newly calculated one.
	door.SetPosition(doorPosFromAttachPoint)
	// Apply the rotation on the y axis.
	door.RotateZ(rZ4)
	door.RotateX(rX4)
	door.RotateY(rY4)

	if r.currentAnimationTime >= doorAnimationTime {
		r.doorState = (r.doorState + 1) % 4
	}
}

// returns angles
func (r *Room) matrixToAngles(m mgl32.Mat4) (float32, float32, float32) {
	var x, y, z float32
	if m.At(1, 2) < 1 {
		if m.At(1, 2) > -1 {
			x = float32(math.Asin(-float64(m.At(1, 2))))
			y = float32(math.Atan2(float64(m.At(0, 2)), float64(m.At(2, 2))))
			z = float32(math.Atan2(float64(m.At(1, 0)), float64(m.At(1, 1))))
		} else {
			x = math.Pi / 2
			y = -float32(math.Atan2(-float64(m.At(0, 1)), float64(m.At(0, 0))))
			z = 0
		}
	} else {
		x = -math.Pi / 2
		y = float32(math.Atan2(-float64(m.At(0, 1)), float64(m.At(0, 0))))
		z = 0
	}
	return mgl32.RadToDeg(x), mgl32.RadToDeg(y), mgl32.RadToDeg(z)
}

// Update function loops over each of the meshes and calls their Update function.
func (r *Room) Update(dt float64) {
	r.animateDoor(dt)
	for i, _ := range r.meshes {
		r.meshes[i].Update(dt)
	}
}
func (r *Room) GetDoor() interfaces.Mesh {
	return r.meshes[2]
}
func (r *Room) GetFloor() interfaces.Mesh {
	return r.meshes[0]
}
