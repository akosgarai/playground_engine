package application

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/transformations"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	DEBUG  = glfw.KeyH
	EXPORT = glfw.KeyP
)

type Camera interface {
	Log() string
	GetViewMatrix() mgl32.Mat4
	GetProjectionMatrix() mgl32.Mat4
	Walk(float32)
	Strafe(float32)
	Lift(float32)
	UpdateDirection(float32, float32)
	GetPosition() mgl32.Vec3
	GetVelocity() float32
	GetRotationStep() float32
	BoundingObjectAfterWalk(float32) *coldet.Sphere
	BoundingObjectAfterStrafe(float32) *coldet.Sphere
	BoundingObjectAfterLift(float32) *coldet.Sphere
}

type Window interface {
	GetCursorPos() (float64, float64)
	SetKeyCallback(glfw.KeyCallback) glfw.KeyCallback
	SetMouseButtonCallback(glfw.MouseButtonCallback) glfw.MouseButtonCallback
	ShouldClose() bool
	SwapBuffers()
	GetSize() (int, int)
	SetShouldClose(bool)
}

type Application struct {
	window    Window
	camera    Camera
	cameraSet bool

	shaderMap  map[interfaces.Shader][]interfaces.Model
	mouseDowns map[glfw.MouseButton]bool
	MousePosX  float64
	MousePosY  float64

	directionalLightSources []DirectionalLightSource
	pointLightSources       []PointLightSource
	spotLightSources        []SpotLightSource

	keyDowns map[glfw.Key]bool

	// it holds the keyMaps for the camera movement
	cameraKeyboardMovementMap map[string]glfw.Key
	// rotateOnEdgeDistance stores the variable
	// that is checked for rotating the camera
	// if the mouse is close to the window edge
	rotateOnEdgeDistance float32

	// uniforms, that needs to be set for every shader.
	uniformFloat  map[string]float32    // map for float32
	uniformVector map[string]mgl32.Vec3 // map for 3 float32

	// closestMesh the mesh that is closest to the mouse position
	// closestModel the model of the closestMesh
	// closestDistance is the distance of the mouse position from the closestMesh
	closestMesh     interfaces.Mesh
	closestModel    interfaces.Model
	closestDistance float32
}

// New returns an application instance
func New() *Application {
	return &Application{
		cameraSet:                 false,
		shaderMap:                 make(map[interfaces.Shader][]interfaces.Model),
		mouseDowns:                make(map[glfw.MouseButton]bool),
		directionalLightSources:   []DirectionalLightSource{},
		pointLightSources:         []PointLightSource{},
		spotLightSources:          []SpotLightSource{},
		keyDowns:                  make(map[glfw.Key]bool),
		cameraKeyboardMovementMap: make(map[string]glfw.Key),
		rotateOnEdgeDistance:      0.0,
		uniformFloat:              make(map[string]float32),
		uniformVector:             make(map[string]mgl32.Vec3),
		closestDistance:           math.MaxFloat32,
	}
}

// Log returns the string representation of this object.
func (a *Application) Log() string {
	logString := "Application:\n"
	if a.cameraSet {
		logString += " - camera : " + a.camera.Log() + "\n"
	}
	return logString
}

// SetCameraMovementMap sets the cameraKeyboardMovementMap variable.
// Currently the following values are supported: 'forward', 'back',
// 'left', 'right', 'up', 'down', 'rotateLeft', 'rotateRight',
// 'rotateUp', 'rotateDown'
func (a *Application) SetCameraMovementMap(m map[string]glfw.Key) {
	a.cameraKeyboardMovementMap = m
}

// SetRotateOnEdgeDistance updates the rotateOnEdgeDistance variable.
// The value has to be in the [0-1] interval. If not, a message is printed to the
// console and the variable update is skipped.
func (a *Application) SetRotateOnEdgeDistance(value float32) {
	// validate value. [0-1]
	if value < 0 || value > 1 {
		fmt.Printf("Skipping rotateOnEdgeDistance variable update, value '%f' invalid.\n", value)
		return
	}
	a.rotateOnEdgeDistance = value
}

