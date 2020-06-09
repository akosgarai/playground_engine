package light

import (
	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/playground_engine/pkg/transformations"
)

type Light struct {
	position mgl32.Vec3

	ambient  mgl32.Vec3
	diffuse  mgl32.Vec3
	specular mgl32.Vec3

	// in case of directional lights it's important.
	direction mgl32.Vec3

	// in case of point light sources we have to know the terms.
	constantTerm  float32
	linearTerm    float32
	quadraticTerm float32

	// spotlights
	cutoff      float32
	outerCutoff float32
}

// NewPointLight returns a Light with point light settings. The vectorComponent
// [4]mgl32.Vec3 input has to contain the 'position', 'ambient', 'diffuse', 'specular'
// component vectors in this order. The terms [3]float32 input has to
// contain the 'constant', 'linear', 'quadratic' term components in this order.
func NewPointLight(vectorComponents [4]mgl32.Vec3, terms [3]float32) *Light {
	return &Light{
		position: vectorComponents[0],
		ambient:  vectorComponents[1],
		diffuse:  vectorComponents[2],
		specular: vectorComponents[3],

		constantTerm:  terms[0],
		linearTerm:    terms[1],
		quadraticTerm: terms[2],
	}
}

// NewDirectionalLight returns a Light with directional light settings.
// The vectorComponent [4]mgl32.Vec3 input has to contain the 'direction',
// 'ambient', 'diffuse', 'specular' components in this order.
func NewDirectionalLight(vectorComponents [4]mgl32.Vec3) *Light {
	return &Light{
		direction: vectorComponents[0],
		ambient:   vectorComponents[1],
		diffuse:   vectorComponents[2],
		specular:  vectorComponents[3],
	}
}

// NewSpotLight returns a Light with spot light settings. The vectorComponent
// [5]mgl32.Vec3 input has to contain the 'position', 'direction', 'ambient',
// 'diffuse', 'specular' components in this order. The terms[5]float32 input has
// to contain the 'constant', 'linear', 'quadratic' terms, 'cutoff' and the
// 'outerCutoff' components in this order.
func NewSpotLight(vectorComponents [5]mgl32.Vec3, terms [5]float32) *Light {
	return &Light{
		position:  vectorComponents[0],
		direction: vectorComponents[1],
		ambient:   vectorComponents[2],
		diffuse:   vectorComponents[3],
		specular:  vectorComponents[4],

		constantTerm:  terms[0],
		linearTerm:    terms[1],
		quadraticTerm: terms[2],
		cutoff:        terms[3],
		outerCutoff:   terms[4],
	}
}

// Log returns the current state of the object
func (l *Light) Log() string {
	logString := "Light\n"
	logString += " - Position: Vector{" + transformations.Vec3ToString(l.position) + "}\n"
	logString += " - Ambient: Vector{" + transformations.Vec3ToString(l.ambient) + "}\n"
	logString += " - Diffuse: Vector{" + transformations.Vec3ToString(l.diffuse) + "}\n"
	logString += " - Specualar: Vector{" + transformations.Vec3ToString(l.specular) + "}\n"
	return logString
}

// GetDiffuse returns the diffuse color of the material
func (l *Light) GetAmbient() mgl32.Vec3 {
	return l.ambient
}

// GetDiffuse returns the diffuse color of the material
func (l *Light) GetDiffuse() mgl32.Vec3 {
	return l.diffuse
}

// GetSpecular returns the specular color of the material
func (l *Light) GetSpecular() mgl32.Vec3 {
	return l.specular
}

// GetPosition returns the shininess of the material
func (l *Light) GetPosition() mgl32.Vec3 {
	return l.position
}

// SetPosition updates the position of the light
func (l *Light) SetPosition(pos mgl32.Vec3) {
	l.position = pos
}

// GetConstantTerm returns the constant term component of the light
func (l *Light) GetConstantTerm() float32 {
	return l.constantTerm
}

// GetLinearTerm returns the linear term component of the light
func (l *Light) GetLinearTerm() float32 {
	return l.linearTerm
}

// GetQuadraticTerm returns the quadratic term component of the light
func (l *Light) GetQuadraticTerm() float32 {
	return l.quadraticTerm
}

// GetDirection returns the direction of the light
func (l *Light) GetDirection() mgl32.Vec3 {
	return l.direction
}

// GetCutoff returns the cutoff component of the light
func (l *Light) GetCutoff() float32 {
	return l.cutoff
}

// GetOuterCutoff returns the outerCutoff component of the light
func (l *Light) GetOuterCutoff() float32 {
	return l.outerCutoff
}
