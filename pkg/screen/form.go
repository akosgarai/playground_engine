package screen

import (
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/model"
	"github.com/akosgarai/playground_engine/pkg/primitives/rectangle"
	"github.com/akosgarai/playground_engine/pkg/shader"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	frameWidth = float32(0.02)
	FontFile   = "/assets/fonts/Frijole/Frijole-Regular.ttf"
)

type FormScreen struct {
	*ScreenBase
	charset *model.Charset
	frame   *material.Material
	header  string
}

func frameRectangle(width, length float32, position mgl32.Vec3, mat *material.Material, wrapper interfaces.GLWrapper) *mesh.MaterialMesh {
	v, i, _ := rectangle.NewExact(width, length).MeshInput()
	frameMesh := mesh.NewMaterialMesh(v, i, mat, wrapper)
	frameMesh.RotateX(90)
	frameMesh.SetPosition(position)
	return frameMesh
}
func charset(wrapper interfaces.GLWrapper) *model.Charset {
	cs, err := model.LoadCharset(baseDirScreen()+FontFile, 32, 127, 40.0, 72, wrapper)
	if err != nil {
		panic(err)
	}
	cs.SetTransparent(true)
	return cs
}

// NewFormScreen returns a FormScreen. The screen contains a material Frame.
func NewFormScreen(frame *material.Material, label string, wrapper interfaces.GLWrapper) *FormScreen {
	s := newScreenBase()
	bgShaderApplication := shader.NewMenuBackgroundShader(wrapper)
	fgShaderApplication := shader.NewFontShader(wrapper)
	s.AddShader(bgShaderApplication)
	s.AddShader(fgShaderApplication)
	charset := charset(wrapper)
	s.AddModelToShader(charset, fgShaderApplication)
	background := model.New()
	// create frame here.
	bottomFrame := frameRectangle(2.0, frameWidth, mgl32.Vec3{0.0, -0.99, 0.0}, material.Chrome, wrapper)
	background.AddMesh(bottomFrame)
	s.AddModelToShader(background, bgShaderApplication)
	return &FormScreen{
		ScreenBase: s,
		charset:    charset,
		frame:      frame,
		header:     label,
	}
}

// Update loops on the shaderMap, and calls Update function on every Model.
// It also handles the camera movement and rotation, if the camera is set.
func (f *FormScreen) Update(dt, posX, posY float64, keyStore interfaces.RoKeyStore, buttonStore interfaces.RoButtonStore) {
}
