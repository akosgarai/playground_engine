package screen

import (
	"math"
	"strconv"
	"strings"

	"github.com/akosgarai/playground_engine/pkg/config"
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
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
	DefaultFontFile    = "/assets/fonts/Frijole/Frijole-Regular.ttf"
	TopLeftFrameWidth  = float32(0.1)
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
	DefaultFrameWidth  = float32(2.0) // the full width of the screen.
	DefaultFrameLength = float32(0.02)
)

var (
	DirectionalLightDirection = mgl32.Vec3{0.0, 0.0, 1.0}.Normalize()
	DirectionalLightAmbient   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightDiffuse   = mgl32.Vec3{0.5, 0.5, 0.5}
	DirectionalLightSpecular  = mgl32.Vec3{1.0, 1.0, 1.0}

	DefaultFormItemMaterial   = material.Whiteplastic
	HighlightFormItemMaterial = material.Ruby

	DefaultFormScreenHeaderLabelColor   = mgl32.Vec3{0, 0, 1}
	DefaultFormScreenFormItemLabelColor = mgl32.Vec3{0, 0, 1}
	DefaultFormScreenFormItemInputColor = mgl32.Vec3{0, 0.5, 0}
	DefaultClearColor                   = mgl32.Vec3{0.55, 0.55, 0.55}
)

type FormScreenBuilder struct {
	*ScreenWithFrameBuilder
	headerLabel               string
	formItemMaterial          *material.Material
	formItemHighlightMaterial *material.Material
	config                    config.Config
	configOrder               []string
	charset                   *model.Charset
	lastItemState             string
	offsetY                   float32
	headerLabelColor          mgl32.Vec3
	formItemLabelColor        mgl32.Vec3
	formItemInputColor        mgl32.Vec3
	clearColor                mgl32.Vec3
}

func NewFormScreenBuilder() *FormScreenBuilder {
	swfb := NewScreenWithFrameBuilder()
	swfb.SetFrameSize(DefaultFrameWidth, DefaultFrameLength, TopLeftFrameWidth)
	swfb.SetFrameMaterial(HighlightFormItemMaterial)
	swfb.SetDetailContentBoxMaterial(DefaultFormItemMaterial)
	swfb.SetDetailContentBoxHeight(0.3)
	return &FormScreenBuilder{
		ScreenWithFrameBuilder:    swfb,
		headerLabel:               "Default label",
		charset:                   nil,
		config:                    config.New(),
		lastItemState:             "F",
		offsetY:                   0.9,
		formItemMaterial:          DefaultFormItemMaterial,
		formItemHighlightMaterial: HighlightFormItemMaterial,
		headerLabelColor:          DefaultFormScreenHeaderLabelColor,
		formItemLabelColor:        DefaultFormScreenFormItemLabelColor,
		formItemInputColor:        DefaultFormScreenFormItemInputColor,
		clearColor:                DefaultClearColor,
	}
}

// GetFullWidth returns the width of the drawable screen. (width - 2*length)
func (f *FormScreenBuilder) GetFullWidth() float32 {
	return f.frameWidth - (2 * f.frameLength)
}

// SetHeaderLabel sets the value for the header of the form, that is displayed
// on the top left of the form.
func (b *FormScreenBuilder) SetHeaderLabel(l string) {
	b.headerLabel = l
}

// SetFormItemMaterial sets the material that is used for the form items and the detailcontentbox.
func (b *FormScreenBuilder) SetFormItemMaterial(m *material.Material) {
	b.formItemMaterial = m
}

// SetFormItemHighlightMaterial sets the material that is used for the hovered form items.
func (b *FormScreenBuilder) SetFormItemHighlightMaterial(m *material.Material) {
	b.formItemHighlightMaterial = m
}

// SetConfig sets the config of the form.
func (b *FormScreenBuilder) SetConfig(c config.Config) {
	b.config = c
}

// SetConfigOrder setsh the order of the config items.
func (b *FormScreenBuilder) SetConfigOrder(o []string) {
	b.configOrder = o
}

