package screen

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"

	"github.com/akosgarai/playground_engine/pkg/interfaces"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type SetupFunction func(wrapper interfaces.GLWrapper)

type Screen struct {
	camera    interfaces.Camera
	cameraSet bool

	shaderMap map[interfaces.Shader][]interfaces.Model

	directionalLightSources []DirectionalLightSource
	pointLightSources       []PointLightSource
	spotLightSources        []SpotLightSource

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
	// Setup function is called right before drawing.
	setupFunction SetupFunction
}

// New returns a screen instance
func New() *Screen {
	return &Screen{
		cameraSet:                 false,
		shaderMap:                 make(map[interfaces.Shader][]interfaces.Model),
		directionalLightSources:   []DirectionalLightSource{},
		pointLightSources:         []PointLightSource{},
		spotLightSources:          []SpotLightSource{},
		cameraKeyboardMovementMap: make(map[string]glfw.Key),
		rotateOnEdgeDistance:      0.0,
		uniformFloat:              make(map[string]float32),
		uniformVector:             make(map[string]mgl32.Vec3),
		closestDistance:           math.MaxFloat32,
		setupFunction:             nil,
	}
}

// Log returns the string representation of this object.
func (s *Screen) Log() string {
	logString := "Screen:\n"
	if s.cameraSet {
		logString += " - camera : " + s.camera.Log() + "\n"
		logString += " - movement map:\n" + fmt.Sprintf("\t%#v\n", s.cameraKeyboardMovementMap)
		logString += " - rotate distance:\n" + fmt.Sprintf("\t%#v\n", s.rotateOnEdgeDistance)
	}
	return logString
}

// SetCameraMovementMap sets the cameraKeyboardMovementMap variable.
// Currently the following values are supported: 'forward', 'back',
// 'left', 'right', 'up', 'down', 'rotateLeft', 'rotateRight',
// 'rotateUp', 'rotateDown'
func (s *Screen) SetCameraMovementMap(m map[string]glfw.Key) {
	s.cameraKeyboardMovementMap = m
}

// SetRotateOnEdgeDistance updates the rotateOnEdgeDistance variable.
// The value has to be in the [0-1] interval. If not, a message is printed to the
// console and the variable update is skipped.
func (s *Screen) SetRotateOnEdgeDistance(value float32) {
	// validate value. [0-1]
	if value < 0 || value > 1 {
		fmt.Printf("Skipping rotateOnEdgeDistance variable update, value '%f' invalid.\n", value)
		return
	}
	s.rotateOnEdgeDistance = value
}

// SetCamera updates the camera with the new one.
func (s *Screen) SetCamera(c interfaces.Camera) {
	s.cameraSet = true
	s.camera = c
}

// GetCamera returns the current camera of the screen.
func (s *Screen) GetCamera() interfaces.Camera {
	return s.camera
}

// AddShader method inserts the new shader to the shaderMap
func (s *Screen) AddShader(sh interfaces.Shader) {
	s.shaderMap[sh] = []interfaces.Model{}
}

// AddModelToShader attaches the model to a shader.
func (s *Screen) AddModelToShader(m interfaces.Model, sh interfaces.Shader) {
	s.shaderMap[sh] = append(s.shaderMap[sh], m)
}

// GetClosestModelMeshDistance returns the closest model, mesh and its distance
// from the mouse position.
func (s *Screen) GetClosestModelMeshDistance() (interfaces.Model, interfaces.Mesh, float32) {
	return s.closestModel, s.closestMesh, s.closestDistance
}

// SetUniformFloat sets the given float value to the given string key in
// the uniformFloat map.
func (s *Screen) SetUniformFloat(key string, value float32) {
	s.uniformFloat[key] = value
}

// SetUniformVector sets the given mgl32.Vec3 value to the given string key in
// the uniformVector map.
func (s *Screen) SetUniformVector(key string, value mgl32.Vec3) {
	s.uniformVector[key] = value
}

