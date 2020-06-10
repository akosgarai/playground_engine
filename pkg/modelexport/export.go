package modelexport

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/mesh"
	trans "github.com/akosgarai/playground_engine/pkg/transformations"

	"github.com/go-gl/mathgl/mgl32"
)

type Mtl struct {
	// Ambient color [0-1] 'Ka' prefix
	Ka [3]float32
	// Diffuse color [0-1] 'Kd' pefix
	Kd [3]float32
	// Specular color [0-1] 'Ks' prefix
	Ks [3]float32
	// Specular exponent [0-1000] 'Ns' prefix
	Ns float32
	// Transparency (Tr = 1 - d) [0-1] 'd' and 'Tr' prefix
	D  float32
	Tr float32
	// Optical density aka index of refraction [0.001-1000] 'Ni' prefix
	Ni float32
	// Illumination model enum [1-10] 'illum' prefix more details : https://en.wikipedia.org/wiki/Wavefront_.obj_file
	Illum float32
	// Ambient texture map (filename) 'map_Ka' prefix
	MapKa string
	// Diffuse texture map (filename) 'map_Kd' prefix
	MapKd string
	// Specular texture map (filename) 'map_Ks' prefix
	MapKs string
	// Specular highlight (filename) 'map_Ns' prefix
	MapNs string
	// Alpha texture map (filename) ? 'map_d' prefix
	MapD string
	// Bump map aka mapBump (filename) 'bump' prefix or 'map_bump'
	Bump string
	// Displacement map (filename) 'disp' prefix
	Disp string
	// Stencil decal texture (filename) 'decal' prefix
	Decal string
	// The name of the material 'newmtl' prefix
	Name string
}

type Obj struct {
	// the name of the object 'o' perfix
	Name string
	// Vertices array 'v' prefix
	V [][3]float32
	// Normal array 'vn' prefix
	Normal [][3]float32
	// Texture coords array 'vt' prefix
	TexCoord [][2]float32
	// material name identifier 'usemtl' prefix
	MaterialName string
	// Indices - string array - it has to be previously constructed.
	Indices []string
	// should be printed the indexes as faces (f) or point (p)
	HasFaces bool
}

type Export struct {
	meshes []interfaces.Mesh
	// The directory path. Files will be written here.
	directory        string
	materials        []Mtl
	objects          []Obj
	positionMaxIndex int
	normalMaxIndex   int
	tcMaxIndex       int
}

// New gets the meshes as input and returns the Export populated with the given meshes.
func New(meshes []interfaces.Mesh) *Export {
	return &Export{
		meshes:           meshes,
		positionMaxIndex: 0,
		normalMaxIndex:   0,
		tcMaxIndex:       0,
	}
}

// Export gets a filepath as input. The files will be written into this directory.
func (e *Export) Export(path string) error {
	// check that the directory exists. if not, return error
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	e.directory = path
	for _, m := range e.meshes {
		switch m.(type) {
		case *mesh.ColorMesh:
			e.processColorMesh(m.(*mesh.ColorMesh))
			break
		case *mesh.MaterialMesh:
			e.processMaterialMesh(m.(*mesh.MaterialMesh))
			break
		case *mesh.TexturedMesh:
			e.processTextureMesh(m.(*mesh.TexturedMesh))
			break
		case *mesh.TexturedColoredMesh:
			e.processTexturedColorMesh(m.(*mesh.TexturedColoredMesh))
			break
		case *mesh.TexturedMaterialMesh:
			e.processTexturedMaterialMesh(m.(*mesh.TexturedMaterialMesh))
			break
		case *mesh.PointMesh:
			e.processPointMesh(m.(*mesh.PointMesh))
			break
		default:
			continue
			break
		}
	}
	objectFileName := "object.obj"
	materialFileName := "material.mat"
	if len(e.materials) > 0 {
		f, err := os.Create(filepath.Join(path, materialFileName))
		if err != nil {
			return err
		}
		content := e.materialExport()
		_, err = f.WriteString(content)
		f.Close()
		if err != nil {
			return err
		}
	}
	f, err := os.Create(filepath.Join(path, objectFileName))
	if err != nil {
		return err
	}
	content := e.objectExport()
	_, err = f.WriteString(content)
	f.Close()
	if err != nil {
		return err
	}
	return nil
}

