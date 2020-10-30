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
	"github.com/akosgarai/playground_engine/pkg/model/ui"
	"github.com/akosgarai/playground_engine/pkg/primitives/cuboid"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	defaultLeftMonitorRotationAngle  = float32(60)
	defaultRightMonitorRotationAngle = float32(-60)
	defaultMonitorTextureName        = "crt_monitor_1280.png"
	defaultBackgroundTextureName     = "rusty_iron_1280.jpg"
)

var (
	defaultAssetsDirectory       string
	defaultMiddleMonitorPosition = mgl32.Vec3{2.0, 0, 0}
	defaultMonitorSize           = mgl32.Vec2{2.0, 2.0}
	defaultMiddleScreenPosition  = mgl32.Vec3{-0.02, 0.0, 0.07} //relative from the middle monitor position.
	defaultLeftScreenPosition    = mgl32.TransformCoordinate(defaultMiddleScreenPosition, mgl32.HomogRotate3DZ(mgl32.DegToRad(defaultLeftMonitorRotationAngle)))
	defaultRightScreenPosition   = mgl32.TransformCoordinate(defaultMiddleScreenPosition, mgl32.HomogRotate3DZ(mgl32.DegToRad(defaultRightMonitorRotationAngle)))
	defaultScreenSize            = mgl32.Vec3{1.5, 1.3, 0}
	defaultScreenMaterial        = material.Emerald
	defaultTableMaterial         = material.Chrome
	defaultTableSize             = mgl32.Vec3{2.5, 6, 0.05}
	defaultTablePosition         = mgl32.Vec3{1.5, 0, -1}
	defaultBackgroundSize        = mgl32.Vec3{8, 8, 8}
	defaultControlPoints         = []mgl32.Vec3{mgl32.Vec3{-3.5, 0, 1.2}, mgl32.Vec3{0, 0, 1.2}, mgl32.Vec3{0, 0, 0}}
	defaultClearColor            = mgl32.Vec3{0.0, 0.25, 0.5}
	defaultLightDir              = mgl32.Vec3{0.0, 0.0, -1.0}
	defaultLightAmbient          = mgl32.Vec3{1.0, 1.0, 1.0}
	defaultLightDiffuse          = mgl32.Vec3{1.0, 1.0, 1.0}
	defaultLightSpecular         = mgl32.Vec3{1.0, 1.0, 1.0}
)

func init() {
	defaultAssetsDirectory = baseDirScreen() + "/assets/"
}

type CubeFormScreenBuilder struct {
	// The position in world coordinates of the middle form container.
	middleMonitorPosition mgl32.Vec3
	// The rotation degree of the left form container.
	// The actual position will be the middle one's rotated with this angle.
	leftMonitorRotationAngle float32
	// The rotation degree of the right form container.
	// The actual position will be the middle one's rotated with this angle.
	rightMonitorRotationAngle float32
	// The size of the form containers.
	middleMonitorSize mgl32.Vec2
	leftMonitorSize   mgl32.Vec2
	rightMonitorSize  mgl32.Vec2
	// The position (relative from the monitor) of the screen.
	middleScreenPosition mgl32.Vec3
	leftScreenPosition   mgl32.Vec3
	rightScreenPosition  mgl32.Vec3
	// The size vector of the screen.
	middleScreenSize mgl32.Vec3
	leftScreenSize   mgl32.Vec3
	rightScreenSize  mgl32.Vec3
	// The path to the assets dir. If it is empty, the default value will be used instead of this.
	assetsDirectory string
	// The files of the container textures.
	middleMonitorTexture string
	leftMonitorTexture   string
	rightMonitorTexture  string
	// the material of the screen
	screenMaterial *material.Material
	// The mesh that represents the table.
	tableMaterial *material.Material
	tableSize     mgl32.Vec3
	tablePosition mgl32.Vec3
	// The background
	backgroundSize    mgl32.Vec3
	backgroundTexture string
	// The size of the window
	windowWidth  float32
	windowHeight float32
	// gl wrapper
	wrapper interfaces.GLWrapper
	// camera
	camera interfaces.Camera
	// control points for the initial camera movement animation
	controlPoints []mgl32.Vec3
	// clear color of the screen
	clearColor mgl32.Vec3
	// directional light source
	lightDir      mgl32.Vec3
	lightAmbient  mgl32.Vec3
	lightDiffuse  mgl32.Vec3
	lightSpecular mgl32.Vec3
}

