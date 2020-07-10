package mesh

import (
	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
	"github.com/akosgarai/playground_engine/pkg/material"
	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"
	"github.com/akosgarai/playground_engine/pkg/primitives/vertex"
	"github.com/akosgarai/playground_engine/pkg/texture"
	"github.com/akosgarai/playground_engine/pkg/transformations"

	"github.com/go-gl/mathgl/mgl32"
)

type Mesh struct {
	Vertices vertex.Vertices
	Indices  []uint32

	vbo uint32
	vao uint32

	// the center position of the mesh. the model transformation is calculated based on this.
	position mgl32.Vec3
	// movement paramteres
	direction mgl32.Vec3
	velocity  float32
	// rotation parameters - YPR - these values are not in radian
	yaw   float32 // rotation on Y axis
	pitch float32 // rotation on X axis
	roll  float32 // rotation on Z axis
	// for scaling - if a want to make other rectangles than unit ones.
	// This vector contains the scale factor for each axis.
	scale mgl32.Vec3
	// For calling gl functions.
	wrapper interfaces.GLWrapper
	// parent-child hierarchy
	parent    interfaces.Mesh
	parentSet bool

	bo                *boundingobject.BoundingObject
	boundingObjectSet bool
}

// SetScale updates the scale of the mesh.
func (m *Mesh) SetScale(s mgl32.Vec3) {
	m.scale = s
}

// SetPosition updates the position of the mesh.
func (m *Mesh) SetPosition(p mgl32.Vec3) {
	m.position = p
}

// SetDirection updates the move direction of the mesh.
func (m *Mesh) SetDirection(p mgl32.Vec3) {
	m.direction = p
}

// SetSpeed updates the velocity of the mesh.
func (m *Mesh) SetSpeed(a float32) {
	m.velocity = a
}

// GetPosition returns the current position of the mesh.
func (m *Mesh) GetPosition() mgl32.Vec3 {
	return m.position
}

// GetDirection returns the current direction of the mesh.
func (m *Mesh) GetDirection() mgl32.Vec3 {
	return m.direction
}

// SetBoundingObject sets the bounding object of the mesh.
func (m *Mesh) SetBoundingObject(bo *boundingobject.BoundingObject) {
	m.boundingObjectSet = true
	m.bo = bo
}

// GetBoundingObject returns the bounding object of the mesh. It calculates the
// transformed values of the params.
func (m *Mesh) GetBoundingObject() *boundingobject.BoundingObject {
	transformedParams := make(map[string]float32)
	boType := m.bo.Type()
	boParams := m.bo.Params()
	if boType == "Sphere" {
		// In this case, we only need to handle the radius, that only needs
		// to be scaled. Because the paren't scale also counts, the scaling
		// tr will applied to a vector that is based on the rad.
		rad := mgl32.Vec3{boParams["radius"], 0, 0}
		transformedRad := mgl32.TransformCoordinate(rad, m.ScaleTransformation())
		transformedParams["radius"] = transformedRad.X()
	} else if boType == "AABB" {
		// Here we need to handle the side lengths. The scale & the rotation
		// counts here.
		sideLengths := mgl32.Vec3{boParams["width"], boParams["height"], boParams["length"]}
		rotatedSideLengths := mgl32.TransformCoordinate(sideLengths, m.RotationTransformation())
		scaledSideLengths := mgl32.TransformCoordinate(rotatedSideLengths, m.ScaleTransformation())
		transformedParams["width"] = transformations.Float32Abs(scaledSideLengths.X())
		transformedParams["height"] = transformations.Float32Abs(scaledSideLengths.Y())
		transformedParams["length"] = transformations.Float32Abs(scaledSideLengths.Z())
	}
	return boundingobject.New(boType, transformedParams)
}

// SetParent sets the given mesh to the parent. It also sets the
// parentSet variable true, to make this state trackable.
func (m *Mesh) SetParent(msh interfaces.Mesh) {
	m.parentSet = true
	m.parent = msh
}

// GetParentTranslationTransformation returns the translation transformation
// of the parent mesh. If parent is not set, then ident. matrix is returned.
func (m *Mesh) GetParentTranslationTransformation() mgl32.Mat4 {
	if m.parentSet {
		return m.parent.TranslationTransformation()
	}
	return mgl32.Ident4()
}