// SetCharset sets the charset of the form screen.
func (b *FormScreenBuilder) SetCharset(m *model.Charset) {
	b.charset = m
}

// SetHeaderLabelColor sets the color of the header label.
func (b *FormScreenBuilder) SetHeaderLabelColor(c mgl32.Vec3) {
	b.headerLabelColor = c
}

// SetFormItemLabelColor sets the color of the FormItem labels.
func (b *FormScreenBuilder) SetFormItemLabelColor(c mgl32.Vec3) {
	b.formItemLabelColor = c
}

// SetFormItemInputColor sets the color of the FormItem labels.
func (b *FormScreenBuilder) SetFormItemInputColor(c mgl32.Vec3) {
	b.formItemInputColor = c
}

// SetClearColor sets the color of the background.
func (b *FormScreenBuilder) SetClearColor(c mgl32.Vec3) {
	b.clearColor = c
}

// AddConfigBool is for adding a bool config item to the configs. It returns the key of the config.
func (b *FormScreenBuilder) AddConfigBool(formLabel, formDescription string, defaultValue bool) string {
	key := "bool_config_item_" + strconv.Itoa(len(b.config))
	b.config.AddConfig(key, formLabel, formDescription, defaultValue, nil)
	return key
}

// AddConfigInt is for adding an integer config item to the configs. It returns the key of the config.
func (b *FormScreenBuilder) AddConfigInt(formLabel, formDescription string, defaultValue int, validator model.IntValidator) string {
	key := "int_config_item_" + strconv.Itoa(len(b.config))
	b.config.AddConfig(key, formLabel, formDescription, defaultValue, validator)
	return key
}

// AddConfigFloat is for adding a float config item to the configs. It returns the key of the config.
func (b *FormScreenBuilder) AddConfigFloat(formLabel, formDescription string, defaultValue float32, validator model.FloatValidator) string {
	key := "float_config_item_" + strconv.Itoa(len(b.config))
	b.config.AddConfig(key, formLabel, formDescription, defaultValue, validator)
	return key
}

// AddConfigInt64 is for adding an int64 config item to the configs. It returns the key of the config.
func (b *FormScreenBuilder) AddConfigInt64(formLabel, formDescription string, defaultValue int64, validator model.Int64Validator) string {
	key := "int64_config_item_" + strconv.Itoa(len(b.config))
	b.config.AddConfig(key, formLabel, formDescription, defaultValue, validator)
	return key
}

// AddConfigVector is for adding a vector config item to the configs. It returns the key of the config.
func (b *FormScreenBuilder) AddConfigVector(formLabel, formDescription string, defaultValue mgl32.Vec3, validator model.FloatValidator) string {
	key := "vector_config_item_" + strconv.Itoa(len(b.config))
	b.config.AddConfig(key, formLabel, formDescription, defaultValue, validator)
	return key
}

// AddConfigText is for adding a text config item to the configs. It returns the key of the config.
func (b *FormScreenBuilder) AddConfigText(formLabel, formDescription string, defaultValue string, validator model.StringValidator) string {
	key := "text_config_item_" + strconv.Itoa(len(b.config))
	b.config.AddConfig(key, formLabel, formDescription, defaultValue, validator)
	return key
}