// It returns a builder instance that holds the default values.
func NewCubeFormScreenBuilder() *CubeFormScreenBuilder {
	return &CubeFormScreenBuilder{
		middleMonitorPosition:     defaultMiddleMonitorPosition,
		leftMonitorRotationAngle:  defaultLeftMonitorRotationAngle,
		rightMonitorRotationAngle: defaultRightMonitorRotationAngle,
		middleMonitorSize:         defaultMonitorSize,
		leftMonitorSize:           defaultMonitorSize,
		rightMonitorSize:          defaultMonitorSize,
		middleScreenPosition:      defaultMiddleScreenPosition,
		leftScreenPosition:        defaultLeftScreenPosition,
		rightScreenPosition:       defaultRightScreenPosition,
		middleScreenSize:          defaultScreenSize,
		leftScreenSize:            defaultScreenSize,
		rightScreenSize:           defaultScreenSize,
		assetsDirectory:           defaultAssetsDirectory,
		middleMonitorTexture:      defaultMonitorTextureName,
		leftMonitorTexture:        defaultMonitorTextureName,
		rightMonitorTexture:       defaultMonitorTextureName,
		screenMaterial:            defaultScreenMaterial,
		tableMaterial:             defaultTableMaterial,
		tableSize:                 defaultTableSize,
		tablePosition:             defaultTablePosition,
		backgroundSize:            defaultBackgroundSize,
		backgroundTexture:         defaultBackgroundTextureName,
		wrapper:                   nil,
		camera:                    nil,
		controlPoints:             defaultControlPoints,
		clearColor:                defaultClearColor,
		lightDir:                  defaultLightDir,
		lightAmbient:              defaultLightAmbient,
		lightDiffuse:              defaultLightDiffuse,
		lightSpecular:             defaultLightSpecular,
	}
}

// SetMiddleMonitorPosition updates the middleMonitorPosition with the new value.
func (b *CubeFormScreenBuilder) SetMiddleMonitorPosition(p mgl32.Vec3) {
	b.middleMonitorPosition = p
}

// SetMonitorRotationAngles updates the left and right monitor rotation angles.
func (b *CubeFormScreenBuilder) SetMonitorRotationAngles(left, right float32) {
	b.leftMonitorRotationAngle = left
	b.rightMonitorRotationAngle = right
}

// SetScreenPositions updates the positions of the screens.
func (b *CubeFormScreenBuilder) SetScreenPositions(left, middle, right mgl32.Vec3) {
	b.middleScreenPosition = middle
	b.leftScreenPosition = left
	b.rightScreenPosition = right
}

// SetScreenSizes updates the size of the screens.
func (b *CubeFormScreenBuilder) SetScreenSizes(left, middle, right mgl32.Vec3) {
	b.middleScreenSize = middle
	b.leftScreenSize = left
	b.rightScreenSize = right
}

// SetAssetsDirectory updates the assetsDirectory to the given one.
func (b *CubeFormScreenBuilder) SetAssetsDirectory(path string) {
	b.assetsDirectory = path
}

// SetMonitorTextureNames updates the texture names of the monitors.
func (b *CubeFormScreenBuilder) SetMonitorTextureNames(pathLeft, pathMiddle, pathRight string) {
	b.middleMonitorTexture = pathMiddle
	b.leftMonitorTexture = pathLeft
	b.rightMonitorTexture = pathRight
}

// SetScreenMaterial updates the material of the screen surface.
func (b *CubeFormScreenBuilder) SetScreenMaterial(m *material.Material) {
	b.screenMaterial = m
}

// SetTableMaterial updates the material of the table surface.
func (b *CubeFormScreenBuilder) SetTableMaterial(m *material.Material) {
	b.tableMaterial = m
}