// GetParentRotationTransformation returns the rotation transformation
// of the parent mesh. If parent is not set, then ident. matrix is returned.
func (m *Mesh) GetParentRotationTransformation() mgl32.Mat4 {
	if m.parentSet {
		return m.parent.RotationTransformation()
	}
	return mgl32.Ident4()
}

// GetParentScaleTransformation returns the scale transformation
// of the parent mesh. If parent is not set, then ident. matrix is returned.
func (m *Mesh) GetParentScaleTransformation() mgl32.Mat4 {
	if m.parentSet {
		return m.parent.ScaleTransformation()
	}
	return mgl32.Ident4()
}

// Update calulates the position change. It's input is the delta since the current draw circle.
// The movement is calculated from the direction, velocity and delta.
// motion = motionVector * (delta * velocity)
// new position = current position + motion
func (m *Mesh) Update(dt float64) {
	delta := float32(dt)
	motionVector := m.direction
	if motionVector.Len() > 0 {
		motionVector = motionVector.Normalize().Mul(delta * m.velocity)
	}
	m.position = m.position.Add(motionVector)
}

// ModelTransformation returns the transformation that we can
// use as the model transformation of this mesh.
// The matrix is calculated from the position (translate), the rotation (rotate)
// and from the scale (scale) patameters.
func (m *Mesh) ModelTransformation() mgl32.Mat4 {
	return m.TranslationTransformation().Mul4(
		m.RotationTransformation()).Mul4(
		m.ScaleTransformation())
}

// ScaleTransformation returns the scale part of the model transformation.
func (m *Mesh) ScaleTransformation() mgl32.Mat4 {
	return mgl32.Scale3D(m.scale.X(), m.scale.Y(), m.scale.Z()).Mul4(m.GetParentScaleTransformation())
}

// TranslateTransformation returns the translation part of the model transformation.
func (m *Mesh) TranslationTransformation() mgl32.Mat4 {
	return mgl32.Translate3D(m.position.X(), m.position.Y(), m.position.Z()).Mul4(m.GetParentTranslationTransformation())
}

// RotationTransformation returns the rotation part of the model transformation.
// It is used in the export module, where we have to handle the normal vectors also.
func (m *Mesh) RotationTransformation() mgl32.Mat4 {
	return mgl32.HomogRotate3DY(mgl32.DegToRad(m.yaw)).Mul4(
		mgl32.HomogRotate3DX(mgl32.DegToRad(m.pitch))).Mul4(
		mgl32.HomogRotate3DZ(mgl32.DegToRad(m.roll))).Mul4(m.GetParentRotationTransformation())
}

// RotateY adds the given input angle to the yaw. It also updates the direction vector.
func (m *Mesh) RotateY(angleDeg float32) {
	m.yaw = m.yaw + angleDeg
	m.rotateDirection(angleDeg, mgl32.Vec3{0.0, 1.0, 0.0})
}

// RotateX adds the given input angle to the pitch. It also updates the direction vector.
func (m *Mesh) RotateX(angleDeg float32) {
	m.pitch = m.pitch + angleDeg
	m.rotateDirection(angleDeg, mgl32.Vec3{1.0, 0.0, 0.0})
}

// RotateZ adds the given input angle to the roll. It also updates the direction vector.
func (m *Mesh) RotateZ(angleDeg float32) {
	m.roll = m.roll + angleDeg
	m.rotateDirection(angleDeg, mgl32.Vec3{0.0, 0.0, 1.0})
}

// RotatePosition rotates the position. The transformation matrix is constructed
// as rotation matrix based on the input angle and axis.
func (m *Mesh) RotatePosition(angleDeg float32, axisVector mgl32.Vec3) {
	trMat := mgl32.HomogRotate3D(mgl32.DegToRad(angleDeg), axisVector)
	m.position = mgl32.TransformCoordinate(m.position, trMat)
}
func (m *Mesh) rotateDirection(angleDeg float32, axisVector mgl32.Vec3) {
	trMat := mgl32.HomogRotate3D(mgl32.DegToRad(angleDeg), axisVector)
	m.direction = mgl32.TransformNormal(m.direction, trMat)
}

