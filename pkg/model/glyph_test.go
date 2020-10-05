package model

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/go-gl/mathgl/mgl32"
)

func TestCharsetLoad(t *testing.T) {
	_, err := LoadCharsetDebug("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during debug load: %s\n", err.Error())
	}
	_, err = LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during load: %s\n", err.Error())
	}
}
func TestCharsetPrintTo(t *testing.T) {
	cols := []mgl32.Vec3{mgl32.Vec3{0, 0, 1}}
	fonts, err := LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during load: %s\n", err.Error())
	}
	msh := mesh.NewPointMesh(wrapperMock)
	fonts.PrintTo("Hello", 0, 0, 0.0, 1.0, wrapperMock, msh, cols)
	fonts.Debug = true
	fonts.PrintTo("Hello", 0, 0, 0.0, 1.0, wrapperMock, msh, cols)
}
func TestCharsetCleanSurface(t *testing.T) {
	cols := []mgl32.Vec3{mgl32.Vec3{0, 0, 1}}
	fonts, err := LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during load: %s\n", err.Error())
	}
	msh := mesh.NewPointMesh(wrapperMock)
	fonts.PrintTo("Hello", 0, 0, 0.0, 1.0, wrapperMock, msh, cols)
	if len(fonts.meshes) != 5 {
		t.Errorf("Invalid number of meshes. Instead of '%d', we have '%d'.", 5, len(fonts.meshes))
	}
	fonts.CleanSurface(msh)
	if len(fonts.meshes) != 0 {
		t.Errorf("Invalid number of meshes. Instead of '%d', we have '%d'.", 0, len(fonts.meshes))
	}
}
func TestCharsetTextWidth(t *testing.T) {
	fonts, err := LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during load: %s\n", err.Error())
	}
	testData := []struct {
		text  string
		scale float32
		width float32
	}{
		{"a", 1, 25},
		{"a", 2, 50},
		{"a", 0.5, 12.5},
		{"b", 1, 22},
		{"1", 1, 15},
		{"b1", 1, 37},
	}
	for _, tt := range testData {
		width := fonts.TextWidth(tt.text, tt.scale)
		if width != tt.width {
			t.Errorf("Invalid text width for '%s'. Instead of '%f', we have '%f'.", tt.text, tt.width, width)
		}
	}
}
func TestCharsetTextContainerSize(t *testing.T) {
	fonts, err := LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40, 72, wrapperMock)
	if err != nil {
		t.Errorf("Error during load: %s\n", err.Error())
	}
	testData := []struct {
		text   string
		scale  float32
		width  float32
		height float32
	}{
		{"a", 1, 25, 18},
		{"a", 2, 50, 36},
		{"a", 0.5, 12.5, 9},
		{"b", 1, 22, 30},
		{"1", 1, 15, 21},
		{"b1", 1, 37, 30},
		{"Spotlight editor", float32(3) / float32(800), 1.128750, 0.15},
	}
	for _, tt := range testData {
		width, height := fonts.TextContainerSize(tt.text, tt.scale)
		if !testhelper.Float32ApproxEqual(width, tt.width, 0.00001) {
			t.Errorf("Invalid text width for '%s'. Instead of '%f', we have '%f'.", tt.text, tt.width, width)
		}
		if !testhelper.Float32ApproxEqual(height, tt.height, 0.00001) {
			t.Errorf("Invalid text height for '%s'. Instead of '%f', we have '%f'.", tt.text, tt.height, height)
		}
	}
}
