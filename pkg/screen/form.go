package screen

import (
	"fmt"
	"math"

	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	FontFile          = "/assets/fonts/Frijole/Frijole-Regular.ttf"
	BottomFrameWidth  = float32(2.0)
	BottomFrameLength = float32(0.02)
	SideFrameLength   = float32(1.98)
	SideFrameWidth    = float32(0.02)
	TopLeftFrameWidth = float32(0.1)
	CameraMoveSpeed   = 0.005
)

var (
	DirectionalLightDirection = (mgl32.Vec3{0.0, 0.0, 1.0}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightDiffuse   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightSpecular  = mgl32.Vec3{1.0, 1.0, 1.0}
)

type FormScreen struct {
	*ScreenBase
	charset    *model.Charset
	background *model.BaseModel
	frame      *material.Material
	header     string
	formItems  []*model.FormItemBool
	bgShader   *shader.Shader
}

func frameRectangle(width, length float32, position mgl32.Vec3, mat *material.Material, wrapper interfaces.GLWrapper) *mesh.TexturedMaterialMesh {
	v, i, _ := rectangle.NewExact(width, length).MeshInput()
	var tex texture.Textures
	tex.TransparentTexture(1, 1, 128, "tex.diffuse", wrapper)
	tex.TransparentTexture(1, 1, 128, "tex.specular", wrapper)
	frameMesh := mesh.NewTexturedMaterialMesh(v, i, tex, mat, wrapper)
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
	wrapper.ClearColor(0.7, 0.7, 0.7, 1.0)
	wrapper.Enable(glwrapper.DEPTH_TEST)
	wrapper.DepthFunc(glwrapper.LESS)
	wrapper.Enable(glwrapper.BLEND)
	wrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
}

// NewFormScreen returns a FormScreen. The screen contains a material Frame.
func NewFormScreen(frame *material.Material, label string, wrapper interfaces.GLWrapper, wW, wH float32) *FormScreen {
	s := newScreenBase()
	s.SetCamera(createCamera(wW / wH))
	s.SetCameraMovementMap(cameraMovementMap())
	s.Setup(setup)
	bgShaderApplication := shader.NewTextureMatShaderBlending(wrapper)
	frameShaderApplication := shader.NewMaterialShader(wrapper)
	fgShaderApplication := shader.NewFontShader(wrapper)
	s.AddShader(bgShaderApplication)
	s.AddShader(fgShaderApplication)
	s.AddShader(frameShaderApplication)
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
	frameModel := model.New()
	// create frame.
	bottomFrame := frameRectangle(BottomFrameWidth, BottomFrameLength, mgl32.Vec3{0.0, -0.99, 0.0}, frame, wrapper)
	leftFrame := frameRectangle(SideFrameWidth, SideFrameLength, mgl32.Vec3{-0.99, 0.0, 0.0}, frame, wrapper)
	rightFrame := frameRectangle(SideFrameWidth, SideFrameLength, mgl32.Vec3{0.99, 0.0, 0.0}, frame, wrapper)
	topLeftFrame := frameRectangle(TopLeftFrameWidth, BottomFrameLength, mgl32.Vec3{0.95, 0.99, 0.0}, frame, wrapper)
	textWidth := chars.TextWidth("Settings", 3.0/wW)
	textContainer := frameRectangle(textWidth, 0.15, mgl32.Vec3{1 - TopLeftFrameWidth - textWidth/2, 0.925, 0}, frame, wrapper)
	textContainer.RotateX(-180)
	textContainer.RotateY(180)
	chars.PrintTo("Settings", -textWidth/2, -0.05, -0.01, 3.0/wW, wrapper, textContainer, []mgl32.Vec3{mgl32.Vec3{0, 0, 1}})
	topRightFrame := frameRectangle(2.0-TopLeftFrameWidth-textWidth, BottomFrameLength, mgl32.Vec3{(-TopLeftFrameWidth - textWidth) / 2, 0.99, 0.0}, frame, wrapper)
	frameModel.AddMesh(bottomFrame)
	frameModel.AddMesh(leftFrame)
	frameModel.AddMesh(rightFrame)
	frameModel.AddMesh(topLeftFrame)
	frameModel.AddMesh(topRightFrame)
	s.AddModelToShader(frameModel, frameShaderApplication)
	s.AddModelToShader(background, bgShaderApplication)
	return &FormScreen{
		ScreenBase: s,
		charset:    chars,
		background: background,
		bgShader:   bgShaderApplication,
		frame:      frame,
		header:     label,
	}
}

// Update loops on the shaderMap, and calls Update function on every Model.
// It also handles the camera movement and rotation, if the camera is set.
func (f *FormScreen) Update(dt, posX, posY float64, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
	cursorX := float32(posX)
	cursorY := float32(posY)
	if f.cameraSet {
		f.cameraKeyboardMovement("up", "down", "Lift", dt, keyStore)
		cursorX = cursorX + f.GetCamera().GetPosition().X()
		cursorY = cursorY + f.GetCamera().GetPosition().Y()
	}

	coords := mgl32.Vec3{cursorX, cursorY, 0.0}
	closestDistance := float32(math.MaxFloat32)
	var closestMesh interfaces.Mesh
	var closestModel interfaces.Model
	for s, _ := range f.shaderMap {
		for index, _ := range f.shaderMap[s] {
			f.shaderMap[s][index].Update(dt)
			msh, dist := f.shaderMap[s][index].ClosestMeshTo(coords)
			if dist < closestDistance {
				closestDistance = dist
				closestMesh = msh
				closestModel = f.shaderMap[s][index]
			}
		}
	}
	// Update the material in case of hover state.
	f.closestMesh = closestMesh
	f.closestDistance = closestDistance
	f.closestModel = closestModel

	switch f.closestModel.(type) {
	case *model.FormItemBool:
		tmMesh := f.closestMesh.(*mesh.TexturedMaterialMesh)
		tmMesh.Material = material.Whiteplastic
		minDiff := float32(0.0)
		if closestDistance <= minDiff+0.01 {
			fmt.Printf("closestDistance: %f, coords: %#v\n", closestDistance, coords)
			tmMesh.Material = material.Ruby
		}
		break
	}
}

// AddFormItemBool is for adding a bool form item to the form.
func (f *FormScreen) AddFormItemBool(formLabel string, wrapper interfaces.GLWrapper, wW float32) {
	// calculate the position of the option:
	// - bottom of the header: 0.85
	// - formItem: 0.1
	// - first form item Y: 0.80
	// - left col X: 0.49
	// - right col X: -0.49
	lenItems := len(f.formItems)
	posX := model.FormItemWidth / 2
	if lenItems%2 == 1 {
		posX = -1.0 * posX
	}
	posY := 0.80 - float32((lenItems/2))*0.1
	fi := model.NewFormItemBool(formLabel, material.Whiteplastic, mgl32.Vec3{posX, posY, 0}, wrapper)
	fi.RotateX(-90)
	fi.RotateY(180)
	f.AddModelToShader(fi, f.bgShader)
	f.charset.PrintTo(fi.GetLabel(), -0.48, -0.03, -0.01, 1.0/wW, wrapper, fi.GetSurface(), []mgl32.Vec3{mgl32.Vec3{0, 0, 1}})
	f.formItems = append(f.formItems, fi)
}
