package model

import (
	"testing"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"
)

func TestBug(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	bottomPosition := mgl32.Vec3{-1.5, 0.0, 0.0}
	eyeBase := float32(0.6350853)
	eye1Position := mgl32.Vec3{eyeBase, eyeBase, eyeBase}
	eye2Position := mgl32.Vec3{eyeBase, eyeBase, -eyeBase}
	scale := mgl32.Vec3{1.0, 1.0, 1.0}

	bug := NewBug(position, scale, wrapperMock)

	if bug.GetBodyPosition() != position {
		t.Errorf("Invalid body position. Instead of '%v', we have '%v'.", position, bug.GetBodyPosition())
	}
	if bug.GetBottomPosition() != bottomPosition {
		t.Errorf("Invalid bottom position. Instead of '%v', we have '%v'.", bottomPosition, bug.GetBottomPosition())
	}
	if bug.GetEye1Position() != eye1Position {
		t.Errorf("Invalid eye1 position. Instead of '%v', we have '%v'.", eye1Position, bug.GetEye1Position())
	}
	if bug.GetEye2Position() != eye2Position {
		t.Errorf("Invalid eye2 position. Instead of '%v', we have '%v'.", eye2Position, bug.GetEye2Position())
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Update shouldn't have panic.")
			}
		}()
		bug.Update(10)
	}()
	testData := []struct {
		position  [3]float32
		radius    float32
		intersect bool
		msg       string
	}{
		{[3]float32{0, 0, 0}, 0.5, true, "Should intersect at x=-0.5."},
		{[3]float32{-1.5, 1.3, 0.0}, 1.0, true, "Should intersect at y=1."},
		{[3]float32{-2, -2, -2}, 1.5, false, "Shouldn't intersect."},
	}

	for _, tt := range testData {
		base := coldet.NewBoundingSphere(tt.position, tt.radius)
		result := bug.CollideTestWithSphere(base)
		if result != tt.intersect {
			t.Errorf("%s expected: '%v', result: '%v'.", tt.msg, tt.intersect, result)
		}
	}
}
