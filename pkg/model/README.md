# Model

The purpose of this package, to gather the mashes that are connected to the same object (eg a composite object - a lamp with pole, and bulb). The meshes of a model are moving together, they are rotating, in the same time. The model contains a transparency flag, that can be used to prevent the early drawing.
The base model has been extended with collision detection support. Now it can return a nearest mesh and its distance from a given point. The `Clear` function deletes the current meshes from the model.

## Bug model

This predefined model represents a "composite object". This model contains 4 material squere, one for the bottom, one for the body and 2 for the eyes. It could contaion attach points (point meshes) and wings. It has getter functions for the body parts. The wings are moving up-down direction.

### BugBuilder

This is a builder for the Bug model.

- `position` - the position of the body mesh.
- `scale` - the scale of the body mesh.
- `wrapper` - the wrapper pkg (interfaces.GLWrapper) for the gl functions. It can be set with the `SetWrapper` function.
- `bodyMaterial` - the material of the body mesh.
- `bottomMaterial` - the material of the bottom mesh.
- `eyeMaterial` - the material of the eye meshes.
- `rotationX` - the rotation (x axis) of the pole mesh.
- `rotationY` - the rotation (y axis) of the pole mesh.
- `rotationZ` - the rotation (z axis) of the pole mesh.
- `spherePrecision` - the precision of the spheres.
- `lightAmbient` - the ambient color component of the lightsource.
- `lightDiffuse` - the diffuse color component of the lightsource.
- `lightSpecular` - the specular color component of the lightsource.
- `constantTerm` - the constant term of the lightsource.
- `linearTerm` - the linear term of the lightsource.
- `quadraticTerm` - the quadratic term of the lightsource.
- `withLight` - if this flag set, the bug will contain a point lightsource.
- `velocity` - the initial speed of the bug.
- `direction` - the direction of the bug.
- `movementRotationAngle` - the rotation angle of the movement.
- `movementRotationAxis` - the rotation axis of the movement.
- `sameDirectionTime` - the bug goes to the same direction this time (ms).
- `withWings` - if this flag set true, the bug will contain wing meshes.
- `wingStrikeTime` - the wings strike time (ms).

## StreetLamp model

This predefined model represents a "composite object". This model contains 2 cuboids or cylinders, 'Pole' and 'Top' and a sphere 'Bulb'. It has 2 implementation. One is from plane materials, and another one that is textured. In the material version, the pole and the top is cuboid, in the textured one, they are cylinders.

### StreetLampBuilder

This is a builder for the StreetLamp model.

- `position` - the bottom position of the pole mesh.
- `poleLength` - the length of the pole mesh.
- `rotationX` - the rotation (x axis) of the pole mesh.
- `rotationY` - the rotation (y axis) of the pole mesh.
- `rotationZ` - the rotation (z axis) of the pole mesh.
- `assetsBaseDir` - In case of textured room, we have to know where are the assets.
- `wrapper` the wrapper pkg (interfaces.GLWrapper) for the gl functions. It can be set with the `SetGlWrapper` function.
- `bulbMaterial` - the material of the bulb. It is also used as the color components of the spot lightsource.
- `constantTerm` - the constant term of the lightsource.
- `linearTerm` - the linear term of the lightsource.
- `quadraticTerm` - the quadratic term of the lightsource.
- `cutoff` - the cutoff parameter of the spot lightsource.
- `outerCutOff` - the outer cutoff parameter of the spot lightsource.
- `lampOn` - if this flag is on, the lamp will be turned on.

For modeling a spotlight, we need to know the position and the direction of the lightsource. The position could be calculated from the position and the length of the pole mesh, and it's rotations. The direction could be also calculated from the rotations.

## Room model

It is a cuboid object that contains a door like surface. The door could be opened and closed. It's movement is animated.

### RoomBuilder

This is a builder for the room model.

