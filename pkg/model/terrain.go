package model

import (
	"errors"
	"fmt"
	"math/rand"
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/vertex"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	defaultTerrainWidth  = 10
	defaultTerrainLength = 10
	defaultIterations    = 1
	defaultMinHeight     = float32(0.0)
	defaultMaxHeight     = float32(0.0)
	defaultSeed          = int64(0)
	distanceTolerance    = float32(0.01)
)

var (
	ErrorNotAboveTheSurface = errors.New("Not above the surface")
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Terrain struct {
	Model
	heightMap     [][]float32
	width, length int
	debugMode     bool
}

// GetTerrain returns the terrain mesh
func (t *Terrain) GetTerrain() interfaces.Mesh {
	return t.meshes[0]
}

// HeightAtPos returns the height value at the given position, and nil. In case of the position
// is not under or above the surface, it returns 01 and error.
func (t *Terrain) HeightAtPos(pos mgl32.Vec3) (float32, error) {
	tMesh := t.GetTerrain()
	scaleTr := tMesh.ScaleTransformation()
	scaleX := scaleTr[0]
	scaleZ := scaleTr[10]
	posX := pos.X() / scaleX
	posZ := pos.Z() / scaleZ
	// Exclude points that are not above or under the mesh.
	if posX > float32(t.width)/2.0 || posX < float32(-t.width)/2 || posZ > float32(t.length)/2.0 || posZ < float32(-t.length)/2.0 {
		return -1, ErrorNotAboveTheSurface
	}
	var posX1, posX2, posZ1, posZ2 int
	var mapIndexX1, mapIndexX2, mapIndexZ1, mapIndexZ2 int
	posX1 = int(posX)
	if float32(posX1) > posX {
		posX1 = posX1 - 1
	}
	posZ1 = int(posZ)
	if float32(posZ1) > posZ {
		posZ1 = posZ1 - 1
	}
	posX2 = posX1 + 1
	posZ2 = posZ1 + 1

	mapIndexX1 = int(t.width/2 + posX1)
	mapIndexX2 = int(t.width/2 + posX2)
	mapIndexZ1 = int(t.length/2 + posZ1)
	mapIndexZ2 = int(t.length/2 + posZ2)

	wX := posX - float32(int(posX))
	if wX < 0 {
		wX = 1.0 + wX
	}
	wZ := posZ - float32(int(posZ))
	if wZ < 0 {
		wZ = 1.0 + wZ
	}
	if t.debugMode {
		fmt.Printf("Terrain.HeightAtPos:\n\twx:\t%f\n\twZ:\t%f\nMapIndex:\n\tZ1:\t%d\n\tZ2:\t%d\n\tX1:\t%d\n\tX2:\t%d\n", wX, wZ, mapIndexZ1, mapIndexZ2, mapIndexX1, mapIndexX2)
	}
	height := (t.heightMap[mapIndexZ2][mapIndexX1]*(1.0-wX)+t.heightMap[mapIndexZ2][mapIndexX2]*wX)*wZ + (t.heightMap[mapIndexZ1][mapIndexX1]*(1.0-wX)+t.heightMap[mapIndexZ1][mapIndexX2]*wX)*(1-wZ)
	return height, nil
}

// CollideTestWithSphere is the collision detection function for heightmap vs sphere.
func (t *Terrain) CollideTestWithSphere(boundingSphere *coldet.Sphere) bool {
	height, err := t.HeightAtPos(mgl32.Vec3{boundingSphere.X(), boundingSphere.Y(), boundingSphere.Z()})
	if err != nil {
		return false
	}
	bpPos := t.GetTerrain().GetPosition().Add(mgl32.Vec3{boundingSphere.X(), height, boundingSphere.Z()})
	boundingPoint := coldet.NewBoundingPoint([3]float32{bpPos.X(), bpPos.Y(), bpPos.Z()})
	return coldet.CheckPointInSphere(*boundingPoint, *boundingSphere)
}

// Update function does nothing.
func (t *Terrain) Update(dt float64) {
}

// TerrainBuilder is a helper structure for generating terrain. It has a fluid API,
// so that the settings could be chained.
type TerrainBuilder struct {
	width, length, iterations int
	minH, maxH                float32
	seed                      int64
	heightMap                 [][]float32
	minHIsDefault             bool
	peakProbability           int
	cliffProbability          int
	wrapper                   interfaces.GLWrapper
	tex                       texture.Textures
	scale                     mgl32.Vec3
	debugMode                 bool
}

// NewTerrainBuilder returns a TerrainBuilder with default settings.
func NewTerrainBuilder() *TerrainBuilder {
	return &TerrainBuilder{
		width:            defaultTerrainWidth,
		length:           defaultTerrainLength,
		iterations:       defaultIterations,
		minH:             defaultMinHeight,
		maxH:             defaultMaxHeight,
		seed:             defaultSeed,
		minHIsDefault:    false,
		peakProbability:  0,
		cliffProbability: 0,
		scale:            mgl32.Vec3{1, 1, 1},
		debugMode:        false,
	}
}

// SetWidth updates the width.
func (t *TerrainBuilder) SetWidth(width int) {
	t.width = width
}

// SetLength updates the length.
func (t *TerrainBuilder) SetLength(length int) {
	t.length = length
}

// SetIterations updates the iterations.
func (t *TerrainBuilder) SetIterations(iterations int) {
	t.iterations = iterations
}

// SetMinHeight updates the minH.
func (t *TerrainBuilder) SetMinHeight(height float32) {
	t.minH = height
}

// SetMaxHeight updates the maxH.
func (t *TerrainBuilder) SetMaxHeight(height float32) {
	t.maxH = height
}

// SetSeed updates the seed.
func (t *TerrainBuilder) SetSeed(seed int64) {
	t.seed = seed
}

// RandomSeeds sets up a random seed value.
func (t *TerrainBuilder) RandomSeed() {
	t.seed = time.Now().UnixNano()
}

// SetPeakProbability sets the peakProbability value.
func (t *TerrainBuilder) SetPeakProbability(p int) {
	t.peakProbability = p
}

// SetCliffProbability sets the cliffProbability value.
func (t *TerrainBuilder) SetCliffProbability(p int) {
	t.cliffProbability = p
}

// MinHeightIsDefault sets the minHIsDefault flag.
func (t *TerrainBuilder) MinHeightIsDefault(f bool) {
	t.minHIsDefault = f
}

// SetGlWrapper sets the wrapper.
func (t *TerrainBuilder) SetGlWrapper(w interfaces.GLWrapper) {
	t.wrapper = w
}

// SetScale sets the scale.
func (t *TerrainBuilder) SetScale(s mgl32.Vec3) {
	t.scale = s
}

//SetDebugMode updates the debug flag
func (t *TerrainBuilder) SetDebugMode(v bool) {
	t.debugMode = v
}

// GrassTexture sets the texture to grass.
func (t *TerrainBuilder) GrassTexture() {
	_, filename, _, _ := runtime.Caller(1)
	fileDir := path.Dir(filename)
	t.tex.AddTexture(fileDir+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", t.wrapper)
	t.tex.AddTexture(fileDir+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", t.wrapper)
}

// It sets the default values to the map. By default, it is 0, but if we set the
// minHIsDefault flag then the minH value is used as default.
func (t *TerrainBuilder) initHeightMap() {
	defaultHeight := float32(0.0)
	if t.minHIsDefault {
		defaultHeight = t.minH
	}
	t.heightMap = make([][]float32, t.length+1)
	for l := 0; l <= t.length; l++ {
		t.heightMap[l] = make([]float32, t.width+1)
	}
	for l := 0; l <= t.length; l++ {
		for w := 0; w <= t.width; w++ {
			t.heightMap[l][w] = defaultHeight
		}
	}
}

// buildHeightMap sets the final values of the height map.
func (t *TerrainBuilder) buildHeightMap() {
	iterationStep := (t.maxH - t.minH) / float32(t.iterations)
	if t.debugMode {
		fmt.Printf("Setup random seed to '%d'.\n", t.seed)
	}
	rand.Seed(t.seed)
	defaultHeight := float32(0.0)
	if t.minHIsDefault {
		defaultHeight = t.minH
	}
	if t.debugMode {
		fmt.Printf("TerrainBuilder.buildHeightMap defaultHeight '%f'.\n", defaultHeight)
		fmt.Printf("TerrainBuilder.buildHeightMap iterations '%d'.\n", t.iterations)
		fmt.Printf("TerrainBuilder.buildHeightMap peakProbability '%d'.\n", t.peakProbability)
	}
	for i := 0; i < t.iterations; i++ {
		height := t.minH + float32(i)*iterationStep
		for l := 0; l <= t.length; l++ {
			for w := 0; w <= t.width; w++ {
				if t.heightMap[l][w] != defaultHeight {
					continue
				}
				random := rand.Intn(100)
				if t.debugMode {
					fmt.Printf("TerrainBuilder.buildHeightMap current random: '%d'.\n", random)
				}
				if t.adjacentElevation(w, l, height-iterationStep) || random < t.peakProbability {
					t.heightMap[l][w] = height
				}
			}
		}
	}
}
func (t *TerrainBuilder) adjacentElevation(cW, cL int, elevation float32) bool {
	for l := max(0, cL-1); l <= min(t.length-1, cL+1); l++ {
		for w := max(0, cW-1); w <= min(t.width-1, cW+1); w++ {
			if t.heightMap[l][w] == elevation {
				return rand.Intn(100) > t.cliffProbability
			}
		}
	}
	return false
}
func (t *TerrainBuilder) vertices() []vertex.Vertex {
	textureCoords := [4]mgl32.Vec2{
		{0.0, 1.0},
		{1.0, 1.0},
		{1.0, 0.0},
		{0.0, 0.0},
	}
	var vertices vertex.Vertices
	for l := 0; l <= t.length; l++ {
		for w := 0; w <= t.width; w++ {
			texIndex := (w % 2) + (l%2)*2
			vertices = append(vertices, vertex.Vertex{
				Position:  mgl32.Vec3{-float32(t.width)/2.0 + float32(w), t.heightMap[l][w], -float32(t.length)/2.0 + float32(l)},
				Normal:    mgl32.Vec3{0, 1, 0},
				TexCoords: textureCoords[texIndex],
			})
		}
	}

	return vertices
}
func (t *TerrainBuilder) indices() []uint32 {
	var indices []uint32
	for w := 0; w < t.width; w++ {
		for l := 0; l < t.length; l++ {
			i0 := uint32(w*(t.length+1) + l)
			i1 := uint32(1) + i0
			i2 := uint32(t.length+1) + i0
			i3 := uint32(1) + i2
			indices = append(indices, i0)
			indices = append(indices, i1)
			indices = append(indices, i2)

			indices = append(indices, i2)
			indices = append(indices, i1)
			indices = append(indices, i3)
		}
	}
	return indices
}

// Build returns a Terrain
func (t *TerrainBuilder) Build() *Terrain {
	t.initHeightMap()
	if t.debugMode {
		fmt.Printf("TerrainBuilder.heightMap after init:\n'%v'\n", t.heightMap)
	}
	t.buildHeightMap()
	if t.debugMode {
		fmt.Printf("TerrainBuilder.heightMap after build:\n'%v'\n", t.heightMap)
	}
	v := t.vertices()
	i := t.indices()
	terrainMesh := mesh.NewTexturedMesh(v, i, t.tex, t.wrapper)
	terrainMesh.SetScale(t.scale)
	m := newModel()
	m.AddMesh(terrainMesh)
	return &Terrain{Model: *m, heightMap: t.heightMap, width: t.width, length: t.length, debugMode: t.debugMode}
}
