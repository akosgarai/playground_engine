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

// NewTextureShaderLiquid returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used.
func NewTextureShaderLiquid(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texture_liquid.vert", baseDirShaders()+"texture_liquid.frag", wrapper)
}

// NewTextureShaderLiquidWithFog returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used.
func NewTextureShaderLiquidWithFog(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texture_liquid.vert", baseDirShaders()+"texture_liquid_with_fog.frag", wrapper)
}

// NewTextureShaderWithFog returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used. In this application, the `Fog`
// structure has to be filled. The `fog.minDistance`, `fog.maxDistance` floats and the `fog.color` mgl32.Vec3.
func NewTextureShaderWithFog(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texture.vert", baseDirShaders()+"texture_with_fog.frag", wrapper)
}

// NewTextureShaderBlending returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used.
func NewTextureShaderBlending(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texture.vert", baseDirShaders()+"texture_blending.frag", wrapper)
}

// NewTextureShaderBlendingWithFog returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used. In this application, the `Fog`
// structure has to be filled. The `fog.minDistance`, `fog.maxDistance` floats and the `fog.color` mgl32.Vec3.
func NewTextureShaderBlendingWithFog(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texture.vert", baseDirShaders()+"texture_blending_with_fog.frag", wrapper)
}

// NewMaterialShader returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used.
func NewMaterialShader(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"material.vert", baseDirShaders()+"material.frag", wrapper)
}

// NewMaterialShaderWithFog returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used. In this application, the `Fog`
// structure has to be filled. The `fog.minDistance`, `fog.maxDistance` floats and the `fog.color` mgl32.Vec3.
func NewMaterialShaderWithFog(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"material.vert", baseDirShaders()+"material_with_fog.frag", wrapper)
}

// NewTextureMatShader returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used.
func NewTextureMatShader(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texturemat.vert", baseDirShaders()+"texturemat.frag", wrapper)
}

// NewTextureMatShaderBlending returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used.
func NewTextureMatShaderBlending(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texturemat.vert", baseDirShaders()+"texturemat_blending.frag", wrapper)
}

// NewTextureMatShaderWithFog returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used. In this application, the `Fog`
// structure has to be filled. The `fog.minDistance`, `fog.maxDistance` floats and the `fog.color` mgl32.Vec3.
func NewTextureMatShaderWithFog(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texturemat.vert", baseDirShaders()+"texturemat_with_fog.frag", wrapper)
}

// NewTextureMatShaderBlendingWithFog returns a Shader, that uses the default texture vertex & fragment shaders.
// It works the same as NewShader, but the internal shader files are used. In this application, the `Fog`
// structure has to be filled. The `fog.minDistance`, `fog.maxDistance` floats and the `fog.color` mgl32.Vec3.
func NewTextureMatShaderBlendingWithFog(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"texturemat.vert", baseDirShaders()+"texturemat_blending_with_fog.frag", wrapper)
}

// NewFontShader returns a Shader, that could be user for rendering fonts.
func NewFontShader(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"font.vert", baseDirShaders()+"font.frag", wrapper)
}

// NewMenuBackgroundShader returns a Shader, that could be used for rendering the background
// for the menu screen.
func NewMenuBackgroundShader(wrapper interfaces.GLWrapper) *Shader {
	return NewShader(baseDirShaders()+"menu-background.vert", baseDirShaders()+"menu-background.frag", wrapper)
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
