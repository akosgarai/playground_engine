package theme

import (
	"github.com/akosgarai/playground_engine/pkg/material"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	Default = &Theme{
		frameWidth:              float32(2.0),
		frameLength:             float32(0.02),
		frameTopLeftWidth:       float32(0.1),
		detailContentBoxHeight:  float32(0.5),
		frameMaterial:           material.Jade,
		menuItemDefaultMaterial: material.Whiteplastic,
		menuItemHoverMaterial:   material.Ruby,
		menuItemSurfaceTexture:  nil,
		headerLabelColor:        mgl32.Vec3{0, 0, 1},
		labelColor:              mgl32.Vec3{0, 0, 1},
		inputColor:              mgl32.Vec3{0, 0.5, 0},
		backgroundColor:         mgl32.Vec3{0.55, 0.55, 0.55},
	}
)
