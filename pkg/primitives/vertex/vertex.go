package vertex

import "github.com/go-gl/mathgl/mgl32"

const (
	POSITION_NORMAL          = 1
	POSITION_NORMAL_TEXCOORD = 2
	POSITION_COLOR_SIZE      = 3
	POSITION_COLOR           = 4
	POSITION_COLOR_TEXCOORD  = 5
)

type Vertex struct {
	// Position vector
	Position mgl32.Vec3
	// Normal vector
	Normal mgl32.Vec3
	// Texture coordinates for textured objects.
	// As the textures are 2D, we need vec2 for storing the coordinates.
	TexCoords mgl32.Vec2
	// Color vector
	Color mgl32.Vec3
	// Point size for points
	PointSize float32
}

type Vertices []Vertex

// Get returns the vertex array object as []float32. It's input describes the
// fields and the order that we want to add to the vao.
// Modes: 'POSITION_NORMAL' (1) the position and the normal vectors are returned in this order.
// 'POSITION_NORMAL_TEXCOORD' (2) the position, normal and the tex coords vectors are
// returned in this order. 'POSITION_COLOR_SIZE' (3) the position, the color and the point size
// are retuned in this order. 'POSITION_COLOR' (4) the position and the color vectors are
// returned in this order. 'POSITION_COLOR_TEXCOORD' (5) the position, the color and the tex
// coords vectors are returned in this order.
func (v Vertices) Get(resultMode int) []float32 {
	if resultMode == POSITION_COLOR_SIZE {
		return v.getPoint()
	}
	if resultMode == POSITION_COLOR {
		return v.positionColor()
	}
	if resultMode == POSITION_COLOR_TEXCOORD {
		return v.positionColorTexCoords()
	}
	var vao []float32
	for _, vertex := range v {
		vao = append(vao, vertex.Position.X())
		vao = append(vao, vertex.Position.Y())
		vao = append(vao, vertex.Position.Z())

		vao = append(vao, vertex.Normal.X())
		vao = append(vao, vertex.Normal.Y())
		vao = append(vao, vertex.Normal.Z())

		if resultMode == POSITION_NORMAL_TEXCOORD {
			vao = append(vao, vertex.TexCoords.X())
			vao = append(vao, vertex.TexCoords.Y())
		}
	}

	return vao
}
func (v Vertices) getPoint() []float32 {
	var vao []float32
	for _, vertex := range v {
		vao = append(vao, vertex.Position.X())
		vao = append(vao, vertex.Position.Y())
		vao = append(vao, vertex.Position.Z())

		vao = append(vao, vertex.Color.X())
		vao = append(vao, vertex.Color.Y())
		vao = append(vao, vertex.Color.Z())

		vao = append(vao, vertex.PointSize)
	}

	return vao
}
func (v Vertices) positionColor() []float32 {
	var vao []float32
	for _, vertex := range v {
		vao = append(vao, vertex.Position.X())
		vao = append(vao, vertex.Position.Y())
		vao = append(vao, vertex.Position.Z())

		vao = append(vao, vertex.Color.X())
		vao = append(vao, vertex.Color.Y())
		vao = append(vao, vertex.Color.Z())
	}

	return vao
}
func (v Vertices) positionColorTexCoords() []float32 {
	var vao []float32
	for _, vertex := range v {
		vao = append(vao, vertex.Position.X())
		vao = append(vao, vertex.Position.Y())
		vao = append(vao, vertex.Position.Z())

		vao = append(vao, vertex.Color.X())
		vao = append(vao, vertex.Color.Y())
		vao = append(vao, vertex.Color.Z())

		vao = append(vao, vertex.TexCoords.X())
		vao = append(vao, vertex.TexCoords.Y())
	}

	return vao
}

// Add appends the given vertex to the vertices.
func (v *Vertices) Add(ver Vertex) {
	*v = append(*v, ver)
}