// AddDirectionalLightSource sets up a directional light source.
// It takes a DirectionalLight input that contains the model related info,
// and it also takes a [4]string, with the uniform names that are used in the shader applications
// the 'DirectionUniformName', 'AmbientUniformName', 'DiffuseUniformName', 'SpecularUniformName'.
// They has to be in this order.
func (s *Screen) AddDirectionalLightSource(lightSource interfaces.DirectionalLight, uniformNames [4]string) {
	var dSource DirectionalLightSource
	dSource.LightSource = lightSource
	dSource.DirectionUniformName = uniformNames[0]
	dSource.AmbientUniformName = uniformNames[1]
	dSource.DiffuseUniformName = uniformNames[2]
	dSource.SpecularUniformName = uniformNames[3]

	s.directionalLightSources = append(s.directionalLightSources, dSource)
}

// AddPointLightSource sets up a point light source. It takes a PointLight
// input that contains the model related info, and it also containt the uniform names in [7]string format.
// The order has to be the following: 'PositionUniformName', 'AmbientUniformName', 'DiffuseUniformName',
// 'SpecularUniformName', 'ConstantTermUniformName', 'LinearTermUniformName', 'QuadraticTermUniformName'.
func (s *Screen) AddPointLightSource(lightSource interfaces.PointLight, uniformNames [7]string) {
	var pSource PointLightSource
	pSource.LightSource = lightSource
	pSource.PositionUniformName = uniformNames[0]
	pSource.AmbientUniformName = uniformNames[1]
	pSource.DiffuseUniformName = uniformNames[2]
	pSource.SpecularUniformName = uniformNames[3]
	pSource.ConstantTermUniformName = uniformNames[4]
	pSource.LinearTermUniformName = uniformNames[5]
	pSource.QuadraticTermUniformName = uniformNames[6]

	s.pointLightSources = append(s.pointLightSources, pSource)
}

// AddSpotLightSource sets up a spot light source. It takes a SpotLight input
// that contains the model related info, and it also contains the uniform names in [10]string format.
// The order has to be the following: 'PositionUniformName', 'DirectionUniformName', 'AmbientUniformName',
// 'DiffuseUniformName', 'SpecularUniformName', 'ConstantTermUniformName', 'LinearTermUniformName',
// 'QuadraticTermUniformName', 'CutoffUniformName'.
func (s *Screen) AddSpotLightSource(lightSource interfaces.SpotLight, uniformNames [10]string) {
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

	s.spotLightSources = append(s.spotLightSources, sSource)
}

// Setup light related uniforms.
func (s *Screen) lightHandler(sh interfaces.Shader) {
	s.setupDirectionalLightForShader(sh)
	s.setupPointLightForShader(sh)
	s.setupSpotLightForShader(sh)
}

