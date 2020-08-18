package model

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"
)

func TestNewStreetLampBuilder(t *testing.T) {
	builder := NewStreetLampBuilder()
	defaultPosition := mgl32.Vec3{0.0, 0.0, 0.0}
	if builder.position != defaultPosition {
		t.Error("Invalid position")
	}
	if builder.wrapper != nil {
		t.Error("Wrapper should be nil")
	}
	if builder.rotationX != 0.0 {
		t.Errorf("Invalid rotationX. Instead of '0.0', it is '%f'.", builder.rotationX)
	}
	if builder.rotationY != 0.0 {
		t.Errorf("Invalid rotationY. Instead of '0.0', it is '%f'.", builder.rotationY)
	}
	if builder.rotationZ != 0.0 {
		t.Errorf("Invalid rotationZ. Instead of '0.0', it is '%f'.", builder.rotationZ)
	}
	if builder.constantTerm != 0.0 {
		t.Errorf("Invalid constantTerm. Instead of '0.0', it is '%f'.", builder.constantTerm)
	}
	if builder.linearTerm != 0.0 {
		t.Errorf("Invalid linearTerm. Instead of '0.0', it is '%f'.", builder.linearTerm)
	}
	if builder.quadraticTerm != 0.0 {
		t.Errorf("Invalid quadraticTerm. Instead of '0.0', it is '%f'.", builder.quadraticTerm)
	}
	if builder.cutoff != 0.0 {
		t.Errorf("Invalid cutoff. Instead of '0.0', it is '%f'.", builder.cutoff)
	}
	if builder.outerCutoff != 0.0 {
		t.Errorf("Invalid outerCutoff. Instead of '0.0', it is '%f'.", builder.outerCutoff)
	}
	if !builder.lampOn {
		t.Error("Invalid lampOn flag.")
	}
}
func TestStreetLampBuilderSetPosition(t *testing.T) {
	builder := NewStreetLampBuilder()
	pos := mgl32.Vec3{0.0, 1.0, 0.0}
	builder.SetPosition(pos)
	if builder.position != pos {
		t.Error("Invalid position")
	}
}
func TestStreetLampBuilderSetRotation(t *testing.T) {
	builder := NewStreetLampBuilder()
	x := float32(10)
	y := float32(20)
	z := float32(30)
	builder.SetRotation(x, y, z)
	if builder.rotationX != 10.0 {
		t.Errorf("Invalid rotationX. Instead of '10.0', it is '%f'.", builder.rotationX)
	}
	if builder.rotationY != 20.0 {
		t.Errorf("Invalid rotationY. Instead of '20.0', it is '%f'.", builder.rotationY)
	}
	if builder.rotationZ != 30.0 {
		t.Errorf("Invalid rotationZ. Instead of '30.0', it is '%f'.", builder.rotationZ)
	}
}
func TestStreetLampBuilderSetPoleLength(t *testing.T) {
	builder := NewStreetLampBuilder()
	len := float32(10)
	builder.SetPoleLength(len)
	if builder.poleLength != len {
		t.Errorf("Invalid poleLength. Instead of '%f', it is '%f'.", len, builder.poleLength)
	}
}
func TestStreetLampBuilderSetBulbMaterial(t *testing.T) {
	builder := NewStreetLampBuilder()
	mat := material.Jade
	builder.SetBulbMaterial(mat)
	if builder.bulbMaterial != mat {
		t.Errorf("Invalid bulbMaterial")
	}
}
func TestStreetLampBuilderSetLightTerms(t *testing.T) {
	builder := NewStreetLampBuilder()
	con := float32(1)
	lin := float32(0.5)
	qua := float32(0.09)
	builder.SetLightTerms(con, lin, qua)
	if builder.constantTerm != con {
		t.Errorf("Invalid constantTerm. Instead of '%f', it is '%f'.", con, builder.constantTerm)
	}
	if builder.linearTerm != lin {
		t.Errorf("Invalid linearTerm. Instead of '%f', it is '%f'.", lin, builder.rotationY)
	}
	if builder.quadraticTerm != qua {
		t.Errorf("Invalid quadraticTerm. Instead of '%f', it is '%f'.", qua, builder.rotationZ)
	}
}
func TestStreetLampBuilderSetCutoff(t *testing.T) {
	builder := NewStreetLampBuilder()
	c := float32(1)
	oc := float32(0.5)
	builder.SetCutoff(c, oc)
	if builder.cutoff != c {
		t.Errorf("Invalid cutoff. Instead of '%f', it is '%f'.", c, builder.cutoff)
	}
	if builder.outerCutoff != oc {
		t.Errorf("Invalid outerCutoff. Instead of '%f', it is '%f'.", oc, builder.outerCutoff)
	}
}
func TestMaterialStreetLamp(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	bulbPosition := mgl32.Vec3{0.34999996, 0, -0.20000002}
	polePosition := mgl32.Vec3{0.0, 3.0, 0.0}
	topPosition := mgl32.Vec3{0.54999995, 0, 3.2}
	scale := float32(6.0)
	builder := NewStreetLampBuilder()
	builder.SetPosition(position)
	builder.SetPoleLength(scale)
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Without wrapper, we should have panic.")
			}
		}()
		builder.BuildMaterial()
	}()
	builder.SetWrapper(wrapperMock)

	lamp := builder.BuildMaterial()

	if lamp.meshes[0].GetPosition() != polePosition {
		t.Errorf("Invalid pole position. Instead of '%v', we have '%v'.", polePosition, lamp.meshes[0].GetPosition())
	}
	if !lamp.meshes[1].GetPosition().ApproxEqualThreshold(topPosition, 0.0001) {
		t.Errorf("Invalid top position. Instead of '%v', we have '%v'.", topPosition, lamp.meshes[1].GetPosition())
	}
	if !lamp.meshes[2].GetPosition().ApproxEqualThreshold(bulbPosition, 0.0001) {
		t.Errorf("Invalid bulb position. Instead of '%v', we have '%v'.", bulbPosition, lamp.meshes[2].GetPosition())
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
func TestMaterialStreetLampWithoutLight(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	bulbPosition := mgl32.Vec3{0.34999996, 0, -0.20000002}
	polePosition := mgl32.Vec3{0.0, 3.0, 0.0}
	topPosition := mgl32.Vec3{0.54999995, 0, 3.2}
	scale := float32(6.0)
	builder := NewStreetLampBuilder()
	builder.SetPosition(position)
	builder.SetPoleLength(scale)
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Without wrapper, we should have panic.")
			}
		}()
		builder.BuildMaterial()
	}()
	builder.SetWrapper(wrapperMock)
	builder.SetLampOn(false)

	lamp := builder.BuildMaterial()
	lamp.TurnLampOn()

	if lamp.meshes[0].GetPosition() != polePosition {
		t.Errorf("Invalid pole position. Instead of '%v', we have '%v'.", polePosition, lamp.meshes[0].GetPosition())
	}
	if !lamp.meshes[1].GetPosition().ApproxEqualThreshold(topPosition, 0.0001) {
		t.Errorf("Invalid top position. Instead of '%v', we have '%v'.", topPosition, lamp.meshes[1].GetPosition())
	}
	if !lamp.meshes[2].GetPosition().ApproxEqualThreshold(bulbPosition, 0.0001) {
		t.Errorf("Invalid bulb position. Instead of '%v', we have '%v'.", bulbPosition, lamp.meshes[2].GetPosition())
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
	bulbPosition := mgl32.Vec3{0, 0.34999993, 0}
	polePosition := mgl32.Vec3{0.0, 3.0, 0.0}
	topPosition := mgl32.Vec3{0.54999995, 0, 3}
	scale := float32(6.0)

	builder := NewStreetLampBuilder()
	builder.SetPosition(position)
	builder.SetPoleLength(scale)
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Without wrapper, we should have panic.")
			}
		}()
		builder.BuildTexture()
	}()
	builder.SetWrapper(wrapperMock)

	lamp := builder.BuildTexture()

	if lamp.meshes[0].GetPosition() != polePosition {
		t.Errorf("Invalid pole position. Instead of '%v', we have '%v'.", polePosition, lamp.meshes[0].GetPosition())
	}
	if !lamp.meshes[1].GetPosition().ApproxEqualThreshold(topPosition, 0.0001) {
		t.Errorf("Invalid top position. Instead of '%v', we have '%v'.", topPosition, lamp.meshes[1].GetPosition())
	}
	if !lamp.meshes[2].GetPosition().ApproxEqualThreshold(bulbPosition, 0.001) {
		t.Errorf("Invalid bulb position. Instead of '%v', we have '%v'.", bulbPosition, lamp.meshes[2].GetPosition())
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
func TestTexturedStreetLampWithoutLight(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	bulbPosition := mgl32.Vec3{0, 0.34999993, 0}
	polePosition := mgl32.Vec3{0.0, 3.0, 0.0}
	topPosition := mgl32.Vec3{0.54999995, 0, 3}
	scale := float32(6.0)

	builder := NewStreetLampBuilder()
	builder.SetPosition(position)
	builder.SetPoleLength(scale)
	builder.SetLampOn(false)
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Without wrapper, we should have panic.")
			}
		}()
		builder.BuildTexture()
	}()
	builder.SetWrapper(wrapperMock)

	lamp := builder.BuildTexture()
	lamp.TurnLampOn()

	if lamp.meshes[0].GetPosition() != polePosition {
		t.Errorf("Invalid pole position. Instead of '%v', we have '%v'.", polePosition, lamp.meshes[0].GetPosition())
	}
	if !lamp.meshes[1].GetPosition().ApproxEqualThreshold(topPosition, 0.0001) {
		t.Errorf("Invalid top position. Instead of '%v', we have '%v'.", topPosition, lamp.meshes[1].GetPosition())
	}
	if !lamp.meshes[2].GetPosition().ApproxEqualThreshold(bulbPosition, 0.001) {
		t.Errorf("Invalid bulb position. Instead of '%v', we have '%v'.", bulbPosition, lamp.meshes[2].GetPosition())
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
