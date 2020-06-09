package cylinder

import (
	"testing"

	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/go-gl/mathgl/mgl32"
)

func TestNew(t *testing.T) {
	testData := []struct {
		input                   mgl32.Vec3
		expectedLengthPoint     int
		expectedLengthNormal    int
		expectedLengthIndices   int
		expectedLengthTexCoords int
	}{
		{mgl32.Vec3{1, 2, 1}, 12, 12, 24, 12},
	}
	for _, tt := range testData {
		cyl := New(tt.input.X(), int(tt.input.Y()), tt.input.Z())
		if len(cyl.Points) != tt.expectedLengthPoint {
			t.Errorf("Invalid point length. instead of '%d', we have '%d'\n", tt.expectedLengthPoint, len(cyl.Points))
		}
		if len(cyl.Normals) != tt.expectedLengthNormal {
			t.Errorf("Invalid normals length. instead of '%d', we have '%d'\n", tt.expectedLengthNormal, len(cyl.Normals))
		}
		if len(cyl.Indices) != tt.expectedLengthIndices {
			t.Errorf("Invalid indices length. instead of '%d', we have '%d'\n", tt.expectedLengthIndices, len(cyl.Indices))
		}
		if len(cyl.TexCoords) != tt.expectedLengthTexCoords {
			t.Errorf("Invalid tex coords length. instead of '%d', we have '%d'\n", tt.expectedLengthTexCoords, len(cyl.TexCoords))
		}
	}
}
func TestNewHalfCircleBased(t *testing.T) {
	testData := []struct {
		input                   mgl32.Vec3
		expectedLengthPoint     int
		expectedLengthNormal    int
		expectedLengthIndices   int
		expectedLengthTexCoords int
	}{
		{mgl32.Vec3{1, 2, 1}, 12, 12, 24, 12},
	}
	for _, tt := range testData {
		cyl := NewHalfCircleBased(tt.input.X(), int(tt.input.Y()), tt.input.Z())
		if len(cyl.Points) != tt.expectedLengthPoint {
			t.Errorf("Invalid point length. instead of '%d', we have '%d'\n", tt.expectedLengthPoint, len(cyl.Points))
		}
		if len(cyl.Normals) != tt.expectedLengthNormal {
			t.Errorf("Invalid normals length. instead of '%d', we have '%d'\n", tt.expectedLengthNormal, len(cyl.Normals))
		}
		if len(cyl.Indices) != tt.expectedLengthIndices {
			t.Errorf("Invalid indices length. instead of '%d', we have '%d'\n", tt.expectedLengthIndices, len(cyl.Indices))
		}
		if len(cyl.TexCoords) != tt.expectedLengthTexCoords {
			t.Errorf("Invalid tex coords length. instead of '%d', we have '%d'\n", tt.expectedLengthTexCoords, len(cyl.TexCoords))
		}
	}
}
func TestMaterialMeshInput(t *testing.T) {
	testData := []struct {
		input                  mgl32.Vec3
		expectedLengthVertices int
		expectedLengthIndices  int
	}{
		{mgl32.Vec3{1, 2, 1}, 12, 24},
	}
	for _, tt := range testData {
		cyl := New(tt.input.X(), int(tt.input.Y()), tt.input.Z())
		v, i, _ := cyl.MaterialMeshInput()
		if len(i) != tt.expectedLengthIndices {
			t.Errorf("Invalid indices length. instead of '%d', we have '%d'\n", tt.expectedLengthIndices, len(i))
		}
		if len(v) != tt.expectedLengthVertices {
			t.Errorf("Invalid vertices length. instead of '%d', we have '%d'\n", tt.expectedLengthVertices, len(v))
		}
	}
}
func TestColoredMeshInput(t *testing.T) {
	testData := []struct {
		input                  mgl32.Vec3
		expectedLengthVertices int
		expectedLengthIndices  int
	}{
		{mgl32.Vec3{1, 2, 1}, 12, 24},
	}
	for _, tt := range testData {
		cyl := New(tt.input.X(), int(tt.input.Y()), tt.input.Z())
		v, i, _ := cyl.ColoredMeshInput([]mgl32.Vec3{mgl32.Vec3{1, 1, 1}})
		if len(i) != tt.expectedLengthIndices {
			t.Errorf("Invalid indices length. instead of '%d', we have '%d'\n", tt.expectedLengthIndices, len(i))
		}
		if len(v) != tt.expectedLengthVertices {
			t.Errorf("Invalid vertices length. instead of '%d', we have '%d'\n", tt.expectedLengthVertices, len(v))
		}
	}
}
func TestTexturedMeshInput(t *testing.T) {
	testData := []struct {
		input                  mgl32.Vec3
		expectedLengthVertices int
		expectedLengthIndices  int
	}{
		{mgl32.Vec3{1, 2, 1}, 12, 24},
	}
	for _, tt := range testData {
		cyl := New(tt.input.X(), int(tt.input.Y()), tt.input.Z())
		v, i, _ := cyl.TexturedMeshInput()
		if len(i) != tt.expectedLengthIndices {
			t.Errorf("Invalid indices length. instead of '%d', we have '%d'\n", tt.expectedLengthIndices, len(i))
		}
		if len(v) != tt.expectedLengthVertices {
			t.Errorf("Invalid vertices length. instead of '%d', we have '%d'\n", tt.expectedLengthVertices, len(v))
		}
	}
}
func TestCircleWithRadius(t *testing.T) {
	testData := []struct {
		radius         float32
		precision      int
		expectedOutput []float32
	}{
		{0.75, 4, []float32{0.75, 0, 0, 0, 0.75, 0, -0.75, 0, 0, 0, -0.75, 0, 0.75, 0, 0}},
		{1.0, 4, []float32{1.0, 0, 0, 0, 1.0, 0, -1.0, 0, 0, 0, -1.0, 0, 1.0, 0, 0}},
	}
	for _, tt := range testData {
		points := circleWithRadius(tt.radius, tt.precision)
		if len(points) != len(tt.expectedOutput) {
			t.Error("Invalid number of generated points.")
		}
		for index := range points {
			if !testhelper.Float32ApproxEqual(points[index], tt.expectedOutput[index], 0.001) {
				t.Errorf("The given points are different. Instead of '%f', we have '%f'.\n", tt.expectedOutput[index], points[index])
			}
		}
	}
}
func TestCalculatePointsNormalsTexCoordsForTheSides(t *testing.T) {
	testData := []struct {
		precision    int
		length       float32
		baseVertices []float32
	}{
		{4, 1.0, circleWithRadius(1.0, 4)},
	}
	for _, tt := range testData {
		var cylinder Cylinder
		cylinder.calculatePointsNormalsTexCoordsForTheSides(tt.baseVertices, tt.precision, tt.length)
		expectedLength := tt.precision*2 + 2
		if len(cylinder.Points) != expectedLength {
			t.Errorf("Invalid number of Points after side calculation. instead of '%d', we have '%d'.\n", expectedLength, len(cylinder.Points))
			t.Log(cylinder.Points)
		}
		if len(cylinder.Normals) != expectedLength {
			t.Errorf("Invalid number of Normals after side calculation. instead of '%d', we have '%d'.\n", expectedLength, len(cylinder.Normals))
			t.Log(cylinder.Normals)
		}
		if len(cylinder.TexCoords) != expectedLength {
			t.Errorf("Invalid number of TexCoords after side calculation. instead of '%d', we have '%d'.\n", expectedLength, len(cylinder.TexCoords))
			t.Log(cylinder.TexCoords)
		}
	}
}
func TestCalculatePointsNormalsTexCoordsForTheBases(t *testing.T) {
	testData := []struct {
		precision    int
		length       float32
		rad          float32
		baseVertices []float32
	}{
		{4, 1.0, 1.0, circleWithRadius(1.0, 4)},
	}
	for _, tt := range testData {
		var cylinder Cylinder
		cylinder.calculatePointsNormalsTexCoordsForTheBases(tt.baseVertices, tt.rad, tt.precision, tt.length)
		expectedLength := tt.precision*2 + 2
		if len(cylinder.Points) != expectedLength {
			t.Errorf("Invalid number of Points after base calculation. instead of '%d', we have '%d'.\n", expectedLength, len(cylinder.Points))
			t.Log(cylinder.Points)
		}
		if len(cylinder.Normals) != expectedLength {
			t.Errorf("Invalid number of Normals after base calculation. instead of '%d', we have '%d'.\n", expectedLength, len(cylinder.Normals))
			t.Log(cylinder.Normals)
		}
		if len(cylinder.TexCoords) != expectedLength {
			t.Errorf("Invalid number of TexCoords after base calculation. instead of '%d', we have '%d'.\n", expectedLength, len(cylinder.TexCoords))
			t.Log(cylinder.TexCoords)
		}
	}
}
func TestCalculateIndicesForTheSideSurface(t *testing.T) {
	testData := []int{3, 4, 5, 6, 7, 7, 10, 100}
	for _, prec := range testData {
		var cylinder Cylinder
		cylinder.calculateIndicesForTheSideSurface(prec)
		if len(cylinder.Indices) != 6*prec {
			t.Errorf("Invalid number of indices. instead of '%d', we have '%d'.\n", 6*prec, len(cylinder.Indices))
		}
	}
}
func TestCalculateIndicesForTheBaseSurface(t *testing.T) {
	testData := []int{3, 4, 5, 6, 7, 7, 10, 100}
	baseCenterIndex := 0
	for _, prec := range testData {
		var cylinder Cylinder
		cylinder.calculateIndicesForTheBaseSurface(baseCenterIndex, prec)
		if len(cylinder.Indices) != 3*prec {
			t.Errorf("Invalid number of indices. instead of '%d', we have '%d'.\n", 3*prec, len(cylinder.Indices))
		}
	}
}
func TestCalculateIndicesForTheTopSurface(t *testing.T) {
	testData := []int{3, 4, 5, 6, 7, 7, 10, 100}
	topCenterIndex := 0
	for _, prec := range testData {
		var cylinder Cylinder
		cylinder.calculateIndicesForTheTopSurface(topCenterIndex, prec)
		if len(cylinder.Indices) != 3*prec {
			t.Errorf("Invalid number of indices. instead of '%d', we have '%d'.\n", 3*prec, len(cylinder.Indices))
		}
	}
}
func TestResults(t *testing.T) {
	rad := float32(1.0)
	prec := 4
	length := float32(1)
	expectedBaseCenter := mgl32.Vec3{0, 0, -0.5}
	expectedtopCenter := mgl32.Vec3{0, 0, 0.5}
	c := Cylinder{
		Points:    []mgl32.Vec3{},
		Normals:   []mgl32.Vec3{},
		Indices:   []uint32{},
		TexCoords: []mgl32.Vec2{},
	}
	circleVertices := circleWithRadius(rad, prec)
	c.calculatePointsNormalsTexCoordsForTheSides(circleVertices, prec, length)
	// At this time we should have prec*2=2 point.
	if len(c.Points) != prec*2+2 {
		t.Errorf("Invalid number of points after the first calculation. Instead of '%d', we have '%d'.\n", prec*2+2, len(c.Points))
	}
	// store the expected center points.
	baseCenterIndex := len(c.Points)
	topCenterIndex := baseCenterIndex + prec + 1
	c.calculatePointsNormalsTexCoordsForTheBases(circleVertices, rad, prec, length)
	if c.Points[baseCenterIndex] != expectedBaseCenter {
		t.Error("Invalid center point for the base")
		t.Log(c.Points[baseCenterIndex])
	}
	if c.Points[topCenterIndex] != expectedtopCenter {
		t.Errorf("Invalid center point for the top. index:'%d'\n", topCenterIndex)
		t.Log(c.Points[topCenterIndex])
	}
	// lengths should be fine.
	if len(c.Points) != prec*4+4 {
		t.Errorf("Invalid number of points after the second calculation. Instead of '%d', we have '%d'.\n", prec*2+2, len(c.Points))
	}
	// check indices.
	c.calculateIndicesForTheSideSurface(prec)
	if len(c.Indices) != 6*prec {
		t.Errorf("Invalid number of indices after side setup. Instead of '%d', we have '%d'.\n", 6*prec, len(c.Indices))
	}
	c.calculateIndicesForTheBaseSurface(baseCenterIndex, prec)
	if len(c.Indices) != 9*prec {
		t.Errorf("Invalid number of indices after base setup. Instead of '%d', we have '%d'.\n", 9*prec, len(c.Indices))
	}
	c.calculateIndicesForTheTopSurface(topCenterIndex, prec)
	if len(c.Indices) != 12*prec {
		t.Errorf("Invalid number of indices after top setup. Instead of '%d', we have '%d'.\n", 12*prec, len(c.Indices))
	}
}
func TestResultsHalfCircle(t *testing.T) {
	rad := float32(1.0)
	prec := 4
	length := float32(1)
	expectedBaseCenter := mgl32.Vec3{0, 0, -0.5}
	expectedtopCenter := mgl32.Vec3{0, 0, 0.5}
	c := Cylinder{
		Points:    []mgl32.Vec3{},
		Normals:   []mgl32.Vec3{},
		Indices:   []uint32{},
		TexCoords: []mgl32.Vec2{},
	}
	halfCircleVertices := halfCircleWithRadius(rad, prec)
	c.calculatePointsNormalsTexCoordsForTheSides(halfCircleVertices, prec, length)
	// At this time we should have prec*2=2 point.
	if len(c.Points) != prec*2+2 {
		t.Errorf("Invalid number of points after the first calculation. Instead of '%d', we have '%d'.\n", prec*2+2, len(c.Points))
	}
	// store the expected center points.
	baseCenterIndex := len(c.Points)
	topCenterIndex := baseCenterIndex + prec + 1
	c.calculatePointsNormalsTexCoordsForTheBases(halfCircleVertices, rad, prec, length)
	if c.Points[baseCenterIndex] != expectedBaseCenter {
		t.Error("Invalid center point for the base")
		t.Log(c.Points[baseCenterIndex])
	}
	if c.Points[topCenterIndex] != expectedtopCenter {
		t.Errorf("Invalid center point for the top. index:'%d'\n", topCenterIndex)
		t.Log(c.Points[topCenterIndex])
	}
	// lengths should be fine.
	if len(c.Points) != prec*4+4 {
		t.Errorf("Invalid number of points after the second calculation. Instead of '%d', we have '%d'.\n", prec*2+2, len(c.Points))
	}
	// check indices.
	c.calculateIndicesForTheSideSurface(prec)
	if len(c.Indices) != 6*prec {
		t.Errorf("Invalid number of indices after side setup. Instead of '%d', we have '%d'.\n", 6*prec, len(c.Indices))
	}
	c.calculateIndicesForTheBaseSurface(baseCenterIndex, prec)
	if len(c.Indices) != 9*prec {
		t.Errorf("Invalid number of indices after base setup. Instead of '%d', we have '%d'.\n", 9*prec, len(c.Indices))
	}
	c.calculateIndicesForTheTopSurface(topCenterIndex, prec)
	if len(c.Indices) != 12*prec {
		t.Errorf("Invalid number of indices after top setup. Instead of '%d', we have '%d'.\n", 12*prec, len(c.Indices))
	}
}
