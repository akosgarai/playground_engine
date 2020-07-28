# Screen package

This package is written for handling different needs in the same application. For example we need a menu screen without camera and a world screen with a camera in the same application. This package provides solution for this need.
The models and shaders are handled in the screen package also the camera, light, movement maps, closest objects and the screen level uniforms.

It holds:
- a `camera`, that has to implement the `Camera` interface
- `shaderMap`, that makes the connection between the models and the shaders
- `directionalLightSources`, for storing the directional lights.
- `pointLightSources`, for storing the point lights.
- `spotLightSources`, for storing the spot lights.
- `cameraKeyboardMovementMap`, makes connection between the keyboard buttons and the camera state updates.
- `rotateOnEdgeDistance`, for the mouse rotations.
- `uniformFloat`, for storing the float uniforms that needs to be set for every shader.
- `uniformVector`, for storing the vector uniforms that needs to be set for every shader.
- `setupFunction`, for gl setup, like Enable stuff, setup clear color, etc.

## Functions

**New**

This function returns an initialized screen instance.

**Log**

Log returns the string representation of this object.

**SetCameraMovementMap**

SetCameraMovementMap sets the cameraKeyboardMovementMap variable. Currently the following values are supported: `forward`, `back`, `left`, `right`, `up`, `down`, `rotateLeft`, `rotateRight`, `rotateUp`, `rotateDown`

**SetRotateOnEdgeDistance**

SetRotateOnEdgeDistance updates the rotateOnEdgeDistance variable. The value has to be in the [0-1] interval. If not, a message is printed to the console and the variable update is skipped.

**SetCamera**

SetCamera updates the camera with the new one.

**GetCamera**

GetCamera returns the current camera of the screen.

**AddShader**

AddShader method inserts the new shader to the shaderMap

**AddModelToShader**

AddModelToShader attaches the model to a shader.

**GetClosestModelMeshDistance**

GetClosestModelMeshDistance returns the closest model, mesh and its distance from the mouse position.

**SetUniformFloat**

SetUniformFloat sets the given float value to the given string key in the uniformFloat map.

**SetUniformVector**

SetUniformVector sets the given mgl32.Vec3 value to the given string key in the uniformVector map.

**AddDirectionalLightSource**

AddDirectionalLightSource sets up a directional light source. It takes a DirectionalLight input that contains the model related info, and it also takes a [4]string, with the uniform names that are used in the shader applications the `DirectionUniformName`, `AmbientUniformName`, `DiffuseUniformName`, `SpecularUniformName`. They has to be in this order.

**AddPointLightSource**

AddPointLightSource sets up a point light source. It takes a PointLight input that contains the model related info, and it also containt the uniform names in [7]string format. The order has to be the following: `PositionUniformName`, `AmbientUniformName`, `DiffuseUniformName`, `SpecularUniformName`, `ConstantTermUniformName`, `LinearTermUniformName`, `QuadraticTermUniformName`.

**AddSpotLightSource**

AddSpotLightSource sets up a spot light source. It takes a SpotLight input that contains the model related info, and it also contains the uniform names in [10]string format. The order has to be the following: `PositionUniformName`, `DirectionUniformName`, `AmbientUniformName`, `DiffuseUniformName`, `SpecularUniformName`, `ConstantTermUniformName`, `LinearTermUniformName`, `QuadraticTermUniformName`, `CutoffUniformName`.

**Draw**

Draw calls Draw function in every drawable item. It calls the setupFunction, then it loops on the shaderMap (shaders). For each shader, first set it to used state, setup camera realted uniforms, then setup light related uniformsi and custom uniforms. Then we can pass the shader to the Model for drawing.

**Update**

Update loops on the shaderMap, and calls Update function on every Model. It also handles the camera movement and rotation, if the camera is set.

**Export**

Export creates a directory for the screen and calls Export function on the models.

## Screens

Some screens are provided by the engine.

### ScreenWithFrame

