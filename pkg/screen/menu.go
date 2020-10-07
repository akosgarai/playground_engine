package screen

import (
	"math"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	DefaultMenuItemMaterial   = material.Whiteplastic
	HighlightMenuItemMaterial = material.Ruby
	MenuItemFontColor         = mgl32.Vec3{1, 1, 1}
)

type DisplayFunction func(map[string]bool) bool
type EventFunction func()

type Option struct {
	label            string
	displayCondition DisplayFunction
	clickEvent       EventFunction
	surface          interfaces.Mesh
}

// NewMenuScreenOption returns an Option. The label, displayCondition and clickEvent are
// set to the input values.
func NewMenuScreenOption(l string, dc DisplayFunction, ce EventFunction) *Option {
	return &Option{
		label:            l,
		displayCondition: dc,
		clickEvent:       ce,
	}
}

// SetSurface sets the surface mesh to the given one.
func (o *Option) SetSurface(surface interfaces.Mesh) {
	o.surface = surface
}

// DisplayCondition calls the display function and returns the result
func (o *Option) DisplayCondition(state map[string]bool) bool {
	return o.displayCondition(state)
}

type MenuScreenBuilder struct {
	*ScreenWithFrameBuilder
	charset                 *model.Charset
	menuItemSurfaceTexture  texture.Textures
	menuItemDefaultMaterial *material.Material
	menuItemHoverMaterial   *material.Material
	options                 []Option
	menuItemFontColor       mgl32.Vec3
	backgroundColor         mgl32.Vec3
	state                   map[string]bool
}

// NewMenuScreenBuilder returns a MenuScreenBuilder that we cen use for creating a MenuScreen.
func NewMenuScreenBuilder() *MenuScreenBuilder {
	swfb := NewScreenWithFrameBuilder()
	swfb.SetFrameSize(DefaultFrameWidth, DefaultFrameLength, TopLeftFrameWidth)
	swfb.SetFrameMaterial(HighlightMenuItemMaterial)
	return &MenuScreenBuilder{
		ScreenWithFrameBuilder:  swfb,
		charset:                 nil,
		menuItemDefaultMaterial: DefaultMenuItemMaterial,
		menuItemHoverMaterial:   HighlightMenuItemMaterial,
		menuItemFontColor:       MenuItemFontColor,
		backgroundColor:         MenuItemFontColor,
		state:                   make(map[string]bool),
	}
}

func (b *MenuScreenBuilder) defaultCharset() {
	cs, err := model.LoadCharset(baseDirScreen()+DefaultFontFile, 32, 127, 40.0, 72, b.wrapper)
	if err != nil {
		panic(err)
	}
	cs.SetTransparent(true)
	b.charset = cs
}

// SetMenuItemSurfaceTexture sets the texture of the menu item.
func (b *MenuScreenBuilder) SetMenuItemSurfaceTexture(tex texture.Textures) {
	b.menuItemSurfaceTexture = tex
}

// SetMenuItemDefaultMaterial sets the default material of the menu item.
func (b *MenuScreenBuilder) SetMenuItemDefaultMaterial(mat *material.Material) {
	b.menuItemDefaultMaterial = mat
}

// SetMenuItemHighlightMaterial sets the hover material of the menu item.
func (b *MenuScreenBuilder) SetMenuItemHighlightMaterial(mat *material.Material) {
	b.menuItemHoverMaterial = mat
}

// SetMenuItemFontColor sets the font color of the menu item.
func (b *MenuScreenBuilder) SetMenuItemFontColor(c mgl32.Vec3) {
	b.menuItemFontColor = c
}

// SetBackgroundColor sets the background color of the menu screen.
func (b *MenuScreenBuilder) SetBackgroundColor(c mgl32.Vec3) {
	b.backgroundColor = c
}

// SetState updates the current state of the screen.
func (b *MenuScreenBuilder) SetState(m map[string]bool) {
	b.state = m
}

// SetCharset sets the charset of the form screen.
func (b *MenuScreenBuilder) SetCharset(m *model.Charset) {
	b.charset = m
}

