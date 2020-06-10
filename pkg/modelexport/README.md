# Model Export

This package is responsible for exporting models. It creates wavefront object and material files and duplicates the textures. All this stuff is moved to a common directory. The generated files fupposed to be valid, and could be opened with [Blender](https://www.blender.org/) application.

This package contains `Mtl` struct, that holds the parameters for the materials. It also contains an `Obj` struct, that holds the mesh parameters. The `Export` struct holds the original meshes, and the `Mtl`, `Obj` structures, and some internal variables.

## The export process.

The `New` function gets the meshes as inputs, setups an `Export` and returns it. The process starts when we call the `Export` function. It gets the destination directory as input. First it checks that the given directory exists. If not, it returns error. Then it iterates over the meshes and setups the `Obj` and `Mtl` structures based on the meshes. Currently it can export the following meshes: `ColorMesh`, `MaterialMesh`, `TexturedMesh`, `TexturedColoredMesh`, `TexturedMaterialMesh`, `PointMesh`. The last step is the output generation, it writes the data to files in the given directory, and also creates a copy from the textures (if we used them).