// SetWindow updates the window with the new one.
func (a *Application) SetWindow(w Window) {
	a.window = w
}

// GetWindow returns the current window of the application.
func (a *Application) GetWindow() Window {
	return a.window
}

// SetCamera updates the camera with the new one.
func (a *Application) SetCamera(c Camera) {
	a.cameraSet = true
	a.camera = c
}

// GetCamera returns the current camera of the application.
func (a *Application) GetCamera() Camera {
	return a.camera
}

// AddShader method inserts the new shader to the shaderMap
func (a *Application) AddShader(s interfaces.Shader) {
	a.shaderMap[s] = []interfaces.Model{}
}

// AddModelToShader attaches the model to a shader.
func (a *Application) AddModelToShader(m interfaces.Model, s interfaces.Shader) {
	a.shaderMap[s] = append(a.shaderMap[s], m)
}

// GetClosestModelMeshDistance returns the closest model, mesh and its distance
// from the mouse position.
func (a *Application) GetClosestModelMeshDistance() (interfaces.Model, interfaces.Mesh, float32) {
	return a.closestModel, a.closestMesh, a.closestDistance
}

// cameraKeyboardMovement is responsible for handling a movement for a specific direction.
// The direction is described by the key strings. The handler function name is also added
// as input to be able to call it. For the movement we also need to know the delta time,
// that is also added as function input. In case of invalid function name,
// it prints out some message to the console.
func (a *Application) cameraKeyboardMovement(directionKey, oppositeKey, handlerName string, delta float64) {
	keyStateDirection := false
	keyStateOpposite := false
	if val, ok := a.cameraKeyboardMovementMap[directionKey]; ok {
		keyStateDirection = a.GetKeyState(val)
	}
	if val, ok := a.cameraKeyboardMovementMap[oppositeKey]; ok {
		keyStateOpposite = a.GetKeyState(val)
	}
	step := float32(0.0)
	if keyStateDirection && !keyStateOpposite {
		step = float32(delta) * a.camera.GetVelocity()
	} else if keyStateOpposite && !keyStateDirection {
		step = -float32(delta) * a.camera.GetVelocity()
	}
	if step != 0 {
		// Collision detection. The function for the test is prefixed by 'BoundingObjectAfter'
		boundingObjectFunc := reflect.ValueOf(a.camera).MethodByName("BoundingObjectAfter" + handlerName)
		var inputParams []reflect.Value
		inputParams = append(inputParams, reflect.ValueOf(step))
		collide := false

		if !boundingObjectFunc.IsValid() || boundingObjectFunc.IsZero() {
			fmt.Printf("Invalid method name '%s' was given for collisison detection. Skipping it.\n", "BoundingObjectAfter"+handlerName)
		} else {
			// Call the function for getting the bounding object after the step.
			functionResult := boundingObjectFunc.Call(inputParams)
			// Instead of []reflect.Value, i need coldet.Sphere, so it needs to be casted
			// to its type.
			testInput := functionResult[0].Interface().(*coldet.Sphere)
			collide = a.cameraCollisionTest(testInput)
		}

		if !collide {
			method := reflect.ValueOf(a.camera).MethodByName(handlerName)
			if !method.IsValid() || method.IsZero() {
				fmt.Printf("Invalid method name '%s' was given for camera movement.\n", handlerName)
				return
			}
			method.Call(inputParams)
		}
	}
}

// cameraCollisionTest is responsible for the camera movement collision testing. It gets the bounding object for the next step.
// Under the hood, it iterates over the shaders, and tests collision for every mesh. It stops the test after the fist
// detected collision and returns true. Without detected collision it returns false.
func (a *Application) cameraCollisionTest(boundingSphere *coldet.Sphere) bool {
	for s, _ := range a.shaderMap {
		for index, _ := range a.shaderMap[s] {
			if a.shaderMap[s][index].CollideTestWithSphere(boundingSphere) {
				return true
			}
		}
	}
	return false
}

