package screen

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/pointer"
	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/store"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

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
	if shader.camera != nil {
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
	screen.camera = cam
	log = screen.Log()
	if log == emptylog {
		t.Errorf("Invalid log for camera. We have the same as before '%s'.", emptylog)
	}
}
func TestSetupCameraGoodOptionsWithoutMovement(t *testing.T) {
	var screen Screen
	goodOptionsForDefaultCamera := map[string]interface{}{
		"mode":                 "default",
		"rotateOnEdgeDistance": float32(0.5),
	}
	screen.SetupCamera(cam, goodOptionsForDefaultCamera)
	if screen.camera != cam {
		t.Error("Invalid camera setup.")
	}
	if screen.cameraMode != goodOptionsForDefaultCamera["mode"].(string) {
		t.Error("Invalid camera mode.")
	}
	if screen.rotateOnEdgeDistance != goodOptionsForDefaultCamera["rotateOnEdgeDistance"].(float32) {
		t.Error("Invalid rotateOnEdgeDistance.")
	}
	if len(screen.cameraKeyboardMovementMap) != 0 {
		t.Error("Movement should not be set.")
	}
}
func TestSetupCameraGoodOptionsWithMovement(t *testing.T) {
	var screen Screen
	screen.cameraKeyboardMovementMap = make(map[string][]glfw.Key)
	goodOptionsForDefaultCamera := map[string]interface{}{
		"mode":                 "default",
		"rotateOnEdgeDistance": float32(0.5),
		"forward":              []glfw.Key{glfw.KeyW, glfw.KeyI},
	}
	screen.SetupCamera(cam, goodOptionsForDefaultCamera)
	if screen.camera != cam {
		t.Error("Invalid camera setup.")
	}
	if screen.cameraMode != goodOptionsForDefaultCamera["mode"].(string) {
		t.Error("Invalid camera mode.")
	}
	if screen.rotateOnEdgeDistance != goodOptionsForDefaultCamera["rotateOnEdgeDistance"].(float32) {
		t.Error("Invalid rotateOnEdgeDistance.")
	}
	if len(screen.cameraKeyboardMovementMap) != 1 {
		t.Error("Movement should not be set.")
	}
}
func TestSetupCameraMissingMode(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing mode key in options.")
			}
		}()
		var screen Screen
		screen.cameraKeyboardMovementMap = make(map[string][]glfw.Key)
		optionsForDefaultCamera := map[string]interface{}{
			"rotateOnEdgeDistance": float32(0.5),
			"forward":              []glfw.Key{glfw.KeyW, glfw.KeyI},
		}
		screen.SetupCamera(cam, optionsForDefaultCamera)
	}()
}
func TestSetupCameraInvalidMode(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the invalid mode key in options.")
			}
		}()
		var screen Screen
		screen.cameraKeyboardMovementMap = make(map[string][]glfw.Key)
		optionsForDefaultCamera := map[string]interface{}{
			"mode":                 "wrongValue",
			"rotateOnEdgeDistance": float32(0.5),
			"forward":              []glfw.Key{glfw.KeyW, glfw.KeyI},
		}
		screen.SetupCamera(cam, optionsForDefaultCamera)
	}()
}
func TestSetupCameraDefaultModeMissingRotateValue(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing rotateOnEdgeDistance key in options.")
			}
		}()
		var screen Screen
		screen.cameraKeyboardMovementMap = make(map[string][]glfw.Key)
		optionsForDefaultCamera := map[string]interface{}{
			"mode":    "default",
			"forward": []glfw.Key{glfw.KeyW, glfw.KeyI},
		}
		screen.SetupCamera(cam, optionsForDefaultCamera)
	}()
}
func TestGetCamera(t *testing.T) {
	var screen Screen
	screen.cameraKeyboardMovementMap = make(map[string][]glfw.Key)
	optionsForDefaultCamera := map[string]interface{}{
		"mode":                 "default",
		"rotateOnEdgeDistance": float32(0.5),
		"forward":              []glfw.Key{glfw.KeyW, glfw.KeyI},
	}
	screen.SetupCamera(cam, optionsForDefaultCamera)

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
		optionsForDefaultCamera := map[string]interface{}{
			"mode":                 "default",
			"rotateOnEdgeDistance": float32(0.5),
			"forward":              []glfw.Key{glfw.KeyW, glfw.KeyI},
		}
		screen.SetupCamera(cam, optionsForDefaultCamera)
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
		screen.Update(10, pointer.New(0.0, 0.0, 0.0, 0.0), kst, bst)
		// with camera
		optionsForDefaultCamera := map[string]interface{}{
			"mode":                 "default",
			"rotateOnEdgeDistance": float32(0.1),
			"forward":              []glfw.Key{glfw.KeyW, glfw.KeyI},
		}
		screen.SetupCamera(cam, optionsForDefaultCamera)
		screen.Update(10, pointer.New(0.0, 0.0, 0.0, 0.0), kst, bst)
		// with shader & mesh
		screen.AddShader(sm)
		mod := model.New()
		screen.AddModelToShader(mod, sm)
		msh := mesh.NewPointMesh(wrapperMock)
		mod.AddMesh(msh)
		screen.Update(10, pointer.New(0.0, 0.0, 0.0, 0.0), kst, bst)
	}()
}
func TestCameraKeyboardMovement(t *testing.T) {
	screen := New()
	optionsForDefaultCamera := map[string]interface{}{
		"mode":                 "default",
		"rotateOnEdgeDistance": float32(0.1),
		"forward":              []glfw.Key{glfw.KeyW},
		"back":                 []glfw.Key{glfw.KeyS},
	}
	screen.SetupCamera(cam, optionsForDefaultCamera)
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
func TestCameraKeyboardRotationWithoutKeymap(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
			}
		}()
		screen := New()
		optionsForDefaultCamera := map[string]interface{}{
			"mode":                 "default",
			"rotateOnEdgeDistance": float32(0.1),
		}
		screen.SetupCamera(cam, optionsForDefaultCamera)
		st := store.NewGlfwKeyStore()
		st.Set(glfw.KeyW, false)
		st.Set(glfw.KeyS, true)
		// with keymap
		screen.cameraKeyboardRotation(10, st)
	}()
}
func TestCameraKeyboardRotationWithKeymap(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
			}
		}()
		screen := New()
		optionsForDefaultCamera := map[string]interface{}{
			"mode":                 "default",
			"rotateOnEdgeDistance": float32(0.1),
			"rotateUp":             []glfw.Key{glfw.KeyW},
			"rotateDown":           []glfw.Key{glfw.KeyS},
			"rotateLeft":           []glfw.Key{glfw.KeyA},
			"rotateRight":          []glfw.Key{glfw.KeyD},
		}
		screen.SetupCamera(cam, optionsForDefaultCamera)
		st := store.NewGlfwKeyStore()
		// down
		st.Set(glfw.KeyW, false)
		st.Set(glfw.KeyS, true)
		screen.cameraKeyboardRotation(10, st)
		// up
		st.Set(glfw.KeyS, false)
		st.Set(glfw.KeyW, true)
		screen.cameraKeyboardRotation(10, st)
		// left
		st.Set(glfw.KeyW, false)
		st.Set(glfw.KeyD, false)
		st.Set(glfw.KeyA, true)
		screen.cameraKeyboardRotation(10, st)
		// right
		st.Set(glfw.KeyD, true)
		st.Set(glfw.KeyA, false)
		screen.cameraKeyboardRotation(10, st)
	}()
}
func TestCameraMouseRotationWithoutDistance(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Log(r)
				t.Error("Shouldn't have panic.")
			}
		}()
		screen := New()
		optionsForDefaultCamera := map[string]interface{}{
			"mode":                 "default",
			"rotateOnEdgeDistance": float32(0.0),
		}
		screen.SetupCamera(cam, optionsForDefaultCamera)
		screen.cameraMouseRotationDefault(10, 10, 10)
	}()
}
func TestCameraMouseRotationWithDistance(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Log(r)
				t.Error("Shouldn't have panic.")
			}
		}()
		screen := New()
		optionsForDefaultCamera := map[string]interface{}{
			"mode":                 "default",
			"rotateOnEdgeDistance": float32(1.0),
		}
		screen.SetupCamera(cam, optionsForDefaultCamera)
		screen.cameraMouseRotationDefault(10, 30, 40)
		screen.cameraMouseRotationDefault(10, -1, -1)
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
func TestCleanSpotLightSource(t *testing.T) {
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
		screen.CleanSpotLightSources()
		if len(screen.spotLightSources) != 0 {
			t.Errorf("Invalid number if spot. length. Instead of '0', it is '%d'.", len(screen.spotLightSources))
		}
	}()
}
func TestCleanPointLightSource(t *testing.T) {
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
		screen.CleanPointLightSources()
		if len(screen.pointLightSources) != 0 {
			t.Errorf("Invalid number if point. length. Instead of '0', it is '%d'.", len(screen.pointLightSources))
		}
	}()
}
func TestSetWindowSize(t *testing.T) {
	var screen Screen
	testData := []struct {
		x float32
		y float32
	}{
		{10.0, 10.0},
		{100.0, 100.0},
		{600.0, 600.0},
	}
	for _, tt := range testData {
		screen.SetWindowSize(tt.x, tt.y)
		if screen.windowWidth != tt.x {
			t.Errorf("Invalid windowWidth. Instead of '%f', we have '%f'.", tt.x, screen.windowWidth)
		}
		if screen.windowHeight != tt.y {
			t.Errorf("Invalid windowHeight. Instead of '%f', we have '%f'.", tt.y, screen.windowHeight)
		}
	}
}
func TestGetWindowSize(t *testing.T) {
	var screen Screen
	testData := []struct {
		x float32
		y float32
	}{
		{10.0, 10.0},
		{100.0, 100.0},
		{600.0, 600.0},
	}
	for _, tt := range testData {
		screen.SetWindowSize(tt.x, tt.y)
		x, y := screen.GetWindowSize()
		if x != tt.x {
			t.Errorf("Invalid windowWidth. Instead of '%f', we have '%f'.", tt.x, x)
		}
		if y != tt.y {
			t.Errorf("Invalid windowHeight. Instead of '%f', we have '%f'.", tt.y, y)
		}
	}
}
func TestGetAspectRatio(t *testing.T) {
	var screen Screen
	testData := []struct {
		x      float32
		y      float32
		aspect float32
	}{
		{10.0, 10.0, 1.0},
		{100.0, 100.0, 1.0},
		{600.0, 600.0, 1.0},
		{1000.0, 500.0, 2.0},
		{500.0, 1000.0, 0.5},
	}
	for _, tt := range testData {
		screen.SetWindowSize(tt.x, tt.y)
		aspect := screen.GetAspectRatio()
		if aspect != tt.aspect {
			t.Errorf("Invalid aspect. Instead of '%f', we have '%f'.", tt.aspect, aspect)
		}
	}
}
