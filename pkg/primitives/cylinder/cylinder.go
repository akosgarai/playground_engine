package cylinder

import (
	"math"

	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/primitives/vertex"

	"github.com/go-gl/mathgl/mgl32"
)

type Cylinder struct {
	Points    []mgl32.Vec3
	Normals   []mgl32.Vec3
	Indices   []uint32
	TexCoords []mgl32.Vec2
	BB        *boundingobject.BoundingObject
}

// New function returns a cylinder.
// Based on the following example:
// http://www.songho.ca/opengl/gl_cylinder.html
// 'rad' - the radius of the circle. 'prec' - the precision of the circle.
// 'length' - the length of the body of the cylinder.
func New(rad float32, prec int, length float32) *Cylinder {
	circleVertices := circleWithRadius(rad, prec)
	params := make(map[string]float32)
	params["width"] = rad
	params["length"] = length
	params["height"] = rad
	c := Cylinder{
		Points:    []mgl32.Vec3{},
		Normals:   []mgl32.Vec3{},
		Indices:   []uint32{},
		TexCoords: []mgl32.Vec2{},
		BB:        boundingobject.New("AABB", params),
	}
	c.buildFromBaseShape(circleVertices, rad, prec, length)
	return &c
}

// NewHalfCircleBased function returns a half circle based cylinder.
// Based on the following example:
// http://www.songho.ca/opengl/gl_cylinder.html
// 'rad' - the radius of the circle. 'prec' - the precision of the circle.
// 'length' - the length of the body of the cylinder.
func NewHalfCircleBased(rad float32, prec int, length float32) *Cylinder {
	baseShapeVertices := halfCircleWithRadius(rad, prec)
	params := make(map[string]float32)
	params["width"] = rad
	params["length"] = length
	params["height"] = rad
	c := Cylinder{
		Points:    []mgl32.Vec3{},
		Normals:   []mgl32.Vec3{},
		Indices:   []uint32{},
		TexCoords: []mgl32.Vec2{},
		BB:        boundingobject.New("AABB", params),
	}
	c.buildFromBaseShape(baseShapeVertices, rad, prec, length)
	return &c
}

// Based on the following example:
// http://www.songho.ca/opengl/gl_cylinder.html
// 'rad' - the radius of the base shape. 'prec' - the precision of the base shape.
// 'length' - the length of the body of the cylinder.
func (c *Cylinder) buildFromBaseShape(baseVertices []float32, rad float32, prec int, length float32) {
	c.calculatePointsNormalsTexCoordsForTheSides(baseVertices, prec, length)

	// Precalculate the center points of the top/bottom shapes.
	baseCenterIndex := len(c.Points)
	topCenterIndex := baseCenterIndex + prec + 1
	c.calculatePointsNormalsTexCoordsForTheBases(baseVertices, rad, prec, length)

	c.calculateIndicesForTheSideSurface(prec)
	c.calculateIndicesForTheBaseSurface(baseCenterIndex, prec)
	c.calculateIndicesForTheTopSurface(topCenterIndex, prec)
}
func (c *Cylinder) calculatePointsNormalsTexCoordsForTheSides(baseVertices []float32, prec int, length float32) {
	precFloat := float32(prec)
	for i := 0; i < 2; i++ {
		height := -length/2 + float32(i)*length
		texCoord := float32(1.0 - i)
		k := 0
		for j := 0; j <= prec; j++ {
			uX := baseVertices[k]
			uY := baseVertices[k+1]
			uZ := baseVertices[k+2]
			// position vector
			c.Points = append(c.Points, mgl32.Vec3{uX, uY, height})
			// normal vectors
			c.Normals = append(c.Normals, mgl32.Vec3{uX, uY, uZ})
			// texture coordinate
			c.TexCoords = append(c.TexCoords, mgl32.Vec2{float32(j) / precFloat, texCoord})
			k = k + 3
		}
	}
}
func (c *Cylinder) calculatePointsNormalsTexCoordsForTheBases(baseVertices []float32, rad float32, prec int, length float32) {
	for i := 0; i < 2; i++ {
		height := -length/2 + float32(i)*length
		normal := float32(-1 + i*2)

		// center point
		c.Points = append(c.Points, mgl32.Vec3{0.0, 0.0, height})
		c.Normals = append(c.Normals, mgl32.Vec3{0.0, 0.0, normal})
		c.TexCoords = append(c.TexCoords, mgl32.Vec2{0.5, 0.5})

		k := 0
		for j := 0; j < prec; j++ {
			// position vector
			c.Points = append(c.Points, mgl32.Vec3{baseVertices[k], baseVertices[k+1], height})
			// normal vectors
			c.Normals = append(c.Normals, mgl32.Vec3{0.0, 0.0, normal})
			// texture coordinate
			c.TexCoords = append(c.TexCoords, mgl32.Vec2{-baseVertices[k]/rad*0.5 + 0.5, -baseVertices[k+1]/rad*0.5 + 0.5})
			k = k + 3
		}
	}
}
func (c *Cylinder) calculateIndicesForTheSideSurface(prec int) {
	k1 := 0        // first vertex index of the bottom
	k2 := prec + 1 // first vertex index of the top
	for i := 0; i < prec; i++ {
		c.Indices = append(c.Indices, uint32(k1))
		c.Indices = append(c.Indices, uint32(k1+1))
		c.Indices = append(c.Indices, uint32(k2))

		c.Indices = append(c.Indices, uint32(k2))
		c.Indices = append(c.Indices, uint32(k1+1))
		c.Indices = append(c.Indices, uint32(k2+1))

		k1 = k1 + 1
		k2 = k2 + 1
	}
}
func (c *Cylinder) calculateIndicesForTheBaseSurface(baseCenterIndex, prec int) {
	// indices for the base surface
	k := baseCenterIndex + 1
	for i := 0; i < prec; i++ {
		if i < prec-1 {
			c.Indices = append(c.Indices, uint32(baseCenterIndex))
			c.Indices = append(c.Indices, uint32(k+1))
			c.Indices = append(c.Indices, uint32(k))
		} else {
			// the last triangle
			c.Indices = append(c.Indices, uint32(baseCenterIndex))
			c.Indices = append(c.Indices, uint32(baseCenterIndex+1))
			c.Indices = append(c.Indices, uint32(k))
		}
		k = k + 1
	}
}
func (c *Cylinder) calculateIndicesForTheTopSurface(topCenterIndex, prec int) {
	// indices for the base surface
	k := topCenterIndex + 1
	for i := 0; i < prec; i++ {
		if i < prec-1 {
			c.Indices = append(c.Indices, uint32(topCenterIndex))
			c.Indices = append(c.Indices, uint32(k+1))
			c.Indices = append(c.Indices, uint32(k))
		} else {
			c.Indices = append(c.Indices, uint32(topCenterIndex))
			c.Indices = append(c.Indices, uint32(topCenterIndex+1))
			c.Indices = append(c.Indices, uint32(k))
		}
		k = k + 1
	}
}