// IsParentMesh returns true if parent mesh is not set to this mesh.
func (m *Mesh) IsParentMesh() bool {
	return !m.parentSet
}

// GetParent returns the parent mesh of this mesh
func (m *Mesh) GetParent() interfaces.Mesh {
	return m.parent
}

// IsBoundingObjectSet returns true, if the bounding object is set to this mesh.
func (m *Mesh) IsBoundingObjectSet() bool {
	return m.boundingObjectSet
}

type TexturedMesh struct {
	Mesh
	Indices  []uint32
	Textures texture.Textures
	ebo      uint32
}

func (m *TexturedMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()
	m.ebo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Vertices.Get(vertex.POSITION_NORMAL_TEXCOORD))

	m.wrapper.BindBuffer(glwrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	m.wrapper.ElementBufferData(m.Indices)

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(0))
	// setup normals
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*3))
	// setup texture position
	m.wrapper.VertexAttribPointer(2, 2, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*6))

	// close
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. It's input is a shader.
// First it binds the textures with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function). Then it sets up the model uniform,
// and the shininess.
// Then it binds the vertex array and draws the mesh with triangles. Finally it cleans up.
func (m *TexturedMesh) Draw(shader interfaces.Shader) {
	for _, item := range m.Textures {
		item.Bind()
		shader.SetUniform1i(item.UniformName, int32(item.Id-glwrapper.TEXTURE0))
	}
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	shader.SetUniform1f("material.shininess", float32(32))
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawTriangleElements(int32(len(m.Indices)))

	m.Textures.UnBind()
	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}

// NewTexturedMesh gets the vertices, indices, textures, glwrapper as inputs and makes the
// necessary setup for a standing (not moving) textured mesh before returning it.
// The vbo, vao, ebo is also set.
func NewTexturedMesh(v []vertex.Vertex, i []uint32, t texture.Textures, wrapper interfaces.GLWrapper) *TexturedMesh {
	mesh := &TexturedMesh{
		Mesh: Mesh{
			Vertices: v,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			yaw:       0,
			pitch:     0,
			roll:      0,
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,

			boundingObjectSet: false,
		},
		Indices:  i,
		Textures: t,
	}
	mesh.setup()
	return mesh
}

type MaterialMesh struct {
	Mesh
	Indices  []uint32
	Material *material.Material
	ebo      uint32
}

// NewMaterialMesh gets the vertices, indices, material, glwrapper as inputs and makes the
// necessary setup for a standing (not moving) material mesh before returning it.
// The vbo, vao, ebo is also set.
func NewMaterialMesh(v []vertex.Vertex, i []uint32, mat *material.Material, wrapper interfaces.GLWrapper) *MaterialMesh {
	mesh := &MaterialMesh{
		Mesh: Mesh{
			Vertices: v,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			yaw:       0,
			pitch:     0,
			roll:      0,
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,

			boundingObjectSet: false,
		},
		Indices:  i,
		Material: mat,
	}
	mesh.setup()
	return mesh
}
func (m *MaterialMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()
	m.ebo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Vertices.Get(vertex.POSITION_NORMAL))

	m.wrapper.BindBuffer(glwrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	m.wrapper.ElementBufferData(m.Indices)

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*6, m.wrapper.PtrOffset(0))
	// setup normal vector
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*6, m.wrapper.PtrOffset(4*3))

	// close
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. It's input is a shader.
// First it binds the material with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function). It also sets up the model uniform.
// Then it binds the vertex array and draws the mesh with triangles. Finally it cleans up.
func (m *MaterialMesh) Draw(shader interfaces.Shader) {
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	diffuse := m.Material.GetDiffuse()
	ambient := m.Material.GetAmbient()
	specular := m.Material.GetSpecular()
	shininess := m.Material.GetShininess()
	shader.SetUniform3f("material.diffuse", diffuse.X(), diffuse.Y(), diffuse.Z())
	shader.SetUniform3f("material.ambient", ambient.X(), ambient.Y(), ambient.Z())
	shader.SetUniform3f("material.specular", specular.X(), specular.Y(), specular.Z())
	shader.SetUniform1f("material.shininess", shininess)
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawTriangleElements(int32(len(m.Indices)))

	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}

