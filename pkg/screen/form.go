package screen

import (
	"math"
	"strconv"
	"strings"

	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/config"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/transformations"

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
	FullWidth          = BottomFrameWidth - 2*SideFrameWidth // The full width of the usable area
	EventEpsilon       = 200
	BACK_SPACE         = glfw.KeyBackspace
	KEY_UP             = glfw.KeyUp
	KEY_DOWN           = glfw.KeyDown
	LabelFontScale     = float32(1.0)
	InputTextFontScale = float32(0.80)
	ZFrame             = float32(0.0)
	ZText              = float32(-0.01)
	ZBackground        = float32(0.02)
	FormItemMoveSpeed  = float32(0.005)
	formItemMinY       = float32(0.0)
)

var (
	DirectionalLightDirection = mgl32.Vec3{0.0, 0.0, 1.0}.Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightDiffuse   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightSpecular  = mgl32.Vec3{1.0, 1.0, 1.0}

	DefaultFormItemMaterial   = material.Whiteplastic
	HighlightFormItemMaterial = material.Ruby
)

type FormScreen struct {
	*ScreenBase
	charset         *model.Charset
	formItems       []interfaces.FormItem
	formItemShader  *shader.Shader
	sinceLastClick  float64
	sinceLastDelete float64
	underEdit       interfaces.CharFormItem
	// Item position
	currentY         float32
	formItemCurrentY float32
	lastItemState    string
	// Info box for displaying the details of the form items.
	detailContentBox interfaces.Mesh
	// configuration package
	configuration config.Config
	// map for formItem-configItem
	formItemToConf map[interfaces.FormItem]*config.ConfigItem
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
	return camera
}

func setup(wrapper interfaces.GLWrapper) {
	wrapper.ClearColor(0.55, 0.55, 0.55, 1.0)
	wrapper.Enable(glwrapper.DEPTH_TEST)
	wrapper.DepthFunc(glwrapper.LESS)
	wrapper.Enable(glwrapper.BLEND)
	wrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
}

// NewFormScreenFromConfig returns a FormScreen. The formItems are also set, based on the input config, and config order.
func NewFormScreenFromConfig(frame *material.Material, label string, wrapper interfaces.GLWrapper, wW, wH float32, conf config.Config, formItemOrder []string) *FormScreen {
	formScreen := NewFormScreen(frame, label, wrapper, wW, wH)
	formScreen.configuration = conf
	for i := 0; i < len(formItemOrder); i++ {
		key := formItemOrder[i]
		if _, ok := formScreen.configuration[key]; ok {
			switch formScreen.configuration[key].GetValueType() {
			case config.ValueTypeInt:
				formScreen.addFormItemFromConfigInt(conf[key], wrapper)
				break
			case config.ValueTypeInt64:
				formScreen.addFormItemFromConfigInt64(conf[key], wrapper)
				break
			case config.ValueTypeFloat:
				formScreen.addFormItemFromConfigFloat(conf[key], wrapper)
				break
			case config.ValueTypeText:
				formScreen.addFormItemFromConfigText(conf[key], wrapper)
				break
			case config.ValueTypeBool:
				formScreen.addFormItemFromConfigBool(conf[key], wrapper)
				break
			case config.ValueTypeVector:
				formScreen.addFormItemFromConfigVector(conf[key], wrapper)
				break
			}
		}
	}
	return formScreen
}