// circleWithRadius returns the position vectors of
// a circle on XY plane,
func circleWithRadius(radius float32, precision int) []float32 {
	var positionVectors []float32
	sectorStep := float64(2*math.Pi) / float64(precision)
	for i := 0; i <= precision; i++ {
		sectorAngle := float64(i) * sectorStep
		positionVectors = append(positionVectors, float32(math.Cos(sectorAngle))*radius)
		positionVectors = append(positionVectors, float32(math.Sin(sectorAngle))*radius)
		positionVectors = append(positionVectors, 0)
	}
	return positionVectors
}

// halfCircleWithRadius returns the position vectors of
// a halfcircle on XY plane.
func halfCircleWithRadius(radius float32, precision int) []float32 {
	var positionVectors []float32
	// we are iterating from 0deg to 180deg.
	sectorStep := float64(2*math.Pi) / float64(precision)
	for i := 0; i <= precision/2; i++ {
		sectorAngle := float64(i) * sectorStep
		positionVectors = append(positionVectors, float32(math.Cos(sectorAngle))*radius)
		positionVectors = append(positionVectors, float32(math.Sin(sectorAngle))*radius)
		positionVectors = append(positionVectors, 0)
	}
	// Draw a line from the end of the curve to the start. (v{-1,0,0} -> v{1,0,0})
	lineStep := float32(2.0) / float32(precision/2)
	for i := 0; i <= precision/2; i++ {
		current := -lineStep*float32(i) + 1.0
		positionVectors = append(positionVectors, current*radius)
		positionVectors = append(positionVectors, 0)
		positionVectors = append(positionVectors, 0)
	}
	return positionVectors
}

// MaterialMeshInput method returns the vertices, indices, bounding object (AABB) inputs for the NewMaterialMesh function.
func (c *Cylinder) MaterialMeshInput() (vertex.Vertices, []uint32, *boundingobject.BoundingObject) {
	var vertices vertex.Vertices
	for i := 0; i < len(c.Points); i++ {
		vertices = append(vertices, vertex.Vertex{
			Position: c.Points[i],
			Normal:   c.Normals[i],
		})
	}
	return vertices, c.Indices, c.BB
}

// ColorMeshInput method returns the vertices, indices, bounding object (AABB) inputs for the NewColorMesh function.
func (c *Cylinder) ColoredMeshInput(col []mgl32.Vec3) (vertex.Vertices, []uint32, *boundingobject.BoundingObject) {
	var vertices vertex.Vertices
	for i := 0; i < len(c.Points); i++ {
		vertices = append(vertices, vertex.Vertex{
			Position: c.Points[i],
			Color:    col[i%len(col)],
		})
	}
	return vertices, c.Indices, c.BB
}

// TexturedMeshInput method returns the vertices, indices, bounding object (AABB) inputs for the NewTexturedMesh function.
func (c *Cylinder) TexturedMeshInput() (vertex.Vertices, []uint32, *boundingobject.BoundingObject) {
	var vertices vertex.Vertices
	for i := 0; i < len(c.Points); i++ {
		vertices = append(vertices, vertex.Vertex{
			Position:  c.Points[i],
			Normal:    c.Normals[i],
			TexCoords: c.TexCoords[i],
		})
	}
	return vertices, c.Indices, c.BB
}
