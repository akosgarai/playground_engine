package model

import (
	"math"
	"path"
	"runtime"

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
	position   mgl32.Vec3 // the position of the room (center point of the floor mesh)
	worldUp    mgl32.Vec3 // the up direction in the world.
	width      float32    // the length of the usable area in the x axis
	height     float32    // the lenght of the usable area in the y axis
	length     float32    // the length of the usable area in the z axis
	wallWidth  float32    // the width of the walls
	doorWidth  float32    // the width of the door that is on the right side of the front wall.
	doorHeight float32    // the height of the door that is on the right side of the front wall.
	rotationX  float32
	rotationY  float32
	rotationZ  float32
	wrapper    interfaces.GLWrapper
}

func NewRoomBuilder() *RoomBuilder {
	return &RoomBuilder{
		position:   mgl32.Vec3{0.0, 0.0, 0.0},
		worldUp:    mgl32.Vec3{0.0, 1.0, 0.0},
		width:      1.0,
		height:     1.0,
		length:     1.0,
		wallWidth:  0.005,
		doorWidth:  0.4,
		doorHeight: 0.6,
		rotationX:  0.0,
		rotationY:  0.0,
		rotationZ:  0.0,
		wrapper:    nil,
	}
}

// SetPosition sets the position.
func (b *RoomBuilder) SetPosition(p mgl32.Vec3) {
	b.position = p
}

