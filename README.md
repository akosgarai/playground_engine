# Playgorund Engine

![Tests](https://github.com/akosgarai/playground_engine/workflows/Tests/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/akosgarai/playground_engine/badge.svg?branch=master)](https://coveralls.io/github/akosgarai/playground_engine?branch=master)

This engine is written for playing around with the opengle v4.x features. It contains basic primitives, supports camera, light sources, materials, textures, (wavefront) model import / export solutions.

## CLI

This package provides a cli tool for downloading the necessary assets or shaders for an application.

### Usage

First of all, You have to install the cli-application.

```
go get github.com/akosgarai/playground_engine
cd $GOPATH/src/github.com/akosgarai/playground_engine
go install
```

After this step, You should be able to use the cli-tool. It's easy to use, it gets a command and a parameter as arguments. Without arguments it prints out the help menu. The most important command is the install, that You can use for downloading model textures or shaders.

The following command creates a `shaders` directory (if not exists) and downloads the shader apps into the `shaders` dir.

```
playground_engine install shaders
```

This command creates an `assets` directory (if not exists) and downloads the assest files into this directory.

```
playground_engine install models
```

## Sample

![Sample gif from outer space](https://github.com/akosgarai/go_opengl_playground/blob/master/examples/07-textured-spheres/sample/sample.gif)

## Useful links

The following tutorials / documentations / explanations were the basics of the engine.

- [Godoc glfw](https://godoc.org/github.com/go-gl/glfw/v3.3/glfw)
- [Godoc mgl32](https://godoc.org/github.com/go-gl/mathgl/mgl32)
- [Godoc gl](https://godoc.org/github.com/go-gl/gl/v4.1-core/gl)
- [Learnopengl](https://learnopengl.com/) - good explanations and cpp examples.
- [About glsl](https://www.khronos.org/opengl/wiki/OpenGL_Shading_Language)
- A tutorial [first part](https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl) and [second part](https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-2-drawing-the-game-board)
- [Other tutorial](https://medium.com/@drgomesp/opengl-and-golang-getting-started-abcd3d96f3db)
- [About transformations](http://www.codinglabs.net/article_world_view_projection_matrix.aspx)

## Possible issues ubuntu.

- Opengl version mismatch.

The applications are using the opengl 4.1 package. If your version is same or higher, the appliactions should run without issues.
To check your opengl version just run the following command in terminal (based on [this](https://askubuntu.com/questions/47062/what-is-terminal-command-that-can-show-opengl-version)):

```bash
glxinfo | grep "OpenGL version"
```

The output is something like: `OpenGL version string: 4.6.0 NVIDIA 440.82`.

## I want to use different gl version.

In this case, you have to modify the wrapper package. The gl lib is included there. If you updated it (eg to `v3.3-core`), you have to update the `GL_MAJOR_VERSION`, `GL_MINOR_VERSION` version constants also. Unfortunately updating the shaders is a manual step. The versions are hardcoded to the shader applications.
