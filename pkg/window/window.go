package window

import (
	"fmt"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// InitGlfw returns a *glfw.Windows instance.
func InitGlfw(windowWidth, windowHeight int, windowTitle string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, glwrapper.GL_MAJOR_VERSION)
	glfw.WindowHint(glfw.ContextVersionMinor, glwrapper.GL_MINOR_VERSION)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(windowWidth, windowHeight, windowTitle, nil, nil)

	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	window.MakeContextCurrent()

	return window
}

// InitGlfwFullSize returns a *glfw.Window instance.
// The width and height of the window is based on the size of
// the workarea.
func InitGlfwFullSize(windowTitle string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}

	glfw.WindowHint(glfw.ContextVersionMajor, glwrapper.GL_MAJOR_VERSION)
	glfw.WindowHint(glfw.ContextVersionMinor, glwrapper.GL_MINOR_VERSION)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	monitor := glfw.GetPrimaryMonitor()
	_, _, waw, wah := monitor.GetWorkarea()

	window, err := glfw.CreateWindow(waw, wah, windowTitle, nil, nil)

	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	window.MakeContextCurrent()

	return window
}

// DummyKeyCallback is responsible for the keyboard event handling with log.
// So this function does nothing but printing out the input parameters.
func DummyKeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("KeyCallback has been called with the following options: key: '%d', scancode: '%d', action: '%d'!, mods: '%d'\n", key, scancode, action, mods)
}

// DummyMouseButtonCallback is responsible for the mouse button event handling with log.
// So this function does nothing but printing out the input parameters.
func DummyMouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("MouseButtonCallback has been called with the following options: button: '%d', action: '%d'!, mods: '%d'\n", button, action, mods)
}
