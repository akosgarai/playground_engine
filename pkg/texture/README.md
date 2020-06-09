# Texture

This package contains the texture related data
- **Id** The id of the texture. eg: `gl.TEXTURE0`
- **UniformName** The Uniform name of the texture
- **TextureName** The generated name that was given by the GenTextures command
- **TargetId** The target that we use for BindTexture. eg: `TEXTURE_2D`
- **Wrapper** The glwrapper, that can be used for calling gl functions.
- **FilePath** The path of the file that is used for the texture. Good for the export process.

We can bind or unbind the textures with the `Bind` or `UnBind` methods.

## Textures

It contains Texture objects. Its `AddTexture` method creates a new Texture and adds it to itself. The `UnBind` helper method calls UnBind on each texture that it contains.
