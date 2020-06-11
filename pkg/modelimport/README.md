# Model import

This package is responsible for importing models / meshes from wavefront object files. For the fileparsing process, it uses the [gwob](https://github.com/udhos/gwob) package.

## Import process

We can create an `Import` structure with the basic setup with the `New` function. It gets the directory and the object file name and the gl wrapper as input. The import process could be started with the `Import` function. Under the hood, it setups the gwob, reads the files, and creates the meshes from the gwob structures. Finally, we can get the meshes with the `GetMeshes` function.
