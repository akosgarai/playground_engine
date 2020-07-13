package application

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/store"
	"github.com/akosgarai/playground_engine/pkg/transformations"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DEBUG  = glfw.KeyH
	EXPORT = glfw.KeyP
	ESCAPE = glfw.KeyEscape
)

type Camera interface {
	Log() string
	GetViewMatrix() mgl32.Mat4
	GetProjectionMatrix() mgl32.Mat4
	Walk(float32)
	Strafe(float32)
	Lift(float32)
	UpdateDirection(float32, float32)
	GetPosition() mgl32.Vec3
	GetVelocity() float32
	GetRotationStep() float32
	BoundingObjectAfterWalk(float32) *coldet.Sphere
	BoundingObjectAfterStrafe(float32) *coldet.Sphere
	BoundingObjectAfterLift(float32) *coldet.Sphere
}

type Window interface {
	GetCursorPos() (float64, float64)
	SetKeyCallback(glfw.KeyCallback) glfw.KeyCallback
	SetMouseButtonCallback(glfw.MouseButtonCallback) glfw.MouseButtonCallback
	SetCharCallback(glfw.CharCallback) glfw.CharCallback
	ShouldClose() bool
	SwapBuffers()
	GetSize() (int, int)
	SetShouldClose(bool)
}

type Application struct {
	window Window

	mouseDowns interfaces.ButtonStore
	MousePosX  float64
	MousePosY  float64

	keyDowns interfaces.KeyStore

	// screens
	activeScreen interfaces.Screen
	screens      []interfaces.Screen
	menuScreen   interfaces.Screen
	menuSet      bool

	// wrapper for char callback
	wrapper interfaces.GLWrapper
}

// New returns an application instance
func New(wrapper interfaces.GLWrapper) *Application {
	return &Application{
		mouseDowns: store.NewGlfwMouseStore(),
		keyDowns:   store.NewGlfwKeyStore(),
		menuSet:    false,
		wrapper:    wrapper,
	}
}

// Log returns the string representation of this object.
func (a *Application) Log() string {
	logString := "Application:\n"
	logString += fmt.Sprintf("\tKeyDowns: %#v\n", a.keyDowns)
	logString += fmt.Sprintf("\tMouseDowns: %#v\n", a.mouseDowns)
	if a.activeScreen != nil {
		logString += fmt.Sprintf("\tactiveScreen: %s\n", a.activeScreen.Log())
	}
	return logString
}

// SetWrapper updates the wrapper with the new one.
func (a *Application) SetWrapper(w interfaces.GLWrapper) {
	a.wrapper = w
}

// GetWrapper returns the current wrapper of the application.
func (a *Application) GetWrapper() interfaces.GLWrapper {
	return a.wrapper
}

// SetWindow updates the window with the new one.
func (a *Application) SetWindow(w Window) {
	a.window = w
}

// GetWindow returns the current window of the application.
func (a *Application) GetWindow() Window {
	return a.window
}

// GetCamera returns the camera of the activeScreen.
func (a *Application) GetCamera() interfaces.Camera {
	return a.activeScreen.GetCamera()
}

// GetClosestModelMeshDistance calls GetClosestModelMeshDistance on the activeScreen.
func (a *Application) GetClosestModelMeshDistance() (interfaces.Model, interfaces.Mesh, float32) {
	return a.activeScreen.GetClosestModelMeshDistance()
}

// Update sets up the mouse coordinates and calls Update function on the activeScreen.
func (a *Application) Update(dt float64) {
	MousePosX, MousePosY := a.window.GetCursorPos()
	WindowWidth, WindowHeight := a.window.GetSize()
	mX, mY := transformations.MouseCoordinates(MousePosX, MousePosY, float64(WindowWidth), float64(WindowHeight))
	a.activeScreen.Update(dt, mX, mY, a.keyDowns, a.mouseDowns)
}

// AddScreen appends the screen to screens.
func (a *Application) AddScreen(s interfaces.Screen) {
	if a.window != nil {
		WindowWidth, WindowHeight := a.window.GetSize()
		s.SetWindowSize(float32(WindowWidth), float32(WindowHeight))
		s.SetWrapper(a.wrapper)
	}
	a.screens = append(a.screens, s)
}

