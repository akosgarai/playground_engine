# Store package

This package represent a key - value store.

It provides Set and Get methods.

## Implemented stores

The following kind of stores are implemented:

- `GlfwKeyStore` stores the keyboard button states. It maps glfw.Keys to bool values, that tells us, that the key is pressed or not.
- `GlfwMouseStore` stores the mouse button states. It maps glfw.MouseButtons to bool values, that tells us that the button is pressed or not.
