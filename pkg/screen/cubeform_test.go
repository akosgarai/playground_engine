package screen

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewCubeFormScreenBuilder(t *testing.T) {
	builder := NewCubeFormScreenBuilder()
	if builder.wrapper != nil {
		t.Error("Wrapper should be nil by default.")
	}
	if builder.camera != nil {
		t.Error("Camera should be nil by default.")
	}
	if builder.middleMonitorPosition != defaultMiddleMonitorPosition {
		t.Errorf("Invalid middleMonitorPosition. Instead of '%#v', it is '%#v'.", defaultMiddleMonitorPosition, builder.middleMonitorPosition)
	}
	if builder.leftMonitorRotationAngle != defaultLeftMonitorRotationAngle {
		t.Errorf("Invalid leftMonitorRotationAngle. Instead of '%f', it is '%f'.", defaultLeftMonitorRotationAngle, builder.leftMonitorRotationAngle)
	}
	if builder.rightMonitorRotationAngle != defaultRightMonitorRotationAngle {
		t.Errorf("Invalid rightMonitorRotationAngle. Instead of '%f', it is '%f'.", defaultRightMonitorRotationAngle, builder.rightMonitorRotationAngle)
	}
	if builder.middleMonitorSize != defaultMonitorSize {
		t.Errorf("Invalid middleMonitorSize. Instead of '%#v', it is '%#v'.", defaultMonitorSize, builder.middleMonitorSize)
	}
	if builder.leftMonitorSize != defaultMonitorSize {
		t.Errorf("Invalid leftMonitorSize. Instead of '%#v', it is '%#v'.", defaultMonitorSize, builder.leftMonitorSize)
	}
	if builder.rightMonitorSize != defaultMonitorSize {
		t.Errorf("Invalid rightMonitorSize. Instead of '%#v', it is '%#v'.", defaultMonitorSize, builder.rightMonitorSize)
	}
	if builder.middleScreenPosition != defaultMiddleScreenPosition {
		t.Errorf("Invalid middleScreenPosition. Instead of '%#v', it is '%#v'.", defaultMiddleScreenPosition, builder.middleScreenPosition)
	}
	if builder.leftScreenPosition != defaultLeftScreenPosition {
		t.Errorf("Invalid leftScreenPosition. Instead of '%#v', it is '%#v'.", defaultLeftScreenPosition, builder.leftScreenPosition)
	}
	if builder.rightScreenPosition != defaultRightScreenPosition {
		t.Errorf("Invalid rightScreenPosition. Instead of '%#v', it is '%#v'.", defaultRightScreenPosition, builder.rightScreenPosition)
	}
	if builder.middleScreenSize != defaultScreenSize {
		t.Errorf("Invalid middleScreenSize. Instead of '%#v', it is '%#v'.", defaultScreenSize, builder.middleScreenSize)
	}
	if builder.leftScreenSize != defaultScreenSize {
		t.Errorf("Invalid leftScreenSize. Instead of '%#v', it is '%#v'.", defaultScreenSize, builder.leftScreenSize)
	}
	if builder.rightScreenSize != defaultScreenSize {
		t.Errorf("Invalid rightScreenSize. Instead of '%#v', it is '%#v'.", defaultScreenSize, builder.rightScreenSize)
	}
	if builder.middleMonitorTexture != defaultMonitorTextureName {
		t.Errorf("Invalid middleMonitorTexture. Instead of '%s', it is '%s'.", defaultMonitorTextureName, builder.middleMonitorTexture)
	}
	if builder.leftMonitorTexture != defaultMonitorTextureName {
		t.Errorf("Invalid leftMonitorTexture. Instead of '%s', it is '%s'.", defaultMonitorTextureName, builder.leftMonitorTexture)
	}
	if builder.rightMonitorTexture != defaultMonitorTextureName {
		t.Errorf("Invalid rightMonitorTexture. Instead of '%s', it is '%s'.", defaultMonitorTextureName, builder.rightMonitorTexture)
	}
	if builder.screenMaterial != defaultScreenMaterial {
		t.Errorf("Invalid screenMaterial. Instead of '%#v', it is '%#v'.", defaultScreenMaterial, builder.screenMaterial)
	}
	if builder.tableMaterial != defaultTableMaterial {
		t.Errorf("Invalid tableMaterial. Instead of '%#v', it is '%#v'.", defaultTableMaterial, builder.tableMaterial)
	}
	if builder.tableSize != defaultTableSize {
		t.Errorf("Invalid tableSize. Instead of '%#v', it is '%#v'.", defaultTableSize, builder.tableSize)
	}
	if builder.tablePosition != defaultTablePosition {
		t.Errorf("Invalid tablePosition. Instead of '%#v', it is '%#v'.", defaultTablePosition, builder.tablePosition)
	}
	if builder.backgroundSize != defaultBackgroundSize {
		t.Errorf("Invalid backgroundSize. Instead of '%#v', it is '%#v'.", defaultBackgroundSize, builder.backgroundSize)
	}
	if builder.backgroundTexture != defaultBackgroundTextureName {
		t.Errorf("Invalid backgroundTexture. Instead of '%s', it is '%s'.", defaultBackgroundTextureName, builder.backgroundTexture)
	}
	if !reflect.DeepEqual(builder.controlPoints, defaultControlPoints) {
		t.Errorf("Invalid controlPoints. Instead of '%#v', it is '%#v'.", defaultControlPoints, builder.controlPoints)
	}
	if builder.clearColor != defaultClearColor {
		t.Errorf("Invalid clearColor. Instead of '%#v', it is '%#v'.", defaultClearColor, builder.clearColor)
	}
	if builder.lightDir != defaultLightDir {
		t.Errorf("Invalid lightDir. Instead of '%#v', it is '%#v'.", defaultLightDir, builder.lightDir)
	}
	if builder.lightAmbient != defaultLightAmbient {
		t.Errorf("Invalid lightAmbient. Instead of '%#v', it is '%#v'.", defaultLightAmbient, builder.lightAmbient)
	}
	if builder.lightDiffuse != defaultLightDiffuse {
		t.Errorf("Invalid lightDiffuse. Instead of '%#v', it is '%#v'.", defaultLightDiffuse, builder.lightDiffuse)
	}
	if builder.lightSpecular != defaultLightSpecular {
		t.Errorf("Invalid lightSpecular. Instead of '%#v', it is '%#v'.", defaultLightSpecular, builder.lightSpecular)
	}
}
func TestCubeFormScreenBuilderSetMiddleMonitorPosition(t *testing.T) {
	testData := []struct {
		value mgl32.Vec3
	}{
		{mgl32.Vec3{1.0, 0.0, 0.0}},
		{mgl32.Vec3{2.0, 0.0, 0.0}},
		{mgl32.Vec3{2.0, 0.0, 2.0}},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetMiddleMonitorPosition(tt.value)
		if builder.middleMonitorPosition != tt.value {
			t.Errorf("Invalid middleMonitorPosition. Instead of '%#v', it is '%#v'.", tt.value, builder.middleMonitorPosition)
		}
	}
}
func TestCubeFormScreenBuilderSetMonitorRotationAngles(t *testing.T) {
	testData := []struct {
		left  float32
		right float32
	}{
		{10, -10},
		{20, -20},
		{60, -60},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetMonitorRotationAngles(tt.left, tt.right)
		if builder.leftMonitorRotationAngle != tt.left {
			t.Errorf("Invalid leftMonitorRotationAngle. Instead of '%f', it is '%f'.", tt.left, builder.leftMonitorRotationAngle)
		}
		if builder.rightMonitorRotationAngle != tt.right {
			t.Errorf("Invalid rightMonitorRotationAngle. Instead of '%f', it is '%f'.", tt.right, builder.rightMonitorRotationAngle)
		}
	}
}
func TestCubeFormScreenBuilderSetScreenPositions(t *testing.T) {
	testData := []struct {
		left   mgl32.Vec3
		middle mgl32.Vec3
		right  mgl32.Vec3
	}{
		{mgl32.Vec3{1.0, 0.0, 0.0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{-1, 0, 0}},
		{mgl32.Vec3{2.0, 0.0, 0.0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{-2, 0, 0}},
		{mgl32.Vec3{2.0, 0.0, 2.0}, mgl32.Vec3{0, 0, 2}, mgl32.Vec3{-2, 0, 2}},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetScreenPositions(tt.left, tt.middle, tt.right)
		if builder.middleScreenPosition != tt.middle {
			t.Errorf("Invalid middleScreenPosition. Instead of '%#v', it is '%#v'.", tt.middle, builder.middleScreenPosition)
		}
		if builder.leftScreenPosition != tt.left {
			t.Errorf("Invalid leftScreenPosition. Instead of '%#v', it is '%#v'.", tt.left, builder.leftScreenPosition)
		}
		if builder.rightScreenPosition != tt.right {
			t.Errorf("Invalid rightScreenPosition. Instead of '%#v', it is '%#v'.", tt.right, builder.rightScreenPosition)
		}
	}
}
func TestCubeFormScreenBuilderSetScreenSizes(t *testing.T) {
	testData := []struct {
		left   mgl32.Vec3
		middle mgl32.Vec3
		right  mgl32.Vec3
	}{
		{mgl32.Vec3{1.0, 0.0, 0.0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{-1, 0, 0}},
		{mgl32.Vec3{2.0, 0.0, 0.0}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{-2, 0, 0}},
		{mgl32.Vec3{2.0, 0.0, 2.0}, mgl32.Vec3{0, 0, 2}, mgl32.Vec3{-2, 0, 2}},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetScreenSizes(tt.left, tt.middle, tt.right)
		if builder.middleScreenSize != tt.middle {
			t.Errorf("Invalid middleScreenSize. Instead of '%#v', it is '%#v'.", tt.middle, builder.middleScreenSize)
		}
		if builder.leftScreenSize != tt.left {
			t.Errorf("Invalid leftScreenSize. Instead of '%#v', it is '%#v'.", tt.left, builder.leftScreenSize)
		}
		if builder.rightScreenSize != tt.right {
			t.Errorf("Invalid rightScreenSize. Instead of '%#v', it is '%#v'.", tt.right, builder.rightScreenSize)
		}
	}
}
func TestCubeFormScreenBuilderSetAssetsDirectory(t *testing.T) {
	testData := []struct {
		value string
	}{
		{"testData/00"},
		{"testData/01"},
		{"testData/02"},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetAssetsDirectory(tt.value)
		if builder.assetsDirectory != tt.value {
			t.Errorf("Invalid assetsDirectory. Instead of '%s', it is '%s'.", tt.value, builder.assetsDirectory)
		}
	}
}
func TestCubeFormScreenBuilderSetMonitorTextureNames(t *testing.T) {
	testData := []struct {
		left   string
		middle string
		right  string
	}{
		{"testData.png", "testData02.png", "testData04.png"},
		{"testData.png", "testData03.png", "testData06.png"},
		{"testData.png", "testData05.png", "testData07.png"},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetMonitorTextureNames(tt.left, tt.middle, tt.right)
		if builder.middleMonitorTexture != tt.middle {
			t.Errorf("Invalid middleMonitorTexture. Instead of '%s', it is '%s'.", tt.middle, builder.middleMonitorTexture)
		}
		if builder.leftMonitorTexture != tt.left {
			t.Errorf("Invalid leftMonitorTexture. Instead of '%s', it is '%s'.", tt.left, builder.leftMonitorTexture)
		}
		if builder.rightMonitorTexture != tt.right {
			t.Errorf("Invalid rightMonitorTexture. Instead of '%s', it is '%s'.", tt.right, builder.rightMonitorTexture)
		}
	}
}
