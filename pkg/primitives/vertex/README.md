# Vertex

A `Vertex` is the container of the possible values that we could use as vertex buffer data. Now it has support for the following options:

- **Position**: It stores the position coordinate. (Where is the vertex)
- **Normal**: The normal vector of the surface in the given position.
- **TexCoords**: If we use textures, we use it for storing the texture coordinates.
- **Color**: The color of the surface in the given position. We can use it instead of textures.
- **PointSize**: If we draw points, we can modify it's size with this float number.

# Vertices

In general we want to draw more than one vertex. The `Vertices` structure gives help for this need.
It has an `Add` method, that we can use for adding new items to the vertices. It's useful, if our application draws points.
The `Get` method returns the `[]float32` that we can use as vertex array object data. We can manage the result set with the function input.
Modes:
- `POSITION_NORMAL` (1) the position and the normal vectors are returned in this order.
- `POSITION_NORMAL_TEXCOORD` (2) the position, normal and the tex coords vectors are returned in this order.
- `POSITION_COLOR_SIZE` (3) the position, the color and the point size are retuned in this order.
- `POSITION_COLOR` (4) the position and the color vectors are returned in this order.
- `POSITION_COLOR_TEXCOORD` (5) the position, the color and the tex coords vectors are returned in this order.
