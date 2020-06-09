package light

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	DefaultLightPosition     = mgl32.Vec3{0, 0, 0}
	DefaultLightDirection    = mgl32.Vec3{0, 1, 0}
	DefaultAmbientComponent  = mgl32.Vec3{1, 1, 1}
	DefaultDiffuseComponent  = mgl32.Vec3{0.2, 0.2, 0.2}
	DefaultSpecularComponent = mgl32.Vec3{0.4, 0.4, 0.4}

	DefaultConstantTerm  = float32(1.0)
	DefaultLinearTerm    = float32(0.5)
	DefaultQuadraticTerm = float32(0.05)
	DefaultCutoff        = float32(15)
	DefaultOuterCutoff   = float32(20)
)

func TestNew(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	l := NewDirectionalLight(vectorComponent)
	if l.position != DefaultLightPosition {
		t.Error("Invalid light position")
	}
	if l.ambient != DefaultAmbientComponent {
		t.Error("Invalid ambient component")
	}
	if l.diffuse != DefaultDiffuseComponent {
		t.Error("Invalid diffuse component")
	}
	if l.specular != DefaultSpecularComponent {
		t.Error("Invalid specular component")
	}
}
func TestLog(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	l := NewDirectionalLight(vectorComponent)
	log := l.Log()
	if len(log) < 10 {
		t.Error("Log too short")
	}
}
func TestGetAmbient(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	l := NewDirectionalLight(vectorComponent)
	if l.GetAmbient() != DefaultAmbientComponent {
		t.Error("Invalid ambient color")
	}
}
func TestGetDiffuse(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	l := NewDirectionalLight(vectorComponent)
	if l.GetDiffuse() != DefaultDiffuseComponent {
		t.Error("Invalid diffuse color")
	}
}
func TestGetSpecular(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	l := NewDirectionalLight(vectorComponent)
	if l.GetSpecular() != DefaultSpecularComponent {
		t.Error("Invalid specular color")
	}
}
func TestGetPosition(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	l := NewDirectionalLight(vectorComponent)
	if l.GetPosition() != DefaultLightPosition {
		t.Error("Invalid position vector")
	}
}
func TestSetPosition(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	l := NewDirectionalLight(vectorComponent)
	newPosition := mgl32.Vec3{2, 2, 2}
	l.SetPosition(newPosition)
	if l.GetPosition() != newPosition {
		t.Error("Invalid position vector")
	}
}
func TestGetDirection(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	l := NewDirectionalLight(vectorComponent)
	if l.GetDirection() != DefaultLightDirection {
		t.Error("Invalid direction component")
	}
}
func TestGetConstantTerm(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightPosition, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	termComponent := [3]float32{DefaultConstantTerm, DefaultLinearTerm, DefaultQuadraticTerm}
	l := NewPointLight(vectorComponent, termComponent)
	if l.GetConstantTerm() != DefaultConstantTerm {
		t.Errorf("Invalid constant term component. Instead of '%f', We have '%f'.", DefaultConstantTerm, l.GetConstantTerm())
	}
}
func TestGetLinearTerm(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightPosition, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	termComponent := [3]float32{DefaultConstantTerm, DefaultLinearTerm, DefaultQuadraticTerm}
	l := NewPointLight(vectorComponent, termComponent)
	if l.GetLinearTerm() != DefaultLinearTerm {
		t.Errorf("Invalid linear term component. Instead of '%f', We have '%f'.", DefaultLinearTerm, l.GetLinearTerm())
	}
}
func TestGetQuadraticTerm(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightPosition, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	termComponent := [3]float32{DefaultConstantTerm, DefaultLinearTerm, DefaultQuadraticTerm}
	l := NewPointLight(vectorComponent, termComponent)
	if l.GetQuadraticTerm() != DefaultQuadraticTerm {
		t.Errorf("Invalid quadratic term component. Instead of '%f', We have '%f'.", DefaultQuadraticTerm, l.GetQuadraticTerm())
	}
}
func TestGetCutOff(t *testing.T) {
	vectorComponent := [5]mgl32.Vec3{DefaultLightPosition, DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	termComponent := [5]float32{DefaultConstantTerm, DefaultLinearTerm, DefaultQuadraticTerm, DefaultCutoff}
	l := NewSpotLight(vectorComponent, termComponent)
	if l.GetCutoff() != DefaultCutoff {
		t.Errorf("Invalid cutoff component. Instead of '%f', We have '%f'.", DefaultCutoff, l.GetCutoff())
	}
}
func TestGetOuterCutOff(t *testing.T) {
	vectorComponent := [5]mgl32.Vec3{DefaultLightPosition, DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	termComponent := [5]float32{DefaultConstantTerm, DefaultLinearTerm, DefaultQuadraticTerm, DefaultCutoff, DefaultOuterCutoff}
	l := NewSpotLight(vectorComponent, termComponent)
	if l.GetOuterCutoff() != DefaultOuterCutoff {
		t.Errorf("Invalid cutoff component. Instead of '%f', We have '%f'.", DefaultOuterCutoff, l.GetOuterCutoff())
	}
}
func TestNewDirectionalLight(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	l := NewDirectionalLight(vectorComponent)
	if l.direction != DefaultLightDirection {
		t.Error("Invalid direction component")
	}
	if l.ambient != DefaultAmbientComponent {
		t.Error("Invalid ambient component")
	}
	if l.diffuse != DefaultDiffuseComponent {
		t.Error("Invalid diffuse component")
	}
	if l.specular != DefaultSpecularComponent {
		t.Error("Invalid specular component")
	}
}
func TestNewPointLight(t *testing.T) {
	vectorComponent := [4]mgl32.Vec3{DefaultLightPosition, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	termComponent := [3]float32{DefaultConstantTerm, DefaultLinearTerm, DefaultQuadraticTerm}
	l := NewPointLight(vectorComponent, termComponent)
	if l.position != DefaultLightPosition {
		t.Error("Invalid light position")
	}
	if l.ambient != DefaultAmbientComponent {
		t.Error("Invalid ambient component")
	}
	if l.diffuse != DefaultDiffuseComponent {
		t.Error("Invalid diffuse component")
	}
	if l.specular != DefaultSpecularComponent {
		t.Error("Invalid specular component")
	}
	if l.constantTerm != DefaultConstantTerm {
		t.Errorf("Invalid constant term component. Instead of '%f', We have '%f'.", DefaultConstantTerm, l.constantTerm)
	}
	if l.linearTerm != DefaultLinearTerm {
		t.Errorf("Invalid linear term component. Instead of '%f', We have '%f'.", DefaultLinearTerm, l.linearTerm)
	}
	if l.quadraticTerm != DefaultQuadraticTerm {
		t.Errorf("Invalid quadratic term component. Instead of '%f', We have '%f'.", DefaultQuadraticTerm, l.quadraticTerm)
	}
}
func TestNewSpotLight(t *testing.T) {
	vectorComponent := [5]mgl32.Vec3{DefaultLightPosition, DefaultLightDirection, DefaultAmbientComponent, DefaultDiffuseComponent, DefaultSpecularComponent}
	termComponent := [5]float32{DefaultConstantTerm, DefaultLinearTerm, DefaultQuadraticTerm, DefaultCutoff, DefaultOuterCutoff}
	l := NewSpotLight(vectorComponent, termComponent)
	if l.position != DefaultLightPosition {
		t.Error("Invalid light position")
	}
	if l.direction != DefaultLightDirection {
		t.Error("Invalid direction component")
	}
	if l.ambient != DefaultAmbientComponent {
		t.Error("Invalid ambient component")
	}
	if l.diffuse != DefaultDiffuseComponent {
		t.Error("Invalid diffuse component")
	}
	if l.specular != DefaultSpecularComponent {
		t.Error("Invalid specular component")
	}
	if l.constantTerm != DefaultConstantTerm {
		t.Errorf("Invalid constant term component. Instead of '%f', We have '%f'.", DefaultConstantTerm, l.constantTerm)
	}
	if l.linearTerm != DefaultLinearTerm {
		t.Errorf("Invalid linear term component. Instead of '%f', We have '%f'.", DefaultLinearTerm, l.linearTerm)
	}
	if l.quadraticTerm != DefaultQuadraticTerm {
		t.Errorf("Invalid quadratic term component. Instead of '%f', We have '%f'.", DefaultQuadraticTerm, l.quadraticTerm)
	}
	if l.cutoff != DefaultCutoff {
		t.Errorf("Invalid cutoff component. Instead of '%f', We have '%f'.", DefaultCutoff, l.cutoff)
	}
	if l.outerCutoff != DefaultOuterCutoff {
		t.Errorf("Invalid couterCutoff component. Instead of '%f', We have '%f'.", DefaultOuterCutoff, l.outerCutoff)
	}
}
