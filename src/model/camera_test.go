package model

import (
	"testing"
)

func TestCameraPixelRay(t *testing.T) {
	for i, tt := range []struct {
		c    Camera
		x, y int
		want Ray
	}{
		{
			// For uneven resolution, we can point at the exact mid pixel
			// and expect (0,0,-1) (0,0,-1) as per camera specs
			c:    NewCamera(641, 481),
			x:    320,
			y:    240,
			want: NewRay(Vector{0, 0, -1}, Vector{0, 0, -1}),
		},
	} {
		got := tt.c.PixelRay(tt.x, tt.y)
		if got != tt.want {
			t.Errorf("%d) got %v want %v", i, got, tt.want)
		}
	}
}
