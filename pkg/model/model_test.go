package model

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	InvalidFilename = "not-existing-file.obj"
	ValidFilename   = "testdata/test_cube.obj"
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
}
