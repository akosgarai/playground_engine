# User Interface Components

This package provides components that can be used for implementing a UI.

## Components

The following components are implemented.

### Button

```
-----------------
|               |
| ------------- |
| | labeltext | |
| ------------- |
|               |
-----------------
```

A button is a 2 state switch, as it could be activated (clicked), or not (untouched or clicked again). The state is toggled by the click event. The button itself is displayed as 2 meshes. The label setup is also hold by the button, but displaying the label text is not a button responsibility.

**Label**

The text, color, position and size are stored for displaying the label.

**Materials**

The button stores 3 kind of materials. One for the frame and the off state (without hover) (defaultMaterial), one for the hover action (hoverMaterial) and one for the on state (onStateMaterial).

### ButtonBuilder

This tool is provided for building buttons. For building a button, we have to setup the label params, the materials and the size parameters.