- `position` - the position of the room (center point of the floor mesh)
- `width` - the length of the usable area in the x axis
- `height` - the length of the usable area in the y axis
- `length` - the length of the usable area in the z axis
- `wallWidth` - the width of the walls
- `doorWidth` - the width of the door that is on the right side of the front wall.
- `doorHeight` - the height of the door that is on the right side of the front wall.
- `rotationX` - the rotation (x axis) of the floor mesh.
- `rotationY` - the rotation (y axis) of the floor mesh.
- `rotationZ` - the rotation (z axis) of the floor mesh.
- `assetsBaseDir` - In case of textured room, we have to know where are the assets.
- `frontWindow` - if this flag is set, the room will contain a window on the front wall.
- `backWindow` - if this flag is set, the room will contain a window on the back wall. currently this feature is unimplemented.
- `leftWindow` - if this flag is set, the room will contain a window on the left wall. currently this feature is unimplemented.
- `rightWindow` - if this flag is set, the room will contain a window on the right wall. currently this feature is unimplemented.
- `doorOpened` - if this flag is set, the front door will be in opened state.
- `windowWidth` - the width of the windows that we could set on the textured rooms.
- `windowHeight` - the height of the windows that we could set on the textured rooms.
- `wrapper` the wrapper pkg (interfaces.GLWrapper) for the gl functions. It can be set with the `SetGlWrapper` function.

The Builder provides material (BuildMaterial) or textured (BuildTexture) room solutions. The window feature only works with the textured rooms.

## Terrain model

This model could be used for generating terrain surfaces. The package provides a `TerrainBuilder`, for the terrain generation.

### TerrainBuilder

This is the builder for the terrain.

- `width` the width (integer) of the terrain (the length in the `X` axis) It can be updated with the `SetWidth` function.
- `length` the length (integer) of the terrain (the length in the `Z` axis) It can be updated with the `SetLength` function.
- `iterations` the number of iterations (integer) of the randomization process. It can be updated with the `SetIterations` function.
- `minH` the lowest possible height value (float32) in the map. It can be updated with the `SetMinHeight` function.
- `maxH` the maximum possible height value (float32) in the map. It can be updated with the `SetMaxHeight` function.
- `seed` this (int64) value is uses as seed for the random function. It can be updated with the `SetSeed` function.
- `minHIsDefault` if this flag is set true, the minH value is used as default height in the map, otherwise 0.0. By default it is false. It can be updated with the `MinHeightIsDefault` function.
- `peakProbability` this value (integer) represents the percentage of the chance of the generated height will be a peak. It can be updated with the `SetPeekProbability` function.
- `cliffProbability` this value (integer) represents the percentage of the chance of the generated height will be a cliff. It can be updated with the `SetCliffProbability` function.
- `wrapper` the wrapper pkg (interfaces.GLWrapper) for the gl functions. It can be set with the `SetGlWrapper` function.
- `tex` the texture container (texture.Textures) for the surface textures. It can be set to the grass textures with the `GrassTexture` function.
- `scale` the scale vector (mgl32.Vec3) of the terrain mesh. It can be updated with the `SetScale` function.
- `debugMode` is this flag is set true, useful information will be printed to the console. It can be set with the `SetDebugMode` function.

The `Build` function returns the `Terrain` model, that is generated with the given setup.

### Terrain

The Terrain represents the surface, the ground, whatever. It contains the `heightMap` that the builder generated and also the width, length, debugModes. The `HeightAtPos` function returns the height value in a given position. The calculation is based on a basic interpolation algorithm.

**The interpolation**

- If the given point is not above or below the surface, it returns error. Otherwise it is based on the following pseudo algorithm:

```
The interpolation algorithm:
Let wX = position.X() - int(position.X()). If wX less than 0, increase it with 1.0.
Let wZ = position.Z() - int(position.Z()). If wZ less than 0, increase it with 1.0.

Y+
^
|  A  B
|  ----
|  |  |
|  ----
|  D  C
|
-------->X+

Let heightAtTheGivenPosition = (heightA*(1-wX) + heightB*(wX)) * wZ + (heightD*(1-wX) + heightC*(wX)) * (1-wZ)
```

## Charset model

This model is useful for displaying texts on the screen.

### Glyph

The glyph is responsible for holding the parameters for a given glyph.

- Width - the width of the glyph in pixels.
- Height - the height of the glyph in pixels.
- BearingX - the offset of the glyph in pixels on the `X` axis from the origin.
- BearingY - the offset of the glyph in pixels on the `Y` axis from the origin.
- Advance - the offset of the next origin in pixels on the `X` axis from the current origin.
- tex - the texture of the glyph.
- Debug - if it is true, it prints out debug informations.

