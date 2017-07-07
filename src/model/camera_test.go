package model

import (
	"testing"
)

func TestCameraPixelRay(t *testing.T) {
	for i, tt := range []struct {
		c             Camera
		x, y          int
		wantOrigin    *Vector
		wantDirection *Vector
	}{
		{
			// For uneven resolution, we can point at the exact mid pixel
			// and expect (0,0,-1) (0,0,-1) as per camera specs
			c:             NewCamera(101, 101),
			x:             50,
			y:             50,
			wantOrigin:    &Vector{0, 0, -1},
			wantDirection: &Vector{0, 0, -1},
		},
		{
			// ULHC should be at (-1, 1, -1),
			// but moved 0.5 * pixelwidth/height to the center
			c:          NewCamera(100, 100),
			x:          0,
			y:          0,
			wantOrigin: &Vector{-0.99, 0.99, -1},
		},
	} {
		got := tt.c.PixelRay(tt.x, tt.y)
		if tt.wantOrigin != nil && got.Origin != *tt.wantOrigin {
			t.Errorf("%d) got origin %v want %v", i, got.Origin, *tt.wantOrigin)
		}
		if tt.wantDirection != nil && got.Direction != *tt.wantDirection {
			t.Errorf("%d) got direction %v want %v", i, got.Direction, *tt.wantDirection)
		}
	}
}
