package screen

import (
	"fmt"

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
	frameWidth             float32 // the size on the x axis
	frameLength            float32 // the size on the y axis
	detailContentBox       interfaces.Mesh
	detailContentBoxHeight float32
}

// tmp function for testing
func (f *ScreenWithFrame) Update(dt, posX, posY float64, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
}

// tmp function for testing
func (f *ScreenWithFrame) CharCallback(char rune, wrapper interfaces.GLWrapper) {
}

// GetFullWidth returns the width of the drawable screen. (width - 2*length)
func (f *ScreenWithFrame) GetFullWidth() float32 {
	return f.frameWidth - (2 * f.frameLength)
}

type ScreenWithFrameBuilder struct {
	wrapper                  interfaces.GLWrapper
	windowWidth              float32
	windowHeight             float32
	frameMaterial            *material.Material
	frameWidth               float32 // the size on the x axis
	frameLength              float32 // the size on the y axis
	frameTopLeftWidth        float32 // for supporting the header text, there is an option to give a padding here.
	fov                      float32
	labelWidth               float32
	detailContentBoxHeight   float32
	detailContentBoxMaterial *material.Material
}

// NewScreenWithFrameBuilder returns a builder instance.
func NewScreenWithFrameBuilder() *ScreenWithFrameBuilder {
	return &ScreenWithFrameBuilder{
		frameWidth:               DefaultFrameWidth,
		frameLength:              DefaultFrameLength,
		frameTopLeftWidth:        TopLeftFrameWidth,
		fov:                      float32(45),
		labelWidth:               float32(0.0),
		detailContentBoxHeight:   float32(0.0),
		detailContentBoxMaterial: DefaultFormItemMaterial,
		wrapper:                  nil,
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

// SetDetailContentBoxMaterial sets the material that is used for the  detailcontentbox.
func (b *ScreenWithFrameBuilder) SetDetailContentBoxMaterial(m *material.Material) {
	b.detailContentBoxMaterial = m
}

// SetLabelText sets the label text width of the screen.
func (b *ScreenWithFrameBuilder) SetLabelWidth(w float32) {
	b.labelWidth = w
}

// SetDetailContentBoxHeight sets the height of the detailContentBox.
func (b *ScreenWithFrameBuilder) SetDetailContentBoxHeight(h float32) {
	b.detailContentBoxHeight = h
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
	framePositionVertical := halfWidth - (b.frameLength / 2)
	framePositionHorizontal := halfWidth - (b.frameLength / 2)
	fullWithoutFrame := b.frameWidth - (2 * b.frameLength)
	// variables for aspect ratio.
	aspWidth := float32(1.0)
	aspHeight := float32(1.0)
	if b.windowWidth > b.windowHeight {
		aspWidth = float32(b.windowHeight) / float32(b.windowWidth)
	}
	if b.windowWidth < b.windowHeight {
		aspHeight = float32(b.windowWidth) / float32(b.windowHeight)
	}

	// bottom frame. it supposed to be full width long.
	frameModel.AddMesh(b.frameRectangle(b.frameWidth, b.frameLength*aspHeight, mgl32.Vec3{0.0, -1, ZFrame}))
	// left, right
	frameModel.AddMesh(b.frameRectangle(b.frameLength*aspWidth, b.frameWidth-b.frameLength*aspHeight, mgl32.Vec3{-framePositionHorizontal, 0.0, ZFrame}))
	frameModel.AddMesh(b.frameRectangle(b.frameLength*aspWidth, b.frameWidth-b.frameLength*aspHeight, mgl32.Vec3{framePositionHorizontal, 0.0, ZFrame}))
	// top
	frameModel.AddMesh(b.frameRectangle(b.frameTopLeftWidth*aspWidth, b.frameLength*aspHeight, mgl32.Vec3{(halfWidth - (b.frameTopLeftWidth / 2)) * aspWidth, framePositionVertical * aspHeight, ZFrame}))
	frameModel.AddMesh(b.frameRectangle((b.frameWidth-b.frameTopLeftWidth-b.labelWidth)*aspWidth, b.frameLength*aspHeight, mgl32.Vec3{((-b.frameTopLeftWidth - b.labelWidth) / 2) * aspWidth, framePositionVertical * aspHeight, ZFrame}))
	var detailContentBox interfaces.Mesh
	if b.detailContentBoxHeight > 0.0 {
		detailContainerPosition := mgl32.Vec3{0.0, (-halfWidth + b.frameLength + b.detailContentBoxHeight/2) * aspHeight, ZFrame}
		detailContentBox = b.frameRectangleWithMaterial(fullWithoutFrame*aspWidth, b.detailContentBoxHeight*aspHeight, detailContainerPosition, b.detailContentBoxMaterial)
		detailContentBox.RotateX(-180)
		detailContentBox.RotateY(180)
		frameModel.AddMesh(detailContentBox)
	} else {
		detailContentBox = nil
	}
	s.AddModelToShader(frameModel, frameShaderApplication)
	directionalLightSource := light.NewDirectionalLight([4]mgl32.Vec3{
		DirectionalLightDirection,
		DirectionalLightAmbient,
		DirectionalLightDiffuse,
		DirectionalLightSpecular,
	})
	s.AddDirectionalLightSource(directionalLightSource, [4]string{"dirLight[0].direction", "dirLight[0].ambient", "dirLight[0].diffuse", "dirLight[0].specular"})

	return &ScreenWithFrame{
		ScreenBase:             s,
		frameWidth:             b.frameWidth,
		frameLength:            b.frameLength,
		detailContentBox:       detailContentBox,
		detailContentBoxHeight: b.detailContentBoxHeight,
	}
}

// It creates a new camera with the necessary setup
func (b *ScreenWithFrameBuilder) defaultCamera() *camera.DefaultCamera {
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
	fmt.Printf("Frame mesh (%f * %f) to %v position.\n", width, length, position)
	v, i, _ := rectangle.NewExact(width, length).MeshInput()
	var tex texture.Textures
	tex.TransparentTexture(1, 1, 128, "tex.diffuse", b.wrapper)
	tex.TransparentTexture(1, 1, 128, "tex.specular", b.wrapper)
	frameMesh := mesh.NewTexturedMaterialMesh(v, i, tex, mat, b.wrapper)
	frameMesh.RotateX(90)
	frameMesh.SetPosition(position)
	return frameMesh
}
