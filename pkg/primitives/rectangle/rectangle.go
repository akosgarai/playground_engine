package rectangle

import (
	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/primitives/vertex"

	"github.com/go-gl/mathgl/mgl32"
)

type Rectangle struct {
	Points [4]mgl32.Vec3
	Normal mgl32.Vec3
	BB     *boundingobject.BoundingObject
}

// NewSquare creates a rectangle with origo as middle point.
// The normal points to -Y.
// The longest side is scaled to one, and the same downscale is done with the other edge.
// - width represents the length on the X axis.
// - height represents the length on the Z axis.
// ratio = width / length
// ratio == 1 => return NewSquare.
// ratio > 1 => width is the longer -> X [-0.5, 0.5], Y [-1/(ratio*2), 1/(ratio*2)].
// ratio < 1 => length is the longer -> X [-ratio/2, ratio/2], Y [-0.5, 0.5].
func New(width, height float32) *Rectangle {
	ratio := width / height
	if ratio == 1 {
		return NewSquare()
	} else if ratio > 1 {
		return NewExact(1, 1/ratio)
	} else {
		return NewExact(ratio, 1)
	}
}

//NewExact works like the New, but without scaling.
func NewExact(width, height float32) *Rectangle {
	normal := mgl32.Vec3{0, -1, 0}
	points := [4]mgl32.Vec3{
		mgl32.Vec3{-width / 2, 0, -height / 2},
		mgl32.Vec3{width / 2, 0, -height / 2},
		mgl32.Vec3{width / 2, 0, height / 2},
		mgl32.Vec3{-width / 2, 0, height / 2},
	}
	params := make(map[string]float32)
	params["width"] = width
	params["length"] = height
	params["height"] = float32(0.0)
	return &Rectangle{
		Points: points,
		Normal: normal,
		BB:     boundingobject.New("AABB", params),
	}
}

// NewSquare creates a rectangle with origo as middle point.
// Each side is 1 unit long, and it's plane is the X-Z plane.
// The normal points to -Y.
func NewSquare() *Rectangle {
	normal := mgl32.Vec3{0, -1, 0}
	points := [4]mgl32.Vec3{
		mgl32.Vec3{-0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, -0.5},
		mgl32.Vec3{0.5, 0, 0.5},
		mgl32.Vec3{-0.5, 0, 0.5},
	}
	params := make(map[string]float32)
	params["width"] = 1.0
	params["length"] = 1.0
	params["height"] = float32(0.0)
	return &Rectangle{
		Points: points,
		Normal: normal,
		BB:     boundingobject.New("AABB", params),
	}
}

// MeshInput method returns the vertices, indices, bounding object (AABB) inputs for the New Mesh function.
func (r *Rectangle) MeshInput() (vertex.Vertices, []uint32, *boundingobject.BoundingObject) {
	textureCoords := [4]mgl32.Vec2{
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
		{0.0, 0.0},
	}
	indices := []uint32{0, 1, 2, 0, 2, 3}
	var vertices vertex.Vertices
	for i := 0; i < 4; i++ {
		vertices = append(vertices, vertex.Vertex{
			Position:  r.Points[i],
			Normal:    r.Normal,
			TexCoords: textureCoords[i],
		})
	}
	return vertices, indices, r.BB
}

// ColoredMeshInput method returns the vertices, indices, bounding object (AABB) inputs for the New Mesh function.
func (r *Rectangle) ColoredMeshInput(col []mgl32.Vec3) (vertex.Vertices, []uint32, *boundingobject.BoundingObject) {
	indices := []uint32{0, 1, 2, 0, 2, 3}
	var vertices vertex.Vertices
	for i := 0; i < 4; i++ {
		vertices = append(vertices, vertex.Vertex{
			Position: r.Points[i],
			Color:    col[i%len(col)],
		})
	}
	return vertices, indices, r.BB
}

// TexturedColoredMeshInput method returns the vertices, indices, bounding object (AABB) inputs for the NewTexturedColoredMesh function.
func (r *Rectangle) TexturedColoredMeshInput(col []mgl32.Vec3) (vertex.Vertices, []uint32, *boundingobject.BoundingObject) {
	textureCoords := [4]mgl32.Vec2{
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
		{0.0, 0.0},
	}
	indices := []uint32{0, 1, 2, 0, 2, 3}
	var vertices vertex.Vertices
	for i := 0; i < 4; i++ {
		vertices = append(vertices, vertex.Vertex{
			Position:  r.Points[i],
			Color:     col[i%len(col)],
			TexCoords: textureCoords[i],
		})
	}
	return vertices, indices, r.BB
}