// It builds a form screen with the given setup. In case of missing wrapper, it panics.
func (b *FormScreenBuilder) Build() *FormScreen {
	if b.wrapper == nil {
		panic("Wrapper is missing for the build process.")
	}
	if b.charset == nil {
		b.defaultCharset()
	}

	// variables for aspect ratio.
	aspRatio := float32(b.windowWidth) / float32(b.windowHeight)

	textWidth := float32(0.0)
	textHeight := float32(0.0)
	if b.headerLabel != "" {
		textWidth, textHeight = b.charset.TextContainerSize(b.headerLabel, 3.0/b.windowWidth*aspRatio)
		b.SetLabelWidth(textWidth)
	}
	s := b.ScreenWithFrameBuilder.Build()

	bgShaderApplication := shader.NewTextureMatShaderBlending(b.wrapper)
	fgShaderApplication := shader.NewFontShader(b.wrapper)
	s.AddShader(bgShaderApplication)
	s.AddShader(fgShaderApplication)

	s.AddModelToShader(b.charset, fgShaderApplication)

	if b.headerLabel != "" {
		textContainerPosition := mgl32.Vec3{b.frameWidth/2 - b.frameTopLeftWidth - textWidth/2, (b.frameWidth/2/aspRatio - textHeight/2), ZFrame}
		textContainer := b.frameRectangle(textWidth, textHeight, textContainerPosition)
		textContainer.RotateX(-180)
		textContainer.RotateY(180)
		b.charset.PrintTo(b.headerLabel, -textWidth/2, -textHeight/2, ZText, 3.0/b.windowWidth*aspRatio, b.wrapper, textContainer, []mgl32.Vec3{b.headerLabelColor})
	}

	formScreen := &FormScreen{
		ScreenWithFrame:           s,
		charset:                   b.charset,
		formItemShader:            bgShaderApplication,
		sinceLastClick:            0,
		currentScrollOffset:       float32(0.0),
		formItemToConf:            make(map[interfaces.FormItem]*config.ConfigItem),
		formItemLabelColor:        b.formItemLabelColor,
		formItemInputColor:        b.formItemInputColor,
		formItemDefaultMaterial:   b.formItemMaterial,
		formItemHighlightMaterial: b.formItemHighlightMaterial,
		clearColor:                b.clearColor,
	}
	s.Setup(formScreen.setupFormScreen)
	b.offsetY = (b.frameWidth/2 - 0.1) / aspRatio
	for i := 0; i < len(b.configOrder); i++ {
		key := b.configOrder[i]
		if _, ok := b.config[key]; ok {
			switch b.config[key].GetValueType() {
			case config.ValueTypeInt:
				formScreen.addFormItemFromConfigInt(b.config[key], b.itemPosition(model.ITEM_WIDTH_HALF, formScreen.GetFullWidth()*model.ITEM_HEIGHT_MULTIPLIER/aspRatio))
				break
			case config.ValueTypeInt64:
				formScreen.addFormItemFromConfigInt64(b.config[key], b.itemPosition(model.ITEM_WIDTH_LONG, formScreen.GetFullWidth()*model.ITEM_HEIGHT_MULTIPLIER/aspRatio))
				break
			case config.ValueTypeFloat:
				formScreen.addFormItemFromConfigFloat(b.config[key], b.itemPosition(model.ITEM_WIDTH_HALF, formScreen.GetFullWidth()*model.ITEM_HEIGHT_MULTIPLIER/aspRatio))
				break
			case config.ValueTypeText:
				formScreen.addFormItemFromConfigText(b.config[key], b.itemPosition(model.ITEM_WIDTH_FULL, formScreen.GetFullWidth()*model.ITEM_HEIGHT_MULTIPLIER/aspRatio))
				break
			case config.ValueTypeBool:
				formScreen.addFormItemFromConfigBool(b.config[key], b.itemPosition(model.ITEM_WIDTH_SHORT, formScreen.GetFullWidth()*model.ITEM_HEIGHT_MULTIPLIER/aspRatio))
				break
			case config.ValueTypeVector:
				formScreen.addFormItemFromConfigVector(b.config[key], b.itemPosition(model.ITEM_WIDTH_FULL, formScreen.GetFullWidth()*model.ITEM_HEIGHT_MULTIPLIER/aspRatio))
				break
			}
		}
	}
	// bottom pos (-width/2) + length of the frame mesh + length of the detail content box + length of one form item - offsetY. if this value is negative, we could us 0 instead.
	formScreen.maxScrollOffset = (-(b.frameWidth/2.0)+b.frameLength+b.detailContentBoxHeight+(formScreen.GetFullWidth()*model.ITEM_HEIGHT_MULTIPLIER))/aspRatio - b.offsetY
	if formScreen.maxScrollOffset < 0 {
		formScreen.maxScrollOffset = 0.0
	}

	return formScreen
}

