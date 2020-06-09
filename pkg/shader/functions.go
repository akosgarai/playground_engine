package shader

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/akosgarai/playground_engine/pkg/glwrapper"
	"github.com/akosgarai/playground_engine/pkg/interfaces"
)

// LoadShaderFromFile takes a filepath string arguments.
// It loads the file and returns it as a '\x00' terminated string.
// It returns an error also.
func LoadShaderFromFile(path string) (string, error) {
	shaderCode, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	result := string(shaderCode) + "\x00"
	return result, nil
}

// CompileShader creeates a shader, compiles the shader source, and returns
// the uint32 identifier of the shader and nil. If the compile fails, it returns
// an error and 0 as shader id.
func CompileShader(source string, shaderType uint32, wrapper interfaces.GLWrapper) (uint32, error) {
	shader := wrapper.CreateShader(shaderType)

	csources, free := wrapper.Strs(source)
	wrapper.ShaderSource(shader, 1, csources, nil)
	free()
	wrapper.CompileShader(shader)

	var status int32
	wrapper.GetShaderiv(shader, glwrapper.COMPILE_STATUS, &status)
	if status == glwrapper.FALSE {
		var logLength int32
		wrapper.GetShaderiv(shader, glwrapper.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		wrapper.GetShaderInfoLog(shader, logLength, nil, wrapper.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
