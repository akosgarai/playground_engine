package glwrapper

import (
	"runtime"
	"strings"
	"testing"

	"github.com/akosgarai/playground_engine/pkg/testhelper"
)

const (
	VertexShader = `
#version 410
layout(location = 0) in vec3 vVertex;
uniform mat4 model;
uniform mat3 matrixname;
uniform float floatName;
uniform int intName;
uniform vec3 uniVec;
void main()
{
 float val = floatName * 2;
 mat3 tmp = matrixname * val;
 gl_Position = model * vec4(vVertex,1);
}
 ` + "\x00"
	FragmentShader = `
# version 410
out vec4 FragColor;
void main()
{
    FragColor = vec4(1.0);
}
 ` + "\x00"
)

var w Wrapper

func setup() {
	runtime.LockOSThread()
	testhelper.GlfwInit()
	w.InitOpenGL()
}

func TestGenVertexArrays(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		w.GenVertexArrays()
	}()
}
func TestGenBuffers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		w.GenBuffers()
	}()
}
func TestBindVertexArray(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		id := w.GenVertexArrays()
		w.BindVertexArray(id)
	}()
}
func TestBindBuffer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		id := w.GenBuffers()
		w.BindBuffer(ELEMENT_ARRAY_BUFFER, id)
	}()
}
func TestArrayBufferData(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		id := w.GenBuffers()
		w.BindBuffer(ARRAY_BUFFER, id)
		data := []float32{0.0, 1.0, 2.0}
		w.ArrayBufferData(data)
	}()
}
func TestElementBufferData(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		id := w.GenBuffers()
		w.BindBuffer(ELEMENT_ARRAY_BUFFER, id)
		data := []uint32{0, 1, 2}
		w.ElementBufferData(data)
	}()
}
func TestVertexAttribPointer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		id := w.GenBuffers()
		w.BindBuffer(ARRAY_BUFFER, id)
		data := []float32{0.0, 1.0, 2.0}
		w.ArrayBufferData(data)
		w.VertexAttribPointer(0, 3, FLOAT, false, 4*3, w.PtrOffset(0))
	}()
}
func TestActiveTexture(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		var id uint32
		w.GenTextures(1, &id)
		w.ActiveTexture(TEXTURE0)
	}()
}
func TestBindTexture(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		w.ActiveTexture(TEXTURE0)
		var id uint32
		w.GenTextures(1, &id)
		w.BindTexture(TEXTURE_2D, id)
	}()
}
func TestDrawTriangleElements(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		fs := w.CreateShader(FRAGMENT_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
		source, free = w.Strs(FragmentShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(fs)
		w.AttachShader(prog, fs)
		w.LinkProgram(prog)
		w.UseProgram(prog)
		vao := w.GenVertexArrays()
		vbo := w.GenBuffers()
		ebo := w.GenBuffers()
		w.BindVertexArray(vao)
		w.BindBuffer(ARRAY_BUFFER, vbo)
		data := []float32{0, 0, 0, 0, 1, 0, 1, 1, 0}
		w.ArrayBufferData(data)
		w.BindBuffer(ELEMENT_ARRAY_BUFFER, ebo)
		indices := []uint32{1, 2, 3, 1, 3, 4}
		w.ElementBufferData(indices)
		w.VertexAttribPointer(0, 3, FLOAT, false, 4*3, w.PtrOffset(0))
		w.DrawTriangleElements(int32(6))
	}()
}
func TestUseProgram(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		fs := w.CreateShader(FRAGMENT_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
		source, free = w.Strs(FragmentShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(fs)
		w.AttachShader(prog, fs)
		w.LinkProgram(prog)
		w.UseProgram(prog)
	}()
}
func TestGetUniformLocation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		fs := w.CreateShader(FRAGMENT_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
		source, free = w.Strs(FragmentShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(fs)
		w.AttachShader(prog, fs)
		w.LinkProgram(prog)
		location := w.GetUniformLocation(prog, "some-uniform-name")
		if location != -1 {
			t.Errorf("Invalid location id for not existing uniform name. '%v'.\n", location)
		}
	}()
}
func TestUniform1i(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		fs := w.CreateShader(FRAGMENT_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
		source, free = w.Strs(FragmentShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(fs)
		w.AttachShader(prog, fs)
		w.LinkProgram(prog)
		w.UseProgram(prog)
		w.Uniform1i(w.GetUniformLocation(prog, "intName"), 16.0)
	}()
}
func TestCreateProgram(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		w.CreateProgram()
	}()
}
func TestAttachShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
	}()
}
func TestLinkProgram(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		fs := w.CreateShader(FRAGMENT_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
		source, free = w.Strs(FragmentShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(fs)
		w.AttachShader(prog, fs)
		w.LinkProgram(prog)
	}()
}
func TestUniformMatrix4fv(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		fs := w.CreateShader(FRAGMENT_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
		source, free = w.Strs(FragmentShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(fs)
		w.AttachShader(prog, fs)
		w.LinkProgram(prog)
		w.UseProgram(prog)
		mat := [16]float32{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
		w.UniformMatrix4fv(w.GetUniformLocation(prog, "model"), 16, false, &mat[0])
	}()
}
func TestCreateShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		w.CreateShader(VERTEX_SHADER)
	}()
}
func TestStrs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		input := "test"
		_, free := w.Strs(input)
		free()
	}()
}
func TestShaderSource(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		shader := w.CreateShader(VERTEX_SHADER)
		input := "test"
		source, free := w.Strs(input)
		w.ShaderSource(shader, 1, source, nil)
		free()
	}()
}
func TestCompileShader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		shader := w.CreateShader(VERTEX_SHADER)
		input := "test"
		source, free := w.Strs(input)
		w.ShaderSource(shader, 1, source, nil)
		free()
		w.CompileShader(shader)
	}()
}
func TestGetShaderiv(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		shader := w.CreateShader(VERTEX_SHADER)
		input := "test"
		source, free := w.Strs(input)
		w.ShaderSource(shader, 1, source, nil)
		free()
		w.CompileShader(shader)
		var status int32
		w.GetShaderiv(shader, COMPILE_STATUS, &status)
		if status != FALSE {
			t.Error("It should be failed, due to the wrong input.")
		}
	}()
}
func TestGetShaderInfoLog(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		shader := w.CreateShader(VERTEX_SHADER)
		input := "test"
		source, free := w.Strs(input)
		w.ShaderSource(shader, 1, source, nil)
		free()
		w.CompileShader(shader)
		var logLength int32
		w.GetShaderiv(shader, INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		w.GetShaderInfoLog(shader, logLength, nil, w.Str(log))
		if !strings.Contains(log, "syntax error") {
			t.Errorf("The error '%s' should contains 'syntax error'.", log)
		}
	}()
}
func TestStr(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		shader := w.CreateShader(VERTEX_SHADER)
		input := "test"
		source, free := w.Strs(input)
		w.ShaderSource(shader, 1, source, nil)
		free()
		w.CompileShader(shader)
		var logLength int32
		w.GetShaderiv(shader, INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		w.GetShaderInfoLog(shader, logLength, nil, w.Str(log))
		if !strings.Contains(log, "syntax error") {
			t.Errorf("The error '%s' should contains 'syntax error'.", log)
		}
	}()
}
func TestInitOpenGL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		w.InitOpenGL()
	}()
}
func TestTexImage2D(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		var id uint32
		w.GenTextures(1, &id)
		w.ActiveTexture(TEXTURE0)
		w.BindTexture(TEXTURE_2D, id)
		w.TexParameteri(TEXTURE_2D, TEXTURE_WRAP_R, CLAMP_TO_EDGE)
		w.TexParameteri(TEXTURE_2D, TEXTURE_WRAP_S, CLAMP_TO_EDGE)
		w.TexParameteri(TEXTURE_2D, TEXTURE_MIN_FILTER, LINEAR)
		w.TexParameteri(TEXTURE_2D, TEXTURE_MAG_FILTER, LINEAR)
		w.TexImage2D(TEXTURE_2D, 0, RGBA, int32(1), int32(1), 0, RGBA, uint32(UNSIGNED_BYTE), w.Ptr([]uint8{0, 0, 0}))

	}()
}
func TestPtr(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		var id uint32
		w.GenTextures(1, &id)
		w.ActiveTexture(TEXTURE0)
		w.BindTexture(TEXTURE_2D, id)
		w.TexParameteri(TEXTURE_2D, TEXTURE_WRAP_R, CLAMP_TO_EDGE)
		w.TexParameteri(TEXTURE_2D, TEXTURE_WRAP_S, CLAMP_TO_EDGE)
		w.TexParameteri(TEXTURE_2D, TEXTURE_MIN_FILTER, LINEAR)
		w.TexParameteri(TEXTURE_2D, TEXTURE_MAG_FILTER, LINEAR)
		w.TexImage2D(TEXTURE_2D, 0, RGBA, int32(1), int32(1), 0, RGBA, uint32(UNSIGNED_BYTE), w.Ptr([]uint8{0, 0, 0}))

	}()
}
func TestGenerateMipmap(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		var id uint32
		w.GenTextures(1, &id)
		w.GenerateMipmap(id)
	}()
}
func TestGenTextures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		var id uint32
		w.GenTextures(1, &id)
	}()
}
func TestUniformMatrix3fv(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		fs := w.CreateShader(FRAGMENT_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
		source, free = w.Strs(FragmentShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(fs)
		w.AttachShader(prog, fs)
		w.LinkProgram(prog)
		w.UseProgram(prog)
		mat := [9]float32{1, 0, 0, 0, 1, 0, 0, 0, 1}
		w.UniformMatrix3fv(w.GetUniformLocation(prog, "matrixname"), 9, false, &mat[0])
	}()
}
func TestUniform3f(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		fs := w.CreateShader(FRAGMENT_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
		source, free = w.Strs(FragmentShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(fs)
		w.AttachShader(prog, fs)
		w.LinkProgram(prog)
		w.UseProgram(prog)
		w.Uniform3f(w.GetUniformLocation(prog, "uniVec"), 1, 1, 1)
	}()
}
func TestUniform1f(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		fs := w.CreateShader(FRAGMENT_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
		source, free = w.Strs(FragmentShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(fs)
		w.AttachShader(prog, fs)
		w.LinkProgram(prog)
		w.UseProgram(prog)
		w.Uniform1f(w.GetUniformLocation(prog, "floatName"), 16.0)
	}()
}
func TestPtrOffset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		w.PtrOffset(0)
	}()
}
func TestDisableVertexAttribArray(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		id := w.GenBuffers()
		w.BindBuffer(ARRAY_BUFFER, id)
		data := []float32{0.0, 1.0, 2.0}
		w.ArrayBufferData(data)
		w.VertexAttribPointer(0, 3, FLOAT, false, 4*3, w.PtrOffset(0))
		w.DisableVertexAttribArray(0)
	}()
}
func TestDrawArrays(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		vs := w.CreateShader(VERTEX_SHADER)
		fs := w.CreateShader(FRAGMENT_SHADER)
		prog := w.CreateProgram()
		source, free := w.Strs(VertexShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(vs)
		w.AttachShader(prog, vs)
		source, free = w.Strs(FragmentShader)
		w.ShaderSource(vs, 1, source, nil)
		free()
		w.CompileShader(fs)
		w.AttachShader(prog, fs)
		w.LinkProgram(prog)
		w.UseProgram(prog)
		vao := w.GenVertexArrays()
		vbo := w.GenBuffers()
		w.BindVertexArray(vao)
		w.BindBuffer(ARRAY_BUFFER, vbo)
		data := []float32{0, 0, 0, 0, 1, 0, 1, 1, 0}
		w.ArrayBufferData(data)
		w.VertexAttribPointer(0, 3, FLOAT, false, 4*3, w.PtrOffset(0))
		w.DrawArrays(TRIANGLES, 0, int32(3))
	}()
}
func TestTexParameteri(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Error("It shouldn't fail.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		var id uint32
		w.GenTextures(1, &id)
		w.ActiveTexture(TEXTURE0)
		w.BindTexture(TEXTURE_2D, id)
		w.TexParameteri(TEXTURE_2D, TEXTURE_WRAP_R, CLAMP_TO_EDGE)
		w.TexParameteri(TEXTURE_2D, TEXTURE_WRAP_S, CLAMP_TO_EDGE)
		w.TexParameteri(TEXTURE_2D, TEXTURE_MIN_FILTER, LINEAR)
		w.TexParameteri(TEXTURE_2D, TEXTURE_MAG_FILTER, LINEAR)
		w.TexImage2D(TEXTURE_2D, 0, RGBA, int32(1), int32(1), 0, RGBA, uint32(UNSIGNED_BYTE), w.Ptr([]uint8{0, 0, 0}))

	}()
}
func TestTexParameterfv(t *testing.T) {
	t.Skip("Needs to be implemented.")
}
func TestClearColor(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		w.ClearColor(0, 0, 0, 1)
	}()
}
func TestClear(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		w.Clear(COLOR_BUFFER_BIT)
	}()
}
func TestEnable(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		w.Enable(DEPTH_TEST)
	}()
}
func TestDepthFunc(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		w.Enable(DEPTH_TEST)
		w.DepthFunc(LESS)
	}()
}
func TestViewport(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping it in short mode")
	}
	func() {
		defer func() {
			if r := recover(); r != nil {
				defer testhelper.GlfwTerminate()
				t.Error("It shouldn't fail with setup.")
			}
		}()
		setup()
		defer testhelper.GlfwTerminate()
		w.Viewport(0, 0, 800, 800)
	}()
}