// AddOption appens the new option to the end of the option list.
func (b *MenuScreenBuilder) AddOption(o Option) {
	b.options = append(b.options, o)
}
func (b *MenuScreenBuilder) Build() *MenuScreen {
	if b.wrapper == nil {
		panic("Wrapper is missing for the build process.")
	}
	if b.charset == nil {
		b.defaultCharset()
	}
	s := b.ScreenWithFrameBuilder.Build()
	bgShaderApplication := shader.NewMenuBackgroundShader(b.wrapper)
	fgShaderApplication := shader.NewFontShader(b.wrapper)
	s.AddShader(bgShaderApplication)
	s.AddShader(fgShaderApplication)
	s.AddModelToShader(b.charset, fgShaderApplication)
	background := model.New()
	s.AddModelToShader(background, bgShaderApplication)
	menuScreen := &MenuScreen{
		ScreenWithFrame:  s,
		surfaceTexture:   b.menuItemSurfaceTexture,
		defaultMaterial:  b.menuItemDefaultMaterial,
		hoverMaterial:    b.menuItemHoverMaterial,
		charset:          b.charset,
		background:       background,
		options:          b.options,
		backgroundShader: bgShaderApplication,
		fontShader:       fgShaderApplication,
		fontColor:        []mgl32.Vec3{b.menuItemFontColor},
		backgroundColor:  b.backgroundColor,
		state:            b.state,
		surfaceToOption:  make(map[interfaces.Mesh]Option),
	}
	s.Setup(menuScreen.setupMenu)
	menuScreen.BuildScreen()
	return menuScreen
}

type MenuScreen struct {
	*ScreenWithFrame
	charset             *model.Charset
	background          interfaces.Model
	surfaceTexture      texture.Textures
	defaultMaterial     *material.Material
	hoverMaterial       *material.Material
	options             []Option
	backgroundShader    interfaces.Shader
	fontShader          interfaces.Shader
	fontColor           []mgl32.Vec3
	backgroundColor     mgl32.Vec3
	state               map[string]bool
	surfaceToOption     map[interfaces.Mesh]Option
	maxScrollOffset     float32
	currentScrollOffset float32
}

// SetState is the state maintainer function
func (m *MenuScreen) SetState(key string, value bool) {
	m.state[key] = value
}

// BuildScreen function sets the screen up based on the option conditions.
func (m *MenuScreen) BuildScreen() {
	// clear prev. screen
	m.charset.Clear()
	m.background.Clear()
	optionsToDisplay := m.getOptionsToDisplay()
	// variables for aspect ratio.
	aspRatio := m.GetAspectRatio()
	windowWidth, _ := m.GetWindowSize()

	positionY := -0.8 * m.frameWidth / 2.0 / aspRatio
	topOfTheBottomForegroundArea := (-(m.frameWidth / 2.0) + m.frameLength + m.detailContentBoxHeight)
	positionOfTheBottomMenuItem := -0.8 * m.frameWidth / 4.0
	bottomOffset := (positionOfTheBottomMenuItem - topOfTheBottomForegroundArea) * aspRatio
	if bottomOffset > 0 {
		bottomOffset = 0
	}
	positionX := positionOfTheBottomMenuItem
	surfaceToOption := make(map[interfaces.Mesh]Option)
	for i := len(optionsToDisplay) - 1; i >= 0; i-- {
		surface := m.menuSurface(mgl32.Vec3{0.0, positionY, ZBackground})
		optionsToDisplay[i].SetSurface(surface)
		surfaceToOption[surface] = optionsToDisplay[i]
		_, textHeight := m.charset.TextContainerSize(optionsToDisplay[i].label, 3.0/windowWidth*aspRatio)
		m.charset.PrintTo(optionsToDisplay[i].label, positionX, -textHeight/2.0, ZText, 3.0/windowWidth*aspRatio, m.wrapper, optionsToDisplay[i].surface, m.fontColor)
		positionY += m.frameWidth * 0.2 / aspRatio
	}
	bottomOfTheTopForegroundArea := (m.frameWidth/2.0 - m.frameLength) / aspRatio
	topOffset := positionY - bottomOfTheTopForegroundArea
	if topOffset < 0 {
		topOffset = 0
	}
	m.maxScrollOffset = bottomOffset + topOffset
	m.currentScrollOffset = bottomOffset
	m.surfaceToOption = surfaceToOption
}

