package application

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/screen"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	cam         = camera.NewCamera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}, 0, 0)
	wm          testhelper.WindowMock
	wrapperMock testhelper.GLWrapperMock
)

func TestNew(t *testing.T) {
	app := New(wrapperMock)
	if len(app.screens) != 0 {
		t.Error("Invalid application - screens should be empty")
	}
}
func TestLog(t *testing.T) {
	app := New(wrapperMock)
	log := app.Log()
	if len(log) < 10 {
		t.Error("Log too short.")
	}
}
func TestSetWindow(t *testing.T) {
	app := New(wrapperMock)
	app.SetWindow(wm)

	if app.window != wm {
		t.Error("Invalid window setup.")
	}
}
func TestGetWindow(t *testing.T) {
	app := New(wrapperMock)
	app.SetWindow(wm)

	if app.GetWindow() != wm {
		t.Error("Invalid window setup.")
	}
}
func TestGetCamera(t *testing.T) {
	s := screen.New()
	s.SetCamera(cam)
	app := New(wrapperMock)
	app.AddScreen(s)
	app.ActivateScreen(s)
	if app.GetCamera() != cam {
		t.Error("Invalid camera")
	}
}
func TestSetKeyState(t *testing.T) {
	app := New(wrapperMock)
	app.SetKeyState(glfw.KeyW, glfw.Press)
	if !app.keyDowns.Get(glfw.KeyW) {
		t.Error("W should be pressed")
	}
	app.SetKeyState(glfw.KeyW, glfw.Release)
	if app.keyDowns.Get(glfw.KeyW) {
		t.Error("W should be released")
	}
}
func TestSetButtonState(t *testing.T) {
	app := New(wrapperMock)
	app.SetButtonState(glfw.MouseButtonLeft, glfw.Press)
	if !app.mouseDowns.Get(glfw.MouseButtonLeft) {
		t.Error("LMB should be pressed")
	}
	app.SetButtonState(glfw.MouseButtonLeft, glfw.Release)
	if app.mouseDowns.Get(glfw.MouseButtonLeft) {
		t.Error("LMB should be released")
	}
}
func TestGetMouseButtonState(t *testing.T) {
	app := New(wrapperMock)
	app.SetButtonState(glfw.MouseButtonLeft, glfw.Press)
	if !app.GetMouseButtonState(glfw.MouseButtonLeft) {
		t.Error("LMB should be pressed")
	}
	app.SetButtonState(glfw.MouseButtonLeft, glfw.Release)
	if app.GetMouseButtonState(glfw.MouseButtonLeft) {
		t.Error("LMB should be released")
	}
}
func TestGetKeyState(t *testing.T) {
	app := New(wrapperMock)
	app.SetKeyState(glfw.KeyW, glfw.Press)
	if !app.GetKeyState(glfw.KeyW) {
		t.Error("W should be pressed")
	}
	app.SetKeyState(glfw.KeyW, glfw.Release)
	if app.GetKeyState(glfw.KeyW) {
		t.Error("W should be released")
	}
}
func TestAddScreen(t *testing.T) {
	app := New(wrapperMock)
	s := screen.New()
	app.AddScreen(s)

	if len(app.screens) != 1 {
		t.Errorf("Invalid screen length. Instead of '1', we have '%d'.", len(app.screens))
	}
}
func TestActiveScreen(t *testing.T) {
	app := New(wrapperMock)
	s := screen.New()
	app.AddScreen(s)
	app.ActivateScreen(s)

	if app.activeScreen != s {
		t.Error("Invalid active screen.")
	}
}
func TestUpdateWoScreen(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing activeScreen.")
			}
		}()
		app := New(wrapperMock)
		app.SetWindow(wm)
		// wo everything
		app.Update(10)
	}()
}
func TestUpdateWithScreen(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
			}
		}()
		app := New(wrapperMock)
		app.SetWindow(wm)
		s := screen.New()
		app.AddScreen(s)
		app.ActivateScreen(s)
		// wo everything
		app.Update(10)
	}()
}
func TestSetUniformFloat(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New(wrapperMock)
		s := screen.New()
		app.AddScreen(s)
		key := "testname"
		value := float32(2.2)
		app.SetUniformFloat(key, value)
		app.ActivateScreen(s)
		app.SetUniformFloat(key, value)
	}()
}
func TestSetUniformVector(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New(wrapperMock)
		s := screen.New()
		app.AddScreen(s)
		key := "testname"
		value := mgl32.Vec3{42.0, 0.0, 0.0}
		app.SetUniformVector(key, value)
		app.ActivateScreen(s)
		app.SetUniformVector(key, value)
	}()
}
func TestDraw(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New(wrapperMock)
		app.SetWindow(wm)
		app.SetUniformFloat("testFloat", float32(42.0))
		app.SetUniformVector("testVector", mgl32.Vec3{1.0, 2.0, 3.0})
		s := screen.New()
		app.AddScreen(s)
		app.ActivateScreen(s)
		app.Draw(wrapperMock)
	}()
}
func TestGetClosestModelMeshDistanceWoActiveScreen(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic.")
			}
		}()
		app := New(wrapperMock)
		s := screen.New()
		app.AddScreen(s)
		_, _, _ = app.GetClosestModelMeshDistance()
	}()
}
func TestGetClosestModelMeshDistanceWithActiveScreen(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New(wrapperMock)
		s := screen.New()
		app.AddScreen(s)
		app.ActivateScreen(s)
		_, _, _ = app.GetClosestModelMeshDistance()
	}()
}
func TestMenuScreen(t *testing.T) {
	s := screen.New()
	s.SetCamera(cam)
	app := New(wrapperMock)
	app.AddScreen(s)
	app.MenuScreen(s)
	if app.menuSet != true {
		t.Error("Invalid set value for menu")
	}
	if app.menuScreen != s {
		t.Error("Invalid menu screen")
	}
}
func TestSetWrapper(t *testing.T) {
	app := New(wrapperMock)
	app.SetWrapper(wrapperMock)
	if app.wrapper != wrapperMock {
		t.Error("Invalid wrapper mock")
	}
}
func TestGetWrapper(t *testing.T) {
	app := New(wrapperMock)
	app.SetWrapper(wrapperMock)
	if app.GetWrapper() != wrapperMock {
		t.Error("Invalid wrapper mock")
	}
}
