package screen

import (
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
	FontFile           = "/assets/fonts/Frijole/Frijole-Regular.ttf"
	BottomFrameWidth   = float32(2.0)
	BottomFrameLength  = float32(0.02)
	SideFrameLength    = float32(1.98)
	SideFrameWidth     = float32(0.02)
	TopLeftFrameWidth  = float32(0.1)
	CameraMoveSpeed    = 0.005
	LightConstantTerm  = float32(1.0)
	LightLinearTerm    = float32(0.14)
	LightQuadraticTerm = float32(0.07)
	EventEpsilon       = 200
	BACK_SPACE         = glfw.KeyBackspace
)

var (
	DirectionalLightDirection = (mgl32.Vec3{0.0, 0.0, 1.0}).Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightDiffuse   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightSpecular  = mgl32.Vec3{1.0, 1.0, 1.0}

	DefaultFormItemMaterial = material.Whiteplastic
	SpotLightAmbient        = mgl32.Vec3{1, 1, 1}
	SpotLightDiffuse        = mgl32.Vec3{1, 1, 1}
	SpotLightSpecular       = mgl32.Vec3{1, 1, 1}
	SpotLightDirection      = (mgl32.Vec3{0, 0, -1}).Normalize()
	SpotLightCutoff         = float32(0.05)
	SpotLightOuterCutoff    = float32(0.07)
)

type FormScreen struct {
	*ScreenBase
	charset         *model.Charset
	background      *model.BaseModel
	frame           *material.Material
	header          string
	formItems       []interfaces.FormItem
	bgShader        *shader.Shader
	sinceLastClick  float64
	sinceLastDelete float64
	underEdit       interfaces.CharFormItem
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
	s.SetWindowSize(wW, wH)
	s.SetWrapper(wrapper)
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
		ScreenBase:     s,
		charset:        chars,
		background:     background,
		bgShader:       bgShaderApplication,
		frame:          frame,
		header:         label,
		sinceLastClick: 0,
	}
}

// initMaterialForTheFormItems sets the material to the default of the form items.
// It could be used un the update loop to make all of them to default state.
func (f *FormScreen) initMaterialForTheFormItems() {
	for s, _ := range f.shaderMap {
		for index, _ := range f.shaderMap[s] {
			switch f.shaderMap[s][index].(type) {
			case *model.FormItemInt:
				fi := f.shaderMap[s][index].(*model.FormItemInt)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = DefaultFormItemMaterial
				break
			case *model.FormItemBool:
				fi := f.shaderMap[s][index].(*model.FormItemBool)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = DefaultFormItemMaterial
				lightMesh := fi.GetLight().(*mesh.TexturedMaterialMesh)
				if fi.GetValue() {
					lightMesh.Material = material.Ruby
				} else {
					lightMesh.Material = material.Whiteplastic
				}
				break
			}
		}
	}
}

// deleteCursor removes the cursor from the text form inputs.
func (f *FormScreen) deleteCursor() {
	switch f.underEdit.(type) {
	case *model.FormItemInt:
		fi := f.underEdit.(*model.FormItemInt)
		fi.DeleteCursor()
		f.underEdit = nil
		break
	}
}

