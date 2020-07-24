package screen

import (
	"github.com/akosgarai/playground_engine/pkg/camera"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/shader"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

type ScreenWithFrame struct {
	*ScreenBase
}

// tmp function for testing
func (f *ScreenWithFrame) Update(dt, posX, posY float64, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
}

// tmp function for testing
func (f *ScreenWithFrame) CharCallback(char rune, wrapper interfaces.GLWrapper) {
}

type ScreenWithFrameBuilder struct {
	wrapper       interfaces.GLWrapper
	windowWidth   float32
	windowHeight  float32
	frameMaterial *material.Material
}

// NewScreenWithFrameBuilder returns a builder instance.
func NewScreenWithFrameBuilder() *ScreenWithFrameBuilder {
	return &ScreenWithFrameBuilder{}
}

// SetWindowSize sets the windowWidth and the windowHeight parameters.
func (b *ScreenWithFrameBuilder) SetWindowSize(w, h float32) {
	b.windowWidth = w
	b.windowHeight = h
}

// SetWrapper sets the wrapper.
func (b *ScreenWithFrameBuilder) SetWrapper(w interfaces.GLWrapper) {
	b.wrapper = w
}

// SetFrameMaterial sets the material that is used for the frame of the screen.
func (b *ScreenWithFrameBuilder) SetFrameMaterial(m *material.Material) {
	b.frameMaterial = m
}

func (b *ScreenWithFrameBuilder) Build() *ScreenWithFrame {
	if b.wrapper == nil {
		panic("Wrapper is missing for the build process.")
	}
	s := newScreenBase()
	s.SetWrapper(b.wrapper)
	s.SetCamera(b.defaultCamera())
	s.SetWindowSize(b.windowWidth, b.windowHeight)

	frameShaderApplication := shader.NewMaterialShader(b.wrapper)
	s.AddShader(frameShaderApplication)

	frameModel := model.New()
	frameModel.AddMesh(b.frameRectangle(BottomFrameWidth, BottomFrameLength, mgl32.Vec3{0.0, -0.99, ZFrame}))
	frameModel.AddMesh(b.frameRectangle(SideFrameWidth, SideFrameLength, mgl32.Vec3{-0.99, 0.0, ZFrame}))
	frameModel.AddMesh(b.frameRectangle(SideFrameWidth, SideFrameLength, mgl32.Vec3{0.99, 0.0, ZFrame}))
	frameModel.AddMesh(b.frameRectangle(TopLeftFrameWidth, BottomFrameLength, mgl32.Vec3{0.95, 0.99, ZFrame}))
	frameModel.AddMesh(b.frameRectangle(2.0-TopLeftFrameWidth, BottomFrameLength, mgl32.Vec3{(-TopLeftFrameWidth) / 2, 0.99, ZFrame}))
	s.AddModelToShader(frameModel, frameShaderApplication)
	directionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
	})
	s.AddDirectionalLightSource(directionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})

	return &ScreenWithFrame{
		ScreenBase: s,
	}
}

// It creates a new camera with the necessary setup
func (b *ScreenWithFrameBuilder) defaultCamera() *camera.Camera {
	cam := camera.NewCamera(mgl32.Vec3{0, 0, -1.8}, mgl32.Vec3{0, -1, 0}, 90.0, 0.0)
	cam.SetupProjection(45, b.windowWidth/b.windowHeight, 0.001, 10.0)
	return cam
}

// It creates a rectangle for the screen frame.
func (b *ScreenWithFrameBuilder) frameRectangle(width, length float32, position mgl32.Vec3) *mesh.TexturedMaterialMesh {
	return b.frameRectangleWithMaterial(width, length, position, b.frameMaterial)
}
func (b *ScreenWithFrameBuilder) frameRectangleWithMaterial(width, length float32, position mgl32.Vec3, mat *material.Material) *mesh.TexturedMaterialMesh {
	v, i, _ := rectangle.NewExact(width, length).MeshInput()
	var tex texture.Textures
	tex.TransparentTexture(1, 1, 128, "tex.diffuse", b.wrapper)
	tex.TransparentTexture(1, 1, 128, "tex.specular", b.wrapper)
	frameMesh := mesh.NewTexturedMaterialMesh(v, i, tex, mat, b.wrapper)
	frameMesh.RotateX(90)
	frameMesh.SetPosition(position)
	return frameMesh
}
