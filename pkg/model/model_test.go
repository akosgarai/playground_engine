package model

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	InvalidFilename      = "not-existing-file.obj"
	ValidFilename        = "testdata/test_cube.obj"
	DefaultFormItemLabel = "form item label"
)

var (
	wrapperMock testhelper.GLWrapperMock
	shaderMock  testhelper.ShaderMock
)

func TestNew(t *testing.T) {
	model := New()
	if len(model.meshes) != 0 {
		t.Errorf("Invalid number of meshes. Instead of '0', we have '%d'.", len(model.meshes))
	}
}
func TestAddMesh(t *testing.T) {
	model := New()
	for i := 0; i < 10; i++ {
		msh := mesh.NewPointMesh(wrapperMock)
		if len(model.meshes) != i {
			t.Errorf("Invalid number of meshes before adding. Instead of '%d', we have '%d'.", i, len(model.meshes))
		}
		model.AddMesh(msh)
		if len(model.meshes) != i+1 {
			t.Errorf("Invalid number of meshes after adding. Instead of '%d', we have '%d'.", i+1, len(model.meshes))
		}
	}
}
func TestSetTransparent(t *testing.T) {
	model := New()
	if model.transparent != false {
		t.Error("Invalid initial transparent value. It should be false.")
	}
	model.SetTransparent(true)
	if model.transparent != true {
		t.Error("Invalid updated transparent value. It should be true.")
	}
	model.SetTransparent(false)
	if model.transparent != false {
		t.Error("Invalid updated transparent value. It should be false.")
	}
}
func TestIsTransparent(t *testing.T) {
	model := New()
	if model.IsTransparent() != false {
		t.Error("Invalid initial transparent value. It should be false.")
	}
	model.SetTransparent(true)
	if model.IsTransparent() != true {
		t.Error("Invalid updated transparent value. It should be true.")
	}
	model.SetTransparent(false)
	if model.IsTransparent() != false {
		t.Error("Invalid updated transparent value. It should be false.")
	}
}
func TestDraw(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Draw shouldn't have paniced.")
			}
		}()
		model := New()
		model.Draw(shaderMock)
		msh := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(msh)
		model.Draw(shaderMock)
		model.SetUniformFloat("name", float32(0.1))
		model.SetUniformVector("name", mgl32.Vec3{0.0, 1.0, 1.0})
		model.Draw(shaderMock)
	}()
}
func TestSetSpeed(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("SetSpeed shouldn't have paniced.")
			}
		}()
		model := New()
		for i := 0; i < 10; i++ {
			msh := mesh.NewPointMesh(wrapperMock)
			model.AddMesh(msh)
		}
		model.SetSpeed(2)
	}()
}
func TestUpdate(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Update shouldn't have paniced.")
			}
		}()
		delta := 10.0
		model := New()
		model.Update(delta)
		msh := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(msh)
		model.Update(delta)
	}()
}
func TestSetDirection(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("SetDirection shouldn't have paniced.")
			}
		}()
		model := New()
		dir := mgl32.Vec3{1.0, 0.0, 0.0}
		for i := 0; i < 10; i++ {
			msh := mesh.NewPointMesh(wrapperMock)
			model.AddMesh(msh)
		}
		model.SetDirection(dir)
	}()
}
func TestRotateX(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("RotateX shouldn't have paniced.")
			}
		}()
		model := New()
		parent := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(parent)

		for i := 0; i < 10; i++ {
			msh := mesh.NewPointMesh(wrapperMock)
			if i%2 == 1 {
				msh.SetParent(parent)
			}
			model.AddMesh(msh)
		}
		model.RotateX(90)
	}()
}
func TestRotateY(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("RotateY shouldn't have paniced.")
			}
		}()
		model := New()
		parent := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(parent)

		for i := 0; i < 10; i++ {
			msh := mesh.NewPointMesh(wrapperMock)
			if i%2 == 1 {
				msh.SetParent(parent)
			}
			model.AddMesh(msh)
		}
		model.RotateY(90)
	}()
}
func TestRotateZ(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("RotateZ shouldn't have paniced.")
			}
		}()
		model := New()
		parent := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(parent)

		for i := 0; i < 10; i++ {
			msh := mesh.NewPointMesh(wrapperMock)
			if i%2 == 1 {
				msh.SetParent(parent)
			}
			model.AddMesh(msh)
		}
		model.RotateZ(90)
	}()
}
func TestCollideTestWithSphere(t *testing.T) {
	model := New()
	meshWoBo := mesh.NewPointMesh(wrapperMock)
	model.AddMesh(meshWoBo)
	meshSphere := mesh.NewPointMesh(wrapperMock)
	sphereParams := make(map[string]float32)
	sphereParams["radius"] = float32(1.0)
	sphereBo := boundingobject.New("Sphere", sphereParams)
	meshSphere.SetBoundingObject(sphereBo)
	meshSphere.SetPosition(mgl32.Vec3{2.0, 2.0, 2.0})
	model.AddMesh(meshSphere)
	meshCube := mesh.NewPointMesh(wrapperMock)
	cubeParams := make(map[string]float32)
	cubeParams["width"] = float32(1.0)
	cubeParams["height"] = float32(1.0)
	cubeParams["length"] = float32(1.0)
	cubeBo := boundingobject.New("AABB", cubeParams)
	meshCube.SetBoundingObject(cubeBo)
	meshCube.SetPosition(mgl32.Vec3{-3, -3, -3})
	model.AddMesh(meshCube)
	parent := mesh.NewPointMesh(wrapperMock)
	parent.SetBoundingObject(sphereBo)
	parent.SetPosition(mgl32.Vec3{5.0, 5.0, 5.0})
	model.AddMesh(parent)
	child := mesh.NewPointMesh(wrapperMock)
	child.SetBoundingObject(sphereBo)
	child.SetPosition(mgl32.Vec3{5.0, 5.0, 5.0})
	child.SetParent(parent)
	model.AddMesh(child)

	testData := []struct {
		position  [3]float32
		radius    float32
		intersect bool
		msg       string
	}{
		{[3]float32{0, 0, 0}, 0.5, false, "Shouldn't intersect."},
		{[3]float32{2, 1, 2}, 1.0, true, "Should intersect with the sphere."},
		{[3]float32{-2, -2, -2}, 1.5, true, "Should intersect with the cube."},
	}

	for _, tt := range testData {
		base := coldet.NewBoundingSphere(tt.position, tt.radius)
		result := model.CollideTestWithSphere(base)
		if result != tt.intersect {
			t.Errorf("%s expected: '%v', result: '%v'.", tt.msg, tt.intersect, result)
		}
	}
}
func TestExport(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Export shouldn't have panic.")
			}
		}()
		model := New()
		model.Export("invalid-path")
	}()
}
func TestSetUniformFloat(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("SetUniformFloat shouldn't have panic.")
			}
		}()
		model := New()
		testData := []struct {
			key   string
			value float32
		}{
			{"testName", float32(0.0)},
			{"testName2", float32(1.0)},
			{"testName", float32(1.0)},
		}
		for _, tt := range testData {
			model.SetUniformFloat(tt.key, tt.value)
			if model.uniformFloat[tt.key] != tt.value {
				t.Errorf("Invalud uniform value for key '%s'. Instead of '%f', we have '%f'.", tt.key, tt.value, model.uniformFloat[tt.key])
			}
		}

	}()
}
func TestSetUniformVector(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("SetUniformVector shouldn't have panic.")
			}
		}()
		model := New()
		testData := []struct {
			key   string
			value mgl32.Vec3
		}{
			{"testName", mgl32.Vec3{0.0, 0.0, 0.0}},
			{"testName2", mgl32.Vec3{0.0, 1.0, 0.0}},
			{"testName", mgl32.Vec3{1.0, 0.0, 1.0}},
		}
		for _, tt := range testData {
			model.SetUniformVector(tt.key, tt.value)
			if model.uniformVector[tt.key] != tt.value {
				t.Errorf("Invalud uniform value for key '%s'. Instead of '%v', we have '%v'.", tt.key, tt.value, model.uniformVector[tt.key])
			}
		}

	}()
}
func TestClosestMeshTo(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("SetUniformVector shouldn't have panic.")
			}
		}()
		refPoint := mgl32.Vec3{0, 0, 0}
		model := New()
		// Get closest mesh without mesh
		msh, dist := model.ClosestMeshTo(refPoint)
		if dist != float32(math.MaxFloat32) {
			t.Errorf("Invalid distance. Instead of maxfloat it is '%f', msh: %#v", dist, msh)
		}
		// Get closest mesh with one mesh

		cubeParams := make(map[string]float32)
		cubeParams["width"] = float32(1.0)
		cubeParams["height"] = float32(1.0)
		cubeParams["length"] = float32(1.0)
		cubeBo := boundingobject.New("AABB", cubeParams)
		mshOne := mesh.NewPointMesh(wrapperMock)
		mshOne.SetBoundingObject(cubeBo)
		mshOne.SetPosition(mgl32.Vec3{0, 1, 0})
		model.AddMesh(mshOne)
		msh, dist = model.ClosestMeshTo(refPoint)
		if dist != float32(0.5) {
			t.Errorf("Invalid distance. Instead of 0.5 it is '%f', msh: %#v", dist, msh)
		}
		// Get closest mesh with multiple mesh
		cubeParams = make(map[string]float32)
		cubeParams["width"] = float32(1.0)
		cubeParams["height"] = float32(1.0)
		cubeParams["length"] = float32(1.0)
		cubeBo = boundingobject.New("AABB", cubeParams)

		mshTwo := mesh.NewPointMesh(wrapperMock)
		mshTwo.SetBoundingObject(cubeBo)
		mshTwo.SetPosition(mgl32.Vec3{0.75, 0, 0})
		model.AddMesh(mshTwo)
		msh, dist = model.ClosestMeshTo(refPoint)
		if dist != float32(0.25) {
			t.Errorf("Invalid distance. Instead of 0.5 it is '%f', msh: %#v", dist, msh)
		}
	}()
}
func TestClear(t *testing.T) {
	model := New()
	mshOne := mesh.NewPointMesh(wrapperMock)
	model.AddMesh(mshOne)
	model.Clear()
	if len(model.meshes) != 0 {
		t.Errorf("Invalid number of meshes after clear. Instead of 0, we have %d.", len(model.meshes))
	}
}
func TestBug(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	bottomPosition := mgl32.Vec3{-1.5, 0.0, 0.0}
	eyeBase := float32(0.6350853)
	eye1Position := mgl32.Vec3{eyeBase, eyeBase, eyeBase}
	eye2Position := mgl32.Vec3{eyeBase, eyeBase, -eyeBase}
	scale := mgl32.Vec3{1.0, 1.0, 1.0}

	bug := NewBug(position, scale, wrapperMock)

	if bug.GetBodyPosition() != position {
		t.Errorf("Invalid body position. Instead of '%v', we have '%v'.", position, bug.GetBodyPosition())
	}
	if bug.GetBottomPosition() != bottomPosition {
		t.Errorf("Invalid bottom position. Instead of '%v', we have '%v'.", bottomPosition, bug.GetBottomPosition())
	}
	if bug.GetEye1Position() != eye1Position {
		t.Errorf("Invalid eye1 position. Instead of '%v', we have '%v'.", eye1Position, bug.GetEye1Position())
	}
	if bug.GetEye2Position() != eye2Position {
		t.Errorf("Invalid eye2 position. Instead of '%v', we have '%v'.", eye2Position, bug.GetEye2Position())
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Update shouldn't have panic.")
			}
		}()
		bug.Update(10)
	}()
	testData := []struct {
		position  [3]float32
		radius    float32
		intersect bool
		msg       string
	}{
		{[3]float32{0, 0, 0}, 0.5, true, "Should intersect at x=-0.5."},
		{[3]float32{-1.5, 1.3, 0.0}, 1.0, true, "Should intersect at y=1."},
		{[3]float32{-2, -2, -2}, 1.5, false, "Shouldn't intersect."},
	}

	for _, tt := range testData {
		base := coldet.NewBoundingSphere(tt.position, tt.radius)
		result := bug.CollideTestWithSphere(base)
		if result != tt.intersect {
			t.Errorf("%s expected: '%v', result: '%v'.", tt.msg, tt.intersect, result)
		}
	}
}
func TestMaterialStreetLamp(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	bulbPosition := mgl32.Vec3{0.9349, 3.0, 2.98}
	polePosition := mgl32.Vec3{0.0, 3.0, 0.0}
	topPosition := mgl32.Vec3{1.1, 3.0, 6.4}
	scale := float32(6.0)

	lamp := NewMaterialStreetLamp(position, scale, wrapperMock)

	if lamp.GetPolePosition() != polePosition {
		t.Errorf("Invalid pole position. Instead of '%v', we have '%v'.", polePosition, lamp.GetPolePosition())
	}
	if !lamp.GetTopPosition().ApproxEqualThreshold(topPosition, 0.0001) {
		t.Errorf("Invalid top position. Instead of '%v', we have '%v'.", topPosition, lamp.GetTopPosition())
	}
	if !lamp.GetBulbPosition().ApproxEqualThreshold(bulbPosition, 0.0001) {
		t.Errorf("Invalid bulb position. Instead of '%v', we have '%v'.", bulbPosition, lamp.GetBulbPosition())
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Update shouldn't have panic.")
			}
		}()
		lamp.Update(10)
	}()
	testData := []struct {
		position  [3]float32
		radius    float32
		intersect bool
		msg       string
	}{
		{[3]float32{-0.6, 3, 0}, 0.5, true, "Should intersect at x."},
		{[3]float32{-2, -2, -2}, 1.5, false, "Shouldn't intersect."},
	}

	for _, tt := range testData {
		base := coldet.NewBoundingSphere(tt.position, tt.radius)
		result := lamp.CollideTestWithSphere(base)
		if result != tt.intersect {
			t.Errorf("%s expected: '%v', result: '%v'.", tt.msg, tt.intersect, result)
		}
	}
}
func TestTexturedStreetLamp(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	bulbPosition := mgl32.Vec3{0.9, 3.035, 3.0}
	polePosition := mgl32.Vec3{0.0, 3.0, 0.0}
	topPosition := mgl32.Vec3{3.55, 3.55, 3.0}
	scale := float32(6.0)

	lamp := NewTexturedStreetLamp(position, scale, wrapperMock)

	if lamp.GetPolePosition() != polePosition {
		t.Errorf("Invalid pole position. Instead of '%v', we have '%v'.", polePosition, lamp.GetPolePosition())
	}
	if !lamp.GetTopPosition().ApproxEqualThreshold(topPosition, 0.0001) {
		t.Errorf("Invalid top position. Instead of '%v', we have '%v'.", topPosition, lamp.GetTopPosition())
	}
	if !lamp.GetBulbPosition().ApproxEqualThreshold(bulbPosition, 0.0001) {
		t.Errorf("Invalid bulb position. Instead of '%v', we have '%v'.", bulbPosition, lamp.GetBulbPosition())
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Update shouldn't have panic.")
			}
		}()
		lamp.Update(10)
	}()
}
func CheckDefaultRoomOptions(room *Room, t *testing.T) {
	doorPosition := mgl32.Vec3{-0.4975, 0.7, 0.6975}
	if room.GetDoor().GetPosition() != doorPosition {
		t.Errorf("Invalid door position. Instead of '%v', we have '%v'.", doorPosition, room.GetDoor().GetPosition())
	}
	if room.doorState != _DOOR_OPENED {
		t.Errorf("Invalid initial door state. Instead of '%d', we have '%d'.", _DOOR_OPENED, room.doorState)
	}
	room.PushDoorState()
	if room.doorState != _DOOR_CLOSING {
		t.Errorf("Invalid next door state. Instead of '%d', we have '%d'.", _DOOR_CLOSING, room.doorState)
	}
	room.PushDoorState()
	if room.doorState != _DOOR_CLOSING {
		t.Errorf("Invalid door state. Instead of '%d', we have '%d'.", _DOOR_CLOSING, room.doorState)
	}

	room.doorState = _DOOR_CLOSED
	room.PushDoorState()
	if room.doorState != _DOOR_OPENING {
		t.Errorf("Invalid next door state. Instead of '%d', we have '%d'.", _DOOR_OPENING, room.doorState)
	}
	room.PushDoorState()
	if room.doorState != _DOOR_OPENING {
		t.Errorf("Invalid next door state. Instead of '%d', we have '%d'.", _DOOR_OPENING, room.doorState)
	}
	if room.currentAnimationTime != 0.0 {
		t.Errorf("Invalid initial animation time. Instead of '0.0', it is '%f'.", room.currentAnimationTime)
	}
	room.animateDoor(100)
	if room.currentAnimationTime != 100.0 {
		t.Errorf("Invalid animation time. Instead of '100.0', it is '%f'.", room.currentAnimationTime)
	}
	room.animateDoor(950)
	if room.doorState != _DOOR_OPENED {
		t.Errorf("Invalid next door state. Instead of '%d', we have '%d'.", _DOOR_OPENED, room.doorState)
	}
	room.PushDoorState()
	room.animateDoor(100)
	if room.currentAnimationTime != 100.0 {
		t.Errorf("Invalid animation time. Instead of '100.0', it is '%f'.", room.currentAnimationTime)
	}
	room.animateDoor(950)
	if room.doorState != _DOOR_CLOSED {
		t.Errorf("Invalid next door state. Instead of '%d', we have '%d'.", _DOOR_CLOSED, room.doorState)
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Update shouldn't have panic.")
			}
		}()
		room.Update(100)
	}()
}
func TestMaterialRoom(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}

	room := NewMaterialRoom(position, wrapperMock)

	CheckDefaultRoomOptions(room, t)
}
func TestTexturedRoom(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}

	room := NewTextureRoom(position, wrapperMock)

	CheckDefaultRoomOptions(room, t)
}
func TestTerrainBuilderNewTerrainBuilder(t *testing.T) {
	terr := NewTerrainBuilder()
	if terr.width != defaultTerrainWidth {
		t.Errorf("Invalid default width. Instead of '%d', we have '%d'.\n", defaultTerrainWidth, terr.width)
	}
	if terr.length != defaultTerrainLength {
		t.Errorf("Invalid default length. Instead of '%d', we have '%d'.\n", defaultTerrainLength, terr.length)
	}
	if terr.iterations != defaultIterations {
		t.Errorf("Invalid default iterations. Instead of '%d', we have '%d'.\n", defaultIterations, terr.iterations)
	}
	if terr.minH != defaultMinHeight {
		t.Errorf("Invalid default minH. Instead of '%f', we have '%f'.\n", defaultMinHeight, terr.minH)
	}
	if terr.maxH != defaultMaxHeight {
		t.Errorf("Invalid default maxH. Instead of '%f', we have '%f'.\n", defaultMaxHeight, terr.maxH)
	}
	if terr.seed != defaultSeed {
		t.Errorf("Invalid default seed. Instead of '%d', we have '%d'.\n", defaultSeed, terr.seed)
	}
	if terr.minHIsDefault != false {
		t.Error("Invalid  minHIsDefault. Not false")
	}
	if terr.cliffProbability != 0 {
		t.Errorf("Invalid default cliffProbability. Instead of '%d', we have '%d'.\n", 0, terr.cliffProbability)
	}
	if terr.peakProbability != 0 {
		t.Errorf("Invalid default peakProbability. Instead of '%d', we have '%d'.\n", 0, terr.peakProbability)
	}
}
func TestTerrainBuilderSetWidth(t *testing.T) {
	width := defaultTerrainWidth + 1
	terr := NewTerrainBuilder()
	terr.SetWidth(width)
	if terr.width != width {
		t.Errorf("Invalid width. Instead of '%d', we have '%d'.", width, terr.width)
	}
}
func TestTerrainBuilderSetLength(t *testing.T) {
	length := defaultTerrainLength + 1
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	if terr.length != length {
		t.Errorf("Invalid length. Instead of '%d', we have '%d'.", length, terr.length)
	}
}
func TestTerrainBuilderSetIterations(t *testing.T) {
	iterations := defaultIterations + 1
	terr := NewTerrainBuilder()
	terr.SetIterations(iterations)
	if terr.iterations != iterations {
		t.Errorf("Invalid iterations. Instead of '%d', we have '%d'.", iterations, terr.iterations)
	}
}
func TestTerrainBuilderSetMinHeight(t *testing.T) {
	height := defaultMinHeight + 1.0
	terr := NewTerrainBuilder()
	terr.SetMinHeight(height)
	if terr.minH != height {
		t.Errorf("Invalid minH. Instead of '%f', we have '%f'.", height, terr.minH)
	}
}
func TestTerrainBuilderSetMaxHeight(t *testing.T) {
	height := defaultMaxHeight + 1.0
	terr := NewTerrainBuilder()
	terr.SetMaxHeight(height)
	if terr.maxH != height {
		t.Errorf("Invalid maxH. Instead of '%f', we have '%f'.", height, terr.maxH)
	}
}
func TestTerrainBuilderSetSeed(t *testing.T) {
	seed := defaultSeed + 1
	terr := NewTerrainBuilder()
	terr.SetSeed(seed)
	if terr.seed != seed {
		t.Errorf("Invalid seed. Instead of '%d', we have '%d'.", seed, terr.seed)
	}
}
func TestTerrainBuilderSetScale(t *testing.T) {
	scale := mgl32.Vec3{2, 2, 2}
	terr := NewTerrainBuilder()
	terr.SetScale(scale)
	if terr.scale != scale {
		t.Errorf("Invalid scale. Instead of '%v', we have '%v'.", scale, terr.scale)
	}
}
func TestTerrainBuilderSetPosition(t *testing.T) {
	position := mgl32.Vec3{2, 2, 2}
	terr := NewTerrainBuilder()
	terr.SetPosition(position)
	if terr.position != position {
		t.Errorf("Invalid position. Instead of '%v', we have '%v'.", position, terr.position)
	}
}
func TestTerrainBuilderRandomSeed(t *testing.T) {
	before := time.Now().UnixNano()
	terr := NewTerrainBuilder()
	terr.RandomSeed()
	after := time.Now().UnixNano()
	if terr.seed < before || terr.seed > after {
		t.Errorf("Invalid random seed '%d'. It supposed to be beetween '%d' and '%d'.", terr.seed, before, after)
	}
}
func TestTerrainBuilderSetPeakProbability(t *testing.T) {
	prob := 1
	terr := NewTerrainBuilder()
	terr.SetPeakProbability(prob)
	if terr.peakProbability != prob {
		t.Errorf("Invalid peak prob. Instead of '%d', we have '%d'.", prob, terr.peakProbability)
	}
}
func TestTerrainBuilderSetCliffProbability(t *testing.T) {
	prob := 1
	terr := NewTerrainBuilder()
	terr.SetCliffProbability(prob)
	if terr.cliffProbability != prob {
		t.Errorf("Invalid cliff prob. Instead of '%d', we have '%d'.", prob, terr.cliffProbability)
	}
}
func TestTerrainBuilderMinHeightIsDefault(t *testing.T) {
	terr := NewTerrainBuilder()
	terr.MinHeightIsDefault(true)
	if !terr.minHIsDefault {
		t.Error("Invalid minHIsDefault. It should be true.")
	}
	terr.MinHeightIsDefault(false)
	if terr.minHIsDefault {
		t.Error("Invalid minHIsDefault. It should be false.")
	}
}
func TestTerrainBuilderSetGlWrapper(t *testing.T) {
	var wrapper testhelper.GLWrapperMock
	terr := NewTerrainBuilder()
	terr.SetGlWrapper(wrapper)
	if terr.wrapper != wrapper {
		t.Error("Invalid gl wrapper")
	}
}
func TestTerrainBuilderSurfaceTextureGrass(t *testing.T) {
	var wrapper testhelper.GLWrapperMock
	terr := NewTerrainBuilder()
	terr.SetGlWrapper(wrapper)
	if len(terr.tex) != 0 {
		t.Errorf("Invalid texture length. Instead of '0', we have '%d'.", len(terr.tex))
	}
	terr.SurfaceTextureGrass()
	if len(terr.tex) != 2 {
		t.Errorf("Invalid texture length. Instead of '2', we have '%d'.", len(terr.tex))
	}
}
func TestTerrainBuilderInitHeightMap(t *testing.T) {
	length := 2
	width := 2
	minH := float32(-2)
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	terr.SetWidth(width)
	terr.SetMinHeight(minH)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, 0.0)
	for l := 0; l < terr.length; l++ {
		for w := 0; w < terr.width; w++ {
			if terr.heightMap[l][w] != 0.0 {
				t.Errorf("Invalid heightMap. Instead of '0.0', we have '%f' at (%d,%d)", terr.heightMap[l][w], l, w)
			}
		}
	}
	terr.MinHeightIsDefault(true)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, terr.minH)
	for l := 0; l < terr.length; l++ {
		for w := 0; w < terr.width; w++ {
			if terr.heightMap[l][w] != -2.0 {
				t.Errorf("Invalid heightMap. Instead of '-2.0', we have '%f' at (%d,%d)", terr.heightMap[l][w], l, w)
			}
		}
	}
}
func TestTerrainBuilderBuildHeightMap(t *testing.T) {
	length := 4
	width := 4
	iteration := 10
	minH := float32(-1.0)
	maxH := float32(3.0)
	seed := int64(0)
	peakProb := 5
	cliffProb := 5
	expected := [][]float32{
		{0, 0, 1.8, 1, 0},
		{0, -0.19999999, 0, 1.4000001, 0},
		{-0.19999999, -0.6, -0.19999999, 2.6000001, 0},
		{-0.19999999, -0.19999999, -0.19999999, 0, 0},
		{0, 0, -0.6, 0, -1},
	}
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	terr.SetWidth(width)
	terr.SetMinHeight(minH)
	terr.SetMaxHeight(maxH)
	terr.SetIterations(iteration)
	terr.SetSeed(seed)
	terr.SetPeakProbability(peakProb)
	terr.SetCliffProbability(cliffProb)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, 0.0)
	terr.buildTerrainHeightMap()
	if !reflect.DeepEqual(terr.heightMap, expected) {
		t.Error("Invalid heightmap")
		t.Log(terr.heightMap)
		t.Log(expected)
	}
	terr.MinHeightIsDefault(true)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, terr.minH)
	terr.buildTerrainHeightMap()
}
func TestTerrainBuilderAdjacentElevation(t *testing.T) {
	length := 4
	width := 4
	iteration := 10
	minH := float32(-1.0)
	maxH := float32(3.0)
	seed := int64(0)
	peakProb := 5
	cliffProb := 5
	expected := [][]float32{
		{0, 0, 1.8, 1, 0},
		{0, -0.19999999, 0, 1.4000001, 0},
		{-0.19999999, -0.6, -0.19999999, 2.6000001, 0},
		{-0.19999999, -0.19999999, -0.19999999, 0, 0},
		{0, 0, -0.6, 0, -1},
	}
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	terr.SetWidth(width)
	terr.SetMinHeight(minH)
	terr.SetMaxHeight(maxH)
	terr.SetIterations(iteration)
	terr.SetSeed(seed)
	terr.SetPeakProbability(peakProb)
	terr.SetCliffProbability(cliffProb)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, 0.0)
	terr.buildTerrainHeightMap()
	if !reflect.DeepEqual(terr.heightMap, expected) {
		t.Error("Invalid heightmap")
		t.Log(terr.heightMap)
		t.Log(expected)
	}
}
func TestTerrainBuilderVertices(t *testing.T) {
	length := 4
	width := 4
	iteration := 10
	minH := float32(-1.0)
	maxH := float32(3.0)
	seed := int64(0)
	peakProb := 5
	cliffProb := 5
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	terr.SetWidth(width)
	terr.SetMinHeight(minH)
	terr.SetMaxHeight(maxH)
	terr.SetIterations(iteration)
	terr.SetSeed(seed)
	terr.SetPeakProbability(peakProb)
	terr.SetCliffProbability(cliffProb)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, 0.0)
	terr.buildTerrainHeightMap()
	v := terr.vertices(terr.width, terr.length, 1, 1, terr.heightMap)
	if len(v) != (length+1)*(width+1) {
		t.Errorf("Invalid vertices length. Instead of '%d', we have '%d'.", length*width, len(v))
	}
}
func TestTerrainBuilderIndices(t *testing.T) {
	length := 5
	width := 5
	iteration := 10
	minH := float32(-1.0)
	maxH := float32(3.0)
	seed := int64(0)
	peakProb := 5
	cliffProb := 5
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	terr.SetWidth(width)
	terr.SetMinHeight(minH)
	terr.SetMaxHeight(maxH)
	terr.SetIterations(iteration)
	terr.SetSeed(seed)
	terr.SetPeakProbability(peakProb)
	terr.SetCliffProbability(cliffProb)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, 0.0)
	terr.buildTerrainHeightMap()
	i := terr.indices(terr.width, terr.length)
	if len(i) != length*width*6 {
		t.Errorf("Invalid indices length. Instead of '%d', we have '%d'.", length*width*6, len(i))
	}
}
func TestTerrainBuilderBuild(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing textures.")
			}
		}()
		length := 5
		width := 5
		iteration := 10
		minH := float32(-1.0)
		maxH := float32(3.0)
		seed := int64(0)
		peakProb := 5
		cliffProb := 5
		tb := NewTerrainBuilder()
		tb.SetLength(length)
		tb.SetWidth(width)
		tb.SetMinHeight(minH)
		tb.SetMaxHeight(maxH)
		tb.SetIterations(iteration)
		tb.SetSeed(seed)
		tb.SetPeakProbability(peakProb)
		tb.SetCliffProbability(cliffProb)
		tb.SetDebugMode(true)
		terr := tb.Build()
		t.Log(terr)
	}()
}
func TestTerrainBuilderBuildTerrain(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing textures.")
			}
		}()
		length := 5
		width := 5
		iteration := 10
		minH := float32(-1.0)
		maxH := float32(3.0)
		seed := int64(0)
		peakProb := 5
		cliffProb := 5
		tb := NewTerrainBuilder()
		tb.SetLength(length)
		tb.SetWidth(width)
		tb.SetMinHeight(minH)
		tb.SetMaxHeight(maxH)
		tb.SetIterations(iteration)
		tb.SetSeed(seed)
		tb.SetPeakProbability(peakProb)
		tb.SetCliffProbability(cliffProb)
		tb.SetDebugMode(true)
		terr := tb.buildTerrain()
		t.Log(terr)
	}()
}
func TestTerrainBuilderBuildWithLiquid(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing textures.")
			}
		}()
		length := 5
		width := 5
		iteration := 10
		minH := float32(-1.0)
		maxH := float32(3.0)
		seed := int64(0)
		peakProb := 5
		cliffProb := 5
		tb := NewTerrainBuilder()
		tb.SetLength(length)
		tb.SetWidth(width)
		tb.SetMinHeight(minH)
		tb.SetMaxHeight(maxH)
		tb.SetIterations(iteration)
		tb.SetSeed(seed)
		tb.SetPeakProbability(peakProb)
		tb.SetCliffProbability(cliffProb)
		tb.SetDebugMode(true)
		terr, water := tb.BuildWithLiquid()
		t.Log(terr, water)
	}()
}
func TestTerrainBuilderBuildLiquid(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing textures.")
			}
		}()
		length := 5
		width := 5
		iteration := 10
		minH := float32(-1.0)
		maxH := float32(3.0)
		seed := int64(0)
		peakProb := 5
		cliffProb := 5
		tb := NewTerrainBuilder()
		tb.SetLength(length)
		tb.SetWidth(width)
		tb.SetMinHeight(minH)
		tb.SetMaxHeight(maxH)
		tb.SetIterations(iteration)
		tb.SetSeed(seed)
		tb.SetPeakProbability(peakProb)
		tb.SetCliffProbability(cliffProb)
		tb.SetDebugMode(true)
		water := tb.buildLiquid()
		t.Log(water)
	}()
}
func TestTerrainHeightAtPos(t *testing.T) {
	tb := NewTerrainBuilder()
	tb.SetScale(mgl32.Vec3{2, 1, 2})
	tb.SetGlWrapper(wrapperMock)
	tb.SurfaceTextureGrass()
	terrain := tb.Build()
	terrain.heightMap = [][]float32{
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	testData := []struct {
		position  mgl32.Vec3
		height    float32
		debugMode bool
		err       error
	}{
		{mgl32.Vec3{15, 0, 15}, -1, false, ErrorNotAboveTheSurface},
		{mgl32.Vec3{11, 0, 10}, -1, false, ErrorNotAboveTheSurface},
		{mgl32.Vec3{10.2, 0, 10}, -1, false, ErrorNotAboveTheSurface},
		{mgl32.Vec3{10.000002, 0, 10}, -1, false, ErrorNotAboveTheSurface},
		{mgl32.Vec3{2, 0, 8}, 0, false, nil},
		{mgl32.Vec3{0, 0, 0}, 0, false, nil},
		{mgl32.Vec3{-10, 0, -10}, 0, false, nil},
		{mgl32.Vec3{-10, 0, -8}, 0, false, nil},
		{mgl32.Vec3{-9, 0, -9}, 0.25, false, nil},
		{mgl32.Vec3{-8, 0, -8.01}, 0.995, false, nil},
		{mgl32.Vec3{-8.01, 0, -8.01}, 0.990025, false, nil},
		{mgl32.Vec3{-8.01, 0, -8}, 0.995, false, nil},
		{mgl32.Vec3{-8, 0, -8}, 1.0, false, nil},
		{mgl32.Vec3{-7.99, 0, -8}, 1.0, false, nil},
		{mgl32.Vec3{-8, 0, -7.99}, 1.0, false, nil},
		{mgl32.Vec3{-7.99, 0, -7.99}, 1.0, false, nil},
		{mgl32.Vec3{-7, 0, -7}, 1.0, false, nil},
		{mgl32.Vec3{-6, 0, -6}, 1.0, false, nil},
		{mgl32.Vec3{-5, 0, -5}, 0.25, true, nil},
		{mgl32.Vec3{-10, 0, -7}, 0.0, false, nil},
	}
	for _, v := range testData {
		terrain.debugMode = v.debugMode
		x, err := terrain.HeightAtPos(v.position)
		if !testhelper.Float32ApproxEqual(x, v.height, 0.00001) || err != v.err {
			t.Errorf("Invalid height results.\n - At position '%v',\n\texpected height: '%f', error: '%v',\n\tgiven height: '%f', error: '%v'.", v.position, v.height, v.err, x, err)
			t.Log(terrain.heightMap)
		}
	}
}
func TestTerrainCollideTestWithSphere(t *testing.T) {
	tb := NewTerrainBuilder()
	tb.SetScale(mgl32.Vec3{2, 1, 2})
	tb.SetGlWrapper(wrapperMock)
	tb.SurfaceTextureGrass()
	terrain := tb.Build()
	terrain.heightMap = [][]float32{
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	testData := []struct {
		position [3]float32
		radius   float32
		collide  bool
	}{
		{[3]float32{15, 0, 15}, 1, false},
		{[3]float32{11, 0, 10}, 1, false},
		{[3]float32{10.2, 0, 10}, 1, false},
		{[3]float32{10.000002, 0, 10}, 1, false},
		{[3]float32{2, 1.1, 8}, 1, false},
		{[3]float32{2, 0.9, 8}, 1, true},
	}
	for _, v := range testData {
		bs := coldet.NewBoundingSphere(v.position, v.radius)
		result := terrain.CollideTestWithSphere(bs)
		if result != v.collide {
			t.Errorf("Invalid collision.\nPosition:\t'%v'\nRadius:\t'%f'\nExpected:\t'%v'", v.position, v.radius, v.collide)
		}
	}
	terrain.GetTerrain().SetPosition(mgl32.Vec3{0, 1, 0})
	testData = []struct {
		position [3]float32
		radius   float32
		collide  bool
	}{
		{[3]float32{15, 0, 15}, 1, false},
		{[3]float32{11, 0, 10}, 1, false},
		{[3]float32{10.2, 0, 10}, 1, false},
		{[3]float32{10.000002, 0, 10}, 1, false},
		{[3]float32{2, 2.1, 8}, 1, false},
		{[3]float32{2, 1.9, 8}, 1, true},
	}
	for _, v := range testData {
		bs := coldet.NewBoundingSphere(v.position, v.radius)
		result := terrain.CollideTestWithSphere(bs)
		if result != v.collide {
			t.Errorf("Invalid collision.\nPosition:\t'%v'\nRadius:\t'%f'\nExpected:\t'%v'", v.position, v.radius, v.collide)
		}
	}
}
func TestLiquidCollideTestWithSphere(t *testing.T) {
	tb := NewTerrainBuilder()
	tb.SetScale(mgl32.Vec3{2, 1, 2})
	tb.SetGlWrapper(wrapperMock)
	tb.SurfaceTextureGrass()
	tb.LiquidTextureWater()
	_, water := tb.BuildWithLiquid()
	testData := []struct {
		position [3]float32
		radius   float32
		collide  bool
	}{
		{[3]float32{15, 0, 15}, 1, false},
		{[3]float32{11, 0, 10}, 1, false},
		{[3]float32{10.2, 0, 10}, 1, false},
		{[3]float32{10.000002, 0, 10}, 1, false},
		{[3]float32{2, 2.1, 8}, 1, false},
		{[3]float32{2, 1.9, 8}, 1, false},
	}
	for _, v := range testData {
		bs := coldet.NewBoundingSphere(v.position, v.radius)
		result := water.CollideTestWithSphere(bs)
		if result != v.collide {
			t.Errorf("Invalid collision.\nPosition:\t'%v'\nRadius:\t'%f'\nExpected:\t'%v'", v.position, v.radius, v.collide)
		}
	}
}
func TestTerrainBuilderSetLiquidEta(t *testing.T) {
	eta := float32(1)
	terr := NewTerrainBuilder()
	terr.SetLiquidEta(eta)
	if terr.liquidEta != eta {
		t.Errorf("Invalid liquid eta. Instead of '%f', we have '%f'.", eta, terr.liquidEta)
	}
}
func TestTerrainBuilderSetLiquidAmplitude(t *testing.T) {
	ampl := float32(1)
	terr := NewTerrainBuilder()
	terr.SetLiquidAmplitude(ampl)
	if terr.liquidAmplitude != ampl {
		t.Errorf("Invalid liquid amplitude. Instead of '%f', we have '%f'.", ampl, terr.liquidAmplitude)
	}
}
func TestTerrainBuilderSetLiquidFrequency(t *testing.T) {
	freq := float32(1)
	terr := NewTerrainBuilder()
	terr.SetLiquidFrequency(freq)
	if terr.liquidFrequency != freq {
		t.Errorf("Invalid liquid frequency. Instead of '%f', we have '%f'.", freq, terr.liquidFrequency)
	}
}
func TestTerrainBuilderSetLiquidWaterLevel(t *testing.T) {
	wl := float32(1)
	terr := NewTerrainBuilder()
	terr.SetLiquidWaterLevel(wl)
	if terr.liquidWaterLevel != wl {
		t.Errorf("Invalid liquid waterlevel. Instead of '%f', we have '%f'.", wl, terr.liquidWaterLevel)
	}
}
func TestTerrainBuilderSetLiquidDetailMultiplier(t *testing.T) {
	dm := 3
	terr := NewTerrainBuilder()
	terr.SetLiquidDetailMultiplier(dm)
	if terr.liquidDetailMultiplier != dm {
		t.Errorf("Invalid liquid detail multiplier. Instead of '%d', we have '%d'.", dm, terr.liquidDetailMultiplier)
	}
}
func TestCharsetLoad(t *testing.T) {
	_, err := LoadCharsetDebug("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during debug load: %s\n", err.Error())
	}
	_, err = LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during load: %s\n", err.Error())
	}
}
func TestCharsetPrintTo(t *testing.T) {
	cols := []mgl32.Vec3{mgl32.Vec3{0, 0, 1}}
	fonts, err := LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during load: %s\n", err.Error())
	}
	msh := mesh.NewPointMesh(wrapperMock)
	fonts.PrintTo("Hello", 0, 0, 0.0, 1.0, wrapperMock, msh, cols)
	fonts.Debug = true
	fonts.PrintTo("Hello", 0, 0, 0.0, 1.0, wrapperMock, msh, cols)
}
func TestCharsetCleanSurface(t *testing.T) {
	cols := []mgl32.Vec3{mgl32.Vec3{0, 0, 1}}
	fonts, err := LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during load: %s\n", err.Error())
	}
	msh := mesh.NewPointMesh(wrapperMock)
	fonts.PrintTo("Hello", 0, 0, 0.0, 1.0, wrapperMock, msh, cols)
	if len(fonts.meshes) != 5 {
		t.Errorf("Invalid number of meshes. Instead of '%d', we have '%d'.", 5, len(fonts.meshes))
	}
	fonts.CleanSurface(msh)
	if len(fonts.meshes) != 0 {
		t.Errorf("Invalid number of meshes. Instead of '%d', we have '%d'.", 0, len(fonts.meshes))
	}
}
func TestCharsetTextWidth(t *testing.T) {
	fonts, err := LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during load: %s\n", err.Error())
	}
	testData := []struct {
		text  string
		scale float32
		width float32
	}{
		{"a", 1, 25},
		{"a", 2, 50},
		{"a", 0.5, 12.5},
		{"b", 1, 22},
		{"1", 1, 15},
		{"b1", 1, 37},
	}
	for _, tt := range testData {
		width := fonts.TextWidth(tt.text, tt.scale)
		if width != tt.width {
			t.Errorf("Invalid text width for '%s'. Instead of '%f', we have '%f'.", tt.text, tt.width, width)
		}
	}
}
func testFormItemBool(t *testing.T) *FormItemBool {
	mat := material.Chrome
	pos := mgl32.Vec3{0, 0, 0}
	fi := NewFormItemBool(DefaultFormItemLabel, mat, pos, wrapperMock)

	if fi.label != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.label)
	}
	return fi
}
func TestNewFormItemBool(t *testing.T) {
	_ = testFormItemBool(t)
}
func TestFormItemBoolGetLabel(t *testing.T) {
	fi := testFormItemBool(t)
	if fi.GetLabel() != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.GetLabel())
	}

}
func TestFormItemBoolGetValue(t *testing.T) {
	fi := testFormItemBool(t)
	val := true
	fi.value = val
	if fi.GetValue() != val {
		t.Errorf("Invalid form item value. Instead of '%v', it is '%v'.", val, fi.GetValue())
	}
}
func TestFormItemBoolSetValue(t *testing.T) {
	fi := testFormItemBool(t)
	val := true
	fi.value = val
	fi.SetValue(!val)
	if fi.GetValue() != !val {
		t.Errorf("Invalid form item value. Instead of '%v', it is '%v'.", !val, fi.GetValue())
	}
}
func TestFormItemBoolGetSurface(t *testing.T) {
	fi := testFormItemBool(t)
	if fi.GetSurface() != fi.meshes[0] {
		t.Error("Invalid surface mesh")
	}
}
func TestFormItemBoolGetLight(t *testing.T) {
	fi := testFormItemBool(t)
	if fi.GetLight() != fi.meshes[1] {
		t.Error("Invalid light mesh")
	}
}
func testFormItemInt(t *testing.T) *FormItemInt {
	mat := material.Chrome
	pos := mgl32.Vec3{0, 0, 0}
	fi := NewFormItemInt(DefaultFormItemLabel, mat, pos, wrapperMock)

	if fi.label != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.label)
	}
	return fi
}
func TestNewFormItemInt(t *testing.T) {
	_ = testFormItemInt(t)
}
func TestFormItemIntGetLabel(t *testing.T) {
	fi := testFormItemInt(t)
	if fi.GetLabel() != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.GetLabel())
	}

}
func TestFormItemIntGetValue(t *testing.T) {
	fi := testFormItemInt(t)
	val := 3
	fi.value = val
	if fi.GetValue() != val {
		t.Errorf("Invalid form item value. Instead of '%d', it is '%d'.", val, fi.GetValue())
	}
}
func TestFormItemIntSetValue(t *testing.T) {
	fi := testFormItemInt(t)
	val := 3
	fi.value = val
	fi.SetValue(2 * val)
	if fi.GetValue() != 2*val {
		t.Errorf("Invalid form item value. Instead of '%d', it is '%d'.", 2*val, fi.GetValue())
	}
}
func TestFormItemIntGetSurface(t *testing.T) {
	fi := testFormItemInt(t)
	if fi.GetSurface() != fi.meshes[0] {
		t.Error("Invalid surface mesh")
	}
}
func TestFormItemIntGetTarget(t *testing.T) {
	fi := testFormItemInt(t)
	if fi.GetTarget() != fi.meshes[1] {
		t.Error("Invalid target mesh")
	}
}
func TestFormItemIntAddCursor(t *testing.T) {
	fi := testFormItemInt(t)
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemIntDeleteCursor(t *testing.T) {
	fi := testFormItemInt(t)
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
	fi.DeleteCursor()
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemIntCharCallback(t *testing.T) {
	fi := testFormItemInt(t)
	fi.AddCursor()
	fi.CharCallback('b', 0.1)
	if fi.value != 0 {
		t.Errorf("Invalid value. Instead of '0', we have '%d'.", fi.value)
	}
	fi.CharCallback('0', 0.1)
	if fi.value != 0 {
		t.Errorf("Invalid value. Instead of '0', we have '%d'.", fi.value)
	}
	fi.CharCallback('-', 0.1)
	if fi.value != 0 {
		t.Errorf("Invalid value. Instead of '0', we have '%d'.", fi.value)
	}
	if fi.isNegative != true {
		t.Errorf("Invalid negative flag")
	}
	fi.CharCallback('1', 0.1)
	if fi.value != -1 {
		t.Errorf("Invalid value. Instead of '-1', we have '%d'.", fi.value)
	}
	fi = testFormItemInt(t)
	fi.AddCursor()
	fi.CharCallback('1', 0.1)
	if fi.value != 1 {
		t.Errorf("Invalid value. Instead of '1', we have '%d'.", fi.value)
	}
	fi.value = math.MaxInt32 - 10
	fi.CharCallback('1', 0.1)
	if fi.value != math.MaxInt32-10 {
		t.Errorf("Invalid value. Instead of '%d', we have '%d'.", math.MaxInt32-10, fi.value)
	}
}
func TestFormItemIntValueToString(t *testing.T) {
	fi := testFormItemInt(t)
	val := fi.ValueToString()
	if val != "0" {
		t.Errorf("Invalid value. Instead of '0', we have '%s'.", val)
	}
	fi.isNegative = true
	val = fi.ValueToString()
	if val != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", val)
	}
	fi.value = -3
	val = fi.ValueToString()
	if val != "-3" {
		t.Errorf("Invalid value. Instead of '-3', we have '%s'.", val)
	}
	fi.isNegative = false
	fi.value = 3
	val = fi.ValueToString()
	if val != "3" {
		t.Errorf("Invalid value. Instead of '3', we have '%s'.", val)
	}
}
func TestFormItemIntDeleteLastCharacter(t *testing.T) {
	fi := testFormItemInt(t)
	fi.value = 12345
	fi.charOffsets = []float32{0.1, 0.1, 0.1, 0.1, 0.1}
	fi.cursorOffsetX = float32(0.5)
	fi.DeleteLastCharacter()
	if fi.value != 1234 {
		t.Errorf("Invalid value. Instead of '1234', we have '%d'.", fi.value)
	}
	fi.DeleteLastCharacter()
	if fi.value != 123 {
		t.Errorf("Invalid value. Instead of '123', we have '%d'.", fi.value)
	}
	fi.isNegative = true
	fi.value = -2
	fi.DeleteLastCharacter()
	if fi.value != 0 {
		t.Errorf("Invalid value. Instead of '0', we have '%d'.", fi.value)
	}
	if fi.isNegative != true {
		t.Error("Invalid isNegative flag")
	}
	fi.DeleteLastCharacter()
	if fi.isNegative != false {
		t.Error("Invalid isNegative flag")
	}
	fi.DeleteLastCharacter()
	if fi.isNegative != false {
		t.Error("Invalid isNegative flag")
	}
}
func testFormItemFloat(t *testing.T) *FormItemFloat {
	mat := material.Chrome
	pos := mgl32.Vec3{0, 0, 0}
	fi := NewFormItemFloat(DefaultFormItemLabel, mat, pos, wrapperMock)

	if fi.label != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.label)
	}
	return fi
}
func TestNewFormItemFloat(t *testing.T) {
	_ = testFormItemFloat(t)
}
func TestFormItemFloatGetLabel(t *testing.T) {
	fi := testFormItemFloat(t)
	if fi.GetLabel() != DefaultFormItemLabel {
		t.Errorf("Invalid form item label. Instead of '%s', we have '%s'.", DefaultFormItemLabel, fi.GetLabel())
	}

}
func TestFormItemFloatGetValue(t *testing.T) {
	fi := testFormItemFloat(t)
	valI := 3
	fi.value = "3"
	if fi.GetValue() != float32(valI) {
		t.Errorf("Invalid form item value. Instead of '%d', it is '%f'.", valI, fi.GetValue())
	}
	fi.value = "3.2"
	if fi.GetValue() != 3.2 {
		t.Errorf("Invalid form item value. Instead of '3.2', it is '%f'.", fi.GetValue())
	}
}
func TestFormItemFloatSetValue(t *testing.T) {
	fi := testFormItemFloat(t)
	fi.SetValue(3.3)
	if fi.GetValue() != 3.3 {
		t.Errorf("Invalid form item value. Instead of '%f', it is '%f'.", 3.3, fi.GetValue())
	}
	fi.SetValue(float32(3.000))
	if fi.GetValue() != 3.0 {
		t.Errorf("Invalid form item value. Instead of '%f', it is '%f'.", 3.0, fi.GetValue())
	}
	fi.SetValue(float32(1234567890.000))
	if fi.GetValue() != 3.0 {
		t.Errorf("Invalid form item value. Instead of '%f', it is '%f'.", 3.0, fi.GetValue())
	}
	val := 12.3456700001
	fi.SetValue(float32(val))
	if fi.GetValue() != 12.345670 {
		t.Errorf("Invalid form item value. Instead of '%f', it is '%f'.", 12.345670, fi.GetValue())
	}
	fi.SetValue(float32(10000000.5752))
	if fi.GetValue() != 12.345670 {
		t.Errorf("Invalid form item value. Instead of '%f', it is '%f'.", 12.345670, fi.GetValue())
	}
}
func TestFormItemFloatGetSurface(t *testing.T) {
	fi := testFormItemFloat(t)
	if fi.GetSurface() != fi.meshes[0] {
		t.Error("Invalid surface mesh")
	}
}
func TestFormItemFloatGetTarget(t *testing.T) {
	fi := testFormItemFloat(t)
	if fi.GetTarget() != fi.meshes[1] {
		t.Error("Invalid target mesh")
	}
}
func TestFormItemFloatAddCursor(t *testing.T) {
	fi := testFormItemFloat(t)
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemFloatDeleteCursor(t *testing.T) {
	fi := testFormItemFloat(t)
	fi.AddCursor()
	if len(fi.meshes) != 3 {
		t.Error("Invalid number of target mesh")
	}
	fi.DeleteCursor()
	if len(fi.meshes) != 2 {
		t.Error("Invalid number of target mesh")
	}
}
func TestFormItemFloatCharCallback(t *testing.T) {
	fi := testFormItemFloat(t)
	fi.AddCursor()
	// start with 0
	fi.CharCallback('0', 0.1)
	if fi.value != "0" {
		t.Errorf("Invalid value. Instead of '0', we have '%s'.", fi.value)
	}
	if fi.typeState != "P0" {
		t.Errorf("Invalid typeState. Instead of 'P0', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "0" {
		t.Errorf("Invalid value. Instead of '0', we have '%s'.", fi.value)
	}
	if fi.typeState != "P0" {
		t.Errorf("Invalid typeState. Instead of 'P0', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "0." {
		t.Errorf("Invalid value. Instead of '0.', we have '%s'.", fi.value)
	}
	if fi.typeState != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('2', 0.1)
	if fi.value != "0.2" {
		t.Errorf("Invalid value. Instead of '0.2', we have '%s'.", fi.value)
	}
	if fi.typeState != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('8', 0.1)
	if fi.value != "0.28" {
		t.Errorf("Invalid value. Instead of '0.28', we have '%s'.", fi.value)
	}
	if fi.typeState != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeState)
	}
	// start with .
	fi = testFormItemFloat(t)
	fi.AddCursor()
	fi.CharCallback('.', 0.1)
	if fi.value != "" {
		t.Errorf("Invalid value. Instead of '', we have '%s'.", fi.value)
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "1" {
		t.Errorf("Invalid value. Instead of '1', we have '%s'.", fi.value)
	}
	if fi.typeState != "PI" {
		t.Errorf("Invalid typeState. Instead of 'PI', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('2', 0.1)
	if fi.value != "12" {
		t.Errorf("Invalid value. Instead of '12', we have '%s'.", fi.value)
	}
	if fi.typeState != "PI" {
		t.Errorf("Invalid typeState. Instead of 'PI', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "12." {
		t.Errorf("Invalid value. Instead of '12.', we have '%s'.", fi.value)
	}
	if fi.typeState != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('8', 0.1)
	if fi.value != "12.8" {
		t.Errorf("Invalid value. Instead of '12.8', we have '%s'.", fi.value)
	}
	if fi.typeState != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeState)
	}
	// start with -
	fi = testFormItemFloat(t)
	fi.AddCursor()
	fi.CharCallback('-', 0.1)
	if fi.value != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.value)
	}
	if fi.typeState != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.value)
	}
	if fi.typeState != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('1', 0.1)
	if fi.value != "-1" {
		t.Errorf("Invalid value. Instead of '-1', we have '%s'.", fi.value)
	}
	if fi.typeState != "NI" {
		t.Errorf("Invalid typeState. Instead of 'NI', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('2', 0.1)
	if fi.value != "-12" {
		t.Errorf("Invalid value. Instead of '-12', we have '%s'.", fi.value)
	}
	if fi.typeState != "NI" {
		t.Errorf("Invalid typeState. Instead of 'NI', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "-12." {
		t.Errorf("Invalid value. Instead of '-12.', we have '%s'.", fi.value)
	}
	if fi.typeState != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('8', 0.1)
	if fi.value != "-12.8" {
		t.Errorf("Invalid value. Instead of '-12.8', we have '%s'.", fi.value)
	}
	if fi.typeState != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeState)
	}
	// start with -0.
	fi = testFormItemFloat(t)
	fi.AddCursor()
	fi.CharCallback('-', 0.1)
	if fi.value != "-" {
		t.Errorf("Invalid value. Instead of '-', we have '%s'.", fi.value)
	}
	if fi.typeState != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('0', 0.1)
	if fi.value != "-0" {
		t.Errorf("Invalid value. Instead of '-0', we have '%s'.", fi.value)
	}
	if fi.typeState != "N0" {
		t.Errorf("Invalid typeState. Instead of 'N0', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('.', 0.1)
	if fi.value != "-0." {
		t.Errorf("Invalid value. Instead of '-0.', we have '%s'.", fi.value)
	}
	if fi.typeState != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('8', 0.1)
	if fi.value != "-0.8" {
		t.Errorf("Invalid value. Instead of '-0.8', we have '%s'.", fi.value)
	}
	if fi.typeState != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeState)
	}
}
func TestFormItemFloatValueToString(t *testing.T) {
	fi := testFormItemFloat(t)
	fi.AddCursor()
	fi.value = "-"
	if fi.ValueToString() != "-" {
		t.Errorf("Invalid valuestring. instead of '-', we have '%s'.", fi.ValueToString())
	}
	fi.value = "-3"
	if fi.ValueToString() != "-3" {
		t.Errorf("Invalid valuestring. instead of '-3', we have '%s'.", fi.ValueToString())
	}
	fi.value = "-3."
	if fi.ValueToString() != "-3." {
		t.Errorf("Invalid valuestring. instead of '-3.', we have '%s'.", fi.ValueToString())
	}
	fi.value = "-3.3"
	if fi.ValueToString() != "-3.3" {
		t.Errorf("Invalid valuestring. instead of '-3.3', we have '%s'.", fi.ValueToString())
	}

}
func TestFormItemFloatDeleteLastCharacter(t *testing.T) {
	fi := testFormItemFloat(t)
	fi.AddCursor()
	fi.CharCallback('-', 0.1)
	fi.CharCallback('3', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('3', 0.1)
	if fi.ValueToString() != "-3.3" {
		t.Errorf("Invalid valuestring. instead of '-3.3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-3." {
		t.Errorf("Invalid valuestring. instead of '-3.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-3" {
		t.Errorf("Invalid valuestring. instead of '-3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "NI" {
		t.Errorf("Invalid typeState. Instead of 'NI', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-" {
		t.Errorf("Invalid valuestring. instead of '-', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('-', 0.1)
	fi.CharCallback('0', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('2', 0.1)
	if fi.ValueToString() != "-0.2" {
		t.Errorf("Invalid valuestring. instead of '-0.2', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "NF" {
		t.Errorf("Invalid typeState. Instead of 'NF', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-0." {
		t.Errorf("Invalid valuestring. instead of '-0.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "N." {
		t.Errorf("Invalid typeState. Instead of 'N.', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-0" {
		t.Errorf("Invalid valuestring. instead of '-0', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "N0" {
		t.Errorf("Invalid typeState. Instead of 'N0', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "-" {
		t.Errorf("Invalid valuestring. instead of '-', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "N" {
		t.Errorf("Invalid typeState. Instead of 'N', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('3', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('3', 0.1)
	if fi.ValueToString() != "3.3" {
		t.Errorf("Invalid valuestring. instead of '3.3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "3." {
		t.Errorf("Invalid valuestring. instead of '3.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "3" {
		t.Errorf("Invalid valuestring. instead of '3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "PI" {
		t.Errorf("Invalid typeState. Instead of 'PI', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
	fi.CharCallback('0', 0.1)
	fi.CharCallback('.', 0.1)
	fi.CharCallback('3', 0.1)
	if fi.ValueToString() != "0.3" {
		t.Errorf("Invalid valuestring. instead of '0.3', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "PF" {
		t.Errorf("Invalid typeState. Instead of 'PF', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "0." {
		t.Errorf("Invalid valuestring. instead of '0.', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P." {
		t.Errorf("Invalid typeState. Instead of 'P.', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "0" {
		t.Errorf("Invalid valuestring. instead of '0', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P0" {
		t.Errorf("Invalid typeState. Instead of 'P0', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
	fi.DeleteLastCharacter()
	if fi.ValueToString() != "" {
		t.Errorf("Invalid valuestring. instead of '', we have '%s'.", fi.ValueToString())
	}
	if fi.typeState != "P" {
		t.Errorf("Invalid typeState. Instead of 'P', we have '%s'.", fi.typeState)
	}
}