// cameraKeyboardRotation is responsible for handling the rotation events generated by the keyboard.
// The rotation(Up|Down|Left|Right) keys are checked from the maps
func (a *Application) cameraKeyboardRotation(delta float64) {
	rotateUp := false
	rotateDown := false
	rotateLeft := false
	rotateRight := false
	if val, ok := a.cameraKeyboardMovementMap["rotateUp"]; ok {
		rotateUp = a.GetKeyState(val)
	}
	if val, ok := a.cameraKeyboardMovementMap["rotateDown"]; ok {
		rotateDown = a.GetKeyState(val)
	}
	if val, ok := a.cameraKeyboardMovementMap["rotateLeft"]; ok {
		rotateLeft = a.GetKeyState(val)
	}
	if val, ok := a.cameraKeyboardMovementMap["rotateRight"]; ok {
		rotateRight = a.GetKeyState(val)
	}
	a.applyMouseRotation(rotateLeft, rotateRight, rotateUp, rotateDown, delta)
}

// applyMouseRotation calls the camera's UpdateDirection function if necessary.
func (a *Application) applyMouseRotation(rotateLeft, rotateRight, rotateUp, rotateDown bool, delta float64) {
	dX := float32(0.0)
	dY := float32(0.0)

	if rotateLeft && !rotateRight {
		dX = -a.camera.GetRotationStep() * float32(delta)
	} else if rotateRight && !rotateLeft {
		dX = a.camera.GetRotationStep() * float32(delta)
	}
	if rotateUp && !rotateDown {
		dY = -a.camera.GetRotationStep() * float32(delta)
	} else if rotateDown && !rotateUp {
		dY = a.camera.GetRotationStep() * float32(delta)
	}
	if dX != 0.0 || dY != 0.0 {
		a.camera.UpdateDirection(dX, dY)
	}
}

// cameraMouseRotation function is responsible for the rotation generated by the mouse
// position. If it is close to the edges, it triggers movement.
func (a *Application) cameraMouseRotation(delta float64) {
	currX, currY := a.window.GetCursorPos()
	windowWidth, windowHeight := a.window.GetSize()
	tX, tY := transformations.MouseCoordinates(currX, currY, float64(windowWidth), float64(windowHeight))
	rotateUp := false
	rotateDown := false
	rotateLeft := false
	rotateRight := false
	x := float32(tX)
	y := float32(tY)
	if y > 1.0-a.rotateOnEdgeDistance && y < 1.0 {
		rotateUp = true
	}
	if y > -1.0 && y < -1.0+a.rotateOnEdgeDistance {
		rotateDown = true
	}
	if x > -1.0 && x < -1.0+a.rotateOnEdgeDistance {
		rotateLeft = true
	}
	if x > 1.0-a.rotateOnEdgeDistance && x < 1.0 {
		rotateRight = true
	}

	a.applyMouseRotation(rotateLeft, rotateRight, rotateUp, rotateDown, delta)
}

// Update loops on the shaderMap, and calls Update function on every Model.
// It also handles the camera movement and rotation, if the camera is set.
func (a *Application) Update(dt float64) {
	MousePosX, MousePosY := a.window.GetCursorPos()
	WindowWidth, WindowHeight := a.window.GetSize()
	mX, mY := transformations.MouseCoordinates(MousePosX, MousePosY, float64(WindowWidth), float64(WindowHeight))
	TransformationMatrix := mgl32.Ident4()
	if a.cameraSet {
		a.cameraKeyboardMovement("forward", "back", "Walk", dt)
		a.cameraKeyboardMovement("right", "left", "Strafe", dt)
		a.cameraKeyboardMovement("up", "down", "Lift", dt)
		a.cameraKeyboardRotation(dt)
		if a.rotateOnEdgeDistance > 0.0 {
			a.cameraMouseRotation(dt)
		}
		TransformationMatrix = (a.camera.GetProjectionMatrix().Mul4(a.camera.GetViewMatrix())).Inv()
	}
	coords := mgl32.TransformCoordinate(mgl32.Vec3{float32(mX), float32(mY), 0.0}, TransformationMatrix)
	closestDistance := float32(math.MaxFloat32)
	var closestMesh interfaces.Mesh
	var closestModel interfaces.Model

	for s, _ := range a.shaderMap {
		for index, _ := range a.shaderMap[s] {
			// The collision detection between the moving meshes supposed to be implemented somewhere here.
			// It could be the same as the collision detection between the camera and the other objects.
			// If a mesh is a moving object (in this case, the moving means movement on the direction),
			// then it needs to be tested against other objects. If collision is not found, the update
			// could be applied, otherwise it needs to be skipped. In the future, the collusion effect
			// also could be handled here.
			a.shaderMap[s][index].Update(dt)
			msh, dist := a.shaderMap[s][index].ClosestMeshTo(coords)
			if dist < closestDistance {
				closestDistance = dist
				closestMesh = msh
				closestModel = a.shaderMap[s][index]
			}
		}
	}
	a.closestMesh = closestMesh
	a.closestModel = closestModel
	a.closestDistance = closestDistance
}

