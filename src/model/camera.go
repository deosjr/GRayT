package model

import "math"

type Camera struct {
	Ray
	FieldOfView float64
	Image       Image
}

// Aligned to POSITIVE z-axis for now
// change direction of Ray to (0,0,-1)
// and pixel (line 34) to (px, py, -1)
// for negative z-axis align
func NewCamera(w, h uint) Camera {
	img := newImage(w, h)
	r := NewRay(Vector{0, 0, 0}, Vector{0, 0, 1})
	return Camera{
		Ray:         r,
		FieldOfView: 0.5 * math.Pi,
		Image:       img,
	}
}

// 3d translation of 2d point on view
// assumes width >= height
// view size is 2x2 in world, from -1,1 to 1,-1
func (c Camera) PixelRay(x, y int) Ray {
	w, h := float64(c.Image.width), float64(c.Image.height)
	tan := math.Tan(c.FieldOfView / 2)
	px := (w / h) * (2*((float64(x)+0.5)/w) - 1) * tan
	py := (1 - 2*((float64(y)+0.5)/h)) * tan

	pixel := Vector{px, py, 1}
	direction := VectorFromTo(c.Origin, pixel)
	return NewRay(pixel, direction)
}
