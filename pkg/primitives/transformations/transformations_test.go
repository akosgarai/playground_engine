package transformations

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestMouseCoordinates(t *testing.T) {
	testData := []struct {
		x      float64
		y      float64
		width  float64
		height float64
		rX     float64
		rY     float64
	}{
		{0, 0, 600, 600, -1, 1},
		{600, 600, 600, 600, 1, -1},
		{300, 300, 600, 600, 0, 0},
	}
	for _, tt := range testData {
		x, y := MouseCoordinates(tt.x, tt.y, tt.width, tt.height)
		if x != tt.rX {
			t.Errorf("Invalid x Component instead of %f, got %f", tt.rX, x)
		}
		if y != tt.rY {
			t.Errorf("Invalid y Component instead of %f, got %f", tt.rY, y)
		}
	}
}
func TestVec3ToString(t *testing.T) {
	testData := []struct {
		v mgl32.Vec3
		s string
	}{
		{mgl32.Vec3{0, 0, 0}, "X : 0.0000000000, Y : 0.0000000000, Z : 0.0000000000"},
		{mgl32.Vec3{1, 1, 1}, "X : 1.0000000000, Y : 1.0000000000, Z : 1.0000000000"},
		{mgl32.Vec3{0.5, 0.5, 0.5}, "X : 0.5000000000, Y : 0.5000000000, Z : 0.5000000000"},
	}
	for _, tt := range testData {
		str := Vec3ToString(tt.v)
		if str != tt.s {
			t.Errorf("Invalid string representation. Instead of '%s', got '%s'", tt.s, str)
		}
	}
}
func TestFloat64ToString(t *testing.T) {
	testData := []struct {
		v float64
		s string
	}{
		{0, "0.0000000000"},
		{1, "1.0000000000"},
		{0.5, "0.5000000000"},
	}
	for _, tt := range testData {
		str := Float64ToString(tt.v)
		if str != tt.s {
			t.Errorf("Invalid string representation. Instead of '%s', got '%s'", tt.s, str)
		}
	}
}
func TestFloat32ToString(t *testing.T) {
	testData := []struct {
		v float32
		s string
	}{
		{0, "0.0000000000"},
		{1, "1.0000000000"},
		{0.5, "0.5000000000"},
	}
	for _, tt := range testData {
		str := Float32ToString(tt.v)
		if str != tt.s {
			t.Errorf("Invalid string representation. Instead of '%s', got '%s'", tt.s, str)
		}
	}
}
func TestIntegerToString(t *testing.T) {
	testData := []struct {
		v int
		s string
	}{
		{0, "0"},
		{1, "1"},
		{5, "5"},
	}
	for _, tt := range testData {
		str := IntegerToString(tt.v)
		if str != tt.s {
			t.Errorf("Invalid string representation. Instead of '%s', got '%s'", tt.s, str)
		}
	}
}