// SetWorldUpDirection sets the worldUp.
func (b *RoomBuilder) SetWorldUpDirection(p mgl32.Vec3) {
	b.worldUp = p.Normalize()
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

// Build returns a material room that is constructed from the given setup.
func (b *RoomBuilder) Build() *Room {
	if b.wrapper == nil {
		panic("Wrapper is missing.")
	}
	// rotation calculation:
	// - get the mul4 product of the 3 component (translation transform).
	// - except the first one, where only the rotations has to be set,
	// both the rotation and position has to be transformed.
	rotationTransformationMatrix := mgl32.HomogRotate3DY(mgl32.DegToRad(b.rotationY)).Mul4(
		mgl32.HomogRotate3DX(mgl32.DegToRad(b.rotationX))).Mul4(
		mgl32.HomogRotate3DZ(mgl32.DegToRad(b.rotationZ)))
	// floor + ceiling
	basementSizeCuboid := cuboid.New(b.width, b.length, b.wallWidth)
	basementV, basementI, bo := basementSizeCuboid.MaterialMeshInput()

	floor := mesh.NewMaterialMesh(basementV, basementI, material.Chrome, b.wrapper)
	floor.SetPosition(b.position)
	floor.SetBoundingObject(bo)

	ceiling := mesh.NewMaterialMesh(basementV, basementI, material.Chrome, b.wrapper)
	ceilingPosition := mgl32.TransformCoordinate(mgl32.Vec3{0.0, b.height, 0.0}, rotationTransformationMatrix)
	ceiling.SetPosition(ceilingPosition)
	ceiling.SetParent(floor)
	ceiling.SetBoundingObject(bo)

	// back wall
	backWallSizeCuboid := cuboid.New(b.width, b.wallWidth, b.height)
	backWallV, backWallI, bo := backWallSizeCuboid.MaterialMeshInput()

	backWall := mesh.NewMaterialMesh(backWallV, backWallI, material.Chrome, b.wrapper)
	backWallPosition := mgl32.TransformCoordinate(mgl32.Vec3{0.0, (b.height - b.wallWidth) / 2, -(b.length + b.wallWidth) / 2}, rotationTransformationMatrix)
	backWall.SetPosition(backWallPosition)
	backWall.SetParent(floor)
	backWall.SetBoundingObject(bo)

	// side wall
	sideWallSizeCuboid := cuboid.New(b.wallWidth, b.length, b.height)
	sideWallV, sideWallI, bo := sideWallSizeCuboid.MaterialMeshInput()

	rightWall := mesh.NewMaterialMesh(sideWallV, sideWallI, material.Chrome, b.wrapper)
	rightWallPosition := mgl32.TransformCoordinate(mgl32.Vec3{-(b.width - b.wallWidth) / 2, b.height / 2, 0.0}, rotationTransformationMatrix)
	rightWall.SetPosition(rightWallPosition)
	rightWall.SetParent(floor)
	rightWall.SetBoundingObject(bo)

	leftWall := mesh.NewMaterialMesh(sideWallV, sideWallI, material.Chrome, b.wrapper)
	leftWallPosition := mgl32.TransformCoordinate(mgl32.Vec3{(b.width - b.wallWidth) / 2, b.height / 2, 0.0}, rotationTransformationMatrix)
	leftWall.SetPosition(leftWallPosition)
	leftWall.SetParent(floor)
	leftWall.SetBoundingObject(bo)

	// front wall parts

	frontCuboid := cuboid.New(b.width-b.doorWidth, b.wallWidth, b.height)
	V, I, bo := frontCuboid.MaterialMeshInput()

	frontWallMain := mesh.NewMaterialMesh(V, I, material.Chrome, b.wrapper)
	frontWallMainPosition := mgl32.TransformCoordinate(mgl32.Vec3{-b.doorWidth / 2, b.height / 2, (b.length - b.wallWidth) / 2}, rotationTransformationMatrix)
	frontWallMain.SetPosition(frontWallMainPosition)
	frontWallMain.SetParent(floor)
	frontWallMain.SetBoundingObject(bo)

	frontTopCuboid := cuboid.New(b.doorWidth, b.wallWidth, b.height-b.doorHeight)
	V, I, bo = frontTopCuboid.MaterialMeshInput()

	frontWallRest := mesh.NewMaterialMesh(V, I, material.Chrome, b.wrapper)
	frontWallRestPosition := mgl32.TransformCoordinate(mgl32.Vec3{((b.width - b.doorWidth) / 2), (b.height + b.doorHeight) / 2, (b.length - b.wallWidth) / 2}, rotationTransformationMatrix)
	frontWallRest.SetPosition(frontWallRestPosition)
	frontWallRest.SetParent(floor)
	frontWallRest.SetBoundingObject(bo)

	doorCuboid := cuboid.New(b.doorWidth, b.wallWidth, b.doorHeight)
	V, I, bo = doorCuboid.MaterialMeshInput()

	door := mesh.NewMaterialMesh(V, I, material.Bronze, b.wrapper)
	doorPosition := mgl32.TransformCoordinate(mgl32.Vec3{((b.width - b.doorWidth) / 2), b.doorHeight / 2, (b.length - b.wallWidth) / 2}, rotationTransformationMatrix)
	door.SetPosition(doorPosition)
	door.SetParent(floor)
	door.SetBoundingObject(bo)

	m := newCDModel()
	m.AddMesh(floor)
	m.AddMesh(ceiling)
	m.AddMesh(backWall)
	m.AddMesh(rightWall)
	m.AddMesh(leftWall)
	m.AddMesh(frontWallMain)
	m.AddMesh(frontWallRest)
	m.AddMesh(door)
	return &Room{BaseCollisionDetectionModel: *m, doorState: _DOOR_CLOSED, currentAnimationTime: 0}
}

type Room struct {
	BaseCollisionDetectionModel
	doorState            int
	currentAnimationTime float64
}

// NewMaterialRoom returns a Room that is based on material meshes.
// The position is the center point of the floor of the room.
// The initial orientation of the floor is the xy plane.
// The floor, the roof, the back wall, left wall, right wall are 1 * 1 * 0.05 cuboids.
// The front wall holds a door that is different colored.
func NewMaterialRoom(position mgl32.Vec3, glWrapper interfaces.GLWrapper) *Room {
	floorCuboid := cuboid.New(1.0, 1.0, 0.005)
	floorV, floorI, bo := floorCuboid.MaterialMeshInput()

	floor := mesh.NewMaterialMesh(floorV, floorI, material.Chrome, glWrapper)
	floor.SetPosition(position)
	floor.SetBoundingObject(bo)

	ceiling := mesh.NewMaterialMesh(floorV, floorI, material.Chrome, glWrapper)
	ceiling.SetPosition(mgl32.Vec3{0.0, 1.0, 0.0})
	ceiling.SetParent(floor)
	ceiling.SetBoundingObject(bo)

	backWall := mesh.NewMaterialMesh(floorV, floorI, material.Chrome, glWrapper)
	backWall.SetPosition(mgl32.Vec3{0.0, 0.5, -0.4975})
	backWall.RotateX(90)
	backWall.SetParent(floor)
	backWall.SetBoundingObject(bo)

	rightWall := mesh.NewMaterialMesh(floorV, floorI, material.Chrome, glWrapper)
	rightWall.SetPosition(mgl32.Vec3{-0.4975, 0.5, 0.0})
	rightWall.RotateZ(90)
	rightWall.SetParent(floor)
	rightWall.SetBoundingObject(bo)

	leftWall := mesh.NewMaterialMesh(floorV, floorI, material.Chrome, glWrapper)
	leftWall.SetPosition(mgl32.Vec3{0.4975, 0.5, 0.0})
	leftWall.RotateZ(90)
	leftWall.SetParent(floor)
	leftWall.SetBoundingObject(bo)

	// front wall parts

	frontCuboid := cuboid.New(0.6, 1.0, 0.005)
	V, I, bo := frontCuboid.MaterialMeshInput()
	frontWallMain := mesh.NewMaterialMesh(V, I, material.Chrome, glWrapper)
	frontWallMain.SetPosition(mgl32.Vec3{0.2, 0.5, 0.4975})
	frontWallMain.RotateX(90)
	frontWallMain.SetParent(floor)
	frontWallMain.SetBoundingObject(bo)
	frontTopCuboid := cuboid.New(0.4, 0.4, 0.005)
	V, I, bo = frontTopCuboid.MaterialMeshInput()
	frontWallRest := mesh.NewMaterialMesh(V, I, material.Chrome, glWrapper)
	frontWallRest.SetPosition(mgl32.Vec3{-0.3, 0.2, 0.4975})
	frontWallRest.RotateX(90)
	frontWallRest.SetParent(floor)
	frontWallRest.SetBoundingObject(bo)
	doorCuboid := cuboid.New(0.4, 0.005, 0.6)
	V, I, bo = doorCuboid.MaterialMeshInput()
	door := mesh.NewMaterialMesh(V, I, material.Bronze, glWrapper)
	door.SetPosition(mgl32.Vec3{-0.4975, 0.7, 0.6975})
	door.RotateY(90)
	door.SetParent(floor)
	door.SetBoundingObject(bo)

	m := newCDModel()
	m.AddMesh(floor)
	m.AddMesh(ceiling)
	m.AddMesh(backWall)
	m.AddMesh(rightWall)
	m.AddMesh(leftWall)
	m.AddMesh(frontWallMain)
	m.AddMesh(frontWallRest)
	m.AddMesh(door)
	return &Room{BaseCollisionDetectionModel: *m, doorState: _DOOR_OPENED, currentAnimationTime: 0}
}
func NewTextureRoom(position mgl32.Vec3, glWrapper interfaces.GLWrapper) *Room {
	var concreteTexture texture.Textures
	_, filename, _, _ := runtime.Caller(1)
	fileDir := path.Dir(filename)
	concreteTexture.AddTexture(fileDir+"/assets/concrete-wall.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	concreteTexture.AddTexture(fileDir+"/assets/concrete-wall.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)
	var doorTexture texture.Textures
	doorTexture.AddTexture(fileDir+"/assets/door.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	doorTexture.AddTexture(fileDir+"/assets/door.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	var windowTexture texture.Textures
	windowTexture.AddTexture(fileDir+"/assets/window.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", glWrapper)
	windowTexture.AddTexture(fileDir+"/assets/window.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", glWrapper)

	floorCuboid := cuboid.New(1.0, 1.0, 0.005)
	floorV, floorI, bo := floorCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)

	floor := mesh.NewTexturedMaterialMesh(floorV, floorI, concreteTexture, material.Chrome, glWrapper)
	floor.SetPosition(position)
	floor.SetBoundingObject(bo)

	ceiling := mesh.NewTexturedMaterialMesh(floorV, floorI, concreteTexture, material.Chrome, glWrapper)
	ceiling.SetPosition(mgl32.Vec3{0.0, 1.0, 0.0})
	ceiling.SetParent(floor)
	ceiling.SetBoundingObject(bo)

	backWall := mesh.NewTexturedMaterialMesh(floorV, floorI, concreteTexture, material.Chrome, glWrapper)
	backWall.SetPosition(mgl32.Vec3{0.0, 0.5, -0.4975})
	backWall.RotateX(90)
	backWall.SetParent(floor)
	backWall.SetBoundingObject(bo)

	rightWall := mesh.NewTexturedMesh(floorV, floorI, concreteTexture, glWrapper)
	rightWall.SetPosition(mgl32.Vec3{-0.4975, 0.5, 0.0})
	rightWall.RotateZ(90)
	rightWall.SetParent(floor)
	rightWall.SetBoundingObject(bo)

	leftWall := mesh.NewTexturedMesh(floorV, floorI, concreteTexture, glWrapper)
	leftWall.SetPosition(mgl32.Vec3{0.4975, 0.5, 0.0})
	leftWall.RotateZ(90)
	leftWall.SetParent(floor)
	leftWall.SetBoundingObject(bo)

	// front wall parts

	frontLongCuboid := cuboid.New(0.2, 1.0, 0.005)
	V, I, bo := frontLongCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	frontWallMain1 := mesh.NewTexturedMesh(V, I, concreteTexture, glWrapper)
	frontWallMain1.SetPosition(mgl32.Vec3{0.0, 0.5, 0.4975})
	frontWallMain1.RotateX(90)
	frontWallMain1.SetParent(floor)
	frontWallMain1.SetBoundingObject(bo)

	frontWallMain2 := mesh.NewTexturedMesh(V, I, concreteTexture, glWrapper)
	frontWallMain2.SetPosition(mgl32.Vec3{0.4, 0.5, 0.4975})
	frontWallMain2.RotateX(90)
	frontWallMain2.SetParent(floor)
	frontWallMain2.SetBoundingObject(bo)

	frontSmallCuboid := cuboid.New(0.2, 0.3, 0.005)
	V, I, bo = frontSmallCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	frontWallMain3 := mesh.NewTexturedMesh(V, I, concreteTexture, glWrapper)
	frontWallMain3.SetPosition(mgl32.Vec3{0.2, 0.15, 0.4975})
	frontWallMain3.RotateX(90)
	frontWallMain3.SetParent(floor)
	frontWallMain3.SetBoundingObject(bo)

	frontWallMain4 := mesh.NewTexturedMesh(V, I, concreteTexture, glWrapper)
	frontWallMain4.SetPosition(mgl32.Vec3{0.2, 0.85, 0.4975})
	frontWallMain4.RotateX(90)
	frontWallMain4.SetParent(floor)
	frontWallMain4.SetBoundingObject(bo)

	windowCuboid := cuboid.New(0.2, 0.4, 0.005)
	V, I, bo = windowCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	window := mesh.NewTexturedMesh(V, I, windowTexture, glWrapper)
	window.SetPosition(mgl32.Vec3{0.2, 0.5, 0.4975})
	window.RotateX(90)
	window.SetParent(floor)
	window.SetBoundingObject(bo)

	frontTopCuboid := cuboid.New(0.4, 0.4, 0.005)
	V, I, bo = frontTopCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	frontWallRest := mesh.NewTexturedMesh(V, I, concreteTexture, glWrapper)
	frontWallRest.SetPosition(mgl32.Vec3{-0.3, 0.2, 0.4975})
	frontWallRest.RotateX(90)
	frontWallRest.SetParent(floor)
	frontWallRest.SetBoundingObject(bo)
	doorCuboid := cuboid.New(0.4, 0.005, 0.6)
	V, I, bo = doorCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_SAME)
	door := mesh.NewTexturedMesh(V, I, doorTexture, glWrapper)
	door.SetPosition(mgl32.Vec3{-0.4975, 0.7, 0.6975})
	door.RotateY(-90)
	door.SetParent(floor)
	door.SetBoundingObject(bo)

	m := newCDModel()
	m.AddMesh(floor)
	m.AddMesh(ceiling)
	m.AddMesh(backWall)
	m.AddMesh(rightWall)
	m.AddMesh(leftWall)
	m.AddMesh(frontWallMain1)
	m.AddMesh(frontWallRest)
	m.AddMesh(door)
	m.AddMesh(frontWallMain2)
	m.AddMesh(frontWallMain3)
	m.AddMesh(frontWallMain4)
	m.AddMesh(window)
	m.SetTransparent(true)
	return &Room{BaseCollisionDetectionModel: *m, doorState: _DOOR_OPENED, currentAnimationTime: 0}
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
	maxDelta := math.Min(dt, doorAnimationTime-r.currentAnimationTime+dt)
	r.currentAnimationTime += maxDelta

	door := r.GetDoor()
	currentPos := door.GetPosition()
	var rotationDeg float32
	if r.doorState == _DOOR_OPENING {
		rotationDeg = float32(90.0 / doorAnimationTime * maxDelta)
	}
	if r.doorState == _DOOR_CLOSING {
		rotationDeg = float32(-90.0 / doorAnimationTime * maxDelta)
	}
	sinDeg := float32(math.Sin(float64(mgl32.DegToRad(rotationDeg))))
	cosDeg := float32(math.Cos(float64(mgl32.DegToRad(90 - rotationDeg))))
	door.SetPosition(mgl32.Vec3{currentPos.X() - sinDeg*0.125, currentPos.Y(), currentPos.Z() + cosDeg*0.125})
	door.RotateY(-rotationDeg)
	if r.currentAnimationTime >= doorAnimationTime {
		r.doorState = (r.doorState + 1) % 4
	}
}

// Update function loops over each of the meshes and calls their Update function.
func (r *Room) Update(dt float64) {
	r.animateDoor(dt)
	for i, _ := range r.meshes {
		r.meshes[i].Update(dt)
	}
}
func (r *Room) GetDoor() interfaces.Mesh {
	return r.meshes[7]
}
