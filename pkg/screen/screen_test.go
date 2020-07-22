package screen

import (
	"reflect"
	"runtime"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/camera"
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		if form.header != screenLabel {
			t.Errorf("Invalid header. Instead of '%s', we have '%s'.", screenLabel, form.header)
		}
		if form.frame != frameMat {
			t.Error("Invalid material.")
		}
		if len(form.configuration) != 0 {
			t.Errorf("Invalid initial configuration length. '%d'.", len(form.configuration))
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		ks := store.NewGlfwKeyStore()
		ms := store.NewGlfwMouseStore()
		form.Update(10, 0.5, 0.5, ks, ms)
		// add option
		index := form.AddFormItemBool("label bool", DefaultFormItemDescription, wrapperReal, true)
		form.Update(10, -0.4, 0.79, ks, ms)
		form.AddFormItemInt("label int", DefaultFormItemDescription, wrapperReal, 1, nil)
		form.Update(10, -0.4, 0.79, ks, ms)
		ms.Set(LEFT_MOUSE_BUTTON, true)
		form.sinceLastClick = 201
		form.Update(10, -0.4, 0.79, ks, ms)
		if form.GetFormItem(index).(*model.FormItemBool).GetValue() != false {
			t.Error("FormItemBool value should be toggled.")
		}
		form.sinceLastClick = 201
		form.Update(10, 0.4, 0.79, ks, ms)
		form.sinceLastClick = 201
		form.sinceLastDelete = 201
		ks.Set(BACK_SPACE, true)
		form.Update(10, 0.4, 0.79, ks, ms)
		// further options
		form.AddFormItemInt64("label int64", DefaultFormItemDescription, wrapperReal, 10, nil)
		form.AddFormItemFloat("label float", DefaultFormItemDescription, wrapperReal, 0.44, nil)
		form.AddFormItemText("label text", DefaultFormItemDescription, wrapperReal, "sample", nil)
		form.AddFormItemVector("label vector", DefaultFormItemDescription, wrapperReal, mgl32.Vec3{0.01, 0.02, 0.03}, nil)
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
func TestFormScreenAddFormItemBool(t *testing.T) {
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		labels := []string{"label1", "label2", "label3", "label4"}
		for i := 0; i < len(labels); i++ {
			index := form.AddFormItemBool(labels[i], DefaultFormItemDescription, wrapperReal, true)
			if index != i {
				t.Error("Invalid index.")
			}
			if len(form.formItems) != i+1 {
				t.Error("Invalid number of form items.")
			}
		}
	}()
}
func TestFormScreenAddFormItemInt(t *testing.T) {
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		labels := []string{"label1", "label2", "label3", "label4"}
		for i := 0; i < len(labels); i++ {
			index := form.AddFormItemInt(labels[i], DefaultFormItemDescription, wrapperReal, 1, nil)
			if index != i {
				t.Error("Invalid index.")
			}
			if len(form.formItems) != i+1 {
				t.Error("Invalid number of form items.")
			}
		}
	}()
}
func TestFormScreenAddFormItemFloat(t *testing.T) {
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		labels := []string{"label1", "label2", "label3", "label4"}
		for i := 0; i < len(labels); i++ {
			index := form.AddFormItemFloat(labels[i], DefaultFormItemDescription, wrapperReal, 0.2, nil)
			if index != i {
				t.Error("Invalid index.")
			}
			if len(form.formItems) != i+1 {
				t.Error("Invalid number of form items.")
			}
		}
	}()
}
func TestFormScreenAddFormItemInt64(t *testing.T) {
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		labels := []string{"label1", "label2", "label3", "label4"}
		for i := 0; i < len(labels); i++ {
			index := form.AddFormItemInt64(labels[i], DefaultFormItemDescription, wrapperReal, 33366699900, nil)
			if index != i {
				t.Error("Invalid index.")
			}
			if len(form.formItems) != i+1 {
				t.Error("Invalid number of form items.")
			}
		}
	}()
}
func TestFormScreenAddFormItemText(t *testing.T) {
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		labels := []string{"label1", "label2", "label3", "label4"}
		for i := 0; i < len(labels); i++ {
			index := form.AddFormItemText(labels[i], DefaultFormItemDescription, wrapperReal, "sample text", nil)
			if index != i {
				t.Error("Invalid index.")
			}
			if len(form.formItems) != i+1 {
				t.Error("Invalid number of form items.")
			}
		}
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		form.AddFormItemText("text", DefaultFormItemDescription, wrapperReal, "", nil)
		form.AddFormItemInt("int", DefaultFormItemDescription, wrapperReal, 1, nil)
		form.AddFormItemInt64("int64", DefaultFormItemDescription, wrapperReal, 2, nil)
		form.AddFormItemFloat("float", DefaultFormItemDescription, wrapperReal, 0.0, nil)
		form.underEdit = form.formItems[0].(*model.FormItemText)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "1" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.underEdit = form.formItems[1].(*model.FormItemInt)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "11" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.underEdit = form.formItems[2].(*model.FormItemInt64)
		form.CharCallback('1', wrapperReal)
		if form.underEdit.ValueToString() != "21" {
			t.Errorf("Invalid value: '%s'.", form.underEdit.ValueToString())
		}
		form.underEdit = form.formItems[3].(*model.FormItemFloat)
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		index := form.AddFormItemText("text", DefaultFormItemDescription, wrapperReal, "", nil)
		fi := form.GetFormItem(index)
		if fi.ValueToString() != "" {
			t.Error("Invalid form item initial value")
		}
		index = form.AddFormItemInt("int", DefaultFormItemDescription, wrapperReal, 0, nil)
		fi = form.GetFormItem(index)
		if fi.ValueToString() != "0" {
			t.Error("Invalid form item initial value")
		}
		index = form.AddFormItemInt64("int64", DefaultFormItemDescription, wrapperReal, 3, nil)
		fi = form.GetFormItem(index)
		if fi.ValueToString() != "3" {
			t.Error("Invalid form item initial value")
		}
		index = form.AddFormItemFloat("float", DefaultFormItemDescription, wrapperReal, 0.0, nil)
		fi = form.GetFormItem(index)
		if fi.ValueToString() != "0" {
			t.Errorf("Invalid form item initial value. '%s'.", fi.ValueToString())
		}
		index = form.AddFormItemBool("bool", DefaultFormItemDescription, wrapperReal, false)
		fi = form.GetFormItem(index)
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		index := form.AddFormItemText("text", DefaultFormItemDescription, wrapperReal, "", func(t string) bool { return true })
		fi := form.GetFormItem(index)
		if fi.ValueToString() != "" {
			t.Error("Invalid form item initial value")
		}
		index = form.AddFormItemInt("int", DefaultFormItemDescription, wrapperReal, 0, func(i int) bool { return true })
		fi = form.GetFormItem(index)
		if fi.ValueToString() != "0" {
			t.Error("Invalid form item initial value")
		}
		index = form.AddFormItemInt64("int64", DefaultFormItemDescription, wrapperReal, 3, func(i int64) bool { return true })
		fi = form.GetFormItem(index)
		if fi.ValueToString() != "3" {
			t.Error("Invalid form item initial value")
		}
		index = form.AddFormItemFloat("float", DefaultFormItemDescription, wrapperReal, 0.0, func(f float32) bool { return true })
		fi = form.GetFormItem(index)
		if fi.ValueToString() != "0" {
			t.Errorf("Invalid form item initial value. '%s'.", fi.ValueToString())
		}
		index = form.AddFormItemBool("bool", DefaultFormItemDescription, wrapperReal, false)
		fi = form.GetFormItem(index)
		if fi.ValueToString() != "false" {
			t.Error("Invalid form item initial value")
		}
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
		frameMat := material.Chrome
		screenLabel := "test-label"
		wW := float32(800)
		wH := float32(800)
		runtime.LockOSThread()
		testhelper.GlfwInit()
		wrapperReal.InitOpenGL()
		form := NewFormScreen(frameMat, screenLabel, wrapperReal, wW, wH)
		defer testhelper.GlfwTerminate()
		index := form.AddFormItemText("text", DefaultFormItemDescription, wrapperReal, "txt", nil)
		_ = form.GetFormItem(index + 2)
	}()
}