// SetTableSize updates the tableSize.
func (b *CubeFormScreenBuilder) SetTableSize(s mgl32.Vec3) {
	b.tableSize = s
}

// SetTablePosition updates the tablePosition.
func (b *CubeFormScreenBuilder) SetTablePosition(p mgl32.Vec3) {
	b.tablePosition = p
}

// SetBackgroundSize updates the backgroundSize.
func (b *CubeFormScreenBuilder) SetBackgroundSize(s mgl32.Vec3) {
	b.backgroundSize = s
}

// SetBackgroundTextureName updates the backgroundTexture to the given one.
func (b *CubeFormScreenBuilder) SetBackgroundTextureName(n string) {
	b.backgroundTexture = n
}

// SetWindowSize sets the windowWidth and the windowHeight parameters.
func (b *CubeFormScreenBuilder) SetWindowSize(w, h float32) {
	b.windowWidth = w
	b.windowHeight = h
}

// SetWrapper sets the wrapper.
func (b *CubeFormScreenBuilder) SetWrapper(w interfaces.GLWrapper) {
	b.wrapper = w
}

// SetControlPoints updates the control points of the initial animation.
func (b *CubeFormScreenBuilder) SetControlPoints(cp []mgl32.Vec3) {
	b.controlPoints = cp
}

// SetClearColor sets the color of the background.
func (b *CubeFormScreenBuilder) SetClearColor(c mgl32.Vec3) {
	b.clearColor = c
}

