package interfaces

import (
	"unsafe"

	"github.com/akosgarai/playground_engine/pkg/primitives/boundingobject"

	"github.com/akosgarai/coldet"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Shader interface {
	Use()
	SetUniformMat4(string, mgl32.Mat4)
	GetId() uint32
	SetUniform3f(string, float32, float32, float32)
	SetUniform1f(string, float32)
	SetUniform1i(string, int32)
}
type DirectionalLight interface {
	GetDirection() mgl32.Vec3
	GetAmbient() mgl32.Vec3
	GetDiffuse() mgl32.Vec3
	GetSpecular() mgl32.Vec3
}
type PointLight interface {
	GetPosition() mgl32.Vec3
	GetAmbient() mgl32.Vec3
	GetDiffuse() mgl32.Vec3
	GetSpecular() mgl32.Vec3
	GetConstantTerm() float32
	GetLinearTerm() float32
	GetQuadraticTerm() float32
}
type SpotLight interface {
	GetPosition() mgl32.Vec3
	GetDirection() mgl32.Vec3
	GetAmbient() mgl32.Vec3
	GetDiffuse() mgl32.Vec3
	GetSpecular() mgl32.Vec3
	GetConstantTerm() float32
	GetLinearTerm() float32
	GetQuadraticTerm() float32
	GetCutoff() float32
	GetOuterCutoff() float32
}

type GLWrapper interface {
	GenVertexArrays() uint32
	GenBuffers() uint32
	BindVertexArray(vao uint32)
	BindBuffer(bufferType, vbo uint32)
	ArrayBufferData(bufferData []float32)
	ElementBufferData(bufferData []uint32)
	VertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, pointer unsafe.Pointer)
	ActiveTexture(id uint32)
	BindTexture(id, textureId uint32)
	DrawTriangleElements(count int32)
	UseProgram(id uint32)
	GetUniformLocation(shaderProgramId uint32, uniformName string) int32
	Uniform1i(location int32, value int32)
	CreateProgram() uint32
	AttachShader(program, shader uint32)
	LinkProgram(program uint32)
	UniformMatrix4fv(location int32, count int32, transpose bool, value *float32)
	CreateShader(shaderType uint32) uint32
	Strs(strs string) (**uint8, func())
	ShaderSource(shader uint32, count int32, xstring **uint8, length *int32)
	CompileShader(id uint32)
	GetShaderiv(shader uint32, pname uint32, params *int32)
	GetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *uint8)
	Str(str string) *uint8
	InitOpenGL()
	TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, pixels unsafe.Pointer)
	Ptr(data interface{}) unsafe.Pointer
	GenerateMipmap(target uint32)
	GenTextures(n int32, textures *uint32)
	UniformMatrix3fv(location int32, count int32, transpose bool, value *float32)
	Uniform3f(location int32, v0 float32, v1 float32, v2 float32)
	Uniform1f(location int32, v0 float32)
	PtrOffset(offset int) unsafe.Pointer
	DisableVertexAttribArray(index uint32)
	DrawArrays(mode uint32, first int32, count int32)
	TexParameteri(target uint32, pname uint32, param int32)
	TexParameterfv(target uint32, pname uint32, params *float32)
	ClearColor(red float32, green float32, blue float32, alpha float32)
	Clear(mask uint32)
	Enable(cap uint32)
	DepthFunc(xfunc uint32)
	Viewport(x int32, y int32, width int32, height int32)
	BlendFunc(sfactor uint32, dfactor uint32)
}

type Mesh interface {
	Draw(Shader)
	Update(float64)
	SetSpeed(float32)
	SetDirection(mgl32.Vec3)
	GetPosition() mgl32.Vec3
	SetPosition(mgl32.Vec3)
	ModelTransformation() mgl32.Mat4
	TranslationTransformation() mgl32.Mat4
	GetParentTranslationTransformation() mgl32.Mat4
	RotationTransformation() mgl32.Mat4
	ScaleTransformation() mgl32.Mat4
	IsParentMesh() bool
	RotateX(float32)
	RotateY(float32)
	RotateZ(float32)
	RotatePosition(float32, mgl32.Vec3)
	IsBoundingObjectSet() bool
	GetBoundingObject() *boundingobject.BoundingObject
	GetParent() Mesh
}
type Model interface {
	Draw(Shader)
	Update(float64)
	Export(string)
	CollideTestWithSphere(*coldet.Sphere) bool
	IsTransparent() bool
	ClosestMeshTo(mgl32.Vec3) (Mesh, float32)
	Clear()
	AddMesh(Mesh)
}
type FormItem interface {
	Model
	GetLabel() string
	ValueToString() string
}
type CharFormItem interface {
	FormItem
	DeleteLastCharacter()
	GetCursorInitialPosition() mgl32.Vec3
	GetTarget() Mesh
}

type KeyStore interface {
	Get(glfw.Key) bool
	Set(glfw.Key, bool)
}
type RoKeyStore interface {
	Get(glfw.Key) bool
}
type ButtonStore interface {
	Get(glfw.MouseButton) bool
	Set(glfw.MouseButton, bool)
}
type RoButtonStore interface {
	Get(glfw.MouseButton) bool
}

type Camera interface {
	Log() string
	GetViewMatrix() mgl32.Mat4
	GetProjectionMatrix() mgl32.Mat4
	Walk(float32)
	Strafe(float32)
	Lift(float32)
	UpdateDirection(float32, float32)
	GetPosition() mgl32.Vec3
	GetVelocity() float32
	GetRotationStep() float32
	BoundingObjectAfterWalk(float32) *coldet.Sphere
	BoundingObjectAfterStrafe(float32) *coldet.Sphere
	BoundingObjectAfterLift(float32) *coldet.Sphere
}
type Screen interface {
	Log() string
	Draw(GLWrapper)
	Update(float64, float64, float64, RoKeyStore, RoButtonStore)
	Export(string)
	GetCamera() Camera
	GetClosestModelMeshDistance() (Model, Mesh, float32)
	SetUniformFloat(string, float32)
	SetUniformVector(string, mgl32.Vec3)
	CharCallback(rune, GLWrapper)
	SetWindowSize(float32, float32)
	SetWrapper(GLWrapper)
}