// It transforms the color to material, and saves it as material mesh, but without normal vectors.
func (e *Export) processColorMesh(m *mesh.ColorMesh) {
	var mtl Mtl
	avgColor := mgl32.Vec3{0, 0, 0}
	for i := 0; i < len(m.Color); i++ {
		avgColor = avgColor.Add(m.Color[i])
	}
	avgColor = avgColor.Mul(1.0 / float32(len(m.Color)))
	mtl.Name = fmt.Sprintf("Color_Material_%d", len(e.materials))
	ka := avgColor
	mtl.Ka = [3]float32{ka.X(), ka.Y(), ka.Z()}
	kd := avgColor
	mtl.Kd = [3]float32{kd.X(), kd.Y(), kd.Z()}
	ks := mgl32.Vec3{1, 1, 1}
	mtl.Ks = [3]float32{ks.X(), ks.Y(), ks.Z()}
	mtl.Ns = float32(32)
	e.materials = append(e.materials, mtl)

	objIndexPositionMap := make(map[mgl32.Vec3]int)

	for _, vert := range m.Vertices {
		if _, ok := objIndexPositionMap[vert.Position]; !ok {
			mapLen := len(objIndexPositionMap)
			objIndexPositionMap[vert.Position] = mapLen
		}
	}
	var obj Obj
	obj.HasFaces = true
	obj.Name = fmt.Sprintf("Color_Material_Object_%d", len(e.objects))
	obj.MaterialName = mtl.Name
	for _, indexValue := range m.Indices {
		index := fmt.Sprintf("%d", e.positionMaxIndex+objIndexPositionMap[m.Vertices[indexValue].Position]+1)
		obj.Indices = append(obj.Indices, index)
	}
	e.positionMaxIndex = e.positionMaxIndex + len(obj.Indices)
	InverseMap := make(map[int][3]float32)
	for origPos, val := range objIndexPositionMap {
		trMatrix := m.ModelTransformation()
		pos := mgl32.TransformCoordinate(origPos, trMatrix)
		InverseMap[val] = [3]float32{pos.X(), pos.Y(), pos.Z()}
	}
	for i := 0; i < len(InverseMap); i++ {
		obj.V = append(obj.V, InverseMap[i])
	}
	e.objects = append(e.objects, obj)
}
func (e *Export) processMaterialMesh(m *mesh.MaterialMesh) {
	var mtl Mtl
	mtl.Name = fmt.Sprintf("Material_%d", len(e.materials))
	ka := m.Material.GetAmbient()
	mtl.Ka = [3]float32{ka.X(), ka.Y(), ka.Z()}
	kd := m.Material.GetDiffuse()
	mtl.Kd = [3]float32{kd.X(), kd.Y(), kd.Z()}
	ks := m.Material.GetSpecular()
	mtl.Ks = [3]float32{ks.X(), ks.Y(), ks.Z()}
	mtl.Ns = m.Material.GetShininess()
	e.materials = append(e.materials, mtl)

	objIndexPositionMap := make(map[mgl32.Vec3]int)
	objIndexNormalMap := make(map[mgl32.Vec3]int)
	for _, vert := range m.Vertices {
		if _, ok := objIndexPositionMap[vert.Position]; !ok {
			mapLen := len(objIndexPositionMap)
			objIndexPositionMap[vert.Position] = mapLen
		}
		if _, ok := objIndexNormalMap[vert.Normal]; !ok {
			mapLen := len(objIndexNormalMap)
			objIndexNormalMap[vert.Normal] = mapLen
		}
	}
	var obj Obj
	obj.HasFaces = true
	obj.Name = fmt.Sprintf("Material_Object_%d", len(e.objects))
	obj.MaterialName = mtl.Name
	for _, indexValue := range m.Indices {
		index := fmt.Sprintf("%d//%d", e.positionMaxIndex+objIndexPositionMap[m.Vertices[indexValue].Position]+1, e.normalMaxIndex+objIndexNormalMap[m.Vertices[indexValue].Normal]+1)
		obj.Indices = append(obj.Indices, index)
	}
	e.positionMaxIndex = e.positionMaxIndex + len(objIndexPositionMap)
	e.normalMaxIndex = e.normalMaxIndex + len(objIndexNormalMap)
	orderedPos := make(map[int][3]float32)
	for origPos, val := range objIndexPositionMap {
		trMatrix := m.ModelTransformation()
		pos := mgl32.TransformCoordinate(origPos, trMatrix)
		orderedPos[val] = [3]float32{pos.X(), pos.Y(), pos.Z()}
	}
	for i := 0; i < len(orderedPos); i++ {
		obj.V = append(obj.V, orderedPos[i])
	}
	orderedNorm := make(map[int][3]float32)
	for origNormal, val := range objIndexNormalMap {
		trMatrix := m.RotationTransformation()
		normal := mgl32.TransformCoordinate(origNormal, trMatrix)
		orderedNorm[val] = [3]float32{normal.X(), normal.Y(), normal.Z()}
	}
	for i := 0; i < len(orderedNorm); i++ {
		obj.Normal = append(obj.Normal, orderedNorm[i])
	}
	e.objects = append(e.objects, obj)
}
func (e *Export) processTextureMesh(m *mesh.TexturedMesh) {
	var mtl Mtl
	mtl.Name = fmt.Sprintf("Texture_Material_%d", len(e.materials))
	mtl.Ka = [3]float32{1, 1, 1}
	mtl.Kd = [3]float32{1, 1, 1}
	mtl.Ks = [3]float32{1, 1, 1}
	mtl.Ns = float32(32)
	if len(m.Textures) > 0 {
		for _, tex := range m.Textures {
			if strings.Contains(tex.UniformName, "diffuse") {
				mtl.MapKd = tex.FilePath
				mtl.MapKa = tex.FilePath
			} else if strings.Contains(tex.UniformName, "specular") {
				mtl.MapKs = tex.FilePath
			} else if strings.Contains(tex.UniformName, "ambient") {
				mtl.MapKa = tex.FilePath
			}
		}
		if mtl.MapKa == "" && mtl.MapKs == "" {
			mtl.MapKd = m.Textures[0].FilePath
			mtl.MapKa = m.Textures[0].FilePath
			mtl.MapKs = m.Textures[0].FilePath
		} else if mtl.MapKs == "" {
			mtl.MapKs = mtl.MapKa
		} else if mtl.MapKa == "" {
			mtl.MapKd = mtl.MapKs
			mtl.MapKa = mtl.MapKs
		}
	}
	e.materials = append(e.materials, mtl)
	objIndexPositionMap := make(map[mgl32.Vec3]int)
	objIndexNormalMap := make(map[mgl32.Vec3]int)
	objIndexTexCoordMap := make(map[mgl32.Vec2]int)
	for _, vert := range m.Vertices {
		if _, ok := objIndexPositionMap[vert.Position]; !ok {
			mapLen := len(objIndexPositionMap)
			objIndexPositionMap[vert.Position] = mapLen
		}
		if _, ok := objIndexNormalMap[vert.Normal]; !ok {
			mapLen := len(objIndexNormalMap)
			objIndexNormalMap[vert.Normal] = mapLen
		}
		if _, ok := objIndexTexCoordMap[vert.TexCoords]; !ok {
			mapLen := len(objIndexTexCoordMap)
			objIndexTexCoordMap[vert.TexCoords] = mapLen
		}
	}
	var obj Obj
	obj.HasFaces = true
	obj.Name = fmt.Sprintf("Material_Object_%d", len(e.objects))
	obj.MaterialName = mtl.Name
	for _, indexValue := range m.Indices {
		index := fmt.Sprintf("%d/%d/%d", e.positionMaxIndex+objIndexPositionMap[m.Vertices[indexValue].Position]+1, e.tcMaxIndex+objIndexTexCoordMap[m.Vertices[indexValue].TexCoords]+1, e.normalMaxIndex+objIndexNormalMap[m.Vertices[indexValue].Normal]+1)
		obj.Indices = append(obj.Indices, index)
	}
	e.positionMaxIndex = e.positionMaxIndex + len(objIndexPositionMap)
	e.normalMaxIndex = e.normalMaxIndex + len(objIndexNormalMap)
	e.tcMaxIndex = e.tcMaxIndex + len(objIndexTexCoordMap)
	orderedPos := make(map[int][3]float32)
	for origPos, val := range objIndexPositionMap {
		trMatrix := m.ModelTransformation()
		pos := mgl32.TransformCoordinate(origPos, trMatrix)
		orderedPos[val] = [3]float32{pos.X(), pos.Y(), pos.Z()}
	}
	for i := 0; i < len(orderedPos); i++ {
		obj.V = append(obj.V, orderedPos[i])
	}
	orderedNorm := make(map[int][3]float32)
	for origNormal, val := range objIndexNormalMap {
		trMatrix := m.RotationTransformation()
		normal := mgl32.TransformCoordinate(origNormal, trMatrix)
		orderedNorm[val] = [3]float32{normal.X(), normal.Y(), normal.Z()}
	}
	for i := 0; i < len(orderedNorm); i++ {
		obj.Normal = append(obj.Normal, orderedNorm[i])
	}
	orderedTexCoord := make(map[int][2]float32)
	for tc, val := range objIndexTexCoordMap {
		orderedTexCoord[val] = [2]float32{tc.X(), tc.Y()}
	}
	for i := 0; i < len(orderedTexCoord); i++ {
		obj.TexCoord = append(obj.TexCoord, orderedTexCoord[i])
	}
	e.objects = append(e.objects, obj)
}
func (e *Export) processTexturedMaterialMesh(m *mesh.TexturedMaterialMesh) {
	var mtl Mtl
	mtl.Name = fmt.Sprintf("Textured_Color_Material_%d", len(e.materials))
	ka := m.Material.GetAmbient()
	mtl.Ka = [3]float32{ka.X(), ka.Y(), ka.Z()}
	kd := m.Material.GetDiffuse()
	mtl.Kd = [3]float32{kd.X(), kd.Y(), kd.Z()}
	ks := m.Material.GetSpecular()
	mtl.Ks = [3]float32{ks.X(), ks.Y(), ks.Z()}
	mtl.Ns = m.Material.GetShininess()
	if len(m.Textures) > 0 {
		for _, tex := range m.Textures {
			if strings.Contains(tex.UniformName, "diffuse") {
				mtl.MapKd = tex.FilePath
				mtl.MapKa = tex.FilePath
			} else if strings.Contains(tex.UniformName, "specular") {
				mtl.MapKs = tex.FilePath
			} else if strings.Contains(tex.UniformName, "ambient") {
				mtl.MapKa = tex.FilePath
			}
		}
		if mtl.MapKa == "" && mtl.MapKs == "" {
			mtl.MapKd = m.Textures[0].FilePath
			mtl.MapKa = m.Textures[0].FilePath
			mtl.MapKs = m.Textures[0].FilePath
		} else if mtl.MapKs == "" {
			mtl.MapKs = mtl.MapKa
		} else if mtl.MapKs == "" {
			mtl.MapKd = mtl.MapKs
			mtl.MapKa = mtl.MapKs
		}
	}
	e.materials = append(e.materials, mtl)

	objIndexPositionMap := make(map[mgl32.Vec3]int)
	objIndexTexCoordMap := make(map[mgl32.Vec2]int)
	for _, vert := range m.Vertices {
		if _, ok := objIndexPositionMap[vert.Position]; !ok {
			mapLen := len(objIndexPositionMap)
			objIndexPositionMap[vert.Position] = mapLen
		}
		if _, ok := objIndexTexCoordMap[vert.TexCoords]; !ok {
			mapLen := len(objIndexTexCoordMap)
			objIndexTexCoordMap[vert.TexCoords] = mapLen
		}
	}
	var obj Obj
	obj.HasFaces = true
	obj.Name = fmt.Sprintf("Texture_Color_Material_Object_%d", len(e.objects))
	obj.MaterialName = mtl.Name
	for _, indexValue := range m.Indices {
		index := fmt.Sprintf("%d/%d", e.positionMaxIndex+objIndexPositionMap[m.Vertices[indexValue].Position]+1, e.tcMaxIndex+objIndexTexCoordMap[m.Vertices[indexValue].TexCoords]+1)
		obj.Indices = append(obj.Indices, index)
	}
	e.positionMaxIndex = e.positionMaxIndex + len(objIndexPositionMap)
	e.tcMaxIndex = e.tcMaxIndex + len(objIndexTexCoordMap)
	orderedPos := make(map[int][3]float32)
	for origPos, val := range objIndexPositionMap {
		trMatrix := m.ModelTransformation()
		pos := mgl32.TransformCoordinate(origPos, trMatrix)
		orderedPos[val] = [3]float32{pos.X(), pos.Y(), pos.Z()}
	}
	for i := 0; i < len(orderedPos); i++ {
		obj.V = append(obj.V, orderedPos[i])
	}
	orderedTexCoord := make(map[int][2]float32)
	for tc, val := range objIndexTexCoordMap {
		orderedTexCoord[val] = [2]float32{tc.X(), tc.Y()}
	}
	for i := 0; i < len(orderedTexCoord); i++ {
		obj.TexCoord = append(obj.TexCoord, orderedTexCoord[i])
	}
	e.objects = append(e.objects, obj)
}
func (e *Export) processTexturedColorMesh(m *mesh.TexturedColoredMesh) {
	var mtl Mtl
	avgColor := mgl32.Vec3{0, 0, 0}
	for i := 0; i < len(m.Color); i++ {
		avgColor = avgColor.Add(m.Color[i])
	}
	avgColor = avgColor.Mul(1.0 / float32(len(m.Color)))
	mtl.Name = fmt.Sprintf("Textured_Color_Material_%d", len(e.materials))
	ka := avgColor
	mtl.Ka = [3]float32{ka.X(), ka.Y(), ka.Z()}
	kd := avgColor
	mtl.Kd = [3]float32{kd.X(), kd.Y(), kd.Z()}
	ks := mgl32.Vec3{1, 1, 1}
	mtl.Ks = [3]float32{ks.X(), ks.Y(), ks.Z()}
	mtl.Ns = float32(32)
	if len(m.Textures) > 0 {
		for _, tex := range m.Textures {
			if strings.Contains(tex.UniformName, "diffuse") {
				mtl.MapKd = tex.FilePath
				mtl.MapKa = tex.FilePath
			} else if strings.Contains(tex.UniformName, "specular") {
				mtl.MapKs = tex.FilePath
			} else if strings.Contains(tex.UniformName, "ambient") {
				mtl.MapKa = tex.FilePath
			}
		}
		if mtl.MapKa == "" && mtl.MapKs == "" {
			mtl.MapKd = m.Textures[0].FilePath
			mtl.MapKa = m.Textures[0].FilePath
			mtl.MapKs = m.Textures[0].FilePath
		} else if mtl.MapKs == "" {
			mtl.MapKs = mtl.MapKa
		} else if mtl.MapKs == "" {
			mtl.MapKd = mtl.MapKs
			mtl.MapKa = mtl.MapKs
		}
	}
	e.materials = append(e.materials, mtl)

	objIndexPositionMap := make(map[mgl32.Vec3]int)
	objIndexTexCoordMap := make(map[mgl32.Vec2]int)
	for _, vert := range m.Vertices {
		if _, ok := objIndexPositionMap[vert.Position]; !ok {
			mapLen := len(objIndexPositionMap)
			objIndexPositionMap[vert.Position] = mapLen
		}
		if _, ok := objIndexTexCoordMap[vert.TexCoords]; !ok {
			mapLen := len(objIndexTexCoordMap)
			objIndexTexCoordMap[vert.TexCoords] = mapLen
		}
	}
	var obj Obj
	obj.HasFaces = true
	obj.Name = fmt.Sprintf("Texture_Color_Material_Object_%d", len(e.objects))
	obj.MaterialName = mtl.Name
	for _, indexValue := range m.Indices {
		index := fmt.Sprintf("%d/%d", e.positionMaxIndex+objIndexPositionMap[m.Vertices[indexValue].Position]+1, e.tcMaxIndex+objIndexTexCoordMap[m.Vertices[indexValue].TexCoords]+1)
		obj.Indices = append(obj.Indices, index)
	}
	e.positionMaxIndex = e.positionMaxIndex + len(objIndexPositionMap)
	e.tcMaxIndex = e.tcMaxIndex + len(objIndexTexCoordMap)
	orderedPos := make(map[int][3]float32)
	for origPos, val := range objIndexPositionMap {
		trMatrix := m.ModelTransformation()
		pos := mgl32.TransformCoordinate(origPos, trMatrix)
		orderedPos[val] = [3]float32{pos.X(), pos.Y(), pos.Z()}
	}
	for i := 0; i < len(orderedPos); i++ {
		obj.V = append(obj.V, orderedPos[i])
	}
	orderedTexCoord := make(map[int][2]float32)
	for tc, val := range objIndexTexCoordMap {
		orderedTexCoord[val] = [2]float32{tc.X(), tc.Y()}
	}
	for i := 0; i < len(orderedTexCoord); i++ {
		obj.TexCoord = append(obj.TexCoord, orderedTexCoord[i])
	}
	e.objects = append(e.objects, obj)
}
func (e *Export) processPointMesh(m *mesh.PointMesh) {
	var obj Obj
	obj.HasFaces = false
	obj.Name = fmt.Sprintf("Point_Object_%d", len(e.objects))
	for _, vert := range m.Vertices {
		trMatrix := m.ModelTransformation()
		pos := mgl32.TransformCoordinate(vert.Position, trMatrix)
		obj.V = append(obj.V, [3]float32{pos.X(), pos.Y(), pos.Z()})
	}
	for i := 0; i < len(obj.V); i++ {
		obj.Indices = append(obj.Indices, fmt.Sprintf("%d", e.positionMaxIndex+i+1))
	}
	e.positionMaxIndex = e.positionMaxIndex + len(obj.Indices)
	e.objects = append(e.objects, obj)
}