// Build returns the CubeFormScreen. In case of missing wrapper, it panics.
func (b *CubeFormScreenBuilder) Build() *CubeFormScreen {
	// check the wrapper
	if b.wrapper == nil {
		panic("Wrapper is missing for the build process.")
	}
	// new screen base with camera
	sb := newScreenBase()
	sb.SetupCamera(b.defaultCamera(), b.defaultCameraOptions())
	// the models
	monitors := model.New()
	desk := model.New()
	rustySurface := model.New()
	screens := model.New()
	// setup the middle monitor
	mm := b.createTexturedRectangle(b.getTexture(b.assetsDirectory+"/"+b.middleMonitorTexture), b.middleMonitorSize.X(), b.middleMonitorSize.Y())
	mm.SetPosition(b.middleMonitorPosition)
	mm.RotateZ(90)
	monitors.AddMesh(mm)
	// middle screen
	middleMonitorScreen := b.createMaterialCube(b.screenMaterial, b.middleScreenSize)
	middleMonitorScreen.SetParent(mm)
	middleMonitorScreen.SetPosition(b.middleScreenPosition)
	screens.AddMesh(middleMonitorScreen)
	// left monitor
	lm := b.createTexturedRectangle(b.getTexture(b.assetsDirectory+"/"+b.leftMonitorTexture), b.leftMonitorSize.X(), b.leftMonitorSize.Y())
	lm.SetPosition(mgl32.TransformCoordinate(b.middleMonitorPosition, mgl32.HomogRotate3DZ(mgl32.DegToRad(b.leftMonitorRotationAngle))))
	lm.RotateZ(90 + b.leftMonitorRotationAngle)
	monitors.AddMesh(lm)
	// left screen
	leftMonitorScreen := b.createMaterialCube(b.screenMaterial, b.leftScreenSize)
	leftMonitorScreen.SetParent(lm)
	leftMonitorScreen.SetPosition(b.leftScreenPosition)
	screens.AddMesh(leftMonitorScreen)
	// right monitor
	rm := b.createTexturedRectangle(b.getTexture(b.assetsDirectory+"/"+b.rightMonitorTexture), b.rightMonitorSize.X(), b.rightMonitorSize.Y())
	rm.SetPosition(mgl32.TransformCoordinate(b.middleMonitorPosition, mgl32.HomogRotate3DZ(mgl32.DegToRad(b.rightMonitorRotationAngle))))
	rm.RotateZ(90 + b.rightMonitorRotationAngle)
	monitors.AddMesh(rm)
	// right screen
	rightMonitorScreen := b.createMaterialCube(b.screenMaterial, b.rightScreenSize)
	rightMonitorScreen.SetParent(rm)
	rightMonitorScreen.SetPosition(b.rightScreenPosition)
	screens.AddMesh(rightMonitorScreen)
	// table surface
	tableSurface := b.createMaterialCube(b.tableMaterial, b.tableSize)
	tableSurface.SetPosition(b.tablePosition)
	tableSurface.RotateX(90)
	desk.AddMesh(tableSurface)
	// background - the room
	rustySurface.AddMesh(b.createTexturedCube(b.getTexture(b.assetsDirectory+"/"+b.backgroundTexture), b.backgroundSize))
	monitors.SetTransparent(true)
	// define the shader
	shaderAppMonitors := shader.NewTextureShaderBlending(b.wrapper)
	sb.AddShader(shaderAppMonitors)
	shaderAppMaterial := shader.NewMaterialShader(b.wrapper)
	sb.AddShader(shaderAppMaterial)
	shaderAppTexture := shader.NewTextureShader(b.wrapper)
	sb.AddShader(shaderAppTexture)
	shaderAppButtons := shader.NewTextureMatShaderBlending(b.wrapper)
	sb.AddShader(shaderAppButtons)
	// Attach the models to the shaders
	sb.AddModelToShader(monitors, shaderAppMonitors)
	sb.AddModelToShader(desk, shaderAppMaterial)
	sb.AddModelToShader(screens, shaderAppMaterial)
	sb.AddModelToShader(rustySurface, shaderAppTexture)
	// The buttons below are only for testing purposes. The final layout will be designed later.
	bb := ui.NewUIButtonBuilder(b.wrapper)
	bb.SetupSize(0.5, 1.0, 0.05)
	bb.SetupMaterials(material.Chrome, material.Ruby, material.Emerald)
	button1 := bb.Build()
	button1.AttachToScreen(middleMonitorScreen, mgl32.Vec3{-0.002, 0.35, 0.0})
	sb.AddModelToShader(button1, shaderAppButtons)
	// Rotate right button
	MiddleMonitorGoRightButton, _ := button1.GetMeshByIndex(0)
	// Rotate left button
	button2 := bb.Build()
	button2.AttachToScreen(middleMonitorScreen, mgl32.Vec3{-0.002, -0.35, 0.0})
	sb.AddModelToShader(button2, shaderAppButtons)
	MiddleMonitorGoLeftButton, _ := button2.GetMeshByIndex(0)
	// right go back button
	button3 := bb.Build()
	button3.AttachToScreen(rightMonitorScreen, mgl32.Vec3{-0.002, 0.35, 0.0})
	sb.AddModelToShader(button3, shaderAppButtons)
	RightMonitorGoLeftButton, _ := button3.GetMeshByIndex(0)
	// left go back button
	button4 := bb.Build()
	button4.AttachToScreen(leftMonitorScreen, mgl32.Vec3{-0.002, -0.35, 0.0})
	sb.AddModelToShader(button4, shaderAppButtons)
	LeftMonitorGoRightButton, _ := button4.GetMeshByIndex(0)
	s := &CubeFormScreen{
		ScreenBase:                 sb,
		MiddleMonitorGoRightButton: MiddleMonitorGoRightButton.(*mesh.TexturedMaterialMesh),
		MiddleMonitorGoLeftButton:  MiddleMonitorGoLeftButton.(*mesh.TexturedMaterialMesh),
		RightMonitorGoLeftButton:   RightMonitorGoLeftButton.(*mesh.TexturedMaterialMesh),
		LeftMonitorGoRightButton:   LeftMonitorGoRightButton.(*mesh.TexturedMaterialMesh),
		RotationToLeftAngle:        b.leftMonitorRotationAngle,
		RotationToRightAngle:       b.rightMonitorRotationAngle,
		SumOfRotation:              float32(0.0),
		middleMonitorPosition:      b.middleMonitorPosition,
		currentRotation:            float32(0.0),
		state:                      "initial",
		controlPoints:              b.controlPoints,
		clearColor:                 b.clearColor,
	}
	DirectionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		b.lightDir,
		b.lightAmbient,
		b.lightDiffuse,
		b.lightSpecular,
	})
	// Add the lightources to the application
	sb.AddDirectionalLightSource(DirectionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})
	sb.Setup(s.setupCubeFormScreen)
	return s
}

