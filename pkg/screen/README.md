# Screen package

This package is written for handling different needs in the same application. For example we need a menu screen without camera and a world screen with a camera in the same application. This package provides solution for this need.
The models and shaders are handled in the screen package also the camera, light, movement maps, closest objects and the screen level uniforms.

It holds:
- a `camera`, that has to implement the `Camera` interface
- `shaderMap`, that makes the connection between the models and the shaders
- `directionalLightSources`, for storing the directional lights.
- `pointLightSources`, for storing the point lights.
- `spotLightSources`, for storing the spot lights.
- `cameraKeyboardMovementMap`, makes connection between the keyboard buttons and the camera state updates.
- `rotateOnEdgeDistance`, for the mouse rotations.
- `uniformFloat`, for storing the float uniforms that needs to be set for every shader.
- `uniformVector`, for storing the vector uniforms that needs to be set for every shader.
- `setupFunction`, for gl setup, like Enable stuff, setup clear color, etc.

## Functions

**New**

This function returns an initialized screen instance.

**Log**

Log returns the string representation of this object.

**SetCameraMovementMap**

SetCameraMovementMap sets the cameraKeyboardMovementMap variable. Currently the following values are supported: `forward`, `back`, `left`, `right`, `up`, `down`, `rotateLeft`, `rotateRight`, `rotateUp`, `rotateDown`

**SetRotateOnEdgeDistance**

SetRotateOnEdgeDistance updates the rotateOnEdgeDistance variable. The value has to be in the [0-1] interval. If not, a message is printed to the console and the variable update is skipped.

**SetCamera**

SetCamera updates the camera with the new one.

**GetCamera**

GetCamera returns the current camera of the screen.

**AddShader**

AddShader method inserts the new shader to the shaderMap

**AddModelToShader**

AddModelToShader attaches the model to a shader.

**GetClosestModelMeshDistance**

GetClosestModelMeshDistance returns the closest model, mesh and its distance from the mouse position.

**SetUniformFloat**

SetUniformFloat sets the given float value to the given string key in the uniformFloat map.

**SetUniformVector**

SetUniformVector sets the given mgl32.Vec3 value to the given string key in the uniformVector map.

**AddDirectionalLightSource**

AddDirectionalLightSource sets up a directional light source. It takes a DirectionalLight input that contains the model related info, and it also takes a [4]string, with the uniform names that are used in the shader applications the `DirectionUniformName`, `AmbientUniformName`, `DiffuseUniformName`, `SpecularUniformName`. They has to be in this order.

**AddPointLightSource**

AddPointLightSource sets up a point light source. It takes a PointLight input that contains the model related info, and it also containt the uniform names in [7]string format. The order has to be the following: `PositionUniformName`, `AmbientUniformName`, `DiffuseUniformName`, `SpecularUniformName`, `ConstantTermUniformName`, `LinearTermUniformName`, `QuadraticTermUniformName`.

**AddSpotLightSource**

AddSpotLightSource sets up a spot light source. It takes a SpotLight input that contains the model related info, and it also contains the uniform names in [10]string format. The order has to be the following: `PositionUniformName`, `DirectionUniformName`, `AmbientUniformName`, `DiffuseUniformName`, `SpecularUniformName`, `ConstantTermUniformName`, `LinearTermUniformName`, `QuadraticTermUniformName`, `CutoffUniformName`.

**Draw**

Draw calls Draw function in every drawable item. It calls the setupFunction, then it loops on the shaderMap (shaders). For each shader, first set it to used state, setup camera realted uniforms, then setup light related uniformsi and custom uniforms. Then we can pass the shader to the Model for drawing.

**Update**

Update loops on the shaderMap, and calls Update function on every Model. It also handles the camera movement and rotation, if the camera is set.

**Export**

Export creates a directory for the screen and calls Export function on the models.
