package sphere

import (
	"math"

	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/primitives/vertex"

	"github.com/go-gl/mathgl/mgl32"
)

type Sphere struct {
	Points    []mgl32.Vec3
	Indices   []uint32
	TexCoords []mgl32.Vec2
	BO        *boundingobject.BoundingObject
}

// based on this: http://www.songho.ca/opengl/gl_sphere.html
func New(precision int) *Sphere {
	var points []mgl32.Vec3
	var indices []uint32
	var texCoords []mgl32.Vec2
	sectorStep := (2 * math.Pi) / float64(precision)
	stackStep := (math.Pi) / float64(precision)

	for i := 0; i <= precision; i++ {
		stackAngle := (math.Pi / 2) - (float64(i) * stackStep)
		xy := math.Cos(stackAngle)
		z := math.Sin(stackAngle)

		k1 := i * (precision + 1)
		k2 := k1 + precision + 1

		for j := 0; j <= precision; j++ {
			sectorAngle := float64(j) * sectorStep
			x := xy * math.Cos(sectorAngle)
			y := xy * math.Sin(sectorAngle)

			points = append(points, mgl32.Vec3{float32(x), float32(y), float32(z)})
			// textures [0-1], s,t
			s := float32(j) / float32(precision)
			t := float32(i) / float32(precision)
			texCoords = append(texCoords, mgl32.Vec2{s, t})
			// indices
			if !(i == precision || j == precision) {
				if i != 0 {
					indices = append(indices, uint32(k1))
					indices = append(indices, uint32(k2))
					indices = append(indices, uint32(k1+1))
				}
				if i != precision-1 {
					indices = append(indices, uint32(k1+1))
					indices = append(indices, uint32(k2))
					indices = append(indices, uint32(k2+1))
				}

			}

			k1 = k1 + 1
			k2 = k2 + 1
		}
	}
	params := make(map[string]float32)
	params["radius"] = 1.0
	return &Sphere{
		Points:    points,
		Indices:   indices,
		TexCoords: texCoords,
		BO:        boundingobject.New("Sphere", params),
	}
}

// MaterialMeshInput method returns the vertices, indices, bounding object (Sphere) inputs for the NewMaterialMesh function.
func (s *Sphere) MaterialMeshInput() (vertex.Vertices, []uint32, *boundingobject.BoundingObject) {
	var vertices vertex.Vertices
	for i := 0; i < len(s.Points); i++ {
		vertices = append(vertices, vertex.Vertex{
			Position: s.Points[i],
			Normal:   s.Points[i],
		})
	}
	return vertices, s.Indices, s.BO
}

// ColorMeshInput method returns the vertices, indices, bounding object (Sphere) inputs for the NewColorMesh function.
func (s *Sphere) ColoredMeshInput(col []mgl32.Vec3) (vertex.Vertices, []uint32, *boundingobject.BoundingObject) {
	var vertices vertex.Vertices
	for i := 0; i < len(s.Points); i++ {
		vertices = append(vertices, vertex.Vertex{
			Position: s.Points[i],
			Color:    col[i%len(col)],
		})
	}
	return vertices, s.Indices, s.BO
}

// TexturedMeshInput method returns the vertices, indices, bounding object (Sphere) inputs for the NewTexturedMesh function.
func (s *Sphere) TexturedMeshInput() (vertex.Vertices, []uint32, *boundingobject.BoundingObject) {
	var vertices vertex.Vertices
	for i := 0; i < len(s.Points); i++ {
		vertices = append(vertices, vertex.Vertex{
			Position:  s.Points[i],
			Normal:    s.Points[i],
			TexCoords: s.TexCoords[i],
		})
	}
	return vertices, s.Indices, s.BO
}
