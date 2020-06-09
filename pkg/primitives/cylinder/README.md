# Cylinder

This package represents a cylinder. The inputs are the radius (float32) of the circle, the precision (int) of the circle, and the length (float32) of the object.
The math behind the calculation is based on [this](http://www.songho.ca/opengl/gl_cylinder.html) very good explanation. The bounding object of the cylinder (and also of the half circle based cylinder) is AABB.

## Circle based cylinder

This kind of cylinder is circle based. It means, that the base shape is a circle, that could be other kind of shape, depending of the precision.

## Half circle based cylinder

This kind of cylinder is half circle based.

## New

New function returns a cylinder.

- `rad` - the radius of the circle.
- `prec` - the precision of the circle.
- `length` - the length of the body of the cylinder.

## NewHalfCircleBased

NewHalfCircleBased function returns a half circle based cylinder.
- `rad` - the radius of the circle.
- `prec` - the precision of the circle.
- `length` - the length of the body of the cylinder.

## MaterialMeshInput

MaterialMeshInput method returns the vertices, indices, bounding object (AABB) inputs for the NewMaterialMesh function.

## ColorMeshInput

ColorMeshInput method returns the vertices, indices, bounding object (AABB) inputs for the NewColorMesh function.

## TexturedMeshInput

TexturedMeshInput method returns the vertices, indices, bounding object (AABB) inputs for the NewTexturedMesh function.
