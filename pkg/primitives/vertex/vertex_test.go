package vertex

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestAdd(t *testing.T) {
	vert := Vertex{
		Position: mgl32.Vec3{0, 0, 0},
	}
	var vertices Vertices
	vertices.Add(vert)
	if len(vertices) != 1 {
		t.Error("Invalid length")
	}
}
func TestGet(t *testing.T) {
	vert := Vertex{
		Position:  mgl32.Vec3{0, 0, 0},
		Normal:    mgl32.Vec3{1, 0, 0},
		TexCoords: mgl32.Vec2{1, 1},
		Color:     mgl32.Vec3{0, 1, 1},
		PointSize: float32(11),
	}
	var vertices Vertices
	vertices.Add(vert)
	expectedTextures := []float32{0, 0, 0, 1, 0, 0, 1, 1}
	expectedMaterial := []float32{0, 0, 0, 1, 0, 0}
	expectedPoint := []float32{0, 0, 0, 0, 1, 1, 11}
	expectedColored := []float32{0, 0, 0, 0, 1, 1}
	expectedTextureColor := []float32{0, 0, 0, 0, 1, 1, 1, 1}
	if !reflect.DeepEqual(vertices.Get(POSITION_NORMAL_TEXCOORD), expectedTextures) {
		t.Error("Invalid texture vao")
	}
	if !reflect.DeepEqual(vertices.Get(POSITION_NORMAL), expectedMaterial) {
		t.Error("Invalid material vao")
	}
	if !reflect.DeepEqual(vertices.Get(POSITION_COLOR_SIZE), expectedPoint) {
		t.Error("Invalid point vao")
	}
	if !reflect.DeepEqual(vertices.Get(POSITION_COLOR), expectedColored) {
		t.Error("Invalid point vao")
	}
	if !reflect.DeepEqual(vertices.Get(POSITION_COLOR_TEXCOORD), expectedTextureColor) {
		t.Error("Invalid point vao")
	}
}