// ActivateScreen sets the given screen to active screen
func (a *Application) ActivateScreen(s interfaces.Screen) {
	a.activeScreen = s
}

// MenuScreen sets the given screen to menu screen and the menuSet
// flag true
func (a *Application) MenuScreen(s interfaces.Screen) {
	a.menuScreen = s
	a.menuSet = true
}

// Draw calls Draw function on the activeScreen.
func (a *Application) Draw(wrapper interfaces.GLWrapper) {
	a.activeScreen.Draw(wrapper)
}

// KeyCallback is responsible for the keyboard event handling.
// After configuring the keyboard well, the esc character seems to
// be working well.
func (a *Application) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case ESCAPE:
		if a.menuSet {
			a.ActivateScreen(a.menuScreen)
		} else {
			a.window.SetShouldClose(true)
		}
		break
	case DEBUG:
		if action != glfw.Release {
			fmt.Printf("%s\n", a.Log())
		}
		break
	case EXPORT:
		if action != glfw.Release {
			a.export()
		}
		break
	default:
		a.SetKeyState(key, action)
		break
	}
}

// MouseButtonCallback is responsible for the mouse button event handling.
func (a *Application) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	a.MousePosX, a.MousePosY = w.GetCursorPos()
	switch button {
	default:
		a.SetButtonState(button, action)
		break
	}
}

// CharCallback is responsible for the character stream input (typing on keyboard)
func (a *Application) CharCallback(w *glfw.Window, char rune) {
	if a.activeScreen != nil {
		a.activeScreen.CharCallback(char, a.wrapper)
	}
}

// SetKeyState setups the keyDowns based on the key and action
func (a *Application) SetKeyState(key glfw.Key, action glfw.Action) {
	var isButtonPressed bool
	if action != glfw.Release {
		isButtonPressed = true
	} else {
		isButtonPressed = false
	}
	a.keyDowns.Set(key, isButtonPressed)
}

// SetButtonState setups the mouseDowns based on the key and action
func (a *Application) SetButtonState(button glfw.MouseButton, action glfw.Action) {
	var isButtonPressed bool
	if action != glfw.Release {
		isButtonPressed = true
	} else {
		isButtonPressed = false
	}
	a.mouseDowns.Set(button, isButtonPressed)
}

// GetMouseButtonState returns the state of the given button
func (a *Application) GetMouseButtonState(button glfw.MouseButton) bool {
	return a.mouseDowns.Get(button)
}

// GetKeyState returns the state of the given key
func (a *Application) GetKeyState(key glfw.Key) bool {
	return a.keyDowns.Get(key)
}

// SetUniformFloat loops on screens and sets the given float value to the given string key in
// the uniformFloat map of the screen.
func (a *Application) SetUniformFloat(key string, value float32) {
	for i, _ := range a.screens {
		a.screens[i].SetUniformFloat(key, value)
	}
}

// SetUniformVector loops on screens and sets the given mgl32.Vec3 value to the given string key in
// the uniformVector map of the screen.
func (a *Application) SetUniformVector(key string, value mgl32.Vec3) {
	for i, _ := range a.screens {
		a.screens[i].SetUniformVector(key, value)
	}
}

// This function is called for starting the export process. It is attached to a key callback.
func (a *Application) export() {
	ExportBaseDir := "./exports"
	Directory := time.Now().Format("20060102150405")
	err := os.Mkdir(ExportBaseDir+"/"+Directory, os.ModeDir|os.ModePerm)
	if err != nil {
		fmt.Printf("Cannot create export directory. '%s'\n", err.Error())
	}
	i := 0
	for s, _ := range a.screens {
		modelDir := strconv.Itoa(i)
		err := os.Mkdir(ExportBaseDir+"/"+Directory+"/"+modelDir, os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Printf("Cannot create model directory. '%s'\n", err.Error())
		}
		a.screens[s].Export(ExportBaseDir + "/" + Directory + "/" + modelDir)
		i++
	}
}
