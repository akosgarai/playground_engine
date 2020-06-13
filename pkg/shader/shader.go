package shader

import (
	"path"
	"runtime"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	defaultAppPath = "/shaders/"
)

type Shader struct {
	id      uint32
	wrapper interfaces.GLWrapper
}

// NewShader returns a Shader. It's inputs are the filenames of the shaders.
// It reads the files and compiles them. The shaders are attached to the shader program.
func NewShader(vertexShaderPath, fragmentShaderPath string, wrapper interfaces.GLWrapper) *Shader {
	vertexShaderSource, err := LoadShaderFromFile(vertexShaderPath)
	if err != nil {
		panic(err)
	}
	vertexShader, err := CompileShader(vertexShaderSource, glwrapper.VERTEX_SHADER, wrapper)
	if err != nil {
		panic(err)
	}
	fragmentShaderSource, err := LoadShaderFromFile(fragmentShaderPath)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := CompileShader(fragmentShaderSource, glwrapper.FRAGMENT_SHADER, wrapper)
	if err != nil {
		panic(err)
	}

	program := wrapper.CreateProgram()
	wrapper.AttachShader(program, vertexShader)
	wrapper.AttachShader(program, fragmentShader)
	wrapper.LinkProgram(program)

	return &Shader{
		id:      program,
		wrapper: wrapper,
	}
}
func baseDirShaders() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename) + defaultAppPath
}

// NewTextureShader returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used.
func NewTextureShader(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texture.vert", baseDirShaders()+"texture.frag", wrapper)
}

// NewMaterialShader returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used.
func NewMaterialShader(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"material.vert", baseDirShaders()+"material.frag", wrapper)
}

// NewTextureMatShader returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used.
func NewTextureMatShader(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texturemat.vert", baseDirShaders()+"texturemat.frag", wrapper)
}

// Use is a wrapper for gl.UseProgram
func (s *Shader) Use() {
	s.wrapper.UseProgram(s.id)
}

// GetId returns the program identifier of the shader.
func (s *Shader) GetId() uint32 {
	return s.id
}

// SetUniformMat4 gets an uniform name string and the value matrix as input and
// calls the gl.UniformMatrix4fv function
func (s *Shader) SetUniformMat4(uniformName string, mat mgl32.Mat4) {
	location := s.wrapper.GetUniformLocation(s.id, uniformName)
	s.wrapper.UniformMatrix4fv(location, 1, false, &mat[0])
}

// SetUniform3f gets an uniform name string and 3 float values as input and
// calls the gl.Uniform3f function
func (s *Shader) SetUniform3f(uniformName string, v1, v2, v3 float32) {
	location := s.wrapper.GetUniformLocation(s.id, uniformName)
	s.wrapper.Uniform3f(location, v1, v2, v3)
}

// SetUniform1f gets an uniform name string and a float value as input and
// calls the gl.Uniform1f function
func (s *Shader) SetUniform1f(uniformName string, v1 float32) {
	location := s.wrapper.GetUniformLocation(s.id, uniformName)
	s.wrapper.Uniform1f(location, v1)
}

// SetUniform1i gets an uniform name string and an integer value as input and
// calls the gl.Uniform1i function
func (s *Shader) SetUniform1i(uniformName string, v1 int32) {
	location := s.wrapper.GetUniformLocation(s.id, uniformName)
	s.wrapper.Uniform1i(location, v1)
}