// getTexture has a filepath as input and returns the Textures that we can use for the meshes later.
func (b *CubeFormScreenBuilder) getTexture(path string) texture.Textures {
	var t texture.Textures
	t.AddTexture(path, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.diffuse", b.wrapper)
	t.AddTexture(path, glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "tex.specular", b.wrapper)
	return t
}

// Create a textured rectangle, that represents the monitors. The texture and the size are the inputs.
// It returns the textured mesh without bounding object setup.
func (b *CubeFormScreenBuilder) createTexturedRectangle(t texture.Textures, w, h float32) *mesh.TexturedMesh {
	r := rectangle.NewExact(w, h)
	V, I, _ := r.MeshInput()
	return mesh.NewTexturedMesh(V, I, t, b.wrapper)
}

// Create material cube, that represents the table and the screens.
func (b *CubeFormScreenBuilder) createMaterialCube(m *material.Material, size mgl32.Vec3) *mesh.MaterialMesh {
	r := cuboid.New(size.X(), size.Y(), size.Z())
	V, I, _ := r.MaterialMeshInput()
	return mesh.NewMaterialMesh(V, I, m, b.wrapper)
}

// Create a textured rectangle, that represents the walls.
func (b *CubeFormScreenBuilder) createTexturedCube(t texture.Textures, size mgl32.Vec3) *mesh.TexturedMesh {
	r := cuboid.New(size.X(), size.Y(), size.Z())
	V, I, _ := r.TexturedMeshInput(cuboid.TEXTURE_ORIENTATION_DEFAULT)
	return mesh.NewTexturedMesh(V, I, t, b.wrapper)
}

// It creates a new camera with the necessary setup
func (b *CubeFormScreenBuilder) defaultCamera() *camera.DefaultCamera {
	cam := camera.NewCamera(b.controlPoints[0], mgl32.Vec3{0, 0, -1}, 0.0, 0.0)
	cam.SetupProjection(45, b.windowWidth/b.windowHeight, 0.001, 10.0)
	cam.SetVelocity(float32(0.005))
	cam.SetRotationStep(float32(0.125))
	return cam
}

// Setup options for the camera.
// TODO: The movement supposed to be disabled after the screen is tested out well.
func (b *CubeFormScreenBuilder) defaultCameraOptions() map[string]interface{} {
	cm := make(map[string]interface{})
	cm["mode"] = "default"
	cm["rotateOnEdgeDistance"] = float32(0.0)
	cm["forward"] = []glfw.Key{glfw.KeyW}
	cm["back"] = []glfw.Key{glfw.KeyS}
	cm["up"] = []glfw.Key{glfw.KeyQ}
	cm["down"] = []glfw.Key{glfw.KeyE}
	cm["left"] = []glfw.Key{glfw.KeyA}
	cm["right"] = []glfw.Key{glfw.KeyD}
	return cm
}

type CubeFormScreen struct {
	*ScreenBase
	MiddleMonitorGoRightButton *mesh.TexturedMaterialMesh
	MiddleMonitorGoLeftButton  *mesh.TexturedMaterialMesh
	RightMonitorGoLeftButton   *mesh.TexturedMaterialMesh
	LeftMonitorGoRightButton   *mesh.TexturedMaterialMesh
	// Rotation to the monitors
	RotationToLeftAngle  float32
	RotationToRightAngle float32
	SumOfRotation        float32
	// for the click handling, i want to be able to detect the screen plane.
	middleMonitorPosition mgl32.Vec3
	currentRotation       float32
	// The state of the screen. Possible values: 'initial', 'move-to-place',
	// 'form', 'rotate-left', 'rotate-right'.
	state string
	// the control points for the initial animation.
	controlPoints     []mgl32.Vec3
	controlPointIndex int
	// clear color of the screen
	clearColor mgl32.Vec3
}

func (f *CubeFormScreen) setupCubeFormScreen(wrapper interfaces.GLWrapper) {
	col := f.clearColor
	wrapper.ClearColor(col.X(), col.Y(), col.Z(), 1.0)
	wrapper.Enable(glwrapper.DEPTH_TEST)
	wrapper.DepthFunc(glwrapper.LESS)
	wrapper.Enable(glwrapper.BLEND)
	wrapper.BlendFunc(glwrapper.SRC_APLHA, glwrapper.ONE_MINUS_SRC_ALPHA)
	wrapper.Viewport(0, 0, int32(f.windowWidth), int32(f.windowHeight))
}

// Update
// if the main animation is running, we move from the start point to the finish point.
// I will define 'control points' for the movements. The direction will be based on the
// control points. The last control point is the destination.
// Lets define states for the screen, to make the animations easily distinguishable.
// The 'initial' state is the start state. If the current state is this, it changes it inmediately to
// the 'move-to-place' state, that moves the camera from control point to control point.
// If the last point is reached, it changes the state to 'form' where we can use the screens.
// The left rotation is triggered in the 'rotate-left-from-*' state and the right rotation is triggered
// in the 'rotate-right-from-*' state.
func (f *CubeFormScreen) Update(dt float64, p interfaces.Pointer, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
	f.handleState()
	f.cameraAnimation(dt)
	// movement handler. It will be unnecessary. Only for testing purposes in form.
	if f.state == "form" {
		f.cameraKeyboardMovement("forward", "back", "Walk", dt, keyStore)
		f.cameraKeyboardMovement("right", "left", "Strafe", dt, keyStore)
		f.cameraKeyboardMovement("up", "down", "Lift", dt, keyStore)
	}
	posX, posY := p.GetCurrent()
	TransformationMatrix := (f.camera.GetProjectionMatrix().Mul4(f.camera.GetViewMatrix())).Inv()
	currentMonitorPosition := mgl32.TransformCoordinate(f.middleMonitorPosition, mgl32.HomogRotate3DZ(mgl32.DegToRad(f.currentRotation)))
	coords := currentMonitorPosition.Add(mgl32.TransformCoordinate(mgl32.Vec3{float32(posX), float32(posY), 0.0}, TransformationMatrix))

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
	f.closestMesh = closestMesh
	f.closestDistance = closestDistance
	f.closestModel = closestModel
	// movement handler. It will be unnecessary. Only for testing purposes in form.
	if f.state == "form" {
		f.cameraKeyboardMovement("forward", "back", "Walk", dt, keyStore)
		f.cameraKeyboardMovement("right", "left", "Strafe", dt, keyStore)
		f.cameraKeyboardMovement("up", "down", "Lift", dt, keyStore)
		// check the buttons also.
		switch f.closestModel.(type) {
		case *ui.UIButton:
			btn := f.closestModel.(*ui.UIButton)
			btn.Hover()
			if buttonStore.Get(LEFT_MOUSE_BUTTON) {
				btn.OnState()
			} else {
				btn.OnState()
			}
			break
		}
		switch f.closestMesh.(type) {
		case *mesh.TexturedMaterialMesh:
			mMesh := f.closestMesh.(*mesh.TexturedMaterialMesh)
			if buttonStore.Get(LEFT_MOUSE_BUTTON) {
				if mMesh == f.MiddleMonitorGoRightButton {
					fmt.Println("MiddleMonitorGoRightButton button has been pressed.")
					f.state = "rotate-right-from-middle"
				} else if mMesh == f.MiddleMonitorGoLeftButton {
					fmt.Println("MiddleMonitorGoLeftButton button has been pressed.")
					f.state = "rotate-left-from-middle"
				} else if mMesh == f.RightMonitorGoLeftButton {
					fmt.Println("RightMonitorGoLeftButton button has been pressed.")
					f.state = "rotate-left-from-right"
				} else if mMesh == f.LeftMonitorGoRightButton {
					fmt.Println("LeftMonitorGoRightButton button has been pressed.")
					f.state = "rotate-right-from-left"
				}
			}
			break
		}
	}
}

// CharCallback
func (f *CubeFormScreen) CharCallback(char rune, wrapper interfaces.GLWrapper) {
}

// This function is responsible for the camera movement.
func (f *CubeFormScreen) cameraAnimation(dt float64) {
	dY := float32(0.0)
	switch f.state {
	case "move-to-place":
		// In this state, the camera moves from controlpoint to controlpoint.
		// From this initial position to the first control point, then from
		// the first one to the second one, etc.
		direction := f.controlPoints[f.controlPointIndex].Sub(f.controlPoints[f.controlPointIndex-1]).Normalize()
		currentPosition := f.GetCamera().GetPosition()
		newPosition := currentPosition.Add(direction.Mul(f.GetCamera().GetVelocity() * float32(dt)))
		currentDiffFromTheTarget := currentPosition.Sub(f.controlPoints[f.controlPointIndex]).Len()
		nextDiffFromTarget := newPosition.Sub(f.controlPoints[f.controlPointIndex]).Len()
		if nextDiffFromTarget < currentDiffFromTheTarget {
			// if the new position is closer to the control point, we move the camera to there.
			f.GetCamera().SetPosition(newPosition)
		} else {
			// move the camera to the exact position, and increase the index.
			f.GetCamera().SetPosition(f.controlPoints[f.controlPointIndex])
			f.controlPointIndex = f.controlPointIndex + 1
		}
		break
	case "form":
		// under the testing period, i want to make the movement possible.
		//f.GetCamera().SetPosition(f.controlPoints[len(f.controlPoints)-1])
		break
	case "rotate-left-from-right":
		dY = f.GetCamera().GetRotationStep() * float32(dt)
		f.SumOfRotation = f.SumOfRotation - dY
		if f.SumOfRotation < f.RotationToRightAngle {
			dY = dY - (f.RotationToRightAngle - f.SumOfRotation)
		}
		break
	case "rotate-left-from-middle":
		dY = f.GetCamera().GetRotationStep() * float32(dt)
		f.SumOfRotation = f.SumOfRotation + dY
		if f.SumOfRotation > f.RotationToLeftAngle {
			dY = dY - (f.SumOfRotation - f.RotationToLeftAngle)
		}
		break
	case "rotate-right-from-middle":
		dY = -f.GetCamera().GetRotationStep() * float32(dt)
		f.SumOfRotation = f.SumOfRotation + dY
		if f.SumOfRotation < f.RotationToRightAngle {
			dY = dY - (f.RotationToRightAngle - f.SumOfRotation)
		}
		break
	case "rotate-right-from-left":
		dY = -f.GetCamera().GetRotationStep() * float32(dt)
		f.SumOfRotation = f.SumOfRotation - dY
		if f.SumOfRotation > f.RotationToLeftAngle {
			dY = dY - (f.SumOfRotation - f.RotationToLeftAngle)
		}
		break
	}
	f.GetCamera().UpdateDirection(0.0, dY)
}

// State handler function. It updates the state, if the update conditions are reached.
func (f *CubeFormScreen) handleState() {
	switch f.state {
	case "initial":
		if len(f.controlPoints) > 1 {
			f.state = "move-to-place"
			f.controlPointIndex = 1
		} else {
			f.state = "form"
		}
		break
	case "move-to-place":
		if f.controlPointIndex == len(f.controlPoints) {
			f.state = "form"
		}
		break
	case "rotate-left-from-right":
		if f.SumOfRotation < f.RotationToRightAngle {
			f.state = "form"
			f.SumOfRotation = 0
			f.currentRotation = 0
		}
		break
	case "rotate-left-from-middle":
		if f.SumOfRotation > f.RotationToLeftAngle {
			f.state = "form"
			f.SumOfRotation = 0
			f.currentRotation = f.RotationToLeftAngle
		}
		break
	case "rotate-right-from-middle":
		if f.SumOfRotation < f.RotationToRightAngle {
			f.state = "form"
			f.SumOfRotation = 0
			f.currentRotation = f.RotationToRightAngle
		}
		break
	case "rotate-right-from-left":
		if f.SumOfRotation > f.RotationToLeftAngle {
			f.state = "form"
			f.SumOfRotation = 0
			f.currentRotation = 0
		}
		break
	case "form":
		break
	}
}
