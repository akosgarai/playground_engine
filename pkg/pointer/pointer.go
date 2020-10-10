package pointer

type Pointer struct {
	cx float64
	cy float64
	dx float64
	dy float64
}

// New returns a pointer.
func New(cx, cy, dx, dy float64) *Pointer {
	return &Pointer{
		cx: cx,
		cy: cy,
		dx: dx,
		dy: dy,
	}
}

// GetCurrent returns the current position (cx, cy) of the pointer.
func (p *Pointer) GetCurrent() (float64, float64) {
	return p.cx, p.cy
}

// GetDelta returns the delta (dx, dy) of the pointer.
func (p *Pointer) GetDelta() (float64, float64) {
	return p.dx, p.dy
}
