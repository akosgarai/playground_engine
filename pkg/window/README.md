# Window

This package is written for the glfw intialization and window creation. It provides a Builder tool for configuring the window.

## WindowBuilder

This tool can be used for creating windows. You can get a builder instance with calling the `NewWindowBuilder` function, that calls glfw.Init function. The builder provides functions for getting some data about the current monitor.

- Resolution (width, height) in pixels with the `GetCurrentMonitorResolution` function.
- `GetCurrentMonitorPhysicalSize` returns the size, in millimeters, of the display area of the monitor.
- `GetCurrentMonitorContentScale` function retrieves the content scale for the specified monitor. The content scale is the ratio between the current DPI and the platform's default DPI. If you scale all pixel dimensions by this scale then your content should appear at an appropriate size. This is especially important for text and any UI elements.
- `GetCurrentMonitorWorkarea` returns the position, in screen coordinates, of the upper-left corner of the work area of the specified monitor along with the work area size in screen coordinates. The work area is defined as the area of the monitor not occluded by the operating system task bar where present. If no task bar exists then the work area is the monitor resolution in screen coordinates.

## InitGlfw

This function calls glfw.Init, sets up the necessary windowhints, creates a new window, makes it the current context and returns it.

## DummyKeyCallback

This function is a dummy key callback, that could be attached to the window. It only prints the inputs to the console.

## DummyMouseButtonCallback

This function is a dummy mouse button callback, that could be attached to the window. It only prints out the inputs to the console.
