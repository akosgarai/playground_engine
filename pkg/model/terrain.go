package model

import (
	"math/rand"
	"path"
	"runtime"
	"time"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/vertex"
	"github.com/akosgarai/playground_engine/pkg/texture"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	defaultTerrainWidth  = 10
	defaultTerrainLength = 10
	defaultIterations    = 5
	defaultMinHeight     = float32(0.0)
	defaultMaxHeight     = float32(3.0)
	defaultSeed          = int64(0)
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
	}
}

// SetWidth updates the width.
func (t *TerrainBuilder) SetWidth(width int) *TerrainBuilder {
	t.width = width
	return t
}

// SetLength updates the length.
func (t *TerrainBuilder) SetLength(length int) *TerrainBuilder {
	t.length = length
	return t
}

// SetIterations updates the iterations.
func (t *TerrainBuilder) SetIterations(iterations int) *TerrainBuilder {
	t.iterations = iterations
	return t
}

// SetMinHeight updates the minH.
func (t *TerrainBuilder) SetMinHeight(height float32) *TerrainBuilder {
	t.minH = height
	return t
}

// SetMaxHeight updates the maxH.
func (t *TerrainBuilder) SetMaxHeight(height float32) *TerrainBuilder {
	t.maxH = height
	return t
}

// SetSeed updates the seed.
func (t *TerrainBuilder) SetSeed(seed int64) *TerrainBuilder {
	t.seed = seed
	return t
}

// RandomSeeds sets up a random seed value.
func (t *TerrainBuilder) RandomSeed() *TerrainBuilder {
	t.seed = time.Now().UnixNano()
	return t
}

// SetPeekProbability sets the peakProbability value.
func (t *TerrainBuilder) SetPeekProbability(p int) *TerrainBuilder {
	t.peakProbability = p
	return t
}

// SetCliffProbability sets the cliffProbability value.
func (t *TerrainBuilder) SetCliffProbability(p int) *TerrainBuilder {
	t.cliffProbability = p
	return t
}

// MinHeightIsDefault sets the minHIsDefault flag.
func (t *TerrainBuilder) MinHeightIsDefault(f bool) *TerrainBuilder {
	t.minHIsDefault = f
	return t
}

// SetGlWrapper sets the wrapper.
func (t *TerrainBuilder) SetGlWrapper(w interfaces.GLWrapper) *TerrainBuilder {
	t.wrapper = w
	return t
}

// GrassTexture sets the texture to grass.
func (t *TerrainBuilder) GrassTexture() *TerrainBuilder {
	_, filename, _, _ := runtime.Caller(1)
	fileDir := path.Dir(filename)
	t.tex.AddTexture(fileDir+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", t.wrapper)
	t.tex.AddTexture(fileDir+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", t.wrapper)
	return t
}

// It sets the default values to the map. By default, it is 0, but if we set the
// minHIsDefault flag then the minH value is used as default.
func (t *TerrainBuilder) initHeightMap() {
	defaultHeight := float32(0.0)
	if t.minHIsDefault {
		defaultHeight = t.minH
	}
	t.heightMap = make([][]float32, t.length)
	for l := 0; l < t.length; l++ {
		t.heightMap[l] = make([]float32, t.width)
	}
	for l := 0; l < t.length; l++ {
		for w := 0; w < t.width; w++ {
			t.heightMap[l][w] = defaultHeight
		}
	}
}

// buildHeightMap sets the final values of the height map.
func (t *TerrainBuilder) buildHeightMap() {
	iterationStep := (t.maxH - t.minH) / float32(t.iterations)
	rand.Seed(t.seed)
	defaultHeight := float32(0.0)
	if t.minHIsDefault {
		defaultHeight = t.minH
	}
	for i := 0; i <= t.iterations; i++ {
		height := t.minH + float32(i)*iterationStep
		for l := 0; l < t.length; l++ {
			for w := 0; w < t.width; w++ {
				if t.heightMap[l][w] != defaultHeight {
					continue
				}
				random := rand.Intn(100)
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
	for l := 0; l < t.length; l++ {
		for w := 0; w < t.width; w++ {
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
	for w := 0; w <= t.width-1; w++ {
		for l := 0; l <= t.length-1; l++ {
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
	t.buildHeightMap()
	v := t.vertices()
	i := t.indices()
	terrainMesh := mesh.NewTexturedMesh(v, i, t.tex, t.wrapper)
	m := newModel()
	m.AddMesh(terrainMesh)
	return &Terrain{Model: *m}
}