// Update loops on the shaderMap, and calls Update function on every Model.
// It also handles the camera movement and rotation, if the camera is set.
func (f *FormScreen) Update(dt, posX, posY float64, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
	f.sinceLastClick = f.sinceLastClick + dt
	f.sinceLastDelete = f.sinceLastDelete + dt
	cursorX := float32(-posX)
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

	f.initMaterialForTheFormItems()
	switch f.closestModel.(type) {
	case *model.FormItemBool:
		tmMesh := f.closestMesh.(*mesh.TexturedMaterialMesh)
		minDiff := float32(0.0)
		if closestDistance <= minDiff+0.01 {
			tmMesh.Material = material.Ruby
			if buttonStore.Get(LEFT_MOUSE_BUTTON) {
				formModel := f.closestModel.(*model.FormItemBool)
				if f.sinceLastClick > EventEpsilon {
					f.deleteCursor()
					formModel.SetValue(!formModel.GetValue())
					f.sinceLastClick = 0
				}
			}
		}
		break
	case *model.FormItemInt:
		tmMesh := f.closestMesh.(*mesh.TexturedMaterialMesh)
		minDiff := float32(0.0)
		if closestDistance <= minDiff+0.01 {
			tmMesh.Material = material.Ruby
			if buttonStore.Get(LEFT_MOUSE_BUTTON) {
				formModel := f.closestModel.(*model.FormItemInt)
				if f.sinceLastClick > EventEpsilon {
					f.deleteCursor()
					formModel.AddCursor()
					f.sinceLastClick = 0
					f.underEdit = formModel
				}
			}
		}
		break
	case *model.FormItemFloat:
		tmMesh := f.closestMesh.(*mesh.TexturedMaterialMesh)
		minDiff := float32(0.0)
		if closestDistance <= minDiff+0.01 {
			tmMesh.Material = material.Ruby
			if buttonStore.Get(LEFT_MOUSE_BUTTON) {
				formModel := f.closestModel.(*model.FormItemFloat)
				if f.sinceLastClick > EventEpsilon {
					f.deleteCursor()
					formModel.AddCursor()
					f.sinceLastClick = 0
					f.underEdit = formModel
				}
			}
		}
		break
	}
	if keyStore.Get(BACK_SPACE) {
		if f.sinceLastDelete > EventEpsilon {
			f.underEdit.DeleteLastCharacter()
			f.sinceLastDelete = 0
			f.charset.CleanSurface(f.underEdit.GetTarget())
			f.charset.PrintTo(f.underEdit.ValueToString(), -model.CursorInitX, -0.015, -0.01, 1.0/f.windowWindth, f.wrapper, f.underEdit.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
		}
	}
}

// AddFormItemBool is for adding a bool form item to the form.
func (f *FormScreen) AddFormItemBool(formLabel string, wrapper interfaces.GLWrapper) {
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
	f.charset.PrintTo(fi.GetLabel(), -0.48, -0.03, -0.01, 1.0/f.windowWindth, wrapper, fi.GetSurface(), []mgl32.Vec3{mgl32.Vec3{0, 0, 1}})
	f.formItems = append(f.formItems, fi)
}

// AddFormItemInt is for adding an integer form item to the form.
func (f *FormScreen) AddFormItemInt(formLabel string, wrapper interfaces.GLWrapper) {
	lenItems := len(f.formItems)
	posX := model.FormItemWidth / 2
	if lenItems%2 == 1 {
		posX = -1.0 * posX
	}
	posY := 0.80 - float32((lenItems/2))*0.1
	fi := model.NewFormItemInt(formLabel, material.Whiteplastic, mgl32.Vec3{posX, posY, 0}, wrapper)
	fi.RotateX(-90)
	fi.RotateY(180)
	f.AddModelToShader(fi, f.bgShader)
	f.charset.PrintTo(fi.GetLabel(), -0.48, -0.03, -0.01, 1.0/f.windowWindth, wrapper, fi.GetSurface(), []mgl32.Vec3{mgl32.Vec3{0, 0, 1}})
	f.formItems = append(f.formItems, fi)
}

// AddFormItemFloat is for adding a float form item to the form.
func (f *FormScreen) AddFormItemFloat(formLabel string, wrapper interfaces.GLWrapper) {
	lenItems := len(f.formItems)
	posX := model.FormItemWidth / 2
	if lenItems%2 == 1 {
		posX = -1.0 * posX
	}
	posY := 0.80 - float32((lenItems/2))*0.1
	fi := model.NewFormItemFloat(formLabel, material.Whiteplastic, mgl32.Vec3{posX, posY, 0}, wrapper)
	fi.RotateX(-90)
	fi.RotateY(180)
	f.AddModelToShader(fi, f.bgShader)
	f.charset.PrintTo(fi.GetLabel(), -0.48, -0.03, -0.01, 1.0/f.windowWindth, wrapper, fi.GetSurface(), []mgl32.Vec3{mgl32.Vec3{0, 0, 1}})
	f.formItems = append(f.formItems, fi)
}

// CharCallback is the character stream input handler
func (f *FormScreen) CharCallback(char rune, wrapper interfaces.GLWrapper) {
	if f.underEdit != nil {
		switch f.underEdit.(type) {
		case *model.FormItemInt:
			fi := f.underEdit.(*model.FormItemInt)
			// offset for the current character has to be calculated.
			offsetX := f.charset.TextWidth(string(char), 1.0/f.windowWindth)
			fi.CharCallback(char, offsetX)
			f.charset.CleanSurface(fi.GetTarget())
			f.charset.PrintTo(fi.ValueToString(), -model.CursorInitX, -0.015, -0.01, 1.0/f.windowWindth, wrapper, fi.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
			break
		}
	}
}
