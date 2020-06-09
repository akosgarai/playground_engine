# BoundingObject

This package contains the necessary setup for initializing a bounding object. It holds a type name (`AABB` or `Sphere`) and the params, that are necessary to know. For the `AABB`, we need to have a `width` a `height` and a `length` key in the params. For the `Sphere` We only need to have a `radius` key. The package provides getter functions for the parameters. The **Type()** returns the typeName, and the **Params()** returns the params map.

## What is a bounding object

This is used for collision detection. Let's imagine that we have a tetrahedron in our 3D world. If we want to know that an other object collided with the tetrahedron or not, we have to calculate intersection on the tetrahedron. But what if we use bounding object. For example for the tetrahedron, we could use cube. What does it mean? When we want to check the intersection for the terahedron, we check the intersection against the cube. Calculating with the cube is much easier than with the tetrahedron. So this is bounding object for. Making the coillision detection easier for us. But this is a bounding object, not the object itself, so this calculation is just an approximation, but it's enough precise for us.

## Further reading

- [learnopengle.com](https://learnopengl.com/In-Practice/2D-Game/Collisions/Collision-detection) tutorial about the collision detection in 2D
- [3D collision detection](https://developer.mozilla.org/en-US/docs/Games/Techniques/3D_collision_detection) & bounding boxes & spheres.
