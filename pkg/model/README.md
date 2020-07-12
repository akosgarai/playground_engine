# Model

The purpose of this package, to gather the mashes that are connected to the same object (eg a composite object - a lamp with pole, and bulb). The meshes of a model are moving together, they are rotating, in the same time. The model contains a transparency flag, that can be used to prevent the early drawing.
The base model has been extended with collision detection support. Now it can return a nearest mesh and its distance from a given point. The `Clear` function deletes the current meshes from the model.

## Bug model

This predefined model represents a "composite object". This model contains 4 material squere, one for the bottom, one for the body and 2 for the eyes. It has getter functions for the center point of the body parts. It's initial position, and scale are set during it's construction with the NewBug function.

## StreetLamp model

This predefined model represents a "composite object". This model contains 2 cuboids or cylinders, 'Pole' and 'Top' and a sphere 'Bulb'. It has 2 implementation. One is from plane materials, and another one that is textured. In the material version, the pole and the top is cuboid, in the textured one, they are cylinders.

## Room model

This predefined model represents a "composite object". This model contains 8 cuboids, floor, ceiling, 3 full wall, 1 door, a half wall next to the door, and a half wall above the door. It has 2 implementation. One is from plane materials and another one that is textured.

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

For the form screen, we need a couple of models.

### Form item bool

This model represents a form item for maintaining a bool value. The texture was downloaded from [here](https://pixabay.com/hu/vectors/lekerek%C3%ADtett-v%C3%B6r%C3%B6s-t%C3%BCnetek-vezetett-26377/).

### Form item int

This model represents a form item for maintaining integer values.

### Form item float

This model represents a form item for maintaining float values. The input is handled with different states.

- Positive state (**P**) - The input field is empty. (valueInt = 0, valueFloat = 0, floatPosition = -1, isNegative = false)
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
