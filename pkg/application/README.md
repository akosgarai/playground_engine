# Application package

The common application related stuff goes here. It holds

- a `window`, that has to implement the `Window` interface
- mouse position `MousePosX`, `MousePosY`, set with the mouse button callback
- `mouseDowns`, that stores the mouse button states.
- `keyDowns`, that stores the keyboard button states.
- `screens`, the screens added to this application
- `activeScreen`, the screen that currently used
- `menuScreen`, the screen that is connected to the application menu.
- `menuSet`, this flag is true if the menuScreen has been set.
- `wrapper`, the interface for calling gl commands.

## New

This function returns an initialized application instance. It's only input is the wrapper.

## Log

Log returns the string representation of this object.

## SetWindow

SetWindow updates the window with the new one.

## GetWindow

GetWindow returns the current window of the application.

## GetCamera

GetCamera returns the current camera of the application.

## GetClosestModelMeshDistance

This function calls GetClosestModelMeshDistance on the activeScreen.

## Update

Update calls Update function on the activeScreen.

## AddScreen

AddScreen appends a screen to the screens.

## ActivateScreen

ActivateScreen sets the given screen to active screen

## MenuScreen

MenuScreen sets the given screen to menu screen. It also sets the menuSet variable to true.

## Draw

Draw calls Draw function on the activeScreen.

## KeyCallback

KeyCallback is responsible for the keyboard event handling.

## MouseButtonCallback

MouseButtonCallback is responsible for the mouse button event handling.

## SetKeyState

SetKeyState setups the keyDowns based on the key and action

## SetButtonState

SetButtonState setups the mouseDowns based on the key and action

## GetMouseButtonState

GetMouseButtonState returns the state of the given button

## GetKeyState

GetKeyState returns the state of the given key

## SetUniformFloat

SetUniformFloat sets the given float value to the given string key in the uniformFloat map.

# SetUniformVector

SetUniformVector sets the given mgl32.Vec3 value to the given string key in the uniformVector map.
