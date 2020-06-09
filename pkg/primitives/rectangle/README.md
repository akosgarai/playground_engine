# Rectangle

This package represents a rectangle, that is described with 4 points and its normal vector. The structure also contains indices, that is used for populating the element array buffer and boungind object, that is also a rectangle (cuboid with 0 height). 

## New

New creates a rectangle with origo as middle point. The normal points to -Y. The longest side is scaled to one, and the same downscale is done with the other edge.

- `width` represents the length on the `X` axis.
- `height` represents the length on the `Z` axis.

The math for the calculation:

```
ratio = width / length
ratio == 1 => return NewSquare.
ratio > 1 => width is the longer -> X [-0.5, 0.5], Y [-1/(ratio*2), 1/(ratio*2)].
ratio < 1 => length is the longer -> X [-ratio/2, ratio/2], Y [-0.5, 0.5].
```

## NewSquare

NewSquare creates a rectangle with origo as middle point. Each side is 1 unit long, and it's plane is the X-Z plane. The normal points to -Y.

## MeshInput

MeshInput method returns the vertices, indices, bounding object (AABB) - inputs for the New Mesh function.

## ColoredMeshInput

ColoredMeshInput method returns the vertices, indices, bounding object (AABB) - inputs for the New Mesh function.

## TexturedColoredMeshInput

This method returns the vertices, indices, bounding object (AABB) - inputs for the NewTexturedColoredMesh function.
