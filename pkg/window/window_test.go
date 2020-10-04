package window

import (
	"reflect"
	"testing"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	WindowWidth  = 800
	WindowHeight = 800
	WindowTitle  = "Test title"
)

func TestWindowBuilderNew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("NewWindowBuilder shouldn't be failed.")
			}
		}()
		builder := NewWindowBuilder()
		defer glfw.Terminate()
		monitor := glfw.GetPrimaryMonitor()
		if !reflect.DeepEqual(monitor, builder.primaryMonitor) {
			t.Errorf("Invalid primary monitor. Instead of %#v, we have %#v.", monitor, builder.primaryMonitor)
		}
		// initial values
		if builder.width != DefaultWindowWidth {
			t.Errorf("Invalid default width. Instead of '%d', we have '%d'.", DefaultWindowWidth, builder.width)
		}
		if builder.height != DefaultWindowHeight {
			t.Errorf("Invalid default height. Instead of '%d', we have '%d'.", DefaultWindowHeight, builder.height)
		}
		if builder.title != DefaultWindowTitle {
			t.Errorf("Invalid default title. Instead of '%s', we have '%s'.", DefaultWindowTitle, builder.title)
		}
		if builder.fullScreen != DefaultWindowFullScreen {
			t.Errorf("Invalid default fullScreen. Instead of '%v', we have '%v'.", DefaultWindowFullScreen, builder.fullScreen)
		}
		if builder.decorated != DefaultWindowDecorated {
			t.Errorf("Invalid default decorated. Instead of '%v', we have '%v'.", DefaultWindowDecorated, builder.decorated)
		}
	}()
}
func TestWindowBuilderSetTitle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("SetTitle shouldn't be failed.")
			}
		}()
		builder := NewWindowBuilder()
		defer glfw.Terminate()
		testData := []string{"new title 01", "New title 02", "Last title", "", "longtitlelongtitlelongtitlelongtitlelongtitlelongtitle"}
		for i := 0; i < len(testData); i++ {
			builder.SetTitle(testData[i])
			if builder.title != testData[i] {
				t.Error("Failed set title")
			}
		}
	}()
}
func TestWindowBuilderSetFullScreen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("SetFullScreen shouldn't be failed.")
			}
		}()
		builder := NewWindowBuilder()
		defer glfw.Terminate()
		testData := []bool{true, true, false, true, false, false, false, true}
		for i := 0; i < len(testData); i++ {
			builder.SetFullScreen(testData[i])
			if builder.fullScreen != testData[i] {
				t.Error("Failed set fullScreen")
			}
		}
	}()
}
func TestWindowBuilderSetDecorated(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("SetDecorated shouldn't be failed.")
			}
		}()
		builder := NewWindowBuilder()
		defer glfw.Terminate()
		testData := []bool{true, true, false, true, false, false, false, true}
		for i := 0; i < len(testData); i++ {
			builder.SetDecorated(testData[i])
			if builder.decorated != testData[i] {
				t.Error("Failed set decorated")
			}
		}
	}()
}
func TestWindowBuilderSetWindowSize(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("SetWindowSize shouldn't be failed.")
			}
		}()
		builder := NewWindowBuilder()
		defer glfw.Terminate()
		testData := []struct {
			w int
			h int
		}{
			{100, 100},
			{200, 100},
			{300, 300},
			{700, 700},
			{800, 800},
			{1000, 800},
		}
		for _, tt := range testData {
			builder.SetWindowSize(tt.w, tt.h)
			if builder.width != tt.w {
				t.Error("Failed set width")
			}
			if builder.height != tt.h {
				t.Error("Failed set height")
			}
		}
	}()
}
func TestWindowBuilderGetCurrentMonitorResolution(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("GetCurrentMonitorResolution shouldn't be failed.")
			}
		}()
		builder := NewWindowBuilder()
		defer glfw.Terminate()
		builder.GetCurrentMonitorResolution()
	}()
}
func TestWindowBuilderGetCurrentMonitorPhysicalSize(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("GetCurrentMonitorPhysicalSize shouldn't be failed.")
			}
		}()
		builder := NewWindowBuilder()
		defer glfw.Terminate()
		builder.GetCurrentMonitorPhysicalSize()
	}()
}
func TestWindowBuilderGetCurrentMonitorContentScale(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("GetCurrentMonitorContentScale shouldn't be failed.")
			}
		}()
		builder := NewWindowBuilder()
		defer glfw.Terminate()
		builder.GetCurrentMonitorContentScale()
	}()
}
func TestWindowBuilderGetCurrentMonitorWorkarea(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("GetCurrentMonitorWorkarea shouldn't be failed.")
			}
		}()
		builder := NewWindowBuilder()
		defer glfw.Terminate()
		builder.GetCurrentMonitorWorkarea()
	}()
}
func TestWindowBuilderPrintCurrentMonitorData(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("PrintCurrentMonitorData shouldn't be failed.")
			}
		}()
		builder := NewWindowBuilder()
		defer glfw.Terminate()
		builder.PrintCurrentMonitorData()
	}()
}
func TestWindowBuilderBuild(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Error("Build shouldn't be failed.")
			}
		}()
		testData := []struct {
			w int
			h int
			d bool
			f bool
			t string
		}{
			{100, 100, false, false, "test 01"},
			{100, 100, false, true, "test 02"},
			{100, 100, true, true, "test 03"},
			{100, 100, true, false, "test 04"},
		}
		for _, tt := range testData {
			builder := NewWindowBuilder()
			defer glfw.Terminate()
			builder.SetWindowSize(tt.w, tt.h)
			builder.SetDecorated(tt.d)
			builder.SetFullScreen(tt.f)
			builder.SetTitle(tt.t)
			window := builder.Build()
			window.SetShouldClose(true)
			window.Destroy()
		}
	}()
}
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
