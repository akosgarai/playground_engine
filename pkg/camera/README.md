# Camera

It represents the camera or eye. We see our model world from the camera's point of view. The implementation was based on the [learnopengl.com](https://learnopengl.com/Getting-started/Camera) tutorial.
Now it supports 2 kind of cameras, the `DefaultCamera` - the one that we already had - and the `FPSCamera` - the one that is similar to the camera that we can see in FPS games.

## NewCamera

Returns a new DefaultCamera with the given setup.

- `position` - the camera or eye position
- `worldUp` - the up direction in the world coordinate system
- `yaw` - the rotation in `Z` axis
- `pitch` - the rotation in `Y` axis

## NewFPSCamera

Returns a new FPSCamera with the given setup.

- `position` - the camera or eye position
- `worldUp` - the up direction in the world coordinate system
- `yaw` - the rotation in `Z` axis
- `pitch` - the rotation in `Y` axis, also the value of the frontDirPitch.

## Walk

It updates the position (forward, back directions).

## Strafe

It updates the position (left, right directions).

## Lift

It updates the position (up, down directions).

## SetupProjection

It sets the projection related variables.

- `fov` - field of view
- `aspectRatio` - windowWidth/windowHeight
- `near` - near clip plane
- `far` - far clip plane

## GetProjectionMatrix

It returns the projectionMatrix of the camera. It setups a perspective transformation.

## GetViewMatrix

It gets the matrix to transform from world coordinates to this cameras coordinates. It returns the viewMatrix of the camera.

## UpdateDirection

It updates the pitch and yaw values.

## GetBoundingObject

It returns the bounding object of the camera. Now it is defined as a sphere. The position is the current position of the camera. The radius is hardcoded to 0.1 (defaultCameraRadius constant).

## BoundingObjectAfterWalk

BoundingObjectAfterWalk returns the bounding object of the new position. It is used for collision detection. The step is forbidden if it leads to collision.

## BoundingObjectAfterStrafe

BoundingObjectAfterStrafe returns the bounding object of the new position. It is used for collision detection. The step is forbidden if it leads to collision.

## BoundingObjectAfterLift

BoundingObjectAfterLift returns the bounding object of the new position. It is used for collision detection. The step is forbidden if it leads to collision.
