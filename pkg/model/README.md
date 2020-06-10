# Model

The purpose of this package, to gather the mashes that are connected to the same object (eg a composite object - a lamp with pole, and bulb). The meshes of a model are moving together, they are rotating, in the same time.

## Bug model

This predefined model represents a "composite object". This model contains 4 material squere, one for the bottom, one for the body and 2 for the eyes. It has getter functions for the center point of the body parts. It's initial position, and scale are set during it's construction with the NewBug function.

## StreetLamp model

This predefined model represents a "composite object". This model contains 2 cuboids or cylinders, 'Pole' and 'Top' and a sphere 'Bulb'. It has 2 implementation. One is from plane materials, and another one that is textured. In the material version, the pole and the top is cuboid, in the textured one, they are cylinders.

## Room model

This predefined model represents a "composite object". This model contains 8 cuboids, floor, ceiling, 3 full wall, 1 door, a half wall next to the door, and a half wall above the door. It has 2 implementation. One is from plane materials and another one that is textured.
