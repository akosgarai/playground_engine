package window

import (
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Test title"
)

func TestInitGlfw(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("Init shouldn't be failed.")
			}
		}()
		InitGlfw(WindowWidth, WindowHeight, WindowTitle)
		defer glfw.Terminate()
	}()
}
func TestInitGlfwFullSize(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("Init shouldn't be failed.")
			}
		}()
		InitGlfwFullSize(WindowTitle)
		defer glfw.Terminate()
	}()
}
func TestDummyKeyCallback(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("It shouldn't be failed.")
			}
		}()
		w := InitGlfw(WindowWidth, WindowHeight, WindowTitle)
		defer glfw.Terminate()
		DummyKeyCallback(w, glfw.KeyQ, 0, glfw.Release, 0)
	}()
}
func TestDummyMouseButtonCallback(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("It shouldn't be failed.")
			}
		}()
		w := InitGlfw(WindowWidth, WindowHeight, WindowTitle)
		defer glfw.Terminate()
		DummyMouseButtonCallback(w, glfw.MouseButtonLeft, glfw.Release, 0)
	}()
}
