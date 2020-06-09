package testhelper

import (
	"fmt"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Test title"
)

// This function returns true, if the given a, b is almost equal,
// the difference between them is less than epsilon.
func Float32ApproxEqual(a, b, epsilon float32) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}

func GlfwInit() {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, WindowTitle, nil, nil)

	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	window.MakeContextCurrent()
}
func GlfwTerminate() {
	glfw.Terminate()
}
