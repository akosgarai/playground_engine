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

	frontCuboid := cuboid.New(0.6, 1.0, 0.005)
	V, I, bo := frontCuboid.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	frontWallMain := mesh.NewTexturedMesh(V, I, concreteTexture, glWrapper)
	frontWallMain.SetPosition(mgl32.Vec3{0.2, 0.5, 0.4975})
	frontWallMain.RotateX(90)
	frontWallMain.SetParent(floor)
	frontWallMain.SetBoundingObject(bo)
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
	m.AddMesh(frontWallMain)
	m.AddMesh(frontWallRest)
	m.AddMesh(door)
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
