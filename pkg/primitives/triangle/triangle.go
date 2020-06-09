package triangle

import (
	"math"

	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/primitives/vertex"

	"github.com/go-gl/mathgl/mgl32"
)

type Triangle struct {
	Points [3]mgl32.Vec3
	Normal mgl32.Vec3
	BO     *boundingobject.BoundingObject
}

// New returns a Triangle. The inputs are the angles of the triangle.
// The angles are in degrees, not in radian.
func New(alpha, beta, gamma float32) *Triangle {
	sortedAngles := sortAngles(alpha, beta, gamma)
	sinAlpha := float32(math.Sin(float64(mgl32.DegToRad(sortedAngles[0]))))
	sinGamma := float32(math.Sin(float64(mgl32.DegToRad(sortedAngles[2]))))
	lengthSideC := sinGamma / sinAlpha
	pointToRotate := mgl32.Vec3{lengthSideC, 0, 0}
	rotationMatrix := mgl32.HomogRotate3D(mgl32.DegToRad(sortedAngles[1]), mgl32.Vec3{0, -1, 0})
	point := mgl32.TransformCoordinate(pointToRotate, rotationMatrix)

	params := make(map[string]float32)
	params["width"] = 1.0
	params["length"] = point.Z()
	params["height"] = point.Y()
	return &Triangle{
		[3]mgl32.Vec3{
			mgl32.Vec3{-0.5, 0, 0},
			mgl32.Vec3{-0.5 + point.X(), point.Y(), point.Z()},
			mgl32.Vec3{0.5, 0, 0},
		},
		mgl32.Vec3{0, -1, 0},
		boundingobject.New("AABB", params),
	}
}

// This function gets 3 float32 input and returns them in descending order
func sortAngles(a, b, c float32) [3]float32 {
	var smallest, middle, greatest float32
	// sort first 2
	if a <= b {
		smallest = a
		greatest = b
	} else {
		smallest = b
		greatest = a
	}
	if c > greatest {
		middle = greatest
		greatest = c
	} else if c < smallest {
		middle = smallest
		smallest = c
	} else {
		middle = c
	}
	return [3]float32{greatest, middle, smallest}
}

// ColoredMeshInput method returns the vertices, indices, bounding object (AABB) inputs for the New Mesh function.
func (t *Triangle) ColoredMeshInput(col []mgl32.Vec3) (vertex.Vertices, []uint32, *boundingobject.BoundingObject) {
	indices := []uint32{0, 1, 2}
	var vertices vertex.Vertices
	for i := 0; i < 3; i++ {
		vertices = append(vertices, vertex.Vertex{
			Position: t.Points[i],
			Color:    col[i%len(col)],
		})
	}
	return vertices, indices, t.BO
}
