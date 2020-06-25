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

type Liquid struct {
	Model
	heightMap     [][]float32
	width, length int
	debugMode     bool
}

// GetLiquid returns the liquid mesh
func (l *Liquid) GetLiquid() interfaces.Mesh {
	return l.meshes[0]
}

// CollideTestWithSphere is the collision detection function for liquid vs sphere.
func (l *Liquid) CollideTestWithSphere(boundingSphere *coldet.Sphere) bool {
	return false
}

// Update function does nothing.
func (l *Liquid) Update(dt float64) {
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
	liquidTex                 texture.Textures
	scale                     mgl32.Vec3
	debugMode                 bool
	position                  mgl32.Vec3
	liquidAmplitude           float32
	liquidFrequency           float32
	liquidEta                 float32
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
		position:         mgl32.Vec3{0, 0, 0},
		liquidEta:        0.0,
		liquidAmplitude:  0.0,
		liquidFrequency:  0.0,
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

// SetPosition sets the position.
func (t *TerrainBuilder) SetPosition(p mgl32.Vec3) {
	t.position = p
}

//SetDebugMode updates the debug flag
func (t *TerrainBuilder) SetDebugMode(v bool) {
	t.debugMode = v
}

// SetLiquidEta sets the liquidEta.
func (t *TerrainBuilder) SetLiquidEta(e float32) {
	t.liquidEta = e
}

// SetLiquidAmplitude sets the liquidAmplitude.
func (t *TerrainBuilder) SetLiquidAmplitude(a float32) {
	t.liquidAmplitude = a
}

// SetLiquidFrequency sets the liquidFrequency.
func (t *TerrainBuilder) SetLiquidFrequency(f float32) {
	t.liquidFrequency = f
}

// SurfaceTextureGrass sets the surface texture to grass.
func (t *TerrainBuilder) SurfaceTextureGrass() {
	_, filename, _, _ := runtime.Caller(1)
	fileDir := path.Dir(filename)
	t.tex.AddTexture(fileDir+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", t.wrapper)
	t.tex.AddTexture(fileDir+"/assets/grass.jpg", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", t.wrapper)
}

// LiquidTextureWater sets the liquid texture to water.
func (t *TerrainBuilder) LiquidTextureWater() {
	_, filename, _, _ := runtime.Caller(1)
	fileDir := path.Dir(filename)
	t.liquidTex.AddTexture(fileDir+"/assets/water.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.diffuse", t.wrapper)
	t.liquidTex.AddTexture(fileDir+"/assets/water.png", glwrapper.CLAMP_TO_EDGE, glwrapper.CLAMP_TO_EDGE, glwrapper.LINEAR, glwrapper.LINEAR, "material.specular", t.wrapper)
}

// It sets the default values to the map. By default, it is 0, but if we set the
// minHIsDefault flag then the minH value is used as default.
func (t *TerrainBuilder) initHeightMap(width, length int, defaultHeight float32) [][]float32 {
	heightMap := make([][]float32, length+1)
	for l := 0; l <= length; l++ {
		heightMap[l] = make([]float32, width+1)
	}
	for l := 0; l <= length; l++ {
		for w := 0; w <= width; w++ {
			heightMap[l][w] = defaultHeight
		}
	}
	return heightMap
}

// buildTerrainHeightMap sets the final values of the height map.
func (t *TerrainBuilder) buildTerrainHeightMap() {
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
		fmt.Printf("TerrainBuilder.buildTerrainHeightMap defaultHeight '%f'.\n", defaultHeight)
		fmt.Printf("TerrainBuilder.buildTerrainHeightMap iterations '%d'.\n", t.iterations)
		fmt.Printf("TerrainBuilder.buildTerrainHeightMap peakProbability '%d'.\n", t.peakProbability)
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
					fmt.Printf("TerrainBuilder.buildTerrainHeightMap current random: '%d'.\n", random)
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
func (t *TerrainBuilder) vertices(width, length, textureWidth, textureHeight int, heightMap [][]float32) []vertex.Vertex {
	textureStepX := float32(1.0) / float32(textureWidth)
	textureStepY := float32(1.0) / float32(textureHeight)
	var vertices vertex.Vertices
	for l := 0; l <= length; l++ {
		for w := 0; w <= width; w++ {
			textureModX := w % (textureWidth + 1)
			textureModY := l % (textureHeight + 1)
			var iL, iW int
			if l == length {
				iL = l
			} else {
				iL = l + 1
			}
			if w == width {
				iW = w
			} else {
				iW = w + 1
			}
			currentPos := mgl32.Vec3{
				-float32(width)/2.0 + float32(w),
				heightMap[l][w],
				-float32(length)/2.0 + float32(l)}
			nextPosX := mgl32.Vec3{
				-float32(width)/2.0 + float32(iW),
				heightMap[l][iW],
				-float32(length)/2.0 + float32(l)}
			nextPosY := mgl32.Vec3{
				-float32(width)/2.0 + float32(w),
				heightMap[iL][w],
				-float32(length)/2.0 + float32(iL)}
			normal := nextPosX.Sub(currentPos).Cross(nextPosY.Sub(currentPos)).Normalize()
			vertices = append(vertices, vertex.Vertex{
				Position:  currentPos,
				Normal:    normal,
				TexCoords: mgl32.Vec2{float32(textureModX) * textureStepX, float32(textureModY) * textureStepY},
			})
		}
	}

	return vertices
}
func (t *TerrainBuilder) indices(width, length int) []uint32 {
	var indices []uint32
	for w := 0; w < width; w++ {
		for l := 0; l < length; l++ {
			i0 := uint32(w*(length+1) + l)
			i1 := uint32(1) + i0
			i2 := uint32(length+1) + i0
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
	return t.buildTerrain()
}

// BuildWithLiquid returns a Terrain and a Liquid that is generated for the terrain.
func (t *TerrainBuilder) BuildWithLiquid() (*Terrain, *Liquid) {
	return t.buildTerrain(), t.buildLiquid()
}
func (t *TerrainBuilder) buildTerrain() *Terrain {
	defaultHeight := float32(0.0)
	if t.minHIsDefault {
		defaultHeight = t.minH
	}
	t.heightMap = t.initHeightMap(t.width, t.length, defaultHeight)
	if t.debugMode {
		fmt.Printf("TerrainBuilder.buildTerrain.heightMap after init:\n'%v'\n", t.heightMap)
	}
	t.buildTerrainHeightMap()
	if t.debugMode {
		fmt.Printf("TerrainBuilder.heightMap after build:\n'%v'\n", t.heightMap)
	}
	v := t.vertices(t.width, t.length, 1, 1, t.heightMap)
	i := t.indices(t.width, t.length)
	terrainMesh := mesh.NewTexturedMesh(v, i, t.tex, t.wrapper)
	terrainMesh.SetScale(t.scale)
	terrainMesh.SetPosition(t.position)
	m := newModel()
	m.AddMesh(terrainMesh)
	return &Terrain{Model: *m, heightMap: t.heightMap, width: t.width, length: t.length, debugMode: t.debugMode}
}

func (t *TerrainBuilder) buildLiquid() *Liquid {
	waterTopLevel := float32(0.0)
	waterMultiplier := 10
	waterWidth := t.width * int(t.scale.X()) * waterMultiplier
	waterLength := t.length * int(t.scale.Z()) * waterMultiplier
	waterHeightMap := t.initHeightMap(waterWidth, waterLength, waterTopLevel)
	if t.debugMode {
		fmt.Printf("TerrainBuilder.buildLiquid.heightMap after init:\n'%v'\n", waterHeightMap)
	}
	v := t.vertices(waterWidth, waterLength, waterWidth, waterLength, waterHeightMap)
	i := t.indices(waterWidth, waterLength)
	liquidMesh := mesh.NewTexturedMesh(v, i, t.liquidTex, t.wrapper)
	scaleValue := float32(1.0) / float32(waterMultiplier)
	liquidMesh.SetScale(mgl32.Vec3{scaleValue, 1.0, scaleValue})
	liquidMesh.SetPosition(mgl32.Vec3{t.position.X(), t.position.Y() + waterTopLevel, t.position.Z()})

	m := newModel()
	m.AddMesh(liquidMesh)
	m.SetTransparent(true)
	m.SetUniformFloat("Eta", t.liquidEta)
	m.SetUniformFloat("amplitude", t.liquidAmplitude)
	m.SetUniformFloat("frequency", t.liquidFrequency)

	return &Liquid{
		Model:     *m,
		heightMap: waterHeightMap,
		width:     waterWidth,
		length:    waterLength,
		debugMode: t.debugMode,
	}
}
