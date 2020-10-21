package camera

import (
	"math"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/akosgarai/playground_engine/pkg/transformations"
)

const (
	defaultCameraRadius = float32(0.1)
)

type Camera struct {
	// Camera options

	// Euler Angles
	pitch float32
	yaw   float32

	// Camera attributes
	cameraPosition       mgl32.Vec3
	cameraFrontDirection mgl32.Vec3
	cameraUpDirection    mgl32.Vec3
	cameraRightDirection mgl32.Vec3
	worldUp              mgl32.Vec3
	// Projection options.
	projectionOptions struct {
		fov         float32
		aspectRatio float32

		far  float32
		near float32
	}
	// velocity is used for the camera movement
	velocity float32
	// it is used for the camera rotation. rotation step deg.
	// has to be rotated in 1 sec.
	rotationStep float32
}

type DefaultCamera struct {
	Camera
}
type FPSCamera struct {
	Camera
	frontDirPitch float32
}

// Log returns the string representation of this object.
func (c *Camera) Log() string {
	logString := "cameraPosition: Vector{" + transformations.Vec3ToString(c.cameraPosition) + "}\n"
	logString += "worldUp: Vector{" + transformations.Vec3ToString(c.worldUp) + "}\n"
	logString += "cameraFrontDirection: Vector{" + transformations.Vec3ToString(c.cameraFrontDirection) + "}\n"
	logString += "cameraUpDirection: Vector{" + transformations.Vec3ToString(c.cameraUpDirection) + "}\n"
	logString += "cameraRightDirection: Vector{" + transformations.Vec3ToString(c.cameraRightDirection) + "}\n"
	logString += "yaw : " + transformations.Float32ToString(c.yaw) + "\n"
	logString += "pitch : " + transformations.Float32ToString(c.pitch) + "\n"
	logString += "velocity : " + transformations.Float32ToString(c.velocity) + "\n"
	logString += "rotationStep : " + transformations.Float32ToString(c.rotationStep) + "\n"
	logString += "ProjectionOptions:\n"
	logString += " - fov : " + transformations.Float32ToString(c.projectionOptions.fov) + "\n"
	logString += " - aspectRatio : " + transformations.Float32ToString(c.projectionOptions.aspectRatio) + "\n"
	logString += " - far : " + transformations.Float32ToString(c.projectionOptions.far) + "\n"
	logString += " - near : " + transformations.Float32ToString(c.projectionOptions.near) + "\n"
	return logString
}

// Returns a new camera with the given setup
// position - the camera or eye position
// worldUp - the up direction in the world coordinate system
// yaw - the rotation in z
// pitch - the rotation in y
func NewCamera(position, worldUp mgl32.Vec3, yaw, pitch float32) *DefaultCamera {
	cam := Camera{
		pitch:             pitch,
		yaw:               yaw,
		cameraPosition:    position,
		cameraUpDirection: mgl32.Vec3{0, 1, 0},
		worldUp:           worldUp,
		velocity:          0,
		rotationStep:      0,
	}

	cam.updateVectors()
	return &DefaultCamera{
		cam,
	}
}
func NewFPSCamera(position, worldUp mgl32.Vec3, yaw, pitch float32) *FPSCamera {
	cam := Camera{
		pitch:             pitch,
		yaw:               yaw,
		cameraPosition:    position,
		cameraUpDirection: mgl32.Vec3{0, 1, 0},
		worldUp:           worldUp,
		velocity:          0,
		rotationStep:      0,
	}

	cam.updateVectors()
	return &FPSCamera{
		cam,
		pitch,
	}
}

// BoundingObjectAfterWalk returns the bounding object of the new position.
func (c *FPSCamera) BoundingObjectAfterWalk(amount float32) *coldet.Sphere {
	// Front direction in the world system
	radPitch := float64(mgl32.DegToRad(c.frontDirPitch))
	radYaw := float64(mgl32.DegToRad(c.yaw))
	worldFrontDir := mgl32.Vec3{
		float32(math.Cos(radPitch) * math.Cos(radYaw)),
		float32(math.Sin(radPitch)),
		float32(math.Cos(radPitch) * math.Sin(radYaw)),
	}.Normalize()
	np := c.cameraPosition.Add(worldFrontDir.Mul(amount))
	return coldet.NewBoundingSphere([3]float32{np.X(), np.Y(), np.Z()}, 0.1)
}

// BoundingObjectAfterStrafe returns the bounding object of the new position.
func (c *FPSCamera) BoundingObjectAfterStrafe(amount float32) *coldet.Sphere {
	np := c.cameraPosition.Add(c.cameraRightDirection.Mul(amount))
	return coldet.NewBoundingSphere([3]float32{np.X(), np.Y(), np.Z()}, 0.1)
}

// BoundingObjectAfterLift returns the bounding object of the new position.
func (c *FPSCamera) BoundingObjectAfterLift(amount float32) *coldet.Sphere {
	np := c.cameraPosition.Add(c.cameraUpDirection.Mul(amount))
	return coldet.NewBoundingSphere([3]float32{np.X(), np.Y(), np.Z()}, 0.1)
}

// Walk updates the position (forward, back directions)
func (c *FPSCamera) Walk(amount float32) {
	// Front direction in the world system
	radPitch := float64(mgl32.DegToRad(c.frontDirPitch))
	radYaw := float64(mgl32.DegToRad(c.yaw))
	worldFrontDir := mgl32.Vec3{
		float32(math.Cos(radPitch) * math.Cos(radYaw)),
		float32(math.Sin(radPitch)),
		float32(math.Cos(radPitch) * math.Sin(radYaw)),
	}.Normalize()
	c.cameraPosition = c.cameraPosition.Add(worldFrontDir.Mul(amount))
}