// NewFormScreen returns a FormScreen. The screen contains a material Frame.
func NewFormScreen(frame *material.Material, label string, wrapper interfaces.GLWrapper, wW, wH float32) *FormScreen {
	s := newScreenBase()
	s.SetCamera(createCamera(wW / wH))
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
	frameModel := model.New()
	// create frame.
	bottomFrame := frameRectangle(BottomFrameWidth, BottomFrameLength, mgl32.Vec3{0.0, -0.99, ZFrame}, frame, wrapper)
	leftFrame := frameRectangle(SideFrameWidth, SideFrameLength, mgl32.Vec3{-0.99, 0.0, ZFrame}, frame, wrapper)
	rightFrame := frameRectangle(SideFrameWidth, SideFrameLength, mgl32.Vec3{0.99, 0.0, ZFrame}, frame, wrapper)
	topLeftFrame := frameRectangle(TopLeftFrameWidth, BottomFrameLength, mgl32.Vec3{0.95, 0.99, ZFrame}, frame, wrapper)
	textWidth := chars.TextWidth(label, 3.0/wW)
	textContainer := frameRectangle(textWidth, 0.15, mgl32.Vec3{1 - TopLeftFrameWidth - textWidth/2, 0.925, ZFrame}, frame, wrapper)
	textContainer.RotateX(-180)
	textContainer.RotateY(180)
	chars.PrintTo(label, -textWidth/2, -0.05, ZText, 3.0/wW, wrapper, textContainer, []mgl32.Vec3{mgl32.Vec3{0, 0, 1}})
	topRightFrame := frameRectangle(2.0-TopLeftFrameWidth-textWidth, BottomFrameLength, mgl32.Vec3{(-TopLeftFrameWidth - textWidth) / 2, 0.99, ZFrame}, frame, wrapper)
	detailContainer := frameRectangle(FullWidth, 0.3, mgl32.Vec3{0.0, -1.0 + BottomFrameLength + 0.15, ZFrame}, DefaultFormItemMaterial, wrapper)
	detailContainer.RotateX(-180)
	detailContainer.RotateY(180)
	frameModel.AddMesh(bottomFrame)
	frameModel.AddMesh(leftFrame)
	frameModel.AddMesh(rightFrame)
	frameModel.AddMesh(topLeftFrame)
	frameModel.AddMesh(topRightFrame)
	frameModel.AddMesh(detailContainer)
	s.AddModelToShader(frameModel, frameShaderApplication)
	return &FormScreen{
		ScreenBase:       s,
		charset:          chars,
		formItemShader:   bgShaderApplication,
		sinceLastClick:   0,
		currentY:         0.9,
		formItemCurrentY: float32(0.0),
		lastItemState:    "F",
		detailContentBox: detailContainer,
		configuration:    config.New(),
		formItemToConf:   make(map[interfaces.FormItem]*config.ConfigItem),
	}
}

// initMaterialForTheFormItems sets the material to the default of the form items.
// It could be used un the update loop to make all of them to default state.
func (f *FormScreen) initMaterialForTheFormItems() {
	f.charset.CleanSurface(f.detailContentBox)
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
					lightMesh.Material = HighlightFormItemMaterial
				} else {
					lightMesh.Material = DefaultFormItemMaterial
				}
				break
			case *model.FormItemInt64:
				fi := f.shaderMap[s][index].(*model.FormItemInt64)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = DefaultFormItemMaterial
				break
			case *model.FormItemVector:
				fi := f.shaderMap[s][index].(*model.FormItemVector)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = DefaultFormItemMaterial
				break
			}
		}
	}
}

// deleteCursor removes the cursor from the text form inputs.
func (f *FormScreen) deleteCursor() {
	if f.underEdit != nil {
		f.underEdit.DeleteCursor()
		f.underEdit = nil
	}
}
func (f *FormScreen) wrapTextToLines(desc string, scale, maxLineWidth float32) []string {
	words := strings.Split(desc, " ")
	var result []string
	line := ""
	for i := 0; i < len(words); i++ {
		lineWithNextWorld := line + words[i]
		width := f.charset.TextWidth(lineWithNextWorld, scale)
		if width > maxLineWidth {
			result = append(result, line)
			line = words[i]
			continue
		}
		if i < len(words)-1 {
			lineWithNextWorld = line + words[i] + " "
			width = f.charset.TextWidth(lineWithNextWorld, scale)
			if width > maxLineWidth {
				result = append(result, line+words[i])
				line = ""
				continue
			}
		}
		line = lineWithNextWorld
	}
	result = append(result, line)
	return result
}

