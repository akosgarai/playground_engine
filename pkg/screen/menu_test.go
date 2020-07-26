package screen

import (
	"runtime"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/store"
	"github.com/akosgarai/playground_engine/pkg/testhelper"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

func NewTestOption(t *testing.T) *Option {
	dcb := func(m map[string]bool) bool { return true }
	ecb := func() {}
	label := "option-label"
	opt := NewMenuScreenOption(label, dcb, ecb)
	if opt.label != label {
		t.Errorf("Invalid label. Instead of '%s', we have '%s'.", label, opt.label)
	}
	return opt
}
func TestNewMenuScreenOption(t *testing.T) {
	_ = NewTestOption(t)
}
func TestOptionSetSurface(t *testing.T) {
	opt := NewTestOption(t)
	var msh interfaces.Mesh
	opt.SetSurface(msh)
	if opt.surface != msh {
		t.Error("Invalid surface")
	}
}
func TestOptionDisplayCondition(t *testing.T) {
	dcb := func(m map[string]bool) bool { return m["test"] }
	ecb := func() {}
	label := "option-label"
	opt := NewMenuScreenOption(label, dcb, ecb)
	mp := make(map[string]bool)
	if opt.DisplayCondition(mp) != false {
		t.Error("Invalid display condition result")
	}
	mp["test"] = true
	if opt.DisplayCondition(mp) != true {
		t.Error("Invalid display condition result")
	}
	mp["test"] = false
	if opt.DisplayCondition(mp) != false {
		t.Error("Invalid display condition result")
	}
}
func NewTestMenuScreen(t *testing.T) *MenuScreen {
	testhelper.GlfwInit()
	var tex texture.Textures
	wrapperReal.InitOpenGL()
	charset, err := model.LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40.0, 72, wrapperReal)
	if err != nil {
		t.Errorf("Error during load charset: %#v.", err)
	}
	dMat := material.Jade
	hMat := material.Ruby
	fontColor := mgl32.Vec3{0, 0, 0}
	bgColor := mgl32.Vec3{1, 1, 1}
	menu := NewMenuScreen(tex, dMat, hMat, charset, fontColor, bgColor, wrapperReal)

	if menu.defaultMaterial != dMat {
		t.Errorf("Invalid default material. Instead of '%#v', we have '%#v'.", dMat, menu.defaultMaterial)
	}
	if menu.hoverMaterial != hMat {
		t.Errorf("Invalid hover material. Instead of '%#v', we have '%#v'.", hMat, menu.hoverMaterial)
	}
	if menu.charset != charset {
		t.Error("Invalid charset")
	}
	if menu.fontColor[0] != fontColor {
		t.Error("Invalid font color")
	}
	if menu.backgroundColor != bgColor {
		t.Error("Invalid background color")
	}
	return menu
}
func TestNewMenuScreen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	runtime.LockOSThread()
	_ = NewTestMenuScreen(t)
	defer testhelper.GlfwTerminate()
}
func TestMenuScreenBuildScreen(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		runtime.LockOSThread()
		menu := NewTestMenuScreen(t)
		defer testhelper.GlfwTerminate()
		menu.BuildScreen(wrapperReal, 1.0)
		option := NewTestOption(t)
		menu.AddOption(*option)
		menu.BuildScreen(wrapperReal, 1.0)
	}()
}
func TestMenuScreenAddOption(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	runtime.LockOSThread()
	menu := NewTestMenuScreen(t)
	defer testhelper.GlfwTerminate()
	if len(menu.options) != 0 {
		t.Error("Invalid initial option")
	}
	var option Option
	menu.AddOption(option)
	if len(menu.options) != 1 {
		t.Errorf("Invalid number of options (%d)", len(menu.options))
	}
}
func TestMenuScreenUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		runtime.LockOSThread()
		menu := NewTestMenuScreen(t)
		defer testhelper.GlfwTerminate()
		menu.BuildScreen(wrapperReal, 1.0)
		option := NewTestOption(t)
		menu.AddOption(*option)
		menu.BuildScreen(wrapperReal, 1.0)
		ks := store.NewGlfwKeyStore()
		ms := store.NewGlfwMouseStore()
		menu.Update(10, 0, 0, ks, ms)
	}()
}
func TestMenuSetState(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	runtime.LockOSThread()
	defer testhelper.GlfwTerminate()
	menu := NewTestMenuScreen(t)
	if menu.state["world-started"] != false {
		t.Error("Invalid initial state")
	}
	menu.SetState("world-started", true)
	if menu.state["world-started"] != true {
		t.Error("Invalid state")
	}
}
