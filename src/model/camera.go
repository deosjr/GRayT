package model

import (
	"math"
)

type Camera struct {
	Ray
	FieldOfView   float64
	Width, Height float64
	aspectRatio   float64
	tan           float64
}

func NewCamera(w, h uint) Camera {
	r := NewRay(Vector{0, 0, 0}, Vector{0, 0, -1})
	wf := float64(w)
	hf := float64(h)
	fov := 0.5 * math.Pi
	return Camera{
		Ray: r,
		// precompute some constants used in pixelray() method
		FieldOfView: fov,
		Width:       wf,
		Height:      hf,
		aspectRatio: wf / hf,
		tan:         math.Tan(fov / 2),
	}
}

// 3d translation of 2d point on view
// assumes width >= height
// view size is 2x2 in world, from -1,1 to 1,-1
func (c Camera) PixelRay(x, y int) Ray {
	xNDC := (float64(x) + 0.5) / c.Width
	yNDC := (float64(y) + 0.5) / c.Height

	xScreen := 2*xNDC - 1
	yScreen := 1 - 2*yNDC

	px := xScreen * c.aspectRatio * c.tan
	py := yScreen * c.tan

	pixel := Vector{px, py, -1}
	direction := VectorFromTo(c.Origin, pixel)
	return NewRay(pixel, direction)
}