This screen extension is a base screen, that contains a frame on its edge. It supports a label text on the top of the screen. The frame isn't set behind the label. It also provides an option for activating a surface on the bottom of the screen. This surface could be used for displaying text. A builder is provided for creating this screen. This screen holds the followings:

- frameWidth - The width of the the visible screen. It is the width (length on the `X` axis) of the frames.
- frameLength - The length (on the `Y` axis) of the frames.
- detailContentBox - This is the mesh that could be used for displaying the text.
- detailContentBoxHeight - This is the height of the detailContentBox mesh.

It provides a function (GetFullWidth) for getting the width of the usable area.

### ScreenWithFrameBuilder

This is a handy tool for creating ScreenWithFrames. It provides Set... and Build functions.

- `SetWrapper` function sets the wrapper for further usage.
- `SetWindowSize` function sets the windowWidth and windowHeight variables.
- `SetLabelWidth` function could be used for creating the frame padding to the label text.
- `SetFrameSize` function sets the frameWidth, frameLength, and frameTopLeftWidth variables.
- `SetFrameMaterial` function sets the material of the frames.
- `SetDetailContentBoxMaterial` function sets the material of the detailContentBox.
- `SetDetailContentBoxHeight` function sets the height of the detailContentBox. If this value is 0, the building step of this box will be skipped.
- `Build` function returns a ScreenWithFrame with the given setup. If the detailContentBoxHeight is 0, then the detailContentBox will be set to `nil`.

### MenuScreen

This screen is for displaying menus. It is a `ScreenWithFrame` extension. It holds the followings:

- charset - the charset model for displaying text.
- background - the background model of the screen.
- surfaceTexture - the texture that we use for covering the menu items.
- defaultMaterial - the default material of the menu items.
- hoverMaterial - the highlight material of the menu items.
- options - the list of the options that could be displayed.
- backgroundShader - the shader that is used for rendering the menu items.
- fontShader - the shader that is used for rendering the texts.
- fontColor - the color of the glyphs.
- backgroundColor - the clear color of the screen.
- state - the current state of the screen.
- surfaceToOption - it maps the surface meshes to the options.
- maxScrollOffset - in case of long menu, the scrolling is supported. It limits the scrolling.
- currentScrollOffset - The current offset of the form items in the Y axis. (move the screen with up / down cursors.)

```
 _______________
 |	    - 0.1
 | text     - 0.2 - Continue
 |	    - 0.2
 | text     - 0.2 - Start / Restart
 |	    - 0.2
 | text     - 0.2 - Options
 |	    - 0.2
 | text     - 0.2 - Save state
 |	    - 0.2
 | text     - 0.2 - Exit
 |	    - 0.1
```

**Features**

If the mouse cursor is above a form item, it triggers a hover effect, that does the followings:

- Updates the menu item material to the Highlight value.

This screen supports scrolling on the `Y` axis. This feature is active, if we have more form items that we can display on the visible area. The scolling can be triggered with the `up` and `down` arrow keys.

Clicking on a hovered item.

- It triggers a clickEvent function call on the option.

#### MenuScreenBuilder

This tool is provided for creating menus. It is a `ScreenWithFrameBuilder` extension. The following variables could be set during the construction.

- charset - the charset model for displaying text.
- menuItemSurfaceTexture - the texture that we use for covering the menu items.
- defaultMaterial - the default material of the menu items.
- hoverMaterial - the highlight material of the menu items.
- options - the list of the options that could be displayed.
- menuItemFontColor - the color of the glyphs.
- backgroundColor - the clear color of the screen.
- state - the initial state of the screen.

Set... functions are provided for the `charset`, `menuItemSurfaceTexture`, `defaultMaterial`, `hoverMaterial`, `options`, `menuItemFontColor`, `backgroundColor`, `state`.

### FormScreen

This screen is for displaying forms, like a settings page. It is a `ScreenWithFrame` extension. It holds the followings:

