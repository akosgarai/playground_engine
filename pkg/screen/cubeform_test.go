package screen

import (
	"reflect"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/material"

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
func TestCubeFormScreenBuilderSetScreenMaterial(t *testing.T) {
	testData := []struct {
		value *material.Material
	}{
		{material.Jade},
		{material.Ruby},
		{material.Gold},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetScreenMaterial(tt.value)
		if builder.screenMaterial != tt.value {
			t.Errorf("Invalid screenMaterial. Instead of '%#v', it is '%#v'.", tt.value, builder.screenMaterial)
		}
	}
}
func TestCubeFormScreenBuilderSetTableMaterial(t *testing.T) {
	testData := []struct {
		value *material.Material
	}{
		{material.Jade},
		{material.Ruby},
		{material.Gold},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetTableMaterial(tt.value)
		if builder.tableMaterial != tt.value {
			t.Errorf("Invalid tableMaterial. Instead of '%#v', it is '%#v'.", tt.value, builder.tableMaterial)
		}
	}
}
func TestCubeFormScreenBuilderSetTableSize(t *testing.T) {
	testData := []struct {
		value mgl32.Vec3
	}{
		{mgl32.Vec3{1.0, 0.0, 0.0}},
		{mgl32.Vec3{2.0, 0.0, 0.0}},
		{mgl32.Vec3{2.0, 0.0, 2.0}},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetTableSize(tt.value)
		if builder.tableSize != tt.value {
			t.Errorf("Invalid tableSize. Instead of '%#v', it is '%#v'.", tt.value, builder.tableSize)
		}
	}
}
func TestCubeFormScreenBuilderSetTablePosition(t *testing.T) {
	testData := []struct {
		value mgl32.Vec3
	}{
		{mgl32.Vec3{1.0, 0.0, 0.0}},
		{mgl32.Vec3{2.0, 0.0, 0.0}},
		{mgl32.Vec3{2.0, 0.0, 2.0}},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetTablePosition(tt.value)
		if builder.tablePosition != tt.value {
			t.Errorf("Invalid tablePosition. Instead of '%#v', it is '%#v'.", tt.value, builder.tablePosition)
		}
	}
}
func TestCubeFormScreenBuilderSetBackgroundSize(t *testing.T) {
	testData := []struct {
		value mgl32.Vec3
	}{
		{mgl32.Vec3{1.0, 0.0, 0.0}},
		{mgl32.Vec3{2.0, 0.0, 0.0}},
		{mgl32.Vec3{2.0, 0.0, 2.0}},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetBackgroundSize(tt.value)
		if builder.backgroundSize != tt.value {
			t.Errorf("Invalid backgroundSize. Instead of '%#v', it is '%#v'.", tt.value, builder.backgroundSize)
		}
	}
}
func TestCubeFormScreenBuilderSetBackgroundTextureName(t *testing.T) {
	testData := []struct {
		value string
	}{
		{"testData00.png"},
		{"testData01.png"},
		{"testData02.png"},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetBackgroundTextureName(tt.value)
		if builder.backgroundTexture != tt.value {
			t.Errorf("Invalid backgroundTexture. Instead of '%s', it is '%s'.", tt.value, builder.backgroundTexture)
		}
	}
}
func TestCubeFormScreenBuilderSetWindowSize(t *testing.T) {
	testData := []struct {
		w float32
		h float32
	}{
		{10, 10},
		{20, 20},
		{60, 60},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetWindowSize(tt.w, tt.h)
		if builder.windowWidth != tt.w {
			t.Errorf("Invalid windowWidth. Instead of '%f', it is '%f'.", tt.w, builder.windowWidth)
		}
		if builder.windowHeight != tt.h {
			t.Errorf("Invalid windowHeight. Instead of '%f', it is '%f'.", tt.h, builder.windowHeight)
		}
	}
}
func TestCubeFormScreenBuilderSetWrapper(t *testing.T) {
	builder := NewCubeFormScreenBuilder()
	builder.SetWrapper(wrapperMock)
	if builder.wrapper != wrapperMock {
		t.Error("Invalid wrapper")
	}
}
func TestCubeFormScreenBuilderSetControlPoints(t *testing.T) {
	testData := []struct {
		value []mgl32.Vec3
	}{
		{[]mgl32.Vec3{mgl32.Vec3{1.0, 0.0, 0.0}, mgl32.Vec3{0, 0, 0}}},
		{[]mgl32.Vec3{mgl32.Vec3{0.0, 1.0, 0.0}, mgl32.Vec3{0, 0, 0}}},
		{[]mgl32.Vec3{mgl32.Vec3{0.0, 0.0, 1.0}, mgl32.Vec3{0, 0, 0}}},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetControlPoints(tt.value)
		if !reflect.DeepEqual(builder.controlPoints, tt.value) {
			t.Errorf("Invalid controlPoints. Instead of '%#v', it is '%#v'.", tt.value, builder.controlPoints)
		}
	}
}
func TestCubeFormScreenBuilderSetClearColor(t *testing.T) {
	testData := []struct {
		value mgl32.Vec3
	}{
		{mgl32.Vec3{1.0, 0.0, 0.0}},
		{mgl32.Vec3{0.0, 1.0, 0.0}},
		{mgl32.Vec3{0.0, 0.0, 1.0}},
	}
	builder := NewCubeFormScreenBuilder()
	for _, tt := range testData {
		builder.SetClearColor(tt.value)
		if builder.clearColor != tt.value {
			t.Errorf("Invalid clearColor. Instead of '%#v', it is '%#v'.", tt.value, builder.clearColor)
		}
	}
}
