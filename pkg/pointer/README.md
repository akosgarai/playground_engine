# Pointer package

It holds the mouse movement related informations.

- The current mouse position (cx,cy)
- The movement since the last Update. (dx,dy)

The `New` returns a Pointer instance. The `GetCurrent` returns the current position (cx, cy) of the pointer. The `GetDelta` returns the delta (dx, dy) of the pointer.
