package screen

import (
	"reflect"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/store"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	cam = camera.NewCamera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}, 0, 0)

	sm          testhelper.ShaderMock
	wrapperMock testhelper.GLWrapperMock
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
func TestNewMenuScreenOption(t *testing.T) {
	t.Error("Unimplemented")
}
func TestOptionSetSurface(t *testing.T) {
	t.Error("Unimplemented")
}
func TestOptionDisplayCondition(t *testing.T) {
	t.Error("Unimplemented")
}
func TestNewMenuScreen(t *testing.T) {
	t.Error("Unimplemented")
}
func TestMenuScreenBuildScreen(t *testing.T) {
	t.Error("Unimplemented")
}
func TestMenuScreenAddOption(t *testing.T) {
	t.Error("Unimplemented")
}
func TestMenuScreenUpdate(t *testing.T) {
	t.Error("Unimplemented")
}
