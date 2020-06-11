# Application package

The common application related stuff goes here. It holds

- a `window`, that has to implement the `Window` interface
- a `camera`, that has to implement the `Camera` interface
- mouse position `MousePosX`, `MousePosY`, set with the mouse button callback
- `shaderMap`, that makes the connection between she models and the shaders
- `mouseDowns`, that stores the mouse button states.
- `keyDowns`, that stores the keyboard button states.
- `directionalLightSources`, for storing the directional lights.
- `pointLightSources`, for storing the point lights.
- `spotLightSources`, for storing the spot lights.
- `cameraKeyboardMovementMap`, makes connection between the keyboard buttons and the camera state updates.
- `rotateOnEdgeDistance`, for the mouse rotations.

## New

This function returns an initialized application instance.

## Log

Log returns the string representation of this object.

## SetCameraMovementMap

SetCameraMovementMap sets the cameraKeyboardMovementMap variable. Currently the following values are supported: `forward`, `back`, `left`, `right`, `up`, `down`, `rotateLeft`, `rotateRight`, `rotateUp`, `rotateDown`

## SetRotateOnEdgeDistance

SetRotateOnEdgeDistance updates the rotateOnEdgeDistance variable. The value has to be in the [0-1] interval. If not, a message is printed to the console and the variable update is skipped.

## SetWindow

SetWindow updates the window with the new one.

## GetWindow

GetWindow returns the current window of the application.

## SetCamera

SetCamera updates the camera with the new one.

## GetCamera

GetCamera returns the current camera of the application.

## AddShader

AddShader method inserts the new shader to the shaderMap

## AddModelToShader

AddModelToShader attaches the model to a shader.

## Update

Update loops on the shaderMap, and calls Update function on every Model. It also handles the camera movement and rotation, if the camera is set.

## Draw

Draw calls Draw function in every drawable item. It loops on the shaderMap (shaders). For each shader, first set it to used state, setup camera realted uniforms, then setup light related uniforms. Then we can pass the shader to the Model for drawing.

## KeyCallback

KeyCallback is responsible for the keyboard event handling.

## MouseButtonCallback

MouseButtonCallback is responsible for the mouse button event handling.

## SetKeyState

SetKeyState setups the keyDowns based on the key and action

## SetKeyState

SetKeyState setups the keyDowns based on the key and action

## GetMouseButtonState

GetMouseButtonState returns the state of the given button

## GetKeyState

GetKeyState returns the state of the given key

## AddDirectionalLightSource

AddDirectionalLightSource sets up a directional light source. It takes a DirectionalLight input that contains the model related info, and it also takes a [4]string, with the uniform names that are used in the shader applications the `DirectionUniformName`, `AmbientUniformName`, `DiffuseUniformName`, `SpecularUniformName`. They has to be in this order.

## AddPointLightSource

AddPointLightSource sets up a point light source. It takes a PointLight input that contains the model related info, and it also containt the uniform names in [7]string format. The order has to be the following: `PositionUniformName`, `AmbientUniformName`, `DiffuseUniformName`, `SpecularUniformName`, `ConstantTermUniformName`, `LinearTermUniformName`, `QuadraticTermUniformName`.

## AddSpotLightSource

AddSpotLightSource sets up a spot light source. It takes a SpotLight input that contains the model related info, and it also contains the uniform names in [10]string format. The order has to be the following: `PositionUniformName`, `DirectionUniformName`, `AmbientUniformName`, `DiffuseUniformName`, `SpecularUniformName`, `ConstantTermUniformName`, `LinearTermUniformName`, `QuadraticTermUniformName`, `CutoffUniformName`.
