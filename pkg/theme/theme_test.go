package theme

import (
	"reflect"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

func TestThemeFrameWidth(t *testing.T) {
	var th Theme
	th.frameWidth = 2.0
	if th.frameWidth != 2.0 {
		t.Error("Invalid frameWidth.")
	}
	newVal := float32(3.0)
	th.SetFrameWidth(newVal)
	if th.frameWidth != newVal {
		t.Error("Invalid frameWidth has been set.")
	}
	if th.GetFrameWidth() != newVal {
		t.Error("Invalid frameWidth has been get.")
	}
}
func TestThemeFrameLength(t *testing.T) {
	var th Theme
	th.frameLength = 2.0
	if th.frameLength != 2.0 {
		t.Error("Invalid frameLength.")
	}
	newVal := float32(3.0)
	th.SetFrameLength(newVal)
	if th.frameLength != newVal {
		t.Error("Invalid frameLength has been set.")
	}
	if th.GetFrameLength() != newVal {
		t.Error("Invalid frameLength has been get.")
	}
}
func TestThemeFrameTopLeftWidth(t *testing.T) {
	var th Theme
	th.frameTopLeftWidth = 2.0
	if th.frameTopLeftWidth != 2.0 {
		t.Error("Invalid frameTopLeftWidth.")
	}
	newVal := float32(3.0)
	th.SetFrameTopLeftWidth(newVal)
	if th.frameTopLeftWidth != newVal {
		t.Error("Invalid frameTopLeftWidth has been set.")
	}
	if th.GetFrameTopLeftWidth() != newVal {
		t.Error("Invalid frameTopLeftWidth has been get.")
	}
}
func TestThemeDetailContentBoxHeight(t *testing.T) {
	var th Theme
	th.detailContentBoxHeight = 2.0
	if th.detailContentBoxHeight != 2.0 {
		t.Error("Invalid detailContentBoxHeight.")
	}
	newVal := float32(3.0)
	th.SetDetailContentBoxHeight(newVal)
	if th.detailContentBoxHeight != newVal {
		t.Error("Invalid detailContentBoxHeight has been set.")
	}
	if th.GetDetailContentBoxHeight() != newVal {
		t.Error("Invalid detailContentBoxHeight has been get.")
	}
}
func TestThemeFrameMaterial(t *testing.T) {
	var th Theme
	defaultMaterial := material.Jade
	newMaterial := material.Ruby
	th.frameMaterial = defaultMaterial
	if th.frameMaterial != defaultMaterial {
		t.Error("Invalid frameMaterial.")
	}
	th.SetFrameMaterial(newMaterial)
	if th.frameMaterial != newMaterial {
		t.Error("Invalid frameMaterial has been set.")
	}
	if th.GetFrameMaterial() != newMaterial {
		t.Error("Invalid frameMaterial has been get.")
	}
}
func TestThemeMenuItemDefaultMaterial(t *testing.T) {
	var th Theme
	defaultMaterial := material.Jade
	newMaterial := material.Ruby
	th.menuItemDefaultMaterial = defaultMaterial
	if th.menuItemDefaultMaterial != defaultMaterial {
		t.Error("Invalid menuItemDefaultMaterial.")
	}
	th.SetMenuItemDefaultMaterial(newMaterial)
	if th.menuItemDefaultMaterial != newMaterial {
		t.Error("Invalid menuItemDefaultMaterial has been set.")
	}
	if th.GetMenuItemDefaultMaterial() != newMaterial {
		t.Error("Invalid menuItemDefaultMaterial has been get.")
	}
}
func TestThemeMenuItemHoverMaterial(t *testing.T) {
	var th Theme
	defaultMaterial := material.Jade
	newMaterial := material.Ruby
	th.menuItemHoverMaterial = defaultMaterial
	if th.menuItemHoverMaterial != defaultMaterial {
		t.Error("Invalid menuItemHoverMaterial.")
	}
	th.SetMenuItemHoverMaterial(newMaterial)
	if th.menuItemHoverMaterial != newMaterial {
		t.Error("Invalid menuItemHoverMaterial has been set.")
	}
	if th.GetMenuItemHoverMaterial() != newMaterial {
		t.Error("Invalid menuItemHoverMaterial has been get.")
	}
}
func TestThemeMenuItemSurfaceTexture(t *testing.T) {
	var th Theme
	var newTexture, defaultTexture texture.Textures
	th.menuItemSurfaceTexture = defaultTexture
	if !reflect.DeepEqual(th.menuItemSurfaceTexture, defaultTexture) {
		t.Error("Invalid menuItemSurfaceTexture.")
	}
	th.SetMenuItemSurfaceTexture(newTexture)
	if !reflect.DeepEqual(th.menuItemSurfaceTexture, newTexture) {
		t.Error("Invalid menuItemSurfaceTexture has been set.")
	}
	if !reflect.DeepEqual(th.GetMenuItemSurfaceTexture(), newTexture) {
		t.Error("Invalid menuItemSurfaceTexture has been get.")
	}
}
func TestThemeHeaderLabelColor(t *testing.T) {
	var th Theme
	defaultColor := mgl32.Vec3{0, 1, 0}
	newColor := mgl32.Vec3{1, 1, 0}
	th.headerLabelColor = defaultColor
	if th.headerLabelColor != defaultColor {
		t.Error("Invalid headerLabelColor.")
	}
	th.SetHeaderLabelColor(newColor)
	if th.headerLabelColor != newColor {
		t.Error("Invalid headerLabelColor has been set.")
	}
	if th.GetHeaderLabelColor() != newColor {
		t.Error("Invalid headerLabelColor has been get.")
	}
}
func TestThemeLabelColor(t *testing.T) {
	var th Theme
	defaultColor := mgl32.Vec3{0, 1, 0}
	newColor := mgl32.Vec3{1, 1, 0}
	th.labelColor = defaultColor
	if th.labelColor != defaultColor {
		t.Error("Invalid labelColor.")
	}
	th.SetLabelColor(newColor)
	if th.labelColor != newColor {
		t.Error("Invalid labelColor has been set.")
	}
	if th.GetLabelColor() != newColor {
		t.Error("Invalid labelColor has been get.")
	}
}
func TestThemeInputColor(t *testing.T) {
	var th Theme
	defaultColor := mgl32.Vec3{0, 1, 0}
	newColor := mgl32.Vec3{1, 1, 0}
	th.inputColor = defaultColor
	if th.inputColor != defaultColor {
		t.Error("Invalid inputColor.")
	}
	th.SetInputColor(newColor)
	if th.inputColor != newColor {
		t.Error("Invalid inputColor has been set.")
	}
	if th.GetInputColor() != newColor {
		t.Error("Invalid inputColor has been get.")
	}
}
func TestThemeBackgroundColor(t *testing.T) {
	var th Theme
	defaultColor := mgl32.Vec3{0, 1, 0}
	newColor := mgl32.Vec3{1, 1, 0}
	th.backgroundColor = defaultColor
	if th.backgroundColor != defaultColor {
		t.Error("Invalid backgroundColor.")
	}
	th.SetBackgroundColor(newColor)
	if th.backgroundColor != newColor {
		t.Error("Invalid backgroundColor has been set.")
	}
	if th.GetBackgroundColor() != newColor {
		t.Error("Invalid backgroundColor has been get.")
	}
}