// It creates a rectangle for the screen frame.
func (b *FormScreenBuilder) frameRectangle(width, length float32, position mgl32.Vec3) *mesh.TexturedMaterialMesh {
	return b.frameRectangleWithMaterial(width, length, position, b.frameMaterial)
}
func (b *FormScreenBuilder) frameRectangleWithMaterial(width, length float32, position mgl32.Vec3, mat *material.Material) *mesh.TexturedMaterialMesh {
	v, i, _ := rectangle.NewExact(width, length).MeshInput()
	var tex texture.Textures
	tex.TransparentTexture(1, 1, 128, "tex.diffuse", b.wrapper)
	tex.TransparentTexture(1, 1, 128, "tex.specular", b.wrapper)
	frameMesh := mesh.NewTexturedMaterialMesh(v, i, tex, mat, b.wrapper)
	frameMesh.RotateX(90)
	frameMesh.SetPosition(position)
	return frameMesh
}

func (b *FormScreenBuilder) defaultCharset() {
	cs, err := model.LoadCharset(baseDirScreen()+DefaultFontFile, 32, 127, 40.0, 72, b.wrapper)
	if err != nil {
		panic(err)
	}
	cs.SetTransparent(true)
	b.charset = cs
}

func (b *FormScreenBuilder) itemPosition(itemWidth, itemHeight float32) mgl32.Vec3 {
	b.pushState(itemWidth)
	var x float32
	switch b.lastItemState {
	case "F":
		b.offsetY = b.offsetY - itemHeight
		x = 0.0
		break
	case "LH":
		b.offsetY = b.offsetY - itemHeight
		x = b.GetFullWidth() / 4
		break
	case "LL":
		b.offsetY = b.offsetY - itemHeight
		x = b.GetFullWidth() / 6
		break
	case "LS":
		b.offsetY = b.offsetY - itemHeight
		x = b.GetFullWidth() / 3
		break
	case "RH":
		x = -b.GetFullWidth() / 4
		break
	case "RL":
		x = -b.GetFullWidth() / 6
		break
	case "RS":
		x = -b.GetFullWidth() / 3
		break
	case "MS":
		x = 0.0
		break
	}

	return mgl32.Vec3{x, b.offsetY, ZBackground}
}
func (b *FormScreenBuilder) pushState(itemWidth float32) {
	switch b.lastItemState {
	case "F", "RH", "RL", "RS":
		switch itemWidth {
		case model.ITEM_WIDTH_FULL:
			b.lastItemState = "F"
			break
		case model.ITEM_WIDTH_HALF:
			b.lastItemState = "LH"
			break
		case model.ITEM_WIDTH_LONG:
			b.lastItemState = "LL"
			break
		case model.ITEM_WIDTH_SHORT:
			b.lastItemState = "LS"
			break
		default:
			panic("Unhandled state.")
		}
		break
	case "LH":
		switch itemWidth {
		case model.ITEM_WIDTH_FULL:
			b.lastItemState = "F"
			break
		case model.ITEM_WIDTH_HALF:
			b.lastItemState = "RH"
			break
		case model.ITEM_WIDTH_LONG:
			b.lastItemState = "LL"
			break
		case model.ITEM_WIDTH_SHORT:
			b.lastItemState = "RS"
			break
		default:
			panic("Unhandled state.")
		}
		break
	case "LL":
		switch itemWidth {
		case model.ITEM_WIDTH_FULL:
			b.lastItemState = "F"
			break
		case model.ITEM_WIDTH_HALF:
			b.lastItemState = "LH"
			break
		case model.ITEM_WIDTH_LONG:
			b.lastItemState = "LL"
			break
		case model.ITEM_WIDTH_SHORT:
			b.lastItemState = "RS"
			break
		default:
			panic("Unhandled state.")
		}
		break
	case "LS":
		switch itemWidth {
		case model.ITEM_WIDTH_FULL:
			b.lastItemState = "F"
			break
		case model.ITEM_WIDTH_HALF:
			b.lastItemState = "RH"
			break
		case model.ITEM_WIDTH_LONG:
			b.lastItemState = "RL"
			break
		case model.ITEM_WIDTH_SHORT:
			b.lastItemState = "MS"
			break
		default:
			panic("Unhandled state.")
		}
		break
	case "MS":
		switch itemWidth {
		case model.ITEM_WIDTH_FULL:
			b.lastItemState = "F"
			break
		case model.ITEM_WIDTH_HALF:
			b.lastItemState = "LH"
			break
		case model.ITEM_WIDTH_LONG:
			b.lastItemState = "LL"
			break
		case model.ITEM_WIDTH_SHORT:
			b.lastItemState = "RS"
			break
		default:
			panic("Unhandled state.")
		}
		break
	default:
		panic("Unhandled state.")
	}
}

