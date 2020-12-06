package ui

import (
	"reflect"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	wrapperMock testhelper.GLWrapperMock
)

func TestNewUIButtonBuilder(t *testing.T) {
	builder := NewUIButtonBuilder(wrapperMock)
	if builder.labelText != defaultUIButtonLabelText {
		t.Errorf("Invalid labelText. Instead of '%s', it is '%s'.", defaultUIButtonLabelText, builder.labelText)
	}
	if builder.labelColor != defaultUIButtonLabelColor {
		t.Errorf("Invalid labelColor. Instead of '%#v', it is '%#v'.", defaultUIButtonLabelColor, builder.labelColor)
	}
	if builder.labelPosition != defaultUIButtonLabelPosition {
		t.Errorf("Invalid labelPosition. Instead of '%#v', it is '%#v'.", defaultUIButtonLabelPosition, builder.labelPosition)
	}
	if builder.labelSize != defaultUIButtonLabelSize {
		t.Errorf("Invalid labelSize. Instead of '%f', it is '%f'.", defaultUIButtonLabelSize, builder.labelSize)
	}
	if builder.defaultMaterial != defaultUIButtonDefaultMaterial {
		t.Errorf("Invalid defaultMaterial. Instead of '%#v', it is '%#v'.", defaultUIButtonDefaultMaterial, builder.defaultMaterial)
	}
	if builder.hoverMaterial != defaultUIButtonHoverMaterial {
		t.Errorf("Invalid hoverMaterial. Instead of '%#v', it is '%#v'.", defaultUIButtonHoverMaterial, builder.hoverMaterial)
	}
	if builder.onStateMaterial != defaultUIButtonOnStateMaterial {
		t.Errorf("Invalid onStateMaterial. Instead of '%#v', it is '%#v'.", defaultUIButtonOnStateMaterial, builder.onStateMaterial)
	}
	if builder.buttonWidth != defaultUIButtonWidth {
		t.Errorf("Invalid buttonWidth. Instead of '%f', it is '%f'.", defaultUIButtonWidth, builder.buttonWidth)
	}
	if builder.buttonHeight != defaultUIButtonHeight {
		t.Errorf("Invalid buttonHeight. Instead of '%f', it is '%f'.", defaultUIButtonHeight, builder.buttonHeight)
	}
	if builder.frameWidth != defaultUIButtonFrameWidth {
		t.Errorf("Invalid frameWidth. Instead of '%f', it is '%f'.", defaultUIButtonFrameWidth, builder.frameWidth)
	}
}
func TestNewUIButtonBuilderSetupLabel(t *testing.T) {
	testData := []struct {
		text     string
		color    mgl32.Vec3
		position mgl32.Vec3
		size     float32
	}{
		{"test label", mgl32.Vec3{1, 1, 1}, mgl32.Vec3{0.5, 0, 0}, 2.0},
		{"test label 02", mgl32.Vec3{1, 1, 1}, mgl32.Vec3{-0.5, 0, 0}, 2.0},
		{"test label 03", mgl32.Vec3{1, 0, 0}, mgl32.Vec3{-0.5, 0, 0}, 1.3},
	}
	builder := NewUIButtonBuilder(wrapperMock)
	for _, tt := range testData {
		builder.SetupLabel(tt.text, tt.color, tt.position, tt.size)
		if builder.labelText != tt.text {
			t.Errorf("Invalid labelText. Instead of '%s', it is '%s'.", tt.text, builder.labelText)
		}
		if builder.labelColor != tt.color {
			t.Errorf("Invalid labelColor. Instead of '%#v', it is '%#v'.", tt.color, builder.labelColor)
		}
		if builder.labelPosition != tt.position {
			t.Errorf("Invalid labelPosition. Instead of '%#v', it is '%#v'.", tt.position, builder.labelPosition)
		}
		if builder.labelSize != tt.size {
			t.Errorf("Invalid labelSize. Instead of '%f', it is '%f'.", tt.size, builder.labelSize)
		}
	}
}
func TestNewUIButtonBuilderSetupMaterials(t *testing.T) {
	testData := []struct {
		d *material.Material
		h *material.Material
		o *material.Material
	}{
		{material.Jade, material.Emerald, material.Ruby},
		{material.Emerald, material.Jade, material.Ruby},
		{material.Ruby, material.Jade, material.Emerald},
	}
	builder := NewUIButtonBuilder(wrapperMock)
	for _, tt := range testData {
		builder.SetupMaterials(tt.d, tt.h, tt.o)
		if builder.defaultMaterial != tt.d {
			t.Errorf("Invalid defaultMaterial. Instead of '%#v', it is '%#v'.", tt.d, builder.defaultMaterial)
		}
		if builder.hoverMaterial != tt.h {
			t.Errorf("Invalid hoverMaterial. Instead of '%#v', it is '%#v'.", tt.h, builder.hoverMaterial)
		}
		if builder.onStateMaterial != tt.o {
			t.Errorf("Invalid onStateMaterial. Instead of '%#v', it is '%#v'.", tt.o, builder.onStateMaterial)
		}
	}
}
func TestNewUIButtonBuilderSetupSize(t *testing.T) {
	testData := []struct {
		w  float32
		h  float32
		fw float32
	}{
		{1.0, 1.0, 0.1},
		{2.0, 2.0, 0.2},
		{2.0, 0.5, 0.05},
	}
	builder := NewUIButtonBuilder(wrapperMock)
	for _, tt := range testData {
		builder.SetupSize(tt.w, tt.h, tt.fw)
		if builder.buttonWidth != tt.w {
			t.Errorf("Invalid buttonWidth. Instead of '%f', it is '%f'.", tt.w, builder.buttonWidth)
		}
		if builder.buttonHeight != tt.h {
			t.Errorf("Invalid buttonHeight. Instead of '%f', it is '%f'.", tt.h, builder.buttonHeight)
		}
		if builder.frameWidth != tt.fw {
			t.Errorf("Invalid frameWidth. Instead of '%f', it is '%f'.", tt.fw, builder.frameWidth)
		}
	}
}
func TestNewUIButtonBuilderBuildWithoutWrapper(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("It should fail due to the missing wrapper")
			}
		}()
		builder := NewUIButtonBuilder(wrapperMock)
		builder.wrapper = nil
		builder.Build()
	}()
}
func TestNewUIButtonBuilderBuild(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("It shouldn't panic. '%#v'.", r)
			}
		}()
		builder := NewUIButtonBuilder(wrapperMock)
		builder.Build()
	}()
}
func TestUIButtonLabelSurface(t *testing.T) {
	builder := NewUIButtonBuilder(wrapperMock)
	btn := builder.Build()
	surface := btn.LabelSurface()
	// it has to be textured material mesh
	if _, ok := surface.(*mesh.TexturedMaterialMesh); !ok {
		t.Error("The surface supposed to be textured material mesh.")
	}
}
func TestUIButtonLabelColor(t *testing.T) {
	newColor := mgl32.Vec3{0, 0, 1}
	builder := NewUIButtonBuilder(wrapperMock)
	builder.SetupLabel(defaultUIButtonLabelText, newColor, defaultUIButtonLabelPosition, defaultUIButtonLabelSize)
	btn := builder.Build()
	cols := []mgl32.Vec3{newColor}
	if !reflect.DeepEqual(btn.LabelColor(), cols) {
		t.Errorf("Invalid LabelColor. Instead of '%#v', it is '%#v'.", newColor, btn.LabelColor())
	}
}
func TestUIButtonLabelText(t *testing.T) {
	newText := "Button label new one."
	builder := NewUIButtonBuilder(wrapperMock)
	builder.SetupLabel(newText, defaultUIButtonLabelColor, defaultUIButtonLabelPosition, defaultUIButtonLabelSize)
	btn := builder.Build()
	if btn.LabelText() != newText {
		t.Errorf("Invalid LabelText. Instead of '%s', it is '%s'.", newText, btn.LabelText())
	}
}
func TestUIButtonLabelPosition(t *testing.T) {
	newPosition := mgl32.Vec3{0, 0, 1}
	builder := NewUIButtonBuilder(wrapperMock)
	builder.SetupLabel(defaultUIButtonLabelText, defaultUIButtonLabelColor, newPosition, defaultUIButtonLabelSize)
	btn := builder.Build()
	if btn.LabelPosition() != newPosition {
		t.Errorf("Invalid LabelPosition. Instead of '%#v', it is '%#v'.", newPosition, btn.LabelPosition())
	}
}
func TestUIButtonLabelSize(t *testing.T) {
	newSize := float32(5)
	builder := NewUIButtonBuilder(wrapperMock)
	builder.SetupLabel(defaultUIButtonLabelText, defaultUIButtonLabelColor, defaultUIButtonLabelPosition, newSize)
	btn := builder.Build()
	if btn.LabelSize() != newSize {
		t.Errorf("Invalid LabelSize. Instead of '%f', it is '%f'.", newSize, btn.LabelSize())
	}
}
func TestUIButtonHover(t *testing.T) {
	builder := NewUIButtonBuilder(wrapperMock)
	btn := builder.Build()
	btn.Hover()
	surface := btn.LabelSurface()
	if surface.(*mesh.TexturedMaterialMesh).Material != defaultUIButtonHoverMaterial {
		t.Error("Invalid material in hover state.")
	}
}
func TestUIButtonClear(t *testing.T) {
	builder := NewUIButtonBuilder(wrapperMock)
	btn := builder.Build()
	surface := btn.LabelSurface()
	btn.Hover()
	if surface.(*mesh.TexturedMaterialMesh).Material != defaultUIButtonHoverMaterial {
		t.Error("Invalid material in hover state.")
	}
	btn.Clear()
	if surface.(*mesh.TexturedMaterialMesh).Material != defaultUIButtonDefaultMaterial {
		t.Error("Invalid material in default state.")
	}
}
func TestUIButtonOnState(t *testing.T) {
	builder := NewUIButtonBuilder(wrapperMock)
	btn := builder.Build()
	btn.OnState()
	if btn.state != "on" {
		t.Error("Invalid state.")
	}
	surface := btn.LabelSurface()
	if surface.(*mesh.TexturedMaterialMesh).Material != defaultUIButtonOnStateMaterial {
		t.Error("Invalid material in on state.")
	}
}
func TestUIButtonOffState(t *testing.T) {
	builder := NewUIButtonBuilder(wrapperMock)
	btn := builder.Build()
	btn.OnState()
	if btn.state != "on" {
		t.Error("Invalid state.")
	}
	surface := btn.LabelSurface()
	if surface.(*mesh.TexturedMaterialMesh).Material != defaultUIButtonOnStateMaterial {
		t.Error("Invalid material in on state.")
	}
	btn.OffState()
	if btn.state != "off" {
		t.Error("Invalid state.")
	}
}