// highlightFormAction updates the material of the closest mesh.
// It also prints the details of the form item to the detail content box.
func (f *FormScreen) highlightFormAction() {
	tmMesh := f.closestMesh.(*mesh.TexturedMaterialMesh)
	tmMesh.Material = HighlightFormItemMaterial
	desc := f.closestModel.(interfaces.FormItem).GetDescription()
	lines := f.wrapTextToLines(desc, InputTextFontScale/f.windowWindth, FullWidth)
	for i := 0; i < len(lines); i++ {
		f.charset.PrintTo(lines[i], -FullWidth/2, 0.12-float32(i)*0.075, ZText, InputTextFontScale/f.windowWindth, f.wrapper, f.detailContentBox, []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
	}
}

// Update loops on the shaderMap, and calls Update function on every Model.
// It also handles the camera movement and rotation, if the camera is set.
func (f *FormScreen) Update(dt, posX, posY float64, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
	f.sinceLastClick = f.sinceLastClick + dt
	f.sinceLastDelete = f.sinceLastDelete + dt
	cursorX := float32(-posX)
	cursorY := float32(posY)
	direction := mgl32.Vec3{0, 0, 0}
	// If the Up key is pressed => direction: up, velocity: c
	// If the Down key is pressed => direction: down, velocity: c
	// Otherwise => direction: null, velocity c
	if keyStore.Get(KEY_UP) && !keyStore.Get(KEY_DOWN) {
		direction = mgl32.Vec3{0, -1, 0}
	} else if keyStore.Get(KEY_DOWN) && !keyStore.Get(KEY_UP) {
		direction = mgl32.Vec3{0, 1, 0}
	}
	newYPos := f.formItemCurrentY + FormItemMoveSpeed*float32(dt)*direction.Y()
	maxYValue := (-1 + BottomFrameLength + 0.5 - f.currentY)
	if newYPos > formItemMinY && newYPos < maxYValue {
		f.formItemCurrentY = newYPos
	} else {
		direction = mgl32.Vec3{0, 0, 0}
	}
	for m, _ := range f.shaderMap[f.formItemShader] {
		f.shaderMap[f.formItemShader][m].SetDirection(direction)
	}
	coords := mgl32.Vec3{cursorX, cursorY, ZBackground}
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
	minDiff := float32(0.0)
	if closestDistance <= minDiff+0.01 {
		if buttonStore.Get(LEFT_MOUSE_BUTTON) && f.sinceLastClick > EventEpsilon {
			f.deleteCursor()
			f.sinceLastClick = 0
			switch f.closestModel.(type) {
			case *model.FormItemBool:
				formModel := f.closestModel.(*model.FormItemBool)
				formModel.SetValue(!formModel.GetValue())
				f.formItemToConf[formModel].SetCurrentValue(formModel.GetValue())
				break
			case *model.FormItemInt:
				formModel := f.closestModel.(*model.FormItemInt)
				formModel.AddCursor()
				f.underEdit = formModel
				break
			case *model.FormItemFloat:
				formModel := f.closestModel.(*model.FormItemFloat)
				formModel.AddCursor()
				f.underEdit = formModel
				break
			case *model.FormItemText:
				formModel := f.closestModel.(*model.FormItemText)
				formModel.AddCursor()
				f.underEdit = formModel
				break
			case *model.FormItemInt64:
				formModel := f.closestModel.(*model.FormItemInt64)
				formModel.AddCursor()
				f.underEdit = formModel
				break
			case *model.FormItemVector:
				formModel := f.closestModel.(*model.FormItemVector)
				msh, _ := formModel.ClosestMeshTo(mgl32.Vec3{coords.X(), coords.Y(), coords.Z() - 0.01})
				index := formModel.GetIndex(msh)
				if index > -1 {
					formModel.SetTarget(index - 1)
				}
				formModel.AddCursor()
				f.underEdit = formModel
				break
			}
		}
		f.highlightFormAction()
	}
	if keyStore.Get(BACK_SPACE) && f.sinceLastDelete > EventEpsilon {
		f.sinceLastDelete = 0
		if f.underEdit != nil {
			f.underEdit.DeleteLastCharacter()
			f.syncFormItemValuesToConfigValue()
			f.charset.CleanSurface(f.underEdit.GetTarget())
			switch f.underEdit.(type) {
			case *model.FormItemVector:
				f.charset.PrintTo(f.underEdit.ValueToString(), -f.underEdit.(*model.FormItemVector).GetVectorCursorInitialPosition().X(), -0.015, ZText, InputTextFontScale/f.windowWindth, f.wrapper, f.underEdit.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
				break
			default:
				f.charset.PrintTo(f.underEdit.ValueToString(), -f.underEdit.GetCursorInitialPosition().X(), -0.015, ZText, InputTextFontScale/f.windowWindth, f.wrapper, f.underEdit.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
				break
			}
		}
	}
}
func (f *FormScreen) syncFormItemValuesToConfigValue() {
	switch f.underEdit.(type) {
	case *model.FormItemInt:
		f.formItemToConf[f.underEdit].SetCurrentValue(f.underEdit.(*model.FormItemInt).GetValue())
		break
	case *model.FormItemFloat:
		f.formItemToConf[f.underEdit].SetCurrentValue(f.underEdit.(*model.FormItemFloat).GetValue())
		break
	case *model.FormItemText:
		f.formItemToConf[f.underEdit].SetCurrentValue(f.underEdit.(*model.FormItemText).GetValue())
		break
	case *model.FormItemInt64:
		f.formItemToConf[f.underEdit].SetCurrentValue(f.underEdit.(*model.FormItemInt64).GetValue())
		break
	case *model.FormItemVector:
		f.formItemToConf[f.underEdit].SetCurrentValue(f.underEdit.(*model.FormItemVector).GetValue())
		break
	}
}

func (f *FormScreen) addFormItem(fi interfaces.FormItem, wrapper interfaces.GLWrapper, defaultValue interface{}) int {
	fi.RotateX(-90)
	fi.RotateY(180)
	fi.SetSpeed(FormItemMoveSpeed)
	f.AddModelToShader(fi, f.formItemShader)
	f.charset.PrintTo(fi.GetLabel(), -(fi.GetFormItemWidth()/2)*0.999, -0.03, ZText, LabelFontScale/f.windowWindth, wrapper, fi.GetSurface(), []mgl32.Vec3{mgl32.Vec3{0, 0, 1}})
	f.formItems = append(f.formItems, fi)
	f.SetFormItemValue(len(f.formItems)-1, defaultValue, wrapper)
	return len(f.formItems) - 1
}

// addFormItemFromConfigBool sets up a FormItemBool from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigBool(configItem *config.ConfigItem, wrapper interfaces.GLWrapper) int {
	pos := f.itemPosition(model.ITEM_WIDTH_SHORT, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemBool(FullWidth, model.ITEM_WIDTH_SHORT, configItem.GetLabel(), configItem.GetDescription(), material.Whiteplastic, pos, wrapper)
	f.formItemToConf[fi] = configItem
	return f.addFormItem(fi, wrapper, configItem.GetDefaultValue())
}

// AddFormItemBool is for adding a bool form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemBool(formLabel, formDescription string, wrapper interfaces.GLWrapper, defaultValue bool) int {
	key := "bool_form_item_" + strconv.Itoa(len(f.configuration))
	f.configuration.AddConfig(key, formLabel, formDescription, defaultValue, nil)
	return f.addFormItemFromConfigBool(f.configuration[key], wrapper)
}

// addFormItemFromConfigInt sets up a FormItemInt from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigInt(configItem *config.ConfigItem, wrapper interfaces.GLWrapper) int {
	pos := f.itemPosition(model.ITEM_WIDTH_HALF, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemInt(FullWidth, model.ITEM_WIDTH_HALF, configItem.GetLabel(), configItem.GetDescription(), material.Whiteplastic, pos, wrapper)
	if configItem.GetValidatorFunction() != nil {
		fi.SetValidator(configItem.GetValidatorFunction().(model.IntValidator))
	}
	f.formItemToConf[fi] = configItem
	return f.addFormItem(fi, wrapper, transformations.IntegerToString(configItem.GetDefaultValue().(int)))
}

// AddFormItemInt is for adding an integer form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemInt(formLabel, formDescription string, wrapper interfaces.GLWrapper, defaultValue int, validator model.IntValidator) int {
	key := "int_form_item_" + strconv.Itoa(len(f.configuration))
	f.configuration.AddConfig(key, formLabel, formDescription, defaultValue, validator)
	return f.addFormItemFromConfigInt(f.configuration[key], wrapper)
}

// addFormItemFromConfigFloat sets up a FormItemFloat from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigFloat(configItem *config.ConfigItem, wrapper interfaces.GLWrapper) int {
	pos := f.itemPosition(model.ITEM_WIDTH_HALF, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemFloat(FullWidth, model.ITEM_WIDTH_HALF, configItem.GetLabel(), configItem.GetDescription(), material.Whiteplastic, pos, wrapper)
	if configItem.GetValidatorFunction() != nil {
		fi.SetValidator(configItem.GetValidatorFunction().(model.FloatValidator))
	}
	f.formItemToConf[fi] = configItem
	return f.addFormItem(fi, wrapper, transformations.Float32ToStringExact(configItem.GetDefaultValue().(float32)))
}

// AddFormItemFloat is for adding a float form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemFloat(formLabel, formDescription string, wrapper interfaces.GLWrapper, defaultValue float32, validator model.FloatValidator) int {
	key := "float_form_item_" + strconv.Itoa(len(f.configuration))
	f.configuration.AddConfig(key, formLabel, formDescription, defaultValue, validator)
	return f.addFormItemFromConfigFloat(f.configuration[key], wrapper)
}

