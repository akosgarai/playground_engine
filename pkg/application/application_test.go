package application

import (
	"reflect"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	wm          testhelper.WindowMock
	wrapperMock testhelper.GLWrapperMock
	cm          = camera.NewCamera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}, 0, 0)
	sm          testhelper.ShaderMock
)

func TestNew(t *testing.T) {
	app := New()
	if len(app.shaderMap) != 0 {
		t.Error("Invalid application - shadermap should be empty")
	}
	if app.cameraSet {
		t.Error("Camera shouldn't be set")
	}
}
func TestLog(t *testing.T) {
	app := New()
	log := app.Log()
	if len(log) < 10 {
		t.Error("Log too short.")
	}
	app.SetCamera(cm)
	log = app.Log()
	if len(log) < 10 {
		t.Error("Log too short.")
	}
}
func TestSetWindow(t *testing.T) {
	app := New()
	app.SetWindow(wm)

	if app.window != wm {
		t.Error("Invalid window setup.")
	}
}
func TestGetWindow(t *testing.T) {
	app := New()
	app.SetWindow(wm)

	if app.GetWindow() != wm {
		t.Error("Invalid window setup.")
	}
}
func TestSetCamera(t *testing.T) {
	app := New()
	app.SetCamera(cm)

	if app.camera != cm {
		t.Error("Invalid camera setup.")
	}
}
func TestGetCamera(t *testing.T) {
	app := New()
	app.SetCamera(cm)

	if app.GetCamera() != cm {
		t.Error("Invalid camera setup.")
	}
}
func TestSetKeyState(t *testing.T) {
	app := New()
	app.SetKeyState(glfw.KeyW, glfw.Press)
	if !app.keyDowns[glfw.KeyW] {
		t.Error("W should be pressed")
	}
	app.SetKeyState(glfw.KeyW, glfw.Release)
	if app.keyDowns[glfw.KeyW] {
		t.Error("W should be released")
	}
}
func TestSetButtonState(t *testing.T) {
	app := New()
	app.SetButtonState(glfw.MouseButtonLeft, glfw.Press)
	if !app.mouseDowns[glfw.MouseButtonLeft] {
		t.Error("LMB should be pressed")
	}
	app.SetButtonState(glfw.MouseButtonLeft, glfw.Release)
	if app.mouseDowns[glfw.MouseButtonLeft] {
		t.Error("LMB should be released")
	}
}
func TestGetMouseButtonState(t *testing.T) {
	app := New()
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
	app := New()
	app.SetKeyState(glfw.KeyW, glfw.Press)
	if !app.GetKeyState(glfw.KeyW) {
		t.Error("W should be pressed")
	}
	app.SetKeyState(glfw.KeyW, glfw.Release)
	if app.GetKeyState(glfw.KeyW) {
		t.Error("W should be released")
	}
}
func TestSetCameraMovementMap(t *testing.T) {
	app := New()
	cmMap := make(map[string]glfw.Key)
	app.SetCameraMovementMap(cmMap)

	if !reflect.DeepEqual(app.cameraKeyboardMovementMap, cmMap) {
		t.Error("Invalid camera movement map has been set.")
	}
}
func TestCameraKeyboardMovement(t *testing.T) {
	app := New()
	app.SetCamera(cm)
	cmMap := make(map[string]glfw.Key)
	cmMap["forward"] = glfw.KeyW
	cmMap["back"] = glfw.KeyS
	app.SetCameraMovementMap(cmMap)

	// first dir
	app.SetKeyState(glfw.KeyW, glfw.Press)
	app.cameraKeyboardMovement("forward", "back", "Lift", 10)

	// second dir
	app.SetKeyState(glfw.KeyW, glfw.Release)
	app.SetKeyState(glfw.KeyS, glfw.Press)
	app.cameraKeyboardMovement("forward", "back", "Lift", 10)

	// with velocity
	cm.SetVelocity(10)
	app.cameraKeyboardMovement("forward", "back", "Lift", 10)

	// Wrong handler name
	app.cameraKeyboardMovement("forward", "back", "Wrong", 10)

}
func TestAddModelToShader(t *testing.T) {
	app := New()
	app.AddShader(sm)
	if len(app.shaderMap[sm]) != 0 {
		t.Errorf("Invalid model length. Instead of '0', it is '%d'.\n", len(app.shaderMap[sm]))
	}
	mod := model.New()
	app.AddModelToShader(mod, sm)
	if len(app.shaderMap[sm]) != 1 {
		t.Errorf("Invalid model length. Instead of '1', it is '%d'.\n", len(app.shaderMap[sm]))
	}
}
func TestAddShader(t *testing.T) {
	app := New()
	if len(app.shaderMap) != 0 {
		t.Errorf("Invalid shader map length. Instead of '0', it is '%d'.\n", len(app.shaderMap))
	}
	app.AddShader(sm)
	if len(app.shaderMap) != 1 {
		t.Errorf("Invalid shader map length. Instead of '1', it is '%d'.\n", len(app.shaderMap))
	}
}
func TestSetRotateOnEdgeDistance(t *testing.T) {
	app := New()
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
		app.SetRotateOnEdgeDistance(tt.input)
		if app.rotateOnEdgeDistance != tt.expected {
			t.Errorf("Invalid rotateOnEdgeDistance. Instead of '%f', we have '%f'.", tt.expected, app.rotateOnEdgeDistance)
		}
	}
}
func TestCameraKeyboardRotation(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
			}
		}()
		app := New()
		app.SetCamera(cm)
		keyMap := make(map[string]glfw.Key)
		keyMap["rotateUp"] = glfw.KeyW
		keyMap["rotateDown"] = glfw.KeyS
		keyMap["rotateLeft"] = glfw.KeyA
		keyMap["rotateRight"] = glfw.KeyD
		// wo keymap
		app.cameraKeyboardRotation(10)
		// with keymap
		app.SetCameraMovementMap(keyMap)
		app.cameraKeyboardRotation(10)
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
		app := New()
		app.SetWindow(wm)
		app.SetCamera(cm)
		app.cameraMouseRotation(10)
		app.SetRotateOnEdgeDistance(1.0)
		wm.SetCursorPos(200, 200)
		app.SetWindow(wm)
		app.cameraMouseRotation(10)
		wm.SetCursorPos(600, 600)
		app.SetWindow(wm)
		app.cameraMouseRotation(10)
	}()
}
func TestApplyMouseRotation(t *testing.T) {
	func() {
		defer func() {
			cm.SetRotationStep(0)
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
			}
		}()
		app := New()
		cm.SetRotationStep(1)
		app.SetCamera(cm)
		testData := []struct {
			rL, rR, rU, rD bool
			delta          float64
		}{
			{true, true, true, true, 90},
			{true, false, true, true, 90},
			{false, true, true, true, 90},
			{false, false, true, true, 90},
			{false, false, false, true, 90},
			{false, false, false, false, 90},
			{false, false, true, false, 90},
		}
		for _, tt := range testData {
			app.applyMouseRotation(tt.rL, tt.rR, tt.rU, tt.rD, tt.delta)
		}
	}()
}
func TestCameraCollisionTest(t *testing.T) {
	app := New()
	app.AddShader(sm)
	mod := model.New()
	app.AddModelToShader(mod, sm)
	bs := coldet.NewBoundingSphere([3]float32{0, 0, 0}, 1)
	result := app.cameraCollisionTest(bs)
	if result != false {
		t.Error("Shouldn't collide.")
	}
	msh := mesh.NewPointMesh(wrapperMock)
	boparams := make(map[string]float32)
	boparams["radius"] = 1
	msh.SetPosition(mgl32.Vec3{0, 1, 0})
	msh.SetBoundingObject(boundingobject.New("Sphere", boparams))
	mod.AddMesh(msh)
	result = app.cameraCollisionTest(bs)
	if result != true {
		t.Error("Should collide.")
	}
}
func TestUpdate(t *testing.T) {
	func() {
		defer func() {
			cm.SetRotationStep(0)
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New()
		app.SetWindow(wm)
		// wo everything
		app.Update(10)
		// with camera
		app.SetCamera(cm)
		app.Update(10)
		// with rotate on distance
		app.SetRotateOnEdgeDistance(0.1)
		app.Update(10)
		// with shader & mesh
		app.AddShader(sm)
		mod := model.New()
		app.AddModelToShader(mod, sm)
		msh := mesh.NewPointMesh(wrapperMock)
		mod.AddMesh(msh)
		app.Update(10)
	}()
}
func TestAddDirectionalLightSource(t *testing.T) {
	func() {
		defer func() {
			cm.SetRotationStep(0)
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New()
		if len(app.directionalLightSources) != 0 {
			t.Errorf("Invalid number if dir. lenght. Instead of '0', it is '%d'.", len(app.directionalLightSources))
		}
		ds := light.NewDirectionalLight([4]mgl32.Vec3{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
		})
		uniforms := [4]string{"u1", "u2", "u3", "u4"}
		app.AddDirectionalLightSource(ds, uniforms)
		if len(app.directionalLightSources) != 1 {
			t.Errorf("Invalid number if dir. lenght. Instead of '1', it is '%d'.", len(app.directionalLightSources))
		}
		app.AddDirectionalLightSource(ds, uniforms)
		if len(app.directionalLightSources) != 2 {
			t.Errorf("Invalid number if dir. lenght. Instead of '2', it is '%d'.", len(app.directionalLightSources))
		}
	}()
}
func TestAddPointLightSource(t *testing.T) {
	func() {
		defer func() {
			cm.SetRotationStep(0)
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New()
		if len(app.pointLightSources) != 0 {
			t.Errorf("Invalid number if point. lenght. Instead of '0', it is '%d'.", len(app.pointLightSources))
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
		app.AddPointLightSource(ps, uniforms)
		if len(app.pointLightSources) != 1 {
			t.Errorf("Invalid number if point. lenght. Instead of '1', it is '%d'.", len(app.pointLightSources))
		}
		app.AddPointLightSource(ps, uniforms)
		if len(app.pointLightSources) != 2 {
			t.Errorf("Invalid number if point. lenght. Instead of '2', it is '%d'.", len(app.pointLightSources))
		}
	}()
}
func TestAddSpotLightSource(t *testing.T) {
	func() {
		defer func() {
			cm.SetRotationStep(0)
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New()
		if len(app.spotLightSources) != 0 {
			t.Errorf("Invalid number if spot. lenght. Instead of '0', it is '%d'.", len(app.spotLightSources))
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
		app.AddSpotLightSource(ss, uniforms)
		if len(app.spotLightSources) != 1 {
			t.Errorf("Invalid number if spot. lenght. Instead of '1', it is '%d'.", len(app.spotLightSources))
		}
		app.AddSpotLightSource(ss, uniforms)
		if len(app.spotLightSources) != 2 {
			t.Errorf("Invalid number if spot. lenght. Instead of '2', it is '%d'.", len(app.spotLightSources))
		}
	}()
}
func TestSetupDirectionalLightForShader(t *testing.T) {
	func() {
		defer func() {
			cm.SetRotationStep(0)
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New()
		ds := light.NewDirectionalLight([4]mgl32.Vec3{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
		})
		uniforms := [4]string{"u1", "u2", "u3", "u4"}
		app.AddDirectionalLightSource(ds, uniforms)
		if len(app.directionalLightSources) != 1 {
			t.Errorf("Invalid number if dir. lenght. Instead of '1', it is '%d'.", len(app.directionalLightSources))
		}
		app.setupDirectionalLightForShader(sm)
	}()
}
func TestSetupPointLightForShader(t *testing.T) {
	func() {
		defer func() {
			cm.SetRotationStep(0)
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New()
		ps := light.NewPointLight([4]mgl32.Vec3{
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
			mgl32.Vec3{0, 0, 0},
		},
			[3]float32{1, 0.5, 0.05},
		)
		uniforms := [7]string{"u1", "u2", "u3", "u4", "u5", "u6", "u7"}
		app.AddPointLightSource(ps, uniforms)
		app.setupPointLightForShader(sm)
	}()
}
func TestSetupSpotLightForShader(t *testing.T) {
	func() {
		defer func() {
			cm.SetRotationStep(0)
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New()
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
		app.AddSpotLightSource(ss, uniforms)
		app.setupSpotLightForShader(sm)
	}()
}
func TestDraw(t *testing.T) {
	func() {
		defer func() {
			cm.SetRotationStep(0)
			if r := recover(); r != nil {
				t.Error("Shouldn't have panic.")
				t.Log(r)
			}
		}()
		app := New()
		app.SetWindow(wm)
		app.AddShader(sm)
		mod := model.New()
		app.AddModelToShader(mod, sm)
		msh := mesh.NewPointMesh(wrapperMock)
		mod.AddMesh(msh)
		// wo camera
		app.Draw()
		// with camera
		app.SetCamera(cm)
		app.Draw()
	}()
}