// It returns the options that needs to be displayed in the current state.
func (m *MenuScreen) getOptionsToDisplay() []Option {
	var result []Option
	for i := 0; i < len(m.options); i++ {
		if m.options[i].displayCondition != nil {
			if r := m.options[i].displayCondition(m.state); r {
				result = append(result, m.options[i])
			}
		}
	}
	return result
}

// menuSurface creates a rectangle for the menu option.
func (m *MenuScreen) menuSurface(position mgl32.Vec3) interfaces.Mesh {
	aspRatio := m.GetAspectRatio()
	menuWidth := m.frameWidth / 2.0
	menuHeight := menuWidth / 5 / aspRatio
	rect := rectangle.NewExact(menuWidth, menuHeight)
	v, i, bo := rect.MeshInput()
	msh := mesh.NewTexturedMaterialMesh(v, i, m.surfaceTexture, m.defaultMaterial, m.wrapper)
	msh.SetBoundingObject(bo)
	msh.SetPosition(position)
	msh.RotateX(-90)
	msh.RotateY(180)
	m.background.AddMesh(msh)
	return msh
}
func (m *MenuScreen) setupMenu(glWrapper interfaces.GLWrapper) {
	glWrapper.Enable(glwrapper.DEPTH_TEST)
	glWrapper.DepthFunc(glwrapper.LESS)
	glWrapper.Enable(glwrapper.BLEND)
	glWrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	col := m.backgroundColor
	glWrapper.ClearColor(col.X(), col.Y(), col.Z(), 1.0)
	glWrapper.Viewport(0, 0, int32(m.windowWidth), int32(m.windowHeight))
}

// Update loops on the shaderMap, and calls Update function on every Model.
// It also handles the camera movement and rotation, if the camera is set.
func (s *MenuScreen) Update(dt, posX, posY float64, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
	aspRatio := s.GetAspectRatio()
	cursorX := float32(-posX) * s.frameWidth / 2
	cursorY := float32(posY) / aspRatio * s.frameWidth / 2
	direction := mgl32.Vec3{0, 0, 0}
	if keyStore.Get(KEY_UP) && !keyStore.Get(KEY_DOWN) {
		direction = mgl32.Vec3{0, -1, 0}
	} else if keyStore.Get(KEY_DOWN) && !keyStore.Get(KEY_UP) {
		direction = mgl32.Vec3{0, 1, 0}
	}
	newScrollOffset := s.currentScrollOffset + FormItemMoveSpeed*float32(dt)*direction.Y()
	if newScrollOffset > float32(0.0) && newScrollOffset < s.maxScrollOffset {
		s.currentScrollOffset = newScrollOffset
	} else {
		direction = mgl32.Vec3{0, 0, 0}
	}
	for m, _ := range s.shaderMap[s.backgroundShader] {
		s.shaderMap[s.backgroundShader][m].SetDirection(direction)
	}
	coords := mgl32.Vec3{cursorX, cursorY, ZBackground}

	closestDistance := float32(math.MaxFloat32)
	var closestMesh interfaces.Mesh
	s.closestModel = s.background
	// Here we only need to check the background shader, get the closest stuff, and check the distance
	for index, _ := range s.shaderMap[s.backgroundShader] {
		s.shaderMap[s.backgroundShader][index].Update(dt)
		msh, dist := s.shaderMap[s.backgroundShader][index].ClosestMeshTo(coords)
		if dist < closestDistance {
			closestDistance = dist
			closestMesh = msh
		}
	}
	// Update the material in case of hover state.
	s.closestMesh = closestMesh
	s.closestDistance = closestDistance

	switch s.closestMesh.(type) {
	case *mesh.TexturedMaterialMesh:
		tmMesh := s.closestMesh.(*mesh.TexturedMaterialMesh)
		if s.closestDistance < 0.01 {
			tmMesh.Material = s.hoverMaterial
			if buttonStore.Get(LEFT_MOUSE_BUTTON) && s.surfaceToOption[s.closestMesh].clickEvent != nil {
				s.surfaceToOption[s.closestMesh].clickEvent()
			}
		} else {
			tmMesh.Material = s.defaultMaterial
		}
	}
}

// CharCallback is the character stream input handler
func (s *MenuScreen) CharCallback(char rune, w interfaces.GLWrapper) {
}