// addFormItemFromConfigText sets up a FormItemText from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigText(configItem *config.ConfigItem, wrapper interfaces.GLWrapper) int {
	pos := f.itemPosition(model.ITEM_WIDTH_FULL, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemText(FullWidth, model.ITEM_WIDTH_HALF, configItem.GetLabel(), configItem.GetDescription(), material.Whiteplastic, pos, wrapper)
	if configItem.GetValidatorFunction() != nil {
		fi.SetValidator(configItem.GetValidatorFunction().(model.StringValidator))
	}
	f.formItemToConf[fi] = configItem
	return f.addFormItem(fi, wrapper, configItem.GetDefaultValue())
}

// AddFormItemText is for adding a text form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemText(formLabel, formDescription string, wrapper interfaces.GLWrapper, defaultValue string, validator model.StringValidator) int {
	key := "text_form_item_" + strconv.Itoa(len(f.configuration))
	f.configuration.AddConfig(key, formLabel, formDescription, defaultValue, validator)
	return f.addFormItemFromConfigText(f.configuration[key], wrapper)
}

// addFormItemFromConfigInt64 sets up a FormItemInt64 from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigInt64(configItem *config.ConfigItem, wrapper interfaces.GLWrapper) int {
	pos := f.itemPosition(model.ITEM_WIDTH_HALF, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemInt64(FullWidth, model.ITEM_WIDTH_HALF, configItem.GetLabel(), configItem.GetDescription(), material.Whiteplastic, pos, wrapper)
	if configItem.GetValidatorFunction() != nil {
		fi.SetValidator(configItem.GetValidatorFunction().(model.Int64Validator))
	}
	f.formItemToConf[fi] = configItem
	return f.addFormItem(fi, wrapper, transformations.Integer64ToString(configItem.GetDefaultValue().(int64)))
}

