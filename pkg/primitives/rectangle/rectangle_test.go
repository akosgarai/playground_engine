package rectangle

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewSquare(t *testing.T) {
	square := NewSquare()

	ExpectedPoints := [4]mgl32.Vec3{
		mgl32.Vec3{-0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, 0.5},
		mgl32.Vec3{-0.5, 0, 0.5},
	}
	ExpectedNormal := mgl32.Vec3{0, -1, 0}
	if square.Points != ExpectedPoints {
		t.Errorf("Invalid points. Instead of '%v', we have '%v'.\n", ExpectedPoints, square.Points)
	}
	if square.Normal != ExpectedNormal {
		t.Errorf("Invalid normal vector. Instead of '%v', we have '%v'.\n", ExpectedNormal, square.Normal)
	}
	v1 := square.Points[1].Sub(square.Points[0])
	v2 := square.Points[3].Sub(square.Points[0])
	calculatedNormal := v1.Cross(v2).Normalize()
	if calculatedNormal != square.Normal {
		t.Errorf("Invalid normal vs calculated normal. Instead of '%v', we have '%v'.\n", calculatedNormal, square.Normal)
	}
}
func TestNewOneAsScale(t *testing.T) {
	width := float32(1)
	height := float32(1)
	square := New(width, height)

	ExpectedPoints := [4]mgl32.Vec3{
		mgl32.Vec3{-0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, 0.5},
		mgl32.Vec3{-0.5, 0, 0.5},
	}
	ExpectedNormal := mgl32.Vec3{0, -1, 0}
	if square.Points != ExpectedPoints {
		t.Errorf("Invalid points. Instead of '%v', we have '%v'.\n", ExpectedPoints, square.Points)
	}
	if square.Normal != ExpectedNormal {
		t.Errorf("Invalid normal vector. Instead of '%v', we have '%v'.\n", ExpectedNormal, square.Normal)
	}
	v1 := square.Points[1].Sub(square.Points[0])
	v2 := square.Points[3].Sub(square.Points[0])
	calculatedNormal := v1.Cross(v2).Normalize()
	if calculatedNormal != square.Normal {
		t.Errorf("Invalid normal vs calculated normal. Instead of '%v', we have '%v'.\n", calculatedNormal, square.Normal)
	}
}
func TestNewLowScale(t *testing.T) {
	width := float32(2)
	height := float32(1)
	square := New(width, height)

	ExpectedPoints := [4]mgl32.Vec3{
		mgl32.Vec3{-0.5, 0, -0.25},
		mgl32.Vec3{0.5, 0, -0.25},
		mgl32.Vec3{0.5, 0, 0.25},
		mgl32.Vec3{-0.5, 0, 0.25},
	}
	ExpectedNormal := mgl32.Vec3{0, -1, 0}
	if square.Points != ExpectedPoints {
		t.Errorf("Invalid points. Instead of '%v', we have '%v'.\n", ExpectedPoints, square.Points)
	}
	if square.Normal != ExpectedNormal {
		t.Errorf("Invalid normal vector. Instead of '%v', we have '%v'.\n", ExpectedNormal, square.Normal)
	}
	v1 := square.Points[1].Sub(square.Points[0])
	v2 := square.Points[3].Sub(square.Points[0])
	calculatedNormal := v1.Cross(v2).Normalize()
	if calculatedNormal != square.Normal {
		t.Errorf("Invalid normal vs calculated normal. Instead of '%v', we have '%v'.\n", calculatedNormal, square.Normal)
	}
}
func TestNewHighScale(t *testing.T) {
	width := float32(1)
	height := float32(2)
	square := New(width, height)

	ExpectedPoints := [4]mgl32.Vec3{
		mgl32.Vec3{-0.25, 0, -0.5},
		mgl32.Vec3{0.25, 0, -0.5},
		mgl32.Vec3{0.25, 0, 0.5},
		mgl32.Vec3{-0.25, 0, 0.5},
	}
	ExpectedNormal := mgl32.Vec3{0, -1, 0}
	if square.Points != ExpectedPoints {
		t.Errorf("Invalid points. Instead of '%v', we have '%v'.\n", ExpectedPoints, square.Points)
	}
	if square.Normal != ExpectedNormal {
		t.Errorf("Invalid normal vector. Instead of '%v', we have '%v'.\n", ExpectedNormal, square.Normal)
	}
	v1 := square.Points[1].Sub(square.Points[0])
	v2 := square.Points[3].Sub(square.Points[0])
	calculatedNormal := v1.Cross(v2).Normalize()
	if calculatedNormal != square.Normal {
		t.Errorf("Invalid normal vs calculated normal. Instead of '%v', we have '%v'.\n", calculatedNormal, square.Normal)
	}
}
func TestMeshInput(t *testing.T) {
	square := NewSquare()
	vertices, indices, _ := square.MeshInput()
	expectedIndices := []uint32{0, 1, 2, 0, 2, 3}
	if !reflect.DeepEqual(expectedIndices, indices) {
		t.Errorf("Invalid indices. Instead of '%v', we have '%v'.\n", expectedIndices, indices)
	}
	if len(vertices) != 4 {
		t.Errorf("Invalud number of vertices. Instead of '4', we have '%d' as '%v'.\n", len(vertices), vertices)
	}
}
func TestColoredMeshInput(t *testing.T) {
	square := NewSquare()
	color := []mgl32.Vec3{mgl32.Vec3{1, 1, 1}}
	vertices, indices, _ := square.ColoredMeshInput(color)
	expectedIndices := []uint32{0, 1, 2, 0, 2, 3}
	if !reflect.DeepEqual(expectedIndices, indices) {
		t.Errorf("Invalid indices. Instead of '%v', we have '%v'.\n", expectedIndices, indices)
	}
	if len(vertices) != 4 {
		t.Errorf("Invalud number of vertices. Instead of '4', we have '%d' as '%v'.\n", len(vertices), vertices)
	}
}
func TestColoredTexturesMeshInput(t *testing.T) {
	square := NewSquare()
	color := []mgl32.Vec3{mgl32.Vec3{1, 1, 1}}
	vertices, indices, _ := square.TexturedColoredMeshInput(color)
	expectedIndices := []uint32{0, 1, 2, 0, 2, 3}
	if !reflect.DeepEqual(expectedIndices, indices) {
		t.Errorf("Invalid indices. Instead of '%v', we have '%v'.\n", expectedIndices, indices)
	}
	if len(vertices) != 4 {
		t.Errorf("Invalud number of vertices. Instead of '4', we have '%d' as '%v'.\n", len(vertices), vertices)
	}
}
