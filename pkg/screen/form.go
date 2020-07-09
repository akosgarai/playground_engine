package screen

import (
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/shader"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	FontFile          = "/assets/fonts/Frijole/Frijole-Regular.ttf"
	BottomFrameWidth  = float32(2.0)
	BottomFrameLength = float32(0.02)
	SideFrameLength   = float32(1.98)
	SideFrameWidth    = float32(0.02)
	CameraMoveSpeed   = 0.005
)

var (
	DirectionalLightDirection = (mgl32.Vec3{0.0, 0.0, -1.0}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightDiffuse   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightSpecular  = mgl32.Vec3{0.5, 0.5, 0.5}
	WhiteMaterial             = material.New(
		mgl32.Vec3{1, 1, 1},
		mgl32.Vec3{1, 1, 1},
		mgl32.Vec3{1, 1, 1},
		36)
)

type FormScreen struct {
	*ScreenBase
	charset *model.Charset
	frame   *material.Material
	header  string
}

func frameRectangle(width, length float32, position mgl32.Vec3, mat *material.Material, wrapper interfaces.GLWrapper) *mesh.MaterialMesh {
	v, i, _ := rectangle.NewExact(width, length).MeshInput()
	frameMesh := mesh.NewMaterialMesh(v, i, mat, wrapper)
	frameMesh.RotateX(90)
	frameMesh.SetPosition(position)
	return frameMesh
}
func charset(wrapper interfaces.GLWrapper) *model.Charset {
	cs, err := model.LoadCharset(baseDirScreen()+FontFile, 32, 127, 40.0, 72, wrapper)
	if err != nil {
		panic(err)
	}
	cs.SetTransparent(true)
	return cs
}

// It creates a new camera with the necessary setup
func createCamera(ratio float32) *camera.Camera {
	camera := camera.NewCamera(mgl32.Vec3{0, 0, -1.8}, mgl32.Vec3{0, -1, 0}, 90.0, 0.0)
	camera.SetupProjection(45, ratio, 0.001, 10.0)
	camera.SetVelocity(CameraMoveSpeed)
	return camera
}

// Setup keymap for the camera movement
func cameraMovementMap() map[string]glfw.Key {
	cm := make(map[string]glfw.Key)
	cm["up"] = glfw.KeyUp
	cm["down"] = glfw.KeyDown
	return cm
}
func setup(wrapper interfaces.GLWrapper) {
	wrapper.ClearColor(1.0, 1.0, 1.0, 1.0)
	wrapper.Enable(glwrapper.DEPTH_TEST)
	wrapper.DepthFunc(glwrapper.LESS)
}

// NewFormScreen returns a FormScreen. The screen contains a material Frame.
func NewFormScreen(frame *material.Material, label string, wrapper interfaces.GLWrapper, wW, wH float32) *FormScreen {
	s := newScreenBase()
	s.SetCamera(createCamera(wW / wH))
	s.SetCameraMovementMap(cameraMovementMap())
	s.Setup(setup)
	bgShaderApplication := shader.NewMaterialShader(wrapper)
	fgShaderApplication := shader.NewFontShader(wrapper)
	s.AddShader(bgShaderApplication)
	s.AddShader(fgShaderApplication)
	LightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
	})
	s.AddDirectionalLightSource(LightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})

	chars := charset(wrapper)
	s.AddModelToShader(chars, fgShaderApplication)
	background := model.New()
	// create frame here.
	bottomFrame := frameRectangle(BottomFrameWidth, BottomFrameLength, mgl32.Vec3{0.0, -0.99, 0.0}, material.Chrome, wrapper)
	leftFrame := frameRectangle(SideFrameWidth, SideFrameLength, mgl32.Vec3{-0.99, 0.0, 0.0}, material.Chrome, wrapper)
	rightFrame := frameRectangle(SideFrameWidth, SideFrameLength, mgl32.Vec3{0.99, 0.0, 0.0}, material.Chrome, wrapper)
	textContainer := frameRectangle(chars.TextWidth("Options", 3.0/wW), 0.2, mgl32.Vec3{0, 0, 0}, material.Chrome, wrapper)
	background.AddMesh(bottomFrame)
	background.AddMesh(leftFrame)
	background.AddMesh(rightFrame)
	background.AddMesh(textContainer)
	s.AddModelToShader(background, bgShaderApplication)
	return &FormScreen{
		ScreenBase: s,
		charset:    chars,
		frame:      frame,
		header:     label,
	}
}

// Update loops on the shaderMap, and calls Update function on every Model.
// It also handles the camera movement and rotation, if the camera is set.
func (f *FormScreen) Update(dt, posX, posY float64, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
	if f.cameraSet {
		f.cameraKeyboardMovement("up", "down", "Lift", dt, keyStore)
	}
}