// Setup directional light related uniforms. It iterates over the directional sources
// and setups each uniform, where the name is not empty.
func (s *Screen) setupDirectionalLightForShader(sh interfaces.Shader) {
	for _, source := range s.directionalLightSources {
		if source.DirectionUniformName != "" {
			direction := source.LightSource.GetDirection()
			sh.SetUniform3f(source.DirectionUniformName, direction.X(), direction.Y(), direction.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			sh.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			sh.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			sh.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
	}
	sh.SetUniform1i("NumberOfDirectionalLightSources", int32(len(s.directionalLightSources)))
}

// Setup point light relates uniforms. It iterates over the point light sources and sets
// up every uniform, where the name is not empty.
func (s *Screen) setupPointLightForShader(sh interfaces.Shader) {
	for _, source := range s.pointLightSources {
		if source.PositionUniformName != "" {
			position := source.LightSource.GetPosition()
			sh.SetUniform3f(source.PositionUniformName, position.X(), position.Y(), position.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			sh.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			sh.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			sh.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
		if source.ConstantTermUniformName != "" {
			sh.SetUniform1f(source.ConstantTermUniformName, source.LightSource.GetConstantTerm())
		}
		if source.LinearTermUniformName != "" {
			sh.SetUniform1f(source.LinearTermUniformName, source.LightSource.GetLinearTerm())
		}
		if source.QuadraticTermUniformName != "" {
			sh.SetUniform1f(source.QuadraticTermUniformName, source.LightSource.GetQuadraticTerm())
		}
	}
	sh.SetUniform1i("NumberOfPointLightSources", int32(len(s.pointLightSources)))
}

// Setup spot light related uniforms. It iterates over the spot light sources and sets up
// every uniform, where the name is not empty.
func (s *Screen) setupSpotLightForShader(sh interfaces.Shader) {
	for _, source := range s.spotLightSources {
		if source.DirectionUniformName != "" {
			direction := source.LightSource.GetDirection()
			sh.SetUniform3f(source.DirectionUniformName, direction.X(), direction.Y(), direction.Z())
		}
		if source.PositionUniformName != "" {
			position := source.LightSource.GetPosition()
			sh.SetUniform3f(source.PositionUniformName, position.X(), position.Y(), position.Z())
		}
		if source.AmbientUniformName != "" {
			ambient := source.LightSource.GetAmbient()
			sh.SetUniform3f(source.AmbientUniformName, ambient.X(), ambient.Y(), ambient.Z())
		}
		if source.DiffuseUniformName != "" {
			diffuse := source.LightSource.GetDiffuse()
			sh.SetUniform3f(source.DiffuseUniformName, diffuse.X(), diffuse.Y(), diffuse.Z())
		}
		if source.SpecularUniformName != "" {
			specular := source.LightSource.GetSpecular()
			sh.SetUniform3f(source.DiffuseUniformName, specular.X(), specular.Y(), specular.Z())
		}
		if source.ConstantTermUniformName != "" {
			sh.SetUniform1f(source.ConstantTermUniformName, source.LightSource.GetConstantTerm())
		}
		if source.LinearTermUniformName != "" {
			sh.SetUniform1f(source.LinearTermUniformName, source.LightSource.GetLinearTerm())
		}
		if source.QuadraticTermUniformName != "" {
			sh.SetUniform1f(source.QuadraticTermUniformName, source.LightSource.GetQuadraticTerm())
		}
		if source.CutoffUniformName != "" {
			sh.SetUniform1f(source.CutoffUniformName, source.LightSource.GetCutoff())
		}
		if source.OuterCutoffUniformName != "" {
			sh.SetUniform1f(source.OuterCutoffUniformName, source.LightSource.GetOuterCutoff())
		}
	}
	sh.SetUniform1i("NumberOfSpotLightSources", int32(len(s.spotLightSources)))
}

// Setup custom uniforms for the shader application.
func (s *Screen) customUniforms(sh interfaces.Shader) {
	for name, value := range s.uniformFloat {
		sh.SetUniform1f(name, value)
	}
	for name, value := range s.uniformVector {
		sh.SetUniform3f(name, value.X(), value.Y(), value.Z())
	}
}

// Draw calls Draw function in every drawable item. It calls the setupFunction then
// it loops on the shaderMap (shaders).
// For each shader, first set it to used state, setup camera realted uniforms,
// then setup light related uniformsi and custom uniforms. Then we can pass the shader to the Model for drawing.
func (s *Screen) Draw(wrapper interfaces.GLWrapper) {
	if s.setupFunction != nil {
		s.setupFunction(wrapper)
	}
	// Draw the non transparent models first
	for sh, _ := range s.shaderMap {
		sh.Use()
		if s.cameraSet {
			sh.SetUniformMat4("view", s.camera.GetViewMatrix())
			sh.SetUniformMat4("projection", s.camera.GetProjectionMatrix())
			cameraPos := s.camera.GetPosition()
			sh.SetUniform3f("viewPosition", cameraPos.X(), cameraPos.Y(), cameraPos.Z())
		} else {
			sh.SetUniformMat4("view", mgl32.Ident4())
			sh.SetUniformMat4("projection", mgl32.Ident4())
		}
		s.lightHandler(sh)
		// custom uniform setup.
		s.customUniforms(sh)
		for index, _ := range s.shaderMap[sh] {
			if !s.shaderMap[sh][index].IsTransparent() {
				s.shaderMap[sh][index].Draw(sh)
			}
		}
	}
	// Draw transparent models
	for sh, _ := range s.shaderMap {
		sh.Use()
		if s.cameraSet {
			sh.SetUniformMat4("view", s.camera.GetViewMatrix())
			sh.SetUniformMat4("projection", s.camera.GetProjectionMatrix())
			cameraPos := s.camera.GetPosition()
			sh.SetUniform3f("viewPosition", cameraPos.X(), cameraPos.Y(), cameraPos.Z())
		} else {
			sh.SetUniformMat4("view", mgl32.Ident4())
			sh.SetUniformMat4("projection", mgl32.Ident4())
		}
		s.lightHandler(sh)
		// custom uniform setup.
		s.customUniforms(sh)
		for index, _ := range s.shaderMap[sh] {
			if s.shaderMap[sh][index].IsTransparent() {
				s.shaderMap[sh][index].Draw(sh)
			}
		}
	}
}

// Update loops on the shaderMap, and calls Update function on every Model.
// It also handles the camera movement and rotation, if the camera is set.
func (s *Screen) Update(dt, posX, posY float64, store interfaces.RoKeyStore) {
	TransformationMatrix := mgl32.Ident4()
	if s.cameraSet {
		s.cameraKeyboardMovement("forward", "back", "Walk", dt, store)
		s.cameraKeyboardMovement("right", "left", "Strafe", dt, store)
		s.cameraKeyboardMovement("up", "down", "Lift", dt, store)
		s.cameraKeyboardRotation(dt, store)
		if s.rotateOnEdgeDistance > 0.0 {
			s.cameraMouseRotation(dt, posX, posY)
		}
		TransformationMatrix = (s.camera.GetProjectionMatrix().Mul4(s.camera.GetViewMatrix())).Inv()
	}
	coords := mgl32.TransformCoordinate(mgl32.Vec3{float32(posX), float32(posY), 0.0}, TransformationMatrix)
	closestDistance := float32(math.MaxFloat32)
	var closestMesh interfaces.Mesh
	var closestModel interfaces.Model

	for sh, _ := range s.shaderMap {
		for index, _ := range s.shaderMap[sh] {
			// The collision detection between the moving meshes supposed to be implemented somewhere here.
			// It could be the same as the collision detection between the camera and the other objects.
			// If a mesh is a moving object (in this case, the moving means movement on the direction),
			// then it needs to be tested against other objects. If collision is not found, the update
			// could be applied, otherwise it needs to be skipped. In the future, the collusion effect
			// also could be handled here.
			s.shaderMap[sh][index].Update(dt)
			msh, dist := s.shaderMap[sh][index].ClosestMeshTo(coords)
			if dist < closestDistance {
				closestDistance = dist
				closestMesh = msh
				closestModel = s.shaderMap[sh][index]
			}
		}
	}
	s.closestMesh = closestMesh
	s.closestModel = closestModel
	s.closestDistance = closestDistance
}

// cameraKeyboardMovement is responsible for handling a movement for a specific direction.
// The direction is described by the key strings. The handler function name is also added
// as input to be able to call it. For the movement we also need to know the delta time,
// that is also added as function input. In case of invalid function name,
// it prints out some message to the console.
func (s *Screen) cameraKeyboardMovement(directionKey, oppositeKey, handlerName string, delta float64, store interfaces.RoKeyStore) {
	keyStateDirection := false
	keyStateOpposite := false
	if val, ok := s.cameraKeyboardMovementMap[directionKey]; ok {
		keyStateDirection = store.Get(val)
	}
	if val, ok := s.cameraKeyboardMovementMap[oppositeKey]; ok {
		keyStateOpposite = store.Get(val)
	}
	step := float32(0.0)
	if keyStateDirection && !keyStateOpposite {
		step = float32(delta) * s.camera.GetVelocity()
	} else if keyStateOpposite && !keyStateDirection {
		step = -float32(delta) * s.camera.GetVelocity()
	}
	if step != 0 {
		// Collision detection. The function for the test is prefixed by 'BoundingObjectAfter'
		boundingObjectFunc := reflect.ValueOf(s.camera).MethodByName("BoundingObjectAfter" + handlerName)
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
			collide = s.cameraCollisionTest(testInput)
		}

		if !collide {
			method := reflect.ValueOf(s.camera).MethodByName(handlerName)
			if !method.IsValid() || method.IsZero() {
				fmt.Printf("Invalid method name '%s' was given for camera movement.\n", handlerName)
				return
			}
			method.Call(inputParams)
		}
	}
}

// cameraKeyboardRotation is responsible for handling the rotation events generated by the keyboard.
// The rotation(Up|Down|Left|Right) keys are checked from the maps
func (s *Screen) cameraKeyboardRotation(delta float64, store interfaces.RoKeyStore) {
	rotateUp := false
	rotateDown := false
	rotateLeft := false
	rotateRight := false
	if val, ok := s.cameraKeyboardMovementMap["rotateUp"]; ok {
		rotateUp = store.Get(val)
	}
	if val, ok := s.cameraKeyboardMovementMap["rotateDown"]; ok {
		rotateDown = store.Get(val)
	}
	if val, ok := s.cameraKeyboardMovementMap["rotateLeft"]; ok {
		rotateLeft = store.Get(val)
	}
	if val, ok := s.cameraKeyboardMovementMap["rotateRight"]; ok {
		rotateRight = store.Get(val)
	}
	s.applyMouseRotation(rotateLeft, rotateRight, rotateUp, rotateDown, delta)
}

// applyMouseRotation calls the camera's UpdateDirection function if necessary.
func (s *Screen) applyMouseRotation(rotateLeft, rotateRight, rotateUp, rotateDown bool, delta float64) {
	dX := float32(0.0)
	dY := float32(0.0)

	if rotateLeft && !rotateRight {
		dX = -s.camera.GetRotationStep() * float32(delta)
	} else if rotateRight && !rotateLeft {
		dX = s.camera.GetRotationStep() * float32(delta)
	}
	if rotateUp && !rotateDown {
		dY = -s.camera.GetRotationStep() * float32(delta)
	} else if rotateDown && !rotateUp {
		dY = s.camera.GetRotationStep() * float32(delta)
	}
	if dX != 0.0 || dY != 0.0 {
		s.camera.UpdateDirection(dX, dY)
	}
}

// cameraMouseRotation function is responsible for the rotation generated by the mouse
// position. If it is close to the edges, it triggers movement.
func (s *Screen) cameraMouseRotation(delta, posX, posY float64) {
	rotateUp := false
	rotateDown := false
	rotateLeft := false
	rotateRight := false
	x := float32(posX)
	y := float32(posY)
	if y > 1.0-s.rotateOnEdgeDistance && y < 1.0 {
		rotateUp = true
	}
	if y > -1.0 && y < -1.0+s.rotateOnEdgeDistance {
		rotateDown = true
	}
	if x > -1.0 && x < -1.0+s.rotateOnEdgeDistance {
		rotateLeft = true
	}
	if x > 1.0-s.rotateOnEdgeDistance && x < 1.0 {
		rotateRight = true
	}

	s.applyMouseRotation(rotateLeft, rotateRight, rotateUp, rotateDown, delta)
}

// cameraCollisionTest is responsible for the camera movement collision testing. It gets the bounding object for the next step.
// Under the hood, it iterates over the shaders, and tests collision for every mesh. It stops the test after the fist
// detected collision and returns true. Without detected collision it returns false.
func (s *Screen) cameraCollisionTest(boundingSphere *coldet.Sphere) bool {
	for sh, _ := range s.shaderMap {
		for index, _ := range s.shaderMap[sh] {
			if s.shaderMap[sh][index].CollideTestWithSphere(boundingSphere) {
				return true
			}
		}
	}
	return false
}

// Export creates a directory for the screen and calls Export function on the models.
func (s *Screen) Export(basePath string) {
	i := 0
	for sh, _ := range s.shaderMap {
		modelDir := strconv.Itoa(i)
		err := os.Mkdir(basePath+"/"+modelDir, os.ModeDir|os.ModePerm)
		if err != nil {
			fmt.Printf("Cannot create model directory. '%s'\n", err.Error())
		}
		for index, _ := range s.shaderMap[sh] {
			s.shaderMap[sh][index].Export(basePath + "/" + modelDir)
		}
		i++
	}
}

// Setup function sets the setupFunction to the given one
func (s *Screen) Setup(f SetupFunction) {
	s.setupFunction = f
}