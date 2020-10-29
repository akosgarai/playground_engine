package ui

import (
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	defaultUIButtonLabelText  = "Label"
	defaultUIButtonLabelSize  = float32(3.0)
	defaultUIButtonWidth      = float32(0.5)
	defaultUIButtonHeight     = float32(0.5)
	defaultUIButtonFrameWidth = float32(0.05)
)

var (
	defaultUIButtonLabelColor    = mgl32.Vec3{0, 0, 0}
	defaultUIButtonLabelPosition = mgl32.Vec3{0, 0, 0}

	defaultUIButtonDefaultMaterial = material.Chrome
	defaultUIButtonHoverMaterial   = material.Ruby
	defaultUIButtonOnStateMaterial = material.Emerald
)

type UIButtonBuilder struct {
	labelText       string               // This value is printed to the surface mesh
	labelColor      mgl32.Vec3           // The value will be printed with this color.
	labelPosition   mgl32.Vec3           // The position of the text (relative from the surface).
	labelSize       float32              // The size of the text.
	defaultMaterial *material.Material   // The material of the frame and the default state.
	hoverMaterial   *material.Material   // The material of the hovered state.
	onStateMaterial *material.Material   // The material of the on state of the button.
	buttonWidth     float32              // The full width of the button.
	buttonHeight    float32              // The full height of the button.
	frameWidth      float32              // The width of the frame. -> label surface width: buttonWidth - 2*frameWidth, label surface height: buttonHeight - 2*frameWidth.
	wrapper         interfaces.GLWrapper // The glwrapper that we need for the textures.
	transparent     texture.Textures     // The transparent texture for the meshes.
}

func NewUIButtonBuilder(w interfaces.GLWrapper) *UIButtonBuilder {
	var t texture.Textures
	return &UIButtonBuilder{
		labelText:       defaultUIButtonLabelText,
		labelColor:      defaultUIButtonLabelColor,
		labelPosition:   defaultUIButtonLabelPosition,
		labelSize:       defaultUIButtonLabelSize,
		defaultMaterial: defaultUIButtonDefaultMaterial,
		hoverMaterial:   defaultUIButtonHoverMaterial,
		onStateMaterial: defaultUIButtonOnStateMaterial,
		buttonWidth:     defaultUIButtonWidth,
		buttonHeight:    defaultUIButtonHeight,
		frameWidth:      defaultUIButtonFrameWidth,
		wrapper:         w,
		transparent:     t,
	}
}

// Setup label updates the label related parameters. The params are updated together.
func (b *UIButtonBuilder) SetupLabel(text string, color, position mgl32.Vec3, size float32) {
	b.labelText = text
	b.labelColor = color
	b.labelPosition = position
	b.labelSize = size
}

// SetupMaterials updates the material related parameters. The params are updated together.
// order of the inputs: defaultMaterial, hoverMaterial, onStateMaterial
func (b *UIButtonBuilder) SetupMaterials(def, hov, on *material.Material) {
	b.defaultMaterial = def
	b.hoverMaterial = hov
	b.onStateMaterial = on
}

// SetupSize updates the size related parameters. The paras are updated together.
// Order of the inputs: buttonWidth, buttonHeight, frameWidth.
func (b *UIButtonBuilder) SetupSize(w, h, fw float32) {
	b.buttonWidth = w
	b.buttonHeight = h
	b.frameWidth = fw
}

func (b *UIButtonBuilder) setupTransparentTexture() {
	b.transparent.TransparentTexture(1, 1, 128, "tex.diffuse", b.wrapper)
	b.transparent.TransparentTexture(1, 1, 128, "tex.specular", b.wrapper)
}
func (b *UIButtonBuilder) frameMesh() *mesh.TexturedMaterialMesh {
	rect := rectangle.NewExact(b.buttonWidth, b.buttonHeight)
	V, I, BO := rect.MeshInput()
	msh := mesh.NewTexturedMaterialMesh(V, I, b.transparent, b.defaultMaterial, b.wrapper)
	msh.SetBoundingObject(BO)
	return msh
}
func (b *UIButtonBuilder) surfaceMesh() *mesh.TexturedMaterialMesh {
	rect := rectangle.NewExact(b.buttonWidth-2*b.frameWidth, b.buttonHeight-2*b.frameWidth)
	V, I, _ := rect.MeshInput()
	return mesh.NewTexturedMaterialMesh(V, I, b.transparent, b.defaultMaterial, b.wrapper)
}

// Build returns the UIButton instance. On case of missing wrapper, it panics.
func (b *UIButtonBuilder) Build() *UIButton {
	if b.wrapper == nil {
		panic("Missing wrapper.")
	} else {
		b.setupTransparentTexture()
	}
	// frame mesh: tex. (transp.) mat. (defaultMaterial) rectangle bw*bh,
	frame := b.frameMesh()
	// label surface mesh: tex. (transp.) mat. (defaultMaterial) rectangle (buttonWidth - 2*frameWidth) * (buttonHeight - 2*frameWidth)
	surface := b.surfaceMesh()
	// base model
	m := model.New()
	m.AddMesh(frame)
	m.AddMesh(surface)
	return &UIButton{
		Model:           m,
		labelText:       b.labelText,
		labelColor:      b.labelColor,
		labelPosition:   b.labelPosition,
		labelSize:       b.labelSize,
		defaultMaterial: b.defaultMaterial,
		hoverMaterial:   b.hoverMaterial,
		onStateMaterial: b.onStateMaterial,
	}
}

type UIButton struct {
	Model         *model.BaseModel
	labelText     string     // This value is printed to the surface mesh
	labelColor    mgl32.Vec3 // The value will be printed with this color.
	labelPosition mgl32.Vec3 // The position of the text (relative from the surface).
	labelSize     float32    // The size of the text.
	state         string
	// materials
	defaultMaterial *material.Material
	hoverMaterial   *material.Material
	onStateMaterial *material.Material
}

// LabelSurface returns the mesh that is the target surface of the label.
// meshes[0] - background mesh - the frame of the button
// meshes[1] - the surface mesh. Its material is updatable.
func (m *UIButton) LabelSurface() interfaces.Mesh {
	msh, _ := m.Model.GetMeshByIndex(1)
	return msh
}

// LabelColor returns the colors that are used as font colors.
func (m *UIButton) LabelColor() []mgl32.Vec3 {
	return []mgl32.Vec3{m.labelColor}
}

// LabelText returns the text of the label.
func (m *UIButton) LabelText() string {
	return m.labelText
}

// LabelPosition returns the position of the label text.
func (m *UIButton) LabelPosition() mgl32.Vec3 {
	return m.labelPosition
}

// LabelSize returns the size of the label.
func (m *UIButton) LabelSize() float32 {
	return m.labelSize
}

// Hover event, the label surface material is changed to the hover value.
func (m *UIButton) Hover() {
	msh := (m.LabelSurface()).(*mesh.TexturedMaterialMesh)
	msh.Material = m.hoverMaterial
}

// Clear makes the label surface material to be changed to the default value.
func (m *UIButton) Clear() {
	if m.state != "on" {
		msh := (m.LabelSurface()).(*mesh.TexturedMaterialMesh)
		msh.Material = m.defaultMaterial
	}
}

// OnState updates the state to 'on' value and the material to onStateMaterial value.
func (m *UIButton) OnState() {
	m.state = "on"
	msh := (m.LabelSurface()).(*mesh.TexturedMaterialMesh)
	msh.Material = m.onStateMaterial
}

// OffState updates the state to 'off' value.
func (m *UIButton) OffState() {
	m.state = "off"
}
