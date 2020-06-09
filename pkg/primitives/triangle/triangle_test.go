package triangle

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNew(t *testing.T) {
	triangle := New(60, 60, 60)
	expectedNormal := mgl32.Vec3{0, -1, 0}

	if triangle.Normal != expectedNormal {
		t.Error("Invalid normal vector")
	}
}
func TestSortAngles(t *testing.T) {
	testData := []struct {
		Orig    [3]float32
		Ordered [3]float32
	}{
		{[3]float32{1, 2, 3}, [3]float32{3, 2, 1}},
		{[3]float32{1, 3, 2}, [3]float32{3, 2, 1}},
		{[3]float32{3, 1, 2}, [3]float32{3, 2, 1}},
		{[3]float32{3, 2, 1}, [3]float32{3, 2, 1}},
		{[3]float32{2, 3, 1}, [3]float32{3, 2, 1}},
		{[3]float32{2, 1, 3}, [3]float32{3, 2, 1}},
		{[3]float32{2, 2, 2}, [3]float32{2, 2, 2}},
	}
	for _, tt := range testData {
		ordered := sortAngles(tt.Orig[0], tt.Orig[1], tt.Orig[2])
		if ordered != tt.Ordered {
			t.Error("Invalid order")
			t.Log(ordered)
			t.Log(tt.Ordered)
			t.Log(tt.Orig)
		}
	}
}

func TestColoredMeshInput(t *testing.T) {
	triangle := New(60, 60, 60)
	vert, ind, _ := triangle.ColoredMeshInput([]mgl32.Vec3{mgl32.Vec3{1, 1, 1}})
	if ind[0] != 0 || ind[1] != 1 || ind[2] != 2 {
		t.Error("Invalid indices")
		t.Log(vert)
		t.Log(ind)
	}
}