A glyph could be setup with its `Build` function.

### Charset

The Charset model is responsible for holdinng the glyphs for a given charset. For the file loading and initialization it provides 2 functions. The `LoadCharsetDebug` is for debugging, it prints out a bunch of useful stuff, the `LoadCharset` is for silent load. It has a `PrintTo` function that puts the text to the given surface. For the tests the [`Desyrel`](https://www.dafont.com/desyrel.font) fonts are used.

## Form items

For the form screen, we need a couple of models. The form items provide 4 kind of width sizes.

- `Full` - It is the longest width. It fits to the screens full width.
- `Half` - Half width. It fits to the screens half width.
- `Long` - 2/3 width. It fits to the screens 2/3 width.
- `Short` - 1/3 width. It fits to the screens 1/3 width.

```
---------
|       | - full width
---------
|   |   | - half - half width
---------
|  |    | - short - long width
---------
|    |  | - long - short width
---------
|   ||  | - half - short width
---------
|  ||   | - short - half width
---------
```

**Height**

The full width is given as input. The height (or length, if we talk about the xz plane) of the input boxes is calculated with the following function:

```
height = full / 1.96
```

### Form base

This is the base model of the form items. It calculates the sizes of the items. The calculation is based on the width of the drawable screen. It also create the surface mesh.

### Form item bool

This model represents a form item for maintaining a bool value. The texture was downloaded from [here](https://pixabay.com/hu/vectors/lekerek%C3%ADtett-v%C3%B6r%C3%B6s-t%C3%BCnetek-vezetett-26377/).

### Form item int

This model represents a form item for maintaining integer values. The input is handled with different states.

- Positive state (**P**) - The input field is empty.
- Positive state (**P0**) - The input field contains 0 as input value.
- Positive integer state (**PI**) - The input field contains positive integer value.
- Negative state (**N**) - The input field contains the '-' token.
- Negative integer state (**NI**) - The input field contains negativ integer value.

**Transitions**

current state -next token-> next state

```
P       --(i)-->    PI
P       --(0)-->    P0
P       --(-)-->    N
PI      -(i/0)->    PI
N       --(i)-->    NI
NI      -(i/0)->    NI
```

### Form item float

This model represents a form item for maintaining float values. The input is handled with different states.

- Positive state (**P**) - The input field is empty.
- Positive zero state (**P0**) - The input field contains a '0' value the only valid token after this state is the '.'
- Positive integer state (**PI**) - The input field contains positive integer value.
- Positive dot state (**P.**) - the '.' character is pressed from a negative state.
- Positive float state (**PF**) - the next state after the dot.
- Negative state (**N**) - The input field contains the '-' token.
- Negative zero state (**N0**) - The input field contains the '-0' tokens. the only valid token after this is the '.'
- Negative integer state (**NI**) - The input field contains negativ integer value.
- Negative dot state (**N.**) - the '.' character is pressed from a negative state.
- Negative float state (**NF**) - the next state after the dot.

**Transitions**

current state -next token-> next state

```
P       --(0)-->    P0
P       --(i)-->    PI
P0      --(.)-->    P.
PI      -(i/0)->    PI
PI      --(.)-->    P.
P.      -(i/0)->    PF
PF      -(i/0)->    PF
P       --(-)-->    N
N       --(0)-->    N0
N       --(i)-->    NI
N0      --(.)-->    N.
NI      -(i/0)->    NI
NI      --(.)-->    N.
N.      -(i/0)->    NF
NF      -(i/0)->    NF
```

### Form item text

This model represents a form item for maintaining text values.

### Form item int64

This model represents a form item for maintaining int64 values. The input is handled with different states.

- Positive state (**P**) - The input field is empty.
- Positive integer state (**PI**) - The input field contains positive integer value.
- Negative state (**N**) - The input field contains the '-' token.
- Negative integer state (**NI**) - The input field contains negativ integer value.

**Transitions**

current state -next token-> next state

```
P       --(i)-->    PI
P       --(-)-->    N
PI      -(i/0)->    PI
N       --(i)-->    NI
NI      -(i/0)->    NI
```
