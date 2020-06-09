# Sphere

This package represents a unit sphere, that is described with the points of its surface. The points, and the indices are calculated based in the precision value.
The math behind the calculation is based on [this](http://www.songho.ca/opengl/gl_sphere.html) very good explanation. 

## New

This function returns the unit sphere. The precision is the only input of this function.

## MaterialMeshInput

MaterialMeshInput method returns the vertices, indices, bounding object (Sphere) - inputs for the NewMaterialMesh function.

## ColoredMeshInput

ColoredMeshInput method returns the vertices, indices, bounding object (Sphere) - inputs for the NewColorMesh function.

## TexturedMeshInput

TexturedMeshInput method returns the vertices, indices, bounding object (Sphere) - inputs for the NewTexturedMesh function.
