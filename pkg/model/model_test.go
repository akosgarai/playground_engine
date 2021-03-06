package model

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/akosgarai/playground_engine/pkg/mesh"
	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	InvalidFilename            = "not-existing-file.obj"
	ValidFilename              = "testdata/test_cube.obj"
	DefaultFormItemLabel       = "form item label"
	DefaultFormItemDescription = "form item description."
	DefaultMaxWidth            = float32(1.96)
)

var (
	wrapperMock testhelper.GLWrapperMock
	shaderMock  testhelper.ShaderMock
)

func TestNew(t *testing.T) {
	model := New()
	if len(model.meshes) != 0 {
		t.Errorf("Invalid number of meshes. Instead of '0', we have '%d'.", len(model.meshes))
	}
}
func TestAddMesh(t *testing.T) {
	model := New()
	for i := 0; i < 10; i++ {
		msh := mesh.NewPointMesh(wrapperMock)
		if len(model.meshes) != i {
			t.Errorf("Invalid number of meshes before adding. Instead of '%d', we have '%d'.", i, len(model.meshes))
		}
		model.AddMesh(msh)
		if len(model.meshes) != i+1 {
			t.Errorf("Invalid number of meshes after adding. Instead of '%d', we have '%d'.", i+1, len(model.meshes))
		}
	}
}
func TestGetMeshByIndex(t *testing.T) {
	model := New()
	_, err := model.GetMeshByIndex(2)
	if err != emptyMeshesError {
		t.Errorf("Invalid error for empty meshes. Instead of '%#v', it is '%#v'.", emptyMeshesError, err)
	}
	maxIndex := 10
	for i := 0; i < maxIndex; i++ {
		msh := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(msh)
		m, err := model.GetMeshByIndex(i)
		if err != nil {
			t.Errorf("Should be fine, but we have the following error: '%#v'.", err)
		}
		if m != msh {
			t.Error("Invalid mesh.")
		}
	}
	_, err = model.GetMeshByIndex(-1)
	if err == nil {
		t.Error("Negative index should fail.")
	}
	_, err = model.GetMeshByIndex(maxIndex)
	if err == nil {
		t.Error("Big index should fail.")
	}
}
func TestSetTransparent(t *testing.T) {
	model := New()
	if model.transparent != false {
		t.Error("Invalid initial transparent value. It should be false.")
	}
	model.SetTransparent(true)
	if model.transparent != true {
		t.Error("Invalid updated transparent value. It should be true.")
	}
	model.SetTransparent(false)
	if model.transparent != false {
		t.Error("Invalid updated transparent value. It should be false.")
	}
}
func TestIsTransparent(t *testing.T) {
	model := New()
	if model.IsTransparent() != false {
		t.Error("Invalid initial transparent value. It should be false.")
	}
	model.SetTransparent(true)
	if model.IsTransparent() != true {
		t.Error("Invalid updated transparent value. It should be true.")
	}
	model.SetTransparent(false)
	if model.IsTransparent() != false {
		t.Error("Invalid updated transparent value. It should be false.")
	}
}
func TestDraw(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Draw shouldn't have paniced.")
			}
		}()
		model := New()
		model.Draw(shaderMock)
		msh := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(msh)
		model.Draw(shaderMock)
		model.SetUniformFloat("name", float32(0.1))
		model.SetUniformVector("name", mgl32.Vec3{0.0, 1.0, 1.0})
		model.Draw(shaderMock)
	}()
}
func TestSetSpeed(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("SetSpeed shouldn't have paniced.")
			}
		}()
		model := New()
		for i := 0; i < 10; i++ {
			msh := mesh.NewPointMesh(wrapperMock)
			model.AddMesh(msh)
		}
		model.SetSpeed(2)
	}()
}
func TestUpdate(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Update shouldn't have paniced.")
			}
		}()
		delta := 10.0
		model := New()
		model.Update(delta)
		msh := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(msh)
		model.Update(delta)
	}()
}
func TestSetDirection(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("SetDirection shouldn't have paniced.")
			}
		}()
		model := New()
		dir := mgl32.Vec3{1.0, 0.0, 0.0}
		for i := 0; i < 10; i++ {
			msh := mesh.NewPointMesh(wrapperMock)
			model.AddMesh(msh)
		}
		model.SetDirection(dir)
	}()
}
func TestRotateX(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("RotateX shouldn't have paniced.")
			}
		}()
		model := New()
		parent := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(parent)

		for i := 0; i < 10; i++ {
			msh := mesh.NewPointMesh(wrapperMock)
			if i%2 == 1 {
				msh.SetParent(parent)
			}
			model.AddMesh(msh)
		}
		model.RotateX(90)
	}()
}
func TestRotateY(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("RotateY shouldn't have paniced.")
			}
		}()
		model := New()
		parent := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(parent)

		for i := 0; i < 10; i++ {
			msh := mesh.NewPointMesh(wrapperMock)
			if i%2 == 1 {
				msh.SetParent(parent)
			}
			model.AddMesh(msh)
		}
		model.RotateY(90)
	}()
}
func TestRotateZ(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("RotateZ shouldn't have paniced.")
			}
		}()
		model := New()
		parent := mesh.NewPointMesh(wrapperMock)
		model.AddMesh(parent)

		for i := 0; i < 10; i++ {
			msh := mesh.NewPointMesh(wrapperMock)
			if i%2 == 1 {
				msh.SetParent(parent)
			}
			model.AddMesh(msh)
		}
		model.RotateZ(90)
	}()
}
func TestCollideTestWithSphere(t *testing.T) {
	model := New()
	meshWoBo := mesh.NewPointMesh(wrapperMock)
	model.AddMesh(meshWoBo)
	meshSphere := mesh.NewPointMesh(wrapperMock)
	sphereParams := make(map[string]float32)
	sphereParams["radius"] = float32(1.0)
	sphereBo := boundingobject.New("Sphere", sphereParams)
	meshSphere.SetBoundingObject(sphereBo)
	meshSphere.SetPosition(mgl32.Vec3{2.0, 2.0, 2.0})
	model.AddMesh(meshSphere)
	meshCube := mesh.NewPointMesh(wrapperMock)
	cubeParams := make(map[string]float32)
	cubeParams["width"] = float32(1.0)
	cubeParams["height"] = float32(1.0)
	cubeParams["length"] = float32(1.0)
	cubeBo := boundingobject.New("AABB", cubeParams)
	meshCube.SetBoundingObject(cubeBo)
	meshCube.SetPosition(mgl32.Vec3{-3, -3, -3})
	model.AddMesh(meshCube)
	parent := mesh.NewPointMesh(wrapperMock)
	parent.SetBoundingObject(sphereBo)
	parent.SetPosition(mgl32.Vec3{5.0, 5.0, 5.0})
	model.AddMesh(parent)
	child := mesh.NewPointMesh(wrapperMock)
	child.SetBoundingObject(sphereBo)
	child.SetPosition(mgl32.Vec3{5.0, 5.0, 5.0})
	child.SetParent(parent)
	model.AddMesh(child)

	testData := []struct {
		position  [3]float32
		radius    float32
		intersect bool
		msg       string
	}{
		{[3]float32{0, 0, 0}, 0.5, false, "Shouldn't intersect."},
		{[3]float32{2, 1, 2}, 1.0, true, "Should intersect with the sphere."},
		{[3]float32{-2, -2, -2}, 1.5, true, "Should intersect with the cube."},
	}

	for _, tt := range testData {
		base := coldet.NewBoundingSphere(tt.position, tt.radius)
		result := model.CollideTestWithSphere(base)
		if result != tt.intersect {
			t.Errorf("%s expected: '%v', result: '%v'.", tt.msg, tt.intersect, result)
		}
	}
}
func TestExport(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("Export shouldn't have panic.")
			}
		}()
		model := New()
		model.Export("invalid-path")
	}()
}
func TestSetUniformFloat(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("SetUniformFloat shouldn't have panic.")
			}
		}()
		model := New()
		testData := []struct {
			key   string
			value float32
		}{
			{"testName", float32(0.0)},
			{"testName2", float32(1.0)},
			{"testName", float32(1.0)},
		}
		for _, tt := range testData {
			model.SetUniformFloat(tt.key, tt.value)
			if model.uniformFloat[tt.key] != tt.value {
				t.Errorf("Invalud uniform value for key '%s'. Instead of '%f', we have '%f'.", tt.key, tt.value, model.uniformFloat[tt.key])
			}
		}

	}()
}
func TestSetUniformVector(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("SetUniformVector shouldn't have panic.")
			}
		}()
		model := New()
		testData := []struct {
			key   string
			value mgl32.Vec3
		}{
			{"testName", mgl32.Vec3{0.0, 0.0, 0.0}},
			{"testName2", mgl32.Vec3{0.0, 1.0, 0.0}},
			{"testName", mgl32.Vec3{1.0, 0.0, 1.0}},
		}
		for _, tt := range testData {
			model.SetUniformVector(tt.key, tt.value)
			if model.uniformVector[tt.key] != tt.value {
				t.Errorf("Invalud uniform value for key '%s'. Instead of '%v', we have '%v'.", tt.key, tt.value, model.uniformVector[tt.key])
			}
		}

	}()
}
func TestClosestMeshTo(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("SetUniformVector shouldn't have panic.")
			}
		}()
		refPoint := mgl32.Vec3{0, 0, 0}
		model := New()
		// Get closest mesh without mesh
		msh, dist := model.ClosestMeshTo(refPoint)
		if dist != float32(math.MaxFloat32) {
			t.Errorf("Invalid distance. Instead of maxfloat it is '%f', msh: %#v", dist, msh)
		}
		// Get closest mesh with one mesh

		cubeParams := make(map[string]float32)
		cubeParams["width"] = float32(1.0)
		cubeParams["height"] = float32(1.0)
		cubeParams["length"] = float32(1.0)
		cubeBo := boundingobject.New("AABB", cubeParams)
		mshOne := mesh.NewPointMesh(wrapperMock)
		mshOne.SetBoundingObject(cubeBo)
		mshOne.SetPosition(mgl32.Vec3{0, 1, 0})
		model.AddMesh(mshOne)
		msh, dist = model.ClosestMeshTo(refPoint)
		if dist != float32(0.5) {
			t.Errorf("Invalid distance. Instead of 0.5 it is '%f', msh: %#v", dist, msh)
		}
		// Get closest mesh with multiple mesh
		cubeParams = make(map[string]float32)
		cubeParams["width"] = float32(1.0)
		cubeParams["height"] = float32(1.0)
		cubeParams["length"] = float32(1.0)
		cubeBo = boundingobject.New("AABB", cubeParams)

		mshTwo := mesh.NewPointMesh(wrapperMock)
		mshTwo.SetBoundingObject(cubeBo)
		mshTwo.SetPosition(mgl32.Vec3{0.75, 0, 0})
		model.AddMesh(mshTwo)
		msh, dist = model.ClosestMeshTo(refPoint)
		if dist != float32(0.25) {
			t.Errorf("Invalid distance. Instead of 0.5 it is '%f', msh: %#v", dist, msh)
		}
	}()
}
func TestClear(t *testing.T) {
	model := New()
	mshOne := mesh.NewPointMesh(wrapperMock)
	model.AddMesh(mshOne)
	model.Clear()
	if len(model.meshes) != 0 {
		t.Errorf("Invalid number of meshes after clear. Instead of 0, we have %d.", len(model.meshes))
	}
}
func TestTerrainBuilderNewTerrainBuilder(t *testing.T) {
	terr := NewTerrainBuilder()
	if terr.width != defaultTerrainWidth {
		t.Errorf("Invalid default width. Instead of '%d', we have '%d'.\n", defaultTerrainWidth, terr.width)
	}
	if terr.length != defaultTerrainLength {
		t.Errorf("Invalid default length. Instead of '%d', we have '%d'.\n", defaultTerrainLength, terr.length)
	}
	if terr.iterations != defaultIterations {
		t.Errorf("Invalid default iterations. Instead of '%d', we have '%d'.\n", defaultIterations, terr.iterations)
	}
	if terr.minH != defaultMinHeight {
		t.Errorf("Invalid default minH. Instead of '%f', we have '%f'.\n", defaultMinHeight, terr.minH)
	}
	if terr.maxH != defaultMaxHeight {
		t.Errorf("Invalid default maxH. Instead of '%f', we have '%f'.\n", defaultMaxHeight, terr.maxH)
	}
	if terr.seed != defaultSeed {
		t.Errorf("Invalid default seed. Instead of '%d', we have '%d'.\n", defaultSeed, terr.seed)
	}
	if terr.minHIsDefault != false {
		t.Error("Invalid  minHIsDefault. Not false")
	}
	if terr.cliffProbability != 0 {
		t.Errorf("Invalid default cliffProbability. Instead of '%d', we have '%d'.\n", 0, terr.cliffProbability)
	}
	if terr.peakProbability != 0 {
		t.Errorf("Invalid default peakProbability. Instead of '%d', we have '%d'.\n", 0, terr.peakProbability)
	}
}
func TestTerrainBuilderSetWidth(t *testing.T) {
	width := defaultTerrainWidth + 1
	terr := NewTerrainBuilder()
	terr.SetWidth(width)
	if terr.width != width {
		t.Errorf("Invalid width. Instead of '%d', we have '%d'.", width, terr.width)
	}
}
func TestTerrainBuilderSetLength(t *testing.T) {
	length := defaultTerrainLength + 1
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	if terr.length != length {
		t.Errorf("Invalid length. Instead of '%d', we have '%d'.", length, terr.length)
	}
}
func TestTerrainBuilderSetIterations(t *testing.T) {
	iterations := defaultIterations + 1
	terr := NewTerrainBuilder()
	terr.SetIterations(iterations)
	if terr.iterations != iterations {
		t.Errorf("Invalid iterations. Instead of '%d', we have '%d'.", iterations, terr.iterations)
	}
}
func TestTerrainBuilderSetMinHeight(t *testing.T) {
	height := defaultMinHeight + 1.0
	terr := NewTerrainBuilder()
	terr.SetMinHeight(height)
	if terr.minH != height {
		t.Errorf("Invalid minH. Instead of '%f', we have '%f'.", height, terr.minH)
	}
}
func TestTerrainBuilderSetMaxHeight(t *testing.T) {
	height := defaultMaxHeight + 1.0
	terr := NewTerrainBuilder()
	terr.SetMaxHeight(height)
	if terr.maxH != height {
		t.Errorf("Invalid maxH. Instead of '%f', we have '%f'.", height, terr.maxH)
	}
}
func TestTerrainBuilderSetSeed(t *testing.T) {
	seed := defaultSeed + 1
	terr := NewTerrainBuilder()
	terr.SetSeed(seed)
	if terr.seed != seed {
		t.Errorf("Invalid seed. Instead of '%d', we have '%d'.", seed, terr.seed)
	}
}
func TestTerrainBuilderSetScale(t *testing.T) {
	scale := mgl32.Vec3{2, 2, 2}
	terr := NewTerrainBuilder()
	terr.SetScale(scale)
	if terr.scale != scale {
		t.Errorf("Invalid scale. Instead of '%v', we have '%v'.", scale, terr.scale)
	}
}
func TestTerrainBuilderSetPosition(t *testing.T) {
	position := mgl32.Vec3{2, 2, 2}
	terr := NewTerrainBuilder()
	terr.SetPosition(position)
	if terr.position != position {
		t.Errorf("Invalid position. Instead of '%v', we have '%v'.", position, terr.position)
	}
}
func TestTerrainBuilderRandomSeed(t *testing.T) {
	before := time.Now().UnixNano()
	terr := NewTerrainBuilder()
	seed := terr.RandomSeed()
	after := time.Now().UnixNano()
	if seed < before || seed > after || seed != terr.seed {
		t.Errorf("Invalid random seed '%d'. It supposed to be beetween '%d' and '%d'.", terr.seed, before, after)
	}
}
func TestTerrainBuilderSetPeakProbability(t *testing.T) {
	prob := 1
	terr := NewTerrainBuilder()
	terr.SetPeakProbability(prob)
	if terr.peakProbability != prob {
		t.Errorf("Invalid peak prob. Instead of '%d', we have '%d'.", prob, terr.peakProbability)
	}
}
func TestTerrainBuilderSetCliffProbability(t *testing.T) {
	prob := 1
	terr := NewTerrainBuilder()
	terr.SetCliffProbability(prob)
	if terr.cliffProbability != prob {
		t.Errorf("Invalid cliff prob. Instead of '%d', we have '%d'.", prob, terr.cliffProbability)
	}
}
func TestTerrainBuilderMinHeightIsDefault(t *testing.T) {
	terr := NewTerrainBuilder()
	terr.MinHeightIsDefault(true)
	if !terr.minHIsDefault {
		t.Error("Invalid minHIsDefault. It should be true.")
	}
	terr.MinHeightIsDefault(false)
	if terr.minHIsDefault {
		t.Error("Invalid minHIsDefault. It should be false.")
	}
}
func TestTerrainBuilderSetGlWrapper(t *testing.T) {
	var wrapper testhelper.GLWrapperMock
	terr := NewTerrainBuilder()
	terr.SetGlWrapper(wrapper)
	if terr.wrapper != wrapper {
		t.Error("Invalid gl wrapper")
	}
}
func TestTerrainBuilderSurfaceTextureGrass(t *testing.T) {
	var wrapper testhelper.GLWrapperMock
	terr := NewTerrainBuilder()
	terr.SetGlWrapper(wrapper)
	if len(terr.tex) != 0 {
		t.Errorf("Invalid texture length. Instead of '0', we have '%d'.", len(terr.tex))
	}
	terr.SurfaceTextureGrass()
	if len(terr.tex) != 2 {
		t.Errorf("Invalid texture length. Instead of '2', we have '%d'.", len(terr.tex))
	}
}
func TestTerrainBuilderInitHeightMap(t *testing.T) {
	length := 2
	width := 2
	minH := float32(-2)
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	terr.SetWidth(width)
	terr.SetMinHeight(minH)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, 0.0)
	for l := 0; l < terr.length; l++ {
		for w := 0; w < terr.width; w++ {
			if terr.heightMap[l][w] != 0.0 {
				t.Errorf("Invalid heightMap. Instead of '0.0', we have '%f' at (%d,%d)", terr.heightMap[l][w], l, w)
			}
		}
	}
	terr.MinHeightIsDefault(true)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, terr.minH)
	for l := 0; l < terr.length; l++ {
		for w := 0; w < terr.width; w++ {
			if terr.heightMap[l][w] != -2.0 {
				t.Errorf("Invalid heightMap. Instead of '-2.0', we have '%f' at (%d,%d)", terr.heightMap[l][w], l, w)
			}
		}
	}
}
func TestTerrainBuilderBuildHeightMap(t *testing.T) {
	length := 4
	width := 4
	iteration := 10
	minH := float32(-1.0)
	maxH := float32(3.0)
	seed := int64(0)
	peakProb := 5
	cliffProb := 5
	expected := [][]float32{
		{0, 0, 1.8, 1, 0},
		{0, -0.19999999, 0, 1.4000001, 0},
		{-0.19999999, -0.6, -0.19999999, 2.6000001, 0},
		{-0.19999999, -0.19999999, -0.19999999, 0, 0},
		{0, 0, -0.6, 0, -1},
	}
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	terr.SetWidth(width)
	terr.SetMinHeight(minH)
	terr.SetMaxHeight(maxH)
	terr.SetIterations(iteration)
	terr.SetSeed(seed)
	terr.SetPeakProbability(peakProb)
	terr.SetCliffProbability(cliffProb)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, 0.0)
	terr.buildTerrainHeightMap()
	if !reflect.DeepEqual(terr.heightMap, expected) {
		t.Error("Invalid heightmap")
		t.Log(terr.heightMap)
		t.Log(expected)
	}
	terr.MinHeightIsDefault(true)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, terr.minH)
	terr.buildTerrainHeightMap()
}
func TestTerrainBuilderAdjacentElevation(t *testing.T) {
	length := 4
	width := 4
	iteration := 10
	minH := float32(-1.0)
	maxH := float32(3.0)
	seed := int64(0)
	peakProb := 5
	cliffProb := 5
	expected := [][]float32{
		{0, 0, 1.8, 1, 0},
		{0, -0.19999999, 0, 1.4000001, 0},
		{-0.19999999, -0.6, -0.19999999, 2.6000001, 0},
		{-0.19999999, -0.19999999, -0.19999999, 0, 0},
		{0, 0, -0.6, 0, -1},
	}
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	terr.SetWidth(width)
	terr.SetMinHeight(minH)
	terr.SetMaxHeight(maxH)
	terr.SetIterations(iteration)
	terr.SetSeed(seed)
	terr.SetPeakProbability(peakProb)
	terr.SetCliffProbability(cliffProb)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, 0.0)
	terr.buildTerrainHeightMap()
	if !reflect.DeepEqual(terr.heightMap, expected) {
		t.Error("Invalid heightmap")
		t.Log(terr.heightMap)
		t.Log(expected)
	}
}
func TestTerrainBuilderVertices(t *testing.T) {
	length := 4
	width := 4
	iteration := 10
	minH := float32(-1.0)
	maxH := float32(3.0)
	seed := int64(0)
	peakProb := 5
	cliffProb := 5
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	terr.SetWidth(width)
	terr.SetMinHeight(minH)
	terr.SetMaxHeight(maxH)
	terr.SetIterations(iteration)
	terr.SetSeed(seed)
	terr.SetPeakProbability(peakProb)
	terr.SetCliffProbability(cliffProb)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, 0.0)
	terr.buildTerrainHeightMap()
	v := terr.vertices(terr.width, terr.length, 1, 1, terr.heightMap)
	if len(v) != (length+1)*(width+1) {
		t.Errorf("Invalid vertices length. Instead of '%d', we have '%d'.", length*width, len(v))
	}
}
func TestTerrainBuilderIndices(t *testing.T) {
	length := 5
	width := 5
	iteration := 10
	minH := float32(-1.0)
	maxH := float32(3.0)
	seed := int64(0)
	peakProb := 5
	cliffProb := 5
	terr := NewTerrainBuilder()
	terr.SetLength(length)
	terr.SetWidth(width)
	terr.SetMinHeight(minH)
	terr.SetMaxHeight(maxH)
	terr.SetIterations(iteration)
	terr.SetSeed(seed)
	terr.SetPeakProbability(peakProb)
	terr.SetCliffProbability(cliffProb)
	terr.heightMap = terr.initHeightMap(terr.width, terr.length, 0.0)
	terr.buildTerrainHeightMap()
	i := terr.indices(terr.width, terr.length)
	if len(i) != length*width*6 {
		t.Errorf("Invalid indices length. Instead of '%d', we have '%d'.", length*width*6, len(i))
	}
}
func TestTerrainBuilderBuild(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing textures.")
			}
		}()
		length := 5
		width := 5
		iteration := 10
		minH := float32(-1.0)
		maxH := float32(3.0)
		seed := int64(0)
		peakProb := 5
		cliffProb := 5
		tb := NewTerrainBuilder()
		tb.SetLength(length)
		tb.SetWidth(width)
		tb.SetMinHeight(minH)
		tb.SetMaxHeight(maxH)
		tb.SetIterations(iteration)
		tb.SetSeed(seed)
		tb.SetPeakProbability(peakProb)
		tb.SetCliffProbability(cliffProb)
		tb.SetDebugMode(true)
		terr := tb.Build()
		t.Log(terr)
	}()
}
func TestTerrainBuilderBuildTerrain(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing textures.")
			}
		}()
		length := 5
		width := 5
		iteration := 10
		minH := float32(-1.0)
		maxH := float32(3.0)
		seed := int64(0)
		peakProb := 5
		cliffProb := 5
		tb := NewTerrainBuilder()
		tb.SetLength(length)
		tb.SetWidth(width)
		tb.SetMinHeight(minH)
		tb.SetMaxHeight(maxH)
		tb.SetIterations(iteration)
		tb.SetSeed(seed)
		tb.SetPeakProbability(peakProb)
		tb.SetCliffProbability(cliffProb)
		tb.SetDebugMode(true)
		terr := tb.buildTerrain()
		t.Log(terr)
	}()
}
func TestTerrainBuilderBuildWithLiquid(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing textures.")
			}
		}()
		length := 5
		width := 5
		iteration := 10
		minH := float32(-1.0)
		maxH := float32(3.0)
		seed := int64(0)
		peakProb := 5
		cliffProb := 5
		tb := NewTerrainBuilder()
		tb.SetLength(length)
		tb.SetWidth(width)
		tb.SetMinHeight(minH)
		tb.SetMaxHeight(maxH)
		tb.SetIterations(iteration)
		tb.SetSeed(seed)
		tb.SetPeakProbability(peakProb)
		tb.SetCliffProbability(cliffProb)
		tb.SetDebugMode(true)
		terr, water := tb.BuildWithLiquid()
		t.Log(terr, water)
	}()
}
func TestTerrainBuilderBuildLiquid(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Should have panic due to the missing textures.")
			}
		}()
		length := 5
		width := 5
		iteration := 10
		minH := float32(-1.0)
		maxH := float32(3.0)
		seed := int64(0)
		peakProb := 5
		cliffProb := 5
		tb := NewTerrainBuilder()
		tb.SetLength(length)
		tb.SetWidth(width)
		tb.SetMinHeight(minH)
		tb.SetMaxHeight(maxH)
		tb.SetIterations(iteration)
		tb.SetSeed(seed)
		tb.SetPeakProbability(peakProb)
		tb.SetCliffProbability(cliffProb)
		tb.SetDebugMode(true)
		water := tb.buildLiquid()
		t.Log(water)
	}()
}
func TestTerrainHeightAtPos(t *testing.T) {
	tb := NewTerrainBuilder()
	tb.SetScale(mgl32.Vec3{2, 1, 2})
	tb.SetGlWrapper(wrapperMock)
	tb.SurfaceTextureGrass()
	terrain := tb.Build()
	terrain.heightMap = [][]float32{
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	testData := []struct {
		position  mgl32.Vec3
		height    float32
		debugMode bool
		err       error
	}{
		{mgl32.Vec3{15, 0, 15}, -1, false, ErrorNotAboveTheSurface},
		{mgl32.Vec3{11, 0, 10}, -1, false, ErrorNotAboveTheSurface},
		{mgl32.Vec3{10.2, 0, 10}, -1, false, ErrorNotAboveTheSurface},
		{mgl32.Vec3{10.000002, 0, 10}, -1, false, ErrorNotAboveTheSurface},
		{mgl32.Vec3{2, 0, 8}, 0, false, nil},
		{mgl32.Vec3{0, 0, 0}, 0, false, nil},
		{mgl32.Vec3{-10, 0, -10}, 0, false, nil},
		{mgl32.Vec3{-10, 0, -8}, 0, false, nil},
		{mgl32.Vec3{-9, 0, -9}, 0.25, false, nil},
		{mgl32.Vec3{-8, 0, -8.01}, 0.995, false, nil},
		{mgl32.Vec3{-8.01, 0, -8.01}, 0.990025, false, nil},
		{mgl32.Vec3{-8.01, 0, -8}, 0.995, false, nil},
		{mgl32.Vec3{-8, 0, -8}, 1.0, false, nil},
		{mgl32.Vec3{-7.99, 0, -8}, 1.0, false, nil},
		{mgl32.Vec3{-8, 0, -7.99}, 1.0, false, nil},
		{mgl32.Vec3{-7.99, 0, -7.99}, 1.0, false, nil},
		{mgl32.Vec3{-7, 0, -7}, 1.0, false, nil},
		{mgl32.Vec3{-6, 0, -6}, 1.0, false, nil},
		{mgl32.Vec3{-5, 0, -5}, 0.25, true, nil},
		{mgl32.Vec3{-10, 0, -7}, 0.0, false, nil},
	}
	for _, v := range testData {
		terrain.debugMode = v.debugMode
		x, err := terrain.HeightAtPos(v.position)
		if !testhelper.Float32ApproxEqual(x, v.height, 0.00001) || err != v.err {
			t.Errorf("Invalid height results.\n - At position '%v',\n\texpected height: '%f', error: '%v',\n\tgiven height: '%f', error: '%v'.", v.position, v.height, v.err, x, err)
			t.Log(terrain.heightMap)
		}
	}
}
func TestTerrainCollideTestWithSphere(t *testing.T) {
	tb := NewTerrainBuilder()
	tb.SetScale(mgl32.Vec3{2, 1, 2})
	tb.SetGlWrapper(wrapperMock)
	tb.SurfaceTextureGrass()
	terrain := tb.Build()
	terrain.heightMap = [][]float32{
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	testData := []struct {
		position [3]float32
		radius   float32
		collide  bool
	}{
		{[3]float32{15, 0, 15}, 1, false},
		{[3]float32{11, 0, 10}, 1, false},
		{[3]float32{10.2, 0, 10}, 1, false},
		{[3]float32{10.000002, 0, 10}, 1, false},
		{[3]float32{2, 1.1, 8}, 1, false},
		{[3]float32{2, 0.9, 8}, 1, true},
	}
	for _, v := range testData {
		bs := coldet.NewBoundingSphere(v.position, v.radius)
		result := terrain.CollideTestWithSphere(bs)
		if result != v.collide {
			t.Errorf("Invalid collision.\nPosition:\t'%v'\nRadius:\t'%f'\nExpected:\t'%v'", v.position, v.radius, v.collide)
		}
	}
	terrain.GetTerrain().SetPosition(mgl32.Vec3{0, 1, 0})
	testData = []struct {
		position [3]float32
		radius   float32
		collide  bool
	}{
		{[3]float32{15, 0, 15}, 1, false},
		{[3]float32{11, 0, 10}, 1, false},
		{[3]float32{10.2, 0, 10}, 1, false},
		{[3]float32{10.000002, 0, 10}, 1, false},
		{[3]float32{2, 2.1, 8}, 1, false},
		{[3]float32{2, 1.9, 8}, 1, true},
	}
	for _, v := range testData {
		bs := coldet.NewBoundingSphere(v.position, v.radius)
		result := terrain.CollideTestWithSphere(bs)
		if result != v.collide {
			t.Errorf("Invalid collision.\nPosition:\t'%v'\nRadius:\t'%f'\nExpected:\t'%v'", v.position, v.radius, v.collide)
		}
	}
}
func TestLiquidCollideTestWithSphere(t *testing.T) {
	tb := NewTerrainBuilder()
	tb.SetScale(mgl32.Vec3{2, 1, 2})
	tb.SetGlWrapper(wrapperMock)
	tb.SurfaceTextureGrass()
	tb.LiquidTextureWater()
	_, water := tb.BuildWithLiquid()
	testData := []struct {
		position [3]float32
		radius   float32
		collide  bool
	}{
		{[3]float32{15, 0, 15}, 1, false},
		{[3]float32{11, 0, 10}, 1, false},
		{[3]float32{10.2, 0, 10}, 1, false},
		{[3]float32{10.000002, 0, 10}, 1, false},
		{[3]float32{2, 2.1, 8}, 1, false},
		{[3]float32{2, 1.9, 8}, 1, false},
	}
	for _, v := range testData {
		bs := coldet.NewBoundingSphere(v.position, v.radius)
		result := water.CollideTestWithSphere(bs)
		if result != v.collide {
			t.Errorf("Invalid collision.\nPosition:\t'%v'\nRadius:\t'%f'\nExpected:\t'%v'", v.position, v.radius, v.collide)
		}
	}
}
func TestTerrainBuilderSetLiquidEta(t *testing.T) {
	eta := float32(1)
	terr := NewTerrainBuilder()
	terr.SetLiquidEta(eta)
	if terr.liquidEta != eta {
		t.Errorf("Invalid liquid eta. Instead of '%f', we have '%f'.", eta, terr.liquidEta)
	}
}
func TestTerrainBuilderSetLiquidAmplitude(t *testing.T) {
	ampl := float32(1)
	terr := NewTerrainBuilder()
	terr.SetLiquidAmplitude(ampl)
	if terr.liquidAmplitude != ampl {
		t.Errorf("Invalid liquid amplitude. Instead of '%f', we have '%f'.", ampl, terr.liquidAmplitude)
	}
}
func TestTerrainBuilderSetLiquidFrequency(t *testing.T) {
	freq := float32(1)
	terr := NewTerrainBuilder()
	terr.SetLiquidFrequency(freq)
	if terr.liquidFrequency != freq {
		t.Errorf("Invalid liquid frequency. Instead of '%f', we have '%f'.", freq, terr.liquidFrequency)
	}
}
func TestTerrainBuilderSetLiquidWaterLevel(t *testing.T) {
	wl := float32(1)
	terr := NewTerrainBuilder()
	terr.SetLiquidWaterLevel(wl)
	if terr.liquidWaterLevel != wl {
		t.Errorf("Invalid liquid waterlevel. Instead of '%f', we have '%f'.", wl, terr.liquidWaterLevel)
	}
}
func TestTerrainBuilderSetLiquidDetailMultiplier(t *testing.T) {
	dm := 3
	terr := NewTerrainBuilder()
	terr.SetLiquidDetailMultiplier(dm)
	if terr.liquidDetailMultiplier != dm {
		t.Errorf("Invalid liquid detail multiplier. Instead of '%d', we have '%d'.", dm, terr.liquidDetailMultiplier)
	}
}
