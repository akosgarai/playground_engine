package camera

import (
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	DefaultCameraPosition = mgl32.Vec3{0, 0, 0}
	WorldUp               = mgl32.Vec3{0, 1, 0}
	DefaultYaw            = float32(0)
	DefaultPitch          = float32(0)
	DefaultFront          = mgl32.Vec3{1, 0, 0}
	DefaultUp             = mgl32.Vec3{0, -1, 0}
	DefaultRight          = mgl32.Vec3{0, 0, -1}

	DefaultFov      = float32(45)
	DefaultAspRatio = float32(1.0)
	DefaultNear     = float32(0.1)
	DefaultFar      = float32(100)
)

func TestNewCamera(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)

	if cam.pitch != DefaultPitch {
		t.Errorf("Invalid pitch instead of '%f', we have '%f'", DefaultPitch, cam.pitch)
	}
	if cam.yaw != DefaultYaw {
		t.Errorf("Invalid yaw instead of '%f', we have '%f'", DefaultYaw, cam.yaw)
	}
	if cam.cameraPosition != DefaultCameraPosition {
		t.Errorf("Invalid position")
	}
	if cam.GetPosition() != DefaultCameraPosition {
		t.Errorf("Invalid position")
	}
	if cam.worldUp != WorldUp {
		t.Errorf("Invalid worldUp")
	}
	if cam.cameraFrontDirection != DefaultFront {
		t.Error("Invalid front direction")
	}
	if cam.cameraUpDirection != DefaultUp {
		t.Error("Invalid up direction")
	}
	if cam.cameraRightDirection != DefaultRight {
		t.Error("Invalid right direction")
	}
}
func TestLog(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	log := cam.Log()

	if len(log) < 10 {
		t.Errorf("Log too short: '%s'", log)
	}
}
func TestWalk(t *testing.T) {
	amountToMove := float32(2)

	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.Walk(amountToMove)
	// the front is +X, amount is 2 -> position {2,0,0}
	expectedPos := mgl32.Vec3{2, 0, 0}
	if cam.cameraPosition != expectedPos {
		t.Error("Invalid movement")
	}
	if cam.GetPosition() != expectedPos {
		t.Errorf("Invalid position")
	}
	if cam.cameraFrontDirection != DefaultFront {
		t.Error("Invalid front direction")
	}
	if cam.cameraUpDirection != DefaultUp {
		t.Error("Invalid up direction")
	}
	if cam.cameraRightDirection != DefaultRight {
		t.Error("Invalid right direction")
	}
}
func TestStrafe(t *testing.T) {
	amountToMove := float32(2)

	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.Strafe(amountToMove)
	// the right is -Z, amount is 2 -> position {0,0,-2}
	expectedPos := mgl32.Vec3{0, 0, -2}
	if cam.cameraPosition != expectedPos {
		t.Error("Invalid movement")
	}
	if cam.GetPosition() != expectedPos {
		t.Errorf("Invalid position")
	}
	if cam.cameraFrontDirection != DefaultFront {
		t.Error("Invalid front direction")
	}
	if cam.cameraUpDirection != DefaultUp {
		t.Error("Invalid up direction")
	}
	if cam.cameraRightDirection != DefaultRight {
		t.Error("Invalid right direction")
	}
}
func TestLift(t *testing.T) {
	amountToMove := float32(2)

	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.Lift(amountToMove)
	// the up is -Y, amount is 2 -> position {0,-2,0}
	expectedPos := mgl32.Vec3{0, -2, 0}
	if cam.cameraPosition != expectedPos {
		t.Error("Invalid movement")
	}
	if cam.GetPosition() != expectedPos {
		t.Errorf("Invalid position")
	}
	if cam.cameraFrontDirection != DefaultFront {
		t.Error("Invalid front direction")
	}
	if cam.cameraUpDirection != DefaultUp {
		t.Error("Invalid up direction")
	}
	if cam.cameraRightDirection != DefaultRight {
		t.Error("Invalid right direction")
	}
}
func TestSetupProjection(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.SetupProjection(DefaultFov, DefaultAspRatio, DefaultNear, DefaultFar)
	if cam.projectionOptions.fov != DefaultFov {
		t.Errorf("Invalid fov instead of '%f', we have '%f'", DefaultFov, cam.projectionOptions.fov)
	}
	if cam.projectionOptions.aspectRatio != DefaultAspRatio {
		t.Errorf("Invalid aspectRatio instead of '%f', we have '%f'", DefaultAspRatio, cam.projectionOptions.aspectRatio)
	}
	if cam.projectionOptions.near != DefaultNear {
		t.Errorf("Invalid near instead of '%f', we have '%f'", DefaultNear, cam.projectionOptions.near)
	}
	if cam.projectionOptions.far != DefaultFar {
		t.Errorf("Invalid far instead of '%f', we have '%f'", DefaultFar, cam.projectionOptions.far)
	}
}
func TestGetProjectionMatrix(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.SetupProjection(DefaultFov, DefaultAspRatio, DefaultNear, DefaultFar)
	cam.GetProjectionMatrix()
}
func TestGetViewMatrix(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.SetupProjection(DefaultFov, DefaultAspRatio, DefaultNear, DefaultFar)
	cam.GetViewMatrix()
}
func TestUpdateDirection(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.SetupProjection(DefaultFov, DefaultAspRatio, DefaultNear, DefaultFar)
	cam.UpdateDirection(0, 0)
	if cam.cameraPosition != DefaultCameraPosition {
		t.Error("Invalid movement")
	}
	if cam.GetPosition() != DefaultCameraPosition {
		t.Errorf("Invalid position")
	}
	if cam.cameraFrontDirection != DefaultFront {
		t.Error("Invalid front direction")
	}
	if cam.cameraUpDirection != DefaultUp {
		t.Error("Invalid up direction")
	}
	if cam.cameraRightDirection != DefaultRight {
		t.Error("Invalid right direction")
	}
}
func TestGetBoundingObject(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	bo := cam.GetBoundingObject()
	if bo.X() != DefaultCameraPosition.X() {
		t.Errorf("Invalid X coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.X(), bo.X())
	}
	if bo.Y() != DefaultCameraPosition.Y() {
		t.Errorf("Invalid Y coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.Y(), bo.Y())
	}
	if bo.Z() != DefaultCameraPosition.Z() {
		t.Errorf("Invalid Z coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.Z(), bo.Z())
	}
	if bo.Radius() != 0.1 {
		t.Errorf("Invalid radius for the bounding sphere. Instead of 0.1, we have '%f'.\n", bo.Radius())
	}

}
func TestBoundingObjectAfterWalk(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	bo := cam.BoundingObjectAfterWalk(1)
	if bo.X() != DefaultCameraPosition.X()+1 {
		t.Errorf("Invalid X coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.X()+1, bo.X())
	}
	if bo.Y() != DefaultCameraPosition.Y() {
		t.Errorf("Invalid Y coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.Y(), bo.Y())
	}
	if bo.Z() != DefaultCameraPosition.Z() {
		t.Errorf("Invalid Z coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.Z(), bo.Z())
	}
	if bo.Radius() != 0.1 {
		t.Errorf("Invalid radius for the bounding sphere. Instead of 0.1, we have '%f'.\n", bo.Radius())
	}

}
func TestBoundingObjectAfterStrafe(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	bo := cam.BoundingObjectAfterStrafe(-1)
	if bo.X() != DefaultCameraPosition.X() {
		t.Errorf("Invalid X coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.X(), bo.X())
	}
	if bo.Y() != DefaultCameraPosition.Y() {
		t.Errorf("Invalid Y coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.Y(), bo.Y())
	}
	if bo.Z() != DefaultCameraPosition.Z()+1 {
		t.Errorf("Invalid Z coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.Z()+1, bo.Z())
	}
	if bo.Radius() != 0.1 {
		t.Errorf("Invalid radius for the bounding sphere. Instead of 0.1, we have '%f'.\n", bo.Radius())
	}

}
func TestBoundingObjectAfterLift(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	bo := cam.BoundingObjectAfterLift(-1)
	if bo.X() != DefaultCameraPosition.X() {
		t.Errorf("Invalid X coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.X(), bo.X())
	}
	if bo.Y() != DefaultCameraPosition.Y()+1 {
		t.Errorf("Invalid Y coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.Y()+1, bo.Y())
	}
	if bo.Z() != DefaultCameraPosition.Z() {
		t.Errorf("Invalid Z coordinate for the bounding sphere. Instead of '%f', we have '%f'.\n", DefaultCameraPosition.Z(), bo.Z())
	}
	if bo.Radius() != 0.1 {
		t.Errorf("Invalid radius for the bounding sphere. Instead of 0.1, we have '%f'.\n", bo.Radius())
	}
}

func TestSetVelocity(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.SetVelocity(3)
	if cam.velocity != 3 {
		t.Errorf("Invalid velocity. Instead of '3', we have '%f'.\n", cam.velocity)
	}
}
func TestGetVelocity(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.SetVelocity(3)
	if cam.GetVelocity() != 3 {
		t.Errorf("Invalid velocity. Instead of '3', we have '%f'.\n", cam.GetVelocity())
	}
}
func TestSetRotationStep(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.SetRotationStep(3)
	if cam.rotationStep != 3 {
		t.Errorf("Invalid rotationStep. Instead of '3', we have '%f'.\n", cam.velocity)
	}
}
func TestGetRotationStep(t *testing.T) {
	cam := NewCamera(DefaultCameraPosition, WorldUp, DefaultYaw, DefaultPitch)
	cam.SetRotationStep(3)
	if cam.GetRotationStep() != 3 {
		t.Errorf("Invalid rotationStep. Instead of '3', we have '%f'.\n", cam.GetRotationStep())
	}
}