// Draw calls Draw function in every drawable item. It loops on the shaderMap (shaders).
// For each shader, first set it to used state, setup camera realted uniforms,
// then setup light related uniformsi and custom uniforms. Then we can pass the shader to the Model for drawing.
func (a *Application) Draw() {
	// Draw the non transparent models first
	for s, _ := range a.shaderMap {
		s.Use()
		if a.cameraSet {
			s.SetUniformMat4("view", a.camera.GetViewMatrix())
			s.SetUniformMat4("projection", a.camera.GetProjectionMatrix())
			cameraPos := a.camera.GetPosition()
			s.SetUniform3f("viewPosition", cameraPos.X(), cameraPos.Y(), cameraPos.Z())
		} else {
			s.SetUniformMat4("view", mgl32.Ident4())
			s.SetUniformMat4("projection", mgl32.Ident4())
		}
		a.lightHandler(s)
		// custom uniform setup.
		a.customUniforms(s)
		for index, _ := range a.shaderMap[s] {
			if !a.shaderMap[s][index].IsTransparent() {
				a.shaderMap[s][index].Draw(s)
			}
		}
	}
	// Draw transparent models
	for s, _ := range a.shaderMap {
		s.Use()
		if a.cameraSet {
			s.SetUniformMat4("view", a.camera.GetViewMatrix())
			s.SetUniformMat4("projection", a.camera.GetProjectionMatrix())
			cameraPos := a.camera.GetPosition()
			s.SetUniform3f("viewPosition", cameraPos.X(), cameraPos.Y(), cameraPos.Z())
		} else {
			s.SetUniformMat4("view", mgl32.Ident4())
			s.SetUniformMat4("projection", mgl32.Ident4())
		}
		a.lightHandler(s)
		// custom uniform setup.
		a.customUniforms(s)
		for index, _ := range a.shaderMap[s] {
			if a.shaderMap[s][index].IsTransparent() {
				a.shaderMap[s][index].Draw(s)
			}
		}
	}
}

// KeyCallback is responsible for the keyboard event handling.
func (a *Application) KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch key {
	case DEBUG:
		if action != glfw.Release {
			fmt.Printf("%s\n", a.Log())
		}
		break
	case EXPORT:
		if action != glfw.Release {
			a.export()
		}
		break
	default:
		a.SetKeyState(key, action)
		break
	}
}

// MouseButtonCallback is responsible for the mouse button event handling.
func (a *Application) MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	a.MousePosX, a.MousePosY = w.GetCursorPos()
	switch button {
	default:
		a.SetButtonState(button, action)
		break
	}
}

// SetKeyState setups the keyDowns based on the key and action
func (a *Application) SetKeyState(key glfw.Key, action glfw.Action) {
	var isButtonPressed bool
	if action != glfw.Release {
		isButtonPressed = true
	} else {
		isButtonPressed = false
	}
	a.keyDowns[key] = isButtonPressed
}

// SetKeyState setups the keyDowns based on the key and action
func (a *Application) SetButtonState(button glfw.MouseButton, action glfw.Action) {
	var isButtonPressed bool
	if action != glfw.Release {
		isButtonPressed = true
	} else {
		isButtonPressed = false
	}
	a.mouseDowns[button] = isButtonPressed
}

// GetMouseButtonState returns the state of the given button
func (a *Application) GetMouseButtonState(button glfw.MouseButton) bool {
	return a.mouseDowns[button]
}

