package shader

import (
	"os"
	"runtime"
	"testing"

	wrapper "github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/light"
	"github.com/akosgarai/playground_engine/pkg/testhelper"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	ValidFragmentShaderString = `
#version 410
smooth in vec4 vSmoothColor;
layout(location=0) out vec4 vFragColor;
void main()
{
    vFragColor = vSmoothColor;
}
    `
	ValidVertexShaderWithUniformsString = `
#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
const float pointSize = 20.0;
smooth out vec4 vSmoothColor;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
void main()
{
    gl_Position = projection * view * model * vec4(vVertex,1);
    gl_PointSize = pointSize;
    vSmoothColor = vec4(vColor,1);
}
    `
	ValidVertexShaderWithUniformsStringWithTrailingChars = `
#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
const float pointSize = 20.0;
smooth out vec4 vSmoothColor;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
void main()
{
    gl_Position = projection * view * model * vec4(vVertex,1);
    gl_PointSize = pointSize;
    vSmoothColor = vec4(vColor,1);
}
    ` + "\x00"
	ValidVertexShaderWithMat3String = `
#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
smooth out vec4 vSmoothColor;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform mat3 normal;		//normal matrix
void main()
{
    vec3 vNormal = normalize(normal * vVertex);
    gl_Position = projection * view * model * vec4(vNormal,1);
    vSmoothColor = vec4(vColor,1);
}
    `
	ValidVertexShaderWithFloatUniformString = `
#version 410
layout(location = 0) in vec3 vVertex;
layout(location = 1) in vec3 vColor;
smooth out vec4 vSmoothColor;
uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;
uniform float pointSize;
void main()
{
    gl_Position = projection * view * model * vec4(vVertex,1);
    gl_PointSize = pointSize;
    vSmoothColor = vec4(vColor,1);
}
    `
	InvalidShaderString = `
    This string is not valid as
    a shader progrma.
    `
	InvalidShaderStringWithTrailingChars = `
    This string is not valid as
    a shader progrma.
    ` + "\x00"
	ValidTextureVertexShader = `
# version 410
layout (location = 0) in vec3 vVertex;
layout (location = 1) in vec3 vColor;
layout (location = 2) in vec2 vTexCoord;

out vec3 vSmoothColor;
out vec2 vSmoothTexCoord;

void main()
{
    gl_Position = vec4(vVertex, 1.0);
    vSmoothColor = vColor;
    vSmoothTexCoord = vec2(vTexCoord.x, vTexCoord.y);
}
    `
	ValidTextureFragmentShader = `
# version 410
out vec4 FragColor;
  
in vec3 vSmoothColor;
in vec2 vSmoothTexCoord;

uniform sampler2D textureOne;

void main()
{
    FragColor = texture(textureOne, vSmoothTexCoord) * vec4(vSmoothColor, 1.0);
}
    `
	EmptyString            = ""
	FragmentShaderFileName = "fragmentShader.frag"
	VertexShaderFileName   = "vertexShader.vert"
)

var (
	LightPosition      = mgl32.Vec3{0, 0, 0}
	LightDirection     = mgl32.Vec3{0, 1, 0}
	LightAmbient       = mgl32.Vec3{1, 1, 1}
	LightDiffuse       = mgl32.Vec3{1, 1, 1}
	LightSpecular      = mgl32.Vec3{1, 1, 1}
	LightConstantTerm  = float32(1.0)
	LightLinearTerm    = float32(0.5)
	LightQuadraticTerm = float32(0.05)
	LightCutoff        = float32(12.0)
	LightOuterCutoff   = float32(20.0)

	testGlWrapper testhelper.GLWrapperMock
	realGlWrapper wrapper.Wrapper
)

func NewPointLightSource() *light.Light {
	source := light.NewPointLight([4]mgl32.Vec3{LightPosition, LightAmbient, LightDiffuse, LightSpecular}, [3]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm})
	return source
}
func NewDirectionalLightSource() *light.Light {
	source := light.NewDirectionalLight([4]mgl32.Vec3{LightDirection, LightAmbient, LightDiffuse, LightSpecular})
	return source
}
func NewSpotLightSource() *light.Light {
	source := light.NewSpotLight([5]mgl32.Vec3{LightPosition, LightDirection, LightAmbient, LightDiffuse, LightSpecular}, [5]float32{LightConstantTerm, LightLinearTerm, LightQuadraticTerm, LightCutoff, LightOuterCutoff})
	return source
}

