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