// This function is responsible for the material processing.
// Create material file, name materials, write it to material file.
// I expect that the Ka, Kd, Ks, Ns variables are always set.
func (e *Export) materialExport() string {
	materialString := ""
	for _, mat := range e.materials {
		materialString += "newmtl " + mat.Name + "\n"
		materialString += "Ka " + trans.Float32ToString(mat.Ka[0]) + " " + trans.Float32ToString(mat.Ka[1]) + " " + trans.Float32ToString(mat.Ka[2]) + "\n"
		materialString += "Kd " + trans.Float32ToString(mat.Kd[0]) + " " + trans.Float32ToString(mat.Kd[1]) + " " + trans.Float32ToString(mat.Kd[2]) + "\n"
		materialString += "Ks " + trans.Float32ToString(mat.Ks[0]) + " " + trans.Float32ToString(mat.Ks[1]) + " " + trans.Float32ToString(mat.Ks[2]) + "\n"
		materialString += "Ns " + trans.Float32ToString(mat.Ns) + "\n"
		var newValue string
		if mat.MapKa != "" {
			newValue = e.copyFile(mat.MapKa)
			materialString += "map_Ka " + newValue + "\n"
		}
		if mat.MapKd != "" {
			newValue = e.copyFile(mat.MapKd)
			materialString += "map_Kd " + newValue + "\n"
		}
		if mat.MapKs != "" {
			newValue = e.copyFile(mat.MapKs)
			materialString += "map_Ks " + newValue + "\n"
		}
		if mat.MapNs != "" {
			newValue = e.copyFile(mat.MapNs)
			materialString += "map_Ns " + newValue + "\n"
		}
		if mat.MapD != "" {
			newValue = e.copyFile(mat.MapD)
			materialString += "map_d " + newValue + "\n"
		}
		if mat.Bump != "" {
			newValue = e.copyFile(mat.Bump)
			materialString += "bump " + newValue + "\n"
			materialString += "map_bump " + newValue + "\n"
		}
		if mat.Disp != "" {
			newValue = e.copyFile(mat.Disp)
			materialString += "disp " + newValue + "\n"
		}
		if mat.Decal != "" {
			newValue = e.copyFile(mat.Decal)
			materialString += "decal " + newValue + "\n"
		}
		materialString += "\n"
	}
	return materialString
}
func (e *Export) copyFile(src string) string {
	_, err := os.Stat(src)
	if err != nil {
		fmt.Printf("Skipping file export due to the missing file. '%s'\n'%s'\n", src, err.Error())
		return ""
	}
	source, err := os.Open(src)
	if err != nil {
		fmt.Printf("Skipping file export due to the file can not be opened. '%s'\n'%s'\n", src, err.Error())
		return ""
	}
	path := strings.Split(src, "/")
	defer source.Close()
	destination, err := os.Create(e.directory + "/" + path[len(path)-1])
	if err != nil {
		fmt.Printf("Skipping file export due to the wrong destination '%s'. '%s'\n'%s'\n", e.directory+"/"+path[len(path)-1], src, err.Error())
		return ""
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		fmt.Printf("Skipping file export due to the wrong copy '%s'. '%s'\n'%s'\n", e.directory+"/"+path[len(path)-1], src, err.Error())
		return ""
	}
	return path[len(path)-1]
}

