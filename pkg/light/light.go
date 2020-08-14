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

// GetAmbient returns the ambient color component of the light
func (l *Light) GetAmbient() mgl32.Vec3 {
	return l.ambient
}

// SetAmbient updates the ambient color component of the light
func (l *Light) SetAmbient(a mgl32.Vec3) {
	l.ambient = a
}

// GetDiffuse returns the diffuse color component of the light
func (l *Light) GetDiffuse() mgl32.Vec3 {
	return l.diffuse
}

// SetDiffuse updates the diffuse color component of the light
func (l *Light) SetDiffuse(d mgl32.Vec3) {
	l.diffuse = d
}

// GetSpecular returns the specular color component of the light
func (l *Light) GetSpecular() mgl32.Vec3 {
	return l.specular
}

// SetSpecular updates the specular color component of the light
func (l *Light) SetSpecular(s mgl32.Vec3) {
	l.specular = s
}

// GetPosition returns the position of the light
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

// SetConstantTerm updates the constant term component of the light
func (l *Light) SetConstantTerm(c float32) {
	l.constantTerm = c
}

// GetLinearTerm returns the linear term component of the light
func (l *Light) GetLinearTerm() float32 {
	return l.linearTerm
}

// SetLinearTerm updates the linear term component of the light
func (l *Light) SetLinearTerm(lt float32) {
	l.linearTerm = lt
}

// GetQuadraticTerm returns the quadratic term component of the light
func (l *Light) GetQuadraticTerm() float32 {
	return l.quadraticTerm
}

// SetQuadraticTerm updates the quadratic term component of the light
func (l *Light) SetQuadraticTerm(q float32) {
	l.quadraticTerm = q
}

// GetDirection returns the direction of the light
func (l *Light) GetDirection() mgl32.Vec3 {
	return l.direction
}

// SetDirection updates the direction of the light
func (l *Light) SetDirection(d mgl32.Vec3) {
	l.direction = d
}

// GetCutoff returns the cutoff component of the light
func (l *Light) GetCutoff() float32 {
	return l.cutoff
}

// SetCutoff updates the cutoff component of the light
func (l *Light) SetCutoff(c float32) {
	l.cutoff = c
}

// GetOuterCutoff returns the outerCutoff component of the light
func (l *Light) GetOuterCutoff() float32 {
	return l.outerCutoff
}

// SetOuterCutoff updates the outerCutoff component of the light
func (l *Light) SetOuterCutoff(oc float32) {
	l.outerCutoff = oc
}