type PointMesh struct {
	Mesh
}

// NewPointMesh has only a glwrapper input, because it returns an empty mesh (without Vertices).
// Due to this, the vao, vbo setup is unnecessary now.
func NewPointMesh(wrapper interfaces.GLWrapper) *PointMesh {
	mesh := &PointMesh{
		Mesh{
			Vertices: []vertex.Vertex{},

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			yaw:       0,
			pitch:     0,
			roll:      0,
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,

			boundingObjectSet: false,
		},
	}
	return mesh
}
func (m *PointMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Vertices.Get(vertex.POSITION_COLOR_SIZE))

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*7, m.wrapper.PtrOffset(0))
	// setup color vector
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*7, m.wrapper.PtrOffset(4*3))
	// setup point size
	m.wrapper.VertexAttribPointer(2, 1, glwrapper.FLOAT, false, 4*7, m.wrapper.PtrOffset(4*6))

	// close
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. It's input is a shader.
// First it binds the  model uniform with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function).
// Then it binds the vertex array and draws the mesh with arrays (points). Finally it cleans up.
func (m *PointMesh) Draw(shader interfaces.Shader) {
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawArrays(glwrapper.POINTS, 0, int32(len(m.Vertices)))

	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}

// AddVertex inserts a new vertex to the vertices. Then it calls setup,
// because the vertices are changed, so that we have to generate the vao again.
func (m *PointMesh) AddVertex(v vertex.Vertex) {
	m.Vertices.Add(v)
	m.setup()
}

type ColorMesh struct {
	Mesh
	Indices []uint32
	Color   []mgl32.Vec3
	ebo     uint32
}

// NewColorMesh gets the vertices, indices, colors, glwrapper as inputs and makes the
// necessary setup for a standing (not moving) colored mesh before returning it.
// The vbo, vao, ebo is also set.
func NewColorMesh(v []vertex.Vertex, i []uint32, color []mgl32.Vec3, wrapper interfaces.GLWrapper) *ColorMesh {
	mesh := &ColorMesh{
		Mesh: Mesh{
			Vertices: v,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			yaw:       0,
			pitch:     0,
			roll:      0,
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,

			boundingObjectSet: false,
		},
		Indices: i,
		Color:   color,
	}
	mesh.setup()
	return mesh
}
func (m *ColorMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()
	m.ebo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Vertices.Get(vertex.POSITION_COLOR))

	m.wrapper.BindBuffer(glwrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	m.wrapper.ElementBufferData(m.Indices)

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*6, m.wrapper.PtrOffset(0))
	// setup color vector
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*6, m.wrapper.PtrOffset(4*3))

	// closeColorMesh
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. It's input is a shader.
// First it binds the  model uniform with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function).
// Then it binds the vertex array and draws the mesh with arrays (points). Finally it cleans up.
func (m *ColorMesh) Draw(shader interfaces.Shader) {
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawTriangleElements(int32(len(m.Indices)))

	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}

type TexturedColoredMesh struct {
	Mesh
	Indices  []uint32
	Textures texture.Textures
	Color    []mgl32.Vec3
	ebo      uint32
}

// NewTexturedColoredMesh gets the vertices, indices, textures, colors, glwrapper as inputs
// and makes the necessary setup for a standing (not moving) textured colored mesh before
// returning it. The vbo, vao, ebo is also set.
func NewTexturedColoredMesh(v []vertex.Vertex, i []uint32, t texture.Textures, color []mgl32.Vec3, wrapper interfaces.GLWrapper) *TexturedColoredMesh {
	mesh := &TexturedColoredMesh{
		Mesh: Mesh{
			Vertices: v,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			yaw:       0,
			pitch:     0,
			roll:      0,
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,

			boundingObjectSet: false,
		},
		Indices:  i,
		Textures: t,
		Color:    color,
	}
	mesh.setup()
	return mesh
}
func (m *TexturedColoredMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()
	m.ebo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Vertices.Get(vertex.POSITION_COLOR_TEXCOORD))

	m.wrapper.BindBuffer(glwrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	m.wrapper.ElementBufferData(m.Indices)

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(0))
	// setup normals
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*3))
	// setup texture position
	m.wrapper.VertexAttribPointer(2, 2, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*6))

	// close
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. Its input is a shader.
// First it binds the textures with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function). Then it sets up the model uniform.
// Then it binds the vertex array and draws the mesh with triangles. Finally it cleans up.
func (m *TexturedColoredMesh) Draw(shader interfaces.Shader) {
	for _, item := range m.Textures {
		item.Bind()
		shader.SetUniform1i(item.UniformName, int32(item.Id-glwrapper.TEXTURE0))
	}
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawTriangleElements(int32(len(m.Indices)))

	m.Textures.UnBind()
	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}