- charset - the charset model for displaying text.
- formItemShader - the shader that is used for rendering the form items.
- sinceLastClick - the time since the last click event.
- sinceLastDelete - the time since the last character deletion event.
- underEdit - the character based form item, that is currently edited.
- maxScrollOffset - The maximum offset of the scrolling. It is calculated from the lengths of the screen and the form items.
- currentScrollOffset - The current offset of the form items in the Y axis. (move the screen with up / down cursors.)
- formItemToConf - It maps the FormItems to ConfigItems. It is used to sync the values.
- formItemLabelColor - It is the color of the form item labels.
- formItemInputColor - It is the color of the form item inputs.

**Features**

If the mouse cursor is above a form item, it triggers a hover effect, that does the followings:

- Updates the form item material to the Highlight value.
- Displays the form item description in the detailContentBox.

This screen supports scrolling on the `Y` axis. This feature is active, if we have more form items that we can display on the visible area. The scolling can be triggered with the `up` and `down` arrow keys.

Clicking on a hovered item.

- On a bool input, it negates its value.
- On a character based form item, it activates its edit mode, a cursor appeares at the end of the current value.

Editing the value of a character based form item. Aka character callback & backspace key handler.

- A new character can be added to the end of the current value.
- The last character can be deleted.

#### FormScreenBuilder

This tool is provided for creating forms. It is a `ScreenWithFrameBuilder` extension. The following variables could be set during the construction.

- headerLabel - This is displayed on the top left of the screen.
- config - The form items are based on this configuration.
- configOrder - The order of the form items.
- charset - The charset model that will be used for text writing.
- lastItemState - It is the state of the latest inserted form item.
- offsetY - The y component of the position of the last inserted item.
- headerLabelColor - The color of the header label text.
- formItemLabelColor - It is the color of the form item labels.
- formItemInputColor - It is the color of the form item inputs.

Set... functions are provided for the `headerLabel`, `windowHeight`, `config`, `configOrder`, `charset`, `headerLabelColor`, `formItemLabelColor`, `formItemInputColor`. The `lastItemState`, and the `offsetY` is used and maintained during the process of the form item building from the config items.

For displaying stuff, it uses the following system:

- `Full` width items for text / vector inputs. It is the longest width. It fits to the screens full width.
- `Half` width items for int / float inputs. Half width. It fits to the screens half width.
- `Long` width items for int64 inputs. 2/3 width. It fits to the screens 2/3 width.
- `Short` width items for bool inputs. 1/3 width. It fits to the screens 1/3 width.

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
|  | |  | - short - short - short width
---------
```

**8 Possible position**:
- Full item position (**F**)
- Left Half item position (**LH**)
- Right Half item position (**RH**)
- Left Short item position (**LS**)
- Middle Short item position (**MS**)
- Right Short item position (**RS**)
- Left Long item position (**LL**)
- Right Long item position (**RL**)

State machine for position step:
- Initial state: **F**
- Next state is calculated based on the current state & the current item width.

```
F   -(Full)->   F
F   -(Half)->   LH
F   -(Long)->   LL
F   -(Short)->  LS
LH  -(Full)->   F
LH  -(Half)->   RH
LH  -(Long)->   LL
LH  -(Short)->  RS
RH  -(Full)->   F
RH  -(Half)->   LH
RH  -(Long)->   LL
RH  -(Short)->  LS
LL  -(Full)->   F
LL  -(Half)->   LH
LL  -(Long)->   LL
LL  -(Short)->  RS
RL  -(Full)->   F
RL  -(Half)->   LH
RL  -(Long)->   LL
RL  -(Short)->  LS
LS  -(Full)->   F
LS  -(Half)->   RH
LS  -(Long)->   RL
LS  -(Short)->  MS
MS  -(Full)->   F
MS  -(Half)->   LH
MS  -(Long)->   LL
MS  -(Short)->  RS
RS  -(Full)->   F
RS  -(Half)->   LH
RS  -(Long)->   LL
RS  -(Short)->  LS
```

