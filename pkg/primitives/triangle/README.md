# Triangle

It represents a triangle, that is described by its points (3) and its normal vector.

## New

It returns a triangle. The input is the angle vector (Vec3). The longest side is 1 unit long, so that the side lengts are scaled based on this. The origo is the center point of the longest side. The triangle's plane is the X-Y plane.

Math behind the calculations:

- sine rule:

```
a : b : c = sin(alpha) : sin (beta) : sin (gamma)
```

## ColoredMeshInput

ColoredMeshInput method returns the vertices, indices, bounding object (AABB) - inputs for the New Mesh function.