// This function is responsible for the geometry processing.
// Create object file, write the content.
func (e *Export) objectExport() string {
	objectString := ""
	if len(e.materials) > 0 {
		objectString += "mtllib material.mat\n"
	}
	for _, obj := range e.objects {
		objectString += "o " + obj.Name + "\n"
		for _, vert := range obj.V {
			objectString += "v " + trans.Float32ToString(vert[0]) + " " + trans.Float32ToString(vert[1]) + " " + trans.Float32ToString(vert[2]) + "\n"
		}
		for _, norm := range obj.Normal {
			objectString += "vn " + trans.Float32ToString(norm[0]) + " " + trans.Float32ToString(norm[1]) + " " + trans.Float32ToString(norm[2]) + "\n"
		}
		for _, tc := range obj.TexCoord {
			objectString += "vt " + trans.Float32ToString(tc[0]) + " " + trans.Float32ToString(tc[1]) + "\n"
		}
		if obj.MaterialName != "" {
			objectString += "usemtl " + obj.MaterialName + "\n"
		}
		if len(obj.Indices) > 3 {
			objectString += "s off\n"
		}
		if obj.HasFaces {
			for i := 0; i < len(obj.Indices); i += 3 {
				objectString += "f " + obj.Indices[i] + " " + obj.Indices[i+1] + " " + obj.Indices[i+2] + "\n"
			}
		} else {
			for i := 0; i < len(obj.Indices); i++ {
				objectString += "p " + obj.Indices[i] + "\n"
			}
		}
		objectString += "\n"
	}
	return objectString
}
