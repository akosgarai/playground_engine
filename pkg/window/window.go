package window

import (
	"fmt"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	DefaultWindowWidth      = 800
	DefaultWindowHeight     = 800
	DefaultWindowTitle      = "Playground application"
	DefaultWindowFullScreen = false
	DefaultWindowDecorated  = true
)

type WindowBuilder struct {
	primaryMonitor *glfw.Monitor
	width          int
	height         int
	title          string
	fullScreen     bool
	decorated      bool
}

// NewWindowBuilder returns a WindowBuilder. If the glfw lib could not be initialized, it panics.
func NewWindowBuilder() *WindowBuilder {
	if err := glfw.Init(); err != nil {
		panic(fmt.Errorf("could not initialize glfw: %v", err))
	}
	return &WindowBuilder{
		primaryMonitor: glfw.GetPrimaryMonitor(),
		width:          DefaultWindowWidth,
		height:         DefaultWindowHeight,
		title:          DefaultWindowTitle,
		fullScreen:     DefaultWindowFullScreen,
		decorated:      DefaultWindowDecorated,
	}
}

// The title will be used as 3. parameter of the CreateWindow function.
func (b *WindowBuilder) SetTitle(t string) {
	b.title = t
}

// If this flag is set true, the current monitor will be attached as monitor (4. parameter) to the CreateWindow function.
func (b *WindowBuilder) SetFullScreen(f bool) {
	b.fullScreen = f
}

// If this flag is set true, the current monitor will be attached as monitor (4. parameter) to the CreateWindow function.
func (b *WindowBuilder) SetDecorated(d bool) {
	b.decorated = d
}

// The width and height values are used as the first 2 parameter of the CreateWindow function.
func (b *WindowBuilder) SetWindowSize(w, h int) {
	b.width = w
	b.height = h
}

// GetCurrentMonitorResolution returns the current video mode of the monitor. If you are using a full screen window, the return value will therefore depend on whether it is focused.
func (b *WindowBuilder) GetCurrentMonitorResolution() (int, int) {
	mode := b.primaryMonitor.GetVideoMode()
	return mode.Width, mode.Height
}

// GetCurrentMonitorPhysicalSize returns the size, in millimeters, of the display area of the monitor.
func (b *WindowBuilder) GetCurrentMonitorPhysicalSize() (int, int) {
	return b.primaryMonitor.GetPhysicalSize()
}

// GetCurrentMonitorContentScale function retrieves the content scale for the specified monitor. The content scale is the ratio between the current DPI and the platform's default DPI. If you scale all pixel dimensions by this scale then your content should appear at an appropriate size. This is especially important for text and any UI elements
func (b *WindowBuilder) GetCurrentMonitorContentScale() (float32, float32) {
	return b.primaryMonitor.GetContentScale()
}

// GetCurrentMonitorWorkarea returns the position, in screen coordinates, of the upper-left corner of the work area of the specified monitor along with the work area size in screen coordinates. The work area is defined as the area of the monitor not occluded by the operating system task bar where present. If no task bar exists then the work area is the monitor resolution in screen coordinates.
func (b *WindowBuilder) GetCurrentMonitorWorkarea() (int, int, int, int) {
	return b.primaryMonitor.GetWorkarea()
}

// PrintCurrentMonitorData prints the resolution, phisical size, content scale, workarea to the stdout.
func (b *WindowBuilder) PrintCurrentMonitorData() {
	w, h := b.GetCurrentMonitorResolution()
	sx, sy := b.GetCurrentMonitorPhysicalSize()
	cx, cy := b.GetCurrentMonitorContentScale()
	wax, way, waw, wah := b.GetCurrentMonitorWorkarea()
	fmt.Printf("Current monitor video mode: %d * %d\n", w, h)
	fmt.Printf("Current monitor physical size: %d * %d\n", sx, sy)
	fmt.Printf("Current monitor content scale: %f * %f\n", cx, cy)
	fmt.Printf("Current monitor workarea: %d - %d, %d * %d\n", wax, way, waw, wah)
}
func (b *WindowBuilder) windowHints() {
	glfw.WindowHint(glfw.ContextVersionMajor, glwrapper.GL_MAJOR_VERSION)
	glfw.WindowHint(glfw.ContextVersionMinor, glwrapper.GL_MINOR_VERSION)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	// Specifies whether the window will be resizable by the user.
	if b.fullScreen {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	}
	// Specifies whether the window will have window decorations such as a border, a close widget, etc.
	if b.decorated {
		glfw.WindowHint(glfw.Decorated, glfw.True)
	} else {
		glfw.WindowHint(glfw.Decorated, glfw.False)
	}
	// Specified whether the window content area should be resized based on the monitor content scale of any monitor it is placed on.
	// This includes the initial placement when the window is created.
	glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)
}
func (b *WindowBuilder) Build() *glfw.Window {
	b.windowHints()

	var window *glfw.Window
	var err error
	if b.fullScreen {
		window, err = glfw.CreateWindow(b.width, b.height, b.title, b.primaryMonitor, nil)
	} else {
		window, err = glfw.CreateWindow(b.width, b.height, b.title, nil, nil)
	}

	if err != nil {
		panic(fmt.Errorf("could not create opengl renderer: %v", err))
	}

	window.MakeContextCurrent()

	return window

}

// InitGlfw returns a *glfw.Windows instance.
func InitGlfw(windowWidth, windowHeight int, windowTitle string) *glfw.Window {
	builder := NewWindowBuilder()
	builder.SetWindowSize(windowWidth, windowHeight)
	builder.SetFullScreen(false)
	builder.SetDecorated(true)
	builder.SetTitle(windowTitle)
	return builder.Build()
}

// InitGlfwFullSize returns a *glfw.Window instance.
// The width and height of the window is based on the size of
// the workarea.
func InitGlfwFullSize(windowTitle string) *glfw.Window {
	builder := NewWindowBuilder()
	w, h := builder.GetCurrentMonitorResolution()
	builder.SetWindowSize(w, h)
	builder.SetFullScreen(true)
	builder.SetDecorated(false)
	builder.SetTitle(windowTitle)
	return builder.Build()
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
