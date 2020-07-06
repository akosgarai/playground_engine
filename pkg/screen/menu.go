package screen

import (
	"fmt"
	"math"
	"reflect"

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

type MenuScreen struct {
	*ScreenBase
	charset          *model.Charset
	background       interfaces.Model
	surfaceTexture   texture.Textures
	defaultMaterial  *material.Material
	hoverMaterial    *material.Material
	options          []Option
	backgroundShader interfaces.Shader
	fontShader       interfaces.Shader
	fontColor        []mgl32.Vec3
	backgroundColor  mgl32.Vec3
	state            map[string]bool
}

// NewMenuScreen returns a MenuScreen without options.
func NewMenuScreen(surface texture.Textures, defaultMat *material.Material, hoverMat *material.Material, charset *model.Charset, fontColor mgl32.Vec3, backgroundColor mgl32.Vec3, wrapper interfaces.GLWrapper) *MenuScreen {
	s := newScreenBase()
	bgShaderApplication := shader.NewMenuBackgroundShader(wrapper)
	fgShaderApplication := shader.NewFontShader(wrapper)
	s.AddShader(bgShaderApplication)
	s.AddShader(fgShaderApplication)
	s.AddModelToShader(charset, fgShaderApplication)
	background := model.New()
	s.AddModelToShader(background, bgShaderApplication)
	state := make(map[string]bool)
	state["world-started"] = false
	menuScreen := &MenuScreen{
		ScreenBase:       s,
		surfaceTexture:   surface,
		defaultMaterial:  defaultMat,
		hoverMaterial:    hoverMat,
		charset:          charset,
		background:       background,
		options:          []Option{},
		backgroundShader: bgShaderApplication,
		fontShader:       fgShaderApplication,
		fontColor:        []mgl32.Vec3{fontColor},
		backgroundColor:  backgroundColor,
		state:            state,
	}
	s.Setup(menuScreen.setupMenu)
	return menuScreen
}

// BuildScreen function sets the screen up based on the option conditions.
func (m *MenuScreen) BuildScreen(wrapper interfaces.GLWrapper, scale float32) {
	// clear prev. screen
	m.charset.Clear()
	m.background.Clear()
	optionsToDisplay := m.getOptionsToDisplay()
	// The position of the buttons needs to be calculated based on the number
	// of the options. The button has to be maximum 1.5 * height of the fonts.
	// We need a padding from the top and also from the bottom.
	// Width: [-1,1] -> 2.
	// Height: [-1,1] -> 2.
	// Padding from left / right: 0.2-0.5
	// MenuSurfaceWidth: 1.0
	positionY := float32(-0.8)
	positionX := float32(-0.4)
	positionZ := float32(0.01)
	for i := len(optionsToDisplay) - 1; i >= 0; i-- {
		optionsToDisplay[i].SetSurface(m.menuSurface(mgl32.Vec3{0.0, positionY, 0.0}, wrapper))
		m.charset.PrintTo(optionsToDisplay[i].label, positionX, -0.03, positionZ, scale, wrapper, optionsToDisplay[i].surface, m.fontColor)
		positionY += float32(0.4)
	}
}

// AddOption appens the new option to the end of the option list.
func (m *MenuScreen) AddOption(o Option) {
	m.options = append(m.options, o)
}

// It returns the options that needs to be displayed in the current state.
func (m *MenuScreen) getOptionsToDisplay() []Option {
	var result []Option
	for i := 0; i < len(m.options); i++ {
		if r := m.options[i].displayCondition(m.state); r {
			result = append(result, m.options[i])
		}
	}
	return result
}

// menuSurface creates a rectangle for the menu option.
func (m *MenuScreen) menuSurface(position mgl32.Vec3, wrapper interfaces.GLWrapper) interfaces.Mesh {
	menuWidth := float32(1.0)
	menuHeight := float32(0.2)
	rect := rectangle.NewExact(menuWidth, menuHeight)
	v, i, bo := rect.MeshInput()
	msh := mesh.NewTexturedMaterialMesh(v, i, m.surfaceTexture, m.defaultMaterial, wrapper)
	msh.SetBoundingObject(bo)
	msh.SetPosition(position)
	msh.RotateX(-90)
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
}

// Update loops on the shaderMap, and calls Update function on every Model.
// It also handles the camera movement and rotation, if the camera is set.
func (s *MenuScreen) Update(dt, posX, posY float64, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
	TransformationMatrix := mgl32.Ident4()
	if s.cameraSet {
		s.cameraKeyboardMovement("up", "down", "Lift", dt, keyStore)
		TransformationMatrix = (s.camera.GetProjectionMatrix().Mul4(s.camera.GetViewMatrix())).Inv()
	}
	coords := mgl32.TransformCoordinate(mgl32.Vec3{float32(posX), float32(posY), 0.0}, TransformationMatrix)
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

	tmMesh := s.closestMesh.(*mesh.TexturedMaterialMesh)
	if s.closestDistance < 0.01 {
		tmMesh.Material = s.hoverMaterial
		if buttonStore.Get(LEFT_MOUSE_BUTTON) {
			for i, _ := range s.options {
				fmt.Printf("&(s.options[i].surface): %v | &(s.closestMesh): %v | &tmMesh %v\n", &(s.options[i].surface), &(s.closestMesh), &tmMesh)
				fmt.Printf("s.options[i].surface == s.closestMesh (%v) | s.options[i].surface == tmMesh (%v) | s.closestMesh == tmMesh (%v) | reflect.DeepEqual(s.options[i].surface, s.closestMesh) (%v)\n", s.options[i].surface == s.closestMesh, s.options[i].surface == tmMesh, s.closestMesh == tmMesh, reflect.DeepEqual(s.options[i].surface, s.closestMesh))
				if &(s.options[i].surface) == &(s.closestMesh) {
					fmt.Println("Surface has been found.")
					s.options[i].clickEvent()
				}
			}
		}
	} else {
		tmMesh.Material = s.defaultMaterial
	}
}
