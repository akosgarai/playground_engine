# Mesh

It contains everything that we need for drawing a stuff. Now i have 6 kind of meshes above the base one.

## Base mesh

Those parameters that are always needed for managing an object.

- **Vertices** - The verticies of the mesh.
- **vbo** - vertex buffer object.
- **vao** - vertex array object.
- **position** - The center position of the mesh. The model transformation is calculated based on this.
- **direction** - The mesh is moving to this direction. If this value is null vector, then the mes is not moving.
- **velocity** - The mesh is moving to the direction with this speed. if this value is null, then the mash is not moving.
- **yaw** - The rotation angle of the mesh on the 'Y' axis.
- **pitch** - The rotation angle of the mesh on the 'X' axis.
- **roll** - The rotation angle of the mesh on the 'Z' axis.
- **scale** - The mesh is scaled by this vector. The model transformation is calculated based on this.
- **wrapper** - The glwrapper, that we can use for calling gl functions.
- **parent** - The parent mesh of the current one.
- **parentSet** - If the mesh has a parent, set with the relevant function, this flag is updated.
- **bo** - The parameters for the bounding object.
- **boundingObjectSet** - returns true, if the object is set.

It has setter functions for the parameters, and getters for the necessary ones, also one for the model transformation matrix calculation and one for updating the state of the mesh.

## Textured mesh

It is a mesh extension for textured objects. It's parameter list is extended with the followings:

- **Indices** - In the Draw function the gl.DrawElements function is used, so that i have to maintain a buffer for the indices. These are the values that i can pass to the buffer.
- **Textures** - The textures that are used for covering the drew mesh.
- **ebo** - The element buffer object identifier. the indices are stored here.

Its `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with triangles, and then cleans up. The `NewTexturedMesh` function returns a textured mesh.

## Material mesh

It is a mesh extension for material objects. It's parameter list is extended with the followings:

- **Indices** - In the Draw function the gl.DrawElements function is used, so that i have to maintain a buffer for the indices. These are the values that i can pass to the buffer.
- **Material** - The material that is used for calculating the color of the mesh.
- **ebo** - The element buffer object identifier. the indices are stored here.

Its `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with triangles, and then cleans up. The `NewMaterialMesh` function returns a material mesh.

## Point mesh

It is a mesh extension for point objects. Its parameter list isn't extended, but it has an `Add` function for extending the mesh. It's necessary, because the `NewPointMesh` function returns an empty mesh, so that we have to fill it with the vertices one by one. Its `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with points, and then cleans up.

## Color mesh

It is a mesh extension for colored objects. The colors are set in the vertices. Its parameter list is exended with the followings:

- **Indices** - In the Draw function the gl.DrawElements function is used, so that i have to maintain a buffer for the indices. These are the values that i can pass to the buffer.
- **Color** - The colors that was used to generate the color of the object. It is used in the export process.
- **ebo** - The element buffer object identifier. the indices are stored here.

Its `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with triangles and then cleans up. The `NewColorMesh` returns a color mesh.

## Textured color mesh

It is a mesh extension for textured object where the object color also counts. Its parameter list is extended with the followings:

- **Indices** - In the Draw function the gl.DrawElements function is used, so that i have to maintain a buffer for the indices. These are the values that i can pass to the buffer.
- **Textures** - The textures that are used for covering the drew mesh.
- **Color** - The colors that was used to generate the color of the object. It is used in the export process.
- **ebo** - The element buffer object identifier. the indices are stored here.

Its `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with triangles, and then cleans up. The `NewTexturedColoredMesh` function returns a textured colored mesh.

## Textured material mesh

It is a mesh extension for textured object where the object material also counts. Its parameter list is extended with the followings:

- **Indices** - In the Draw function the gl.DrawElements function is used, so that i have to maintain a buffer for the indices. These are the values that i can pass to the buffer.
- **Textures** - The textures that are used for covering the drew mesh.
- **Material** - The material that is used for calculating the color of the mesh.
- **ebo** - The element buffer object identifier. the indices are stored here.

Its `Draw` function gets the Shader as input. It makes the uniform setup, buffer bindings, draws with triangles, and then cleans up. The `NewTexturedMaterialMesh` function returns a textured material mesh.
