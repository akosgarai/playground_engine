package mesh

import (
	"reflect"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/primitives/vertex"
	"github.com/akosgarai/playground_engine/pkg/testhelper"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	wrapperMock testhelper.GLWrapperMock
	shaderMock  testhelper.ShaderMock

	DefaultPosition  = mgl32.Vec3{0.0, 0.0, 0.0}
	DefaultDirection = mgl32.Vec3{0.0, 0.0, 0.0}
	DefaultScale     = mgl32.Vec3{1.0, 1.0, 1.0}
	DefaultYaw       = float32(0.0)
	DefaultPitch     = float32(0.0)
	DefaultRoll      = float32(0.0)
	DefaultVelocity  = float32(0.0)
)

func TestSetScale(t *testing.T) {
	var m Mesh
	scale := mgl32.Vec3{2, 2, 2}
	m.SetScale(scale)
	if m.scale != scale {
		t.Error("Scale mismatch")
	}
}
func TestSetPosition(t *testing.T) {
	var m Mesh
	pos := mgl32.Vec3{0, 1, 2}
	m.SetPosition(pos)
	if m.position != pos {
		t.Error("Position mismatch")
	}
}
func TestSetDirection(t *testing.T) {
	var m Mesh
	dir := mgl32.Vec3{0, 1, 0}
	m.SetDirection(dir)
	if m.direction != dir {
		t.Error("Direction mismatch")
	}
}
func TestSetSpeed(t *testing.T) {
	var m Mesh
	m.SetSpeed(10)
	if m.velocity != 10 {
		t.Error("Speed mismatch")
	}
}
func TestGetPosition(t *testing.T) {
	var m Mesh
	pos := mgl32.Vec3{0, 1, 2}
	m.SetPosition(pos)
	if m.GetPosition() != pos {
		t.Error("Position mismatch")
	}
}
func TestGetDirection(t *testing.T) {
	var m Mesh
	dir := mgl32.Vec3{0, 1, 0}
	m.SetDirection(dir)
	if m.GetDirection() != dir {
		t.Error("Direction mismatch")
	}
}
func TestUpdate(t *testing.T) {
	var m Mesh
	m.SetDirection(mgl32.Vec3{0, 0, 0})
	pos := mgl32.Vec3{0, 1, 2}
	m.SetPosition(pos)
	m.Update(2)
	if m.GetPosition() != pos {
		t.Error("Invalid position after update")
	}
	dir := mgl32.Vec3{0, 1, 0}
	m.SetDirection(dir)
	m.SetSpeed(10)
	m.Update(2)
	expectedPosition := mgl32.Vec3{0, 21, 2}
	if m.GetPosition() != expectedPosition {
		t.Error("Invalid position after update")
	}
}
func TestModelTransformation(t *testing.T) {
	var m Mesh
	pos := mgl32.Vec3{0, 0, 0}
	m.SetPosition(pos)
	scale := mgl32.Vec3{1, 1, 1}
	m.SetScale(scale)
	M := m.ModelTransformation()
	if M != mgl32.Ident4() {
		t.Error("Invalid model matrix")
	}
}
func TestSetParent(t *testing.T) {
	var m Mesh
	var parent interfaces.Mesh
	m.SetParent(parent)
	if m.parentSet != true {
		t.Error("After setting the parent, the flag supposed to be true")
	}
	if m.parent != parent {
		t.Error("The parent supposed to be the same")
	}
}
func TestIsParentMesh(t *testing.T) {
	var m Mesh
	var parent interfaces.Mesh
	if m.IsParentMesh() != true {
		t.Error("Before setting the parent, it should return true")
	}
	m.SetParent(parent)
	if m.IsParentMesh() != false {
		t.Error("After setting the parent, it should return false")
	}
}
func TestTransformationGettersWithParent(t *testing.T) {
	var m Mesh
	parent := NewPointMesh(wrapperMock)
	parent.position = mgl32.Vec3{1.0, 0.0, 0.0}
	m.SetParent(parent)
	modelTr := m.ModelTransformation()
	if modelTr == mgl32.Ident4() {
		t.Error("Model tr shouldn't be ident, if the parent transformation is set.")
	}
}
func TestRotateDirection(t *testing.T) {
	m := NewPointMesh(wrapperMock)
	rotationAngle := float32(90.0)
	rotationAxis := mgl32.Vec3{0.0, 1.0, 0.0}
	m.rotateDirection(rotationAngle, rotationAxis)
	nullVec := mgl32.Vec3{0.0, 0.0, 0.0}
	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	leftDir := mgl32.Vec3{-1.0, 0.0, 0.0}
	frontDir := mgl32.Vec3{0.0, 0.0, 1.0}
	if m.direction != nullVec {
		t.Error("Rotating 0 vec should lead to 0 vec.")
	}
	m.direction = upDir
	m.rotateDirection(rotationAngle, rotationAxis)
	if !m.direction.ApproxEqualThreshold(upDir, 0.0001) {
		t.Log(m.direction)
		t.Error("Rotating the same dir as axis shouldn't change the direction.")
	}
	m.direction = leftDir
	m.rotateDirection(rotationAngle, rotationAxis)
	if !m.direction.ApproxEqualThreshold(frontDir, 0.001) {
		t.Log(m.direction)
		t.Log(frontDir)
		t.Error("Rotating different dir and axis should change the direction.")
	}
}
func TestRotatePosition(t *testing.T) {
	m := NewPointMesh(wrapperMock)
	rotationAngle := float32(90.0)
	rotationAxis := mgl32.Vec3{0.0, 1.0, 0.0}
	m.RotatePosition(rotationAngle, rotationAxis)
	nullVec := mgl32.Vec3{0.0, 0.0, 0.0}
	upPos := mgl32.Vec3{0.0, 1.0, 0.0}
	leftPos := mgl32.Vec3{-1.0, 0.0, 0.0}
	frontPos := mgl32.Vec3{0.0, 0.0, 1.0}
	if m.position != nullVec {
		t.Error("Rotating 0 vec should lead to 0 vec.")
	}
	m.position = upPos
	m.RotatePosition(rotationAngle, rotationAxis)
	if !m.position.ApproxEqualThreshold(upPos, 0.0001) {
		t.Log(m.position)
		t.Error("Rotating the same pos as axis shouldn't change the position.")
	}
	m.position = leftPos
	m.RotatePosition(rotationAngle, rotationAxis)
	if !m.position.ApproxEqualThreshold(frontPos, 0.001) {
		t.Log(m.position)
		t.Log(frontPos)
		t.Error("Rotating different pos and axis should change the position.")
	}
}
func TestRotateY(t *testing.T) {
	m := NewPointMesh(wrapperMock)
	rotationAngle := float32(90.0)
	m.RotateY(rotationAngle)
	nullVec := mgl32.Vec3{0.0, 0.0, 0.0}
	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	leftDir := mgl32.Vec3{-1.0, 0.0, 0.0}
	frontDir := mgl32.Vec3{0.0, 0.0, 1.0}
	if m.yaw != rotationAngle {
		t.Error("RotateY should update the yaw")
	}
	if m.direction != nullVec {
		t.Error("Rotating 0 vec should lead to 0 vec.")
	}
	m.direction = upDir
	m.RotateY(rotationAngle)
	if !m.direction.ApproxEqualThreshold(upDir, 0.0001) {
		t.Log(m.direction)
		t.Error("Rotating the same dir as axis shouldn't change the direction.")
	}
	if m.yaw != rotationAngle*2 {
		t.Error("RotateY should update the yaw")
	}
	m.direction = leftDir
	m.RotateY(rotationAngle)
	if !m.direction.ApproxEqualThreshold(frontDir, 0.001) {
		t.Log(m.direction)
		t.Log(frontDir)
		t.Error("Rotating different dir and axis should change the direction.")
	}
	if m.yaw != rotationAngle*3 {
		t.Error("RotateY should update the yaw")
	}
}
func TestRotateX(t *testing.T) {
	m := NewPointMesh(wrapperMock)
	rotationAngle := float32(90.0)
	m.RotateX(rotationAngle)
	nullVec := mgl32.Vec3{0.0, 0.0, 0.0}
	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	leftDir := mgl32.Vec3{-1.0, 0.0, 0.0}
	frontDir := mgl32.Vec3{0.0, 0.0, 1.0}
	if m.pitch != rotationAngle {
		t.Error("RotateX should update the pitch")
	}
	if m.direction != nullVec {
		t.Error("Rotating 0 vec should lead to 0 vec.")
	}
	m.direction = upDir
	m.RotateX(rotationAngle)
	if !m.direction.ApproxEqualThreshold(frontDir, 0.001) {
		t.Log(m.direction)
		t.Error("Rotating different dir and axis should change the direction.")
	}
	if m.pitch != rotationAngle*2 {
		t.Error("RotateX should update the pitch")
	}
	m.direction = leftDir
	m.RotateX(rotationAngle)
	if !m.direction.ApproxEqualThreshold(leftDir, 0.001) {
		t.Log(m.direction)
		t.Log(frontDir)
		t.Error("Rotating the same dir as axis shouldn't change the direction.")
	}
	if m.pitch != rotationAngle*3 {
		t.Error("RotateX should update the pitch")
	}
}
func TestRotateZ(t *testing.T) {
	m := NewPointMesh(wrapperMock)
	rotationAngle := float32(90.0)
	m.RotateZ(rotationAngle)
	nullVec := mgl32.Vec3{0.0, 0.0, 0.0}
	upDir := mgl32.Vec3{0.0, 1.0, 0.0}
	leftDir := mgl32.Vec3{-1.0, 0.0, 0.0}
	frontDir := mgl32.Vec3{0.0, 0.0, 1.0}
	if m.roll != rotationAngle {
		t.Error("RotateZ should update the roll")
	}
	if m.direction != nullVec {
		t.Error("Rotating 0 vec should lead to 0 vec.")
	}
	m.direction = upDir
	m.RotateZ(rotationAngle)
	if !m.direction.ApproxEqualThreshold(leftDir, 0.001) {
		t.Log(m.direction)
		t.Error("Rotating different dir and axis should change the direction.")
	}
	if m.roll != rotationAngle*2 {
		t.Error("RotateZ should update the roll")
	}
	m.direction = frontDir
	m.RotateZ(rotationAngle)
	if !m.direction.ApproxEqualThreshold(frontDir, 0.001) {
		t.Log(m.direction)
		t.Log(frontDir)
		t.Error("Rotating the same dir as axis shouldn't change the direction.")
	}
	if m.roll != rotationAngle*3 {
		t.Error("RotateZ should update the roll")
	}
}
func TestIsBoundingObjectParamsSet(t *testing.T) {
	var m Mesh
	boParams := make(map[string]float32)
	if m.IsBoundingObjectSet() != false {
		t.Error("Before setting the bo, it should return false")
	}
	m.SetBoundingObject(boundingobject.New("AABB", boParams))
	if m.IsBoundingObjectSet() != true {
		t.Error("After setting the bo, it should return true")
	}
}
func TestGetBoundingObjectSphere(t *testing.T) {
	v := []vertex.Vertex{}
	i := []uint32{}
	params := make(map[string]float32)
	params["radius"] = float32(0.1)
	bo := boundingobject.New("Sphere", params)
	mesh := NewMaterialMesh(v, i, material.Jade, wrapperMock)
	mesh.SetBoundingObject(bo)

	newBo := mesh.GetBoundingObject()
	if !reflect.DeepEqual(bo, newBo) {
		t.Errorf("Invalid bounding object. Instead of '%v', we have '%v'.", bo, newBo)
	}
	scaledParams := make(map[string]float32)
	scaledParams["radius"] = float32(1.0)
	scaledBo := boundingobject.New("Sphere", scaledParams)
	mesh.SetScale(mgl32.Vec3{10.0, 10.0, 10.0})
	newBo = mesh.GetBoundingObject()
	if !reflect.DeepEqual(scaledBo, newBo) {
		t.Errorf("Invalid bounding object. Instead of '%v', we have '%v'.", scaledBo, newBo)
	}
}
func TestGetBoundingObjectAABB(t *testing.T) {
	v := []vertex.Vertex{}
	i := []uint32{}
	params := make(map[string]float32)
	params["height"] = float32(1.0)
	params["width"] = float32(1.0)
	params["length"] = float32(1.0)
	bo := boundingobject.New("AABB", params)
	mesh := NewMaterialMesh(v, i, material.Jade, wrapperMock)
	mesh.SetBoundingObject(bo)

	newBo := mesh.GetBoundingObject()
	if !reflect.DeepEqual(bo, newBo) {
		t.Errorf("Invalid bounding object. Instead of '%v', we have '%v'.", bo, newBo)
	}
	scaledParams := make(map[string]float32)
	scaledParams["height"] = float32(10.0)
	scaledParams["width"] = float32(10.0)
	scaledParams["length"] = float32(10.0)
	scaledBo := boundingobject.New("AABB", scaledParams)
	mesh.SetScale(mgl32.Vec3{10.0, 10.0, 10.0})
	newBo = mesh.GetBoundingObject()
	if !reflect.DeepEqual(scaledBo, newBo) {
		t.Errorf("Invalid bounding object. Instead of '%v', we have '%v'.", scaledBo, newBo)
	}
}
func CheckNewMesh(mesh Mesh, t *testing.T) {
	if mesh.position != DefaultPosition {
		t.Errorf("Invalid position. Instead of '%v', we have '%v'.", DefaultPosition, mesh.position)
	}
	if mesh.direction != DefaultDirection {
		t.Errorf("Invalid direction. Instead of '%v', we have '%v'.", DefaultDirection, mesh.direction)
	}
	if mesh.scale != DefaultScale {
		t.Errorf("Invalid scale. Instead of '%v', we have '%v'.", DefaultScale, mesh.scale)
	}
	if mesh.yaw != DefaultYaw {
		t.Errorf("Invalid yaw. Instead of '%f', we have '%f'.", DefaultYaw, mesh.yaw)
	}
	if mesh.pitch != DefaultPitch {
		t.Errorf("Invalid pitch. Instead of '%f', we have '%f'.", DefaultPitch, mesh.pitch)
	}
	if mesh.roll != DefaultRoll {
		t.Errorf("Invalid roll. Instead of '%f', we have '%f'.", DefaultRoll, mesh.roll)
	}
	if mesh.velocity != DefaultVelocity {
		t.Errorf("Invalid velocity. Instead of '%f', we have '%f'.", DefaultVelocity, mesh.velocity)
	}
	if mesh.parentSet != false {
		t.Error("Invalid parentSet. Should be false.")
	}
	if mesh.boundingObjectSet != false {
		t.Error("Invalid boundingObjectSet. Should be false.")
	}
}
func TestNewMaterialMesh(t *testing.T) {
	v := []vertex.Vertex{}
	i := []uint32{}
	mesh := NewMaterialMesh(v, i, material.Jade, wrapperMock)

	CheckNewMesh(mesh.Mesh, t)
}
func TestMaterialMeshDraw(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panicked.")
			}
		}()
		v := []vertex.Vertex{}
		i := []uint32{}
		mesh := NewMaterialMesh(v, i, material.Jade, wrapperMock)

		CheckNewMesh(mesh.Mesh, t)
		mesh.Draw(shaderMock)
	}()
}
func TestNewPointMesh(t *testing.T) {
	mesh := NewPointMesh(wrapperMock)

	CheckNewMesh(mesh.Mesh, t)
}
func TestPointMeshAddVertex(t *testing.T) {
	mesh := NewPointMesh(wrapperMock)

	CheckNewMesh(mesh.Mesh, t)
	var v vertex.Vertex
	if len(mesh.Vertices) != 0 {
		t.Errorf("Invalid number of vertices. Instead of '0', we have '%d'.", len(mesh.Vertices))
	}
	mesh.AddVertex(v)
	if len(mesh.Vertices) != 1 {
		t.Errorf("Invalid number of vertices. Instead of '1', we have '%d'.", len(mesh.Vertices))
	}
}
func TestPointMeshDraw(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panicked.")
			}
		}()
		mesh := NewPointMesh(wrapperMock)

		CheckNewMesh(mesh.Mesh, t)
		var v vertex.Vertex
		mesh.AddVertex(v)
		mesh.Draw(shaderMock)
	}()
}
func TestNewColorMesh(t *testing.T) {
	v := []vertex.Vertex{}
	i := []uint32{}
	col := []mgl32.Vec3{mgl32.Vec3{0.0, 0.0, 0.0}}
	mesh := NewColorMesh(v, i, col, wrapperMock)

	CheckNewMesh(mesh.Mesh, t)
}
func TestColorMeshDraw(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panicked.")
			}
		}()
		v := []vertex.Vertex{}
		i := []uint32{}
		col := []mgl32.Vec3{mgl32.Vec3{0.0, 0.0, 0.0}}
		mesh := NewColorMesh(v, i, col, wrapperMock)

		CheckNewMesh(mesh.Mesh, t)
		mesh.Draw(shaderMock)
	}()
}
func TestNewTexturedMesh(t *testing.T) {
	v := []vertex.Vertex{}
	i := []uint32{}
	var tex *texture.Texture
	textures := texture.Textures{tex}
	mesh := NewTexturedMesh(v, i, textures, wrapperMock)

	CheckNewMesh(mesh.Mesh, t)
}
func TestTexturedMeshDraw(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Shouldn't have panicked. '%v'", r)
			}
		}()
		v := []vertex.Vertex{}
		i := []uint32{}
		tex := texture.Texture{1, 1, 1, "uniform", wrapperMock, "path"}
		textures := texture.Textures{&tex}
		mesh := NewTexturedMesh(v, i, textures, wrapperMock)

		CheckNewMesh(mesh.Mesh, t)
		mesh.Draw(shaderMock)
	}()
}
func TestNewTexturedColoredMesh(t *testing.T) {
	v := []vertex.Vertex{}
	i := []uint32{}
	col := []mgl32.Vec3{mgl32.Vec3{0.0, 0.0, 0.0}}
	var tex *texture.Texture
	textures := texture.Textures{tex}
	mesh := NewTexturedColoredMesh(v, i, textures, col, wrapperMock)

	CheckNewMesh(mesh.Mesh, t)
}
func TestTexturedColoredMeshDraw(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Shouldn't have panicked. '%v'", r)
			}
		}()
		v := []vertex.Vertex{}
		i := []uint32{}
		col := []mgl32.Vec3{mgl32.Vec3{0.0, 0.0, 0.0}}
		tex := texture.Texture{1, 1, 1, "uniform", wrapperMock, "path"}
		textures := texture.Textures{&tex}
		mesh := NewTexturedColoredMesh(v, i, textures, col, wrapperMock)

		CheckNewMesh(mesh.Mesh, t)
		mesh.Draw(shaderMock)
	}()
}
func TestNewTexturedMaterialMesh(t *testing.T) {
	v := []vertex.Vertex{}
	i := []uint32{}
	var tex *texture.Texture
	textures := texture.Textures{tex}
	mesh := NewTexturedMaterialMesh(v, i, textures, material.Jade, wrapperMock)

	CheckNewMesh(mesh.Mesh, t)
}
func TestTexturedMaterialMeshDraw(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Shouldn't have panicked. '%v'", r)
			}
		}()
		v := []vertex.Vertex{}
		i := []uint32{}
		tex := texture.Texture{1, 1, 1, "uniform", wrapperMock, "path"}
		textures := texture.Textures{&tex}
		mesh := NewTexturedMaterialMesh(v, i, textures, material.Jade, wrapperMock)

		CheckNewMesh(mesh.Mesh, t)
		mesh.Draw(shaderMock)
	}()
}
