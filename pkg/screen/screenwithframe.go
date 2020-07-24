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
	wrapper           interfaces.GLWrapper
	windowWidth       float32
	windowHeight      float32
	frameMaterial     *material.Material
	frameWidth        float32 // the size on the x axis
	frameLength       float32 // the size on the y axis
	frameTopLeftWidth float32 // for supporting the header text, there is an option to give a padding here.
	fov               float32
	labelWidth        float32
}

// NewScreenWithFrameBuilder returns a builder instance.
func NewScreenWithFrameBuilder() *ScreenWithFrameBuilder {
	return &ScreenWithFrameBuilder{
		frameWidth:        BottomFrameWidth,
		frameLength:       BottomFrameLength,
		frameTopLeftWidth: TopLeftFrameWidth,
		fov:               float32(45),
		labelWidth:        float32(0.0),
	}
}

// SetWindowSize sets the windowWidth and the windowHeight parameters.
func (b *ScreenWithFrameBuilder) SetWindowSize(w, h float32) {
	b.windowWidth = w
	b.windowHeight = h
}

// SetFrameSize sets the frameWidth and the frameLength parameters.
func (b *ScreenWithFrameBuilder) SetFrameSize(w, l, r float32) {
	b.frameWidth = w
	b.frameLength = l
	b.frameTopLeftWidth = r
}

// SetWrapper sets the wrapper.
func (b *ScreenWithFrameBuilder) SetWrapper(w interfaces.GLWrapper) {
	b.wrapper = w
}

// SetFrameMaterial sets the material that is used for the frame of the screen.
func (b *ScreenWithFrameBuilder) SetFrameMaterial(m *material.Material) {
	b.frameMaterial = m
}

// SetLabelText sets the label text of the screen.
func (b *ScreenWithFrameBuilder) SetLabelWidth(w float32) {
	b.labelWidth = w
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
	// calculate the positions of the frames:
	halfWidth := b.frameWidth / 2
	halfLength := b.frameLength / 2
	framePosition := halfWidth - halfLength

	frameModel.AddMesh(b.frameRectangle(b.frameWidth, b.frameLength, mgl32.Vec3{0.0, -framePosition, ZFrame}))
	frameModel.AddMesh(b.frameRectangle(b.frameLength, b.frameWidth-b.frameLength, mgl32.Vec3{-framePosition, 0.0, ZFrame}))
	frameModel.AddMesh(b.frameRectangle(b.frameLength, b.frameWidth-b.frameLength, mgl32.Vec3{framePosition, 0.0, ZFrame}))
	frameModel.AddMesh(b.frameRectangle(b.frameTopLeftWidth, b.frameLength, mgl32.Vec3{halfWidth - (b.frameTopLeftWidth / 2), framePosition, ZFrame}))
	frameModel.AddMesh(b.frameRectangle(b.frameWidth-b.frameTopLeftWidth-b.labelWidth, b.frameLength, mgl32.Vec3{(-b.frameTopLeftWidth - b.labelWidth) / 2, framePosition, ZFrame}))
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
	mat := mgl32.Perspective(b.fov, b.windowWidth/b.windowHeight, 0.001, 10)
	cam := camera.NewCamera(mgl32.Vec3{0, 0, -mat[0] * (b.frameWidth / 2)}, mgl32.Vec3{0, -1, 0}, 90.0, 0.0)
	cam.SetupProjection(b.fov, b.windowWidth/b.windowHeight, 0.001, 10.0)
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
