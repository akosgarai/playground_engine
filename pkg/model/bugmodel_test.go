package model

import (
	"reflect"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"
)

func TestNewBugBuilder(t *testing.T) {
	defaultPosition := mgl32.Vec3{0, 0, 0}
	defaultScale := mgl32.Vec3{1, 1, 1}
	defaultLightCC := mgl32.Vec3{1, 1, 1}
	defaultBodyMaterial := material.Greenrubber
	defaultBottomMaterial := material.Emerald
	defaultDirection := mgl32.Vec3{0, 0, 0}
	defaultRotationAxis := mgl32.Vec3{0, 0, 0}
	builder := NewBugBuilder()
	if builder.position != defaultPosition {
		t.Errorf("Invalid position. Instead of '%v', it is '%v'.", defaultPosition, builder.position)
	}
	if builder.scale != defaultScale {
		t.Errorf("Invalid scale. Instead of '%v', it is '%v'.", defaultScale, builder.scale)
	}
	if builder.wrapper != nil {
		t.Error("Invalid default wrapper. It supposed to be nil.")
	}
	if builder.bodyMaterial != defaultBodyMaterial {
		t.Error("Invalid body material.")
	}
	if builder.bottomMaterial != defaultBottomMaterial {
		t.Error("Invalid bottom material.")
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
	if builder.spherePrecision != 20 {
		t.Errorf("Invalid spherePrecision. Instead of '20', it is '%d'.", builder.spherePrecision)
	}
	if builder.lightAmbient != defaultLightCC {
		t.Error("Invalid ambient light component.")
	}
	if builder.lightDiffuse != defaultLightCC {
		t.Error("Invalid diffuse light component.")
	}
	if builder.lightDiffuse != defaultLightCC {
		t.Error("Invalid diffuse light component.")
	}
	if builder.constantTerm != 1.0 {
		t.Errorf("Invalid constantTerm. Instead of '1.0', it is '%f'.", builder.constantTerm)
	}
	if builder.linearTerm != 0.14 {
		t.Errorf("Invalid linearTerm. Instead of '0.14', it is '%f'.", builder.linearTerm)
	}
	if builder.quadraticTerm != 0.07 {
		t.Errorf("Invalid quadraticTerm. Instead of '0.07', it is '%f'.", builder.quadraticTerm)
	}
	if !builder.withLight {
		t.Error("Invalid initial value of withLight")
	}
	if builder.velocity != 0.0 {
		t.Errorf("Invalid velocity. Instead of '0.0', it is '%f'.", builder.velocity)
	}
	if builder.movementRotationAngle != 0.0 {
		t.Errorf("Invalid movementRotationAngle. Instead of '0.0', it is '%f'.", builder.movementRotationAngle)
	}
	if builder.direction != defaultDirection {
		t.Errorf("Invalid direction. Instead of '%v', it is '%v'.", defaultDirection, builder.direction)
	}
	if builder.movementRotationAxis != defaultRotationAxis {
		t.Errorf("Invalid movementRotationAxis. Instead of '%v', it is '%v'.", defaultRotationAxis, builder.movementRotationAxis)
	}
	if builder.withWings {
		t.Error("Invalid initial value of withWings.")
	}
}
func TestBugBuilderSetPosition(t *testing.T) {
	newPosition := mgl32.Vec3{0, 0, 0}
	builder := NewBugBuilder()
	builder.SetPosition(newPosition)
	if builder.position != newPosition {
		t.Errorf("Invalid position. Instead of '%v', it is '%v'.", newPosition, builder.position)
	}
}
func TestBugBuilderSetScale(t *testing.T) {
	newScale := mgl32.Vec3{0, 0, 0}
	builder := NewBugBuilder()
	builder.SetScale(newScale)
	if builder.scale != newScale {
		t.Errorf("Invalid scale. Instead of '%v', it is '%v'.", newScale, builder.scale)
	}
}
func TestBugBuilderSetWrapper(t *testing.T) {
	builder := NewBugBuilder()
	builder.SetWrapper(wrapperMock)
	if !reflect.DeepEqual(builder.wrapper, wrapperMock) {
		t.Error("Invalid wrapper.")
	}
}
func TestBugBuilderSetBodyMaterial(t *testing.T) {
	builder := NewBugBuilder()
	mat := material.Jade
	builder.SetBodyMaterial(mat)
	if builder.bodyMaterial != mat {
		t.Error("Invalid bodyMaterial.")
	}
}
func TestBugBuilderSetBottomMaterial(t *testing.T) {
	builder := NewBugBuilder()
	mat := material.Jade
	builder.SetBottomMaterial(mat)
	if builder.bottomMaterial != mat {
		t.Error("Invalid bottomMaterial.")
	}
}
func TestBugBuilderSetEyeMaterial(t *testing.T) {
	builder := NewBugBuilder()
	mat := material.Jade
	builder.SetEyeMaterial(mat)
	if builder.eyeMaterial != mat {
		t.Error("Invalid eyeMaterial.")
	}
}
func TestBugBuilderSetRotation(t *testing.T) {
	builder := NewBugBuilder()
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
func TestBugBuilderSetLightAmbient(t *testing.T) {
	newAmbient := mgl32.Vec3{0, 0.5, 0}
	builder := NewBugBuilder()
	builder.SetLightAmbient(newAmbient)
	if builder.lightAmbient != newAmbient {
		t.Errorf("Invalid lightAmbient. Instead of '%v', it is '%v'.", newAmbient, builder.lightAmbient)
	}
}
func TestBugBuilderSetLightDiffuse(t *testing.T) {
	newDiffuse := mgl32.Vec3{0, 0.5, 0}
	builder := NewBugBuilder()
	builder.SetLightDiffuse(newDiffuse)
	if builder.lightDiffuse != newDiffuse {
		t.Errorf("Invalid lightDiffuse. Instead of '%v', it is '%v'.", newDiffuse, builder.lightDiffuse)
	}
}
func TestBugBuilderSetLightSpecular(t *testing.T) {
	newSpecular := mgl32.Vec3{0, 0.5, 0}
	builder := NewBugBuilder()
	builder.SetLightSpecular(newSpecular)
	if builder.lightSpecular != newSpecular {
		t.Errorf("Invalid lightSpecular. Instead of '%v', it is '%v'.", newSpecular, builder.lightSpecular)
	}
}
func TestBugBuilderSetLightTerms(t *testing.T) {
	builder := NewBugBuilder()
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
func TestBugBuilderSetWithLight(t *testing.T) {
	builder := NewBugBuilder()
	values := []bool{true, true, false, true, false, false}
	for _, v := range values {
		builder.SetWithLight(v)
		if builder.withLight != v {
			t.Errorf("Invalid withLight flag. it supposed to be '%v'.", v)
		}
	}
}
func TestBugBuilderSetDirection(t *testing.T) {
	newDirection := mgl32.Vec3{0, 1, 0}
	builder := NewBugBuilder()
	builder.SetDirection(newDirection)
	if builder.direction != newDirection {
		t.Errorf("Invalid direction. Instead of '%v', it is '%v'.", newDirection, builder.direction)
	}
}
func TestBugBuilderSetMovementRotationAxis(t *testing.T) {
	newAxis := mgl32.Vec3{0, 1, 0}
	builder := NewBugBuilder()
	builder.SetMovementRotationAxis(newAxis)
	if builder.movementRotationAxis != newAxis {
		t.Errorf("Invalid movementRotationAxis. Instead of '%v', it is '%v'.", newAxis, builder.movementRotationAxis)
	}
}
func TestBugBuilderSetMovementRotationAngle(t *testing.T) {
	newAngle := float32(30.0)
	builder := NewBugBuilder()
	builder.SetMovementRotationAngle(newAngle)
	if builder.movementRotationAngle != newAngle {
		t.Errorf("Invalid movementRotationAngle. Instead of '%f', it is '%f'.", newAngle, builder.movementRotationAngle)
	}
}
func TestBugBuilderSetVelocity(t *testing.T) {
	newVelocity := float32(30.0)
	builder := NewBugBuilder()
	builder.SetVelocity(newVelocity)
	if builder.velocity != newVelocity {
		t.Errorf("Invalid velocity. Instead of '%f', it is '%f'.", newVelocity, builder.velocity)
	}
}
func TestBugBuilderSetSameDirectionTime(t *testing.T) {
	newTime := float32(30.0)
	builder := NewBugBuilder()
	builder.SetSameDirectionTime(newTime)
	if builder.sameDirectionTime != newTime {
		t.Errorf("Invalid sameDirectionTime. Instead of '%f', it is '%f'.", newTime, builder.sameDirectionTime)
	}
}
func TestBugBuilderSetWithWings(t *testing.T) {
	builder := NewBugBuilder()
	values := []bool{true, true, false, true, false, false}
	for _, v := range values {
		builder.SetWithWings(v)
		if builder.withWings != v {
			t.Errorf("Invalid withWings flag. it supposed to be '%v'.", v)
		}
	}
}
func TestBugBuilderBuildMaterialWithoutWrapper(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic.")
			}
		}()
		builder := NewBugBuilder()
		builder.BuildMaterial()
	}()
}
func TestBugBuilderBuildMaterial(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Shouldn't have panic. '%v'.", r)
			}
		}()
		builder := NewBugBuilder()
		builder.SetWrapper(wrapperMock)
		builder.BuildMaterial()
	}()
}