// AddFormItemInt64 is for adding an int64 form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemInt64(formLabel, formDescription string, wrapper interfaces.GLWrapper, defaultValue int64, validator model.Int64Validator) int {
	key := "int64_form_item_" + strconv.Itoa(len(f.configuration))
	f.configuration.AddConfig(key, formLabel, formDescription, defaultValue, validator)
	return f.addFormItemFromConfigInt64(f.configuration[key], wrapper)
}

// addFormItemFromConfigVector sets up a FormItemInt64 from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigVector(configItem *config.ConfigItem, wrapper interfaces.GLWrapper) int {
	pos := f.itemPosition(model.ITEM_WIDTH_FULL, FullWidth*model.ITEM_HEIGHT_MULTIPLIER)
	fi := model.NewFormItemVector(FullWidth, model.ITEM_WIDTH_HALF, configItem.GetLabel(), configItem.GetDescription(), model.CHAR_NUM_FLOAT, material.Whiteplastic, pos, wrapper)
	if configItem.GetValidatorFunction() != nil {
		fi.SetValidator(configItem.GetValidatorFunction().(model.FloatValidator))
	}
	f.formItemToConf[fi] = configItem
	return f.addFormItem(fi, wrapper, transformations.VectorToString(configItem.GetDefaultValue().(mgl32.Vec3)))
}

// AddFormItemVector is for adding a vector form item to the form. It returns the index of the
// inserted item.
func (f *FormScreen) AddFormItemVector(formLabel, formDescription string, wrapper interfaces.GLWrapper, defaultValue mgl32.Vec3, validator model.FloatValidator) int {
	key := "vector_form_item_" + strconv.Itoa(len(f.configuration))
	f.configuration.AddConfig(key, formLabel, formDescription, defaultValue, validator)
	return f.addFormItemFromConfigVector(f.configuration[key], wrapper)
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
		offsetX := f.charset.TextWidth(string(char), InputTextFontScale/f.windowWindth)
		f.underEdit.CharCallback(char, offsetX)
		f.syncFormItemValuesToConfigValue()
		f.charset.CleanSurface(f.underEdit.GetTarget())
		switch f.underEdit.(type) {
		case *model.FormItemVector:
			f.charset.PrintTo(f.underEdit.ValueToString(), -f.underEdit.(*model.FormItemVector).GetVectorCursorInitialPosition().X(), -0.015, ZText, InputTextFontScale/f.windowWindth, wrapper, f.underEdit.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
			break
		default:
			f.charset.PrintTo(f.underEdit.ValueToString(), -f.underEdit.GetCursorInitialPosition().X(), -0.015, ZText, InputTextFontScale/f.windowWindth, wrapper, f.underEdit.GetTarget(), []mgl32.Vec3{mgl32.Vec3{0, 0.5, 0}})
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
	case *model.FormItemVector:
		f.underEdit = item.(*model.FormItemVector)
		newValues := valueNew.([3]string)
		for i := 0; i < 3; i++ {
			item.(*model.FormItemVector).SetTarget(i)
			value := f.underEdit.ValueToString()
			valueLength := len(value) + strings.Count(value, " ")
			for i := 0; i < valueLength; i++ {
				f.underEdit.DeleteLastCharacter()
			}
			f.setDefaultValueChar(newValues[i], wrapper)
		}
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

	return mgl32.Vec3{x, f.currentY, ZBackground}
}