type TexturedMaterialMesh struct {
	Mesh
	Indices  []uint32
	Textures texture.Textures
	Material *material.Material
	ebo      uint32
}

// NewTexturedMaterialMesh gets the vertices, indices, textures, material, glwrapper as
// inputs and makes the necessary setup for a standing (not moving) textured material
// mesh before returning it. The vbo, vao, ebo is also set.
func NewTexturedMaterialMesh(v []vertex.Vertex, i []uint32, t texture.Textures, material *material.Material, wrapper interfaces.GLWrapper) *TexturedMaterialMesh {
	mesh := &TexturedMaterialMesh{
		Mesh: Mesh{
			Vertices: v,

			position:  mgl32.Vec3{0, 0, 0},
			direction: mgl32.Vec3{0, 0, 0},
			velocity:  0,
			yaw:       0,
			pitch:     0,
			roll:      0,
			scale:     mgl32.Vec3{1, 1, 1},
			wrapper:   wrapper,
			parentSet: false,

			boundingObjectSet: false,
		},
		Indices:  i,
		Textures: t,
		Material: material,
	}
	mesh.setup()
	return mesh
}
func (m *TexturedMaterialMesh) setup() {
	m.vao = m.wrapper.GenVertexArrays()
	m.vbo = m.wrapper.GenBuffers()
	m.ebo = m.wrapper.GenBuffers()

	m.wrapper.BindVertexArray(m.vao)

	m.wrapper.BindBuffer(glwrapper.ARRAY_BUFFER, m.vbo)
	m.wrapper.ArrayBufferData(m.Vertices.Get(vertex.POSITION_NORMAL_TEXCOORD))

	m.wrapper.BindBuffer(glwrapper.ELEMENT_ARRAY_BUFFER, m.ebo)
	m.wrapper.ElementBufferData(m.Indices)

	// setup coordinates
	m.wrapper.VertexAttribPointer(0, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(0))
	// setup normals
	m.wrapper.VertexAttribPointer(1, 3, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*3))
	// setup texture position
	m.wrapper.VertexAttribPointer(2, 2, glwrapper.FLOAT, false, 4*8, m.wrapper.PtrOffset(4*6))

	// close
	m.wrapper.BindVertexArray(0)
}

// Draw function is responsible for the actual drawing. Its input is a shader.
// First it binds the textures with the help of the shader (i expect that the shader
// is activated with the UseProgram gl function). Then it binds the material and sets
// up the model uniform. Then it binds the vertex array and draws the mesh with
// triangles. Finally it cleans up.
func (m *TexturedMaterialMesh) Draw(shader interfaces.Shader) {
	for _, item := range m.Textures {
		item.Bind()
		shader.SetUniform1i(item.UniformName, int32(item.Id-glwrapper.TEXTURE0))
	}
	M := m.ModelTransformation()
	shader.SetUniformMat4("model", M)
	diffuse := m.Material.GetDiffuse()
	ambient := m.Material.GetAmbient()
	specular := m.Material.GetSpecular()
	shininess := m.Material.GetShininess()
	shader.SetUniform3f("material.diffuse", diffuse.X(), diffuse.Y(), diffuse.Z())
	shader.SetUniform3f("material.ambient", ambient.X(), ambient.Y(), ambient.Z())
	shader.SetUniform3f("material.specular", specular.X(), specular.Y(), specular.Z())
	shader.SetUniform1f("material.shininess", shininess)
	m.wrapper.BindVertexArray(m.vao)
	m.wrapper.DrawTriangleElements(int32(len(m.Indices)))

	m.Textures.UnBind()
	m.wrapper.BindVertexArray(0)
	m.wrapper.ActiveTexture(0)
}
