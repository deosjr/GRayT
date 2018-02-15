package model

import (
	"math"
	"testing"
)

func TestCameraPixelRay(t *testing.T) {
	for i, tt := range []struct {
		c             Camera
		x, y          int
		from, to, up  Vector
		wantOrigin    Vector
		wantDirection Vector
	}{
		{
			// For uneven resolution, we can point at the exact mid pixel
			// and expect direction (to - from) normalized
			c:             NewPerspectiveCamera(101, 101, 0.5*math.Pi),
			from:          Vector{0, 0, 0},
			to:            Vector{0, 0, 1},
			up:            Vector{0, 1, 0},
			x:             50,
			y:             50,
			wantOrigin:    Vector{0, 0, 0},
			wantDirection: Vector{0, 0, 1},
		},
		{
			// For uneven resolution, we can point at the exact mid pixel
			// and expect direction (to - from) normalized
			c:             NewPerspectiveCamera(101, 101, 0.5*math.Pi),
			from:          Vector{2, 2, 2},
			to:            Vector{3, 3, 3},
			up:            Vector{0, 1, 0},
			x:             50,
			y:             50,
			wantOrigin:    Vector{2, 2, 2},
			wantDirection: Vector{3, 3, 3}.Sub(Vector{2, 2, 2}).Normalize(),
		},
	} {
		tt.c.LookAt(tt.from, tt.to, tt.up)
		got := tt.c.PixelRay(tt.x, tt.y)
		if got.Origin != tt.wantOrigin {
			t.Errorf("%d) got origin %v want %v", i, got.Origin, tt.wantOrigin)
		}
		if got.Direction != tt.wantDirection {
			t.Errorf("%d) got direction %v want %v", i, got.Direction, tt.wantDirection)
		}
	}
}
