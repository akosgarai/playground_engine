package model

import (
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNewRoomBuilder(t *testing.T) {
	b := NewRoomBuilder()
	expectedPos := mgl32.Vec3{0, 0, 0}
	expectedDir := mgl32.Vec3{0, 1, 0}
	if b.position != expectedPos {
		t.Errorf("Invalid position. '%#v'.", b.position)
	}
	if b.worldUp != expectedDir {
		t.Errorf("Invalid worldUp. '%#v'.", b.worldUp)
	}
	if b.width != float32(1.0) {
		t.Errorf("Invalid width. '%f'.", b.width)
	}
	if b.height != float32(1.0) {
		t.Errorf("Invalid height. '%f'.", b.height)
	}
	if b.length != float32(1.0) {
		t.Errorf("Invalid length. '%f'.", b.length)
	}
	if b.wallWidth != float32(0.005) {
		t.Errorf("Invalid wallWidth. '%f'.", b.wallWidth)
	}
	if b.doorHeight != float32(0.6) {
		t.Errorf("Invalid door height. '%f'.", b.doorHeight)
	}
	if b.doorWidth != float32(0.4) {
		t.Errorf("Invalid door width. '%f'.", b.doorWidth)
	}
}
func TestRoomBuilderSetPosition(t *testing.T) {
	b := NewRoomBuilder()
	expectedPos := mgl32.Vec3{1, 1, 1}
	b.SetPosition(expectedPos)
	if b.position != expectedPos {
		t.Errorf("Invalid position. '%#v'.", b.position)
	}
}
func TestRoomBuilderSetWorldUpDirection(t *testing.T) {
	b := NewRoomBuilder()
	expectedDir := mgl32.Vec3{1, 0, 0}
	b.SetWorldUpDirection(expectedDir)
	if b.worldUp != expectedDir {
		t.Errorf("Invalid worldUp. '%#v'.", b.worldUp)
	}
	b.SetWorldUpDirection(expectedDir.Mul(5))
	if b.worldUp != expectedDir {
		t.Errorf("Invalid worldUp. '%#v'.", b.worldUp)
	}
}
func TestRoomBuilderSetSize(t *testing.T) {
	b := NewRoomBuilder()
	expectedW := float32(2.0)
	expectedH := float32(2.1)
	expectedL := float32(2.2)
	b.SetSize(expectedW, expectedH, expectedL)
	if b.width != expectedW {
		t.Errorf("Invalid width. '%f'.", b.width)
	}
	if b.height != expectedH {
		t.Errorf("Invalid height. '%f'.", b.height)
	}
	if b.length != expectedL {
		t.Errorf("Invalid length. '%f'.", b.length)
	}
}
func TestRoomBuilderSetWallWidth(t *testing.T) {
	b := NewRoomBuilder()
	expectedW := float32(2.0)
	b.SetWallWidth(expectedW)
	if b.wallWidth != expectedW {
		t.Errorf("Invalid wallWidth. '%f'.", b.wallWidth)
	}
}
func TestRoomBuilderSetWrapper(t *testing.T) {
	b := NewRoomBuilder()
	b.SetWrapper(wrapperMock)
	if !reflect.DeepEqual(b.wrapper, wrapperMock) {
		t.Error("Invalid wrapper.")
	}
}
func TestRoomBuilderSetDoorSize(t *testing.T) {
	b := NewRoomBuilder()
	expectedW := float32(2.0)
	expectedH := float32(2.1)
	b.SetDoorSize(expectedW, expectedH)
	if b.doorWidth != expectedW {
		t.Errorf("Invalid door width. '%f'.", b.doorWidth)
	}
	if b.doorHeight != expectedH {
		t.Errorf("Invalid door height. '%f'.", b.doorHeight)
	}
}
func TestRoomBuilderBuild(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("It shouldn't have panic. '%#v'.", r)
			}
		}()
		b := NewRoomBuilder()
		b.SetWrapper(wrapperMock)
		_ = b.Build()
	}()
}
func TestRoomBuilderBuildWithoutWrapper(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("It should have paniced.")
			}
		}()
		b := NewRoomBuilder()
		_ = b.Build()
	}()
}

func CheckDefaultRoomOptions(room *Room, t *testing.T) {
	doorPosition := mgl32.Vec3{-0.4975, 0.7, 0.6975}
	if room.GetDoor().GetPosition() != doorPosition {
		t.Errorf("Invalid door position. Instead of '%v', we have '%v'.", doorPosition, room.GetDoor().GetPosition())
	}
	if room.doorState != _DOOR_OPENED {
		t.Errorf("Invalid initial door state. Instead of '%d', we have '%d'.", _DOOR_OPENED, room.doorState)
	}
	room.PushDoorState()
	if room.doorState != _DOOR_CLOSING {
		t.Errorf("Invalid next door state. Instead of '%d', we have '%d'.", _DOOR_CLOSING, room.doorState)
	}
	room.PushDoorState()
	if room.doorState != _DOOR_CLOSING {
		t.Errorf("Invalid door state. Instead of '%d', we have '%d'.", _DOOR_CLOSING, room.doorState)
	}

	room.doorState = _DOOR_CLOSED
	room.PushDoorState()
	if room.doorState != _DOOR_OPENING {
		t.Errorf("Invalid next door state. Instead of '%d', we have '%d'.", _DOOR_OPENING, room.doorState)
	}
	room.PushDoorState()
	if room.doorState != _DOOR_OPENING {
		t.Errorf("Invalid next door state. Instead of '%d', we have '%d'.", _DOOR_OPENING, room.doorState)
	}
	if room.currentAnimationTime != 0.0 {
		t.Errorf("Invalid initial animation time. Instead of '0.0', it is '%f'.", room.currentAnimationTime)
	}
	room.animateDoor(100)
	if room.currentAnimationTime != 100.0 {
		t.Errorf("Invalid animation time. Instead of '100.0', it is '%f'.", room.currentAnimationTime)
	}
	room.animateDoor(950)
	if room.doorState != _DOOR_OPENED {
		t.Errorf("Invalid next door state. Instead of '%d', we have '%d'.", _DOOR_OPENED, room.doorState)
	}
	room.PushDoorState()
	room.animateDoor(100)
	if room.currentAnimationTime != 100.0 {
		t.Errorf("Invalid animation time. Instead of '100.0', it is '%f'.", room.currentAnimationTime)
	}
	room.animateDoor(950)
	if room.doorState != _DOOR_CLOSED {
		t.Errorf("Invalid next door state. Instead of '%d', we have '%d'.", _DOOR_CLOSED, room.doorState)
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Update shouldn't have panic.")
			}
		}()
		room.Update(100)
	}()
}
func TestMaterialRoom(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}

	room := NewMaterialRoom(position, wrapperMock)

	CheckDefaultRoomOptions(room, t)
}
func TestTexturedRoom(t *testing.T) {
	position := mgl32.Vec3{0.0, 0.0, 0.0}

	room := NewTextureRoom(position, wrapperMock)

	CheckDefaultRoomOptions(room, t)
}
