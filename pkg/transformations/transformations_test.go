package transformations

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/testhelper"

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
func TestVectorToString(t *testing.T) {
	testData := []struct {
		v mgl32.Vec3
		s [3]string
	}{
		{mgl32.Vec3{0, 0, 0}, [3]string{"0", "0", "0"}},
		{mgl32.Vec3{0.5, 0.5, 0.5}, [3]string{"0.5", "0.5", "0.5"}},
		{mgl32.Vec3{1.50005, 0.5, 0.5}, [3]string{"1.50005", "0.5", "0.5"}},
	}
	for _, tt := range testData {
		str := VectorToString(tt.v)
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
func TestFloat64ToStringExact(t *testing.T) {
	testData := []struct {
		v float64
		s string
	}{
		{0, "0"},
		{1, "1"},
		{0.5, "0.5"},
	}
	for _, tt := range testData {
		str := Float64ToStringExact(tt.v)
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
func TestFloat32ToStringExact(t *testing.T) {
	testData := []struct {
		v float32
		s string
	}{
		{0, "0"},
		{1, "1"},
		{0.5, "0.5"},
	}
	for _, tt := range testData {
		str := Float32ToStringExact(tt.v)
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
func TestInteger64ToString(t *testing.T) {
	testData := []struct {
		v int64
		s string
	}{
		{0, "0"},
		{1, "1"},
		{5, "5"},
		{5010101010101010111, "5010101010101010111"},
	}
	for _, tt := range testData {
		str := Integer64ToString(tt.v)
		if str != tt.s {
			t.Errorf("Invalid string representation. Instead of '%s', got '%s'", tt.s, str)
		}
	}
}
func TestFloat32Abs(t *testing.T) {
	testData := []struct {
		input float32
		abs   float32
	}{
		{0.0, 0.0},
		{1.0, 1.0},
		{5.0, 5.0},
		{0.1, 0.1},
		{-0.1, 0.1},
		{-1.0, 1.0},
	}
	for _, tt := range testData {
		abs := Float32Abs(tt.input)
		if abs != tt.abs {
			t.Errorf("Invalid f32 abs. Instead of '%f', got '%f'", tt.abs, abs)
		}
	}
}
func TestExtractAngles(t *testing.T) {
	testData := []struct {
		x float32
		y float32
		z float32
	}{
		{0.0, 0.0, 0.0},
		{1.0, 1.0, 1.0},
		{45.0, 0.0, 45.0},
		{-45.0, 0.0, 45.0},
		{30.0, 20.0, 45.0},
		{-90.0, 0.0, 0.0},
		{90.0, 0.0, 0.0},
	}
	for _, tt := range testData {
		rotationMatrix := mgl32.HomogRotate3DY(mgl32.DegToRad(tt.y)).Mul4(mgl32.HomogRotate3DX(mgl32.DegToRad(tt.x))).Mul4(mgl32.HomogRotate3DZ(mgl32.DegToRad(tt.z)))
		x, y, z := ExtractAngles(rotationMatrix)
		if !testhelper.Float32ApproxEqual(x, tt.x, 0.0001) {
			t.Errorf("Invalid x component. Instead of '%f', it is '%f'.", tt.x, x)
		}
		if !testhelper.Float32ApproxEqual(y, tt.y, 0.0001) {
			t.Errorf("Invalid y component. Instead of '%f', it is '%f'.", tt.y, y)
		}
		if !testhelper.Float32ApproxEqual(z, tt.z, 0.0001) {
			t.Errorf("Invalid z component. Instead of '%f', it is '%f'.", tt.z, z)
		}
	}
}