func TestBug(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	scale := mgl32.Vec3{1.0, 1.0, 1.0}
	builder := NewBugBuilder()
	builder.SetPosition(position)
	builder.SetScale(scale)
	builder.SetWrapper(wrapperMock)
	bottomPosition := mgl32.Vec3{-1.0, 0.0, 0.0}
	eyeBase := float32(0.5773503)
	eye1Position := mgl32.Vec3{eyeBase, eyeBase, eyeBase}
	eye2Position := mgl32.Vec3{eyeBase, eyeBase, -eyeBase}

	bug := builder.BuildMaterial()
	if bug.GetLightSource() == nil {
		t.Error("Invalid lightsource")
	}

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
func TestBugWithoutLight(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	scale := mgl32.Vec3{1.0, 1.0, 1.0}
	builder := NewBugBuilder()
	builder.SetPosition(position)
	builder.SetScale(scale)
	builder.SetWrapper(wrapperMock)
	builder.SetWithLight(false)
	bottomPosition := mgl32.Vec3{-1.0, 0.0, 0.0}
	eyeBase := float32(0.5773503)
	eye1Position := mgl32.Vec3{eyeBase, eyeBase, eyeBase}
	eye2Position := mgl32.Vec3{eyeBase, eyeBase, -eyeBase}

	bug := builder.BuildMaterial()
	if bug.GetLightSource() != nil {
		t.Error("Invalid lightsource")
	}

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
				t.Errorf("Update shouldn't have panic. %v", r)
			}
		}()
		for i := 0; i < 15; i++ {
			bug.Update(float64(i * 100))
		}
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
