package pointer

import "testing"

func TestNew(t *testing.T) {
	testData := []struct {
		cx float64
		cy float64
		dx float64
		dy float64
	}{
		{1, 1, 0.5, 0.5},
		{0.5, 0.5, 0.5, 0.5},
		{0.6, 0.6, 0.5, 0.5},
		{0.6, 0.6, 0.0, 0.0},
	}
	for _, tt := range testData {
		p := New(tt.cx, tt.cy, tt.dx, tt.dy)
		if p.cx != tt.cx {
			t.Errorf("Invalid cx. Instead of '%f', it is '%f'.", tt.cx, p.cx)
		}
		if p.cy != tt.cy {
			t.Errorf("Invalid cy. Instead of '%f', it is '%f'.", tt.cy, p.cy)
		}
		if p.dx != tt.dx {
			t.Errorf("Invalid dx. Instead of '%f', it is '%f'.", tt.dx, p.dx)
		}
		if p.dy != tt.dy {
			t.Errorf("Invalid dy. Instead of '%f', it is '%f'.", tt.dy, p.dy)
		}
	}
}
func TestGetCurrent(t *testing.T) {
	testData := []struct {
		cx float64
		cy float64
		dx float64
		dy float64
	}{
		{1, 1, 0.5, 0.5},
		{0.5, 0.5, 0.5, 0.5},
		{0.6, 0.6, 0.5, 0.5},
		{0.6, 0.6, 0.0, 0.0},
	}
	for _, tt := range testData {
		p := New(tt.cx, tt.cy, tt.dx, tt.dy)
		cx, cy := p.GetCurrent()
		if tt.cx != cx {
			t.Errorf("Invalid cx. Instead of '%f', it is '%f'.", tt.cx, cx)
		}
		if tt.cy != cy {
			t.Errorf("Invalid cy. Instead of '%f', it is '%f'.", tt.cy, cy)
		}
	}
}
func TestGetDelta(t *testing.T) {
	testData := []struct {
		cx float64
		cy float64
		dx float64
		dy float64
	}{
		{1, 1, 0.5, 0.5},
		{0.5, 0.5, 0.5, 0.5},
		{0.6, 0.6, 0.5, 0.5},
		{0.6, 0.6, 0.0, 0.0},
	}
	for _, tt := range testData {
		p := New(tt.cx, tt.cy, tt.dx, tt.dy)
		dx, dy := p.GetDelta()
		if tt.dx != dx {
			t.Errorf("Invalid dx. Instead of '%f', it is '%f'.", tt.dx, dx)
		}
		if tt.dy != dy {
			t.Errorf("Invalid dy. Instead of '%f', it is '%f'.", tt.dy, dy)
		}
	}
}