// GetKeyState returns the state of the given key
func (a *Application) GetKeyState(key glfw.Key) bool {
	return a.keyDowns[key]
}

// SetUniformFloat sets the given float value to the given string key in
// the uniformFloat map.
func (a *Application) SetUniformFloat(key string, value float32) {
	a.uniformFloat[key] = value
}

// SetUniformVector sets the given mgl32.Vec3 value to the given string key in
// the uniformVector map.
func (a *Application) SetUniformVector(key string, value mgl32.Vec3) {
	a.uniformVector[key] = value
}

// Setup custom uniforms for the shader application.
func (a *Application) customUniforms(s interfaces.Shader) {
	for name, value := range a.uniformFloat {
		s.SetUniform1f(name, value)
	}
	for name, value := range a.uniformVector {
		s.SetUniform3f(name, value.X(), value.Y(), value.Z())
	}
}

// Setup light related uniforms.
func (a *Application) lightHandler(s interfaces.Shader) {
	a.setupDirectionalLightForShader(s)
	a.setupPointLightForShader(s)
	a.setupSpotLightForShader(s)
}

// Setup directional light related uniforms. It iterates over the directional sources
// and setups each uniform, where the name is not empty.
func (a *Application) setupDirectionalLightForShader(s interfaces.Shader) {
	for _, source := range a.directionalLightSources {
		if source.DirectionUniformName != "" {
			direction := source.LightSource.GetDirection()
			s.SetUniform3f(source.DirectionUniformName, direction.X(), direction.Y(), direction.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			s.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			s.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			s.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
	}
	s.SetUniform1i("NumberOfDirectionalLightSources", int32(len(a.directionalLightSources)))

}

// Setup point light relates uniforms. It iterates over the point light sources and sets
// up every uniform, where the name is not empty.
func (a *Application) setupPointLightForShader(s interfaces.Shader) {
	for _, source := range a.pointLightSources {
		if source.PositionUniformName != "" {
			position := source.LightSource.GetPosition()
			s.SetUniform3f(source.PositionUniformName, position.X(), position.Y(), position.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			s.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			s.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			s.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
		if source.ConstantTermUniformName != "" {
			s.SetUniform1f(source.ConstantTermUniformName, source.LightSource.GetConstantTerm())
		}
		if source.LinearTermUniformName != "" {
			s.SetUniform1f(source.LinearTermUniformName, source.LightSource.GetLinearTerm())
		}
		if source.QuadraticTermUniformName != "" {
			s.SetUniform1f(source.QuadraticTermUniformName, source.LightSource.GetQuadraticTerm())
		}
	}
	s.SetUniform1i("NumberOfPointLightSources", int32(len(a.pointLightSources)))
}

// Setup spot light related uniforms. It iterates over the spot light sources and sets up
// every uniform, where the name is not empty.
func (a *Application) setupSpotLightForShader(s interfaces.Shader) {
	for _, source := range a.spotLightSources {
		if source.DirectionUniformName != "" {
			direction := source.LightSource.GetDirection()
			s.SetUniform3f(source.DirectionUniformName, direction.X(), direction.Y(), direction.Z())
		}
		if source.PositionUniformName != "" {
			position := source.LightSource.GetPosition()
			s.SetUniform3f(source.PositionUniformName, position.X(), position.Y(), position.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			s.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			s.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			s.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
		if source.ConstantTermUniformName != "" {
			s.SetUniform1f(source.ConstantTermUniformName, source.LightSource.GetConstantTerm())
		}
		if source.LinearTermUniformName != "" {
			s.SetUniform1f(source.LinearTermUniformName, source.LightSource.GetLinearTerm())
		}
		if source.QuadraticTermUniformName != "" {
			s.SetUniform1f(source.QuadraticTermUniformName, source.LightSource.GetQuadraticTerm())
		}
		if source.CutoffUniformName != "" {
			s.SetUniform1f(source.CutoffUniformName, source.LightSource.GetCutoff())
		}
		if source.OuterCutoffUniformName != "" {
			s.SetUniform1f(source.OuterCutoffUniformName, source.LightSource.GetOuterCutoff())
		}
	}
	s.SetUniform1i("NumberOfSpotLightSources", int32(len(a.spotLightSources)))
}

// AddDirectionalLightSource sets up a directional light source.
// It takes a DirectionalLight input that contains the model related info,
// and it also takes a [4]string, with the uniform names that are used in the shader applications
// the 'DirectionUniformName', 'AmbientUniformName', 'DiffuseUniformName', 'SpecularUniformName'.
// They has to be in this order.
func (a *Application) AddDirectionalLightSource(lightSource interfaces.DirectionalLight, uniformNames [4]string) {
	var dSource DirectionalLightSource
	dSource.LightSource = lightSource
	dSource.DirectionUniformName = uniformNames[0]
	dSource.AmbientUniformName = uniformNames[1]
	dSource.DiffuseUniformName = uniformNames[2]
	dSource.SpecularUniformName = uniformNames[3]

	a.directionalLightSources = append(a.directionalLightSources, dSource)
}

// AddPointLightSource sets up a point light source. It takes a PointLight
// input that contains the model related info, and it also containt the uniform names in [7]string format.
// The order has to be the following: 'PositionUniformName', 'AmbientUniformName', 'DiffuseUniformName',
// 'SpecularUniformName', 'ConstantTermUniformName', 'LinearTermUniformName', 'QuadraticTermUniformName'.
func (a *Application) AddPointLightSource(lightSource interfaces.PointLight, uniformNames [7]string) {
	var pSource PointLightSource
	pSource.LightSource = lightSource
	pSource.PositionUniformName = uniformNames[0]
	pSource.AmbientUniformName = uniformNames[1]
	pSource.DiffuseUniformName = uniformNames[2]
	pSource.SpecularUniformName = uniformNames[3]
	pSource.ConstantTermUniformName = uniformNames[4]
	pSource.LinearTermUniformName = uniformNames[5]
	pSource.QuadraticTermUniformName = uniformNames[6]

	a.pointLightSources = append(a.pointLightSources, pSource)
}

// AddSpotLightSource sets up a spot light source. It takes a SpotLight input
// that contains the model related info, and it also contains the uniform names in [10]string format.
// The order has to be the following: 'PositionUniformName', 'DirectionUniformName', 'AmbientUniformName',
// 'DiffuseUniformName', 'SpecularUniformName', 'ConstantTermUniformName', 'LinearTermUniformName',
// 'QuadraticTermUniformName', 'CutoffUniformName'.
func (a *Application) AddSpotLightSource(lightSource interfaces.SpotLight, uniformNames [10]string) {
	var sSource SpotLightSource
	sSource.LightSource = lightSource
	sSource.PositionUniformName = uniformNames[0]
	sSource.DirectionUniformName = uniformNames[1]
	sSource.AmbientUniformName = uniformNames[2]
	sSource.DiffuseUniformName = uniformNames[3]
	sSource.SpecularUniformName = uniformNames[4]
	sSource.ConstantTermUniformName = uniformNames[5]
	sSource.LinearTermUniformName = uniformNames[6]
	sSource.QuadraticTermUniformName = uniformNames[7]
	sSource.CutoffUniformName = uniformNames[8]
	sSource.OuterCutoffUniformName = uniformNames[9]

	a.spotLightSources = append(a.spotLightSources, sSource)
}

// This function is called for starting the export process. It is attached to a key callback.
func (a *Application) export() {
	ExportBaseDir := "./exports"
	Directory := time.Now().Format("20060102150405")
	err := os.Mkdir(ExportBaseDir+"/"+Directory, os.ModeDir|os.ModePerm)
	if err != nil {
		fmt.Printf("Cannot create export directory. '%s'\n", err.Error())
	}
	i := 0
	for s, _ := range a.shaderMap {
		modelDir := strconv.Itoa(i)
		err := os.Mkdir(ExportBaseDir+"/"+Directory+"/"+modelDir, os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Printf("Cannot create model directory. '%s'\n", err.Error())
		}
		for index, _ := range a.shaderMap[s] {
			a.shaderMap[s][index].Export(ExportBaseDir + "/" + Directory + "/" + modelDir)
		}
		i++
	}
}
