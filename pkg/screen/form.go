package screen

import (
	"math"
	"strings"

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
	BottomFrameWidth   = float32(2.0) // the full width of the screen.
	BottomFrameLength  = float32(0.02)
	SideFrameLength    = float32(1.98)
	SideFrameWidth     = float32(0.02) // the width of the border frames.
	TopLeftFrameWidth  = float32(0.1)
	FullWidth          = BottomFrameWidth - 2*SideFrameWidth
	CameraMoveSpeed    = 0.005
	LightConstantTerm  = float32(1.0)
	LightLinearTerm    = float32(0.14)
	LightQuadraticTerm = float32(0.07)
	EventEpsilon       = 200
	BACK_SPACE         = glfw.KeyBackspace

	FullX       = 0
	HalfLeftX   = 0.4
	HalfRigthX  = -0.4
	LongLeftX   = 0.4
	LongRightX  = -0.4
	ShortLeftX  = 0.4
	ShortRightX = 0.4
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
	// Item position
	currentY      float32
	lastItemState string
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
		currentY:       0.9,
		lastItemState:  "F",
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
			case *model.FormItemFloat:
				fi := f.shaderMap[s][index].(*model.FormItemFloat)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = DefaultFormItemMaterial
				break
			case *model.FormItemText:
				fi := f.shaderMap[s][index].(*model.FormItemText)
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
			case *model.FormItemInt64:
				fi := f.shaderMap[s][index].(*model.FormItemInt64)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = DefaultFormItemMaterial
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
	case *model.FormItemFloat:
		fi := f.underEdit.(*model.FormItemFloat)
		fi.DeleteCursor()
		f.underEdit = nil
		break
	case *model.FormItemText:
		fi := f.underEdit.(*model.FormItemText)
		fi.DeleteCursor()
		f.underEdit = nil
		break
	case *model.FormItemInt64:
		fi := f.underEdit.(*model.FormItemInt64)
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
	case *model.FormItemText:
		tmMesh := f.closestMesh.(*mesh.TexturedMaterialMesh)
		minDiff := float32(0.0)
		if closestDistance <= minDiff+0.01 {
			tmMesh.Material = material.Ruby
			if buttonStore.Get(LEFT_MOUSE_BUTTON) {
				formModel := f.closestModel.(*model.FormItemText)
				if f.sinceLastClick > EventEpsilon {
					f.deleteCursor()
					formModel.AddCursor()
					f.sinceLastClick = 0
					f.underEdit = formModel
				}
			}
		}
		break
	case *model.FormItemInt64:
		tmMesh := f.closestMesh.(*mesh.TexturedMaterialMesh)
		minDiff := float32(0.0)
		if closestDistance <= minDiff+0.01 {
			tmMesh.Material = material.Ruby
			if buttonStore.Get(LEFT_MOUSE_BUTTON) {
				formModel := f.closestModel.(*model.FormItemInt64)
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
			f.sinceLastDelete = 0
			if f.underEdit != nil {
				f.underEdit.DeleteLastCharacter()
				f.charset.CleanSurface(f.underEdit.GetTarget())
				f.charset.PrintTo(f.underEdit.ValueToString(), -f.underEdit.GetCursorInitialPosition().X(), -0.015, -0.01, 1.0/f.windowWindth, f.wrapper, f.underEdit.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
			}
		}
	}
}

func (f *FormScreen) addFormItem(fi interfaces.FormItem, wrapper interfaces.GLWrapper, defaultValue interface{}) int {
	fi.RotateX(-90)
	fi.RotateY(180)
	f.AddModelToShader(fi, f.bgShader)
	f.charset.PrintTo(fi.GetLabel(), -(fi.GetFormItemWidth()/2)*0.999, -0.03, -0.01, 1.0/f.windowWindth, wrapper, fi.GetSurface(), []mgl32.Vec3{mgl32.Vec3{0, 0, 1}})
	f.formItems = append(f.formItems, fi)
	f.SetFormItemValue(len(f.formItems)-1, defaultValue, wrapper)
	return len(f.formItems) - 1
}

// AddFormItemBool is for adding a bool form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemBool(formLabel string, wrapper interfaces.GLWrapper, defaultValue bool) int {
	pos := f.itemPosition(model.ITEM_WIDTH_SHORT, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemBool(FullWidth, model.ITEM_WIDTH_SHORT, formLabel, material.Whiteplastic, pos, wrapper)
	return f.addFormItem(fi, wrapper, defaultValue)
}

// AddFormItemInt is for adding an integer form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemInt(formLabel string, wrapper interfaces.GLWrapper, defaultValue string) int {
	pos := f.itemPosition(model.ITEM_WIDTH_HALF, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemInt(FullWidth, model.ITEM_WIDTH_HALF, formLabel, material.Whiteplastic, pos, wrapper)
	return f.addFormItem(fi, wrapper, defaultValue)
}

// AddFormItemFloat is for adding a float form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemFloat(formLabel string, wrapper interfaces.GLWrapper, defaultValue string) int {
	pos := f.itemPosition(model.ITEM_WIDTH_HALF, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemFloat(FullWidth, model.ITEM_WIDTH_HALF, formLabel, material.Whiteplastic, pos, wrapper)
	return f.addFormItem(fi, wrapper, defaultValue)
}

// AddFormItemText is for adding a text form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemText(formLabel string, wrapper interfaces.GLWrapper, defaultValue string) int {
	pos := f.itemPosition(model.ITEM_WIDTH_FULL, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemText(FullWidth, model.ITEM_WIDTH_FULL, formLabel, material.Whiteplastic, pos, wrapper)
	return f.addFormItem(fi, wrapper, defaultValue)
}

// AddFormItemInt64 is for adding an int64 form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemInt64(formLabel string, wrapper interfaces.GLWrapper, defaultValue string) int {
	pos := f.itemPosition(model.ITEM_WIDTH_LONG, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemInt64(FullWidth, model.ITEM_WIDTH_LONG, formLabel, material.Whiteplastic, pos, wrapper)
	return f.addFormItem(fi, wrapper, defaultValue)
}
func (f *FormScreen) setDefaultValueChar(input string, wrapper interfaces.GLWrapper) {
	chars := []rune(input)
	for i := 0; i < len(chars); i++ {
		f.CharCallback(chars[i], wrapper)
	}
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
			f.charset.PrintTo(fi.ValueToString(), -f.underEdit.GetCursorInitialPosition().X(), -0.015, -0.01, 1.0/f.windowWindth, wrapper, fi.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
			break
		case *model.FormItemFloat:
			fi := f.underEdit.(*model.FormItemFloat)
			// offset for the current character has to be calculated.
			offsetX := f.charset.TextWidth(string(char), 1.0/f.windowWindth)
			fi.CharCallback(char, offsetX)
			f.charset.CleanSurface(fi.GetTarget())
			f.charset.PrintTo(fi.ValueToString(), -f.underEdit.GetCursorInitialPosition().X(), -0.015, -0.01, 1.0/f.windowWindth, wrapper, fi.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
			break
		case *model.FormItemText:
			fi := f.underEdit.(*model.FormItemText)
			// offset for the current character has to be calculated.
			offsetX := f.charset.TextWidth(string(char), 1.0/f.windowWindth)
			fi.CharCallback(char, offsetX)
			f.charset.CleanSurface(fi.GetTarget())
			f.charset.PrintTo(fi.ValueToString(), -f.underEdit.GetCursorInitialPosition().X(), -0.015, -0.01, 1.0/f.windowWindth, wrapper, fi.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
			break
		case *model.FormItemInt64:
			fi := f.underEdit.(*model.FormItemInt64)
			// offset for the current character has to be calculated.
			offsetX := f.charset.TextWidth(string(char), 1.0/f.windowWindth)
			fi.CharCallback(char, offsetX)
			f.charset.CleanSurface(fi.GetTarget())
			f.charset.PrintTo(fi.ValueToString(), -f.underEdit.GetCursorInitialPosition().X(), -0.015, -0.01, 1.0/f.windowWindth, wrapper, fi.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
			break
		}
	}
}

// GetFormItem gets an index as input, and returns formItem[index] form item.
// In case of invalid index, it panics.
func (f *FormScreen) GetFormItem(index int) interfaces.FormItem {
	if index > len(f.formItems) {
		panic("Invalid form item index.")
	}
	return f.formItems[index]
}

// SetFormItemValue gets an index, a value and and a wrapper and sets the form
// item value of the formItem[index].
func (f *FormScreen) SetFormItemValue(index int, valueNew interface{}, wrapper interfaces.GLWrapper) {
	item := f.GetFormItem(index)
	switch item.(type) {
	case *model.FormItemInt:
		f.underEdit = item.(*model.FormItemInt)
		value := f.underEdit.ValueToString()
		valueLength := len(value) + strings.Count(value, " ")
		for i := 0; i < valueLength; i++ {
			f.underEdit.DeleteLastCharacter()
		}
		f.setDefaultValueChar(valueNew.(string), wrapper)
		break
	case *model.FormItemFloat:
		f.underEdit = item.(*model.FormItemFloat)
		value := f.underEdit.ValueToString()
		valueLength := len(value) + strings.Count(value, " ")
		for i := 0; i < valueLength; i++ {
			f.underEdit.DeleteLastCharacter()
		}
		f.setDefaultValueChar(valueNew.(string), wrapper)
		break
	case *model.FormItemText:
		f.underEdit = item.(*model.FormItemText)
		value := f.underEdit.ValueToString()
		valueLength := len(value) + strings.Count(value, " ")
		for i := 0; i < valueLength; i++ {
			f.underEdit.DeleteLastCharacter()
		}
		f.setDefaultValueChar(valueNew.(string), wrapper)
		break
	case *model.FormItemInt64:
		f.underEdit = item.(*model.FormItemInt64)
		value := f.underEdit.ValueToString()
		valueLength := len(value) + strings.Count(value, " ")
		for i := 0; i < valueLength; i++ {
			f.underEdit.DeleteLastCharacter()
		}
		f.setDefaultValueChar(valueNew.(string), wrapper)
		break
	case *model.FormItemBool:
		fi := item.(*model.FormItemBool)
		fi.SetValue(valueNew.(bool))
		break
	}
}
func (f *FormScreen) pushState(itemWidth float32) {
	switch f.lastItemState {
	case "F", "RH", "RL", "RS":
		switch itemWidth {
		case model.ITEM_WIDTH_FULL:
			f.lastItemState = "F"
			break
		case model.ITEM_WIDTH_HALF:
			f.lastItemState = "LH"
			break
		case model.ITEM_WIDTH_LONG:
			f.lastItemState = "LL"
			break
		case model.ITEM_WIDTH_SHORT:
			f.lastItemState = "LS"
			break
		default:
			panic("Unhandled state.")
		}
		break
	case "LH":
		switch itemWidth {
		case model.ITEM_WIDTH_FULL:
			f.lastItemState = "F"
			break
		case model.ITEM_WIDTH_HALF:
			f.lastItemState = "RH"
			break
		case model.ITEM_WIDTH_LONG:
			f.lastItemState = "LL"
			break
		case model.ITEM_WIDTH_SHORT:
			f.lastItemState = "RS"
			break
		default:
			panic("Unhandled state.")
		}
		break
	case "LL":
		switch itemWidth {
		case model.ITEM_WIDTH_FULL:
			f.lastItemState = "F"
			break
		case model.ITEM_WIDTH_HALF:
			f.lastItemState = "LH"
			break
		case model.ITEM_WIDTH_LONG:
			f.lastItemState = "LL"
			break
		case model.ITEM_WIDTH_SHORT:
			f.lastItemState = "RS"
			break
		default:
			panic("Unhandled state.")
		}
		break
	case "LS":
		switch itemWidth {
		case model.ITEM_WIDTH_FULL:
			f.lastItemState = "F"
			break
		case model.ITEM_WIDTH_HALF:
			f.lastItemState = "RH"
			break
		case model.ITEM_WIDTH_LONG:
			f.lastItemState = "RL"
			break
		case model.ITEM_WIDTH_SHORT:
			f.lastItemState = "MS"
			break
		default:
			panic("Unhandled state.")
		}
		break
	case "MS":
		switch itemWidth {
		case model.ITEM_WIDTH_FULL:
			f.lastItemState = "F"
			break
		case model.ITEM_WIDTH_HALF:
			f.lastItemState = "LH"
			break
		case model.ITEM_WIDTH_LONG:
			f.lastItemState = "LL"
			break
		case model.ITEM_WIDTH_SHORT:
			f.lastItemState = "RS"
			break
		default:
			panic("Unhandled state.")
		}
		break
	default:
		panic("Unhandled state.")
	}
}
func (f *FormScreen) itemPosition(itemWidth, itemHeight float32) mgl32.Vec3 {
	f.pushState(itemWidth)
	var x float32
	switch f.lastItemState {
	case "F":
		f.currentY = f.currentY - itemHeight
		x = 0.0
		break
	case "LH":
		f.currentY = f.currentY - itemHeight
		x = FullWidth / 4
		break
	case "LL":
		f.currentY = f.currentY - itemHeight
		x = FullWidth / 6
		break
	case "LS":
		f.currentY = f.currentY - itemHeight
		x = FullWidth / 3
		break
	case "RH":
		x = -FullWidth / 4
		break
	case "RL":
		x = -FullWidth / 6
		break
	case "RS":
		x = -FullWidth / 3
		break
	case "MS":
		x = 0.0
		break
	}

	return mgl32.Vec3{x, f.currentY, 0}
}
