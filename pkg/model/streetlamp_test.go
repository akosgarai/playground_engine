package model

import (
	"testing"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"
)

func TestMaterialStreetLamp(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	bulbPosition := mgl32.Vec3{0.9349, 3.0, 2.98}
	polePosition := mgl32.Vec3{0.0, 3.0, 0.0}
	topPosition := mgl32.Vec3{1.1, 3.0, 6.4}
	scale := float32(6.0)

	lamp := NewMaterialStreetLamp(position, scale, wrapperMock)

	if lamp.GetPolePosition() != polePosition {
		t.Errorf("Invalid pole position. Instead of '%v', we have '%v'.", polePosition, lamp.GetPolePosition())
	}
	if !lamp.GetTopPosition().ApproxEqualThreshold(topPosition, 0.0001) {
		t.Errorf("Invalid top position. Instead of '%v', we have '%v'.", topPosition, lamp.GetTopPosition())
	}
	if !lamp.GetBulbPosition().ApproxEqualThreshold(bulbPosition, 0.0001) {
		t.Errorf("Invalid bulb position. Instead of '%v', we have '%v'.", bulbPosition, lamp.GetBulbPosition())
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Update shouldn't have panic.")
			}
		}()
		lamp.Update(10)
	}()
	testData := []struct {
		position  [3]float32
		radius    float32
		intersect bool
		msg       string
	}{
		{[3]float32{-0.6, 3, 0}, 0.5, true, "Should intersect at x."},
		{[3]float32{-2, -2, -2}, 1.5, false, "Shouldn't intersect."},
	}

	for _, tt := range testData {
		base := coldet.NewBoundingSphere(tt.position, tt.radius)
		result := lamp.CollideTestWithSphere(base)
		if result != tt.intersect {
			t.Errorf("%s expected: '%v', result: '%v'.", tt.msg, tt.intersect, result)
		}
	}
}
func TestTexturedStreetLamp(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}
	bulbPosition := mgl32.Vec3{0.9, 3.035, 3.0}
	polePosition := mgl32.Vec3{0.0, 3.0, 0.0}
	topPosition := mgl32.Vec3{3.55, 3.55, 3.0}
	scale := float32(6.0)

	lamp := NewTexturedStreetLamp(position, scale, wrapperMock)

	if lamp.GetPolePosition() != polePosition {
		t.Errorf("Invalid pole position. Instead of '%v', we have '%v'.", polePosition, lamp.GetPolePosition())
	}
	if !lamp.GetTopPosition().ApproxEqualThreshold(topPosition, 0.0001) {
		t.Errorf("Invalid top position. Instead of '%v', we have '%v'.", topPosition, lamp.GetTopPosition())
	}
	if !lamp.GetBulbPosition().ApproxEqualThreshold(bulbPosition, 0.0001) {
		t.Errorf("Invalid bulb position. Instead of '%v', we have '%v'.", bulbPosition, lamp.GetBulbPosition())
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Update shouldn't have panic.")
			}
		}()
		lamp.Update(10)
	}()
}
