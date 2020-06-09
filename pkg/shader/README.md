# Shader

This structure helps us to create shader applications. It stores an `id`, that was generated with the gl.CreateProgram command, and the wrapper, that is ued to call gl. functions.


## Functions

The stuff that we can do with the shaders.

### LoadShaderFromFile

LoadShaderFromFile takes a filepath string arguments. It loads the file and returns it as a `\x00` terminated string. It returns an error also.

### CompileShader

It compiles the given shader program string as the given type of shader. In case of success the shader program id and nil returns. On case of error, it's returned instead of the nil.

### NewShader

NewShader returns a Shader. It's inputs are the filenames of the shaders, and the glwrapper instance. It reads the files and compiles them. The shaders are attached to the shader program.

### Use

Use is a wrapper for gl.UseProgram

### GetId

GetId returns the shader program id.

### SetUniformMat4

SetUniformMat4 gets an uniform name string and the value matrix as input and calls the gl.UniformMatrix4fv function through its wrapper.

### SetUniform3f

SetUniform3f gets an uniform name string and 3 float values as input and calls the gl.Uniform3f function through its wrapper.

### SetUniform1f

SetUniform1f gets an uniform name string and a float value as input and calls the gl.Uniform1f function through its wrapper.

### SetUniform1i

SetUniform1i gets an uniform name string and an integer value as input and calls the gl.Uniform1i function through its wrapper.

## Default shader applications

These applications are implemented for decreasing the code duplication. As i developed the example apps, i experienced that i'm reusing the shader applications without changing them. It is a copy-paste step. I want to eliminate this unnecessary copy-paste.

### Naming conventions & variables

For the lightsources the following variables are set:

- 'NumberOfDirectionalLightSources'
- 'NumberOfPointLightSources'
- 'NumberOfSpotLightSources'

They are used for tracking the number of the distinct lighstources. The loop counters are maximized to these values.


### Material

This shader is written to handle material objects. It doesn't support textures, only materials. The maximum number of lighsources is 16. You can add more, but the surplus will not be handled.

### Texture

This shader is written to handle plane textured objects. It doesn't support materials or colors. The maximum number of lighsources is 16. You can add more, but the surplus will not be handled.

### TextureMat

This shader is written to handle textured, material objects. The maximum number of lighsources is 16. You can add more, but the surplus will not be handled.

### TextureColor

This shader is written to handle textured, colored objects. It doesn't support materials or lightsources.

### Point

This shader is written to handle point objects. It doesn't support materials, textures or light sources, but it supports colors and point size.
