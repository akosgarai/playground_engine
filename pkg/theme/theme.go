package theme

import (
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

type Theme struct {
	frameWidth              float32 // the size on the x axis
	frameLength             float32 // the size on the y axis
	frameTopLeftWidth       float32 // for supporting the header text, there is an option to give a padding here.
	detailContentBoxHeight  float32 // the size of the dcb in the y axis.
	frameMaterial           *material.Material
	menuItemDefaultMaterial *material.Material
	menuItemHoverMaterial   *material.Material
	menuItemSurfaceTexture  texture.Textures
	headerLabelColor        mgl32.Vec3
	labelColor              mgl32.Vec3
	inputColor              mgl32.Vec3
	backgroundColor         mgl32.Vec3
}

// SetFrameWidth sets the value of the frameWidth.
func (t *Theme) SetFrameWidth(w float32) {
	t.frameWidth = w
}

// GetFrameWidth returns the value of the frameWidth.
func (t *Theme) GetFrameWidth() float32 {
	return t.frameWidth
}

// SetFrameLength sets the value of the frameLength.
func (t *Theme) SetFrameLength(w float32) {
	t.frameLength = w
}

// GetFrameLength returns the value of the frameLength.
func (t *Theme) GetFrameLength() float32 {
	return t.frameLength
}

// SetFrameTopLeftWidth sets the value of the frameTopLeftWidth.
func (t *Theme) SetFrameTopLeftWidth(w float32) {
	t.frameTopLeftWidth = w
}

// GetFrameTopLeftWidth returns the value of the frameTopLeftWidth.
func (t *Theme) GetFrameTopLeftWidth() float32 {
	return t.frameTopLeftWidth
}

// SetDetailContentBoxHeight sets the value of the detailContentBoxHeight.
func (t *Theme) SetDetailContentBoxHeight(w float32) {
	t.detailContentBoxHeight = w
}

// GetDetailContentBoxHeight returns the value of the detailContentBoxHeight.
func (t *Theme) GetDetailContentBoxHeight() float32 {
	return t.detailContentBoxHeight
}

// SetFrameMaterial sets the frameMaterial.
func (t *Theme) SetFrameMaterial(m *material.Material) {
	t.frameMaterial = m
}

// GetFrameMaterial returns the frameMaterial.
func (t *Theme) GetFrameMaterial() *material.Material {
	return t.frameMaterial
}

// SetMenuItemDefaultMaterial sets the menuItemDefaultMaterial.
func (t *Theme) SetMenuItemDefaultMaterial(m *material.Material) {
	t.menuItemDefaultMaterial = m
}

// GetMenuItemDefaultMaterial returns the menuItemDefaultMaterial.
func (t *Theme) GetMenuItemDefaultMaterial() *material.Material {
	return t.menuItemDefaultMaterial
}

// SetMenuItemHoverMaterial sets the menuItemHoverMaterial.
func (t *Theme) SetMenuItemHoverMaterial(m *material.Material) {
	t.menuItemHoverMaterial = m
}

// GetMenuItemHoverMaterial returns the menuItemHoverMaterial.
func (t *Theme) GetMenuItemHoverMaterial() *material.Material {
	return t.menuItemHoverMaterial
}

// SetMenuItemSurfaceTexture sets the menuItemSurfaceTexture.
func (t *Theme) SetMenuItemSurfaceTexture(tex texture.Textures) {
	t.menuItemSurfaceTexture = tex
}

// GetMenuItemSurfaceTexture returns the menuItemSurfaceTexture.
func (t *Theme) GetMenuItemSurfaceTexture() texture.Textures {
	return t.menuItemSurfaceTexture
}

// SetHeaderLabelColor sets the value of the headerLabelColor.
func (t *Theme) SetHeaderLabelColor(c mgl32.Vec3) {
	t.headerLabelColor = c
}

// GetHeaderLabelColor returns the value of the headerLabelColor.
func (t *Theme) GetHeaderLabelColor() mgl32.Vec3 {
	return t.headerLabelColor
}

// SetLabelColor sets the value of the labelColor.
func (t *Theme) SetLabelColor(c mgl32.Vec3) {
	t.labelColor = c
}

// GetLabelColor returns the value of the labelColor.
func (t *Theme) GetLabelColor() mgl32.Vec3 {
	return t.labelColor
}

// SetInputColor sets the value of the inputColor.
func (t *Theme) SetInputColor(c mgl32.Vec3) {
	t.inputColor = c
}

// GetInputColor returns the value of the inputColor.
func (t *Theme) GetInputColor() mgl32.Vec3 {
	return t.inputColor
}

// SetBackgroundColor sets the value of the backgroundColor.
func (t *Theme) SetBackgroundColor(c mgl32.Vec3) {
	t.backgroundColor = c
}

// GetBackgroundColor returns the value of the backgroundColor.
func (t *Theme) GetBackgroundColor() mgl32.Vec3 {
	return t.backgroundColor
}