type FormScreen struct {
	*ScreenWithFrame
	charset         *model.Charset
	formItemShader  *shader.Shader
	sinceLastClick  float64
	sinceLastDelete float64
	underEdit       interfaces.CharFormItem
	// Item position
	maxScrollOffset     float32
	currentScrollOffset float32
	// map for formItem-configItem
	formItemToConf     map[interfaces.FormItem]*config.ConfigItem
	formItemLabelColor mgl32.Vec3
	formItemInputColor mgl32.Vec3
	clearColor         mgl32.Vec3
	// materials
	formItemDefaultMaterial   *material.Material
	formItemHighlightMaterial *material.Material
}

func (f *FormScreen) setupFormScreen(wrapper interfaces.GLWrapper) {
	col := f.clearColor
	wrapper.ClearColor(col.X(), col.Y(), col.Z(), 1.0)
	wrapper.Enable(glwrapper.DEPTH_TEST)
	wrapper.DepthFunc(glwrapper.LESS)
	wrapper.Enable(glwrapper.BLEND)
	wrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	wrapper.Viewport(0, 0, int32(f.windowWidth), int32(f.windowHeight))
}

// initMaterialForTheFormItems sets the material to the default of the form items.
// It is used in the update loop to push all of them to default state.
func (f *FormScreen) initMaterialForTheFormItems() {
	f.charset.CleanSurface(f.detailContentBox)
	for s, _ := range f.shaderMap {
		for index, _ := range f.shaderMap[s] {
			switch f.shaderMap[s][index].(type) {
			case *model.FormItemInt:
				fi := f.shaderMap[s][index].(*model.FormItemInt)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = f.formItemDefaultMaterial
				break
			case *model.FormItemFloat:
				fi := f.shaderMap[s][index].(*model.FormItemFloat)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = f.formItemDefaultMaterial
				break
			case *model.FormItemText:
				fi := f.shaderMap[s][index].(*model.FormItemText)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = f.formItemDefaultMaterial
				break
			case *model.FormItemBool:
				fi := f.shaderMap[s][index].(*model.FormItemBool)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = f.formItemDefaultMaterial
				lightMesh := fi.GetLight().(*mesh.TexturedMaterialMesh)
				if fi.GetValue() {
					lightMesh.Material = f.formItemHighlightMaterial
				} else {
					lightMesh.Material = f.formItemDefaultMaterial
				}
				break
			case *model.FormItemInt64:
				fi := f.shaderMap[s][index].(*model.FormItemInt64)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = f.formItemDefaultMaterial
				break
			case *model.FormItemVector:
				fi := f.shaderMap[s][index].(*model.FormItemVector)
				surfaceMesh := fi.GetSurface().(*mesh.TexturedMaterialMesh)
				surfaceMesh.Material = f.formItemDefaultMaterial
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
	tmMesh.Material = f.formItemHighlightMaterial
	if f.detailContentBox == nil {
		return
	}
	desc := f.closestModel.(interfaces.FormItem).GetDescription()
	aspRatio := f.GetAspectRatio()
	textScale := InputTextFontScale / f.windowWidth * aspRatio
	// Get the width and height of the 'W' character. The height of the full screen, will be decreased with the width.
	// The height will be used for the vertical positioning.
	wW, hW := f.charset.TextContainerSize("W", textScale)
	lines := f.wrapTextToLines(desc, textScale, f.GetFullWidth()-2*wW)
	for i := 0; i < len(lines); i++ {
		lineVerticalPosition := f.detailContentBoxHeight/aspRatio/2 - float32(i+1)*1.5*hW
		f.charset.PrintTo(lines[i], (-f.GetFullWidth()+wW)/2, lineVerticalPosition, ZText, textScale, f.wrapper, f.detailContentBox, []mgl32.Vec3{f.formItemLabelColor})
	}
}

// Update function increases the time since the last events with dt. Handles the up / down scroll events.
// Sets the form items to their initial state and then handles the mouse events and also the deletion events.
func (f *FormScreen) Update(dt, posX, posY float64, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
	f.sinceLastClick = f.sinceLastClick + dt
	f.sinceLastDelete = f.sinceLastDelete + dt
	aspRatio := f.GetAspectRatio()
	cursorX := float32(-posX) * f.frameWidth / 2
	cursorY := float32(posY) / aspRatio * f.frameWidth / 2
	direction := mgl32.Vec3{0, 0, 0}
	if keyStore.Get(KEY_UP) && !keyStore.Get(KEY_DOWN) {
		direction = mgl32.Vec3{0, -1, 0}
	} else if keyStore.Get(KEY_DOWN) && !keyStore.Get(KEY_UP) {
		direction = mgl32.Vec3{0, 1, 0}
	}
	// If the Up key is pressed => direction: up, velocity: c
	// If the Down key is pressed => direction: down, velocity: c
	// Otherwise => direction: null, velocity c
	newScrollOffset := f.currentScrollOffset + FormItemMoveSpeed*float32(dt)*direction.Y()
	if newScrollOffset > float32(0.0) && newScrollOffset < f.maxScrollOffset {
		f.currentScrollOffset = newScrollOffset
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
				// I know that the input mesh indexes are 1,2,3. I will implement the following logic:
				// If the index gt -1, a field has to be set, if the index is 0, it will use the first index.
				if index > -1 {
					if index == 0 {
						index++
					}
					formModel.SetTarget(index - 1)
					formModel.AddCursor()
					f.underEdit = formModel
				}
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
			aspRatio := f.GetAspectRatio()
			textScale := InputTextFontScale / f.windowWidth * aspRatio
			_, hW := f.charset.TextContainerSize("W", textScale)
			switch f.underEdit.(type) {
			case *model.FormItemVector:
				f.charset.PrintTo(f.underEdit.ValueToString(), -f.underEdit.(*model.FormItemVector).GetVectorCursorInitialPosition().X(), -hW/2, ZText, textScale, f.wrapper, f.underEdit.GetTarget(), []mgl32.Vec3{f.formItemInputColor})
				break
			default:
				f.charset.PrintTo(f.underEdit.ValueToString(), -f.underEdit.GetCursorInitialPosition().X(), -hW/2, ZText, textScale, f.wrapper, f.underEdit.GetTarget(), []mgl32.Vec3{f.formItemInputColor})
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

func (f *FormScreen) addFormItem(fi interfaces.FormItem, defaultValue interface{}) {
	fi.RotateX(-90)
	fi.RotateY(180)
	fi.SetSpeed(FormItemMoveSpeed)
	f.AddModelToShader(fi, f.formItemShader)
	aspRatio := f.GetAspectRatio()
	textScale := LabelFontScale / f.windowWidth * aspRatio
	wW, hW := f.charset.TextContainerSize("W", textScale)
	f.charset.PrintTo(fi.GetLabel(), (-fi.GetFormItemWidth()+wW)/2, -hW/2, ZText, textScale, f.wrapper, fi.GetSurface(), []mgl32.Vec3{f.formItemLabelColor})
	f.SetFormItemValue(fi, defaultValue)
}

// addFormItemFromConfigBool sets up a FormItemBool from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigBool(configItem *config.ConfigItem, pos mgl32.Vec3) {
	asp := 1.0 / f.GetAspectRatio()
	fi := model.NewFormItemBool(f.GetFullWidth(), model.ITEM_WIDTH_SHORT, asp, configItem.GetLabel(), configItem.GetDescription(), f.formItemDefaultMaterial, pos, f.wrapper)
	f.formItemToConf[fi] = configItem
	f.addFormItem(fi, configItem.GetDefaultValue())
}

// addFormItemFromConfigInt sets up a FormItemInt from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigInt(configItem *config.ConfigItem, pos mgl32.Vec3) {
	asp := 1.0 / f.GetAspectRatio()
	fi := model.NewFormItemInt(f.GetFullWidth(), model.ITEM_WIDTH_HALF, asp, configItem.GetLabel(), configItem.GetDescription(), f.formItemDefaultMaterial, pos, f.wrapper)
	if configItem.GetValidatorFunction() != nil {
		fi.SetValidator(configItem.GetValidatorFunction().(model.IntValidator))
	}
	f.formItemToConf[fi] = configItem
	f.addFormItem(fi, transformations.IntegerToString(configItem.GetDefaultValue().(int)))
}

// addFormItemFromConfigFloat sets up a FormItemFloat from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigFloat(configItem *config.ConfigItem, pos mgl32.Vec3) {
	asp := 1.0 / f.GetAspectRatio()
	fi := model.NewFormItemFloat(f.GetFullWidth(), model.ITEM_WIDTH_HALF, asp, configItem.GetLabel(), configItem.GetDescription(), f.formItemDefaultMaterial, pos, f.wrapper)
	if configItem.GetValidatorFunction() != nil {
		fi.SetValidator(configItem.GetValidatorFunction().(model.FloatValidator))
	}
	f.formItemToConf[fi] = configItem
	f.addFormItem(fi, transformations.Float32ToStringExact(configItem.GetDefaultValue().(float32)))
}

// addFormItemFromConfigText sets up a FormItemText from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigText(configItem *config.ConfigItem, pos mgl32.Vec3) {
	asp := 1.0 / f.GetAspectRatio()
	fi := model.NewFormItemText(f.GetFullWidth(), model.ITEM_WIDTH_FULL, asp, configItem.GetLabel(), configItem.GetDescription(), f.formItemDefaultMaterial, pos, f.wrapper)
	if configItem.GetValidatorFunction() != nil {
		fi.SetValidator(configItem.GetValidatorFunction().(model.StringValidator))
	}
	f.formItemToConf[fi] = configItem
	f.addFormItem(fi, configItem.GetDefaultValue())
}

// addFormItemFromConfigInt64 sets up a FormItemInt64 from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigInt64(configItem *config.ConfigItem, pos mgl32.Vec3) {
	asp := 1.0 / f.GetAspectRatio()
	fi := model.NewFormItemInt64(f.GetFullWidth(), model.ITEM_WIDTH_LONG, asp, configItem.GetLabel(), configItem.GetDescription(), f.formItemDefaultMaterial, pos, f.wrapper)
	if configItem.GetValidatorFunction() != nil {
		fi.SetValidator(configItem.GetValidatorFunction().(model.Int64Validator))
	}
	f.formItemToConf[fi] = configItem
	f.addFormItem(fi, transformations.Integer64ToString(configItem.GetDefaultValue().(int64)))
}

// addFormItemFromConfigVector sets up a FormItemInt64 from a ConfigItem structure.
func (f *FormScreen) addFormItemFromConfigVector(configItem *config.ConfigItem, pos mgl32.Vec3) {
	asp := 1.0 / f.GetAspectRatio()
	fi := model.NewFormItemVector(f.GetFullWidth(), model.ITEM_WIDTH_FULL, asp, configItem.GetLabel(), configItem.GetDescription(), model.CHAR_NUM_FLOAT, f.formItemDefaultMaterial, pos, f.wrapper)
	if configItem.GetValidatorFunction() != nil {
		fi.SetValidator(configItem.GetValidatorFunction().(model.FloatValidator))
	}
	f.formItemToConf[fi] = configItem
	f.addFormItem(fi, transformations.VectorToString(configItem.GetDefaultValue().(mgl32.Vec3)))
}

func (f *FormScreen) setDefaultValueChar(input string) {
	chars := []rune(input)
	for i := 0; i < len(chars); i++ {
		f.CharCallback(chars[i], f.wrapper)
	}
}

// CharCallback is the character stream input handler
func (f *FormScreen) CharCallback(char rune, wrapper interfaces.GLWrapper) {
	if f.underEdit != nil {
		aspRatio := f.GetAspectRatio()
		textScale := InputTextFontScale / f.windowWidth * aspRatio
		offsetX := f.charset.TextWidth(string(char), textScale)
		f.underEdit.CharCallback(char, offsetX)
		f.syncFormItemValuesToConfigValue()
		f.charset.CleanSurface(f.underEdit.GetTarget())
		_, hW := f.charset.TextContainerSize("W", textScale)
		switch f.underEdit.(type) {
		case *model.FormItemVector:
			f.charset.PrintTo(f.underEdit.ValueToString(), -f.underEdit.(*model.FormItemVector).GetVectorCursorInitialPosition().X(), -hW/2, ZText, textScale, wrapper, f.underEdit.GetTarget(), []mgl32.Vec3{f.formItemInputColor})
			break
		default:
			f.charset.PrintTo(f.underEdit.ValueToString(), -f.underEdit.GetCursorInitialPosition().X(), -hW/2, ZText, textScale, wrapper, f.underEdit.GetTarget(), []mgl32.Vec3{f.formItemInputColor})
			break
		}
	}
}

// GetFormItem gets a configkey as input, and returns the form item tha connects to the config.
// In case of invalid key, it panics.
func (f *FormScreen) GetFormItem(configKey string) interfaces.FormItem {
	for fi, v := range f.formItemToConf {
		if v.IsConfigOf(configKey) {
			return fi
		}
	}
	panic("Invalid config key:" + configKey)
}

// SetFormItemValue gets a FormItem, a value and sets the value of the form item to the new one.
func (f *FormScreen) SetFormItemValue(item interfaces.FormItem, valueNew interface{}) {
	switch item.(type) {
	case *model.FormItemInt:
		f.underEdit = item.(*model.FormItemInt)
		value := f.underEdit.ValueToString()
		valueLength := len(value) + strings.Count(value, " ")
		for i := 0; i < valueLength; i++ {
			f.underEdit.DeleteLastCharacter()
		}
		f.setDefaultValueChar(valueNew.(string))
		break
	case *model.FormItemFloat:
		f.underEdit = item.(*model.FormItemFloat)
		value := f.underEdit.ValueToString()
		valueLength := len(value) + strings.Count(value, " ")
		for i := 0; i < valueLength; i++ {
			f.underEdit.DeleteLastCharacter()
		}
		f.setDefaultValueChar(valueNew.(string))
		break
	case *model.FormItemText:
		f.underEdit = item.(*model.FormItemText)
		value := f.underEdit.ValueToString()
		valueLength := len(value) + strings.Count(value, " ")
		for i := 0; i < valueLength; i++ {
			f.underEdit.DeleteLastCharacter()
		}
		f.setDefaultValueChar(valueNew.(string))
		break
	case *model.FormItemInt64:
		f.underEdit = item.(*model.FormItemInt64)
		value := f.underEdit.ValueToString()
		valueLength := len(value) + strings.Count(value, " ")
		for i := 0; i < valueLength; i++ {
			f.underEdit.DeleteLastCharacter()
		}
		f.setDefaultValueChar(valueNew.(string))
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
			f.setDefaultValueChar(newValues[i])
		}
		break
	}
}
