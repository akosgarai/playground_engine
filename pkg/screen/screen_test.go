package screen

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/config"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/store"
	"github.com/akosgarai/playground_engine/pkg/testhelper"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DefaultFormItemDescription = "form item description."
)

var (
	cam = camera.NewCamera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}, 0, 0)

	sm          testhelper.ShaderMock
	wrapperMock testhelper.GLWrapperMock
	wrapperReal glwrapper.Wrapper
)

func TestNew(t *testing.T) {
	shader := New()
	if len(shader.shaderMap) != 0 {
		t.Error("Invalid screen - shadermap should be empty")
	}
	if shader.cameraSet {
		t.Error("Camera shouldn't be set")
	}
}
func TestLog(t *testing.T) {
	emptylog := "Screen:\n"
	var screen Screen
	log := screen.Log()
	if log != emptylog {
		t.Errorf("Invalid log for empty screen. Instead of '%s', we have '%s'.", emptylog, log)
	}
	screen.cameraSet = true
	screen.camera = cam
	log = screen.Log()
	if log == emptylog {
		t.Errorf("Invalid log for camera. We have the same as before '%s'.", emptylog)
	}
}
func TestSetCameraMovementMap(t *testing.T) {
	var screen Screen
	cmMap := make(map[string]glfw.Key)
	screen.SetCameraMovementMap(cmMap)

	if !reflect.DeepEqual(screen.cameraKeyboardMovementMap, cmMap) {
		t.Error("Invalid camera movement map has been set.")
	}
}
func TestSetRotateOnEdgeDistance(t *testing.T) {
	var screen Screen
	testData := []struct {
		input    float32
		expected float32
	}{
		{0.5, 0.5},
		{0.0, 0.0},
		{-0.1, 0.0},
		{0.9, 0.9},
		{1.1, 0.9},
		{1.0, 1.0},
	}
	for _, tt := range testData {
		screen.SetRotateOnEdgeDistance(tt.input)
		if screen.rotateOnEdgeDistance != tt.expected {
			t.Errorf("Invalid rotateOnEdgeDistance. Instead of '%f', we have '%f'.", tt.expected, screen.rotateOnEdgeDistance)
		}
	}
}
func TestSetCamera(t *testing.T) {
	var screen Screen
	screen.SetCamera(cam)

	if screen.camera != cam {
		t.Error("Invalid camera setup.")
	}
}
func TestGetCamera(t *testing.T) {
	var screen Screen
	screen.SetCamera(cam)

	if screen.GetCamera() != cam {
		t.Error("Invalid camera setup.")
	}
}
func TestAddShader(t *testing.T) {
	screen := New()
	if len(screen.shaderMap) != 0 {
		t.Errorf("Invalid shader map length. Instead of '0', it is '%d'.\n", len(screen.shaderMap))
	}
	screen.AddShader(sm)
	if len(screen.shaderMap) != 1 {
		t.Errorf("Invalid shader map length. Instead of '1', it is '%d'.\n", len(screen.shaderMap))
	}
}
func TestAddModelToShader(t *testing.T) {
	screen := New()
	screen.AddShader(sm)
	if len(screen.shaderMap[sm]) != 0 {
		t.Errorf("Invalid model length. Instead of '0', it is '%d'.\n", len(screen.shaderMap[sm]))
	}
	mod := model.New()
	screen.AddModelToShader(mod, sm)
	if len(screen.shaderMap[sm]) != 1 {
		t.Errorf("Invalid model length. Instead of '1', it is '%d'.\n", len(screen.shaderMap[sm]))
	}
}
func TestGetClosestModelMeshDistance(t *testing.T) {
	screen := New()
	mod := model.New()
	screen.closestModel = mod
	screen.closestDistance = 1.1
	cmod, _, dst := screen.GetClosestModelMeshDistance()
	if cmod != mod {
		t.Errorf("The mod should be the closest, but we have %#v.", cmod)
	}
	if dst != 1.1 {
		t.Errorf("The destination supposed to be 1.1, but we have '%f'", dst)
	}
}
func TestSetUniformFloat(t *testing.T) {
	screen := New()
	key := "testname"
	value := float32(42.0)
	if _, ok := screen.uniformFloat[key]; ok {
		t.Errorf("Key '%s' shouldn't be set.", key)
	}
	screen.SetUniformFloat(key, value)
	if _, ok := screen.uniformFloat[key]; !ok {
		t.Errorf("Key '%s' should be set.", key)
	}
	val := screen.uniformFloat[key]
	if val != value {
		t.Errorf("Invalud value for key '%s'. Instead of '%f', we have '%f'.", key, value, val)
	}
}
func TestSetUniformVector(t *testing.T) {
	screen := New()
	key := "testname"
	value := mgl32.Vec3{42.0, 0.0, 0.0}
	if _, ok := screen.uniformVector[key]; ok {
		t.Errorf("Key '%s' shouldn't be set.", key)
	}
	screen.SetUniformVector(key, value)
	if _, ok := screen.uniformVector[key]; !ok {
		t.Errorf("Key '%s' should be set.", key)
	}
	val := screen.uniformVector[key]
	if val != value {
		t.Errorf("Invalud value for key '%s'. Instead of '%v', we have '%v'.", key, value, val)
	}
}
func TestAddDirectionalLightSource(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		screen := New()
		if len(screen.directionalLightSources) != 0 {
			t.Errorf("Invalid number if dir. length. Instead of '0', it is '%d'.", len(screen.directionalLightSources))
		}
		ds := light.NewDirectionalLight([4]mgl32.Vec3{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
		})
		uniforms := [4]string{"u1", "u2", "u3", "u4"}
		screen.AddDirectionalLightSource(ds, uniforms)
		if len(screen.directionalLightSources) != 1 {
			t.Errorf("Invalid number if dir. length. Instead of '1', it is '%d'.", len(screen.directionalLightSources))
		}
		screen.AddDirectionalLightSource(ds, uniforms)
		if len(screen.directionalLightSources) != 2 {
			t.Errorf("Invalid number if dir. length. Instead of '2', it is '%d'.", len(screen.directionalLightSources))
		}
	}()
}
func TestAddPointLightSource(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		screen := New()
		if len(screen.pointLightSources) != 0 {
			t.Errorf("Invalid number if point. length. Instead of '0', it is '%d'.", len(screen.pointLightSources))
		}
		ps := light.NewPointLight([4]mgl32.Vec3{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
		},
			[3]float32{1, 0.5, 0.05},
		)
		uniforms := [7]string{"u1", "u2", "u3", "u4", "u5", "u6", "u7"}
		screen.AddPointLightSource(ps, uniforms)
		if len(screen.pointLightSources) != 1 {
			t.Errorf("Invalid number if point. length. Instead of '1', it is '%d'.", len(screen.pointLightSources))
		}
		screen.AddPointLightSource(ps, uniforms)
		if len(screen.pointLightSources) != 2 {
			t.Errorf("Invalid number if point. length. Instead of '2', it is '%d'.", len(screen.pointLightSources))
		}
	}()
}
func TestAddSpotLightSource(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		screen := New()
		if len(screen.spotLightSources) != 0 {
			t.Errorf("Invalid number if spot. length. Instead of '0', it is '%d'.", len(screen.spotLightSources))
		}
		ss := light.NewSpotLight([5]mgl32.Vec3{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
		},
			[5]float32{1, 0.5, 0.05, 4, 5},
		)
		uniforms := [10]string{"u1", "u2", "u3", "u4", "u5", "u6", "u7", "u8", "u9", "u10"}
		screen.AddSpotLightSource(ss, uniforms)
		if len(screen.spotLightSources) != 1 {
			t.Errorf("Invalid number if spot. length. Instead of '1', it is '%d'.", len(screen.spotLightSources))
		}
		screen.AddSpotLightSource(ss, uniforms)
		if len(screen.spotLightSources) != 2 {
			t.Errorf("Invalid number if spot. length. Instead of '2', it is '%d'.", len(screen.spotLightSources))
		}
	}()
}
func TestSetupDirectionalLightForShader(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		screen := New()
		ds := light.NewDirectionalLight([4]mgl32.Vec3{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
		})
		uniforms := [4]string{"u1", "u2", "u3", "u4"}
		screen.AddDirectionalLightSource(ds, uniforms)
		if len(screen.directionalLightSources) != 1 {
			t.Errorf("Invalid number if dir. length. Instead of '1', it is '%d'.", len(screen.directionalLightSources))
		}
		screen.setupDirectionalLightForShader(sm)
	}()
}
func TestSetupPointLightForShader(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		screen := New()
		ps := light.NewPointLight([4]mgl32.Vec3{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
		},
			[3]float32{1, 0.5, 0.05},
		)
		uniforms := [7]string{"u1", "u2", "u3", "u4", "u5", "u6", "u7"}
		screen.AddPointLightSource(ps, uniforms)
		screen.setupPointLightForShader(sm)
	}()
}
func TestSetupSpotLightForShader(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		screen := New()
		ss := light.NewSpotLight([5]mgl32.Vec3{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
		},
			[5]float32{1, 0.5, 0.05, 4, 5},
		)
		uniforms := [10]string{"u1", "u2", "u3", "u4", "u5", "u6", "u7", "u8", "u9", "u10"}
		screen.AddSpotLightSource(ss, uniforms)
		screen.setupSpotLightForShader(sm)
	}()
}
func TestCustomUniforms(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		screen := New()
		keyVec := "testname"
		valueVec := mgl32.Vec3{42.0, 0.0, 0.0}
		screen.SetUniformVector(keyVec, valueVec)
		key := "testNameOther"
		value := float32(3.0)
		screen.SetUniformFloat(key, value)
		screen.customUniforms(sm)
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
		screen := New()
		screen.AddShader(sm)
		mod := model.New()
		screen.AddModelToShader(mod, sm)
		msh := mesh.NewPointMesh(wrapperMock)
		mod.AddMesh(msh)
		// wo camera
		screen.Draw(wrapperMock)
		// with camera
		screen.SetCamera(cam)
		screen.Draw(wrapperMock)
		// with custom uniforms
		screen.SetUniformFloat("testFloat", float32(42.0))
		screen.SetUniformVector("testVector", mgl32.Vec3{1.0, 2.0, 3.0})
		screen.Draw(wrapperMock)
		setupFunc := func(w interfaces.GLWrapper) { w.Enable(glwrapper.DEPTH_TEST) }
		screen.Setup(setupFunc)
		screen.Draw(wrapperMock)
	}()
}
func TestUpdate(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		screen := New()
		// wo everything
		kst := store.NewGlfwKeyStore()
		bst := store.NewGlfwMouseStore()
		screen.Update(10, 0, 0, kst, bst)
		// with camera
		screen.SetCamera(cam)
		screen.Update(10, 0, 0, kst, bst)
		// with rotate on distance
		screen.SetRotateOnEdgeDistance(0.1)
		screen.Update(10, 0, 0, kst, bst)
		// with shader & mesh
		screen.AddShader(sm)
		mod := model.New()
		screen.AddModelToShader(mod, sm)
		msh := mesh.NewPointMesh(wrapperMock)
		mod.AddMesh(msh)
		screen.Update(10, 0, 0, kst, bst)
	}()
}
func TestCameraKeyboardMovement(t *testing.T) {
	screen := New()
	screen.SetCamera(cam)
	cmMap := make(map[string]glfw.Key)
	cmMap["forward"] = glfw.KeyW
	cmMap["back"] = glfw.KeyS
	screen.SetCameraMovementMap(cmMap)
	st := store.NewGlfwKeyStore()

	// first dir
	st.Set(glfw.KeyW, true)
	screen.cameraKeyboardMovement("forward", "back", "Lift", 10, st)

	// second dir
	st.Set(glfw.KeyW, false)
	st.Set(glfw.KeyS, true)
	screen.cameraKeyboardMovement("forward", "back", "Lift", 10, st)

	// with velocity
	cam.SetVelocity(10)
	screen.cameraKeyboardMovement("forward", "back", "Lift", 10, st)

	// Wrong handler name
	screen.cameraKeyboardMovement("forward", "back", "Wrong", 10, st)

}
func TestCameraKeyboardRotation(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
			}
		}()
		screen := New()
		screen.SetCamera(cam)
		keyMap := make(map[string]glfw.Key)
		keyMap["rotateUp"] = glfw.KeyW
		keyMap["rotateDown"] = glfw.KeyS
		keyMap["rotateLeft"] = glfw.KeyA
		keyMap["rotateRight"] = glfw.KeyD
		st := store.NewGlfwKeyStore()
		st.Set(glfw.KeyW, false)
		st.Set(glfw.KeyS, true)
		// wo keymap
		screen.cameraKeyboardRotation(10, st)
		// with keymap
		screen.SetCameraMovementMap(keyMap)
		screen.cameraKeyboardRotation(10, st)
	}()
}
func TestCameraMouseRotation(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Log(r)
				t.Error("Shouldn't have panic.")
			}
		}()
		screen := New()
		screen.SetCamera(cam)
		screen.cameraMouseRotation(10, 10, 10)
		screen.SetRotateOnEdgeDistance(1.0)
		screen.cameraMouseRotation(10, 30, 40)
		screen.cameraMouseRotation(10, -1, -1)
	}()
}
func TestCameraCollisionTest(t *testing.T) {
	screen := New()
	screen.AddShader(sm)
	mod := model.New()
	screen.AddModelToShader(mod, sm)
	bs := coldet.NewBoundingSphere([3]float32{0, 0, 0}, 1)
	result := screen.cameraCollisionTest(bs)
	if result != false {
		t.Error("Shouldn't collide.")
	}
	msh := mesh.NewPointMesh(wrapperMock)
	boparams := make(map[string]float32)
	boparams["radius"] = 1
	msh.SetPosition(mgl32.Vec3{0, 1, 0})
	msh.SetBoundingObject(boundingobject.New("Sphere", boparams))
	mod.AddMesh(msh)
	result = screen.cameraCollisionTest(bs)
	if result != true {
		t.Error("Should collide.")
	}
}
func TestExport(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestSetup(t *testing.T) {
	screen := New()
	setupFunc := func(w interfaces.GLWrapper) { w.Enable(glwrapper.DEPTH_TEST) }
	screen.Setup(setupFunc)

	if screen.setupFunction == nil {
		t.Errorf("Invalid function.")
	}

}
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
func newFormScreen() *FormScreen {
	builder := NewFormScreenBuilder()
	builder.SetWrapper(wrapperReal)
	builder.SetWindowSize(800, 800)
	return builder.Build()
}
func TestNewFormScreen(t *testing.T) {
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
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := newFormScreen()
		defer testhelper.GlfwTerminate()
		if len(form.formItemToConf) != 0 {
			t.Errorf("Invalid initial formItemToConf length. '%d'.", len(form.formItemToConf))
		}
	}()
}
func TestFormScreenUpdate(t *testing.T) {
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
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		defer testhelper.GlfwTerminate()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		builder.SetConfigOrder([]string{
			builder.AddConfigBool("label bool", DefaultFormItemDescription, true),
			builder.AddConfigInt("label int", DefaultFormItemDescription, 1, nil),
			builder.AddConfigInt64("label int64", DefaultFormItemDescription, 10, nil),
			builder.AddConfigFloat("label float", DefaultFormItemDescription, 0.44, nil),
			builder.AddConfigText("label text", DefaultFormItemDescription, "sample", nil),
			builder.AddConfigVector("label vector", DefaultFormItemDescription, mgl32.Vec3{0.01, 0.02, 0.03}, nil),
		})
		form := builder.Build()
		ks := store.NewGlfwKeyStore()
		ms := store.NewGlfwMouseStore()
		form.Update(10, 0.5, 0.5, ks, ms)
		form.Update(10, -0.4, 0.79, ks, ms)
		ms.Set(LEFT_MOUSE_BUTTON, true)
		form.sinceLastClick = 201
		form.Update(10, -0.4, 0.79, ks, ms)
		form.sinceLastClick = 201
		form.Update(10, 0.4, 0.79, ks, ms)
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		ks.Set(BACK_SPACE, true)
		form.Update(10, 0.4, 0.79, ks, ms)
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, 0.4, 0.69, ks, ms)
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, -0.4, 0.69, ks, ms)
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, -0.4, 0.59, ks, ms)
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, -0.4, 0.69, ks, ms)
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		form.Update(10, 0.2, 0.49, ks, ms)
	}()
}
func TestFormScreenCharCallback(t *testing.T) {
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
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		textKey := builder.AddConfigText("text", DefaultFormItemDescription, "", nil)
		intKey := builder.AddConfigInt("int", DefaultFormItemDescription, 1, nil)
		int64Key := builder.AddConfigInt64("int64", DefaultFormItemDescription, 2, nil)
		floatKey := builder.AddConfigFloat("float", DefaultFormItemDescription, 0.0, nil)
		builder.SetConfigOrder([]string{textKey, intKey, int64Key, floatKey})
		form := builder.Build()
		defer testhelper.GlfwTerminate()
		form.underEdit = form.GetFormItem(textKey).(*model.FormItemText)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "1" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.underEdit = form.GetFormItem(intKey).(*model.FormItemInt)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "11" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.underEdit = form.GetFormItem(int64Key).(*model.FormItemInt64)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "21" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.underEdit = form.GetFormItem(floatKey).(*model.FormItemFloat)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "0" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.CharCallback('.', wrapperReal)
		if form.underEdit.ValueToString() != "0." {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "0.1" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
	}()
}
func TestFormScreenSetup(t *testing.T) {
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
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		setup(wrapperReal)
	}()
}
func TestFormGetFormItemValidIndex(t *testing.T) {
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
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		defer testhelper.GlfwTerminate()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		textKey := builder.AddConfigText("text", DefaultFormItemDescription, "", nil)
		intKey := builder.AddConfigInt("int", DefaultFormItemDescription, 0, nil)
		int64Key := builder.AddConfigInt64("int64", DefaultFormItemDescription, 3, nil)
		floatKey := builder.AddConfigFloat("float", DefaultFormItemDescription, 0.0, nil)
		boolKey := builder.AddConfigBool("bool", DefaultFormItemDescription, false)
		vectorKey := builder.AddConfigVector("vector", DefaultFormItemDescription, mgl32.Vec3{0, 1, 0}, nil)
		builder.SetConfigOrder([]string{textKey, intKey, int64Key, floatKey, boolKey, vectorKey})
		form := builder.Build()
		fi := form.GetFormItem(textKey)
		if fi.ValueToString() != "" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(intKey)
		if fi.ValueToString() != "0" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(int64Key)
		if fi.ValueToString() != "3" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(floatKey)
		if fi.ValueToString() != "0" {
			t.Errorf("Invalid form item initial value. '%s'.", fi.ValueToString())
		}
		fi = form.GetFormItem(boolKey)
		if fi.ValueToString() != "false" {
			t.Error("Invalid form item initial value")
		}
	}()
}
func TestFormGetFormItemValidIndexValidators(t *testing.T) {
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
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		defer testhelper.GlfwTerminate()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		textKey := builder.AddConfigText("text", DefaultFormItemDescription, "", func(t string) bool { return true })
		intKey := builder.AddConfigInt("int", DefaultFormItemDescription, 0, func(i int) bool { return true })
		int64Key := builder.AddConfigInt64("int64", DefaultFormItemDescription, 3, func(i int64) bool { return true })
		floatKey := builder.AddConfigFloat("float", DefaultFormItemDescription, 0.0, func(f float32) bool { return true })
		boolKey := builder.AddConfigBool("bool", DefaultFormItemDescription, false)
		vectorKey := builder.AddConfigVector("vector", DefaultFormItemDescription, mgl32.Vec3{0, 1, 0}, func(f float32) bool { return true })
		builder.SetConfigOrder([]string{textKey, intKey, int64Key, floatKey, boolKey, vectorKey})
		form := builder.Build()
		fi := form.GetFormItem(textKey)
		if fi.ValueToString() != "" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(intKey)
		if fi.ValueToString() != "0" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(int64Key)
		if fi.ValueToString() != "3" {
			t.Error("Invalid form item initial value")
		}
		fi = form.GetFormItem(floatKey)
		if fi.ValueToString() != "0" {
			t.Errorf("Invalid form item initial value. '%s'.", fi.ValueToString())
		}
		fi = form.GetFormItem(boolKey)
		if fi.ValueToString() != "false" {
			t.Error("Invalid form item initial value")
		}
		_ = form.GetFormItem(vectorKey)
	}()
}
func TestFormGetFormItemInvalidIndex(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer testhelper.GlfwTerminate()
				t.Error("Should have panic.")
			}
		}()
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		builder := NewFormScreenBuilder()
		builder.SetWrapper(wrapperReal)
		builder.SetWindowSize(wW, wH)
		form := builder.Build()
		defer testhelper.GlfwTerminate()
		_ = form.GetFormItem("invalidkey")
	}()
}
func TestNewFormScreenBuilder(t *testing.T) {
	builder := NewFormScreenBuilder()
	if builder.headerLabel != "Default label" {
		t.Error("Invalid default label")
	}
	if builder.wrapper != nil {
		t.Error("Wrapper supposed to be nil by default")
	}
	if builder.charset != nil {
		t.Error("Charset supposed to be nil by default")
	}
	if builder.camera != nil {
		t.Error("camera supposed to be nil by default")
	}
}
func TestFormScreenBuilderBuild(t *testing.T) {
	t.Skip("Unimplemented")
}
func TestFormScreenBuilderSetHeaderLabel(t *testing.T) {
	builder := NewFormScreenBuilder()
	label := "new label"
	builder.SetHeaderLabel(label)
	if builder.headerLabel != label {
		t.Errorf("Invalid header label. Instead of '%s', we have '%s'.", label, builder.headerLabel)
	}
}
func TestFormScreenBuilderSetWrapper(t *testing.T) {
	builder := NewFormScreenBuilder()
	builder.SetWrapper(wrapperMock)
	if builder.wrapper != wrapperMock {
		t.Error("Invalid wrapper")
	}
}
func TestFormScreenBuilderSetWindowSize(t *testing.T) {
	builder := NewFormScreenBuilder()
	wW := float32(800)
	wH := float32(800)
	builder.SetWindowSize(wW, wH)
	if builder.windowWidth != wW {
		t.Errorf("Invalid window width. Instead of '%f', we have '%f'.", wW, builder.windowWidth)
	}
	if builder.windowHeight != wH {
		t.Errorf("Invalid window height. Instead of '%f', we have '%f'.", wH, builder.windowHeight)
	}
}
func TestFormScreenBuilderSetConfig(t *testing.T) {
	conf := config.New()
	builder := NewFormScreenBuilder()
	builder.SetConfig(conf)
	if !reflect.DeepEqual(builder.config, conf) {
		t.Error("Invalid configuration.")
	}
}
func TestFormScreenBuilderSetConfigOrder(t *testing.T) {
	builder := NewFormScreenBuilder()
	order := []string{"o1", "p1", "l1"}
	builder.SetConfigOrder(order)
	if !reflect.DeepEqual(builder.configOrder, order) {
		t.Error("Invalid order.")
	}
}
func TestFormScreenBuilderSetCharset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	runtime.LockOSThread()
	testhelper.GlfwInit()
	defer testhelper.GlfwTerminate()
	wrapperReal.InitOpenGL()
	builder := NewFormScreenBuilder()
	charset, err := model.LoadCharset("./assets/fonts/Desyrel/desyrel.ttf", 32, 127, 40.0, 72, wrapperReal)
	if err != nil {
		t.Errorf("Error during load charset: %#v.", err)
	}
	builder.SetCharset(charset)
	if builder.charset != charset {
		t.Error("Invalid charset.")
	}
}
func TestFormScreenBuilderAddConfigBool(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := true
	builder.AddConfigBool(confLabel, confDesc, value)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderAddConfigInt(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := 3
	builder.AddConfigInt(confLabel, confDesc, value, nil)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderAddConfigInt64(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := int64(3)
	builder.AddConfigInt64(confLabel, confDesc, value, nil)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderAddConfigFloat(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := float32(3.3)
	builder.AddConfigFloat(confLabel, confDesc, value, nil)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderAddConfigText(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := "test"
	builder.AddConfigText(confLabel, confDesc, value, nil)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
func TestFormScreenBuilderAddConfigVector(t *testing.T) {
	builder := NewFormScreenBuilder()
	if len(builder.config) != 0 {
		t.Error("Invalid initial config length.")
	}
	confLabel := "label"
	confDesc := "desc"
	value := mgl32.Vec3{0, 1, 0}
	builder.AddConfigVector(confLabel, confDesc, value, nil)
	if len(builder.config) != 1 {
		t.Error("Invalid config length.")
	}
}
