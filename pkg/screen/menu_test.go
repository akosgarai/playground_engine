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
func TestNewMenuScreenBuilder(t *testing.T) {
	b := NewMenuScreenBuilder()
	if b.menuItemDefaultMaterial != DefaultMenuItemMaterial {
		t.Error("Invalid default material.")
	}
	if b.menuItemHoverMaterial != HighlightMenuItemMaterial {
		t.Error("Invalid hover material.")
	}
	if b.menuItemFontColor != MenuItemFontColor {
		t.Error("Invalid font color.")
	}
	if b.backgroundColor != MenuItemFontColor {
		t.Error("Invalid background color.")
	}
}
func TestMenuScreenBuilderSetMenuItemSurfaceTexture(t *testing.T) {
	b := NewMenuScreenBuilder()
	var tex texture.Textures
	b.SetMenuItemSurfaceTexture(tex)
	if len(b.menuItemSurfaceTexture) != len(tex) {
		t.Error("Invalid texture has been set.")
	}
}
func TestMenuScreenBuilderSetMenuItemDefaultMaterial(t *testing.T) {
	b := NewMenuScreenBuilder()
	mat := material.Jade
	b.SetMenuItemDefaultMaterial(mat)
	if b.menuItemDefaultMaterial != mat {
		t.Error("Invalid default material has been set.")
	}
}
func TestMenuScreenBuilderSetMenuItemHighlightMaterial(t *testing.T) {
	b := NewMenuScreenBuilder()
	mat := material.Jade
	b.SetMenuItemHighlightMaterial(mat)
	if b.menuItemHoverMaterial != mat {
		t.Error("Invalid hover material has been set.")
	}
}
func TestMenuScreenBuilderSetMenuItemFontColor(t *testing.T) {
	b := NewMenuScreenBuilder()
	col := mgl32.Vec3{1, 1, 1}
	b.SetMenuItemFontColor(col)
	if b.menuItemFontColor != col {
		t.Error("Invalid font color has been set.")
	}
}
func TestMenuScreenBuilderSetBackgroundColor(t *testing.T) {
	b := NewMenuScreenBuilder()
	col := mgl32.Vec3{1, 1, 1}
	b.SetBackgroundColor(col)
	if b.backgroundColor != col {
		t.Error("Invalid background color has been set.")
	}
}
func TestMenuScreenBuilderSetState(t *testing.T) {
	b := NewMenuScreenBuilder()
	state := make(map[string]bool)
	state["foo"] = false
	state["bar"] = true
	b.SetState(state)
	if _, ok := b.state["foo"]; !ok {
		t.Error("Missing state key.")
	}
	if b.state["foo"] != state["foo"] {
		t.Error("Invalid state has been set.")
	}
	if _, ok := b.state["bar"]; !ok {
		t.Error("Missing state key.")
	}
	if b.state["bar"] != state["bar"] {
		t.Error("Invalid state has been set.")
	}
}
func TestMenuScreenBuilderAddOption(t *testing.T) {
	b := NewMenuScreenBuilder()
	if len(b.options) != 0 {
		t.Error("Invalid option length")
	}
	o := NewMenuScreenOption("test", nil, nil)
	b.AddOption(*o)
	if len(b.options) != 1 {
		t.Error("Invalid option length")
	}
}
func TestMenuScreenBuilderBuild(t *testing.T) {
	runtime.LockOSThread()
	testhelper.GlfwInit()
	wrapperReal.InitOpenGL()
	b := NewMenuScreenBuilder()
	defer testhelper.GlfwTerminate()
	// wo wrapper, it should panic
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("It should panic without wrapper.")
			}
		}()
		b.Build()
	}()
	b.SetWrapper(wrapperReal)
	b.SetWindowSize(100, 100)
	b.Build()
}
func NewTestMenuScreen(t *testing.T) *MenuScreen {
	testhelper.GlfwInit()
	wrapperReal.InitOpenGL()
	charset, err := model.LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40.0, 72, wrapperReal)
	if err != nil {
		t.Errorf("Error during load charset: %#v.", err)
	}
	dMat := material.Jade
	hMat := material.Ruby
	fontColor := mgl32.Vec3{0, 0, 0}
	bgColor := mgl32.Vec3{1, 1, 1}
	b := NewMenuScreenBuilder()
	b.SetWrapper(wrapperReal)
	b.SetWindowSize(100, 100)
	var tex texture.Textures
	b.SetMenuItemSurfaceTexture(tex)
	b.SetMenuItemDefaultMaterial(dMat)
	b.SetMenuItemHighlightMaterial(hMat)
	b.SetMenuItemFontColor(fontColor)
	b.SetBackgroundColor(bgColor)
	b.SetCharset(charset)
	menu := b.Build()

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
func NewTestMenuScreenWithOptions(t *testing.T) *MenuScreen {
	testhelper.GlfwInit()
	wrapperReal.InitOpenGL()
	charset, err := model.LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40.0, 72, wrapperReal)
	if err != nil {
		t.Errorf("Error during load charset: %#v.", err)
	}
	dMat := material.Jade
	hMat := material.Ruby
	fontColor := mgl32.Vec3{0, 0, 0}
	bgColor := mgl32.Vec3{1, 1, 1}
	b := NewMenuScreenBuilder()
	b.SetWrapper(wrapperReal)
	b.SetWindowSize(100, 100)
	var tex texture.Textures
	b.SetMenuItemSurfaceTexture(tex)
	b.SetMenuItemDefaultMaterial(dMat)
	b.SetMenuItemHighlightMaterial(hMat)
	b.SetMenuItemFontColor(fontColor)
	b.SetBackgroundColor(bgColor)
	b.SetCharset(charset)
	option := NewTestOption(t)
	b.AddOption(*option)
	b.AddOption(*option)
	menu := b.Build()

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
		menu.BuildScreen()
	}()
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
		menu := NewTestMenuScreenWithOptions(t)
		defer testhelper.GlfwTerminate()
		menu.BuildScreen()
		ks := store.NewGlfwKeyStore()
		ms := store.NewGlfwMouseStore()
		menu.Update(10, 0, -0.79, ks, ms)
		ms.Set(LEFT_MOUSE_BUTTON, true)
		menu.Update(10, 0, -0.79, ks, ms)
		ms.Set(LEFT_MOUSE_BUTTON, false)
		menu.Update(10, 0, -0.63, ks, ms)
		ks.Set(KEY_UP, true)
		menu.Update(10, 0, -0.74, ks, ms)
		ks.Set(KEY_UP, false)
		ks.Set(KEY_DOWN, true)
		menu.Update(10, 0, -0.74, ks, ms)
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
func TestMenuSetupMenu(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	runtime.LockOSThread()
	defer testhelper.GlfwTerminate()
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Errorf("Shouldn't have panic, %#v.", r)
			}
		}()
		runtime.LockOSThread()
		menu := NewTestMenuScreen(t)
		menu.setupMenu(wrapperReal)
	}()
}