// Strafe updates the position (left, right directions)
func (c *FPSCamera) Strafe(amount float32) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraRightDirection.Mul(amount))
}

// Lift updates the position (up, down directions)
func (c *FPSCamera) Lift(amount float32) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraUpDirection.Mul(amount))
}

// BoundingObjectAfterWalk returns the bounding object of the new position.
func (c *DefaultCamera) BoundingObjectAfterWalk(amount float32) *coldet.Sphere {
	np := c.cameraPosition.Add(c.cameraFrontDirection.Mul(amount))
	return coldet.NewBoundingSphere([3]float32{np.X(), np.Y(), np.Z()}, 0.1)
}

// BoundingObjectAfterStrafe returns the bounding object of the new position.
func (c *DefaultCamera) BoundingObjectAfterStrafe(amount float32) *coldet.Sphere {
	np := c.cameraPosition.Add(c.cameraRightDirection.Mul(amount))
	return coldet.NewBoundingSphere([3]float32{np.X(), np.Y(), np.Z()}, 0.1)
}

// BoundingObjectAfterLift returns the bounding object of the new position.
func (c *DefaultCamera) BoundingObjectAfterLift(amount float32) *coldet.Sphere {
	np := c.cameraPosition.Add(c.cameraUpDirection.Mul(amount))
	return coldet.NewBoundingSphere([3]float32{np.X(), np.Y(), np.Z()}, 0.1)
}

// Walk updates the position (forward, back directions)
func (c *DefaultCamera) Walk(amount float32) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraFrontDirection.Mul(amount))
}

// Strafe updates the position (left, right directions)
func (c *DefaultCamera) Strafe(amount float32) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraRightDirection.Mul(amount))
}

// Lift updates the position (up, down directions)
func (c *DefaultCamera) Lift(amount float32) {
	c.cameraPosition = c.cameraPosition.Add(c.cameraUpDirection.Mul(amount))
}

// SetupProjection sets the projection related variables
// fov - field of view
// aspectRatio - windowWidth/windowHeight
// near - near clip plane
// far - far clip plane
func (c *Camera) SetupProjection(fov, aspRatio, near, far float32) {
	c.projectionOptions.fov = fov
	c.projectionOptions.aspectRatio = aspRatio
	c.projectionOptions.near = near
	c.projectionOptions.far = far
}

// GetProjectionMatrix returns the projectionMatrix of the camera
func (c *Camera) GetProjectionMatrix() mgl32.Mat4 {
	return mgl32.Perspective(
		c.projectionOptions.fov,
		c.projectionOptions.aspectRatio,
		c.projectionOptions.near,
		c.projectionOptions.far)
}

// GetViewMatrix gets the matrix to transform from world coordinates to
// this camera's coordinates.
// GetViewMatrix returns the viewMatrix of the camera
func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(
		c.cameraPosition,
		c.cameraPosition.Add(c.cameraFrontDirection),
		c.cameraUpDirection)
}

func (c *Camera) updateVectors() {
	radPitch := float64(mgl32.DegToRad(c.pitch))
	radYaw := float64(mgl32.DegToRad(c.yaw))
	c.cameraFrontDirection = mgl32.Vec3{
		float32(math.Cos(radPitch) * math.Cos(radYaw)),
		float32(math.Sin(radPitch)),
		float32(math.Cos(radPitch) * math.Sin(radYaw)),
	}.Normalize()
	// Gram-Schmidt process to figure out right and up vectors
	c.cameraRightDirection = c.worldUp.Cross(c.cameraFrontDirection).Normalize()
	c.cameraUpDirection = c.cameraRightDirection.Cross(c.cameraFrontDirection).Normalize()
}

// UpdateDirection updates the pitch and yaw values.
func (c *Camera) UpdateDirection(amountX, amountY float32) {
	c.pitch = float32(math.Mod(float64(c.pitch+amountY), 360))
	c.yaw = float32(math.Mod(float64(c.yaw+amountX), 360))
	c.updateVectors()
}

// GetPosition returns the current position of the camera.
func (c *Camera) GetPosition() mgl32.Vec3 {
	return c.cameraPosition
}

// SetPosition updates the position of the camera.
func (c *Camera) SetPosition(p mgl32.Vec3) {
	c.cameraPosition = p
}

// GetBoundingObject returns the bounding object of the camera. It is defined as a sphere.
func (c *Camera) GetBoundingObject() *coldet.Sphere {
	return coldet.NewBoundingSphere(
		[3]float32{c.cameraPosition.X(), c.cameraPosition.Y(), c.cameraPosition.Z()},
		defaultCameraRadius)
}

// GetVelocity returns the current velocity of the camera.
func (c *Camera) GetVelocity() float32 {
	return c.velocity
}

// SetVelocity updates the current velocity of the camera.
func (c *Camera) SetVelocity(v float32) {
	c.velocity = v
}

// GetRotationStep returns the rotationStep of the camera.
func (c *Camera) GetRotationStep() float32 {
	return c.rotationStep
}

// SetRotationStep updates the rotationStep of the camera.
func (c *Camera) SetRotationStep(v float32) {
	c.rotationStep = v
}