func CreateFileWithContent(name, content string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
func DeleteFile(name string) error {
	return os.Remove(name)
}
func InitGlfw() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(600, 600, "Test-window", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
}
func NewTestShader(t *testing.T, validFragmentShaderContent, validVertexShaderContent string) *Shader {
	CreateFileWithContent(FragmentShaderFileName, validFragmentShaderContent)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, validVertexShaderContent)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	testGlWrapper.InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName, testGlWrapper)
	if shader.id == 0 {
		t.Error("Invalid shader program id")
		t.Fail()
	}
	return shader
}
func TestLoadShaderFromFile(t *testing.T) {
	// Create tmp file with a known content.
	// call function with
	// - bad filename, that doesn't exist, so that we should have an error.
	// - good filename, that exists and we know it's content
	wrongFileName := "badfile.name"
	content, err := LoadShaderFromFile(wrongFileName)
	if err == nil {
		t.Error("Wrong filename should return error")
	}
	if content != EmptyString {
		t.Errorf("Wrong filename should return empty content. We got: '%s'", content)
	}
	goodFileName := "goodfile.name"
	CreateFileWithContent(goodFileName, InvalidShaderString)
	defer DeleteFile(goodFileName)
	content, err = LoadShaderFromFile(goodFileName)
	if err != nil {
		t.Error("Good file shouldn't return error")
	}
	if content == InvalidShaderString {
		t.Error("Good file content should have the trailing '\\x00'")
	}
	if content != InvalidShaderString+"\x00" {
		t.Error("Good file content should be the same")
	}
}
func TestCompileShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	realGlWrapper.InitOpenGL()
	_, err := CompileShader(InvalidShaderStringWithTrailingChars, wrapper.VERTEX_SHADER, realGlWrapper)
	if err == nil {
		t.Error("Compile should fail with wrong content.")
	}
	prog, err := CompileShader(ValidVertexShaderWithUniformsStringWithTrailingChars, wrapper.VERTEX_SHADER, realGlWrapper)
	if err != nil {
		t.Error(err)
	}
	if prog == 0 {
		t.Error("Invalid shader program id")
	}
}
func TestNewShaderPanicOnVertexContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer glfw.Terminate()
				t.Errorf("NewShader should have panicked due to the invalid content!")
			}
		}()
		CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
		defer DeleteFile(FragmentShaderFileName)
		CreateFileWithContent(VertexShaderFileName, InvalidShaderString)
		defer DeleteFile(VertexShaderFileName)
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		NewShader(VertexShaderFileName, FragmentShaderFileName, realGlWrapper)
	}()
}
func TestNewShaderPanicOnFragmentContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer glfw.Terminate()
				t.Errorf("NewShader should have panicked due to the invalid content!")
			}
		}()
		CreateFileWithContent(FragmentShaderFileName, InvalidShaderString)
		defer DeleteFile(FragmentShaderFileName)
		CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
		defer DeleteFile(VertexShaderFileName)
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		NewShader(VertexShaderFileName, FragmentShaderFileName, realGlWrapper)
	}()
}
func TestNewShaderPanicOnFragmentFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer glfw.Terminate()
				t.Errorf("NewShader should have panicked due to the missing file!")
			}
		}()
		CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
		defer DeleteFile(VertexShaderFileName)
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		NewShader(VertexShaderFileName, FragmentShaderFileName, realGlWrapper)
	}()
}
func TestNewShaderPanicOnVertexFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				defer glfw.Terminate()
				t.Errorf("NewShader should have panicked due to the missing file!")
			}
		}()
		CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
		defer DeleteFile(FragmentShaderFileName)
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		NewShader(VertexShaderFileName, FragmentShaderFileName, realGlWrapper)
	}()
}
func TestNewShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	realGlWrapper.InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName, realGlWrapper)
	if shader.id == 0 || shader.GetId() == 0 {
		t.Error("Invalid shader program id")
	}
}
func TestNewTextureShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Errorf("NewTextureShader shouldn't have panicked!")
			}
		}()
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		shader := NewTextureShader(realGlWrapper)
		if shader.id == 0 || shader.GetId() == 0 {
			t.Error("Invalid shader program id")
		}
	}()
}
func TestNewTextureShaderBlending(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Errorf("NewTextureShaderBlending shouldn't have panicked!")
			}
		}()
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		shader := NewTextureShaderBlending(realGlWrapper)
		if shader.id == 0 || shader.GetId() == 0 {
			t.Error("Invalid shader program id")
		}
	}()
}
func TestNewTextureShaderWithFog(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Errorf("NewTextureShaderWithFog shouldn't have panicked!")
			}
		}()
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		shader := NewTextureShaderWithFog(realGlWrapper)
		if shader.id == 0 || shader.GetId() == 0 {
			t.Error("Invalid shader program id")
		}
	}()
}
func TestNewTextureShaderBlendingWithFog(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Errorf("NewTextureShaderBlendingWithFog shouldn't have panicked!")
			}
		}()
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		shader := NewTextureShaderBlendingWithFog(realGlWrapper)
		if shader.id == 0 || shader.GetId() == 0 {
			t.Error("Invalid shader program id")
		}
	}()
}
func TestNewMaterialShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Errorf("NewMaterialShader shouldn't have panicked!")
			}
		}()
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		shader := NewMaterialShader(realGlWrapper)
		if shader.id == 0 || shader.GetId() == 0 {
			t.Error("Invalid shader program id")
		}
	}()
}
func TestNewMaterialShaderWithFog(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Errorf("NewMaterialShaderWithFog shouldn't have panicked!")
			}
		}()
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		shader := NewMaterialShaderWithFog(realGlWrapper)
		if shader.id == 0 || shader.GetId() == 0 {
			t.Error("Invalid shader program id")
		}
	}()
}
func TestNewTextureMatShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Errorf("NewTextureMatShader shouldn't have panicked!. %v", r)
			}
		}()
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		shader := NewTextureMatShader(realGlWrapper)
		if shader.id == 0 || shader.GetId() == 0 {
			t.Error("Invalid shader program id")
		}
	}()
}
func TestNewTextureMatShaderBlending(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Errorf("NewTextureMatShaderBlending shouldn't have panicked!. %v", r)
			}
		}()
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		shader := NewTextureMatShaderBlending(realGlWrapper)
		if shader.id == 0 || shader.GetId() == 0 {
			t.Error("Invalid shader program id")
		}
	}()
}
func TestNewTextureMatShaderWithFog(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Errorf("NewTextureMatShaderWithFog shouldn't have panicked!. %v", r)
			}
		}()
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		shader := NewTextureMatShaderWithFog(realGlWrapper)
		if shader.id == 0 || shader.GetId() == 0 {
			t.Error("Invalid shader program id")
		}
	}()
}
func TestNewTextureMatShaderBlendingWithFog(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer glfw.Terminate()
				t.Errorf("NewTextureMatShaderBlendingWithFog shouldn't have panicked!. %v", r)
			}
		}()
		runtime.LockOSThread()
		InitGlfw()
		defer glfw.Terminate()
		realGlWrapper.InitOpenGL()
		shader := NewTextureMatShaderBlendingWithFog(realGlWrapper)
		if shader.id == 0 || shader.GetId() == 0 {
			t.Error("Invalid shader program id")
		}
	}()
}
func TestUse(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	realGlWrapper.InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName, realGlWrapper)
	if shader.id == 0 || shader.GetId() == 0 {
		t.Error("Invalid shader program id")
	}
	shader.Use()
}
func TestSetUniformMat4(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithUniformsString)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	realGlWrapper.InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName, realGlWrapper)
	if shader.id == 0 || shader.GetId() == 0 {
		t.Error("Invalid shader program id")
	}
	shader.Use()
	shader.SetUniformMat4("model", mgl32.Ident4())
}
func TestSetUniform3f(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithMat3String)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	realGlWrapper.InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName, realGlWrapper)
	if shader.id == 0 || shader.GetId() == 0 {
		t.Error("Invalid shader program id")
	}
	shader.Use()
	shader.SetUniform3f("ambientColor", 1, 1, 1)
}
func TestSetUniform1f(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithFloatUniformString)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	realGlWrapper.InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName, realGlWrapper)
	if shader.id == 0 || shader.GetId() == 0 {
		t.Error("Invalid shader program id")
	}
	var valueToSet float32
	valueToSet = 20
	shader.Use()
	shader.SetUniform1f("pointSize", valueToSet)
}
func TestSetUniform1i(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	CreateFileWithContent(FragmentShaderFileName, ValidFragmentShaderString)
	defer DeleteFile(FragmentShaderFileName)
	CreateFileWithContent(VertexShaderFileName, ValidVertexShaderWithFloatUniformString)
	defer DeleteFile(VertexShaderFileName)
	runtime.LockOSThread()
	InitGlfw()
	defer glfw.Terminate()
	realGlWrapper.InitOpenGL()
	shader := NewShader(VertexShaderFileName, FragmentShaderFileName, realGlWrapper)
	if shader.id == 0 || shader.GetId() == 0 {
		t.Error("Invalid shader program id")
	}
	var valueToSet int32
	valueToSet = 20
	shader.Use()
	shader.SetUniform1i("pointSize", valueToSet)
}
