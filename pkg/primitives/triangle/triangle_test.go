package triangle

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNew(t *testing.T) {
	triangle := New(60, 60, 60)
	expectedNormal := mgl32.Vec3{0, -1, 0}
	expectedWidth := float32(1.0)
	expectedHeight := float32(0.0)
	expectedLength := float32(0.86602545) // sqrt(1^2 - 0.5^2)

	// normal
	if triangle.Normal != expectedNormal {
		t.Error("Invalid normal vector")
	}
	// bounding object
	boParams := triangle.BO.Params()
	if boParams["height"] != expectedHeight {
		t.Errorf("Invalid BO height. Instead of '%f', we have '%f'.\n", expectedHeight, boParams["height"])
	}
	if boParams["length"] != expectedLength {
		t.Errorf("Invalid BO length. Instead of '%f', we have '%f'.\n", expectedLength, boParams["length"])
	}
	if boParams["width"] != expectedWidth {
		t.Errorf("Invalid BO width. Instead of '%f', we have '%f'.\n", expectedWidth, boParams["width"])
	}
	// position
	leftPoint := mgl32.Vec3{-0.5, 0.0, 0.0}
	rightPoint := mgl32.Vec3{0.5, 0.0, 0.0}
	topPoint := mgl32.Vec3{0.0, 0.0, expectedLength}
	if !leftPoint.ApproxEqual(triangle.Points[0]) {
		t.Errorf("Invalid left point. Insead of '%v', we have '%v'.\n", leftPoint, triangle.Points[0])
	}
	if !topPoint.ApproxEqualThreshold(triangle.Points[1], 0.001) {
		t.Errorf("Invalid top point. Insead of '%v', we have '%v'.\n", topPoint, triangle.Points[1])
	}
	if !rightPoint.ApproxEqual(triangle.Points[2]) {
		t.Errorf("Invalid top right. Insead of '%v', we have '%v'.\n", rightPoint, triangle.Points[2])
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
